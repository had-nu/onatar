import { render, screen } from '@testing-library/svelte'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { _resetContent } from '../content.svelte'
import type { Content as ContentData } from '../types'
import ContentView from './Content.svelte'

const fixture: ContentData = {
  classes: [
    {
      id: 'sorcerer',
      name: 'Sorcerer',
      hitDie: 'd6',
      spellcaster: true,
      subclassLevel: 3,
      suggestedSpecies: ['tiefling'],
      suggestedBackgrounds: [],
      data: { description: 'Inherent magic.', primaryAbility: 'CHA' },
    },
  ],
  subclasses: [],
  species: [
    {
      id: 'tiefling',
      name: 'Tiefling',
      data: { description: 'Infernal blood.', traits: [{ name: 'Darkvision' }] },
    },
  ],
  backgrounds: [
    { id: 'sage', name: 'Sage', data: { description: 'Bookish.', skills: ['arcana'] } },
  ],
  spells: [],
  feats: [],
  features: [],
}

beforeEach(() => {
  _resetContent()
  localStorage.clear()
  vi.unstubAllGlobals()
  vi.stubGlobal(
    'fetch',
    vi.fn(() =>
      Promise.resolve(
        new Response(JSON.stringify(fixture), {
          status: 200,
          headers: { 'Content-Type': 'application/json' },
        })
      )
    )
  )
})

describe('Content view', () => {
  it('renders class cards with suggestions', async () => {
    render(ContentView)
    expect(await screen.findByText('Sorcerer')).toBeTruthy()
    expect(screen.getByText('Espécies sugeridas')).toBeTruthy()
    expect(screen.getByText('tiefling')).toBeTruthy()
  })

  it('switches to the species tab', async () => {
    render(ContentView)
    await screen.findByText('Sorcerer')
    const tab = screen.getByRole('button', { name: 'Espécies' })
    await tab.click()
    expect(await screen.findByText('Tiefling')).toBeTruthy()
  })

  it('shows an error state when the request fails', async () => {
    vi.unstubAllGlobals()
    vi.stubGlobal(
      'fetch',
      vi.fn(() => Promise.reject(new TypeError('network down')))
    )
    render(ContentView)
    expect(await screen.findByText(/não foi possível carregar/i)).toBeTruthy()
  })
})
