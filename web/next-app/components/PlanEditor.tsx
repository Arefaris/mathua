'use client'

import { useCallback, useEffect, useMemo, useState } from 'react'
import { getDestinations, getEstimate, getScores, savePlan, getCurrentPlan, type DestinationStatus, type EstimateRes, type StudyPlan } from '../lib/api'
import { getUserInfo, getGuestId } from '../lib/auth'
import { clampGoal, daysFor, monthLabel, DEFAULT_DAILY_GOAL, GOAL_PRESETS } from '../lib/plan'
import { formatNumber, useT } from '../lib/i18n'

// Learning planner: destination + (daily effort XOR deadline) → estimate.
// Dates render as months (estimates, never promises); quiz eligibility and
// mastery are independent of every number here. Shared by /plan (full) and
// Settings (compact): compact hides the workload detail list, the what-if
// slider, and the plan-delta line, keeping destination/mode/inputs/save.
export default function PlanEditor({ compact = false, initialDestination, onSaved }: { compact?: boolean; initialDestination?: string; onSaved?: () => void }) {
  const t = useT()
  const [dests, setDests] = useState<DestinationStatus[]>([])
  const [destId, setDestId] = useState('')
  const [mode, setMode] = useState<'effort' | 'deadline'>('effort')
  const [dailyGoal, setDailyGoal] = useState(10)
  const [deadlineDays, setDeadlineDays] = useState(90)
  const [restDays, setRestDays] = useState(1)
  const [whatIf, setWhatIf] = useState(10)
  const [est, setEst] = useState<EstimateRes | null>(null)
  const [saved, setSaved] = useState<StudyPlan | null>(null)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    const info = getUserInfo()
    const sid = info?.student_id || getGuestId() || ''
    // Destinations failing is a real, reportable state: it used to be swallowed into an
    // empty array, which rendered a <select> with no options and no message, so the page
    // looked broken rather than unavailable. The score and plan fetches are advisory — a
    // failure there only means we fall back to a default — so they stay non-fatal but are
    // no longer silently indistinguishable from "you have no plan".
    Promise.all([
      getDestinations().then(d => ({ dests: d, err: null as string | null })).catch(e => ({
        dests: [] as DestinationStatus[],
        err: e instanceof Error ? e.message : t('goals.planner.destinationsFailed'),
      })),
      (sid ? getScores(sid) : Promise.resolve(null)).catch(() => null),
      getCurrentPlan().catch(() => ({ plan: null as StudyPlan | null })),
    ]).then(([d, s, p]) => {
      setDests(d.dests)
      // A link from the assistant names a destination; preselecting it is what
      // makes that link land on the thing it named instead of the first option.
      const wanted = d.dests.find(x => x.id === initialDestination)?.id ?? d.dests[0]?.id
      if (wanted) setDestId(wanted)
      if (d.err) setError(d.err)
      const guestGoal = (() => {
        try {
          const v = Number(window.localStorage.getItem('mathua_daily_goal'))
          return v > 0 ? v : null
        } catch {
          return null
        }
      })()
      const g = guestGoal ?? s?.daily_xp_goal ?? DEFAULT_DAILY_GOAL
      setDailyGoal(g)
      setWhatIf(g)
      if (p.plan) setSaved(p.plan)
      setLoading(false)
    })
  }, [initialDestination, t])

  const fetchEst = useCallback(async (id: string, goal: number, deadline: number, rest: number) => {
    if (!id) return
    setError('')
    try {
      const r = await getEstimate(id, {
        daily_goal: clampGoal(goal),
        deadline_days: mode === 'deadline' ? deadline : undefined,
        rest_days: rest,
        diagnostic_min: 10,
      })
      setEst(r)
    } catch (e) {
      setEst(null)
      setError(e instanceof Error ? e.message : t('goals.planner.estimateFailed'))
    }
  }, [mode, t])

  useEffect(() => {
    if (destId) fetchEst(destId, dailyGoal, deadlineDays, restDays)
  }, [destId, dailyGoal, deadlineDays, restDays, fetchEst])

  const whatIfDays = useMemo(
    () => (est ? daysFor(est.estimate.xp_remaining, whatIf, restDays) : -1),
    [est, whatIf, restDays],
  )

  async function handleSave() {
    if (!destId) return
    setSaving(true)
    try {
      const p = await savePlan({ destination: destId, daily_goal: clampGoal(dailyGoal), deadline_days: mode === 'deadline' ? deadlineDays : 0, rest_days: restDays })
      setSaved(p)
      // The caller owns what to say after a save. This component cannot: the
      // next step is a link to another route, and Settings embeds the editor
      // compact for editing a plan that already exists, where that advice is
      // noise. A callback keeps one copy of the sentence in the page that can
      // show it.
      onSaved?.()
    } catch (e) {
      setError(e instanceof Error ? e.message : t('goals.planner.saveFailed'))
    } finally {
      setSaving(false)
    }
  }

  if (loading) return <div className="p-6 font-mono text-xs text-mathua-muted">{t('goals.planner.loading')}</div>

  const e = est?.estimate
  return (
    <div>
      <label htmlFor="plan-dest" className="mt-5 block font-mono text-[11px] uppercase tracking-wider text-mathua-muted">{t('goals.planner.destinationLabel')}</label>
      {dests.length === 0 ? (
        <p role="status" className="mt-1 border border-mathua-border bg-mathua-surface px-3 py-2.5 font-mono text-xs text-mathua-secondary">
          {t('goals.planner.noPaths')}
        </p>
      ) : (
        <select id="plan-dest" value={destId} onChange={ev => setDestId(ev.target.value)} className="mt-1 w-full bg-mathua-surface border border-mathua-border px-3 py-2.5 font-mono text-sm text-mathua-primary min-h-[44px]">
          {dests.map(d => <option key={d.id} value={d.id}>{t('goals.planner.destinationOption', { name: d.name, percent: Math.round(d.pct * 100) })}</option>)}
        </select>
      )}

      <div className="mt-4 flex gap-2" role="group" aria-label={t('goals.planner.direction')}>
        {(['effort', 'deadline'] as const).map(m => (
          <button key={m} type="button" onClick={() => setMode(m)} aria-pressed={mode === m} className={`flex-1 border px-4 py-2.5 font-mono text-xs min-h-[44px] transition-colors ${mode === m ? 'border-mathua-blue text-mathua-blue' : 'border-mathua-border text-mathua-secondary hover:border-mathua-blue'}`}>
            {m === 'effort' ? t('goals.planner.byEffort') : t('goals.planner.byDeadline')}
          </button>
        ))}
      </div>

      {mode === 'effort' ? (
        <div className="mt-3">
          <label htmlFor="plan-goal" className="font-mono text-[11px] uppercase tracking-wider text-mathua-muted">{t('goals.planner.dailyTarget')}</label>
          <div className="mt-1 flex gap-2">
            {GOAL_PRESETS.map(p => (
              <button key={p} type="button" onClick={() => setDailyGoal(p)} aria-pressed={dailyGoal === p} className={`flex-1 border px-2 py-2 font-mono text-xs min-h-[40px] ${dailyGoal === p ? 'border-mathua-blue text-mathua-blue' : 'border-mathua-border text-mathua-secondary'}`}>{p}</button>
            ))}
            <input id="plan-goal" value={dailyGoal} onChange={ev => setDailyGoal(clampGoal(Number(ev.target.value) || 0))} inputMode="numeric" className="w-20 bg-mathua-surface border border-mathua-border px-2 py-2 font-mono text-xs text-mathua-primary text-center" aria-label={t('goals.planner.customTarget')} />
          </div>
        </div>
      ) : (
        <div className="mt-3">
          <label htmlFor="plan-deadline" className="font-mono text-[11px] uppercase tracking-wider text-mathua-muted">{t('goals.planner.deadlineLabel')}</label>
          <div className="mt-1 flex gap-2">
            {[30, 90, 180, 365].map(p => (
              <button key={p} type="button" onClick={() => setDeadlineDays(p)} aria-pressed={deadlineDays === p} className={`flex-1 border px-2 py-2 font-mono text-xs min-h-[40px] ${deadlineDays === p ? 'border-mathua-blue text-mathua-blue' : 'border-mathua-border text-mathua-secondary'}`}>{p >= 365 ? t('goals.planner.presetYear') : p >= 30 ? t('goals.planner.presetMonths', { count: Math.round(p / 30) }) : `${p}d`}</button>
            ))}
            <input id="plan-deadline" value={deadlineDays} onChange={ev => setDeadlineDays(Math.max(1, Number(ev.target.value) || 1))} inputMode="numeric" className="w-20 bg-mathua-surface border border-mathua-border px-2 py-2 font-mono text-xs text-mathua-primary text-center" aria-label={t('goals.planner.customDeadline')} />
          </div>
        </div>
      )}

      <div className="mt-3">
        <label htmlFor="plan-rest" className="font-mono text-[11px] uppercase tracking-wider text-mathua-muted">{t('goals.planner.restDays', { count: restDays })}</label>
        <input id="plan-rest" type="range" min={0} max={3} value={restDays} onChange={ev => setRestDays(Number(ev.target.value))} className="mt-1 w-full" />
      </div>

      {error && <p className="mt-3 font-mono text-xs text-red-400">{error}</p>}

      {e && (
        <div className="mt-5 border border-mathua-border bg-mathua-surface p-5">
          <div className="flex items-baseline justify-between gap-2 flex-wrap">
            <p className="font-mono text-[11px] uppercase tracking-wider text-mathua-muted">{t('goals.planner.remaining')}</p>
            <p className="font-mono text-sm text-mathua-primary">{t('goals.planner.xpApprox', { xp: formatNumber(Math.round(e.xp_remaining), t.locale) })}</p>
          </div>
          {mode === 'effort' ? (
            <p className="mt-3 font-serif text-2xl text-mathua-primary">
              {e.days > 0 ? t('goals.planner.daysApprox', { count: Math.ceil(e.days) }) : t('goals.planner.daysUnknown')}
              {e.finish_date && <span className="text-mathua-secondary text-lg"> · {monthLabel(e.finish_date)}</span>}
            </p>
          ) : (
            <div className="mt-3">
              <p className="font-serif text-2xl text-mathua-primary">{t('goals.planner.perDayNeeded', { value: Math.ceil(e.required_per_day) })}</p>
              {!e.feasible && <p className="mt-1 font-mono text-xs text-red-400">{t('goals.planner.infeasible')}</p>}
            </div>
          )}
          <p className="mt-1 font-mono text-[11px] text-mathua-muted">
            {est?.pace.source === 'measured'
              ? t('goals.planner.paceMeasured', { rate: est.pace.rate.toFixed(0), active: est.pace.active_days, days: est.pace.trailing_days })
              : t('goals.planner.paceNone')}
          </p>
          {!compact && (
            <dl className="mt-4 space-y-1 font-mono text-xs text-mathua-secondary">
              <div className="flex justify-between"><dt>{t('goals.planner.topicsRemaining')}</dt><dd className="text-mathua-primary">{t('goals.planner.topicsOf', { remaining: e.remaining.length, total: e.total })}</dd></div>
              <div className="flex justify-between"><dt>{t('goals.planner.focusedTime')}</dt><dd className="text-mathua-primary">{t('goals.planner.hoursMinutes', { hours: Math.round(e.time_min_remaining / 60), minutes: Math.round(e.time_min_remaining % 60) })}</dd></div>
              <div className="flex justify-between"><dt>{t('goals.planner.reviewsDue')}</dt><dd className="text-mathua-primary">{e.reviews_due}</dd></div>
              <div className="flex justify-between"><dt>{t('goals.planner.checksAhead')}</dt><dd className="text-mathua-primary">{e.quizzes_ahead}</dd></div>
              <div className="flex justify-between"><dt>{t('goals.planner.diagnosticTime')}</dt><dd className="text-mathua-primary">{t('goals.planner.minutes', { count: e.diagnostic_min })}</dd></div>
              <div className="flex justify-between"><dt>{t('goals.planner.coverage')}</dt><dd className="text-mathua-primary">{t('goals.planner.probes', { count: est?.probes.length ?? 0 })}</dd></div>
            </dl>
          )}
          {est?.fresh_start ? (
            <p className="mt-3 font-mono text-xs text-mathua-secondary">
              {t('goals.planner.freshStart')}
            </p>
          ) : (
            !compact && est?.plan_delta_days !== undefined && est.plan_delta_days !== null && (
              <p className="mt-3 font-mono text-xs text-mathua-blue">
                {t(est.plan_delta_days >= 0 ? 'goals.planner.ahead' : 'goals.planner.behind', { count: Number(Math.abs(est.plan_delta_days).toFixed(0)) })}
              </p>
            )
          )}
          {!compact && (
            <div className="mt-4">
              <label htmlFor="plan-whatif" className="font-mono text-[11px] uppercase tracking-wider text-mathua-muted">{t('goals.planner.whatIf', { value: whatIf })} {whatIfDays > 0 && <span className="text-mathua-primary">{t('goals.planner.whatIfDays', { count: Math.ceil(whatIfDays) })}</span>}</label>
              <input id="plan-whatif" type="range" min={2} max={60} step={1} value={whatIf} onChange={ev => setWhatIf(Number(ev.target.value))} className="mt-1 w-full" />
            </div>
          )}
        </div>
      )}

      {/* The save button used to live inside the `{e && ...}` block above, so a failed or
          missing estimate removed the only way to save a plan from the page. Saving records
          a destination and a pace; it does not depend on the workload estimate rendering, and
          it must not be unreachable because a projection failed. */}
      <div className={e ? 'mt-4' : 'mt-5'}>
        <button type="button" onClick={handleSave} disabled={saving || !destId} className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint px-6 py-2.5 font-mono text-xs min-h-[44px] disabled:opacity-50">
          {saving ? t('common.action.saving') : saved ? t('goals.planner.update') : t('goals.planner.create')}
        </button>
      </div>

      {e ? (
        <p className="mt-3 font-mono text-[11px] text-mathua-muted">{t('goals.planner.illustrative')}</p>
      ) : (
        <p className="mt-3 font-mono text-[11px] text-mathua-muted">
          {t('goals.planner.noEstimate')}
        </p>
      )}
    </div>
  )
}
