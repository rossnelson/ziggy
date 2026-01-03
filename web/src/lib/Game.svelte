<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { ziggyState, mood } from './store';
  import { startSSE, stopSSE, getConfig, aiEnabled } from './api';
  import Background from './Background.svelte';
  import Ziggy from './Ziggy.svelte';
  import Stats from './Stats.svelte';
  import Message from './Message.svelte';
  import Controls from './Controls.svelte';
  import Chat from './Chat.svelte';
  import ChatDrawer from './ChatDrawer.svelte';

  let loading = $state(true);
  let celebrationMessage = $state('');
  let showCelebration = $state(false);

  const celebrationMessages = [
    "WOW! 100 HP! You're the best friend ever! 🎉",
    "I feel AMAZING! Thank you for caring so much! ✨",
    "Maximum power! I'm so happy right now! 💚",
    "You did it! I'm at full health! 🌟",
    "I've never felt better! You're incredible! 💜",
  ];

  function handleMaxHp() {
    celebrationMessage = celebrationMessages[Math.floor(Math.random() * celebrationMessages.length)];
    showCelebration = true;
    setTimeout(() => {
      showCelebration = false;
    }, 5000);
  }

  let isAiEnabled = $derived($aiEnabled);

  onMount(async () => {
    await getConfig();
    startSSE();
    loading = false;
  });

  onDestroy(() => {
    stopSSE();
  });
</script>

{#if loading}
  <div class="flex justify-center items-center min-h-screen p-5 bg-[#0a0a12]">
    <div class="text-green-400/80 font-mono text-sm">Loading...</div>
  </div>
{:else}
  <div class="flex justify-center items-center min-h-screen p-5 bg-[#0a0a12]">
    <div class="flex flex-col sm:flex-row gap-3 items-center sm:items-start">
      <!-- Controls: below canvas on mobile, left side on desktop -->
      <div class="order-2 sm:order-1 flex flex-col justify-center">
        <Controls />
      </div>

      <!-- Game Canvas -->
      <div class="order-1 sm:order-2 game-canvas">
        <Background timeOfDay={$ziggyState.timeOfDay} />

        <div class="absolute inset-0 flex flex-col p-2">
          <div class="z-10">
            <Stats
              fullness={$ziggyState.fullness}
              happiness={$ziggyState.happiness}
              bond={$ziggyState.bond}
              hp={$ziggyState.hp}
              onMaxHp={handleMaxHp}
            />
          </div>

          <div class="absolute top-[105px] left-1/2 -translate-x-1/2 z-10">
            {#if showCelebration}
              <div class="celebration-message">
                {celebrationMessage}
              </div>
            {:else}
              <Message message={$ziggyState.message} />
            {/if}
          </div>

          <div class="absolute top-[175px] left-1/2 -translate-x-1/2 z-5">
            <Ziggy mood={$mood} stage={$ziggyState.stage} />
          </div>
        </div>
      </div>

      <!-- Chat: hidden on mobile (drawer replaces it), visible on desktop -->
      {#if isAiEnabled}
        <div class="order-3 flex flex-col items-start">
          <Chat />
        </div>
      {:else if isAiEnabled === false}
        <div class="order-3 hidden sm:flex w-[280px] h-60 bg-[rgba(26,26,46,0.95)] border-2 border-green-400/30 rounded-lg flex-col items-center justify-center">
          <span class="text-red-400/80 font-mono text-[10px] text-center px-4">Chat unavailable<br/><span class="text-[9px] text-[#a0a0b0]">AI not configured</span></span>
        </div>
      {/if}
    </div>

    <!-- Mobile chat drawer -->
    {#if isAiEnabled}
      <ChatDrawer />
    {/if}
  </div>
{/if}

<style>
  .game-canvas {
    position: relative;
    width: 240px;
    height: 240px;
    border: 2px solid rgba(74, 222, 128, 0.3);
    border-radius: 8px;
    overflow: hidden;
    box-shadow:
      0 0 30px rgba(74, 222, 128, 0.1),
      0 10px 40px rgba(0, 0, 0, 0.5);
  }

  .celebration-message {
    background: linear-gradient(135deg, rgba(74, 222, 128, 0.95), rgba(168, 85, 247, 0.95));
    border-radius: 8px;
    padding: 10px 14px;
    max-width: 220px;
    font-family: monospace;
    font-size: 11px;
    color: white;
    text-align: center;
    box-shadow: 0 0 20px rgba(74, 222, 128, 0.5), 0 0 40px rgba(168, 85, 247, 0.3);
    animation: celebrate-pulse 1s ease-in-out infinite;
  }

  @keyframes celebrate-pulse {
    0%, 100% {
      transform: scale(1);
      box-shadow: 0 0 20px rgba(74, 222, 128, 0.5), 0 0 40px rgba(168, 85, 247, 0.3);
    }
    50% {
      transform: scale(1.02);
      box-shadow: 0 0 30px rgba(74, 222, 128, 0.7), 0 0 60px rgba(168, 85, 247, 0.5);
    }
  }
</style>
