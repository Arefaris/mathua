'use client'

import { useEffect, useMemo, useRef, useState } from 'react'
import { PlayIcon } from '@animateicons/react/lucide'
import Fuse from 'fuse.js'
import type { GraphStatus } from '../../lib/graphStatus'
import { useT } from '../../lib/i18n'
import { domainName } from '../../lib/domainName'
import { STATUS_COLORS, domainColor, statusLabel, type GraphConcept } from './shared'

// The four controls laid over the graph canvas: search, the flow toggle, and the keyboard list.
// They were inside ConceptGraphFlow.tsx; they moved out so the translation did not push a
// 1200-line file further past the 700-line limit.

interface SearchResult {
  id: string
  label: string
  domain: string
}

export function SearchOverlay({
  concepts,
  onSelect,
}: {
  concepts: GraphConcept[]
  onSelect: (id: string) => void
}) {
  const t = useT()
  const [query, setQuery] = useState('')
  const [open, setOpen] = useState(false)
  const boxRef = useRef<HTMLDivElement>(null)

  const fuse = useMemo(
    () =>
      open
        ? new Fuse(concepts, {
            keys: [
              { name: 'label', weight: 2 },
              { name: 'domain', weight: 1 },
            ],
            threshold: 0.4,
            minMatchCharLength: 2,
          })
        : null,
    [concepts, open]
  )

  const results = useMemo<SearchResult[]>(() => {
    if (!fuse || !query.trim()) return []
    return fuse
      .search(query.trim())
      .slice(0, 8)
      .map(r => ({ id: r.item.id, label: r.item.label, domain: r.item.domain }))
  }, [query, fuse])

  useEffect(() => {
    function handleClick(e: MouseEvent) {
      if (boxRef.current && !boxRef.current.contains(e.target as globalThis.Node)) setOpen(false)
    }
    document.addEventListener('mousedown', handleClick)
    return () => document.removeEventListener('mousedown', handleClick)
  }, [])

  return (
    <div
      ref={boxRef}
      className="absolute top-2.5 left-[62px] sm:left-16 z-[5] w-[min(calc(100%-76px),240px)] sm:w-[240px] font-mono"
    >
      <input
        value={query}
        onChange={e => {
          setQuery(e.target.value)
          setOpen(true)
        }}
        onFocus={() => setOpen(true)}
        aria-label={t('explore.graph.flow.searchLabel')}
        onKeyDown={e => {
          if (e.key === 'Enter' && results.length > 0) {
            onSelect(results[0].id)
            setOpen(false)
          }
          if (e.key === 'Escape') setOpen(false)
        }}
        placeholder={t('explore.graph.flow.searchPlaceholder')}
        className="w-full min-h-[36px] h-9 px-2.5 text-[11px] font-mono text-mathua-primary bg-mathua-surface-elevated border-[0.5px] border-mathua-border-strong outline-none"
      />
      {open && results.length > 0 && (
        <div className="mt-1 bg-mathua-surface-elevated border-[0.5px] border-mathua-border overflow-hidden max-h-[min(50dvh,240px)] overflow-y-auto">
          {results.map(r => (
            <button
              key={r.id}
              type="button"
              onClick={() => {
                onSelect(r.id)
                setOpen(false)
                setQuery('')
              }}
              className="flex items-center gap-1.5 w-full text-left px-2.5 py-2 min-h-[36px] text-[11px] font-mono text-mathua-secondary bg-transparent border-none border-b-[0.5px] border-mathua-border last:border-b-0 cursor-pointer hover:bg-mathua-code"
            >
              <span className="w-[3px] h-3.5 shrink-0" style={{ background: domainColor(r.domain) }} />
              <span className="truncate flex-1 min-w-0">{r.label}</span>
            </button>
          ))}
        </div>
      )}
    </div>
  )
}

export function AmbientToggle({
  enabled,
  onToggle,
}: {
  enabled: boolean
  onToggle: () => void
}) {
  const t = useT()
  return (
    <button
      type="button"
      onClick={onToggle}
      title={t('explore.graph.flow.flowToggle')}
      aria-label={t('explore.graph.flow.flowToggle')}
      style={{
        position: 'absolute',
        top: 10,
        right: 10,
        zIndex: 5,
        height: 30,
        padding: '0 12px',
        fontSize: 10,
        textTransform: 'uppercase',
        letterSpacing: '0.05em',
        fontFamily: "var(--font-jetbrains-mono), 'JetBrains Mono', monospace",
        color: enabled ? '#fff' : 'var(--text-muted)',
        background: enabled ? 'var(--accent-blue)' : 'var(--surface-elevated)',
        border: '0.5px solid var(--border-strong)',
        cursor: 'pointer',
      }}
    >
      {/* An icon rather than the play glyph (U+23F5) this replaced, which sat at an
          unpredictable baseline next to the word below. */}
      <PlayIcon
        size={12}
        duration={1}
        color="currentColor"
        className="inline-block align-[-1px]"
        aria-hidden="true"
      />
      <span className="ml-1">{t('explore.graph.flow.flowButton')}</span>
    </button>
  )
}

