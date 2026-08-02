import type {
  AbilityScores, BuildRequest, BuildResponse, CharacterSheet,
  ChoicePoint, BuilderSnapshot, ClassEntry, SpeciesEntry, BackgroundEntry, SpellEntry
} from './types';
import { box } from './box.svelte';
import { loadContent } from './content.svelte';
import { content } from './content.svelte';

const defaultAbilities: AbilityScores = { STR: 10, DEX: 10, CON: 10, INT: 10, WIS: 10, CHA: 10 };

const initialDraft = (): BuildRequest => ({
  name: '',
  classes: [],
  backgroundId: '',
  speciesId: '',
  speciesVariant: undefined,
  level: 1,
  abilityScores: { ...defaultAbilities },
  abilityMethod: 'standard',
  skills: [],
  spells: [],
  feats: [],
});

// ─── State ─────────────────────────────────────────────────
export const draft = box<BuildRequest>(initialDraft());
export const step = box(0);
export const preview = box<CharacterSheet | null>(null);
export const pendingChoices = box<ChoicePoint[]>([]);
export const isLoading = box(false);
export const error = box<string | null>(null);

// Undo / redo
const MAX_HISTORY = 50;
let history = $state<BuilderSnapshot[]>([{ draft: initialDraft(), step: 0, timestamp: Date.now() }]);
let historyIndex = $state(0);

// Content cache (from content store, with mock fallback)
export const classes = box<ClassEntry[]>([]);
export const species = box<SpeciesEntry[]>([]);
export const backgrounds = box<BackgroundEntry[]>([]);
export const spells = box<SpellEntry[]>([]);

// Load content on initialization
async function initializeContent() {
  try {
    const data = await loadContent();
    classes.value = data.classes.map(c => ({
      id: c.id,
      name: c.name,
      hitDie: c.hitDie ? parseInt(c.hitDie.replace('d', '')) : 8,
      spellcaster: c.spellcaster,
      primaryAbility: c.data?.primaryAbility ? [c.data.primaryAbility] : [],
      savingThrows: c.data?.savingThrows || [],
      subclassLevel: c.subclassLevel,
      subClasses: c.subclasses?.filter(sc => sc.classId === c.id).map(sc => ({
        id: sc.id,
        name: sc.name,
        description: sc.data?.description || ''
      })) || [],
      features: c.data?.features?.map(f => ({
        name: f.name,
        level: f.level,
        description: f.description || ''
      })) || [],
      spellcasting: c.spellcaster && c.data?.spellcasting ? {
        ability: c.data.spellcasting.ability,
        preparedSpells: c.data.spellcasting.preparedSpells || [],
        knownSpells: c.data.spellcasting.knownSpells || false,
        slots: c.data.spellcasting.slots || {}
      } : null,
      skillChoices: c.data?.skillChoices || { count: 0, from: [] }
    }));

    species.value = data.species.map(s => ({
      id: s.id,
      name: s.name,
      description: s.data?.description || '',
      traits: s.data?.traits || [],
      abilityBonuses: s.data?.abilityScores || {},
      size: s.data?.size || 'Medium',
      speed: s.data?.speed || 30,
      languages: s.data?.languages || ['Common'],
      variants: s.data?.variants || []
    }));

    backgrounds.value = data.backgrounds.map(b => ({
      id: b.id,
      name: b.name,
      description: b.data?.description || '',
      skillProficiencies: b.data?.skillProficiencies || [],
      toolProficiencies: b.data?.toolProficiencies || [],
      languages: b.data?.languages || 0,
      equipment: b.data?.equipment || [],
      feature: b.data?.feature || { name: '', description: '' }
    }));

    spells.value = data.spells.map(s => ({
      id: s.id,
      name: s.name,
      level: s.level,
      school: s.school,
      description: s.data?.description || ''
    }));
  } catch (err) {
    console.warn('Failed to load content, using mock data:', err);
    // Fallback to mock data is already set as initial values
  }
}

// Initialize content on module load
initializeContent();

