import { describe, it, expect } from 'vitest'
import { checkSources } from '../scripts/i18n-sources.mjs'

// Every translation was made from the English text that is there now. A failure here usually
// means an upstream merge reworded English copy: retranslate the listed entries against the
// current English, then `npm run i18n:sources -- stamp <locale> <entries>`.
describe('translations', () => {
  it('are not stale against the current English', () => {
    for (const [locale, r] of Object.entries(checkSources())) {
      expect({ locale, stale: r.stale, unstamped: r.unstamped }).toEqual({ locale, stale: [], unstamped: [] })
    }
  })
})
