// Which language the learner reads, and where that choice lives.
//
// This module is deliberately free of React so `lib/auth.ts` can ask for the locale on every
// API request without importing a component tree. The React side (`useLocale`, `useT`) is in
// `react.tsx` and subscribes to the store below.
//
// Why a module-level store and not a context provider: the site is a static export with no
// `[locale]` segment, so there is no server-known locale to hand down. A provider would only
// wrap the same client-side value, and it would make every unit test that renders a single
// component mount one first. `useSyncExternalStore` gives the same re-render on change
// without it.

export type Locale = 'en' | 'ru'

export const LOCALES: readonly Locale[] = ['en', 'ru']

/** The language of the source text, and the one the server renders at build time. */
export const DEFAULT_LOCALE: Locale = 'en'

/** Names are written in their own language, so a learner can find theirs in a language they cannot read. */
export const LOCALE_NAMES: Record<Locale, string> = {
  en: 'English',
  ru: 'Русский',
}

export const LOCALE_STORAGE_KEY = 'mathua_locale'

export function isLocale(value: unknown): value is Locale {
  return typeof value === 'string' && (LOCALES as readonly string[]).includes(value)
}

function readStored(): Locale | null {
  try {
    const stored = localStorage.getItem(LOCALE_STORAGE_KEY)
    return isLocale(stored) ? stored : null
  } catch {
    // Private mode or blocked site data: the choice then lasts for the tab only.
    return null
  }
}

/** First visit: Russian for a browser that asks for it, English for everyone else. */
function fromBrowser(): Locale {
  const lang = typeof navigator !== 'undefined' ? navigator.language : ''
  return typeof lang === 'string' && lang.toLowerCase().startsWith('ru') ? 'ru' : DEFAULT_LOCALE
}

let current: Locale | null = null
// Bumped only by an explicit switch. `LocaleRoot` keys the page tree on it, so every screen
// remounts and refetches the data that arrives localized from the server (concept labels,
// courses). It is not bumped when the stored locale is first resolved on load: that would
// remount a Russian learner's page once, straight after it mounted, and fetch everything twice.
let epoch = 0
const listeners = new Set<() => void>()

function notify() {
  listeners.forEach(fn => fn())
}

export function getLocale(): Locale {
  // The server (static prerender) and tests without a window have no learner to ask.
  if (typeof window === 'undefined') return DEFAULT_LOCALE
  if (current === null) current = readStored() ?? fromBrowser()
  return current
}

export function getLocaleEpoch(): number {
  return epoch
}

export function subscribeLocale(listener: () => void): () => void {
  listeners.add(listener)
  return () => {
    listeners.delete(listener)
  }
}

function applyToDocument(locale: Locale) {
  if (typeof document !== 'undefined') document.documentElement.lang = locale
}

export function setLocale(next: Locale) {
  if (!isLocale(next) || next === getLocale()) return
  current = next
  try {
    localStorage.setItem(LOCALE_STORAGE_KEY, next)
  } catch {
    // Not persisted; the in-memory choice still applies.
  }
  epoch += 1
  applyToDocument(next)
  notify()
}

/** Keeps `<html lang>` in step with the locale; called from `LocaleRoot` once it has mounted. */
export function syncDocumentLang() {
  applyToDocument(getLocale())
}

// Another tab changed the choice: follow it, so two open tabs never disagree about the
// language the server is being asked for.
if (typeof window !== 'undefined') {
  window.addEventListener('storage', e => {
    if (e.key !== LOCALE_STORAGE_KEY) return
    const next = isLocale(e.newValue) ? e.newValue : fromBrowser()
    if (next === current) return
    current = next
    epoch += 1
    applyToDocument(next)
    notify()
  })
}

/** Test seam: pin the locale (or forget it, so the next read starts from storage again). */
export function resetLocaleForTests(next: Locale | null = null) {
  current = next
  epoch = 0
  applyToDocument(next ?? DEFAULT_LOCALE)
  notify()
}