// ─── Derived ───────────────────────────────────────────────
export function getTotalLevel() { return draft.value.classes.reduce((s, c) => s + c.level, 0); }
export function getCanUndo() { return historyIndex > 0; }
export function getCanRedo() { return historyIndex < history.length - 1; }
export function getCurrentStepValid() { return validateStep(step.value); }
export function getClassDef() { return classes.value.find(c => c.id === draft.value.classes[0]?.id); }
export function getSpeciesDef() { return species.value.find(s => s.id === draft.value.speciesId); }
export function getBackgroundDef() { return backgrounds.value.find(b => b.id === draft.value.backgroundId); }

// ─── Step validation ───────────────────────────────────────
function validateStep(s: number): boolean {
  switch (s) {
    case 0: return draft.value.classes.length > 0 && draft.value.classes[0].level >= 1;
    case 1: return !!draft.value.backgroundId;
    case 2: return !!draft.value.speciesId;
    case 3: return Object.values(draft.value.abilityScores).every(v => v >= 3 && v <= 20);
    case 4: return true; // equipment optional for now
    case 5: return !!draft.value.name && pendingChoices.value.length === 0;
    default: return false;
  }
}

// ─── History (undo/redo) ───────────────────────────────────
function pushHistory() {
  // Remove future history if we're not at the end
  if (historyIndex < history.length - 1) {
    history = history.slice(0, historyIndex + 1);
  }
  history.push({ draft: structuredClone(draft.value), step: step.value, timestamp: Date.now() });
  if (history.length > MAX_HISTORY) history.shift();
  else historyIndex++;
}

export function undo() {
  if (!getCanUndo()) return;
  historyIndex--;
  const snap = history[historyIndex];
  draft.value = structuredClone(snap.draft);
  step.value = snap.step;
  requestPreview();
}

export function redo() {
  if (!getCanRedo()) return;
  historyIndex++;
  const snap = history[historyIndex];
  draft.value = structuredClone(snap.draft);
  step.value = snap.step;
  requestPreview();
}

// ─── Draft mutations ───────────────────────────────────────
export function setStep(s: number) {
  if (s < 0 || s > 5) return;
  if (s > step.value && !getCurrentStepValid()) return; // block forward if invalid
  pushHistory();
  step.value = s;
}

export function nextStep() { setStep(step.value + 1); }
export function prevStep() { setStep(step.value - 1); }

export function setName(name: string) {
  draft.value.name = name;
  requestPreview();
}

export function selectClass(classId: string) {
  pushHistory();
  draft.value.classes = [{ id: classId, level: 1 }];
  draft.value.level = 1;
  requestPreview();
}

export function setSubclass(classId: string, subclassId: string) {
  pushHistory();
  draft.value.classes = draft.value.classes.map(c =>
    c.id === classId ? { ...c, subclassId } : c
  );
  requestPreview();
}

export function setBackground(id: string) {
  pushHistory();
  draft.value.backgroundId = id;
  // Auto-add background skills
  const bg = backgrounds.value.find(b => b.id === id);
  if (bg) {
    draft.value.skills = [...new Set([...draft.value.skills, ...bg.skillProficiencies])];
  }
  requestPreview();
}

export function setSpecies(id: string, variant?: string) {
  pushHistory();
  draft.value.speciesId = id;
  draft.value.speciesVariant = variant;
  // Apply ability bonuses
  const sp = species.value.find(s => s.id === id);
  if (sp) {
    const bonuses = sp.abilityBonuses;
    draft.value.abilityScores = {
      STR: defaultAbilities.STR + (bonuses.STR || 0),
      DEX: defaultAbilities.DEX + (bonuses.DEX || 0),
      CON: defaultAbilities.CON + (bonuses.CON || 0),
      INT: defaultAbilities.INT + (bonuses.INT || 0),
      WIS: defaultAbilities.WIS + (bonuses.WIS || 0),
      CHA: defaultAbilities.CHA + (bonuses.CHA || 0),
    };
  }
  requestPreview();
}

export function setAbilityScore(ability: keyof AbilityScores, value: number) {
  draft.value.abilityScores[ability] = Math.max(3, Math.min(20, value));
  requestPreview();
}

export function setAbilityMethod(method: BuildRequest['abilityMethod']) {
  pushHistory();
  draft.value.abilityMethod = method;
  // Reset scores based on method
  if (method === 'standard') {
    draft.value.abilityScores = { ...defaultAbilities };
  }
  requestPreview();
}

