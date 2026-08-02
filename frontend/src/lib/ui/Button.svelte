<script lang="ts">
  interface Props {
    variant?: 'primary' | 'secondary' | 'outline' | 'ghost' | 'danger';
    size?: 'sm' | 'md' | 'lg';
    disabled?: boolean;
    type?: 'button' | 'submit';
    children?: import('svelte').Snippet;
    onclick?: (e: MouseEvent) => void;
    class?: string;
  }
  let {
    variant = 'primary',
    size = 'md',
    disabled = false,
    type = 'button',
    children,
    onclick,
    class: className = ''
  }: Props = $props();
</script>

<button
  {type}
  class="on-btn on-btn--{variant} on-btn--{size} {className}"
  class:on-btn--disabled={disabled}
  {disabled}
  onclick={onclick}
>
  {@render children?.()}
</button>

<style>
  .on-btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    border-radius: 8px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.15s ease;
    border: 1px solid transparent;
    background: transparent;
    color: var(--on-text-muted);
    font-family: inherit;
  }
  .on-btn--sm { padding: 6px 14px; font-size: 12px; }
  .on-btn--md { padding: 10px 20px; font-size: 14px; }
  .on-btn--lg { padding: 14px 28px; font-size: 16px; }

  .on-btn--primary {
    background: var(--on-red);
    color: #fff;
    border-color: var(--on-red);
    box-shadow: 0 2px 8px rgba(197, 0, 9, 0.3);
  }
  .on-btn--primary:hover:not(.on-btn--disabled) {
    background: var(--on-red-hover);
    transform: translateY(-1px);
    box-shadow: 0 4px 16px rgba(197, 0, 9, 0.4);
  }

  .on-btn--secondary {
    background: var(--on-bg-elevated);
    color: var(--on-text);
    border-color: var(--on-border);
  }
  .on-btn--secondary:hover:not(.on-btn--disabled) {
    border-color: var(--on-border-light);
    background: var(--on-bg-hover);
  }

  .on-btn--outline {
    border-color: var(--on-border);
    color: var(--on-text-muted);
  }
  .on-btn--outline:hover:not(.on-btn--disabled) {
    border-color: var(--on-red);
    color: var(--on-text);
  }

  .on-btn--ghost:hover:not(.on-btn--disabled) {
    background: var(--on-bg-hover);
    color: var(--on-text);
  }

  .on-btn--danger {
    background: rgba(197, 0, 9, 0.12);
    color: #ff6b6b;
    border-color: rgba(197, 0, 9, 0.3);
  }
  .on-btn--danger:hover:not(.on-btn--disabled) {
    background: rgba(197, 0, 9, 0.2);
  }

  .on-btn--disabled {
    opacity: 0.4;
    cursor: not-allowed;
    transform: none !important;
  }
</style>
