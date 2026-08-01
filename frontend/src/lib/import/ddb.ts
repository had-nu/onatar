// D&D Beyond PDF/JSON import (RF-08, Fase 2 v1.1).
//
// Strategy: client-side parsing with pdfjs-dist. Text items are flattened into
// column-aware lines, then a tolerant keyword parser extracts fields into a
// normalized ParsedDDB. Name matching runs against the cached rule content so
// the parser does not hardcode D&D spell/feat/class lists. Extraction is a
// starting point — the Import view forces a manual review before creating a
// draft (Fase 2 changelog #3).
import type { Ability, BuildRequest, ClassInput, Content } from '../types'
import { ABILITIES } from '../types'

export const MAX_IMPORT_BYTES = 5 * 1024 * 1024

export interface ParsedDDB {
  name?: string
  classes: ClassInput[]
  speciesId?: string
  backgroundId?: string
  abilityScores: Partial<Record<Ability, number>>
  hpMax?: number
  ac?: number
  spellIds: string[]
  featIds: string[]
  skillIds: string[]
  /** Names that appeared in context but could not be matched to content. */
  unmapped: string[]
}

const ABILITY_LABELS: Record<string, Ability> = {
  STRENGTH: 'STR',
  DEXTERITY: 'DEX',
  CONSTITUTION: 'CON',
  INTELLIGENCE: 'INT',
  WISDOM: 'WIS',
  CHARISMA: 'CHA',
  STR: 'STR',
  DEX: 'DEX',
  CON: 'CON',
  INT: 'INT',
  WIS: 'WIS',
  CHA: 'CHA',
}

const SKILL_NAMES = [
  'Acrobatics',
  'Animal Handling',
  'Arcana',
  'Athletics',
  'Deception',
  'History',
  'Insight',
  'Intimidation',
  'Investigation',
  'Medicine',
  'Nature',
  'Perception',
  'Performance',
  'Persuasion',
  'Religion',
  'Sleight of Hand',
  'Stealth',
  'Survival',
]

function escapeRegExp(s: string): string {
  return s.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
}

/** Match whole multi-word names (case-insensitive) and return them in order. */
function matchNames(names: string[], text: string): string[] {
  const found: string[] = []
  for (const n of names) {
    const re = new RegExp(`\\b${escapeRegExp(n)}\\b`, 'i')
    if (re.test(text)) found.push(n)
  }
  return found
}

/** Match a content entity (class/species/background/subclass) by label context. */
function matchByLabel(text: string, labels: string[], contentNames: string[]): string | undefined {
  for (const label of labels) {
    const re = new RegExp(`${label}\\s*[:]?\\s*([A-Za-z][A-Za-z0-9 '()\\-]{1,40})`, 'i')
    const m = re.exec(text)
    if (!m) continue
    const phrase = m[1].trim()
    const hit = contentNames
      .filter((n) => phrase.toLowerCase().startsWith(n.toLowerCase()))
      .sort((a, b) => b.length - a.length)[0]
    if (hit) return hit
  }
  return undefined
}

function parseAbilityScores(text: string): Partial<Record<Ability, number>> {
  const scores: Partial<Record<Ability, number>> = {}

  // Primary layout: a label row (STRENGTH DEXTERITY ...) followed by a score
  // row (17 10 15 ...). Pair labels with scores positionally.
  const lines = text.split('\n').map((l) => l.trim())
  for (let i = 0; i < lines.length - 1; i++) {
    const labels = lines[i]
      .split(/\s+/)
      .map((t) => t.replace(/[^A-Za-z]/g, '').toUpperCase())
      .filter((t) => t in ABILITY_LABELS)
    if (labels.length < 3) continue
    const nums = (lines[i + 1].match(/\d{1,2}/g) ?? []).map(Number)
    if (nums.length < labels.length) continue
    labels.forEach((l, idx) => {
      const ability = ABILITY_LABELS[l]
      const v = nums[idx]
      if (ability && v >= 1 && v <= 30 && !(ability in scores)) scores[ability] = v
    })
    if (Object.keys(scores).length === 6) return scores
  }

  // Fallback: inline abbreviations anywhere in the text (STR 17, DEX 10, ...).
  for (const abbr of ['STR', 'DEX', 'CON', 'INT', 'WIS', 'CHA']) {
    const re = new RegExp(`\\b${abbr}\\s*[:]?\\s*(\\d{1,2})\\b`, 'i')
    const m = re.exec(text)
    if (m) {
      const v = Number(m[1])
      if (v >= 1 && v <= 30) scores[ABILITY_LABELS[abbr]] = v
    }
  }
  return scores
}

