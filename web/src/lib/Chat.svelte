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
  let messagesContainer: HTMLDivElement;

  // Subscribe to chat stores
  let messages = $derived($chatMessages);
  let isLoading = $derived($chatLoading);
  let mystery = $derived($mysteryStatus);

  async function loadMysteries() {
    const result = await getAvailableMysteries('fun');
    if (result.success && result.data) {
      mysteries = result.data;
    }
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

    await startMystery(mystery.id);
    await sendChatMessage(`let's solve the mystery "${mystery.title}"`);

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

<div class="chat-panel">
  <div class="chat-header">
    <span>Chat with Ziggy</span>
    <button class="mystery-btn" onclick={() => (showMysteries = !showMysteries)}>
      {showMysteries ? 'Hide' : 'Mysteries'}
    </button>
  </div>

  {#if mystery?.active}
    <div class="mystery-progress">
      <span class="mystery-label">{mystery.mystery?.title}</span>
      <div class="progress-bar">
        <div
          class="progress-fill"
          style="width: {(mystery.progress / mystery.totalHints) * 100}%"
        ></div>
      </div>
      <span class="progress-text">{mystery.progress}/{mystery.totalHints} hints</span>
    </div>
  {:else if showMysteries}
    <div class="mysteries-list">
      {#each mysteries as m}
        <button class="mystery-item" onclick={() => handleStartMystery(m)}>
          <span class="mystery-title">{m.title}</span>
          <span class="mystery-desc">{m.description}</span>
        </button>
      {/each}
      {#if mysteries.length === 0}
        <div class="no-mysteries">No mysteries available</div>
      {/if}
    </div>
  {/if}

  <div class="messages" bind:this={messagesContainer}>
    {#each messages as message}
      <div class="message {message.role}">
        <span class="role">{message.role === 'user' ? 'You' : 'Ziggy'}</span>
        <span class="content">{message.content}</span>
      </div>
    {/each}
    {#if messages.length === 0 && !isLoading}
      <div class="empty-chat">Say hi to Ziggy!</div>
    {/if}
    {#if isLoading}
      <div class="message ziggy loading">
        <span class="role">Ziggy</span>
        <span class="content typing">
          <span class="dot"></span>
          <span class="dot"></span>
          <span class="dot"></span>
        </span>
      </div>
    {/if}
  </div>

  <div class="input-area">
    <input
      type="text"
      placeholder="Type a message..."
      bind:value={inputValue}
      onkeydown={handleKeydown}
      disabled={isLoading}
    />
    <button class="send-btn" onclick={handleSend} disabled={isLoading || !inputValue.trim()}>
      {isLoading ? '...' : 'Send'}
    </button>
  </div>
</div>

<style>
  .chat-panel {
    width: 280px;
    height: 240px;
    background: rgba(26, 26, 46, 0.95);
    border: 2px solid rgba(74, 222, 128, 0.3);
    border-radius: 8px;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    position: relative;
  }

  .chat-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 8px 12px;
    border-bottom: 1px solid rgba(74, 222, 128, 0.2);
    color: #4ade80;
    font-family: monospace;
    font-size: 11px;
    font-weight: bold;
    flex-shrink: 0;
  }

  .mystery-btn {
    background: rgba(74, 222, 128, 0.1);
    border: 1px solid rgba(74, 222, 128, 0.3);
    border-radius: 4px;
    color: #4ade80;
    font-family: monospace;
    font-size: 9px;
    padding: 4px 8px;
    cursor: pointer;
  }

  .mystery-btn:hover {
    background: rgba(74, 222, 128, 0.2);
  }

  .mysteries-list {
    position: absolute;
    top: 30px;
    left: 12px;
    right: 12px;
    max-height: 200px;
    overflow-y: auto;
    background: #1a1a2c;
    border: 1px solid rgba(74, 222, 128, 0.2);
    border-radius: 4px;
    z-index: 10;
  }

  .mystery-item {
    display: flex;
    flex-direction: column;
    width: 100%;
    padding: 8px 12px;
    background: transparent;
    border: none;
    border-bottom: 1px solid rgba(74, 222, 128, 0.1);
    color: #d0d0e0;
    font-family: monospace;
    text-align: left;
    cursor: pointer;
  }

  .mystery-item:hover {
    background: rgba(74, 222, 128, 0.1);
  }

  .mystery-title {
    font-size: 10px;
    font-weight: bold;
    color: #4ade80;
  }

  .mystery-desc {
    font-size: 9px;
    color: #a0a0b0;
    margin-top: 2px;
  }

  .no-mysteries {
    padding: 12px;
    text-align: center;
    color: #a0a0b0;
    font-size: 10px;
    font-family: monospace;
  }

  .mystery-progress {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 12px;
    background: rgba(168, 85, 247, 0.1);
    border-bottom: 1px solid rgba(168, 85, 247, 0.2);
    flex-shrink: 0;
  }

  .mystery-label {
    font-family: monospace;
    font-size: 9px;
    font-weight: bold;
    color: #a855f7;
    white-space: nowrap;
  }

  .progress-bar {
    flex: 1;
    height: 6px;
    background: rgba(0, 0, 0, 0.3);
    border-radius: 3px;
    overflow: hidden;
  }

  .progress-fill {
    height: 100%;
    background: linear-gradient(90deg, #a855f7, #4ade80);
    border-radius: 3px;
    transition: width 0.3s ease;
  }

  .progress-text {
    font-family: monospace;
    font-size: 8px;
    color: #a0a0b0;
    white-space: nowrap;
  }

  .messages {
    flex: 1;
    overflow-y: auto;
    padding: 8px;
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  .message {
    display: flex;
    flex-direction: column;
    gap: 2px;
    padding: 6px 8px;
    border-radius: 6px;
    font-family: monospace;
    font-size: 10px;
  }

  .message.user {
    background: rgba(74, 222, 128, 0.15);
    align-self: flex-end;
    max-width: 85%;
  }

  .message.ziggy {
    background: rgba(147, 51, 234, 0.15);
    align-self: flex-start;
    max-width: 85%;
  }

  .message .role {
    font-size: 8px;
    font-weight: bold;
    color: #a0a0b0;
    text-transform: uppercase;
  }

  .message.user .role {
    color: #4ade80;
  }

  .message.ziggy .role {
    color: #a855f7;
  }

  .message .content {
    color: #e0e0e0;
    word-wrap: break-word;
    white-space: pre-wrap;
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

  .dot:nth-child(1) {
    animation-delay: -0.32s;
  }

  .dot:nth-child(2) {
    animation-delay: -0.16s;
  }

  @keyframes bounce {
    0%,
    80%,
    100% {
      transform: scale(0);
      opacity: 0.5;
    }
    40% {
      transform: scale(1);
      opacity: 1;
    }
  }

  .empty-chat {
    text-align: center;
    color: #a0a0b0;
    font-family: monospace;
    font-size: 10px;
    padding: 20px;
  }

  .input-area {
    display: flex;
    gap: 6px;
    padding: 8px;
    border-top: 1px solid rgba(74, 222, 128, 0.2);
    flex-shrink: 0;
  }

  .input-area input {
    flex: 1;
    padding: 6px 8px;
    background: rgba(0, 0, 0, 0.3);
    border: 1px solid rgba(74, 222, 128, 0.3);
    border-radius: 4px;
    color: #e0e0e0;
    font-family: monospace;
    font-size: 10px;
  }

  .input-area input:focus {
    outline: none;
    border-color: rgba(74, 222, 128, 0.6);
  }

  .input-area input::placeholder {
    color: #606080;
  }

  .input-area input:disabled {
    opacity: 0.6;
  }

  .send-btn {
    padding: 6px 12px;
    background: rgba(74, 222, 128, 0.2);
    border: 1px solid rgba(74, 222, 128, 0.4);
    border-radius: 4px;
    color: #4ade80;
    font-family: monospace;
    font-size: 10px;
    font-weight: bold;
    cursor: pointer;
    transition: all 0.2s;
    min-width: 50px;
  }

  .send-btn:hover:not(:disabled) {
    background: rgba(74, 222, 128, 0.3);
    border-color: rgba(74, 222, 128, 0.6);
  }

  .send-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
</style>
