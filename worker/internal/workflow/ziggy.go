package workflow

import (
	"math/rand"
	"time"

	"go.temporal.io/sdk/workflow"
)

const (
	SignalFeed = "feed"
	SignalPlay = "play"
	SignalPet  = "pet"
	SignalWake = "wake"

	QueryState = "state"

	PoolRegenerationInterval = 6 * time.Hour
)

type ZiggyInput struct {
	Owner      string `json:"owner"`
	Timezone   string `json:"timezone"`
	Generation int    `json:"generation"`
}

func ZiggyWorkflow(ctx workflow.Context, input ZiggyInput) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Ziggy workflow started", "owner", input.Owner, "generation", input.Generation)

	timezone := input.Timezone
	if timezone == "" {
		timezone = "America/Los_Angeles"
	}

	state := NewZiggyState(timezone)
	state.Generation = input.Generation
	if state.Generation == 0 {
		state.Generation = 1
	}

	err := workflow.SetQueryHandler(ctx, QueryState, func() (ZiggyState, error) {
		return state, nil
	})
	if err != nil {
		return err
	}

	regeneratePool := func(reason string) {
		triggerPoolRegeneration(ctx, &state, logger, reason)
	}

	// Generate initial pool on startup
	regeneratePool("startup")

	feedCh := workflow.GetSignalChannel(ctx, SignalFeed)
	playCh := workflow.GetSignalChannel(ctx, SignalPlay)
	petCh := workflow.GetSignalChannel(ctx, SignalPet)
	wakeCh := workflow.GetSignalChannel(ctx, SignalWake)

	// Timer for periodic pool regeneration
	timerFuture := workflow.NewTimer(ctx, PoolRegenerationInterval)

	for {
		selector := workflow.NewSelector(ctx)
		prevPersonality := state.Personality

		selector.AddReceive(feedCh, func(c workflow.ReceiveChannel, more bool) {
			var signal struct{}
			c.Receive(ctx, &signal)
			now := workflow.Now(ctx)
			state = state.CalculateCurrentState(now)
			handleFeed(&state, now, logger)
			state.LastUpdateTime = now
		})

		selector.AddReceive(playCh, func(c workflow.ReceiveChannel, more bool) {
			var signal struct{}
			c.Receive(ctx, &signal)
			now := workflow.Now(ctx)
			state = state.CalculateCurrentState(now)
			handlePlay(&state, now, logger)
			state.LastUpdateTime = now
		})

		selector.AddReceive(petCh, func(c workflow.ReceiveChannel, more bool) {
			var signal struct{}
			c.Receive(ctx, &signal)
			now := workflow.Now(ctx)
			state = state.CalculateCurrentState(now)
			handlePet(&state, now, logger)
			state.LastUpdateTime = now
		})

		selector.AddReceive(wakeCh, func(c workflow.ReceiveChannel, more bool) {
			var signal struct{}
			c.Receive(ctx, &signal)
			now := workflow.Now(ctx)
			state = state.CalculateCurrentState(now)
			handleWake(&state, now, logger)
			state.LastUpdateTime = now
		})

		selector.AddFuture(timerFuture, func(f workflow.Future) {
			logger.Info("Pool regeneration timer fired")
			regeneratePool("scheduled")
			timerFuture = workflow.NewTimer(ctx, PoolRegenerationInterval)
		})

		selector.Select(ctx)

		// Check for personality change after interaction
		if state.Personality != prevPersonality {
			logger.Info("Personality changed", "from", prevPersonality, "to", state.Personality)
			regeneratePool("personality_change")
		}

		if workflow.GetInfo(ctx).GetCurrentHistoryLength() > 10000 {
			logger.Info("Continuing as new due to history length")
			return workflow.NewContinueAsNewError(ctx, ZiggyWorkflow, ZiggyInput{
				Owner:      input.Owner,
				Timezone:   input.Timezone,
				Generation: state.Generation + 1,
			})
		}
	}
}

func triggerPoolRegeneration(ctx workflow.Context, state *ZiggyState, logger interface{ Info(string, ...interface{}) }, reason string) {
	logger.Info("Triggering pool regeneration", "reason", reason, "personality", state.Personality)

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
	}
	actCtx := workflow.WithActivityOptions(ctx, ao)

	age := workflow.Now(ctx).Sub(state.CreatedAt).Seconds()
	input := PoolRegenerationInput{
		Personality: state.Personality,
		Stage:       GetStageForAge(age),
		Bond:        state.Bond,
	}

	var output PoolRegenerationOutput
	err := workflow.ExecuteActivity(actCtx, "RegeneratePool", input).Get(ctx, &output)
	if err != nil {
		logger.Info("Pool regeneration failed, using fallback", "error", err.Error())
		return
	}

	state.RuntimePool = output.Pool
	state.PoolGeneratedAt = output.GeneratedAt
	logger.Info("Pool regenerated successfully")
}

func getPoolSelector(state *ZiggyState) *PoolSelector {
	fallback := GetFallbackPool(state.Personality)
	generic := GetFallbackPool(PersonalityStoic)
	return NewPoolSelector(state.RuntimePool, fallback, generic)
}

