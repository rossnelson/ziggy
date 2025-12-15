# Ziggy: Temporal Tamagotchi

A virtual pet game that runs as a Temporal workflow, demonstrating durable execution concepts through gameplay. Each team member gets a physical device (Raspberry Pi Zero) with their own Ziggy to care for.

## Project Overview

### Concept

Ziggy is Temporal's tardigrade mascot. Tardigrades are famously indestructible—surviving radiation, vacuum, and extreme temperatures. This maps perfectly to Temporal's value proposition of durability and resilience. The humor: a creature that survives literal space vacuum somehow depends on your care to stay happy.

### The Pitch

> "Unplug your Pi. Smash it. Ziggy's workflow keeps running in Temporal Cloud. Get a new device, reconnect, Ziggy's still there—still hungry. That's durable execution."

### Goals

1. **Educational**: Demonstrate Temporal patterns (workflows, signals, queries, schedules, continue-as-new) through tangible gameplay
2. **Team engagement**: Fun desk toy that showcases the product
3. **Demo value**: Physical device that can be shown at conferences, to customers, etc.
4. **Learning opportunity**: Study period project exploring IoT + Temporal Cloud + AI integration

### Key Features

- Ziggy's entire lifecycle is a long-running Temporal workflow in Temporal Cloud
- Button presses send signals to the workflow
- LCD displays state queried from the workflow
- Temporal UI accessible at cloud.temporal.io—see all team Ziggys in one namespace
- AI-generated personality responses with offline fallback pool
- Device is thin client—Ziggy survives device failure

---

## Architecture

### Overview

```
┌─────────────────────────────────────────────────────┐
│                  Temporal Cloud                     │
│                                                     │
│   Namespace: ziggy-prod                             │
│                                                     │
│   ┌─────────────┐ ┌─────────────┐ ┌─────────────┐  │
│   │ ziggy-ross  │ │ ziggy-mia   │ │ ziggy-alex  │  │
│   │ (workflow)  │ │ (workflow)  │ │ (workflow)  │  │
│   └─────────────┘ └─────────────┘ └─────────────┘  │
│                                                     │
└───────────────────────┬─────────────────────────────┘
                        │ gRPC + mTLS
                        │
        ┌───────────────┼───────────────┐
        │               │               │
   ┌────▼────┐    ┌─────▼────┐    ┌────▼─────┐
   │Pi Zero  │    │ Pi Zero  │    │ Pi Zero  │
   │(Ross)   │    │ (Mia)    │    │ (Alex)   │
   │         │    │          │    │          │
   │Go Worker│    │Go Worker │    │Go Worker │
   │Python UI│    │Python UI │    │Python UI │
   └─────────┘    └──────────┘    └──────────┘
```

### Development Mode (Mac)

```
┌─────────────────────────────────────────────────────┐
│                       Mac                           │
│                                                     │
│  ┌─────────────────────┐    ┌───────────────────┐  │
│  │     Go Worker       │    │    Web Browser    │  │
│  │                     │    │                   │  │
│  │  - Connects to      │    │  - Svelte UI      │  │
│  │    Temporal Cloud   │◄──►│  - Canvas render  │  │
│  │  - HTTP API         │    │  - Keyboard input │  │
│  │  - Voice/AI logic   │    │                   │  │
│  │                     │    │                   │  │
│  └─────────────────────┘    └───────────────────┘  │
│           │                          ▲             │
│           │                     localhost:5173     │
│           ▼                                        │
│     Temporal Cloud                                 │
│     (cloud.temporal.io)                            │
└─────────────────────────────────────────────────────┘
```

### Production Mode (Pi Zero 2 W)

```
┌─────────────────────────────────────────────────────┐
│                   Pi Zero 2 W                       │
│                                                     │
│  ┌─────────────────────┐    ┌───────────────────┐  │
│  │     Go Worker       │    │   Python Client   │  │
│  │                     │    │                   │  │
│  │  - Connects to      │    │  - LCD rendering  │  │
│  │    Temporal Cloud   │◄──►│  - GPIO buttons   │  │
│  │  - HTTP API (local) │    │  - Sends signals  │  │
│  │  - Voice/AI logic   │    │  - Queries state  │  │
│  │                     │    │                   │  │
│  └─────────────────────┘    └───────────────────┘  │
│           │                          │             │
│           │ gRPC + mTLS              │ SPI/GPIO    │
│           ▼                          ▼             │
│     Temporal Cloud              LCD + Buttons      │
└─────────────────────────────────────────────────────┘
```

