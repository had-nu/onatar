import { fireEvent, render, screen } from '@testing-library/svelte'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import Campaigns from './Campaigns.svelte'
import { _resetCampaigns, listCampaigns } from '../campaigns.svelte'
import { _resetCharacters } from '../characters.svelte'

beforeEach(() => {
  _resetCampaigns()
  _resetCharacters()
  localStorage.clear()
  window.confirm = vi.fn(() => true)
})

describe('Campaigns view', () => {
  it('shows the empty state', () => {
    render(Campaigns)
    expect(screen.getByText(/ainda não há campanhas/i)).toBeTruthy()
  })

  it('creates a campaign from the input', async () => {
    render(Campaigns)
    const input = screen.getByPlaceholderText('Nome da campanha') as HTMLInputElement
    await fireEvent.input(input, { target: { value: 'Avernus' } })
    await fireEvent.click(screen.getByRole('button', { name: 'Criar campanha' }))
    expect(screen.getByText('Avernus')).toBeTruthy()
    expect(listCampaigns()).toHaveLength(1)
  })

  it('deletes a campaign after confirmation', async () => {
    render(Campaigns)
    const input = screen.getByPlaceholderText('Nome da campanha') as HTMLInputElement
    await fireEvent.input(input, { target: { value: 'Avernus' } })
    await fireEvent.click(screen.getByRole('button', { name: 'Criar campanha' }))
    await fireEvent.click(screen.getByRole('button', { name: 'Apagar Avernus' }))
    expect(window.confirm).toHaveBeenCalled()
    expect(screen.queryByText('Avernus')).toBeNull()
  })
})
