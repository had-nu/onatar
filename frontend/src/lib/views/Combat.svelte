<script lang="ts">
  // RF-09 (Fase 2 v1.2): combat tracker — initiative, HP, conditions.
  import {
    addCharacterToCombat,
    addCombatant,
    damage,
    deleteSession,
    getSession,
    heal,
    listSessions,
    newCombatSession,
    nextTurn,
    prevTurn,
    removeCombatant,
    setInitiative,
    sortByInitiative,
    toggleCondition,
    type CombatSession,
  } from '../combat.svelte'
  import { listCharacters } from '../characters.svelte'
  import { CONDITIONS } from '../types'

  let selectedId = $state('')
  let newName = $state('')
  let combatants = $state<CombatSession | undefined>(undefined)

  function open(id: string) {
    selectedId = id
    combatants = getSession(id)
  }

  function create() {
    const s = newCombatSession(newName)
    newName = ''
    open(s.id)
  }

  function addFromList(characterId: string) {
    if (!selectedId || !characterId) return
    addCharacterToCombat(selectedId, characterId)
    combatants = getSession(selectedId)
  }

  function addRaw() {
    if (!selectedId) return
    addCombatant(selectedId, {
      name: 'Novo',
      initiative: 0,
      hpCurrent: 10,
      hpMax: 10,
      conditions: [],
    })
    combatants = getSession(selectedId)
  }

  function apply(fn: (id: string) => void) {
    if (!selectedId) return
    fn(selectedId)
    combatants = getSession(selectedId)
  }

  function hpOf(c: { hpCurrent: number; hpMax: number }): string {
    return `${c.hpCurrent}/${c.hpMax}`
  }
</script>

<div class="page-head">
  <h1>Combate</h1>
  <p class="muted">Regista iniciativa, HP e condições durante os encontros (local).</p>
</div>

