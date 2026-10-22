import { test, expect } from './fixtures'

// The public board answers with a page, not an array, and the avatar is a shape
// the client turns into a URL. These fixtures mirror that contract; mocking the
// old array shape made the page sit on its error state.
const WEEKLY = {
  view: 'weekly',
  page: 1,
  page_size: 25,
  total: 6,
  total_pages: 1,
  week_start: '2026-10-19T00:00:00Z',
  week_end: '2026-10-26T00:00:00Z',
  entries: [
    {
      rank: 1,
      id: 'ada',
      name: 'Ada',
      avatar: { kind: 'generated', style: 'bottts' },
      xp: 120,
      xp_total: 500,
      mastered: 5,
      level: 'Novice',
    },
  ],
}

const MY_LEAGUE = {
  week: '2026-10-19',
  week_start: '2026-10-19T00:00:00Z',
  week_end: '2026-10-26T00:00:00Z',
  tier: 'bronze',
  league_no: 1,
  size: 6,
  promote: 2,
  demote: 2,
  members: [
    { rank: 1, id: 'ada', name: 'Ada', avatar: { kind: 'generated', style: 'bottts' }, xp: 120, is_me: true },
    { rank: 2, id: 'bob', name: 'Bob', xp: 90 },
    { rank: 3, id: 'cy', name: 'Cy', xp: 70 },
    { rank: 4, id: 'dee', name: 'Dee', xp: 50 },
    { rank: 5, id: 'eli', name: 'Eli', xp: 30 },
    { rank: 6, id: 'fay', name: 'Fay', xp: 10 },
  ],
  me: { rank: 1, xp: 120 },
}

test('shows names with avatars in the board and in your league', async ({ page }) => {
  // The board is fetched with a query string, so the glob must allow one.
  await page.route('**/api/leaderboard*', route => route.fulfill({ json: WEEKLY }))
  await page.route('**/api/leagues', route => route.fulfill({ json: MY_LEAGUE }))
  await page.goto('/leaderboard')

  const row = page.locator('tbody tr', { hasText: 'Ada' })
  await expect(row).toBeVisible({ timeout: 30_000 })
  await expect(row.getByRole('img', { name: 'Ada avatar' })).toBeVisible()

  // The league is its own view, and it states the rule the server applies.
  await page.getByRole('tab', { name: 'My league' }).click()
  await expect(page.getByText('bronze league 1')).toBeVisible()
  await expect(page.getByText(/Top 2 move up a tier/)).toBeVisible()
  const leagueRow = page.locator('tbody tr', { hasText: 'Ada' })
  await expect(leagueRow.getByRole('img', { name: 'Ada avatar' })).toBeVisible()
  await expect(leagueRow.getByText('you')).toBeVisible()
})

test('blank names fall back to the public id, never a blank cell', async ({ page }) => {
  await page.route('**/api/leaderboard*', route =>
    route.fulfill({
      json: {
        ...WEEKLY,
        entries: [{ ...WEEKLY.entries[0], name: '   ', id: 'mystery_fox42' }],
      },
    }),
  )
  await page.route('**/api/leagues', route => route.fulfill({ json: MY_LEAGUE }))
  await page.goto('/leaderboard')
  await expect(page.locator('tbody tr', { hasText: 'mystery_fox42' })).toBeVisible({ timeout: 30_000 })
})
