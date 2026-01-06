package chat

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	z "ziggy/internal/ziggy"
)

const (
	SignalSendMessage  = "send_message"
	SignalStartMystery = "start_mystery"

	QueryChatHistory  = "chat_history"
	QueryMysteryStatus = "mystery_status"

	MaxMessages = 50
)

type Input struct {
	Owner   string `json:"owner"`
	ZiggyID string `json:"ziggyId"`
	Track   string `json:"track"`

	RecentMessages  []Message `json:"recentMessages,omitempty"`
	ActiveMystery   *Mystery  `json:"activeMystery,omitempty"`
	MysteryProgress int       `json:"mysteryProgress,omitempty"`
	HintsGiven      []string  `json:"hintsGiven,omitempty"`
	Solved          []string  `json:"solved,omitempty"`
}

type SendMessageSignal struct {
	Content string `json:"content"`
}

type StartMysterySignal struct {
	MysteryID string `json:"mysteryId"`
	Track     string `json:"track"`
}

func Workflow(ctx workflow.Context, input Input) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Chat workflow started", "owner", input.Owner, "track", input.Track)

	track := input.Track
	if track == "" {
		track = "fun"
	}

	state := NewState(input.Owner)

	if len(input.RecentMessages) > 0 {
		state.Messages = input.RecentMessages
	}
	if input.ActiveMystery != nil {
		state.ActiveMystery = input.ActiveMystery
		state.MysteryProgress = input.MysteryProgress
		state.HintsGiven = input.HintsGiven
	}
	if len(input.Solved) > 0 {
		state.Solved = input.Solved
	}

	err := workflow.SetQueryHandler(ctx, QueryChatHistory, func() (HistoryResponse, error) {
		mysteryStatus := state.GetMysteryStatus()
		return HistoryResponse{
			Messages:      state.Messages,
			MysteryStatus: &mysteryStatus,
			IsTyping:      state.IsTyping,
		}, nil
	})
	if err != nil {
		return err
	}

	err = workflow.SetQueryHandler(ctx, QueryMysteryStatus, func() (MysteryStatus, error) {
		return state.GetMysteryStatus(), nil
	})
	if err != nil {
		return err
	}

	messageCh := workflow.GetSignalChannel(ctx, SignalSendMessage)
	mysteryCh := workflow.GetSignalChannel(ctx, SignalStartMystery)

	activityOpts := workflow.ActivityOptions{
		StartToCloseTimeout: 2 * time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 3,
		},
	}
	actCtx := workflow.WithActivityOptions(ctx, activityOpts)

	queryOpts := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    500 * time.Millisecond,
			BackoffCoefficient: 2.0,
			MaximumInterval:    10 * time.Second,
			MaximumAttempts:    5,
		},
	}
	queryCtx := workflow.WithActivityOptions(ctx, queryOpts)

	for {
		selector := workflow.NewSelector(ctx)

		selector.AddReceive(messageCh, func(c workflow.ReceiveChannel, more bool) {
			var signal SendMessageSignal
			c.Receive(ctx, &signal)
			now := workflow.Now(ctx)

			responseTrack := track
			if state.ActiveMystery != nil && state.ActiveMystery.Track != "" {
				responseTrack = state.ActiveMystery.Track
			}

			if responseTrack == "educational" && state.ActiveMystery != nil {
				state.AddMessage("ziggy", "Searching the Temporal docs...", now)
				state.IsTyping = true
			}

			var ziggyState *z.State
			err := workflow.ExecuteActivity(queryCtx, "QueryZiggyState", input.ZiggyID).Get(ctx, &ziggyState)
			if err != nil {
				logger.Info("Failed to query Ziggy state", "error", err.Error())
			}

			processInput := ProcessMessageInput{
				State:      state,
				Content:    signal.Content,
				ZiggyState: ziggyState,
				Track:      track,
				Now:        now,
			}

			var output ProcessMessageOutput
			err = workflow.ExecuteActivity(actCtx, "ProcessChatMessage", processInput).Get(ctx, &output)
			if err != nil {
				logger.Info("ProcessChatMessage failed", "error", err.Error())
				return
			}
			state = output.State
		})

		selector.AddReceive(mysteryCh, func(c workflow.ReceiveChannel, more bool) {
			var signal StartMysterySignal
			c.Receive(ctx, &signal)
			logger.Info("Starting mystery", "mysteryID", signal.MysteryID, "track", signal.Track)

			mysteryTrack := signal.Track
			if mysteryTrack == "" {
				mysteryTrack = track
			}

			mystery := GetMystery(signal.MysteryID, mysteryTrack)
			if mystery != nil {
				state.ActiveMystery = mystery
				state.MysteryProgress = 0
				state.HintsGiven = []string{}
			}
		})

		selector.Select(ctx)

		if len(state.Messages) >= MaxMessages {
			logger.Info("Continuing as new due to message limit")

			recentMessages := state.Messages
			if len(recentMessages) > 20 {
				recentMessages = recentMessages[len(recentMessages)-20:]
			}

			return workflow.NewContinueAsNewError(ctx, Workflow, Input{
				Owner:           input.Owner,
				ZiggyID:         input.ZiggyID,
				Track:           input.Track,
				RecentMessages:  recentMessages,
				ActiveMystery:   state.ActiveMystery,
				MysteryProgress: state.MysteryProgress,
				HintsGiven:      state.HintsGiven,
				Solved:          state.Solved,
			})
		}
	}
}
