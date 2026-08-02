<script lang="ts">
  // RF-08 (Phase 2 v1.1): import a D&D Beyond PDF (or JSON) into an editable
  // draft. Extraction is tolerant; the review step is mandatory before create.
  import { onMount } from 'svelte'
  import { navigate } from '../router.svelte'
  import { loadContent } from '../content.svelte'
  import { buildDraft, createCharacter } from '../characters.svelte'
  import { ABILITIES, type Ability } from '../types'
  import { extractTextFromPDF, parseDDBText, toBuildRequest, type ParsedDDB } from '../import/ddb'

  let status = $state<'idle' | 'parsing' | 'review' | 'saving'>('idle')
  let error = $state('')
  let parsed = $state<ParsedDDB | null>(null)
  let saving = $state(false)
  let dragOver = $state(false)
  let contentReady = $state(false)
  let fileInput = $state<HTMLInputElement>()
  const pendingAbility = $state<Partial<Record<Ability, number>>>({})
  let reviewChecked = $state(false)

  onMount(async () => {
    try {
      await loadContent()
      contentReady = true
    } catch {
      error = 'Failed to load rules content. Check server connection.'
    }
  })

  function applyParsed(p: ParsedDDB) {
    parsed = p
    Object.assign(pendingAbility, p.abilityScores)
    reviewChecked = false
  }

  async function handleFile(file: File | undefined) {
    if (!file) return
    error = ''
    status = 'parsing'
    try {
      const isPdf = file.type === 'application/pdf' || /\.pdf$/i.test(file.name)
      const isJson = /\.json$/i.test(file.name)
      if (!isPdf && !isJson) throw new Error('Only PDF or JSON files are supported.')

      const content = contentReady ? await loadContent() : null
      if (!content) throw new Error('Rules content unavailable.')

      if (isJson) {
        const { parseDDBJSON } = await import('../import/ddb')
        applyParsed(await parseDDBJSON(file, content))
      } else {
        const text = await extractTextFromPDF(file)
        applyParsed(parseDDBText(text, content))
      }
      status = 'review'
    } catch (e) {
      status = 'idle'
      error = e instanceof Error ? e.message : 'Failed to read file.'
    }
  }

  async function create() {
    if (!parsed) return
    status = 'saving'
    saving = true
    const draft = toBuildRequest({ ...parsed, abilityScores: pendingAbility })
    try {
      const sheet = await buildDraft(draft)
      const c = createCharacter({ ...draft })
      c.sheet = sheet
      navigate(`/characters/${c.id}`)
    } catch {
      status = 'review'
      saving = false
      error = 'Failed to save character. Review values and try again.'
    }
  }
</script>

<div class="page-head">
  <h1>Import</h1>
  <p class="muted">Import a D&D Beyond character (PDF or JSON) and create an editable draft.</p>
</div>

