# Ziggy: Temporal Tamagotchi

A virtual pet game that runs as a Temporal workflow, demonstrating durable execution concepts through gameplay. Each team member gets a physical device (Raspberry Pi) with their own Ziggy to care for.

## Project Overview

### Concept

Ziggy is Temporal's tardigrade mascot. Tardigrades are famously indestructible—surviving radiation, vacuum, and extreme temperatures. This maps perfectly to Temporal's value proposition of durability and resilience. The humor: a creature that survives literal space vacuum somehow depends on your care to stay happy.

### Goals

1. **Educational**: Demonstrate Temporal patterns (workflows, signals, queries, schedules, continue-as-new) through tangible gameplay
2. **Team engagement**: Fun desk toy that showcases the product
3. **Demo value**: Physical device that can be shown at conferences, to customers, etc.
4. **Learning opportunity**: Study period project exploring IoT + Temporal + AI integration

### Key Features

- Ziggy's entire lifecycle is a long-running Temporal workflow
- Button presses send signals to the workflow
- LCD displays state queried from the workflow
- Temporal UI accessible via browser to see workflow internals
- AI-generated personality responses with offline fallback pool

---

## Architecture

### Development Mode (Mac)

```
┌─────────────────────────────────────────────────────┐
│                     Mac                             │
│                                                     │
│  ┌─────────────────────┐    ┌───────────────────┐  │
│  │     Go Binary       │    │    Web Browser    │  │
│  │                     │    │                   │  │
│  │  - Temporal Worker  │    │  - Svelte UI      │  │
│  │  - Ziggy Workflow   │◄──►│  - Canvas render  │  │
│  │  - HTTP API         │    │  - Keyboard input │  │
│  │  - Voice/AI logic   │    │                   │  │
│  │                     │    │                   │  │
│  └─────────────────────┘    └───────────────────┘  │
│           │                          ▲             │
│           ▼                     localhost:5173     │
│     Temporal Dev Server              │             │
│     (gRPC :7233, UI :8233)           │             │
│                                      │             │
│     Go serves: localhost:8080 ───────┘             │
└─────────────────────────────────────────────────────┘
```

### Production Mode (Raspberry Pi)

```
┌─────────────────────────────────────────────────────┐
│                      Pi 4                           │
│                                                     │
│  ┌─────────────────────┐    ┌───────────────────┐  │
│  │     Go Binary       │    │   Python Client   │  │
│  │                     │    │                   │  │
│  │  - Temporal Worker  │    │  - LCD rendering  │  │
│  │  - Ziggy Workflow   │◄──►│  - GPIO buttons   │  │
│  │  - HTTP API         │    │  - Query state    │  │
│  │  - Voice/AI logic   │    │  - Send signals   │  │
│  │                     │    │                   │  │
│  └─────────────────────┘    └───────────────────┘  │
│           │                          │             │
│           ▼                          ▼             │
│     Temporal Dev Server         LCD + Buttons      │
│     (UI exposed on network)                        │
└─────────────────────────────────────────────────────┘
```

### Why This Split?

- **Go for Temporal**: Lean, fast, compiles to single binary, familiar SDK
- **Python for hardware**: Best-in-class Pi ecosystem (GPIO, SPI displays)
- **Svelte for dev UI**: Fast iteration, your expertise, swappable later
- **Clean boundary**: Python/Web are thin clients, all game logic in Go

---

## Tech Stack

### Go Backend

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
| Temporal client | `temporalio` |
| Display | `st7789` + `Pillow` |
| GPIO/Buttons | `gpiozero` |
| Async | `asyncio` |

### Hardware BOM

| Part | Specific Model | ~Cost |
|------|----------------|-------|
| Pi | Raspberry Pi 4 Model B 2GB+ | $45 |
| Display | Pimoroni 1.3" SPI LCD (240x240, ST7789) | $15 |
| Buttons | 3x tactile buttons + caps | $5 |
| Speaker | (optional) PAM8403 amp + small speaker | $8 |
| Power | Official Pi 4 USB-C PSU | $10 |
| SD Card | 32GB | $10 |
| Case | 3D printed | $0 |
| Misc | Jumper wires, standoffs | $5 |

**Total: ~$100 per unit**

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
- **Sleep Cycle**: Ziggy sleeps at night (configurable hours). Different rules apply—waking Ziggy is bad for happiness.
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

---

## AI Integration

### Strategy: API with Pool Fallback

1. When online, call Claude API for fresh responses
2. When offline, fall back to pre-generated response pool
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

### Pool Entry Structure

```json
{
  "category": "feed_success",
  "entries": [
    {
      "text": "Mmm, perfect.\nI can survive\nanother eon.",
      "min_hunger": 30,
      "max_hunger": 70
    },
    {
      "text": "I was mass-\nextinction levels\nof hungry.",
      "min_hunger": 0,
      "max_hunger": 30
    }
  ]
}
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

Responses must be formatted for this constraint.

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
│   │   ├── ziggy_happy.png
│   │   ├── ziggy_neutral.png
│   │   ├── ziggy_hungry.png
│   │   ├── ziggy_sad.png
│   │   ├── ziggy_sleeping.png
│   │   ├── ziggy_critical.png
│   │   └── ziggy_tun.png
│   └── backgrounds/
│       ├── night.png
│       ├── dawn.png
│       ├── day.png
│       └── dusk.png
│
└── scripts/
    ├── generate_pool.go          # One-time pool generation
    └── install_pi.sh             # Pi setup script
```

---

## API Endpoints

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `GET /api/state` | GET | Query current Ziggy state |
| `POST /api/signal/feed` | POST | Send feed signal |
| `POST /api/signal/play` | POST | Send play signal |
| `POST /api/signal/pet` | POST | Send pet signal |
| `GET /api/health` | GET | Health check |
| `GET /` | GET | Serve web UI (production) |
| `GET /assets/*` | GET | Serve sprites/backgrounds |

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