export function ListToggleButton({
  open,
  onToggle,
}: {
  open: boolean
  onToggle: () => void
}) {
  const t = useT()
  return (
    <button
      type="button"
      onClick={onToggle}
      title={t('explore.graph.flow.listTitle')}
      aria-expanded={open}
      aria-label={open ? t('explore.graph.flow.listClose') : t('explore.graph.flow.listOpen')}
      className={`absolute top-2.5 left-2.5 z-[6] min-h-[36px] h-9 px-2.5 text-[10px] uppercase tracking-[0.05em] font-mono border-[0.5px] border-mathua-border-strong cursor-pointer transition-colors ${open ? 'bg-mathua-teal text-white' : 'bg-mathua-surface-elevated text-mathua-muted hover:text-mathua-primary'}`}
    >
      {t('explore.graph.flow.listButton')}
    </button>
  )
}

export function ListView({
  concepts,
  conceptStatuses,
  selectedId,
  onSelect,
  onClose,
}: {
  concepts: GraphConcept[]
  conceptStatuses?: Record<string, GraphStatus>
  selectedId: string | null
  onSelect: (id: string) => void
  onClose: () => void
}) {
  const t = useT()
  const [query, setQuery] = useState('')
  const inputRef = useRef<HTMLInputElement>(null)

  useEffect(() => {
    inputRef.current?.focus()
  }, [])

  const filtered = useMemo(() => {
    const q = query.trim().toLowerCase()
    if (!q) return concepts
    return concepts.filter(
      c =>
        c.label.toLowerCase().includes(q) ||
        c.id.toLowerCase().includes(q) ||
        c.domain.toLowerCase().includes(q)
    )
  }, [concepts, query])

  const groups = useMemo(() => {
    const order: string[] = []
    const byDomain = new Map<string, GraphConcept[]>()
    for (const c of filtered) {
      if (!byDomain.has(c.domain)) {
        byDomain.set(c.domain, [])
        order.push(c.domain)
      }
      byDomain.get(c.domain)!.push(c)
    }
    return order.map(d => ({ domain: d, items: byDomain.get(d)! }))
  }, [filtered])

  return (
    <div
      role="dialog"
      aria-label={t('explore.graph.flow.listDialog')}
      style={{
        position: 'absolute',
        inset: '48px 8px 8px 8px',
        zIndex: 6,
        display: 'flex',
        flexDirection: 'column',
        background: 'var(--surface-elevated)',
        border: '0.5px solid var(--border)',
        overflow: 'hidden',
      }}
    >
      <div style={{ padding: 8 }}>
        <input
          ref={inputRef}
          value={query}
          onChange={e => setQuery(e.target.value)}
          onKeyDown={e => {
            if (e.key === 'Escape') {
              e.preventDefault()
              onClose()
            }
          }}
          placeholder={t('explore.graph.flow.filterPlaceholder')}
          aria-label={t('explore.graph.flow.filterLabel')}
          style={{
            width: '100%',
            height: 28,
            padding: '0 8px',
            fontSize: 11,
            fontFamily: 'inherit',
            color: 'var(--text-primary)',
            background: 'var(--surface-elevated)',
            border: '0.5px solid var(--border-strong)',
            outline: 'none',
          }}
        />
      </div>
      <div style={{ flex: 1, overflowY: 'auto', padding: '0 8px 8px' }}>
        {groups.map(g => (
          <div key={g.domain}>
            <div
              style={{
                margin: '8px 0 4px',
                fontSize: 9,
                textTransform: 'uppercase',
                letterSpacing: '0.08em',
                color: 'var(--text-muted)',
                fontFamily: "var(--font-jetbrains-mono), 'JetBrains Mono', monospace",
              }}
            >
              {domainName(g.domain)} · {g.items.length}
            </div>
            {g.items.map(c => {
              const status = conceptStatuses?.[c.id]
              return (
                <button
                  key={c.id}
                  type="button"
                  onClick={() => {
                    onSelect(c.id)
                    onClose()
                  }}
                  aria-current={c.id === selectedId ? 'true' : undefined}
                  style={{
                    display: 'flex',
                    alignItems: 'center',
                    gap: 8,
                    width: '100%',
                    textAlign: 'left',
                    padding: '5px 8px',
                    fontSize: 11,
                    fontFamily: 'inherit',
                    color:
                      c.id === selectedId
                        ? 'var(--text-primary)'
                        : 'var(--text-secondary)',
                    background:
                      c.id === selectedId
                        ? 'var(--surface-elevated)'
                        : 'transparent',
                    border: 'none',
                    borderBottom: '0.5px solid var(--border)',
                    cursor: 'pointer',
                  }}
                >
                  <span
                    style={{
                      width: 7,
                      height: 7,
                      borderRadius: '50%',
                      flexShrink: 0,
                      background:
                        status != null ? STATUS_COLORS[status] : STATUS_COLORS.unseen,
                    }}
                    aria-hidden
                  />
                  <span
                    style={{
                      whiteSpace: 'nowrap',
                      overflow: 'hidden',
                      textOverflow: 'ellipsis',
                      flex: 1,
                    }}
                  >
                    {c.label}
                  </span>
                  {status && (
                    <span style={{ color: 'var(--text-muted)', fontSize: 9 }}>
                      {statusLabel(t, status)}
                    </span>
                  )}
                </button>
              )
            })}
          </div>
        ))}
        {groups.length === 0 && (
          <div
            style={{
              padding: 12,
              fontSize: 11,
              color: 'var(--text-muted)',
              fontFamily: "var(--font-jetbrains-mono), 'JetBrains Mono', monospace",
            }}
          >
            {t('explore.graph.flow.noMatch')}
          </div>
        )}
      </div>
    </div>
  )
}

