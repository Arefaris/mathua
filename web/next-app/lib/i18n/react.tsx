'use client'

import { Fragment, useEffect, useMemo, useSyncExternalStore, type ReactNode } from 'react'
import {
  DEFAULT_LOCALE,
  getLocale,
  getLocaleEpoch,
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

const translators = new Map<Locale, TFunction>()

/** One stable function per locale, so `t` is safe in a dependency array. */
export function getTranslator(locale: Locale): TFunction {
  let t = translators.get(locale)
  if (!t) {
    const fn = ((key: MessageKey, params?: Params) => translate(locale, key, params)) as TFunction
    Object.defineProperty(fn, 'locale', { value: locale })
    fn.dynamic = (key, params, fallback) => resolve(locale, key, params) ?? fallback ?? key
    t = fn
    translators.set(locale, t)
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
  return getTranslator(getLocale())(key, params)
}

export function useLocale(): Locale {
  // The server snapshot is English: it is what the static export is built with, and React
  // hydrates with it before switching to the learner's own choice.
  return useSyncExternalStore(subscribeLocale, getLocale, () => DEFAULT_LOCALE)
}

export function useT(): TFunction {
  const locale = useLocale()
  return useMemo(() => getTranslator(locale), [locale])
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
  useEffect(() => {
    syncDocumentLang()
  }, [locale])
  return <Fragment key={epoch}>{children}</Fragment>
}
