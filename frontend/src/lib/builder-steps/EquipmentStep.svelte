<script lang="ts">
  import { draft, getBackgroundDef } from '$lib/builder.svelte';
  import Card from '$lib/ui/Card.svelte';

  // Simplified equipment — auto-assigned from background for MVP
  const backgroundDef = $derived(getBackgroundDef());
  const equipment = $derived(backgroundDef?.equipment ?? []);
</script>

<div class="step-equip">
  <h2 class="step-title">Equipamento</h2>
  <p class="step-desc">O teu background fornece o equipamento inicial. No futuro, poderás personalizar armas, armaduras e itens.</p>

  <div class="equip-grid">
    {#each equipment as item, i}
      <Card variant="default" class="equip-card">
        <div class="equip-num">{String(i + 1).padStart(2, '0')}</div>
        <div class="equip-name">{item}</div>
      </Card>
    {/each}
  </div>

  <div class="equip-note">
    <p>💡 <strong>Nota:</strong> A gestão detalhada de equipamento (armas, armaduras, consumíveis) será implementada numa versão futura. Por agora, o equipamento base é herdado do background.</p>
  </div>
</div>

<style>
  .step-equip { animation: fadeIn 0.3s ease; }
  @keyframes fadeIn { from { opacity: 0; transform: translateY(8px); } to { opacity: 1; transform: translateY(0); } }

  .step-title { font-size: 28px; font-weight: 500; margin: 0 0 8px 0; }
  .step-desc { font-size: 14px; color: var(--on-text-dim); margin: 0 0 24px 0; }

  .equip-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 12px;
    margin-bottom: 24px;
  }
  :global(.equip-card) {
    padding: 16px;
    display: flex;
    align-items: center;
    gap: 12px;
  }
  .equip-num {
    font-size: 20px;
    font-weight: 600;
    color: var(--on-text-dim);
    font-variant-numeric: tabular-nums;
    min-width: 28px;
  }
  .equip-name { font-size: 14px; color: var(--on-text-muted); }

  .equip-note {
    padding: 16px;
    background: rgba(201, 169, 78, 0.05);
    border: 1px solid rgba(201, 169, 78, 0.15);
    border-radius: 10px;
    font-size: 13px;
    color: var(--on-text-muted);
    line-height: 1.6;
  }
  .equip-note strong { color: var(--on-gold); }
</style>