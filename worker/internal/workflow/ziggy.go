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

	feedCh := workflow.GetSignalChannel(ctx, SignalFeed)
	playCh := workflow.GetSignalChannel(ctx, SignalPlay)
	petCh := workflow.GetSignalChannel(ctx, SignalPet)
	wakeCh := workflow.GetSignalChannel(ctx, SignalWake)

	for {
		selector := workflow.NewSelector(ctx)

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

		selector.Select(ctx)

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

func getPoolSelector(state *ZiggyState) *PoolSelector {
	fallback := GetFallbackPool(state.Personality)
	generic := GetFallbackPool(PersonalityStoic)
	return NewPoolSelector(state.RuntimePool, fallback, generic)
}

func handleFeed(state *ZiggyState, now time.Time, logger interface{ Info(string, ...interface{}) }) {
	pool := getPoolSelector(state)
	state.CareMetrics.RecordInteraction(state.Fullness, state.Bond, now)
	state.Personality = DerivePersonality(state.CareMetrics, state.Bond, now)

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
	state.CareMetrics.RecordInteraction(state.Fullness, state.Bond, now)
	state.Personality = DerivePersonality(state.CareMetrics, state.Bond, now)

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
	state.CareMetrics.RecordInteraction(state.Fullness, state.Bond, now)
	state.Personality = DerivePersonality(state.CareMetrics, state.Bond, now)

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
