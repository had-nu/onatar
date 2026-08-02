// Types matching the backend API contract (PRD §3.5 GET /content, POST /build).

export type Ability = 'STR' | 'DEX' | 'CON' | 'INT' | 'WIS' | 'CHA'

export interface User {
  id: string
  github_id: number
  login: string
  name: string | null
  avatar_url: string | null
  email: string | null
}

export const ABILITIES: Ability[] = ['STR', 'DEX', 'CON', 'INT', 'WIS', 'CHA']

export interface Content {
  classes: Class[]
  subclasses: Subclass[]
  species: Species[]
  backgrounds: Background[]
  spells: Spell[]
  feats: Feat[]
  features: Feature[]
}

export interface Class {
  id: string
  name: string
  hitDie: string
  spellcaster: boolean
  subclassLevel: number
  suggestedSpecies: string[]
  suggestedBackgrounds: string[]
  data: Record<string, unknown>
}

export interface Subclass {
  id: string
  classId: string
  name: string
  levelRequired: number
  data: Record<string, unknown>
}

export interface Species {
  id: string
  name: string
  data: Record<string, unknown>
}

export interface Background {
  id: string
  name: string
  data: Record<string, unknown>
}

export interface Spell {
  id: string
  name: string
  level: number
  school: string
  data: Record<string, unknown>
}

export interface Feat {
  id: string
  name: string
  prerequisites: Record<string, unknown>
  data: Record<string, unknown>
}

export interface Feature {
  id: string
  classId: string
  subclassId: string
  name: string
  level: number
  data: Record<string, unknown>
}

export interface ClassInput {
  id: string
  level: number
  subclassId?: string
}

export interface BuildRequest {
  name: string
  classes: ClassInput[]
  speciesId?: string
  backgroundId?: string
  abilityScores: Record<Ability, number>
  abilityMethod?: string
  skills?: string[]
  spells?: string[]
  feats?: string[]
  equipment?: string[]
  isNpc?: boolean
}

export interface AbilityScore {
  score: number
  mod: number
}

export interface SheetFeature {
  name: string
  level: number
  description: string
}

export interface PendingChoice {
  type: string
  description: string
}

export interface Sheet {
  level: number
  hp: { max: number; current: number }
  ac: number
  proficiencyBonus: number
  spellSlots: number[]
  abilities: Record<Ability, AbilityScore>
  features: SheetFeature[]
  pendingChoices: PendingChoice[]
}

export interface BuildResponse {
  sheet: Sheet
}

/** Live/editable sheet state kept alongside the computed sheet (Sprint 4). */
export interface SheetLive {
  hpCurrent: number
  slotsUsed: number[]
  conditions: string[]
  resources: Record<string, number>
}

export function emptyLive(sheet: Sheet | undefined): SheetLive {
  return {
    hpCurrent: sheet?.hp.current ?? sheet?.hp.max ?? 0,
    slotsUsed: Array.from({ length: 9 }, () => 0),
    conditions: [],
    resources: {},
  }
}

export const CONDITIONS = [
  'blinded',
  'charmed',
  'deafened',
  'exhaustion',
  'frightened',
  'grappled',
  'incapacitated',
  'poisoned',
  'prone',
  'restrained',
  'stunned',
  'unconscious',
] as const

export interface Character {
  id: string
  name: string
  isNpc: boolean
  campaignId?: string
  draft: BuildRequest
  sheet?: Sheet
  live?: SheetLive
  createdAt: number
  updatedAt: number
}

// data.description / data.primaryAbility live inside the flexible `data` JSON.
export function dataString(data: Record<string, unknown> | undefined, key: string): string {
  const v = data?.[key]
  return typeof v === 'string' ? v : ''
}

/* ──────────────────────────────────────────────────────────────
 * New Arcanum Builder Types (from zip) — for offline mock data
 * and richer CharacterSheet with savingThrows, spellSlots, features[]
 * ────────────────────────────────────────────────────────────── */

export interface AbilityScores {
  STR: number
  DEX: number
  CON: number
  INT: number
  WIS: number
  CHA: number
}

export interface ClassEntry {
  id: string
  name: string
  hitDie: number
  spellcaster: boolean
  primaryAbility: string[]
  savingThrows: string[]
  subclassLevel?: number
  subClasses?: SubClassEntry[]
  features?: ClassFeature[]
  spellcasting?: SpellcastingEntry | null
  skillChoices?: { count: number; from: string[] }
}

export interface SubClassEntry {
  id: string
  name: string
  description: string
}

export interface ClassFeature {
  name: string
  level: number
  description: string
}

export interface SpellcastingEntry {
  ability: string
  preparedSpells: number[]
  knownSpells: boolean
  slots: Record<string, number[]>
}

export interface SpeciesEntry {
  id: string
  name: string
  description: string
  traits: TraitEntry[]
  abilityBonuses: Partial<AbilityScores>
  size: string
  speed: number
  languages: string[]
  variants: SpeciesVariant[]
}

export interface SpeciesVariant {
  id: string
  name: string
  description: string
}

export interface TraitEntry {
  name: string
  description: string
}

export interface BackgroundEntry {
  id: string
  name: string
  description: string
  skillProficiencies: string[]
  toolProficiencies: string[]
  languages: number[]
  equipment: string[]
  feature: { name: string; description: string }
}

export interface SpellEntry {
  id: string
  name: string
  level: number
  school: string
  description: string
}

export interface FeatEntry {
  id: string
  name: string
  prerequisites: Record<string, unknown>
  description: string
}

export interface FeatureEntry {
  id: string
  classId?: string
  name: string
  level: number
  description: string
}

export interface BuildRequestNew {
  name: string
  classes: ClassReq[]
  backgroundId: string
  speciesId: string
  speciesVariant?: string
  level: number
  abilityScores: AbilityScores
  abilityMethod: 'standard' | 'point-buy' | 'roll'
  skills: string[]
  spells: string[]
  feats: string[]
}

export interface ClassReq {
  id: string
  level: number
  subclassId?: string
}

export interface BuildResponseNew {
  valid: boolean
  errors?: string[]
  sheet?: CharacterSheet
}

export interface CharacterSheet {
  name: string
  level: number
  hp: { max: number; current?: number }
  ac: number
  proficiencyBonus: number
  abilities: AbilityScores
  savingThrows: Record<string, { proficient: boolean; modifier: number }>
  skills: Record<string, { proficient: boolean; modifier: number }>
  features: string[]
  spells?: string[]
  spellSlots?: Record<string, number>
}

export interface ChoicePoint {
  type: 'subclass' | 'spell' | 'ability-improvement' | 'skill'
  classId?: string
  level: number
  name: string
  description: string
  options: Array<{ id: string; name: string; description: string }>
}

export interface BuilderSnapshot {
  draft: BuildRequestNew
  step: number
  timestamp: number
}