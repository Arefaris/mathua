import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { render } from '@testing-library/react'
import ScrollToTopOnNavigate from '../components/ScrollToTopOnNavigate'

const nav = vi.hoisted(() => ({ pathname: '/learn' }))
vi.mock('next/navigation', () => ({ usePathname: () => nav.pathname }))

let frames: FrameRequestCallback[] = []

/** Run every queued frame, including frames queued by the frames just run. */
function flushFrames() {
  for (let pass = 0; pass < 4; pass++) {
    const queued = frames
    frames = []
    for (const cb of queued) cb(0)
  }
}

function go(pathname: string) {
  nav.pathname = pathname
  rerender(<ScrollToTopOnNavigate />)
}

let rerender: (ui: React.ReactElement) => void
let scrollTo: ReturnType<typeof vi.fn>

beforeEach(() => {
  nav.pathname = '/learn'
  frames = []
  vi.stubGlobal('requestAnimationFrame', (cb: FrameRequestCallback) => {
    frames.push(cb)
    return frames.length
  })
  vi.stubGlobal('cancelAnimationFrame', (id: number) => {
    if (frames[id - 1]) frames[id - 1] = () => {}
  })
  scrollTo = vi.fn()
  vi.spyOn(window, 'scrollTo').mockImplementation(scrollTo as never)
  window.history.replaceState({}, '', '/learn')
})

afterEach(() => {
  vi.unstubAllGlobals()
  vi.restoreAllMocks()
})

describe('ScrollToTopOnNavigate', () => {
  it('does not move on the first paint, then resets the top on the next route', () => {
    const view = render(<ScrollToTopOnNavigate />)
    rerender = view.rerender
    flushFrames()
    expect(scrollTo).not.toHaveBeenCalled()

    go('/domains')
    flushFrames()
    expect(scrollTo).toHaveBeenCalledWith({ top: 0, left: 0, behavior: 'auto' })
  })

  it('leaves back and forward to the browser', () => {
    const view = render(<ScrollToTopOnNavigate />)
    rerender = view.rerender
    flushFrames()

    window.dispatchEvent(new Event('popstate'))
    go('/learn')
    flushFrames()
    expect(scrollTo).not.toHaveBeenCalled()
  })

  it('does not fight a deliberate in-page anchor', () => {
    const view = render(<ScrollToTopOnNavigate />)
    rerender = view.rerender
    flushFrames()

    window.history.replaceState({}, '', '/glossary#linear')
    go('/glossary')
    flushFrames()
    expect(scrollTo).not.toHaveBeenCalled()
  })
})
