// Characters store (Fase 3): CRUD with dual-mode persistence.
// - Guest: localStorage (RF-03)
// - Authenticated: MariaDB via API + localStorage cache for offline (RNF-03)
// Migration: on first login, localStorage chars are pushed to API.
import { box } from './box.svelte'
import { isAuthenticated, getUser } from './auth.svelte'
import type { BuildRequest, Character, Sheet, SheetLive, User } from './types'

export type { BuildRequest, Character, Sheet, SheetLive } from './types'

const KEY = 'onatar.characters'
const SYNC_KEY = 'onatar.characters.synced' // tracks which local IDs have been synced

const API_BASE = '/api/v1'

interface ApiCharacter {
  id: string
  name: string
  isNpc: boolean
  campaignId?: string
  draft: BuildRequest
  sheet?: Sheet
  live?: SheetLive
  createdAt: number
  updatedAt: number
}

function readAll(): Character[] {
  try {
    const raw = localStorage.getItem(KEY)
    if (!raw) return []
    const arr = JSON.parse(raw) as unknown
    return Array.isArray(arr) ? (arr as Character[]) : []
  } catch {
    return []
  }
}

function readSyncedIds(): Set<string> {
  try {
    const raw = localStorage.getItem(SYNC_KEY)
    if (!raw) return new Set()
    const arr = JSON.parse(raw) as unknown
    return Array.isArray(arr) ? new Set(arr as string[]) : new Set()
  } catch {
    return new Set()
  }
}

function writeSyncedIds(ids: Set<string>) {
  try {
    localStorage.setItem(SYNC_KEY, JSON.stringify(Array.from(ids)))
  } catch {
    /* ignore */
  }
}

export const characters = box<Character[]>(readAll())

let currentUser: User | null = null
let syncInProgress = false

function persist() {
  try {
    localStorage.setItem(KEY, JSON.stringify(characters.value))
  } catch {
    /* storage unavailable — ignore */
  }
}

function toLocal(apiChar: ApiCharacter): Character {
  return {
    id: apiChar.id,
    name: apiChar.name,
    isNpc: apiChar.isNpc,
    campaignId: apiChar.campaignId,
    draft: apiChar.draft,
    sheet: apiChar.sheet,
    live: apiChar.live,
    createdAt: apiChar.createdAt,
    updatedAt: apiChar.updatedAt,
  }
}

async function fetchJson<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(path, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options?.headers,
    },
    credentials: 'include',
  })
  if (!res.ok) {
    let msg = `HTTP ${res.status}`
    try {
      const body = (await res.json()) as { error?: { message?: string } }
      if (body?.error?.message) msg = body.error.message
    } catch {
      /* ignore */
    }
    throw new Error(msg)
  }
  return res.json() as Promise<T>
}

export async function loadFromApi(): Promise<void> {
  if (syncInProgress) return
  syncInProgress = true
  try {
    const apiChars = await fetchJson<ApiCharacter[]>(`${API_BASE}/characters`)
    const localChars = characters.value
    const localById = new Map(localChars.map((c) => [c.id, c]))

    // Merge: API is source of truth for cloud, but preserve local unsynced changes
    const merged = apiChars.map((ac) => {
      const local = localById.get(ac.id)
      if (local && local.updatedAt > ac.updatedAt) {
        // Local is newer — will be pushed on next save
        return { ...toLocal(ac), _pendingPush: true }
      }
      return toLocal(ac)
    })

    // Add local-only characters (not yet synced)
    const syncedIds = readSyncedIds()
    for (const local of localChars) {
      if (!syncedIds.has(local.id) && !apiChars.some((ac) => ac.id === local.id)) {
        merged.push({ ...local, _localOnly: true })
      }
    }

    characters.value = merged
    persist()
  } catch (err) {
    console.warn('Failed to load characters from API:', err)
  } finally {
    syncInProgress = false
  }
}

export async function migrateLocalToApi(): Promise<void> {
  const user = getUser()
  if (!user) return

  const localChars = characters.value
  const syncedIds = readSyncedIds()

  for (const char of localChars) {
    if (syncedIds.has(char.id)) continue

    try {
      const draftBytes = JSON.stringify(char.draft)
      const sheetBytes = char.sheet ? JSON.stringify(char.sheet) : undefined
      const liveBytes = char.live ? JSON.stringify(char.live) : undefined

      const res = await fetch(`${API_BASE}/characters`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({
          name: char.name,
          draft: char.draft,
          isNpc: char.isNpc,
        }),
      })

      if (!res.ok) {
        console.warn(`Failed to migrate character ${char.id}:`, await res.text())
        continue
      }

      const created = (await res.json()) as ApiCharacter
      // Update local ID to match server ID if different
      if (created.id !== char.id) {
        characters.value = characters.value.map((c) =>
          c.id === char.id ? { ...c, id: created.id } : c
        )
      }
      syncedIds.add(created.id)
      writeSyncedIds(syncedIds)
    } catch (err) {
      console.warn(`Migration error for ${char.id}:`, err)
    }
  }

  // Reload from API to get canonical state
  await loadFromApi()
}

