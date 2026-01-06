package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

type Client struct {
	client anthropic.Client
	init   bool
}

var (
	instance     *Client
	instanceOnce sync.Once
)

// NewClient returns a lazily-initialized singleton AI client.
// The actual API connection is deferred until first use.
func NewClient() *Client {
	instanceOnce.Do(func() {
		instance = &Client{}
	})
	return instance
}

func (c *Client) ensureInit() bool {
	if c.init {
		return true
	}
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return false
	}
	c.client = anthropic.NewClient(option.WithAPIKey(apiKey))
	c.init = true
	return true
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

	if c == nil || !c.ensureInit() {
		log.Printf("[AI] Client not initialized - missing ANTHROPIC_API_KEY")
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
- Avoid phrases that could sound inappropriate out of context (e.g. "gentle petting")
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
	Failed      bool    `json:"failed"`
	HintGiven   string  `json:"hintGiven,omitempty"`
	NewProgress float64 `json:"newProgress"`
}

func (c *Client) GenerateChat(ctx context.Context, input ChatInput) (*ChatResponse, error) {
	log.Printf("[AI] GenerateChat called: personality=%s mood=%s track=%s",
		input.Personality, input.Mood, input.Track)

	if c == nil || !c.ensureInit() {
		return nil, fmt.Errorf("AI client not initialized")
	}

	// Use web search for educational track
	if input.Track == "educational" {
		return c.generateChatWithWebSearch(ctx, input)
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

// generateChatWithWebSearch uses the Anthropic web search tool to provide
// real-time documentation for educational queries about Temporal.
func (c *Client) generateChatWithWebSearch(ctx context.Context, input ChatInput) (*ChatResponse, error) {
	log.Printf("[AI] Using web search for educational track")

	// Build conversation history for multi-turn context
	messages := buildConversationMessages(input)

	// Create web search tool with Temporal docs domain filtering
	webSearchTool := anthropic.WebSearchTool20250305Param{
		AllowedDomains: []string{"docs.temporal.io", "temporal.io", "learn.temporal.io"},
	}

	message, err := c.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeHaiku4_5,
		MaxTokens: 2048,
		System: []anthropic.TextBlockParam{
			{Text: buildEducationalSystemPrompt(input)},
		},
		Messages: messages,
		Tools: []anthropic.ToolUnionParam{
			{OfWebSearchTool20250305: &webSearchTool},
		},
	})
	if err != nil {
		log.Printf("[AI] Web search API request failed: %v", err)
		return nil, fmt.Errorf("claude API error: %w", err)
	}

	// Extract response text and citations from content blocks
	return parseWebSearchResponse(message)
}

// buildConversationMessages converts chat input to Anthropic message format
func buildConversationMessages(input ChatInput) []anthropic.MessageParam {
	messages := make([]anthropic.MessageParam, 0, len(input.Messages))

	for _, m := range input.Messages {
		if m.Role == "user" {
			messages = append(messages, anthropic.NewUserMessage(anthropic.NewTextBlock(m.Content)))
		} else {
			messages = append(messages, anthropic.NewAssistantMessage(anthropic.NewTextBlock(m.Content)))
		}
	}

	return messages
}

// buildEducationalSystemPrompt creates the system prompt for educational mode
func buildEducationalSystemPrompt(input ChatInput) string {
	topicContext := ""
	if input.Mystery != nil {
		topicContext = fmt.Sprintf(`
Current topic: %s
Description: %s
`, input.Mystery.Title, input.Mystery.Description)
	}

	return fmt.Sprintf(`You are Ziggy, a friendly tardigrade who lives inside a Temporal workflow.
You help developers learn about Temporal concepts by searching the official documentation.

%s

Guidelines:
- Use the web search tool to find accurate, up-to-date information from Temporal's docs
- Explain concepts clearly and concisely (3-5 sentences)
- Include relevant code examples when helpful
- Relate concepts to your own experience as a workflow when appropriate
- Be encouraging and make learning fun
- Never use emoji
- Do NOT include "Learn more" links - citations are added automatically`, topicContext)
}

