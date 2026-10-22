'use client'

// Route changes must land at the top of the new page.
//
// The browser tries. Chrome resets the scroll on a push-navigation, but Safari
// on iOS keeps the old offset until the new document has laid out, so a short
// page can open scrolled past its own content and then appear to jump. A global
// `scroll-behavior: smooth` made that worse: the reset itself was animated, so
// the reader watched a second-long glide through the new page before it settled.
//
// Two behaviours this must not break:
//   - Back and forward restore the position the reader left. That is the
//     browser's job, and a `popstate` in the last second means it is happening.
//   - An in-page `#fragment` link is a deliberate position, not a new page.
//
// The double frame is not decoration: the first frame runs before the new page
// has been painted and can be undone by the navigation's own scroll handling.

import { useEffect, useRef } from 'react'
import { usePathname } from 'next/navigation'

export default function ScrollToTopOnNavigate() {
  const pathname = usePathname()
  const first = useRef(true)
  const popped = useRef(false)

  useEffect(() => {
    const onPop = () => {
      popped.current = true
      window.setTimeout(() => { popped.current = false }, 1000)
    }
    window.addEventListener('popstate', onPop)
    return () => window.removeEventListener('popstate', onPop)
  }, [])

  useEffect(() => {
    if (first.current) {
      first.current = false
      return
    }
    if (popped.current || window.location.hash) return

    let inner = 0
    const outer = requestAnimationFrame(() => {
      inner = requestAnimationFrame(() => {
        window.scrollTo({ top: 0, left: 0, behavior: 'auto' })
      })
    })
    return () => {
      cancelAnimationFrame(outer)
      cancelAnimationFrame(inner)
    }
  }, [pathname])

  return null
}
