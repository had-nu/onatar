<script lang="ts">
  import { onMount } from 'svelte'
  import { buildCharacter, getCharacter, setCampaign, setLive } from '../characters.svelte'
  import { campaigns } from '../campaigns.svelte'
  import { content } from '../content.svelte'
  import { navigate } from '../router.svelte'
  import { characterToJSON, downloadJSON, exportCharacterPDF } from '../export'
  import { ABILITIES, CONDITIONS, emptyLive } from '../types'
  import type { Sheet, SheetLive } from '../types'

  let { id } = $props<{ id: string }>()

  const character = $derived(getCharacter(id))
  let sheet = $state<Sheet | null>(null)
  let status = $state<'loading' | 'ready' | 'error'>('loading')
  let error = $state('')
  let live = $state<SheetLive>(emptyLive(undefined))
  let sheetNode = $state<HTMLElement>()

  onMount(() => {
    const c = getCharacter(id)
    if (!c) {
      status = 'error'
      error = 'Character not found.'
      return
    }
    sheet = c.sheet ?? null
    live = c.live ?? emptyLive(sheet ?? undefined)
    setLive(id, live)
    status = c.sheet ? 'ready' : 'loading'
    if (!sheet) void compute()
  })

  async function compute() {
    if (!character) return
    status = 'loading'
    error = ''
    try {
      const s = await buildCharacter(character)
      sheet = s
      live = { ...emptyLive(s), ...live }
      setLive(character.id, live)
      status = 'ready'
    } catch (err) {
      status = 'error'
      error = err instanceof Error ? err.message : String(err)
    }
  }

  function retry() {
    void compute()
  }

  function modSign(mod: number): string {
    return mod >= 0 ? `+${mod}` : `${mod}`
  }

  function persistLive() {
    if (character) setLive(character.id, live)
  }

  function hpStepper(delta: number) {
    if (!sheet) return
    live = {
      ...live,
      hpCurrent: Math.max(0, Math.min(sheet.hp.max, live.hpCurrent + delta)),
    }
    persistLive()
  }

  function toggleSlot(index: number) {
    if (!sheet) return
    const max = sheet.spellSlots[index] ?? 0
    const used = live.slotsUsed[index] ?? 0
    const next = (used + 1) % (max + 1)
    const slotsUsed = live.slotsUsed.slice()
    slotsUsed[index] = next
    live = { ...live, slotsUsed }
    persistLive()
  }

  function toggleCondition(c: string) {
    const has = live.conditions.includes(c)
    live = {
      ...live,
      conditions: has ? live.conditions.filter((x) => x !== c) : [...live.conditions, c],
    }
    persistLive()
  }

  function resourceNames(): string[] {
    const classes = content.value?.classes ?? []
    const out: string[] = []
    for (const ci of character?.draft.classes ?? []) {
      const k = classes.find((c) => c.id === ci.id)
      const res = k?.data?.resources
      if (Array.isArray(res)) {
        for (const r of res) {
          if (
            r &&
            typeof r === 'object' &&
            typeof r.name === 'string' &&
            typeof r.max === 'number'
          ) {
            if (!out.includes(r.name)) out.push(r.name)
          }
        }
      }
    }
    return out
  }

  function resourceMax(name: string): number {
    const classes = content.value?.classes ?? []
    for (const ci of character?.draft.classes ?? []) {
      const k = classes.find((c) => c.id === ci.id)
      const res = k?.data?.resources
      if (Array.isArray(res)) {
        for (const r of res) {
          if (r && typeof r === 'object' && r.name === name && typeof r.max === 'number')
            return r.max
        }
      }
    }
    return 0
  }

  function resourceStepper(name: string, delta: number) {
    const max = resourceMax(name)
    const cur = live.resources[name] ?? max
    live = {
      ...live,
      resources: { ...live.resources, [name]: Math.max(0, Math.min(max, cur + delta)) },
    }
    persistLive()
  }

  function classNames(): string {
    if (!character) return ''
    return character.draft.classes.map((x) => x.id).join(' + ') || '—'
  }

  function onExportJSON() {
    if (!character) return
    downloadJSON(`${character.name || 'character'}.json`, characterToJSON(character))
  }

  async function onExportPDF() {
    if (!sheetNode) return
    try {
      await exportCharacterPDF(sheetNode, `${character?.name || 'character'}.pdf`)
    } catch (err) {
      error = err instanceof Error ? err.message : String(err)
      status = 'error'
    }
  }

  function onCampaignChange(event: Event) {
    const value = (event.currentTarget as HTMLSelectElement).value
    if (character) setCampaign(character.id, value === '' ? undefined : value)
  }

  function resourceCurrent(name: string): number {
    return live.resources[name] ?? resourceMax(name)
  }
</script>

