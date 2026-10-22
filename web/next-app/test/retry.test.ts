import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { authedFetch } from '../lib/auth'
import {
  MAX_ATTEMPTS,
  retryBudget,
  retryDelay,
  retryableError,
  retryableStatus,
} from '../lib/retry'

// The reported freezes came from failures the client either retried blindly or
// never retried at all. These tests pin the policy: reads retry, mutations only
// on request, 500 never, and a caller's abort is final.

describe('retry policy', () => {
  it('retries only the statuses that mean try again', () => {
    expect([429, 502, 503, 504].every(retryableStatus)).toBe(true)
    expect([400, 401, 403, 404, 409, 422, 500].some(retryableStatus)).toBe(false)
  })

  it('retries reads by default and mutations only when the caller opts in', () => {
    expect(retryBudget('GET')).toBe(MAX_ATTEMPTS - 1)
    expect(retryBudget('HEAD')).toBe(MAX_ATTEMPTS - 1)
    expect(retryBudget('POST')).toBe(0)
    expect(retryBudget('POST', 1)).toBe(MAX_ATTEMPTS - 1)
    // An opt-in cannot raise the ceiling.
    expect(retryBudget('POST', 9)).toBe(MAX_ATTEMPTS - 1)
  })

  it('honours Retry-After and otherwise backs off within bounds', () => {
    const first = retryDelay(0)
    expect(first).toBeGreaterThanOrEqual(300)
    expect(first).toBeLessThanOrEqual(400)
    const second = retryDelay(1)
    expect(second).toBeGreaterThanOrEqual(900)
    expect(second).toBeLessThanOrEqual(1000)
    expect(retryDelay(0, '2')).toBe(2000)
    expect(retryDelay(0, '999')).toBe(10000)
    expect(retryDelay(0, 'soon')).toBeGreaterThanOrEqual(300)
  })

  it('treats our timeout and a network failure as retryable, a caller abort as final', () => {
    expect(retryableError(new TypeError('failed to fetch'), false)).toBe(true)
    const timeout = new DOMException('aborted', 'AbortError')
    expect(retryableError(timeout, false)).toBe(true)
    expect(retryableError(timeout, true)).toBe(false)
    expect(retryableError(new Error('bad json'), false)).toBe(false)
  })
})

function respond(status: number, body = '{}'): Response {
  return new Response(body, { status, headers: { 'Content-Type': 'application/json' } })
}

describe('authedFetch retries', () => {
  beforeEach(() => {
    localStorage.clear()
  })
  afterEach(() => {
    vi.unstubAllGlobals()
    vi.useRealTimers()
  })

  it('retries a read that hits 503 and returns the later response', async () => {
    const fetchMock = vi
      .fn()
      .mockResolvedValueOnce(respond(503))
      .mockResolvedValueOnce(respond(200, '{"ok":true}'))
    vi.stubGlobal('fetch', fetchMock)
    vi.useFakeTimers()

    const pending = authedFetch('https://api.test/api/config')
    await vi.advanceTimersByTimeAsync(500)
    const res = await pending
    expect(res.status).toBe(200)
    expect(fetchMock).toHaveBeenCalledTimes(2)
  })

  it('does not retry a mutation unless the caller opts in', async () => {
    const fetchMock = vi.fn().mockResolvedValue(respond(503))
    vi.stubGlobal('fetch', fetchMock)

    const res = await authedFetch('https://api.test/api/study/answer', { method: 'POST', body: '{}' })
    expect(res.status).toBe(503)
    expect(fetchMock).toHaveBeenCalledTimes(1)
  })

  it('repeats the same body when a mutation opts in, so an attempt id is stable', async () => {
    const fetchMock = vi
      .fn()
      .mockResolvedValueOnce(respond(504))
      .mockResolvedValueOnce(respond(200))
    vi.stubGlobal('fetch', fetchMock)
    vi.useFakeTimers()

    const body = JSON.stringify({ attempt_id: 'a1' })
    const pending = authedFetch('https://api.test/api/goal/diagnostic/answer', {
      method: 'POST',
      body,
      retryAttempts: 1,
    })
    await vi.advanceTimersByTimeAsync(500)
    await pending

    expect(fetchMock).toHaveBeenCalledTimes(2)
    expect((fetchMock.mock.calls[1][1] as RequestInit).body).toBe(body)
  })

  it('does not retry a 500, because the same answer cannot change it', async () => {
    const fetchMock = vi.fn().mockResolvedValue(respond(500))
    vi.stubGlobal('fetch', fetchMock)
    const res = await authedFetch('https://api.test/api/config')
    expect(res.status).toBe(500)
    expect(fetchMock).toHaveBeenCalledTimes(1)
  })

  it('does not retry after the caller aborts', async () => {
    const ctrl = new AbortController()
    const fetchMock = vi.fn().mockImplementation(() => {
      ctrl.abort()
      return Promise.reject(Object.assign(new Error('aborted'), { name: 'AbortError' }))
    })
    vi.stubGlobal('fetch', fetchMock)
    await expect(
      authedFetch('https://api.test/api/config', { signal: ctrl.signal }),
    ).rejects.toThrow()
    expect(fetchMock).toHaveBeenCalledTimes(1)
  })
})
