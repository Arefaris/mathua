import { test, expect } from '@playwright/test'
import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'

/**
 * The figures a lesson body references must actually render.
 *
 * Nothing above this catches a blank image. The Go extractor proves a reference is
 * well-formed, `TestEveryAssetResolvesToAShippedDiagram` proves the file is on disk, and
 * `TestPlacedFiguresAreInTheSectionTheirPlacementNamed` proves each figure sits in the section
 * it was placed in — all three are satisfied by a diagram a browser declines to draw. An HTTP
 * 200 carrying an SVG with no viewBox, a zero-height `<img>`, a figure the renderer drops: each
 * passes every other gate and shows a learner an empty box where the instruction should be.
 *
 * Two things make this a measurement rather than a decoration:
 *
 *  - It serves the real lesson files from `data/lessons/`, not fixtures. A fixture would test
 *    the harness; the corpus file tests the corpus.
 *
 *  - It carries a negative control. A check that cannot fail is worse than no check, because
 *    it is read as coverage, so the control asserts the same measurement reports zero for a
 *    figure-less lesson and for a path nothing serves. If the control ever passed, the
 *    measurement would be vacuous and this spec would prove nothing.
 *
 * Per ADR-029, `npx playwright` serves the prebuilt `out/`, so `next build` must run first or
 * this exercises the previous build and passes regardless of the change under test.
 *
 * It used to drive `/study`, which rendered every lesson body in full. ADR-047 closed that
 * route and ADR-054 restored it as a context-bound reference page; either way this spec drives
 * `/learn`, whose reference panel is the renderer that resolves a figure reference out of
 * markdown — which is why the spec moved rather than being deleted. See `stubLesson` for the
 * configuration it measures now.
 */

const ROOT = resolve(__dirname, '../../..')

/** Real bodies, chosen because they now carry figures and because one does not. */
const CASES = [
  { concept: 'discrete.sets.operations', file: 'algebrica/sets-and-numbers/sets.md', expectFigures: 6 },
  { concept: 'alg.log.concept', file: 'algebrica/powers-radicals-logarithms/logarithms.md', expectFigures: 3 },
  { concept: 'alg.exp.concept', file: 'algebrica/functions/exponential-function.md', expectFigures: 3 },
  { concept: 'calc.deriv.cauchy_mvt', file: 'algebrica/differential-calculus-theorems/cauchy-theorem.md', expectFigures: 1 },
  { concept: 'calc.seq.euler', file: 'algebrica/sequences/euler-number-limit-sequence.md', expectFigures: 1 },
  // A reclaimed OpenStax figure, and the reason the selector below is `/diagrams/` rather than
  // `/diagrams/algebrica/`: it was vendored-only, so it counted nothing in this directory and a
  // figure that failed to draw here would have been reported as a lesson with no figures at all.
  // Seven, not one: this body carries Figures 1.2 and 1.4-1.9, the arithmetic reclaims.
  { concept: 'arith.round.tens', file: 'teaching/arith.round.tens.md', expectFigures: 7 },
  // The other direction: a lesson that carried no figure at all until Figure 1.12 was drawn
  // into it. One, so what is measured is a reclaimed figure drawing in a body that had nothing
  // else in it — the arithmetic case above cannot show this, since it counts figures that have
  // been drawing alongside each other since before this case existed.
  { concept: 'arith.mult.tables', file: 'teaching/arith.mult.tables.md', expectFigures: 1 },
  // Two, the one count the cases above miss, and the count the fraction reclaims produce: 4.2
  // and 4.3 are cited by the same seven lessons, so this is a body with two reclaimed figures
  // and nothing else. Together the cases now cover one, two, three, six and seven figures.
  { concept: 'frac.basics.parts', file: 'teaching/frac.basics.parts.md', expectFigures: 2 },
  // The other kind of two: the percent pair is a table rather than a drawing, so this case
  // proves a reclaimed figure that is ruled cells and text renders, which the lines, blocks and
  // bars above cannot show. It is also the batch where upstream shipped one image for both
  // figures, so what is measured here is the figure the lesson needs and not the one upstream
  // linked twice.
  { concept: 'pct.convert.to_dec', file: 'teaching/pct.convert.to_dec.md', expectFigures: 2 },
  // One, and the one primitive the cases above do not carry: Figure 5.5's measurement bracket
  // is an unfilled path with turned-down ends, so this asserts that a reclaimed figure made of
  // paths draws at all and not only figures made of lines, rects and text.
  { concept: 'dec.ops.sub', file: 'teaching/dec.ops.sub.md', expectFigures: 1 },
  // One, with the other primitive that needs a sanitizer to be generous: Figure 5.6's
  // `Circumference` is thirteen `<text>` elements carrying a `transform`, so this asserts a
  // figure reaches the page at all when its labels are placed by rotation and not by `x`/`y`.
  { concept: 'frac.ops.to_decimal', file: 'teaching/frac.ops.to_decimal.md', expectFigures: 1 },
]

const FIGURELESS = { concept: 'arith.add.word', file: 'teaching/arith.add.word.md' }

function bodyFor(file: string): string {
  return readFileSync(resolve(ROOT, 'data/lessons', file), 'utf8')
}

/**
 * Stub the lesson API with a real body, and progress with nothing.
 *
 * The knowledge-point shard is stubbed empty on purpose. The figures this spec measures are
 * referenced inline in the *article*, and the article is what the reference panel falls back to
 * for a concept with no shards — so an empty shard is the one configuration that renders it.
 *
 * That is not the common path and the spec says so. A concept with shards is taught through
 * them: `/learn` asks the server which figures belong to which step and shows those, so the
 * figure appears with the instruction that needs it rather than in a wall of prose. This spec
 * covers the other renderer, which is the one that resolves `/diagrams/…` references out of
 * markdown, and it is the renderer most likely to go quietly untested after `/study` closed.
 */
