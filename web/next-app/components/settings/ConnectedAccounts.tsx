'use client'

import { useState, useEffect } from 'react'
import { getIdentities, deleteIdentity, createLinkToken, startOAuthLogin, getConfig, OAUTH_LABELS, type OAuthProvider } from '../../lib/api'
import { useT } from '../../lib/i18n'

/** The Connected accounts block of /settings: which sign-in providers are linked, and the rest to add. */
export default function ConnectedAccounts() {
  const t = useT()
  const [identities, setIdentities] = useState<{ provider: string; email: string }[]>([])
  const [idBusy, setIdBusy] = useState(false)
  const [idMsg, setIdMsg] = useState('')
  const [providers, setProviders] = useState<OAuthProvider[]>([])

  useEffect(() => {
    getIdentities().then(setIdentities).catch(() => {})
    getConfig().then(cfg => {
      if (cfg && Array.isArray(cfg.providers)) {
        setProviders(cfg.providers.filter((p): p is OAuthProvider => p in OAUTH_LABELS))
      }
    }).catch(() => {})
  }, [])

  const refreshIdentities = () => {
    getIdentities().then(setIdentities).catch(() => {})
  }

  const handleConnect = async (provider: string) => {
    setIdBusy(true)
    setIdMsg('')
    try {
      const linkToken = await createLinkToken()
      startOAuthLogin(provider, { intent: 'link', linkToken })
    } catch {
      setIdMsg(t('account.settings.connectFailed'))
      setIdBusy(false)
    }
  }

  const handleDisconnect = async (provider: string) => {
    setIdBusy(true)
    setIdMsg('')
    try {
      await deleteIdentity(provider)
      refreshIdentities()
    } catch (e: unknown) {
      setIdMsg(e instanceof Error ? e.message : t('account.settings.disconnectFailed'))
    } finally {
      setIdBusy(false)
    }
  }

  // Single-pass lookup for the connect-button list below.
  const connectedSet = new Set(identities.map(id => id.provider))

  return (
    <div className="border-t border-mathua-border pt-6 min-w-0">
      <span className="font-mono text-sm text-mathua-primary">{t('account.settings.connectedTitle')}</span>
      <p className="text-mathua-muted text-xs mt-1">
        {t('account.settings.connectedHint')}
      </p>
      <div className="mt-3 space-y-2">
        {identities.map(id => (
          <div key={id.provider} className="flex items-center justify-between gap-3 border border-mathua-border px-3 h-10 min-w-0">
            <span className="font-mono text-xs text-mathua-primary truncate">
              {OAUTH_LABELS[id.provider as OAuthProvider] ?? id.provider}
              {id.email && <span className="text-mathua-muted"> · {id.email}</span>}
            </span>
            <button
              type="button"
              onClick={() => handleDisconnect(id.provider)}
              disabled={idBusy}
              className="shrink-0 font-mono text-[11px] text-mathua-muted hover:text-mathua-red disabled:opacity-50"
            >
              {t('account.settings.disconnect')}
            </button>
          </div>
        ))}
        {providers.flatMap(p => connectedSet.has(p) ? [] : [
          <button
            type="button"
            key={p}
            onClick={() => handleConnect(p)}
            disabled={idBusy}
            className="w-full border border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue rounded-none px-4 h-10 text-xs font-mono disabled:opacity-50"
          >
            {t('account.settings.connect', { provider: OAUTH_LABELS[p] })}
          </button>
        ])}
        {providers.length === 0 && (
          <p className="font-mono text-[11px] text-mathua-muted">{t('account.settings.noProviders')}</p>
        )}
      </div>
      {idMsg && <p className="font-mono text-[11px] text-mathua-secondary mt-2">{idMsg}</p>}
    </div>
  )
}
