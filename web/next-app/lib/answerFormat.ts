import { tNow } from './i18n'

// Expected-answer-form microcopy for typed inputs. The server passes the
// concept's grading_type alongside every question; this table turns it into
// a one-line hint plus the right mobile keyboard. Unknown types fall back
// to a generic line, never blank, never a guess at the format.

export interface AnswerFormat {
  hint: string
  inputMode: 'numeric' | 'text'
}

// The server sends the hint and input mode the generator itself declared,
// because the grading type describes the grader and not the answer: "5 R 2",
// "2 sqrt(5)", "4:7", "(0, -2)" and "[[1,2],[3,4]]" are all numeric concepts,
// and a digits-only keypad cannot produce any of them. The type is the
// fallback for payloads that predate the declaration.
export interface QuestionFormatFields {
  grading_type?: string
  answer_hint?: string
  answer_input_mode?: string
}

export function formatForQuestion(q: QuestionFormatFields): AnswerFormat {
  const fallback = formatForGradingType(q.grading_type)
  const hint = q.answer_hint?.trim()
  const mode =
    q.answer_input_mode === 'numeric' || q.answer_input_mode === 'text'
      ? q.answer_input_mode
      : fallback.inputMode
  return { hint: hint || fallback.hint, inputMode: mode }
}

export function formatForGradingType(gradingType?: string): AnswerFormat {
  switch (gradingType) {
    case 'numeric':
      return { hint: tNow('learn.answerFormat.numeric'), inputMode: 'numeric' }
    case 'expression':
    case 'polynomial':
      return { hint: tNow('learn.answerFormat.expression'), inputMode: 'text' }
    case 'symbolic':
      return { hint: tNow('learn.answerFormat.symbolic'), inputMode: 'text' }
    case 'multiple_choice':
      return { hint: tNow('learn.answerFormat.multipleChoice'), inputMode: 'text' }
    case 'comparison':
      return { hint: tNow('learn.answerFormat.comparison'), inputMode: 'text' }
    case 'ordering':
      return { hint: tNow('learn.answerFormat.ordering'), inputMode: 'text' }
    case 'tuple':
      return { hint: tNow('learn.answerFormat.tuple'), inputMode: 'text' }
    case 'complex':
      return { hint: tNow('learn.answerFormat.complex'), inputMode: 'text' }
    default:
      return { hint: tNow('learn.answerFormat.default'), inputMode: 'text' }
  }
}
