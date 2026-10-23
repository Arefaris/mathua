import { describe, expect, it, vi, afterEach } from 'vitest'
import { classifyFailure, failureHeadline } from '../lib/failure'

function withStatus(status: number, message = 'boom') {
  return Object.assign(new Error(message), { status, serverMessage: message })
}

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('classifyFailure', () => {
  it('treats edge timeouts as retryable and application errors as not', () => {
    expect(classifyFailure(withStatus(504))).toMatchObject({ kind: 'gateway', retryable: true })
    expect(classifyFailure(withStatus(502))).toMatchObject({ kind: 'gateway', retryable: true })
    expect(classifyFailure(withStatus(500))).toMatchObject({ kind: 'server', retryable: false })
  })

  it('separates a spent credential from a throttled one', () => {
    expect(classifyFailure(withStatus(401))).toMatchObject({ kind: 'auth', retryable: false })
    expect(classifyFailure(withStatus(429))).toMatchObject({ kind: 'throttled', retryable: true })
  })

  it('reads an abort as a timeout', () => {
    const abort = new Error('aborted')
    abort.name = 'AbortError'
    expect(classifyFailure(abort).kind).toBe('timeout')
  })

  it('reads a failed fetch as offline when the browser says so, network otherwise', () => {
    vi.stubGlobal('navigator', { onLine: false })
    expect(classifyFailure(new TypeError('fetch failed')).kind).toBe('offline')
    vi.stubGlobal('navigator', { onLine: true })
    expect(classifyFailure(new TypeError('fetch failed')).kind).toBe('network')
  })

  it('keeps a parse failure distinct from a transport one', () => {
    expect(classifyFailure(new SyntaxError('unexpected token')).kind).toBe('parse')
  })
})

describe('failureHeadline', () => {
  it('never promises a retry that cannot help', () => {
    expect(failureHeadline(classifyFailure(withStatus(500)))).toMatch(/will not change it/)
    expect(failureHeadline(classifyFailure(withStatus(503)))).toMatch(/Retry in a moment/)
  })
})
