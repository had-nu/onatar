<script lang="ts">
  import { characters, createCharacter, deleteCharacter, starterDraft } from '../characters.svelte'
  import { navigate } from '../router.svelte'
  import type { Character } from '../types'

  function open(c: Character) {
    navigate(`/characters/${c.id}`)
  }

  function create() {
    navigate(`/characters/${createCharacter(starterDraft()).id}`)
  }

  function remove(c: Character) {
    if (window.confirm(`Apagar "${c.name}"?`)) {
      deleteCharacter(c.id)
    }
  }

  function classSummary(c: Character): string {
    const levels = c.draft.classes.map((x) => `${x.id} ${x.level}`).join(', ')
    return levels || 'Sem classes'
  }
</script>

<div class="page-head">
  <h1>Os meus personagens</h1>
  <button class="btn primary" onclick={create}>Novo personagem</button>
</div>

{#if characters.value.length === 0}
  <div class="empty">
    <p>Ainda não tens personagens.</p>
    <button class="btn primary" onclick={create}>Criar o primeiro</button>
  </div>
{:else}
  <ul class="grid">
    {#each characters.value as c (c.id)}
      <li>
        <article class="card">
          <div class="card-head">
            <h2>{c.name}</h2>
            {#if c.isNpc}<span class="badge">NPC</span>{/if}
          </div>
          <p class="muted">{classSummary(c)}</p>
          <div class="actions">
            <button class="btn" onclick={() => open(c)}>Abrir</button>
            <button class="btn danger" onclick={() => remove(c)}>Apagar</button>
          </div>
        </article>
      </li>
    {/each}
  </ul>
{/if}

<style>
  .page-head {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
    margin-bottom: 1.5rem;
  }
  h1 {
    margin: 0;
    color: var(--text-h);
  }
  .grid {
    list-style: none;
    margin: 0;
    padding: 0;
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(16rem, 1fr));
    gap: 1rem;
  }
  .card-head {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }
  h2 {
    margin: 0;
    font-size: 1.1rem;
    color: var(--text-h);
  }
  .muted {
    opacity: 0.7;
  }
  .actions {
    display: flex;
    gap: 0.5rem;
    margin-top: 1rem;
  }
  .empty {
    text-align: center;
    padding: 3rem 0;
    opacity: 0.8;
  }
</style>