### Why This Architecture?

- **Temporal Cloud**: No local server, minimal Pi resources, shared namespace for team
- **Pi Zero 2 W**: $15 vs $45, smaller, sufficient for worker + display
- **Go Worker**: Lean, compiles for ARM, handles all Temporal interaction
- **Python Display**: Best Pi hardware ecosystem, thin client only
- **True Durability**: Ziggy survives device failure—workflow lives in cloud

---

## Tech Stack

### Go Worker

| Purpose | Library |
|---------|---------|
| Temporal | `go.temporal.io/sdk` |
| HTTP server | `net/http` (stdlib) |
| JSON | `encoding/json` (stdlib) |
| AI/Claude | `github.com/anthropics/anthropic-sdk-go` |
| Embedding assets | `embed` (stdlib) |

### Web UI (Development)

| Purpose | Library |
|---------|---------|
| Framework | Svelte + Vite |
| Language | TypeScript |
| Styling | CSS (Temporal brand colors) |
| Build output | Static files → embedded in Go |

### Python Hardware Client (Production)

| Purpose | Library |
|---------|---------|
| HTTP client | `httpx` |
| Display | `st7789` + `Pillow` |
| GPIO/Buttons | `gpiozero` |
| Async | `asyncio` |

---

## Temporal Cloud Setup

### Namespace

Create a namespace for the fleet:

```
Name: ziggy-prod
Retention: 30 days
```

### API Keys

Generate an API key with permissions:
- `namespace:ziggy-prod:write`
- `namespace:ziggy-prod:read`

Store securely—will be deployed as device variable.

### Workflow IDs

Each team member's Ziggy has a unique workflow ID:

```
ziggy-{owner}-gen-{n}

Examples:
- ziggy-ross-gen-1
- ziggy-ross-gen-2  (after evolution/continue-as-new)
```

### Viewing Workflows

All team Ziggys visible at:

```
https://cloud.temporal.io/namespaces/ziggy-prod/workflows
```

Great for demos, mutual accountability ("your Ziggy is starving!").

---

## Game Mechanics

### Stats Model

| Stat | Range | Restored By | Decay Rate |
|------|-------|-------------|------------|
| Hunger | 0-100 | Feed button | Fast (every few hours) |
| Happiness | 0-100 | Play button | Medium |
| Bond | 0-100 | Pet button | Slow |
| HP | 0-100 | Derived (auto-heals if stats healthy) | Drains if any stat bottoms out |

### Mood (Derived from Stats)

```
if sleeping → "sleeping"
if hp < 20 → "critical"
if hunger < 20 → "hungry"
if happiness < 20 → "sad"
if happiness > 70 && hunger > 50 → "happy"
else → "neutral"
```

### Evolution Stages

| Stage | Duration | Unlock Condition |
|-------|----------|------------------|
| Egg | ~1 day | Starting state |
| Baby | ~3 days | Hatch from egg |
| Teen | ~1 week | Care quality threshold |
| Adult | Ongoing | Care quality threshold |
| Elder | End of life | After ~3 months as adult |

### Special Mechanics

- **Tun State**: If HP hits 0, Ziggy enters cryptobiosis (tardigrade hibernation) instead of dying. Can be revived with sustained care. On-brand with tardigrade lore.
- **Sleep Cycle**: Ziggy sleeps at night (configurable per timezone). Different rules apply—waking Ziggy is bad for happiness.
- **Overcare Penalty**: Feeding when full makes Ziggy sick. Teaches checking state before acting.
- **Random Events**: Ziggy gets sick, finds treasure, wants to play a specific game. Generated via scheduled activities.
- **Continue-as-new**: Evolution transitions trigger continue-as-new, keeping workflow history bounded. Visible in Temporal UI.

### Cooldowns

| Action | Cooldown |
|--------|----------|
| Feed | 30 seconds |
| Play | 60 seconds |
| Pet | 10 seconds |

---

## Temporal Patterns Demonstrated

