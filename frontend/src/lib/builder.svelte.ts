// Builder wizard store (PRD §3.4 / RF-01). Holds the in-progress BuildRequest,
// the active step, per-step validation, and an undo/redo history.
import { box } from './box.svelte'
import { content } from './content.svelte'
import { createCharacter } from './characters.svelte'
import type { Ability, BuildRequest, Character, ClassInput } from './types'

export type StepId = 'class' | 'background' | 'species' | 'abilities' | 'equipment' | 'review'
export type AbilityMethod = 'standard-array' | 'point-buy' | 'rolled'

export interface Step {
  id: StepId
  label: string
}

export const STEPS: Step[] = [
  { id: 'class', label: 'Classe' },
  { id: 'background', label: 'Background' },
  { id: 'species', label: 'Espécie' },
  { id: 'abilities', label: 'Atributos' },
  { id: 'equipment', label: 'Equipamento' },
  { id: 'review', label: 'Revisão' },
]

export const STANDARD_ARRAY = [15, 14, 13, 12, 10, 8] as const
export const POINT_BUY_MIN = 8
export const POINT_BUY_MAX = 15
export const POINT_BUY_BUDGET = 27
export const POINT_BUY_COST: Record<number, number> = {
  8: 0,
  9: 1,
  10: 2,
  11: 3,
  12: 4,
  13: 5,
  14: 7,
  15: 9,
}

export interface AbilityState {
  score: number
  method: AbilityMethod
  assigned: Partial<Record<Ability, number>>
  rolled: number[]
  pointBuy: Record<Ability, number>
}

export interface BuilderState {
  stepIndex: number
  draft: BuildRequest
  abilities: AbilityState
  equipment: string[]
  name: string
  history: Snapshot[]
  future: Snapshot[]
}

interface Snapshot {
  stepIndex: number
  draft: BuildRequest
  abilities: AbilityState
  equipment: string[]
  name: string
}

const ABILITIES: Ability[] = ['STR', 'DEX', 'CON', 'INT', 'WIS', 'CHA']

export function blankDraft(): BuildRequest {
  return {
    name: '',
    classes: [],
    abilityScores: { STR: 8, DEX: 8, CON: 8, INT: 8, WIS: 8, CHA: 8 },
    abilityMethod: 'standard-array',
    skills: [],
    spells: [],
    feats: [],
    isNpc: false,
  }
}

export const builder = box<BuilderState>({
  stepIndex: 0,
  draft: blankDraft(),
  abilities: {
    score: 0,
    method: 'standard-array',
    assigned: {},
    rolled: [],
    pointBuy: { STR: 8, DEX: 8, CON: 8, INT: 8, WIS: 8, CHA: 8 },
  },
  equipment: [],
  name: '',
  history: [],
  future: [],
})

function snapshot(): Snapshot {
  const s = builder.value
  return {
    stepIndex: s.stepIndex,
    draft: s.draft,
    abilities: s.abilities,
    equipment: s.equipment,
    name: s.name,
  }
}

/** Record history for undo, then apply a mutation to the current state. */
function mutate(fn: (s: BuilderState) => BuilderState) {
  const s = builder.value
  const prev: Snapshot = snapshot()
  const next: BuilderState = fn(s)
  const history = [...s.history, prev].slice(-50)
  builder.value = { ...next, history, future: [] }
}

/** Test helper: reset the wizard between tests. */
export function _resetBuilder() {
  builder.value = {
    stepIndex: 0,
    draft: blankDraft(),
    abilities: {
      score: 0,
      method: 'standard-array',
      assigned: {},
      rolled: [],
      pointBuy: { STR: 8, DEX: 8, CON: 8, INT: 8, WIS: 8, CHA: 8 },
    },
    equipment: [],
    name: '',
    history: [],
    future: [],
  }
}

// --- steps ---

export function step(): Step {
  return STEPS[builder.value.stepIndex]
}

