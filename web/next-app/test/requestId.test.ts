import { afterEach, describe, expect, it, vi } from 'vitest'
import { authedFetch, REQUEST_ID_HEADER } from '../lib/auth'

function capture(): { seen: RequestInit[] } {
  const seen: RequestInit[] = []
  vi.stubGlobal('fetch', vi.fn(async (_input: RequestInfo | URL, init?: RequestInit) => {
    seen.push(init ?? {})
    return new Response('{}', { status: 200, headers: { 'Content-Type': 'application/json' } })
  }))
  return { seen }
}

afterEach(() => {
  vi.unstubAllGlobals()
  localStorage.clear()
})

// Without an id, a learner's failure report and the server's log line can only
// be matched by guessing at timestamps — which is exactly what happened when a
// diagnostic answer "just froze".
describe('request correlation', () => {
  it('sends a request id the server can echo', async () => {
    const { seen } = capture()
    await authedFetch('https://api.example/x')
    const headers = new Headers(seen[0].headers)
    expect(headers.get(REQUEST_ID_HEADER)).toMatch(/[A-Za-z0-9-]{8,}/)
  })

  it('gives two calls different ids', async () => {
    const { seen } = capture()
    await authedFetch('https://api.example/x')
    await authedFetch('https://api.example/x')
    const first = new Headers(seen[0].headers).get(REQUEST_ID_HEADER)
    const second = new Headers(seen[1].headers).get(REQUEST_ID_HEADER)
    expect(first).not.toEqual(second)
  })

  it('keeps an id a caller supplied, so a retry stays one traceable request', async () => {
    const { seen } = capture()
    await authedFetch('https://api.example/x', { headers: { [REQUEST_ID_HEADER]: 'retry-abc-123' } })
    expect(new Headers(seen[0].headers).get(REQUEST_ID_HEADER)).toBe('retry-abc-123')
  })
})
