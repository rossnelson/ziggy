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

<div class="flex flex-row sm:flex-col gap-2">
  <button
    class="action-btn group hover:border-amber-500"
    class:warning={isFull && !isSleeping}
    onclick={handleFeed}
    disabled={feedCooldown > 0 || isSleeping}
  >
    <span class="text-base">🍖</span>
    <span class="font-bold uppercase">Feed</span>
    <span class="shortcut">F</span>
    {#if isSleeping}
      <span class="status-badge">💤</span>
    {:else if feedCooldown > 0}
      <span class="status-badge text-amber-500">{formatCooldown(feedCooldown)}</span>
    {:else if isFull}
      <span class="status-badge text-red-500 font-bold">FULL</span>
    {/if}
  </button>

  <button
    class="action-btn group hover:border-green-400"
    onclick={handlePlay}
    disabled={playCooldown > 0 || isSleeping}
  >
    <span class="text-base">⚽</span>
    <span class="font-bold uppercase">Play</span>
    <span class="shortcut">P</span>
    {#if isSleeping}
      <span class="status-badge">💤</span>
    {:else if playCooldown > 0}
      <span class="status-badge text-amber-500">{formatCooldown(playCooldown)}</span>
    {/if}
  </button>

  <button
    class="action-btn group hover:border-pink-500"
    onclick={handlePet}
    disabled={petCooldown > 0 || isSleeping}
  >
    <span class="text-base">✋</span>
    <span class="font-bold uppercase">Pet</span>
    <span class="shortcut">T</span>
    {#if isSleeping}
      <span class="status-badge">💤</span>
    {:else if petCooldown > 0}
      <span class="status-badge text-amber-500">{formatCooldown(petCooldown)}</span>
    {/if}
  </button>

  {#if isSleeping}
    <button
      class="action-btn border-amber-400/50 bg-amber-400/10 hover:border-amber-400 hover:bg-amber-400/20"
      onclick={handleWake}
    >
      <span class="text-base">☀️</span>
      <span class="font-bold uppercase">Wake</span>
      <span class="shortcut">W</span>
      <span class="status-badge text-red-500 text-[8px] whitespace-nowrap">-10 HAP</span>
    </button>
  {/if}
</div>

<style>
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

  .action-btn:disabled:hover {
    transform: translateY(0);
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

  .status-badge {
    position: absolute;
    left: -28px;
    top: 50%;
    transform: translateY(-50%);
    font-size: 9px;
    white-space: nowrap;
  }

  .warning {
    border-color: rgba(239, 68, 68, 0.6);
    background: rgba(239, 68, 68, 0.1);
  }

  .warning:hover:not(:disabled) {
    border-color: #ef4444;
  }
</style>
