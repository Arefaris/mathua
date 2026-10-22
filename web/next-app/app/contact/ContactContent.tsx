'use client'

import Link from 'next/link'
import Header from '../../components/Header'
import Footer from '../../components/Footer'
import BottomTabs from '../../components/BottomTabs'
import ContactBox from '../../components/ContactBox'
import RichText from '../../components/RichText'
import { useT } from '../../lib/i18n'

const monoFont = "var(--font-jetbrains-mono), 'JetBrains Mono', monospace"
const bodyFont = "var(--font-space-grotesk), 'Space Grotesk', serif"

const textStyle: React.CSSProperties = {
  fontFamily: bodyFont,
  fontSize: '0.95rem',
  color: 'var(--text-secondary)',
  lineHeight: 1.85,
  marginBottom: '1rem',
}

const h2Style: React.CSSProperties = {
  fontFamily: bodyFont,
  fontWeight: 400,
  fontSize: '1.3rem',
  color: 'var(--text-primary)',
  marginTop: '2.5rem',
  marginBottom: '1rem',
  borderBottom: '0.5px solid var(--border)',
  paddingBottom: '0.5rem',
}

const linkStyle: React.CSSProperties = {
  color: 'var(--text-primary)',
  textDecoration: 'underline',
}

// The body of the contact page. It is a client component because it reads the locale; the page
// itself stays a server component so it can export its metadata.
// The public route to a person.
//
// Three problems arrive here and only one of them belongs in the in-app report box:
// a wrong expected answer needs the captured question to be reproducible at all, a
// broken account needs a channel that is not public, and a code bug is best served
// by an issue tracker where a fix can be referenced. Saying which is which is the
// whole job of this page, because a learner who sends an account problem to a public
// tracker has published it.
//
// No email address: an address on a page is harvested. The box has no separate
// inbox either -- it writes into the report queue, the page the administrator opens
// whenever a question is reported, so a message cannot quietly pile up somewhere
// nobody looks. X and the issue tracker are the two-way channels.
export default function ContactContent() {
  const t = useT()
  return (
    <>
      <Header />
      <div className="max-w-[680px] mx-auto px-4 sm:px-6 py-12 sm:py-20 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-20 overflow-x-hidden min-w-0">
        <Link href="/" style={{ fontFamily: monoFont, fontSize: '13px', color: 'var(--text-muted)' }}>
          {t('shell.contact.back')}
        </Link>
        <div style={{ fontFamily: monoFont, fontSize: '12px', color: 'var(--text-muted)', marginTop: '1.5rem', marginBottom: '0.5rem' }}>
          {t('shell.contact.eyebrow')}
        </div>
        <h1
          style={{
            fontFamily: "var(--font-space-grotesk), 'Space Grotesk', serif",
            fontWeight: 400,
            fontSize: '1.9rem',
            color: 'var(--text-primary)',
            marginBottom: '2rem',
          }}
        >
          {t('shell.contact.title')}
        </h1>

        <p style={textStyle}>{t('shell.contact.intro')}</p>

        <h2 style={h2Style}>{t('shell.contact.message.title')}</h2>
        <p style={textStyle}>{t('shell.contact.message.body')}</p>
        <ContactBox />

        <h2 style={h2Style}>{t('shell.contact.discord.title')}</h2>
        <p style={textStyle}>
          <RichText
            text={t('shell.contact.discord.body')}
            slots={{
              link: (
                <a href="https://discord.gg/qXPeQ2jys" target="_blank" rel="noreferrer" style={linkStyle}>
                  {t('shell.contact.discord.link')}
                </a>
              ),
            }}
          />
        </p>

        <h2 style={h2Style}>{t('shell.contact.wrong.title')}</h2>
        <p style={textStyle}>
          <RichText
            text={t('shell.contact.wrong.body')}
            slots={{ report: <em>{t('shell.contact.wrong.report')}</em> }}
          />
        </p>
        <p style={textStyle}>{t('shell.contact.wrong.body2')}</p>

        <h2 style={h2Style}>{t('shell.contact.broken.title')}</h2>
        <p style={textStyle}>
          <RichText
            text={t('shell.contact.broken.body')}
            slots={{
              x: (
                <a
                  href="https://x.com/wisdomsglow"
                  target="_blank"
                  rel="noreferrer"
                  style={linkStyle}
                >
                  @wisdomsglow
                </a>
              ),
            }}
          />
        </p>
        <p style={textStyle}>{t('shell.contact.broken.body2')}</p>

        <h2 style={h2Style}>{t('shell.contact.bug.title')}</h2>
        <p style={textStyle}>
          <RichText
            text={t('shell.contact.bug.body')}
            slots={{
              link: (
                <a
                  href="https://github.com/chuma-beep/mathua/issues"
                  target="_blank"
                  rel="noreferrer"
                  style={linkStyle}
                >
                  github.com/chuma-beep/mathua/issues
                </a>
              ),
            }}
          />
        </p>

        <h2 style={h2Style}>{t('shell.contact.notSend.title')}</h2>
        <p style={textStyle}>{t('shell.contact.notSend.body')}</p>

        <h2 style={h2Style}>{t('shell.contact.next.title')}</h2>
        <p style={textStyle}>{t('shell.contact.next.body')}</p>
      </div>
      <Footer />
      <BottomTabs />
    </>
  )
}
