import { describe, it, expect, vi, beforeEach } from 'vitest'
import { act, fireEvent, render, screen, waitFor } from '@testing-library/react'
import ProfileHeader from '../components/ProfileHeader'
import StrugglesSection from '../components/StrugglesSection'
import AttemptList from '../components/AttemptList'
import DangerZone from '../components/DangerZone'
import LeaderboardPage from '../app/leaderboard/page'
import { setLocale } from '../lib/i18n'
import type { AttemptRecord, LeaderboardBoard, MyLeague, Scores, WeaknessRes } from '../lib/api'
import type { CatalogueConcept } from '../lib/progress'

// The account area in Russian: real strings on screen, correct plural forms, and no raw
// dictionary key left where a sentence should be. The English suite elsewhere covers the same
// components, so these only assert what changes with the language.

const { getLeaderboard, getMyLeague, resetAccount } = vi.hoisted(() => ({
  getLeaderboard: vi.fn(),
  getMyLeague: vi.fn(),
  resetAccount: vi.fn(),
}))

vi.mock('../lib/api', async (importOriginal) => {
  const actual = await importOriginal<typeof import('../lib/api')>()
  return { ...actual, getLeaderboard, getMyLeague, resetAccount, getAttempts: vi.fn() }
})
vi.mock('next/navigation', () => ({
  usePathname: () => '/leaderboard',
  useRouter: () => ({ push: vi.fn() }),
}))
vi.mock('next/link', () => ({
  default: ({ children, href }: { children: React.ReactNode; href: string }) => <a href={href}>{children}</a>,
}))
vi.mock('../hooks/useTheme', () => ({ useTheme: () => ({ mounted: true }) }))
vi.mock('../components/Header', () => ({ default: () => null }))
vi.mock('../components/Footer', () => ({ default: () => null }))
vi.mock('../components/BottomTabs', () => ({ default: () => null }))

const catalogue: CatalogueConcept[] = [
  { id: 'a', domain: 'arithmetic', prerequisites: [] },
  { id: 'b', domain: 'arithmetic', prerequisites: ['a'] },
  { id: 'c', domain: 'arithmetic', prerequisites: ['a'] },
]
const scoresFor = (streak: number) =>
  ({
    lifetime_points: 100, weekly_score: 10, speed_bonus: 0, concepts_mastered: 1,
    current_streak: streak, level: 'Novice', xp_total: 120, xp_today: 6, daily_xp_goal: 10,
  }) as Scores

beforeEach(() => {
  getLeaderboard.mockReset()
  getMyLeague.mockReset()
  resetAccount.mockReset()
  act(() => setLocale('ru'))
})

const noRawKeys = (root: HTMLElement) => expect(root.textContent).not.toMatch(/\baccount\.[a-z]+\.[a-zA-Z]+/)

describe('profile header in Russian', () => {
  it('says where the learner is, in Russian, with the level translated', () => {
    const { container } = render(
      <ProfileHeader name="Ада" scores={scoresFor(2)} catalogue={catalogue} progress={{ a: { status: 'MASTERED', streak: 3 } }} />,
    )
    const region = screen.getByRole('region', { name: 'Где вы сейчас' })
    expect(region.textContent).toContain('Ада · Ур. 01 Новичок')
    expect(region.textContent).toContain('1 из 3 доступных вам')
    expect(region.textContent).toContain('Освоено: 1 · Открыто: 3 · В процессе: 0')
    noRawKeys(container)
  })

  it.each([
    [1, '1 день подряд'],
    [2, '2 дня подряд'],
    [5, '5 дней подряд'],
    [11, '11 дней подряд'],
    [21, '21 день подряд'],
    [22, '22 дня подряд'],
  ])('streak of %i reads %s', (streak, expected) => {
    render(<ProfileHeader name="Ада" scores={scoresFor(streak)} catalogue={catalogue} progress={{}} />)
    expect(screen.getByRole('region', { name: 'Где вы сейчас' }).textContent).toContain(expected)
  })
})

describe('struggles in Russian', () => {
  it('counts domains and topics with the right plural forms', () => {
    const weaknesses = {
      by_domain: {
        arithmetic: [{ id: 'a', label: 'A', weakness: 0.9 }, { id: 'b', label: 'B', weakness: 0.2 }],
        fractions: [{ id: 'c', label: 'C', weakness: 0.7 }, { id: 'd', label: 'D', weakness: 0.1 }],
        algebra: [{ id: 'e', label: 'E', weakness: 0.3 }],
      },
    } as WeaknessRes
    const { container } = render(<StrugglesSection weaknesses={weaknesses} />)
    expect(screen.getByText('Нужно внимание')).toBeInTheDocument()
    expect(container.textContent).toContain('Трудно: 2')
    expect(container.textContent).toContain('5 тем')
    expect(container.textContent).toContain('3 раздела')
    noRawKeys(container)
  })
})

