'use client'

import { useCallback, useMemo, useState } from 'react'
import Header from '../../components/Header'
import BottomTabs from '../../components/BottomTabs'
import SectionHeader from '../../components/SectionHeader'
import Footer from '../../components/Footer'
import Loading from '../../components/Loading'
import GlossarySearch, { type GlossaryCategoryFilter } from '../../components/GlossarySearch'
import GlossaryEntryCard from '../../components/GlossaryEntry'
import { useTheme } from '../../hooks/useTheme'
import { GLOSSARY_ENTRIES } from '../../lib/glossary/data'
import { GLOSSARY_CATEGORIES } from '../../lib/glossary/categories'
import { filterEntries } from '../../lib/glossary/search'

/**
 * The mathematical glossary: a standalone reference for symbols, notation
 * and terminology. It reads the static bundle only — no API, no auth, no
 * XP, no progress — so looking anything up cannot affect learning state.
 */
export default function GlossaryPage() {
  const { mounted } = useTheme()
  const [query, setQuery] = useState('')
  const [category, setCategory] = useState<GlossaryCategoryFilter>('all')

  const results = useMemo(
    () => filterEntries(GLOSSARY_ENTRIES, { query, category }),
    [query, category],
  )

  // With no filtering the page browses by category; once the learner types
  // or picks a category, results flatten into relevance order.
  const browsing = query.trim() === '' && category === 'all'
  const grouped = useMemo(() => {
    if (!browsing) return []
    return GLOSSARY_CATEGORIES.map(c => ({
      category: c,
      entries: GLOSSARY_ENTRIES.filter(e => e.category === c.id),
    })).filter(g => g.entries.length > 0)
  }, [browsing])

  function clear() {
    setQuery('')
    setCategory('all')
  }

  const selectRelated = useCallback((id: string) => {
    clear()
    // Let the full list render before scrolling to the card.
    requestAnimationFrame(() => {
      document.getElementById(`glossary-${id}`)?.scrollIntoView({ behavior: 'smooth', block: 'start' })
      document.getElementById(`glossary-${id}-name`)?.focus?.()
    })
  }, [])

  if (!mounted) return <Loading label="Loading glossary" full />

  return (
    <>
      <Header />
      <main className="mx-auto w-full max-w-[820px] min-w-0 px-4 sm:px-6 py-8 sm:py-12 overflow-x-hidden">
        <SectionHeader label="Reference" title="Mathematical glossary" />
        <p className="font-mono text-[11px] text-mathua-secondary mb-6">
          Look up symbols, Greek letters, notation and terminology — from
          arithmetic through university mathematics. Browsing here never
          affects your progress, XP or mastery.
        </p>

        <GlossarySearch
          query={query}
          category={category}
          resultCount={results.length}
          onQueryChange={setQuery}
          onCategoryChange={setCategory}
          onClear={clear}
        />

        <div className="mt-6">
          {results.length === 0 ? (
            <div className="border border-mathua-border bg-mathua-surface px-4 py-8 text-center">
              <p className="font-mono text-xs text-mathua-primary">
                No entries match{query.trim() && <> “{query.trim()}”</>}.
              </p>
              <p className="mt-2 font-mono text-[11px] text-mathua-secondary">
                Try a symbol (∀, ∫, λ), a name (integral, lambda) or a word
                from a definition.
              </p>
              <button
                type="button"
                onClick={clear}
                className="mt-4 border border-mathua-blue px-4 py-2 font-mono text-[11px] text-mathua-blue hover:bg-mathua-blue-faint min-h-[44px]"
              >
                Clear search and filters
              </button>
            </div>
          ) : browsing ? (
            <div className="space-y-8">
              {grouped.map(g => (
                <section key={g.category.id} aria-labelledby={`glossary-cat-${g.category.id}`}>
                  <h2
                    id={`glossary-cat-${g.category.id}`}
                    className="mb-1 font-mono text-[11px] uppercase tracking-[0.15em] text-mathua-muted"
                  >
                    {g.category.label} · {g.entries.length}
                  </h2>
                  <p className="mb-3 font-mono text-[11px] text-mathua-secondary">
                    {g.category.description}
                  </p>
                  <div className="space-y-4">
                    {g.entries.map(e => (
                      <GlossaryEntryCard key={e.id} entry={e} onSelectRelated={selectRelated} />
                    ))}
                  </div>
                </section>
              ))}
            </div>
          ) : (
            <div className="space-y-4">
              {results.map(e => (
                <GlossaryEntryCard key={e.id} entry={e} onSelectRelated={selectRelated} />
              ))}
            </div>
          )}
        </div>
      </main>
      <Footer />
      <BottomTabs />
    </>
  )
}
