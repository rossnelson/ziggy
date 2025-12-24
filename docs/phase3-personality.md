# Phase 3: AI Integration with Personality System

## Overview

Add personality-based message pools to Ziggy with:
- Embedded fallback pools per personality (in binary)
- Runtime pools stored in workflow state (regenerated via Claude API)
- Graceful fallback hierarchy: runtime → fallback → generic

## Personality Types

| Type | Traits | Voice Example |
|------|--------|---------------|
| Stoic | Deadpan, philosophical, references mass extinctions | "I've survived worse. The Permian extinction, for instance." |
| Dramatic | Theatrical, everything is life or death | "FINALLY! I was PERISHING! You've saved me... this time." |
| Cheerful | Optimistic, encouraging, wholesome | "Yay, food! You're the best friend a tardigrade could have!" |
| Sassy | Sarcastic, playful roasts, demanding | "Oh, you remembered I exist? How generous of you." |
| Shy | Quiet, tentative, warms up with high bond | "...thank you... *small wiggle*" |

## Architecture

```
┌─────────────────────────────────────────────────────┐
│                   Message Selection                  │
├─────────────────────────────────────────────────────┤
│  1. Check runtime pool (workflow state)             │
│     └─ If empty or missing category → fallback      │
│  2. Check fallback pool (embedded, personality)     │
│     └─ If missing → generic fallback                │
│  3. Generic fallback (current messages.go)          │
└─────────────────────────────────────────────────────┘
```

## Data Structures

### Personality Type
```go
type Personality string

const (
    PersonalityStoic    Personality = "stoic"
    PersonalityDramatic Personality = "dramatic"
    PersonalityCheerful Personality = "cheerful"
    PersonalitySassy    Personality = "sassy"
    PersonalityShy      Personality = "shy"
)
```

### Message Pool Structure
```go
type MessagePool struct {
    // Action responses
    FeedSuccess  []string `json:"feedSuccess"`
    FeedFull     []string `json:"feedFull"`
    FeedHungry   []string `json:"feedHungry"`
    FeedSleeping []string `json:"feedSleeping"`
    FeedTun      []string `json:"feedTun"`

    PlaySuccess  []string `json:"playSuccess"`
    PlayTired    []string `json:"playTired"`
    PlayHappy    []string `json:"playHappy"`
    PlaySleeping []string `json:"playSleeping"`
    PlayTun      []string `json:"playTun"`

    PetSuccess   []string `json:"petSuccess"`
    PetMaxBond   []string `json:"petMaxBond"`
    PetLowMood   []string `json:"petLowMood"`
    PetSleeping  []string `json:"petSleeping"`
    PetTun       []string `json:"petTun"`

    Reviving     []string `json:"reviving"`

    // Idle by mood
    IdleHappy    []string `json:"idleHappy"`
    IdleNeutral  []string `json:"idleNeutral"`
    IdleHungry   []string `json:"idleHungry"`
    IdleSad      []string `json:"idleSad"`
    IdleLonely   []string `json:"idleLonely"`
    IdleCritical []string `json:"idleCritical"`
    IdleTun      []string `json:"idleTun"`
    IdleSleeping []string `json:"idleSleeping"`
}
```

### Updated ZiggyState
```go
type ZiggyState struct {
    // Existing fields...
    Fullness       float64   `json:"fullness"`
    Happiness      float64   `json:"happiness"`
    Bond           float64   `json:"bond"`
    HP             float64   `json:"hp"`
    LastUpdateTime time.Time `json:"lastUpdateTime"`
    CreatedAt      time.Time `json:"createdAt"`
    Sleeping       bool      `json:"sleeping"`
    Stage          Stage     `json:"stage"`
    Message        string    `json:"message"`
    LastAction     Action    `json:"lastAction,omitempty"`
    Timezone       string    `json:"timezone"`
    Generation     int       `json:"generation"`

    // New personality fields
    Personality     Personality  `json:"personality"`
    CareMetrics     CareMetrics  `json:"careMetrics"`
    RuntimePool     *MessagePool `json:"runtimePool,omitempty"`
    PoolGeneratedAt time.Time    `json:"poolGeneratedAt,omitempty"`
}
```

## Files to Create/Modify

### New Files
| File | Purpose |
|------|---------|
| `internal/workflow/personality.go` | Personality types and derivation logic |
| `internal/workflow/pool.go` | MessagePool struct, selection logic, fallback hierarchy |
| `internal/workflow/pools_fallback.go` | Embedded fallback pools per personality |
| `internal/workflow/activities.go` | Claude API activity for pool regeneration |
| `internal/ai/client.go` | Claude API client wrapper |

