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

// The learn flow, read in Russian. Each test pins the locale before rendering; the suite's setup
// puts it back to English afterwards.

function expectNoRawKeys(container: HTMLElement) {
  const text = container.textContent ?? ''
  expect(text).not.toMatch(/\b(learn|common)\.[a-z][A-Za-z]+\./)
}

const shelf: Shelf = {
  next: { kind: 'new', badge: 'Новая тема', title: 'Сложение дробей', detail: 'Следующий шаг', href: '/learn?concept=a', cta: 'Начать', xp: 12 },
  alternatives: [
    { kind: 'review', badge: 'Повторение', title: 'Дроби', detail: '', href: '/review', cta: 'Повторить', xp: 0 },
    { kind: 'weakness', badge: 'Практика', title: 'Проценты', detail: '', href: '/learn?concept=b', cta: 'Практика', xp: 0 },
  ],
}

describe('learn flow in Russian', () => {
  it('renders the next-up card', () => {
    resetLocaleForTests('ru')
    const { container } = render(<NextUpCard shelf={shelf} />)
    expect(screen.getByRole('heading', { name: 'Что дальше' })).toBeInTheDocument()
    expect(screen.getByText('Или выберите другое (2)')).toBeInTheDocument()
    // Server-provided wording and XP pass through untouched.
    expect(screen.getByText('Новая тема')).toBeInTheDocument()
    expect(screen.getByText('12 XP')).toBeInTheDocument()
    expectNoRawKeys(container)
  })

  it('switches an already mounted card without a reload', () => {
    render(<NextUpCard shelf={shelf} />)
    expect(screen.getByText('Or pick something else (2)')).toBeInTheDocument()
    act(() => setLocale('ru'))
    expect(screen.getByText('Или выберите другое (2)')).toBeInTheDocument()
  })

  it('renders the briefings the quiz and review pages share', () => {
    resetLocaleForTests('ru')
    const { container } = render(
      <>
        <BriefingCard eyebrow="Перед началом" items={REVIEW_BRIEFING} />
        <BriefingCard eyebrow="Перед началом" items={QUIZ_BRIEFING} />
      </>,
    )
    expect(screen.getByText('Только то, что пора повторить')).toBeInTheDocument()
    expect(screen.getByText('До 5 вопросов · на время · без подсказок')).toBeInTheDocument()
    expect(screen.getByText(/Разобранное решение появится сразу/)).toBeInTheDocument()
    expectNoRawKeys(container)
  })

  it('renders the prerequisite list with Russian state words', () => {
    resetLocaleForTests('ru')
    const prereqs: EligibilityPrereq[] = [
      { conceptId: 'a', title: 'Дроби', state: 'mastered', met: true },
      { conceptId: 'b', title: 'Проценты', state: 'unlocked', met: false },
      { conceptId: 'c', title: 'Степени', state: 'due_for_review', met: true },
    ] as EligibilityPrereq[]
    const { container } = render(<PrerequisitePanel prerequisites={prereqs} onDismiss={() => {}} />)
    expect(screen.getByRole('note', { name: 'Опорные темы' })).toBeInTheDocument()
    expect(screen.getByText('Опорные темы (3)')).toBeInTheDocument()
    expect(screen.getByText(/Освоена · Обязательная опорная тема/)).toBeInTheDocument()
    expect(screen.getByText(/Не начата · Обязательная опорная тема/)).toBeInTheDocument()
    expect(screen.getByText(/Освоена · пора проверить · Обязательная опорная тема/)).toBeInTheDocument()
    expect(screen.getAllByText('Изучить →')).toHaveLength(3)
    expect(screen.getByRole('button', { name: 'Закрыть' })).toBeInTheDocument()
    expectNoRawKeys(container)
  })

  it('renders submit errors', () => {
    resetLocaleForTests('ru')
    const { container } = render(
      <SubmitErrorBlock error={{ message: 'boom', status: 500 }} onRetry={() => {}} onSkip={() => {}} skipsLeft={0} onRestart={() => {}} restartLabel="Начать заново" />,
    )
    expect(screen.getByText('Не удалось отправить ответ — он остался в поле.')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Попробовать ещё раз' })).toBeInTheDocument()
    expect(screen.getByText('Пропусков не осталось — начните заново.')).toBeInTheDocument()
    expectNoRawKeys(container)
  })

  it('labels the multiple-choice group', () => {
    resetLocaleForTests('ru')
    render(<ChoiceOptions question={'Выберите.\n\nA) да\n\nB) нет'} value="" onPick={() => {}} />)
    expect(screen.getByRole('group', { name: 'Варианты ответа' })).toBeInTheDocument()
  })

  it('renders the review intro', () => {
    resetLocaleForTests('ru')
    const { container } = render(<ReviewHost />)
    expect(screen.getByText('Пора повторить')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Начать повторение →' })).toBeInTheDocument()
    expectNoRawKeys(container)
  })

  it('renders the completion card', () => {
    resetLocaleForTests('ru')
    const { container } = render(
      <LearnDoneCard
        conceptId="a"
        correct={5}
        answered={6}
        xp={18}
        band="Хороший прогресс"
        back={{ href: '/study?concept=x' }}
        returnToLabel="Дроби"
        nextShelf={null}
        nextLoading={false}
        shelfProgress={{}}
        onPracticeAgain={() => {}}
      />,
    )
    expect(screen.getByText('✓ Готово — 5/6 верно · +18 XP · Хороший прогресс')).toBeInTheDocument()
    expect(screen.getByText('← Назад: Дроби')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Тренироваться снова' })).toBeInTheDocument()
    expect(screen.getByText('Назад в профиль')).toBeInTheDocument()
    expectNoRawKeys(container)
  })
})

describe('text produced outside components', () => {
  it('translates the answer-format hint but not the grading rule behind it', () => {
    resetLocaleForTests('ru')
    expect(formatForGradingType('numeric').hint).toBe('Ответьте числом, дробью или смешанным числом (например, 4 1/10)')
    // The keypad choice is a behaviour and does not depend on the language.
    expect(formatForGradingType('numeric').inputMode).toBe('numeric')
    expect(formatForGradingType('comparison').hint).toBe('Ответьте: <, > или =')
  })

  it('translates the "all learned" fallback shelf', () => {
    resetLocaleForTests('ru')
    const s = recommendationResToShelf({ primary: null, alternatives: [], generatedAt: '' })
    expect(s.next.badge).toBe('Всё изучено')
    expect(s.next.href).toBe('/graph')
  })

  it('translates the math keyboard captions and keeps the inserted LaTeX', () => {
    resetLocaleForTests('ru')
    const [first] = layoutsForMode('arithmetic')
    expect(first.tooltip).toBe('Числа, четыре действия, дроби и степени')
    const asides = layoutsForMode('arithmetic').flatMap(l => l.layers.flatMap(layer => layer.rows.flat()))
      .flatMap(k => (typeof k === 'string' ? [] : [k.aside, ...((Array.isArray(k.variants) ? k.variants : []).map(v => (typeof v === 'string' ? undefined : v.aside)))]))
      .filter(Boolean)
    expect(asides).toContain('пробел')
    expect(asides).toContain('квадратный корень')
    expect(asides).not.toContain('square root')
    const latex = layoutsForMode('arithmetic').flatMap(l => l.layers.flatMap(layer => layer.rows.flat()))
      .flatMap(k => (typeof k === 'string' ? [] : [k.latex]))
    expect(latex).toContain('\\sqrt{#0}')
  })

  it('uses the shared plural forms for counts of concepts', () => {
    const t = getTranslator('ru')
    expect(t('common.count.concepts', { count: 1 })).toBe('1 тема')
    expect(t('common.count.concepts', { count: 3 })).toBe('3 темы')
    expect(t('common.count.concepts', { count: 5 })).toBe('5 тем')
    expect(t('common.count.concepts', { count: 21 })).toBe('21 тема')
  })
})
