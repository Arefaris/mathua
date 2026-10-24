'use client'

import type { RefObject } from 'react'
import Link from 'next/link'
import KatexContent from '../../components/KatexContent'
import SectionHeader from '../../components/SectionHeader'
import ProgressBar from '../../components/ProgressBar'
import ReportButton from '../../components/ReportButton'
import BriefingCard, { DIAGNOSTIC_BRIEFING } from '../../components/BriefingCard'
import SubmitErrorBlock, { type SubmitError } from '../../components/SubmitErrorBlock'
import DiagnosticResults from '../../components/DiagnosticResults'
import Loading from '../../components/Loading'
import type { DiagnosticProgress, GoalPlanRes } from '../../lib/api'
import type { AnswerFormat } from '../../lib/answerFormat'
import ChoiceOptions from '../../components/ChoiceOptions'
import { stemWithoutChoices } from '../../lib/choices'
import { domainLabels, type DomainInfo } from './domains'
import { useT } from '../../lib/i18n'
import { MathAnswerInput, type MathFocusHandle } from '@/components/math/MathInput'
import { focusClearOfTabs } from '../../lib/focusClearOfTabs'

export function WelcomeStep({
  domains,
  loading,
  hasPaused,
  selectedCount,
  confirming,
  onToggle,
  onSelectAll,
  onStart,
  onResume,
  onBegin,
  onCancel,
  needsName,
  demoName,
  onDemoName,
  onStartDemo,
  onCancelDemo,
}: {
  domains: DomainInfo[]
  loading: boolean
  hasPaused: boolean
  selectedCount: number
  confirming: boolean
  onToggle: (name: string) => void
  onSelectAll: () => void
  onStart: () => void
  onResume: () => void
  onBegin: () => void
  onCancel: () => void
  needsName: boolean
  demoName: string
  onDemoName: (name: string) => void
  onStartDemo: () => void
  onCancelDemo: () => void
}) {
  const t = useT()
  if (confirming) {
    const selected = domains.filter(d => d.selected)
    const shown = selected.slice(0, 3).map(d => domainLabels(d.name))
    const extra = selected.length > 3 ? t('goals.onboard.moreAreas', { count: selected.length - 3 }) : ''
    return (
      <div className="max-w-2xl mx-auto px-2">
        <div className="text-center mb-4 mt-6 sm:mt-8 min-w-0">
          <div className="font-mono text-[11px] uppercase text-mathua-muted mb-3">{t('goals.onboard.confirmEyebrow')}</div>
          <h1 className="font-serif text-2xl sm:text-3xl font-medium text-mathua-primary px-2">
            {t('goals.onboard.confirmTitle')}
          </h1>
          <p className="text-mathua-secondary text-sm mt-3 px-2">
            {t('goals.onboard.confirmScope', { areas: shown.join(', '), extra, count: selectedCount })}
          </p>
          <p className="text-mathua-muted text-xs mt-2 px-2">
            {t('goals.onboard.confirmDuration')}
          </p>
          <p className="text-mathua-muted text-xs mt-1 px-2">
            {t('goals.onboard.confirmFoundations')}
          </p>
        </div>
        <BriefingCard eyebrow={t('goals.onboard.expect')} items={DIAGNOSTIC_BRIEFING} />
        <div className="mt-6 flex flex-col sm:flex-row gap-3 justify-center px-2">
          <button
            type="button"
            onClick={onBegin}
            disabled={loading}
            className="border border-mathua-blue bg-mathua-blue text-white hover:opacity-90 rounded-none h-12 min-h-[44px] px-6 sm:px-10 font-medium text-sm disabled:opacity-50"
          >
            {loading ? (<><Loading inline size={13} /> {t('common.status.loading')}</>) : t('goals.onboard.begin')}
          </button>
          <button
            type="button"
            onClick={onCancel}
            disabled={loading}
            className="font-mono text-xs text-mathua-muted hover:text-mathua-primary min-h-[44px] px-4 disabled:opacity-50"
          >
            {t('goals.onboard.back')}
          </button>
        </div>
      </div>
    )
  }
  // A guest who has no token yet. The name is not an account: the server makes
  // a student row so a session has an owner, and this screen says plainly that
  // the row cannot be signed back into, rather than letting a guest believe
  // their run was saved.
  if (needsName) {
    return (
      <div className="max-w-2xl mx-auto px-2">
        <div className="text-center mb-4 mt-6 sm:mt-8 min-w-0">
          <div className="font-mono text-[11px] uppercase text-mathua-muted mb-3">{t('goals.onboard.guestEyebrow')}</div>
          <h1 className="font-serif text-2xl sm:text-3xl font-medium text-mathua-primary px-2">
            {t('goals.onboard.guestTitle')}
          </h1>
          <p className="text-mathua-secondary text-sm mt-3 px-2">
            {t('goals.onboard.guestBody')}
          </p>
        </div>
        <div className="border border-mathua-border bg-mathua-surface p-4 sm:p-6">
          <label htmlFor="demo-name" className="block font-mono text-[11px] uppercase text-mathua-muted mb-2">
            {t('goals.onboard.nameLabel')}
          </label>
          <input
            id="demo-name"
            type="text"
            value={demoName}
            onChange={e => onDemoName(e.target.value)}
            maxLength={40}
            autoComplete="nickname"
            placeholder={t('goals.onboard.namePlaceholder')}
            className="w-full h-11 px-3 bg-mathua-bg border border-mathua-border text-mathua-primary text-sm rounded-none focus:outline-none focus:border-mathua-blue"
          />
          <p className="text-mathua-muted text-xs mt-2">
            {t('goals.onboard.nameNote')}
          </p>
        </div>
        <div className="mt-6 flex flex-col sm:flex-row gap-3 justify-center px-2">
          <button
            type="button"
            onClick={onStartDemo}
            disabled={loading || demoName.trim() === ''}
            className="border border-mathua-blue bg-mathua-blue text-white hover:opacity-90 rounded-none h-12 min-h-[44px] px-6 sm:px-10 font-medium text-sm disabled:opacity-50"
          >
            {loading ? (<><Loading inline size={13} /> {t('common.status.loading')}</>) : t('goals.onboard.runDemo')}
          </button>
          <button
            type="button"
            onClick={onCancelDemo}
            disabled={loading}
            className="font-mono text-xs text-mathua-muted hover:text-mathua-primary min-h-[44px] px-4 disabled:opacity-50"
          >
            {t('goals.onboard.back')}
          </button>
        </div>
        <p className="text-center mt-4">
          <Link href="/login" className="font-mono text-xs text-mathua-blue hover:underline">
            {t('goals.onboard.signInSaved')}
          </Link>
        </p>
      </div>
    )
  }
  return (
    <div className="max-w-4xl mx-auto">
      <div className="text-center mb-4 mt-6 sm:mt-8 px-2 min-w-0">
        <div className="font-mono text-[11px] uppercase text-mathua-muted mb-3">{t('goals.onboard.welcomeEyebrow')}</div>
        <h1 className="font-serif text-2xl sm:text-4xl font-medium text-mathua-primary px-2">
          {t('goals.onboard.welcomeTitle')}
        </h1>
        <p className="text-mathua-secondary text-sm mt-3 max-w-[500px] mx-auto px-2">
          {t('goals.onboard.welcomeIntro')}
        </p>
      </div>

      <div className="flex gap-3 justify-center mb-6 px-2">
        <button
          type="button"
          onClick={onSelectAll}
          aria-pressed={domains.length > 0 && domains.every(d => d.selected)}
          className="bg-mathua-surface border border-mathua-border rounded-none h-10 min-h-[36px] px-6 text-sm text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue"
        >
          {domains.length > 0 && domains.every(d => d.selected) ? t('goals.onboard.deselectAll') : t('goals.onboard.selectAll')}
        </button>
      </div>

      <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-2 sm:gap-3 mb-8 min-w-0">
        {domains.map(d => {
          const label = domainLabels(d.name)
          return (
            <button
              type="button"
              key={d.name}
              onClick={() => onToggle(d.name)}
              className={`rounded-none p-3 sm:p-4 text-left transition-all text-sm min-h-[60px] min-w-0 overflow-hidden ${
                d.selected
                  ? 'bg-mathua-blue text-white'
                  : 'bg-mathua-surface border border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue'
              }`}
            >
              <div className="font-medium truncate sm:whitespace-normal sm:line-clamp-2 break-words text-[13px] sm:text-sm">{label}</div>
              <div className={`font-mono text-[10px] mt-1 truncate ${d.selected ? 'text-white/70' : 'text-mathua-muted'}`}>
                {t('common.count.concepts', { count: d.count })}
              </div>
            </button>
          )
        })}
      </div>

      <div className="text-center px-4">
        {hasPaused && (
          <button
            type="button"
            onClick={onResume}
            disabled={loading}
            className="border border-mathua-blue bg-mathua-blue text-white hover:opacity-90 rounded-none h-12 min-h-[44px] px-6 sm:px-10 font-medium text-sm disabled:opacity-50 max-w-full mb-3"
          >
            {loading ? (<><Loading inline size={13} /> {t('common.status.loading')}</>) : t('goals.onboard.continue')}
          </button>
        )}
        <button
          type="button"
          onClick={onStart}
          disabled={selectedCount === 0 || loading}
          className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint rounded-none h-12 min-h-[44px] px-6 sm:px-10 font-medium text-sm disabled:opacity-50 max-w-full"
        >
          {loading ? (<><Loading inline size={13} /> {t('common.status.loading')}</>) : t('goals.onboard.start', { count: selectedCount })}
        </button>
      </div>

      {/* The way out, and the reason it is a link rather than a flag.
          A learner can start any topic whose prerequisites are met without ever
          taking the diagnostic: /learn?concept=… is gated server-side, not by
          this flag, and /domains is already browsable. Nothing here records a
          result, seeds placement, or grants credit, it only makes an existing
          path findable from the screen a learner is actually looking at. */}
      <div className="text-center px-4 mt-6">
        <Link
          href="/domains"
          className="font-mono text-xs text-mathua-muted hover:text-mathua-blue underline underline-offset-2 min-h-[44px] inline-flex items-center"
        >
          {t('goals.onboard.browse')}
        </Link>
      </div>
    </div>
  )
}

