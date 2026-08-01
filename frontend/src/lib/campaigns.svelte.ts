// Campaigns store (RF-05, minimal): campaigns are just `id` + `name`, persisted
// in localStorage. Characters link via `campaignId` (see types.Character).
import { box } from './box.svelte'
import { newId } from './characters.svelte'

export interface Campaign {
  id: string
  name: string
  createdAt: number
}

const KEY = 'onatar.campaigns'

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

export const campaigns = box<Campaign[]>(readAll())

function persist() {
  try {
    localStorage.setItem(KEY, JSON.stringify(campaigns.value))
  } catch {
    /* storage unavailable — ignore */
  }
}

export function listCampaigns(): Campaign[] {
  return campaigns.value
}

export function getCampaign(id: string): Campaign | undefined {
  return campaigns.value.find((c) => c.id === id)
}

export function createCampaign(name: string): Campaign {
  const trimmed = name.trim()
  const campaign: Campaign = {
    id: newId(),
    name: trimmed || 'Nova campanha',
    createdAt: Date.now(),
  }
  campaigns.value = [...campaigns.value, campaign]
  persist()
  return campaign
}

export function deleteCampaign(id: string) {
  campaigns.value = campaigns.value.filter((c) => c.id !== id)
  persist()
}

/** Test helper: reset the store between tests. */
export function _resetCampaigns() {
  campaigns.value = []
  try {
    localStorage.removeItem(KEY)
  } catch {
    /* ignore */
  }
}
