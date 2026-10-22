'use client'

import { useT } from '@/lib/i18n'

// Root error boundary: replaces the whole layout when the root itself
// throws, so it must render its own <html>/<body> and cannot rely on the
// root layout's fonts, theme script, or Toaster.
export default function GlobalError({
  error,
  reset,
}: {
  error: Error & { digest?: string }
  reset: () => void
}) {
  // This boundary replaces the root layout, so LocaleRoot is not mounted above it; the locale
  // store works without one.
  const t = useT()
  return (
    <html lang={t.locale} className="dark">
      <body style={{ background: '#0b0f1a', color: '#e8e2d5', fontFamily: 'monospace' }}>
        <main style={{ display: 'flex', minHeight: '100vh', alignItems: 'center', justifyContent: 'center', padding: '24px' }}>
          <section style={{ width: '100%', maxWidth: '640px' }}>
            <div style={{ fontSize: '10px', textTransform: 'uppercase', letterSpacing: '0.18em', opacity: 0.6 }}>
              {t('shell.globalError.crumb')}
            </div>
            <h1 style={{ fontSize: '2rem', marginTop: '24px' }}>{t('shell.globalError.heading')}</h1>
            <p role="alert" style={{ marginTop: '12px', fontSize: '14px', opacity: 0.7 }}>
              {error?.message || t('shell.errorPage.fallbackMessage')}
              {error?.digest ? ` ${t('shell.errorPage.ref', { digest: error.digest })}` : ''}
            </p>
            <div style={{ marginTop: '32px', display: 'flex', gap: '12px' }}>
              <button
                type="button"
                onClick={() => reset()}
                style={{ border: '1px solid #c8a96e', color: '#c8a96e', background: 'transparent', padding: '10px 20px', fontSize: '12px', cursor: 'pointer' }}
              >
                {t('shell.errorPage.tryAgain')}
              </button>
              <a href="/" style={{ border: '1px solid #444', color: '#e8e2d5', padding: '10px 20px', fontSize: '12px', textDecoration: 'none' }}>
                {t('shell.errorPage.returnHome')}
              </a>
            </div>
          </section>
        </main>
      </body>
    </html>
  )
}
