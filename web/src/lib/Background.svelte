<script lang="ts">
  import type { TimeOfDay } from './store';

  interface Props {
    timeOfDay: TimeOfDay;
  }

  let { timeOfDay }: Props = $props();

  const bgPositions: Record<TimeOfDay, { x: number; y: number }> = {
    night: { x: 0, y: 0 },
    dusk: { x: 240, y: 0 },
    day: { x: 0, y: 240 },
    dawn: { x: 240, y: 240 },
  };

  let pos = $derived(bgPositions[timeOfDay]);
</script>

<div
  class="background"
  style:background-image="url('/assets/backgrounds.png')"
  style:background-position="-{pos.x}px -{pos.y}px"
></div>

<style>
  .background {
    position: absolute;
    inset: 0;
    background-size: 480px 480px;
    background-repeat: no-repeat;
    image-rendering: pixelated;
  }
</style>
