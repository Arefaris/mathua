import { afterEach, describe, expect, it, vi } from 'vitest'
import { listAdminReports } from '../lib/api'

// A contact message has no question behind it, so every answer field is the Go
// zero value. `answer_input_mode` is a plain string without omitempty, which
// means the key arrives as "" and used to fail an enum of ['numeric','text'].
// One such row made the whole Reports queue unreadable for the administrator.
function report(overrides: Record<string, unknown> = {}) {
  return {
    id: 1,
    reporter_id: 'reporter',
    concept_id: '',
    kind: 'message',
    question: '',
    expected: '',
    explanation: '',
    lesson_id: '',
    source: 'contact',
    session_id: '',
    attempt_id: '',
    answer_given: '',
    grading_type: '',
    answer_hint: '',
    answer_input_mode: '',
    reason: '',
    detail: 'The site would not let me answer.',
    status: 'open',
    resolution: '',
    resolved_by: '',
    resolved_at: '',
    created_at: '2026-10-22T00:00:00Z',
    ...overrides,
  }
}

function respond(body: unknown) {
  vi.stubGlobal('fetch', vi.fn(async () => new Response(JSON.stringify(body), {
    status: 200,
    headers: { 'Content-Type': 'application/json' },
  })))
}

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('the admin reports payload', () => {
  it('accepts a pre-capture report whose answer mode is empty', async () => {
    respond({ reports: [report()] })
    const { reports } = await listAdminReports()
    expect(reports).toHaveLength(1)
    expect(reports[0].answer_input_mode).toBeUndefined()
  })

  it('keeps the mode when one was captured', async () => {
    respond({ reports: [report({ kind: 'wrong_question', answer_input_mode: 'numeric' })] })
    const { reports } = await listAdminReports()
    expect(reports[0].answer_input_mode).toBe('numeric')
  })
})
