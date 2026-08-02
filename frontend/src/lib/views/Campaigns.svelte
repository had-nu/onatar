<script lang="ts">
  // RF-05 minimal campaigns: create/list/delete (id + name), linked via
  // `characters.campaignId`.
  import { createCampaign, deleteCampaign, listCampaigns } from '../campaigns.svelte'
  import { characters } from '../characters.svelte'

  let newName = $state('')

  function add() {
    createCampaign(newName)
    newName = ''
  }

  function membersOf(campaignId: string): number {
    return characters.value.filter((c) => c.campaignId === campaignId).length
  }
</script>

<div class="page-head">
  <h1>Campaigns</h1>
  <p class="muted">Group your characters by campaign (minimal structure).</p>
</div>

<form class="create" onsubmit={(e) => e.preventDefault()}>
  <input bind:value={newName} placeholder="Campaign name" />
  <button class="btn primary" onclick={add}>Create Campaign</button>
</form>

{#if listCampaigns().length === 0}
  <p class="muted empty">No campaigns yet. Create the first one above.</p>
{:else}
  <ul class="grid">
    {#each listCampaigns() as c (c.id)}
      <li>
        <article class="card">
          <header>
            <h2>{c.name}</h2>
            <button
              class="btn danger"
              onclick={() => {
                if (confirm(`Delete campaign "${c.name}"?`)) deleteCampaign(c.id)
              }}
              aria-label={`Delete ${c.name}`}
            >
              Delete
            </button>
          </header>
          <p class="muted">{membersOf(c.id)} characters</p>
        </article>
      </li>
    {/each}
  </ul>
{/if}

<p class="muted back"><a href="#/characters">← Back to characters</a></p>

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
  .create {
    display: flex;
    gap: 0.5rem;
    margin-bottom: 1.5rem;
  }
  input {
    font: inherit;
    color: var(--text-h);
    background: var(--bg);
    border: 1px solid var(--border);
    border-radius: 8px;
    padding: 0.45rem 0.75rem;
    flex: 1;
  }
  .empty {
    margin-bottom: 1rem;
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
    background: var(--code-bg);
    border: 1px solid var(--border);
    border-radius: 10px;
    padding: 1rem 1.25rem;
  }
  header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 0.5rem;
  }
  h2 {
    margin: 0;
    color: var(--text-h);
  }
  .back {
    margin-top: 1.5rem;
  }
</style>