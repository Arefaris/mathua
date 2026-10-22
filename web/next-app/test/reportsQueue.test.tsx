import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import AdminReportsPage from '../app/admin/reports/page'
import { listAdminReports, bulkModerateReports } from '../lib/api'
import type { QuestionReport } from '../lib/api'

vi.mock('../lib/api', async (importOriginal) => ({
  ...(await importOriginal<typeof import('../lib/api')>()),
  listAdminReports: vi.fn(),
  moderateReport: vi.fn(),
  bulkModerateReports: vi.fn(),
}))

// A moderator, set per test: reports.read without reports.moderate is a real
// account, and it must see the queue but no decision control.
const access = { permissions: ['admin.access', 'reports.read', 'reports.moderate'] as string[] }

vi.mock('../hooks/useAdminMe', () => ({
  useAdminMe: () => ({
    me: null,
    loading: false,
    can: (p: string) => access.permissions.includes(p),
    isStaff: true,
  }),
}))

vi.mock('../components/Header', () => ({ default: () => <header /> }))
vi.mock('../components/Footer', () => ({ default: () => <footer /> }))
vi.mock('../components/BottomTabs', () => ({ default: () => <nav /> }))
vi.mock('../components/KatexContent', () => ({
  default: ({ children }: { children: React.ReactNode }) => <span>{children}</span>,
}))
vi.mock('../components/ReportReplay', () => ({ default: () => <div>replay</div> }))

function report(over: Partial<QuestionReport> & { id: number }): QuestionReport {
  return {
    reporter_id: 'learner',
    concept_id: 'algebra.linear',
    kind: 'question',
    question: '2x = 4',
    expected: '2',
    explanation: '',
    lesson_id: '',
    source: 'practice',
    session_id: '',
    attempt_id: '',
    reason: 'wrong_answer',
    detail: '',
    status: 'open',
    resolution: '',
    resolved_by: '',
    resolved_at: '',
    created_at: '2026-10-20T00:00:00Z',
    ...over,
  }
}

beforeEach(() => {
  vi.clearAllMocks()
  access.permissions = ['admin.access', 'reports.read', 'reports.moderate']
})

describe('the report queue groups by concept', () => {
  it('shows one section per concept with its report count', async () => {
    vi.mocked(listAdminReports).mockResolvedValue({
      reports: [
        report({ id: 1 }),
        report({ id: 2, status: 'reviewing' }),
        report({ id: 3, concept_id: 'calculus.limits' }),
      ],
    })
    render(<AdminReportsPage />)
    expect(await screen.findByRole('heading', { name: 'algebra.linear' })).toBeTruthy()
    expect(screen.getByRole('heading', { name: 'calculus.limits' })).toBeTruthy()
    // The section says how many complaints it holds, so "one bug" is visible
    // before any decision is made.
    expect(screen.getByText('2 reports')).toBeTruthy()
    expect(screen.getByText('1 report')).toBeTruthy()
  })

  // The point of bulk close: one fix, one reason, every report about that bug.
  it('closes every unsettled report for a concept with one shared reason', async () => {
    vi.mocked(listAdminReports).mockResolvedValue({
      reports: [report({ id: 1 }), report({ id: 2, status: 'reviewing' })],
    })
    vi.mocked(bulkModerateReports).mockResolvedValue({ updated: 2, skipped: 0, failed: [] })
    render(<AdminReportsPage />)

    fireEvent.click(await screen.findByRole('button', { name: 'Resolve all' }))
    fireEvent.change(screen.getByLabelText(/Reason for closing 2 reports/), {
      target: { value: 'Fixed in f7ed1804: the input format is declared.' },
    })
    fireEvent.click(screen.getByRole('button', { name: 'Confirm resolved' }))

    await waitFor(() => expect(bulkModerateReports).toHaveBeenCalledWith(
      { concept_id: 'algebra.linear' },
      'resolved',
      'Fixed in f7ed1804: the input format is declared.',
    ))
    // The queue is reloaded so the closed reports leave the Open filter.
    await waitFor(() => expect(listAdminReports).toHaveBeenCalledTimes(2))
    expect(await screen.findByRole('status')).toHaveTextContent('Closed 2 reports.')
  })

  it('refuses a bulk decision without a reason', async () => {
    vi.mocked(listAdminReports).mockResolvedValue({
      reports: [report({ id: 1 }), report({ id: 2 })],
    })
    render(<AdminReportsPage />)
    fireEvent.click(await screen.findByRole('button', { name: 'Dismiss all' }))
    fireEvent.click(screen.getByRole('button', { name: 'Confirm dismissed' }))
    expect(await screen.findByRole('alert')).toHaveTextContent('A reason is required')
    expect(bulkModerateReports).not.toHaveBeenCalled()
  })

  it('offers no decision control to a read-only moderator', async () => {
    access.permissions = ['admin.access', 'reports.read']
    vi.mocked(listAdminReports).mockResolvedValue({
      reports: [report({ id: 1 }), report({ id: 2 })],
    })
    render(<AdminReportsPage />)
    await screen.findByRole('heading', { name: 'algebra.linear' })
    expect(screen.queryByRole('button', { name: 'Resolve all' })).toBeNull()
    expect(screen.queryByRole('button', { name: 'Dismiss all' })).toBeNull()
  })
})
