<script lang="ts">
  import { draft, classes, selectClass, setSubclass, getClassDef, setName, setLevel, addSpell, removeSpell } from '$lib/builder.svelte';
  import Card from '$lib/ui/Card.svelte';
  import Tag from '$lib/ui/Tag.svelte';
  import InfoPopup from '$lib/components/InfoPopup.svelte';
  import Button from '$lib/ui/Button.svelte';
  import type { ClassEntry, SubClassEntry, SpellEntry } from '$lib/types';

  const selectedClassId = $derived(draft.value?.classes?.[0]?.id);
  const selectedSubclassId = $derived(draft.value?.classes?.[0]?.subclassId);
  const classDef = $derived(getClassDef());

  // Character name
  const characterName = $derived(draft.value?.name || '');

  // Level (1-20)
  const characterLevel = $derived(draft.value?.classes?.[0]?.level || 1);

  // Popup state
  let popupOpen = $state(false);
  let popupType = $state<'class' | 'subclass'>('class');
  let popupData = $state<any>(null);

  // Avatar selection (16 placeholders)
  const avatarOptions = Array.from({ length: 16 }, (_, i) => ({
    id: i + 1,
    url: `https://ui-avatars.com/api/?name=Hero+${i + 1}&background=8b0000&color=fff&size=128`
  }));

  // Available spells for the selected class (if spellcaster)
  const availableSpells = $derived(() => {
    if (!classDef?.spellcasting) return [];
    // Return all spells that match the class
    // This would come from the content store in a real implementation
    return [] as SpellEntry[];
  });

  // Selected spells for the character
  const selectedSpells = $derived(draft.value?.spells || []);

  function openClassPopup(cls: any) {
    popupType = 'class';
    popupData = cls;
    popupOpen = true;
  }

  function openSubclassPopup(sc: any) {
    popupType = 'subclass';
    popupData = sc;
    popupOpen = true;
  }

  function closePopup() {
    popupOpen = false;
    popupData = null;
  }

  // Avatar selection
  function selectAvatar(url: string) {
    // Store avatar URL in draft or character data
    draft.value = { ...draft.value, avatarUrl: url };
  }

  // Name change
  function handleNameChange(e: Event) {
    const value = (e.target as HTMLInputElement).value;
    setName(value);
  }

  // Level change
  function handleLevelChange(e: Event) {
    const value = parseInt((e.target as HTMLInputElement).value, 10);
    if (!isNaN(value) && value >= 1 && value <= 20) {
      setLevel(value);
    }
  }

  // Spell selection
  function toggleSpell(spellId: string) {
    const hasSpell = selectedSpells.includes(spellId);
    if (hasSpell) {
      removeSpell(spellId);
    } else {
      addSpell(spellId);
    }
  }
</script>