export function setStep(i: number) {
  const s = builder.value
  if (i < 0 || i >= STEPS.length) return
  builder.value = { ...s, stepIndex: i }
}

export function nextStep() {
  const s = builder.value
  if (!canGoNext()) return
  builder.value = { ...s, stepIndex: Math.min(s.stepIndex + 1, STEPS.length - 1) }
}

export function prevStep() {
  const s = builder.value
  builder.value = { ...s, stepIndex: Math.max(s.stepIndex - 1, 0) }
}

// --- undo / redo ---

export function undo() {
  const s = builder.value
  const prev = s.history.at(-1)
  if (!prev) return
  const history = s.history.slice(0, -1)
  builder.value = {
    ...s,
    stepIndex: prev.stepIndex,
    draft: prev.draft,
    abilities: prev.abilities,
    equipment: prev.equipment,
    name: prev.name,
    history,
    future: [snapshot(), ...s.future].slice(0, 50),
  }
}

export function redo() {
  const s = builder.value
  const next = s.future[0]
  if (!next) return
  const future = s.future.slice(1)
  builder.value = {
    ...s,
    stepIndex: next.stepIndex,
    draft: next.draft,
    abilities: next.abilities,
    equipment: next.equipment,
    name: next.name,
    history: [...s.history, snapshot()].slice(-50),
    future,
  }
}

// --- validation ---

export function validateStep(id: StepId): boolean {
  const s = builder.value
  switch (id) {
    case 'class':
      return s.draft.classes.length > 0 && s.draft.classes.every((c) => c.level >= 1)
    case 'background':
      return Boolean(s.draft.backgroundId)
    case 'species':
      return Boolean(s.draft.speciesId)
    case 'abilities':
      return abilitiesValid(s)
    case 'equipment':
      return true
    case 'review':
      return Boolean(s.name.trim())
  }
}

function abilitiesValid(s: BuilderState): boolean {
  const a = s.abilities
  if (a.method === 'point-buy') {
    const spent = Object.values(a.pointBuy).reduce((acc, v) => acc + POINT_BUY_COST[v], 0)
    return spent <= POINT_BUY_BUDGET && Object.values(a.pointBuy).every((v) => v >= POINT_BUY_MIN)
  }
  const assigned = Object.values(a.assigned)
  return assigned.length === ABILITIES.length && assigned.every((v) => typeof v === 'number')
}

export function canGoNext(): boolean {
  return validateStep(step().id)
}

// --- draft builders ---

function abilityScoresFrom(state: AbilityState): Record<Ability, number> {
  if (state.method === 'point-buy') {
    return { ...state.pointBuy }
  }
  const out = { STR: 8, DEX: 8, CON: 8, INT: 8, WIS: 8, CHA: 8 } as Record<Ability, number>
  for (const ab of ABILITIES) {
    if (typeof state.assigned[ab] === 'number') out[ab] = state.assigned[ab] as number
  }
  return out
}

function withAbilityScores(s: BuilderState): BuilderState {
  return {
    ...s,
    draft: {
      ...s.draft,
      abilityScores: abilityScoresFrom(s.abilities),
      abilityMethod: s.abilities.method,
    },
  }
}

// --- class step ---

export function selectClass(id: string) {
  mutate((s) => {
    const existing = s.draft.classes.some((c) => c.id === id)
    if (existing) return s
    const classes = [...s.draft.classes, { id, level: 1 } as ClassInput]
    return withAbilityScores({ ...s, draft: { ...s.draft, classes } })
  })
}

export function deselectClass(id: string) {
  mutate((s) => {
    const classes = s.draft.classes.filter((c) => c.id !== id)
    return withAbilityScores({ ...s, draft: { ...s.draft, classes } })
  })
}

export function setClassLevel(id: string, level: number) {
  mutate((s) => {
    const classes = s.draft.classes.map((c) => (c.id === id ? { ...c, level } : c))
    return withAbilityScores({ ...s, draft: { ...s.draft, classes } })
  })
}

