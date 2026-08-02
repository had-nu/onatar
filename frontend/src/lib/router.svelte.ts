// Minimal hash router (no dependencies). Routes are derived from
// `window.location.hash`; the whole SPA lives under one document, so only the
// hash changes (PRD §6 C2: Svelte SPA + Go API).
import { tick } from 'svelte'

export type RouteName =
  | 'home'
  | 'login'
  | 'characters'
  | 'character'
  | 'content'
  | 'builder'
  | 'campaigns'
  | 'import'
  | 'combat'
  | 'notfound'

export interface MatchedRoute {
  name: RouteName
  params: Record<string, string>
}

const patterns: Record<RouteName, RegExp> = {
  home: /^\/$/,
  login: /^\/login$/,
  characters: /^\/characters$/,
  character: /^\/characters\/([^/]+)$/,
  content: /^\/content$/,
  builder: /^\/builder$/,
  campaigns: /^\/campaigns$/,
  import: /^\/import$/,
  combat: /^\/combat$/,
  notfound: /$^/,
}

export function matchRoute(path: string): MatchedRoute {
  for (const name of Object.keys(patterns) as RouteName[]) {
    if (name === 'notfound') continue
    const m = patterns[name].exec(path)
    if (m) {
      const params: Record<string, string> = {}
      if (name === 'character' && m[1]) params.id = decodeURIComponent(m[1])
      return { name, params }
    }
  }
  return { name: 'notfound', params: {} }
}

function currentHashPath(): string {
  if (typeof window === 'undefined') return '/'
  return (window.location.hash.replace(/^#/, '') || '/').split('?')[0]
}

export const route = $state<MatchedRoute>(matchRoute(currentHashPath()))

/** Navigate to a path; updates the hash (and thus `route`). */
export function navigate(path: string) {
  if (currentHashPath() === path) return
  window.location.hash = path
}

/** Programmatic navigation used for tests / non-anchor flows. */
export function replaceRoute(m: MatchedRoute) {
  route.name = m.name
  route.params = m.params
}

/** Wire hashchange -> route. Returns an unsubscribe function. */
export function initRouter(): () => void {
  const update = () => {
    route.name = matchRoute(currentHashPath()).name
    route.params = matchRoute(currentHashPath()).params
  }
  window.addEventListener('hashchange', update)
  // Ensure initial route matches even if the URL was set before mount.
  void tick().then(update)
  return () => window.removeEventListener('hashchange', update)
}
