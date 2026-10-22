// Message lookup, `{name}` interpolation and plural selection.
//
// Pure functions over a dictionary, so they are tested without a DOM. The dictionaries
// themselves are in `messages/`; `translate` is told which locale to read.

import { DEFAULT_LOCALE, type Locale } from './locale'
import { MESSAGES } from './messages'

export type Params = Record<string, string | number>

const PLURAL_FORMS = ['zero', 'one', 'two', 'few', 'many', 'other'] as const
export type PluralForm = (typeof PLURAL_FORMS)[number]

/** A message that changes with `count`: one entry per CLDR plural category the language uses. */
export type PluralMessage = { other: string } & Partial<Record<PluralForm, string>>

type Node = string | PluralMessage | { [key: string]: Node }

export function isPluralMessage(node: unknown): node is PluralMessage {
  if (node === null || typeof node !== 'object') return false
  const record = node as Record<string, unknown>
  const keys = Object.keys(record)
  return (
    keys.includes('other') &&
    keys.every(k => (PLURAL_FORMS as readonly string[]).includes(k) && typeof record[k] === 'string')
  )
}

function lookup(tree: unknown, key: string): string | PluralMessage | undefined {
  let node = tree as Node | undefined
  for (const part of key.split('.')) {
    if (node === undefined || typeof node === 'string' || isPluralMessage(node)) return undefined
    node = (node as { [key: string]: Node })[part]
  }
  if (typeof node === 'string' || isPluralMessage(node)) return node
  return undefined
}

const rules = new Map<Locale, Intl.PluralRules>()
function pluralRules(locale: Locale): Intl.PluralRules {
  let r = rules.get(locale)
  if (!r) {
    r = new Intl.PluralRules(locale)
    rules.set(locale, r)
  }
  return r
}

/** Russian: 1, 21, 101 are `one`; 2-4, 22-24 are `few`; 0, 5-20, 25-30 are `many`. */
export function pluralForm(locale: Locale, count: number): PluralForm {
  return pluralRules(locale).select(count) as PluralForm
}

function pick(message: string | PluralMessage, locale: Locale, params?: Params): string {
  if (typeof message === 'string') return message
  const count = Number(params?.count)
  if (!Number.isFinite(count)) return message.other
  // English has no `zero` category, so "No concepts yet" would be unreachable without this.
  if (count === 0 && message.zero !== undefined) return message.zero
  return message[pluralForm(locale, count)] ?? message.other
}

function interpolate(text: string, params?: Params): string {
  if (!params) return text
  // An unknown name is left as written: a visible `{name}` in the UI is a bug someone will
  // report, and an empty string would hide it.
  return text.replace(/\{(\w+)\}/g, (whole, name: string) =>
    name in params ? String(params[name]) : whole,
  )
}

/**
 * Builds the lookup over a given set of dictionaries. The app uses the real ones through
 * `resolve` and `translate` below; tests hand in small fixtures so they can name every plural
 * form without depending on what the product happens to say today.
 */
export function createResolver(messages: Record<Locale, unknown>) {
  return function resolve(locale: Locale, key: string, params?: Params): string | undefined {
    const message = lookup(messages[locale], key)
    if (message !== undefined) return interpolate(pick(message, locale, params), params)
    if (locale === DEFAULT_LOCALE) return undefined
    const fallback = lookup(messages[DEFAULT_LOCALE], key)
    // The plural form follows the language the text is actually in.
    return fallback === undefined ? undefined : interpolate(pick(fallback, DEFAULT_LOCALE, params), params)
  }
}

/** Looks in `locale`, then in English. `undefined` means neither has the key. */
export const resolve = createResolver(MESSAGES)

/** Looks a key up in `locale`, then in English, then returns the key itself. */
export function translate(locale: Locale, key: string, params?: Params): string {
  return resolve(locale, key, params) ?? key
}
