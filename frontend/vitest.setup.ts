// vitest setup — mock fetch for API endpoints used in tests
import { vi } from 'vitest'

// Mock content fixture for tests
const mockContent = {
  classes: [
    {
      id: 'fighter',
      name: 'Fighter',
      hitDie: 'd10',
      spellcaster: false,
      subclassLevel: 3,
      suggestedSpecies: ['human', 'dwarf'],
      suggestedBackgrounds: ['soldier', 'folk-hero'],
      data: { primaryAbility: 'STR', description: 'A master of martial combat' },
    },
    {
      id: 'sorcerer',
      name: 'Sorcerer',
      hitDie: 'd6',
      spellcaster: true,
      subclassLevel: 3,
      suggestedSpecies: ['tiefling', 'human'],
      suggestedBackgrounds: ['charlatan', 'sage'],
      data: { primaryAbility: 'CHA', description: 'Magic from within' },
    },
  ],
  subclasses: [
    { id: 'champion', classId: 'fighter', name: 'Champion', levelRequired: 3, data: {} },
    { id: 'draconic', classId: 'sorcerer', name: 'Draconic Bloodline', levelRequired: 3, data: {} },
  ],
  species: [
    { id: 'human', name: 'Human', data: { traits: ['Versatile'], abilityScores: { STR: 1 } } },
    {
      id: 'tiefling',
      name: 'Tiefling',
      data: { traits: ['Darkvision', 'Hellish Resistance'], abilityScores: { CHA: 2, INT: 1 } },
    },
  ],
  backgrounds: [
    { id: 'sage', name: 'Sage', data: { skills: ['arcana', 'history'], feature: 'Researcher' } },
    {
      id: 'soldier',
      name: 'Soldier',
      data: { skills: ['athletics', 'intimidation'], feature: 'Military Rank' },
    },
  ],
  spells: [
    { id: 'fire-bolt', name: 'Fire Bolt', level: 0, school: 'evocation', data: {} },
    { id: 'magic-missile', name: 'Magic Missile', level: 1, school: 'evocation', data: {} },
    { id: 'shield', name: 'Shield', level: 1, school: 'abjuration', data: {} },
  ],
  feats: [
    { id: 'war-caster', name: 'War Caster', prerequisites: {}, data: {} },
    {
      id: 'ability-score-improvement',
      name: 'Ability Score Improvement',
      prerequisites: {},
      data: {},
    },
  ],
  features: [],
}

// Global fetch mock
const originalFetch = global.fetch

beforeEach(() => {
  vi.useFakeTimers()
  global.fetch = vi.fn(async (url: string | URL | Request) => {
    const urlStr = typeof url === 'string' ? url : url.toString()

    if (urlStr.includes('/api/v1/content')) {
      return new Response(JSON.stringify(mockContent), {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
      })
    }

    if (urlStr.includes('/api/v1/build')) {
      return new Response(
        JSON.stringify({
          sheet: {
            level: 1,
            hp: { max: 10, current: 10 },
            ac: 12,
            proficiencyBonus: 2,
            spellSlots: [2, 0, 0, 0, 0, 0, 0, 0, 0],
            abilities: {
              STR: { score: 8, mod: -1 },
              DEX: { score: 14, mod: 2 },
              CON: { score: 13, mod: 1 },
              INT: { score: 12, mod: 1 },
              WIS: { score: 10, mod: 0 },
              CHA: { score: 15, mod: 2 },
            },
            features: [],
            pendingChoices: [],
          },
        }),
        { status: 200, headers: { 'Content-Type': 'application/json' } }
      )
    }

    // Fallback to original fetch for other requests
    return originalFetch(url)
  })
})

afterEach(() => {
  vi.useRealTimers()
  global.fetch = originalFetch
})
