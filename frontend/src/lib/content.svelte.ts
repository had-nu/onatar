// Content store: loads GET /api/v1/content once, keeps an in-memory copy and a
// localStorage cache so the app degrades gracefully offline (PRD RNF-03).
// Falls back to mockData if API unavailable.
import { box } from './box.svelte'
import type { Class, Content, Species } from './types'
import { mockClasses, mockSpecies, mockBackgrounds, mockSpells, mockFeats, mockFeatures } from './mockData'

const KEY = 'onatar.content'

export const content = box<Content | null>(null)
export const contentError = box<string>('')
interface CacheEntry {
  at: number
  data: Content
}

function readCache(): Content | null {
  try {
    const raw = localStorage.getItem(KEY)
    if (!raw) return null
    const entry = JSON.parse(raw) as CacheEntry
    if (!entry?.data?.classes) return null
    return entry.data
  } catch {
    return null
  }
}

function writeCache(data: Content) {
  try {
    const entry: CacheEntry = { at: Date.now(), data }
    localStorage.setItem(KEY, JSON.stringify(entry))
  } catch {
    /* storage unavailable — ignore */
  }
}

function buildMockContent(): Content {
  return {
    classes: mockClasses.map(c => ({
      id: c.id,
      name: c.name,
      hitDie: `d${c.hitDie}`,
      spellcaster: c.spellcaster,
      subclassLevel: c.subclassLevel ?? 3,
      suggestedSpecies: [],
      suggestedBackgrounds: [],
      data: {
        ...c,
        primaryAbility: c.primaryAbility[0],
      },
    })),
    subclasses: mockClasses.flatMap(c =>
      (c.subClasses ?? []).map(sc => ({
        id: sc.id,
        classId: c.id,
        name: sc.name,
        levelRequired: c.subclassLevel ?? 3,
        data: { description: sc.description },
      }))
    ),
    species: mockSpecies.map(s => ({
      id: s.id,
      name: s.name,
      data: {
        traits: s.traits,
        abilityScores: s.abilityBonuses,
        description: s.description,
        size: s.size,
        speed: s.speed,
        languages: s.languages,
        variants: s.variants,
      },
    })),
    backgrounds: mockBackgrounds.map(b => ({
      id: b.id,
      name: b.name,
      data: {
        skillProficiencies: b.skillProficiencies,
        toolProficiencies: b.toolProficiencies,
        languages: b.languages,
        equipment: b.equipment,
        feature: b.feature,
        description: b.description,
      },
    })),
    spells: mockSpells.map(s => ({
      id: s.id,
      name: s.name,
      level: s.level,
      school: s.school,
      data: { description: s.description },
    })),
    feats: mockFeats,
    features: mockFeatures,
  }
}

/** Load content from the API, falling back to the localStorage cache, then mockData. */
export async function loadContent(force = false): Promise<Content> {
  if (content.value && !force) return content.value
  contentError.value = ''
  try {
    const res = await fetch('/api/v1/content')
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    const data = (await res.json()) as Content
    content.value = data
    writeCache(data)
    return data
  } catch (err) {
    const cached = readCache()
    if (cached) {
      content.value = cached
      return cached
    }
    // Fallback to mock data
    const mock = buildMockContent()
    content.value = mock
    writeCache(mock)
    return mock
  }
}

/** Synchronous access to whatever content is already loaded/cached. */
export function cachedContent(): Content | null {
  if (content.value) return content.value
  const cached = readCache()
  if (cached) content.value = cached
  return cached
}

export function classById(id: string): Class | undefined {
  return (content.value ?? cachedContent())?.classes.find((c) => c.id === id)
}

export function speciesById(id: string): Species | undefined {
  return (content.value ?? cachedContent())?.species.find((s) => s.id === id)
}

/** Test helper: reset the module singleton between tests. */
export function _resetContent() {
  content.value = null
  contentError.value = ''
}