export function parseDDBText(text: string, content: Content): ParsedDDB {
  const norm = text.replace(/\r/g, '').replace(/[ \t]+/g, ' ')
  const scores = parseAbilityScores(norm)

  const classNames = content.classes.map((c) => c.name)
  const subclassNames = content.subclasses.map((s) => s.name)
  const speciesNames = content.species.map((s) => s.name)
  const backgroundNames = content.backgrounds.map((b) => b.name)

  // Class: "Fighter Level 3" / "Level 3 Fighter" / "Fighter" near "Level N".
  const classes: ClassInput[] = []
  const levelRe = /Level\s+(\d{1,2})/gi
  let lm: RegExpExecArray | null
  while ((lm = levelRe.exec(norm)) !== null) {
    const level = Number(lm[1])
    const windowText = norm.slice(Math.max(0, lm.index - 60), lm.index + 60)
    const cls = classNames
      .filter((n) => new RegExp(`\\b${escapeRegExp(n)}\\b`, 'i').test(windowText))
      .sort((a, b) => b.length - a.length)[0]
    if (cls && !classes.some((c) => c.id === content.classes.find((x) => x.name === cls)?.id)) {
      const id = content.classes.find((x) => x.name === cls)?.id
      if (id) classes.push({ id, level })
    }
  }
  if (classes.length === 0) {
    const first = classNames
      .filter((n) => new RegExp(`\\b${escapeRegExp(n)}\\b`, 'i').test(norm))
      .sort((a, b) => b.length - a.length)[0]
    const id = first ? content.classes.find((x) => x.name === first)?.id : undefined
    if (id) classes.push({ id, level: 1 })
  }

  // Subclass: attach to the first class if any.
  const subclass = matchByLabel(
    norm,
    ['Subclass', 'Archetype', 'College', 'Domain', 'School', 'Origin'],
    subclassNames
  )
  if (subclass && classes.length > 0) {
    const sid = content.subclasses.find((s) => s.name === subclass)?.id
    if (sid) classes[0] = { ...classes[0], subclassId: sid }
  }

  const species = matchByLabel(norm, ['Species', 'Race', 'Species:'], speciesNames)
  const background = matchByLabel(norm, ['Background', 'Background:'], backgroundNames)

  const spellIds = matchNames(
    content.spells.map((s) => s.name),
    norm
  )
    .map((n) => content.spells.find((s) => s.name === n)?.id)
    .filter((id): id is string => Boolean(id))

  const featIds = matchNames(
    content.feats.map((f) => f.name),
    norm
  )
    .map((n) => content.feats.find((f) => f.name === n)?.id)
    .filter((id): id is string => Boolean(id))

  const skillIds = matchNames(SKILL_NAMES, norm).map((n) => n.toLowerCase().replace(/\s+/g, '-'))

  // Name: first non-empty line that is not a label/known entity name.
  const name = norm
    .split('\n')
    .map((l) => l.trim())
    .filter(
      (l) =>
        l.length > 1 &&
        !/^(str|dex|con|int|wis|cha)/i.test(l) &&
        !/level|species|background|class/i.test(l)
    )[0]

  const hpMatch = /Hit Points\s*[:]?\s*(\d{1,4})/i.exec(norm)
  const acMatch = /Armor Class\s*[:]?\s*(\d{1,3})/i.exec(norm)

  return {
    name: name && name.length <= 60 ? name : undefined,
    classes,
    speciesId: species ? content.species.find((s) => s.name === species)?.id : undefined,
    backgroundId: background
      ? content.backgrounds.find((b) => b.name === background)?.id
      : undefined,
    abilityScores: scores,
    hpMax: hpMatch ? Number(hpMatch[1]) : undefined,
    ac: acMatch ? Number(acMatch[1]) : undefined,
    spellIds,
    featIds,
    skillIds,
    unmapped: [],
  }
}