{#if status === 'idle' || status === 'parsing'}
  <label
    class="drop"
    class:dragging={dragOver}
    ondragover={(e) => {
      e.preventDefault()
      dragOver = true
    }}
    ondragleave={() => (dragOver = false)}
    ondrop={(e) => {
      e.preventDefault()
      dragOver = false
      handleFile(e.dataTransfer?.files[0])
    }}
  >
    <input
      type="file"
      accept="application/pdf,.pdf,application/json,.json"
      bind:this={fileInput}
      onchange={(e) => handleFile(e.currentTarget.files?.[0])}
      hidden
    />
    <strong>{status === 'parsing' ? 'Reading file…' : 'Drag a PDF/JSON here'}</strong>
    <span class="muted"
      >or
      <button class="btn-link" onclick={() => fileInput?.click()}>select file</button>
      (max 5 MiB)</span
    >
  </label>
{/if}

{#if error}
  <p class="error" role="alert">{error}</p>
{/if}

{#if status === 'review' && parsed}
  <section class="review">
    <h2>Review</h2>
    <p class="muted">Confirm extracted values before creating the character.</p>

    <div class="field">
      <label for="imp-name">Name</label>
      <input id="imp-name" bind:value={parsed.name} />
    </div>

    <div class="field">
      <label for="imp-class">Class</label>
      <select id="imp-class" bind:value={parsed.classes[0].id}>
        <option value={parsed.classes[0].id}>{parsed.classes[0].id}</option>
      </select>
      <input
        type="number"
        min="1"
        max="20"
        aria-label="Level"
        bind:value={parsed.classes[0].level}
      />
    </div>

    <div class="field">
      <label id="abilities-label">Abilities</label>
      <div class="scores" aria-labelledby="abilities-label">
        {#each ABILITIES as a (a)}
          <div class="score">
            <span>{a}</span>
            <input type="number" min="1" max="30" bind:value={pendingAbility[a]} />
          </div>
        {/each}
      </div>
    </div>

    <div class="field">
      <label id="spells-label">Spells</label>
      <ul class="chips" aria-labelledby="spells-label">
        {#each parsed.spellIds as id (id)}
          <li>{id}</li>
        {/each}
      </ul>
    </div>

    <div class="field">
      <label id="feat-label">Feat</label>
      <ul class="chips" aria-labelledby="feat-label">
        {#each parsed.featIds as id (id)}
          <li>{id}</li>
        {/each}
      </ul>
    </div>

    <label class="check">
      <input type="checkbox" bind:checked={reviewChecked} />
      I confirm I have reviewed the values above.
    </label>

    <div class="actions">
      <button class="btn" onclick={() => (status = 'idle')}>Cancel</button>
      <button class="btn primary" disabled={!reviewChecked || saving} onclick={create}>
        {saving ? 'Saving…' : 'Create Character'}
      </button>
    </div>
  </section>
{/if}

<p class="muted back"><a href="#/characters">← Back to characters</a></p>

<style>
  .page-head {
    margin-bottom: 1rem;
  }
  h1 {
    margin: 0;
    color: var(--text-h);
  }
  h2 {
    margin-top: 1.5rem;
    color: var(--text-h);
  }
  .muted {
    opacity: 0.7;
  }
  .drop {
    display: block;
    border: 2px dashed var(--border);
    border-radius: 12px;
    padding: 2rem;
    text-align: center;
    cursor: pointer;
  }
  .drop.dragging {
    border-color: var(--accent);
    background: var(--accent-bg);
  }
  .drop strong {
    display: block;
    color: var(--text-h);
    margin-bottom: 0.25rem;
  }
  .error {
    color: var(--danger);
    margin-top: 1rem;
  }
  .review {
    margin-top: 1.5rem;
  }
  .field {
    margin-bottom: 1rem;
  }
  .field label {
    display: block;
    font-size: 0.9rem;
    color: var(--text-h);
    margin-bottom: 0.25rem;
  }
  input,
  select {
    font: inherit;
    color: var(--text-h);
    background: var(--bg);
    border: 1px solid var(--border);
    border-radius: 8px;
    padding: 0.45rem 0.75rem;
  }
  #imp-name {
    width: 100%;
  }
  .scores {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(4.5rem, 1fr));
    gap: 0.5rem;
  }
  .score {
    display: flex;
    align-items: center;
    gap: 0.4rem;
  }
  .score span {
    font-weight: 600;
    color: var(--text-h);
  }
  .score input {
    width: 4rem;
  }
  .chips {
    list-style: none;
    display: flex;
    flex-wrap: wrap;
    gap: 0.4rem;
    margin: 0;
    padding: 0;
  }
  .chips li {
    background: var(--code-bg);
    border: 1px solid var(--border);
    border-radius: 999px;
    padding: 0.15rem 0.7rem;
    font-size: 0.85rem;
  }
  .check {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    margin: 1rem 0;
    color: var(--text-h);
  }
  .actions {
    display: flex;
    gap: 0.5rem;
  }
  .back {
    margin-top: 1.5rem;
  }
</style>