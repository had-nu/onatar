<script lang="ts">
  // Live preview sidebar: debounced POST /build while the wizard draft changes.
  import { buildDraft } from '../../characters.svelte'
  import { builder, validateStep } from '../../builder.svelte'
  import type { Sheet } from '../../types'

  let sheet = $state<Sheet | null>(null)
  let status = $state<'idle' | 'loading' | 'ready' | 'error'>('idle')
  let error = $state('')
  let timer: ReturnType<typeof setTimeout> | undefined

  $effect(() => {
    const draft = builder.value.draft
    clearTimeout(timer)
    if (!validateStep('class')) {
      sheet = null
      status = 'idle'
      return
    }
    status = 'loading'
    timer = setTimeout(async () => {
      try {
        sheet = await buildDraft(draft)
        status = 'ready'
      } catch (err) {
        status = 'error'
        error = err instanceof Error ? err.message : String(err)
      }
    }, 350)
    return () => clearTimeout(timer)
  })

  function modSign(mod: number): string {
    return mod >= 0 ? `+${mod}` : `${mod}`
  }
</script>

<h2 class="side-title">Preview</h2>

{#if status === 'idle'}
  <p class="muted">Select at least a class to see live sheet.</p>
{:else if status === 'loading'}
  <p class="muted">Calculating sheet…</p>
{:else if status === 'error'}
  <p class="error-text">{error}</p>
{:else if sheet}
  <dl class="stats">
    <div>
      <dt>Level</dt>
      <dd>{sheet.level}</dd>
    </div>
    <div>
      <dt>HP</dt>
      <dd>{sheet.hp.max}</dd>
    </div>
    <div>
      <dt>AC</dt>
      <dd>{sheet.ac}</dd>
    </div>
    <div>
      <dt>Prof.</dt>
      <dd>{modSign(sheet.proficiencyBonus)}</dd>
    </div>
  </dl>

  {#if sheet.spellSlots.some((n) => n > 0)}
    <div class="slots">
      {#each sheet.spellSlots as n, i (i)}
        {#if n > 0}<span class="chip">{i + 1}º ×{n}</span>{/if}
      {/each}
    </div>
  {/if}

  {#if sheet.features.length > 0}
    <ul class="features">
      {#each sheet.features as f (f.name + f.level)}
        <li>
          <span class="feature-name">{f.name}</span>
          <span class="muted">level {f.level}</span>
        </li>
      {/each}
    </ul>
  {/if}

  {#if sheet.pendingChoices.length > 0}
    <p class="chips-label">Pending Choices</p>
    <ul class="features">
      {#each sheet.pendingChoices as pc (pc.type + pc.description)}
        <li>{pc.description}</li>
      {/each}
    </ul>
  {/if}
{/if}

<style>
  .side-title {
    margin: 0 0 0.75rem;
    color: var(--text-h);
  }
  .muted {
    opacity: 0.7;
  }
  .error-text {
    color: var(--danger);
  }
  .stats {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 0.5rem;
    margin: 0 0 0.75rem;
  }
  .stats > div {
    background: var(--code-bg);
    border: 1px solid var(--border);
    border-radius: 10px;
    padding: 0.75rem;
    text-align: center;
  }
  .stats dt {
    font-size: 0.72rem;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    opacity: 0.6;
  }
  .stats dd {
    margin: 0;
    font-size: 1.4rem;
    font-weight: 700;
    color: var(--text-h);
  }
  .slots {
    display: flex;
    flex-wrap: wrap;
    gap: 0.35rem;
    margin-bottom: 0.75rem;
  }
  .features {
    list-style: none;
    margin: 0;
    padding: 0;
    display: grid;
    gap: 0.4rem;
  }
  .features li {
    font-size: 0.85rem;
    display: flex;
    justify-content: space-between;
    gap: 0.5rem
  }
  .feature-name {
    color: var(--text-h);
  }
  .chips-label {
    margin: 0.75rem 0 0.25rem;
    font-size: 0.72rem;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    opacity: 0.6;
  }
</style>