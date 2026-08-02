<script lang="ts">
  import { draft, backgrounds, setBackground, getBackgroundDef } from '$lib/builder.svelte';
  import Card from '$lib/ui/Card.svelte';
  import Tag from '$lib/ui/Tag.svelte';

  const selectedId = $derived(draft.backgroundId);
  const backgroundDef = $derived(getBackgroundDef());
</script>

<div class="step-bg">
  <h2 class="step-title">Choose Your Background</h2>
  <p class="step-desc">Your background defines your history, proficiencies, and starting equipment.</p>

  <div class="bg-grid">
    {#each backgrounds as bg}
      <Card
        variant={selectedId === bg.id ? 'selected' : 'interactive'}
        onclick={() => setBackground(bg.id)}
        class="bg-card"
      >
        <div class="bg-name">{bg.name}</div>
        <p class="bg-desc">{bg.description}</p>
        <div class="bg-skills">
          {#each bg.skillProficiencies as sk}
            <Tag variant="green">{sk}</Tag>
          {/each}
          {#if bg.languages[0] > 0}
            <Tag variant="gold">+{bg.languages[0]} language(s)</Tag>
          {/if}
        </div>
        {#if bg.feature}
          <div class="bg-feature">
            <span class="bg-feature-label">Feature:</span>
            <span class="bg-feature-name">{bg.feature.name}</span>
          </div>
        {/if}
      </Card>
    {/each}
  </div>

  {#if backgroundDef}
    <div class="bg-detail">
      <h3>Starting Equipment</h3>
      <ul>
        {#each backgroundDef.equipment as item}
          <li>{item}</li>
        {/each}
      </ul>
      {#if backgroundDef.feature}
        <div class="bg-feature-box">
          <h4>{backgroundDef.feature.name}</h4>
          <p>{backgroundDef.feature.description}</p>
        </div>
      {/if}
    </div>
  {/if}
</div>

<style>
  .step-bg { animation: fadeIn 0.3s ease; }
  @keyframes fadeIn { from { opacity: 0; transform: translateY(8px); } to { opacity: 1; transform: translateY(0); } }

  .step-title { font-size: 28px; font-weight: 500; margin: 0 0 8px 0; }
  .step-desc { font-size: 14px; color: var(--on-text-dim); margin: 0 0 24px 0; }

  .bg-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
    gap: 16px;
    margin-bottom: 32px;
  }
  :global(.bg-card) { padding: 20px; }
  .bg-name { font-size: 18px; font-weight: 500; color: var(--on-text); margin-bottom: 6px; }
  .bg-desc { font-size: 13px; color: var(--on-text-dim); margin: 0 0 12px 0; line-height: 1.5; }
  .bg-skills { display: flex; gap: 6px; flex-wrap: wrap; margin-bottom: 10px; }
  .bg-feature {
    font-size: 12px;
    padding: 8px 10px;
    background: var(--on-bg-root);
    border-radius: 6px;
    border: 1px solid var(--on-border);
  }
  .bg-feature-label { color: var(--on-text-dim); }
  .bg-feature-name { color: var(--on-gold); font-weight: 500; }

  .bg-detail {
    padding: 20px;
    background: linear-gradient(135deg, var(--on-bg-surface) 0%, var(--on-bg-elevated) 100%);
    border: 1px solid var(--on-border);
    border-radius: 12px;
  }
  .bg-detail h3 { font-size: 14px; text-transform: uppercase; letter-spacing: 1px; color: var(--on-text-dim); margin: 0 0 12px 0; }
  .bg-detail ul { margin: 0 0 16px 0; padding-left: 18px; color: var(--on-text-muted); font-size: 13px; line-height: 1.8; }
  .bg-feature-box {
    padding: 14px;
    background: var(--on-bg-root);
    border-radius: 8px;
    border-left: 3px solid var(--on-gold);
  }
  .bg-feature-box h4 { font-size: 14px; color: var(--on-gold); margin: 0 0 4px 0; }
  .bg-feature-box p { font-size: 13px; color: var(--on-text-muted); margin: 0; line-height: 1.5; }
</style>