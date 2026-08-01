import { beforeEach, describe, expect, it } from 'vitest'
import { applyTheme, cycleTheme, setTheme, theme } from './theme.svelte'

beforeEach(() => {
  setTheme('system')
})

describe('theme', () => {
  it('applies data-theme attribute', () => {
    setTheme('dark')
    expect(document.documentElement.getAttribute('data-theme')).toBe('dark')
    setTheme('light')
    expect(document.documentElement.getAttribute('data-theme')).toBe('light')
  })

  it('system removes the attribute (follow OS)', () => {
    setTheme('dark')
    setTheme('system')
    expect(document.documentElement.hasAttribute('data-theme')).toBe(false)
  })

  it('persists the choice in localStorage', () => {
    setTheme('dark')
    expect(localStorage.getItem('onatar.theme')).toBe('dark')
  })

  it('cycles light -> dark -> system -> light', () => {
    setTheme('light')
    cycleTheme()
    expect(theme.value).toBe('dark')
    cycleTheme()
    expect(theme.value).toBe('system')
    cycleTheme()
    expect(theme.value).toBe('light')
  })

  it('applyTheme handles each mode', () => {
    applyTheme('dark')
    expect(document.documentElement.dataset.theme).toBe('dark')
    applyTheme('light')
    expect(document.documentElement.dataset.theme).toBe('light')
    applyTheme('system')
    expect(document.documentElement.dataset.theme).toBeUndefined()
  })
})
