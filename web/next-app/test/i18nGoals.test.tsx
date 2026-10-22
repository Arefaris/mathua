import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import QuizHost from '../app/goals/QuizHost'
import OnboardPage from '../app/onboard/page'
import PlanEditor from '../components/PlanEditor'
import DiagnosticResults from '../components/DiagnosticResults'
import DailyGoalControl from '../components/DailyGoalControl'
import MonthlyCards from '../components/MonthlyCards'
import { startQuizSession, getDestinations, getEstimate, getScores, getCurrentPlan } from '../lib/api'
import type { GoalPlanRes } from '../lib/api'
import { resetLocaleForTests } from '../lib/i18n/locale'
import { getTranslator } from '../lib/i18n'
import { monthLabel } from '../lib/plan'

// The goals, diagnostic, onboarding and planner screens, read in Russian. The English suites
// cover the same screens; this one pins that the Russian text arrives, that counts take the
// right plural form, and that no raw key id reaches the page.

vi.mock('../lib/api', async (importOriginal) => {
  const mod = await importOriginal<typeof import('../lib/api')>()
  return {
    ...mod,
    startQuizSession: vi.fn(),
    submitQuizAnswer: vi.fn(),
    skipQuizQuestion: vi.fn(),
    startGoalDiagnostic: vi.fn(),
    getGoalPlan: vi.fn(),
    getDestinations: vi.fn(),
    getEstimate: vi.fn(),
    getScores: vi.fn(),
    getCurrentPlan: vi.fn(),
    savePlan: vi.fn(),
  }
})

vi.mock('next/navigation', () => ({
  useRouter: () => ({ push: vi.fn(), replace: vi.fn() }),
}))
vi.mock('next/link', () => ({
  default: ({ href, children, ...rest }: { href: string; children: React.ReactNode }) => (
    <a href={href} {...rest}>{children}</a>
  ),
}))
vi.mock('sonner', () => ({ toast: { error: vi.fn(), success: vi.fn() } }))
vi.mock('../hooks/useTheme', () => ({ useTheme: () => ({ mounted: true }) }))
vi.mock('../components/Header', () => ({ default: () => null }))
vi.mock('../components/Footer', () => ({ default: () => null }))
vi.mock('../components/BottomTabs', () => ({ default: () => null }))
vi.mock('../components/Loading', () => ({ default: () => <div>Загрузка</div> }))
vi.mock('../components/KatexContent', () => ({ default: ({ children }: { children: React.ReactNode }) => <span>{children}</span> }))
vi.mock('../components/math/MathInput', () => ({ MathAnswerInput: () => <input /> }))
vi.mock('../lib/auth', async (importOriginal) => {
  const mod = await importOriginal<typeof import('../lib/auth')>()
  return { ...mod, getUserInfo: vi.fn(() => ({ student_id: 's1' })), getGuestId: vi.fn(() => '') }
})

// The key ids look like `goals.quiz.start`. If a lookup misses, the id is what renders.
function expectNoKeyLeak() {
  expect(document.body.textContent).not.toMatch(/\b(goals|common)\.[a-zA-Z]+\.[a-zA-Z]/)
}

beforeEach(() => {
  resetLocaleForTests('ru')
  sessionStorage.clear()
})

describe('plural forms in the goals dictionary', () => {
  it.each([
    [1, 'Пройти диагностику (в охвате 1 тема)'],
    [2, 'Пройти диагностику (в охвате 2 темы)'],
    [5, 'Пройти диагностику (в охвате 5 тем)'],
    [11, 'Пройти диагностику (в охвате 11 тем)'],
    [21, 'Пройти диагностику (в охвате 21 тема)'],
    [22, 'Пройти диагностику (в охвате 22 темы)'],
  ])('scope of %i concepts', (count, expected) => {
    expect(getTranslator('ru')('goals.onboard.start', { count })).toBe(expected)
  })

  it('counts days, ahead and behind with the right form', () => {
    const t = getTranslator('ru')
    expect(t('goals.planner.daysApprox', { count: 3 })).toBe('≈ 3 дня')
    expect(t('goals.planner.daysApprox', { count: 12 })).toBe('≈ 12 дней')
    expect(t('goals.planner.ahead', { count: 21 })).toBe('Опережение сохранённого плана: ≈21 день')
    expect(t('goals.planner.behind', { count: 4 })).toBe('Отставание от сохранённого плана: ≈4 дня')
  })
})

