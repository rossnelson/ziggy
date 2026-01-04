package workflow

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

const SignalPoolResult = "pool_result"

type PoolRegeneratorInput struct {
	ZiggyWorkflowID string      `json:"ziggyWorkflowId"`
	Personality     Personality `json:"personality"`
	Stage           Stage       `json:"stage"`
	Bond            float64     `json:"bond"`
}

func PoolRegeneratorWorkflow(ctx workflow.Context, input PoolRegeneratorInput) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting pool regeneration", "personality", input.Personality, "stage", input.Stage)

	actCtx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 2 * time.Minute,
	})

	var output PoolRegenerationOutput
	err := workflow.ExecuteActivity(actCtx, "RegeneratePool", PoolRegenerationInput{
		Personality: input.Personality,
		Stage:       input.Stage,
		Bond:        input.Bond,
	}).Get(ctx, &output)

	if err != nil {
		logger.Info("Pool regeneration activity failed", "error", err.Error())
		return err
	}

	// Signal result back to main workflow
	err = workflow.SignalExternalWorkflow(ctx, input.ZiggyWorkflowID, "", SignalPoolResult, output).Get(ctx, nil)
	if err != nil {
		logger.Info("Failed to signal pool result", "error", err.Error())
		return err
	}

	logger.Info("Pool regeneration complete, signaled result")
	return nil
}
