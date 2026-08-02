<script lang="ts">
  import { onMount } from 'svelte'
  import Shell from './lib/Shell.svelte'
  import { initRouter, route } from './lib/router.svelte'
  import { initTheme } from './lib/theme.svelte'
  import { checkAuth } from './lib/auth.svelte'
  import { loadFromApi as loadCharacters, migrateLocalToApi as migrateCharacters } from './lib/characters.svelte'
  import { loadFromApi as loadCampaigns, migrateLocalToApi as migrateCampaigns } from './lib/campaigns.svelte'
  import Landing from './lib/views/Landing.svelte'
  import Login from './lib/views/Login.svelte'
  import Characters from './lib/views/Characters.svelte'
  import CharacterView from './lib/views/CharacterView.svelte'
  import Content from './lib/views/Content.svelte'
  import Builder from './lib/views/Builder.svelte'
  import Campaigns from './lib/views/Campaigns.svelte'
  import Import from './lib/views/Import.svelte'
  import Combat from './lib/views/Combat.svelte'

  onMount(async () => {
    initTheme()
    const cleanupRouter = initRouter()

    const user = await checkAuth()
    if (user) {
      // Authenticated: load from API and migrate any local data
      await Promise.all([loadCharacters(), loadCampaigns()])
      await Promise.all([migrateCharacters(), migrateCampaigns()])
    } else {
      // Guest: data already loaded from localStorage via box initialization
    }

    return cleanupRouter
  })
</script>

<Shell>
  {#if route.name === 'home'}
    <Landing />
  {:else if route.name === 'login'}
    <Login />
  {:else if route.name === 'characters'}
    <Characters />
  {:else if route.name === 'character'}
    <CharacterView id={route.params.id} />
  {:else if route.name === 'content'}
    <Content />
  {:else if route.name === 'builder'}
    <Builder />
  {:else if route.name === 'campaigns'}
    <Campaigns />
  {:else if route.name === 'import'}
    <Import />
  {:else if route.name === 'combat'}
    <Combat />
  {:else}
    <div class="notfound">
      <h1>404</h1>
      <p>Página não encontrada.</p>
      <a class="btn" href="#/">Voltar ao início</a>
    </div>
  {/if}
</Shell>

<style>
  .notfound {
    text-align: center;
    padding: 3rem 0;
  }
</style>