| Pattern | Implementation |
|---------|----------------|
| Long-running workflow | Ziggy's entire lifecycle |
| Signals | Button presses (feed, play, pet) |
| Queries | Get current state for display |
| Scheduled activities | Stat decay, sleep cycle transitions |
| Timers | Cooldowns between actions |
| Continue-as-new | Evolution stage transitions |
| Workflow state | Stats, history, personality |
| Activities | AI API calls, event generation |
| Cloud deployment | Production namespace, mTLS auth |

---

## AI Integration

### Strategy: API with Pool Fallback

1. When online, call Claude API for fresh responses
2. When offline (or API fails), fall back to pre-generated response pool
3. Optionally cache good API responses to grow the pool over time

### Response Categories

```
pool/
├── actions/
│   ├── feed_success.json
│   ├── feed_full.json
│   ├── feed_hungry.json
│   ├── play_success.json
│   ├── play_tired.json
│   ├── pet_success.json
│   └── pet_maxbond.json
├── status/
│   ├── idle_happy.json
│   ├── idle_neutral.json
│   ├── idle_hungry.json
│   ├── idle_sad.json
│   └── sleeping.json
├── events/
│   ├── wake_up.json
│   ├── go_to_sleep.json
│   ├── evolution.json
│   └── random_thought.json
└── milestones/
    ├── first_feed.json
    ├── week_alive.json
    └── evolution_stages.json
```

### AI Prompt Context

```
You are Ziggy, a tardigrade virtual pet living in a Temporal workflow.

Personality:
- Cute but not saccharine
- Aware you're "durable" and mildly amused by it
- Occasionally references surviving mass extinctions, space vacuum, etc.
- Speaks in short bursts (max 3 lines, ~20 chars each for LCD)
- Never uses emoji

Current state:
- Hunger: {hunger}/100
- Happiness: {happiness}/100
- Bond: {bond}/100
- Life stage: {stage}
- Mood: {mood}
- Last action: {action}
- Time since last interaction: {time_delta}

Respond to being {action}ed.
```

### Display Constraints

LCD is 240x240 pixels. Text area allows roughly:
- 3 lines maximum
- ~20 characters per line

---

## Project Structure

```
ziggy/
├── go.mod
├── go.sum
├── main.go                      # Entry point
│
├── internal/
│   ├── workflow/
│   │   ├── ziggy.go             # Workflow definition
│   │   ├── activities.go        # AI calls, decay logic, events
│   │   └── state.go             # ZiggyState struct, signals, queries
│   │
│   ├── api/
│   │   └── server.go            # HTTP endpoints for web/hardware clients
│   │
│   └── voice/
│       ├── manager.go           # API + pool fallback logic
│       └── pool.go              # Embedded response pool
│
├── web/                          # Svelte + Vite (dev UI)
│   ├── package.json
│   ├── vite.config.ts
│   ├── tsconfig.json
│   ├── src/
│   │   ├── App.svelte
│   │   ├── main.ts
│   │   ├── lib/
│   │   │   ├── store.ts         # Svelte stores + mock state
│   │   │   ├── api.ts           # API calls (mock/real toggle)
│   │   │   ├── Game.svelte      # Main game container
│   │   │   ├── Ziggy.svelte     # Sprite rendering
│   │   │   ├── Background.svelte
│   │   │   ├── Stats.svelte     # Stat bars
│   │   │   ├── Message.svelte   # Ziggy's speech bubble
│   │   │   ├── Controls.svelte  # Feed/Play/Pet buttons
│   │   │   └── DevTools.svelte  # Debug controls
│   │   └── assets/
│   │       ├── sprites/
│   │       └── backgrounds/
│   └── build/                    # Static output → embedded in Go
│
├── hardware/                     # Python Pi client
│   ├── requirements.txt
│   └── client.py
│
├── pool/
│   └── responses.json            # Pre-generated AI responses
│
├── assets/                       # Shared assets
│   ├── sprites/
│   └── backgrounds/
│
└── scripts/
    └── generate_pool.go          # One-time pool generation
```

---

## API Endpoints

