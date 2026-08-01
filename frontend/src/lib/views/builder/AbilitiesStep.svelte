<script lang="ts">
  // Step 4 — Abilities: standard array, point buy, or rolled (4d6 drop lowest).
  import { ABILITIES } from '../../types'
  import type { Ability } from '../../types'
  import {
    POINT_BUY_BUDGET,
    POINT_BUY_COST,
    POINT_BUY_MAX,
    POINT_BUY_MIN,
    STANDARD_ARRAY,
    assignAbility,
    builder,
    rollScores,
    setMethod,
    setPointBuy,
  } from '../../builder.svelte'
  import type { AbilityMethod } from '../../builder.svelte'

  const METHODS: { id: AbilityMethod; label: string; hint: string }[] = [
    { id: 'standard-array', label: 'Array padrão', hint: '15, 14, 13, 12, 10, 8' },
    { id: 'point-buy', label: 'Point buy', hint: `orçamento ${POINT_BUY_BUDGET}` },
    { id: 'rolled', label: 'Rolado', hint: '4d6, descarta o menor' },
  ]

  const method = $derived(builder.value.abilities.method)

  function pool(): number[] {
    if (method === 'rolled') return builder.value.abilities.rolled
    return [...STANDARD_ARRAY]
  }

  function poolCount(v: number): number {
    return pool().filter((x) => x === v).length
  }

  function optionsFor(ability: Ability): number[] {
    const used: Record<number, number> = {}
    for (const ab of ABILITIES) {
      if (ab === ability) continue
      const v = builder.value.abilities.assigned[ab]
      if (typeof v === 'number') used[v] = (used[v] ?? 0) + 1
    }
    return pool().filter((v) => (used[v] ?? 0) < poolCount(v))
  }

  function pointBuySpent(): number {
    return Object.values(builder.value.abilities.pointBuy).reduce(
      (acc, v) => acc + POINT_BUY_COST[v],
      0
    )
  }
</script>

<h2>Atributos</h2>

<div class="tabs">
  {#each METHODS as m (m.id)}
    <button class:active={method === m.id} onclick={() => setMethod(m.id)}>
      {m.label} <span class="hint">{m.hint}</span>
    </button>
  {/each}
</div>

{#if method === 'rolled'}
  <div class="roll-row">
    <button class="btn" onclick={rollScores}>Rolar 4d6 (drop lowest)</button>
    {#if builder.value.abilities.rolled.length > 0}
      <div class="chips">
        {#each builder.value.abilities.rolled as v, i (i)}
          <span class="chip">{v}</span>
        {/each}
      </div>
    {/if}
  </div>
{/if}

{#if method === 'point-buy'}
  <p class:over={pointBuySpent() > POINT_BUY_BUDGET} class="budget">
    Gasto: {pointBuySpent()} / {POINT_BUY_BUDGET}
  </p>
{/if}

<ul class="abilities">
  {#each ABILITIES as ab (ab)}
    {@const score = builder.value.abilities.assigned[ab]}
    <li>
      <span class="name">{ab}</span>
      {#if method === 'point-buy'}
        <div class="stepper">
          <button
            class="btn"
            onclick={() => setPointBuy(ab, builder.value.abilities.pointBuy[ab] - 1)}>−</button
          >
          <span class="val">{builder.value.abilities.pointBuy[ab]}</span>
          <button
            class="btn"
            onclick={() => setPointBuy(ab, builder.value.abilities.pointBuy[ab] + 1)}>+</button
          >
        </div>
      {:else}
        <select
          value={typeof score === 'number' ? score : ''}
          onchange={(e) => {
            const v = Number(e.currentTarget.value)
            assignAbility(ab, e.currentTarget.value === '' ? null : v)
          }}
        >
          <option value="">—</option>
          {#each optionsFor(ab) as v (v)}
            <option value={v}>{v}</option>
          {/each}
        </select>
      {/if}
    </li>
  {/each}
</ul>

{#if method === 'point-buy'}
  <p class="muted">
    Intervalo {POINT_BUY_MIN}–{POINT_BUY_MAX}. Custo: 8→0, 9→1, 10→2, 11→3, 12→4, 13→5, 14→7, 15→9.
  </p>
{/if}

<style>
  .muted {
    opacity: 0.7;
  }
  .tabs {
    display: flex;
    gap: 0.5rem;
    margin-bottom: 1rem;
    flex-wrap: wrap;
  }
  .tabs button {
    font: inherit;
    padding: 0.4rem 1rem;
    border-radius: 999px;
    border: 1px solid var(--border);
    background: var(--code-bg);
    color: var(--text);
    cursor: pointer;
  }
  .tabs button.active {
    color: var(--accent);
    background: var(--accent-bg);
    border-color: var(--accent-border);
  }
  .hint {
    opacity: 0.6;
    font-size: 0.75rem;
    margin-left: 0.3rem;
  }
  .roll-row {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    margin-bottom: 0.75rem;
  }
  .chips {
    display: flex;
    flex-wrap: wrap;
    gap: 0.35rem;
  }
  .budget {
    margin: 0 0 0.75rem;
    font-weight: 600;
    color: var(--text-h);
  }
  .budget.over {
    color: var(--danger);
  }
  .abilities {
    list-style: none;
    margin: 0;
    padding: 0;
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(10rem, 1fr));
    gap: 0.75rem;
  }
  .abilities li {
    background: var(--code-bg);
    border: 1px solid var(--border);
    border-radius: 10px;
    padding: 0.9rem 1rem;
    display: grid;
    gap: 0.5rem;
    justify-items: center;
  }
  .name {
    font-size: 0.8rem;
    font-weight: 700;
    opacity: 0.7;
  }
  .stepper {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }
  .val {
    font-weight: 700;
    font-size: 1.3rem;
    min-width: 2rem;
    text-align: center;
    color: var(--text-h);
  }
  select {
    font: inherit;
    color: var(--text-h);
    background: var(--bg);
    border: 1px solid var(--border);
    border-radius: 8px;
    padding: 0.3rem 0.6rem;
  }
</style>
