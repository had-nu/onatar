// Export/import helpers (PRD RF-04). JSON is plain Blob downloads; PDF uses
// jsPDF + html2canvas to snapshot the rendered sheet DOM node.
import type { Character } from './types'

const REQUIRED = ['id', 'name', 'draft', 'createdAt', 'updatedAt']

export function characterToJSON(c: Character): string {
  return JSON.stringify(
    {
      app: 'onatar',
      version: 1,
      character: c,
    },
    null,
    2
  )
}

export function downloadJSON(filename: string, content: string) {
  const blob = new Blob([content], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = filename
  document.body.appendChild(a)
  a.click()
  a.remove()
  URL.revokeObjectURL(url)
}

/** Validate a parsed JSON payload and return the Character or null. */
export function parseCharacterJSON(text: string): Character | null {
  try {
    const parsed = JSON.parse(text) as unknown
    if (!parsed || typeof parsed !== 'object') return null
    const payload = (parsed as { character?: Character }).character
    if (!payload || typeof payload !== 'object') return null
    for (const key of REQUIRED) {
      if (!(key in payload)) return null
    }
    if (!Array.isArray(payload.draft.classes)) return null
    return payload
  } catch {
    return null
  }
}

/** Render a DOM node to a PDF snapshot (best-effort, requires the node in DOM). */
export async function exportCharacterPDF(node: HTMLElement, filename: string): Promise<void> {
  const [{ jsPDF }, html2canvasModule] = await Promise.all([import('jspdf'), import('html2canvas')])
  const html2canvas = (
    html2canvasModule as unknown as {
      default: (el: HTMLElement, opts?: Record<string, unknown>) => Promise<HTMLCanvasElement>
    }
  ).default
  if (!html2canvas) throw new Error('html2canvas indisponível')
  const canvas = await html2canvas(node, {
    scale: 2,
    useCORS: true,
    backgroundColor: '#ffffff',
  })
  const img = canvas.toDataURL('image/png')
  const pdf = new jsPDF({ orientation: 'portrait', unit: 'mm', format: 'a4' })
  const pageWidth = pdf.internal.pageSize.getWidth()
  const imgWidth = pageWidth
  const imgHeight = (canvas.height * imgWidth) / canvas.width
  let heightLeft = imgHeight
  let position = 0
  pdf.addImage(img, 'PNG', 0, position, imgWidth, imgHeight)
  heightLeft -= pdf.internal.pageSize.getHeight()
  while (heightLeft > 0) {
    position -= pdf.internal.pageSize.getHeight()
    pdf.addPage()
    pdf.addImage(img, 'PNG', 0, position, imgWidth, imgHeight)
    heightLeft -= pdf.internal.pageSize.getHeight()
  }
  pdf.save(filename)
}