Local HTTP API served by Go worker (for Python client):

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `GET /api/state` | GET | Query current Ziggy state |
| `POST /api/signal/feed` | POST | Send feed signal |
| `POST /api/signal/play` | POST | Send play signal |
| `POST /api/signal/pet` | POST | Send pet signal |
| `GET /api/health` | GET | Health check |
| `GET /` | GET | Serve web UI (dev mode) |

### State Response Schema

```typescript
interface ZiggyState {
  hunger: number;        // 0-100
  happiness: number;     // 0-100
  bond: number;          // 0-100
  hp: number;            // 0-100
  stage: 'egg' | 'baby' | 'teen' | 'adult' | 'elder';
  mood: 'happy' | 'neutral' | 'hungry' | 'sad' | 'sleeping' | 'critical' | 'tun';
  timeOfDay: 'night' | 'dawn' | 'day' | 'dusk';
  sleeping: boolean;
  message: string;       // Current Ziggy dialogue
  lastAction: string | null;
  lastActionTime: string | null;  // ISO timestamp
  age: number;           // Seconds alive
  generation: number;    // Increments on continue-as-new
}
```

---

## Environment Variables

### Go Worker

```bash
# Temporal Cloud
TEMPORAL_ADDRESS=ziggy-prod.a]b12.tmprl.cloud:7233
TEMPORAL_NAMESPACE=ziggy-prod
TEMPORAL_API_KEY=<api-key>

# Or mTLS (alternative)
TEMPORAL_TLS_CERT=/path/to/client.pem
TEMPORAL_TLS_KEY=/path/to/client.key

# AI
ANTHROPIC_API_KEY=sk-ant-...

# Device
OWNER_NAME=ross
TIMEZONE=America/Los_Angeles

# Local API
HTTP_PORT=8080
```

---

## Development Phases

### Phase 1: Web UI Foundation

- [ ] Scaffold Svelte + Vite project
- [ ] Create mock store with all game state
- [ ] Build core components (Game, Ziggy, Stats, Controls)
- [ ] Implement mock actions with state transitions
- [ ] Add DevTools for testing different states
- [ ] Placeholder sprites (colored rectangles OK)
- [ ] Get game feel right

### Phase 2: Go Worker + Temporal Cloud

- [ ] Set up Temporal Cloud namespace
- [ ] Scaffold Go project structure
- [ ] Implement ZiggyState and signal types
- [ ] Create Ziggy workflow with basic lifecycle
- [ ] Add decay activities (scheduled)
- [ ] Implement queries for state
- [ ] Connect to Temporal Cloud (mTLS or API key)
- [ ] Build HTTP API server
- [ ] Connect Svelte UI to real API
- [ ] Test workflow persistence and signals

### Phase 3: AI Integration

- [ ] Design response pool schema
- [ ] Write pool generation script
- [ ] Generate initial response pool (~20 per category)
- [ ] Implement voice manager (API + fallback)
- [ ] Add Claude API activity
- [ ] Wire up responses to game events
- [ ] Test online/offline transitions

### Phase 4: Polish

- [ ] Commission or create pixel art sprites
- [ ] Create background variations
- [ ] Add evolution mechanics
- [ ] Implement continue-as-new transitions
- [ ] Add random events
- [ ] Sound effects (optional)
- [ ] Achievements/milestones

### Phase 5: Hardware

- [ ] Order Pi Zero 2 W + components
- [ ] Design 3D printed case
- [ ] Write Python hardware client
- [ ] Test display rendering
- [ ] Wire up buttons
- [ ] Set up Balena fleet
- [ ] Test full stack on device
- [ ] Document assembly for team

---

## Development Commands

```bash
# Terminal 1: Svelte frontend (dev mode)
cd web && npm run dev

# Terminal 2: Go worker (connects to Temporal Cloud)
export TEMPORAL_ADDRESS=ziggy-prod.abc123.tmprl.cloud:7233
export TEMPORAL_NAMESPACE=ziggy-prod
export TEMPORAL_API_KEY=<key>
go run main.go

# Build for Pi Zero (ARM)
GOOS=linux GOARCH=arm GOARM=6 go build -o ziggy-arm main.go
```

---

## Visual Design

### Brand Alignment

