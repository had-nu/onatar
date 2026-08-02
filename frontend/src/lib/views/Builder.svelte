<script lang="ts">
  import { step, nextStep, prevStep, getCurrentStepValid, undo, getCanUndo, redo, getCanRedo } from '$lib/builder.svelte';
  import StepsBar from '$lib/components/StepsBar.svelte';
  import PreviewSidebar from '$lib/components/PreviewSidebar.svelte';
  import ClassStep from '$lib/builder-steps/ClassStep.svelte';
  import BackgroundStep from '$lib/builder-steps/BackgroundStep.svelte';
  import SpeciesStep from '$lib/builder-steps/SpeciesStep.svelte';
  import AbilitiesStep from '$lib/builder-steps/AbilitiesStep.svelte';
  import EquipmentStep from '$lib/builder-steps/EquipmentStep.svelte';
  import ReviewStep from '$lib/builder-steps/ReviewStep.svelte';
  import Button from '$lib/ui/Button.svelte';

  const steps = [
    { label: 'Class', component: ClassStep },
    { label: 'Background', component: BackgroundStep },
    { label: 'Species', component: SpeciesStep },
    { label: 'Abilities', component: AbilitiesStep },
    { label: 'Equipment', component: EquipmentStep },
    { label: 'Review', component: ReviewStep },
  ];

  const StepComponent = $derived(steps[step].component);
</script>

<div class="builder-page">
  <StepsBar {steps} />

  <div class="builder-layout">
    <main class="builder-content">
      <StepComponent />

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
          >
            Next →
          </Button>
        {/if}
      </div>
    </main>

    <PreviewSidebar />
  </div>
</div>

<style>
  .builder-page {
    min-height: 100vh;
    background: var(--on-bg-root);
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