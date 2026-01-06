package need_updater

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	z "ziggy/internal/ziggy"
)

const (
	SignalUpdateNeedMessage = "updateNeedMessage"

	NeedUpdateInterval = 30 * time.Second
	NeedMessageDelay   = 30 * time.Second
	MaxIterations      = 100
)

type Input struct {
	ZiggyWorkflowID string `json:"ziggyWorkflowId"`
	Iteration       int    `json:"iteration"`
}

type UpdateNeedMessageSignal struct {
	Message     string        `json:"message"`
	Personality z.Personality `json:"personality,omitempty"`
}

func Workflow(ctx workflow.Context, input Input) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("NeedUpdater started", "ziggyWorkflowId", input.ZiggyWorkflowID, "iteration", input.Iteration)

	iteration := input.Iteration

	for {
		if err := workflow.Sleep(ctx, NeedUpdateInterval); err != nil {
			return err
		}

		iteration++

		state := queryZiggyState(ctx, input.ZiggyWorkflowID, logger)
		if state == nil {
			logger.Info("Could not query Ziggy state, will retry")
			continue
		}

		now := workflow.Now(ctx)
		lastAction := state.GetMostRecentActionTime()

		if !lastAction.IsZero() && now.Sub(lastAction) < NeedMessageDelay {
			continue
		}

		current := state.CalculateCurrentState(now)

		personality := z.DerivePersonality(state.CareMetrics, current.Bond, now)

		need := current.GetMostUrgentNeed()

		if personality != state.Personality {
			signalZiggyUpdate(ctx, input.ZiggyWorkflowID, "", personality, logger)
		}

		if need == z.NeedNone {
			continue
		}

		current.Personality = personality
		message := pickNeedMessage(&current, need)
		if message == "" {
			continue
		}

		signalZiggyUpdate(ctx, input.ZiggyWorkflowID, message, personality, logger)

		if iteration >= MaxIterations {
			logger.Info("Continuing as new", "iterations", iteration)
			return workflow.NewContinueAsNewError(ctx, Workflow, Input{
				ZiggyWorkflowID: input.ZiggyWorkflowID,
				Iteration:       0,
			})
		}
	}
}

func queryZiggyState(ctx workflow.Context, ziggyID string, logger interface{ Info(string, ...interface{}) }) *z.State {
	if ziggyID == "" {
		return nil
	}

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 2,
		},
	}
	actCtx := workflow.WithActivityOptions(ctx, ao)

	var state z.State
	err := workflow.ExecuteActivity(actCtx, "QueryZiggyState", ziggyID).Get(ctx, &state)
	if err != nil {
		logger.Info("Failed to query Ziggy state", "error", err.Error())
		return nil
	}

	return &state
}

func pickNeedMessage(state *z.State, need z.NeedType) string {
	fallback := z.GetFallbackPool(state.Personality)
	generic := z.GetFallbackPool(z.PersonalityStoic)
	selector := z.NewPoolSelector(state.RuntimePool, fallback, generic)
	return selector.Pick(string(need))
}

func signalZiggyUpdate(ctx workflow.Context, workflowID string, message string, personality z.Personality, logger interface{ Info(string, ...interface{}) }) {
	signal := UpdateNeedMessageSignal{
		Message:     message,
		Personality: personality,
	}

	err := workflow.SignalExternalWorkflow(ctx, workflowID, "", SignalUpdateNeedMessage, signal).Get(ctx, nil)
	if err != nil {
		logger.Info("Failed to signal Ziggy", "error", err.Error())
	}
}
