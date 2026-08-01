import { beforeEach, describe, expect, it } from 'vitest'
import { initRouter, matchRoute, navigate, replaceRoute, route } from './router.svelte'

beforeEach(() => {
  window.location.hash = ''
  replaceRoute({ name: 'home', params: {} })
})

describe('matchRoute', () => {
  it('matches home', () => {
    expect(matchRoute('/')).toEqual({ name: 'home', params: {} })
  })

  it('matches a character route with id', () => {
    expect(matchRoute('/characters/abc')).toEqual({ name: 'character', params: { id: 'abc' } })
  })

  it('decodes the id', () => {
    expect(matchRoute('/characters/a%20b').params.id).toBe('a b')
  })

  it('maps content and builder', () => {
    expect(matchRoute('/content').name).toBe('content')
    expect(matchRoute('/builder').name).toBe('builder')
  })

  it('unknown paths fall back to notfound', () => {
    expect(matchRoute('/nope').name).toBe('notfound')
  })
})

describe('hash navigation', () => {
  it('navigate sets the hash', () => {
    navigate('/content')
    expect(window.location.hash).toBe('#/content')
  })

  it('hashchange updates the route', () => {
    const cleanup = initRouter()
    window.location.hash = '#/characters'
    window.dispatchEvent(new Event('hashchange'))
    expect(route.name).toBe('characters')
    cleanup()
  })
})