describe('/onboard in Russian', () => {
  it('shows the welcome step with Russian domain names and a Russian start button', async () => {
    render(<OnboardPage />)
    await waitFor(() => expect(screen.getByText('Арифметика')).toBeTruthy())
    expect(screen.getByText('Найдём, с чего вам начать')).toBeTruthy()
    expect(screen.getByRole('button', { name: 'Выбрать всё' })).toBeTruthy()
    // Per-domain concept counts use the plural form, not "N concepts".
    expect(screen.queryByText(/\d+ concepts/)).toBeNull()
    fireEvent.click(screen.getByText('Арифметика'))
    const start = screen.getByRole('button', { name: /^Пройти диагностику \(в охвате \d+ тем/ })
    expect(start).not.toBeDisabled()
    expectNoKeyLeak()
  })
})

describe('QuizHost in Russian', () => {
  it('shows the intro and a Russian failure screen', async () => {
    vi.mocked(startQuizSession).mockRejectedValue(new Error('network down'))
    render(<QuizHost />)
    expect(screen.getByText('Проверка освоения')).toBeTruthy()
    fireEvent.click(screen.getByRole('button', { name: /Начать контрольную/ }))
    const alert = await screen.findByRole('alert')
    expect(alert.textContent).toContain('Контрольная не запустилась')
    expect(screen.getByRole('button', { name: /Попробовать ещё раз/ })).toBeTruthy()
    expectNoKeyLeak()
  })
})

describe('DiagnosticResults in Russian', () => {
  const plan: GoalPlanRes = {
    readiness: 0.7,
    total_tested: 20,
    correct_count: 14,
    weak_areas: { fractions: [{ id: 'frac.add', label: 'Сложение дробей' }] },
    strong_areas: { arithmetic: ['arith.add'] },
    frontier_label: 'Сложение дробей',
    frontier_conditional: true,
    conditionally_completed: ['a', 'b', 'c'],
  }

  it('renders the banner, the score and plural topic counts', () => {
    render(<DiagnosticResults plan={plan} onStartPractice={() => {}} />)
    expect(screen.getByText(/Начните с темы: Сложение дробей/)).toBeTruthy()
    expect(screen.getByText(/предварительно/)).toBeTruthy()
    expect(screen.getByText(/Предварительно: 3 темы считаются известными/)).toBeTruthy()
    expect(screen.getByText(/результат диагностики \(14\/20 верно\)/)).toBeTruthy()
    expect(screen.getByText('Дроби')).toBeTruthy()
    expect(screen.getByText('повторить: 1 · уверенно: 0')).toBeTruthy()
    expect(screen.getByRole('button', { name: /Приступить к практике/ })).toBeTruthy()
    expectNoKeyLeak()
  })
})

describe('PlanEditor in Russian', () => {
  const dests = [{ id: 'alg', name: 'Алгебра', description: '', total: 10, mastered: 2, pct: 0.2 }]
  const estimate = {
    destination: 'alg',
    estimate: {
      total: 5, mastered: 1, remaining: ['a', 'b'], xp_remaining: 1234, time_min_remaining: 125,
      reviews_due: 1, quizzes_ahead: 1, assessment_min: 10, diagnostic_min: 5, days: 12,
      finish_date: '2026-06-01', required_per_day: 10, feasible: true,
      lines: { learning: 80, assessment: 15, diagnostic: 5 },
    },
    pace: { rate: 10, source: 'measured', adherence: 1, active_days: 5, trailing_days: 7 },
    probes: ['a', 'b'],
    plan_delta_days: 2,
  }

  beforeEach(() => {
    vi.mocked(getDestinations).mockResolvedValue(dests)
    vi.mocked(getScores).mockResolvedValue({ daily_xp_goal: 10 } as Awaited<ReturnType<typeof getScores>>)
    vi.mocked(getEstimate).mockResolvedValue(estimate as Awaited<ReturnType<typeof getEstimate>>)
    vi.mocked(getCurrentPlan).mockResolvedValue({ plan: null })
  })

  it('renders the estimate in Russian with a localized month', async () => {
    render(<PlanEditor />)
    expect(await screen.findByText('Оставшийся объём')).toBeTruthy()
    expect(screen.getByRole('option', { name: 'Алгебра (пройдено 20%)' })).toBeTruthy()
    expect(screen.getByRole('button', { name: 'План по сроку' })).toBeTruthy()
    // 12 days is "many"; the date is the Russian month and year, from Intl.
    expect(screen.getByText(/≈ 12 дней/)).toBeTruthy()
    expect(screen.getByText(/июнь 2026/)).toBeTruthy()
    // The number is grouped the Russian way, with a no-break space.
    expect(screen.getByText(/≈ 1\s234 XP/)).toBeTruthy()
    expect(screen.getByText('Опережение сохранённого плана: ≈2 дня')).toBeTruthy()
    expect(screen.getByRole('button', { name: 'Спланировать обучение' })).toBeTruthy()
    expectNoKeyLeak()
  })

  it('monthLabel follows the locale', () => {
    expect(monthLabel('2026-03-27')).toMatch(/март 2026/)
  })
})

describe('DailyGoalControl in Russian', () => {
  it('labels the target and the custom field', () => {
    render(<DailyGoalControl current={30} onSaved={() => {}} />)
    expect(screen.getByText('Цель на день, XP · сейчас 30 в день')).toBeTruthy()
    expect(screen.getByLabelText('Своя цель на день, XP')).toBeTruthy()
    expect(screen.getByRole('button', { name: 'Задать' })).toBeTruthy()
    expectNoKeyLeak()
  })
})

describe('MonthlyCards in Russian', () => {
  it('names the month and weekdays in Russian and counts days with the plural form', () => {
    const now = new Date()
    const iso = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-02`
    const { container } = render(
      <MonthlyCards data={[{ date: iso, questions: 5, correct: 4, concepts: ['a', 'b', 'c', 'd', 'e'] }]} />,
    )
    const strip = container.querySelector('[data-testid="month-strip"]')!
    const heading = strip.querySelector('p')!.textContent ?? ''
    expect(heading).toMatch(/^[А-Я][а-я]+ \d{4}1 день$/)
    expect(strip.textContent).toContain('Пн')
    expect(strip.textContent).toContain('5 вопросов')
    expect(strip.textContent).toContain('5 тем')
    expect(strip.textContent).toContain('80% верно')
    expect(screen.getByRole('button', { name: 'Предыдущий месяц' })).toBeTruthy()
    expectNoKeyLeak()
  })
})
