// Content store: loads GET /api/v1/content once, keeps an in-memory copy and a
// localStorage cache so the app degrades gracefully offline (PRD RNF-03).
import { box } from './box.svelte'
import type { Class, Content, Species } from './types'

const KEY = 'onatar.content'

export const content = box<Content | null>(null)
export const contentError = box<string>('')
interface CacheEntry {
  at: number
  data: Content
}

function readCache(): Content | null {
  try {
    const raw = localStorage.getItem(KEY)
    if (!raw) return null
    const entry = JSON.parse(raw) as CacheEntry
    if (!entry?.data?.classes) return null
    return entry.data
  } catch {
    return null
  }
}

function writeCache(data: Content) {
  try {
    const entry: CacheEntry = { at: Date.now(), data }
    localStorage.setItem(KEY, JSON.stringify(entry))
  } catch {
    /* storage unavailable — ignore */
  }
}

/** Load content from the API, falling back to the localStorage cache. */
export async function loadContent(force = false): Promise<Content> {
  if (content.value && !force) return content.value
  contentError.value = ''
  try {
    const res = await fetch('/api/v1/content')
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    const data = (await res.json()) as Content
    content.value = data
    writeCache(data)
    return data
  } catch (err) {
    const cached = readCache()
    if (cached) {
      content.value = cached
      return cached
    }
    contentError.value = err instanceof Error ? err.message : String(err)
    throw err
  }
}

/** Synchronous access to whatever content is already loaded/cached. */
export function cachedContent(): Content | null {
  if (content.value) return content.value
  const cached = readCache()
  if (cached) content.value = cached
  return cached
}

export function classById(id: string): Class | undefined {
  return (content.value ?? cachedContent())?.classes.find((c) => c.id === id)
}

export function speciesById(id: string): Species | undefined {
  return (content.value ?? cachedContent())?.species.find((s) => s.id === id)
}

/** Test helper: reset the module singleton between tests. */
export function _resetContent() {
  content.value = null
  contentError.value = ''
}
