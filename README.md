# Ziggy

A virtual pet tardigrade powered by Temporal workflows.

```
    ___
   (o o)    *wiggle*
   ( : )    I've survived worse.
    \_/     But barely.
```

## Overview

Ziggy is a Tamagotchi-style virtual pet that demonstrates durable workflow patterns using [Temporal](https://temporal.io). Your tardigrade persists across server restarts, browser refreshes, and even enters a dormant "tun" state when neglected—just like real tardigrades survive extreme conditions through cryptobiosis.

## Features

- **Durable State** - Ziggy's stats persist via Temporal workflow, surviving restarts and crashes
- **Life Stages** - Evolves through Egg → Baby → Teen → Adult → Elder based on age
- **Personality System** - Develops personality based on care patterns (Stoic, Dramatic, Cheerful, Sassy, Shy)
- **AI Chat** - Converse with Ziggy using Claude AI integration
- **Mystery Games** - Fun track (riddles) and Educational track (learn Temporal concepts)
- **Dynamic Cooldowns** - Action cooldowns scale with stat urgency
- **Day/Night Cycle** - Automatic sleep based on timezone
- **Tun State** - Enter cryptobiosis when HP reaches zero, revive with care

## Quick Start

### Prerequisites

- [Go 1.21+](https://go.dev/dl/)
- [Node.js 20+](https://nodejs.org/)
- [Temporal CLI](https://docs.temporal.io/cli)
- [Task](https://taskfile.dev/) (optional)

### Setup

```bash
git clone https://github.com/rossnelson/ziggy.git
cd ziggy

# Install dependencies
task setup
# or manually:
cd web && npm install && cd ../worker && go mod tidy
```

### Run Development Servers

```bash
# Start all services
task dev

# Or individually:
task dev:temporal  # Temporal server on :7233, UI on :8233
task dev:worker    # Go worker process
task dev:api       # API server on :8080
task dev:web       # Vite dev server on :5173
```

### Optional: AI Integration

```bash
export ANTHROPIC_API_KEY=sk-ant-...
```

Without the API key, Ziggy uses embedded fallback message pools.

---

# Architecture

## System Overview

Ziggy demonstrates how to build interactive, long-running applications with Temporal. The architecture separates concerns into three layers: a reactive frontend, a stateless API server, and durable Temporal workflows.

```
┌─────────────────────────────────────────────────────────────────────────┐
│                              BROWSER                                     │
│  ┌───────────────────────────────────────────────────────────────────┐  │
│  │                     Svelte 5 Frontend                              │  │
│  │  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────────┐  │  │
│  │  │ Ziggy   │ │ Stats   │ │Controls │ │ Message │ │    Chat     │  │  │
│  │  │ Sprite  │ │  Bars   │ │ Buttons │ │ Bubble  │ │  Interface  │  │  │
│  │  └─────────┘ └─────────┘ └─────────┘ └─────────┘ └─────────────┘  │  │
│  │                              │                                     │  │
│  │                     Svelte Stores (State + Cooldowns)              │  │
│  └──────────────────────────────┼────────────────────────────────────┘  │
└─────────────────────────────────┼───────────────────────────────────────┘
                                  │ HTTP + SSE
                                  ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                         Go HTTP API Server (:8080)                       │
│  ┌────────────────┐  ┌────────────────┐  ┌────────────────────────────┐ │
│  │ Signal Routes  │  │ Query Routes   │  │ SSE Event Stream           │ │
│  │ POST /signal/* │  │ GET /state     │  │ GET /events (polls 1s)     │ │
│  └───────┬────────┘  └───────┬────────┘  └─────────────┬──────────────┘ │
└──────────┼───────────────────┼─────────────────────────┼────────────────┘
           │                   │                         │
           └───────────────────┼─────────────────────────┘
                               │ Temporal SDK
                               ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                      Temporal Server (:7233)                             │
│                   Persists all workflow state                            │
└─────────────────────────────────────────────────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                         Go Worker Process                                │
│  ┌──────────────────────────────────────────────────────────────────┐   │
│  │                         WORKFLOWS                                 │   │
│  │  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────────┐  │   │
│  │  │  ZiggyWorkflow  │  │  ChatWorkflow   │  │ NeedUpdater      │  │   │
│  │  │  Main pet state │  │  Conversations  │  │ Periodic needs   │  │   │
│  │  │  + interactions │  │  + mysteries    │  │ messages         │  │   │
│  │  └─────────────────┘  └─────────────────┘  └──────────────────┘  │   │
│  └──────────────────────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────────────────────┐   │
│  │                         ACTIVITIES                                │   │
│  │  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────────┐  │   │
│  │  │ RegeneratePool  │  │ GenerateChat    │  │ QueryZiggyState  │  │   │
│  │  │ AI msg pools    │  │ AI responses    │  │ Cross-workflow   │  │   │
│  │  └────────┬────────┘  └────────┬────────┘  └──────────────────┘  │   │
│  └───────────┼────────────────────┼─────────────────────────────────┘   │
└──────────────┼────────────────────┼─────────────────────────────────────┘
               │                    │
               └────────────────────┘
                        │
                        ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                      Claude API (Haiku model)                            │
│                   Personality dialogue + Chat responses                  │
└─────────────────────────────────────────────────────────────────────────┘
```

## Tech Stack

| Layer | Technology | Purpose |
|-------|------------|---------|
| Frontend | Svelte 5 + TypeScript + Tailwind | Reactive UI with real-time updates |
| API | Go + net/http | Stateless HTTP server bridging UI to Temporal |
| Orchestration | Temporal | Durable workflow execution and state persistence |
| AI | Claude API (Haiku) | Personality-driven dialogue and chat |

---

# Temporal Features

Ziggy showcases Temporal's core capabilities for building durable applications. Here's every Temporal feature used and why.

## Signals

**What**: Asynchronous messages sent to a running workflow.

**Used For**: All player interactions—feed, play, pet, wake, chat messages, starting mysteries.

**Why**: Signals are fire-and-forget, non-blocking, and queue automatically. When a player clicks "Feed", the API sends a signal and returns immediately. The workflow processes signals in order, ensuring no actions are lost even under load.

```go
// API sends signal
registry.SignalWorkflow(ctx, "ziggy-dev", "feed", struct{}{})

// Workflow receives via channel
feedCh := workflow.GetSignalChannel(ctx, "feed")
selector.AddReceive(feedCh, func(c workflow.ReceiveChannel, more bool) {
    c.Receive(ctx, &signal)
    handleFeed(&state, now, logger)
})
```

## Queries

**What**: Synchronous read-only requests to inspect workflow state.

**Used For**: Reading Ziggy's current stats, mood, personality; fetching chat history and mystery status.

**Why**: Queries don't modify state or appear in workflow history. Perfect for UI polling—the SSE handler queries state every second without generating history events.

```go
// Register query handler in workflow
workflow.SetQueryHandler(ctx, "state", func() (ZiggyState, error) {
    return state, nil
})

// API queries workflow
result, _ := registry.QueryWorkflow(ctx, "ziggy-dev", "state")
```

## Activities

**What**: Functions that execute outside the workflow, allowing side effects.

**Used For**: AI API calls (Claude), cross-workflow queries.

**Why**: Workflows must be deterministic—they can't call external APIs directly. Activities isolate non-deterministic operations. If an activity fails, Temporal retries it automatically.

| Activity | Purpose |
|----------|---------|
| `RegeneratePool` | Calls Claude API to generate personality-specific message pools |
| `GenerateChatResponse` | Calls Claude API to generate chat responses |
| `QueryZiggyState` | Queries ZiggyWorkflow from ChatWorkflow (workflows can't query each other directly) |

```go
// Execute activity with timeout and retry
ao := workflow.ActivityOptions{
    StartToCloseTimeout: 2 * time.Minute,
}
actCtx := workflow.WithActivityOptions(ctx, ao)
workflow.ExecuteActivity(actCtx, "RegeneratePool", input).Get(ctx, &output)
```

## Timers

**What**: Durable sleeps that survive workflow restarts.

**Used For**: 6-hour pool regeneration intervals, 30-second need checker intervals.

**Why**: `workflow.NewTimer()` is durable—if the worker crashes mid-sleep, the timer resumes where it left off. Used for periodic AI message regeneration without accumulating history.

```go
// Pool regeneration every 6 hours
timerFuture := workflow.NewTimer(ctx, 6 * time.Hour)
selector.AddFuture(timerFuture, func(f workflow.Future) {
    regeneratePool("scheduled")
    timerFuture = workflow.NewTimer(ctx, 6 * time.Hour) // Reset
})
```

## Selectors

**What**: Multiplexing construct to wait on multiple channels/futures simultaneously.

**Used For**: Main workflow loop—waits for any signal OR timer to fire.

**Why**: A single selector handles feed, play, pet, wake, need messages, AND the regeneration timer. Clean event loop without blocking.

```go
selector := workflow.NewSelector(ctx)
selector.AddReceive(feedCh, handleFeed)
selector.AddReceive(playCh, handlePlay)
selector.AddReceive(petCh, handlePet)
selector.AddReceive(wakeCh, handleWake)
selector.AddFuture(timerFuture, handleTimer)
selector.Select(ctx) // Blocks until one fires
```

## Continue-As-New

**What**: Atomically complete current workflow and start a new execution with fresh history.

**Used For**: Preventing unbounded history growth in long-running workflows.

**Why**: Temporal records every event. Ziggy runs indefinitely, accumulating signals. Without continue-as-new, history would grow forever. We trigger it at 10,000 events (ZiggyWorkflow), 50 messages (ChatWorkflow), or 100 iterations (NeedUpdater).

```go
if workflow.GetInfo(ctx).GetCurrentHistoryLength() > 10000 {
    return workflow.NewContinueAsNewError(ctx, ZiggyWorkflow, ZiggyInput{
        Owner:      input.Owner,
        Generation: state.Generation + 1,
        CreatedAt:  state.CreatedAt, // Preserve birth time
    })
}
```

## Async Goroutines (workflow.Go)

**What**: Spawn concurrent work within a workflow.

**Used For**: Non-blocking AI pool regeneration.

**Why**: Pool regeneration takes 10-30 seconds. Running it synchronously would block signal processing. `workflow.Go()` spawns it asynchronously—the main loop continues handling feed/play/pet while the pool generates in the background.

```go
workflow.Go(ctx, func(ctx workflow.Context) {
    var output PoolRegenerationOutput
    workflow.ExecuteActivity(actCtx, "RegeneratePool", input).Get(ctx, &output)
    state.RuntimePool = output.Pool // Update when ready
})
```

## Signal External Workflow

**What**: Send a signal from one workflow to another.

**Used For**: NeedUpdaterWorkflow signaling need messages to ZiggyWorkflow.

**Why**: NeedUpdater runs independently, checking Ziggy's needs every 30 seconds. When it detects hunger/boredom/loneliness, it signals Ziggy to update the displayed message. Decouples scheduling from the main workflow.

```go
// NeedUpdaterWorkflow signals ZiggyWorkflow
workflow.SignalExternalWorkflow(ctx, ziggyWorkflowID, "", "updateNeedMessage", signal)
```

## Workflow Replay (Implicit)

**What**: Temporal reconstructs workflow state by replaying history.

**Design Consideration**: All workflow code must be deterministic. Same inputs = same decisions.

**Why We Avoided**:
- No `time.Now()` in workflows—use `workflow.Now(ctx)`
- No random in workflows—use `workflow.SideEffect()` or move to activities
- No direct API calls—use activities

State decay is calculated on-demand when signals arrive, not via background timers, keeping the workflow deterministic.

---

# Design Decisions

## Why On-Demand State Decay?

**Problem**: Stats (fullness, happiness, bond) decay over time. Traditional approach: background timer ticks every N seconds.

**Issue**: Timers generate history events. A tick every 10 seconds = 8,640 events/day just for decay.

**Solution**: Calculate decay when state is accessed. When a signal arrives or query executes, compute elapsed time since last update and apply decay retroactively.

```go
func (s *ZiggyState) CalculateCurrentState(now time.Time) ZiggyState {
    elapsed := now.Sub(s.LastUpdateTime)
    ticks := int(elapsed / DecayInterval)
    // Apply ticks worth of decay
    for i := 0; i < ticks; i++ {
        s.applyDecay()
    }
    return s
}
```

**Benefit**: Zero history events for decay. Workflow stays lightweight indefinitely.

## Why Separate NeedUpdaterWorkflow?

**Problem**: We want to show "I'm hungry" messages periodically when stats are low.

**Issue**: Adding a timer to ZiggyWorkflow for this would:
1. Generate timer events in history
2. Complicate the main event loop
3. Couple scheduling with game logic

**Solution**: Spawn a separate workflow that:
1. Sleeps 30 seconds (its own history, doesn't pollute Ziggy's)
2. Queries Ziggy's state via activity
3. Signals Ziggy if a need message should display
4. Continues-as-new every 100 iterations

**Benefit**: Clean separation. Main workflow handles interactions; child handles scheduling.

## Why Client-Side Cooldown Countdown?

**Problem**: Cooldowns (30s for feed, etc.) need UI feedback.

**Issue**: Polling the API every 100ms for remaining cooldown wastes resources.

**Solution**:
1. API returns `lastFeedTime`, `lastPlayTime`, etc. in state
2. Client calculates remaining time locally
3. Re-syncs on action completion

```typescript
function getCooldownRemaining(action: string): number {
    const lastTime = cooldownTimestamps[action];
    const baseCooldown = getCooldownBase(action);
    return Math.max(0, baseCooldown - (Date.now() - lastTime));
}
```

**Benefit**: Smooth countdown animation without API spam.

## Why Three-Tier Message Pools?

**Problem**: AI-generated messages are best but API may fail/be unavailable.

**Solution**: Three fallback layers:
1. **RuntimePool**: AI-generated for current personality/stage
2. **FallbackPool**: Hardcoded per personality
3. **GenericPool**: Universal fallback (Stoic personality)

```go
selector := NewPoolSelector(runtimePool, fallbackPool, genericPool)
message := selector.Pick("feedSuccess")
// Tries each layer until finding non-empty slice
```

**Benefit**: Graceful degradation. No crashes if AI unavailable.

## Why SSE Instead of WebSockets?

**Problem**: UI needs real-time state updates.

**Options**: WebSockets (bidirectional) or SSE (server-push only).

**Decision**: SSE because:
1. Simpler—HTTP-based, works through proxies
2. Auto-reconnect built into browser EventSource API
3. All client-to-server communication already uses POST endpoints
4. Sufficient for our use case (server pushes state changes)

```go
// Server polls and pushes changes
ticker := time.NewTicker(time.Second)
for {
    select {
    case <-ticker.C:
        state := queryWorkflow()
        if stateChanged(state, lastState) {
            sendSSE(w, state)
        }
    }
}
```

---

# Application Components

## Frontend (web/)

Built with Svelte 5, leveraging runes (`$state`, `$derived`, `$effect`) for reactivity.

| Component | Purpose |
|-----------|---------|
| `Game.svelte` | Main layout container (responsive grid) |
| `Ziggy.svelte` | Sprite renderer (mood + stage determine appearance) |
| `Stats.svelte` | Health bars (HP, Fullness, Happiness, Bond) |
| `Controls.svelte` | Action buttons with cooldown indicators |
| `Message.svelte` | Speech bubble for Ziggy's dialogue |
| `Chat.svelte` | Desktop chat interface |
| `ChatDrawer.svelte` | Mobile chat drawer |
| `Background.svelte` | Time-of-day backgrounds |
| `api.ts` | Fetch client + SSE handler |
| `store.ts` | Svelte stores for state management |

## API Server (worker/internal/api/)

Stateless HTTP server using Go's standard library.

| Route | Method | Purpose |
|-------|--------|---------|
| `/api/state` | GET | Query current Ziggy state |
| `/api/signal/{feed\|play\|pet\|wake}` | POST | Send interaction signal |
| `/api/events` | GET | SSE stream for real-time updates |
| `/api/chat/history` | GET | Get chat messages |
| `/api/chat/message` | POST | Send chat message |
| `/api/chat/mysteries` | GET | List available mysteries |
| `/api/chat/mystery/start` | POST | Start a mystery |

## Workflows (worker/internal/workflow/)

| Workflow | Purpose | Continue-as-new Trigger |
|----------|---------|------------------------|
| `ZiggyWorkflow` | Main pet state, interactions, personality | 10,000 history events |
| `ChatWorkflow` | Conversation history, mysteries, AI responses | 50 messages |
| `NeedUpdaterWorkflow` | Periodic need message updates | 100 iterations |

## Activities (worker/internal/workflow/)

| Activity | Purpose |
|----------|---------|
| `RegeneratePool` | Generate AI message pool for personality |
| `GenerateChatResponse` | Generate AI chat response |
| `QueryZiggyState` | Query Ziggy from Chat workflow |

---

# Game Mechanics

## Stats & Decay

| Stat | Awake Decay | Asleep Decay | Critical Threshold |
|------|-------------|--------------|-------------------|
| Fullness | -2/tick | -1/tick | < 20 (Hungry mood) |
| Happiness | -1/tick | +0.5/tick | < 20 (Sad mood) |
| Bond | -0.5/tick | 0/tick | < 20 (Lonely mood) |
| HP | Moves toward average | Faster recovery | 0 (Tun state), < 20 (Critical) |

*Tick interval: 10 seconds*

**Bond Protection**: High bond (> 50) reduces fullness/happiness decay rate.

## Dynamic Cooldowns

Cooldowns shorten when the relevant stat is critical:

| Stat Level | Cooldown Multiplier |
|------------|---------------------|
| < 20 | 0.25x (desperate) |
| < 40 | 0.50x |
| < 60 | 0.75x |
| ≥ 60 | 1.0x (normal) |

| Action | Base Cooldown | When Critical |
|--------|--------------|---------------|
| Feed | 30s | 7.5s |
| Play | 60s | 15s |
| Pet | 10s | 2.5s |

## Personality System

Personality affects dialogue tone and message pools.

| Personality | Trigger |
|-------------|---------|
| **Shy** | Low bond (< 30) |
| **Sassy** | Neglected (2h+ no interaction) AND low bond |
| **Dramatic** | Neglected with any bond level |
| **Cheerful** | High bond (> 70) AND well-fed (> 60) |
| **Stoic** | Default/balanced state |

## Life Stages

Based on age since creation:

| Stage | Age |
|-------|-----|
| Egg | 0-1 minute |
| Baby | 1-5 minutes |
| Teen | 5-15 minutes |
| Adult | 15 minutes - 1 hour |
| Elder | 1+ hour |

## Tun State (Cryptobiosis)

When HP reaches 0, Ziggy enters tun state (tardigrade dormancy):
- Cannot play
- Feeding gives +15 fullness, +5 HP
- Petting gives +5 bond, +2 HP
- Revives when HP ≥ 20

---

# Chat System

## Two Tracks

**Fun Track** (🎮): Riddle-based mystery games with progressive hints.

**Educational Track** (📚): Direct teaching of Temporal concepts:
- Signals & Queries
- Continue-as-New
- Activities
- Workflow Replay
- Child Workflows

## Mystery Progression

1. Select track and topic
2. AI presents the concept/riddle
3. Ask questions or request hints
4. Progress tracked (hints given / total hints)
5. Celebrate completion when solved/learned

---

# Development

## Commands

```bash
task dev        # Start all services
task build      # Build all components
task test       # Run all tests
task lint       # Run all linters
task format     # Format all code
```

## Project Structure

```
ziggy/
├── web/                    # Svelte 5 frontend
│   └── src/lib/            # Components, stores, API client
├── worker/                 # Go backend
│   ├── cmd/                # CLI commands (serve, worker)
│   └── internal/
│       ├── api/            # HTTP handlers
│       ├── ai/             # Claude API client
│       ├── temporal/       # Temporal registry
│       └── workflow/       # Workflows, activities, state
└── Taskfile.yml            # Task runner config
```

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `ANTHROPIC_API_KEY` | No | Enables AI-generated dialogue |
| `TEMPORAL_ADDRESS` | No | Temporal server (default: localhost:7233) |
| `TEMPORAL_NAMESPACE` | No | Namespace (default: default) |

---

# License

MIT
