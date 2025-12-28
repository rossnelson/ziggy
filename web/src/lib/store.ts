import { writable, derived } from 'svelte/store';

export type Stage = 'egg' | 'baby' | 'teen' | 'adult' | 'elder';
export type Mood =
  | 'happy'
  | 'neutral'
  | 'hungry'
  | 'sad'
  | 'lonely'
  | 'sleeping'
  | 'critical'
  | 'tun';
export type TimeOfDay = 'night' | 'dawn' | 'day' | 'dusk';
export type Action = 'feed' | 'play' | 'pet' | 'wake';

export interface ZiggyState {
  fullness: number;
  happiness: number;
  bond: number;
  hp: number;
  stage: Stage;
  timeOfDay: TimeOfDay;
  sleeping: boolean;
  message: string;
  lastAction: Action | null;
  age: number;
  generation: number;
  feedCooldown: number;
  playCooldown: number;
  petCooldown: number;
}

const initialState: ZiggyState = {
  fullness: 70,
  happiness: 70,
  bond: 50,
  hp: 100,
  stage: 'egg',
  timeOfDay: 'day',
  sleeping: false,
  message: 'Loading...',
  lastAction: null,
  age: 0,
  generation: 1,
  feedCooldown: 0,
  playCooldown: 0,
  petCooldown: 0,
};

export const ziggyState = writable<ZiggyState>(initialState);

// Track when cooldowns were last synced from API for local countdown
let cooldownSyncedAt = 0;

export function syncCooldownTimestamp() {
  cooldownSyncedAt = Date.now();
}

export function getCooldownRemaining(action: Action): number {
  let cooldownSeconds = 0;
  ziggyState.subscribe((state) => {
    cooldownSeconds =
      action === 'feed'
        ? state.feedCooldown
        : action === 'play'
          ? state.playCooldown
          : action === 'pet'
            ? state.petCooldown
            : 0;
  })();

  const elapsedMs = Date.now() - cooldownSyncedAt;
  const remainingMs = cooldownSeconds * 1000 - elapsedMs;
  return Math.max(0, remainingMs);
}

// Derive mood from state
export const mood = derived(ziggyState, ($state): Mood => {
  if ($state.hp === 0) return 'tun';
  if ($state.sleeping) return 'sleeping';
  if ($state.hp < 20) return 'critical';
  if ($state.fullness < 20) return 'hungry';
  if ($state.happiness < 20) return 'sad';
  if ($state.bond < 20) return 'lonely';
  if ($state.happiness > 70 && $state.fullness > 50) return 'happy';
  return 'neutral';
});
