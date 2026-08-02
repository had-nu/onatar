import { beforeEach, describe, expect, it, vi } from 'vitest'
import { characterToJSON, downloadJSON, parseCharacterJSON } from './export'
import { _resetCharacters, createCharacter, starterDraft } from './characters.svelte'

beforeEach(() => {
  _resetCharacters()
  localStorage.clear()
  vi.unstubAllGlobals()
})

describe('JSON export', () => {
  it('round-trips a character through toJSON/parseJSON', async () => {
    const c = await createCharacter(starterDraft())
    const parsed = parseCharacterJSON(characterToJSON(c))
    expect(parsed).not.toBeNull()
    expect(parsed!.id).toBe(c.id)
    expect(parsed!.name).toBe(c.name)
    expect(parsed!.draft.classes[0].id).toBe('fighter')
  })

  it('rejects malformed payloads', () => {
    expect(parseCharacterJSON('not json')).toBeNull()
    expect(parseCharacterJSON('{"foo":1}')).toBeNull()
    expect(parseCharacterJSON('{"character":{"name":"x"}}')).toBeNull()
    expect(
      parseCharacterJSON(
        '{"character":{"id":"a","name":"x","draft":{},"createdAt":1,"updatedAt":1}}'
      )
    ).toBeNull()
  })

  it('downloadJSON triggers an anchor click', () => {
    const click = vi.fn()
    const createObjectURL = vi.fn(() => 'blob:x')
    const revokeObjectURL = vi.fn()
    vi.stubGlobal('URL', { createObjectURL, revokeObjectURL })
    const anchor = { click, remove: vi.fn() }
    const appendChild = vi.fn()
    document.body.appendChild = appendChild
    vi.stubGlobal('document', {
      ...document,
      createElement: (tag: string) => (tag === 'a' ? anchor : document.createElement(tag)),
      body: { appendChild },
    })
    downloadJSON('x.json', '{}')
    expect(click).toHaveBeenCalled()
    expect(revokeObjectURL).toHaveBeenCalled()
  })
})