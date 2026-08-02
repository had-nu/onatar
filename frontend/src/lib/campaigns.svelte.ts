// Campaigns store (RF-05, Fase 3): campaigns with dual-mode persistence.
// - Guest: localStorage (minimal id + name)
// - Authenticated: MariaDB via API + localStorage cache for offline (RNF-03)
// Migration: on first login, localStorage campaigns are pushed to API.
import { box } from './box.svelte'
import { isAuthenticated, getUser } from './auth.svelte'
import { newId } from './characters.svelte'
import type { User } from './types'

export interface Campaign {
  id: string
  name: string
  description?: string
  ownerId?: string
  createdAt: number
  updatedAt: number
}

interface ApiCampaign {
  id: string
  name: string
  description?: string
  ownerId: string
  createdAt: number
  updatedAt: number
}

const KEY = 'onatar.campaigns'
const SYNC_KEY = 'onatar.campaigns.synced'

const API_BASE = '/api/v1'

function readAll(): Campaign[] {
  try {
    const raw = localStorage.getItem(KEY)
    if (!raw) return []
    const arr = JSON.parse(raw) as unknown
    return Array.isArray(arr) ? (arr as Campaign[]) : []
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

export const campaigns = box<Campaign[]>(readAll())

let currentUser: User | null = null
let syncInProgress = false

function persist() {
  try {
    localStorage.setItem(KEY, JSON.stringify(campaigns.value))
  } catch {
    /* storage unavailable — ignore */
  }
}

function toLocal(apiCamp: ApiCampaign): Campaign {
  return {
    id: apiCamp.id,
    name: apiCamp.name,
    description: apiCamp.description,
    ownerId: apiCamp.ownerId,
    createdAt: apiCamp.createdAt,
    updatedAt: apiCamp.updatedAt,
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
    const apiCamps = await fetchJson<ApiCampaign[]>(`${API_BASE}/campaigns`)
    const localCamps = campaigns.value
    const localById = new Map(localCamps.map((c) => [c.id, c]))

    const merged = apiCamps.map((ac) => {
      const local = localById.get(ac.id)
      if (local && local.updatedAt > ac.updatedAt) {
        return { ...toLocal(ac), _pendingPush: true }
      }
      return toLocal(ac)
    })

    const syncedIds = readSyncedIds()
    for (const local of localCamps) {
      if (!syncedIds.has(local.id) && !apiCamps.some((ac) => ac.id === local.id)) {
        merged.push({ ...local, _localOnly: true })
      }
    }

    campaigns.value = merged
    persist()
  } catch (err) {
    console.warn('Failed to load campaigns from API:', err)
  } finally {
    syncInProgress = false
  }
}

export async function migrateLocalToApi(): Promise<void> {
  const user = getUser()
  if (!user) return

  const localCamps = campaigns.value
  const syncedIds = readSyncedIds()

  for (const camp of localCamps) {
    if (syncedIds.has(camp.id)) continue

    try {
      const res = await fetch(`${API_BASE}/campaigns`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({
          name: camp.name,
          description: camp.description,
        }),
      })

      if (!res.ok) {
        console.warn(`Failed to migrate campaign ${camp.id}:`, await res.text())
        continue
      }

      const created = (await res.json()) as ApiCampaign
      if (created.id !== camp.id) {
        campaigns.value = campaigns.value.map((c) =>
          c.id === camp.id ? { ...c, id: created.id } : c
        )
      }
      syncedIds.add(created.id)
      writeSyncedIds(syncedIds)
    } catch (err) {
      console.warn(`Migration error for ${camp.id}:`, err)
    }
  }

  await loadFromApi()
}

export function setCurrentUser(user: User | null) {
  currentUser = user
}

export function isSyncInProgress(): boolean {
  return syncInProgress
}

export function listCampaigns(): Campaign[] {
  return campaigns.value
}

export function getCampaign(id: string): Campaign | undefined {
  return campaigns.value.find((c) => c.id === id)
}

export async function createCampaign(name: string, description?: string): Promise<Campaign> {
  const trimmed = name.trim()
  const isAuthed = isAuthenticated()

  if (isAuthed) {
    try {
      const res = await fetch(`${API_BASE}/campaigns`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({
          name: trimmed || 'Nova campanha',
          description: description?.trim(),
        }),
      })

      if (!res.ok) throw new Error(await res.text())

      const created = (await res.json()) as ApiCampaign
      const local = toLocal(created)
      campaigns.value = [...campaigns.value, local]
      persist()
      writeSyncedIds(new Set([...readSyncedIds(), created.id]))
      return local
    } catch (err) {
      console.warn('API create campaign failed, falling back to local:', err)
    }
  }

  const campaign: Campaign = {
    id: newId(),
    name: trimmed || 'Nova campanha',
    description: description?.trim(),
    createdAt: Date.now(),
    updatedAt: Date.now(),
  }
  campaigns.value = [...campaigns.value, campaign]
  persist()
  return campaign
}

export async function deleteCampaign(id: string): Promise<void> {
  campaigns.value = campaigns.value.filter((c) => c.id !== id)
  persist()

  const isAuthed = isAuthenticated()
  const syncedIds = readSyncedIds()

  if (isAuthed && syncedIds.has(id)) {
    try {
      await fetch(`${API_BASE}/campaigns/${id}`, {
        method: 'DELETE',
        credentials: 'include',
      })
      const newSynced = new Set(syncedIds)
      newSynced.delete(id)
      writeSyncedIds(newSynced)
    } catch (err) {
      console.warn('API delete campaign failed:', err)
    }
  } else {
    const newSynced = new Set(syncedIds)
    newSynced.delete(id)
    writeSyncedIds(newSynced)
  }
}

export function _resetCampaigns() {
  campaigns.value = []
  try {
    localStorage.removeItem(KEY)
    localStorage.removeItem(SYNC_KEY)
  } catch {
    /* ignore */
  }
}