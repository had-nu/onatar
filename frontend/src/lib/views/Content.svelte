<script lang="ts">
  import { onMount } from 'svelte'
  import { contentError, loadContent } from '../content.svelte'
  import { dataString } from '../types'
  import type { Background, Class, Species } from '../types'

  let status = $state<'loading' | 'ready' | 'error'>('loading')
  let tab = $state<'classes' | 'species' | 'backgrounds'>('classes')

  onMount(async () => {
    try {
      await loadContent()
      status = 'ready'
    } catch {
      status = 'error'
    }
  })

  function descriptionOf(data: Record<string, unknown> | undefined): string {
    return dataString(data, 'description')
  }

  function primaryAbilityOf(c: Class): string {
    const v = dataString(c.data, 'primaryAbility')
    return v ? ` · ${v}` : ''
  }

  function traitNames(s: Species): string[] {
    const traits = s.data?.traits
    if (Array.isArray(traits)) {
      return traits
        .map((t) => (t && typeof t === 'object' && 'name' in t ? String(t.name) : ''))
        .filter(Boolean)
    }
    return []
  }

  function skillNames(b: Background): string[] {
    const skills = b.data?.skills
    return Array.isArray(skills) ? skills.map(String) : []
  }
</script>

<div class="page-head">
  <h1>Content</h1>
  <p class="muted">Classes, species, and backgrounds rules (SRD 5.2)</p>
</div>

<div class="tabs">
  <button class:active={tab === 'classes'} onclick={() => (tab = 'classes')}>Classes</button>
  <button class:active={tab === 'species'} onclick={() => (tab = 'species')}>Species</button>
  <button class:active={tab === 'backgrounds'} onclick={() => (tab = 'backgrounds')}>Backgrounds</button>
</div>

{#if status === 'loading'}
  <p>Loading content…</p>
{:else if status === 'error'}
  <div class="error-box">
    <p>
      Failed to load content{contentError.value ? `: ${contentError.value}` : ''}.
    </p>
    <button
      class="btn"
      onclick={() =>
        void loadContent(true)
          .then(() => (status = 'ready'))
          .catch(() => undefined)}
    >
      Try Again
    </button>
  </div>
{:else}
  {#await loadContent()}
    <p>Loading content…</p>
  {:then data}
    {#if tab === 'classes'}
      <ul class="grid">
        {#each data.classes as c (c.id)}
          <li>
            <article class="card">
              <h2>{c.name}</h2>
              <p class="badges">
                <span class="badge">{c.hitDie}</span>
                <span class="badge">{c.spellcaster ? 'spellcaster' : 'martial'}</span>
              </p>
              <p class="muted">{descriptionOf(c.data)}{primaryAbilityOf(c)}</p>
              {#if c.suggestedSpecies.length > 0}
                <p class="chips-label">Suggested Species</p>
                <p class="chips">
                  {#each c.suggestedSpecies as s (s)}
                    <span class="chip">{s}</span>
                  {/each}
                </p>
              {/if}
            </article>
          </li>
        {/each}
      </ul>
    {:else if tab === 'species'}
      <ul class="grid">
        {#each data.species as s (s.id)}
          <li>
            <article class="card">
              <h2>{s.name}</h2>
              <p class="muted">{descriptionOf(s.data)}</p>
              {#if traitNames(s).length > 0}
                <p class="chips-label">Traits</p>
                <p class="chips">
                  {#each traitNames(s) as t (t)}
                    <span class="chip">{t}</span>
                  {/each}
                </p>
              {/if}
            </article>
          </li>
        {/each}
      </ul>
    {:else}
      <ul class="grid">
        {#each data.backgrounds as b (b.id)}
          <li>
            <article class="card">
              <h2>{b.name}</h2>
              <p class="muted">{descriptionOf(b.data)}</p>
              {#if skillNames(b).length > 0}
                <p class="chips-label">Skills</p>
                <p class="chips">
                  {#each skillNames(b) as s (s)}
                    <span class="chip">{s}</span>
                  {/each}
                </p>
              {/if}
            </article>
          </li>
        {/each}
      </ul>
    {/if}
  {/await}
{/if}

<style>
  .page-head {
    margin-bottom: 1rem;
  }
  h1 {
    margin: 0;
    color: var(--text-h);
  }
  .muted {
    opacity: 0.7;
  }
  .tabs {
    display: flex;
    gap: 0.5rem;
    margin-bottom: 1.25rem;
  }
  .tabs button {
    font: inherit;
    padding: 0.4rem 1rem
    border-radius: 999px
    border: 1px solid var(--border)
    background: var(--code-bg)
    color: var(--text)
    cursor: pointer
  }
  .tabs button.active {
    color: var(--accent)
    background: var(--accent-bg)
    border-color: var(--accent-border)
  }
  .grid {
    list-style: none
    margin: 0
    padding: 0
    display: grid
    grid-template-columns: repeat(auto-fill, minmax(16rem, 1fr))
    gap: 1rem
  }
  h2 {
    margin: 0 0 0.5rem
    color: var(--text-h)
  }
  .badges,
  .chips {
    display: flex
    flex-wrap: wrap
    gap: 0.35rem
    margin: 0.4rem 0
  }
  .chips-label {
    margin: 0.75rem 0 0.25rem
    font-size: 0.75rem
    text-transform: uppercase
    letter-spacing: 0.05em
    opacity: 0.6
  }
  .error-box {
    border: 1px solid var(--danger-border, #e5484d)
    background: var(--danger-bg, rgba(229, 72, 77, 0.08))
    border-radius: 8px
    padding: 1.25rem
    display: grid
    gap: 0.75rem
    justify-items: start
  }
</style>