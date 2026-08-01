// Characters store (RF-03): CRUD local in localStorage. A character holds a
// BuildRequest draft plus the last computed Sheet (cached for offline view).
import { box } from './box.svelte'
import type { BuildRequest, Character, Sheet, SheetLive } from './types'

export type { BuildRequest, Character, Sheet, SheetLive } from './types'

const KEY = 'onatar.characters'

function readAll(): Character[] {
  try {
    const raw = localStorage.getItem(KEY)
    if (!raw) return []
    const arr = JSON.parse(raw) as unknown
    return Array.isArray(arr) ? (arr as Character[]) : []
  } catch {
    return []
  }
}

export const characters = box<Character[]>(readAll())

function persist() {
  try {
    localStorage.setItem(KEY, JSON.stringify(characters.value))
  } catch {
    /* storage unavailable — ignore */
  }
}

export function newId(): string {
  if (typeof crypto !== 'undefined' && typeof crypto.randomUUID === 'function') {
    return crypto.randomUUID()
  }
  return `c_${Date.now()}_${Math.random().toString(36).slice(2, 8)}`
}

export function listCharacters(): Character[] {
  return characters.value
}

export function getCharacter(id: string): Character | undefined {
  return characters.value.find((c) => c.id === id)
}

/** A sensible starter build so a new character renders a real sheet. */
export function starterDraft(): BuildRequest {
  return {
    name: 'Novo personagem',
    classes: [{ id: 'fighter', level: 1 }],
    speciesId: 'human',
    backgroundId: 'sage',
    abilityScores: { STR: 15, DEX: 13, CON: 14, INT: 10, WIS: 12, CHA: 8 },
    abilityMethod: 'standard-array',
    skills: [],
    spells: [],
    feats: [],
    isNpc: false,
  }
}

export function createCharacter(draft: BuildRequest): Character {
  const now = Date.now()
  const c: Character = {
    id: newId(),
    name: draft.name || 'Novo personagem',
    isNpc: draft.isNpc ?? false,
    draft,
    createdAt: now,
    updatedAt: now,
  }
  characters.value = [...characters.value, c]
  persist()
  return c
}

export function saveCharacter(c: Character) {
  const updated: Character = { ...c, updatedAt: Date.now() }
  characters.value = [...characters.value.filter((x) => x.id !== c.id), updated]
  persist()
}

export function deleteCharacter(id: string) {
  characters.value = characters.value.filter((c) => c.id !== id)
  persist()
}

/** Update the live sheet state (HP, slots, conditions, resources). */
export function setLive(id: string, live: SheetLive) {
  const c = getCharacter(id)
  if (!c) return
  saveCharacter({ ...c, live })
}

/** Assign/unassign a campaign link. */
export function setCampaign(id: string, campaignId: string | undefined) {
  const c = getCharacter(id)
  if (!c) return
  saveCharacter({ ...c, campaignId })
}

/**
 * POST /build for an arbitrary draft (used by the builder live preview).
 * Throws with the API error message on failure.
 */
export async function buildDraft(draft: BuildRequest): Promise<Sheet> {
  const res = await fetch('/api/v1/build', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(draft),
  })
  if (!res.ok) {
    let msg = `HTTP ${res.status}`
    try {
      const body = (await res.json()) as { error?: { message?: string } }
      if (body?.error?.message) msg = body.error.message
    } catch {
      /* non-JSON error body */
    }
    throw new Error(msg)
  }
  const data = (await res.json()) as { sheet: Sheet }
  return data.sheet
}

/**
 * Compute the character sheet via POST /build. Falls back to the last cached
 * sheet when the network is unavailable (RNF-03); throws only if there is no
 * cache either.
 */
export async function buildCharacter(c: Character): Promise<Sheet> {
  try {
    const sheet = await buildDraft(c.draft)
    saveCharacter({ ...c, sheet })
    return sheet
  } catch (err) {
    if (c.sheet) return c.sheet
    throw err
  }
}

/** Test helper: reset the module singleton between tests. */
export function _resetCharacters() {
  characters.value = []
  try {
    localStorage.removeItem(KEY)
  } catch {
    /* ignore */
  }
}
