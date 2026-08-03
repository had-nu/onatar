<script lang="ts">
  import { onMount } from 'svelte';
  import { step, nextStep, prevStep, getCurrentStepValid, undo, getCanUndo, redo, getCanRedo, classes } from '$lib/builder.svelte';
  import StepsBar from '$lib/components/StepsBar.svelte';
  import PreviewSidebar from '$lib/components/PreviewSidebar.svelte';
  import ClassStep from '$lib/builder-steps/ClassStep.svelte';
  import BackgroundStep from '$lib/builder-steps/BackgroundStep.svelte';
  import SpeciesStep from '$lib/builder-steps/SpeciesStep.svelte';
  import AbilitiesStep from '$lib/builder-steps/AbilitiesStep.svelte';
  import EquipmentStep from '$lib/builder-steps/EquipmentStep.svelte';
  import ReviewStep from '$lib/builder-steps/ReviewStep.svelte';
  import Button from '$lib/ui/Button.svelte';

  const stepComponents = [
    ClassStep,
    BackgroundStep,
    SpeciesStep,
    AbilitiesStep,
    EquipmentStep,
    ReviewStep,
  ];

  let StepComponent = stepComponents[step.value];

  $effect(() => {
    StepComponent = stepComponents[step.value];
  });

  const steps = [
    { label: 'Classe', component: ClassStep },
    { label: 'Background', component: BackgroundStep },
    { label: 'Espécie', component: SpeciesStep },
    { label: 'Atributos', component: AbilitiesStep },
    { label: 'Equipamento', component: EquipmentStep },
    { label: 'Revisão', component: ReviewStep },
  ];

  // Expose step info for debugging (runs on client side)
  onMount(() => {
    if (typeof window !== 'undefined') {
      (window as any).__builder_step_value = step.value;
      (window as any).__builder_step_component = 'ClassStep'; // Default step is 0 = ClassStep
      // Update when step changes
      const updateDebug = () => {
        (window as any).__builder_step_value = step.value;
        const stepNames = ['ClassStep', 'BackgroundStep', 'SpeciesStep', 'AbilitiesStep', 'EquipmentStep', 'ReviewStep'];
        (window as any).__builder_step_component = stepNames[step.value] || 'unknown';
      };
      $effect(() => {
        updateDebug();
      });
    }
  });
</script>

<div class="builder-page">
  {#if classes.value.length === 0}
    <div class="builder-loading">
      <div class="spinner"></div>
      <p>Loading content…</p>
    </div>
  {:else}
    <StepsBar {steps} />

    <div class="builder-layout">
      <main class="builder-content">
        <svelte:component this={StepComponent} />

        <div class="builder-nav">
          <Button
            variant="ghost"
            onclick={prevStep}
            disabled={step === 0}
          >
            ← Previous
          </Button>

          <div class="builder-nav-meta">
            {#if getCanUndo()}
              <button class="nav-undo" onclick={undo} title="Undo">↩</button>
            {/if}
            {#if getCanRedo()}
              <button class="nav-redo" onclick={redo} title="Redo">↪</button>
            {/if}
          </div>

          {#if step < 5}
            <Button
              variant="primary"
              onclick={nextStep}
              disabled={!getCurrentStepValid()}
              data-testid="next-btn"
            >
              Next →
            </Button>
          {/if}
        </div>
      </main>

      <PreviewSidebar />
    </div>
  {/if}
</div>

<style>
  .builder-page {
    min-height: 100vh;
    background: var(--on-bg-root);
  }
  .builder-loading {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    min-height: 50vh;
    gap: 1rem;
  }
  .spinner {
    width: 2.5rem;
    height: 2.5rem;
    border: 3px solid var(--border);
    border-top-color: var(--accent);
    border-radius: 50%;
    animation: spin 1s linear infinite;
  }
  @keyframes spin {
    to { transform: rotate(360deg); }
  }
  .builder-layout {
    display: grid;
    grid-template-columns: 1fr 300px;
    gap: 24px;
    max-width: 1200px;
    margin: 0 auto;
    padding: 24px 16px;
    align-items: start;
  }
  .builder-content {
    min-width: 0;
  }
  .builder-nav {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: 32px;
    padding-top: 20px;
    border-top: 1px solid var(--on-border);
  }
  .builder-nav-meta {
    display: flex;
    gap: 8px;
  }
  .nav-undo, .nav-redo {
    width: 32px;
    height: 32px;
    border-radius: 6px;
    border: 1px solid var(--on-border);
    background: transparent;
    color: var(--on-text-muted);
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 14px;
    transition: all 0.15s;
  }
  .nav-undo:hover, .nav-redo:hover {
    background: var(--on-bg-hover);
    color: var(--on-text);
  }

  @media (max-width: 900px) {
    .builder-layout {
      grid-template-columns: 1fr;
    }
  }
</style>
