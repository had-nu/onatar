import { beforeEach, describe, expect, it, vi } from 'vitest'
import {
  POINT_BUY_BUDGET,
  POINT_BUY_MAX,
  POINT_BUY_MIN,
  _resetBuilder,
  assignAbility,
  builder,
  canGoNext,
  deselectClass,
  nextStep,
  prevStep,
  redo,
  rollScores,
  saveCharacterFromWizard,
  selectBackground,
  selectClass,
  selectSpecies,
  setClassLevel,
  setMethod,
  setPointBuy,
  setStep,
  step,
  toggleEquipment,
  toggleFeat,
  toggleSkill,
  toggleSpell,
  undo,
  validateStep,
} from './builder.svelte'

// Mock characters store for tests
const {
  mockCharacters,
  mockCreateCharacter,
  mockSaveCharacter,
  mockStarterDraft,
  mockBuildDraft,
  mockResetCharacters,
} = vi.hoisted(() => {
  const mockCharacters = { value: [] as import('./characters.svelte').Character[] }
  return {
    mockCharacters,
    mockCreateCharacter: vi.fn((draft: import('./characters.svelte').BuildRequest) => {
      const c = {
        id: 'test-id',
        name: draft.name || 'Novo personagem',
        isNpc: draft.isNpc ?? false,
        draft,
        createdAt: Date.now(),
        updatedAt: Date.now(),
      }
      mockCharacters.value.push(c)
      return c
    }),
    mockSaveCharacter: vi.fn((c: import('./characters.svelte').Character) => c),
    mockStarterDraft: vi.fn(() => ({
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
    })),
    mockBuildDraft: vi.fn().mockResolvedValue({
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
    }),
    mockResetCharacters: vi.fn(() => {
      mockCharacters.value = []
    }),
  }
})

// Mock buildDraft to avoid network calls in tests
vi.mock('./characters.svelte', () => ({
  buildDraft: mockBuildDraft,
  _resetCharacters: mockResetCharacters,
  createCharacter: mockCreateCharacter,
  saveCharacter: mockSaveCharacter,
  starterDraft: mockStarterDraft,
  characters: mockCharacters,
}))
import { _resetContent } from './content.svelte'
import { starterDraft } from './characters.svelte'
import type { Content } from './types'

const fixture: Content = {
  classes: [
    {
      id: 'sorcerer',
      name: 'Sorcerer',
      hitDie: 'd6',
      spellcaster: true,
      subclassLevel: 3,
      suggestedSpecies: ['tiefling'],
      suggestedBackgrounds: ['sage'],
      data: { primaryAbility: 'CHA', recommendedSpells: ['magic-missile'] },
    },
    {
      id: 'fighter',
      name: 'Fighter',
      hitDie: 'd10',
      spellcaster: false,
      subclassLevel: 3,
      suggestedSpecies: [],
      suggestedBackgrounds: [],
      data: {},
    },
  ],
  subclasses: [
    { id: 'aberrant', classId: 'sorcerer', name: 'Aberrant', levelRequired: 3, data: {} },
  ],
  species: [{ id: 'tiefling', name: 'Tiefling', data: {} }],
  backgrounds: [{ id: 'sage', name: 'Sage', data: { startingGear: ['Book'] } }],
  spells: [{ id: 'magic-missile', name: 'Magic Missile', level: 1, school: 'ev', data: {} }],
  feats: [{ id: 'war-caster', name: 'War Caster', prerequisites: {}, data: {} }],
  features: [],
}

function seedContent() {
  _resetContent()
  localStorage.setItem('onatar.content', JSON.stringify({ at: Date.now(), data: fixture }))
}

beforeEach(() => {
  seedContent()
  _resetBuilder()
  mockResetCharacters()
})

describe('steps', () => {
  it('starts on the class step', () => {
    expect(step().id).toBe('class')
    expect(validateStep('class')).toBe(false)
    expect(canGoNext()).toBe(false)
  })

  it('goes forward and backward', () => {
    selectClass('fighter')
    nextStep()
    expect(step().id).toBe('background')
    prevStep()
    expect(step().id).toBe('class')
  })

  it('does not advance when the step is invalid', () => {
    nextStep()
    expect(step().id).toBe('class')
  })

  it('jumps to a reached step', () => {
    setStep(5)
    expect(step().id).toBe('review')
    setStep(-1)
    expect(step().id).toBe('review')
  })
})

describe('class step', () => {
  it('selects and deselects a class', () => {
    selectClass('fighter')
    expect(builder.value.draft.classes).toHaveLength(1)
    expect(builder.value.draft.classes[0].id).toBe('fighter')
    deselectClass('fighter')
    expect(builder.value.draft.classes).toHaveLength(0)
  })

  it('sets class level', () => {
    selectClass('fighter')
    setClassLevel('fighter', 5)
    expect(builder.value.draft.classes[0].level).toBe(5)
  })

  it('becomes valid once a class is chosen', () => {
    selectClass('fighter')
    expect(validateStep('class')).toBe(true)
    expect(canGoNext()).toBe(true)
  })
})

