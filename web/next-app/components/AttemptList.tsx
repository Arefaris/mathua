'use client'

import Link from 'next/link'
import KatexContent from './KatexContent'
import type { AttemptRecord } from '../lib/api'
import { formatDate, useLocale, useT, type Locale } from '../lib/i18n'
import { sourceName } from '../lib/accountNames'

function shortDate(iso: string, locale: Locale): string {
  if (!iso) return ''
  return formatDate(iso, locale, { month: 'short', day: 'numeric' })
}

// One reviewed answer: the question, what you gave, what was right.
export default function AttemptList({ attempts }: { attempts: AttemptRecord[] }) {
  const t = useT()
  const locale = useLocale()
  if (attempts.length === 0) {
    return (
      <p className="font-mono text-xs text-mathua-muted text-center py-8">
        {t('account.attempts.none')}
      </p>
    )
  }
  return (
    <div className="space-y-3">
      {attempts.map((a, i) => (
        <div
          key={`${a.timestamp}-${a.concept_id}-${i}`}
          className={`border bg-mathua-surface p-4 min-w-0 ${a.correct ? 'border-mathua-border' : 'border-mathua-red'}`}
        >
          <div className="flex items-center justify-between gap-2 mb-2">
            <span className="font-mono text-[10px] uppercase tracking-wider text-mathua-muted">
              {sourceName(a.source || 'practice')}
            </span>
            <span className="font-mono text-[10px] text-mathua-muted shrink-0">
              {shortDate(a.timestamp, locale)}
            </span>
          </div>
          {a.question ? (
            <KatexContent className="text-sm text-mathua-primary font-mono whitespace-pre-wrap break-words">
              {a.question}
            </KatexContent>
          ) : (
            <p className="font-mono text-xs text-mathua-muted">{t('account.attempts.unavailable')}</p>
          )}
          <p className="mt-2 font-mono text-xs break-words [overflow-wrap:anywhere]">
            <span className={a.correct ? 'text-mathua-green' : 'text-mathua-red'}>
              {a.correct ? '✓' : '✗'} {t('account.attempts.you', { answer: a.answer === '' ? '—' : a.answer })}
            </span>
            {!a.correct && (
              <span className="text-mathua-secondary"> · {t('account.attempts.correct', { expected: a.expected })}</span>
            )}
          </p>
          {!a.correct && a.explanation && (
            <KatexContent className="mt-1 text-[11px] font-mono text-mathua-secondary whitespace-pre-wrap break-words">
              {a.explanation}
            </KatexContent>
          )}
          <Link
            href={`/learn?concept=${encodeURIComponent(a.concept_id)}`}
            className="mt-2 inline-block font-mono text-[11px] text-mathua-blue hover:text-mathua-blue-hover"
          >
            {a.concept_name || a.concept_id} →
          </Link>
        </div>
      ))}
    </div>
  )
}
