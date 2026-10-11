import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import ProfileHeader from '../components/ProfileHeader'
import type { CatalogueConcept } from '../lib/progress'
import type { ConceptProgress, Scores } from '../lib/api'

// The learner header is one card: identity and XP above, "where you are" below. It replaced two
// stacked bordered boxes that printed MASTERED twice, so the guarded property is that mastery is
// named once and the reach headline is still there.

const catalogue: CatalogueConcept[] = [
  { id: 'a', domain: 'arithmetic', prerequisites: [] },
  { id: 'b', domain: 'arithmetic', prerequisites: ['a'] },
  { id: 'c', domain: 'arithmetic', prerequisites: ['a'] },
]

const scores = {
  lifetime_points: 100, weekly_score: 10, speed_bonus: 0, concepts_mastered: 1,
  current_streak: 2, level: 'Novice', xp_total: 120, xp_today: 6, daily_xp_goal: 10,
} as Scores

const progress = { a: { status: 'MASTERED', streak: 3 } } satisfies Record<string, ConceptProgress>

describe('ProfileHeader', () => {
  it('shows identity, XP stats and where you are in one block', () => {
    render(
      <ProfileHeader name="Tester" scores={scores} catalogue={catalogue} progress={progress} />,
    )
    expect(screen.getByText('Tester')).toBeTruthy()
    for (const label of ['XP TODAY', 'TOTAL XP', 'STREAK']) {
      expect(screen.getByText(label)).toBeTruthy()
    }
    expect(screen.getByText(/WEEKLY/)).toBeTruthy()
    const position = screen.getByRole('region', { name: 'Where you are' })
    expect(position.textContent).toMatch(/of \d+ you can reach/)
  })

  it('names mastery once, not in both the strip and the block', () => {
    render(
      <ProfileHeader name="Tester" scores={scores} catalogue={catalogue} progress={progress} />,
    )
    // The strip's MASTERED tile is suppressed; the position breakdown owns the fact.
    expect(screen.getAllByText(/^mastered$/i)).toHaveLength(1)
    expect(screen.queryByText('MASTERED')).toBeNull()
  })
})