describe('attempt list in Russian', () => {
  const attempt = {
    session_id: 's', student_id: 'u', concept_id: 'arith.add.single', concept_name: 'Сложение',
    answer: '5', expected: '4', correct: false, elapsed_seconds: 3,
    timestamp: '2026-09-10T10:00:00Z', question: '2 + 2 = ?', source: 'quiz', explanation: '',
  } as AttemptRecord

  it('translates the source and the answer labels', () => {
    const { container } = render(<AttemptList attempts={[attempt]} />)
    expect(container.textContent).toContain('контрольная')
    expect(container.textContent).toContain('Вы: 5')
    expect(container.textContent).toContain('Верно: 4')
    noRawKeys(container)
  })

  it('has a Russian empty state', () => {
    render(<AttemptList attempts={[]} />)
    expect(screen.getByText('Вы ещё не ответили ни на один вопрос.')).toBeInTheDocument()
  })
})

describe('reset progress in Russian', () => {
  it('asks for the phrase in Russian and still sends the canonical English one', async () => {
    resetAccount.mockResolvedValue(undefined)
    const { container } = render(<DangerZone />)
    expect(screen.getByText('Сбросить прогресс')).toBeInTheDocument()
    expect(container.textContent).toContain('сбросить мой прогресс')
    const button = screen.getByRole('button', { name: 'Сбросить всё перечисленное' })
    expect(button).toBeDisabled()

    fireEvent.change(screen.getByLabelText(/Для подтверждения введите/), { target: { value: 'Сбросить мой прогресс' } })
    expect(button).not.toBeDisabled()
    fireEvent.click(button)
    await waitFor(() => expect(resetAccount).toHaveBeenCalledWith('reset my progress'))
    noRawKeys(container)
  })

  it('explains a mismatch in Russian', () => {
    render(<DangerZone />)
    fireEvent.change(screen.getByLabelText(/Для подтверждения введите/), { target: { value: 'сбросить' } })
    expect(screen.getByRole('status').textContent).toContain('Не совпадает')
  })
})

describe('leaderboard in Russian', () => {
  const board: LeaderboardBoard = {
    view: 'weekly', page: 1, page_size: 25, total: 3, total_pages: 1,
    week_start: '2026-10-19T00:00:00Z', week_end: '2026-10-26T00:00:00Z',
    entries: [
      { rank: 1, id: 'ada', name: 'Ada', xp: 120, xp_total: 900, mastered: 4, level: 'Scholar' },
      { rank: 2, id: 'bob', name: 'Bob', xp: 80, xp_total: 400, mastered: 2, level: 'Student', is_me: true },
      { rank: 3, id: 'cy', name: 'Cy', xp: 40, xp_total: 100, mastered: 1, level: 'Novice' },
    ],
  }
  const league: MyLeague = {
    week: '2026-10-19', week_start: '2026-10-19T00:00:00Z', week_end: '2026-10-26T00:00:00Z',
    tier: 'bronze', league_no: 1, size: 2, promote: 5, demote: 5,
    members: [{ rank: 1, id: 'ada', name: 'Ada', xp: 30, is_me: true }, { rank: 2, id: 'bob', name: 'Bob', xp: 20 }],
    me: { rank: 1, xp: 30 },
  }

  it('translates the headings, the tabs and the level names from the server', async () => {
    getLeaderboard.mockResolvedValue(board)
    const { container } = render(<LeaderboardPage />)
    expect(screen.getByRole('heading', { name: 'Рейтинг' })).toBeInTheDocument()
    expect(screen.getByRole('tab', { name: 'Эта неделя' })).toBeInTheDocument()
    expect(screen.getByRole('tab', { name: 'Моя лига' })).toBeInTheDocument()
    await screen.findAllByText('Учёный')
    expect(screen.getAllByText('вы').length).toBeGreaterThan(0)
    noRawKeys(container)
  })

  it('names the league tier and pluralises its size', async () => {
    getLeaderboard.mockResolvedValue(board)
    getMyLeague.mockResolvedValue(league)
    const { container } = render(<LeaderboardPage />)
    fireEvent.click(screen.getByRole('tab', { name: 'Моя лига' }))
    expect(await screen.findByText('Лига «Бронза», группа 1')).toBeInTheDocument()
    expect(container.textContent).toContain('2 участника')
    noRawKeys(container)
  })
})
