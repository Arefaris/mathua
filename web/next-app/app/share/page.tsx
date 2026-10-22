'use client'

import { useEffect, useState, Suspense } from 'react'
import { useSearchParams } from 'next/navigation'
import Link from 'next/link'
import Header from '../../components/Header'
import BottomTabs from '../../components/BottomTabs'
import SectionHeader from '../../components/SectionHeader'
import Footer from '../../components/Footer'
import Loading from '../../components/Loading'
import ProfileStats from '../../components/ProfileStats'
import ActivityHeatmap from '../../components/ActivityHeatmap'
import DomainTable from '../../components/DomainTable'
import { concepts } from '../../lib/conceptData'
import { buildDomainRows } from '../../lib/domainRows'
import AttemptList from '../../components/AttemptList'
import { getShareReport, type ShareReport } from '../../lib/api'
import { useT } from '../../lib/i18n'

function ShareContent() {
  const t = useT()
  const searchParams = useSearchParams()
  const token = searchParams.get('token') || ''
  const [report, setReport] = useState<ShareReport | null>(null)
  const [error, setError] = useState('')

  useEffect(() => {
    if (!token) {
      setError(t('account.share.missingToken'))
      return
    }
    getShareReport(token)
      .then(r => r ? setReport(r) : setError(t('account.share.invalid')))
      .catch(() => setError(t('account.share.failed')))
  }, [token, t])

  if (error) {
    return (
      <>
        <Header />
        <div className="max-w-container mx-auto px-4 sm:px-6 pt-[var(--chrome-top)] pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 text-center overflow-x-hidden min-w-0">
          <p className="text-mathua-muted text-sm">{error}</p>
          {!token && (
            <p className="text-mathua-muted text-xs font-mono mt-2">{t('account.share.askEnable')}</p>
          )}
          <Link href="/" className="text-mathua-blue text-sm hover:underline mt-4 inline-block">
            {t('account.share.backHome')}
          </Link>
        </div>
        <Footer />
        <BottomTabs />
      </>
    )
  }

  if (!report) {
    return (
      <>
        <Header />
        <div className="max-w-container mx-auto px-4 sm:px-6 pt-[var(--chrome-top)] pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 text-center overflow-x-hidden min-w-0">
          <Loading label={t('account.share.loading')} />
        </div>
        <Footer />
        <BottomTabs />
      </>
    )
  }

  return (
    <>
      <Header />
      <div className="max-w-container mx-auto px-4 sm:px-6 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 overflow-x-hidden min-w-0">
        <section className="pt-8 min-w-0 overflow-hidden">
          <SectionHeader label={t('account.share.label')} title={t('account.share.title', { name: report.name })} />
          <p className="text-mathua-muted text-xs font-mono text-center mt-3">
            {t('account.share.sharedBy')}
          </p>

          <div className="mx-auto w-full max-w-[820px] min-w-0 px-0 sm:px-4 py-8 overflow-hidden">
            <ProfileStats name={report.name} scores={report.scores} />

            <section className="mt-8 flex min-w-0 flex-col items-stretch">
              <h2 className="font-serif text-[1.05rem] font-normal text-mathua-primary mb-4 w-full">
                {t('account.share.activity')}
              </h2>
              <div className="w-full max-w-full min-w-0 flex justify-center overflow-hidden">
                <div className="w-full max-w-full min-w-0">
                  <ActivityHeatmap data={report.activity} />
                </div>
              </div>
            </section>

            <section className="mt-8 min-w-0">
              {/* No links and no expand: this is a shared read-only snapshot, so /learn is not
                  the viewer's to be sent to and there is nothing behind a concept here. */}
              <DomainTable rows={buildDomainRows(concepts, report.progress)} links={false} />
            </section>

            {(report.attempts?.length ?? 0) > 0 && (
              <section className="mt-8 min-w-0">
                <h2 className="font-serif text-[1.05rem] font-normal text-mathua-primary mb-4 w-full">
                  {t('account.share.progressCard')}
                </h2>
                <AttemptList attempts={(report.attempts ?? []).slice(0, 100)} />
                {(report.attempts ?? []).length > 100 && (
                  <p className="font-mono text-[11px] text-mathua-muted mt-3 text-center">
                    {t('account.share.latest', { total: (report.attempts ?? []).length })}
                  </p>
                )}
              </section>
            )}
          </div>
        </section>
      </div>

      <Footer />
      <BottomTabs />
    </>
  )
}

export default function SharePage() {
  const t = useT()
  return (
    <Suspense fallback={
      <><Header /><div className="max-w-container mx-auto px-4 sm:px-6 pt-[var(--chrome-top)] pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 text-center overflow-x-hidden min-w-0"><Loading label={t('account.share.loading')} /></div></>
    }>
      <ShareContent />
    </Suspense>
  )
}
