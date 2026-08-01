<script lang="ts">
  import { route } from './router.svelte'
  import { cycleTheme, theme } from './theme.svelte'
  import type { RouteName } from './router.svelte'

  interface NavLink {
    path: string
    label: string
    match: RouteName[]
  }

  const links: NavLink[] = [
    { path: '/', label: 'Início', match: ['home'] },
    { path: '/characters', label: 'Personagens', match: ['characters', 'character'] },
    { path: '/content', label: 'Conteúdo', match: ['content'] },
    { path: '/campaigns', label: 'Campanhas', match: ['campaigns'] },
  ]

  function isActive(match: RouteName[]): boolean {
    return match.includes(route.name)
  }

  function themeLabel(): string {
    return theme.value === 'system' ? 'Sistema' : theme.value === 'light' ? 'Claro' : 'Escuro'
  }
</script>

<header class="topbar">
  <a class="brand" href="#/">Onatar</a>
  <nav aria-label="Navegação principal">
    {#each links as link (link.path)}
      <a href={link.path} class:active={isActive(link.match)}>{link.label}</a>
    {/each}
  </nav>
  <button class="theme-toggle" onclick={cycleTheme} aria-label="Alternar tema">
    {themeLabel()}
  </button>
</header>

<main class="view">
  <slot />
</main>

<style>
  .topbar {
    position: sticky;
    top: 0;
    z-index: 10;
    display: flex;
    align-items: center;
    gap: 1.5rem;
    padding: 0.75rem 1.5rem;
    background: var(--bg);
    border-bottom: 1px solid var(--border);
  }
  .brand {
    font-weight: 700;
    font-size: 1.1rem;
    color: var(--text-h);
    text-decoration: none;
  }
  nav {
    display: flex;
    gap: 1rem;
    flex: 1;
  }
  nav a {
    color: var(--text);
    text-decoration: none;
    padding: 0.25rem 0.5rem;
    border-radius: 6px;
  }
  nav a:hover {
    color: var(--text-h);
  }
  nav a.active {
    color: var(--accent);
    background: var(--accent-bg);
  }
  .theme-toggle {
    font: inherit;
    color: var(--text-h);
    background: var(--code-bg);
    border: 1px solid var(--border);
    border-radius: 999px;
    padding: 0.3rem 0.9rem;
    cursor: pointer;
  }
  .theme-toggle:hover {
    border-color: var(--accent-border);
  }
  .view {
    max-width: 56rem;
    margin: 0 auto;
    padding: 2rem 1.5rem;
  }
</style>
