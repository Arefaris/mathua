import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import ContactBox from '../components/ContactBox'
import * as api from '../lib/api'

vi.mock('../lib/api', async (importOriginal) => {
  const mod = await importOriginal<typeof import('../lib/api')>()
  return { ...mod, submitReport: vi.fn() }
})

describe('the contact box', () => {
  beforeEach(() => {
    vi.mocked(api.submitReport).mockReset()
  })

  it('will not send a message too short to act on', () => {
    render(<ContactBox />)
    const send = screen.getByRole('button', { name: 'Send message' })
    expect(send).toBeDisabled()

    fireEvent.change(screen.getByLabelText('Your message'), { target: { value: 'hi' } })
    expect(send).toBeDisabled()

    fireEvent.change(screen.getByLabelText('Your message'), {
      target: { value: 'I cannot sign in and the reset link does nothing.' },
    })
    expect(send).toBeEnabled()
  })

  it('writes into the report queue and confirms, because a box that silently succeeds is a lie', async () => {
    vi.mocked(api.submitReport).mockResolvedValue({ id: 7, status: 'open' })
    render(<ContactBox />)

    fireEvent.change(screen.getByLabelText('Your message'), {
      target: { value: 'I cannot sign in and the reset link does nothing.' },
    })
    fireEvent.click(screen.getByRole('button', { name: 'Send message' }))

    await waitFor(() => expect(screen.getByRole('status')).toHaveTextContent('Message received'))
    expect(api.submitReport).toHaveBeenCalledWith(
      expect.objectContaining({
        kind: 'message',
        source: 'contact',
        detail: 'I cannot sign in and the reset link does nothing.',
      }),
    )
  })

  it('names the rate limit instead of blaming the network', async () => {
    vi.mocked(api.submitReport).mockRejectedValue(new Error('Report submit failed: 429'))
    render(<ContactBox />)

    fireEvent.change(screen.getByLabelText('Your message'), {
      target: { value: 'I cannot sign in and the reset link does nothing.' },
    })
    fireEvent.click(screen.getByRole('button', { name: 'Send message' }))

    await waitFor(() => expect(screen.getByRole('alert')).toHaveTextContent('Too many messages'))
  })
})
