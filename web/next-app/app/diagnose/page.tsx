'use client'

import dynamic from 'next/dynamic'
import { useTheme } from '../../hooks/useTheme'
import Link from 'next/link'
import Header from '../../components/Header'
import BottomTabs from '../../components/BottomTabs'
import SectionHeader from '../../components/SectionHeader'
import Footer from '../../components/Footer'
import Loading from '../../components/Loading'
import { useT } from '../../lib/i18n'

// A component rather than an inline arrow, because `dynamic` calls it outside any render the
// page controls and the label needs the locale hook.
function FlowLoading() {
  const t = useT()
  return (
    <div style={{
      height: 420, border: '0.5px solid var(--border)',
      display: 'flex', alignItems: 'center', justifyContent: 'center',
      color: 'var(--text-muted)', fontFamily: "var(--font-jetbrains-mono), 'JetBrains Mono', monospace", fontSize: '13px',
    }}>
      <Loading label={t('goals.diagnose.loading')} />
    </div>
  )
}

const DiagnosticFlow = dynamic(() => import('../../components/DiagnosticFlow'), {
  ssr: false,
  loading: () => <FlowLoading />,
})

const STEP_KEYS = [
  'goals.diagnose.step1',
  'goals.diagnose.step2',
  'goals.diagnose.step3',
  'goals.diagnose.step4',
  'goals.diagnose.step5',
] as const

export default function DiagnosePage() {
  const t = useT()
  const { mounted } = useTheme()

  if (!mounted) return <div style={{ background: 'var(--bg)', minHeight: '100vh' }} />

  return (
    <>
      <Header />
      <div className="max-w-container mx-auto px-4 sm:px-6 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 overflow-x-hidden min-w-0">
      <section className="pt-8 min-w-0 overflow-hidden">
        <span className="flex mb-4">
          <Link href="/" className="text-mathua-secondary text-sm hover:text-mathua-blue">
            {t('goals.diagnose.back')}
          </Link>
        </span>

        <SectionHeader label={t('goals.diagnose.label')} title={t('goals.diagnose.title')} />

        <p className="text-mathua-secondary text-sm leading-relaxed text-center max-w-[600px] mx-auto mt-4">
          {t('goals.diagnose.lead')}
        </p>
      </section>


      {/* How it works */}
      <section className="py-20 max-sm:py-12">
        <SectionHeader label={t('goals.diagnose.methodLabel')} title={t('goals.diagnose.howTitle')} />

        <div className="flex justify-center mt-8">
          <DiagnosticFlow />
        </div>

        <ul className="max-w-[600px] mx-auto mt-8 space-y-3 list-none">
           {STEP_KEYS.map((stepKey) => (
            <li
              key={stepKey}
              className="text-mathua-secondary text-[0.95rem] leading-[1.7] pl-9 relative before:content-[counter(step)] before:absolute before:left-0 before:text-mathua-blue before:font-mono before:text-[13px]"
              style={{ counterIncrement: 'step-counter 1' }}
            >
              {t(stepKey)}
            </li>
          ))}
        </ul>
        <style>{`ul { counter-reset: step-counter; }`}</style>
      </section>


      {/* Actions */}
      <section className="py-20 max-sm:py-12 text-center">
        <SectionHeader label={t('goals.diagnose.readyLabel')} title={t('goals.diagnose.beginTitle')} />

        <p className="text-mathua-secondary text-sm leading-relaxed max-w-[500px] mx-auto mt-4 mb-8">
          {t('goals.diagnose.retake')}
        </p>

        <div className="flex gap-3 justify-center max-sm:flex-col max-sm:items-center min-w-0">
          <Link
            href="/onboard"
            className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint rounded-none h-10 min-h-[36px] px-6 text-[13px] font-medium transition-colors flex items-center justify-center w-full sm:w-auto max-w-[280px]"
          >
            {t('goals.diagnose.start')}
          </Link>
          <Link
            href="/docs/system-design#cat-diagnostic"
            className="border border-mathua-border-strong text-mathua-primary hover:border-mathua-blue hover:text-mathua-blue rounded-none h-10 min-h-[36px] px-6 text-[13px] font-medium transition-colors flex items-center justify-center w-full sm:w-auto max-w-[280px]"
          >
            {t('goals.diagnose.readAlgorithm')}
          </Link>
        </div>
      </section>

    </div>

    <Footer />
    <BottomTabs />
    </>
  )
}
