import { beforeEach, describe, expect, it } from 'vitest'
import {
  addCombatant,
  addCharacterToCombat,
  damage,
  deleteSession,
  getSession,
  heal,
  newCombatSession,
  nextTurn,
  prevTurn,
  removeCombatant,
  sessions,
  setInitiative,
  sortByInitiative,
  toggleCondition,
  _resetCombat,
} from './combat.svelte'
import { _resetCharacters, createCharacter, getCharacter } from './characters.svelte'
import type { BuildRequest } from './types'

const draft: BuildRequest = {
  name: 'E2E Hero',
  classes: [{ id: 'fighter', level: 1 }],
  abilityScores: { STR: 15, DEX: 13, CON: 14, INT: 10, WIS: 12, CHA: 8 },
}

beforeEach(() => {
  _resetCombat()
  _resetCharacters()
})

function seedSession() {
  const s = newCombatSession('Teste')
  addCombatant(s.id, { name: 'A', initiative: 0, hpCurrent: 10, hpMax: 10, conditions: [] })
  addCombatant(s.id, { name: 'B', initiative: 0, hpCurrent: 20, hpMax: 20, conditions: [] })
  return getSession(s.id)!
}

describe('combat store', () => {
  it('creates and lists sessions', () => {
    const s = newCombatSession('C1')
    expect(getSession(s.id)?.name).toBe('C1')
    expect(sessions.value).toHaveLength(1)
  })

  it('sorts by initiative descending and resets turn', () => {
    const s = seedSession()
    setInitiative(s.id, s.combatants[0].id, 5)
    setInitiative(s.id, s.combatants[1].id, 18)
    sortByInitiative(s.id)
    const got = getSession(s.id)!
    expect(got.combatants[0].name).toBe('B')
    expect(got.combatants[1].name).toBe('A')
    expect(got.turnIndex).toBe(0)
  })

  it('advances turns and wraps the round', () => {
    const s = seedSession()
    nextTurn(s.id)
    expect(getSession(s.id)!.turnIndex).toBe(1)
    nextTurn(s.id)
    const got = getSession(s.id)!
    expect(got.turnIndex).toBe(0)
    expect(got.round).toBe(2)
  })

  it('moves back a turn and clamps round to 1', () => {
    const s = seedSession()
    prevTurn(s.id)
    const got = getSession(s.id)!
    expect(got.turnIndex).toBe(1)
    expect(got.round).toBe(1)
  })

  it('applies damage clamped at 0 and heal clamped at max', () => {
    const s = seedSession()
    const id = s.combatants[0].id
    damage(s.id, id, 12)
    expect(getSession(s.id)!.combatants[0].hpCurrent).toBe(0)
    heal(s.id, id, 100)
    expect(getSession(s.id)!.combatants[0].hpCurrent).toBe(10)
  })

  it('toggles conditions', () => {
    const s = seedSession()
    const id = s.combatants[0].id
    toggleCondition(s.id, id, 'cego')
    expect(getSession(s.id)!.combatants[0].conditions).toEqual(['cego'])
    toggleCondition(s.id, id, 'cego')
    expect(getSession(s.id)!.combatants[0].conditions).toEqual([])
  })

  it('links a character and syncs HP back to its live sheet', () => {
    const c = createCharacter(draft)
    const s = seedSession()
    addCharacterToCombat(s.id, c.id)
    const linked = getSession(s.id)!.combatants.find((x) => x.characterId === c.id)!
    expect(linked.name).toBe('E2E Hero')
    damage(s.id, linked.id, 4)
    expect(getSession(s.id)!.combatants.find((x) => x.id === linked.id)!.hpCurrent).toBe(6)
    expect(getCharacter(c.id)?.live?.hpCurrent).toBe(6)
  })

  it('removes combatants and sessions', () => {
    const s = seedSession()
    removeCombatant(s.id, s.combatants[0].id)
    expect(getSession(s.id)!.combatants).toHaveLength(1)
    deleteSession(s.id)
    expect(getSession(s.id)).toBeUndefined()
  })
})
