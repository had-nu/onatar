<script lang="ts">
  import { onMount } from 'svelte';
  import type { ClassEntry, SubClassEntry, SpellEntry, SpeciesEntry, BackgroundEntry } from '../types';

  type PopupType = 'class' | 'subclass' | 'spell' | 'species' | 'background';
  
  interface Props {
    type: PopupType;
    data: ClassEntry | SubClassEntry | SpellEntry | SpeciesEntry | BackgroundEntry;
    isOpen: boolean;
    onClose: () => void;
  }

  let { type, data, isOpen, onClose }: Props = $props();

  // Close on Escape key
  onMount(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape' && isOpen) {
        onClose();
      }
    };
    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  });

  function getAccentColor(): string {
    switch (type) {
      case 'class': return 'var(--on-red)';
      case 'subclass': return 'var(--on-gold)';
      case 'spell': return 'var(--on-blue)';
      case 'species': return 'var(--on-green)';
      case 'background': return 'var(--on-red)';
      default: return 'var(--on-red)';
    }
  }

  function getAccentBg(): string {
    switch (type) {
      case 'class': return 'var(--on-red-glow)';
      case 'subclass': return 'rgba(201, 169, 78, 0.15)';
      case 'spell': return 'rgba(91, 141, 239, 0.15)';
      case 'species': return 'rgba(0, 184, 122, 0.15)';
      case 'background': return 'var(--on-red-glow)';
      default: return 'var(--on-red-glow)';
    }
  }

  function getAccentBorder(): string {
    switch (type) {
      case 'class': return 'var(--on-red)';
      case 'subclass': return 'var(--on-gold)';
      case 'spell': return 'var(--on-blue)';
      case 'species': return 'var(--on-green)';
      case 'background': return 'var(--on-red)';
      default: return 'var(--on-red)';
    }
  }
</script>

