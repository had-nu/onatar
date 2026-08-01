import { beforeEach, describe, expect, it, vi } from 'vitest'
import {
  _resetCharacters,
  buildCharacter,
  createCharacter,
  deleteCharacter,
  getCharacter,
  listCharacters,
  saveCharacter,
  setCampaign,
  setLive,
  starterDraft,
} from './characters.svelte'
import type { Sheet } from './types'

const sheet: Sheet = {
  level: 1,
  hp: { max: 10, current: 10 },
  ac: 13,
  proficiencyBonus: 2,
  spellSlots: [0, 0, 0, 0, 0, 0, 0, 0, 0],
  abilities: {
    STR: { score: 15, mod: 2 },
    DEX: { score: 13, mod: 1 },
    CON: { score: 14, mod: 2 },
    INT: { score: 10, mod: 0 },
    WIS: { score: 12, mod: 1 },
    CHA: { score: 8, mod: -1 },
  },
  features: [],
  pendingChoices: [],
}

beforeEach(() => {
  _resetCharacters()
  localStorage.clear()
  vi.unstubAllGlobals()
})

describe('characters CRUD', () => {
  it('creates, lists and gets a character', () => {
    const c = createCharacter(starterDraft())
    expect(listCharacters()).toHaveLength(1)
    expect(getCharacter(c.id)?.name).toBe('Novo personagem')
    expect(getCharacter('missing')).toBeUndefined()
  })

  it('persists to localStorage', () => {
    const c = createCharacter(starterDraft())
    const raw = JSON.parse(localStorage.getItem('onatar.characters') ?? '[]') as Array<{
      id: string
    }>
    expect(raw).toHaveLength(1)
    expect(raw[0].id).toBe(c.id)
  })

  it('saveCharacter updates fields', () => {
    const c = createCharacter(starterDraft())
    saveCharacter({ ...c, name: 'Bruxa' })
    expect(getCharacter(c.id)?.name).toBe('Bruxa')
  })

  it('deleteCharacter removes the character', () => {
    const c = createCharacter(starterDraft())
    deleteCharacter(c.id)
    expect(listCharacters()).toHaveLength(0)
  })

  it('setLive stores editable sheet state', () => {
    const c = createCharacter(starterDraft())
    const live = {
      hpCurrent: 7,
      slotsUsed: [1, 0, 0, 0, 0, 0, 0, 0, 0],
      conditions: ['cego'],
      resources: {},
    }
    setLive(c.id, live)
    expect(getCharacter(c.id)?.live?.hpCurrent).toBe(7)
    expect(getCharacter(c.id)?.live?.conditions).toEqual(['cego'])
  })

  it('setCampaign links a character to a campaign', () => {
    const c = createCharacter(starterDraft())
    setCampaign(c.id, 'camp-1')
    expect(getCharacter(c.id)?.campaignId).toBe('camp-1')
    setCampaign(c.id, undefined)
    expect(getCharacter(c.id)?.campaignId).toBeUndefined()
  })
})

describe('buildCharacter', () => {
  it('posts the draft and caches the sheet', async () => {
    const c = createCharacter(starterDraft())
    const fetchMock = vi.fn(() =>
      Promise.resolve(new Response(JSON.stringify({ sheet }), { status: 200 }))
    )
    vi.stubGlobal('fetch', fetchMock)

    const result = await buildCharacter(c)
    expect(result.ac).toBe(13)
    expect(fetchMock).toHaveBeenCalledWith(
      '/api/v1/build',
      expect.objectContaining({ method: 'POST' })
    )
    expect(getCharacter(c.id)?.sheet?.level).toBe(1)
  })

  it('falls back to the cached sheet when offline', async () => {
    const c = createCharacter(starterDraft())
    saveCharacter({ ...c, sheet })
    vi.stubGlobal(
      'fetch',
      vi.fn(() => Promise.reject(new TypeError('offline')))
    )

    const result = await buildCharacter(getCharacter(c.id)!)
    expect(result.ac).toBe(13)
  })

  it('throws when offline and no cached sheet exists', async () => {
    const c = createCharacter(starterDraft())
    vi.stubGlobal(
      'fetch',
      vi.fn(() => Promise.reject(new TypeError('offline')))
    )
    await expect(buildCharacter(c)).rejects.toThrow()
  })

  it('surfaces the API error message', async () => {
    const c = createCharacter(starterDraft())
    vi.stubGlobal(
      'fetch',
      vi.fn(() =>
        Promise.resolve(
          new Response(JSON.stringify({ error: { message: 'unknown class "wizard"' } }), {
            status: 422,
          })
        )
      )
    )
    await expect(buildCharacter(c)).rejects.toThrow('unknown class')
  })
})