export function toggleSkill(skillId: string) {
  const set = new Set(draft.value.skills);
  if (set.has(skillId)) set.delete(skillId);
  else set.add(skillId);
  draft.value.skills = Array.from(set);
  requestPreview();
}

export function addSpell(spellId: string) {
  if (!draft.value.spells.includes(spellId)) {
    draft.value.spells = [...draft.value.spells, spellId];
    requestPreview();
  }
}

export function removeSpell(spellId: string) {
  draft.value.spells = draft.value.spells.filter(s => s !== spellId);
  requestPreview();
}

// ─── Debounced Preview ─────────────────────────────────────
let debounceTimer: ReturnType<typeof setTimeout> | null = null;

export function requestPreview() {
  if (debounceTimer) clearTimeout(debounceTimer);
  debounceTimer = setTimeout(async () => {
    await computePreview();
  }, 250);
}

async function computePreview() {
  // Minimal validation
  if (!draft.value.classes.length || !draft.value.backgroundId || !draft.value.speciesId) {
    preview.value = null;
    pendingChoices.value = [];
    return;
  }

  isLoading.value = true;
  error.value = null;

  try {
    // Try backend first
    const res = await fetch('/api/v1/build', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(draft.value),
    });

    if (res.ok) {
      const data: BuildResponse = await res.json();
      preview.value = data.sheet ?? null;
    } else {
      // Fallback to local computation
      preview.value = computeLocalPreview();
    }
  } catch {
    // Offline fallback
    preview.value = computeLocalPreview();
  }

  deriveChoices();
  isLoading.value = false;
}

// ─── Local Preview Computation (offline mode) ──────────────
function computeLocalPreview(): CharacterSheet {
  const cls = getClassDef();
  const lvl = getTotalLevel();
  const prof = Math.floor((lvl + 7) / 4);

  // HP: hit die + CON mod at level 1, then average + CON mod per level
  const conMod = Math.floor((draft.value.abilityScores.CON - 10) / 2);
  const hd = cls?.hitDie || 8;
  const hpMax = hd + conMod + (lvl - 1) * (Math.floor(hd / 2) + 1 + conMod);

  // AC: base 10 + DEX mod (simplified — no armor in MVP)
  const dexMod = Math.floor((draft.value.abilityScores.DEX - 10) / 2);
  const ac = 10 + dexMod;

  // Saving throws
  const saves: CharacterSheet['savingThrows'] = {};
  const abilities = ['STR', 'DEX', 'CON', 'INT', 'WIS', 'CHA'] as const;
  for (const ab of abilities) {
    const mod = Math.floor((draft.value.abilityScores[ab] - 10) / 2);
    const proficient = cls?.savingThrows.includes(ab) ?? false;
    saves[ab] = { proficient, modifier: mod + (proficient ? prof : 0) };
  }

  // Features
  const features: string[] = [];
  if (cls?.features) {
    for (const f of cls.features) {
      if (f.level <= lvl) features.push(f.name);
    }
  }
  if (getSpeciesDef()?.traits) {
    for (const t of getSpeciesDef()!.traits) features.push(t.name);
  }
  if (getBackgroundDef()?.feature) features.push(getBackgroundDef()!.feature.name);

  // Spell slots (simplified full caster progression)
  const spellSlots: Record<string, number> = {};
  if (cls?.spellcaster && cls.spellcasting) {
    const slots = cls.spellcasting.slots;
    for (const [level, arr] of Object.entries(slots)) {
      spellSlots[level] = arr[Math.min(lvl, 20) - 1] || 0;
    }
  }

  return {
    name: draft.value.name || 'Unnamed',
    level: lvl,
    hp: { max: Math.max(1, hpMax) },
    ac: Math.max(10, ac),
    proficiencyBonus: prof,
    abilities: { ...draft.value.abilityScores },
    savingThrows: saves,
    skills: {}, // simplified
    features,
    spells: draft.value.spells,
    spellSlots,
  };
}

