import { render, screen } from '@testing-library/svelte'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { _resetCharacters, createCharacter, starterDraft } from '../characters.svelte'
import Characters from './Characters.svelte'

beforeEach(() => {
  _resetCharacters()
  localStorage.clear()
  vi.restoreAllMocks()
})

describe('Characters view', () => {
  it('shows the empty state', () => {
    render(Characters)
    expect(screen.getByText(/ainda não tens personagens/i)).toBeTruthy()
    expect(screen.getByRole('button', { name: 'Criar o primeiro' })).toBeTruthy()
  })

  it('lists saved characters with a class summary', () => {
    createCharacter({ ...starterDraft(), name: 'Bruxa' })
    render(Characters)
    expect(screen.getByText('Bruxa')).toBeTruthy()
    expect(screen.getByText(/fighter 1/i)).toBeTruthy()
  })

  it('deletes a character after confirmation', async () => {
    createCharacter({ ...starterDraft(), name: 'Bruxa' })
    window.confirm = vi.fn(() => true)
    render(Characters)

    const remove = screen.getByRole('button', { name: 'Apagar' })
    await remove.click()

    expect(window.confirm).toHaveBeenCalled()
    expect(screen.queryByText('Bruxa')).toBeNull()
  })

  it('keeps the character when confirmation is denied', async () => {
    createCharacter({ ...starterDraft(), name: 'Bruxa' })
    window.confirm = vi.fn(() => false)
    render(Characters)

    const remove = screen.getByRole('button', { name: 'Apagar' })
    await remove.click()

    expect(screen.getByText('Bruxa')).toBeTruthy()
  })
})