export function setCurrentUser(user: User | null) {
  currentUser = user
}

export function isSyncInProgress(): boolean {
  return syncInProgress
}

export function newId(): string {
  if (typeof crypto !== 'undefined' && typeof crypto.randomUUID === 'function') {
    return crypto.randomUUID()
  }
  return `c_${Date.now()}_${Math.random().toString(36).slice(2, 8)}`
}

export function listCharacters(): Character[] {
  return characters.value
}

export function getCharacter(id: string): Character | undefined {
  return characters.value.find((c) => c.id === id)
}

export function starterDraft(): BuildRequest {
  return {
    name: 'New character',
    classes: [{ id: 'fighter', level: 1 }],
    speciesId: 'human',
    backgroundId: 'sage',
    abilityScores: { STR: 15, DEX: 13, CON: 14, INT: 10, WIS: 12, CHA: 8 },
    abilityMethod: 'standard-array',
    skills: [],
    spells: [],
    feats: [],
    isNpc: false,
  }
}

export async function createCharacter(draft: BuildRequest): Promise<Character> {
  const now = Date.now()
  const isAuthed = isAuthenticated()

  if (isAuthed) {
    try {
      const res = await fetch(`${API_BASE}/characters`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({
          name: draft.name || 'Novo personagem',
          draft,
          isNpc: draft.isNpc ?? false,
        }),
      })

      if (!res.ok) throw new Error(await res.text())

      const created = (await res.json()) as ApiCharacter
      const local = toLocal(created)
      characters.value = [...characters.value, local]
      persist()
      writeSyncedIds(new Set([...readSyncedIds(), created.id]))
      return local
    } catch (err) {
      console.warn('API create failed, falling back to local:', err)
      // Fall through to local
    }
  }

  // Local-only (guest or API failure fallback)
  const c: Character = {
    id: newId(),
    name: draft.name || 'Novo personagem',
    isNpc: draft.isNpc ?? false,
    draft,
    createdAt: now,
    updatedAt: now,
  }
  characters.value = [...characters.value, c]
  persist()
  return c
}

export async function saveCharacter(c: Character): Promise<void> {
  const updated: Character = { ...c, updatedAt: Date.now() }
  characters.value = [...characters.value.filter((x) => x.id !== c.id), updated]
  persist()

  const isAuthed = isAuthenticated()
  const syncedIds = readSyncedIds()

  if (isAuthed && syncedIds.has(c.id)) {
    try {
      const draftBytes = JSON.stringify(updated.draft)
      const sheetBytes = updated.sheet ? JSON.stringify(updated.sheet) : undefined
      const liveBytes = updated.live ? JSON.stringify(updated.live) : undefined

      await fetch(`${API_BASE}/characters/${c.id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({
          name: updated.name,
          draft: updated.draft,
          isNpc: updated.isNpc,
        }),
      })
    } catch (err) {
      console.warn('API save failed, local only:', err)
    }
  }
}

export async function deleteCharacter(id: string): Promise<void> {
  characters.value = characters.value.filter((c) => c.id !== id)
  persist()

  const isAuthed = isAuthenticated()
  const syncedIds = readSyncedIds()

  if (isAuthed && syncedIds.has(id)) {
    try {
      await fetch(`${API_BASE}/characters/${id}`, {
        method: 'DELETE',
        credentials: 'include',
      })
      const newSynced = new Set(syncedIds)
      newSynced.delete(id)
      writeSyncedIds(newSynced)
    } catch (err) {
      console.warn('API delete failed:', err)
    }
  } else {
    const newSynced = new Set(syncedIds)
    newSynced.delete(id)
    writeSyncedIds(newSynced)
  }
}

export async function setLive(id: string, live: SheetLive): Promise<void> {
  const c = getCharacter(id)
  if (!c) return
  await saveCharacter({ ...c, live })
}

export async function setCampaign(id: string, campaignId: string | undefined): Promise<void> {
  const c = getCharacter(id)
  if (!c) return
  await saveCharacter({ ...c, campaignId })
}

export async function buildDraft(draft: BuildRequest): Promise<Sheet> {
  const res = await fetch('/api/v1/build', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(draft),
  })
  if (!res.ok) {
    let msg = `HTTP ${res.status}`
    try {
      const body = (await res.json()) as { error?: { message?: string } }
      if (body?.error?.message) msg = body.error.message
    } catch {
      /* non-JSON error body */
    }
    throw new Error(msg)
  }
  const data = (await res.json()) as { sheet: Sheet }
  return data.sheet
}

export async function buildCharacter(c: Character): Promise<Sheet> {
  try {
    const sheet = await buildDraft(c.draft)
    await saveCharacter({ ...c, sheet })
    return sheet
  } catch (err) {
    if (c.sheet) return c.sheet
    throw err
  }
}

export function _resetCharacters() {
  characters.value = []
  try {
    localStorage.removeItem(KEY)
    localStorage.removeItem(SYNC_KEY)
  } catch {
    /* ignore */
  }
}