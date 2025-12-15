package workflow

import (
	"math/rand"
	"time"

	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/workflow"
)

const (
	SignalFeed = "feed"
	SignalPlay = "play"
	SignalPet  = "pet"
	SignalWake = "wake"

	QueryState = "state"

	DecayInterval = 10 * time.Second
)

type ZiggyInput struct {
	Owner      string `json:"owner"`
	Generation int    `json:"generation"`
}

func ZiggyWorkflow(ctx workflow.Context, input ZiggyInput) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Ziggy workflow started", "owner", input.Owner, "generation", input.Generation)

	state := NewZiggyState()
	state.Generation = input.Generation
	if state.Generation == 0 {
		state.Generation = 1
	}

	// Set up query handler for state
	err := workflow.SetQueryHandler(ctx, QueryState, func() (ZiggyState, error) {
		return state, nil
	})
	if err != nil {
		return err
	}

	// Signal channels
	feedCh := workflow.GetSignalChannel(ctx, SignalFeed)
	playCh := workflow.GetSignalChannel(ctx, SignalPlay)
	petCh := workflow.GetSignalChannel(ctx, SignalPet)
	wakeCh := workflow.GetSignalChannel(ctx, SignalWake)

	// Decay timer
	decayTimer := workflow.NewTimer(ctx, DecayInterval)

	for {
		selector := workflow.NewSelector(ctx)

		// Handle feed signal
		selector.AddReceive(feedCh, func(c workflow.ReceiveChannel, more bool) {
			var signal struct{}
			c.Receive(ctx, &signal)
			handleFeed(&state, logger)
		})

		// Handle play signal
		selector.AddReceive(playCh, func(c workflow.ReceiveChannel, more bool) {
			var signal struct{}
			c.Receive(ctx, &signal)
			handlePlay(&state, logger)
		})

		// Handle pet signal
		selector.AddReceive(petCh, func(c workflow.ReceiveChannel, more bool) {
			var signal struct{}
			c.Receive(ctx, &signal)
			handlePet(&state, logger)
		})

		// Handle wake signal
		selector.AddReceive(wakeCh, func(c workflow.ReceiveChannel, more bool) {
			var signal struct{}
			c.Receive(ctx, &signal)
			handleWake(&state, logger)
		})

		// Handle decay timer
		selector.AddFuture(decayTimer, func(f workflow.Future) {
			handleDecay(&state, logger)
			decayTimer = workflow.NewTimer(ctx, DecayInterval)
		})

		selector.Select(ctx)

		// Check for continue-as-new (e.g., after evolution or history limit)
		if workflow.GetInfo(ctx).GetCurrentHistoryLength() > 10000 {
			logger.Info("Continuing as new due to history length")
			return workflow.NewContinueAsNewError(ctx, ZiggyWorkflow, ZiggyInput{
				Owner:      input.Owner,
				Generation: state.Generation,
			})
		}
	}
}

func handleFeed(state *ZiggyState, logger log.Logger) {
	if state.Sleeping {
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
		state.Message = pickRandom(messagesFeedFull)
	} else if state.Fullness < 30 {
		state.Fullness += 25
		state.Happiness += 5
		state.Message = pickRandom(messagesFeedHungry)
	} else {
		state.Fullness += 25
		state.Happiness += 5
		state.Message = pickRandom(messagesFeedSuccess)
	}

	state.LastAction = ActionFeed
	state.Clamp()
	logger.Info("Fed Ziggy", "fullness", state.Fullness, "happiness", state.Happiness)
}

func handlePlay(state *ZiggyState, logger log.Logger) {
	if state.Sleeping {
		logger.Info("Cannot play - Ziggy is sleeping")
		return
	}

	tooTired := state.Fullness < 20 || state.HP < 30
	if tooTired {
		state.Happiness += 5
		state.Fullness -= 5
		state.Message = pickRandom(messagesPlayTired)
	} else {
		state.Happiness += 20
		state.Fullness -= 10
		state.Bond += 5
		if state.GetMood() == MoodHappy {
			state.Message = pickRandom(messagesPlayHappy)
		} else {
			state.Message = pickRandom(messagesPlaySuccess)
		}
	}

	state.LastAction = ActionPlay
	state.Clamp()
	logger.Info("Played with Ziggy", "happiness", state.Happiness, "fullness", state.Fullness)
}

func handlePet(state *ZiggyState, logger log.Logger) {
	if state.Sleeping {
		state.Bond += 3
		state.Message = pickRandom(messagesPetSleeping)
	} else if state.Bond > 90 {
		state.Bond += 10
		state.Happiness += 5
		state.Message = pickRandom(messagesPetMaxBond)
	} else {
		state.Bond += 10
		state.Happiness += 5
		mood := state.GetMood()
		if mood == MoodSad || mood == MoodHungry {
			state.Message = pickRandom(messagesPetLowMood)
		} else {
			state.Message = pickRandom(messagesPetSuccess)
		}
	}

	state.LastAction = ActionPet
	state.Clamp()
	logger.Info("Petted Ziggy", "bond", state.Bond, "happiness", state.Happiness)
}

func handleWake(state *ZiggyState, logger log.Logger) {
	if !state.Sleeping {
		return
	}

	state.Sleeping = false
	state.Happiness -= 10
	state.Message = "*yawn*\nI was having\nsuch a nice dream..."
	state.Clamp()
	logger.Info("Woke Ziggy", "happiness", state.Happiness)
}

func handleDecay(state *ZiggyState, logger log.Logger) {
	if state.Sleeping {
		// During sleep: slow fullness decay, slow happiness recovery
		state.Fullness -= 1
		state.Happiness += 0.5

		targetHP := (state.Fullness + state.Happiness + state.Bond) / 3
		if state.HP < targetHP {
			state.HP += 1
		}
	} else {
		// Normal decay with bond protection
		bondProtection := 0.0
		if state.Bond > 50 {
			bondProtection = (state.Bond - 50) / 100
		}

		state.Fullness -= 2 * (1 - bondProtection)
		state.Happiness -= 1 * (1 - bondProtection)
		state.Bond -= 0.5

		// HP trends toward average of other stats
		targetHP := (state.Fullness + state.Happiness + state.Bond) / 3
		if state.HP > targetHP {
			state.HP -= 2
		} else if state.HP < targetHP {
			state.HP += 1
		}
	}

	state.Age += DecayInterval.Seconds()
	state.Clamp()

	// Update idle message periodically
	mood := state.GetMood()
	if msgs, ok := messagesIdle[mood]; ok {
		state.Message = pickRandom(msgs)
	}

	logger.Debug("Decay tick", "fullness", state.Fullness, "happiness", state.Happiness, "hp", state.HP)
}

func pickRandom(messages []string) string {
	if len(messages) == 0 {
		return ""
	}
	return messages[rand.Intn(len(messages))]
}
