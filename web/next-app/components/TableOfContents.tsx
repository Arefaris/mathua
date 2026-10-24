'use client'

import { useState } from 'react'
import { useActiveHeading } from '../hooks/useActiveHeading'
import { Sheet, SheetContent, SheetTitle, SheetTrigger } from './ui/sheet'
import { Layers } from './icons'
import { useT } from '../lib/i18n'

// The table of contents, in one place for the docs pages and app/how-it-works.
//
// Desktop: a sticky sidebar that stays beside the reader. Mobile: a floating button that opens
// the contents in a bottom sheet. The mobile form used to be a full-width horizontal bar pinned
// under the header (`fixed top-[57px] z-[99]`), which sat beneath an 80px header, overlapped it,
// and asked a reader to scroll nine labels sideways on a phone. A sheet gives each entry a full
// row and a real tap target.

export interface TocItem {
  id: string
  text: string
}

const mono = "var(--font-jetbrains-mono), 'JetBrains Mono', monospace"

export default function TableOfContents({
  headings,
  label,
}: {
  headings: TocItem[]
  label?: string
}) {
  const t = useT()
  const heading = label ?? t('shell.toc.onThisPage')
  const active = useActiveHeading(headings.map(h => h.id))
  const [open, setOpen] = useState(false)

  if (headings.length === 0) return null

  // Close the sheet first, then move. `scrollIntoView` rather than a bare hash link because a
  // Radix dialog holds focus while open and returns it on close, which fights the browser's own
  // hash jump.
  const jump = (id: string) => {
    setOpen(false)
    document.getElementById(id)?.scrollIntoView({ behavior: 'smooth', block: 'start' })
  }

  return (
    <>
      <nav
        aria-label={t('shell.toc.label')}
        className="hidden lg:block w-[220px] flex-shrink-0 sticky top-[100px] h-fit max-h-[calc(100vh-130px)] overflow-y-auto"
      >
        <div className="mb-3 font-mono text-[10px] uppercase tracking-[0.2em]" style={{ color: 'var(--text-muted)' }}>
          {heading}
        </div>
        {headings.map(h => {
          const isActive = active === h.id
          return (
            <a
              key={h.id}
              href={`#${h.id}`}
              className="block py-1.5 pr-2 text-[12px] leading-snug transition-colors"
              style={{
                fontFamily: mono,
                color: isActive ? 'var(--accent-blue)' : 'var(--text-muted)',
                borderLeft: isActive ? '2px solid var(--accent-blue)' : '2px solid transparent',
                paddingLeft: '0.75rem',
                textDecoration: 'none',
              }}
            >
              {h.text}
            </a>
          )
        })}
      </nav>

      <Sheet open={open} onOpenChange={setOpen}>
        <SheetTrigger asChild>
          <button
            type="button"
            aria-label={t('shell.toc.label')}
            className="lg:hidden fixed right-4 z-40 flex items-center gap-2 border border-mathua-border bg-mathua-surface px-3 font-mono text-[11px] text-mathua-secondary min-h-[44px]"
            // Clear the bottom tab bar, which is itself fixed at the safe-area bottom,
            // and clear the assistant trigger above it: both float in this corner, and
            // stacking them here is cheaper than teaching either component about the
            // other. 136px is the trigger's 80px offset plus its 44px height plus a
            // 12px gap.
            style={{ bottom: 'calc(136px + env(safe-area-inset-bottom))' }}
          >
            <Layers size={16} />
            {t('shell.toc.contents')}
          </button>
        </SheetTrigger>
        {/* A card at the bottom-right, not a full-width bar: nine entries used to fill the whole
            width and read as a page takeover. The card is 240px wide, so a strip of the page
            stays visible and every label sits on one line, and each row keeps a 44px tap
            target. */}
        <SheetContent
          side="bottom"
          className="inset-x-auto right-3 bottom-3 w-[min(15rem,calc(100vw-1.5rem))] max-h-[60vh] overflow-y-auto rounded-xl border px-4 pt-3 pb-[calc(0.75rem+env(safe-area-inset-bottom))] shadow-lg"
        >
          <SheetTitle style={{ fontFamily: mono, fontSize: '11px', letterSpacing: '0.15em', textTransform: 'uppercase', color: 'var(--text-muted)' }}>
            {heading}
          </SheetTitle>
          <nav aria-label={t('shell.toc.label')} className="mt-1 flex flex-col">
            {headings.map(h => (
              <button
                key={h.id}
                type="button"
                onClick={() => jump(h.id)}
                className="min-h-[44px] border-b border-mathua-border py-1.5 text-left font-mono text-[12px] leading-tight"
                style={{ color: active === h.id ? 'var(--accent-blue)' : 'var(--text-secondary)' }}
              >
                {h.text}
              </button>
            ))}
          </nav>
        </SheetContent>
      </Sheet>
    </>
  )
}
