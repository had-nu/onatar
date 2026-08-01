import { describe, expect, it } from 'vitest'
import {
  extractTextFromPDF,
  itemsToLines,
  MAX_IMPORT_BYTES,
  parseDDBJSON,
  parseDDBText,
  toBuildRequest,
} from './ddb'
import type { Content } from '../types'

const content: Content = {
  classes: [
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
    {
      id: 'sorcerer',
      name: 'Sorcerer',
      hitDie: 'd6',
      spellcaster: true,
      subclassLevel: 3,
      suggestedSpecies: [],
      suggestedBackgrounds: [],
      data: {},
    },
  ],
  subclasses: [
    {
      id: 'champion',
      classId: 'fighter',
      name: 'Champion',
      levelRequired: 3,
      data: {},
    },
  ],
  species: [
    { id: 'human', name: 'Human', data: {} },
    { id: 'tiefling', name: 'Tiefling', data: {} },
  ],
  backgrounds: [
    { id: 'sage', name: 'Sage', data: {} },
    { id: 'charlatan', name: 'Charlatan', data: {} },
  ],
  spells: [
    { id: 'fire-bolt', name: 'Fire Bolt', level: 0, school: 'Evocation', data: {} },
    { id: 'magic-missile', name: 'Magic Missile', level: 1, school: 'Evocation', data: {} },
    { id: 'shield', name: 'Shield', level: 1, school: 'Abjuration', data: {} },
  ],
  feats: [{ id: 'war-caster', name: 'War Caster', prerequisites: {}, data: {} }],
  features: [],
}

// Layout modelled on the D&D Beyond 2024 sheet (label row + score row).
const DDB_STYLE_TEXT = `Onatar Strongarm
Human
Fighter Level 3
Species: Human
Background: Sage
Armor Class: 14
Hit Points: 28
STRENGTH	DEXTERITY	CONSTITUTION	INTELLIGENCE	WISDOM	CHARISMA
17		10		15		8		12		13
Skills: Arcana, History
Spells: Fire Bolt, Magic Missile, Shield
Feats: War Caster`

describe('parseDDBText', () => {
  it('extracts ability scores from the label+score rows', () => {
    const parsed = parseDDBText(DDB_STYLE_TEXT, content)
    expect(parsed.abilityScores).toEqual({ STR: 17, DEX: 10, CON: 15, INT: 8, WIS: 12, CHA: 13 })
  })

  it('extracts class, level, species and background', () => {
    const parsed = parseDDBText(DDB_STYLE_TEXT, content)
    expect(parsed.classes).toEqual([{ id: 'fighter', level: 3, subclassId: undefined }])
    expect(parsed.speciesId).toBe('human')
    expect(parsed.backgroundId).toBe('sage')
    expect(parsed.name).toBe('Onatar Strongarm')
  })

  it('resolves spells, feats and skills to ids', () => {
    const parsed = parseDDBText(DDB_STYLE_TEXT, content)
    expect(parsed.spellIds).toEqual(
      expect.arrayContaining(['fire-bolt', 'magic-missile', 'shield'])
    )
    expect(parsed.featIds).toEqual(['war-caster'])
    expect(parsed.skillIds).toEqual(['arcana', 'history'])
  })

  it('parses HP and AC', () => {
    const parsed = parseDDBText(DDB_STYLE_TEXT, content)
    expect(parsed.hpMax).toBe(28)
    expect(parsed.ac).toBe(14)
  })

  it('falls back to inline abbreviation matching', () => {
    const flat = 'STR 17 DEX 10 CON 15 INT 8 WIS 12 CHA 13'
    const parsed = parseDDBText(flat, content)
    expect(parsed.abilityScores).toEqual({ STR: 17, DEX: 10, CON: 15, INT: 8, WIS: 12, CHA: 13 })
  })

  it('leaves classes empty when nothing matches content', () => {
    const parsed = parseDDBText('Bard Level 5', content)
    expect(parsed.classes).toEqual([])
  })
})

describe('toBuildRequest', () => {
  it('maps a parsed sheet into a valid draft with defaults', () => {
    const parsed = parseDDBText(DDB_STYLE_TEXT, content)
    const draft = toBuildRequest(parsed)
    expect(draft.name).toBe('Onatar Strongarm')
    expect(draft.classes[0].id).toBe('fighter')
    expect(draft.abilityScores.STR).toBe(17)
    expect(draft.spells).toContain('fire-bolt')
    expect(draft.abilityMethod).toBe('manual')
  })

  it('fills missing ability scores with 10', () => {
    const draft = toBuildRequest({
      name: 'X',
      classes: [],
      abilityScores: { STR: 15 },
      spellIds: [],
      featIds: [],
      skillIds: [],
      unmapped: [],
    })
    expect(draft.abilityScores.DEX).toBe(10)
    expect(draft.abilityScores.CON).toBe(10)
  })
})

describe('parseDDBJSON', () => {
  it('maps a D&D Beyond JSON export into a parsed sheet', async () => {
    const file = new File(
      [
        JSON.stringify({
          character: {
            name: 'Kara',
            race: { name: 'Tiefling' },
            background: { name: 'Sage' },
            classes: [{ name: 'Sorcerer', level: 6, subclass: { name: 'Champion' } }],
            stats: { str: 8, dex: 14, con: 15, int: 10, wis: 12, cha: 18 },
            spells: [{ name: 'Magic Missile' }, { name: 'Shield' }],
            feats: [{ name: 'War Caster' }],
          },
        }),
      ],
      'character.json',
      { type: 'application/json' }
    )
    const parsed = await parseDDBJSON(file, content)
    expect(parsed.name).toBe('Kara')
    expect(parsed.classes).toEqual([{ id: 'sorcerer', level: 6 }])
    expect(parsed.speciesId).toBe('tiefling')
    expect(parsed.backgroundId).toBe('sage')
    expect(parsed.abilityScores).toEqual({ STR: 8, DEX: 14, CON: 15, INT: 10, WIS: 12, CHA: 18 })
    expect(parsed.spellIds).toEqual(['magic-missile', 'shield'])
    expect(parsed.featIds).toEqual(['war-caster'])
  })

  it('rejects malformed JSON', async () => {
    const file = new File(['not json'], 'bad.json', { type: 'application/json' })
    await expect(parseDDBJSON(file, content)).rejects.toThrow(/JSON/)
  })
})

describe('itemsToLines', () => {
  it('groups items into lines and separates columns by x-gap', () => {
    const items = [
      { str: 'STRENGTH', transform: [1, 0, 0, 1, 40, 800] },
      { str: 'DEXTERITY', transform: [1, 0, 0, 1, 200, 800] },
      { str: '17', transform: [1, 0, 0, 1, 40, 785] },
      { str: '10', transform: [1, 0, 0, 1, 200, 785] },
    ]
    const out = itemsToLines(items)
    expect(out.split('\n')).toHaveLength(2)
    expect(out.split('\n')[0]).toContain('STRENGTH\tDEXTERITY')
  })
})

describe('extractTextFromPDF', () => {
  it('rejects oversized files before touching pdf.js', async () => {
    const fake = {
      size: MAX_IMPORT_BYTES + 1,
      arrayBuffer: async () => new Uint8Array(0),
    } as unknown as File
    await expect(extractTextFromPDF(fake)).rejects.toThrow(/grande/)
  })
})