/** Convert a parsed sheet into a BuildRequest draft (missing fields get sane defaults). */
export function toBuildRequest(parsed: ParsedDDB): BuildRequest {
  const abilityScores = { ...(parsed.abilityScores as Record<Ability, number>) }
  for (const a of ABILITIES) if (!abilityScores[a]) abilityScores[a] = 10

  const classes: ClassInput[] =
    parsed.classes.length > 0
      ? parsed.classes
      : [{ id: 'fighter', level: parsed.classes.length === 0 ? (parsed.level ?? 1) : 1 }]

  return {
    name: parsed.name ?? 'Importado',
    classes,
    speciesId: parsed.speciesId ?? 'human',
    backgroundId: parsed.backgroundId ?? 'sage',
    abilityScores,
    abilityMethod: 'manual',
    skills: parsed.skillIds,
    spells: parsed.spellIds,
    feats: parsed.featIds,
    isNpc: false,
  }
}

// ---- D&D Beyond JSON import (bonus adapter) --------------------------------

/** Parse a D&D Beyond JSON export (extension-exporters style) into ParsedDDB. */
export async function parseDDBJSON(file: File, content: Content): Promise<ParsedDDB> {
  const text = await file.text()
  let data: unknown
  try {
    data = JSON.parse(text)
  } catch {
    throw new Error('JSON inválido.')
  }
  const obj = (data && typeof data === 'object' ? data : {}) as Record<string, unknown>
  const root = (obj.character ?? obj) as Record<string, unknown>

  const classes: ClassInput[] = []
  const rawClasses = root.classes
  if (Array.isArray(rawClasses)) {
    for (const rc of rawClasses) {
      const cls = rc as Record<string, unknown>
      const name = typeof cls.name === 'string' ? cls.name : ''
      const id = content.classes.find((c) => c.name.toLowerCase() === name.toLowerCase())?.id
      const level = Number(cls.level ?? 1)
      if (id) classes.push({ id, level })
    }
  }

  const speciesName = getString(root, ['species', 'race'], 'name')
  const speciesId = speciesName
    ? content.species.find((s) => s.name.toLowerCase() === speciesName.toLowerCase())?.id
    : undefined
  const backgroundName = getString(root, ['background'], 'name')
  const backgroundId = backgroundName
    ? content.backgrounds.find((b) => b.name.toLowerCase() === backgroundName.toLowerCase())?.id
    : undefined

  const abilityScores: Partial<Record<Ability, number>> = {}
  const stats = root.stats ?? root.scores ?? root.abilities
  if (stats && typeof stats === 'object') {
    for (const [k, v] of Object.entries(stats as Record<string, unknown>)) {
      const ability = abilityFromKey(k)
      if (!ability) continue
      const n =
        typeof v === 'object' && v ? Number((v as Record<string, unknown>).value) : Number(v)
      if (Number.isFinite(n) && n >= 1 && n <= 30) abilityScores[ability] = n
    }
  }

  const spellIds = collectNames(root.spells).flatMap((n) => {
    const id = content.spells.find((s) => s.name.toLowerCase() === n.toLowerCase())?.id
    return id ? [id] : []
  })
  const featIds = collectNames(root.feats).flatMap((n) => {
    const id = content.feats.find((f) => f.name.toLowerCase() === n.toLowerCase())?.id
    return id ? [id] : []
  })

  return {
    name: typeof root.name === 'string' ? root.name : undefined,
    classes,
    speciesId,
    backgroundId,
    abilityScores,
    spellIds,
    featIds,
    skillIds: [],
    unmapped: [],
  }
}

