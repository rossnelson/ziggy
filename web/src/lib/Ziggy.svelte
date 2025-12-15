<script lang="ts">
  import type { Mood, Stage } from './store';

  interface Props {
    mood: Mood;
    stage: Stage;
  }

  let { mood, stage }: Props = $props();

  const SPRITE_WIDTH = 64;
  const SPRITE_HEIGHT = 64;
  const _COLS = 3; // Available for future use

  const moodToSprite: Record<Mood, { col: number; row: number }> = {
    happy: { col: 0, row: 0 },
    neutral: { col: 1, row: 0 },
    sad: { col: 2, row: 0 },
    hungry: { col: 2, row: 0 },
    lonely: { col: 2, row: 0 },
    sleeping: { col: 2, row: 1 },
    critical: { col: 0, row: 2 },
    tun: { col: 0, row: 2 },
  };

  const stageOverrides: Partial<Record<Stage, { col: number; row: number }>> = {
    egg: { col: 0, row: 3 },
    baby: { col: 1, row: 3 },
    elder: { col: 2, row: 3 },
  };

  let sprite = $derived(() => {
    if (stage === 'egg') return stageOverrides.egg!;
    if (stage === 'baby' && mood === 'sleeping') return { col: 2, row: 2 };
    if (stageOverrides[stage]) return stageOverrides[stage]!;
    return moodToSprite[mood];
  });

  let pos = $derived(sprite());

  let animationClass = $derived(
    mood === 'happy'
      ? 'bounce'
      : mood === 'sleeping'
        ? 'sleep'
        : mood === 'sad' || mood === 'hungry' || mood === 'lonely'
          ? 'droop'
          : mood === 'tun'
            ? 'curled'
            : 'idle'
  );

  let scale = $derived(stage === 'egg' ? 0.9 : stage === 'baby' ? 0.85 : 1);
</script>

<div
  class="ziggy {animationClass}"
  class:grayscale={mood === 'tun'}
  style:--sprite-x="-{pos.col * SPRITE_WIDTH}px"
  style:--sprite-y="-{pos.row * SPRITE_HEIGHT}px"
  style:--scale={scale}
></div>

<style>
  .ziggy {
    width: 64px;
    height: 64px;
    background-image: url('/assets/sprite.png');
    background-position: var(--sprite-x) var(--sprite-y);
    background-size: calc(64px * 3) calc(64px * 4);
    background-repeat: no-repeat;
    image-rendering: pixelated;
    transform: scale(var(--scale, 1));
  }

  .grayscale {
    filter: grayscale(0.8) brightness(0.7);
  }

  .bounce {
    animation: bounce 0.8s ease-in-out infinite;
  }

  @keyframes bounce {
    0%,
    100% {
      transform: scale(var(--scale, 1)) translateY(0);
    }
    50% {
      transform: scale(var(--scale, 1)) translateY(-5px);
    }
  }

  .sleep {
    animation: sleep-bob 2s ease-in-out infinite;
  }

  @keyframes sleep-bob {
    0%,
    100% {
      transform: scale(var(--scale, 1)) rotate(0deg);
    }
    50% {
      transform: scale(var(--scale, 1)) rotate(2deg);
    }
  }

  .droop {
    animation: droop 1.5s ease-in-out infinite;
  }

  @keyframes droop {
    0%,
    100% {
      transform: scale(var(--scale, 1)) translateY(0);
    }
    50% {
      transform: scale(var(--scale, 1)) translateY(3px);
    }
  }

  .curled {
    transform: scale(calc(var(--scale, 1) * 0.85));
  }

  .idle {
    animation: idle-sway 3s ease-in-out infinite;
  }

  @keyframes idle-sway {
    0%,
    100% {
      transform: scale(var(--scale, 1)) rotate(0deg);
    }
    25% {
      transform: scale(var(--scale, 1)) rotate(-1deg);
    }
    75% {
      transform: scale(var(--scale, 1)) rotate(1deg);
    }
  }
</style>
