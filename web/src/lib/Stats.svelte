<script lang="ts">
  interface Props {
    fullness: number;
    happiness: number;
    bond: number;
    hp: number;
  }

  let { fullness, happiness, bond, hp }: Props = $props();

  const stats = $derived([
    { label: 'HP', value: hp, color: '#ef4444' },
    { label: 'FUL', value: fullness, color: '#f59e0b' },
    { label: 'HAP', value: happiness, color: '#4ade80' },
    { label: 'BND', value: bond, color: '#ec4899' },
  ]);
</script>

<div class="stats">
  {#each stats as stat}
    <div class="stat-row">
      <span class="label">{stat.label}</span>
      <div class="bar-bg">
        <div
          class="bar-fill"
          style:width="{stat.value}%"
          style:background={stat.color}
          class:low={stat.value < 20}
          class:critical={stat.value < 10}
        ></div>
      </div>
      <span class="value">{Math.round(stat.value)}</span>
    </div>
  {/each}
</div>

<style>
  .stats {
    display: flex;
    flex-direction: column;
    gap: 4px;
    padding: 8px;
    background: rgba(26, 26, 46, 0.8);
    border-radius: 6px;
    font-family: monospace;
    font-size: 10px;
  }

  .stat-row {
    display: flex;
    align-items: center;
    gap: 6px;
  }

  .label {
    width: 28px;
    color: #a0a0b0;
    text-transform: uppercase;
  }

  .bar-bg {
    flex: 1;
    height: 8px;
    background: rgba(255, 255, 255, 0.1);
    border-radius: 4px;
    overflow: hidden;
  }

  .bar-fill {
    height: 100%;
    border-radius: 4px;
    transition:
      width 0.3s ease,
      background 0.3s;
  }

  .bar-fill.low {
    animation: pulse 1s ease-in-out infinite;
  }

  .bar-fill.critical {
    animation: pulse-fast 0.5s ease-in-out infinite;
  }

  @keyframes pulse {
    0%,
    100% {
      opacity: 1;
    }
    50% {
      opacity: 0.6;
    }
  }

  @keyframes pulse-fast {
    0%,
    100% {
      opacity: 1;
    }
    50% {
      opacity: 0.4;
    }
  }

  .value {
    width: 24px;
    text-align: right;
    color: #d0d0e0;
  }
</style>
