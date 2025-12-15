import { writable, derived, get } from 'svelte/store';

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
export type Action = 'feed' | 'play' | 'pet';

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
  lastActionTime: number | null;
  age: number;
  generation: number;
}

interface Cooldowns {
  feed: number | null;
  play: number | null;
  pet: number | null;
}

const COOLDOWN_MS = {
  feed: 30_000,
  play: 60_000,
  pet: 10_000,
};

const DECAY_INTERVAL_MS = 2_000; // Fast for testing (normally 10_000)
const DECAY_AMOUNTS = {
  fullness: 2,
  happiness: 1,
  bond: 0.5,
};

function getInitialTimeOfDay(): TimeOfDay {
  return getTimeOfDayFromHour(new Date().getHours());
}

const initialState: ZiggyState = {
  fullness: 70,
  happiness: 70,
  bond: 50,
  hp: 100,
  stage: 'adult',
  timeOfDay: getInitialTimeOfDay(),
  sleeping: getInitialTimeOfDay() === 'night',
  message:
    getInitialTimeOfDay() === 'night'
      ? 'Zzz... cosmic\ndreams... zzz'
      : "I've survived\nworse than this.\nBut barely.",
  lastAction: null,
  lastActionTime: null,
  age: 0,
  generation: 1,
};

export const ziggyState = writable<ZiggyState>(initialState);
export const cooldowns = writable<Cooldowns>({ feed: null, play: null, pet: null });

let decayTimer: ReturnType<typeof setInterval> | null = null;
let tickCount = 0;

function getTimeOfDayFromHour(hour: number): TimeOfDay {
  if (hour >= 22 || hour < 5) return 'night';
  if (hour >= 5 && hour < 8) return 'dawn';
  if (hour >= 8 && hour < 18) return 'day';
  return 'dusk';
}

function updateTimeOfDay() {
  const hour = new Date().getHours();
  const newTimeOfDay = getTimeOfDayFromHour(hour);
  const shouldSleep = newTimeOfDay === 'night';

  ziggyState.update((state) => {
    if (state.timeOfDay === newTimeOfDay && state.sleeping === shouldSleep) {
      return state;
    }
    return {
      ...state,
      timeOfDay: newTimeOfDay,
      sleeping: shouldSleep,
    };
  });
}

function getMoodFromState(state: ZiggyState): Mood {
  if (state.hp === 0) return 'tun';
  if (state.sleeping) return 'sleeping';
  if (state.hp < 20) return 'critical';
  if (state.fullness < 20) return 'hungry';
  if (state.happiness < 20) return 'sad';
  if (state.bond < 20) return 'lonely';
  if (state.happiness > 70 && state.fullness > 50) return 'happy';
  return 'neutral';
}

export function startDecay() {
  if (decayTimer) return;
  updateTimeOfDay();
  decayTimer = setInterval(() => {
    tickCount++;
    updateTimeOfDay();
    ziggyState.update((state) => {
      if (state.sleeping) {
        // During sleep: slow fullness decay, slow happiness recovery
        const newFullness = Math.max(0, state.fullness - DECAY_AMOUNTS.fullness * 0.5);
        const newHappiness = Math.min(100, state.happiness + 0.5);
        const newBond = state.bond; // Bond doesn't change during sleep

        const targetHp = Math.round((newFullness + newHappiness + newBond) / 3);
        let newHp = state.hp;
        if (state.hp < targetHp) {
          newHp = Math.min(100, state.hp + 1); // Recover HP faster during sleep
        }

        return {
          ...state,
          fullness: newFullness,
          happiness: newHappiness,
          hp: newHp,
          age: state.age + DECAY_INTERVAL_MS / 1000,
        };
      }

      const bondProtection = state.bond > 50 ? (state.bond - 50) / 100 : 0;
      const fullnessDecay = DECAY_AMOUNTS.fullness * (1 - bondProtection);
      const happinessDecay = DECAY_AMOUNTS.happiness * (1 - bondProtection);

      const newFullness = Math.max(0, state.fullness - fullnessDecay);
      const newHappiness = Math.max(0, state.happiness - happinessDecay);
      const newBond = Math.max(0, state.bond - DECAY_AMOUNTS.bond);

      const targetHp = Math.round((newFullness + newHappiness + newBond) / 3);
      let newHp = state.hp;

      if (state.hp > targetHp) {
        newHp = Math.max(0, state.hp - 2);
      } else if (state.hp < targetHp) {
        newHp = Math.min(100, state.hp + 1);
      }

      const newState = {
        ...state,
        fullness: newFullness,
        happiness: newHappiness,
        bond: newBond,
        hp: newHp,
        age: state.age + DECAY_INTERVAL_MS / 1000,
      };

      const currentMood = getMoodFromState(newState);

      // Always update message on critical state changes or every 3 ticks
      const enteredTun = newHp === 0 && state.hp > 0;
      const enteredCritical = newHp < 20 && state.hp >= 20;

      if (enteredTun || enteredCritical || tickCount % 3 === 0) {
        const idleMessages = MESSAGES.idle[currentMood] ?? MESSAGES.idle.neutral;
        newState.message = pickRandom(idleMessages);
      }

      return newState;
    });
  }, DECAY_INTERVAL_MS);
}

