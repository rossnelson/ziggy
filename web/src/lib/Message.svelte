<script lang="ts">
  interface Props {
    message: string;
  }

  let { message }: Props = $props();

  let lines = $derived(message.split('\n').slice(0, 3));
</script>

<div class="relative bg-white/95 rounded-lg py-2 px-3 max-w-[220px] min-h-12 shadow-lg">
  <div class="flex flex-col gap-px font-mono text-[11px] text-[#1a1a2e] leading-tight">
    {#each lines as line}
      <span class="block max-w-[28ch] overflow-hidden text-ellipsis whitespace-nowrap">{line}</span>
    {/each}
  </div>
  <div class="tail"></div>
</div>

<style>
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
