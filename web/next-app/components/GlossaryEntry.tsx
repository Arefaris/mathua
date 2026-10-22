'use client'

import KatexContent from './KatexContent'
import { categoryLabel } from '../lib/glossary/categories'
import { glossaryById } from '../lib/glossary/data'
import { useT } from '../lib/i18n'
import type { GlossaryEntry } from '../lib/glossary/types'

interface GlossaryEntryProps {
  entry: GlossaryEntry
  onSelectRelated: (id: string) => void
}

// One glossary result. The symbol is always paired with its spoken name —
// never symbol-alone — so the meaning survives without the glyph rendering.
export default function GlossaryEntryCard({ entry, onSelectRelated }: GlossaryEntryProps) {
  const t = useT()
  return (
    <article
      id={`glossary-${entry.id}`}
      aria-labelledby={`glossary-${entry.id}-name`}
      className="border border-mathua-border bg-mathua-surface px-4 py-4 scroll-mt-24"
    >
      <div className="flex flex-wrap items-baseline gap-x-4 gap-y-1">
        <span
          role="img"
          aria-label={t('explore.glossary.entry.symbolLabel', { name: entry.name, symbol: entry.symbol })}
          className="font-mono text-2xl text-mathua-blue"
        >
          {entry.symbol}
        </span>
        <h3
          id={`glossary-${entry.id}-name`}
          className="font-mono text-sm text-mathua-primary"
        >
          {entry.name}
        </h3>
      </div>

      <p className="mt-2 font-mono text-xs leading-relaxed text-mathua-secondary">
        {entry.definition}
      </p>

      {entry.reading && (
        <p className="mt-1 font-mono text-[11px] text-mathua-muted">
          {t('explore.glossary.entry.readAs', { reading: entry.reading })}
        </p>
      )}

      <div className="mt-3 space-y-2">
        {entry.examples.map((ex, i) => (
          <div key={i} className="bg-mathua-code px-3 py-2">
            <KatexContent>{`$${ex.tex}$`}</KatexContent>
            <p className="mt-1 font-mono text-[11px] text-mathua-secondary">
              {ex.explanation}
            </p>
          </div>
        ))}
      </div>

      {entry.contexts && entry.contexts.length > 0 && (
        <dl className="mt-3 space-y-1">
          {entry.contexts.map((c, i) => (
            <div key={i} className="flex flex-col gap-0.5 sm:flex-row sm:gap-2">
              <dt className="shrink-0 font-mono text-[11px] uppercase tracking-wider text-mathua-muted">
                {c.context}
              </dt>
              <dd className="font-mono text-[11px] text-mathua-secondary">{c.meaning}</dd>
            </div>
          ))}
        </dl>
      )}

      {entry.commonMistakes && entry.commonMistakes.length > 0 && (
        <ul className="mt-3 space-y-1">
          {entry.commonMistakes.map((m, i) => (
            <li key={i} className="font-mono text-[11px] text-mathua-secondary">
              <span className="text-mathua-muted">{t('explore.glossary.entry.mistake')}</span>
              {m}
            </li>
          ))}
        </ul>
      )}

      {entry.relatedSymbols && entry.relatedSymbols.length > 0 && (
        <div className="mt-3 flex flex-wrap items-center gap-2">
          <span className="font-mono text-[11px] uppercase tracking-wider text-mathua-muted">
            {t('explore.glossary.entry.related')}
          </span>
          {entry.relatedSymbols.map(id => {
            const target = glossaryById(id)
            if (!target) return null
            return (
              <button
                key={id}
                type="button"
                onClick={() => onSelectRelated(id)}
                aria-label={t('explore.glossary.entry.show', { name: target.name })}
                className="border border-mathua-border px-2 py-1 font-mono text-[11px] text-mathua-blue hover:bg-mathua-blue-faint min-h-[36px]"
              >
                {target.symbol} {target.name}
              </button>
            )
          })}
        </div>
      )}

      <p className="mt-3 font-mono text-[10px] text-mathua-muted">
        {categoryLabel(entry.category)}
        {entry.latex && <span> · LaTeX: <code>{entry.latex}</code></span>}
        {t('explore.glossary.entry.sources')}
        {entry.sources.map((s, i) => (
          <span key={i}>
            {i > 0 && '; '}
            <a
              href={s.url}
              target="_blank"
              rel="noreferrer"
              title={s.license}
              className="underline hover:text-mathua-blue"
            >
              {s.label}
            </a>
          </span>
        ))}
      </p>

      {/* The license was suppressed here because the full reuse terms are long and repeat on
          every card. Attribution is a condition of several of these sources, so it belongs on
          the card and not only on /licenses — the note is the source and its license, once per
          source, in the smallest type the card uses. */}
      {entry.sources.some(s => s.license) && (
        <p className="mt-1 font-mono text-[10px] text-mathua-muted">
          {entry.sources.filter(s => s.license).map((s, i) => (
            <span key={i}>
              {i > 0 && '; '}
              {s.label}: {s.license}
            </span>
          ))}
        </p>
      )}

      {entry.notes && (
        <p className="mt-1 font-mono text-[10px] text-mathua-muted">{entry.notes}</p>
      )}
    </article>
  )
}
