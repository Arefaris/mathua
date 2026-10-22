'use client'

import dynamic from 'next/dynamic'
import Image from 'next/image'
import {
  bodyStyle,
  codeBlockStyle,
  h2Style,
  inlineCodeStyle,
  monoFont,
  mutedCodeStyle,
  roman,
  stepLiStyle,
  tableCellStyle,
  tableHeaderStyle,
  bodyFont,
} from './styles'
import { CONCEPT_COUNT, DOMAIN_COUNT } from '../../lib/curriculumFacts'
import { useT } from '../../lib/i18n'
import SlotText from '../../components/explore/SlotText'

// Shown while a diagram's chunk loads. A component rather than inline JSX so it can read the
// language: `dynamic`'s `loading` option runs outside the page's own render.
function FlowPlaceholder({ height }: { height: number }) {
  const t = useT()
  return (
    <div style={{ height, border: '0.5px solid var(--border)', display: 'flex', alignItems: 'center', justifyContent: 'center', color: 'var(--text-muted)', fontFamily: monoFont, fontSize: '13px' }}>
      {t('common.status.loading')}
    </div>
  )
}

const StudentModelFlow = dynamic(() => import('../../components/StudentModelFlow'), {
  ssr: false,
  loading: () => <FlowPlaceholder height={320} />,
})

const DiagnosticFlow = dynamic(() => import('../../components/DiagnosticFlow'), {
  ssr: false,
  loading: () => <FlowPlaceholder height={420} />,
})

const SchedulerFlow = dynamic(() => import('../../components/SchedulerFlow'), {
  ssr: false,
  loading: () => <FlowPlaceholder height={320} />,
})

const GradingFlow = dynamic(() => import('../../components/GradingFlow'), {
  ssr: false,
  loading: () => <FlowPlaceholder height={300} />,
})

export function ConceptGraphSection() {
  const t = useT()
  return (
    <section id="concept-graph" className="mt-12 pt-6 scroll-mt-28">
      <h2 style={h2Style}>{t('explore.howItWorks.conceptGraph.title')}</h2>
      <p style={bodyStyle}>
        {t('explore.howItWorks.conceptGraph.p1')}
      </p>
      <p style={bodyStyle}>
        {t('explore.howItWorks.conceptGraph.p2', {
          concepts: t('explore.howItWorks.conceptGraph.concepts', { count: CONCEPT_COUNT }),
          domains: t('explore.howItWorks.conceptGraph.domains', { count: DOMAIN_COUNT }),
        })}
      </p>
      <pre style={codeBlockStyle}>
{`{
  "id": "frac.add.diff",
  "label": "Add fractions with different denominators",
  "domain": "fractions",
  "subdomain": "fractions.addition",
  "prerequisites": ["frac.add.same", "arith.factor.lcm"],
  "mastery_threshold": { "streak": 5, "avg_time_seconds": 18.0 }
}`}
      </pre>
      <p style={bodyStyle}>
        <SlotText
          text={t('explore.howItWorks.conceptGraph.p3')}
          slots={{ file: <code style={inlineCodeStyle}>data/concepts.json</code> }}
        />
      </p>
      <div className="flex items-center gap-2 flex-wrap my-6">
        <span style={{ fontFamily: monoFont, fontSize: '13px', color: 'var(--accent-blue)' }}>arith.factor.gcf</span>
        <span style={{ color: 'var(--border-strong)', fontFamily: monoFont }}>→</span>
        <span style={{ fontFamily: monoFont, fontSize: '13px', color: 'var(--accent-blue)' }}>arith.factor.lcm</span>
        <span style={{ color: 'var(--border-strong)', fontFamily: monoFont }}>→</span>
        <span style={{ fontFamily: monoFont, fontSize: '13px', color: 'var(--accent-blue)' }}>frac.add.diff</span>
        <span style={{ color: 'var(--border-strong)', fontFamily: monoFont }}>→</span>
        <span style={{ fontFamily: monoFont, fontSize: '13px', color: 'var(--accent-blue)' }}>frac.mixed.add</span>
      </div>
    </section>
  )
}

