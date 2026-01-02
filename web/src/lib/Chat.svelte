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

  // Subscribe to chat stores
  let messages = $derived($chatMessages);
  let isLoading = $derived($chatLoading);
  let mystery = $derived($mysteryStatus);

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

  async function handleStartMystery(mystery: Mystery) {
    showMysteries = false;
    chatLoading.set(true);

    await startMystery(mystery.id, track);
    const message = track === 'educational'
      ? `let's learn about "${mystery.title}"`
      : `let's solve the mystery "${mystery.title}"`;
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

  function scrollToBottom() {
    if (messagesContainer) {
      setTimeout(() => {
        messagesContainer.scrollTop = messagesContainer.scrollHeight;
      }, 50);
    }
  }

  // Scroll to bottom when messages change
  $effect(() => {
    if (messages.length > 0) {
      scrollToBottom();
    }
  });

  onMount(() => {
    loadMysteries();
  });
</script>

<div class="hidden sm:flex w-[280px] h-60 bg-[rgba(26,26,46,0.95)] border-2 border-green-400/30 rounded-lg flex-col overflow-hidden relative">
  <div class="flex justify-between items-center px-3 py-2 border-b border-green-400/20 text-green-400 font-mono text-[11px] font-bold shrink-0">
    <span>Chat with Ziggy</span>
    <div class="flex gap-1">
      <button
        class="bg-green-400/10 border border-green-400/30 rounded px-2 py-1 text-green-400 font-mono text-[9px] cursor-pointer hover:bg-green-400/20"
        onclick={toggleTrack}
        title={track === 'fun' ? 'Switch to educational mode' : 'Switch to fun mode'}
      >
        {track === 'fun' ? '🎮' : '📚'}
      </button>
      <button
        class="border rounded px-2 py-1 text-green-400 font-mono text-[9px] cursor-pointer transition-colors {showMysteries ? 'bg-green-400/30 border-green-400/60' : 'bg-green-400/10 border-green-400/30 hover:bg-green-400/20'}"
        onclick={() => (showMysteries = !showMysteries)}
      >
        {track === 'educational' ? 'Topics' : 'Mysteries'}
      </button>
    </div>
  </div>

  {#if mystery?.active}
    <div class="flex items-center gap-2 px-3 py-2 bg-purple-500/10 border-b border-purple-500/20 shrink-0">
      <span class="font-mono text-[9px] font-bold text-purple-500 whitespace-nowrap">{mystery.mystery?.title}</span>
      <div class="flex-1 h-1.5 bg-black/30 rounded-sm overflow-hidden">
        <div
          class="h-full bg-gradient-to-r from-purple-500 to-green-400 rounded-sm transition-all duration-300"
          style="width: {(mystery.progress / mystery.totalHints) * 100}%"
        ></div>
      </div>
      <span class="font-mono text-[8px] text-[#a0a0b0] whitespace-nowrap">{mystery.progress}/{mystery.totalHints} hints</span>
    </div>
  {:else if showMysteries}
    <div class="absolute top-[30px] left-3 right-3 max-h-[200px] overflow-y-auto bg-[#1a1a2c] border border-green-400/20 rounded z-10">
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

  <div class="flex-1 overflow-y-auto p-2 flex flex-col gap-1.5" bind:this={messagesContainer}>
    {#each messages as message}
      <div class="flex flex-col gap-0.5 px-2 py-1.5 rounded-md font-mono text-[10px] max-w-[85%] {message.role === 'user' ? 'bg-green-400/15 self-end' : 'bg-purple-600/15 self-start'}">
        <span class="text-[8px] font-bold uppercase {message.role === 'user' ? 'text-green-400' : 'text-purple-500'}">
          {message.role === 'user' ? 'You' : 'Ziggy'}
        </span>
        <span class="text-[#e0e0e0] break-words whitespace-pre-wrap">{message.content}</span>
      </div>
    {/each}
    {#if messages.length === 0 && !isLoading}
      <div class="text-center text-[#a0a0b0] font-mono text-[10px] py-5">Say hi to Ziggy!</div>
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

<style>
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
