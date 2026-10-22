'use client'

import { useTheme } from '../../hooks/useTheme'
import { useState, useEffect, useCallback } from 'react'
import Link from 'next/link'
import Header from '../../components/Header'
import BottomTabs from '../../components/BottomTabs'
import SectionHeader from '../../components/SectionHeader'
import Footer from '../../components/Footer'
import Avatar from '../../components/Avatar'
import { LeaderboardTableSkeleton } from '../../components/skeletons/LeaderboardSkeleton'
import {
  getLeaderboard,
  getMyLeague,
  type LeaderboardBoard,
  type LeaderboardEntry,
  type LeaderboardView,
  type MyLeague,
} from '../../lib/api'
import { resolveAvatar } from '../../lib/dicebear'
import RichText from '../../components/RichText'
import { formatDate, useLocale, useT, type Locale, type MessageKey } from '../../lib/i18n'
import { levelName, tierName } from '../../lib/accountNames'

type Tab = LeaderboardView | 'league'

const TABS: { id: Tab; label: MessageKey }[] = [
  { id: 'weekly', label: 'account.leaderboard.tabWeekly' },
  { id: 'all_time', label: 'account.leaderboard.tabAllTime' },
  { id: 'league', label: 'account.leaderboard.tabLeague' },
]

const PAGE_SIZE = 25

function computeCountdown(): { days: number; hours: number; minutes: number } {
  const now = new Date()
  const nextMonday = new Date(now)
  nextMonday.setUTCDate(now.getUTCDate() + ((7 - now.getUTCDay() + 1) % 7 || 7))
  nextMonday.setUTCHours(0, 0, 0, 0)
  const diff = nextMonday.getTime() - now.getTime()
  const d = Math.floor(diff / 86400000)
  const h = Math.floor((diff % 86400000) / 3600000)
  const m = Math.floor((diff % 3600000) / 60000)
  return { days: d, hours: h, minutes: m }
}

// The server decides the avatar shape; the client only has to turn it into a
// URL. A generated avatar is seeded from the public id, never from the account.
function avatarURL(row: { id: string; avatar?: { kind: string; url?: string; style?: string; version?: number } }): string | undefined {
  const av = row.avatar
  if (!av) return undefined
  if (av.kind === 'custom' && av.url) {
    return av.version ? `${av.url}?v=${av.version}` : av.url
  }
  return resolveAvatar(
    { avatar_url: '' },
    { avatar_dicebear: av.style ? { style: av.style, seed: row.id } : null },
  ).url
}

// The server resolves a display name, but a blank cell is a silent failure and
// costs one line to prevent: fall back to the public id, which is never empty.
function nameOf(row: { name: string; id: string }): string {
  return row.name.trim() || row.id
}

function shortWeek(iso: string, locale: Locale): string {
  return formatDate(iso, locale, { month: 'short', day: 'numeric', timeZone: 'UTC' })
}

function Podium({ entries }: { entries: LeaderboardEntry[] }) {
  if (entries.length < 3) return null
  return (
    <ol className="grid grid-cols-1 sm:grid-cols-3 gap-3 mb-6 list-none p-0">
      {entries.slice(0, 3).map((row, i) => (
        <li
          key={`podium-${row.id}`}
          className={`border border-mathua-border bg-mathua-surface px-4 py-3 flex items-center gap-3 min-w-0 ${i === 0 ? 'sm:order-2' : i === 1 ? 'sm:order-1' : 'sm:order-3'}`}
        >
          <span className="font-mono text-lg text-mathua-blue w-5 shrink-0">{row.rank}</span>
          <Avatar seed={row.id} name={nameOf(row)} size={30} url={avatarURL(row)} className="shrink-0" />
          <span className="min-w-0 flex-1">
            <span className="block text-sm text-mathua-primary truncate" title={nameOf(row)}>{nameOf(row)}</span>
            <span className="block font-mono text-[10px] uppercase tracking-wider text-mathua-muted">
              {row.xp} XP
            </span>
          </span>
        </li>
      ))}
    </ol>
  )
}