// ─── Derive Pending Choices ────────────────────────────────
function deriveChoices() {
  const choices: ChoicePoint[] = [];
  const lvl = getTotalLevel();

  for (const cls of draft.value.classes) {
    const def = classes.value.find(c => c.id === cls.id);
    if (!def) continue;

    // Subclass choice
    if (def.subclassLevel && lvl >= def.subclassLevel && !cls.subclassId) {
      choices.push({
        type: 'subclass',
        classId: cls.id,
        level: def.subclassLevel,
        name: 'Choose Subclass',
        description: `Select a subclass for ${def.name}`,
        options: (def.subClasses ?? []).map(sc => ({
          id: sc.id, name: sc.name, description: sc.description
        })),
      });
    }

    // Skill proficiencies (if not fully chosen)
    if (def.skillChoices) {
      const chosen = draft.value.skills.filter(s => def.skillChoices!.from.includes(s)).length;
      const remaining = def.skillChoices.count - chosen;
      if (remaining > 0) {
        choices.push({
          type: 'skill',
          classId: cls.id,
          level: 1,
          name: 'Skill Proficiencies',
          description: `Choose ${remaining} skill(s) from ${def.name}`,
          options: def.skillChoices.from
            .filter(s => !draft.value.skills.includes(s))
            .map(s => ({ id: s, name: s, description: '' })),
        });
      }
    }

    // Spell preparation (simplified)
    if (def.spellcaster && def.spellcasting?.preparedSpells) {
      const maxPrepared = def.spellcasting.preparedSpells[lvl - 1] || 0;
      const current = draft.value.spells.length;
      if (current < maxPrepared) {
        choices.push({
          type: 'spell',
          classId: cls.id,
          level: lvl,
          name: 'Prepare Spells',
          description: `Choose ${maxPrepared - current} prepared spell(s)`,
          options: spells.value
            .filter(s => s.level <= 1 && !draft.value.spells.includes(s.id))
            .map(s => ({ id: s.id, name: s.name, description: s.description })),
        });
      }
    }
  }

  // ASI at levels 4, 8, 12, 16, 19
  if (lvl >= 4 && lvl % 4 === 0) {
    choices.push({
      type: 'ability-improvement',
      level: lvl,
      name: 'Ability Score Improvement',
      description: 'Increase one ability by +2 or two by +1, or choose a feat.',
      options: [],
    });
  }

  pendingChoices.value = choices;
}

// ─── Save / Reset ──────────────────────────────────────────
export async function saveCharacter(): Promise<boolean> {
  if (!draft.value.name) {
    error.value = 'Character needs a name.';
    return false;
  }
  if (pendingChoices.value.length > 0) {
    error.value = 'There are pending choices.';
    return false;
  }

  const payload = {
    request: { ...draft.value },
    sheet: preview.value,
  };

  try {
    const res = await fetch('/api/v1/characters', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    });
    if (!res.ok) throw new Error('Save failed');

    // Also save to localStorage as fallback
    const existing = JSON.parse(localStorage.getItem('onatar-characters') || '[]');
    existing.push(payload);
    localStorage.setItem('onatar-characters', JSON.stringify(existing));

    return true;
  } catch {
    // Fallback: save locally
    const existing = JSON.parse(localStorage.getItem('onatar-characters') || '[]');
    existing.push(payload);
    localStorage.setItem('onatar-characters', JSON.stringify(existing));
    return true;
  }
}

export function resetBuilder() {
  pushHistory();
  draft.value = initialDraft();
  step.value = 0;
  preview.value = null;
  pendingChoices.value = [];
}

export function loadFromCharacter(character: any) {
  draft.value = {
    name: character.name || '',
    classes: (character.classes || []).map((c: any) => ({
      id: c.id, level: c.level, subclassId: c.subclassId
    })),
    backgroundId: character.backgroundId || '',
    speciesId: character.speciesId || '',
    speciesVariant: character.speciesVariant,
    level: character.level || 1,
    abilityScores: character.abilities || { ...defaultAbilities },
    abilityMethod: character.abilityMethod || 'standard',
    skills: character.skills || [],
    spells: character.spells || [],
    feats: character.feats || [],
  };
  step.value = 0;
  requestPreview();
}

/** Test helper: reset the module singleton between tests. */
export function _resetBuilder() {
  draft.value = initialDraft();
  step.value = 0;
  preview.value = null;
  pendingChoices.value = [];
  isLoading.value = false;
  error.value = null;
  history = [{ draft: initialDraft(), step: 0, timestamp: Date.now() }];
  historyIndex = 0;
}