package chat

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"ziggy/internal/ai"
	"ziggy/internal/registry"
	z "ziggy/internal/ziggy"
)

type Activities struct {
	aiClient *ai.Client
}

func NewActivities(aiClient *ai.Client) *Activities {
	return &Activities{aiClient: aiClient}
}

const QueryZiggyState = "state"

func (a *Activities) QueryZiggyState(ctx context.Context, ziggyID string) (*z.State, error) {
	log.Printf("[ChatActivity] Querying Ziggy state for workflow: %s", ziggyID)

	result, err := registry.Get().QueryWorkflow(ctx, ziggyID, QueryZiggyState)
	if err != nil {
		log.Printf("[ChatActivity] Failed to query Ziggy state: %v", err)
		return nil, err
	}

	data, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	var state z.State
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}

	log.Printf("[ChatActivity] Got Ziggy state: personality=%s mood=%s", state.Personality, state.GetMood())
	return &state, nil
}

type ProcessMessageInput struct {
	State      State     `json:"state"`
	Content    string    `json:"content"`
	ZiggyState *z.State  `json:"ziggyState,omitempty"`
	Track      string    `json:"track"`
	Now        time.Time `json:"now"`
}

type ProcessMessageOutput struct {
	State State `json:"state"`
}

func (a *Activities) ProcessChatMessage(ctx context.Context, input ProcessMessageInput) (*ProcessMessageOutput, error) {
	state := input.State
	now := input.Now

	state.AddMessage("user", input.Content, now)
	log.Printf("[ChatActivity] User message received: %s", input.Content)

	responseTrack := input.Track
	if state.ActiveMystery != nil && state.ActiveMystery.Track != "" {
		responseTrack = state.ActiveMystery.Track
	}

	response := a.generateResponse(ctx, &state, input.ZiggyState, responseTrack)

	a.processMysteryUpdate(&state, &response)

	state.AddMessage("ziggy", response.Response, now)
	state.IsTyping = false

	return &ProcessMessageOutput{State: state}, nil
}

type chatResponse struct {
	Response      string
	MysteryUpdate *MysteryUpdate
}

type MysteryUpdate struct {
	Solved      bool   `json:"solved"`
	Failed      bool   `json:"failed"`
	HintGiven   string `json:"hintGiven,omitempty"`
	NewProgress int    `json:"newProgress"`
}

func (a *Activities) generateResponse(ctx context.Context, chatState *State, ziggyState *z.State, track string) chatResponse {
	if a.aiClient == nil {
		return chatResponse{Response: getFallbackResponse(ziggyState)}
	}

	aiInput := ai.ChatInput{
		Messages: convertMessages(chatState.Messages),
		Track:    track,
	}

	if ziggyState != nil {
		aiInput.Personality = string(ziggyState.Personality)
		aiInput.Mood = string(ziggyState.GetMood())
		aiInput.Stage = string(z.GetStageForAge(time.Since(ziggyState.CreatedAt).Seconds()))
		aiInput.Bond = ziggyState.Bond
	}

	if chatState.ActiveMystery != nil {
		aiInput.Mystery = &ai.MysteryContext{
			Title:       chatState.ActiveMystery.Title,
			Description: chatState.ActiveMystery.Description,
			Concept:     chatState.ActiveMystery.Concept,
			Hints:       chatState.ActiveMystery.Hints,
			HintsGiven:  chatState.HintsGiven,
			Progress:    chatState.MysteryProgress,
			Solution:    chatState.ActiveMystery.Solution,
			Summary:     chatState.ActiveMystery.Summary,
		}
	}

	result, err := a.aiClient.GenerateChat(ctx, aiInput)
	if err != nil {
		log.Printf("[ChatActivity] AI error: %v, using fallback", err)
		return chatResponse{Response: getFallbackResponse(ziggyState)}
	}

	resp := chatResponse{Response: result.Response}

	if result.MysteryUpdate != nil {
		newProgress := chatState.MysteryProgress
		if result.MysteryUpdate.HintGiven != "" {
			newProgress++
		}
		resp.MysteryUpdate = &MysteryUpdate{
			Solved:      result.MysteryUpdate.Solved,
			Failed:      result.MysteryUpdate.Failed,
			HintGiven:   result.MysteryUpdate.HintGiven,
			NewProgress: newProgress,
		}
	}

	return resp
}

func (a *Activities) processMysteryUpdate(state *State, resp *chatResponse) {
	if resp.MysteryUpdate == nil {
		return
	}

	update := resp.MysteryUpdate

	if update.HintGiven != "" {
		state.HintsGiven = append(state.HintsGiven, update.HintGiven)
	}
	state.MysteryProgress = update.NewProgress

	if state.ActiveMystery != nil && state.MysteryProgress > len(state.ActiveMystery.Hints) {
		state.MysteryProgress = len(state.ActiveMystery.Hints)
	}

	if update.Solved && state.ActiveMystery != nil {
		state.Solved = append(state.Solved, state.ActiveMystery.ID)
		state.ActiveMystery = nil
		state.MysteryProgress = 0
		state.HintsGiven = nil
	}

	if update.Failed && state.ActiveMystery != nil {
		state.ActiveMystery = nil
		state.MysteryProgress = 0
		state.HintsGiven = nil
	}

	if state.ActiveMystery != nil &&
		state.MysteryProgress >= len(state.ActiveMystery.Hints) &&
		update.HintGiven == "" {
		solution := state.ActiveMystery.Solution
		state.ActiveMystery = nil
		state.MysteryProgress = 0
		state.HintsGiven = nil
		resp.Response = resp.Response + "\n\n*wiggles sympathetically*\nThe answer was: " + solution + "\n\nNice try! Want to try another mystery?"
	}
}

func convertMessages(messages []Message) []ai.ChatMessage {
	result := make([]ai.ChatMessage, len(messages))
	for i, m := range messages {
		result[i] = ai.ChatMessage{
			Role:    m.Role,
			Content: m.Content,
		}
	}
	return result
}

func getFallbackResponse(ziggyState *z.State) string {
	if ziggyState == nil {
		return "*wiggle wiggle*\nHi there!"
	}

	mood := ziggyState.GetMood()
	switch mood {
	case z.MoodHappy:
		return "*happy wiggle*\nGreat to chat!"
	case z.MoodSad:
		return "*slow wiggle*\nI'm feeling\na bit down..."
	case z.MoodHungry:
		return "*tummy rumble*\nSo hungry..."
	case z.MoodSleeping:
		return "*snore*\nzzz..."
	default:
		return "*wiggle*\nHello!"
	}
}