async function stubLesson(page: import('@playwright/test').Page, concept: string, file: string) {
  const body = bodyFor(file)
  const title = body.split('\n')[0].replace(/^#\s*/, '')
  await page.route('**/api/lessons/body*', route =>
    route.fulfill({ json: { title, body, source: file.split('/')[0], assets: [] } }),
  )
  await page.route('**/api/lessons', route =>
    route.fulfill({ json: { lessons: { [file.split('/')[0]]: [{ title, body: '', concepts: [concept] }] } } }),
  )
  await page.route('**/api/lessons/*/kp*', route =>
    route.fulfill({ json: { concept_id: concept, kps: [] } }),
  )
  await page.route('**/api/lessons/*/readiness*', route =>
    route.fulfill({ json: { concept_id: concept, ready: true, weak: [], missing: [] } }),
  )
  await page.route('**/api/lessons/*/practice*', route =>
    route.fulfill({ json: { concept_id: concept, questions: [] } }),
  )
  await page.route('**/api/progress/**', route => route.fulfill({ json: {} }))
}

/**
 * Open `/learn` for a concept and expand the reference panel that carries the article.
 *
 * The panel is a disclosure on a page that already has a learning feed above it, so its figures
 * start well outside the viewport — and `KatexContent` renders them with `loading="lazy"`, so
 * they are not fetched, and therefore have no intrinsic size, until they are near. That is right
 * for a closed panel and a problem for a measurement: a figure that has not loaded reports 0x0,
 * which is the same reading as a broken one. So each figure is scrolled into view and awaited
 * before it is measured, which is also what a reader does.
 */
async function openReference(page: import('@playwright/test').Page, concept: string) {
  await page.goto(`/learn?concept=${encodeURIComponent(concept)}`)
  const summary = page.getByText('Reference — the full lesson')
  await expect(summary).toBeVisible({ timeout: 30_000 })
  await summary.click()
}

/** Scroll a figure into view and wait for it to decode, then report its drawn size. */
async function measuredSize(locator: import('@playwright/test').Locator) {
  await locator.scrollIntoViewIfNeeded()
  await locator.evaluate(
    (e: HTMLImageElement) =>
      new Promise<void>((done) => {
        if (e.complete) return done()
        e.addEventListener('load', () => done(), { once: true })
        e.addEventListener('error', () => done(), { once: true })
      }),
  )
  return locator.evaluate((e: HTMLImageElement) => ({ w: e.naturalWidth, h: e.naturalHeight }))
}

const FIGURE = 'img[src*="/diagrams/"]'

test.describe('lesson figures render', () => {
  test('every referenced figure draws at non-zero size', async ({ page }) => {
    const problems: string[] = []

    for (const c of CASES) {
      await stubLesson(page, c.concept, c.file)
      const notFound: string[] = []
      const onResponse = (r: import('@playwright/test').Response) => {
        if (r.status() === 404 && r.url().includes('/diagrams/')) notFound.push(r.url())
      }
      page.on('response', onResponse)

      await openReference(page, c.concept)
      await expect(page.locator('body')).toContainText(bodyFor(c.file).split('\n')[0].replace(/^#\s*/, ''), {
        timeout: 30_000,
      })
      const imgs = page.locator(FIGURE)
      await expect(imgs).toHaveCount(c.expectFigures)

      for (let i = 0; i < (await imgs.count()); i++) {
        const el = imgs.nth(i)
        const src = await el.getAttribute('src')
        const size = await measuredSize(el)
        if (size.w === 0 || size.h === 0) problems.push(`${c.file}: ${src} drew at ${size.w}x${size.h}`)
      }
      if (notFound.length) problems.push(`${c.file}: 404 on ${notFound.join(', ')}`)
      page.off('response', onResponse)
    }

    expect(problems, problems.join('\n')).toEqual([])
  })

  test('negative control: a lesson with no figures reports none', async ({ page }) => {
    await stubLesson(page, FIGURELESS.concept, FIGURELESS.file)
    await openReference(page, FIGURELESS.concept)
    await expect(page.locator('body')).toContainText(/\w/, { timeout: 30_000 })
    // The page really does render lesson prose, so zero here means "no figures", not "no
    // lesson" — which is the distinction that would otherwise make the test above vacuous.
    await expect(page.locator(FIGURE)).toHaveCount(0)
  })

  test('negative control: a missing diagram measures as zero, not as drawn', async ({ page }) => {
    await stubLesson(page, CASES[0].concept, CASES[0].file)
    await openReference(page, CASES[0].concept)
    const first = page.locator(FIGURE).first()
    await first.scrollIntoViewIfNeeded()
    await expect(first).toBeVisible({ timeout: 30_000 })

    // The same measurement the passing test relies on, against a path nothing serves. If this
    // ever reported a non-zero size, "the figures drew" would mean nothing.
    const size = await page.evaluate(async () => {
      const img = new Image()
      img.src = '/diagrams/algebrica/definitely-not-a-real-diagram-xyz.svg'
      await new Promise((r) => {
        img.onload = r
        img.onerror = r
      })
      return { w: img.naturalWidth, h: img.naturalHeight }
    })
    expect(size.w === 0 || size.h === 0).toBe(true)
  })
})