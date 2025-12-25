<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { sendFeed, sendPlay, sendPet, sendWake } from './api';
  import { getCooldownRemaining, ziggyState } from './store';

  let feedCooldown = $state(0);
  let isFull = $derived($ziggyState.fullness > 90);
  let isSleeping = $derived($ziggyState.sleeping);
  let playCooldown = $state(0);
  let petCooldown = $state(0);

  let cooldownInterval: ReturnType<typeof setInterval> | null = null;

  function updateCooldowns() {
    feedCooldown = getCooldownRemaining('feed');
    playCooldown = getCooldownRemaining('play');
    petCooldown = getCooldownRemaining('pet');
  }

  async function handleFeed() {
    if (feedCooldown > 0) return;
    await sendFeed();
    updateCooldowns();
  }

  async function handlePlay() {
    if (playCooldown > 0) return;
    await sendPlay();
    updateCooldowns();
  }

  async function handlePet() {
    if (petCooldown > 0) return;
    await sendPet();
    updateCooldowns();
  }

  async function handleWake() {
    await sendWake();
  }

  function handleKeydown(event: KeyboardEvent) {
    if (event.target instanceof HTMLInputElement) return;

    if (event.key.toLowerCase() === 'w' && isSleeping) {
      handleWake();
      return;
    }

    if (isSleeping) return;

    switch (event.key.toLowerCase()) {
      case 'f':
        handleFeed();
        break;
      case 'p':
        handlePlay();
        break;
      case 't':
        handlePet();
        break;
    }
  }

  function formatCooldown(ms: number): string {
    if (ms <= 0) return '';
    const seconds = Math.ceil(ms / 1000);
    return `${seconds}s`;
  }

  onMount(() => {
    window.addEventListener('keydown', handleKeydown);
    cooldownInterval = setInterval(updateCooldowns, 100);
  });

  onDestroy(() => {
    window.removeEventListener('keydown', handleKeydown);
    if (cooldownInterval) clearInterval(cooldownInterval);
  });
</script>

<div class="controls">
  <button
    class="action-btn feed"
    class:warning={isFull && !isSleeping}
    onclick={handleFeed}
    disabled={feedCooldown > 0 || isSleeping}
  >
    <span class="icon">🍖</span>
    <span class="label">Feed</span>
    <span class="shortcut">F</span>
    {#if isSleeping}
      <span class="sleep-text">💤</span>
    {:else if feedCooldown > 0}
      <span class="cooldown">{formatCooldown(feedCooldown)}</span>
    {:else if isFull}
      <span class="warning-text">FULL</span>
    {/if}
  </button>

  <button class="action-btn play" onclick={handlePlay} disabled={playCooldown > 0 || isSleeping}>
    <span class="icon">⚽</span>
    <span class="label">Play</span>
    <span class="shortcut">P</span>
    {#if isSleeping}
      <span class="sleep-text">💤</span>
    {:else if playCooldown > 0}
      <span class="cooldown">{formatCooldown(playCooldown)}</span>
    {/if}
  </button>

  <button class="action-btn pet" onclick={handlePet} disabled={petCooldown > 0 || isSleeping}>
    <span class="icon">✋</span>
    <span class="label">Pet</span>
    <span class="shortcut">T</span>
    {#if isSleeping}
      <span class="sleep-text">💤</span>
    {:else if petCooldown > 0}
      <span class="cooldown">{formatCooldown(petCooldown)}</span>
    {/if}
  </button>

  {#if isSleeping}
    <button class="action-btn wake" onclick={handleWake}>
      <span class="icon">☀️</span>
      <span class="label">Wake</span>
      <span class="shortcut">W</span>
      {#if isSleeping}
        <span class="penalty">-10 HAP</span>
      {/if}
    </button>
  {/if}
</div>

<style>
  .controls {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .action-btn {
    position: relative;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 2px;
    padding: 8px 12px;
    min-width: 60px;
    background: rgba(26, 26, 46, 0.9);
    border: 2px solid rgba(74, 222, 128, 0.3);
    border-radius: 8px;
    color: #d0d0e0;
    font-family: monospace;
    font-size: 10px;
    cursor: pointer;
    transition: all 0.2s;
  }

  .action-btn:hover:not(:disabled) {
    border-color: rgba(74, 222, 128, 0.6);
    background: rgba(26, 26, 46, 1);
    transform: translateY(-2px);
  }

  .action-btn:active:not(:disabled) {
    transform: translateY(0);
  }

  .action-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .icon {
    font-size: 16px;
  }

  .label {
    font-weight: bold;
    text-transform: uppercase;
  }

  .shortcut {
    position: absolute;
    top: 4px;
    right: 4px;
    font-size: 8px;
    color: rgba(74, 222, 128, 0.6);
    background: rgba(74, 222, 128, 0.1);
    padding: 1px 4px;
    border-radius: 3px;
  }

  .cooldown {
    position: absolute;
    left: -28px;
    top: 50%;
    transform: translateY(-50%);
    font-size: 9px;
    color: #f59e0b;
    white-space: nowrap;
  }

  .feed:hover:not(:disabled) {
    border-color: #f59e0b;
  }

  .play:hover:not(:disabled) {
    border-color: #4ade80;
  }

  .pet:hover:not(:disabled) {
    border-color: #ec4899;
  }

  .warning {
    border-color: rgba(239, 68, 68, 0.6);
    background: rgba(239, 68, 68, 0.1);
  }

  .warning:hover:not(:disabled) {
    border-color: #ef4444;
  }

  .warning-text {
    position: absolute;
    left: -32px;
    top: 50%;
    transform: translateY(-50%);
    font-size: 8px;
    color: #ef4444;
    font-weight: bold;
  }

  .sleep-text {
    position: absolute;
    left: -24px;
    top: 50%;
    transform: translateY(-50%);
    font-size: 10px;
  }

  .wake {
    border-color: rgba(251, 191, 36, 0.5);
    background: rgba(251, 191, 36, 0.1);
  }

  .wake:hover {
    border-color: #fbbf24;
    background: rgba(251, 191, 36, 0.2);
  }

  .penalty {
    position: absolute;
    left: -48px;
    top: 50%;
    transform: translateY(-50%);
    font-size: 8px;
    color: #ef4444;
    white-space: nowrap;
  }
</style>
