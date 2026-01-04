<script lang="ts">
  import { onMount } from 'svelte';
  import {
    sendChatMessage,
    getAvailableMysteries,
    startMystery,
    chatMessages,
    chatLoading,
    mysteryStatus,
    type Mystery,
  } from './api';

  let inputValue = $state('');
  let mysteries = $state<Mystery[]>([]);
  let showMysteries = $state(false);
  let track = $state<'fun' | 'educational'>('fun');
  let messagesContainer: HTMLDivElement;

  let messages = $derived($chatMessages);
  let isLoading = $derived($chatLoading);
  let mystery = $derived($mysteryStatus);

  type DrawerState = 'collapsed' | 'peek' | 'half' | 'full';
  let drawerState = $state<DrawerState>('collapsed');
  let isDragging = $state(false);
  let drawerHeight = $state(48);
  let startY = 0;
  let startHeight = 0;

  function getClosestState(height: number): DrawerState {
    const vh = typeof window !== 'undefined' ? window.innerHeight : 800;
    if (height < 80) return 'collapsed';
    if (height < vh * 0.3) return 'peek';
    if (height < vh * 0.65) return 'half';
    return 'full';
  }

  function snapTo(state: DrawerState) {
    drawerState = state;
    const vh = typeof window !== 'undefined' ? window.innerHeight : 800;
    drawerHeight = state === 'collapsed' ? 48
      : state === 'peek' ? 120
      : state === 'half' ? vh * 0.5
      : vh * 0.85;
  }

  function onTouchStart(e: TouchEvent) {
    isDragging = true;
    startY = e.touches[0].clientY;
    startHeight = drawerHeight;
  }

  function onTouchMove(e: TouchEvent) {
    if (!isDragging) return;
    const deltaY = startY - e.touches[0].clientY;
    const vh = typeof window !== 'undefined' ? window.innerHeight : 800;
    drawerHeight = Math.max(48, Math.min(vh * 0.9, startHeight + deltaY));
  }

  function onTouchEnd() {
    isDragging = false;
    snapTo(getClosestState(drawerHeight));
  }

  function toggleDrawer() {
    if (drawerState === 'collapsed') {
      snapTo('half');
    } else {
      snapTo('collapsed');
    }
  }

  async function loadMysteries() {
    const result = await getAvailableMysteries(track);
    if (result.success && result.data) {
      mysteries = result.data;
    }
  }

  function toggleTrack() {
    track = track === 'fun' ? 'educational' : 'fun';
    loadMysteries();
  }

  async function handleSend() {
    if (!inputValue.trim() || isLoading) return;
    const content = inputValue.trim();
    inputValue = '';
    chatLoading.set(true);
    await sendChatMessage(content);
    chatLoading.set(false);
    scrollToBottom();
  }

  async function handleStartMystery(m: Mystery) {
    showMysteries = false;
    chatLoading.set(true);
    await startMystery(m.id, track);
    const message = track === 'educational'
      ? `let's learn about "${m.title}"`
      : `let's solve the mystery "${m.title}"`;
    await sendChatMessage(message);
    chatLoading.set(false);
    scrollToBottom();
  }

  function handleKeydown(event: KeyboardEvent) {
    if (event.key === 'Enter' && !event.shiftKey) {
      event.preventDefault();
      handleSend();
    }
  }

  let prevMessageCount = $state(0);

  function scrollToBottom() {
    if (messagesContainer) {
      setTimeout(() => {
        messagesContainer.scrollTop = messagesContainer.scrollHeight;
      }, 50);
    }
  }

  // Scroll to bottom only when new messages are added
  $effect(() => {
    if (messages.length > prevMessageCount) {
      scrollToBottom();
    }
    prevMessageCount = messages.length;
  });

  onMount(() => {
    loadMysteries();
  });

  function formatMessageWithLinks(content: string): string {
    const urlRegex = /(https?:\/\/[^\s]+)/g;
    return content.replace(urlRegex, '<a href="$1" target="_blank" rel="noopener noreferrer" class="text-green-400 underline hover:text-green-300">$1</a>');
  }
</script>

<div
  class="fixed inset-x-0 bottom-0 z-50 sm:hidden bg-[rgba(26,26,46,0.98)] border-t-2 border-green-400/30 rounded-t-xl drawer"
  class:dragging={isDragging}
  style:height="{drawerHeight}px"
