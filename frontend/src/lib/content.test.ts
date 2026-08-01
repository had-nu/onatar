import { beforeEach, describe, expect, it, vi } from 'vitest'
import { _resetContent, content, contentError, loadContent } from './content.svelte'
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
      suggestedBackgrounds: [],
      data: {},
    },
  ],
  subclasses: [],
  species: [],
  backgrounds: [],
  spells: [],
  feats: [],
  features: [],
}

function jsonResponse(body: unknown, status = 200): Response {
  return new Response(JSON.stringify(body), {
    status,
    headers: { 'Content-Type': 'application/json' },
  })
}

beforeEach(() => {
  _resetContent()
  localStorage.clear()
  vi.unstubAllGlobals()
})

describe('loadContent', () => {
  it('fetches from the API and caches it', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn(() => Promise.resolve(jsonResponse(fixture)))
    )
    const data = await loadContent()
    expect(data.classes[0].id).toBe('sorcerer')
    expect(content.value).toEqual(fixture)
    expect(localStorage.getItem('onatar.content')).toBeTruthy()
  })

  it('returns the cached content on a second call without fetching', async () => {
    const fetchMock = vi.fn(() => Promise.resolve(jsonResponse(fixture)))
    vi.stubGlobal('fetch', fetchMock)
    await loadContent()
    await loadContent()
    expect(fetchMock).toHaveBeenCalledTimes(1)
  })

  it('falls back to the localStorage cache when offline', async () => {
    localStorage.setItem('onatar.content', JSON.stringify({ at: Date.now(), data: fixture }))
    vi.stubGlobal(
      'fetch',
      vi.fn(() => Promise.reject(new TypeError('network down')))
    )
    const data = await loadContent()
    expect(data.classes[0].id).toBe('sorcerer')
  })

  it('throws and records the error when offline with no cache', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn(() => Promise.reject(new TypeError('network down')))
    )
    await expect(loadContent()).rejects.toThrow()
    expect(contentError.value).toBeTruthy()
  })

  it('throws on a non-ok response with no cache', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn(() => Promise.resolve(jsonResponse('oops', 500)))
    )
    await expect(loadContent()).rejects.toThrow()
  })
})
