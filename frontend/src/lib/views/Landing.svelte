<script lang="ts">
  import { navigate } from '$lib/router.svelte';
  import Button from '$lib/ui/Button.svelte';
  import Card from '$lib/ui/Card.svelte';

  const stats = [
    { value: 2, label: 'Classes' },
    { value: 2, label: 'Species' },
    { value: 2, label: 'Backgrounds' },
    { value: 0, label: 'Characters Created' },
  ];

  const features = [
    {
      icon: `<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>`,
      title: 'Campaign Management',
      desc: 'Organize your players, track dice rolls, and share content with your group.'
    },
    {
      icon: `<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><circle cx="12" cy="12" r="10"/><path d="M12 6v6l4 2"/></svg>`,
      title: 'Digital Combat',
      desc: 'Initiative tracker, conditions, and HP in real time. Less paper, more play.'
    },
    {
      icon: `<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/><polyline points="3.27 6.96 12 12.01 20.73 6.96"/><line x1="12" y1="22.08" x2="12" y2="12"/></svg>`,
      title: 'Content Library',
      desc: 'Explore classes, species, spells, and items. All indexed and searchable.'
    }
  ];
</script>

<div class="landing">
  <!-- Hero -->
  <section class="hero">
    <div class="hero-badge">
      <span class="badge-dot"></span>
      D&D 2024 · 5.5e · v0.1.0
    </div>
    <h1>Forge Your Hero</h1>
    <p class="hero-sub">D&D 2024 (5.5e) Character Creation</p>
    <p class="hero-desc">Inspired by D&D Beyond, focused on accessibility for new players and simplicity for DMs. No subscriptions, no friction.</p>

    <div class="hero-actions">
      <Button variant="primary" size="lg" onclick={() => navigate('/builder')}>Create Character</Button>
      <Button variant="outline" size="lg" onclick={() => navigate('/characters')}>My Characters</Button>
      <Button variant="ghost" size="lg" onclick={() => navigate('/content')}>Explore Content</Button>
    </div>

    <div class="card-stack">
      <Card variant="default" class="hero-card">
        <div class="hc-icon">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M14.5 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7.5L14.5 2z"/><polyline points="14 2 14 8 20 8"/></svg>
        </div>
        <div class="hc-name">Elara Moonwhisper</div>
        <div class="hc-meta">Elf · Wizard · Level 5</div>
      </Card>
      <Card variant="default" class="hero-card">
        <div class="hc-icon">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M14.5 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7.5L14.5 2z"/><polyline points="14 2 14 8 20 8"/></svg>
        </div>
        <div class="hc-name">Thorin Oakenshield</div>
        <div class="hc-meta">Dwarf · Fighter · Level 3</div>
      </Card>
      <Card variant="default" class="hero-card">
        <div class="hc-icon">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M14.5 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7.5L14.5 2z"/><polyline points="14 2 14 8 20 8"/></svg>
        </div>
        <div class="hc-name">New Character</div>
        <div class="hc-meta">Click to begin</div>
      </Card>
    </div>
  </section>

  <!-- Stats -->
  <div class="stats-bar">
    {#each stats as s}
      <div class="stat">
        <div class="stat-value">{s.value}</div>
        <div class="stat-label">{s.label}</div>
      </div>
    {/each}
  </div>

  <!-- Features -->
  <div class="features">
    {#each features as f}
      <Card variant="elevated" class="feature-card">
        <div class="feat-icon">{@html f.icon}</div>
        <h3>{f.title}</h3>
        <p>{f.desc}</p>
      </Card>
    {/each}
  </div>

  <!-- Teaser -->
  <div class="teaser">
    <Card variant="default" class="teaser-card">
      <h4>Start Playing</h4>
      <p>New to D&D? Onatar guides you step by step through creating your first character without overwhelming you with rules.</p>
      <div class="teaser-tag"><span class="tt-dot"></span> Guide for New Players</div>
    </Card>
    <Card variant="default" class="teaser-card">
      <h4>For Dungeon Masters</h4>
      <p>Import characters, manage encounters, and share content with your group. Everything centralized in one place.</p>
      <div class="teaser-tag"><span class="tt-dot"></span> DM Tools</div>
    </Card>
  </div>
</div>

<style>
  .landing {
    max-width: 1100px;
    margin: 0 auto;
    padding: 0 16px 48px;
  }

  .hero {
    padding: 56px 0 40px;
    text-align: center;
  }
  .hero-badge {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    font-size: 11px;
    font-weight: 500;
    text-transform: uppercase;
    letter-spacing: 2px;
    color: var(--on-gold);
    background: rgba(201, 169, 78, 0.08);
    border: 1px solid rgba(201, 169, 78, 0.2);
    border-radius: 999px;
    padding: 6px 16px;
    margin-bottom: 24px;
  }
  .badge-dot {
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background: var(--on-gold);
    animation: pulse 2s infinite;
  }
  @keyframes pulse { 0%, 100% { opacity: 1; } 50% { opacity: 0.4; } }

  .hero h1 {
    font-size: 48px;
    font-weight: 500;
    line-height: 1.1;
    margin: 0 0 12px 0;
    letter-spacing: -0.5px;
  }
  .hero-sub { font-size: 18px; color: var(--on-text-muted); margin: 0 0 8px 0; }
  .hero-desc { font-size: 14px; color: var(--on-text-dim); line-height: 1.6; margin: 0 auto 32px; max-width: 480px; }

  .hero-actions {
    display: flex;
    gap: 12px;
    justify-content: center;
    flex-wrap: wrap;
    margin-bottom: 40px;
  }

  .card-stack {
    display: flex;
    justify-content: center;
    gap: 16px;
    perspective: 1000px;
    flex-wrap: wrap;
  }
  :global(.hero-card) {
    width: 180px;
    height: 220px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 12px;
    transition: all 0.3s ease;
  }
  :global(.hero-card):nth-child(1) { transform: rotate(-4deg) translateY(-8px); }
  :global(.hero-card):nth-child(2) { transform: rotate(2deg); opacity: 0.7; }
  :global(.hero-card):nth-child(3) { transform: rotate(6deg) translateY(4px); opacity: 0.5; }
  :global(.hero-card):hover {
    transform: rotate(0deg) translateY(-12px) scale(1.03);
    opacity: 1;
    z-index: 10;
  }
  .hc-icon {
    width: 48px;
    height: 48px;
    border-radius: 50%;
    border: 1.5px solid var(--on-border);
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--on-text-muted);
  }
  .hc-name { font-size: 14px; font-weight: 500; color: var(--on-text); }
  .hc-meta { font-size: 11px; color: var(--on-text-dim); }

  .stats-bar {
    display: flex;
    justify-content: center;
    gap: 48px;
    padding: 24px;
    margin-bottom: 40px;
    background: linear-gradient(135deg, var(--on-bg-surface) 0%, var(--on-bg-elevated) 100%);
    border: 1px solid var(--on-border);
    border-radius: 12px;
    flex-wrap: wrap;
  }
  .stat { text-align: center; }
  .stat-value { font-size: 32px; font-weight: 500; color: var(--on-gold); line-height: 1; margin-bottom: 6px; font-variant-numeric: tabular-nums; }
  .stat-label { font-size: 11px; color: var(--on-text-dim); text-transform: uppercase; letter-spacing: 1.5px; }

  .features {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 16px;
    margin-bottom: 40px;
  }
  :global(.feature-card) { padding: 24px; }
  :global(.feature-card):hover {
    border-color: var(--on-red);
    transform: translateY(-2px);
    box-shadow: 0 4px 16px var(--on-red-glow);
  }
  .feat-icon {
    width: 40px;
    height: 40px;
    border-radius: 10px;
    border: 1px solid var(--on-border);
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: 16px;
    color: var(--on-text-muted);
  }
  :global(.feature-card) h3 { font-size: 16px; font-weight: 500; margin: 0 0 6px 0; }
  :global(.feature-card) p { font-size: 13px; color: var(--on-text-dim); line-height: 1.5; margin: 0; }

  .teaser {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 16px;
  }
  :global(.teaser-card) { padding: 28px; }
  :global(.teaser-card) h4 {
    font-size: 15px;
    font-weight: 500;
    margin: 0 0 8px 0;
    display: flex;
    align-items: center;
    gap: 8px;
  }
  :global(.teaser-card) h4::before {
    content: '';
    display: inline-block;
    width: 4px;
    height: 16px;
    background: var(--on-red);
    border-radius: 2px;
  }
  :global(.teaser-card) p { font-size: 13px; color: var(--on-text-dim); margin: 0; line-height: 1.5; }
  .teaser-tag {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    font-size: 11px;
    color: var(--on-text-dim);
    margin-top: 12px;
  }
  .tt-dot { width: 6px; height: 6px; border-radius: 50%; background: var(--on-green); }

  @media (max-width: 800px) {
    .hero h1 { font-size: 32px; }
    .features { grid-template-columns: 1fr; }
    .teaser { grid-template-columns: 1fr; }
    .stats-bar { gap: 24px; }
  }
</style>