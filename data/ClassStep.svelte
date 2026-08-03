<script lang="ts">
  import {
    classes, spells, draft,
    selectClass, setSubclass, requestPreview, getClassDef
  } from '$lib/builder.svelte';
  import { loadContent } from '$lib/content.svelte';
  import InfoPopup from '$lib/components/InfoPopup.svelte';
  import Button from '$lib/ui/Button.svelte';

  // ── Local state ───────────────────────────────────────
  let selectedAvatar = $state<string | null>(null);
  let characterName = $state(draft.value.name || '');
  let selectedClassId = $state(draft.value.classes[0]?.id ?? '');
  let selectedSubclassId = $state(draft.value.classes[0]?.subclassId ?? '');
  let selectedLevel = $state(draft.value.level || 1);
  let selectedSpells = $state<Set<string>>(new Set(draft.value.spells));
  let popupOpen = $state(false);
  let popupData = $state<{
    title: string; subtitle?: string; description: string;
    meta: { label: string; value: string }[];
    tags: string[]; features: string[]; color: string;
  } | null>(null);

  // ── Derived ───────────────────────────────────────────
  const classDef = $derived(getClassDef());
  const showSubclasses = $derived(
    classDef?.subclassLevel ? selectedLevel >= classDef.subclassLevel : false
  );
  const showSpells = $derived(classDef?.spellcaster ?? false);
  const spellList = $derived(() => {
    if (!classDef?.spellcasting?.preparedSpells) return [];
    const maxPrepared = classDef.spellcasting.preparedSpells[Math.min(selectedLevel, 20) - 1] || 0;
    return spells.value
      .filter(s => s.level <= selectedLevel)
      .slice(0, Math.max(maxPrepared, 8));
  });

  // ── Ensure content loaded ─────────────────────────────
  $effect(() => {
    if (classes.value.length === 0) {
      loadContent();
    }
  });

  // ── Actions ───────────────────────────────────────────
  function handleSelectClass(classId: string) {
    selectedClassId = classId;
    selectedSubclassId = '';
    selectedSpells = new Set();
    selectClass(classId);
  }

  function handleSelectSubclass(subclassId: string) {
    selectedSubclassId = subclassId;
    setSubclass(selectedClassId, subclassId);
  }

  function handleLevelChange(level: number) {
    selectedLevel = level;
    draft.value.level = level;
    draft.value.classes = draft.value.classes.map(c =>
      c.id === selectedClassId ? { ...c, level } : c
    );
    requestPreview();
  }

  function toggleSpell(spellId: string) {
    const next = new Set(selectedSpells);
    if (next.has(spellId)) next.delete(spellId);
    else next.add(spellId);
    selectedSpells = next;
    draft.value.spells = Array.from(next);
    requestPreview();
  }

  function setName(name: string) {
    characterName = name;
    draft.value.name = name;
    requestPreview();
  }

  // ── Popup helpers ─────────────────────────────────────
  function openClassPopup(cls: typeof classes.value[number]) {
    popupData = {
      title: cls.name,
      subtitle: `d${cls.hitDie} · ${cls.spellcaster ? 'Conjurador' : 'Marcial'}`,
      description: cls.features?.map(f => f.description).join(' ') || '',
      meta: [
        { label: 'Dado de Vida', value: `d${cls.hitDie}` },
        { label: 'Salvaguardas', value: cls.savingThrows?.join(', ') || '—' },
        { label: 'Habilidade Principal', value: cls.primaryAbility?.join(' / ') || '—' },
        { label: 'Subclasse', value: `Nível ${cls.subclassLevel || 3}` },
      ],
      tags: cls.spellcaster ? ['Conjurador', 'Arcano'] : ['Marcial', 'Combate'],
      features: cls.features?.map(f => `${f.name} (nv.${f.level})`) || [],
      color: 'var(--on-red)',
    };
    popupOpen = true;
  }

  function openSubclassPopup(sub: typeof classDef.subClasses[number]) {
    popupData = {
      title: sub.name,
      subtitle: 'Subclasse',
      description: sub.description || '',
      meta: [],
      tags: ['Subclasse'],
      features: [],
      color: 'var(--on-gold)',
    };
    popupOpen = true;
  }

  function openSpellPopup(spell: typeof spells.value[number]) {
    popupData = {
      title: spell.name,
      subtitle: `${spell.school} · Nível ${spell.level}`,
      description: spell.description || '',
      meta: [
        { label: 'Nível', value: String(spell.level) },
        { label: 'Escola', value: spell.school },
      ],
      tags: [spell.school, spell.level === 0 ? 'Truque' : `Nível ${spell.level}`],
      features: [],
      color: 'var(--on-gold)',
    };
    popupOpen = true;
  }

  // ── Avatar list ───────────────────────────────────────
  const AVATARS = [
    { id: 'warrior-m', label: 'Guerreiro' }, { id: 'warrior-f', label: 'Guerreira' },
    { id: 'mage-m', label: 'Mago' }, { id: 'mage-f', label: 'Maga' },
    { id: 'rogue-m', label: 'Ladino' }, { id: 'rogue-f', label: 'Ladina' },
    { id: 'cleric-m', label: 'Clérigo' }, { id: 'cleric-f', label: 'Clériga' },
    { id: 'ranger-m', label: 'Patrulheiro' }, { id: 'ranger-f', label: 'Patrulheira' },
    { id: 'bard-m', label: 'Bardo' }, { id: 'bard-f', label: 'Barda' },
    { id: 'paladin-m', label: 'Paladino' }, { id: 'paladin-f', label: 'Paladina' },
    { id: 'druid-m', label: 'Druida' }, { id: 'druid-f', label: 'Druida' },
  ];
