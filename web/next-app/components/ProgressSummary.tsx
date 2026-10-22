import type { Scores } from '../lib/api'
import { domainName } from '../lib/domainName'
import { useT } from '../lib/i18n'
import { levelName } from '../lib/accountNames'

interface ProgressSummaryProps {
  scores: Scores
  weakByDomain?: Record<string, { id: string; label: string }[]>
}



export default function ProgressSummary({ scores, weakByDomain }: ProgressSummaryProps) {
  const t = useT()
  const stats = [
    { label: t('account.progressSummary.mastered'), value: scores.concepts_mastered, color: 'var(--accent-blue)' },
    { label: t('account.progressSummary.streak'), value: scores.current_streak, color: 'var(--accent-teal)' },
    { label: t('account.progressSummary.level'), value: levelName(scores.level), color: 'var(--accent-blue)' },
    { label: t('account.progressSummary.weekly'), value: scores.weekly_score, color: 'var(--accent-blue)' },
  ]

  const weakDomains = weakByDomain
    ? Object.entries(weakByDomain).filter(([, concepts]) => concepts.length > 0).sort((a, b) => b[1].length - a[1].length).slice(0, 5)
    : []

  return (
    <>
      <div className="flex flex-wrap gap-4 justify-center mb-6">
        {stats.map((s) => (
          <div
            key={s.label}
            className="bg-mathua-surface border border-mathua-border p-5 flex-1 min-w-[140px] max-w-[200px] text-center"
          >
            <div
              className="font-mono text-2xl font-light mt-1"
              style={{ color: s.color }}
            >
              {s.value}
            </div>
            <div className="font-mono text-[10px] uppercase tracking-[0.1em] text-mathua-muted mt-1">
              {s.label}
            </div>
          </div>
        ))}
      </div>
      {weakDomains.length > 0 && (
        <div className="flex flex-wrap gap-3 justify-center mb-6">
          {weakDomains.map(([domain, concepts]) => {
            const label = domainName(domain)
            return (
              <div key={domain} className="bg-mathua-surface border border-mathua-border px-4 py-2 text-center min-w-[120px]">
                <div className="font-mono text-[10px] uppercase tracking-[0.1em] text-mathua-muted">{label}</div>
                <div className="font-mono text-sm text-mathua-red mt-0.5">{t('account.progressSummary.needReview', { count: concepts.length })}</div>
              </div>
            )
          })}
        </div>
      )}
    </>
  )
}
