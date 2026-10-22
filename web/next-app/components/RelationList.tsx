'use client'

import Link from 'next/link'
import { studyConceptPath, withReturnStack, type ReturnEntry } from '../lib/learnReturn'
import { useT } from '../lib/i18n'

/**
 * One DAG-neighbourhood list on the reference page: the direct prerequisites ("Depends on") or
 * the direct dependents ("Unlocks").
 *
 * Every row is a link, and a `Locked` state means "not eligible to practise" — never "cannot
 * read". A row that looked disabled because its concept is locked would be the exact dead end
 * ADR-054 removed, so the row's state is a label and the row itself is always interactive.
 *
 * This is a display of the DAG, not a second model of it: the caller supplies the direct
 * neighbours and their state, and the component only decides how they read.
 */

export interface RelationItem {
  id: string
  title: string
  /** Learner-facing state word, when the backend (or the learner's own history) reports one. */
  state?: string
}

interface Props {
  id?: string
  heading: string
  items: RelationItem[]
  /** Neutral copy shown when the list is empty, so an empty section is not a hole in the page. */
  emptyText: string
  /** The list is not known yet. A loading line, never the empty text, so "none" is not a guess. */
  loading?: boolean
  returnStack?: ReturnEntry[]
}

export default function RelationList({ id, heading, items, emptyText, loading = false, returnStack = [] }: Props) {
  const t = useT()
  return (
    <section id={id} aria-label={heading} className="min-w-0 scroll-mt-24">
      <h2 className="font-mono text-[10px] uppercase tracking-wider text-mathua-muted">
        {heading}
      </h2>
      {loading ? (
        <p className="mt-2 font-mono text-[11px] text-mathua-muted">{t('common.status.loading')}</p>
      ) : items.length === 0 ? (
        <p className="mt-2 font-mono text-[11px] text-mathua-muted">{emptyText}</p>
      ) : (
        <ul className="mt-2 border-[0.5px] border-mathua-border">
          {items.map((it, i) => (
            <li
              key={it.id}
              className={i < items.length - 1 ? 'border-b-[0.5px] border-mathua-border' : ''}
            >
              <Link
                href={withReturnStack(studyConceptPath(it.id), returnStack)}
                className="group flex items-center gap-3 px-3 py-2.5 min-w-0 hover:bg-mathua-surface-highlight focus-visible:bg-mathua-surface-highlight focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-inset focus-visible:ring-mathua-blue"
              >
                <span className="min-w-0 flex-1">
                  <span className="block font-mono text-[11px] text-mathua-primary break-words [overflow-wrap:anywhere]">
                    {it.title}
                  </span>
                  {it.state && (
                    <span className="mt-0.5 block font-mono text-[10px] uppercase tracking-wider text-mathua-muted">
                      {it.state}
                    </span>
                  )}
                </span>
                <span className="shrink-0 font-mono text-[11px] text-mathua-blue whitespace-nowrap group-hover:underline">
                  {t('learn.relation.study')}
                </span>
              </Link>
            </li>
          ))}
        </ul>
      )}
    </section>
  )
}