func handleFeed(state *ZiggyState, now time.Time, logger interface{ Info(string, ...interface{}) }) {
	pool := getPoolSelector(state)

	// Check cooldown (shorter when hungry)
	effectiveCooldown := state.GetEffectiveCooldown(ActionFeed)
	if !state.LastFeedTime.IsZero() && now.Sub(state.LastFeedTime) < effectiveCooldown {
		state.Message = pool.Pick("feedCooldown")
		logger.Info("Feed on cooldown", "remaining", effectiveCooldown-now.Sub(state.LastFeedTime))
		return
	}

	state.CareMetrics.RecordInteraction(state.Fullness, state.Bond, now)
	state.Personality = DerivePersonality(state.CareMetrics, state.Bond, now)
	state.LastFeedTime = now

	// Tun state: feeding helps revival
	if state.HP == 0 {
		state.Fullness += 15
		state.HP += 5 // Start revival
		state.Message = pool.Pick("feedTun")
		state.LastAction = ActionFeed
		state.Clamp()

		// Check if revived
		if state.HP >= 20 {
			state.Message = pool.Pick("reviving")
		}
		logger.Info("Fed Ziggy in tun state - reviving", "hp", state.HP)
		return
	}

	if state.Sleeping {
		state.Message = pool.Pick("feedSleeping")
		logger.Info("Cannot feed - Ziggy is sleeping")
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
		state.Fullness += 25
		state.Happiness += 5
		state.Message = pool.Pick("feedHungry")
	} else {
		state.Fullness += 25
		state.Happiness += 5
		state.Message = pool.Pick("feedSuccess")
	}

	state.LastAction = ActionFeed
	state.Clamp()
	logger.Info("Fed Ziggy", "fullness", state.Fullness, "happiness", state.Happiness)
}

func handlePlay(state *ZiggyState, now time.Time, logger interface{ Info(string, ...interface{}) }) {
	pool := getPoolSelector(state)

	// Check cooldown (shorter when unhappy)
	effectiveCooldown := state.GetEffectiveCooldown(ActionPlay)
	if !state.LastPlayTime.IsZero() && now.Sub(state.LastPlayTime) < effectiveCooldown {
		state.Message = pool.Pick("playCooldown")
		logger.Info("Play on cooldown", "remaining", effectiveCooldown-now.Sub(state.LastPlayTime))
		return
	}

	state.CareMetrics.RecordInteraction(state.Fullness, state.Bond, now)
	state.Personality = DerivePersonality(state.CareMetrics, state.Bond, now)
	state.LastPlayTime = now

	// Tun state: can't play
	if state.HP == 0 {
		state.Message = pool.Pick("playTun")
		state.LastAction = ActionPlay
		logger.Info("Cannot play - Ziggy is in tun state")
		return
	}

	if state.Sleeping {
		state.Message = pool.Pick("playSleeping")
		logger.Info("Cannot play - Ziggy is sleeping")
		return
	}

	tooTired := state.Fullness < 20 || state.HP < 30
	if tooTired {
		state.Happiness += 5
		state.Fullness -= 5
		state.Message = pool.Pick("playTired")
	} else {
		state.Happiness += 20
		state.Fullness -= 10
		state.Bond += 5
		if state.GetMood() == MoodHappy {
			state.Message = pool.Pick("playHappy")
		} else {
			state.Message = pool.Pick("playSuccess")
		}
	}

	state.LastAction = ActionPlay
	state.Clamp()
	logger.Info("Played with Ziggy", "happiness", state.Happiness, "fullness", state.Fullness)
}

func handlePet(state *ZiggyState, now time.Time, logger interface{ Info(string, ...interface{}) }) {
	pool := getPoolSelector(state)

	// Check cooldown (shorter when bond is low)
	effectiveCooldown := state.GetEffectiveCooldown(ActionPet)
	if !state.LastPetTime.IsZero() && now.Sub(state.LastPetTime) < effectiveCooldown {
		state.Message = pool.Pick("petCooldown")
		logger.Info("Pet on cooldown", "remaining", effectiveCooldown-now.Sub(state.LastPetTime))
		return
	}

	state.CareMetrics.RecordInteraction(state.Fullness, state.Bond, now)
	state.Personality = DerivePersonality(state.CareMetrics, state.Bond, now)
	state.LastPetTime = now

	// Tun state: petting helps revival through warmth/bond
	if state.HP == 0 {
		state.Bond += 5
		state.HP += 2 // Slow revival through comfort
		state.Message = pool.Pick("petTun")
		state.LastAction = ActionPet
		state.Clamp()

		// Check if revived
		if state.HP >= 20 {
			state.Message = pool.Pick("reviving")
		}
		logger.Info("Petted Ziggy in tun state - warming up", "hp", state.HP, "bond", state.Bond)
		return
	}

	if state.Sleeping {
		state.Bond += 3
		state.Message = pool.Pick("petSleeping")
	} else if state.Bond > 90 {
		state.Bond += 10
		state.Happiness += 5
		state.Message = pool.Pick("petMaxBond")
	} else {
		state.Bond += 10
		state.Happiness += 5
		mood := state.GetMood()
		if mood == MoodSad || mood == MoodHungry {
			state.Message = pool.Pick("petLowMood")
		} else {
			state.Message = pool.Pick("petSuccess")
		}
	}

	state.LastAction = ActionPet
	state.Clamp()
	logger.Info("Petted Ziggy", "bond", state.Bond, "happiness", state.Happiness)
}

func handleWake(state *ZiggyState, now time.Time, logger interface{ Info(string, ...interface{}) }) {
	if !state.Sleeping {
		return
	}

	state.CareMetrics.RecordInteraction(state.Fullness, state.Bond, now)
	state.Personality = DerivePersonality(state.CareMetrics, state.Bond, now)

	state.Sleeping = false
	state.Happiness -= 10
	state.Message = "*yawn*\nI was having\nsuch a nice dream..."
	state.LastAction = ActionWake
	state.Clamp()
	logger.Info("Woke Ziggy", "happiness", state.Happiness)
}

func GetIdleMessage(mood Mood) string {
	msgs, ok := messagesIdle[mood]
	if !ok {
		msgs = messagesIdle[MoodNeutral]
	}
	return pickRandom(msgs)
}

func pickRandom(messages []string) string {
	if len(messages) == 0 {
		return ""
	}
	return messages[rand.Intn(len(messages))]
}
