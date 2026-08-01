<script lang="ts">
  import { onMount } from 'svelte'
  import { cachedContent } from '../content.svelte'

  let stats = $state<{ classes: number; species: number; backgrounds: number } | null>(null)

  onMount(() => {
    const c = cachedContent()
    if (c) {
      stats = {
        classes: c.classes.length,
        species: c.species.length,
        backgrounds: c.backgrounds.length,
      }
    }
  })
</script>

<section class="hero">
  <h1>Onatar</h1>
  <p class="tagline">Criação de fichas de personagem D&amp;D 2024 (5.5e)</p>
  <p class="muted">
    Inspirado no D&D Beyond, focado em acessibilidade para novos jogadores e simplicidade para
    mestres.
  </p>
  <div class="cta">
    <a class="btn primary" href="#/builder">Criar personagem</a>
    <a class="btn" href="#/characters">Os meus personagens</a>
    <a class="btn" href="#/content">Explorar conteúdo</a>
  </div>
  {#if stats}
    <p class="stats muted">
      {stats.classes} classes · {stats.species} espécies · {stats.backgrounds} backgrounds
    </p>
  {/if}
</section>

<style>
  .hero {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    min-height: 70vh;
    gap: 0.5rem;
    text-align: center;
    padding: 2rem 0;
  }
  h1 {
    font-size: 3.5rem;
    margin: 0;
    color: var(--text-h);
  }
  .tagline {
    font-size: 1.3rem;
    color: var(--text-h);
    margin: 0;
  }
  .muted {
    opacity: 0.7;
    max-width: 34rem;
  }
  .cta {
    display: flex;
    gap: 0.75rem;
    margin-top: 1.5rem;
    flex-wrap: wrap;
    justify-content: center;
  }
  .stats {
    margin-top: 1.5rem;
  }
</style>
