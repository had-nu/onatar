<script lang="ts">
  import { draft, preview, pendingChoices, setSubclass, setName, saveCharacter, resetBuilder, nextStep, step, setStep } from '$lib/builder.svelte';
  import { getClassDef, getSpeciesDef, getBackgroundDef, getTotalLevel } from '$lib/builder.svelte';
  import Button from '$lib/ui/Button.svelte';
  import Card from '$lib/ui/Card.svelte';

  let saving = $state(false);
  let saveError = $state<string | null>(null);
  let saveSuccess = $state(false);

  async function handleSave() {
    saving = true;
    saveError = null;
    const ok = await saveCharacter();
    saving = false;
    if (ok) saveSuccess = true;
    else saveError = 'Failed to save. Try again.';
  }

  function getMod(score: number): string {
    const mod = Math.floor((score - 10) / 2);
    return mod >= 0 ? `+${mod}` : `${mod}`;
  }

  const classDef = $derived(getClassDef());
  const speciesDef = $derived(getSpeciesDef());
  const backgroundDef = $derived(getBackgroundDef());
  const totalLevel = $derived(getTotalLevel());
</script>

<div class="step-review">
  <h2 class="step-title">Final Review</h2>
  <p class="step-desc">Verify all details before forging your character.</p>

  <!-- Name input -->
  <div class="name-section">
    <label for="char-name">Character Name</label>
    <input
      id="char-name"
      type="text"
      value={draft.name}
      oninput={(e) => setName(e.currentTarget.value)}
      placeholder="e.g., Thorin Oakenshield"
      class="name-input"
    />
  </div>

  <!-- Pending choices -->
  {#if pendingChoices.length > 0}
    <div class="pending-section">
      <h3>Pending Choices</h3>
      <div class="pending-list">
        {#each pendingChoices as choice}
          <Card variant="elevated" class="pending-card">
            <div class="pc-type">{choice.type === 'subclass' ? 'Subclass' : choice.type === 'spell' ? 'Spell' : choice.type === 'skill' ? 'Skill' : 'Ability'}</div>
            <div class="pc-name">{choice.name}</div>
            <p class="pc-desc">{choice.description}</p>
            {#if choice.options.length > 0}
              <div class="pc-options">
                {#each choice.options as opt}
                  <button
                    class="pc-opt-btn"
                    onclick={() => {
                      if (choice.type === 'subclass' && choice.classId) {
                        setSubclass(choice.classId, opt.id);
                      }
                    }}
                  >
                    {opt.name}
                  </button>
                {/each}
              </div>
            {/if}
          </Card>
        {/each}
      </div>
    </div>
  {/if}

  <!-- Summary cards -->
  <div class="review-grid">
    <Card variant="default" class="review-card">
      <h4>Class</h4>
      <p class="review-value">{classDef?.name || '—'} {draft.classes[0]?.subclassId ? `(${classDef?.subClasses?.find(s => s.id === draft.classes[0].subclassId)?.name || ''})` : ''}</p>
      <p class="review-meta">Level {totalLevel} · d{classDef?.hitDie || '—'}</p>
    </Card>

    <Card variant="default" class="review-card">
      <h4>Species</h4>
      <p class="review-value">{speciesDef?.name || '—'}</p>
      <p class="review-meta">{speciesDef?.size || '—'} · {speciesDef?.speed || '—'}ft</p>
    </Card>

    <Card variant="default" class="review-card">
      <h4>Background</h4>
      <p class="review-value">{backgroundDef?.name || '—'}</p>
      <p class="review-meta">{backgroundDef?.skillProficiencies.join(', ') || '—'}</p>
    </Card>

    <Card variant="default" class="review-card">
      <h4>Abilities</h4>
      <div class="review-abilities">
        {#each Object.entries(draft.abilityScores) as [ab, score]}
          <div class="ra-item">
            <span class="ra-ab">{ab}</span>
            <span class="ra-score">{score}</span>
            <span class="ra-mod">{getMod(score)}</span>
          </div>
        {/each}
      </div>
    </Card>
  </div>

  <!-- Preview stats -->
  {#if preview}
    <div class="preview-section">
      <h3>Sheet Preview</h3>
      <div class="preview-stats-row">
        <div class="ps-item">
          <span class="ps-value">{preview.hp.max}</span>
          <span class="ps-label">Max HP</span>
        </div>
        <div class="ps-item">
          <span class="ps-value">{preview.ac}</span>
          <span class="ps-label">AC</span>
        </div>
        <div class="ps-item">
          <span class="ps-value">+{preview.proficiencyBonus}</span>
          <span class="ps-label">Proficiency</span>
        </div>
        <div class="ps-item">
          <span class="ps-value">{preview.features.length}</span>
          <span class="ps-label">Features</span>
        </div>
      </div>
    </div>
  {/if}

  <!-- Actions -->
  <div class="review-actions">
    {#if saveSuccess}
      <div class="success-msg">✓ Character saved successfully!</div>
      <Button variant="secondary" onclick={() => { resetBuilder(); setStep(0); }}>Create New Character</Button>
    {:else}
      <Button variant="outline" onclick={() => setStep(0)}>Restart</Button>
      <Button variant="primary" onclick={handleSave} disabled={saving || pendingChoices.length > 0 || !draft.name}>
        {saving ? 'Saving...' : '⚔ Forge Character'}
      </Button>
    {/if}
  </div>

  {#if saveError}
    <div class="save-error">{saveError}</div>
  {/if}
</div>

<style>
  .step-review { animation: fadeIn 0.3s ease; }
  @keyframes fadeIn { from { opacity: 0; transform: translateY(8px); } to { opacity: 1; transform: translateY(0); } }

  .step-title { font-size: 28px; font-weight: 500; margin: 0 0 8px 0; }
  .step-desc { font-size: 14px; color: var(--on-text-dim); margin: 0 0 24px 0; }

  .name-section { margin-bottom: 24px; }
  .name-section label {
    display: block;
    font-size: 12px;
    text-transform: uppercase;
    letter-spacing: 1px;
    color: var(--on-text-dim);
    margin-bottom: 8px;
  }
  .name-input {
    width: 100%;
    max-width: 400px;
    padding: 12px 16px;
    font-size: 16px;
    background: var(--on-bg-surface);
    border: 1px solid var(--on-border);
    border-radius: 10px;
    color: var(--on-text);
    outline: none;
    transition: border-color 0.15s;
  }
  .name-input:focus { border-color: var(--on-red); }
  .name-input::placeholder { color: var(--on-text-dim); }

  .pending-section { margin-bottom: 24px; }
  .pending-section h3 { font-size: 16px; font-weight: 500; margin: 0 0 12px 0; }
  .pending-list { display: flex; flex-direction: column; gap: 10px; }
  :global(.pending-card) { padding: 16px; }
  .pc-type {
    font-size: 10px;
    text-transform: uppercase;
    letter-spacing: 1.5px;
    color: var(--on-red);
    margin-bottom: 4px;
  }
  .pc-name { font-size: 15px; font-weight: 500; color: var(--on-text); margin-bottom: 4px; }
  .pc-desc { font-size: 12px; color: var(--on-text-dim); margin: 0 0 10px 0; }
  .pc-options { display: flex; gap: 8px; flex-wrap: wrap; }
  .pc-opt-btn {
    padding: 6px 14px;
    border-radius: 6px;
    border: 1px solid var(--on-border);
    background: var(--on-bg-root);
    color: var(--on-text-muted);
    font-size: 12px;
    cursor: pointer;
    transition: all 0.15s;
  }
  .pc-opt-btn:hover { border-color: var(--on-red); color: var(--on-text); }

  .review-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 12px;
    margin-bottom: 24px;
  }
  :global(.review-card) { padding: 16px; }
  :global(.review-card) h4 {
    font-size: 11px;
    text-transform: uppercase;
    letter-spacing: 1px;
    color: var(--on-text-dim);
    margin: 0 0 8px 0;
  }
  .review-value { font-size: 16px; font-weight: 500; color: var(--on-text); margin: 0 0 2px 0; }
  .review-meta { font-size: 12px; color: var(--on-text-dim); margin: 0; }
  .review-abilities {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 8px;
  }
  .ra-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 2px;
    padding: 8px;
    background: var(--on-bg-root);
    border-radius: 6px;
    border: 1px solid var(--on-border);
  }
  .ra-ab { font-size: 10px; color: var(--on-text-dim); text-transform: uppercase; }
  .ra-score { font-size: 18px; font-weight: 500; color: var(--on-text); }
  .ra-mod { font-size: 12px; color: var(--on-green); }

  .preview-section { margin-bottom: 24px; }
  .preview-section h3 { font-size: 16px; font-weight: 500; margin: 0 0 12px 0; }
  .preview-stats-row {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 12px;
  }
  .ps-item {
    text-align: center;
    padding: 14px;
    background: var(--on-bg-surface);
    border: 1px solid var(--on-border);
    border-radius: 10px;
  }
  .ps-value { font-size: 24px; font-weight: 500; color: var(--on-gold); line-height: 1; }
  .ps-label { font-size: 10px; text-transform: uppercase; letter-spacing: 1px; color: var(--on-text-dim); margin-top: 4px; display: block; }

  .review-actions {
    display: flex;
    gap: 12px;
    justify-content: flex-end;
    padding-top: 20px;
    border-top: 1px solid var(--on-border);
  }
  .success-msg {
    font-size: 14px;
    color: var(--on-green);
    font-weight: 500;
    display: flex;
    align-items: center;
    gap: 8px;
  }
  .save-error {
    margin-top: 12px;
    padding: 10px 14px;
    background: rgba(197, 0, 9, 0.08);
    border: 1px solid rgba(197, 0, 9, 0.2);
    border-radius: 8px;
    font-size: 13px;
    color: #ff6b6b;
    text-align: right;
  }

  @media (max-width: 600px) {
    .review-grid { grid-template-columns: 1fr; }
    .preview-stats-row { grid-template-columns: repeat(2, 1fr); }
    .review-actions { flex-direction: column; }
  }
</style>