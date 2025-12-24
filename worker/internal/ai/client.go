package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

type Client struct {
	client anthropic.Client
}

func NewClient() *Client {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return nil
	}

	client := anthropic.NewClient(option.WithAPIKey(apiKey))
	return &Client{client: client}
}

type PoolGenerationInput struct {
	Personality     string
	Stage           string
	BondDescription string
}

type MessagePool struct {
	FeedSuccess  []string `json:"feedSuccess"`
	FeedFull     []string `json:"feedFull"`
	FeedHungry   []string `json:"feedHungry"`
	FeedSleeping []string `json:"feedSleeping"`
	FeedTun      []string `json:"feedTun"`
	FeedCooldown []string `json:"feedCooldown"`

	PlaySuccess  []string `json:"playSuccess"`
	PlayTired    []string `json:"playTired"`
	PlayHappy    []string `json:"playHappy"`
	PlaySleeping []string `json:"playSleeping"`
	PlayTun      []string `json:"playTun"`
	PlayCooldown []string `json:"playCooldown"`

	PetSuccess   []string `json:"petSuccess"`
	PetMaxBond   []string `json:"petMaxBond"`
	PetLowMood   []string `json:"petLowMood"`
	PetSleeping  []string `json:"petSleeping"`
	PetTun       []string `json:"petTun"`
	PetCooldown  []string `json:"petCooldown"`

	Reviving []string `json:"reviving"`

	IdleHappy    []string `json:"idleHappy"`
	IdleNeutral  []string `json:"idleNeutral"`
	IdleHungry   []string `json:"idleHungry"`
	IdleSad      []string `json:"idleSad"`
	IdleLonely   []string `json:"idleLonely"`
	IdleCritical []string `json:"idleCritical"`
	IdleTun      []string `json:"idleTun"`
	IdleSleeping []string `json:"idleSleeping"`
}

func (c *Client) GeneratePool(ctx context.Context, input PoolGenerationInput) (*MessagePool, error) {
	if c == nil {
		return nil, fmt.Errorf("AI client not initialized (missing ANTHROPIC_API_KEY)")
	}

	prompt := buildPrompt(input)

	message, err := c.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaude3_5Haiku20241022,
		MaxTokens: 4096,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("claude API error: %w", err)
	}

	if len(message.Content) == 0 {
		return nil, fmt.Errorf("empty response from claude")
	}

	text := message.Content[0].Text
	if text == "" {
		return nil, fmt.Errorf("empty text in response from claude")
	}

	jsonStr := extractJSON(text)
	if jsonStr == "" {
		return nil, fmt.Errorf("no JSON found in response")
	}

	var pool MessagePool
	if err := json.Unmarshal([]byte(jsonStr), &pool); err != nil {
		return nil, fmt.Errorf("failed to parse pool JSON: %w", err)
	}

	return &pool, nil
}

func buildPrompt(input PoolGenerationInput) string {
	return fmt.Sprintf(`You are generating dialogue for Ziggy, a tardigrade virtual pet.

Personality: %s
Life stage: %s
Bond level: %s

Generate 10 short messages (max 3 lines, ~20 chars each) for each category below.

Categories:
- feedSuccess: Successfully fed when hungry/neutral
- feedFull: Overfed (already full)
- feedHungry: Fed when very hungry
- feedSleeping: Tried to feed while sleeping
- feedTun: Fed while in tun/dormant state (helps revival)
- feedCooldown: Fed too soon after last feeding
- playSuccess: Successfully played
- playTired: Too tired to play properly
- playHappy: Playing while already happy
- playSleeping: Tried to play while sleeping
- playTun: Tried to play while dormant
- playCooldown: Played too soon after last play
- petSuccess: Successfully petted
- petMaxBond: Petted when bond is maxed
- petLowMood: Petted when sad/hungry (comfort)
- petSleeping: Petted while sleeping
- petTun: Petted while dormant (helps revival)
- petCooldown: Petted too soon after last pet
- reviving: Waking up from tun/dormant state
- idleHappy: Idle dialogue when happy
- idleNeutral: Idle dialogue when neutral
- idleHungry: Idle dialogue when hungry
- idleSad: Idle dialogue when sad
- idleLonely: Idle dialogue when bond is low
- idleCritical: Idle dialogue when HP is critical
- idleTun: Idle dialogue when dormant
- idleSleeping: Idle dialogue when sleeping

Rules:
- Never use emoji
- Match the %s personality voice consistently
- Reference tardigrade facts occasionally (survive space, radiation, extreme temps, etc.)
- Each message should be max 3 lines, each line ~20 characters
- Use \n for line breaks within messages
- Keep messages appropriate for the context

Return ONLY a valid JSON object matching this structure (no markdown, no explanation):
{
  "feedSuccess": ["msg1", "msg2", ...],
  "feedFull": ["msg1", "msg2", ...],
  ... (all categories)
}`, input.Personality, input.Stage, input.BondDescription, input.Personality)
}

func extractJSON(text string) string {
	start := -1
	depth := 0

	for i, c := range text {
		if c == '{' {
			if start == -1 {
				start = i
			}
			depth++
		} else if c == '}' {
			depth--
			if depth == 0 && start != -1 {
				return text[start : i+1]
			}
		}
	}
	return ""
}
