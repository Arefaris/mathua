'use client'

import { useCallback, useEffect, useMemo, useRef, useState } from 'react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import Header from '../../components/Header'
import BottomTabs from '../../components/BottomTabs'
import SectionHeader from '../../components/SectionHeader'
import Footer from '../../components/Footer'
import Loading from '../../components/Loading'
import PrerequisitePathDialog, { type PrereqDialogStatus } from '../../components/PrerequisitePathDialog'
import { pushReturn, studyConceptPath } from '../../lib/learnReturn'
import { useTheme } from '../../hooks/useTheme'
import {
  getCurriculumDomains,
  getCurriculumDomain,
  getLessonEligibility,
  type DomainSummary,
  type DomainDetail,
  type TopicSummary,
  type EligibilityPrereq,
} from '../../lib/api'
import { DOMAIN_ORDER } from '../../lib/graphDomains'
import { domainName } from '../../lib/domainName'
import { useT, type TFunction } from '../../lib/i18n'

/**
 * The curriculum browser: domain -> topic, driven entirely by the server's
 * curriculum endpoints. The client never derives a state from raw progress; it
 * renders `state` and asks the server before entering the learning loop.
 *
 * Domains enter from `/profile`'s "Choose what to learn", and they are also a
 * bottom tab, so this is reachable without going through Settings or Study.
 */

// The server's state ids, as dictionary keys. An id this list does not know is shown as sent.
const STATE_KEYS: Record<string, string> = {
  mastered: 'mastered',
  due_for_review: 'dueForReview',
  in_progress: 'inProgress',
  unlocked: 'unlocked',
  upcoming: 'upcoming',
  locked: 'locked',
}

function stateLabel(t: TFunction, state: string): string {
  const key = STATE_KEYS[state]
  return key ? t.dynamic(`explore.domains.state.${key}`, undefined, state) : state
}

function available(row: DomainSummary): number {
  return row.unlockedCount + row.inProgressCount
}