{#if !character}
  <div class="error-box">
    <p>{error}</p>
    <button class="btn" onclick={() => navigate('/characters')}>Back to list</button>
  </div>
{:else}
  <div class="page-head">
    <div>
      <h1>{character.name}</h1>
      <p class="muted">
        {classNames()} · Level {sheet?.level ?? '—'}
        {#if character.isNpc}<span class="badge">NPC</span>{/if}
      </p>
    </div>
    <div class="head-actions">
      <button class="btn" onclick={onExportJSON}>Export JSON</button>
      <button class="btn" onclick={onExportPDF} disabled={!sheet}>Export PDF</button>
      <button class="btn" onclick={() => navigate('/characters')}>Back</button>
    </div>
  </div>

  {#if status === 'loading'}
    <p>Calculating sheet…</p>
  {:else if status === 'error'}
    <div class="error-box">
      <p>{error}</p>
      <button class="btn" onclick={retry}>Try again</button>
    </div>
  {:else if sheet}
    <section class="sheet" bind:this={sheetNode}>
      <div class="stats-row">
        <div class="stat">
          <span class="stat-value">
            <button class="mini" onclick={() => hpStepper(-1)} aria-label="Decrease HP">−</button>
            {live.hpCurrent}
            <button class="mini" onclick={() => hpStepper(1)} aria-label="Increase HP">+</button>
          </span>
          <span class="stat-label">HP / {sheet.hp.max}</span>
        </div>
        <div class="stat">
          <span class="stat-value">{sheet.ac}</span><span class="stat-label">AC</span>
        </div>
        <div class="stat">
          <span class="stat-value">{modSign(sheet.proficiencyBonus)}</span><span class="stat-label"
            >Prof.</span
          >
        </div>
        <div class="stat">
          <span class="stat-value">{sheet.level}</span><span class="stat-label">Level</span>
        </div>
      </div>

      <section>
        <h2>Abilities</h2>
        <ul class="abilities">
          {#each ABILITIES as a (a)}
            <li>
              <span class="ability-name">{a}</span>
              <span class="ability-score">{sheet.abilities[a].score}</span>
              <span class="ability-mod">{modSign(sheet.abilities[a].mod)}</span>
            </li>
          {/each}
        </ul>
      </section>

      {#if sheet.spellSlots.some((n) => n > 0)}
        <section>
          <h2>Spell Slots</h2>
          <ul class="slots">
            {#each sheet.spellSlots as max, i (i)}
              {#if max > 0}
                <li>
                  <button
                    class="slot"
                    onclick={() => toggleSlot(i)}
                    aria-label={`Mark level ${i + 1} as used`}
                    title="Click to mark as used"
                  >
                    <span class="slot-lvl">{i + 1}º</span>
                    <span class="slot-count">{live.slotsUsed[i] ?? 0}/{max}</span>
                  </button>
                </li>
              {/if}
            {/each}
          </ul>
        </section>
      {/if}

      {#if resourceNames().length > 0}
        <section>
          <h2>Resources</h2>
          <ul class="resources">
            {#each resourceNames() as r (r)}
              <li class="resource">
                <span class="resource-name">{r}</span>
                <button
                  class="mini"
                  onclick={() => resourceStepper(r, -1)}
                  aria-label={`Decrease ${r}`}>−</button
                >
                <span class="resource-val">{resourceCurrent(r)}/{resourceMax(r)}</span>
                <button
                  class="mini"
                  onclick={() => resourceStepper(r, 1)}
                  aria-label={`Increase ${r}`}>+</button
                >
              </li>
            {/each}
          </ul>
        </section>
      {/if}

      <section>
        <h2>Conditions</h2>
        <div class="chips">
          {#each CONDITIONS as c (c)}
            {@const on = live.conditions.includes(c)}
            <button class:active={on} class="cond" onclick={() => toggleCondition(c)}>
              {c}
            </button>
          {/each}
        </div>
      </section>

      {#if sheet.features.length > 0}
        <section>
          <h2>Features</h2>
          <ul class="features">
            {#each sheet.features as f (f.name + f.level)}
              <li>
                <h3>{f.name} <span class="muted">— level {f.level}</span></h3>
                {#if f.description}<p>{f.description}</p>{/if}
              </li>
            {/each}
          </ul>
        </section>
      {/if}

      {#if sheet.pendingChoices.length > 0}
        <section>
          <h2>Pending Choices</h2>
          <ul class="pending">
            {#each sheet.pendingChoices as pc (pc.type + pc.description)}
              <li><strong>{pc.type}</strong> — {pc.description}</li>
            {/each}
          </ul>
        </section>
      {/if}
    </section>

    <label class="campaign-row">
      Campaign
      <select value={character.campaignId ?? ''} onchange={onCampaignChange}>
        <option value="">— none —</option>
        {#each campaigns.value as c (c.id)}
          <option value={c.id}>{c.name}</option>
        {/each}
      </select>
    </label>
  {/if}
{/if}

<style>
  .page-head {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 1rem;
    margin-bottom: 1.5rem;
  }
  h1 {
    margin: 0;
    color: var(--text-h);
  }
  .muted {
    opacity: 0.7;
  }
  .head-actions {
    display: flex;
    gap: 0.5rem;
    flex-wrap: wrap;
  }
  .stats-row {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 1rem;
    margin-bottom: 2rem;
  }
  .stat {
    background: var(--code-bg);
    border: 1px solid var(--border);
    border-radius: 8px;
    padding: 1rem;
    text-align: center;
  }
  .stat-value {
    display: block;
    font-size: 1.8rem;
    font-weight: 700;
    color: var(--text-h);
  }
  .stat-label {
    font-size: 0.8rem;
    opacity: 0.7;
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }
  .mini {
    font: inherit;
    font-size: 0.9rem;
    width: 1.6rem;
    height: 1.6rem;
    line-height: 1;
    border-radius: 999px;
    border: 1px solid var(--border);
    background: var(--bg);
    color: var(--text-h);
    cursor: pointer;
    vertical-align: middle;
  }
  h2 {
    color: var(--text-h);
    margin: 0 0 0.75rem;
  }
  .abilities {
    list-style: none;
    margin: 0 0 1.5rem;
    padding: 0;
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(7rem, 1fr));
    gap: 0.5rem;
  }
  .abilities li {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.15rem;
    background: var(--code-bg);
    border: 1px solid var(--border);
    border-radius: 8px;
    padding: 0.6rem;
  }
  .ability-name {
    font-size: 0.75rem;
    font-weight: 700;
    opacity: 0.7;
  }
  .ability-score {
    font-size: 1.4rem;
    font-weight: 700;
    color: var(--text-h);
  }
  .ability-mod {
    color: var(--accent);
  }
  .slots {
    list-style: none;
    margin: 0 0 1.5rem;
    padding: 0;
    display: flex;
    gap: 0.5rem;
    flex-wrap: wrap;
  }
  .slot {
    font: inherit;
    display: flex;
    align-items: baseline;
    gap: 0.5rem
    background: var(--code-bg)
    border: 1px solid var(--border)
    border-radius: 8px
    padding: 0.5rem 0.9rem
    cursor: pointer
    color: var(--text)
  }
  .slot:hover {
    border-color: var(--accent-border)
  }
  .slot-lvl {
    font-size: 0.75rem;
    opacity: 0.7
  }
  .slot-count {
    font-size: 1.3rem
    font-weight: 700
    color: var(--text-h)
  }
  .resources {
    list-style: none
    margin: 0 0 1.5rem
    padding: 0
    display: grid
    gap: 0.5rem
  }
  .resource {
    display: flex
    align-items: center
    gap: 0.75rem
    background: var(--code-bg)
    border: 1px solid var(--border)
    border-radius: 8px
    padding: 0.5rem 0.9rem
  }
  .resource-name {
    font-weight: 600
    color: var(--text-h)
    flex: 1
  }
  .resource-val {
    font-weight: 700
    color: var(--text-h)
  }
  .chips {
    display: flex
    flex-wrap: wrap
    gap: 0.35rem
    margin-bottom: 1.5rem
  }
  .cond {
    font: inherit
    font-size: 0.8rem
    color: var(--text)
    background: var(--code-bg)
    border: 1px solid var(--border)
    border-radius: 999px
    padding: 0.15rem 0.7rem
    cursor: pointer
  }
  .cond.active {
    color: var(--accent)
    background: var(--accent-bg)
    border-color: var(--accent-border)
  }
  .features,
  .pending {
    list-style: none
    margin: 0 0 1.5rem
    padding: 0
    display: grid
    gap: 0.75rem
  }
  .features h3 {
    margin: 0
    color: var(--text-h)
  }
  .features p {
    margin: 0.25rem 0 0
    opacity: 0.85
  }
  .campaign-row {
    display: flex
    align-items: center
    gap: 0.5rem
    margin-top: 0.5rem
    font-weight: 600
    color: var(--text-h)
  }
  select {
    font: inherit
    color: var(--text-h)
    background: var(--bg)
    border: 1px solid var(--border)
    border-radius: 8px
    padding: 0.3rem 0.6rem
  }
  .error-box {
    border: 1px solid var(--danger-border)
    background: var(--danger-bg)
    color: var(--text-h)
    border-radius: 8px
    padding: 1.25rem
    display: grid
    gap: 0.75rem
    justify-items: start
  }
  .error-box p {
    margin: 0
  }
</style>