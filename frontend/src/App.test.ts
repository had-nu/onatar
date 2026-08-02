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
    expect(heading.textContent).toBe('Onatar')
    expect(screen.getByRole('navigation', { name: 'Navegação principal' })).toBeTruthy()
    // When not authenticated, "Personagens" is hidden, "Conteúdo" is always visible
    expect(screen.getByRole('link', { name: 'Conteúdo' })).toBeTruthy()
    // CTA buttons on landing page
    expect(screen.getByRole('link', { name: 'Criar personagem' })).toBeTruthy()
    expect(screen.getByRole('link', { name: 'Os meus personagens' })).toBeTruthy()
    expect(screen.getByRole('link', { name: 'Explorar conteúdo' })).toBeTruthy()
  })
})