export default function DomainsPage() {
  const t = useT()
  const { mounted } = useTheme()
  const router = useRouter()

  const [domains, setDomains] = useState<DomainSummary[] | null>(null)
  const [selected, setSelected] = useState<string | null>(null)
  const [detail, setDetail] = useState<DomainDetail | null>(null)
  const [detailLoading, setDetailLoading] = useState(false)
  const [busy, setBusy] = useState<string | null>(null)
  const [dialog, setDialog] = useState<{ conceptId: string; title: string } | null>(null)
  const [prereqs, setPrereqs] = useState<EligibilityPrereq[]>([])
  const [dlgStatus, setDlgStatus] = useState<PrereqDialogStatus>('loading')
  const [dlgError, setDlgError] = useState('')
  const [error, setError] = useState('')
  // Monotonic request id: a late response for concept A must never overwrite
  // the dialog now showing concept B.
  const seqRef = useRef(0)
  // Element to return focus to when the dialog closes.
  const triggerRef = useRef<HTMLButtonElement | null>(null)

  useEffect(() => {
    if (!mounted) return
    getCurriculumDomains()
      .then(setDomains)
      .catch(() => {
        setDomains([])
        setError(t('explore.domains.loadError'))
      })
  }, [mounted, t])

  const ordered = useMemo(() => {
    if (!domains) return []
    const byId = new Map(domains.map((d) => [d.id, d]))
    const out: DomainSummary[] = []
    for (const id of DOMAIN_ORDER) {
      const d = byId.get(id)
      if (d) out.push(d)
    }
    // Any domain the corpus grows that DOMAIN_ORDER does not yet name still
    // appears, rather than silently vanishing from the picker.
    for (const d of domains) if (!out.includes(d)) out.push(d)
    return out
  }, [domains])

  const openDomain = useCallback(async (id: string) => {
    setSelected(id)
    setDialog(null)
    setDetail(null)
    setDetailLoading(true)
    try {
      setDetail(await getCurriculumDomain(id))
    } finally {
      setDetailLoading(false)
    }
  }, [])

  async function baseEligibility(conceptId: string): Promise<boolean | EligibilityPrereq[] | null> {
    const el = await getLessonEligibility(conceptId)
    if (el.eligible === false) return el.prerequisites ?? []
    return true
  }

  async function chooseTopic(topic: TopicSummary) {
    setError('')
    setDialog(null)
    if (topic.state === 'mastered' || topic.state === 'due_for_review') {
      router.push('/review')
      return
    }
    setBusy(topic.conceptId)
    try {
      const result = await baseEligibility(topic.conceptId)
      if (result === true) {
        router.push(`/learn?concept=${encodeURIComponent(topic.conceptId)}`)
      } else if (Array.isArray(result)) {
        seqRef.current++
        setDialog({ conceptId: topic.conceptId, title: topic.title })
        setPrereqs(result)
        setDlgStatus(result.length === 0 ? 'empty' : 'ready')
        setDlgError('')
      }
    } catch {
      setError(t('explore.domains.startError'))
    } finally {
      setBusy(null)
    }
  }

  // Selecting a locked topic explains the prerequisite path instead of starting.
  // The dialog opens immediately in a loading state so the action is visible
  // regardless of where the row sits in a long list; the fetch then fills it.
  // Only the latest request may write, so rapid clicks on different rows can
  // never leave stale prerequisites on screen.
  async function explainTopic(topic: TopicSummary, trigger?: HTMLButtonElement | null) {
    setError('')
    setDlgError('')
    triggerRef.current = trigger ?? null
    const seq = ++seqRef.current
    setDialog({ conceptId: topic.conceptId, title: topic.title })
    setPrereqs([])
    setDlgStatus('loading')
    setBusy(topic.conceptId)
    try {
      const el = await getLessonEligibility(topic.conceptId)
      if (seqRef.current !== seq) return
      const list = el.prerequisites ?? []
      setPrereqs(list)
      setDlgStatus(list.length === 0 ? 'empty' : 'ready')
    } catch {
      if (seqRef.current !== seq) return
      setDlgStatus('error')
      setDlgError(t('explore.domains.prereqError'))
    } finally {
      if (seqRef.current === seq) setBusy(null)
    }
  }

  function retryDialog() {
    if (!dialog) return
    void explainTopic(
      { conceptId: dialog.conceptId, title: dialog.title, state: 'locked', depth: 0 },
      triggerRef.current,
    )
  }

  function closeDialog() {
    seqRef.current++
    setDialog(null)
    setBusy(null)
    triggerRef.current?.focus()
    triggerRef.current = null
  }

  if (!mounted || domains === null) return <Loading label={t('explore.domains.loading')} full />

  const selectedRow = ordered.find((d) => d.id === selected) ?? null

  return (
    <>
      <Header />
      <main className="mx-auto w-full max-w-[820px] min-w-0 px-4 sm:px-6 py-8 sm:py-12 overflow-x-hidden">
        {selected === null ? (
          <>
            <SectionHeader label={t('explore.domains.label')} title={t('explore.domains.title')} />
            <p className="font-mono text-[11px] text-mathua-secondary mb-6">
              {t('explore.domains.intro')}
            </p>

            {error && (
              <p role="alert" className="mb-4 font-mono text-[11px] text-mathua-secondary">
                {error}
              </p>
            )}

            <ul className="border border-mathua-border divide-y divide-mathua-border bg-mathua-surface">
              {ordered.map((d) => (
                <li key={d.id}>
                  <button
                    type="button"
                    onClick={() => void openDomain(d.id)}
                    className="w-full text-left px-4 py-3 hover:bg-mathua-blue-faint transition-colors"
                  >
                    <span className="flex items-baseline justify-between gap-3 min-w-0">
                      <span className="font-mono text-xs text-mathua-primary truncate">
                        {domainName(d.id)}
                      </span>
                      <span className="shrink-0 font-mono text-[11px] text-mathua-secondary">
                        {t('explore.domains.available', { available: available(d), count: d.conceptCount })}
                      </span>
                    </span>
                  </button>
                </li>
              ))}
            </ul>
          </>
        ) : (
          <>
            <button
              type="button"
              onClick={() => { setSelected(null); setDetail(null); setDialog(null) }}
              className="font-mono text-xs text-mathua-secondary hover:text-mathua-blue"
            >
              {t('explore.domains.back')}
            </button>
            <div className="mt-4">
              <SectionHeader
                label={t('explore.domains.label')}
                title={selectedRow ? domainName(selectedRow.id) : selected}
              />
            </div>

            {selectedRow && (
              <p className="font-mono text-[11px] text-mathua-secondary mb-6">
                {t('explore.domains.progress', { available: available(selectedRow), count: selectedRow.conceptCount })}
                {selectedRow.masteredCount > 0 && t('explore.domains.masteredSuffix', { count: selectedRow.masteredCount })}
              </p>
            )}

            {detailLoading && <p className="font-mono text-xs text-mathua-muted">{t('explore.domains.loadingTopics')}</p>}

            {!detailLoading && detail && detail.topics.length === 0 && (
              <p className="font-mono text-xs text-mathua-secondary">
                {t('explore.domains.noTopics')}
              </p>
            )}

            {!detailLoading && detail && detail.topics.length > 0 && (
              <>
                {!detail.topics.some((topic) =>
                  ['unlocked', 'in_progress', 'mastered', 'due_for_review'].includes(topic.state),
                ) && (
                  <p className="mb-4 border border-mathua-border bg-mathua-surface p-3 font-mono text-[11px] text-mathua-secondary">
                    {t('explore.domains.noTopics')}
                  </p>
                )}
                <ul className="border border-mathua-border divide-y divide-mathua-border bg-mathua-surface">
                  {detail.topics.map((topic) => {
                  const learnable = topic.state === 'unlocked' || topic.state === 'in_progress'
                  const review = topic.state === 'mastered' || topic.state === 'due_for_review'
                  return (
                    <li key={topic.conceptId} className="px-4 py-3 min-w-0">
                      <div className="flex flex-col sm:flex-row sm:items-center gap-2 sm:gap-4 min-w-0">
                        <div className="min-w-0 flex-1">
                          <div className="font-mono text-xs text-mathua-primary break-words [overflow-wrap:anywhere]">
                            {topic.title}
                          </div>
                          {/* State is text first, so it survives without color. */}
                          <div className="font-mono text-[10px] uppercase tracking-wider text-mathua-muted">
                            {stateLabel(t, topic.state)}
                          </div>
                        </div>
                        {learnable && (
                          <button
                            type="button"
                            disabled={busy === topic.conceptId}
                            onClick={() => void chooseTopic(topic)}
                            className="w-full sm:w-auto shrink-0 border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint px-4 py-2 font-mono text-[11px] min-h-[36px] inline-flex items-center justify-center disabled:opacity-60"
                          >
                            {busy === topic.conceptId ? t('explore.domains.checking') : t('explore.domains.learn')}
                          </button>
                        )}
                        {review && (
                          <Link
                            href="/review"
                            className="w-full sm:w-auto shrink-0 border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint px-4 py-2 font-mono text-[11px] min-h-[36px] inline-flex items-center justify-center"
                          >
                            {t('explore.domains.review')}
                          </Link>
                        )}
                        {!learnable && !review && (
                          <button
                            type="button"
                            onClick={(e) => void explainTopic(topic, e.currentTarget)}
                            className="w-full sm:w-auto shrink-0 border border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue px-4 py-2 font-mono text-[11px] min-h-[36px] inline-flex items-center justify-center"
                          >
                            {t('explore.domains.prereqPath')}
                          </button>
                        )}
                      </div>
                    </li>
                  )
                })}
              </ul>
              </>
            )}

            <PrerequisitePathDialog
              open={dialog !== null}
              targetTitle={dialog?.title ?? ''}
              prerequisites={prereqs}
              status={dlgStatus}
              error={dlgError}
              returnStack={dialog ? pushReturn([], studyConceptPath(dialog.conceptId)) : []}
              onClose={closeDialog}
              onRetry={retryDialog}
            />

            {error && (
              <p role="alert" className="mt-4 font-mono text-[11px] text-mathua-secondary">
                {error}
              </p>
            )}
          </>
        )}
      </main>
      <Footer />
      <BottomTabs />
    </>
  )
}