## Visual Design

### Brand Alignment

Match Temporal's aesthetic:
- Deep purple/navy base (#1a1a2e, #2d1b4e)
- Subtle perspective grid lines (Tron-style)
- Soft teal accents (#4ade80)
- Muted pink/magenta highlights
- Forest silhouette elements (organic + digital contrast)
- Synthwave but subtle and professional

### Ziggy Character

Base design: Lavender/periwinkle tardigrade with:
- Soft rounded blob-like body
- 8 stubby legs with small gray claws
- Large expressive black eyes with white highlights
- Rosy pink cheeks
- Wide happy mouth (when content)
- Subtle body segment lines

### Sprite States Needed

**Idle states:**
- Happy (default, slight bounce pose)
- Neutral (relaxed)
- Hungry (droopy, pleading look)
- Sad (tears, drooping posture)

**Action states:**
- Eating (mouth open, happy)
- Playing (energetic, wiggling)
- Being petted (eyes closed, blissful)
- Sleeping (curled, eyes closed, "zzz")

**Health states:**
- Sick (green tinge, swirly eyes)
- Critical (pale, weak)
- Tun state (curled ball, grayscale)
- Reviving (uncurling, color returning)

**Evolution stages:**
- Egg
- Baby (tiny, bigger head ratio)
- Adult (standard)
- Elder (wise appearance)

### Background Variations

Four time-of-day versions of same scene:
1. **Night** (default) - deep purple, stars, teal glow
2. **Dawn** - warmer purple, pink horizon
3. **Day** - lighter purples, brighter
4. **Dusk** - orange/pink bleeding into purple

All maintain: subtle grid, forest silhouette, uncluttered center for character.

---

## Development Phases

### Phase 1: Web UI Foundation (Current)

- [ ] Scaffold Svelte + Vite project
- [ ] Create mock store with all game state
- [ ] Build core components (Game, Ziggy, Stats, Controls)
- [ ] Implement mock actions with state transitions
- [ ] Add DevTools for testing different states
- [ ] Placeholder sprites (colored rectangles OK)
- [ ] Get game feel right

### Phase 2: Go Backend

- [ ] Scaffold Go project structure
- [ ] Implement ZiggyState and signal types
- [ ] Create Ziggy workflow with basic lifecycle
- [ ] Add decay activities (scheduled)
- [ ] Implement queries for state
- [ ] Build HTTP API server
- [ ] Connect Svelte UI to real API
- [ ] Test with Temporal dev server

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

- [ ] Order Pi + components
- [ ] Design 3D printed case
- [ ] Write Python hardware client
- [ ] Test display rendering
- [ ] Wire up buttons
- [ ] Create Pi setup script
- [ ] Test full stack on device
- [ ] Document assembly for team

---

## Development Commands

```bash
# Terminal 1: Temporal server
temporal server start-dev

# Terminal 2: Go backend
go run main.go

# Terminal 3: Svelte frontend (dev mode)
cd web && npm run dev

# Build for production
cd web && npm run build
go build -o ziggy main.go
```

---

## Environment Variables

```bash
# Go backend
ANTHROPIC_API_KEY=sk-ant-...     # For AI responses
TEMPORAL_ADDRESS=localhost:7233   # Temporal server
HTTP_PORT=8080                    # API server port

# Optional
ZIGGY_WORKFLOW_ID=ziggy-dev       # Workflow ID
DECAY_INTERVAL=300                # Seconds between decay ticks
SLEEP_START_HOUR=22               # 10 PM
SLEEP_END_HOUR=7                  # 7 AM
```

---

## Workflow ID Naming

```
ziggy-{owner}-gen-{n}

Examples:
- ziggy-ross-gen-1
- ziggy-ross-gen-2  (after evolution/continue-as-new)
```

Visible in Temporal UI, makes lineage clear.

---

## Testing Checklist

### Game Mechanics
- [ ] Stats decay over time
- [ ] Feed increases hunger, caps at 100
- [ ] Play increases happiness, decreases hunger
- [ ] Pet increases bond
- [ ] HP drains when any stat is 0
- [ ] HP recovers when stats are healthy
- [ ] Cooldowns prevent spam
- [ ] Overfeeding causes sickness
- [ ] Sleep cycle activates at night
- [ ] Interacting during sleep has penalty

### Temporal Integration
- [ ] Workflow starts and persists
- [ ] Signals arrive and update state
- [ ] Queries return current state
- [ ] Scheduled activities fire correctly
- [ ] Workflow survives server restart
- [ ] Continue-as-new works for evolution
- [ ] UI shows workflow history accurately

### AI/Voice
- [ ] API calls succeed when online
- [ ] Fallback to pool when offline
- [ ] Responses match current state context
- [ ] Messages fit display constraints
- [ ] Pool responses feel varied

### Hardware (Phase 5)
- [ ] Display renders correctly
- [ ] Buttons register presses
- [ ] Signals sent on button press
- [ ] State queries update display
- [ ] Pi boots and auto-starts Ziggy
- [ ] Temporal UI accessible on network

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

- Temporal Go SDK: https://docs.temporal.io/develop/go
- Temporal Python SDK: https://docs.temporal.io/develop/python
- periph.io (Go hardware): https://periph.io
- ST7789 Python library: https://github.com/pimoroni/st7789-python
- gpiozero: https://gpiozero.readthedocs.io
- Anthropic Go SDK: https://github.com/anthropics/anthropic-sdk-go
- Svelte: https://svelte.dev
- Vite: https://vitejs.dev
