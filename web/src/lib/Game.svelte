<script lang="ts">
  import { ziggyState, mood } from './store';
  import Background from './Background.svelte';
  import Ziggy from './Ziggy.svelte';
  import Stats from './Stats.svelte';
  import Message from './Message.svelte';
  import Controls from './Controls.svelte';
  import DevTools from './DevTools.svelte';
</script>

<div class="game-container">
  <div class="game-wrapper">
    <div class="game-canvas">
      <Background timeOfDay={$ziggyState.timeOfDay} />

      <div class="game-content">
        <div class="top-bar">
          <Stats
            fullness={$ziggyState.fullness}
            happiness={$ziggyState.happiness}
            bond={$ziggyState.bond}
            hp={$ziggyState.hp}
          />
        </div>

        <div class="main-area">
          <div class="message-area">
            <Message message={$ziggyState.message} />
          </div>
        </div>

        <div class="ziggy-area">
          <Ziggy mood={$mood} stage={$ziggyState.stage} />
        </div>
      </div>
    </div>

    <div class="controls-bar">
      <Controls />
    </div>
  </div>

  <DevTools />
</div>

<style>
  .game-container {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 100vh;
    padding: 20px;
    background: #0a0a12;
  }

  .game-wrapper {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .game-canvas {
    position: relative;
    width: 240px;
    height: 240px;
    border-radius: 12px;
    overflow: hidden;
    box-shadow:
      0 0 0 4px rgba(74, 222, 128, 0.2),
      0 0 30px rgba(74, 222, 128, 0.1),
      0 10px 40px rgba(0, 0, 0, 0.5);
  }

  .game-content {
    position: absolute;
    inset: 0;
    display: flex;
    flex-direction: column;
    padding: 8px;
  }

  .top-bar {
    z-index: 10;
  }

  .main-area {
    position: absolute;
    top: 105px;
    left: 50%;
    transform: translateX(-50%);
    z-index: 10;
  }

  .message-area {
    z-index: 10;
  }

  .ziggy-area {
    position: absolute;
    top: 175px;
    left: 50%;
    transform: translateX(-50%);
    z-index: 5;
  }

  .controls-bar {
    display: flex;
    justify-content: center;
  }
</style>
