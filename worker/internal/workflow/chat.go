package workflow

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

const (
	SignalSendMessage = "send_message"
	SignalStartMystery = "start_mystery"

	QueryChatHistory    = "chat_history"
	QueryMysteryStatus  = "mystery_status"

	MaxChatMessages = 50
)

type ChatInput struct {
	Owner    string `json:"owner"`
	ZiggyID  string `json:"ziggyId"` // workflow ID of the ZiggyWorkflow to query
	Track    string `json:"track"`   // "educational" | "fun"

	// Preserved state across continue-as-new
	RecentMessages  []ChatMessage `json:"recentMessages,omitempty"`  // Last N messages for context
	ActiveMystery   *Mystery      `json:"activeMystery,omitempty"`
	MysteryProgress int           `json:"mysteryProgress,omitempty"`
	HintsGiven      []string      `json:"hintsGiven,omitempty"`
	Solved          []string      `json:"solved,omitempty"`
}

type SendMessageSignal struct {
	Content string `json:"content"`
}

type StartMysterySignal struct {
	MysteryID string `json:"mysteryId"`
	Track     string `json:"track"`
}

func ChatWorkflow(ctx workflow.Context, input ChatInput) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Chat workflow started", "owner", input.Owner, "track", input.Track)

	track := input.Track
	if track == "" {
		track = "fun"
	}

	state := NewChatState(input.Owner)

	// Restore preserved state from continue-as-new
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

	// Query handler for chat history (includes mystery status)
	err := workflow.SetQueryHandler(ctx, QueryChatHistory, func() (ChatHistoryResponse, error) {
		mysteryStatus := state.GetMysteryStatus()
		return ChatHistoryResponse{
			Messages:      state.Messages,
			MysteryStatus: &mysteryStatus,
			IsTyping:      state.IsTyping,
		}, nil
	})
	if err != nil {
		return err
	}

	// Query handler for mystery status
	err = workflow.SetQueryHandler(ctx, QueryMysteryStatus, func() (MysteryStatus, error) {
		return state.GetMysteryStatus(), nil
	})
	if err != nil {
		return err
	}

	messageCh := workflow.GetSignalChannel(ctx, SignalSendMessage)
	mysteryCh := workflow.GetSignalChannel(ctx, SignalStartMystery)

	for {
		selector := workflow.NewSelector(ctx)

		selector.AddReceive(messageCh, func(c workflow.ReceiveChannel, more bool) {
			var signal SendMessageSignal
			c.Receive(ctx, &signal)
			now := workflow.Now(ctx)

			// Add user message
			state.AddMessage("user", signal.Content, now)
			logger.Info("User message received", "content", signal.Content)

			// Use mystery's track if one is active, otherwise use workflow track
			responseTrack := track
			if state.ActiveMystery != nil && state.ActiveMystery.Track != "" {
				responseTrack = state.ActiveMystery.Track
			}

			// For educational track with active topic, show a "searching" message while we query docs
			if responseTrack == "educational" && state.ActiveMystery != nil {
				state.AddMessage("ziggy", "Searching the Temporal docs...", workflow.Now(ctx))
				state.IsTyping = true
			}

			// Get Ziggy's current state for personality context
			ziggyState := queryZiggyState(ctx, input.ZiggyID, logger)

			// Generate response via activity
			response := generateChatResponse(ctx, &state, ziggyState, responseTrack, logger)

			// Add Ziggy's response
			state.AddMessage("ziggy", response, workflow.Now(ctx))
			state.IsTyping = false
			logger.Info("Ziggy response generated", "response", response)
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

		// Continue-as-new if message history gets too long
		if len(state.Messages) >= MaxChatMessages {
			logger.Info("Continuing as new due to message limit")

			// Keep last 20 messages for context
			recentMessages := state.Messages
			if len(recentMessages) > 20 {
				recentMessages = recentMessages[len(recentMessages)-20:]
			}

			return workflow.NewContinueAsNewError(ctx, ChatWorkflow, ChatInput{
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

func queryZiggyState(ctx workflow.Context, ziggyID string, logger interface{ Info(string, ...interface{}) }) *ZiggyState {
	if ziggyID == "" {
		logger.Info("No ZiggyID provided, skipping state query")
		return nil
	}

	// Query the ZiggyWorkflow for current state
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    500 * time.Millisecond,
			BackoffCoefficient: 2.0,
			MaximumInterval:    10 * time.Second,
			MaximumAttempts:    5,
		},
	}
	actCtx := workflow.WithActivityOptions(ctx, ao)

	var state ZiggyState
	err := workflow.ExecuteActivity(actCtx, "QueryZiggyState", ziggyID).Get(ctx, &state)
	if err != nil {
		logger.Info("Failed to query Ziggy state, continuing without it", "error", err.Error())
		return nil
	}

	return &state
}

func generateChatResponse(ctx workflow.Context, chatState *ChatState, ziggyState *ZiggyState, track string, logger interface{ Info(string, ...interface{}) }) string {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 2 * time.Minute,
	}
	actCtx := workflow.WithActivityOptions(ctx, ao)

	input := ChatActivityInput{
		Messages:    chatState.Messages,
		Mystery:     chatState.ActiveMystery,
		Progress:    chatState.MysteryProgress,
		HintsGiven:  chatState.HintsGiven,
		Track:       track,
	}

	if ziggyState != nil {
		input.Personality = ziggyState.Personality
		input.Mood = ziggyState.GetMood()
		input.Stage = GetStageForAge(time.Since(ziggyState.CreatedAt).Seconds())
		input.Bond = ziggyState.Bond
	}

	var output ChatActivityOutput
	err := workflow.ExecuteActivity(actCtx, "GenerateChatResponse", input).Get(ctx, &output)
	if err != nil {
		logger.Info("Failed to generate chat response", "error", err.Error())
		return getFallbackChatResponse(ziggyState)
	}

	// Update mystery progress if applicable
	if output.MysteryUpdate != nil {
		if output.MysteryUpdate.HintGiven != "" {
			chatState.HintsGiven = append(chatState.HintsGiven, output.MysteryUpdate.HintGiven)
		}
		chatState.MysteryProgress = output.MysteryUpdate.NewProgress

		// Cap progress at total hints to prevent overflow
		if chatState.ActiveMystery != nil && chatState.MysteryProgress > len(chatState.ActiveMystery.Hints) {
			chatState.MysteryProgress = len(chatState.ActiveMystery.Hints)
		}

		if output.MysteryUpdate.Solved && chatState.ActiveMystery != nil {
			chatState.Solved = append(chatState.Solved, chatState.ActiveMystery.ID)
			chatState.ActiveMystery = nil
			chatState.MysteryProgress = 0
			chatState.HintsGiven = nil
		}
		if output.MysteryUpdate.Failed && chatState.ActiveMystery != nil {
			// Failed - clear mystery but don't mark as solved (can try again)
			chatState.ActiveMystery = nil
			chatState.MysteryProgress = 0
			chatState.HintsGiven = nil
		}

		// Auto-fail if hints exhausted, user made a guess (no new hint), and AI didn't solve it
		if chatState.ActiveMystery != nil &&
			chatState.MysteryProgress >= len(chatState.ActiveMystery.Hints) &&
			output.MysteryUpdate.HintGiven == "" {
			// All hints used, user made final guess and got it wrong
			solution := chatState.ActiveMystery.Solution
			chatState.ActiveMystery = nil
			chatState.MysteryProgress = 0
			chatState.HintsGiven = nil
			return output.Response + "\n\n*wiggles sympathetically*\nThe answer was: " + solution + "\n\nNice try! Want to try another mystery?"
		}
	}

	return output.Response
}

func getFallbackChatResponse(ziggyState *ZiggyState) string {
	if ziggyState == nil {
		return "*wiggle wiggle*\nHi there!"
	}

	mood := ziggyState.GetMood()
	switch mood {
	case MoodHappy:
		return "*happy wiggle*\nGreat to chat!"
	case MoodSad:
		return "*slow wiggle*\nI'm feeling\na bit down..."
	case MoodHungry:
		return "*tummy rumble*\nSo hungry..."
	case MoodSleeping:
		return "*snore*\nzzz..."
	default:
		return "*wiggle*\nHello!"
	}
}
