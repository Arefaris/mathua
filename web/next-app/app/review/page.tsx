'use client'

import Link from 'next/link'
import Header from '../../components/Header'
import BottomTabs from '../../components/BottomTabs'
import Footer from '../../components/Footer'
import ReviewHost from './ReviewHost'
import { useT } from '../../lib/i18n'

export default function ReviewPage() {
  const t = useT()
  return (
    <>
      <Header />
      <div className="max-w-container mx-auto px-4 sm:px-6 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 overflow-x-hidden min-w-0">
        <section className="pt-8 min-w-0 overflow-hidden">
          <span className="flex mb-4">
            <Link href="/profile" className="text-mathua-secondary text-sm hover:text-mathua-primary">
              {t('learn.review.back')}
            </Link>
          </span>
          <ReviewHost />
          <p className="text-mathua-muted text-xs font-mono text-center mt-6">
            {t('learn.review.footnote')}
          </p>
        </section>
      </div>
      <Footer />
      <BottomTabs />
    </>
  )
}