describe('background / species steps', () => {
  it('selects background and species', () => {
    selectBackground('sage')
    expect(builder.value.draft.backgroundId).toBe('sage')
    selectSpecies('tiefling')
    expect(builder.value.draft.speciesId).toBe('tiefling')
    expect(validateStep('background')).toBe(true)
    expect(validateStep('species')).toBe(true)
  })
})

describe('abilities step', () => {
  it('standard-array: requires all six assigned', () => {
    assignAbility('STR', 15)
    expect(validateStep('abilities')).toBe(false)
    assignAbility('DEX', 14)
    assignAbility('CON', 13)
    assignAbility('INT', 12)
    assignAbility('WIS', 10)
    assignAbility('CHA', 8)
    expect(validateStep('abilities')).toBe(true)
    expect(builder.value.draft.abilityScores.STR).toBe(15)
  })

  it('unassigning a value invalidates the step', () => {
    for (const [ab, v] of [
      ['STR', 15],
      ['DEX', 14],
      ['CON', 13],
      ['INT', 12],
      ['WIS', 10],
      ['CHA', 8],
    ] as const) {
      assignAbility(ab, v)
    }
    assignAbility('STR', null)
    expect(validateStep('abilities')).toBe(false)
  })

  it('point-buy enforces the budget', () => {
    setMethod('point-buy')
    for (const ab of ['STR', 'DEX', 'CON', 'INT', 'WIS', 'CHA'] as const) {
      setPointBuy(ab, POINT_BUY_MAX)
    }
    expect(validateStep('abilities')).toBe(false)
    for (const ab of ['STR', 'DEX', 'CON', 'INT', 'WIS', 'CHA'] as const) {
      setPointBuy(ab, POINT_BUY_MIN)
    }
    expect(validateStep('abilities')).toBe(true)
  })

  it('rolled: rollScores produces six scores and can be assigned', () => {
    setMethod('rolled')
    const scores = rollScores()
    expect(scores).toHaveLength(6)
    expect(builder.value.abilities.rolled).toEqual(scores)
    expect(builder.value.abilities.rolled.every((v) => v >= 3 && v <= 18)).toBe(true)
    scores.forEach((v, i) =>
      assignAbility(['STR', 'DEX', 'CON', 'INT', 'WIS', 'CHA'][i] as never, v)
    )
    expect(validateStep('abilities')).toBe(true)
  })
})

describe('toggles', () => {
  it('toggles spells, skills, feats and equipment', () => {
    toggleSpell('magic-missile')
    toggleSpell('magic-missile')
    expect(builder.value.draft.spells).toEqual([])
    toggleSpell('magic-missile')
    expect(builder.value.draft.spells).toEqual(['magic-missile'])

    toggleSkill('arcana')
    expect(builder.value.draft.skills).toEqual(['arcana'])

    toggleFeat('war-caster')
    expect(builder.value.draft.feats).toEqual(['war-caster'])

    toggleEquipment('Book')
    expect(builder.value.equipment).toEqual(['Book'])
    toggleEquipment('Book')
    expect(builder.value.equipment).toEqual([])
  })
})

describe('undo / redo', () => {
  it('undoes and redoes a mutation', () => {
    selectClass('fighter')
    selectSpecies('tiefling')
    expect(builder.value.draft.speciesId).toBe('tiefling')
    undo()
    expect(builder.value.draft.speciesId).toBeUndefined()
    expect(builder.value.draft.classes).toHaveLength(1)
    redo()
    expect(builder.value.draft.speciesId).toBe('tiefling')
  })

  it('undo returns the step index', () => {
    selectClass('fighter')
    nextStep()
    expect(step().id).toBe('background')
    undo()
    expect(step().id).toBe('class')
  })

  it('no-ops when there is no history', () => {
    undo()
    expect(step().id).toBe('class')
    redo()
    expect(step().id).toBe('class')
  })
})

describe('save', () => {
  function buildValidWizard() {
    selectClass('sorcerer')
    setClassLevel('sorcerer', 3)
    selectBackground('sage')
    selectSpecies('tiefling')
    assignAbility('STR', 8)
    assignAbility('DEX', 14)
    assignAbility('CON', 13)
    assignAbility('INT', 12)
    assignAbility('WIS', 10)
    assignAbility('CHA', 15)
    toggleSpell('magic-missile')
    toggleFeat('war-caster')
    toggleEquipment('Book')
    setStep(5)
    builder.value.name = 'Onatar'
  }

  it('saves the character and clears via createCharacter', async () => {
    buildValidWizard()
    expect(validateStep('review')).toBe(true)
    expect(canGoNext()).toBe(true)
    const c = await saveCharacterFromWizard()
    expect(c).not.toBeNull()
    expect(c!.draft.classes[0].id).toBe('sorcerer')
    expect(c!.draft.spells).toEqual(['magic-missile'])
    expect(c!.draft.equipment).toEqual(['Book'])
    expect(mockCharacters.value).toHaveLength(1)
  })

  it('does not save without a name', async () => {
    buildValidWizard()
    builder.value.name = ''
    expect(await saveCharacterFromWizard()).toBeNull()
  })
})

it('starterDraft is a valid prebuilt draft', () => {
  expect(starterDraft().classes[0].id).toBe('fighter')
  expect(Object.keys(starterDraft().abilityScores)).toHaveLength(6)
  expect(POINT_BUY_BUDGET).toBe(27)
})
