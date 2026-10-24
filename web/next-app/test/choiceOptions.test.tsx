import { render, screen } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import ChoiceOptions from '../components/ChoiceOptions'

// A tap must submit the letter. Submitting the option text put prose through a
// math field, which dropped spaces and rewrote symbols, so the server compared a
// rendering to the authored string and marked correct answers wrong (#154, #179).
describe('ChoiceOptions submission', () => {
  const question = 'Which statement is true?\n\nA) 1 − P(A)\n\nB) P(A)\n\nC) down on the left'

  it('sends the letter of the tapped option', async () => {
    const onPick = vi.fn()
    render(<ChoiceOptions question={question} value="" onPick={onPick} />)

    screen.getByRole('button', { name: /C\s*down on the left/ }).click()

    expect(onPick).toHaveBeenCalledWith('C')
  })

  it('marks the option whose letter is held', () => {
    render(<ChoiceOptions question={question} value="c" onPick={() => {}} />)

    const pressed = screen
      .getAllByRole('button')
      .filter(button => button.getAttribute('aria-pressed') === 'true')
    expect(pressed).toHaveLength(1)
    expect(pressed[0].textContent).toContain('down on the left')
  })

  it('does not treat the old option text as a selection', () => {
    render(<ChoiceOptions question={question} value="down on the left" onPick={() => {}} />)

    const pressed = screen
      .getAllByRole('button')
      .filter(button => button.getAttribute('aria-pressed') === 'true')
    expect(pressed).toHaveLength(0)
  })
})