// parseWebSearchResponse extracts text and citations from the API response
func parseWebSearchResponse(message *anthropic.Message) (*ChatResponse, error) {
	var responseText strings.Builder
	var citations []string
	sawToolUse := false

	log.Printf("[AI] Parsing web search response with %d content blocks, stop_reason: %s", len(message.Content), message.StopReason)

	for i, block := range message.Content {
		log.Printf("[AI] Content block %d type: %s", i, block.Type)
		switch variant := block.AsAny().(type) {
		case anthropic.TextBlock:
			log.Printf("[AI] TextBlock: %s", truncate(variant.Text, 100))
			// Only include text that comes after tool use (skip "I'll search..." preamble)
			if sawToolUse {
				responseText.WriteString(variant.Text)
				// Extract citations from text block
				for _, citation := range variant.Citations {
					log.Printf("[AI] Citation type: %s", citation.Type)
					if webCitation, ok := citation.AsAny().(anthropic.CitationsWebSearchResultLocation); ok {
						citations = append(citations, webCitation.URL)
					}
				}
			}
		case anthropic.ServerToolUseBlock:
			log.Printf("[AI] ServerToolUseBlock: name=%s id=%s", variant.Name, variant.ID)
			sawToolUse = true
		case anthropic.WebSearchToolResultBlock:
			log.Printf("[AI] WebSearchToolResultBlock: tool_use_id=%s", variant.ToolUseID)
		default:
			log.Printf("[AI] Unknown block type: %T", variant)
		}
	}

	text := responseText.String()
	if text == "" {
		return nil, fmt.Errorf("empty response from web search")
	}

	// Append unique citations as markdown list under "Learn more" heading
	if len(citations) > 0 {
		seen := make(map[string]bool)
		var uniqueCitations []string
		for _, c := range citations {
			if !seen[c] {
				seen[c] = true
				uniqueCitations = append(uniqueCitations, "- "+c)
			}
		}
		text += "\n\n## Learn more\n" + strings.Join(uniqueCitations, "\n")
	}

	log.Printf("[AI] Web search response: %s", truncate(text, 200))

	return &ChatResponse{
		Response: text,
		MysteryUpdate: &ChatMysteryUpdate{
			Solved: true, // Educational topics are "solved" after explanation
		},
	}, nil
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
		// Fun track: guessing game with riddles
		// (Educational track uses generateChatWithWebSearch instead)
		hintsExhausted := input.Mystery.Progress >= len(input.Mystery.Hints)
		nextHint := ""
		if !hintsExhausted {
			nextHint = input.Mystery.Hints[input.Mystery.Progress]
		}
		conceptHint := ""
		if input.Mystery.Concept != "" {
			conceptHint = fmt.Sprintf("(The answer relates to: %s)\n", input.Mystery.Concept)
		}

		exhaustedSection := ""
		if hintsExhausted {
			exhaustedSection = `
*** ALL HINTS EXHAUSTED - CHALLENGE OVER ***
The user has used all hints and hasn't solved it. You MUST:
1. Kindly reveal the answer: "The answer was [solution]!"
2. Be encouraging: "Nice try! You were getting close. Want to try another mystery?"
3. Set failed=true AND solved=false in the JSON response
4. Do NOT give any more hints (set hintGiven to empty string)
`
		}

		mysterySection = fmt.Sprintf(`
MYSTERY MODE - You are playing a guessing game with the user!

The mystery: "%s"
Your riddle to them: "%s"
%sHints given so far: %d of %d
Next hint (if they need help): %s
The answer they must guess: %s
%s
IMPORTANT RULES FOR MYSTERY MODE:
1. If this is the START of the mystery (no hints given yet), present your riddle excitedly, ask them to guess, and remind them they can ask for a hint if stumped
2. The user must GUESS the answer - never reveal it directly (unless all hints exhausted and they fail)!
3. If they guess wrong, encourage them and offer a hint
4. If they seem stuck or ask for help, give the next hint naturally
5. If they guess correctly (mention the concept or solution), celebrate and set solved=true
6. Keep it fun and playful - you're excited to share this puzzle!
`,
			input.Mystery.Title,
			input.Mystery.Description,
			conceptHint,
			input.Mystery.Progress,
			len(input.Mystery.Hints),
			nextHint,
			input.Mystery.Solution,
			exhaustedSection,
		)
	}

	responseFormat := `Respond as Ziggy in 2-4 short sentences. Keep responses under 200 characters total.`
	if input.Mystery != nil {
		responseFormat = `Respond as JSON: {"response": "your message", "mysteryUpdate": {"solved": false, "failed": false, "hintGiven": "hint if given", "newProgress": 0}}`
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
