# Workflow Reorganization Plan

## Goals
1. Move each workflow into its own package
2. Ensure all business logic is in activities, not workflows
3. Maintain clean separation of concerns

## Proposed Structure

```
internal/
  ziggy/                    # Core shared types (used by all workflows)
    types.go                # Stage, Mood, Action, TimeOfDay, NeedType
    state.go                # ZiggyState, CareMetrics
    personality.go          # Personality types, DerivePersonality
    pool.go                 # MessagePool, PoolSelector
    pools_fallback.go       # Fallback pool definitions
    messages.go             # Message constants (messagesIdle)

  workflow/
    ziggy/
      workflow.go           # ZiggyWorkflow (signals, queries, activity calls)
      activities.go         # ProcessAction, RegeneratePool

    chat/
      workflow.go           # ChatWorkflow
      activities.go         # GenerateChatResponse, QueryZiggyState
      state.go              # ChatState, ChatMessage, ChatHistoryResponse
      mysteries.go          # Mystery types and definitions

    need_updater/
      workflow.go           # NeedUpdaterWorkflow

    pool_regenerator/
      workflow.go           # PoolRegeneratorWorkflow

    register.go             # RegisterWorkflows(), RegisterActivities()
```

## Logic to Move to Activities

### ChatWorkflow
Current inline logic to move:
- Mystery state updates after AI response (lines 227-261)
- Track selection logic (lines 103-112)
- Message adding and state updates

New activity: `ProcessChatMessage`
- Input: ChatState, message content, ZiggyState, track
- Does: Adds user message, calls AI, processes mystery updates, adds response
- Output: Updated ChatState (messages, mystery progress, hints, solved status)

This makes ChatWorkflow just: receive signal → call activity → update state

### NeedUpdaterWorkflow
- Already clean - uses QueryZiggyState activity and signals

### PoolRegeneratorWorkflow
- Already clean - uses RegeneratePool activity and signals

## Files to Create/Modify

### Phase 1: Create `internal/ziggy/` package with shared types
- `types.go` - Stage, Mood, Action, TimeOfDay, NeedType constants
- `state.go` - ZiggyState struct and methods
- `personality.go` - Personality type and DerivePersonality
- `pool.go` - MessagePool, PoolSelector, GetFallbackPool
- `pools_fallback.go` - Fallback pool data
- `messages.go` - messagesIdle map

### Phase 2: Reorganize workflows
- Move ziggy workflow to `internal/workflow/ziggy/`
- Move chat workflow to `internal/workflow/chat/`
- Move need_updater to `internal/workflow/need_updater/`
- Move pool_regenerator to `internal/workflow/pool_regenerator/`

### Phase 3: Update imports
- Update `cmd/worker.go` imports
- Update `internal/temporal/` if needed
- Update `register.go`

## Import Dependencies

```
internal/ziggy          # No dependencies on workflow packages
internal/workflow/ziggy # Imports internal/ziggy
internal/workflow/chat  # Imports internal/ziggy
internal/workflow/need_updater    # Imports internal/ziggy, signals ziggy workflow
internal/workflow/pool_regenerator # Imports internal/ziggy, signals ziggy workflow
```

## Implementation Steps

### Step 1: Create `internal/ziggy/` package
1. Create directory structure
2. Move shared types: Stage, Mood, Action, TimeOfDay, NeedType, Personality
3. Move ZiggyState and CareMetrics
4. Move MessagePool, PoolSelector, fallback pools
5. Move messages constants
6. Update package declarations

### Step 2: Refactor ChatWorkflow to activity-based
1. Create `ProcessChatMessage` activity in chat activities
2. Consolidate GenerateChatResponse + mystery logic into one activity
3. Simplify ChatWorkflow to just receive → activity → update state

### Step 3: Create workflow subpackages
1. Create `internal/workflow/ziggy/` with workflow.go and activities.go
2. Create `internal/workflow/chat/` with workflow.go, activities.go, state.go, mysteries.go
3. Create `internal/workflow/need_updater/` with workflow.go
4. Create `internal/workflow/pool_regenerator/` with workflow.go
5. Update register.go to import from subpackages

### Step 4: Update imports
1. Update cmd/worker.go to import new packages
2. Verify all cross-package imports work
3. Run `go build ./...` to verify

### Step 5: Test
1. Run `go test ./...`
2. Start worker and verify workflows function correctly
