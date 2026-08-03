<script lang="ts">
  import { fade, scale } from 'svelte/transition';

  interface Props {
    open: boolean;
    title: string;
    subtitle?: string;
    description: string;
    meta?: { label: string; value: string }[];
    tags?: string[];
    features?: string[];
    color?: string;
    onClose: () => void;
  }

  let { open, title, subtitle, description, meta = [], tags = [], features = [], color = 'var(--on-red)', onClose }: Props = $props();
</script>

{#if open}
  <!-- svelte-ignore a11y_click_events_have_key_events a11y_no_static_element_interactions -->
  <div class="popup-overlay" transition:fade={{ duration: 150 }} onclick={onClose}>
    <div
      class="popup-card"
      transition:scale={{ start: 0.94, duration: 200 }}
      onclick={(e) => e.stopPropagation()}
      style="--accent: {color}"
    >
      <button class="popup-close" onclick={onClose} aria-label="Fechar">×</button>

      <div class="popup-header">
        <div class="popup-accent-bar" style="background: {color}"></div>
        <h2 class="popup-title">{title}</h2>
        {#if subtitle}
          <p class="popup-subtitle">{subtitle}</p>
        {/if}
      </div>

      <div class="popup-body">
        <p class="popup-description">{description}</p>

        {#if meta.length > 0}
          <div class="popup-meta">
            {#each meta as m}
              <div class="meta-item">
                <span class="meta-label">{m.label}</span>
                <span class="meta-value">{m.value}</span>
              </div>
            {/each}
          </div>
        {/if}

        {#if tags.length > 0}
          <div class="popup-tags">
            {#each tags as tag}
              <span class="tag" style="border-color: {color}; color: {color}">{tag}</span>
            {/each}
          </div>
        {/if}

        {#if features.length > 0}
          <div class="popup-features">
            <h3>Características</h3>
            <ul>
              {#each features as feat}
                <li>{feat}</li>
              {/each}
            </ul>
          </div>
        {/if}
      </div>
    </div>
  </div>
{/if}

<style>
  .popup-overlay {
    position: fixed;
    inset: 0;
    z-index: 1000;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(0, 0, 0, 0.55);
    backdrop-filter: blur(4px);
    padding: 1rem;
  }

  .popup-card {
    position: relative;
    width: 100%;
    max-width: 500px;
    max-height: 85vh;
    overflow-y: auto;
    background: var(--on-bg-surface);
    border: 1px solid var(--on-border);
    border-radius: 14px;
    box-shadow: 0 24px 64px rgba(0,0,0,0.5), 0 0 0 1px rgba(255,255,255,0.03);
    color: var(--on-text);
  }

  .popup-close {
    position: absolute;
    top: 12px;
    right: 12px;
    width: 30px;
    height: 30px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: transparent;
    border: none;
    border-radius: 6px;
    color: var(--on-text-dim);
    font-size: 1.3rem;
    line-height: 1;
    cursor: pointer;
    transition: all 0.15s;
  }

  .popup-close:hover {
    background: var(--on-bg-hover);
    color: var(--on-text);
  }

  .popup-header {
    padding: 22px 22px 0;
    position: relative;
  }

  .popup-accent-bar {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 3px;
    border-radius: 14px 14px 0 0;
  }

  .popup-title {
    margin: 0;
    font-family: var(--font-display);
    font-size: 1.35rem;
    font-weight: 700;
    letter-spacing: -0.02em;
    color: var(--on-text);
  }

  .popup-subtitle {
    margin: 6px 0 0;
    font-size: 0.85rem;
    color: var(--on-text-muted);
    font-weight: 500;
  }

  .popup-body {
    padding: 14px 22px 22px;
  }

  .popup-description {
    margin: 0 0 16px;
    font-size: 0.9rem;
    line-height: 1.6;
    color: var(--on-text-muted);
  }

  .popup-meta {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(130px, 1fr));
    gap: 8px;
    margin-bottom: 16px;
  }

  .meta-item {
    background: var(--on-bg-root);
    border-radius: 8px;
    padding: 10px 12px;
    display: flex;
    flex-direction: column;
    gap: 2px;
    border: 1px solid var(--on-border);
  }

  .meta-label {
    font-size: 0.65rem;
    text-transform: uppercase;
    letter-spacing: 0.06em;
    color: var(--on-text-dim);
    font-weight: 600;
  }

  .meta-value {
    font-size: 0.8rem;
    font-weight: 600;
    color: var(--on-text);
  }

  .popup-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    margin-bottom: 16px;
  }

  .tag {
    font-size: 0.7rem;
    font-weight: 600;
    padding: 4px 10px;
    border-radius: 20px;
    border: 1.5px solid;
    background: transparent;
    opacity: 0.85;
  }

  .popup-features h3 {
    margin: 0 0 10px;
    font-size: 0.75rem;
    text-transform: uppercase;
    letter-spacing: 0.08em;
    color: var(--on-text-dim);
  }

  .popup-features ul {
    margin: 0;
    padding-left: 18px;
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .popup-features li {
    font-size: 0.85rem;
    line-height: 1.5;
    color: var(--on-text-muted);
  }

  .popup-features li::marker {
    color: var(--accent, var(--on-red));
  }

  @media (max-width: 480px) {
    .popup-card {
      max-height: 90vh;
      border-radius: 10px;
    }
    .popup-header, .popup-body {
      padding-left: 16px;
      padding-right: 16px;
    }
  }
</style>