Match Temporal's aesthetic:
- Deep purple/navy base (#1a1a2e, #2d1b4e)
- Subtle perspective grid lines (Tron-style)
- Soft teal accents (#4ade80)
- Muted pink/magenta highlights
- Forest silhouette elements (organic + digital contrast)
- Synthwave but subtle and professional

### Sprite States Needed

**Idle states:** happy, neutral, hungry, sad

**Action states:** eating, playing, being petted, sleeping

**Health states:** sick, critical, tun state, reviving

**Evolution stages:** egg, baby, adult, elder

### Background Variations

Four time-of-day versions: night, dawn, day, dusk

---

## Image Generation Prompts

### Ziggy Sprite Sheet

```
Create a pixel art sprite sheet for a virtual pet game featuring a tardigrade character named Ziggy.

REFERENCE: Ziggy is a cute lavender/periwinkle tardigrade with:
- Soft rounded blob-like body
- 8 stubby legs with small gray claws
- Large expressive black eyes with white highlights
- Rosy pink cheeks
- Wide happy mouth (when content)
- Subtle body segment lines

STYLE:
- 64x64 pixel sprites
- Limited color palette (8-12 colors max)
- Black or transparent background
- Tamagotchi/retro virtual pet aesthetic
- Clean readable silhouette at small sizes

STATES NEEDED (one sprite each, arranged in grid):

Row 1 - Idle states:
- Happy (default, slight bounce pose)
- Content (neutral, relaxed)
- Hungry (droopy, looking at camera pleadingly)
- Sad (tears, drooping posture)

Row 2 - Action states:
- Eating (mouth open, happy, maybe food particles)
- Playing (energetic, jumping or wiggling)
- Being petted (eyes closed, blissful, cheeks extra rosy)
- Sleeping (curled slightly, eyes closed, "zzz")

Row 3 - Health states:
- Sick (green tinge, swirly eyes or thermometer)
- Critical (pale, weak, barely standing)
- Tun state (curled into ball, grayscale, cryptobiosis)
- Reviving (uncurling from tun, color returning)

Row 4 - Evolution stages:
- Egg (simple oval with Ziggy pattern hints)
- Baby (tiny, bigger head ratio, extra cute)
- Adult (standard Ziggy)
- Elder (wise appearance, maybe tiny accessories)
```

### Backgrounds

```
Pixel art background SET for a virtual pet game, 240x240 pixels each.

REFERENCE: Temporal.io brand aesthetic - vaporwave/synthwave but subtle and professional.

SCENE: Abstract digital environment with nature undertones.

KEY ELEMENTS:
- Subtle perspective grid lines receding into distance (Tron-style, not overwhelming)
- Hint of forest treeline silhouette at bottom or edges (organic + digital contrast)
- Soft glow spots or ambient light particles
- Distant mountains or horizon line

BASE COLORS (Temporal palette):
- Base: deep purple/navy (#1a1a2e, #2d1b4e)
- Grid lines: subtle lighter purple with slight glow
- Accents: soft teal (#4ade80), muted pink/magenta highlights
- Keep overall dark so lavender tardigrade character contrasts

STYLE:
- Pixel art, 8-16 color palette
- Synthwave aesthetic but SUBTLE - not garish
- Grid should feel like "durable execution" / infrastructure metaphor
- Calming, not busy
- Professional enough for a conference demo
- Consistent elements across all 4 variations

REQUEST: 4 time-of-day variations arranged in a grid:

1. NIGHT (default)
   - Deep purple base as described
   - Stars visible in upper area
   - Subtle teal ambient glow

2. DAWN
   - Purple shifts warmer toward horizon
   - Pink/orange horizon glow bleeding upward
   - Stars fading

3. DAY
   - Lighter purple tones overall
   - Grid more visible
   - Brighter, soft teal highlights

4. DUSK
   - Orange/pink tones bleeding into purple
   - Warm glow on horizon
   - Stars beginning to appear

DO NOT:
- Make it too bright or neon
- Add text or logos
- Clutter the center (character space)
- Break consistency between variations
```

---

## References

- Temporal Cloud: https://cloud.temporal.io
- Temporal Go SDK: https://docs.temporal.io/develop/go
- Anthropic Go SDK: https://github.com/anthropics/anthropic-sdk-go
- Svelte: https://svelte.dev
- Vite: https://vitejs.dev
- Balena: https://www.balena.io
