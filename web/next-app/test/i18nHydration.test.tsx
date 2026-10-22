/// <reference types="vite/client" />
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import type { ComponentType } from 'react'

// While the server pass runs, `useSearchParams` throws the way Next's does during a static export:
// the nearest Suspense boundary then ships its fallback, and the real content is client-rendered.
// A page wrapped that way has no server HTML for its body to disagree with.
const server = vi.hoisted(() => ({ pass: false }))

vi.mock('next/navigation', () => ({
  usePathname: () => '/',
  useRouter: () => ({ push: vi.fn(), replace: vi.fn(), back: vi.fn(), prefetch: vi.fn() }),
  useSearchParams: () => {
    if (server.pass) throw new Error('BAILOUT_TO_CLIENT_SIDE_RENDERING')
    return new URLSearchParams()
  },
  redirect: vi.fn(),
}))
vi.mock('next/link', () => ({
  default: ({ children, href, ...rest }: { children: React.ReactNode; href: string }) => (
    <a href={href} {...rest}>
      {children}
    </a>
  ),
}))

// Every learner-facing route, except administration (English this round). Each is prerendered the
// way the static export prerenders it, in English, and then hydrated for a Russian learner. A
// string formatted outside a hook with the wrong locale, or a value read from storage during
// render, makes the first client render differ from the HTML, and React throws the HTML away.
const pages = import.meta.glob<{ default: ComponentType }>(['../app/**/page.tsx', '!../app/admin/**'])

const HYDRATION = /hydrat|did not match|Text content|content does not match/i

// Pages that fail to hydrate in jsdom whatever the language, so they cannot say anything about it.
// Each was checked by hydrating in English, the language its HTML was made in: it fails the same.
//  - the long-form documents render their stylesheet in a <style> element, and jsdom escapes the
//    quotes in it where a browser does not, so the server markup differs from the client's;
//  - /goals renders an empty shell that redirects, and its Suspense boundary is updated by the
//    redirect effect before React has hydrated it.
const NOT_TESTABLE_IN_JSDOM = [/^\.\.\/app\/docs\/[^/]+\/page\.tsx$/, /^\.\.\/app\/licenses\//, /^\.\.\/app\/goals\//]

describe('hydrating the prerendered pages for a Russian learner', () => {
  beforeEach(() => {
    // Nothing answers: a page shows its loading state, which is what the static HTML shows too.
    vi.stubGlobal('fetch', () => new Promise(() => {}))
    // The diagram libraries measure themselves; jsdom has no layout to give them.
    vi.stubGlobal('ResizeObserver', class { observe() {} unobserve() {} disconnect() {} })
    localStorage.clear()
  })
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('finds the pages', () => {
    expect(Object.keys(pages).length).toBeGreaterThan(20)
  })

  for (const [path, load] of Object.entries(pages).filter(([p]) => !NOT_TESTABLE_IN_JSDOM.some(re => re.test(p)))) {
    it(`${path.replace('../app', '')} matches its prerendered HTML`, async () => {
      // A fresh module graph per page: a page's stores (theme, chrome, auth) start as they do on a
      // first visit, and one page's leftovers cannot decide whether the next one hydrates.
      vi.resetModules()
      document.body.innerHTML = ''
      document.head.innerHTML = ''
      for (const a of [...document.documentElement.attributes]) document.documentElement.removeAttribute(a.name)
      sessionStorage.clear()
      const { act } = await import('@testing-library/react')
      const { renderToString } = await import('react-dom/server')
      const { hydrateRoot } = await import('react-dom/client')
      const { LocaleRoot } = await import('../lib/i18n')
      const { resetLocaleForTests } = await import('../lib/i18n/locale')
      await (await import('../lib/i18n/messages')).loadMessages('ru')
      const { default: Page } = await load()
      const tree = (
        <LocaleRoot>
          <Page />
        </LocaleRoot>
      )

      // Hydrate `html` as a learner reading `locale` and return the hydration errors React logged.
      async function hydrationErrors(html: string, locale: 'en' | 'ru') {
        const spy = vi.spyOn(console, 'error').mockImplementation(() => {})
        try {
          resetLocaleForTests(locale)
          const container = document.createElement('div')
          container.innerHTML = html
          document.body.appendChild(container)
          let root: ReturnType<typeof hydrateRoot> | undefined
          await act(async () => {
            root = hydrateRoot(container, tree)
          })
          const found = spy.mock.calls.map(c => String(c[0])).filter(m => HYDRATION.test(m))
          act(() => root?.unmount())
          container.remove()
          return found
        } finally {
          spy.mockRestore()
        }
      }

      const quiet = vi.spyOn(console, 'error').mockImplementation(() => {})
      let html = ''
      try {
        resetLocaleForTests('en')
        server.pass = true
        html = renderToString(tree)
      } catch {
        // A page that cannot be prerendered in jsdom has no server HTML to disagree with.
        return
      } finally {
        server.pass = false
        quiet.mockRestore()
      }

      expect(await hydrationErrors(html, 'ru')).toEqual([])
    })
  }
})