function PositionStrip({
  me,
  view,
}: {
  me: NonNullable<LeaderboardBoard['me']>
  view: LeaderboardView
}) {
  const t = useT()
  if (me.guest) {
    return (
      <p className="mb-4 text-sm text-mathua-secondary">
        <RichText
          text={t('account.leaderboard.guest')}
          slots={{ signIn: <Link href="/login" className="text-mathua-blue hover:underline">{t('account.leaderboard.signIn')}</Link> }}
        />
      </p>
    )
  }
  if (!me.eligible || me.rank === 0) {
    return (
      <p className="mb-4 text-sm text-mathua-secondary">
        {t(view === 'weekly' ? 'account.leaderboard.noPlaceWeekly' : 'account.leaderboard.noPlaceAllTime')}
      </p>
    )
  }
  if (me.in_page) return null
  return (
    <p className="mb-4 font-mono text-xs text-mathua-primary border border-mathua-border bg-mathua-blue-faint px-3 py-2">
      <RichText
        text={t('account.leaderboard.position', { xp: me.xp })}
        slots={{ rank: <span className="text-mathua-blue">#{me.rank}</span> }}
      />
    </p>
  )
}

function PageControls({
  board,
  onPage,
  busy,
}: {
  board: LeaderboardBoard
  onPage: (page: number) => void
  busy: boolean
}) {
  const t = useT()
  const start = board.total === 0 ? 0 : (board.page - 1) * board.page_size + 1
  const end = start === 0 ? 0 : start + board.entries.length - 1
  return (
    <div className="flex flex-wrap items-center justify-between gap-3 mt-4">
      <p className="font-mono text-[11px] text-mathua-muted" aria-live="polite">
        {board.total === 0 ? t('account.leaderboard.noLearners') : t('account.leaderboard.range', { start, end, total: board.total })}
      </p>
      <div className="flex items-center gap-2">
        <button
          type="button"
          onClick={() => onPage(board.page - 1)}
          disabled={busy || board.page <= 1}
          className="font-mono text-xs px-3 py-1.5 border border-mathua-border text-mathua-primary disabled:opacity-40 disabled:cursor-not-allowed hover:border-mathua-blue"
        >
          {t('account.leaderboard.previous')}
        </button>
        <span className="font-mono text-[11px] text-mathua-muted">
          {t('account.leaderboard.page', { page: board.page, total: board.total_pages })}
        </span>
        <button
          type="button"
          onClick={() => onPage(board.page + 1)}
          disabled={busy || board.page >= board.total_pages}
          className="font-mono text-xs px-3 py-1.5 border border-mathua-border text-mathua-primary disabled:opacity-40 disabled:cursor-not-allowed hover:border-mathua-blue"
        >
          {t('account.leaderboard.next')}
        </button>
      </div>
    </div>
  )
}

function BoardTable({ board, view }: { board: LeaderboardBoard; view: LeaderboardView }) {
  const t = useT()
  const scoreLabel = t(view === 'weekly' ? 'account.leaderboard.scoreWeek' : 'account.leaderboard.scoreLifetime')
  return (
    <div className="w-[calc(100%+2rem)] -mx-4 px-4 sm:w-full sm:mx-0 sm:px-0 overflow-x-auto overscroll-x-contain">
      <table className="w-full min-w-[320px] border-collapse bg-mathua-surface border border-mathua-border">
        <caption className="sr-only">
          {t(view === 'weekly' ? 'account.leaderboard.captionWeekly' : 'account.leaderboard.captionAllTime')}
        </caption>
        <thead>
          <tr>
            <th scope="col" className="font-mono text-[11px] uppercase text-mathua-muted text-left p-2.5 sm:p-[14px_20px] bg-mathua-surface-elevated border-b border-mathua-border w-10 sm:w-12">#</th>
            <th scope="col" className="font-mono text-[11px] uppercase text-mathua-muted text-left p-2.5 sm:p-[14px_20px] bg-mathua-surface-elevated border-b border-mathua-border">{t('account.leaderboard.colLearner')}</th>
            <th scope="col" className="font-mono text-[11px] uppercase text-mathua-muted text-right p-2.5 sm:p-[14px_20px] bg-mathua-surface-elevated border-b border-mathua-border">{scoreLabel}</th>
            <th scope="col" className="font-mono text-[11px] uppercase text-mathua-muted text-right p-2.5 sm:p-[14px_20px] bg-mathua-surface-elevated border-b border-mathua-border hidden md:table-cell">{t('account.leaderboard.colLevel')}</th>
          </tr>
        </thead>
        <tbody>
          {board.entries.map((row) => (
            <tr
              key={row.id}
              className={`border-b border-mathua-border last:border-b-0 ${row.is_me ? 'bg-mathua-blue-faint' : ''}`}
            >
              <td className="p-2.5 sm:p-[14px_20px] font-mono text-sm text-mathua-muted">{row.rank}</td>
              <td className="p-2.5 sm:p-[14px_20px] text-mathua-primary text-sm font-medium max-w-[160px] sm:max-w-none">
                <span className="flex items-center gap-2 min-w-0">
                  <Avatar seed={row.id} name={nameOf(row)} size={28} url={avatarURL(row)} className="shrink-0" />
                  <span className="block truncate" title={nameOf(row)}>{nameOf(row)}</span>
                  {row.is_me && (
                    <span className="font-mono text-[10px] uppercase text-mathua-blue shrink-0">{t('account.leaderboard.you')}</span>
                  )}
                </span>
              </td>
              <td className="p-2.5 sm:p-[14px_20px] font-mono text-sm text-mathua-green text-right whitespace-nowrap">{row.xp}</td>
              <td className="p-2.5 sm:p-[14px_20px] text-mathua-secondary text-sm text-right hidden md:table-cell truncate max-w-[130px]">{levelName(row.level)}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

export default function LeaderboardPage() {
  const t = useT()
  const locale = useLocale()
  const { mounted } = useTheme()
  const [countdown, setCountdown] = useState(computeCountdown)
  const [tab, setTab] = useState<Tab>('weekly')
  const [page, setPage] = useState(1)
  const [reload, setReload] = useState(0)
  const [board, setBoard] = useState<LeaderboardBoard | null>(null)
  const [league, setLeague] = useState<MyLeague | null>(null)
  const [loading, setLoading] = useState(true)
  const [failed, setFailed] = useState(false)

  useEffect(() => {
    const update = () => setCountdown(computeCountdown())
    update()
    const interval = setInterval(update, 60000)
    return () => clearInterval(interval)
  }, [])

  useEffect(() => {
    let live = true
    setLoading(true)
    setFailed(false)
    if (tab === 'league') {
      getMyLeague()
        .then((data) => { if (live) { setLeague(data); setLoading(false) } })
        .catch(() => { if (live) { setFailed(true); setLoading(false) } })
    } else {
      getLeaderboard(tab, page, PAGE_SIZE)
        .then((data) => { if (live) { setBoard(data); setLoading(false) } })
        .catch(() => { if (live) { setFailed(true); setLoading(false) } })
    }
    return () => { live = false }
  }, [tab, page, reload])

  const selectTab = useCallback((next: Tab) => {
    setTab(next)
    setPage(1)
  }, [])

  const retry = useCallback(() => setReload((n) => n + 1), [])

  if (!mounted) return <div style={{ background: 'var(--bg)', minHeight: '100vh' }} />

  return (
    <>
      <Header />
      <div className="max-w-container mx-auto px-4 sm:px-6 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 overflow-x-hidden min-w-0">
        <section className="pt-8 min-w-0 overflow-hidden">
          <span className="flex mb-4">
            <button type="button" onClick={() => { if (window.history.length > 1) window.history.back() }} className="text-mathua-secondary text-sm hover:text-mathua-blue">
              {t('account.leaderboard.back')}
            </button>
          </span>

          <SectionHeader label={t('account.leaderboard.label')} title={t('account.leaderboard.title')} />

          <p className="text-mathua-secondary text-sm text-center mt-4 mb-6 px-2 max-w-[560px] mx-auto">
            {t('account.leaderboard.intro')}
          </p>

          {tab === 'weekly' && (
            <div className="flex justify-center mb-6 px-2">
              <div className="bg-mathua-surface border border-mathua-border p-4 text-center w-full max-w-[260px] sm:min-w-[200px] sm:w-auto">
                <span className="font-mono text-[10px] uppercase text-mathua-muted">{t('account.leaderboard.resetsIn')}</span>
                <div className="font-mono text-2xl text-mathua-blue mt-1">{t('account.leaderboard.countdown', countdown)}</div>
              </div>
            </div>
          )}

          <div role="tablist" aria-label={t('account.leaderboard.tabsLabel')} className="flex flex-wrap justify-center gap-1 border-b border-mathua-border mb-6">
            {TABS.map((tabDef) => (
              <button
                key={tabDef.id}
                type="button"
                role="tab"
                id={`tab-${tabDef.id}`}
                aria-selected={tab === tabDef.id}
                aria-controls={`panel-${tabDef.id}`}
                onClick={() => selectTab(tabDef.id)}
                className={`font-mono text-xs uppercase tracking-wider px-3 sm:px-4 py-2 border-b-2 -mb-px whitespace-nowrap ${
                  tab === tabDef.id
                    ? 'border-mathua-blue text-mathua-blue'
                    : 'border-transparent text-mathua-muted hover:text-mathua-primary'
                }`}
              >
                {t(tabDef.label)}
              </button>
            ))}
          </div>
        </section>

        <div
          role="tabpanel"
          id={`panel-${tab}`}
          aria-labelledby={`tab-${tab}`}
          className="mb-12 min-w-0"
        >
          {loading && (
            <div className="w-[calc(100%+2rem)] -mx-4 px-4 sm:w-full sm:mx-0 sm:px-0 overflow-x-auto">
              <span role="status" aria-live="polite" className="sr-only">{t('account.leaderboard.loadingBoard')}</span>
              <table className="w-full min-w-[320px] border-collapse bg-mathua-surface border border-mathua-border">
                <tbody><LeaderboardTableSkeleton rows={8} /></tbody>
              </table>
            </div>
          )}

          {!loading && failed && (
            <p className="text-center text-mathua-muted text-sm py-8">
              {t('account.leaderboard.failed')}{' '}
              <button type="button" onClick={retry} className="text-mathua-blue hover:underline">{t('account.leaderboard.retry')}</button>
            </p>
          )}

          {!loading && !failed && tab !== 'league' && board && (
            <>
              {board.page === 1 && <Podium entries={board.entries} />}
              {board.me && <PositionStrip me={board.me} view={tab} />}
              {board.entries.length === 0 ? (
                <p className="text-center text-mathua-muted text-sm py-8 px-2">
                  {t(tab === 'weekly' ? 'account.leaderboard.emptyWeekly' : 'account.leaderboard.emptyAllTime')}
                  {' '}<Link href="/learn" className="text-mathua-blue hover:underline">{t('account.leaderboard.startLearning')}</Link>
                </p>
              ) : (
                <>
                  <BoardTable board={board} view={tab} />
                  <PageControls board={board} onPage={setPage} busy={loading} />
                </>
              )}
            </>
          )}

          {!loading && !failed && tab === 'league' && (
            league === null ? (
              <p className="text-center text-mathua-muted text-sm py-8 px-2">
                <RichText
                  text={t('account.leaderboard.leagueSignIn')}
                  slots={{ signIn: <Link href="/login" className="text-mathua-blue hover:underline">{t('account.leaderboard.signIn')}</Link> }}
                />
              </p>
            ) : league.tier === '' ? (
              <p className="text-center text-mathua-muted text-sm py-8 px-2">
                {t('account.leaderboard.noLeague')}
              </p>
            ) : (
              <>
                <div className="flex flex-wrap items-baseline justify-between gap-2 mb-4">
                  <h2 className="font-mono text-xs uppercase tracking-wider text-mathua-primary">
                    {t('account.leaderboard.leagueTitle', { tier: tierName(league.tier), number: league.league_no })}
                  </h2>
                  <p className="font-mono text-[11px] text-mathua-muted">
                    {t('account.leaderboard.leagueRange', { start: shortWeek(league.week_start, locale), end: shortWeek(league.week_end, locale), size: t('account.leaderboard.learners', { count: league.size }) })}
                  </p>
                </div>
                <p className="text-mathua-secondary text-sm mb-4">
                  {t('account.leaderboard.leagueRules', { promote: league.promote, demote: league.demote })}
                </p>
                <div className="w-[calc(100%+2rem)] -mx-4 px-4 sm:w-full sm:mx-0 sm:px-0 overflow-x-auto">
                  <table className="w-full min-w-[320px] border-collapse bg-mathua-surface border border-mathua-border">
                    <caption className="sr-only">{t('account.leaderboard.leagueCaption')}</caption>
                    <thead>
                      <tr>
                        <th scope="col" className="font-mono text-[11px] uppercase text-mathua-muted text-left p-2.5 sm:p-[14px_20px] bg-mathua-surface-elevated border-b border-mathua-border w-10 sm:w-12">#</th>
                        <th scope="col" className="font-mono text-[11px] uppercase text-mathua-muted text-left p-2.5 sm:p-[14px_20px] bg-mathua-surface-elevated border-b border-mathua-border">{t('account.leaderboard.colLearner')}</th>
                        <th scope="col" className="font-mono text-[11px] uppercase text-mathua-muted text-right p-2.5 sm:p-[14px_20px] bg-mathua-surface-elevated border-b border-mathua-border">{t('account.leaderboard.leagueXp')}</th>
                      </tr>
                    </thead>
                    <tbody>
                      {league.members.map((m) => (
                        <tr key={m.id} className={`border-b border-mathua-border last:border-b-0 ${m.is_me ? 'bg-mathua-blue-faint' : ''}`}>
                          <td className="p-2.5 sm:p-[14px_20px] font-mono text-sm text-mathua-muted">{m.rank}</td>
                          <td className="p-2.5 sm:p-[14px_20px] text-mathua-primary text-sm font-medium max-w-[160px] sm:max-w-none">
                            <span className="flex items-center gap-2 min-w-0">
                              <Avatar seed={m.id} name={nameOf(m)} size={26} url={avatarURL(m)} className="shrink-0" />
                              <span className="block truncate" title={nameOf(m)}>{nameOf(m)}</span>
                              {m.is_me && <span className="font-mono text-[10px] uppercase text-mathua-blue shrink-0">{t('account.leaderboard.you')}</span>}
                            </span>
                          </td>
                          <td className="p-2.5 sm:p-[14px_20px] font-mono text-sm text-mathua-green text-right whitespace-nowrap">{m.xp}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </>
            )
          )}
        </div>
      </div>

      <Footer />
      <BottomTabs />
    </>
  )
}
