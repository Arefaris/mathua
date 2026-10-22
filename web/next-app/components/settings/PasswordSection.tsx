'use client'

import { useState } from 'react'
import { changePassword } from '../../lib/api'
import { useT } from '../../lib/i18n'
import Loading from '../Loading'

/** The Change password block of /settings. */
export default function PasswordSection() {
  const t = useT()
  const [cpCurrent, setCpCurrent] = useState('')
  const [cpNew, setCpNew] = useState('')
  const [cpBusy, setCpBusy] = useState(false)
  const [cpMsg, setCpMsg] = useState('')

  const handleChangePassword = async () => {
    if (!cpCurrent || !cpNew) {
      setCpMsg(t('account.settings.passwordBoth'))
      return
    }
    setCpBusy(true)
    setCpMsg('')
    try {
      await changePassword(cpCurrent, cpNew)
      setCpCurrent('')
      setCpNew('')
      setCpMsg(t('account.settings.passwordChanged'))
    } catch (e: unknown) {
      setCpMsg(e instanceof Error ? e.message : t('account.settings.passwordFailed'))
    } finally {
      setCpBusy(false)
    }
  }

  return (
    <div className="border-t border-mathua-border pt-6 min-w-0">
      <span className="font-mono text-sm text-mathua-primary">{t('account.settings.passwordTitle')}</span>
      <p className="text-mathua-muted text-xs mt-1">
        {t('account.settings.passwordHint')}
      </p>
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-2 mt-3">
        <input
          type="password"
          value={cpCurrent}
          onChange={(e) => setCpCurrent(e.target.value)}
          placeholder={t('account.settings.passwordCurrent')}
          autoComplete="current-password"
          aria-label={t('account.settings.passwordCurrent')}
          className="min-w-0 bg-mathua-code border border-mathua-border rounded-none h-10 px-3 font-mono text-sm text-mathua-primary placeholder:text-mathua-muted focus:outline-none focus:border-mathua-blue"
        />
        <input
          type="password"
          value={cpNew}
          onChange={(e) => setCpNew(e.target.value)}
          placeholder={t('account.login.newPassword')}
          autoComplete="new-password"
          aria-label={t('account.login.newPassword')}
          className="min-w-0 bg-mathua-code border border-mathua-border rounded-none h-10 px-3 font-mono text-sm text-mathua-primary placeholder:text-mathua-muted focus:outline-none focus:border-mathua-blue"
        />
      </div>
      <button
        type="button"
        onClick={handleChangePassword}
        disabled={cpBusy}
        className="mt-3 border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint rounded-none px-4 h-10 text-xs font-mono disabled:opacity-50 inline-flex items-center gap-2"
      >
        {cpBusy ? (<><Loading inline size={11} /> {t('common.action.saving')}</>) : t('account.settings.passwordTitle')}
      </button>
      {cpMsg && <p className="font-mono text-[11px] text-mathua-secondary mt-2">{cpMsg}</p>}
    </div>
  )
}
