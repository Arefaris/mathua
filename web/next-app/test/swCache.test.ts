import { readFileSync } from 'node:fs'
import { join } from 'node:path'
import { describe, expect, it } from 'vitest'

// The service worker is plain JS shipped from public/, so no unit test loads
// it. This guard reads it: API responses are per-learner and must never be
// served from a cache that the next person at this browser would inherit.
describe('the service worker', () => {
  const src = readFileSync(join(process.cwd(), 'public/sw.js'), 'utf8')

  it('never intercepts or caches an API request', () => {
    expect(src).not.toMatch(/pathname\.startsWith\(['"]\/api\//)
    expect(src).not.toMatch(/DATA_CACHE/)
  })

  it('still serves navigations and hashed assets', () => {
    expect(src).toMatch(/navigationNetworkFirst/)
    expect(src).toMatch(/cacheFirst\(request, ASSET_CACHE\)/)
  })
})
