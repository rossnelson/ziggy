package workflow

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

const (
	SignalUpdateNeedMessage = "updateNeedMessage"

	NeedUpdateInterval       = 30 * time.Second
	NeedMessageDelay         = 30 * time.Second // Don't show need messages for 30s after action
	NeedUpdaterMaxIterations = 100              // Continue-as-new after this many iterations
)

type NeedUpdaterInput struct {
	ZiggyWorkflowID string `json:"ziggyWorkflowId"`
	Iteration       int    `json:"iteration"`
}

type UpdateNeedMessageSignal struct {
	Message     string      `json:"message"`
	Personality Personality `json:"personality,omitempty"`
}

// NeedUpdaterWorkflow periodically checks Ziggy's state and signals need messages
func NeedUpdaterWorkflow(ctx workflow.Context, input NeedUpdaterInput) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("NeedUpdater started", "ziggyWorkflowId", input.ZiggyWorkflowID, "iteration", input.Iteration)

	iteration := input.Iteration

	for {
		// Sleep before checking
		if err := workflow.Sleep(ctx, NeedUpdateInterval); err != nil {
			return err
		}

		iteration++

		// Query Ziggy's state via activity
		state := queryZiggyStateForNeed(ctx, input.ZiggyWorkflowID, logger)
		if state == nil {
			// Ziggy workflow might have ended, exit gracefully
			logger.Info("Could not query Ziggy state, exiting")
			return nil
		}

		// Check if we should show a need message
		now := workflow.Now(ctx)
		lastAction := state.GetMostRecentActionTime()

		// Don't show need messages if action was recent
		if !lastAction.IsZero() && now.Sub(lastAction) < NeedMessageDelay {
			continue
		}

		// Calculate current state with decay
		current := state.CalculateCurrentState(now)

		// Derive personality based on current care metrics
		personality := DerivePersonality(state.CareMetrics, current.Bond, now)

		need := current.GetMostUrgentNeed()

		// Always signal if personality changed, even without a need message
		if personality != state.Personality {
			signalZiggyUpdate(ctx, input.ZiggyWorkflowID, "", personality, logger)
		}

		if need == NeedNone {
			continue
		}

		// Pick a message from the pool (use new personality for correct voice)
		current.Personality = personality
		message := pickNeedMessage(&current, need)
		if message == "" {
			continue
		}

		// Signal Ziggy to update the message and personality
		signalZiggyUpdate(ctx, input.ZiggyWorkflowID, message, personality, logger)

		// Continue-as-new to bound history
		if iteration >= NeedUpdaterMaxIterations {
			logger.Info("Continuing as new", "iterations", iteration)
			return workflow.NewContinueAsNewError(ctx, NeedUpdaterWorkflow, NeedUpdaterInput{
				ZiggyWorkflowID: input.ZiggyWorkflowID,
				Iteration:       0,
			})
		}
	}
}

func queryZiggyStateForNeed(ctx workflow.Context, ziggyID string, logger interface{ Info(string, ...interface{}) }) *ZiggyState {
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

	var state ZiggyState
	err := workflow.ExecuteActivity(actCtx, "QueryZiggyState", ziggyID).Get(ctx, &state)
	if err != nil {
		logger.Info("Failed to query Ziggy state", "error", err.Error())
		return nil
	}

	return &state
}

func pickNeedMessage(state *ZiggyState, need NeedType) string {
	fallback := GetFallbackPool(state.Personality)
	generic := GetFallbackPool(PersonalityStoic)
	selector := NewPoolSelector(state.RuntimePool, fallback, generic)
	return selector.Pick(string(need))
}

func signalZiggyUpdate(ctx workflow.Context, workflowID string, message string, personality Personality, logger interface{ Info(string, ...interface{}) }) {
	signal := UpdateNeedMessageSignal{
		Message:     message,
		Personality: personality,
	}

	err := workflow.SignalExternalWorkflow(ctx, workflowID, "", SignalUpdateNeedMessage, signal).Get(ctx, nil)
	if err != nil {
		logger.Info("Failed to signal Ziggy", "error", err.Error())
	}
}
