package workflow

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

const (
	SignalPoolResult     = "pool_result"
	SignalPoolRegenerate = "pool_regenerate"
)

type PoolRegeneratorInput struct {
	ZiggyWorkflowID string `json:"ziggyWorkflowId"`
}

type PoolRegenerateSignal struct {
	Personality Personality `json:"personality"`
	Stage       Stage       `json:"stage"`
	Bond        float64     `json:"bond"`
}

func PoolRegeneratorWorkflow(ctx workflow.Context, input PoolRegeneratorInput) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Pool regenerator started", "ziggyWorkflowId", input.ZiggyWorkflowID)

	actCtx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 2 * time.Minute,
	})

	regenerateCh := workflow.GetSignalChannel(ctx, SignalPoolRegenerate)
	timer := workflow.NewTimer(ctx, PoolRegenerationInterval)

	for {
		selector := workflow.NewSelector(ctx)

		// Handle regenerate signal from Ziggy
		selector.AddReceive(regenerateCh, func(c workflow.ReceiveChannel, more bool) {
			var signal PoolRegenerateSignal
			c.Receive(ctx, &signal)
			logger.Info("Received regenerate signal", "personality", signal.Personality, "stage", signal.Stage)
			regenerateAndSignal(ctx, actCtx, input.ZiggyWorkflowID, signal, logger)
		})

		// Handle scheduled regeneration timer
		selector.AddFuture(timer, func(f workflow.Future) {
			logger.Info("Scheduled regeneration timer fired")
			// Query Ziggy for current state
			var state ZiggyState
			err := workflow.ExecuteActivity(actCtx, "QueryZiggyState", input.ZiggyWorkflowID).Get(ctx, &state)
			if err != nil {
				logger.Info("Failed to query Ziggy state for scheduled regen", "error", err.Error())
			} else {
				signal := PoolRegenerateSignal{
					Personality: state.Personality,
					Stage:       state.Stage,
					Bond:        state.Bond,
				}
				regenerateAndSignal(ctx, actCtx, input.ZiggyWorkflowID, signal, logger)
			}
			timer = workflow.NewTimer(ctx, PoolRegenerationInterval)
		})

		selector.Select(ctx)

		// Continue-as-new to prevent history growth
		if workflow.GetInfo(ctx).GetCurrentHistoryLength() > 5000 {
			logger.Info("Pool regenerator continuing as new")
			return workflow.NewContinueAsNewError(ctx, PoolRegeneratorWorkflow, input)
		}
	}
}

func regenerateAndSignal(ctx, actCtx workflow.Context, ziggyWorkflowID string, signal PoolRegenerateSignal, logger interface{ Info(string, ...interface{}) }) {
	var output PoolRegenerationOutput
	err := workflow.ExecuteActivity(actCtx, "RegeneratePool", PoolRegenerationInput{
		Personality: signal.Personality,
		Stage:       signal.Stage,
		Bond:        signal.Bond,
	}).Get(ctx, &output)

	if err != nil {
		logger.Info("Pool regeneration activity failed", "error", err.Error())
		return
	}

	err = workflow.SignalExternalWorkflow(ctx, ziggyWorkflowID, "", SignalPoolResult, output).Get(ctx, nil)
	if err != nil {
		logger.Info("Failed to signal pool result to Ziggy", "error", err.Error())
		return
	}

	logger.Info("Pool regenerated and signaled to Ziggy")
}
