// Theme state: 'light' | 'dark' | 'system', persisted in localStorage and
// applied via a `data-theme` attribute on <html> (see app.css).
export type Theme = 'light' | 'dark' | 'system'

export const THEMES: Theme[] = ['light', 'dark', 'system']

const KEY = 'onatar.theme'

import { box } from './box.svelte'

function readStored(): Theme {
  try {
    const v = localStorage.getItem(KEY) as Theme | null
    return v && THEMES.includes(v) ? v : 'system'
  } catch {
    return 'system'
  }
}

export const theme = box<Theme>(readStored())

export function setTheme(t: Theme) {
  theme.value = t
  try {
    localStorage.setItem(KEY, t)
  } catch {
    /* storage unavailable — ignore */
  }
  applyTheme(t)
}

export function cycleTheme() {
  setTheme(THEMES[(THEMES.indexOf(theme.value) + 1) % THEMES.length])
}

export function applyTheme(t: Theme) {
  const el = document.documentElement
  if (t === 'system') {
    el.removeAttribute('data-theme')
  } else {
    el.setAttribute('data-theme', t)
  }
}

export function initTheme() {
  applyTheme(theme.value)
}
