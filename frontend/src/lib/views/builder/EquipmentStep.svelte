<script lang="ts">
  // Step 5 — Equipment: starting gear offered by the chosen class/background.
  import { content } from '../../content.svelte'
  import { builder, toggleEquipment } from '../../builder.svelte'

  const contentData = $derived(content.value)

  function gearList(): string[] {
    const classes = contentData?.classes ?? []
    const backgrounds = contentData?.backgrounds ?? []
    const klass = classes.find((c) => builder.value.draft.classes.some((x) => x.id === c.id))
    const bg = backgrounds.find((b) => b.id === builder.value.draft.backgroundId)
    const items: string[] = []
    for (const list of [klass?.data?.startingGear, bg?.data?.startingGear]) {
      if (Array.isArray(list)) {
        for (const it of list) {
          if (typeof it === 'string' && !items.includes(it)) items.push(it)
        }
      }
    }
    return items
  }

  const gear = $derived(gearList())
</script>

<h2>Equipamento</h2>
<p class="muted">Equipamento inicial sugerido pela classe e pelo background.</p>

{#if gear.length === 0}
  <p class="muted">Sem dados de equipamento disponíveis para esta combinação.</p>
{:else}
  <div class="chips">
    {#each gear as item (item)}
      {@const on = builder.value.equipment.includes(item)}
      <button class:active={on} class="chip-btn" onclick={() => toggleEquipment(item)}>
        {on ? '✓ ' : ''}{item}
      </button>
    {/each}
  </div>
{/if}

<style>
  .muted {
    opacity: 0.7;
  }
  .chips {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
  }
  .chip-btn {
    font: inherit;
    font-size: 0.85rem;
    color: var(--text);
    background: var(--code-bg);
    border: 1px solid var(--border);
    border-radius: 999px;
    padding: 0.3rem 0.9rem;
    cursor: pointer;
  }
  .chip-btn:hover {
    border-color: var(--accent-border);
  }
  .chip-btn.active {
    color: var(--accent);
    background: var(--accent-bg);
    border-color: var(--accent-border);
  }
</style>
