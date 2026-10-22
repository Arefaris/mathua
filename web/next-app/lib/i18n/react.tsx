'use client'

import { Fragment, useEffect, useLayoutEffect, useSyncExternalStore, type ReactNode } from 'react'
import {
  DEFAULT_LOCALE,
  getLocale,
  getLocaleEpoch,
  getMessagesVersion,
  getRenderLocale,
  holdRenderLocale,
  subscribeLocale,
  syncDocumentLang,
  type Locale,
} from './locale'
import { resolve, translate, type Params } from './translate'
import type { MessageKey } from './keys'

/**
 * The translate function a component gets from `useT()`.
 *
 * `t('profile.title')` is checked against the English dictionary at compile time, so a typo
 * is a type error. `t.dynamic` takes a key assembled at run time (`domain.${id}`) and a
 * fallback for ids the dictionary does not know, which is the one case a literal cannot cover.
 */
export interface TFunction {
  (key: MessageKey, params?: Params): string
  dynamic(key: string, params?: Params, fallback?: string): string
  readonly locale: Locale
}

// useLayoutEffect warns when React renders on the server, where there is no layout to run it.
const useIsomorphicLayoutEffect = typeof window === 'undefined' ? useEffect : useLayoutEffect

const translators = new Map<string, TFunction>()

/**
 * One stable function per locale, so `t` is safe in a dependency array. A new one when that
 * locale's dictionary arrives, so a memo that depends on `t` recomputes with the real text.
 */
export function getTranslator(locale: Locale): TFunction {
  const cacheKey = `${locale}:${getMessagesVersion()}`
  let t = translators.get(cacheKey)
  if (!t) {
    const fn = ((key: MessageKey, params?: Params) => translate(locale, key, params)) as TFunction
    Object.defineProperty(fn, 'locale', { value: locale })
    fn.dynamic = (key, params, fallback) => resolve(locale, key, params) ?? fallback ?? key
    t = fn
    translators.set(cacheKey, t)
  }
  return t
}

/**
 * Translate outside a component, in whatever language is current.
 *
 * For `lib/` helpers (`domainLabel`, the mastery state names) that return display text. The
 * caller must be a component that already read the locale through `useT` or `useLocale`;
 * otherwise it will not re-render on a switch.
 */
export function tNow(key: MessageKey, params?: Params): string {
  return getTranslator(getRenderLocale())(key, params)
}

export function useLocale(): Locale {
  // The server snapshot is English: it is what the static export is built with, and React
  // hydrates with it before switching to the learner's own choice.
  return useSyncExternalStore(subscribeLocale, getLocale, () => DEFAULT_LOCALE)
}

export function useT(): TFunction {
  const locale = useLocale()
  // The locale string does not change when its dictionary arrives; the version does. Reading it
  // re-renders the component, and getTranslator hands back the function for the new version.
  useSyncExternalStore(subscribeLocale, getMessagesVersion, () => 0)
  return getTranslator(locale)
}

/**
 * Remounts the page whenever the learner switches language.
 *
 * Concept labels, course names and the graph arrive from the server already translated, and
 * a screen that fetched them in English would keep them in English. Remounting makes every
 * screen fetch again with the new `Accept-Language` instead of each one subscribing to the
 * locale separately. The cost is that unsaved local state is dropped, which is acceptable
 * because the switch is made on /settings and nothing else is in flight there.
 */
export function LocaleRoot({ children }: { children: ReactNode }) {
  const locale = useLocale()
  const epoch = useSyncExternalStore(subscribeLocale, getLocaleEpoch, () => 0)
  // `locale` is English while hydrating even for a Russian learner; `getLocale()` is already
  // Russian. Until the first commit, tell the render-time helpers to agree with the hook, so the
  // first client render matches the prerendered HTML. The re-render React schedules right after
  // hydration then switches everything to Russian.
  if (locale !== getLocale()) holdRenderLocale(locale)
  useIsomorphicLayoutEffect(() => {
    holdRenderLocale(null)
  }, [])
  useEffect(() => {
    syncDocumentLang()
  }, [locale])
  return <Fragment key={epoch}>{children}</Fragment>
}
