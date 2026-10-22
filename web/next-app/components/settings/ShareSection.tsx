'use client'

import { useState, useEffect } from 'react'
import { enableShare, disableShare } from '../../lib/api'
import { useT } from '../../lib/i18n'

/** The Share with parent / teacher block of /settings: a read-only link to the learner's report. */
export default function ShareSection() {
  const t = useT()
  const [, setShareToken] = useState('')
  const [shareUrl, setShareUrl] = useState('')
  const [shareBusy, setShareBusy] = useState(false)
  const [copied, setCopied] = useState(false)

  useEffect(() => {
    try {
      const savedShare = localStorage.getItem('mathua_share_url')
      if (savedShare) setShareUrl(savedShare)
    } catch { /* private mode — share link just won't persist */ }
  }, [])

  const handleShare = async (enable: boolean) => {
    setShareBusy(true)
    try {
      if (enable) {
        const res = await enableShare()
        setShareToken(res.token)
        const url = `${window.location.origin}/share?token=${res.token}`
        setShareUrl(url)
        try { localStorage.setItem('mathua_share_url', url) } catch { /* ignore */ }
      } else {
        await disableShare()
        setShareToken('')
        setShareUrl('')
        try { localStorage.removeItem('mathua_share_url') } catch { /* ignore */ }
      }
    } catch {
      console.error('share toggle failed')
    } finally {
      setShareBusy(false)
    }
  }

  return (
    <div className="border-t border-mathua-border pt-6 min-w-0">
      <span className="font-mono text-sm text-mathua-primary">{t('account.settings.shareTitle')}</span>
      <p className="text-mathua-muted text-xs mt-1">
        {t('account.settings.shareHint')}
      </p>
      {shareUrl ? (
        <div className="mt-3 min-w-0">
          <input
            readOnly
            value={shareUrl}
            aria-label={t('account.settings.shareLink')}
            onFocus={e => e.currentTarget.select()}
            className="w-full bg-mathua-code border border-mathua-border rounded-none h-10 px-3 font-mono text-xs text-mathua-primary focus:outline-none focus:border-mathua-blue min-w-0"
          />
          <div className="flex gap-2 mt-2">
            <button
              type="button"
              onClick={() => {
                const done = () => { setCopied(true); setTimeout(() => setCopied(false), 2000) }
                if (navigator.clipboard?.writeText) {
                  navigator.clipboard.writeText(shareUrl).then(done).catch(() => done())
                } else { done() }
              }}
              className="border border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue rounded-none px-4 h-9 text-xs font-mono"
            >
              {copied ? t('common.action.copied') : t('account.settings.copyLink')}
            </button>
            <button
              type="button"
              onClick={() => handleShare(false)}
              disabled={shareBusy}
              className="border border-mathua-red text-mathua-red hover:bg-mathua-red hover:text-white rounded-none px-4 h-9 text-xs font-mono disabled:opacity-50"
            >
              {t('account.settings.disableShare')}
            </button>
          </div>
        </div>
      ) : (
        <button
          type="button"
          onClick={() => handleShare(true)}
          disabled={shareBusy}
          className="mt-3 border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint rounded-none px-4 h-9 text-xs font-mono disabled:opacity-50"
        >
          {shareBusy ? t('account.settings.generating') : t('account.settings.enableShare')}
        </button>
      )}
    </div>
  )
}