<div class="step-class" data-testid="step-class">
  <h2 class="step-title" data-testid="step-title">Escolha Sua Classe</h2>
  <p class="step-desc">Sua classe define seu papel no combate, suas habilidades e seu estilo de jogo.</p>

  <!-- Character Name & Avatar -->
  <div class="character-header">
    <div class="avatar-section">
      <label class="avatar-label">Avatar</label>
      <div class="avatar-grid">
        {#each avatarOptions as avatar}
          <button
            type="button"
            class="avatar-option {draft.value?.avatarUrl === avatar.url ? 'selected' : ''}"
            onclick={() => selectAvatar(avatar.url)}
            aria-label={`Avatar ${avatar.id}`}
          >
            <img src={avatar.url} alt={`Avatar ${avatar.id}`} />
          </button>
        {/each}
      </div>
    </div>

    <div class="name-level-section">
      <div class="name-field">
        <label for="char-name">Nome do Personagem</label>
        <input
          id="char-name"
          type="text"
          value={characterName}
          oninput={handleNameChange}
          placeholder="Ex: Thorin Escudo-de-Carvalho"
          class="name-input"
          data-testid="char-name-input"
        />
      </div>

      <div class="level-field">
        <label for="char-level">Nível (1-20)</label>
        <div class="level-control">
          <input
            id="char-level"
            type="number"
            min="1"
            max="20"
            value={characterLevel}
            oninput={handleLevelChange}
            class="level-input"
            data-testid="char-level-input"
          />
        </div>
      </div>
    </div>
  </div>

  <!-- Class Selection -->
  <div class="class-grid">
    {#each classes.value as cls}
      <Card
        variant={selectedClassId === cls.id ? 'selected' : 'interactive'}
        onclick={() => selectClass(cls.id)}
        class="class-card"
        on:click={() => openClassPopup(cls)}
        id="class-card"
        data-class-id={cls.id}
      >
        <div class="cc-header">
          <span class="cc-name">{cls.name}</span>
          <span class="cc-die">d{cls.hitDie}</span>
        </div>
        <p class="cc-desc">{cls.spellcaster ? 'Conjurador' : 'Marcial'} · Salvaguardas: {cls.savingThrows.join(', ')}</p>
        <div class="cc-tags">
          {#each cls.primaryAbility as ab}
            <Tag variant="primary">{ab}</Tag>
          {/each}
          {#each cls.savingThrows as st}
            <Tag>{st}</Tag>
          {/each}
        </div>
        {#if cls.features}
          <div class="cc-features">
            {#each cls.features.filter(f => f.level === 1) as feat}
              <span class="cc-feat">{feat.name}</span>
            {/each}
          </div>
        {/if}
      </Card>
    {/each}
  </div>

  <!-- Subclass Selection -->
  {#if selectedClassId && classDef?.subClasses && classDef.subClasses.length > 0}
    <div class="subclass-section">
      <h3 class="subclass-title">Subclasse <span class="subclass-optional">(nível {classDef.subclassLevel})</span></h3>
      <div class="subclass-grid">
        {#each classDef.subClasses as sc}
          <Card
            variant={selectedSubclassId === sc.id ? 'selected' : 'interactive'}
            onclick={() => setSubclass(selectedClassId, sc.id)}
            class="subclass-card"
            on:click={() => openSubclassPopup(sc)}
          >
            <div class="sc-name">{sc.name}</div>
            <p class="sc-desc">{sc.description}</p>
          </Card>
        {/each}
      </div>
    </div>
  {/if}

  <!-- Spells Selection (for spellcasters) -->
  {#if selectedClassId && classDef?.spellcasting}
    <div class="spells-section">
      <h3 class="spells-title">Feitiços Iniciais</h3>
      <p class="spells-desc">Escolha seus feitiços iniciais (pode mudar depois).</p>
      <div class="spells-grid">
        {#each availableSpells as spell}
          <Card
            variant={selectedSpells.includes(spell.id) ? 'selected' : 'interactive'}
            onclick={() => toggleSpell(spell.id)}
            class="spell-card"
          >
            <div class="spell-header">
              <span class="spell-name">{spell.name}</span>
              <span class="spell-level">Nível {spell.level}</span>
            </div>
            <p class="spell-school">{spell.school}</p>
          </Card>
        {/each}
      </div>
    </div>
  {/if}

  <!-- Popup -->
  {#if popupOpen}
    <InfoPopup
      type={popupType}
      data={popupData}
      isOpen={popupOpen}
      onClose={() => { popupOpen = false; popupData = null; }}
    />
  {/if}
</div>

<style>
  .step-class { animation: fadeIn 0.3s ease; }
  @keyframes fadeIn { from { opacity: 0; transform: translateY(8px); } to { opacity: 1; transform: translateY(0); } }

  .step-title { font-size: 28px; font-weight: 500; margin: 0 0 8px 0; letter-spacing: -0.3px; }
  .step-desc { font-size: 14px; color: var(--on-text-dim); margin: 0 0 24px 0; }

  /* Character Header */
  .character-header {
    display: grid;
    grid-template-columns: 120px 1fr;
    gap: 24px;
    margin-bottom: 32px;
    padding: 20px;
    background: var(--on-bg-surface);
    border: 1px solid var(--on-border);
    border-radius: 12px;
  }

  .avatar-label { font-size: 12px; color: var(--on-text-dim); margin-bottom: 8px; display: block; }

  .avatar-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 8px;
  }

  .avatar-option {
    aspect-ratio: 1;
    border: 2px solid var(--on-border);
    border-radius: 50%;
    overflow: hidden;
    background: var(--on-bg-root);
    cursor: pointer;
    padding: 0;
    transition: all 0.15s;
  }

  .avatar-option:hover {
    border-color: var(--on-red);
    transform: scale(1.05);
  }

  .avatar-option.selected {
    border-color: var(--on-gold);
    box-shadow: 0 0 0 2px var(--on-gold);
  }

  .avatar-option img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .name-level-section {
    display: flex;
    flex-direction: column;
    gap: 16px;
    justify-content: center;
  }

  .name-field, .level-field { display: flex; flex-direction: column; gap: 6px; }

  .name-field label, .level-field label {
    font-size: 12px;
    color: var(--on-text-dim);
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  .name-input, .level-input {
    padding: 12px 16px;
    font-size: 16px;
    background: var(--on-bg-root);
    border: 1px solid var(--on-border);
    border-radius: 8px;
    color: var(--on-text);
    outline: none;
    transition: border-color 0.15s;
  }

  .name-input:focus, .level-input:focus { border-color: var(--on-red); }

  .level-control { display: flex; align-items: center; gap: 12px; }

  .level-input {
    width: 80px;
    text-align: center;
    font-size: 20px;
    font-weight: 600;
  }

  /* Class Grid */
  .class-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
    gap: 16px;
    margin-bottom: 32px;
  }
  :global(.class-card) { padding: 20px; }
  .cc-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 8px;
  }
  .cc-name { font-size: 18px; font-weight: 500; color: var(--on-text); }
  .cc-die { font-size: 12px; color: var(--on-text-dim); font-family: var(--font-mono); }
  .cc-desc { font-size: 12px; color: var(--on-text-dim); margin: 0 0 12px 0; line-height: 1.5; }
  .cc-tags { display: flex; gap: 6px; flex-wrap: wrap; margin-bottom: 12px; }
  .cc-features { display: flex; flex-wrap: wrap; gap: 6px; }
  .cc-feat {
    font-size: 11px;
    padding: 4px 10px;
    background: var(--on-bg-root);
    border: 1px solid var(--on-border);
    border-radius: 4px;
    color: var(--on-text-muted);
  }

  /* Subclass Section */
  .subclass-section { margin-top: 8px; }
  .subclass-title { font-size: 18px; font-weight: 500; margin: 0 0 16px 0; }
  .subclass-optional { font-size: 13px; color: var(--on-text-dim); font-weight: 400; }
  .subclass-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
    gap: 12px;
  }
  :global(.subclass-card) { padding: 16px; }
  .sc-name { font-size: 15px; font-weight: 500; color: var(--on-text); margin-bottom: 4px; }
  .sc-desc { font-size: 12px; color: var(--on-text-dim); margin: 0; line-height: 1.5; }

  /* Spells Section */
  .spells-section { margin-top: 8px; }
  .spells-title { font-size: 18px; font-weight: 500; margin: 0 0 4px 0; }
  .spells-desc { font-size: 13px; color: var(--on-text-dim); margin: 0 0 16px 0; }
  .spells-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 12px;
  }
  :global(.spell-card) { padding: 16px; }
  .spell-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 4px;
  }
  .spell-name { font-size: 14px; font-weight: 500; color: var(--on-text); }
  .spell-level { font-size: 11px; color: var(--on-gold); background: var(--on-gold-bg); padding: 2px 8px; border-radius: 4px; }
  .spell-school { font-size: 11px; color: var(--on-text-dim); text-transform: capitalize; }

  /* Animations */
  .step-class { animation: fadeIn 0.3s ease; }
  @keyframes fadeIn { from { opacity: 0; transform: translateY(8px); } to { opacity: 1; transform: translateY(0); } }

  .step-title { font-size: 28px; font-weight: 500; margin: 0 0 8px 0; letter-spacing: -0.3px; }
  .step-desc { font-size: 14px; color: var(--on-text-dim); margin: 0 0 24px 0; }
</style>