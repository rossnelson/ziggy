# Ziggy

A virtual pet tardigrade powered by Temporal workflows.

```
    ___
   (o o)    *wiggle*
   ( : )    I've survived worse.
    \_/     But barely.
```

## Overview

Ziggy is a Tamagotchi-style virtual pet that demonstrates durable workflow patterns using [Temporal](https://temporal.io). Your tardigrade persists across server restarts, browser refreshes, and even enters a dormant "tun" state when neglected (just like real tardigrades).

## Features

- **Durable State** - Ziggy's stats persist via Temporal workflow, surviving restarts
- **Life Stages** - Egg, Baby, Teen, Adult, Elder (based on age)
- **Personality System** - Evolves based on care patterns (Stoic, Dramatic, Cheerful, Sassy, Shy)
- **AI-Generated Dialogue** - Optional Claude API integration for personality-specific messages
- **Dynamic Cooldowns** - Action cooldowns scale with stat levels (faster recovery when critical)
- **Day/Night Cycle** - Automatic sleep based on timezone
- **Tun State** - Enter cryptobiosis when HP reaches zero, revive with care

## Tech Stack

| Component | Technology |
|-----------|------------|
| Workflow Engine | [Temporal](https://temporal.io) |
| Worker/API | Go |
| Frontend | Svelte 5 + TypeScript |
| AI (optional) | Claude API (Haiku) |

## Quick Start

### Prerequisites

- [Go 1.21+](https://go.dev/dl/)
- [Node.js 20+](https://nodejs.org/)
- [Temporal CLI](https://docs.temporal.io/cli)
- [Task](https://taskfile.dev/) (optional, for task runner)

### Setup

```bash
# Clone the repo
git clone https://github.com/rossnelson/ziggy.git
cd ziggy

# Install dependencies
task setup
# or manually:
cd web && npm install && cd ../worker && go mod tidy
```

### Run Development Servers

```bash
# Start all services (Temporal, Worker, API, Web)
task dev

# Or start individually:
task dev:temporal  # Temporal server on :7233, UI on :8233
task dev:worker    # Go worker
task dev:api       # API server on :8080
task dev:web       # Vite dev server on :5173
```

### Optional: AI-Generated Messages

Set the `ANTHROPIC_API_KEY` environment variable to enable AI-generated personality dialogue:

```bash
export ANTHROPIC_API_KEY=sk-ant-...
```

Without the API key, Ziggy uses embedded fallback message pools.

## Architecture

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│   Svelte UI     │────▶│   Go API        │────▶│   Temporal      │
│   :5173         │     │   :8080         │     │   :7233         │
└─────────────────┘     └─────────────────┘     └─────────────────┘
                                                        │
                                                        ▼
                                               ┌─────────────────┐
                                               │  ZiggyWorkflow  │
                                               │  - Stats decay  │
                                               │  - Personality  │
                                               │  - AI messages  │
                                               └─────────────────┘
```

### Workflow Signals & Queries

| Signal | Action |
|--------|--------|
| `feed` | Increase fullness |
| `play` | Increase happiness + bond |
| `pet` | Increase bond |
| `wake` | Wake from sleep (costs happiness) |

| Query | Returns |
|-------|---------|
| `state` | Current stats, mood, personality, cooldowns |

## Stats & Mechanics

| Stat | Decay Rate | Effect |
|------|------------|--------|
| Fullness | 2/tick (awake), 1/tick (asleep) | Hunger mood when < 20 |
| Happiness | 1/tick (awake), +0.5/tick (asleep) | Sad mood when < 20 |
| Bond | 0.5/tick (awake only) | Lonely mood when < 20 |
| HP | Moves toward avg of other stats | Tun state at 0, Critical < 20 |

### Cooldowns

| Action | Base | When Critical (stat < 20) |
|--------|------|---------------------------|
| Feed | 30s | 7.5s |
| Play | 60s | 15s |
| Pet | 10s | 2.5s |

### Personality Evolution

| Personality | Trigger |
|-------------|---------|
| Shy | Low bond (< 30) |
| Sassy | Neglected + low bond |
| Dramatic | Neglected + any bond |
| Cheerful | High bond + well-fed |
| Stoic | Default/balanced |

## Development Commands

```bash
task dev        # Start all services
task build      # Build all components
task test       # Run all tests
task lint       # Run all linters
task format     # Format all code
task clean      # Clean build artifacts
```

## Project Structure

```
ziggy/
├── web/                 # Svelte frontend
│   └── src/lib/         # Components & store
├── worker/              # Go backend
│   ├── cmd/             # Entry point
│   └── internal/
│       ├── api/         # HTTP handlers
│       ├── ai/          # Claude API client
│       ├── temporal/    # Temporal registry
│       └── workflow/    # Workflow & activities
├── assets/              # Sprites & backgrounds
├── docs/                # Planning docs
└── Taskfile.yml         # Task runner config
```

## License

MIT
