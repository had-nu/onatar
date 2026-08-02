<script lang="ts">
  import { route } from './router.svelte'
  import { cycleTheme, theme } from './theme.svelte'
  import { isAuthenticated, getUser, logout, isLoading } from './auth.svelte'
  import { navigate } from './router.svelte'
  import type { RouteName } from './router.svelte'

  interface NavLink {
    path: string
    label: string
    match: RouteName[]
    auth?: boolean
  }

  const links: NavLink[] = [
    { path: '/', label: 'Home', match: ['home'] },
    { path: '/characters', label: 'Characters', match: ['characters', 'character'], auth: true },
    { path: '/content', label: 'Content', match: ['content'] },
    { path: '/campaigns', label: 'Campaigns', match: ['campaigns'], auth: true },
    { path: '/import', label: 'Import', match: ['import'], auth: true },
    { path: '/combat', label: 'Combat', match: ['combat'], auth: true },
  ]

  function isActive(match: RouteName[]): boolean {
    return match.includes(route.name)
  }

  function themeLabel(): string {
    return theme.value === 'system' ? 'System' : theme.value === 'light' ? 'Light' : 'Dark'
  }

  function userInitials(name: string | null, login: string): string {
    if (name) {
      const parts = name.trim().split(/\s+/)
      if (parts.length >= 2) return (parts[0][0] + parts[1][0]).toUpperCase()
      return parts[0].slice(0, 2).toUpperCase()
    }
    return login.slice(0, 2).toUpperCase()
  }

  async function handleLogout() {
    await logout()
    navigate('/')
  }
</script>

<header class="topbar">
  <a class="brand" href="#/">Onatar</a>
  <nav aria-label="Main navigation">
    {#each links as link (link.path)}
      {#if !link.auth || isAuthenticated()}
        <a href={link.path} class:active={isActive(link.match)}>{link.label}</a>
      {/if}
    {/each}
  </nav>
  <div class="topbar-actions">
    <button class="theme-toggle" onclick={cycleTheme} aria-label="Toggle theme">
      {themeLabel()}
    </button>
    {#if isLoading()}
      <div class="avatar-loading" aria-busy="true" aria-label="Loading user…"></div>
    {:else if isAuthenticated()}
      <div class="user-menu">
        <button class="avatar-btn" aria-expanded="false" aria-haspopup="true" aria-label="User menu">
          <span class="avatar" style="background-image: url('{getUser()?.avatar_url || `https://ui-avatars.com/api/?name=${encodeURIComponent(userInitials(getUser()?.name, getUser()?.login))}&background=8b0000&color=fff`}')"></span>
        </button>
        <div class="user-dropdown" role="menu">
          <div class="user-info">
            <span class="user-name">{getUser()?.name || getUser()?.login}</span>
            <span class="user-login">@{getUser()?.login}</span>
          </div>
          <hr class="dropdown-divider" />
          <button class="dropdown-item" role="menuitem" onclick={handleLogout}>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="16" height="16" aria-hidden="true">
              <path d="M9 21H5a2 2 0 01-2-2V5a2 2 0 012-2h4"/>
              <polyline points="16 17 21 12 16 7"/>
              <line x1="21" y1="12" x2="9" y2="12"/>
            </svg>
            Sign out
          </button>
        </div>
      </div>
    {:else}
      <a class="btn primary" href="#/login">Sign in</a>
    {/if}
  </div>
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
  .topbar-actions {
    display: flex;
    align-items: center;
    gap: 0.75rem;
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
  .avatar-loading {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    background: linear-gradient(90deg, var(--code-bg) 25%, var(--border) 50%, var(--code-bg) 75%);
    background-size: 200% 100%;
    animation: shimmer 1.5s infinite;
  }
  @keyframes shimmer {
    to { background-position: -200% 0; }
  }
  .user-menu {
    position: relative;
  }
  .avatar-btn {
    padding: 0;
    border: none;
    background: none;
    cursor: pointer;
  }
  .avatar {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    background-size: cover;
    background-position: center;
    border: 2px solid var(--accent-border);
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--accent);
    font-weight: 700;
    font-size: 0.75rem;
  }
  .avatar[style*="ui-avatars"] {
    background-image: var(--avatar-url);
  }
  .user-dropdown {
    position: absolute;
    top: calc(100% + 0.5rem);
    right: 0;
    min-width: 180px;
    background: var(--card-bg);
    border: 1px solid var(--border);
    border-radius: 8px;
    box-shadow: 0 8px 24px rgba(0,0,0,0.3);
    overflow: hidden;
    z-index: 20;
  }
  .user-info {
    padding: 0.75rem 1rem;
    display: flex;
    flex-direction: column;
    gap: 0.15rem;
  }
  .user-name {
    font-weight: 600;
    color: var(--text-h);
    font-size: 0.9rem;
  }
  .user-login {
    font-size: 0.75rem;
    color: var(--text);
    opacity: 0.7;
  }
  .dropdown-divider {
    border: none;
    border-top: 1px solid var(--border);
    margin: 0
  }
  .dropdown-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    width: 100%;
    padding: 0.5rem 1rem;
    border: none;
    background: none;
    color: var(--text);
    font: inherit;
    text-align: left;
    cursor: pointer;
  }
  .dropdown-item:hover {
    background: var(--accent-bg);
    color: var(--accent);
  }
  .dropdown-item svg {
    flex-shrink: 0
  }
  .view {
    max-width: 56rem;
    margin: 0 auto;
    padding: 2rem 1.5rem;
  }
</style>