'use client'

import { useState } from 'react'
import MathInput from './math/MathInput'
import { Plus } from './icons'
import { GLOSSARY_CATEGORIES, categoryDescription, categoryLabel } from '../lib/glossary/categories'
import { useT } from '../lib/i18n'
import type { GlossaryCategoryId } from '../lib/glossary/types'

export type GlossaryCategoryFilter = GlossaryCategoryId | 'all'

export type GlossarySearchMode = 'text' | 'math'

interface GlossarySearchProps {
  query: string
  category: GlossaryCategoryFilter
  resultCount: number
  rangeStart: number
  rangeEnd: number
  mode: GlossarySearchMode
  mathLatex: string
  onQueryChange: (q: string) => void
  onCategoryChange: (c: GlossaryCategoryFilter) => void
  onModeChange: (m: GlossarySearchMode) => void
  onMathChange: (latex: string) => void
  onClear: () => void
}

// Search + category filter + result count + clear. The count is a live
// status so assistive tech announces filtering without moving focus.
//
// Two input modes share one active query: text for names, aliases and
// definition words; math (MathLive) for notation typed as notation. Only one
// is active at a time, MathLive splits prose into spaced letters rather
// than words, so combining the two would silently corrupt word queries.
export default function GlossarySearch({
  query,
  category,
  resultCount,
  rangeStart,
  rangeEnd,
  mode,
  mathLatex,
  onQueryChange,
  onCategoryChange,
  onModeChange,
  onMathChange,
  onClear,
}: GlossarySearchProps) {
  const t = useT()
  const active = query.trim() !== '' || mathLatex.trim() !== '' || category !== 'all'
  // Category chips hide behind a Tags toggle: nine chips eat a full phone
  // screen before a single result. The toggle is a real button (keyboard
  // accessible, aria-expanded announced); the plus rotates into a cross.
  const [tagsOpen, setTagsOpen] = useState(false)
  return (
    <div className="space-y-3">
      <div role="group" aria-label={t('explore.glossary.search.inputType')} className="flex gap-2">
        {(['text', 'math'] as const).map(m => (
          <button
            key={m}
            type="button"
            onClick={() => onModeChange(m)}
            aria-pressed={mode === m}
            className={`border px-4 py-1 font-mono text-[11px] min-h-[44px] transition-colors ${
              mode === m
                ? 'border-mathua-blue text-mathua-blue'
                : 'border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue'
            }`}
          >
            {m === 'text' ? t('explore.glossary.search.text') : t('explore.glossary.search.math')}
          </button>
        ))}
      </div>
      <div className="flex flex-col gap-2 sm:flex-row">
        <div className="min-w-0 flex-1">
          {mode === 'text' ? (
            <>
              <label htmlFor="glossary-search" className="sr-only">
                {t('explore.glossary.search.label')}
              </label>
              <input
                id="glossary-search"
                type="search"
                autoComplete="off"
                spellCheck={false}
                value={query}
                onChange={e => onQueryChange(e.target.value)}
                placeholder={t('explore.glossary.search.placeholder')}
                className="w-full border border-mathua-border bg-mathua-surface px-4 py-2 font-mono text-sm text-mathua-primary placeholder:text-mathua-muted focus:border-mathua-blue focus:outline-none min-h-[44px]"
              />
            </>
          ) : (
            <MathInput
              value={mathLatex}
              mode="algebra"
              ariaLabel={t('explore.glossary.search.mathLabel')}
              placeholder={t('explore.glossary.search.mathPlaceholder')}
              onChange={v => onMathChange(v.latex)}
            />
          )}
        </div>
        {active && (
          <button
            type="button"
            onClick={onClear}
            className="shrink-0 border border-mathua-border px-4 py-2 font-mono text-xs text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue min-h-[44px]"
          >
            {t('explore.glossary.clear')}
          </button>
        )}
      </div>

      <div>
        <button
          type="button"
          onClick={() => setTagsOpen(o => !o)}
          aria-expanded={tagsOpen}
          aria-controls="glossary-tags"
          className="flex items-center gap-2 border border-mathua-border px-4 py-1 font-mono text-[11px] min-h-[44px] text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue transition-colors"
        >
          <Plus
            size={14}
            aria-hidden="true"
            className={`transition-transform duration-200 motion-reduce:transition-none motion-reduce:transform-none ${tagsOpen ? 'rotate-45' : ''}`}
          />
          {t('explore.glossary.search.tags')}
        </button>
        {tagsOpen && (
          <div id="glossary-tags" className="flex flex-wrap gap-2 mt-2" role="group" aria-label={t('explore.glossary.search.filterGroup')}>
        <button
          type="button"
          onClick={() => onCategoryChange('all')}
          aria-pressed={category === 'all'}
          className={`border px-3 py-1 font-mono text-[11px] min-h-[36px] transition-colors ${
            category === 'all'
              ? 'border-mathua-blue text-mathua-blue'
              : 'border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue'
          }`}
        >
          {t('explore.glossary.search.all')}
        </button>
        {GLOSSARY_CATEGORIES.map(c => (
          <button
            key={c.id}
            type="button"
            onClick={() => onCategoryChange(c.id)}
            aria-pressed={category === c.id}
            title={categoryDescription(c.id)}
            className={`border px-3 py-1 font-mono text-[11px] min-h-[36px] transition-colors ${
              category === c.id
                ? 'border-mathua-blue text-mathua-blue'
                : 'border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue'
            }`}
          >
            {categoryLabel(c.id)}
          </button>
        ))}
          </div>
        )}
      </div>

      <p role="status" aria-live="polite" className="font-mono text-[11px] text-mathua-muted">
        {rangeEnd - rangeStart + 1 === resultCount || resultCount === 0
          ? t('explore.glossary.search.entries', { count: resultCount })
          : t('explore.glossary.search.showing', { start: rangeStart, end: rangeEnd, count: resultCount })}
      </p>
    </div>
  )
}
