import { describe, it, expect } from 'vitest'
import { act, render, screen } from '@testing-library/react'
import NextUpCard from '../components/NextUpCard'
import BriefingCard, { QUIZ_BRIEFING, REVIEW_BRIEFING } from '../components/BriefingCard'
import PrerequisitePanel from '../components/PrerequisitePanel'
import SubmitErrorBlock from '../components/SubmitErrorBlock'
import ChoiceOptions from '../components/ChoiceOptions'
import ReviewHost from '../app/review/ReviewHost'
import LearnDoneCard from '../components/LearnDoneCard'
import { layoutsForMode } from '../components/math/layouts'
import { formatForGradingType } from '../lib/answerFormat'
import { recommendationResToShelf } from '../lib/recommendations'
import { getTranslator, setLocale } from '../lib/i18n'
import { resetLocaleForTests } from '../lib/i18n/locale'
import type { Shelf } from '../lib/nextUp'
import type { EligibilityPrereq } from '../lib/api'

// The learn flow, read in Ukrainian. Each test pins the locale before rendering; the suite's
// setup puts it back to English afterwards.

function expectNoRawKeys(container: HTMLElement) {
  const text = container.textContent ?? ''
  expect(text).not.toMatch(/\b(learn|common)\.[a-z][A-Za-z]+\./)
}

const shelf: Shelf = {
  next: { kind: 'new', badge: 'Нова тема', title: 'Додавання дробів', detail: 'Наступний крок', href: '/learn?concept=a', cta: 'Почати', xp: 12 },
  alternatives: [
    { kind: 'review', badge: 'Повторення', title: 'Дроби', detail: '', href: '/review', cta: 'Повторити', xp: 0 },
    { kind: 'weakness', badge: 'Практика', title: 'Відсотки', detail: '', href: '/learn?concept=b', cta: 'Практика', xp: 0 },
  ],
}

