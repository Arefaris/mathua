'use client'

import { useT, type MessageKey } from '../lib/i18n'

// Each item names two dictionary keys rather than carrying text, so the same constants render in
// either language and the callers (the quiz, the diagnostic, reviews) pass them on unchanged.
export interface BriefingItem {
  titleKey: MessageKey
  detailKey: MessageKey
}

// Pre-test expectations, MA-grounded: diagnostics run 25–45 adaptive
// questions (~30–45 min per engine bounds), quizzes are timed + closed-book.
// Stated once up front so the test itself stays distraction-free.
export const DIAGNOSTIC_BRIEFING: BriefingItem[] = [
  { titleKey: 'learn.briefing.diagnostic.pace.title', detailKey: 'learn.briefing.diagnostic.pace.detail' },
  { titleKey: 'learn.briefing.diagnostic.honest.title', detailKey: 'learn.briefing.diagnostic.honest.detail' },
  { titleKey: 'learn.briefing.diagnostic.speed.title', detailKey: 'learn.briefing.diagnostic.speed.detail' },
  { titleKey: 'learn.briefing.diagnostic.pause.title', detailKey: 'learn.briefing.diagnostic.pause.detail' },
]

export const QUIZ_BRIEFING: BriefingItem[] = [
  { titleKey: 'learn.briefing.quiz.format.title', detailKey: 'learn.briefing.quiz.format.detail' },
  { titleKey: 'learn.briefing.quiz.scope.title', detailKey: 'learn.briefing.quiz.scope.detail' },
  { titleKey: 'learn.briefing.quiz.miss.title', detailKey: 'learn.briefing.quiz.miss.detail' },
  { titleKey: 'learn.briefing.quiz.xp.title', detailKey: 'learn.briefing.quiz.xp.detail' },
]

export const REVIEW_BRIEFING: BriefingItem[] = [
  { titleKey: 'learn.briefing.review.due.title', detailKey: 'learn.briefing.review.due.detail' },
  { titleKey: 'learn.briefing.review.honest.title', detailKey: 'learn.briefing.review.honest.detail' },
  { titleKey: 'learn.briefing.review.speed.title', detailKey: 'learn.briefing.review.speed.detail' },
  { titleKey: 'learn.briefing.review.xp.title', detailKey: 'learn.briefing.review.xp.detail' },
]

export default function BriefingCard({ eyebrow, items }: { eyebrow: string; items: BriefingItem[] }) {
  const t = useT()
  return (
    <div className="border border-mathua-border bg-mathua-surface p-4 text-left">
      <p className="font-mono text-[10px] uppercase tracking-wider text-mathua-muted mb-3">
        {eyebrow}
      </p>
      <ul className="space-y-3">
        {items.map(item => (
          <li key={item.titleKey} className="min-w-0">
            <p className="font-mono text-xs text-mathua-primary">{t(item.titleKey)}</p>
            <p className="font-mono text-[11px] text-mathua-muted mt-0.5 break-words [overflow-wrap:anywhere]">
              {t(item.detailKey)}
            </p>
          </li>
        ))}
      </ul>
    </div>
  )
}