{#if !selectedId}
  <form class="create" onsubmit={(e) => e.preventDefault()}>
    <input bind:value={newName} placeholder="Nome do combate" />
    <button class="btn primary" onclick={create}>Novo combate</button>
  </form>

  {#if listSessions().length === 0}
    <p class="muted empty">Ainda não há combates.</p>
  {:else}
    <ul class="sessions">
      {#each listSessions() as s (s.id)}
        <li>
          <button class="btn session" onclick={() => open(s.id)}>{s.name}</button>
        </li>
      {/each}
    </ul>
  {/if}
{:else if combatants}
  <section class="combat">
    <header class="combat-head">
      <h2>{combatants.name}</h2>
      <span class="round">Round {combatants.round}</span>
      <div class="controls">
        <button class="btn" onclick={() => open('')} aria-label="Voltar">←</button>
        <button
          class="btn danger"
          onclick={() => {
            if (confirm(`Terminar o combate "${combatants!.name}"?`)) {
              deleteSession(combatants!.id)
              selectedId = ''
              combatants = undefined
            }
          }}
        >
          Terminar
        </button>
      </div>
    </header>

    <div class="add-row">
      <select
        aria-label="Adicionar personagem"
        onchange={(e) => addFromList(e.currentTarget.value)}
      >
        <option value="">+ Adicionar personagem…</option>
        {#each listCharacters() as c (c.id)}
          <option value={c.id}>{c.name}{c.isNpc ? ' (NPC)' : ''}</option>
        {/each}
      </select>
      <button class="btn" onclick={addRaw}>+ Avulso</button>
      <button class="btn" onclick={() => apply(sortByInitiative)}>Ordenar iniciativa</button>
    </div>

    {#if combatants.combatants.length === 0}
      <p class="muted empty">Adiciona combatentes para começar.</p>
    {:else}
      <ol class="combatants">
        {#each combatants.combatants as c, i (c.id)}
          <li class:current={i === combatants.turnIndex}>
            <div class="row-head">
              <span class="init">
                <input
                  type="number"
                  aria-label={`Iniciativa de ${c.name}`}
                  value={c.initiative}
                  onchange={(e) =>
                    apply(() => setInitiative(combatants!.id, c.id, Number(e.currentTarget.value)))}
                />
              </span>
              <strong>{c.name}</strong>
              <span class="hp">{hpOf(c)}</span>
            </div>
            <div class="row-actions">
              <button class="btn sm" onclick={() => apply(() => damage(combatants!.id, c.id, 1))}
                >−1</button
              >
              <button class="btn sm" onclick={() => apply(() => damage(combatants!.id, c.id, 5))}
                >−5</button
              >
              <button class="btn sm" onclick={() => apply(() => heal(combatants!.id, c.id, 5))}
                >+5</button
              >
              <button class="btn sm" onclick={() => apply(() => heal(combatants!.id, c.id, 10))}
                >+10</button
              >
            </div>
            <div class="cond">
              {#each CONDITIONS as cond (cond)}
                <button
                  class="chip"
                  class:on={c.conditions.includes(cond)}
                  onclick={() => apply(() => toggleCondition(combatants!.id, c.id, cond))}
                >
                  {cond}
                </button>
              {/each}
            </div>
            <button
              class="btn sm remove"
              onclick={() => {
                apply(() => removeCombatant(combatants!.id, c.id))
              }}
            >
              Remover
            </button>
          </li>
        {/each}
      </ol>

      <div class="turn-nav">
        <button class="btn" onclick={() => apply(prevTurn)}>← Anterior</button>
        <span class="muted">Turno {combatants.turnIndex + 1} de {combatants.combatants.length}</span
        >
        <button class="btn primary" onclick={() => apply(nextTurn)}>Seguinte →</button>
      </div>
    {/if}
  </section>
{/if}

<style>
  .page-head {
    margin-bottom: 1rem;
  }
  h1 {
    margin: 0;
    color: var(--text-h);
  }
  h2 {
    margin: 0;
    color: var(--text-h);
  }
  .muted {
    opacity: 0.7;
  }
  .empty {
    margin: 1rem 0;
  }
  .create {
    display: flex;
    gap: 0.5rem;
    margin-bottom: 1.5rem;
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
  .sessions {
    list-style: none;
    margin: 0;
    padding: 0;
    display: flex;
    flex-direction: column;
    gap: 0.4rem;
  }
  .session {
    text-align: left;
  }
  .combat-head {
    display: flex;
    align-items: center;
    gap: 1rem;
    flex-wrap: wrap;
    margin-bottom: 1rem;
  }
  .round {
    font-weight: 600;
    color: var(--accent);
  }
  .controls {
    margin-left: auto;
    display: flex;
    gap: 0.4rem;
  }
  .add-row {
    display: flex;
    gap: 0.5rem;
    flex-wrap: wrap;
    margin-bottom: 1rem;
  }
  .combatants {
    list-style: none;
    margin: 0;
    padding: 0;
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }
  .combatants li {
    background: var(--code-bg);
    border: 1px solid var(--border);
    border-radius: 10px;
    padding: 0.75rem 1rem;
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 0.75rem;
  }
  .combatants li.current {
    border-color: var(--accent-border);
    box-shadow: 0 0 0 1px var(--accent-border);
  }
  .row-head {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    min-width: 14rem;
  }
  .init input {
    width: 3.5rem;
  }
  .hp {
    font-variant-numeric: tabular-nums;
    margin-left: auto;
  }
  .row-actions {
    display: flex;
    gap: 0.3rem;
  }
  .cond {
    display: flex;
    flex-wrap: wrap;
    gap: 0.3rem;
    flex: 1;
    min-width: 10rem;
  }
  .chip {
    font: inherit;
    font-size: 0.8rem;
    color: var(--text);
    background: var(--bg);
    border: 1px solid var(--border);
    border-radius: 999px;
    padding: 0.1rem 0.55rem;
    cursor: pointer;
  }
  .chip.on {
    color: var(--accent);
    background: var(--accent-bg);
    border-color: var(--accent-border);
  }
  .remove {
    margin-left: auto;
  }
  .turn-nav {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 1rem;
    margin-top: 1.25rem;
  }
</style>
