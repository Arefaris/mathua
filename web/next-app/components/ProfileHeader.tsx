'use client'

import type { ConceptProgress, Scores } from '../lib/api'
import type { CatalogueConcept } from '../lib/progress'
import ProfileStats from './ProfileStats'
import PositionBlock from './PositionBlock'

/**
 * One card answering "where am I": identity and XP on top, the reach headline and its breakdown
 * below.
 *
 * These were two stacked bordered boxes. Each was reasonable alone, but together they made the
 * learner read two cards for one picture and printed MASTERED twice — once as a lifetime count,
 * once in the position breakdown. Merged, the strip keeps the XP/streak/weekly numbers and the
 * position block owns mastery, which is the fact its headline is about.
 *
 * `ProfileStats` is still used standalone on `/share`, so it keeps its border and its MASTERED
 * tile; this header asks it for the bare form and hides the duplicate.
 */
export default function ProfileHeader({
  name,
  scores,
  avatarSeed,
  avatarUrl,
  avatarPreset,
  catalogue,
  progress,
  frontierLabel,
}: {
  name: string
  scores: Scores
  avatarSeed?: string
  avatarUrl?: string
  avatarPreset?: number | null
  catalogue: CatalogueConcept[]
  progress: Record<string, ConceptProgress>
  frontierLabel?: string | null
}) {
  return (
    <div className="border-[0.5px] border-mathua-border p-4 sm:p-5 sm:px-6 w-full max-w-full min-w-0 overflow-hidden">
      <ProfileStats
        bare
        hideMastered
        name={name}
        scores={scores}
        avatarSeed={avatarSeed}
        avatarUrl={avatarUrl}
        avatarPreset={avatarPreset}
      />
      <div className="mt-4 border-t border-mathua-border pt-4">
        <PositionBlock
          embedded
          catalogue={catalogue}
          progress={progress}
          frontierLabel={frontierLabel}
        />
      </div>
    </div>
  )
}
