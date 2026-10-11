'use client'

import { GLOSSARY_CATEGORIES } from '../lib/glossary/categories'
import type { GlossaryCategoryId } from '../lib/glossary/types'

export type GlossaryCategoryFilter = GlossaryCategoryId | 'all'

interface GlossarySearchProps {
  query: string
  category: GlossaryCategoryFilter
  resultCount: number
  onQueryChange: (q: string) => void
  onCategoryChange: (c: GlossaryCategoryFilter) => void
  onClear: () => void
}

// Search + category filter + result count + clear. The count is a live
// status so assistive tech announces filtering without moving focus.
export default function GlossarySearch({
  query,
  category,
  resultCount,
  onQueryChange,
  onCategoryChange,
  onClear,
}: GlossarySearchProps) {
  const active = query.trim() !== '' || category !== 'all'
  return (
    <div className="space-y-3">
      <div className="flex flex-col gap-2 sm:flex-row">
        <div className="min-w-0 flex-1">
          <label htmlFor="glossary-search" className="sr-only">
            Search symbols, names or meanings
          </label>
          <input
            id="glossary-search"
            type="search"
            autoComplete="off"
            spellCheck={false}
            value={query}
            onChange={e => onQueryChange(e.target.value)}
            placeholder="Search a symbol, name or meaning — e.g. ∀, integral, lambda"
            className="w-full border border-mathua-border bg-mathua-surface px-4 py-2 font-mono text-sm text-mathua-primary placeholder:text-mathua-muted focus:border-mathua-blue focus:outline-none min-h-[44px]"
          />
        </div>
        {active && (
          <button
            type="button"
            onClick={onClear}
            className="shrink-0 border border-mathua-border px-4 py-2 font-mono text-xs text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue min-h-[44px]"
          >
            Clear search and filters
          </button>
        )}
      </div>

      <div className="flex flex-wrap gap-2" role="group" aria-label="Filter by category">
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
          All
        </button>
        {GLOSSARY_CATEGORIES.map(c => (
          <button
            key={c.id}
            type="button"
            onClick={() => onCategoryChange(c.id)}
            aria-pressed={category === c.id}
            title={c.description}
            className={`border px-3 py-1 font-mono text-[11px] min-h-[36px] transition-colors ${
              category === c.id
                ? 'border-mathua-blue text-mathua-blue'
                : 'border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue'
            }`}
          >
            {c.label}
          </button>
        ))}
      </div>

      <p role="status" aria-live="polite" className="font-mono text-[11px] text-mathua-muted">
        {resultCount === 1 ? '1 entry' : `${resultCount} entries`}
      </p>
    </div>
  )
}