describe('learn flow in Ukrainian', () => {
  it('renders the next-up card', () => {
    resetLocaleForTests('uk')
    const { container } = render(<NextUpCard shelf={shelf} />)
    expect(screen.getByRole('heading', { name: 'Що далі' })).toBeInTheDocument()
    expect(screen.getByText('Або оберіть щось інше (2)')).toBeInTheDocument()
    // Server-provided wording and XP pass through untouched.
    expect(screen.getByText('Нова тема')).toBeInTheDocument()
    expect(screen.getByText('12 XP')).toBeInTheDocument()
    expectNoRawKeys(container)
  })

  it('switches an already mounted card without a reload', () => {
    render(<NextUpCard shelf={shelf} />)
    expect(screen.getByText('Or pick something else (2)')).toBeInTheDocument()
    act(() => setLocale('uk'))
    expect(screen.getByText('Або оберіть щось інше (2)')).toBeInTheDocument()
  })

  it('renders the briefings the quiz and review pages share', () => {
    resetLocaleForTests('uk')
    const { container } = render(
      <>
        <BriefingCard eyebrow="Перш ніж почати" items={REVIEW_BRIEFING} />
        <BriefingCard eyebrow="Перш ніж почати" items={QUIZ_BRIEFING} />
      </>,
    )
    expect(screen.getByText('Лише те, що час повторити')).toBeInTheDocument()
    expect(screen.getByText('До 5 питань · на час · без підказок')).toBeInTheDocument()
    expect(screen.getByText(/Розібране розв’язання з’являється одразу/)).toBeInTheDocument()
    expectNoRawKeys(container)
  })

  it('renders the prerequisite list with Ukrainian state words', () => {
    resetLocaleForTests('uk')
    const prereqs: EligibilityPrereq[] = [
      { conceptId: 'a', title: 'Дроби', state: 'mastered', met: true },
      { conceptId: 'b', title: 'Відсотки', state: 'unlocked', met: false },
      { conceptId: 'c', title: 'Степені', state: 'due_for_review', met: true },
    ] as EligibilityPrereq[]
    const { container } = render(<PrerequisitePanel prerequisites={prereqs} onDismiss={() => {}} />)
    expect(screen.getByRole('note', { name: 'Необхідні теми' })).toBeInTheDocument()
    expect(screen.getByText('Що вивчити спочатку (3)')).toBeInTheDocument()
    expect(screen.getByText(/Засвоєна · Обов’язкова$/)).toBeInTheDocument()
    expect(screen.getByText(/Не розпочата · Обов’язкова$/)).toBeInTheDocument()
    expect(screen.getByText(/Засвоєна · час перевірити · Обов’язкова$/)).toBeInTheDocument()
    expect(screen.getAllByText('Вивчити →')).toHaveLength(3)
    expect(screen.getByRole('button', { name: 'Закрити' })).toBeInTheDocument()
    expectNoRawKeys(container)
  })

  it('renders submit errors', () => {
    resetLocaleForTests('uk')
    const { container } = render(
      <SubmitErrorBlock error={{ message: 'boom', status: 503 }} onRetry={() => {}} onSkip={() => {}} skipsLeft={0} onRestart={() => {}} restartLabel="Почати заново" />,
    )
    expect(screen.getByText('Сервер перезапускається. Спробуйте за хвилину — вашу відповідь збережено.')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Спробувати ще раз' })).toBeInTheDocument()
    expect(screen.getByText('Пропусків не залишилося — почніть заново.')).toBeInTheDocument()
    expectNoRawKeys(container)
  })

  it('labels the multiple-choice group', () => {
    resetLocaleForTests('uk')
    render(<ChoiceOptions question={'Оберіть.\n\nA) так\n\nB) ні'} value="" onPick={() => {}} />)
    expect(screen.getByRole('group', { name: 'Варіанти відповіді' })).toBeInTheDocument()
  })

  it('renders the review intro', () => {
    resetLocaleForTests('uk')
    const { container } = render(<ReviewHost />)
    expect(screen.getByText('Час повторити')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Почати повторення →' })).toBeInTheDocument()
    expectNoRawKeys(container)
  })

  it('renders the completion card', () => {
    resetLocaleForTests('uk')
    const { container } = render(
      <LearnDoneCard
        conceptId="a"
        correct={5}
        answered={6}
        xp={18}
        band="Хороший прогрес"
        back={{ href: '/study?concept=x' }}
        returnToLabel="Дроби"
        nextShelf={null}
        nextLoading={false}
        shelfProgress={{}}
        onPracticeAgain={() => {}}
      />,
    )
    expect(screen.getByText('✓ Завершено — 5/6 правильно · +18 XP · Хороший прогрес')).toBeInTheDocument()
    expect(screen.getByText('← Назад до теми «Дроби»')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Потренуватися ще раз' })).toBeInTheDocument()
    expect(screen.getByText('Назад до профілю')).toBeInTheDocument()
    expectNoRawKeys(container)
  })
})

describe('text produced outside components, in Ukrainian', () => {
  it('translates the answer-format hint but not the grading rule behind it', () => {
    resetLocaleForTests('uk')
    expect(formatForGradingType('numeric').hint).toBe('Відповідь — число, дріб або мішане число (наприклад, 4 1/10)')
    // The keypad choice is a behaviour and does not depend on the language.
    expect(formatForGradingType('numeric').inputMode).toBe('numeric')
    expect(formatForGradingType('comparison').hint).toBe('Відповідь: <, > або =')
  })

  it('translates the "all learned" fallback shelf', () => {
    resetLocaleForTests('uk')
    const s = recommendationResToShelf({ primary: null, alternatives: [], generatedAt: '' })
    expect(s.next.badge).toBe('Усе вивчено')
    expect(s.next.href).toBe('/graph')
  })

  it('translates the math keyboard captions and keeps the inserted LaTeX', () => {
    resetLocaleForTests('uk')
    const [first] = layoutsForMode('arithmetic')
    expect(first.tooltip).toBe('Числа, чотири арифметичні дії, дроби й степені')
    const asides = layoutsForMode('arithmetic').flatMap(l => l.layers.flatMap(layer => layer.rows.flat()))
      .flatMap(k => (typeof k === 'string' ? [] : [k.aside, ...((Array.isArray(k.variants) ? k.variants : []).map(v => (typeof v === 'string' ? undefined : v.aside)))]))
      .filter(Boolean)
    expect(asides).toContain('пробіл')
    expect(asides).toContain('квадратний корінь')
    expect(asides).not.toContain('square root')
    const latex = layoutsForMode('arithmetic').flatMap(l => l.layers.flatMap(layer => layer.rows.flat()))
      .flatMap(k => (typeof k === 'string' ? [] : [k.latex]))
    expect(latex).toContain('\\sqrt{#0}')
  })

  it('uses the shared plural forms for counts of concepts', () => {
    const t = getTranslator('uk')
    expect(t('common.count.concepts', { count: 1 })).toBe('1 тема')
    expect(t('common.count.concepts', { count: 3 })).toBe('3 теми')
    expect(t('common.count.concepts', { count: 5 })).toBe('5 тем')
    expect(t('common.count.concepts', { count: 11 })).toBe('11 тем')
    expect(t('common.count.concepts', { count: 21 })).toBe('21 тема')
    expect(t('common.count.concepts', { count: 22 })).toBe('22 теми')
  })

  it('declines questions, days and attempts', () => {
    const t = getTranslator('uk')
    expect(t('common.count.questions', { count: 1 })).toBe('1 питання')
    expect(t('common.count.questions', { count: 4 })).toBe('4 питання')
    expect(t('common.count.questions', { count: 12 })).toBe('12 питань')
    expect(t('common.count.days', { count: 2 })).toBe('2 дні')
    expect(t('common.count.days', { count: 25 })).toBe('25 днів')
    expect(t('common.count.attempts', { count: 3 })).toBe('3 спроби')
    expect(t('common.count.attempts', { count: 7 })).toBe('7 спроб')
  })
})