### Modified Files
| File | Changes |
|------|---------|
| `internal/workflow/state.go` | Add Personality, CareMetrics, RuntimePool fields |
| `internal/workflow/ziggy.go` | Use pool selection in handlers, track care metrics |
| `internal/workflow/register.go` | Register pool regeneration activity |
| `cmd/serve.go` | Add ANTHROPIC_API_KEY config |

## Personality Evolution Model

Personality evolves based on care patterns tracked over time:

```
┌─────────────────────────────────────────────────────┐
│              Care Pattern Tracking                   │
├─────────────────────────────────────────────────────┤
│  InteractionRate = interactions / time_alive        │
│  AverageFullness = rolling average of fullness      │
│  AverageBond     = rolling average of bond          │
│  NeglectScore    = time_since_last_action (decay)   │
└─────────────────────────────────────────────────────┘
                        │
                        ▼
┌─────────────────────────────────────────────────────┐
│            Personality Derivation                    │
├─────────────────────────────────────────────────────┤
│  High neglect + low bond     → Sassy                │
│  High neglect + any bond     → Dramatic             │
│  High care + high bond       → Cheerful             │
│  Low bond (new or distant)   → Shy                  │
│  Balanced/moderate care      → Stoic                │
└─────────────────────────────────────────────────────┘
```

### Care Metrics

```go
type CareMetrics struct {
    TotalInteractions int       `json:"totalInteractions"`
    LastInteractionAt time.Time `json:"lastInteractionAt"`
    // Rolling averages (updated on each action)
    AvgFullness       float64   `json:"avgFullness"`
    AvgBond           float64   `json:"avgBond"`
}
```

### Derivation Logic

```go
func DerivePersonality(metrics CareMetrics, state ZiggyState, now time.Time) Personality {
    timeSinceInteraction := now.Sub(metrics.LastInteractionAt)
    neglected := timeSinceInteraction > 2*time.Hour || metrics.TotalInteractions < 10

    if neglected && state.Bond < 40 {
        return PersonalitySassy
    }
    if neglected {
        return PersonalityDramatic
    }
    if state.Bond > 70 && metrics.AvgFullness > 60 {
        return PersonalityCheerful
    }
    if state.Bond < 30 {
        return PersonalityShy
    }
    return PersonalityStoic
}
```

## Pool Regeneration

### Activity Signature
```go
type PoolRegenerationInput struct {
    Personality Personality
    Stage       Stage
    Bond        float64
}

type PoolRegenerationOutput struct {
    Pool      MessagePool
    Generated time.Time
}

func RegeneratePoolActivity(ctx context.Context, input PoolRegenerationInput) (*PoolRegenerationOutput, error)
```

### Triggers
1. **On workflow start**: Generate initial pool
2. **Every 6 hours**: Scheduled timer activity
3. **On personality change**: When derived personality differs from current

### Claude Prompt Template
```
You are generating dialogue for Ziggy, a tardigrade virtual pet.

Personality: {personality}
Life stage: {stage}
Bond level: {bond_description}

Generate 10 short messages (max 3 lines, ~20 chars each) for each category:
- feedSuccess, feedFull, feedHungry, feedSleeping, feedTun
- playSuccess, playTired, playHappy, playSleeping, playTun
- petSuccess, petMaxBond, petLowMood, petSleeping, petTun
- reviving
- idleHappy, idleNeutral, idleHungry, idleSad, idleLonely, idleCritical, idleTun, idleSleeping

Rules:
- Never use emoji
- Match the {personality} voice consistently
- Reference tardigrade facts occasionally (survive space, radiation, etc.)
- Keep messages appropriate for the context

Return as JSON matching the MessagePool schema.
```

## Implementation Order

### Phase 3a: Personality System (no API yet)
- [x] Personality types (`personality.go`)
- [x] Care metrics in state (`state.go`)
- [x] Track care metrics on actions (`ziggy.go`)
- [x] Pool structure (`pool.go`)
- [x] Fallback pools per personality (`pools_fallback.go`)
- [x] Pool selection with fallback (`pool.go`)
- [x] Wire up pool selection (`ziggy.go`)

### Phase 3b: AI Integration
- [x] Claude API client (`internal/ai/client.go`)
- [x] Regeneration activity (`activities.go`)
- [x] Register activity (`register.go`)
- [x] Regeneration triggers (`ziggy.go`)
- [x] Environment config (via `os.Getenv("ANTHROPIC_API_KEY")` in ai/client.go)

## Decisions

| Question | Answer |
|----------|--------|
| Regeneration trigger | On start + every 6 hours + on personality change |
| Pool size | 10 messages per category |
| Personality selection | Evolves from care patterns |
