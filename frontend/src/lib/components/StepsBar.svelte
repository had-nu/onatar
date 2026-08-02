<script lang="ts">
  import { step, setStep } from '$lib/builder.svelte';

  interface Step { label: string; component: any; }
  interface Props { steps: Step[]; }
  let { steps }: Props = $props();

  const stepLabels = ['Classe', 'Background', 'Espécie', 'Atributos', 'Equipamento', 'Revisão'];
</script>

<div class="chrome">
  <div class="chrome-inner">
    <div class="chrome-logo">
      <div class="logo-icon">O</div>
      <div class="logo-text">Onatar</div>
    </div>

    <nav class="steps">
      {#each steps as s, i}
        <button
          class="step"
          class:active={step === i}
          class:done={step > i}
          onclick={() => setStep(i)}
        >
          <span class="step-num">
            {#if step > i}
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"/></svg>
            {:else}
              {i + 1}
            {/if}
          </span>
          <span class="step-label">{s.label}</span>
        </button>
        {#if i < steps.length - 1}
          <span class="step-sep">›</span>
        {/if}
      {/each}
    </nav>

    <div class="chrome-actions">
      <button class="ca-btn" onclick={() => window.location.hash = '#/'}>✕</button>
    </div>
  </div>
</div>

<style>
  .chrome {
    position: sticky;
    top: 0;
    z-index: 100;
    background: linear-gradient(180deg, #0a0a0f 0%, #14141a 100%);
    border-bottom: 3px solid var(--on-red);
    box-shadow: 0 2px 12px rgba(197, 0, 9, 0.15);
  }
  .chrome-inner {
    display: flex;
    align-items: center;
    gap: 16px;
    max-width: 1200px;
    margin: 0 auto;
    padding: 10px 16px;
  }
  .chrome-logo {
    display: flex;
    align-items: center;
    gap: 10px;
    padding-right: 16px;
    border-right: 1px solid var(--on-border);
    flex-shrink: 0;
  }
  .logo-icon {
    width: 32px;
    height: 32px;
    background: linear-gradient(135deg, var(--on-red) 0%, #8b0000 100%);
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: white;
    font-size: 16px;
    font-weight: 700;
  }
  .logo-text { font-size: 16px; font-weight: 500; letter-spacing: 1px; }

  .steps {
    display: flex;
    align-items: center;
    gap: 2px;
    flex: 1;
    overflow-x: auto;
    padding: 2px 0;
  }
  .step {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 14px;
    border-radius: 8px;
    border: 1px solid transparent;
    background: transparent;
    color: var(--on-text-dim);
    cursor: pointer;
    transition: all 0.15s;
    white-space: nowrap;
    font-size: 13px;
    font-weight: 500;
  }
  .step:hover { color: var(--on-text-muted); background: var(--on-bg-hover); }
  .step.active {
    background: var(--on-red);
    color: #fff;
    border-color: var(--on-red);
    box-shadow: 0 2px 8px rgba(197, 0, 9, 0.3);
  }
  .step.done {
    color: var(--on-green);
    background: rgba(0, 184, 122, 0.06);
    border-color: rgba(0, 184, 122, 0.15);
  }
  .step-num {
    width: 20px;
    height: 20px;
    border-radius: 50%;
    background: var(--on-border);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 10px;
    font-weight: 600;
    flex-shrink: 0;
  }
  .step.active .step-num { background: rgba(255,255,255,0.2); }
  .step.done .step-num { background: var(--on-green); color: var(--on-bg-root); }
  .step-label { display: none; }
  @media (min-width: 700px) { .step-label { display: inline; } }

  .step-sep {
    color: var(--on-border-light);
    font-size: 12px;
    padding: 0 2px;
    user-select: none;
  }

  .chrome-actions {
    display: flex;
    gap: 8px;
    flex-shrink: 0;
  }
  .ca-btn {
    width: 32px;
    height: 32px;
    border-radius: 6px;
    border: 1px solid var(--on-border);
    background: transparent;
    color: var(--on-text-dim);
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 14px;
    transition: all 0.15s;
  }
  .ca-btn:hover { background: var(--on-bg-hover); color: var(--on-text); }
</style>