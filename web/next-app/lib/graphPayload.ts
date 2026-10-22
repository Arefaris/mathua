// Build-time graph payload and its runtime schema boundary.
// Kept out of page-bundle code: only the idle-loaded graph chunk and the
// build script import it.
import { z } from 'zod'
import type { ConceptRecord } from './conceptData'
import type { Locale } from './i18n/locale'

export const GRAPH_PAYLOAD_VERSION = 2

// [id, label, domain, importance] — index in this array is the edge index.
const GraphPayloadNodeSchema = z.tuple([z.string(), z.string(), z.string(), z.number()])
// [sourceIndex, targetIndex] — prerequisite -> dependent.
const GraphPayloadEdgeSchema = z.tuple([z.number().int(), z.number().int()])

export const GraphPayloadSchema = z.object({
  v: z.number(),
  nodes: z.array(GraphPayloadNodeSchema),
  edges: z.array(GraphPayloadEdgeSchema),
  domains: z.record(z.string(), z.number()).default({}),
})

export type GraphPayloadNode = z.infer<typeof GraphPayloadNodeSchema>
export type GraphPayloadEdge = z.infer<typeof GraphPayloadEdgeSchema>
export type GraphPayload = z.infer<typeof GraphPayloadSchema>

// The value the landing page shows for the mastery-check gate, mirroring Go's xp.QuizGateXP.
// It is the one number build-graph.mjs cannot derive from the corpus, so it is the one that can
// drift — and it did, for as long as the page said 150. internal/xp's
// TestLandingQuizGateMatchesEngine reads the generated artifact and fails when the two disagree.
export const QUIZ_GATE_XP = 50

/**
 * Which artifact holds the labels for a language.
 *
 * English is `concepts.graph.json`; every other language is `concepts.graph.<locale>.json`,
 * written by build-graph.mjs only when a label overlay exists. A consumer asks for the file
 * for the learner's language and falls back to the English one when that request fails, so a
 * language with no translated labels yet degrades to English instead of an empty graph.
 */
export function graphPayloadFile(locale: Locale): string {
  // Type-only import above: build-graph.mjs runs this file under plain Node, which cannot
  // resolve the extensionless runtime imports of the app. English is the default locale.
  return locale === 'en' ? 'concepts.graph.json' : `concepts.graph.${locale}.json`
}

/**
 * The corpus with translated labels swapped in. An id the overlay does not mention keeps its
 * English label, so a partial translation is usable while it is being finished.
 */
export function applyLabelOverlay(
  concepts: ConceptRecord[],
  overlay: Record<string, string>
): ConceptRecord[] {
  return concepts.map(c => {
    const label = overlay[c.id]
    return typeof label === 'string' && label.trim() !== '' ? { ...c, label } : c
  })
}

export interface GraphMeta {
  quizGateXP: number
  conceptCount: number
  connectionCount: number
  domainCount: number
  domainCounts: Record<string, number>
  layoutVersion: number
}

export function buildGraphPayload(
  concepts: ConceptRecord[],
  importance: Record<string, number>
): GraphPayload {
  const indexById = new Map<string, number>()
  concepts.forEach((c, i) => indexById.set(c.id, i))

  const nodes: GraphPayloadNode[] = concepts.map(c => [
    c.id,
    c.label,
    c.domain,
    importance[c.id] ?? 0,
  ])

  const edges: GraphPayloadEdge[] = []
  concepts.forEach((c, target) => {
    for (const prereq of c.prerequisites ?? []) {
      const source = indexById.get(prereq)
      if (source !== undefined && source !== target) edges.push([source, target])
    }
  })

  const domains: Record<string, number> = {}
  for (const c of concepts) domains[c.domain] = (domains[c.domain] ?? 0) + 1

  return { v: GRAPH_PAYLOAD_VERSION, nodes, edges, domains }
}

export function buildGraphMeta(concepts: ConceptRecord[]): GraphMeta {
  const domainCounts: Record<string, number> = {}
  let connectionCount = 0
  for (const c of concepts) {
    domainCounts[c.domain] = (domainCounts[c.domain] ?? 0) + 1
    connectionCount += c.prerequisites?.length ?? 0
  }
  return {
    quizGateXP: QUIZ_GATE_XP,
    conceptCount: concepts.length,
    connectionCount,
    domainCount: Object.keys(domainCounts).length,
    domainCounts,
    layoutVersion: GRAPH_PAYLOAD_VERSION,
  }
}

/** Reject payloads whose edge indices fall outside the node array. */
export function validateGraphPayload(payload: GraphPayload): GraphPayload | null {
  const count = payload.nodes.length
  for (const [source, target] of payload.edges) {
    if (source < 0 || target < 0 || source >= count || target >= count) return null
  }
  return payload
}
