import { test, expect } from '@playwright/test'

// The glossary is a static reference: no API routes are stubbed here on
// purpose. If the page ever needs the network to render, these tests fail,
// which is exactly the regression this feature must not have.
test('glossary loads with categories, search and a result count', async ({ page }) => {
  const errors: string[] = []
  page.on('console', msg => {
    if (msg.type() === 'error') errors.push(msg.text())
  })
  await page.goto('/glossary')
  await expect(page.getByRole('heading', { name: 'Mathematical glossary' })).toBeVisible({ timeout: 30_000 })
  // Category browsing is visible with no query.
  await expect(page.getByRole('heading', { name: /Greek letters/ })).toBeVisible()
  await expect(page.getByRole('searchbox', { name: 'Search symbols, names or meanings' })).toBeVisible()
  await expect(page.getByText(/\d+ entries/)).toBeVisible()
  expect(errors).toEqual([])
})

test('glossary searches by symbol, name and alias', async ({ page }) => {
  await page.goto('/glossary')
  const box = page.getByRole('searchbox', { name: 'Search symbols, names or meanings' })
  await expect(box).toBeVisible({ timeout: 30_000 })

  await box.fill('∀')
  await expect(page.getByRole('heading', { name: 'Universal quantifier' })).toBeVisible()

  await box.fill('integral')
  await expect(page.getByRole('heading', { name: 'Integral', exact: true })).toBeVisible()

  await box.fill('n choose k')
  await expect(page.getByRole('heading', { name: 'Binomial coefficient' })).toBeVisible()
})

test('glossary filters by category, clears, and handles empty results', async ({ page }) => {
  await page.goto('/glossary')
  await expect(page.getByRole('searchbox', { name: 'Search symbols, names or meanings' })).toBeVisible({ timeout: 30_000 })

  await page.getByRole('button', { name: 'Geometry and trigonometry' }).click()
  await expect(page.getByRole('heading', { name: 'Sine', exact: true })).toBeVisible()
  await expect(page.getByRole('heading', { name: 'Integral', exact: true })).toHaveCount(0)

  await page.getByRole('button', { name: 'Clear search and filters' }).click()
  await expect(page.getByRole('heading', { name: /Greek letters/ })).toBeVisible()

  await page.getByRole('searchbox', { name: 'Search symbols, names or meanings' }).fill('zzzqqq')
  await expect(page.getByText('No entries match')).toBeVisible()
  // Two clear controls are visible in the empty state (search bar + empty
  // panel); either resets the page.
  await page.getByRole('button', { name: 'Clear search and filters' }).first().click()
  await expect(page.getByText('No entries match')).toHaveCount(0)
})

test('glossary shows contextual meanings for an ambiguous symbol', async ({ page }) => {
  await page.goto('/glossary')
  const box = page.getByRole('searchbox', { name: 'Search symbols, names or meanings' })
  await expect(box).toBeVisible({ timeout: 30_000 })
  await box.fill('π')
  await expect(page.getByRole('heading', { name: 'Pi', exact: true })).toBeVisible()
  await expect(page.getByRole('heading', { name: 'Pi (circle constant)' })).toBeVisible()
})

test('glossary paginates twenty per page with URL state', async ({ page }) => {
  await page.goto('/glossary')
  await expect(page.getByRole('heading', { name: 'Mathematical glossary' })).toBeVisible({ timeout: 30_000 })
  await expect(page.getByText(/Showing 1–20 of \d+ entries/)).toBeVisible()
  // Page 1 holds Greek letters; page 2 must show different cards.
  await expect(page.getByRole('heading', { name: 'Alpha', exact: true })).toBeVisible()
  await page.getByRole('button', { name: 'Next page' }).click()
  await expect(page).toHaveURL(/page=2/)
  await expect(page.getByText(/Showing 21–40 of \d+ entries/)).toBeVisible()
  await expect(page.getByRole('heading', { name: 'Alpha', exact: true })).toHaveCount(0)
  // Back takes us to page 1 through history.
  await page.goBack()
  await expect(page.getByRole('heading', { name: 'Alpha', exact: true })).toBeVisible()
})

test('glossary clamps bad page params and resets on new queries', async ({ page }) => {
  await page.goto('/glossary?page=999')
  await expect(page.getByRole('heading', { name: 'Mathematical glossary' })).toBeVisible({ timeout: 30_000 })
  // Clamped to the last page: no Next page available.
  await expect(page.getByRole('button', { name: 'Next page' })).toBeDisabled()
  const box = page.getByRole('searchbox', { name: 'Search symbols, names or meanings' })
  await box.fill('integral')
  await expect(page.getByRole('heading', { name: 'Integral', exact: true })).toBeVisible()
  await expect(page).not.toHaveURL(/page=/)
})

test('glossary related links jump across pages', async ({ page }) => {
  await page.goto('/glossary')
  await expect(page.getByRole('heading', { name: 'Mathematical glossary' })).toBeVisible({ timeout: 30_000 })
  // Summation lives far from the letter sigma; follow the cross-link.
  const box = page.getByRole('searchbox', { name: 'Search symbols, names or meanings' })
  await box.fill('sigma')
  await expect(page.getByRole('heading', { name: 'Sigma', exact: true })).toBeVisible()
  await page.getByRole('button', { name: /Show Summation/ }).click()
  await expect(page.getByRole('heading', { name: 'Summation', exact: true })).toBeVisible()
})