</script>

{#if popupData}
  <InfoPopup
    open={popupOpen}
    title={popupData.title}
    subtitle={popupData.subtitle}
    description={popupData.description}
    meta={popupData.meta}
    tags={popupData.tags}
    features={popupData.features}
    color={popupData.color}
    onClose={() => { popupOpen = false; popupData = null; }}
  />
{/if}

<div class="class-step">
  <h2 class="step-heading">Escolhe uma Classe</h2>

  <!-- Avatar -->
  <section class="section">
    <h3 class="section-title">Avatar</h3>
    <div class="avatar-grid">
      {#each AVATARS as avatar}
        <button
          class="avatar-btn"
          class:selected={selectedAvatar === avatar.id}
          onclick={() => selectedAvatar = avatar.id}
          aria-label={avatar.label}
        >
          <div class="avatar-circle">{avatar.label.charAt(0)}</div>
          <span class="avatar-label">{avatar.label}</span>
        </button>
      {/each}
    </div>
  </section>

  <!-- Nome -->
  <section class="section">
    <h3 class="section-title">Nome do Personagem</h3>
    <input
      type="text"
      class="name-input"
      placeholder="Ex: Thrain Coração-de-Ferro"
      bind:value={characterName}
      oninput={(e) => setName(e.currentTarget.value)}
    />
  </section>

  <!-- Classe -->
  <section class="section">
    <h3 class="section-title">Classe</h3>
    {#if classes.value.length === 0}
      <div class="empty-state">
        <div class="spinner"></div>
        <p>A carregar classes...</p>
      </div>
    {:else}
      <div class="class-grid">
        {#each classes.value as cls}
          <button
            class="class-card"
            class:selected={selectedClassId === cls.id}
            onclick={() => handleSelectClass(cls.id)}
            ondblclick={() => openClassPopup(cls)}
          >
            <div class="class-header">
              <span class="class-icon">{cls.spellcaster ? '🔮' : '⚔️'}</span>
              <button
                class="info-btn"
                onclick={(e) => { e.stopPropagation(); openClassPopup(cls); }}
                aria-label="Informações"
              >
                ℹ
              </button>
            </div>
            <div class="class-name">{cls.name}</div>
            <div class="class-meta">d{cls.hitDie} · {cls.primaryAbility?.join(' / ') || ''}</div>
            <div class="class-tags">
              {#each cls.savingThrows.slice(0, 2) as st}
                <span class="tag">{st}</span>
              {/each}
              {#if cls.spellcaster}
                <span class="tag tag--magic">Conjurador</span>
              {/if}
            </div>
          </button>
        {/each}
      </div>
    {/if}
  </section>

  <!-- Nível -->
  <section class="section">
    <h3 class="section-title">
      Nível
      <span class="level-badge">{selectedLevel}</span>
    </h3>
    <div class="level-wrap">
      <input
        type="range"
        class="level-slider"
        min="1"
        max="20"
        bind:value={selectedLevel}
        onchange={() => handleLevelChange(selectedLevel)}
      />
      <div class="level-ticks">
        <span>1</span><span>5</span><span>10</span><span>15</span><span>20</span>
      </div>
    </div>
  </section>

  <!-- Subclasse -->
  {#if showSubclasses && classDef?.subClasses}
    <section class="section">
      <h3 class="section-title">Subclasse <span class="hint">(nível {classDef.subclassLevel}+)</span></h3>
      <div class="subclass-grid">
        {#each classDef.subClasses as sub}
          <button
            class="subclass-card"
            class:selected={selectedSubclassId === sub.id}
            onclick={() => handleSelectSubclass(sub.id)}
            ondblclick={() => openSubclassPopup(sub)}
          >
            <div class="subclass-name">{sub.name}</div>
            {#if sub.description}
              <div class="subclass-desc">{sub.description}</div>
            {/if}
            <button
              class="info-btn info-btn--small"
              onclick={(e) => { e.stopPropagation(); openSubclassPopup(sub); }}
              aria-label="Informações"
            >
              ℹ
            </button>
          </button>
        {/each}
      </div>
    </section>
  {/if}

  <!-- Spells -->
  {#if showSpells}
    <section class="section">
      <h3 class="section-title">
        Spells
        <span class="hint">({selectedSpells.size} selecionados)</span>
      </h3>
      {#if spellList().length === 0}
        <div class="empty-state">
          <p>Sem spells disponíveis para este nível.</p>
        </div>
      {:else}
        <div class="spell-grid">
          {#each spellList() as spell}
            <button
              class="spell-card"
              class:selected={selectedSpells.has(spell.id)}
              onclick={() => toggleSpell(spell.id)}
              ondblclick={() => openSpellPopup(spell)}
            >
              <div class="spell-header">
                <span class="spell-level">{spell.level === 0 ? 'Truque' : `Nível ${spell.level}`}</span>
                {#if selectedSpells.has(spell.id)}
                  <span class="spell-check">✓</span>
                {/if}
              </div>
              <div class="spell-name">{spell.name}</div>
              <div class="spell-school">{spell.school}</div>
              <button
                class="info-btn info-btn--small"
                onclick={(e) => { e.stopPropagation(); openSpellPopup(spell); }}
                aria-label="Informações"
              >
                ℹ
              </button>
            </button>
          {/each}
        </div>
      {/if}
    </section>
  {/if}
</div>

<style>
  .class-step {
    padding: 8px 0 32px;
  }

  .step-heading {
    font-family: var(--font-display);
    font-size: 1.6rem;
    font-weight: 700;
    margin: 0 0 24px;
    color: var(--on-text);
    letter-spacing: -0.02em;
  }

  .section {
    margin-bottom: 28px;
  }

  .section-title {
    font-size: 0.75rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.12em;
    color: var(--on-text-muted);
    margin: 0 0 14px;
    display: flex;
    align-items: center;
    gap: 10px;
  }

  .hint {
    text-transform: none;
    letter-spacing: 0;
    font-weight: 500;
    color: var(--on-text-dim);
    font-size: 0.7rem;
  }

  .level-badge {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    min-width: 28px;
    height: 24px;
    background: var(--on-red);
    color: #fff;
    border-radius: 6px;
    font-size: 0.8rem;
    font-weight: 700;
    padding: 0 8px;
  }

  /* ── Avatar ─────────────────────────────────────────── */
  .avatar-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(72px, 1fr));
    gap: 10px;
  }

  .avatar-btn {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 6px;
    padding: 10px 6px;
    background: var(--on-bg-surface);
    border: 1.5px solid var(--on-border);
    border-radius: 10px;
    cursor: pointer;
    transition: all 0.18s ease;
    color: var(--on-text-muted);
  }

  .avatar-btn:hover {
    border-color: var(--on-red);
    transform: translateY(-2px);
  }

  .avatar-btn.selected {
    border-color: var(--on-red);
    background: var(--accent-bg);
    box-shadow: 0 0 0 3px var(--on-red-glow);
  }

  .avatar-circle {
    width: 44px;
    height: 44px;
    border-radius: 50%;
    background: linear-gradient(135deg, #3a3028, #25201c);
    display: flex;
    align-items: center;
    justify-content: center;
    border: 1.5px solid var(--on-border);
    font-size: 1.1rem;
    font-weight: 700;
    color: var(--on-text);
  }

  .avatar-label {
    font-size: 0.65rem;
    font-weight: 600;
    text-align: center;
    line-height: 1.2;
  }

  /* ── Nome ───────────────────────────────────────────── */
  .name-input {
    width: 100%;
    max-width: 420px;
    padding: 12px 16px;
    font-size: 1rem;
    font-weight: 500;
    background: var(--on-bg-surface);
    border: 1.5px solid var(--on-border);
    border-radius: 8px;
    color: var(--on-text);
    outline: none;
    transition: border-color 0.2s, box-shadow 0.2s;
    font-family: inherit;
  }

  .name-input:focus {
    border-color: var(--on-red);
    box-shadow: 0 0 0 3px var(--on-red-glow);
  }

  .name-input::placeholder {
    color: var(--on-text-dim);
  }

  /* ── Classe ─────────────────────────────────────────── */
  .class-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(170px, 1fr));
    gap: 12px;
  }

  .class-card {
    position: relative;
    display: flex;
    flex-direction: column;
    gap: 6px;
    padding: 18px 14px;
    background: var(--on-bg-surface);
    border: 1.5px solid var(--on-border);
    border-radius: 12px;
    cursor: pointer;
    transition: all 0.2s ease;
    text-align: left;
  }

  .class-card:hover {
    border-color: var(--on-red);
    transform: translateY(-3px);
    box-shadow: var(--shadow);
  }

  .class-card.selected {
    border-color: var(--on-red);
    background: linear-gradient(180deg, rgba(255,255,255,0.02), transparent);
    box-shadow: 0 0 0 3px var(--on-red-glow), var(--shadow);
  }

  .class-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .class-icon {
    font-size: 1.6rem;
    line-height: 1;
  }

  .class-name {
    font-size: 1rem;
    font-weight: 700;
    color: var(--on-text);
  }

  .class-meta {
    font-size: 0.75rem;
    color: var(--on-text-dim);
  }

  .class-tags {
    display: flex;
    gap: 6px;
    flex-wrap: wrap;
    margin-top: 4px;
  }

  .tag {
    font-size: 0.65rem;
    font-weight: 600;
    padding: 3px 8px;
    border-radius: 4px;
    background: var(--on-bg-hover);
    color: var(--on-text-muted);
    text-transform: uppercase;
    letter-spacing: 0.04em;
  }

  .tag--magic {
    background: var(--gold-bg);
    color: var(--on-gold);
    border: 1px solid var(--gold-border);
  }

  /* ── Info button ────────────────────────────────────── */
  .info-btn {
    position: absolute;
    top: 10px;
    right: 10px;
    width: 22px;
    height: 22px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--on-bg-hover);
    border: 1px solid var(--on-border);
    border-radius: 5px;
    color: var(--on-text-dim);
    font-size: 0.75rem;
    font-weight: 700;
    cursor: pointer;
    opacity: 0;
    transition: opacity 0.15s, background 0.15s;
    padding: 0;
    line-height: 1;
  }

  .class-card:hover .info-btn,
  .subclass-card:hover .info-btn,
  .spell-card:hover .info-btn {
    opacity: 1;
  }

  .info-btn:hover {
    background: var(--on-bg-root);
    color: var(--on-text);
    border-color: var(--on-red);
  }

  .info-btn--small {
    width: 18px;
    height: 18px;
    font-size: 0.65rem;
    top: 8px;
    right: 8px;
  }

  /* ── Nível ──────────────────────────────────────────── */
  .level-wrap {
    max-width: 480px;
  }

  .level-slider {
    width: 100%;
    -webkit-appearance: none;
    appearance: none;
    height: 5px;
    background: var(--on-bg-hover);
    border-radius: 3px;
    outline: none;
  }

  .level-slider::-webkit-slider-thumb {
    -webkit-appearance: none;
    appearance: none;
    width: 20px;
    height: 20px;
    border-radius: 50%;
    background: var(--on-red);
    cursor: pointer;
    border: 3px solid var(--on-bg-root);
    box-shadow: 0 2px 8px rgba(197, 0, 9, 0.4);
    transition: transform 0.15s;
  }

  .level-slider::-webkit-slider-thumb:hover {
    transform: scale(1.15);
  }

  .level-slider::-moz-range-thumb {
    width: 20px;
    height: 20px;
    border-radius: 50%;
    background: var(--on-red);
    cursor: pointer;
    border: 3px solid var(--on-bg-root);
    box-shadow: 0 2px 8px rgba(197, 0, 9, 0.4);
  }

  .level-ticks {
    display: flex;
    justify-content: space-between;
    margin-top: 8px;
    font-size: 0.7rem;
    color: var(--on-text-dim);
    font-weight: 600;
    padding: 0 4px;
  }

  /* ── Subclasse ──────────────────────────────────────── */
  .subclass-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
    gap: 10px;
  }

  .subclass-card {
    position: relative;
    display: flex;
    flex-direction: column;
    gap: 4px;
    padding: 16px 12px;
    background: var(--on-bg-surface);
    border: 1.5px solid var(--on-border);
    border-radius: 10px;
    cursor: pointer;
    transition: all 0.2s ease;
    text-align: left;
  }

  .subclass-card:hover {
    border-color: var(--on-gold);
    transform: translateY(-2px);
  }

  .subclass-card.selected {
    border-color: var(--on-gold);
    background: var(--gold-bg);
    box-shadow: 0 0 0 3px rgba(212, 175, 55, 0.12);
  }

  .subclass-name {
    font-size: 0.9rem;
    font-weight: 700;
    color: var(--on-text);
    padding-right: 20px;
  }

  .subclass-desc {
    font-size: 0.75rem;
    color: var(--on-text-muted);
    line-height: 1.4;
  }

  /* ── Spells ─────────────────────────────────────────── */
  .spell-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
    gap: 10px;
  }

  .spell-card {
    position: relative;
    display: flex;
    flex-direction: column;
    gap: 4px;
    padding: 14px 12px;
    background: var(--on-bg-surface);
    border: 1.5px solid var(--on-border);
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.18s ease;
    text-align: left;
  }

  .spell-card:hover {
    border-color: var(--on-gold);
    transform: translateY(-2px);
  }

  .spell-card.selected {
    border-color: var(--on-gold);
    background: var(--gold-bg);
    box-shadow: 0 0 0 3px rgba(212, 175, 55, 0.1);
  }

  .spell-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .spell-level {
    font-size: 0.65rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.06em;
    color: var(--on-text-dim);
  }

  .spell-check {
    width: 18px;
    height: 18px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--on-gold);
    color: var(--on-bg-root);
    border-radius: 50%;
    font-size: 0.7rem;
    font-weight: 700;
  }

  .spell-name {
    font-size: 0.9rem;
    font-weight: 600;
    color: var(--on-text);
    line-height: 1.3;
  }

  .spell-school {
    font-size: 0.75rem;
    color: var(--on-text-muted);
  }

  /* ── Empty state ────────────────────────────────────── */
  .empty-state {
    text-align: center;
    padding: 32px 16px;
    color: var(--on-text-dim);
    font-size: 0.9rem;
  }

  .spinner {
    width: 24px;
    height: 24px;
    border: 2px solid var(--on-border);
    border-top-color: var(--on-red);
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
    margin: 0 auto 12px;
  }

  @keyframes spin {
    from { transform: rotate(0deg); }
    to { transform: rotate(360deg); }
  }

  @media (max-width: 640px) {
    .class-grid { grid-template-columns: repeat(2, 1fr); }
    .subclass-grid { grid-template-columns: repeat(2, 1fr); }
    .spell-grid { grid-template-columns: repeat(2, 1fr); }
    .avatar-grid { grid-template-columns: repeat(4, 1fr); }
  }
</style>