export function selectSubclass(classId: string, subclassId: string) {
  mutate((s) => {
    const classes = s.draft.classes.map((c) => (c.id === classId ? { ...c, subclassId } : c))
    return { ...s, draft: { ...s.draft, classes } }
  })
}

// --- background / species ---

export function selectBackground(id: string) {
  mutate((s) => ({ ...s, draft: { ...s.draft, backgroundId: id } }))
}

export function selectSpecies(id: string) {
  mutate((s) => ({ ...s, draft: { ...s.draft, speciesId: id } }))
}

// --- abilities ---

export function setMethod(method: AbilityMethod) {
  mutate((s) => withAbilityScores({ ...s, abilities: { ...s.abilities, method } }))
}

export function assignAbility(ability: Ability, value: number | null) {
  mutate((s) => {
    const assigned = { ...s.abilities.assigned }
    if (value === null) delete assigned[ability]
    else assigned[ability] = value
    return withAbilityScores({ ...s, abilities: { ...s.abilities, assigned } })
  })
}

export function setPointBuy(ability: Ability, score: number) {
  const clamped = Math.min(POINT_BUY_MAX, Math.max(POINT_BUY_MIN, score))
  mutate((s) =>
    withAbilityScores({
      ...s,
      abilities: { ...s.abilities, pointBuy: { ...s.abilities.pointBuy, [ability]: clamped } },
    })
  )
}

/** Roll 4d6 drop lowest, once per ability. Returns the values. */
export function rollScores(): number[] {
  const scores = ABILITIES.map(() => {
    const dice = Array.from({ length: 4 }, () => 1 + Math.floor(Math.random() * 6))
    dice.sort((a, b) => a - b)
    return dice.slice(1).reduce((a, b) => a + b, 0)
  })
  mutate((s) => withAbilityScores({ ...s, abilities: { ...s.abilities, rolled: scores } }))
  return scores
}

// --- spells / skills / feats (content-driven) ---

export function toggleSpell(id: string) {
  mutate((s) => {
    const spells = s.draft.spells ?? []
    const next = spells.includes(id) ? spells.filter((x) => x !== id) : [...spells, id]
    return { ...s, draft: { ...s.draft, spells: next } }
  })
}

export function toggleSkill(id: string) {
  mutate((s) => {
    const skills = s.draft.skills ?? []
    const next = skills.includes(id) ? skills.filter((x) => x !== id) : [...skills, id]
    return { ...s, draft: { ...s.draft, skills: next } }
  })
}

export function toggleFeat(id: string) {
  mutate((s) => {
    const feats = s.draft.feats ?? []
    const next = feats.includes(id) ? feats.filter((x) => x !== id) : [...feats, id]
    return { ...s, draft: { ...s.draft, feats: next } }
  })
}

// --- equipment ---

export function toggleEquipment(item: string) {
  mutate((s) => {
    const next = s.equipment.includes(item)
      ? s.equipment.filter((x) => x !== item)
      : [...s.equipment, item]
    return { ...s, equipment: next }
  })
}

// --- suggestions (RF-07) ---

export function suggestedSpeciesForClass(): string[] {
  const c = firstClass()
  return c?.suggestedSpecies ?? []
}

export function suggestedBackgroundsForClass(): string[] {
  const c = firstClass()
  return c?.suggestedBackgrounds ?? []
}

function firstClass() {
  const classes = content.value?.classes ?? []
  return classes.find((c) => builder.value.draft.classes.some((x) => x.id === c.id))
}

export function recommendedSpells(): string[] {
  const data = firstClass()?.data?.recommendedSpells
  return Array.isArray(data) ? data.map(String) : []
}

// --- save (RF-03) ---

export function saveCharacterFromWizard(): Character | null {
  const s = builder.value
  if (!validateStep('review')) return null
  const draft: BuildRequest = { ...s.draft, name: s.name, equipment: s.equipment }
  return createCharacter(draft)
}