// Skill key, correct answers in a row, seconds per question. The names are translated.
const STUDENT_ROWS = [
  ['singleDigit', 5, 8],
  ['tables', 7, 6],
  ['fractions', 5, 18],
  ['equations', 5, 12],
  ['division', 5, 20],
] as const

export function StudentModelSection() {
  const t = useT()
  return (
    <section id="student-model" className="mt-12 pt-6 scroll-mt-28">
      <h2 style={h2Style}>{t('explore.howItWorks.studentModel.title')}</h2>
      <p style={bodyStyle}>
        {t('explore.howItWorks.studentModel.intro')}
      </p>

      <StudentModelFlow />

      {/* Mobile: stacked labeled cards */}
      <div className="sm:hidden space-y-2 w-full max-w-[720px] mx-auto min-w-0 my-6">
        {STUDENT_ROWS.map(([skill, streak, seconds]) => (
          <div key={`card-${skill}`} className="border-[0.5px] border-mathua-border p-3 min-w-0 overflow-hidden text-center">
            <div className="font-mono text-[12px] text-mathua-primary mb-1 break-words">{t(`explore.howItWorks.studentModel.rows.${skill}`)}</div>
            <div className="flex flex-wrap justify-center gap-x-4 gap-y-0.5 font-mono text-[11px] text-mathua-muted">
              <span>{t('explore.howItWorks.studentModel.streak', { count: streak })}</span>
              <span>{t('explore.howItWorks.studentModel.seconds', { count: seconds })}</span>
            </div>
          </div>
        ))}
      </div>

      {/* Desktop: real table */}
      <div className="hidden sm:block overflow-x-auto overscroll-x-contain w-[calc(100%+2rem)] -mx-4 px-4 sm:w-full sm:mx-0 sm:px-0 my-6">
        <table className="w-full"
          style={{
            width: '100%',
            borderCollapse: 'collapse',
            border: 'none',
            background: 'transparent',
          }}
        >
          <thead>
            <tr style={{ background: 'transparent' }}>
              <th style={tableHeaderStyle}>{t('explore.howItWorks.studentModel.skill')}</th>
              <th style={tableHeaderStyle}>{t('explore.howItWorks.studentModel.streakHeader')}</th>
              <th style={tableHeaderStyle}>{t('explore.howItWorks.studentModel.timeHeader')}</th>
            </tr>
          </thead>
          <tbody>
            {STUDENT_ROWS.map(([skill, streak, seconds]) => (
              <tr key={`row-${skill}`} style={{ background: 'transparent' }}>
                <td style={tableCellStyle}>{t(`explore.howItWorks.studentModel.rows.${skill}`)}</td>
                <td style={tableCellStyle}>{streak}</td>
                <td style={tableCellStyle}>{t('explore.howItWorks.studentModel.seconds', { count: seconds })}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </section>
  )
}

export function SpacedRepetitionSection() {
  const t = useT()
  return (
    <section id="spaced-repetition" className="mt-12 pt-6 scroll-mt-28">
      <h2 style={h2Style}>{t('explore.howItWorks.spaced.title')}</h2>
      <p style={bodyStyle}>
        {t('explore.howItWorks.spaced.p1')}
      </p>
      <pre style={codeBlockStyle}>
{t('explore.howItWorks.spaced.schedule')}
      </pre>
      <p style={bodyStyle}>
        {t('explore.howItWorks.spaced.p2')}
      </p>
    </section>
  )
}

const DIAGNOSTIC_STEPS = ['sort', 'forward', 'backward', 'probes', 'record'] as const

export function DiagnosticSection() {
  const t = useT()
  return (
    <section id="diagnostic" className="mt-12 pt-6 scroll-mt-28">
      <h2 style={h2Style}>{t('explore.howItWorks.diagnostic.title')}</h2>
      <p style={bodyStyle}>
        {t('explore.howItWorks.diagnostic.p1')}
      </p>
      <ol className="list-none my-4">
        {DIAGNOSTIC_STEPS.map((step, i) => (
          <li key={step} style={stepLiStyle}>
            <span
              style={{
                position: 'absolute',
                left: 0,
                fontFamily: monoFont,
                fontSize: '13px',
                color: 'var(--accent-blue)',
              }}
            >
              {roman(i + 1)}.
            </span>
            {t(`explore.howItWorks.diagnostic.steps.${step}`)}
          </li>
        ))}
      </ol>
      <DiagnosticFlow />
      <p style={bodyStyle}>
        {t('explore.howItWorks.diagnostic.p2', {
          concepts: t('explore.howItWorks.diagnostic.conceptsOf', { count: CONCEPT_COUNT }),
          questions: t('explore.howItWorks.diagnostic.questionsUpTo', { count: CONCEPT_COUNT }),
        })}
      </p>
      <p style={bodyStyle}>
        {t('explore.howItWorks.diagnostic.p3')}
      </p>
      <style>{`ol { counter-reset: step-counter; }`}</style>
    </section>
  )
}

const SCHEDULER_RULES = ['prerequisites', 'repeat', 'mix'] as const

export function SchedulerSection() {
  const t = useT()
  return (
    <section id="scheduler" className="mt-12 pt-6 scroll-mt-28">
      <h2 style={h2Style}>{t('explore.howItWorks.scheduler.title')}</h2>
      <p style={bodyStyle}>
        {t('explore.howItWorks.scheduler.p1')}
      </p>
      <pre style={codeBlockStyle}>
{`priority = (0.7 × days_since_last_seen)
+ (0.3 × (1.0 − mastery_score))
+ 5.0   if status == DECAYING
+ 2.0   if newly unlocked this session`}
      </pre>
      <p style={bodyStyle}>
        {t('explore.howItWorks.scheduler.p2')}
      </p>
      {SCHEDULER_RULES.map((rule) => (
        <div
          key={rule}
          style={{
            fontFamily: bodyFont,
            fontSize: '0.95rem',
            color: 'var(--text-secondary)',
            lineHeight: 1.7,
            borderLeft: '2px solid var(--accent-blue)',
            paddingLeft: '1rem',
            margin: '0.5rem 0',
          }}
        >
          {t(`explore.howItWorks.scheduler.rules.${rule}`)}
        </div>
      ))}
      <SchedulerFlow />
    </section>
  )
}

export function ScoringSection() {
  const t = useT()
  return (
    <section id="scoring" className="mt-12 pt-6 scroll-mt-28">
      <h2 style={h2Style}>{t('explore.howItWorks.scoring.title')}</h2>
      <p style={bodyStyle}>
        {t('explore.howItWorks.scoring.p1')}
      </p>
      <pre style={codeBlockStyle}>
{`topic_score = (mastered_concepts × 100)
+ Σ speed_bonus per mastered concept
+ domain_completion_bonus (200 pts if all concepts mastered)
speed_bonus = max(0, (time_limit − avg_time) / time_limit × 50)`}
      </pre>
      <pre style={codeBlockStyle}>
{`weekly_score = (concepts_mastered_this_week × 100)
+ speed_bonus_this_week
+ (current_day_streak × 10)`}
      </pre>
    </section>
  )
}

export function GeneratorsSection() {
  const t = useT()
  return (
    <section id="generators" className="mt-12 pt-6 scroll-mt-28">
      <h2 style={h2Style}>{t('explore.howItWorks.generators.title')}</h2>
      <p style={bodyStyle}>
        {t('explore.howItWorks.generators.p1')}
      </p>
      <pre style={{ ...codeBlockStyle, whiteSpace: 'pre', overflowX: 'auto' }}>
{`// internal/generator/arithmetic/add.go
type AddSingleGen struct{}

func (g *AddSingleGen) Generate(difficulty float64) generator.Problem {
    max := int(5 + difficulty*4)
    a   := rand.Intn(max) + 2
    b   := rand.Intn(max) + 2
    ans := a + b
    return generator.Problem{
        Question:    fmt.Sprintf("%d + %d = ?", a, b),
        Answer:      strconv.Itoa(ans),
        Explanation: fmt.Sprintf("%d + %d = %d", a, b, ans),
    }
}`}
      </pre>
      <p style={bodyStyle}>
        <SlotText
          text={t('explore.howItWorks.generators.p2')}
          slots={{ code: <code style={inlineCodeStyle}>difficulty</code> }}
        />
      </p>
    </section>
  )
}

export function SymbolicGradingSection() {
  const t = useT()
  return (
    <section id="symbolic-grading" className="mt-12 pt-6 scroll-mt-28">
      <h2 style={h2Style}>{t('explore.howItWorks.grading.title')}</h2>
      <p style={bodyStyle}>
        <SlotText
          text={t('explore.howItWorks.grading.p1')}
          slots={{
            examples: (
              <>
                <code style={mutedCodeStyle}>x = 4</code>,{' '}
                <code style={mutedCodeStyle}>(x+2)(x+3)</code>,{' '}
                <code style={mutedCodeStyle}>2x&#178; + 3x - 5</code>
              </>
            ),
            sympy: <code style={inlineCodeStyle}>sympy</code>,
          }}
        />
      </p>
      <p style={bodyStyle}>
        <SlotText
          text={t('explore.howItWorks.grading.p2')}
          slots={{
            router: <code style={inlineCodeStyle}>Router</code>,
            gradingType: <code style={mutedCodeStyle}>grading_type</code>,
            symbolic: <code style={mutedCodeStyle}>symbolic</code>,
            polynomial: <code style={mutedCodeStyle}>polynomial</code>,
            expression: <code style={mutedCodeStyle}>expression</code>,
            gradeSymPy: <code style={inlineCodeStyle}>gradeSymPy</code>,
          }}
        />
      </p>
      <GradingFlow />
      <pre style={{ ...codeBlockStyle, whiteSpace: 'pre', overflowX: 'auto' }}>
{`// internal/grader/router.go
func (r *Router) Grade(t GradingType, expected, answer string) Result {
    switch t {
    case GradingNumeric:
        return r.numeric.grade(expected, answer)
    case GradingPolynomial, GradingExpression, GradingSymbolic:
        return gradeSymPy(expected, answer) // Python subprocess
    case GradingMultipleChoice:
        return r.choice.grade(expected, answer)
    case GradingComparison:
        return r.comparison.grade(expected, answer)
    case GradingOrdering:
        return r.ordering.grade(expected, answer)
    case GradingTuple:
        return r.tuple.grade(expected, answer)
    case GradingComplex:
        return r.complex.grade(expected, answer)
    }
}

// internal/engine/engine.go — a generator may override the router:
func (e *Engine) gradeAnswer(conceptID, expected, answer string) grader.Result {
    if gen, err := e.registry.Get(conceptID); err == nil {
        if gg, ok := gen.(generator.GradedGenerator); ok {
            return gg.Grade(expected, answer)
        }
    }
    return e.gr.Grade(gradingTypeOf(conceptID), expected, answer)
}`}
      </pre>
      <p style={bodyStyle}>
        <SlotText
          text={t('explore.howItWorks.grading.p3')}
          slots={{ check: <code style={inlineCodeStyle}>simplify(expected - answer) == 0</code> }}
        />
      </p>
      <pre style={codeBlockStyle}>
{`# grading/sympy_service.py
import sys, json
from sympy import simplify, parse_expr, Symbol

for line in sys.stdin:
    req = json.loads(line)
    expected = parse_expr(req["expected"])
    answer   = parse_expr(req["answer"])
    correct  = simplify(expected - answer) == 0
    print(json.dumps({"id": req["id"], "correct": correct}))`}
      </pre>
      <p style={bodyStyle}>
        {t('explore.howItWorks.grading.p4')}
        <span style={{ display: 'block', fontFamily: monoFont, fontSize: '13px', marginTop: '0.5rem', color: 'var(--text-muted)' }}>
          x² + 2x + 1 == (x+1)²{' · '}
          sin²(x) + cos²(x) == 1{' · '}
          xe^x − e^x + C == e^x(x−1) + C{' · '}
          2e^(2x) == 2exp(2x)
        </span>
      </p>
    </section>
  )
}

const METRIC_ROWS = ['firstPass', 'secondPass', 'attempts', 'retention', 'trend'] as const

export function MeasuringLearningSection() {
  const t = useT()
  return (
    <section id="measuring-learning" className="mt-12 pt-6 scroll-mt-28">
      <h2 style={h2Style}>{t('explore.howItWorks.measuring.title')}</h2>
      <p style={bodyStyle}>
        {t('explore.howItWorks.measuring.p1')}
      </p>

      {/* Mobile: stacked labeled cards */}
      <div className="sm:hidden space-y-2 w-full max-w-[720px] mx-auto min-w-0 my-6">
        {METRIC_ROWS.map((metric) => (
          <div key={`card-${metric}`} className="border-[0.5px] border-mathua-border p-3 min-w-0 overflow-hidden">
            <div className="font-mono text-[12px] text-mathua-primary mb-1 break-words">{t(`explore.howItWorks.measuring.rows.${metric}.name`)}</div>
            <div className="font-mono text-[11px] text-mathua-muted">{t(`explore.howItWorks.measuring.rows.${metric}.text`)}</div>
          </div>
        ))}
      </div>

      {/* Desktop: real table */}
      <div className="hidden sm:block overflow-x-auto overscroll-x-contain w-[calc(100%+2rem)] -mx-4 px-4 sm:w-full sm:mx-0 sm:px-0 my-6">
        <table className="w-full"
          style={{
            width: '100%',
            borderCollapse: 'collapse',
            border: 'none',
            background: 'transparent',
          }}
        >
          <thead>
            <tr style={{ background: 'transparent' }}>
              <th style={tableHeaderStyle}>{t('explore.howItWorks.measuring.metric')}</th>
              <th style={tableHeaderStyle}>{t('explore.howItWorks.measuring.meaning')}</th>
            </tr>
          </thead>
          <tbody>
            {METRIC_ROWS.map((metric) => (
              <tr key={`row-${metric}`} style={{ background: 'transparent' }}>
                <td style={tableCellStyle}>{t(`explore.howItWorks.measuring.rows.${metric}.name`)}</td>
                <td style={tableCellStyle}>{t(`explore.howItWorks.measuring.rows.${metric}.text`)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <div style={{
        border: '0.5px solid var(--border)',
        padding: '1rem',
        background: 'var(--surface)',
      }}>
        <Image
          src="/figures/efficacy-overview.svg"
          alt={t('explore.howItWorks.measuring.figureAlt')}
          width={1198}
          height={778}
          style={{ maxWidth: '100%', height: 'auto', display: 'block' }}
        />
      </div>
      <p style={bodyStyle}>
        {t('explore.howItWorks.measuring.p2')}
      </p>
      <p style={{ fontFamily: monoFont, fontSize: '12px', color: 'var(--text-muted)', lineHeight: 1.7 }}>
        <SlotText
          text={t('explore.howItWorks.measuring.p3')}
          slots={{
            profile: <a href="/profile#efficacy" className="link-underline" style={{ color: 'var(--accent-blue)' }}>{t('explore.howItWorks.measuring.profileLink')}</a>,
            docs: <a href="/docs/efficacy" className="link-underline" style={{ color: 'var(--accent-blue)' }}>{t('explore.howItWorks.measuring.docsLink')}</a>,
          }}
        />
      </p>
    </section>
  )
}
