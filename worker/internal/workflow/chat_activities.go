package workflow

import (
	"context"
	"encoding/json"
	"log"

	"ziggy/internal/ai"
	"ziggy/internal/temporal"
)

type ChatActivities struct {
	aiClient *ai.Client
	registry *temporal.Registry
}

func NewChatActivities(aiClient *ai.Client, registry *temporal.Registry) *ChatActivities {
	return &ChatActivities{aiClient: aiClient, registry: registry}
}

func (a *ChatActivities) QueryZiggyState(ctx context.Context, ziggyID string) (*ZiggyState, error) {
	log.Printf("[ChatActivity] Querying Ziggy state for workflow: %s", ziggyID)

	result, err := a.registry.QueryWorkflow(ctx, ziggyID, QueryState)
	if err != nil {
		log.Printf("[ChatActivity] Failed to query Ziggy state: %v", err)
		return nil, err
	}

	data, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	var state ZiggyState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}

	log.Printf("[ChatActivity] Got Ziggy state: personality=%s mood=%s", state.Personality, state.GetMood())
	return &state, nil
}

type ChatActivityInput struct {
	Messages    []ChatMessage `json:"messages"`
	Personality Personality   `json:"personality"`
	Mood        Mood          `json:"mood"`
	Stage       Stage         `json:"stage"`
	Bond        float64       `json:"bond"`
	Mystery     *Mystery      `json:"mystery,omitempty"`
	Progress    int           `json:"progress"`
	HintsGiven  []string      `json:"hintsGiven"`
	Track       string        `json:"track"`
}

type ChatActivityOutput struct {
	Response      string         `json:"response"`
	MysteryUpdate *MysteryUpdate `json:"mysteryUpdate,omitempty"`
}

type MysteryUpdate struct {
	Solved      bool   `json:"solved"`
	HintGiven   string `json:"hintGiven,omitempty"`
	NewProgress int    `json:"newProgress"`
}

func (a *ChatActivities) GenerateChatResponse(ctx context.Context, input ChatActivityInput) (*ChatActivityOutput, error) {
	log.Printf("[ChatActivity] Generating response: personality=%s mood=%s track=%s",
		input.Personality, input.Mood, input.Track)

	if a.aiClient == nil {
		log.Printf("[ChatActivity] No AI client, using fallback")
		return &ChatActivityOutput{
			Response: getFallbackResponse(input),
		}, nil
	}

	aiInput := ai.ChatInput{
		Messages:    convertMessages(input.Messages),
		Personality: string(input.Personality),
		Mood:        string(input.Mood),
		Stage:       string(input.Stage),
		Bond:        input.Bond,
		Track:       input.Track,
	}

	if input.Mystery != nil {
		aiInput.Mystery = &ai.MysteryContext{
			Title:       input.Mystery.Title,
			Description: input.Mystery.Description,
			Concept:     input.Mystery.Concept,
			Hints:       input.Mystery.Hints,
			HintsGiven:  input.HintsGiven,
			Progress:    input.Progress,
			Solution:    input.Mystery.Solution,
		}
	}

	response, err := a.aiClient.GenerateChat(ctx, aiInput)
	if err != nil {
		log.Printf("[ChatActivity] AI error: %v, using fallback", err)
		return &ChatActivityOutput{
			Response: getFallbackResponse(input),
		}, nil
	}

	output := &ChatActivityOutput{
		Response: response.Response,
	}

	if response.MysteryUpdate != nil {
		// AI returns float progress (0.0-1.0), but we track hint count as int
		// Calculate new progress based on hints given
		newProgress := input.Progress
		if response.MysteryUpdate.HintGiven != "" {
			newProgress++
		}
		output.MysteryUpdate = &MysteryUpdate{
			Solved:      response.MysteryUpdate.Solved,
			HintGiven:   response.MysteryUpdate.HintGiven,
			NewProgress: newProgress,
		}
	}

	log.Printf("[ChatActivity] Generated response: %s", response.Response)
	return output, nil
}

func convertMessages(messages []ChatMessage) []ai.ChatMessage {
	result := make([]ai.ChatMessage, len(messages))
	for i, m := range messages {
		result[i] = ai.ChatMessage{
			Role:    m.Role,
			Content: m.Content,
		}
	}
	return result
}

func getFallbackResponse(input ChatActivityInput) string {
	switch input.Mood {
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
