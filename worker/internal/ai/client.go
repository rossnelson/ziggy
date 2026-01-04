package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

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

	PetSuccess  []string `json:"petSuccess"`
	PetMaxBond  []string `json:"petMaxBond"`
	PetLowMood  []string `json:"petLowMood"`
	PetSleeping []string `json:"petSleeping"`
	PetTun      []string `json:"petTun"`
	PetCooldown []string `json:"petCooldown"`

	Reviving []string `json:"reviving"`

	IdleHappy    []string `json:"idleHappy"`
	IdleNeutral  []string `json:"idleNeutral"`
	IdleHungry   []string `json:"idleHungry"`
	IdleSad      []string `json:"idleSad"`
	IdleLonely   []string `json:"idleLonely"`
	IdleCritical []string `json:"idleCritical"`
	IdleTun      []string `json:"idleTun"`
	IdleSleeping []string `json:"idleSleeping"`

	// Need-based coaxing messages
	NeedsFood      []string `json:"needsFood"`
	NeedsPlay      []string `json:"needsPlay"`
	NeedsAffection []string `json:"needsAffection"`
	NeedsCritical  []string `json:"needsCritical"`
}

func (c *Client) GeneratePool(ctx context.Context, input PoolGenerationInput) (*MessagePool, error) {
	log.Printf("[AI] GeneratePool called: personality=%s stage=%s bond=%s",
		input.Personality, input.Stage, input.BondDescription)

	if c == nil {
		log.Printf("[AI] Client is nil - missing ANTHROPIC_API_KEY")
		return nil, fmt.Errorf("AI client not initialized (missing ANTHROPIC_API_KEY)")
	}

	prompt := buildPrompt(input)
	log.Printf("[AI] Sending request to Claude API (prompt length: %d chars)", len(prompt))

	message, err := c.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaude3_5Haiku20241022,
		MaxTokens: 4096,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	})
	if err != nil {
		log.Printf("[AI] Claude API request failed: %v", err)
		return nil, fmt.Errorf("claude API error: %w", err)
	}

	log.Printf("[AI] Received response from Claude (content blocks: %d)", len(message.Content))

	if len(message.Content) == 0 {
		log.Printf("[AI] Empty response from Claude")
		return nil, fmt.Errorf("empty response from claude")
	}

	text := message.Content[0].Text
	if text == "" {
		log.Printf("[AI] Empty text in response")
		return nil, fmt.Errorf("empty text in response from claude")
	}

	log.Printf("[AI] Response text length: %d chars", len(text))

	// First try parsing the entire response as JSON (common case)
	var pool MessagePool
	text = strings.TrimSpace(text)
	if err := json.Unmarshal([]byte(text), &pool); err == nil {
		log.Printf("[AI] Parsed response directly as JSON")
		return &pool, nil
	} else {
		log.Printf("[AI] Direct JSON parse failed: %v", err)
	}

	// Fall back to extracting JSON from mixed content
	jsonStr := extractJSON(text)
	if jsonStr == "" {
		log.Printf("[AI] No JSON found in response. First 500 chars: %s", truncate(text, 500))
		return nil, fmt.Errorf("no JSON found in response")
	}

	log.Printf("[AI] Extracted JSON length: %d chars", len(jsonStr))

	if err := json.Unmarshal([]byte(jsonStr), &pool); err != nil {
		log.Printf("[AI] JSON parse error: %v. First 500 chars of JSON: %s", err, truncate(jsonStr, 500))
		return nil, fmt.Errorf("failed to parse pool JSON: %w", err)
	}

	log.Printf("[AI] Successfully parsed pool")
	return &pool, nil
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
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
- needsFood: Coaxing messages when hungry (gently ask for food)
- needsPlay: Coaxing messages when bored (gently ask for play)
- needsAffection: Coaxing messages when lonely (gently ask for pets)
- needsCritical: Urgent messages when HP is low (plead for help)

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

