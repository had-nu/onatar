<script lang="ts">
  // Step 2 — Background.
  import { content } from '../../content.svelte'
  import { dataString } from '../../types'
  import type { Background } from '../../types'
  import { builder, suggestedBackgroundsForClass, selectBackground } from '../../builder.svelte'

  const contentData = $derived(content.value)

  function skillsOf(b: Background): string[] {
    const skills = b.data?.skills
    return Array.isArray(skills) ? skills.map(String) : []
  }

  function isSuggested(id: string): boolean {
    return suggestedBackgroundsForClass().includes(id)
  }
</script>

<h2>Escolhe o background</h2>
<p class="muted">O background define as origens e skills do personagem.</p>

<ul class="grid">
  {#each contentData?.backgrounds ?? [] as b (b.id)}
    {@const selected = builder.value.draft.backgroundId === b.id}
    <li>
      <button class:selected class="card" onclick={() => selectBackground(b.id)}>
        <h3>
          {b.name}
          {#if isSuggested(b.id)}<span class="badge">sugerido</span>{/if}
        </h3>
        <p class="muted">{dataString(b.data, 'description')}</p>
        {#if skillsOf(b).length > 0}
          <p class="chips-label">Skills</p>
          <p class="chips">
            {#each skillsOf(b) as s (s)}<span class="chip">{s}</span>{/each}
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
