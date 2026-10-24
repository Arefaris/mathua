'use client'

import KatexContent from './KatexContent'
import ChoiceOptions from './ChoiceOptions'
import { stemWithoutChoices } from '../lib/choices'
import { MathAnswerInput } from './math/MathInput'
import type { PracticeQuestion } from '../lib/api'
import type { AnswerFormat } from '../lib/answerFormat'
import { useT } from '../lib/i18n'

// One question in the append-only learning feed. Once answered, the card locks and shows the
// verdict and the explanation inline; scrollback is the performance history.
export interface QEntry {
  key: number
  kind: 'q'
  kpIndex: number
  q: PracticeQuestion
  answer: string
  checking: boolean
  servedAt: number
  feedback: { correct: boolean; text: string; diagnosis?: string; xp: number } | null
}

interface Props {
  entry: QEntry
  conceptId: string
  gradingType?: string
  format: AnswerFormat
  onAnswerChange: (value: string) => void
  onCheck: () => void
}

export default function LearnQuestionCard({ entry: qe, conceptId, gradingType, format, onAnswerChange, onCheck }: Props) {
  const t = useT()
  const locked = qe.feedback !== null
  return (
    <div className={`border bg-mathua-surface p-6 ${qe.feedback?.correct ? 'border-mathua-green-faint' : qe.feedback ? 'border-mathua-red-faint' : 'border-mathua-border'}`}>
      <KatexContent className="text-sm text-mathua-primary font-mono whitespace-pre-wrap">{stemWithoutChoices(qe.q.question)}</KatexContent>
      <ChoiceOptions question={qe.q.question} value={qe.answer} onPick={onAnswerChange} disabled={locked || qe.checking} />
      {!locked && (
        <form onSubmit={ev => { ev.preventDefault(); onCheck() }} className="mt-3 flex flex-col sm:flex-row gap-2">
          <label htmlFor={`learn-answer-${qe.key}`} className="sr-only">{t('learn.answer.label')}</label>
          <MathAnswerInput
            id={`learn-answer-${qe.key}`}
            value={qe.answer}
            onChange={onAnswerChange}
            onSubmit={onCheck}
            gradingType={gradingType}
            conceptId={conceptId}
            disabled={qe.checking}
            status={qe.feedback ? (qe.feedback.correct ? 'correct' : 'incorrect') : 'default'}
            inputMode={format.inputMode}
            placeholder={t('learn.answer.placeholder')}
            className="sm:flex-1"
          />
          <button type="submit" disabled={qe.checking || !qe.answer.trim()} className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint px-6 h-12 text-sm font-mono disabled:opacity-50 w-auto self-end sm:self-auto shrink-0">{t('common.action.check')}</button>
        </form>
      )}
      {!locked && <p className="mt-2 font-mono text-[11px] text-mathua-muted">{format.hint}</p>}
      {locked && qe.feedback && (
        <div className="mt-3">
          <p className={`font-mono text-xs ${qe.feedback.correct ? 'text-green-400' : 'text-red-400'}`}>
            {qe.feedback.correct ? (qe.feedback.xp ? t('learn.stepper.verdictXp', { xp: qe.feedback.xp }) : t('learn.stepper.verdictCorrect')) : t('learn.stepper.verdictWrong')}{' '}
            <span className="text-mathua-muted">{t('learn.stepper.youAnswered', { answer: qe.answer })}</span>
          </p>
          <p className="mt-1.5 font-mono text-xs text-mathua-primary">
            {t('learn.stepper.answerLabel')}<KatexContent>{qe.q.answer}</KatexContent>
          </p>
          {/* The explanation is the solution for *this* instance, shown
              on both verdicts: a correct answer still needs the
              reasoning that made it right. Guarded on a non-empty
              string so a server that returns nothing leaves no empty
              section behind. */}
          {/* The mistake, when it can be named with certainty. It is
              a description of what happened, never a hint about what
              to do, and it is empty far more often than not. */}
          {qe.feedback.diagnosis && (
            <p data-diagnosis className="mt-2 font-mono text-[11px] text-mathua-muted">{qe.feedback.diagnosis}</p>
          )}
          {qe.feedback.text && (
            <div className="mt-1.5">
              <p className="font-mono text-[11px] uppercase tracking-wider text-mathua-muted">
                {qe.feedback.correct ? t('learn.stepper.whyCorrect') : t('learn.stepper.howWrong')}
              </p>
              <KatexContent className="mt-1 text-xs font-mono text-mathua-secondary whitespace-pre-wrap">{qe.feedback.text}</KatexContent>
            </div>
          )}
        </div>
      )}
    </div>
  )
}