{#if isOpen}
  <div 
    class="info-popup-overlay" 
    on:click={onClose}
    style="--accent-color: {getAccentColor()}; --accent-bg: {getAccentBg()}; --accent-border: {getAccentBorder()};"
  >
    <div class="info-popup" on:click={(e) => e.stopPropagation()}>
      <div class="popup-header">
        <div class="popup-title-row">
          <div class="popup-icon" style="background: var(--accent-bg); border-color: var(--accent-border);">
            {getIcon()}
          </div>
          <div class="popup-title-group">
            <h3 class="popup-title">{getTitle()}</h3>
            <p class="popup-subtitle">{getSubtitle()}</p>
          </div>
        </div>
        <button class="popup-close" on:click={onClose} aria-label="Close">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="6" x2="6" y2="18"/>
            <line x1="6" y1="6" x2="18" y2="18"/>
          </svg>
        </button>
      </div>
      
      <div class="popup-content">
        <p class="popup-description">{getDescription()}</p>
        
        <div class="popup-meta">
          {#each getMetaFields() as field}
            <div class="meta-field">
              <span class="meta-label">{field.label}</span>
              <span class="meta-value">{field.value}</span>
            </div>
          {/each}
        </div>

        {#if getTags().length > 0}
          <div class="popup-tags">
            {#each getTags() as tag}
              <span class="popup-tag" style="--tag-color: {tag.color}">{tag.label}</span>
            {/each}
          </div>
        {/if}

        {#if getFeatures().length > 0}
          <div class="popup-features">
            <h4>Features</h4>
            <ul class="features-list">
              {#each getFeatures() as feature}
                <li class="feature-item">
                  <span class="feature-name">{feature.name}</span>
                  {#if feature.level !== undefined}
                    <span class="feature-level">Level {feature.level}</span>
                  {/if}
                  {#if feature.description}
                    <p class="feature-desc">{feature.description}</p>
                  {/if}
                </li>
              {/each}
            </ul>
          </div>
        {/if}

        {#if getSpells().length > 0}
          <div class="popup-spells">
            <h4>Spells</h4>
            <div class="spells-grid">
              {#each getSpells() as spell}
                <div class="spell-card">
                  <span class="spell-name">{spell.name}</span>
                  <span class="spell-level">Level {spell.level}</span>
                  <span class="spell-school">{spell.school}</span>
                </div>
              {/each}
            </div>
          </div>
        {/if}

        {#if getVariants().length > 0}
          <div class="popup-variants">
            <h4>Variants</h4>
            <ul class="variants-list">
              {#each getVariants() as variant}
                <li class="variant-item">
                  <span class="variant-name">{variant.name}</span>
                  {#if variant.description}
                    <p class="variant-desc">{variant.description}</p>
                  {/if}
                </li>
              {/each}
            </ul>
          </div>
        {/if}
      </div>
    </div>
  </div>
{/if}

<style>
  .info-popup-overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.6);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
    padding: 1rem;
    animation: fadeIn 0.2s ease;
  }

  @keyframes fadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
  }

  .info-popup {
    background: var(--on-bg-surface);
    border: 1px solid var(--accent-border);
    border-radius: 16px;
    max-width: 480px;
    width: 100%;
    max-height: 85vh;
    overflow-y: auto;
    background: linear-gradient(135deg, var(--on-bg-surface) 0%, var(--on-bg-elevated) 100%);
    box-shadow: 0 8px 32px var(--accent-bg);
    animation: slideUp 0.3s ease;
  }

  @keyframes slideUp {
    from { opacity: 0; transform: translateY(20px); }
    to { opacity: 1; transform: translateY(0); }
  }

  .popup-header {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    padding: 1.25rem 1.5rem 1rem;
    border-bottom: 1px solid var(--on-border);
    flex-wrap: wrap;
    gap: 1rem;
  }

  .popup-title-row {
    display: flex;
    align-items: flex-start;
    gap: 1rem;
    flex: 1;
  }

  .popup-icon {
    width: 48px;
    height: 48px;
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 20px;
    color: var(--accent-color);
    flex-shrink: 0;
  }

  .popup-title-group {
    min-width: 0;
  }

  .popup-title {
    font-size: 18px;
    font-weight: 600;
    color: var(--on-text);
    margin: 0 0 4px 0;
    line-height: 1.2;
  }

  .popup-subtitle {
    font-size: 12px;
    color: var(--on-text-dim);
    margin: 0;
    text-transform: uppercase;
    letter-spacing: 1px;
  }

  .popup-close {
    width: 32px;
    height: 32px;
    border-radius: 8px;
    border: 1px solid var(--on-border);
    background: transparent;
    color: var(--on-text-dim);
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 16px;
    transition: all 0.15s;
    flex-shrink: 0;
  }

  .popup-close:hover {
    background: var(--on-bg-hover);
    color: var(--on-text);
    border-color: var(--accent-border);
  }

  .popup-content {
    padding: 1.25rem 1.5rem;
    display: flex;
    flex-direction: column;
    gap: 1.25rem;
  }

  .popup-description {
    font-size: 13px;
    color: var(--on-text-dim);
    line-height: 1.6;
    margin: 0;
  }

  .popup-meta {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
    gap: 0.75rem;
  }

  .meta-field {
    display: flex;
    flex-direction: column;
    gap: 2px;
  }

  .meta-label {
    font-size: 9px;
    text-transform: uppercase;
    letter-spacing: 1px;
    color: var(--on-text-dim);
    font-weight: 500;
  }

  .meta-value {
    font-size: 12px;
    color: var(--on-text);
    font-weight: 500;
  }

  .popup-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
  }

  .popup-tag {
    font-size: 10px;
    font-weight: 500;
    padding: 4px 10px;
    border-radius: 999px;
    background: var(--accent-bg);
    border: 1px solid var(--accent-border);
    color: var(--tag-color, var(--accent-color));
    text-transform: uppercase;
    letter-spacing: 0.5px;
    font-size: 9px;
  }

  .popup-features h4,
  .popup-spells h4,
  .popup-variants h4 {
    font-size: 11px;
    text-transform: uppercase;
    letter-spacing: 1.5px;
    color: var(--on-text-dim);
    margin: 0 0 8px 0;
  }

  .features-list,
  .variants-list {
    list-style: none;
    padding: 0;
    margin: 0;
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .feature-item {
    background: var(--on-bg-root);
    border: 1px solid var(--on-border);
    border-radius: 8px;
    padding: 12px;
  }

  .feature-name {
    font-weight: 600;
    color: var(--on-text);
    font-size: 13px;
    margin-bottom: 4px;
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .feature-level {
    font-size: 10px;
    font-weight: 600;
    color: var(--on-gold);
    background: var(--on-gold-bg);
    padding: 2px 6px;
    border-radius: 4px;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  .feature-desc {
    font-size: 11px;
    color: var(--on-text-dim);
    line-height: 1.5;
    margin: 4px 0 0 0;
  }

  .spells-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
    gap: 8px;
  }

  .spell-card {
    background: var(--on-bg-root);
    border: 1px solid var(--on-border);
    border-radius: 8px;
    padding: 8px 10px;
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .spell-name {
    font-weight: 500;
    color: var(--on-text);
    font-size: 12px;
  }

  .spell-level {
    font-size: 9px;
    color: var(--on-gold);
    font-weight: 600;
    text-transform: uppercase;
  }

  .spell-school {
    font-size: 9px;
    color: var(--on-text-dim);
    text-transform: capitalize;
  }

  .variants-list {
    list-style: none;
    padding: 0;
    margin: 0;
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  .variant-item {
    background: var(--on-bg-root);
    border: 1px solid var(--on-border);
    border-radius: 8px;
    padding: 10px;
  }

  .variant-name {
    font-weight: 500;
    color: var(--on-text);
    font-size: 13px;
    margin-bottom: 2px;
  }

  .variant-desc {
    font-size: 11px;
    color: var(--on-text-dim);
    margin: 0;
  }

  @media (max-width: 520px) {
    .info-popup-overlay {
      padding: 0;
      align-items: flex-end;
    }
    
    .info-popup {
      max-height: 90vh;
      border-radius: 16px 16px 0 0;
      max-width: 100%;
    }
    
    .popup-content {
      padding: 1rem;
    }
  }
</style>