<script lang="ts">
  import { draft, species, setSpecies, getSpeciesDef } from '$lib/builder.svelte';
  import Card from '$lib/ui/Card.svelte';
  import Tag from '$lib/ui/Tag.svelte';

  const selectedId = $derived(draft.speciesId);
  const selectedVariant = $derived(draft.speciesVariant);
  const speciesDef = $derived(getSpeciesDef());
</script>

<div class="step-species">
  <h2 class="step-title">Escolhe a tua espécie</h2>
  <p class="step-desc">A espécie define a tua herança, traços raciais e bónus de atributos.</p>

  <div class="species-grid">
    {#each species as sp}
      <Card
        variant={selectedId === sp.id ? 'selected' : 'interactive'}
        onclick={() => setSpecies(sp.id, sp.variants[0]?.id)}
        class="species-card"
      >
        <div class="sp-header">
          <span class="sp-name">{sp.name}</span>
          <span class="sp-size">{sp.size} · {sp.speed}ft</span>
        </div>
        <p class="sp-desc">{sp.description}</p>
        <div class="sp-bonuses">
          {#each Object.entries(sp.abilityBonuses) as [ab, val]}
            {#if val && val !== 0}
              <Tag variant="primary">+{val} {ab}</Tag>
            {/if}
          {/each}
        </div>
        <div class="sp-traits">
          {#each sp.traits as t}
            <span class="sp-trait">{t.name}</span>
          {/each}
        </div>
      </Card>
    {/each}
  </div>

  {#if speciesDef?.variants && speciesDef.variants.length > 0}
    <div class="variant-section">
      <h3 class="variant-title">Variante</h3>
      <div class="variant-grid">
        {#each speciesDef.variants as v}
          <Card
            variant={selectedVariant === v.id ? 'selected' : 'interactive'}
            onclick={() => setSpecies(speciesDef.id, v.id)}
            class="variant-card"
          >
            <div class="v-name">{v.name}</div>
            <p class="v-desc">{v.description}</p>
          </Card>
        {/each}
      </div>
    </div>
  {/if}
</div>

<style>
  .step-species { animation: fadeIn 0.3s ease; }
  @keyframes fadeIn { from { opacity: 0; transform: translateY(8px); } to { opacity: 1; transform: translateY(0); } }

  .step-title { font-size: 28px; font-weight: 500; margin: 0 0 8px 0; }
  .step-desc { font-size: 14px; color: var(--on-text-dim); margin: 0 0 24px 0; }

  .species-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
    gap: 16px;
    margin-bottom: 32px;
  }
  :global(.species-card) { padding: 20px; }
  .sp-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px; }
  .sp-name { font-size: 18px; font-weight: 500; color: var(--on-text); }
  .sp-size { font-size: 11px; color: var(--on-text-dim); text-transform: uppercase; letter-spacing: 1px; }
  .sp-desc { font-size: 13px; color: var(--on-text-dim); margin: 0 0 12px 0; line-height: 1.5; }
  .sp-bonuses { display: flex; gap: 6px; flex-wrap: wrap; margin-bottom: 10px; }
  .sp-traits { display: flex; flex-wrap: wrap; gap: 6px; }
  .sp-trait {
    font-size: 11px;
    padding: 4px 10px;
    background: var(--on-bg-root);
    border: 1px solid var(--on-border);
    border-radius: 4px;
    color: var(--on-text-muted);
  }

  .variant-section { margin-top: 8px; }
  .variant-title { font-size: 18px; font-weight: 500; margin: 0 0 16px 0; }
  .variant-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 12px;
  }
  :global(.variant-card) { padding: 16px; }
  .v-name { font-size: 15px; font-weight: 500; color: var(--on-text); margin-bottom: 4px; }
  .v-desc { font-size: 12px; color: var(--on-text-dim); margin: 0; line-height: 1.5; }
</style>