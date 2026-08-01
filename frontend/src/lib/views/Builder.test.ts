import { fireEvent, render, screen } from '@testing-library/svelte'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import Builder from './Builder.svelte'
import { _resetBuilder } from '../builder.svelte'
import { _resetContent } from '../content.svelte'
import { _resetCharacters } from '../characters.svelte'
import type { Content } from '../types'

const fixture: Content = {
  classes: [
    {
      id: 'fighter',
      name: 'Fighter',
      hitDie: 'd10',
      spellcaster: false,
      subclassLevel: 3,
      suggestedSpecies: [],
      suggestedBackgrounds: [],
      data: { description: 'Master of arms.' },
    },
  ],
  subclasses: [],
  species: [{ id: 'human', name: 'Human', data: { description: 'Versatile.' } }],
  backgrounds: [{ id: 'sage', name: 'Sage', data: { description: 'Bookish.' } }],
  spells: [],
  feats: [],
  features: [],
}

const sheet = {
  level: 1,
  hp: { max: 12, current: 12 },
  ac: 14,
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
  _resetContent()
  _resetBuilder()
  _resetCharacters()
  localStorage.clear()
  vi.unstubAllGlobals()
  vi.stubGlobal(
    'fetch',
    vi.fn((input: RequestInfo | URL) => {
      const url = String(input)
      if (url.includes('/api/v1/content')) {
        return Promise.resolve(
          new Response(JSON.stringify(fixture), {
            status: 200,
            headers: { 'Content-Type': 'application/json' },
          })
        )
      }
      if (url.includes('/api/v1/build')) {
        return Promise.resolve(
          new Response(JSON.stringify({ sheet }), {
            status: 200,
            headers: { 'Content-Type': 'application/json' },
          })
        )
      }
      return Promise.resolve(new Response('{}', { status: 404 }))
    })
  )
})

describe('Builder wizard', () => {
  it('renders the six steps and the class step first', async () => {
    render(Builder)
    expect(await screen.findByText('Classe')).toBeTruthy()
    expect(screen.getByRole('heading', { name: 'Escolhe as classes' })).toBeTruthy()
    expect(screen.getByText('Fighter')).toBeTruthy()
    expect(screen.getByText('Pré-visualização')).toBeTruthy()
  })

  it('keeps Continuar disabled until a class is selected', async () => {
    render(Builder)
    await screen.findByText('Fighter')
    const cont = screen.getByRole('button', { name: 'Continuar →' })
    expect((cont as HTMLButtonElement).disabled).toBe(true)
    await fireEvent.click(screen.getByRole('button', { name: 'Escolher' }))
    expect(
      (screen.getByRole('button', { name: 'Continuar →' }) as HTMLButtonElement).disabled
    ).toBe(false)
  })

  it('navigates to the background step after selecting a class', async () => {
    render(Builder)
    await screen.findByText('Fighter')
    await fireEvent.click(screen.getByRole('button', { name: 'Escolher' }))
    await fireEvent.click(screen.getByRole('button', { name: 'Continuar →' }))
    expect(await screen.findByText('Escolhe o background')).toBeTruthy()
    expect(screen.getByRole('button', { name: '← Anterior' })).toBeTruthy()
  })
})
