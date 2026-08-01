// RF-09 (Fase 2 v1.2): local combat tracker. Sessions persist in localStorage;
// combatants linked to a Character share HP/conditions with the sheet's `live`
// state (Fase 2 changelog #5).
import { box } from './box.svelte'
import { getCharacter, newId, setLive } from './characters.svelte'

const KEY = 'onatar:combat'

export interface Combatant {
  id: string
  characterId?: string
  name: string
  initiative: number
  hpCurrent: number
  hpMax: number
  conditions: string[]
}

export interface CombatSession {
  id: string
  name: string
  round: number
  turnIndex: number
  combatants: Combatant[]
  createdAt: number
  updatedAt: number
}

function readAll(): CombatSession[] {
  if (typeof localStorage === 'undefined') return []
  try {
    return JSON.parse(localStorage.getItem(KEY) ?? '[]') as CombatSession[]
  } catch {
    return []
  }
}

export const sessions = box<CombatSession[]>(readAll())

function persist() {
  localStorage.setItem(KEY, JSON.stringify(sessions.value))
}

function save(s: CombatSession) {
  sessions.value = [...sessions.value.filter((x) => x.id !== s.id), { ...s, updatedAt: Date.now() }]
  persist()
}

export function listSessions(): CombatSession[] {
  return sessions.value
}

export function getSession(id: string): CombatSession | undefined {
  return sessions.value.find((s) => s.id === id)
}

export function newCombatSession(name: string): CombatSession {
  const now = Date.now()
  const s: CombatSession = {
    id: newId(),
    name: name || 'Combate',
    round: 1,
    turnIndex: 0,
    combatants: [],
    createdAt: now,
    updatedAt: now,
  }
  sessions.value = [...sessions.value, s]
  persist()
  return s
}

export function deleteSession(id: string) {
  sessions.value = sessions.value.filter((s) => s.id !== id)
  persist()
}

function mutate(id: string, fn: (s: CombatSession) => CombatSession) {
  const s = getSession(id)
  if (s) save(fn(s))
}

function emptyConditions(c: Combatant | undefined): string[] {
  return c?.conditions ?? []
}

/** Write HP/conditions back to the linked character's live sheet. */
function syncLive(c: Combatant) {
  if (!c.characterId) return
  const ch = getCharacter(c.characterId)
  if (!ch) return
  const live = ch.live ?? {
    hpCurrent: c.hpCurrent,
    slotsUsed: Array.from({ length: 9 }, () => 0),
    conditions: [],
    resources: {},
  }
  setLive(ch.id, { ...live, hpCurrent: c.hpCurrent, conditions: c.conditions })
}

export function addCombatant(sessionId: string, data: Omit<Combatant, 'id'>) {
  mutate(sessionId, (s) => ({
    ...s,
    combatants: [...s.combatants, { ...data, id: newId() }],
  }))
}

/** Link an existing character (or NPC) into the combat, carrying its live HP. */
export function addCharacterToCombat(sessionId: string, characterId: string) {
  const ch = getCharacter(characterId)
  if (!ch) return
  const max = ch.sheet?.hp.max ?? ch.live?.hpCurrent ?? 10
  const current = ch.live?.hpCurrent ?? max
  addCombatant(sessionId, {
    name: ch.name,
    characterId,
    initiative: 0,
    hpCurrent: current,
    hpMax: max,
    conditions: emptyConditions(ch.live),
  })
}

export function removeCombatant(sessionId: string, combatantId: string) {
  mutate(sessionId, (s) => ({
    ...s,
    combatants: s.combatants.filter((c) => c.id !== combatantId),
  }))
}

export function setInitiative(sessionId: string, combatantId: string, value: number) {
  mutate(sessionId, (s) => ({
    ...s,
    combatants: s.combatants.map((c) => (c.id === combatantId ? { ...c, initiative: value } : c)),
  }))
}

export function sortByInitiative(sessionId: string) {
  mutate(sessionId, (s) => ({
    ...s,
    combatants: [...s.combatants].sort((a, b) => b.initiative - a.initiative),
    turnIndex: 0,
  }))
}

export function damage(sessionId: string, combatantId: string, amount: number) {
  mutate(sessionId, (s) => {
    const updated = s.combatants.map((c) =>
      c.id === combatantId ? { ...c, hpCurrent: Math.max(0, c.hpCurrent - amount) } : c
    )
    const target = updated.find((c) => c.id === combatantId)
    if (target) syncLive(target)
    return { ...s, combatants: updated }
  })
}

export function heal(sessionId: string, combatantId: string, amount: number) {
  mutate(sessionId, (s) => {
    const updated = s.combatants.map((c) =>
      c.id === combatantId ? { ...c, hpCurrent: Math.min(c.hpMax, c.hpCurrent + amount) } : c
    )
    const target = updated.find((c) => c.id === combatantId)
    if (target) syncLive(target)
    return { ...s, combatants: updated }
  })
}

export function toggleCondition(sessionId: string, combatantId: string, condition: string) {
  mutate(sessionId, (s) => {
    const updated = s.combatants.map((c) => {
      if (c.id !== combatantId) return c
      const conditions = c.conditions.includes(condition)
        ? c.conditions.filter((x) => x !== condition)
        : [...c.conditions, condition]
      syncLive({ ...c, conditions })
      return { ...c, conditions }
    })
    return { ...s, combatants: updated }
  })
}

export function nextTurn(sessionId: string) {
  mutate(sessionId, (s) => {
    const n = s.combatants.length
    if (n === 0) return s
    const next = s.turnIndex + 1
    return next >= n ? { ...s, turnIndex: 0, round: s.round + 1 } : { ...s, turnIndex: next }
  })
}

export function prevTurn(sessionId: string) {
  mutate(sessionId, (s) => {
    const n = s.combatants.length
    if (n === 0) return s
    const prev = s.turnIndex - 1
    return prev < 0
      ? { ...s, turnIndex: n - 1, round: Math.max(1, s.round - 1) }
      : { ...s, turnIndex: prev }
  })
}

export function _resetCombat() {
  sessions.value = []
  localStorage.removeItem(KEY)
}
