<script lang="ts">
  // Step 1 — Class: pick one or more classes, set levels, choose subclass and
  // optional spells/feats when the class unlocks them.
  import { content } from '../../content.svelte'
  import { dataString } from '../../types'
  import type { Subclass } from '../../types'
  import {
    builder,
    deselectClass,
    recommendedSpells,
    selectClass,
    selectSubclass,
    setClassLevel,
    toggleFeat,
    toggleSpell,
  } from '../../builder.svelte'

  const contentData = $derived(content.value)

  function subclassesFor(classId: string): Subclass[] {
    return (contentData?.subclasses ?? []).filter((s) => s.classId === classId)
  }

  function levelOf(classId: string): number {
    return builder.value.draft.classes.find((c) => c.id === classId)?.level ?? 1
  }

  function isSpellcaster(): boolean {
    return (contentData?.classes ?? []).some(
      (c) => c.spellcaster && builder.value.draft.classes.some((x) => x.id === c.id)
    )
  }

  function recommendedChip(id: string): boolean {
    return recommendedSpells().includes(id)
  }

  function maxSpellLevel(): number {
    const classes = contentData?.classes ?? []
    const maxClassLevel = Math.max(
      0,
      ...builder.value.draft.classes
        .filter((c) => classes.find((x) => x.id === c.id)?.spellcaster)
        .map((c) => c.level)
    )
    return Math.min(9, Math.ceil(maxClassLevel / 2))
  }
</script>

<h2>Escolhe as classes</h2>
<p class="muted">Podes escolher mais do que uma classe (multi-classação) até nível 20.</p>

<ul class="class-grid">
  {#each contentData?.classes ?? [] as c (c.id)}
    {@const selected = builder.value.draft.classes.some((x) => x.id === c.id)}
    <li>
      <article class:selected class="card">
        <header>
          <h3>{c.name}</h3>
          <button
            class:btn={!selected}
            class:btn-danger={selected}
            onclick={() => (selected ? deselectClass(c.id) : selectClass(c.id))}
          >
            {selected ? 'Remover' : 'Escolher'}
          </button>
        </header>
        <p class="badges">
          <span class="badge">{c.hitDie}</span>
          <span class="badge">{c.spellcaster ? 'conjurador' : 'marcial'}</span>
          <span class="badge">nível {c.subclassLevel}</span>
        </p>
        <p class="muted">{dataString(c.data, 'description')}</p>
      </article>
    </li>
  {/each}
</ul>

{#each builder.value.draft.classes as cls (cls.id)}
  {@const klass = (contentData?.classes ?? []).find((c) => c.id === cls.id)}
  {#if klass}
    <section class="card config">
      <h3>{klass.name}</h3>
      <label class="level-row">
        Nível
        <button class="btn" onclick={() => setClassLevel(cls.id, Math.max(1, levelOf(cls.id) - 1))}
          >−</button
        >
        <span class="level-val">{levelOf(cls.id)}</span>
        <button class="btn" onclick={() => setClassLevel(cls.id, Math.min(20, levelOf(cls.id) + 1))}
          >+</button
        >
      </label>

      {#if levelOf(cls.id) >= klass.subclassLevel}
        {@const subs = subclassesFor(cls.id)}
        {#if subs.length > 0}
          <div>
            <p class="chips-label">Subclasse (a partir do nível {klass.subclassLevel})</p>
            <div class="chips">
              {#each subs as sub (sub.id)}
                <button
                  class:active={cls.subclassId === sub.id}
                  class="chip-btn"
                  onclick={() => selectSubclass(cls.id, sub.id)}
                >
                  {sub.name}
                </button>
              {/each}
            </div>
          </div>
        {/if}
      {/if}
    </section>
  {/if}
{/each}

{#if isSpellcaster()}
  {@const maxLvl = maxSpellLevel()}
  <section class="card">
    <h3>Feitiços <span class="muted">(até nível {maxLvl})</span></h3>
    {#if recommendedSpells().length > 0}
      <p class="chips-label">Recomendados</p>
      <div class="chips">
        {#each recommendedSpells() as id (id)}
          {@const spell = (contentData?.spells ?? []).find((s) => s.id === id)}
          <button
            class:active={(builder.value.draft.spells ?? []).includes(id)}
            class="chip-btn"
            onclick={() => toggleSpell(id)}
          >
            {spell?.name ?? id}
          </button>
        {/each}
      </div>
    {/if}
    <p class="chips-label">Todos</p>
    <div class="chips">
      {#each (contentData?.spells ?? []).filter((s) => s.level <= maxLvl) as s (s.id)}
        <button
          class:active={(builder.value.draft.spells ?? []).includes(s.id)}
          class:recommended={recommendedChip(s.id)}
          class="chip-btn"
          onclick={() => toggleSpell(s.id)}
        >
          {s.name} <span class="lvl">({s.level}º)</span>
        </button>
      {/each}
    </div>
  </section>
{/if}

{#if (contentData?.feats ?? []).length > 0}
  <section class="card">
    <h3>Talentos (feats)</h3>
    <div class="chips">
      {#each contentData?.feats ?? [] as f (f.id)}
        <button
          class:active={(builder.value.draft.feats ?? []).includes(f.id)}
          class="chip-btn"
          onclick={() => toggleFeat(f.id)}
        >
          {f.name}
        </button>
      {/each}
    </div>
  </section>
{/if}

<style>
  .muted {
    opacity: 0.7;
  }
  .class-grid {
    list-style: none;
    margin: 0;
    padding: 0;
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(14rem, 1fr));
    gap: 0.75rem;
  }
  .card {
    background: var(--code-bg);
    border: 1px solid var(--border);
    border-radius: 10px;
    padding: 1rem 1.25rem;
  }
  .card.selected {
    border-color: var(--accent-border);
    background: var(--accent-bg);
  }
  header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 0.5rem;
  }
  h3 {
    margin: 0;
    color: var(--text-h);
  }
  .badges {
    display: flex;
    flex-wrap: wrap;
    gap: 0.35rem;
    margin: 0.5rem 0;
  }
  .config {
    display: grid;
    gap: 0.75rem;
  }
  .level-row {
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
  }
  .level-val {
    font-weight: 700;
    min-width: 2rem;
    text-align: center;
    color: var(--text-h);
  }
  .chips-label {
    margin: 0.25rem 0 0.25rem;
    font-size: 0.75rem;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    opacity: 0.6;
  }
  .chips {
    display: flex;
    flex-wrap: wrap;
    gap: 0.35rem;
    margin-bottom: 0.5rem;
  }
  .chip-btn {
    font: inherit;
    font-size: 0.8rem;
    color: var(--text);
    background: var(--bg);
    border: 1px solid var(--border);
    border-radius: 999px;
    padding: 0.15rem 0.7rem;
    cursor: pointer;
  }
  .chip-btn.active {
    color: var(--accent);
    background: var(--accent-bg);
    border-color: var(--accent-border);
  }
  .chip-btn.recommended {
    border-style: dashed;
  }
  .lvl {
    opacity: 0.6;
    font-size: 0.72rem;
  }
</style>
