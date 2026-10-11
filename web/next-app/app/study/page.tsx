'use client'

import { Suspense, useEffect, useMemo, useState } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import Link from 'next/link'
import Header from '../../components/Header'
import BottomTabs from '../../components/BottomTabs'
import Footer from '../../components/Footer'
import PrerequisitePanel from '../../components/PrerequisitePanel'
import ConceptReference from '../../components/ConceptReference'
import { getLessonReadiness, type ReadinessRes } from '../../lib/api'
import { concepts } from '../../lib/conceptData'
import { backTarget, learnConceptPath, parseReturnStack, pushReturn, studyConceptPath, withReturnStack } from '../../lib/learnReturn'
import { hrefConceptId } from '../../lib/nextUp'

// The reference and exploration surface.
//
// ADR-047 closed `/study` as a learner workflow because it was a *catalogue* — a browsable
// library of whole lessons with its own search and domain index, linked from every header, tab
// bar and sidebar. That let a learner read their way to a concept instead of being taught it.
//
// ADR-054 restores it, scoped: this page renders the reference for **one** concept named in the
// URL. It is not a catalogue — no domain index, no lesson list, no search — and it awards
// nothing. Discovery stays structured in `/domains`; practice stays gated in `/learn`. The only
// navigation here is the concept's own DAG neighbourhood (prerequisites and dependents) plus a
// separate *Practise this concept* control that enters the normal learning loop.
//
// Without `?concept=` there is nothing to show, so it forwards to `/learn` rather than presenting
// an index. Static export cannot use next.config redirects, so this is a client forward, the
// same pattern as the deprecated stubs.
function ConceptNotFound() {
  const router = useRouter()
  useEffect(() => { router.replace('/learn') }, [router])
  return null
}

function stateLabel(state: string | undefined): string {
  switch (state) {
    case 'mastered':
      return 'Mastered'
    case 'due_for_review':
      return 'Mastered · due a check'
    case 'in_progress':
      return 'In progress'
    case 'unlocked':
      return 'Not started'
    case 'upcoming':
      return 'Not started'
    case 'locked':
      return 'Locked'
    default:
      return ''
  }
}

function StudyContent() {
  const searchParams = useSearchParams()
  const conceptId = searchParams.get('concept') || ''
  const returnStack = useMemo(() => parseReturnStack(searchParams), [searchParams])

  const [readiness, setReadiness] = useState<ReadinessRes | null>(null)
  const [failed, setFailed] = useState(false)

  useEffect(() => {
    if (!conceptId) return
    let cancelled = false
    setReadiness(null)
    setFailed(false)
    getLessonReadiness(conceptId)
      .then(res => { if (!cancelled) setReadiness(res) })
      .catch(() => { if (!cancelled) setFailed(true) })
    return () => { cancelled = true }
  }, [conceptId])

  if (!conceptId) return <ConceptNotFound />

  const label = concepts.find(c => c.id === conceptId)?.label || conceptId
  const prerequisites = readiness?.prerequisites ?? []
  // Dependents are presentation over the corpus the app already bundles, the same way the Learn
  // done card derives "what this unlocks". It is a display of the DAG, never a second rule.
  const dependents = concepts.filter(c => (c.prerequisites ?? []).includes(conceptId))

  const back = backTarget(returnStack)
  const backLabel = back ? (concepts.find(c => c.id === hrefConceptId(back.href))?.label ?? '') : ''
  const mastered = readiness?.state === 'mastered' || readiness?.state === 'due_for_review'
  const practiseHref = withReturnStack(
    mastered ? '/review' : learnConceptPath(conceptId),
    returnStack,
  )
  const currentPath = studyConceptPath(conceptId)

  return (
    <div className="mt-6 mb-16 px-4">
      <div className="max-w-2xl mx-auto">
        <div className="mb-4 flex items-center justify-between gap-3">
          {back ? (
            <Link
              href={back.href}
              className="font-mono text-xs text-mathua-secondary hover:text-mathua-blue"
            >
              ← Back to {backLabel || 'previous topic'}
            </Link>
          ) : (
            <Link
              href="/profile"
              className="font-mono text-xs text-mathua-secondary hover:text-mathua-blue"
            >
              ← Profile
            </Link>
          )}
          {stateLabel(readiness?.state) && (
            <span className="font-mono text-[10px] uppercase tracking-wider text-mathua-muted">
              {stateLabel(readiness?.state)}
            </span>
          )}
        </div>

        <h1 className="font-mono text-sm text-mathua-primary break-words [overflow-wrap:anywhere]">
          {label}
        </h1>

        {failed && (
          <p role="alert" className="mt-3 font-mono text-xs text-mathua-red">
            Couldn&apos;t refresh this concept&apos;s status. You can still read the material
            below.
          </p>
        )}

        {prerequisites.length > 0 && (
          <PrerequisitePanel
            prerequisites={prerequisites}
            returnStack={pushReturn(returnStack, currentPath)}
            heading={`Depends on (${prerequisites.length})`}
          />
        )}

        <div className="mt-4">
          <ConceptReference conceptId={conceptId} variant="page" />
        </div>

        {dependents.length > 0 && (
          <div className="mt-4 border border-mathua-border bg-mathua-surface p-4">
            <h3 className="font-mono text-[11px] uppercase tracking-wider text-mathua-muted">
              Unlocks ({dependents.length})
            </h3>
            <ul className="mt-3 space-y-2">
              {dependents.map(c => (
                <li key={c.id} className="flex flex-col sm:flex-row sm:items-center gap-2 min-w-0">
                  <span className="min-w-0 flex-1 font-mono text-[11px] text-mathua-primary break-words [overflow-wrap:anywhere]">
                    {c.label}
                  </span>
                  <Link
                    href={withReturnStack(studyConceptPath(c.id), pushReturn(returnStack, currentPath))}
                    className="w-full sm:w-auto shrink-0 border border-mathua-border px-4 py-2 font-mono text-[11px] text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue min-h-[36px] inline-flex items-center justify-center text-center whitespace-nowrap"
                  >
                    Study →
                  </Link>
                </li>
              ))}
            </ul>
          </div>
        )}

        <div className="mt-6 border border-mathua-border bg-mathua-surface p-4">
          <p className="font-mono text-xs text-mathua-primary">Ready to practise?</p>
          <p className="mt-1 font-mono text-[11px] text-mathua-secondary">
            Test your understanding and work toward mastery. Practice is graded and counts toward
            your progress.
          </p>
          <div className="mt-3 flex flex-wrap gap-3">
            <Link
              href={practiseHref}
              className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint px-5 py-2 font-mono text-xs min-h-[40px] inline-flex items-center"
            >
              Practise this concept →
            </Link>
            {back && (
              <Link
                href={back.href}
                className="border border-mathua-border px-5 py-2 font-mono text-xs text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue min-h-[40px] inline-flex items-center"
              >
                Back to {backLabel || 'previous topic'}
              </Link>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}

export default function StudyPage() {
  return (
    <>
      <Header />
      <div className="pt-[var(--chrome-top)] lg:pt-0">
        <Suspense fallback={<div className="p-6 font-mono text-xs text-mathua-muted">Loading…</div>}>
          <StudyContent />
        </Suspense>
      </div>
      <BottomTabs />
      <Footer />
    </>
  )
}
