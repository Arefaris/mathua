'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import Header from '../../components/Header'
import BottomTabs from '../../components/BottomTabs'
import SectionHeader from '../../components/SectionHeader'
import DangerZone from '../../components/DangerZone'
import DeleteAccount from '../../components/DeleteAccount'
import DailyGoalControl from '../../components/DailyGoalControl'
import PlanEditor from '../../components/PlanEditor'
import Footer from '../../components/Footer'
import { getSettings, getScores, updateSettings, type UserSettings } from '../../lib/api'
import { isLoggedIn, getUserInfo } from '../../lib/auth'
import { Switch } from '../../components/ui/switch'
import Loading from '../../components/Loading'
import LanguageSwitcher from '../../components/LanguageSwitcher'
import RichText from '../../components/RichText'
import ProfileSection from '../../components/settings/ProfileSection'
import PasswordSection from '../../components/settings/PasswordSection'
import ConnectedAccounts from '../../components/settings/ConnectedAccounts'
import ShareSection from '../../components/settings/ShareSection'
import { formatDate, useLocale, useT } from '../../lib/i18n'

export default function SettingsPage() {
  const t = useT()
  const locale = useLocale()
  const { push } = useRouter()
  const [settings, setSettings] = useState<UserSettings>({})
  const [dailyGoal, setDailyGoal] = useState(10)
  const [saved, setSaved] = useState(false)
  const [loading, setLoading] = useState(true)
  const [restricted, setRestricted] = useState(false)
  const [saveError, setSaveError] = useState('')

  useEffect(() => {
    if (!isLoggedIn()) {
      setRestricted(true)
      setLoading(false)
      return
    }
    getSettings().then(s => {
      setSettings(s)
      setLoading(false)
    }).catch((e) => {
      console.error('getSettings failed:', e)
      setLoading(false)
    })
    const info = getUserInfo()
    if (info?.student_id) {
      getScores(info.student_id).then(s => setDailyGoal(s.daily_xp_goal)).catch(() => {})
    }
  }, [])

  const handleCheckedChange = async (checked: boolean) => {
    const next = { ...settings, show_timer: checked }
    setSettings(next)
    setSaved(false)
    try {
      await updateSettings(next)
      setSaved(true)
      setTimeout(() => setSaved(false), 2000)
    } catch {
      // Keep optimistic state so the toggle still works when the API is
      // unreachable (e.g. static preview). Persist will apply on next save.
      console.error('updateSettings failed')
      setSaveError(t('account.settings.saveFailed'))
      setTimeout(() => setSaveError(''), 3000)
    }
  }

  const handlePause = async (days: number | null) => {
    const next: UserSettings = { ...settings }
    if (days === null) {
      next.pause_until = null
    } else {
      const d = new Date()
      d.setUTCDate(d.getUTCDate() + days)
      next.pause_until = d.toISOString().slice(0, 10)
    }
    setSettings(next)
    setSaved(false)
    try {
      await updateSettings(next)
      setSaved(true)
      setTimeout(() => setSaved(false), 2000)
    } catch {
      console.error('updateSettings failed')
    }
  }

  const paused = !!settings.pause_until
  const pausedLabel = (() => {
    if (!settings.pause_until) return ''
    const d = new Date(settings.pause_until + 'T00:00:00Z')
    return isNaN(d.getTime()) ? settings.pause_until : formatDate(d, locale, { year: 'numeric', month: 'short', day: 'numeric' })
  })()

  const planLink = (
    <Link href="/plan" className="text-mathua-blue hover:text-mathua-blue-hover">{t('account.settings.planLink')}</Link>
  )

  if (restricted) {
    return (
      <>
        <Header />
        <div className="max-w-container mx-auto px-4 sm:px-6 py-20 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-20 text-center overflow-x-hidden min-w-0">
          <p className="font-mono text-xs text-mathua-secondary mb-4">{t('account.settings.needsAccount')}</p>
          <Link href="/login" className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint px-6 py-2 font-mono text-xs min-h-[36px] inline-flex items-center justify-center">{t('common.action.signIn')}</Link>
          {/* The language belongs to the browser, not the account, so a guest can change it. */}
          <div className="mt-10 flex flex-col items-center gap-2">
            <span className="font-mono text-sm text-mathua-primary">{t('common.language.label')}</span>
            <LanguageSwitcher />
          </div>
        </div>
        <BottomTabs />
        <Footer />
      </>
    )
  }

  if (loading) {
    return (
      <>
        <Header />
        <div className="max-w-container mx-auto px-4 sm:px-6 pt-[var(--chrome-top)] pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 text-center overflow-x-hidden min-w-0">
          <Loading label={t('account.settings.loading')} />
        </div>
        <BottomTabs />
        <Footer />
      </>
    )
  }

  return (
    <>
      <Header />
      <div className="max-w-container mx-auto px-4 sm:px-6 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 overflow-x-hidden min-w-0">
        <section className="pt-8 min-w-0 overflow-hidden">
          <span className="flex justify-between items-center mb-4">
            <button type="button" onClick={() => { if (window.history.length > 1) window.history.back(); else push('/profile') }} className="text-mathua-secondary text-sm hover:text-mathua-primary">
              ← {t('common.action.back')}
            </button>
          </span>

          <div className="max-w-lg mx-auto mt-8 sm:mt-12 w-full max-w-full min-w-0 px-2 sm:px-0">
            <SectionHeader label={t('account.settings.label')} title={t('account.settings.title')} />
            <p className="text-mathua-secondary text-sm text-center mt-2 mb-8">
              {t('account.settings.intro')}
            </p>

            <div className="bg-mathua-surface border border-mathua-border rounded-none p-4 sm:p-6 space-y-6 w-full max-w-full min-w-0 overflow-hidden">
              <div className="min-w-0">
                <span className="font-mono text-sm text-mathua-primary">{t('common.language.label')}</span>
                <p className="text-mathua-muted text-xs mt-1">{t('common.language.hint')}</p>
                <div className="mt-3">
                  <LanguageSwitcher />
                </div>
              </div>

              <div className="border-t border-mathua-border pt-6 flex items-center justify-between gap-4 min-w-0">
                <div className="min-w-0 flex-1">
                  <span className="font-mono text-sm text-mathua-primary">{t('account.settings.timerTitle')}</span>
                  <p className="text-mathua-muted text-xs mt-1">
                    {t('account.settings.timerHint')}
                  </p>
                </div>
                <Switch
                  checked={!!settings.show_timer}
                  onCheckedChange={handleCheckedChange}
                  aria-label={t('account.settings.timerToggle')}
                />
              </div>

              <div className="border-t border-mathua-border pt-6 min-w-0">
                <span className="font-mono text-sm text-mathua-primary">{t('account.settings.dailyTitle')}</span>
                <p className="text-mathua-muted text-xs mt-1">
                  <RichText text={t('account.settings.dailyHint')} slots={{ plan: planLink }} />
                </p>
                <div className="mt-3">
                  <DailyGoalControl current={dailyGoal} onSaved={setDailyGoal} bare />
                </div>
              </div>

              <div className="border-t border-mathua-border pt-6 min-w-0">
                <span className="font-mono text-sm text-mathua-primary">{t('account.settings.planTitle')}</span>
                <p className="text-mathua-muted text-xs mt-1">
                  <RichText text={t('account.settings.planHint')} slots={{ plan: planLink }} />
                </p>
                <div className="mt-2">
                  <PlanEditor compact />
                </div>
              </div>

              <div className="border-t border-mathua-border pt-6 min-w-0">
                <span className="font-mono text-sm text-mathua-primary">{t('account.settings.pauseTitle')}</span>
                <p className="text-mathua-muted text-xs mt-1">
                  {paused
                    ? t('account.settings.pausedHint', { date: pausedLabel })
                    : t('account.settings.pauseHint')}
                </p>
                <div className="flex flex-wrap gap-2 mt-3">
                  {[30, 60, 90].map(d => (
                    <button
                      type="button"
                      key={d}
                      onClick={() => handlePause(d)}
                      className="border border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue rounded-none px-4 h-9 text-xs font-mono disabled:opacity-40"
                    >
                      {t('account.settings.pauseFor', { count: d })}
                    </button>
                  ))}
                  {paused && (
                    <button
                      type="button"
                      onClick={() => handlePause(null)}
                      className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint rounded-none px-4 h-9 text-xs font-mono"
                    >
                      {t('account.settings.resumeNow')}
                    </button>
                  )}
                </div>
              </div>

              <ProfileSection settings={settings} onSettings={setSettings} />
              <PasswordSection />
              <ConnectedAccounts />
              <ShareSection />
            </div>

            {saved && (
              <p className="text-mathua-green text-xs text-center mt-4 font-mono">
                {t('account.settings.saved')}
              </p>
            )}
            {saveError && (
              <p className="text-mathua-red text-xs text-center mt-4 font-mono">
                {saveError}
              </p>
            )}
          </div>

          <div className="mt-8">
            <DangerZone />
          </div>

          <div className="mt-4">
            <DeleteAccount />
          </div>
        </section>
      </div>

      <Footer />
      <BottomTabs />
    </>
  )
}
