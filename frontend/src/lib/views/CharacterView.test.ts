import { fireEvent, render, screen } from '@testing-library/svelte'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import CharacterView from './CharacterView.svelte'
import {
  _resetCharacters,
  createCharacter,
  getCharacter,
  saveCharacter,
  starterDraft,
} from '../characters.svelte'
import { _resetCampaigns, createCampaign } from '../campaigns.svelte'
import type { Sheet } from '../types'

const sheet: Sheet = {
  level: 3,
  hp: { max: 20, current: 20 },
  ac: 15,
  proficiencyBonus: 2,
  spellSlots: [4, 2, 0, 0, 0, 0, 0, 0, 0],
  abilities: {
    STR: { score: 10, mod: 0 },
    DEX: { score: 14, mod: 2 },
    CON: { score: 13, mod: 1 },
    INT: { score: 12, mod: 1 },
    WIS: { score: 10, mod: 0 },
    CHA: { score: 8, mod: -1 },
  },
  features: [{ name: 'Font of Magic', level: 2, description: 'Sorc points.' }],
  pendingChoices: [],
}

async function seededCharacter() {
  _resetCharacters()
  const c = await createCharacter(starterDraft())
  const withSheet = { ...getCharacter(c.id)!, sheet }
  await saveCharacter(withSheet)
  return withSheet
}

beforeEach(() => {
  localStorage.clear()
  _resetCampaigns()
  vi.unstubAllGlobals()
  vi.stubGlobal(
    'fetch',
    vi.fn(() => Promise.reject(new TypeError('offline')))
  )
})

describe('CharacterView interactive sheet', () => {
  it('renders the read-only stats plus editable HP', async () => {
    const c = await seededCharacter()
    render(CharacterView, { props: { id: c.id } })
    expect(screen.getByRole('heading', { name: 'New character' })).toBeTruthy()
    expect(screen.getByText('HP / 20')).toBeTruthy()
    expect(screen.getByText('AC')).toBeTruthy()
    expect(getCharacter(c.id)?.live?.hpCurrent).toBe(20)
  })

  it('decrements and persists HP via the stepper', async () => {
    const c = await seededCharacter()
    render(CharacterView, { props: { id: c.id } })
    await fireEvent.click(screen.getByRole('button', { name: 'Decrease HP' }))
    expect(screen.getByText('19')).toBeTruthy()
    expect(getCharacter(c.id)?.live?.hpCurrent).toBe(19)
  })

  it('marks a spell slot as used and persists it', async () => {
    const c = await seededCharacter()
    render(CharacterView, { props: { id: c.id } })
    await fireEvent.click(screen.getByRole('button', { name: 'Mark level 1 as used' }))
    expect(screen.getByText('1/4')).toBeTruthy()
    expect(getCharacter(c.id)?.live?.slotsUsed[0]).toBe(1)
  })

  it('toggles a condition', async () => {
    const c = await seededCharacter()
    render(CharacterView, { props: { id: c.id } })
    await fireEvent.click(screen.getByRole('button', { name: 'blinded' }))
    expect(getCharacter(c.id)?.live?.conditions).toContain('blinded')
  })

  it('assigns a character to a campaign', async () => {
    const c = await seededCharacter()
    const camp = await createCampaign('Avernus')
    render(CharacterView, { props: { id: c.id } })
    const select = screen.getByRole('combobox', { name: /campaign/i })
    await fireEvent.change(select, { target: { value: camp.id } })
    expect(getCharacter(c.id)?.campaignId).toBe(camp.id)
  })
})