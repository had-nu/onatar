<script lang="ts">
  // Step 3 — Species.
  import { content } from '../../content.svelte'
  import { dataString } from '../../types'
  import type { Species } from '../../types'
  import { builder, selectSpecies, suggestedSpeciesForClass } from '../../builder.svelte'

  const contentData = $derived(content.value)

  function traitsOf(s: Species): string[] {
    const traits = s.data?.traits
    if (Array.isArray(traits)) {
      return traits
        .map((t) => (t && typeof t === 'object' && 'name' in t ? String(t.name) : ''))
        .filter(Boolean)
    }
    return []
  }

  function isSuggested(id: string): boolean {
    return suggestedSpeciesForClass().includes(id)
  }
</script>

<h2>Escolhe a espécie</h2>
<p class="muted">As espécies definem traits raciais. As sugestões combinam com a tua classe.</p>

<ul class="grid">
  {#each contentData?.species ?? [] as s (s.id)}
    {@const selected = builder.value.draft.speciesId === s.id}
    <li>
      <button class:selected class="card" onclick={() => selectSpecies(s.id)}>
        <h3>
          {s.name}
          {#if isSuggested(s.id)}<span class="badge">sugerido</span>{/if}
        </h3>
        <p class="muted">{dataString(s.data, 'description')}</p>
        {#if traitsOf(s).length > 0}
          <p class="chips-label">Traits</p>
          <p class="chips">
            {#each traitsOf(s) as t (t)}<span class="chip">{t}</span>{/each}
          </p>
        {/if}
      </button>
    </li>
  {/each}
</ul>

<style>
  .muted {
    opacity: 0.7;
  }
  .grid {
    list-style: none;
    margin: 0;
    padding: 0;
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(16rem, 1fr));
    gap: 0.75rem;
  }
  .card {
    font: inherit;
    text-align: left;
    width: 100%;
    background: var(--code-bg);
    border: 1px solid var(--border);
    border-radius: 10px;
    padding: 1rem 1.25rem;
    cursor: pointer;
    color: var(--text);
    display: grid;
    gap: 0.35rem;
  }
  .card:hover {
    border-color: var(--accent-border);
  }
  .card.selected {
    border-color: var(--accent-border);
    background: var(--accent-bg);
  }
  h3 {
    margin: 0;
    color: var(--text-h);
  }
  .chips-label {
    margin: 0.5rem 0 0.25rem;
    font-size: 0.75rem;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    opacity: 0.6;
  }
  .chips {
    display: flex;
    flex-wrap: wrap;
    gap: 0.35rem;
  }
</style>
