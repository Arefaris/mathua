import { render, screen, act } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import SubmitErrorBlock from '../components/SubmitErrorBlock'

function renderBlock(error: { message: string; status?: number }, onRetry = vi.fn()) {
  render(
    <SubmitErrorBlock
      error={error}
      onRetry={onRetry}
      onRestart={() => {}}
      restartLabel="Restart"
    />,
  )
  return onRetry
}

afterEach(() => {
  vi.useRealTimers()
})

// A 429 is the server refusing an answer for pace, not a wrong answer. The
// learner should see that distinction and should not have to press Retry for a
// condition the server clears in seconds.
describe('the submit error block', () => {
  it('calls a throttled answer a pause and retries it once by itself', () => {
    vi.useFakeTimers()
    const onRetry = renderBlock({ status: 429, message: 'too many answers at once' })

    expect(screen.getByRole('alert')).toHaveTextContent(/Saving is catching up/)
    expect(screen.queryByText(/still in the box/)).toBeNull()

    act(() => {
      vi.advanceTimersByTime(2000)
    })
    expect(onRetry).toHaveBeenCalledTimes(1)

    // A second failure with the same message must not loop.
    act(() => {
      vi.advanceTimersByTime(10000)
    })
    expect(onRetry).toHaveBeenCalledTimes(1)
  })

  it('keeps an ordinary failure manual and says the answer is still there', () => {
    vi.useFakeTimers()
    const onRetry = renderBlock({ status: 500, message: 'goal answer failed' })

    expect(screen.getByRole('alert')).toHaveTextContent(/still in the box/)
    act(() => {
      vi.advanceTimersByTime(10000)
    })
    expect(onRetry).not.toHaveBeenCalled()
  })
})