>
  <!-- Grab bar with collapsed label -->
  <div
    class="h-12 flex justify-center items-center cursor-grab touch-none relative"
    ontouchstart={onTouchStart}
    ontouchmove={onTouchMove}
    ontouchend={onTouchEnd}
    onclick={toggleDrawer}
    role="button"
    tabindex="0"
    onkeydown={(e) => e.key === 'Enter' && toggleDrawer()}
  >
    <div class="w-10 h-1 bg-green-400/50 rounded-full absolute top-2"></div>
    {#if drawerState === 'collapsed'}
      <span class="text-green-400/70 font-mono text-[10px] mt-2">Chat with Ziggy</span>
    {/if}
  </div>

  <!-- Chat content (hidden when collapsed) -->
  <div class="flex flex-col h-[calc(100%-48px)] overflow-hidden" class:hidden={drawerState === 'collapsed'}>
    <!-- Header -->
    <div class="flex justify-between items-center px-3 py-1.5 border-b border-green-400/20 shrink-0">
      <span class="text-green-400 font-mono text-[11px] font-bold">Chat with Ziggy</span>
      <div class="flex gap-1">
        <button
          class="bg-green-400/10 border border-green-400/30 rounded px-2 py-0.5 text-green-400 font-mono text-[9px] cursor-pointer hover:bg-green-400/20"
          onclick={toggleTrack}
          title={track === 'fun' ? 'Switch to educational mode' : 'Switch to fun mode'}
        >
          {track === 'fun' ? '🎮' : '📚'}
        </button>
        <button
          class="border rounded px-2 py-0.5 text-green-400 font-mono text-[9px] cursor-pointer transition-colors {showMysteries ? 'bg-green-400/30 border-green-400/60' : 'bg-green-400/10 border-green-400/30 hover:bg-green-400/20'}"
          onclick={() => (showMysteries = !showMysteries)}
        >
          {track === 'educational' ? 'Topics' : 'Mysteries'}
        </button>
      </div>
    </div>

    {#if mystery?.active}
      <div class="flex items-center gap-2 px-3 py-1.5 bg-purple-500/10 border-b border-purple-500/20 shrink-0">
        <span class="font-mono text-[9px] font-bold text-purple-500 whitespace-nowrap">{mystery.mystery?.title}</span>
        {#if mystery.mystery?.track !== 'educational'}
          <div class="flex-1 h-1.5 bg-black/30 rounded-sm overflow-hidden">
            <div
              class="h-full bg-gradient-to-r from-purple-500 to-green-400 rounded-sm transition-all duration-300"
              style="width: {(mystery.progress / mystery.totalHints) * 100}%"
            ></div>
          </div>
          <span class="font-mono text-[8px] text-[#a0a0b0] whitespace-nowrap">{mystery.progress}/{mystery.totalHints}</span>
        {/if}
      </div>
    {:else if showMysteries}
      <div class="max-h-32 overflow-y-auto bg-[#1a1a2c] border-b border-green-400/20">
        {#each mysteries as m}
          <button
            class="flex flex-col w-full px-3 py-2 bg-transparent border-none border-b border-green-400/10 text-[#d0d0e0] font-mono text-left cursor-pointer hover:bg-green-400/10"
            onclick={() => handleStartMystery(m)}
          >
            <span class="text-[10px] font-bold text-green-400">{m.title}</span>
            <span class="text-[9px] text-[#a0a0b0] mt-0.5">{m.description}</span>
          </button>
        {/each}
        {#if mysteries.length === 0}
          <div class="p-3 text-center text-[#a0a0b0] text-[10px] font-mono">No mysteries available</div>
        {/if}
      </div>
    {/if}

    <!-- Messages -->
    <div class="flex-1 overflow-y-auto p-2 flex flex-col gap-1.5" bind:this={messagesContainer}>
      {#each messages as message}
        <div class="flex flex-col gap-0.5 px-2 py-1.5 rounded-md font-mono text-[10px] max-w-[85%] {message.role === 'user' ? 'bg-green-400/15 self-end' : 'bg-purple-600/15 self-start'}">
          <span class="text-[8px] font-bold uppercase {message.role === 'user' ? 'text-green-400' : 'text-purple-500'}">
            {message.role === 'user' ? 'You' : 'Ziggy'}
          </span>
          <span class="text-[#e0e0e0] break-words whitespace-pre-wrap">{@html formatMessageWithLinks(message.content)}</span>
        </div>
      {/each}
      {#if messages.length === 0 && !isLoading}
        <div class="text-center text-[#a0a0b0] font-mono text-[10px] py-3">Say hi to Ziggy!</div>
      {/if}
      {#if isLoading}
        <div class="flex flex-col gap-0.5 px-2 py-1.5 rounded-md font-mono text-[10px] bg-purple-600/15 self-start max-w-[85%]">
          <span class="text-[8px] font-bold uppercase text-purple-500">Ziggy</span>
          <span class="typing">
            <span class="dot"></span>
            <span class="dot"></span>
            <span class="dot"></span>
          </span>
        </div>
      {/if}
    </div>

    <!-- Input -->
    <div class="flex gap-1.5 p-2 border-t border-green-400/20 shrink-0">
      <input
        type="text"
        placeholder="Type a message..."
        class="flex-1 px-2 py-1.5 bg-black/30 border border-green-400/30 rounded text-[#e0e0e0] font-mono text-[10px] placeholder:text-[#606080] focus:outline-none focus:border-green-400/60 disabled:opacity-60"
        bind:value={inputValue}
        onkeydown={handleKeydown}
        disabled={isLoading}
      />
      <button
        class="px-3 py-1.5 bg-green-400/20 border border-green-400/40 rounded text-green-400 font-mono text-[10px] font-bold cursor-pointer transition-all min-w-[50px] hover:bg-green-400/30 hover:border-green-400/60 disabled:opacity-50 disabled:cursor-not-allowed"
        onclick={handleSend}
        disabled={isLoading || !inputValue.trim()}
      >
        {isLoading ? '...' : 'Send'}
      </button>
    </div>
  </div>
</div>

<style>
  .drawer {
    transition: height 0.3s cubic-bezier(0.32, 0.72, 0, 1);
  }

  .drawer.dragging {
    transition: none;
  }

  .typing {
    display: flex;
    gap: 4px;
    padding: 4px 0;
  }

  .dot {
    width: 6px;
    height: 6px;
    background: #a855f7;
    border-radius: 50%;
    animation: bounce 1.4s infinite ease-in-out both;
  }

  .dot:nth-child(1) { animation-delay: -0.32s; }
  .dot:nth-child(2) { animation-delay: -0.16s; }

  @keyframes bounce {
    0%, 80%, 100% {
      transform: scale(0);
      opacity: 0.5;
    }
    40% {
      transform: scale(1);
      opacity: 1;
    }
  }
</style>