/** Map a JSON stat key (str/dex/constitution/strength/…) to an Ability. */
function abilityFromKey(key: string): Ability | undefined {
  const norm = key.toLowerCase().replace(/[^a-z]/g, '')
  if (norm === 'strength' || norm === 'str') return 'STR'
  if (norm === 'dexterity' || norm === 'dex') return 'DEX'
  if (norm === 'constitution' || norm === 'con') return 'CON'
  if (norm === 'intelligence' || norm === 'int') return 'INT'
  if (norm === 'wisdom' || norm === 'wis') return 'WIS'
  if (norm === 'charisma' || norm === 'cha') return 'CHA'
  return undefined
}

function getString(root: Record<string, unknown>, keys: string[], sub: string): string | undefined {
  for (const k of keys) {
    const v = root[k]
    if (v && typeof v === 'object') {
      const s = (v as Record<string, unknown>)[sub]
      if (typeof s === 'string') return s
    }
  }
  return undefined
}

function collectNames(list: unknown): string[] {
  if (!Array.isArray(list)) return []
  return list
    .map((item) => {
      if (typeof item === 'string') return item
      if (item && typeof item === 'object') {
        const s = (item as Record<string, unknown>).name
        return typeof s === 'string' ? s : ''
      }
      return ''
    })
    .filter((s) => s.length > 0)
}

// ---- pdf.js extraction -----------------------------------------------------

type TextItem = { str: string; transform: number[] }

let workerConfigured = false

/** Group pdf.js text items into column-aware lines (x is transform[4]). */
export function itemsToLines(items: TextItem[]): string {
  const rows = items.map((it) => ({ x: it.transform[4], y: it.transform[5], str: it.str }))
  rows.sort((a, b) => b.y - a.y || a.x - b.x)
  const lines: { y: number; cols: { x: number; str: string }[] }[] = []
  for (const r of rows) {
    const last = lines[lines.length - 1]
    if (last && Math.abs(last.y - r.y) < 2) last.cols.push({ x: r.x, str: r.str })
    else lines.push({ y: r.y, cols: [{ x: r.x, str: r.str }] })
  }
  return lines
    .map((l) => {
      l.cols.sort((a, b) => a.x - b.x)
      let out = ''
      let prevX: number | null = null
      for (const c of l.cols) {
        if (prevX !== null && c.x - prevX > 12) out += '\t'
        out += c.str
        prevX = c.x
      }
      return out
    })
    .join('\n')
}

/** Extract plain text from a PDF file using pdfjs-dist (lazy-loaded worker). */
export async function extractTextFromPDF(file: File): Promise<string> {
  if (file.size > MAX_IMPORT_BYTES) {
    throw new Error(`Ficheiro demasiado grande (máx. ${MAX_IMPORT_BYTES / 1024 / 1024} MiB)`)
  }
  const pdfjsLib = await import('pdfjs-dist')
  if (!workerConfigured) {
    const workerSrc = (await import('pdfjs-dist/build/pdf.worker.min.mjs?url')).default
    pdfjsLib.GlobalWorkerOptions.workerSrc = workerSrc
    workerConfigured = true
  }
  const data = new Uint8Array(await file.arrayBuffer())
  const doc = await pdfjsLib.getDocument({ data }).promise
  try {
    const chunks: string[] = []
    for (let p = 1; p <= doc.numPages; p++) {
      const page = await doc.getPage(p)
      const textContent = await page.getTextContent()
      chunks.push(itemsToLines(textContent.items as unknown as TextItem[]))
    }
    return chunks.join('\n')
  } finally {
    await doc.destroy()
  }
}
