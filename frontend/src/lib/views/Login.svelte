<script lang="ts">
  import { onMount } from 'svelte'
  import { navigate } from '../router.svelte'
  import { loginWithGitHub, checkAuth, isLoading, getAuthError } from '../auth.svelte'

  let checking = $derived(isLoading())
  let error = $derived(getAuthError())

  onMount(async () => {
    const user = await checkAuth()
    if (user) {
      navigate('/characters')
    }
  })
</script>

<div class="page-head">
  <h1>Entrar no Onatar</h1>
  <p class="muted">Faça login com GitHub para sincronizar personagens na nuvem.</p>
</div>

{#if checking}
  <div class="loading">
    <div class="spinner" aria-busy="true" aria-label="A verificar sessão…"></div>
    <p>A verificar sessão…</p>
  </div>
{:else}
  {#if error}
    <div class="error-box" role="alert">
      <p>{error}</p>
      <button class="btn" onclick={loginWithGitHub}>Tentar novamente</button>
    </div>
  {:else}
    <div class="login-card">
      <div class="brand-icon" aria-hidden="true">
        <svg viewBox="0 0 24 24" fill="currentColor" width="48" height="48">
          <path d="M12 0C5.374 0 0 5.373 0 12c0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23A11.509 11.509 0 0112 5.803c1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576C20.566 21.797 24 17.3 24 12c0-6.627-5.373-12-12-12z"/>
        </svg>
      </div>
      <h2>Login com GitHub</h2>
      <p class="muted">Seguro, sem passwords para gerir, e os seus dados ficam consigo.</p>
      <button class="btn primary github-btn" onclick={loginWithGitHub} disabled={checking}>
        <svg viewBox="0 0 24 24" fill="currentColor" width="20" height="20" aria-hidden="true">
          <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23A11.509 11.509 0 0112 5.803c1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576C20.566 21.797 24 17.3 24 12c0-6.627-5.373-12-12-12z"/>
        </svg>
        Continuar com GitHub
      </button>
      <p class="guest-link">
        <a href="#/characters" onclick={(e) => { e.preventDefault(); navigate('/characters'); }}>
          Continuar como convidado (localStorage)
        </a>
      </p>
    </div>
  {/if}
{/if}

<style>
  .page-head {
    text-align: center;
    margin-bottom: 2rem;
  }
  h1 {
    margin: 0 0 0.5rem;
    color: var(--text-h);
  }
  .muted {
    opacity: 0.7;
  }
  .loading {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem;
    padding: 3rem;
  }
  .spinner {
    width: 2.5rem;
    height: 2.5rem;
    border: 3px solid var(--border);
    border-top-color: var(--accent);
    border-radius: 50%;
    animation: spin 1s linear infinite;
  }
  @keyframes spin {
    to { transform: rotate(360deg); }
  }
  .error-box {
    border: 1px solid var(--danger-border);
    background: var(--danger-bg);
    border-radius: 8px;
    padding: 1.25rem;
    display: grid;
    gap: 0.75rem;
    justify-items: start;
    max-width: 28rem;
    margin: 0 auto;
  }
  .login-card {
    max-width: 28rem;
    margin: 0 auto;
    padding: 2.5rem;
    background: var(--card-bg);
    border: 1px solid var(--border);
    border-radius: 16px;
    text-align: center;
    display: grid;
    gap: 1.25rem;
  }
  .brand-icon {
    color: var(--accent);
  }
  h2 {
    margin: 0;
    color: var(--text-h);
  }
  .github-btn {
    display: inline-flex;
    align-items: center;
    gap: 0.6rem;
    font: inherit;
    font-weight: 600;
    padding: 0.75rem 1.5rem;
    border-radius: 999px;
    border: 1px solid var(--border);
    background: var(--bg);
    color: var(--text-h);
    cursor: pointer;
    transition: border-color 0.15s, background 0.15s;
  }
  .github-btn:hover:not(:disabled) {
    border-color: var(--accent-border);
    background: var(--accent-bg);
  }
  .github-btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
  .guest-link {
    margin-top: 0.5rem;
  }
  .guest-link a {
    color: var(--accent);
    text-decoration: none;
    font-size: 0.9rem;
  }
  .guest-link a:hover {
    text-decoration: underline;
  }
</style>