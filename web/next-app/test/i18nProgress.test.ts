import { describe, it, expect } from 'vitest'
import { countWords, readmeIsCurrent, renderBlock } from '../scripts/i18n-progress.mjs'

// The progress block in README.md follows the translations. A failure here means a translation
// landed without the block being refreshed: run `npm run i18n:progress` and commit README.md.
describe('README translation progress', () => {
  it('is current', () => {
    expect(readmeIsCurrent()).toBe(true)
  })
})

describe('i18n-progress', () => {
  it('counts prose words, not formulas, code or placeholders', () => {
    expect(countWords('Add the fractions \\(\\frac{1}{2}\\) and $x^2$ to get {count} items')).toBe(7)
    expect(countWords('Run `go test ./...` and see [the docs](https://example.com/a)')).toBe(5)
    expect(countWords("Couldn't load it")).toBe(3)
  })

  it('renders a row per area and an overall line for both languages', () => {
    const area = { total: 10, uk: 0, ru: 10 }
    const result = Object.fromEntries(
      ['Interface', 'Topic names', 'Courses', 'Lesson steps', 'Lesson texts', 'Worked solutions', 'Problem texts'].map(a => [a, area]),
    )
    const block = renderBlock(result)
    expect(block).toContain('| Interface | 10 | `░░░░░░░░░░░░░░░░░░░░` 0% | `████████████████████` 100% |')
    expect(block).toContain('Ukrainian-0%25-red')
    expect(block).toContain('Russian-100%25-brightgreen')
  })
})
