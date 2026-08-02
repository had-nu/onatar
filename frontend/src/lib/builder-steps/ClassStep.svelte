<script lang="ts">
  import { draft, classes, selectClass, setSubclass, getClassDef } from '$lib/builder.svelte';
  import Card from '$lib/ui/Card.svelte';
  import Tag from '$lib/ui/Tag.svelte';

  const selectedClassId = $derived(draft.classes[0]?.id);
  const selectedSubclassId = $derived(draft.classes[0]?.subclassId);
  const classDef = $derived(getClassDef());
</script>

<div class="step-class">
  <h2 class="step-title">Choose Your Class</h2>
  <p class="step-desc">Your class defines your role in combat, your abilities, and your playstyle.</p>

  <div class="class-grid">
    {#each classes as cls}
      <Card
        variant={selectedClassId === cls.id ? 'selected' : 'interactive'}
        onclick={() => selectClass(cls.id)}
        class="class-card"
      >
        <div class="cc-header">
          <span class="cc-name">{cls.name}</span>
          <span class="cc-die">d{cls.hitDie}</span>
        </div>
        <p class="cc-desc">{cls.spellcaster ? 'Spellcaster' : 'Martial'} · Saving Throws: {cls.savingThrows.join(', ')}</p>
        <div class="cc-tags">
          {#each cls.primaryAbility as ab}
            <Tag variant="primary">{ab}</Tag>
          {/each}
          {#each cls.savingThrows as st}
            <Tag>{st}</Tag>
          {/each}
        </div>
        {#if cls.features}
          <div class="cc-features">
            {#each cls.features.filter(f => f.level === 1) as feat}
              <span class="cc-feat">{feat.name}</span>
            {/each}
          </div>
        {/if}
      </Card>
    {/each}
  </div>

  {#if selectedClassId && classDef?.subClasses && classDef.subClasses.length > 0}
    <div class="subclass-section">
      <h3 class="subclass-title">Subclass <span class="subclass-optional">(level {classDef.subclassLevel})</span></h3>
      <div class="subclass-grid">
        {#each classDef.subClasses as sc}
          <Card
            variant={selectedSubclassId === sc.id ? 'selected' : 'interactive'}
            onclick={() => setSubclass(selectedClassId, sc.id)}
            class="subclass-card"
          >
            <div class="sc-name">{sc.name}</div>
            <p class="sc-desc">{sc.description}</p>
          </Card>
        {/each}
      </div>
    </div>
  {/if}
</div>

<style>
  .step-class { animation: fadeIn 0.3s ease; }
  @keyframes fadeIn { from { opacity: 0; transform: translateY(8px); } to { opacity: 1; transform: translateY(0); } }

  .step-title { font-size: 28px; font-weight: 500; margin: 0 0 8px 0; letter-spacing: -0.3px; }
  .step-desc { font-size: 14px; color: var(--on-text-dim); margin: 0 0 24px 0; }

  .class-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
    gap: 16px;
    margin-bottom: 32px;
  }
  :global(.class-card) { padding: 20px; }
  .cc-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 8px;
  }
  .cc-name { font-size: 18px; font-weight: 500; color: var(--on-text); }
  .cc-die { font-size: 12px; color: var(--on-text-dim); font-family: var(--font-mono); }
  .cc-desc { font-size: 12px; color: var(--on-text-dim); margin: 0 0 12px 0; line-height: 1.5; }
  .cc-tags { display: flex; gap: 6px; flex-wrap: wrap; margin-bottom: 12px; }
  .cc-features { display: flex; flex-wrap: wrap; gap: 6px; }
  .cc-feat {
    font-size: 11px;
    padding: 4px 10px;
    background: var(--on-bg-root);
    border: 1px solid var(--on-border);
    border-radius: 4px;
    color: var(--on-text-muted);
  }

  .subclass-section { margin-top: 8px; }
  .subclass-title { font-size: 18px; font-weight: 500; margin: 0 0 16px 0; }
  .subclass-optional { font-size: 13px; color: var(--on-text-dim); font-weight: 400; }
  .subclass-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
    gap: 12px;
  }
  :global(.subclass-card) { padding: 16px; }
  .sc-name { font-size: 15px; font-weight: 500; color: var(--on-text); margin-bottom: 4px; }
  .sc-desc { font-size: 12px; color: var(--on-text-dim); margin: 0; line-height: 1.5; }
</style>