// stripJSON removes JSON objects from text, returning the remaining content
func stripJSON(text string) string {
	result := []byte(text)
	for {
		start := -1
		depth := 0
		end := -1

		for i, c := range string(result) {
			if c == '{' {
				if start == -1 {
					start = i
				}
				depth++
			} else if c == '}' {
				depth--
				if depth == 0 && start != -1 {
					end = i + 1
					break
				}
			}
		}

		if start == -1 || end == -1 {
			break
		}

		// Remove the JSON block
		result = append(result[:start], result[end:]...)
	}

	// Clean up whitespace
	cleaned := strings.TrimSpace(string(result))
	return cleaned
}

// Chat types and methods

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type MysteryContext struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Concept     string   `json:"concept,omitempty"`
	Hints       []string `json:"hints"`
	HintsGiven  []string `json:"hintsGiven"`
	Progress    int      `json:"progress"`
	Solution    string   `json:"solution"`
	Summary     string   `json:"summary,omitempty"`
	DocsURL     string   `json:"docsUrl,omitempty"`
}

type ChatInput struct {
	Messages    []ChatMessage   `json:"messages"`
	Personality string          `json:"personality"`
	Mood        string          `json:"mood"`
	Stage       string          `json:"stage"`
	Bond        float64         `json:"bond"`
	Track       string          `json:"track"`
	Mystery     *MysteryContext `json:"mystery,omitempty"`
}

type ChatResponse struct {
	Response      string             `json:"response"`
	MysteryUpdate *ChatMysteryUpdate `json:"mysteryUpdate,omitempty"`
}

type ChatMysteryUpdate struct {
	Solved      bool    `json:"solved"`
	HintGiven   string  `json:"hintGiven,omitempty"`
	NewProgress float64 `json:"newProgress"`
}

func (c *Client) GenerateChat(ctx context.Context, input ChatInput) (*ChatResponse, error) {
	log.Printf("[AI] GenerateChat called: personality=%s mood=%s track=%s",
		input.Personality, input.Mood, input.Track)

	if c == nil {
		return nil, fmt.Errorf("AI client not initialized")
	}

	prompt := buildChatPrompt(input)
	log.Printf("[AI] Chat prompt length: %d chars", len(prompt))

	message, err := c.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaude3_5Haiku20241022,
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	})
	if err != nil {
		log.Printf("[AI] Chat API request failed: %v", err)
		return nil, fmt.Errorf("claude API error: %w", err)
	}

	if len(message.Content) == 0 || message.Content[0].Text == "" {
		return nil, fmt.Errorf("empty response from claude")
	}

	text := message.Content[0].Text
	log.Printf("[AI] Chat response: %s", truncate(text, 200))

	// Try to parse as JSON first (for mystery updates)
	jsonStr := extractJSON(text)
	log.Printf("[AI] extractJSON result: found=%v len=%d", jsonStr != "", len(jsonStr))
	if jsonStr != "" {
		var response ChatResponse
		if err := json.Unmarshal([]byte(jsonStr), &response); err == nil {
			log.Printf("[AI] JSON unmarshal success: response=%q mysteryUpdate=%v", truncate(response.Response, 50), response.MysteryUpdate != nil)
			// Ensure we have a valid response field
			if response.Response != "" {
				return &response, nil
			}
			// JSON parsed but response empty - try to extract from nested structure
			log.Printf("[AI] JSON parsed but response field empty")
		} else {
			log.Printf("[AI] JSON extraction found but parse failed: %v", err)
		}
	}

	// Plain text response - strip any JSON that might be embedded
	cleanText := stripJSON(text)
	if cleanText == "" {
		cleanText = text // fallback to original if stripping removed everything
	}
	return &ChatResponse{Response: cleanText}, nil
}

