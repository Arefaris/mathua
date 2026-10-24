// One vocabulary for the ways a request can fail. Before this, every caller
// collapsed a timeout, an offline browser, a 502 from the edge and a rejected
// credential into "something went wrong", which told the learner nothing and
// told support less: the same toast covered a retry that would work and one
// that never would.
export type FailureKind =
  | 'offline'
  | 'timeout'
  | 'network'
  | 'auth'
  | 'throttled'
  | 'not_found'
  | 'server'
  | 'gateway'
  | 'client'
  | 'parse'

export interface Failure {
  kind: FailureKind
  /** HTTP status when the server answered; absent for transport failures. */
  status?: number
  message: string
  /** Whether repeating the same request can plausibly succeed. */
  retryable: boolean
}

interface ErrorWithStatus {
  status?: number
  serverMessage?: string
  name?: string
}

// A 504 or 502 comes from the edge, not from a broken request: the same call
// usually works a second later. A 500 is the application's own failure and
// retrying it just burns the learner's patience, so it is deliberately not
// retryable.
function classifyStatus(status: number): Pick<Failure, 'kind' | 'retryable'> {
  if (status === 401 || status === 403) return { kind: 'auth', retryable: false }
  if (status === 404) return { kind: 'not_found', retryable: false }
  if (status === 429) return { kind: 'throttled', retryable: true }
  if (status === 502 || status === 503 || status === 504) return { kind: 'gateway', retryable: true }
  if (status >= 500) return { kind: 'server', retryable: false }
  return { kind: 'client', retryable: false }
}

export function classifyFailure(error: unknown): Failure {
  const err = (error ?? {}) as ErrorWithStatus
  const status = typeof err.status === 'number' ? err.status : undefined
  const message = typeof err.serverMessage === 'string' && err.serverMessage
    ? err.serverMessage
    : error instanceof Error ? error.message : ''

  if (status !== undefined) {
    return { status, kind: classifyStatus(status).kind, retryable: classifyStatus(status).retryable, message }
  }
  // A caller-cancelled request is not a failure of the network; the abort is
  // indistinguishable from a timeout here, so both read as a timeout and the
  // caller decides whether to surface it.
  if (err.name === 'AbortError') {
    return { kind: 'timeout', message: message || 'The request took too long.', retryable: true }
  }
  if (typeof navigator !== 'undefined' && navigator.onLine === false) {
    return { kind: 'offline', message: message || 'The browser is offline.', retryable: true }
  }
  if (error instanceof TypeError) {
    return { kind: 'network', message: message || 'The connection dropped.', retryable: true }
  }
  if (error instanceof SyntaxError) {
    return { kind: 'parse', message: message || 'The server sent a response we could not read.', retryable: false }
  }
  return { kind: 'network', message, retryable: true }
}

// What the learner is told. Each line says what happened and whether waiting
// or retrying can help, which is the difference between a usable error and a
// dead end.
export function failureHeadline(f: Failure): string {
  switch (f.kind) {
    case 'offline':
      return 'You are offline. Your answer is kept and will be sent when the connection returns.'
    case 'timeout':
      return 'The server took too long to answer. Your answer is kept.'
    case 'network':
      return 'The connection dropped before the answer arrived. It is still in the box.'
    case 'auth':
      return 'Your session ended. Sign in again and the answer can be sent.'
    case 'throttled':
      return 'Saving is catching up, your answer is kept and will be sent again.'
    case 'not_found':
      return 'The server forgot this run. Restart it to continue.'
    case 'gateway':
      return 'The server is restarting. Retry in a moment, your answer is kept.'
    case 'server':
      return 'The server hit an error while grading. Retrying the same answer will not change it.'
    case 'parse':
      return 'The server answered in a form this page could not read. Reload to continue.'
    default:
      return "Couldn't submit your answer, it's still in the box."
  }
}