export function DiagnosticStep({
  question,
  conceptName,
  progress,
  questionCount,
  answerInput,
  lastResult,
  accuracy,
  loading,
  inputRef,
  conceptId,
  sessionId,
  onInputChange,
  onSubmit,
  onDontKnow,
  onNext,
  done,
  answerFormat,
  submitError,
  planError,
  onSkip,
  skipsLeft,
  onRestart,
  onRetrySubmit,
  onRetryPlan,
  retryAvailable,
  onRetryQuestion,
}: {
  question: string
  conceptName: string
  progress: DiagnosticProgress | null
  questionCount: number
  answerInput: string
  lastResult: { correct: boolean; feedback: string } | null
  accuracy: { correct: number; total: number }
  loading: boolean
  inputRef: RefObject<MathFocusHandle | null>
  conceptId: string
  sessionId: string
  onInputChange: (value: string) => void
  onSubmit: () => void
  onDontKnow: () => void
  onNext: () => void
  done: boolean
  answerFormat: AnswerFormat
  submitError: SubmitError | null
  planError: string
  onSkip: () => void
  skipsLeft: number
  onRestart: () => void
  /** Retry the failed submission as the kind of answer it was. Optional because
   *  a caller that only ever submits graded answers can rely on onSubmit. */
  onRetrySubmit?: () => void
  onRetryPlan: () => void
  retryAvailable: boolean
  onRetryQuestion: () => void
}) {
  const t = useT()
  // The number of the question on screen. `progress.answered` already counts it once
  // its feedback is showing, so adding one then labelled the result of question 1 as
  // "Question 2" before Next had been pressed.
  const shownNumber = progress ? progress.answered + (lastResult ? 0 : 1) : questionCount
  return (
    <>
      <SectionHeader label={t('goals.question.number', { count: shownNumber })} title={conceptName} />
      <div className="max-w-2xl mx-auto px-2 sm:px-0 min-w-0 overflow-hidden">
        <ProgressBar
          answered={progress ? progress.answered : 0}
          budget={progress ? progress.probe_budget : 0}
          question={shownNumber}
        />
        <div className="bg-mathua-surface border border-mathua-border rounded-none p-4 sm:p-6 mb-6 w-full max-w-full min-w-0 overflow-hidden">
          <div className="bg-mathua-code border border-mathua-border rounded-none p-4 sm:p-6 text-center mb-4 w-full max-w-full min-w-0 overflow-hidden">
            <div className="w-full max-w-full min-w-0 overflow-hidden">
              <KatexContent className="text-mathua-primary text-lg font-mono font-light whitespace-pre-wrap break-words">{stemWithoutChoices(question)}</KatexContent>
            </div>
          </div>
          {/* The controls belong to the question that is still unanswered.
              Once a verdict is in, they are not rendered at all: leaving a
              disabled field beside a disabled button looked like a frozen
              page, and the same controls were already showing the next
              question underneath the verdict. */}
          {!lastResult && (
            <>
          <ChoiceOptions question={question} value={answerInput} onPick={onInputChange} disabled={loading || lastResult !== null} />
          <form
            onSubmit={e => { e.preventDefault(); onSubmit() }}
            className="flex flex-col sm:flex-row gap-3 min-w-0"
          >
            <label htmlFor="onboard-answer" className="sr-only">{t('goals.question.yourAnswer')}</label>
            <MathAnswerInput
              id="onboard-answer"
              value={answerInput}
              onChange={onInputChange}
              onSubmit={onSubmit}
              conceptId={conceptId}
              disabled={loading || lastResult !== null}
              status={lastResult ? (lastResult.correct ? 'correct' : 'incorrect') : 'default'}
              inputMode={answerFormat.inputMode}
              placeholder={t('goals.question.placeholder')}
              className="sm:flex-1"
              focusRef={inputRef}
            />
            <button
              type="submit"
              disabled={!answerInput.trim() || loading || lastResult !== null}
              className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint rounded-none h-12 min-h-[36px] px-6 font-medium text-sm disabled:opacity-50 shrink-0 w-auto self-end sm:self-auto"
            >
              {t('goals.question.checkAnswer')}
            </button>
          </form>
            </>
          )}
          <div className="mt-2 flex items-center justify-between gap-2">
            <p className="font-mono text-[11px] text-mathua-muted">{answerFormat.hint}</p>
            {!lastResult && (
              <button
                type="button"
                onClick={onDontKnow}
                disabled={loading}
                className="shrink-0 font-mono text-[11px] text-mathua-muted hover:text-mathua-red underline underline-offset-2 disabled:opacity-50"
              >
                {loading ? t('goals.question.submitting') : t('goals.question.dontKnow')}
              </button>
            )}
          </div>
          <div className="mt-2 flex justify-end">
            <ReportButton
              key={`diag-${shownNumber}-${question}`}
              conceptId={conceptId}
              kind="question"
              question={question}
              source="diagnostic"
              sessionId={sessionId}
              answerGiven={answerInput}
              answerHint={answerFormat.hint}
              answerInputMode={answerFormat.inputMode}
            />
          </div>
        </div>

        {lastResult && (
          <div className={`bg-mathua-surface border rounded-none p-4 mb-4 text-center ${lastResult.correct ? 'border-mathua-green' : 'border-mathua-red'}`}>
            <KatexContent className={lastResult.correct ? 'text-mathua-green' : 'text-mathua-red'}>{lastResult.feedback}</KatexContent>
            {/* A miss offers the teaching surface as a deliberate next step, not
                as the answer. The feedback above is a verdict; this is the only
                route from it to an explanation, and it leaves the assessment
                rather than injecting instruction into it. */}
            {!lastResult.correct && !done && conceptId && (
              <Link
                href={`/learn?concept=${encodeURIComponent(conceptId)}`}
                className="mt-3 inline-flex items-center border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint px-4 py-2 font-mono text-[11px] uppercase tracking-[0.18em]"
              >
                {t('goals.question.learnThis')}
              </Link>
            )}
            {!lastResult.correct && !done && retryAvailable && (
              <button
                type="button"
                onClick={onRetryQuestion}
                disabled={loading}
                className="mt-2 font-mono text-[11px] text-mathua-muted hover:text-mathua-blue underline underline-offset-2 disabled:opacity-50"
              >
                {t('goals.question.retrySlip')}
              </button>
            )}
          </div>
        )}

        {submitError ? (
          <div className="mb-4">
            <SubmitErrorBlock
              error={submitError}
              onRetry={onRetrySubmit ?? onSubmit}
              onSkip={onSkip}
              skipsLeft={skipsLeft}
              onRestart={onRestart}
              restartLabel={t('goals.diagnostic.restart')}
              retrying={loading}
            />
          </div>
        ) : lastResult && !done ? (
          <div className="mb-4 text-center">
            <button
              type="button"
              ref={focusClearOfTabs}
              onClick={onNext}
              className="clear-of-tabs border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint rounded-none h-12 px-8 font-medium text-sm shrink-0"
            >
              {t('goals.question.next')}
            </button>
          </div>
        ) : lastResult && done && planError ? (
          <div className="mb-4 border border-mathua-red bg-mathua-surface p-4 text-left" role="alert">
            <p className="font-mono text-xs text-mathua-red">{t('goals.diagnostic.resultsErrorTitle')}</p>
            <p className="mt-1 font-mono text-[11px] text-mathua-muted break-words [overflow-wrap:anywhere]">{planError}</p>
            <div className="mt-3">
              <button
                type="button"
                onClick={onRetryPlan}
                disabled={loading}
                className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint px-4 h-10 font-mono text-xs disabled:opacity-50"
              >
                {loading ? t('common.status.loading') : t('goals.diagnostic.loadResults')}
              </button>
            </div>
          </div>
        ) : lastResult && done ? (
          <div className="mb-4 text-center text-mathua-muted text-xs font-mono">
            {t('goals.diagnostic.preparingResults')}
          </div>
        ) : null}

        <div className="text-center text-mathua-muted text-xs font-mono">
          {accuracy.total > 0 && t('goals.question.correctOf', { correct: accuracy.correct, total: accuracy.total })}
        </div>
      </div>
    </>
  )
}

export function ResultsStep({ plan, onStartPractice, demo }: { plan: GoalPlanRes; onStartPractice: () => void; demo?: boolean }) {
  const t = useT()
  return (
    <div className="max-w-3xl mx-auto min-w-0 overflow-hidden px-2 sm:px-0">
      <SectionHeader label={t('goals.onboard.resultsLabel')} title={t('goals.onboard.resultsTitle')} />
      {demo && (
        <div className="mt-4 border border-mathua-border bg-mathua-surface px-4 py-3">
          <p className="text-mathua-secondary text-sm">
            {t('goals.onboard.demoNote')}
          </p>
          <p className="mt-2">
            <Link href="/login" className="font-mono text-xs text-mathua-blue hover:underline">
              {t('goals.onboard.demoSignIn')}
            </Link>
          </p>
        </div>
      )}
      <div className="mt-6 min-w-0 overflow-hidden">
        <DiagnosticResults plan={plan} onStartPractice={onStartPractice} />
      </div>
    </div>
  )
}
