'use client'

import Link from 'next/link'
import RichText from '../../components/RichText'
import { useLocale, useT } from '../../lib/i18n'
import type { MessageKey } from '../../lib/i18n/keys'

const monoFont = "var(--font-jetbrains-mono), 'JetBrains Mono', monospace"
const bodyFont = "var(--font-space-grotesk), 'Space Grotesk', serif"

// The list is data, not six copied blocks. Each entry's title and description match the
// generated docs module, so the index and the page it links to cannot disagree. Adding a doc is
// one line here and a route beside the others.
const DOCS = [
  { href: '/docs/architecture', titleKey: 'shell.docsIndex.items.architecture.title', blurbKey: 'shell.docsIndex.items.architecture.blurb' },
  { href: '/docs/system-design', titleKey: 'shell.docsIndex.items.systemDesign.title', blurbKey: 'shell.docsIndex.items.systemDesign.blurb' },
  { href: '/docs/administration', titleKey: 'shell.docsIndex.items.administration.title', blurbKey: 'shell.docsIndex.items.administration.blurb' },
  { href: '/docs/math-input', titleKey: 'shell.docsIndex.items.mathInput.title', blurbKey: 'shell.docsIndex.items.mathInput.blurb' },
  { href: '/docs/contributing', titleKey: 'shell.docsIndex.items.contributing.title', blurbKey: 'shell.docsIndex.items.contributing.blurb' },
  { href: '/docs/efficacy', titleKey: 'shell.docsIndex.items.efficacy.title', blurbKey: 'shell.docsIndex.items.efficacy.blurb' },
] satisfies { href: string; titleKey: MessageKey; blurbKey: MessageKey }[]

export default function DocsIndex() {
  const t = useT()
  const locale = useLocale()
  const textStyle: React.CSSProperties = {
    fontFamily: bodyFont,
    fontSize: '0.95rem',
    color: 'var(--text-secondary)',
    lineHeight: 1.7,
  }

  const rowStyle: React.CSSProperties = {
    display: 'flex',
    gap: '1rem',
    padding: '0.75rem 0',
    borderBottom: '0.5px solid var(--border)',
    color: 'inherit',
    textDecoration: 'none',
  }

  return (
    <div className="max-w-[720px] mx-auto px-4 sm:px-6 py-12 sm:py-20 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-20 overflow-x-hidden min-w-0">
      <Link href="/" style={{ fontFamily: monoFont, fontSize: '13px', color: 'var(--text-muted)' }}>
        {t('shell.docsIndex.back')}
      </Link>
      <h1 style={{
        fontFamily: "var(--font-space-grotesk), 'Space Grotesk', serif",
        fontWeight: 400,
        fontSize: '1.9rem',
        color: 'var(--text-primary)',
        marginBottom: '0.5rem',
        borderBottom: '0.5px solid var(--border)',
        paddingBottom: '0.5rem',
      }}>
        {t('shell.docsIndex.title')}
      </h1>

      <p style={{ ...textStyle, marginBottom: '2.5rem' }}>
        {t('shell.docsIndex.intro')}
      </p>
      {/* The index is translated; the articles it links to are long-form and stay in English
          this round, and a reader should hear that before clicking rather than after. */}
      {locale !== 'en' && (
        <p style={{ ...textStyle, marginTop: '-1.75rem', marginBottom: '2.5rem', color: 'var(--text-muted)' }}>
          {t('shell.docsIndex.englishOnly')}
        </p>
      )}

      <div style={{ marginBottom: '3rem' }}>
        {DOCS.map((doc, i) => (
          <Link key={doc.href} href={doc.href} className="link-underline" style={rowStyle}>
            <span style={{ fontFamily: monoFont, fontSize: '13px', color: 'var(--accent-blue)', minWidth: '1.5rem', paddingTop: '1px' }}>
              {i + 1}.
            </span>
            <div>
              <div style={{ fontFamily: "var(--font-space-grotesk), 'Space Grotesk', serif", fontWeight: 400, fontSize: '1.1rem', color: 'var(--text-primary)', marginBottom: '2px' }}>
                {t(doc.titleKey)}
              </div>
              <p style={{ ...textStyle, margin: 0, fontSize: '0.9rem', color: 'var(--text-muted)' }}>
                {t(doc.blurbKey)}
              </p>
            </div>
          </Link>
        ))}
      </div>

      <div style={{ borderTop: '0.5px solid var(--border)', paddingTop: '3rem' }}>
        <h2 style={{
          fontFamily: "var(--font-space-grotesk), 'Space Grotesk', serif",
          fontWeight: 400,
          fontSize: '1.3rem',
          color: 'var(--text-primary)',
          marginBottom: '1rem',
        }}>
          {t('shell.docsIndex.quickstart')} <span style={{ fontFamily: monoFont, fontSize: '12px', color: 'var(--text-muted)' }}>{t('shell.docsIndex.forContributors')}</span>
        </h2>

        <pre style={{
          fontFamily: monoFont,
          fontSize: '13px',
          color: 'var(--text-secondary)',
          whiteSpace: 'pre',
          overflowX: 'auto',
          lineHeight: 1.6,
          margin: 0,
          width: '100%',
          maxWidth: '100%',
          minWidth: 0,
        }}>
{`git clone https://github.com/chuma-beep/mathua.git
cd mathua
go build -o bin/mathua ./cmd/mathua

# Run the web server (SQLite by default)
./bin/mathua --serve --port 8080`}
        </pre>

        <p style={{
          fontFamily: bodyFont,
          fontStyle: 'italic',
          fontSize: '0.9rem',
          color: 'var(--text-muted)',
          marginTop: '1rem',
          lineHeight: 1.7,
        }}>
          <RichText
            text={t('shell.docsIndex.requirements')}
            slots={{
              howItWorks: <Link href="/how-it-works" className="link-underline" style={{ color: 'var(--accent-blue)' }}>{t('shell.footer.howItWorks')}</Link>,
              licenses: <Link href="/licenses" className="link-underline" style={{ color: 'var(--accent-blue)' }}>{t('shell.footer.licenses')}</Link>,
            }}
          />
        </p>
      </div>
    </div>
  )
}
