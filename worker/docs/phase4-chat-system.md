# Phase 4: Agentic Storytelling & Chat System

## Overview

Add an interactive chat system where users can converse with Ziggy to solve mysteries. Two tracks available via env var:
- **Educational track**: Learn Temporal concepts through mystery-solving
- **Fun track**: Whimsical adventures and puzzles

## Core Features

| Feature | Description |
|---------|-------------|
| Chat workflow | Separate workflow storing conversation history |
| Personality sync | Chat voice matches tamagotchi personality/state |
| Mystery system | Multi-step storylines with hints |
| Track toggle | `ZIGGY_TRACK=educational\|fun` env var |

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         Frontend                                 │
├───────────────────────────────┬─────────────────────────────────┤
│      Tamagotchi Panel         │         Chat Panel              │
│  ┌─────────────────────┐      │   ┌─────────────────────────┐   │
│  │                     │      │   │ Message history         │   │
│  │      [Ziggy]        │      │   │ - Ziggy: *wiggle*       │   │
│  │                     │      │   │ - User: What's up?      │   │
│  │   HP: ████░░ 70%    │      │   │ - Ziggy: I found...     │   │
│  └─────────────────────┘      │   ├─────────────────────────┤   │
│  [Feed] [Play] [Pet]          │   │ [Type message...]  [>]  │   │
└───────────────────────────────┴─────────────────────────────────┘
                │                               │
                ▼                               ▼
┌───────────────────────────────┐   ┌─────────────────────────────┐
│      ZiggyWorkflow            │   │      ChatWorkflow           │
│  (existing - pet state)       │◄──┤  (new - conversation)       │
│                               │   │                             │
│  - Fullness, HP, Bond         │   │  - Messages[]               │
│  - Personality                │   │  - ActiveMystery            │
│  - Stage                      │   │  - MysteryProgress          │
└───────────────────────────────┘   └─────────────────────────────┘
                                                │
                                                ▼
                                    ┌─────────────────────────────┐
                                    │      Claude API             │
                                    │  (chat activity)            │
                                    │                             │
                                    │  - Generate responses       │
                                    │  - Evaluate mystery progress│
                                    │  - Pick next hints          │
                                    └─────────────────────────────┘
```

## Data Structures

### Chat Message
```go
type ChatMessage struct {
    ID        string    `json:"id"`
    Role      string    `json:"role"` // "user" | "ziggy"
    Content   string    `json:"content"`
    Timestamp time.Time `json:"timestamp"`
}
```

### Mystery
```go
type Mystery struct {
    ID          string   `json:"id"`
    Title       string   `json:"title"`
    Description string   `json:"description"`
    Track       string   `json:"track"` // "educational" | "fun"
    Hints       []string `json:"hints"`
    Solution    string   `json:"solution"`
    Concept     string   `json:"concept,omitempty"` // Temporal concept for educational
}
```

### Chat State
```go
type ChatState struct {
    Messages       []ChatMessage `json:"messages"`
    ActiveMystery  *Mystery      `json:"activeMystery,omitempty"`
    MysteryProgress int          `json:"mysteryProgress"` // hints revealed
    HintsGiven     []string      `json:"hintsGiven"`
    Solved         []string      `json:"solved"` // mystery IDs
}
```

## Educational Track Mysteries

| Mystery | Temporal Concept | Example Hint |
|---------|------------------|--------------|
| The Vanishing Signal | Signals & Queries | "Messages can arrive even when I'm busy..." |
| The Eternal Loop | Continue-as-new | "Sometimes to go forward, you start fresh..." |
| The Sleeping Worker | Activity heartbeats | "Even sleeping, I send signs of life..." |
| The Time Traveler | Workflow replay | "What if every moment could be replayed exactly?" |
| The Parallel Paths | Child workflows | "I can be in many places at once..." |

## Fun Track Mysteries

| Mystery | Theme | Example Hint |
|---------|-------|--------------|
| The Missing Snack | Food heist | "The crumbs lead somewhere cold..." |
| The Cosmic Radio | Space signals | "Beep... boop... is anyone out there?" |
| The Dream Maze | Abstract puzzle | "In dreams, time flows differently..." |

## API Endpoints

### New Endpoints
```
POST /api/ziggy/{id}/chat          - Send message, get response
GET  /api/ziggy/{id}/chat          - Get chat history
POST /api/ziggy/{id}/mystery/start - Begin a mystery
GET  /api/ziggy/{id}/mystery       - Get current mystery status
```

## Workflow Signals & Queries

### ChatWorkflow
```go
// Signals
SendMessage(content string)      // User sends message
StartMystery(mysteryID string)   // Begin a mystery
RevealHint()                     // Request next hint

