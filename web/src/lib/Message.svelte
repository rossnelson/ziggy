<script lang="ts">
  interface Props {
    message: string;
  }

  let { message }: Props = $props();

  let lines = $derived(message.split('\n').slice(0, 3));
</script>

<div class="message-bubble">
  <div class="message-content">
    {#each lines as line}
      <span class="line">{line}</span>
    {/each}
  </div>
  <div class="tail"></div>
</div>

<style>
  .message-bubble {
    position: relative;
    background: rgba(255, 255, 255, 0.95);
    border-radius: 8px;
    padding: 8px 12px;
    max-width: 220px;
    min-height: 48px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
  }

  .message-content {
    display: flex;
    flex-direction: column;
    gap: 1px;
    font-family: monospace;
    font-size: 11px;
    color: #1a1a2e;
    line-height: 1.3;
  }

  .line {
    display: block;
    max-width: 28ch;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .tail {
    position: absolute;
    bottom: -8px;
    left: 20px;
    width: 0;
    height: 0;
    border-left: 8px solid transparent;
    border-right: 8px solid transparent;
    border-top: 8px solid rgba(255, 255, 255, 0.95);
  }
</style>