export function stopDecay() {
  if (decayTimer) {
    clearInterval(decayTimer);
    decayTimer = null;
  }
}

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

const MESSAGES: Record<string, Record<string, string[]>> = {
  feed: {
    success: [
      'Mmm, perfect.\nI can survive\nanother eon.',
      'Delicious.\nAlmost as good\nas cosmic dust.',
      'I needed that.\nThanks, human.',
    ],
    full: [
      'Too full!\nNow I feel sick.\n*happiness down*',
      'Ugh... stuffed.\nThat made me\nunhappy.',
      'No more!\nOverfeeding\nhurts me.',
    ],
    hungry: [
      'FINALLY.\nI was mass-\nextinction hungry.',
      "Oh thank you.\nI thought I'd\nstarve forever.",
      'Food! Beautiful\nlife-giving food!',
    ],
    sleeping: ['Zzz... not now...\nzzz...', '*mumbles*\nfive more eons...'],
  },
  play: {
    success: [
      'Wheee!\nThis is fun!',
      'Again! Again!\nI have energy\nfor eons!',
      'Playing is the\nbest survival\nstrategy.',
    ],
    tired: ["I'm too tired...\nmaybe later?", 'Need food first.\nThen play.'],
    happy: ['More playing!\nI love this!', 'Best day since\nthe Permian\nextinction!'],
    sleeping: ['Zzz... playing\nin my dreams...', '*sleep-wiggles*'],
  },
  pet: {
    success: ['*happy wiggle*\nI like you.', "That's nice.\nKeep going.", 'Mmm...\nright there.'],
    maxBond: [
      "We're already\nbest friends!\nBut okay...",
      '*content sigh*\nI trust you\ncompletely.',
    ],
    low_mood: ['Thanks...\nI needed that.', "*small wiggle*\nYou're kind."],
    sleeping: ['Zzz...\n*happy mumble*', '*snuggles closer*'],
  },
  idle: {
    happy: [
      "Life is good.\nI've survived\nworse.",
      'Did you know\nI can live in\nspace? Cool, right?',
      'Just vibing.\nDurably.',
    ],
    neutral: ['...', "I'm fine.\nJust existing.", 'Waiting for\nsomething to\nsurvive.'],
    hungry: ['My stomach\nis a void.', 'Feed me?\nPlease?', "I'm withering\naway here..."],
    sad: ['Nobody loves\na tardigrade...', '*sad wiggle*', "I'm fine.\nEverything is\nfine."],
    lonely: [
      'Is anyone\nthere...?',
      'I miss you.\nCome back soon.',
      '*looks around*\nSo quiet...',
      'Even tardigrades\nneed friends.',
    ],
    critical: ["I don't feel\nso good...", 'Help...', 'Is this how\nit ends?'],
    tun: ['*curled up*\n*not responding*'],
    sleeping: ['Zzz...', '*peaceful snoring*', 'Zzz... cosmic\ndreams... zzz'],
  },
};

function pickRandom<T>(arr: T[]): T {
  return arr[Math.floor(Math.random() * arr.length)];
}

function getMessage(action: Action, state: ZiggyState, currentMood: Mood): string {
  const actionMessages = MESSAGES[action];

  if (state.sleeping) {
    return pickRandom(actionMessages.sleeping);
  }

  if (action === 'feed') {
    if (state.fullness > 90) return pickRandom(actionMessages.full);
    if (state.fullness < 30) return pickRandom(actionMessages.hungry);
    return pickRandom(actionMessages.success);
  }

  if (action === 'play') {
    if (state.fullness < 20 || state.hp < 30) return pickRandom(actionMessages.tired);
    if (currentMood === 'happy') return pickRandom(actionMessages.happy);
    return pickRandom(actionMessages.success);
  }

  if (action === 'pet') {
    if (state.bond > 90) return pickRandom(actionMessages.maxBond);
    if (currentMood === 'sad' || currentMood === 'hungry')
      return pickRandom(actionMessages.low_mood);
    return pickRandom(actionMessages.success);
  }

  return pickRandom(MESSAGES.idle[currentMood] ?? MESSAGES.idle.neutral);
}

