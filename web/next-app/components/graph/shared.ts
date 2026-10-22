import type { GraphStatus } from '../../lib/graphStatus'
import type { TFunction } from '../../lib/i18n'

export interface GraphConcept {
  id: string
  label: string
  domain: string
  prerequisites: string[]
}

export const STATUS_COLORS: Record<GraphStatus, string> = {
  mastered: '#10b981',
  // Darker amber than `learning`'s, and deliberately so. Both states are "attention"
  // and the brief is that decay is signalled by colour alone — no second bar, no
  // pattern, no badge, no shorter fill — so if the two ambers matched, `learning` and
  // `decaying` would be the same picture. A full bar in a deeper amber reads as
  // "learned, come back to this", which is what it is.
  decaying: '#b45309',
  practicing: '#60a5fa',
  learning: '#f59e0b',
  unseen: '#71717a',
  locked: '#52525b',
}

export function domainColor(domain: string): string {
  let hash = 0
  for (let i = 0; i < domain.length; i++) hash = (hash * 31 + domain.charCodeAt(i)) | 0
  const hue = Math.abs(hash) % 360
  return `hsl(${hue} 55% 55%)`
}

/** A node's state, in the learner's language. The graph's own wording: `decaying` reads "Due for review". */
export function statusLabel(t: TFunction, status: GraphStatus): string {
  return t(`explore.graph.flow.status.${status}`)
}
