<script lang="ts">
  import { onMount } from 'svelte'
  import { navigate } from '../router.svelte'
  import { contentError, loadContent } from '../content.svelte'
  import {
    STEPS,
    builder,
    canGoNext,
    nextStep,
    prevStep,
    redo,
    setStep,
    step,
    undo,
  } from '../builder.svelte'
  import ClassStep from './builder/ClassStep.svelte'
  import BackgroundStep from './builder/BackgroundStep.svelte'
  import SpeciesStep from './builder/SpeciesStep.svelte'
  import AbilitiesStep from './builder/AbilitiesStep.svelte'
  import EquipmentStep from './builder/EquipmentStep.svelte'
  import ReviewStep from './builder/ReviewStep.svelte'
  import BuilderPreview from './builder/BuilderPreview.svelte'

  let status = $state<'loading' | 'ready' | 'error'>('loading')

  onMount(async () => {
    try {
      await loadContent()
      status = 'ready'
    } catch {
      status = 'error'
    }
  })

  function retry() {
    status = 'loading'
    void loadContent(true)
      .then(() => (status = 'ready'))
      .catch(() => (status = 'error'))
  }

  const current = $derived(step())
</script>

<div class="page-head">
  <h1>Constructor de personagem</h1>
  <div class="head-actions">
    <button class="btn" onclick={undo} disabled={builder.value.history.length === 0}
      >Desfazer</button
    >
    <button class="btn" onclick={redo} disabled={builder.value.future.length === 0}>Refazer</button>
  </div>
</div>

{#if status === 'loading'}
  <p>Carregar conteúdo…</p>
{:else if status === 'error'}
  <div class="error-box">
    <p>
      Não foi possível carregar o conteúdo{contentError.value ? `: ${contentError.value}` : ''}.
    </p>
    <button class="btn" onclick={retry}>Tentar de novo</button>
  </div>
{:else}
  <nav class="stepper" aria-label="Passos do wizard">
    {#each STEPS as st, i (st.id)}
      <button
        class:active={i === builder.value.stepIndex}
        class:done={i < builder.value.stepIndex}
        onclick={() => setStep(i)}
      >
        <span class="idx">{i + 1}</span>
        {st.label}
      </button>
    {/each}
  </nav>

  <div class="builder-layout">
    <main class="step-panel">
      {#if current.id === 'class'}
        <ClassStep />
      {:else if current.id === 'background'}
        <BackgroundStep />
      {:else if current.id === 'species'}
        <SpeciesStep />
      {:else if current.id === 'abilities'}
        <AbilitiesStep />
      {:else if current.id === 'equipment'}
        <EquipmentStep />
      {:else}
        <ReviewStep />
      {/if}

      <div class="nav-actions">
        <button class="btn" onclick={prevStep} disabled={builder.value.stepIndex === 0}>
          ← Anterior
        </button>
        {#if current.id !== 'review'}
          <button class="btn primary" onclick={nextStep} disabled={!canGoNext()}>
            Continuar →
          </button>
        {/if}
      </div>
    </main>
    <aside class="preview">
      <BuilderPreview />
    </aside>
  </div>
{/if}

<style>
  .page-head {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
    margin-bottom: 1.25rem;
  }
  h1 {
    margin: 0;
    color: var(--text-h);
  }
  .head-actions {
    display: flex;
    gap: 0.5rem;
  }
  .stepper {
    display: flex;
    flex-wrap: wrap;
    gap: 0.4rem;
    margin-bottom: 1.5rem;
  }
  .stepper button {
    font: inherit;
    font-size: 0.85rem;
    display: inline-flex;
    align-items: center;
    gap: 0.4rem;
    padding: 0.35rem 0.85rem;
    border-radius: 999px;
    border: 1px solid var(--border);
    background: var(--code-bg);
    color: var(--text);
    cursor: pointer;
  }
  .stepper .idx {
    font-size: 0.7rem;
    font-weight: 700;
    color: var(--text-h);
  }
  .stepper button.active {
    color: var(--accent);
    background: var(--accent-bg);
    border-color: var(--accent-border);
  }
  .stepper button.done {
    color: var(--text-h);
  }
  .stepper button:disabled {
    cursor: default;
  }
  .builder-layout {
    display: grid;
    grid-template-columns: minmax(0, 1fr) 22rem;
    gap: 1.5rem;
    align-items: start;
  }
  .step-panel {
    display: grid;
    gap: 1rem;
  }
  .preview {
    position: sticky;
    top: 1rem;
  }
  .nav-actions {
    display: flex;
    justify-content: space-between;
    margin-top: 1.5rem;
  }
  .error-box {
    border: 1px solid var(--danger-border);
    background: var(--danger-bg);
    border-radius: 8px;
    padding: 1.25rem;
    display: grid;
    gap: 0.75rem;
    justify-items: start;
  }
  @media (max-width: 1024px) {
    .builder-layout {
      grid-template-columns: 1fr;
    }
    .preview {
      position: static;
    }
  }
</style>
