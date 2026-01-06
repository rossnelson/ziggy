package pool_regenerator

import (
	"time"

	"go.temporal.io/sdk/workflow"

	z "ziggy/internal/ziggy"
)

const (
	SignalPoolResult     = "pool_result"
	SignalPoolRegenerate = "pool_regenerate"

	RegenerationInterval = 6 * time.Hour
)

type Input struct {
	ZiggyWorkflowID string `json:"ziggyWorkflowId"`
}

type RegenerateSignal struct {
	Personality z.Personality `json:"personality"`
	Stage       z.Stage       `json:"stage"`
	Bond        float64       `json:"bond"`
}

type RegenerationInput struct {
	Personality z.Personality `json:"personality"`
	Stage       z.Stage       `json:"stage"`
	Bond        float64       `json:"bond"`
}

type RegenerationOutput struct {
	Pool        *z.MessagePool `json:"pool"`
	GeneratedAt time.Time      `json:"generatedAt"`
}

func Workflow(ctx workflow.Context, input Input) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Pool regenerator started", "ziggyWorkflowId", input.ZiggyWorkflowID)

	actCtx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 2 * time.Minute,
	})

	regenerateCh := workflow.GetSignalChannel(ctx, SignalPoolRegenerate)
	timer := workflow.NewTimer(ctx, RegenerationInterval)

	for {
		selector := workflow.NewSelector(ctx)

		selector.AddReceive(regenerateCh, func(c workflow.ReceiveChannel, more bool) {
			var signal RegenerateSignal
			c.Receive(ctx, &signal)
			logger.Info("Received regenerate signal", "personality", signal.Personality, "stage", signal.Stage)
			regenerateAndSignal(ctx, actCtx, input.ZiggyWorkflowID, signal, logger)
		})

		selector.AddFuture(timer, func(f workflow.Future) {
			logger.Info("Scheduled regeneration timer fired")
			var state z.State
			err := workflow.ExecuteActivity(actCtx, "QueryZiggyState", input.ZiggyWorkflowID).Get(ctx, &state)
			if err != nil {
				logger.Info("Failed to query Ziggy state for scheduled regen", "error", err.Error())
			} else {
				signal := RegenerateSignal{
					Personality: state.Personality,
					Stage:       state.Stage,
					Bond:        state.Bond,
				}
				regenerateAndSignal(ctx, actCtx, input.ZiggyWorkflowID, signal, logger)
			}
			timer = workflow.NewTimer(ctx, RegenerationInterval)
		})

		selector.Select(ctx)

		if workflow.GetInfo(ctx).GetCurrentHistoryLength() > 5000 {
			logger.Info("Pool regenerator continuing as new")
			return workflow.NewContinueAsNewError(ctx, Workflow, input)
		}
	}
}

func regenerateAndSignal(ctx, actCtx workflow.Context, ziggyWorkflowID string, signal RegenerateSignal, logger interface{ Info(string, ...interface{}) }) {
	var output RegenerationOutput
	err := workflow.ExecuteActivity(actCtx, "RegeneratePool", RegenerationInput{
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
