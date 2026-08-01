<script lang="ts">
  // Step 6 — Review: name the character and confirm the build before saving.
  import { content } from '../../content.svelte'
  import { ABILITIES } from '../../types'
  import { builder } from '../../builder.svelte'

  const contentData = $derived(content.value)

  function classNames(): string {
    const classes = contentData?.classes ?? []
    return builder.value.draft.classes
      .map((c) => {
        const k = classes.find((x) => x.id === c.id)
        return `${k?.name ?? c.id} ${c.level}`
      })
      .join(' + ')
  }

  function nameOf(id: string | undefined, list: Array<{ id: string; name: string }>): string {
    if (!id) return '—'
    return list.find((x) => x.id === id)?.name ?? id
  }

  function spells(): string[] {
    const spells = contentData?.spells ?? []
    return (builder.value.draft.spells ?? []).map(
      (id) => spells.find((s) => s.id === id)?.name ?? id
    )
  }
</script>

<h2>Revisão</h2>
<p class="muted">Confere os dados e guarda o personagem.</p>

<label class="name-row">
  Nome
  <input bind:value={builder.value.name} placeholder="Ex.: Onatar, o Eterno" />
</label>

<dl class="summary">
  <div>
    <dt>Classes</dt>
    <dd>{classNames() || '—'}</dd>
  </div>
  <div>
    <dt>Background</dt>
    <dd>{nameOf(builder.value.draft.backgroundId, contentData?.backgrounds ?? [])}</dd>
  </div>
  <div>
    <dt>Espécie</dt>
    <dd>{nameOf(builder.value.draft.speciesId, contentData?.species ?? [])}</dd>
  </div>
  <div>
    <dt>Atributos</dt>
    <dd>
      {#each ABILITIES as ab (ab)}
        {ab}
        {builder.value.draft.abilityScores[ab]}{#if ab !== 'CHA'},
        {/if}
      {/each}
    </dd>
  </div>
  {#if (builder.value.draft.spells ?? []).length > 0}
    <div>
      <dt>Feitiços</dt>
      <dd>{spells().join(', ')}</dd>
    </div>
  {/if}
  {#if (builder.value.draft.feats ?? []).length > 0}
    <div>
      <dt>Talentos</dt>
      <dd>{builder.value.draft.feats?.join(', ')}</dd>
    </div>
  {/if}
  {#if builder.value.equipment.length > 0}
    <div>
      <dt>Equipamento</dt>
      <dd>{builder.value.equipment.join(', ')}</dd>
    </div>
  {/if}
</dl>

<style>
  .muted {
    opacity: 0.7;
  }
  .name-row {
    display: grid;
    gap: 0.35rem;
    margin-bottom: 1rem;
    font-weight: 600;
    color: var(--text-h);
  }
  input {
    font: inherit;
    color: var(--text-h);
    background: var(--bg);
    border: 1px solid var(--border);
    border-radius: 8px;
    padding: 0.5rem 0.75rem;
  }
  .summary {
    display: grid;
    gap: 0.5rem;
    margin: 0;
  }
  .summary > div {
    display: grid;
    grid-template-columns: 10rem 1fr;
    gap: 0.5rem;
    padding: 0.5rem 0;
    border-bottom: 1px solid var(--border);
  }
  dt {
    font-weight: 600;
    color: var(--text-h);
  }
  dd {
    margin: 0;
  }
</style>
