package ziggy

import (
	"context"
	"log"
	"time"

	"ziggy/internal/ai"
	z "ziggy/internal/pet"
)

type Activities struct {
	aiClient *ai.Client
}

func NewActivities(aiClient *ai.Client) *Activities {
	return &Activities{aiClient: aiClient}
}

type ProcessActionInput struct {
	State  z.State  `json:"state"`
	Action z.Action `json:"action"`
	Now    time.Time `json:"now"`
}

type ProcessActionOutput struct {
	State z.State `json:"state"`
}

type PoolRegenerationInput struct {
	Personality z.Personality `json:"personality"`
	Stage       z.Stage       `json:"stage"`
	Bond        float64       `json:"bond"`
}

func (a *Activities) ProcessAction(ctx context.Context, input ProcessActionInput) (*ProcessActionOutput, error) {
	state := input.State
	now := input.Now

	state = state.CalculateCurrentState(now)

	switch input.Action {
	case z.ActionFeed:
		processActionFeed(&state, now)
	case z.ActionPlay:
		processActionPlay(&state, now)
	case z.ActionPet:
		processActionPet(&state, now)
	case z.ActionWake:
		processActionWake(&state, now)
	}

	state.LastUpdateTime = now
	return &ProcessActionOutput{State: state}, nil
}

func processActionFeed(state *z.State, now time.Time) {
	age := now.Sub(state.CreatedAt).Seconds()
	if z.GetStageForAge(age) == z.StageEgg {
		state.Message = "*wiggle*\n*wiggle*\nStill hatching..."
		return
	}

	pool := getPoolSelector(state)

	effectiveCooldown := state.GetEffectiveCooldown(z.ActionFeed)
	if !state.LastFeedTime.IsZero() && now.Sub(state.LastFeedTime) < effectiveCooldown {
		state.Message = pool.Pick("feedCooldown")
		return
	}

	state.CareMetrics.RecordInteraction(state.Fullness, state.Bond, now)
	state.Personality = z.DerivePersonality(state.CareMetrics, state.Bond, now)
	state.LastFeedTime = now

	if state.HP == 0 {
		state.Fullness += 15
		state.HP += 5
		state.Message = pool.Pick("feedTun")
		state.LastAction = z.ActionFeed
		state.Clamp()
		if state.HP >= 20 {
			state.Message = pool.Pick("reviving")
		}
		return
	}

	if state.Sleeping {
		state.Message = pool.Pick("feedSleeping")
		return
	}

	wasOverfed := state.Fullness > 90
	bondProtection := 0.0
	if state.Bond > 50 {
		bondProtection = (state.Bond - 50) / 20
	}

	if wasOverfed {
		state.Fullness += 5
		state.Happiness += -15 + bondProtection
		state.Message = pool.Pick("feedFull")
	} else if state.Fullness < 30 {
		state.Fullness += 30
		state.Happiness += 8
		state.Message = pool.Pick("feedHungry")
	} else {
		state.Fullness += 28
		state.Happiness += 5
		state.Message = pool.Pick("feedSuccess")
	}

	state.LastAction = z.ActionFeed
	state.Clamp()
}

func processActionPlay(state *z.State, now time.Time) {
	age := now.Sub(state.CreatedAt).Seconds()
	if z.GetStageForAge(age) == z.StageEgg {
		state.Message = "*wiggle*\n*wiggle*\nStill hatching..."
		return
	}

	pool := getPoolSelector(state)

	effectiveCooldown := state.GetEffectiveCooldown(z.ActionPlay)
	if !state.LastPlayTime.IsZero() && now.Sub(state.LastPlayTime) < effectiveCooldown {
		state.Message = pool.Pick("playCooldown")
		return
	}

	state.CareMetrics.RecordInteraction(state.Fullness, state.Bond, now)
	state.Personality = z.DerivePersonality(state.CareMetrics, state.Bond, now)
	state.LastPlayTime = now

	if state.HP == 0 {
		state.Message = pool.Pick("playTun")
		state.LastAction = z.ActionPlay
		return
	}

	if state.Sleeping {
		state.Message = pool.Pick("playSleeping")
		return
	}

	tooTired := state.Fullness < 20 || state.HP < 30
	if tooTired {
		state.Happiness += 8
		state.Fullness -= 3
		state.Message = pool.Pick("playTired")
	} else {
		state.Happiness += 25
		state.Fullness -= 8
		state.Bond += 8
		if state.GetMood() == z.MoodHappy {
			state.Message = pool.Pick("playHappy")
		} else {
			state.Message = pool.Pick("playSuccess")
		}
	}

	state.LastAction = z.ActionPlay
	state.Clamp()
}

