// One place that decides where the API lives.
//
// The static export is served by a host that has no API routes, so a build
// without NEXT_PUBLIC_API_URL does not degrade — every request becomes a 404
// on the page's own origin and the site looks broken while loading perfectly.
// Failing the build is the only failure mode a deploy cannot ship.
export function resolveAPIBase(configured: string | undefined): string {
  const base = (configured ?? '').trim().replace(/\/+$/, '')
  if (base) return base
  if (process.env.NODE_ENV === 'production') {
    throw new Error(
      'NEXT_PUBLIC_API_URL is required: without it the browser calls /api/* on the static host, where every request is a 404.',
    )
  }
  return ''
}

export const API_BASE = resolveAPIBase(process.env.NEXT_PUBLIC_API_URL)
