import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { authedFetch } from '../lib/auth'

// A guest bearer expires silently. Before this test existed, authedFetch cleared
// it and returned the 401, so the next call went out with no credential and an
// anonymous learner saw "sign in again" in the middle of a diagnostic. These
// tests pin the reclaim: mint a bearer for the same guest id, repeat once, and
// never touch a registered login or a caller-supplied bearer.

type Call = { url: string; auth: string | null }

const GUEST_ID_KEY = 'mathua_guest_id'
const GUEST_TOKEN_KEY = 'mathua_guest_token'

let calls: Call[] = []

function json(body: unknown, status = 200): Response {
  return new Response(JSON.stringify(body), {
    status,
    headers: { 'content-type': 'application/json' },
  })
}

function installFetch(onGuest: () => Response, onApi: (auth: string | null) => Response) {
  vi.stubGlobal(
    'fetch',
    vi.fn(async (input: RequestInfo | URL, init?: RequestInit) => {
      const url = String(input)
      const auth = new Headers(init?.headers).get('authorization')
      calls.push({ url, auth })
      if (url.endsWith('/api/auth/guest')) return onGuest()
      return onApi(auth)
    }),
  )
}

beforeEach(() => {
  calls = []
  localStorage.clear()
  localStorage.setItem(GUEST_ID_KEY, 'guest-abc')
  localStorage.setItem(GUEST_TOKEN_KEY, 'stale-token')
})

afterEach(() => {
  vi.unstubAllGlobals()
  localStorage.clear()
})

describe('authedFetch guest recovery', () => {
  it('mints a bearer for the same guest and repeats the request', async () => {
    installFetch(
      () => json({ token: 'fresh-token', student_id: 'guest-abc' }),
      (auth) => (auth === 'Bearer fresh-token' ? json({ ok: true }) : json({ error: 'expired' }, 401)),
    )

    const res = await authedFetch('https://api.test/api/next')

    expect(res.status).toBe(200)
    expect(localStorage.getItem(GUEST_TOKEN_KEY)).toBe('fresh-token')
    // The guest endpoint is built from API_BASE, which is empty under test, so
    // match on the path rather than pinning a URL prefix that the environment
    // decides.
    expect(calls).toHaveLength(3)
    expect(calls[0].url).toBe('https://api.test/api/next')
    expect(calls[1].url.endsWith('/api/auth/guest')).toBe(true)
    expect(calls[2].url).toBe('https://api.test/api/next')
    expect(calls[0].auth).toBe('Bearer stale-token')
    expect(calls[2].auth).toBe('Bearer fresh-token')
  })

  it('returns the 401 when the guest cannot be reclaimed', async () => {
    installFetch(
      () => json({ error: 'nope' }, 500),
      () => json({ error: 'expired' }, 401),
    )

    const res = await authedFetch('https://api.test/api/next')

    expect(res.status).toBe(401)
    expect(calls.filter((c) => c.url === 'https://api.test/api/next')).toHaveLength(1)
  })

  it('leaves a caller-supplied bearer alone', async () => {
    installFetch(
      () => json({ token: 'fresh-token', student_id: 'guest-abc' }),
      () => json({ error: 'foreign' }, 401),
    )

    const res = await authedFetch('https://api.test/api/admin', {
      headers: { Authorization: 'Bearer someone-else' },
    })

    expect(res.status).toBe(401)
    expect(calls.some((c) => c.url.endsWith('/api/auth/guest'))).toBe(false)
    expect(localStorage.getItem(GUEST_TOKEN_KEY)).toBe('stale-token')
  })
})
