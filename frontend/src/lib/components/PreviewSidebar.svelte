<script lang="ts">
  import { draft, preview, isLoading, pendingChoices, getClassDef, getSpeciesDef, getBackgroundDef, getTotalLevel } from '$lib/builder.svelte';

  function getMod(score: number): string {
    const mod = Math.floor((score - 10) / 2);
    return mod >= 0 ? `+${mod}` : `${mod}`;
  }
</script>

<aside class="preview-sidebar">
  <div class="ps-header">
    <span class="ps-dot" class:ps-dot-loading={isLoading}></span>
    <h3>Pré-visualização</h3>
  </div>

  <div class="ps-body">
    {#if !preview}
      <div class="ps-empty">
        {#if isLoading}
          <div class="ps-spinner"></div>
          <p>A calcular...</p>
        {:else}
          <p>Completa os passos anteriores para ver a pré-visualização.</p>
        {/if}
      </div>
    {:else}
      <div class="ps-name">{preview.name}</div>
      <div class="ps-sub">
        {getSpeciesDef()?.name || '—'} · {getClassDef()?.name || '—'} · Level {getTotalLevel()}
      </div>

      <div class="ps-stats">
        <div class="ps-stat">
          <div class="ps-stat-value">{preview.hp.max}</div>
          <div class="ps-stat-label">HP</div>
        </div>
        <div class="ps-stat">
          <div class="ps-stat-value">{preview.ac}</div>
          <div class="ps-stat-label">AC</div>
        </div>
        <div class="ps-stat">
          <div class="ps-stat-value">+{preview.proficiencyBonus}</div>
          <div class="ps-stat-label">Prof</div>
        </div>
      </div>

      <div class="ps-divider"></div>

      <div class="ps-section">
        <h4>Atributos</h4>
        <div class="ps-abilities">
          {#each Object.entries(preview.abilities) as [ab, score]}
            <div class="ps-ability">
              <span class="psa-ab">{ab}</span>
              <span class="psa-score">{score}</span>
              <span class="psa-mod">{getMod(score)}</span>
            </div>
          {/each}
        </div>
      </div>

      {#if preview.savingThrows && Object.keys(preview.savingThrows).length > 0}
        <div class="ps-section">
          <h4>Salvaguardas</h4>
          <div class="ps-list">
            {#each Object.entries(preview.savingThrows).filter(([,v]) => v.proficient) as [ab]}
              <div class="ps-item">{ab}</div>
            {/each}
          </div>
        </div>
      {/if}

      {#if preview.features && preview.features.length > 0}
        <div class="ps-section">
          <h4>Features</h4>
          <div class="ps-list">
            {#each preview.features.slice(0, 6) as feat}
              <div class="ps-item">{feat}</div>
            {/each}
            {#if preview.features.length > 6}
              <div class="ps-item ps-item-more">+{preview.features.length - 6} mais</div>
            {/if}
          </div>
        </div>
      {/if}

      {#if preview.spellSlots && Object.keys(preview.spellSlots).length > 0}
        <div class="ps-section">
          <h4>Spell Slots</h4>
          <div class="ps-slots">
            {#each Object.entries(preview.spellSlots).filter(([,c]) => c > 0) as [lvl, count]}
              <div class="ps-slot">
                <span class="ps-slot-lvl">{lvl}º</span>
                <span class="ps-slot-count">{count}</span>
              </div>
            {/each}
          </div>
        </div>
      {/if}

      {#if pendingChoices.length > 0}
        <div class="ps-divider"></div>
        <div class="ps-warnings">
          {#each pendingChoices as choice}
            <div class="ps-warning">
              <span>⚠</span>
              <span>{choice.name}</span>
            </div>
          {/each}
        </div>
      {/if}
    {/if}
  </div>
</aside>

<style>
  .preview-sidebar {
    position: sticky;
    top: 76px;
    background: linear-gradient(135deg, var(--on-bg-surface) 0%, var(--on-bg-elevated) 100%);
    border: 1px solid var(--on-border);
    border-radius: 12px;
    overflow: hidden;
    min-height: 200px;
  }
  .ps-header {
    padding: 14px 18px;
    border-bottom: 1px solid var(--on-border);
    display: flex;
    align-items: center;
    gap: 10px;
  }
  .ps-header h3 {
    font-size: 12px;
    font-weight: 500;
    margin: 0;
    text-transform: uppercase;
    letter-spacing: 1.5px;
    color: var(--on-text-muted);
  }
  .ps-dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: var(--on-green);
    animation: pulse 2s infinite;
  }
  .ps-dot-loading {
    background: var(--on-gold);
    animation: spin 1s linear infinite;
    border-radius: 2px;
  }
  @keyframes pulse { 0%, 100% { opacity: 1; } 50% { opacity: 0.4; } }
  @keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }

  .ps-body { padding: 18px; }
  .ps-empty {
    text-align: center;
    padding: 32px 16px;
    color: var(--on-text-dim);
    font-size: 13px;
  }
  .ps-spinner {
    width: 24px;
    height: 24px;
    border: 2px solid var(--on-border);
    border-top-color: var(--on-red);
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
    margin: 0 auto 12px;
  }

  .ps-name { font-size: 18px; font-weight: 500; color: var(--on-text); margin: 0 0 2px 0; word-break: break-word; }
  .ps-sub { font-size: 12px; color: var(--on-text-dim); margin: 0 0 16px 0; }

  .ps-stats {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 8px;
    margin-bottom: 16px;
  }
  .ps-stat {
    text-align: center;
    padding: 10px 6px;
    background: var(--on-bg-root);
    border-radius: 8px;
    border: 1px solid var(--on-border);
  }
  .ps-stat-value { font-size: 22px; font-weight: 500; color: var(--on-gold); line-height: 1; }
  .ps-stat-label { font-size: 9px; text-transform: uppercase; letter-spacing: 1px; color: var(--on-text-dim); margin-top: 4px; }

  .ps-divider { height: 1px; background: var(--on-border); margin: 14px 0; }

  .ps-section { margin-bottom: 14px; }
  .ps-section h4 {
    font-size: 10px;
    text-transform: uppercase;
    letter-spacing: 1.5px;
    color: var(--on-text-dim);
    margin: 0 0 8px 0;
  }
  .ps-abilities {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 6px;
  }
  .ps-ability {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1px;
    padding: 6px 4px;
    background: var(--on-bg-root);
    border-radius: 6px;
    border: 1px solid var(--on-border);
  }
  .psa-ab { font-size: 9px; color: var(--on-text-dim); text-transform: uppercase; letter-spacing: 0.5px; }
  .psa-score { font-size: 15px; font-weight: 500; color: var(--on-text); line-height: 1; }
  .psa-mod { font-size: 11px; color: var(--on-green); }

  .ps-list { display: flex; flex-direction: column; gap: 4px; }
  .ps-item {
    font-size: 12px;
    color: var(--on-text-muted);
    padding: 5px 8px;
    background: var(--on-bg-root);
    border-radius: 5px;
    border: 1px solid var(--on-border);
  }
  .ps-item-more { color: var(--on-text-dim); font-style: italic; text-align: center; }

  .ps-slots { display: flex; gap: 6px; flex-wrap: wrap; }
  .ps-slot {
    display: flex;
    align-items: center;
    gap: 4px;
    padding: 4px 8px;
    background: var(--on-bg-root);
    border: 1px solid var(--on-border);
    border-radius: 5px;
    font-size: 12px;
  }
  .ps-slot-lvl { color: var(--on-text-dim); }
  .ps-slot-count { color: var(--on-gold); font-weight: 500; }

  .ps-warnings { display: flex; flex-direction: column; gap: 6px; }
  .ps-warning {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 12px;
    color: #ff6b6b;
    padding: 8px 10px;
    background: rgba(197, 0, 9, 0.06);
    border-radius: 6px;
    border: 1px solid rgba(197, 0, 9, 0.15);
  }

  @media (max-width: 900px) {
    .preview-sidebar {
      position: static;
      order: -1;
      margin-bottom: 16px;
    }
  }
</style>
