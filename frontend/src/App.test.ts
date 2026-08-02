import { render, screen } from '@testing-library/svelte'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import App from './App.svelte'
import { _resetCharacters } from './lib/characters.svelte'
import { _resetContent } from './lib/content.svelte'

describe('App', () => {
  beforeEach(() => {
    window.location.hash = ''
    localStorage.clear()
    _resetCharacters()
    _resetContent()
    vi.unstubAllGlobals()
  })

  it('renders the landing page with navigation', () => {
    render(App)
    const heading = screen.getByRole('heading', { level: 1 })
    expect(heading.textContent).toBe('Forge Your Hero')
    expect(screen.getByRole('navigation', { name: 'Main navigation' })).toBeTruthy()
    // When not authenticated, "Characters" is hidden, "Content" is always visible
    expect(screen.getByRole('link', { name: 'Content' })).toBeTruthy()
    // CTA buttons on landing page
    expect(screen.getByRole('link', { name: 'Create Character' })).toBeTruthy()
    expect(screen.getByRole('link', { name: 'My Characters' })).toBeTruthy()
    expect(screen.getByRole('link', { name: 'Explore Content' })).toBeTruthy()
  })
})