func processActionPet(state *z.State, now time.Time) {
	pool := getPoolSelector(state)

	effectiveCooldown := state.GetEffectiveCooldown(z.ActionPet)
	if !state.LastPetTime.IsZero() && now.Sub(state.LastPetTime) < effectiveCooldown {
		state.Message = pool.Pick("petCooldown")
		return
	}

	state.CareMetrics.RecordInteraction(state.Fullness, state.Bond, now)
	state.Personality = z.DerivePersonality(state.CareMetrics, state.Bond, now)
	state.LastPetTime = now

	if state.HP == 0 {
		state.Bond += 5
		state.HP += 2
		state.Message = pool.Pick("petTun")
		state.LastAction = z.ActionPet
		state.Clamp()
		if state.HP >= 20 {
			state.Message = pool.Pick("reviving")
		}
		return
	}

	if state.Sleeping {
		state.Bond += 5
		state.Message = pool.Pick("petSleeping")
	} else if state.Bond > 90 {
		state.Bond += 10
		state.Happiness += 8
		state.Message = pool.Pick("petMaxBond")
	} else {
		state.Bond += 15
		state.Happiness += 8
		mood := state.GetMood()
		if mood == z.MoodSad || mood == z.MoodHungry {
			state.Message = pool.Pick("petLowMood")
		} else {
			state.Message = pool.Pick("petSuccess")
		}
	}

	state.LastAction = z.ActionPet
	state.Clamp()
}

func processActionWake(state *z.State, now time.Time) {
	if !state.Sleeping {
		return
	}

	state.CareMetrics.RecordInteraction(state.Fullness, state.Bond, now)
	state.Personality = z.DerivePersonality(state.CareMetrics, state.Bond, now)

	state.Sleeping = false
	state.Happiness -= 10
	state.Message = "*yawn*\nI was having\nsuch a nice dream..."
	state.LastAction = z.ActionWake
	state.Clamp()
}

func (a *Activities) RegeneratePool(ctx context.Context, input PoolRegenerationInput) (*PoolRegenerationOutput, error) {
	log.Printf("[RegeneratePool] Starting pool regeneration: personality=%s stage=%s bond=%.1f",
		input.Personality, input.Stage, input.Bond)

	if a.aiClient == nil {
		log.Printf("[RegeneratePool] No AI client configured, using fallback pools")
		return &PoolRegenerationOutput{
			Pool:        nil,
			GeneratedAt: time.Now(),
		}, nil
	}

	bondDescription := getBondDescription(input.Bond)
	log.Printf("[RegeneratePool] Calling Claude API with bond description: %s", bondDescription)

	aiInput := ai.PoolGenerationInput{
		Personality:     string(input.Personality),
		Stage:           string(input.Stage),
		BondDescription: bondDescription,
	}

	aiPool, err := a.aiClient.GeneratePool(ctx, aiInput)
	if err != nil {
		log.Printf("[RegeneratePool] Claude API error: %v", err)
		return &PoolRegenerationOutput{
			Pool:        nil,
			GeneratedAt: time.Now(),
		}, nil
	}

	pool := convertAIPool(aiPool)
	log.Printf("[RegeneratePool] Successfully generated pool with %d feedSuccess messages",
		len(pool.FeedSuccess))

	return &PoolRegenerationOutput{
		Pool:        pool,
		GeneratedAt: time.Now(),
	}, nil
}

func getBondDescription(bond float64) string {
	if bond >= 80 {
		return "deeply bonded (best friends)"
	}
	if bond >= 60 {
		return "close bond (good friends)"
	}
	if bond >= 40 {
		return "developing bond (getting to know each other)"
	}
	if bond >= 20 {
		return "new acquaintance (still shy)"
	}
	return "barely met (very timid)"
}

func convertAIPool(aiPool *ai.MessagePool) *z.MessagePool {
	return &z.MessagePool{
		FeedSuccess:  aiPool.FeedSuccess,
		FeedFull:     aiPool.FeedFull,
		FeedHungry:   aiPool.FeedHungry,
		FeedSleeping: aiPool.FeedSleeping,
		FeedTun:      aiPool.FeedTun,
		FeedCooldown: aiPool.FeedCooldown,

		PlaySuccess:  aiPool.PlaySuccess,
		PlayTired:    aiPool.PlayTired,
		PlayHappy:    aiPool.PlayHappy,
		PlaySleeping: aiPool.PlaySleeping,
		PlayTun:      aiPool.PlayTun,
		PlayCooldown: aiPool.PlayCooldown,

		PetSuccess:  aiPool.PetSuccess,
		PetMaxBond:  aiPool.PetMaxBond,
		PetLowMood:  aiPool.PetLowMood,
		PetSleeping: aiPool.PetSleeping,
		PetTun:      aiPool.PetTun,
		PetCooldown: aiPool.PetCooldown,

		Reviving: aiPool.Reviving,

		IdleHappy:    aiPool.IdleHappy,
		IdleNeutral:  aiPool.IdleNeutral,
		IdleHungry:   aiPool.IdleHungry,
		IdleSad:      aiPool.IdleSad,
		IdleLonely:   aiPool.IdleLonely,
		IdleCritical: aiPool.IdleCritical,
		IdleTun:      aiPool.IdleTun,
		IdleSleeping: aiPool.IdleSleeping,

		NeedsFood:      aiPool.NeedsFood,
		NeedsPlay:      aiPool.NeedsPlay,
		NeedsAffection: aiPool.NeedsAffection,
		NeedsCritical:  aiPool.NeedsCritical,
	}
}

func getPoolSelector(state *z.State) *z.PoolSelector {
	fallback := z.GetFallbackPool(state.Personality)
	generic := z.GetFallbackPool(z.PersonalityStoic)
	return z.NewPoolSelector(state.RuntimePool, fallback, generic)
}
