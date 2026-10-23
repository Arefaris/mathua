// Parse the multiple-choice options a generator embeds in the question stem.
//
// Most `multiple_choice` generators emit options as trailing "A) …", "B) …"
// blocks (see generator/precalculus mcProblem). Rendering them as buttons lets
// students answer by tapping instead of typing the exact option text into a
// free-text box. Returns null when the question has no parseable option list,
// so every caller falls back to the text input unchanged.

export interface ChoiceOption {
  letter: string
  text: string
}

const OPTION_RE = /\n\n([A-E])\)\s+([\s\S]*?)(?=\n\n[A-E]\)\s|$)/g

export function parseChoices(question: string): ChoiceOption[] | null {
  if (!question) return null
  const opts: ChoiceOption[] = []
  let m: RegExpExecArray | null
  OPTION_RE.lastIndex = 0
  while ((m = OPTION_RE.exec(question)) !== null) {
    const text = m[2].trim()
    if (text) opts.push({ letter: m[1], text })
  }
  if (opts.length < 2) return null
  // Letters must run A, B, C, … so we never mistake prose for options.
  for (let i = 0; i < opts.length; i++) {
    if (opts[i].letter !== String.fromCharCode(65 + i)) return null
  }
  return opts
}

// The stem without the option block, for display above tappable options.
//
// The generator writes the options into the question text, and the options are
// also drawn as buttons, so a learner reads every choice twice: once as prose
// and once as a tap target (report #232). Only the rendered stem loses the
// list. The full question is still passed to the grader and to the report
// button, because the letter is resolved against the text the server served.
export function stemWithoutChoices(question: string): string {
  if (!parseChoices(question)) return question
  const cut = question.search(/\n\n[A-E]\)\s/)
  return cut === -1 ? question : question.slice(0, cut).trimEnd()
}
