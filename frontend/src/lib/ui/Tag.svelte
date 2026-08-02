<script lang="ts">
  interface Props {
    variant?: 'default' | 'primary' | 'gold' | 'green' | 'blue';
    children?: import('svelte').Snippet;
    class?: string;
  }
  let { variant = 'default', children, class: className = '' }: Props = $props();
</script>

<span class="on-tag on-tag--{variant} {className}">
  {@render children?.()}
</span>

<style>
  .on-tag {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    font-size: 11px;
    padding: 3px 10px;
    border-radius: 4px;
    font-weight: 500;
    border: 1px solid var(--on-border);
    background: var(--on-bg-hover);
    color: var(--on-text-muted);
  }
  .on-tag--primary {
    background: rgba(197, 0, 9, 0.12);
    color: #ff6b6b;
    border-color: rgba(197, 0, 9, 0.3);
  }
  .on-tag--gold {
    background: rgba(201, 169, 78, 0.12);
    color: var(--on-gold);
    border-color: rgba(201, 169, 78, 0.3);
  }
  .on-tag--green {
    background: rgba(0, 184, 122, 0.12);
    color: var(--on-green);
    border-color: rgba(0, 184, 122, 0.3);
  }
  .on-tag--blue {
    background: rgba(91, 141, 239, 0.12);
    color: var(--on-blue);
    border-color: rgba(91, 141, 239, 0.3);
  }
</style>
