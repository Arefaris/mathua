'use client'

import { useRouter } from 'next/navigation'
import { useTheme } from '../hooks/useTheme'
import { ensureGuestId, ensureGuestToken } from '../lib/auth'
import { useHomeRedirect } from '../hooks/useHomeRedirect'
import { useScrollReveal } from '../hooks/useScrollReveal'
import { useT } from '../lib/i18n'
import Header from '../components/Header'
import Footer from '../components/Footer'
import {
  CoverageSection,
  FaqSection,
  FeaturesSection,
  HeroSection,
  RankingSection,
  StatesSection,
  TrustSection,
} from './home/sections'

export default function HomePage() {
  const t = useT()
  const { mounted } = useTheme()
  const { push } = useRouter()
  // Logged-in visits bounce to /profile; explicit Home clicks stay put.
  const homeChecked = useHomeRedirect(mounted)
  // Sections only exist once the redirect check resolves, so gate on it.
  useScrollReveal(mounted && homeChecked)

  if (!mounted || !homeChecked) {
    return <div style={{ background: 'var(--bg)', minHeight: '100vh' }} />
  }

  function handleGetStarted() {
    // Navigate first, mint in the background.
    //
    // The token is only needed for the API calls the next page makes, and
    // chaining the push onto it made this button look dead for the request's
    // full ten-second timeout whenever that request stalled, a filtered or
    // blocked API origin does exactly that. Ten seconds of nothing is
    // indistinguishable from "the page won't load", which is how this was
    // reported. The token still arrives before the first answer in every
    // ordinary case, and if it never does the learn page is anonymously usable.
    ensureGuestId()
    push('/learn')
    void ensureGuestToken()
  }

  return (
    <>
      {/* The landing hardcodes its own header rather than inheriting lib/nav.ts,
          because a marketing page wants Docs and How it works inline. Learn is
          still listed first: the primary destination of the product should not
          vanish on the page that introduces it. */}
      <Header links={[
        { label: t('explore.nav.learn'), href: '/learn' },
        { label: t('explore.nav.howItWorks'), href: '/how-it-works' },
        { label: t('explore.nav.docs'), href: '/docs' },
        { label: t('explore.nav.note'), href: '/note' },
        { label: t('explore.nav.leaderboard'), href: '/leaderboard' },
        { label: t('explore.nav.login'), href: '/login' },
      ]} />
      <main className="max-w-container mx-auto px-4 sm:px-6 overflow-x-clip min-w-0">
        <HeroSection onGetStarted={handleGetStarted} />
        <TrustSection />
        <FeaturesSection />
        <StatesSection />
        <CoverageSection />
        <RankingSection />
        <FaqSection />
      </main>
      <Footer />
    </>
  )
}
