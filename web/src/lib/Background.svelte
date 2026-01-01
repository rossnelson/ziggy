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
  class="absolute inset-0 bg-no-repeat"
  style:background-image="url('/assets/backgrounds.png')"
  style:background-position="-{pos.x}px -{pos.y}px"
  style:background-size="480px 480px"
  style:image-rendering="pixelated"
></div>
