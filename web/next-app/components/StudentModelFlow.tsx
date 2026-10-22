'use client'

import FlowDiagram, { type FlowNodeDef, type FlowEdgeDef } from './FlowDiagram'
import { useLocale, useT } from '../lib/i18n'

export default function StudentModelFlow() {
  const t = useT()
  const locale = useLocale()
  // The five states are drawn in capitals, as they are named in the code.
  const state = (s: 'unseen' | 'learning' | 'practicing' | 'mastered' | 'decaying') =>
    t(`common.state.${s}`).toLocaleUpperCase(locale)
  const nodes: FlowNodeDef[] = [
    { id: 'unseen', label: state('unseen') },
    { id: 'learning', label: state('learning') },
    { id: 'practicing', label: state('practicing') },
    { id: 'mastered', label: state('mastered'), variant: 'terminal' },
    { id: 'decaying', label: state('decaying'), variant: 'error' },
  ]

  // Labels follow internal/mastery/machine.go. Advancement is decided by evidence, so the
  // rungs are named for what the evidence shows, not for a streak. The previous labels were
  // also off by one rung: "first correct" sat on LEARNING -> PRACTICING and "streak reached"
  // on PRACTICING -> MASTERED, which is not the ladder the state diagram draws.
  const edges: FlowEdgeDef[] = [
    { id: 'e-unseen-learning', source: 'unseen', target: 'learning', label: t('explore.howItWorks.flow.firstCorrect'), color: 'green', animated: true },
    { id: 'e-learning-practicing', source: 'learning', target: 'practicing', label: t('explore.howItWorks.flow.consistent'), color: 'green', animated: true },
    { id: 'e-practicing-mastered', source: 'practicing', target: 'mastered', label: t('explore.howItWorks.flow.evidence'), color: 'green', animated: true },
    { id: 'e-mastered-decaying', source: 'mastered', target: 'decaying', label: t('explore.howItWorks.flow.idle'), color: 'blue', animated: true },
    { id: 'e-decaying-practicing', source: 'decaying', target: 'practicing', label: t('explore.howItWorks.flow.reviewCorrect'), dashed: true, color: 'green' },
  ]

  return <FlowDiagram nodes={nodes} edges={edges} direction="TB" height={420} rankSep={150} nodeSep={80} allowZoom />
}
