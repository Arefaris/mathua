'use client'

import Link from 'next/link'
import { upcomingLocked, hrefConceptId, type Shelf, type LockedSuccessor } from '../lib/nextUp'
import type { ConceptProgress } from '../lib/api'
import { concepts } from '../lib/conceptData'
import { useT } from '../lib/i18n'

// The graph helpers read id, label and prerequisites only; the stepper used to carry this
// beside its own state, and it moved here with the one card that uses it.
function catalogEntries() {
  return concepts.map(c => ({ id: c.id, label: c.label, prerequisites: c.prerequisites ?? [], avgTimeSeconds: c.mastery_threshold?.avg_time_seconds }))
}

interface Props {
  conceptId: string
  /** Why the section ended early, when it did. */
  note?: string
  correct: number
  answered: number
  xp: number
  band: string
  /** Where "back" goes, when the learner arrived through a prerequisite chain. */
  back: { href: string } | null
  returnToLabel: string | null
  nextShelf: Shelf | null
  nextLoading: boolean
  shelfProgress: Record<string, ConceptProgress>
  onPracticeAgain: () => void
}

/**
 * The completion card at the end of a concept: the tally, what to do next, and what this
 * concept unlocks. Split out of LearnStepper, which only decides when it appears.
 */
export default function LearnDoneCard({
  conceptId,
  note,
  correct,
  answered,
  xp,
  band,
  back,
  returnToLabel,
  nextShelf,
  nextLoading,
  shelfProgress,
  onPracticeAgain,
}: Props) {
  const t = useT()
  const headId = nextShelf ? hrefConceptId(nextShelf.next.href) : null
  // Only a real concept change may render a navigating <Link>: the
  // stepper keys everything off conceptId, so a same-concept href
  // (however many query params it carries) would scroll to top and
  // strand the learner.
  const canContinue = !nextLoading && nextShelf && headId !== null && headId !== conceptId
  const upcoming: LockedSuccessor[] = upcomingLocked(catalogEntries(), shelfProgress, conceptId)
  return (
    <div className="border border-mathua-green-faint bg-mathua-surface p-6">
      <p className="font-mono text-xs text-green-400">{t('learn.done.complete', { correct, answered, xp, band })}</p>
      {note && <p role="status" className="mt-2 font-mono text-[11px] text-mathua-secondary">{note}</p>}
      <p className="mt-2 font-mono text-[11px] text-mathua-secondary">{t('learn.done.scrollUp')}</p>
      <div className="mt-4 flex flex-wrap gap-2">
        {back && <Link href={back.href} className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint px-5 py-2 font-mono text-xs inline-flex items-center min-h-[40px]">{t('learn.study.backTo', { label: returnToLabel ?? '' })}</Link>}
        {nextLoading && <span className="font-mono text-xs text-mathua-muted inline-flex items-center min-h-[40px]">{t('learn.done.findingNext')}</span>}
        {canContinue && nextShelf && (
          <Link href={nextShelf.next.href} className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint px-5 py-2 font-mono text-xs inline-flex items-center min-h-[40px]">{nextShelf.next.badge}: {nextShelf.next.title} →</Link>
        )}
        {!nextLoading && !nextShelf && (
          <Link href="/profile" className="border border-mathua-border px-5 py-2 font-mono text-xs text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue inline-flex items-center min-h-[40px]">{t('learn.answer.backToProfile')}</Link>
        )}
        <button type="button" onClick={onPracticeAgain} className="border border-mathua-border px-5 py-2 font-mono text-xs text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue inline-flex items-center min-h-[40px]">{t('learn.done.practiceAgain')}</button>
      </div>
      {!nextLoading && nextShelf && nextShelf.alternatives.length > 0 && (
        <details className="mt-3 border border-mathua-border">
          <summary className="font-mono text-[11px] text-mathua-secondary cursor-pointer px-4 py-2.5">
            {t('learn.nextUp.more', { count: nextShelf.alternatives.length })}
          </summary>
          <div className="px-4 pb-4 grid grid-cols-1 sm:grid-cols-2 gap-2">
            {nextShelf.alternatives.map((it, i) => (
              <Link
                key={`${it.kind}:${it.href}:${i}`}
                href={it.href}
                className="border border-mathua-border p-3 hover:border-mathua-blue transition-colors block min-w-0"
              >
                <div className="font-mono text-xs text-mathua-primary truncate">{it.title}</div>
                <div className="mt-1 font-mono text-[11px] text-mathua-blue">{it.cta}</div>
              </Link>
            ))}
          </div>
        </details>
      )}
      {upcoming.length > 0 && (
        <div className="mt-3 border border-mathua-border">
          <p className="font-mono text-[11px] text-mathua-secondary px-4 py-2.5">
            {t('learn.done.comingUp')}
          </p>
          <div className="px-4 pb-4 grid grid-cols-1 sm:grid-cols-2 gap-2">
            {upcoming.map(u => (
              <Link
                key={u.id}
                href={`/learn?concept=${encodeURIComponent(u.id)}`}
                className="border border-mathua-border p-3 hover:border-mathua-blue transition-colors block min-w-0"
              >
                <div className="font-mono text-xs text-mathua-primary truncate">{u.label}</div>
                <div className="mt-1 font-mono text-[11px] text-mathua-muted">{t('learn.done.needs', { names: u.missing.join(', ') })}</div>
              </Link>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}
