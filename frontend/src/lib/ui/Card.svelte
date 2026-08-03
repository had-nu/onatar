<script lang="ts">
  interface Props {
    variant?: 'default' | 'elevated' | 'selected' | 'interactive';
    children?: import('svelte').Snippet;
    onclick?: () => void;
    class?: string;
    id?: string;
    'data-class-id'?: string;
  }
  let { variant = 'default', children, onclick, class: className = '', id, 'data-class-id': classId }: Props = $props();

  const isClickable = !!onclick;
</script>

<div
  class="on-card on-card--{variant} {className}"
  class:on-card--clickable={isClickable}
  role={isClickable ? 'button' : undefined}
  tabindex={isClickable ? 0 : undefined}
  onclick={isClickable ? onclick : undefined}
  id={id}
  data-class-id={classId}
>
  {@render children?.()}
</div>

<style>
  .on-card {
    background: linear-gradient(135deg, var(--on-bg-surface) 0%, var(--on-bg-elevated) 100%);
    border: 1px solid var(--on-border);
    border-radius: 12px;
    overflow: hidden;
    position: relative;
    transition: all 0.2s ease;
  }
  .on-card--clickable { cursor: pointer; }
  .on-card--clickable:hover {
    border-color: var(--on-red);
    transform: translateY(-2px);
    box-shadow: 0 4px 16px var(--on-red-glow);
  }
  .on-card--elevated:hover {
    border-color: var(--on-red);
    transform: translateY(-2px);
    box-shadow: 0 4px 16px var(--on-red-glow);
  }
  .on-card--selected {
    border-color: var(--on-gold);
    box-shadow: 0 0 0 1px var(--on-gold), 0 4px 20px rgba(201, 169, 78, 0.15);
  }
  .on-card--selected::before {
    content: '✓';
    position: absolute;
    top: 10px;
    right: 10px;
    width: 22px;
    height: 22px;
    border-radius: 50%;
    background: var(--on-gold);
    color: var(--on-bg-root);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 12px;
    font-weight: 700;
    z-index: 2;
  }
  .on-card--interactive:active {
    transform: translateY(0) scale(0.98);
  }
</style>