<script lang="ts">
  import {
    ziggyState,
    setStat,
    setTimeOfDay,
    setStage,
    setSleeping,
    startDecay,
    stopDecay,
    resetState,
    type TimeOfDay,
    type Stage,
  } from './store';

  let isOpen = $state(false);
  let decayActive = $state(false);

  const timeOptions: TimeOfDay[] = ['night', 'dawn', 'day', 'dusk'];
  const stageOptions: Stage[] = ['egg', 'baby', 'teen', 'adult', 'elder'];

  function toggleDecay() {
    if (decayActive) {
      stopDecay();
    } else {
      startDecay();
    }
    decayActive = !decayActive;
  }

  function handleReset() {
    stopDecay();
    decayActive = false;
    resetState();
  }
</script>

<div class="devtools" class:open={isOpen}>
  <button class="toggle" onclick={() => isOpen = !isOpen}>
    {isOpen ? '✕' : '⚙'}
  </button>

  {#if isOpen}
    <div class="panel">
      <h3>Dev Tools</h3>

      <div class="section">
        <label>
          Fullness: {$ziggyState.fullness}
          <input
            type="range"
            min="0"
            max="100"
            value={$ziggyState.fullness}
            oninput={(e) => setStat('fullness', Number(e.currentTarget.value))}
          />
        </label>

        <label>
          Happiness: {$ziggyState.happiness}
          <input
            type="range"
            min="0"
            max="100"
            value={$ziggyState.happiness}
            oninput={(e) => setStat('happiness', Number(e.currentTarget.value))}
          />
        </label>

        <label>
          Bond: {$ziggyState.bond}
          <input
            type="range"
            min="0"
            max="100"
            value={$ziggyState.bond}
            oninput={(e) => setStat('bond', Number(e.currentTarget.value))}
          />
        </label>

        <label>
          HP: {$ziggyState.hp}
          <input
            type="range"
            min="0"
            max="100"
            value={$ziggyState.hp}
            oninput={(e) => setStat('hp', Number(e.currentTarget.value))}
          />
        </label>
      </div>

      <div class="section">
        <label>
          Time of Day:
          <select
            value={$ziggyState.timeOfDay}
            onchange={(e) => setTimeOfDay(e.currentTarget.value as TimeOfDay)}
          >
            {#each timeOptions as opt}
              <option value={opt}>{opt}</option>
            {/each}
          </select>
        </label>

        <label>
          Stage:
          <select
            value={$ziggyState.stage}
            onchange={(e) => setStage(e.currentTarget.value as Stage)}
          >
            {#each stageOptions as opt}
              <option value={opt}>{opt}</option>
            {/each}
          </select>
        </label>

        <label class="checkbox">
          <input
            type="checkbox"
            checked={$ziggyState.sleeping}
            onchange={(e) => setSleeping(e.currentTarget.checked)}
          />
          Sleeping
        </label>
      </div>

      <div class="section buttons">
        <button class:active={decayActive} onclick={toggleDecay}>
          {decayActive ? 'Stop Decay' : 'Start Decay'}
        </button>
        <button onclick={handleReset}>Reset</button>
      </div>

      <div class="info">
        <div>Age: {Math.floor($ziggyState.age)}s</div>
        <div>Gen: {$ziggyState.generation}</div>
      </div>
    </div>
  {/if}
</div>

<style>
  .devtools {
    position: fixed;
    top: 10px;
    right: 10px;
    z-index: 1000;
  }

  .toggle {
    width: 32px;
    height: 32px;
    background: rgba(26, 26, 46, 0.9);
    border: 2px solid rgba(74, 222, 128, 0.3);
    border-radius: 6px;
    color: #4ade80;
    font-size: 14px;
    cursor: pointer;
    transition: all 0.2s;
  }

  .toggle:hover {
    border-color: rgba(74, 222, 128, 0.6);
  }

  .panel {
    position: absolute;
    top: 40px;
    right: 0;
    width: 220px;
    background: rgba(26, 26, 46, 0.95);
    border: 2px solid rgba(74, 222, 128, 0.3);
    border-radius: 8px;
    padding: 12px;
    font-family: monospace;
    font-size: 11px;
    color: #d0d0e0;
  }

  h3 {
    margin: 0 0 10px;
    font-size: 12px;
    color: #4ade80;
    text-transform: uppercase;
  }

  .section {
    display: flex;
    flex-direction: column;
    gap: 8px;
    margin-bottom: 12px;
    padding-bottom: 12px;
    border-bottom: 1px solid rgba(74, 222, 128, 0.2);
  }

  .section:last-child {
    border-bottom: none;
    margin-bottom: 0;
    padding-bottom: 0;
  }

  label {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  label.checkbox {
    flex-direction: row;
    align-items: center;
  }

  input[type="range"] {
    width: 100%;
    height: 6px;
    -webkit-appearance: none;
    appearance: none;
    background: rgba(255, 255, 255, 0.1);
    border-radius: 3px;
    outline: none;
  }

  input[type="range"]::-webkit-slider-thumb {
    -webkit-appearance: none;
    width: 14px;
    height: 14px;
    background: #4ade80;
    border-radius: 50%;
    cursor: pointer;
  }

  select {
    background: rgba(255, 255, 255, 0.1);
    border: 1px solid rgba(74, 222, 128, 0.3);
    border-radius: 4px;
    padding: 4px;
    color: #d0d0e0;
    font-family: monospace;
    font-size: 11px;
  }

  input[type="checkbox"] {
    accent-color: #4ade80;
  }

  .buttons {
    flex-direction: row;
    flex-wrap: wrap;
    gap: 6px;
  }

  .buttons button {
    flex: 1;
    padding: 6px 10px;
    background: rgba(74, 222, 128, 0.1);
    border: 1px solid rgba(74, 222, 128, 0.3);
    border-radius: 4px;
    color: #d0d0e0;
    font-family: monospace;
    font-size: 10px;
    cursor: pointer;
    transition: all 0.2s;
  }

  .buttons button:hover {
    background: rgba(74, 222, 128, 0.2);
    border-color: rgba(74, 222, 128, 0.5);
  }

  .buttons button.active {
    background: rgba(239, 68, 68, 0.2);
    border-color: rgba(239, 68, 68, 0.5);
    color: #ef4444;
  }

  .info {
    display: flex;
    justify-content: space-between;
    color: #a0a0b0;
    font-size: 10px;
  }
</style>