function isOnCooldown(action: Action): boolean {
  const cd = get(cooldowns);
  const lastTime = cd[action];
  if (!lastTime) return false;
  return Date.now() - lastTime < COOLDOWN_MS[action];
}

function setCooldown(action: Action) {
  cooldowns.update((cd) => ({ ...cd, [action]: Date.now() }));
}

export function getCooldownRemaining(action: Action): number {
  const cd = get(cooldowns);
  const lastTime = cd[action];
  if (!lastTime) return 0;
  const remaining = COOLDOWN_MS[action] - (Date.now() - lastTime);
  return Math.max(0, remaining);
}

export function feed(): boolean {
  const state = get(ziggyState);
  if (state.sleeping) return false;
  if (isOnCooldown('feed')) return false;

  const currentMood = get(mood);

  const wasOverfed = state.fullness > 90;
  const bondProtection = state.bond > 50 ? (state.bond - 50) / 20 : 0;
  const fullnessGain = wasOverfed ? 5 : 25;
  const happinessChange = wasOverfed ? Math.floor(-15 + bondProtection) : 5;

  ziggyState.update((s) => ({
    ...s,
    fullness: Math.min(100, s.fullness + fullnessGain),
    happiness: Math.max(0, Math.min(100, s.happiness + happinessChange)),
    message: getMessage('feed', s, currentMood),
    lastAction: 'feed',
    lastActionTime: Date.now(),
  }));

  setCooldown('feed');
  return true;
}

export function play(): boolean {
  const state = get(ziggyState);
  if (state.sleeping) return false;
  if (isOnCooldown('play')) return false;

  const currentMood = get(mood);

  const tooTired = state.fullness < 20 || state.hp < 30;
  const happinessGain = tooTired ? 5 : 20;
  const fullnessCost = tooTired ? 5 : 10;

  ziggyState.update((s) => ({
    ...s,
    happiness: Math.min(100, s.happiness + happinessGain),
    fullness: Math.max(0, s.fullness - fullnessCost),
    bond: Math.min(100, s.bond + 5),
    message: getMessage('play', s, currentMood),
    lastAction: 'play',
    lastActionTime: Date.now(),
  }));

  setCooldown('play');
  return true;
}

export function pet(): boolean {
  const state = get(ziggyState);
  if (state.sleeping) return false;
  if (isOnCooldown('pet')) return false;

  const currentMood = get(mood);

  ziggyState.update((s) => ({
    ...s,
    bond: Math.min(100, s.bond + 10),
    happiness: Math.min(100, s.happiness + 5),
    message: getMessage('pet', s, currentMood),
    lastAction: 'pet',
    lastActionTime: Date.now(),
  }));

  setCooldown('pet');
  return true;
}

export function setTimeOfDay(time: TimeOfDay) {
  const shouldSleep = time === 'night';
  ziggyState.update((s) => ({
    ...s,
    timeOfDay: time,
    sleeping: shouldSleep,
    message: shouldSleep
      ? pickRandom(MESSAGES.idle.sleeping)
      : 'Good morning!\nI dreamed of\nsurviving things.',
  }));
}

export function setSleeping(sleeping: boolean) {
  ziggyState.update((s) => ({
    ...s,
    sleeping,
    message: sleeping
      ? pickRandom(MESSAGES.idle.sleeping)
      : 'Good morning!\nI dreamed of\nsurviving things.',
  }));
}

export function setStage(stage: Stage) {
  ziggyState.update((s) => ({ ...s, stage }));
}

export function setStat(stat: 'fullness' | 'happiness' | 'bond' | 'hp', value: number) {
  ziggyState.update((s) => {
    const newState = { ...s, [stat]: Math.max(0, Math.min(100, value)) };
    const newMood = getMoodFromState(newState);
    const idleMessages = MESSAGES.idle[newMood] ?? MESSAGES.idle.neutral;
    newState.message = pickRandom(idleMessages);
    return newState;
  });
}

export function wake(): boolean {
  const state = get(ziggyState);
  if (!state.sleeping) return false;

  ziggyState.update((s) => ({
    ...s,
    sleeping: false,
    happiness: Math.max(0, s.happiness - 10),
    message: '*yawn*\nI was having\nsuch a nice dream...',
  }));
  return true;
}

export function resetState() {
  ziggyState.set(initialState);
  cooldowns.set({ feed: null, play: null, pet: null });
}