// Queries
GetHistory() []ChatMessage       // Get message history
GetMysteryStatus() *MysteryStatus // Current mystery state
```

## Files to Create

| File | Purpose |
|------|---------|
| `internal/workflow/chat.go` | ChatWorkflow definition |
| `internal/workflow/chat_state.go` | ChatState, ChatMessage types |
| `internal/workflow/mysteries.go` | Mystery definitions (embedded) |
| `internal/workflow/chat_activities.go` | Claude chat activity |
| `internal/api/chat_handlers.go` | HTTP handlers for chat |

## Files to Modify

| File | Changes |
|------|---------|
| `internal/workflow/register.go` | Register ChatWorkflow |
| `internal/ai/client.go` | Add GenerateChat method |
| `cmd/serve.go` | Add ZIGGY_TRACK env var |

## Claude Chat Activity

```go
type ChatInput struct {
    Messages    []ChatMessage
    Personality Personality
    Mood        Mood
    Stage       Stage
    Bond        float64
    Mystery     *Mystery
    Progress    int
    Track       string
}

type ChatOutput struct {
    Response      string `json:"response"`
    MysteryUpdate *MysteryUpdate `json:"mysteryUpdate,omitempty"`
}

type MysteryUpdate struct {
    Solved      bool   `json:"solved"`
    HintGiven   string `json:"hintGiven,omitempty"`
    NewProgress int    `json:"newProgress"`
}
```

### Chat Prompt Template
```
You are Ziggy, a tardigrade virtual pet.

Personality: {personality}
Current mood: {mood}
Bond level: {bond_description}
Life stage: {stage}

{if mystery}
You are guiding the user through a mystery:
Title: {mystery.title}
Description: {mystery.description}
{if educational}Temporal concept: {mystery.concept}{/if}
Hints given so far: {hintsGiven}
Next available hint: {mystery.hints[progress]}
Solution (don't reveal directly): {mystery.solution}
{/if}

Conversation so far:
{messages}

User: {latest_message}

Rules:
- Stay in character as a {personality} tardigrade
- Keep responses short (2-4 sentences)
- If mystery active, weave in hints naturally based on user questions
- Never reveal solution directly
- Reference tardigrade facts occasionally
- Match your mood to current state
```

## Implementation Order

### Phase 4a: Chat Infrastructure
1. **Chat state types** (`chat_state.go`)
2. **ChatWorkflow** (`chat.go`) - signals, queries, basic loop
3. **Register workflow** (`register.go`)
4. **HTTP handlers** (`chat_handlers.go`)

### Phase 4b: AI Chat Integration
5. **Chat activity** (`chat_activities.go`)
6. **Claude chat method** (`client.go`)
7. **Connect to workflow** (`chat.go`)

### Phase 4c: Mystery System
8. **Mystery types** (`mysteries.go`)
9. **Embedded mysteries** - educational and fun tracks
10. **Mystery signals/queries** (`chat.go`)
11. **Track toggle** (`cmd/serve.go`)

## Workflow Interaction Pattern

```
┌─────────────┐         ┌─────────────┐         ┌─────────────┐
│   Frontend  │         │ ChatWorkflow │         │ZiggyWorkflow│
└──────┬──────┘         └──────┬──────┘         └──────┬──────┘
       │                       │                       │
       │ POST /chat            │                       │
       │──────────────────────>│                       │
       │                       │                       │
       │                       │ Query: GetState       │
       │                       │──────────────────────>│
       │                       │<──────────────────────│
       │                       │ {personality, mood}   │
       │                       │                       │
       │                       │ Activity: GenerateChat│
       │                       │──────────────────────>│ Claude
       │                       │<──────────────────────│
       │                       │                       │
       │<──────────────────────│                       │
       │ {response, mystery}   │                       │
       └───────────────────────┴───────────────────────┘
```

## Environment Variables

| Variable | Values | Default |
|----------|--------|---------|
| `ZIGGY_TRACK` | `educational`, `fun` | `fun` |
| `ANTHROPIC_API_KEY` | API key | (required for chat) |

## Decisions

| Question | Answer |
|----------|--------|
| Message history limit | 50 messages (continue-as-new after) |
| Mystery trigger | User asks or automatic after bond > 50 |
| Hint reveal rate | One per 3 user messages minimum |
| Fallback without API | Pre-written responses, no mysteries |
