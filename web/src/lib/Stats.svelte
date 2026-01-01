<script lang="ts">
  import confetti from 'canvas-confetti';

  interface Props {
    fullness: number;
    happiness: number;
    bond: number;
    hp: number;
    onMaxHp?: () => void;
  }

  let { fullness, happiness, bond, hp, onMaxHp }: Props = $props();

  let previousHp = $state(hp);
  let hasReachedMax = $state(false);

  const stats = $derived([
    { label: 'HP', value: hp, color: '#ef4444', isHp: true },
    { label: 'FUL', value: fullness, color: '#f59e0b', isHp: false },
    { label: 'HAP', value: happiness, color: '#4ade80', isHp: false },
    { label: 'BND', value: bond, color: '#ec4899', isHp: false },
  ]);

  function celebrate() {
    const duration = 3000;
    const end = Date.now() + duration;

    const colors = ['#4ade80', '#a855f7', '#f59e0b', '#ec4899'];

    (function frame() {
      confetti({
        particleCount: 3,
        angle: 60,
        spread: 55,
        origin: { x: 0, y: 0.7 },
        colors,
      });
      confetti({
        particleCount: 3,
        angle: 120,
        spread: 55,
        origin: { x: 1, y: 0.7 },
        colors,
      });

      if (Date.now() < end) {
        requestAnimationFrame(frame);
      }
    })();

    onMaxHp?.();
  }

  $effect(() => {
    const roundedHp = Math.round(hp);
    const prevRoundedHp = Math.round(previousHp);

    if (roundedHp >= 100 && prevRoundedHp < 100 && !hasReachedMax) {
      hasReachedMax = true;
      celebrate();
    }
    if (roundedHp < 100) {
      hasReachedMax = false;
    }
    previousHp = hp;
  });
</script>

<div class="flex flex-col gap-1 p-2 bg-[rgba(26,26,46,0.8)] rounded-md font-mono text-[10px]">
  {#each stats as stat}
    <div class="flex items-center gap-1.5">
      <span class="w-7 text-[#a0a0b0] uppercase">{stat.label}</span>
      <div class="flex-1 h-2 bg-white/10 rounded overflow-hidden">
        <div
          class="h-full rounded transition-all duration-300"
          class:animate-pulse-slow={stat.value < 20 && stat.value >= 10}
          class:animate-pulse-fast={stat.value < 10}
          class:animate-glow={stat.isHp && Math.round(stat.value) >= 100}
          style:width="{Math.min(stat.value, 100)}%"
          style:background={stat.isHp && Math.round(stat.value) >= 100 ? 'linear-gradient(90deg, #4ade80, #a855f7)' : stat.color}
        ></div>
      </div>
      <span
        class="w-6 text-right"
        class:text-green-400={stat.isHp && Math.round(stat.value) >= 100}
        class:font-bold={stat.isHp && Math.round(stat.value) >= 100}
        class:text-[#d0d0e0]={!(stat.isHp && Math.round(stat.value) >= 100)}
      >{Math.round(stat.value)}</span>
    </div>
  {/each}
</div>

<style>
  @keyframes pulse-slow {
    0%, 100% { opacity: 1; }
    50% { opacity: 0.6; }
  }

  @keyframes pulse-fast {
    0%, 100% { opacity: 1; }
    50% { opacity: 0.4; }
  }

  @keyframes glow {
    0%, 100% {
      box-shadow: 0 0 4px #4ade80, 0 0 8px #a855f7;
    }
    50% {
      box-shadow: 0 0 8px #4ade80, 0 0 16px #a855f7, 0 0 24px #4ade80;
    }
  }

  .animate-pulse-slow {
    animation: pulse-slow 1s ease-in-out infinite;
  }

  .animate-pulse-fast {
    animation: pulse-fast 0.5s ease-in-out infinite;
  }

  .animate-glow {
    animation: glow 2s ease-in-out infinite;
  }
</style>
