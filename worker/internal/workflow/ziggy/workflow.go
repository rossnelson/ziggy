package ziggy

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	z "ziggy/internal/ziggy"
)

const (
	SignalFeed = "feed"
	SignalPlay = "play"
	SignalPet  = "pet"
	SignalWake = "wake"

	QueryState = "state"

	SignalUpdateNeedMessage = "updateNeedMessage"
	SignalPoolResult        = "pool_result"

	PoolRegenerationInterval = 6 * time.Hour
	PoolRegenerationCooldown = 10 * time.Minute
)

type Input struct {
	Owner      string    `json:"owner"`
	Timezone   string    `json:"timezone"`
	Generation int       `json:"generation"`
	CreatedAt  time.Time `json:"createdAt,omitempty"`
}

type UpdateNeedMessageSignal struct {
	Message     string        `json:"message"`
	Personality z.Personality `json:"personality,omitempty"`
}

type PoolRegenerationOutput struct {
	Pool        *z.MessagePool `json:"pool"`
	GeneratedAt time.Time      `json:"generatedAt"`
}

type PoolRegenerateSignal struct {
	Personality z.Personality `json:"personality"`
	Stage       z.Stage       `json:"stage"`
	Bond        float64       `json:"bond"`
}

func Workflow(ctx workflow.Context, input Input) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Ziggy workflow started", "owner", input.Owner, "generation", input.Generation)

	timezone := input.Timezone
	if timezone == "" {
		timezone = "America/Los_Angeles"
	}

	state := z.NewState(timezone)
	state.Generation = input.Generation
	if state.Generation == 0 {
		state.Generation = 1
	}
	if !input.CreatedAt.IsZero() {
		state.CreatedAt = input.CreatedAt
	}

	err := workflow.SetQueryHandler(ctx, QueryState, func() (z.State, error) {
		return state, nil
	})
	if err != nil {
		return err
	}

	regeneratePool := func(reason string) {
		triggerPoolRegeneration(ctx, &state, logger, reason)
	}

	activityOpts := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    500 * time.Millisecond,
			BackoffCoefficient: 2.0,
			MaximumAttempts:    3,
		},
	}
	actCtx := workflow.WithActivityOptions(ctx, activityOpts)

	processAction := func(action z.Action) {
		input := ProcessActionInput{
			State:  state,
			Action: action,
			Now:    workflow.Now(ctx),
		}
		var output ProcessActionOutput
		err := workflow.ExecuteActivity(actCtx, "ProcessAction", input).Get(ctx, &output)
		if err != nil {
			logger.Info("ProcessAction failed", "action", action, "error", err.Error())
			return
		}
		state = output.State
	}

	regeneratePool("startup")

	feedCh := workflow.GetSignalChannel(ctx, SignalFeed)
	playCh := workflow.GetSignalChannel(ctx, SignalPlay)
	petCh := workflow.GetSignalChannel(ctx, SignalPet)
	wakeCh := workflow.GetSignalChannel(ctx, SignalWake)
	needMsgCh := workflow.GetSignalChannel(ctx, SignalUpdateNeedMessage)
	poolResultCh := workflow.GetSignalChannel(ctx, SignalPoolResult)

	for {
		selector := workflow.NewSelector(ctx)
		prevPersonality := state.Personality
		prevStage := z.GetStageForAge(workflow.Now(ctx).Sub(state.CreatedAt).Seconds())

		selector.AddReceive(feedCh, func(c workflow.ReceiveChannel, more bool) {
			var signal struct{}
			c.Receive(ctx, &signal)
			processAction(z.ActionFeed)
		})

		selector.AddReceive(playCh, func(c workflow.ReceiveChannel, more bool) {
			var signal struct{}
			c.Receive(ctx, &signal)
			processAction(z.ActionPlay)
		})

		selector.AddReceive(petCh, func(c workflow.ReceiveChannel, more bool) {
			var signal struct{}
			c.Receive(ctx, &signal)
			processAction(z.ActionPet)
		})

		selector.AddReceive(wakeCh, func(c workflow.ReceiveChannel, more bool) {
			var signal struct{}
			c.Receive(ctx, &signal)
			processAction(z.ActionWake)
		})

		selector.AddReceive(needMsgCh, func(c workflow.ReceiveChannel, more bool) {
			var signal UpdateNeedMessageSignal
			c.Receive(ctx, &signal)
			now := workflow.Now(ctx)

			if signal.Personality != "" {
				state.Personality = signal.Personality
				logger.Info("Updated personality from need updater", "personality", signal.Personality)
			}

			if signal.Message != "" {
				lastAction := state.GetMostRecentActionTime()
				if lastAction.IsZero() || now.Sub(lastAction) > NeedMessageDelay {
					state.Message = signal.Message
					logger.Info("Updated need message", "message", signal.Message)
				}
			}
		})

		selector.AddReceive(poolResultCh, func(c workflow.ReceiveChannel, more bool) {
			var result PoolRegenerationOutput
			c.Receive(ctx, &result)
			if result.Pool != nil {
				state.RuntimePool = result.Pool
				state.PoolGeneratedAt = result.GeneratedAt
				logger.Info("Pool updated from regenerator workflow")
			} else {
				logger.Info("Pool regenerator returned nil, using fallback")
			}
		})

		selector.Select(ctx)

		if state.Personality != prevPersonality {
			logger.Info("Personality changed", "from", prevPersonality, "to", state.Personality)
			regeneratePool("personality_change")
		}

		currentStage := z.GetStageForAge(workflow.Now(ctx).Sub(state.CreatedAt).Seconds())
		if currentStage != prevStage {
			logger.Info("Stage changed", "from", prevStage, "to", currentStage)
			state.Stage = currentStage
			regeneratePool("stage_change")
		}

		if workflow.GetInfo(ctx).GetCurrentHistoryLength() > 10000 {
			logger.Info("Continuing as new due to history length")
			return workflow.NewContinueAsNewError(ctx, Workflow, Input{
				Owner:      input.Owner,
				Timezone:   input.Timezone,
				Generation: state.Generation + 1,
				CreatedAt:  state.CreatedAt,
			})
		}
	}
}

const (
	NeedMessageDelay      = 30 * time.Second
	SignalPoolRegenerate  = "pool_regenerate"
)

func triggerPoolRegeneration(ctx workflow.Context, state *z.State, logger interface{ Info(string, ...interface{}) }, reason string) {
	now := workflow.Now(ctx)

	if !state.PoolGeneratedAt.IsZero() && now.Sub(state.PoolGeneratedAt) < PoolRegenerationCooldown {
		logger.Info("Skipping pool regeneration (cooldown)", "reason", reason, "lastGenerated", state.PoolGeneratedAt)
		return
	}

	logger.Info("Signaling pool regenerator", "reason", reason, "personality", state.Personality)

	state.PoolGeneratedAt = now

	age := now.Sub(state.CreatedAt).Seconds()
	workflowID := workflow.GetInfo(ctx).WorkflowExecution.ID
	poolWorkflowID := workflowID + "-pool-regenerator"

	signal := PoolRegenerateSignal{
		Personality: state.Personality,
		Stage:       z.GetStageForAge(age),
		Bond:        state.Bond,
	}

	workflow.SignalExternalWorkflow(ctx, poolWorkflowID, "", SignalPoolRegenerate, signal)
}
