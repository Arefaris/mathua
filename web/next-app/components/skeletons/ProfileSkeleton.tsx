'use client'

import { Skeleton } from '../ui/skeleton'

// Mirrors the authed /profile layout (max-w-[820px]): one ProfileHeader card (identity + XP +
// where you are), the recommended-task rows, then the Activity heatmap.
export default function ProfileSkeleton() {
  return (
    <div
      role="status"
      aria-live="polite"
      aria-busy="true"
      aria-label="Loading profile"
      className="mx-auto w-full max-w-[820px] min-w-0 px-4 sm:px-6 py-8 sm:py-12 overflow-x-hidden"
    >
      {/* ProfileHeader: identity + XP … */}
      <div className="border-[0.5px] border-mathua-border p-4 sm:p-5 sm:px-6 w-full min-w-0">
        <div className="flex flex-col sm:flex-row flex-wrap gap-4 sm:gap-5 items-start w-full min-w-0">
          <div className="flex flex-row gap-4 sm:contents w-full sm:w-auto">
            <Skeleton className="size-[72px] rounded-full shrink-0" />
            <Skeleton className="h-[86px] min-w-[72px] sm:min-w-[80px] shrink-0" />
          </div>
          <div className="flex-1 min-w-0 w-full">
            <Skeleton className="h-5 w-40 mb-3" />
            <div className="flex flex-wrap gap-x-6 gap-y-3">
              <div className="min-w-0">
                <Skeleton className="h-3 w-16 mb-2" />
                <Skeleton className="h-4 w-24 mb-1" />
                <Skeleton className="h-1 w-full max-w-[200px]" />
              </div>
              {Array.from({ length: 3 }).map((_, i) => (
                <div key={i} className="min-w-0">
                  <Skeleton className="h-3 w-16 mb-2" />
                  <Skeleton className="h-4 w-20" />
                </div>
              ))}
            </div>
          </div>
        </div>

        {/* … then where you are, inside the same card. */}
        <div className="mt-4 border-t border-mathua-border pt-4">
          <Skeleton className="h-3 w-24 mb-3" />
          <Skeleton className="h-8 w-32 mb-3" />
          <Skeleton className="h-1 w-full" />
          <div className="mt-3 grid grid-cols-3 gap-3">
            {Array.from({ length: 3 }).map((_, i) => (
              <Skeleton key={i} className="h-9 w-full" />
            ))}
          </div>
        </div>
      </div>

      {/* Recommended task rows */}
      <div className="mt-6 border-[0.5px] border-mathua-border">
        {Array.from({ length: 3 }).map((_, i) => (
          <div
            key={i}
            className={i < 2 ? 'border-b-[0.5px] border-mathua-border px-3 py-2.5' : 'px-3 py-2.5'}
          >
            <Skeleton className="h-4 w-2/3 mb-1.5" />
            <Skeleton className="h-3 w-1/2" />
          </div>
        ))}
      </div>

      {/* Activity heatmap */}
      <section className="mt-8 flex min-w-0 flex-col items-stretch">
        <Skeleton className="h-5 w-24 mb-4" />
        <div className="w-full border border-mathua-border bg-mathua-surface p-4">
          <div className="grid grid-cols-[repeat(20,minmax(0,1fr))] gap-1">
            {Array.from({ length: 80 }).map((_, i) => (
              <Skeleton key={i} className="h-3 w-full rounded-[2px]" />
            ))}
          </div>
        </div>
      </section>
    </div>
  )
}