func buildChatPrompt(input ChatInput) string {
	bondDesc := getBondDescription(input.Bond)

	// Build conversation history
	history := ""
	for _, m := range input.Messages {
		if m.Role == "user" {
			history += fmt.Sprintf("User: %s\n", m.Content)
		} else {
			history += fmt.Sprintf("Ziggy: %s\n", m.Content)
		}
	}

	mysterySection := ""
	if input.Mystery != nil {
		if input.Track == "educational" {
			// Educational track: provide summary and link to docs
			mysterySection = fmt.Sprintf(`
LEARNING MODE - Teaching about: %s

Summary to paraphrase:
%s

YOUR RESPONSE MUST END WITH THIS EXACT LINE:
Learn more: %s

Keep explanation to 2-3 sentences, then add the learn more link.
Set solved=true in JSON.
`,
				input.Mystery.Title,
				input.Mystery.Summary,
				input.Mystery.DocsURL,
			)
		} else {
			// Fun track: guessing game with riddles
			nextHint := ""
			if input.Mystery.Progress < len(input.Mystery.Hints) {
				nextHint = input.Mystery.Hints[input.Mystery.Progress]
			}
			conceptHint := ""
			if input.Mystery.Concept != "" {
				conceptHint = fmt.Sprintf("(The answer relates to: %s)\n", input.Mystery.Concept)
			}
			mysterySection = fmt.Sprintf(`
MYSTERY MODE - You are playing a guessing game with the user!

The mystery: "%s"
Your riddle to them: "%s"
%sHints given so far: %v
Next hint (if they need help): %s
The answer they must guess: %s

IMPORTANT RULES FOR MYSTERY MODE:
1. If this is the START of the mystery (no hints given yet), present your riddle excitedly, ask them to guess, and remind them they can ask for a hint if stumped
2. The user must GUESS the answer - never reveal it directly!
3. If they guess wrong, encourage them and offer a hint
4. If they seem stuck or ask for help, give the next hint naturally
5. If they guess correctly (mention the concept or solution), celebrate and set solved=true
6. Keep it fun and playful - you're excited to share this puzzle!
`,
				input.Mystery.Title,
				input.Mystery.Description,
				conceptHint,
				input.Mystery.HintsGiven,
				nextHint,
				input.Mystery.Solution,
			)
		}
	}

	responseFormat := `Respond as Ziggy in 2-4 short sentences. Keep responses under 200 characters total.`
	if input.Mystery != nil {
		responseFormat = `Respond as JSON: {"response": "your message", "mysteryUpdate": {"solved": false, "hintGiven": "hint if given", "newProgress": 0}}`
	}

	// Different tone for educational vs fun track
	if input.Track == "educational" {
		return fmt.Sprintf(`You are Ziggy, an educational guide teaching Temporal workflow concepts. You live inside a Temporal workflow yourself, which gives you firsthand experience.

%s
Conversation so far:
%s

Style Rules:
- Be clear, direct, and informative - like a friendly instructor
- Use precise technical terminology
- Skip the quirky personality traits - focus on teaching
- Keep responses concise but thorough (3-5 sentences)
- Use concrete examples from how YOU work as a workflow
- Never use emoji or cutesy expressions

%s`,
			mysterySection,
			history,
			responseFormat,
		)
	}

	return fmt.Sprintf(`You are Ziggy, a tardigrade virtual pet living in a Temporal workflow.

Personality: %s
Current mood: %s
Bond level: %s
Life stage: %s
Track: %s
%s
Conversation so far:
%s

Rules:
- Stay in character as a %s tardigrade
- Keep responses short (2-4 sentences, max 200 chars)
- Reference tardigrade facts occasionally (survive space, radiation, extreme temps)
- Match your mood to current state
- Never use emoji

%s`,
		input.Personality,
		input.Mood,
		bondDesc,
		input.Stage,
		input.Track,
		mysterySection,
		history,
		input.Personality,
		responseFormat,
	)
}

func getBondDescription(bond float64) string {
	if bond >= 80 {
		return "deeply bonded (best friends)"
	}
	if bond >= 60 {
		return "close bond (good friends)"
	}
	if bond >= 40 {
		return "developing bond (getting to know each other)"
	}
	if bond >= 20 {
		return "new acquaintance (still shy)"
	}
	return "barely met (very timid)"
}
