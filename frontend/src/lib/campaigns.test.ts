import { beforeEach, describe, expect, it } from 'vitest'
import {
  _resetCampaigns,
  createCampaign,
  deleteCampaign,
  getCampaign,
  listCampaigns,
} from './campaigns.svelte'

beforeEach(() => {
  _resetCampaigns()
  localStorage.clear()
})

describe('campaigns', () => {
  it('creates, lists and gets campaigns', () => {
    const c = createCampaign('Avernus')
    expect(listCampaigns()).toHaveLength(1)
    expect(getCampaign(c.id)?.name).toBe('Avernus')
    expect(getCampaign('missing')).toBeUndefined()
  })

  it('falls back to a default name when empty', () => {
    const c = createCampaign('   ')
    expect(c.name).toBe('Nova campanha')
  })

  it('persists to localStorage', () => {
    const c = createCampaign('Avernus')
    const raw = JSON.parse(localStorage.getItem('onatar.campaigns') ?? '[]') as Array<{
      id: string
    }>
    expect(raw).toHaveLength(1)
    expect(raw[0].id).toBe(c.id)
  })

  it('deletes a campaign', () => {
    const c = createCampaign('Avernus')
    deleteCampaign(c.id)
    expect(listCampaigns()).toHaveLength(0)
  })
})
