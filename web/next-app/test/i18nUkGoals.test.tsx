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

// The goals, diagnostic, onboarding and planner screens, read in Ukrainian. The English suites
// cover the same screens; this one pins that the Ukrainian text arrives, that counts take the
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
vi.mock('../components/Loading', () => ({ default: () => <div>Завантаження</div> }))
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
  resetLocaleForTests('uk')
  sessionStorage.clear()
})

describe('plural forms in the goals dictionary', () => {
  it.each([
    [1, 'Почати діагностичний тест (охоплено 1 тему)'],
    [2, 'Почати діагностичний тест (охоплено 2 теми)'],
    [5, 'Почати діагностичний тест (охоплено 5 тем)'],
    [11, 'Почати діагностичний тест (охоплено 11 тем)'],
    [21, 'Почати діагностичний тест (охоплено 21 тему)'],
    [22, 'Почати діагностичний тест (охоплено 22 теми)'],
  ])('scope of %i concepts', (count, expected) => {
    expect(getTranslator('uk')('goals.onboard.start', { count })).toBe(expected)
  })

  it('counts days, ahead and behind with the right form', () => {
    const t = getTranslator('uk')
    expect(t('goals.planner.daysApprox', { count: 3 })).toBe('≈ 3 дні')
    expect(t('goals.planner.daysApprox', { count: 12 })).toBe('≈ 12 днів')
    expect(t('goals.planner.ahead', { count: 21 })).toBe('Ви випереджаєте збережений план на ≈21 день')
    expect(t('goals.planner.behind', { count: 4 })).toBe('Ви відстаєте від збереженого плану на ≈4 дні')
  })
})

describe('/onboard in Ukrainian', () => {
  it('shows the welcome step with Ukrainian domain names and a Ukrainian start button', async () => {
    render(<OnboardPage />)
    await waitFor(() => expect(screen.getByText('Арифметика')).toBeTruthy())
    expect(screen.getByText('З’ясуймо, з чого вам почати')).toBeTruthy()
    expect(screen.getByRole('button', { name: 'Обрати все' })).toBeTruthy()
    // Per-domain concept counts use the plural form, not "N concepts".
    expect(screen.queryByText(/\d+ concepts/)).toBeNull()
    fireEvent.click(screen.getByText('Арифметика'))
    const start = screen.getByRole('button', { name: /^Почати діагностичний тест \(охоплено \d+ тем/ })
    expect(start).not.toBeDisabled()
    expectNoKeyLeak()
  })
})

describe('QuizHost in Ukrainian', () => {
  it('shows the intro and a Ukrainian failure screen', async () => {
    vi.mocked(startQuizSession).mockRejectedValue(new Error('network down'))
    render(<QuizHost />)
    expect(screen.getByText('Перевірка засвоєння')).toBeTruthy()
    fireEvent.click(screen.getByRole('button', { name: /Почати контрольну/ }))
    const alert = await screen.findByRole('alert')
    expect(alert.textContent).toContain('Контрольну не вдалося почати')
    expect(screen.getByRole('button', { name: /Спробувати ще раз/ })).toBeTruthy()
    expectNoKeyLeak()
  })
})

describe('DiagnosticResults in Ukrainian', () => {
  const plan: GoalPlanRes = {
    readiness: 0.7,
    total_tested: 20,
    correct_count: 14,
    weak_areas: { fractions: [{ id: 'frac.add', label: 'Додавання дробів' }] },
    strong_areas: { arithmetic: ['arith.add'] },
    frontier_label: 'Додавання дробів',
    frontier_conditional: true,
    conditionally_completed: ['a', 'b', 'c'],
  }

  it('renders the banner, the score and plural topic counts', () => {
    render(<DiagnosticResults plan={plan} onStartPractice={() => {}} />)
    expect(screen.getByText(/Починайте звідси: Додавання дробів/)).toBeTruthy()
    expect(screen.getByText(/попередньо/)).toBeTruthy()
    expect(screen.getByText(/Попередньо: 3 теми вважаємо відомими/)).toBeTruthy()
    expect(screen.getByText(/результат діагностичного тесту \(14\/20 правильно\)/)).toBeTruthy()
    expect(screen.getByText('Дроби')).toBeTruthy()
    expect(screen.getByText('повторити: 1 · упевнено: 0')).toBeTruthy()
    expect(screen.getByRole('button', { name: /Почати практику/ })).toBeTruthy()
    expectNoKeyLeak()
  })
})

describe('PlanEditor in Ukrainian', () => {
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

  it('renders the estimate in Ukrainian with a localized month', async () => {
    render(<PlanEditor />)
    expect(await screen.findByText('Залишок навантаження')).toBeTruthy()
    expect(screen.getByRole('option', { name: 'Алгебра (виконано 20%)' })).toBeTruthy()
    expect(screen.getByRole('button', { name: 'За терміном' })).toBeTruthy()
    // 12 days is "many"; the date is the Ukrainian month and year, from Intl.
    expect(screen.getByText(/≈ 12 днів/)).toBeTruthy()
    expect(screen.getByText(/червень 2026/)).toBeTruthy()
    // The number is grouped the Ukrainian way, with a no-break space.
    expect(screen.getByText(/≈ 1\s234 XP/)).toBeTruthy()
    expect(screen.getByText('Ви випереджаєте збережений план на ≈2 дні')).toBeTruthy()
    expect(screen.getByRole('button', { name: 'Спланувати навчання' })).toBeTruthy()
    expectNoKeyLeak()
  })

  it('monthLabel follows the locale', () => {
    expect(monthLabel('2026-03-27')).toMatch(/березень 2026/)
  })
})

describe('DailyGoalControl in Ukrainian', () => {
  it('labels the target and the custom field', () => {
    render(<DailyGoalControl current={30} onSaved={() => {}} />)
    expect(screen.getByText('Щоденна ціль у XP · зараз 30 на день')).toBeTruthy()
    expect(screen.getByLabelText('Своя щоденна ціль у XP')).toBeTruthy()
    expect(screen.getByRole('button', { name: 'Задати' })).toBeTruthy()
    expectNoKeyLeak()
  })
})

describe('MonthlyCards in Ukrainian', () => {
  it('names the month and weekdays in Ukrainian and counts days with the plural form', () => {
    const now = new Date()
    const iso = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-02`
    const { container } = render(
      <MonthlyCards data={[{ date: iso, questions: 5, correct: 4, concepts: ['a', 'b', 'c', 'd', 'e'] }]} />,
    )
    const strip = container.querySelector('[data-testid="month-strip"]')!
    const heading = strip.querySelector('p')!.textContent ?? ''
    expect(heading).toMatch(/^[а-яіїєґ]+ \d{4}( р\.)?1 день$/i)
    expect(strip.textContent).toContain('Пн')
    expect(strip.textContent).toContain('5 питань')
    expect(strip.textContent).toContain('5 тем')
    expect(strip.textContent).toContain('80% правильно')
    expect(screen.getByRole('button', { name: 'Попередній місяць' })).toBeTruthy()
    expectNoKeyLeak()
  })
})
