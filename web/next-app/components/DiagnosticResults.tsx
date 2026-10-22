'use client'

import { useState } from 'react'
import Link from 'next/link'
import type { GoalPlanRes } from '../lib/api'
import { domainName } from '../lib/domainName'
import { useT } from '../lib/i18n'


interface Props {
  plan: GoalPlanRes
  onStartPractice?: () => void
}

export default function DiagnosticResults({ plan, onStartPractice }: Props) {
  const t = useT()
  const [expanded, setExpanded] = useState<string | null>(null)

  const domains = new Set([...Object.keys(plan.weak_areas), ...Object.keys(plan.strong_areas)])
  const readinessPct = Math.round(plan.readiness * 100)
  const estimates = plan.completion_estimates ? Object.entries(plan.completion_estimates).sort((a, b) => Number(a[0]) - Number(b[0])) : []

  return (
    <div className="max-w-3xl mx-auto">
      {(plan.frontier_label || plan.placement_course_id) && (
        <div className="border border-mathua-blue bg-mathua-surface p-4 mb-6 text-center">
          {plan.frontier_label && (
            <p className="font-serif text-lg text-mathua-primary">
              {t('goals.results.startHere', { label: plan.frontier_label })}
              {plan.frontier_conditional && <span className="font-mono text-xs text-mathua-muted">{t('goals.results.provisional')}</span>}
            </p>
          )}
          {plan.placement_course_id && (
            <p className="font-mono text-xs text-mathua-secondary mt-1">{t('goals.results.suggestedCourse', { course: plan.placement_course_id })}</p>
          )}
          {estimates.length > 0 && (
            <p className="font-mono text-[11px] text-mathua-muted mt-2">
              {t('goals.results.estimated', { list: estimates.map(([xp, when]) => `${xp} XP → ${when}`).join(' · ') })}
            </p>
          )}
          <p className="font-mono text-[10px] text-mathua-muted mt-2">
            {t('goals.results.automaticity')}
          </p>
          {plan.conditionally_completed && plan.conditionally_completed.length > 0 && (
            <p className="font-mono text-[10px] text-mathua-muted mt-1">
              {t('goals.results.provisionalTopics', { count: plan.conditionally_completed.length })}
            </p>
          )}
        </div>
      )}
      <div className="text-center mb-8">
        <div className="relative inline-flex items-center justify-center">
          <svg width="140" height="140" className="-rotate-90">
            <circle cx="70" cy="70" r="60" fill="none" stroke="var(--code-bg)" strokeWidth="10" />
            <circle cx="70" cy="70" r="60" fill="none" stroke="var(--accent-blue)" strokeWidth="10" strokeLinecap="round" strokeDasharray={`${2 * Math.PI * 60}`} strokeDashoffset={`${2 * Math.PI * 60 * (1 - plan.readiness)}`} className="transition-all duration-1000" />
          </svg>
          <span className="absolute text-3xl font-mono font-light text-mathua-primary">{readinessPct}%</span>
        </div>
        <p className="text-mathua-secondary text-sm mt-2">{t('goals.results.score')} {plan.total_tested > 0 && t('goals.results.scoreDetail', { correct: plan.correct_count, total: plan.total_tested })}</p>
      </div>

      <div className="space-y-3 mb-8">
        {Array.from(domains).sort().map((domain) => {
          const label = domainName(domain)
          const weak = plan.weak_areas[domain] || []
          const strong = plan.strong_areas[domain] || []
          const total = weak.length + strong.length
          const weakPct = total > 0 ? Math.round((weak.length / total) * 100) : 0
          const isExpanded = expanded === domain
          return (
            <div key={domain} className="border border-mathua-border bg-mathua-surface">
              <button type="button" onClick={() => setExpanded(isExpanded ? null : domain)} aria-expanded={isExpanded} className="w-full flex items-center justify-between p-3 text-left">
                <span className="font-mono text-xs text-mathua-primary">{label}</span>
                <span className="font-mono text-[11px] text-mathua-muted">{t('goals.results.toReview', { weak: weak.length, strong: strong.length })}</span>
              </button>
              <div className="px-3 pb-2">
                <div className="h-1.5 bg-mathua-code overflow-hidden">
                  <div className="h-full bg-mathua-blue transition-all" style={{ width: `${weakPct}%` }} />
                </div>
              </div>
              {isExpanded && (
                <div className="px-3 pb-3 border-t border-mathua-border mt-1 pt-2">
                  {weak.length > 0 && (
                    <div className="mb-2">
                      <p className="font-mono text-[10px] uppercase text-mathua-red mb-1">{t('goals.results.needsReview')}</p>
                      <div className="flex flex-wrap gap-1.5">
                        {weak.map((c) => (
                          <Link key={c.id} href={`/learn?concept=${encodeURIComponent(c.id)}`} className="text-xs font-mono border border-mathua-border px-2 py-1 hover:border-mathua-blue hover:text-mathua-blue">
                            {c.label}
                          </Link>
                        ))}
                      </div>
                    </div>
                  )}
                  {strong.length > 0 && (
                    <div>
                      <p className="font-mono text-[10px] uppercase text-mathua-green mb-1">{t('goals.results.strong')}</p>
                      <p className="font-mono text-xs text-mathua-secondary">{strong.join(', ')}</p>
                    </div>
                  )}
                </div>
              )}
            </div>
          )
        })}
      </div>

      {onStartPractice && (
        <div className="text-center">
          <button type="button" onClick={onStartPractice} className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint px-10 h-12 font-mono text-sm">
            {t('goals.results.startPracticing')}
          </button>
          <div className="mt-3">
            <Link href="/history" className="font-mono text-xs text-mathua-muted hover:text-mathua-primary">
              {t('goals.question.historyLink')}
            </Link>
          </div>
        </div>
      )}
    </div>
  )
}
