import { describe, it, expect } from 'vitest'
import { parseChoices, stemWithoutChoices } from '../lib/choices'

describe('parseChoices', () => {
  it('parses trailing "A) …" option blocks', () => {
    const q = 'Which shape?\n\nA) circle\n\nB) square\n\nC) triangle'
    expect(parseChoices(q)).toEqual([
      { letter: 'A', text: 'circle' },
      { letter: 'B', text: 'square' },
      { letter: 'C', text: 'triangle' },
    ])
  })

  it('returns null when there is no option list', () => {
    expect(parseChoices('Simplify: \\((5^2)^2\\)')).toBeNull()
    expect(parseChoices('2 + 2 = ?')).toBeNull()
  })

  it('returns null when letters are not sequential from A', () => {
    expect(parseChoices('Q\n\nA) x\n\nC) y')).toBeNull()
  })

  it('preserves math in option text', () => {
    const q = 'Pick one\n\nA) $x^2$\n\nB) $x^3$'
    expect(parseChoices(q)).toEqual([
      { letter: 'A', text: '$x^2$' },
      { letter: 'B', text: '$x^3$' },
    ])
  })
})

describe('stemWithoutChoices', () => {
  it('drops the option block the buttons already show', () => {
    const q = 'Which shape?\n\nA) circle\n\nB) square'
    expect(stemWithoutChoices(q)).toBe('Which shape?')
  })

  it('leaves a question with no options untouched', () => {
    const q = 'Simplify: \\((5^2)^2\\)'
    expect(stemWithoutChoices(q)).toBe(q)
  })

  it('keeps a multi-line stem that only mentions a letter in passing', () => {
    const q = 'Let A be the area.\nWhat is A?'
    expect(stemWithoutChoices(q)).toBe(q)
  })

  it('does not cut a stem whose list starts at B', () => {
    const q = 'Which?\n\nB) circle\n\nC) square'
    expect(stemWithoutChoices(q)).toBe(q)
  })
})
