<script lang="ts">
  import { onMount } from 'svelte'
  import Shell from './lib/Shell.svelte'
  import { initRouter, route } from './lib/router.svelte'
  import { initTheme } from './lib/theme.svelte'
  import Landing from './lib/views/Landing.svelte'
  import Characters from './lib/views/Characters.svelte'
  import CharacterView from './lib/views/CharacterView.svelte'
  import Content from './lib/views/Content.svelte'
  import Builder from './lib/views/Builder.svelte'
  import Campaigns from './lib/views/Campaigns.svelte'

  onMount(() => {
    initTheme()
    return initRouter()
  })
</script>

<Shell>
  {#if route.name === 'home'}
    <Landing />
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
