// Retry policy for the browser's API calls.
//
// The rule this file exists to keep: a retry must never be able to cause a
// second side effect. A read is safe to repeat. A mutation is repeated only
// when its caller says so, and the only caller that does is an answer
// submission, which carries its attempt id and is replayed by the engine
// rather than graded twice.
//
// Nothing retries a 500: the server has already decided the request was bad,
// and repeating it just makes the learner wait twice.

/** Backoff before the first and second retry. No jitter is added here. */
export const RETRY_DELAYS_MS = [300, 900]

/** Two attempts total: the original and one retry. */
export const MAX_ATTEMPTS = 2

/** Statuses worth repeating. 500 is deliberately absent. */
export function retryableStatus(status: number): boolean {
  return status === 429 || status === 502 || status === 503 || status === 504
}

/** A read may retry by itself; a write only when its caller opted in. */
export function retryBudget(method: string | undefined, optIn?: number): number {
  const m = (method ?? 'GET').toUpperCase()
  const allowed = m === 'GET' || m === 'HEAD'
  if (optIn != null) return Math.max(0, Math.min(optIn, MAX_ATTEMPTS - 1))
  return allowed ? MAX_ATTEMPTS - 1 : 0
}

/** Full jitter of up to 100ms keeps retries from synchronising. */
export function retryDelay(attempt: number, retryAfter?: string | null): number {
  const header = retryAfterSeconds(retryAfter)
  if (header != null) return header * 1000
  const base = RETRY_DELAYS_MS[Math.min(attempt, RETRY_DELAYS_MS.length - 1)]
  return base + Math.floor(Math.random() * 100)
}

function retryAfterSeconds(value?: string | null): number | null {
  if (!value) return null
  const seconds = Number(value)
  return Number.isFinite(seconds) && seconds >= 0 ? Math.min(seconds, 10) : null
}

/**
 * A failed fetch is repeatable when it never reached the server, which is what
 * a TypeError from `fetch` means. An AbortError is only repeatable when our own
 * timeout fired; a caller's abort must stay an abort.
 */
export function retryableError(err: unknown, callerAborted: boolean): boolean {
  if (callerAborted) return false
  if (err instanceof TypeError) return true
  if (err instanceof DOMException && err.name === 'AbortError') return true
  return false
}

/** The promise a caller awaits between attempts. */
export function wait(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms))
}
