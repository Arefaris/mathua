import type { Locale } from '../locale'
import { en } from './en'

export type EnglishMessages = typeof en

// English is bundled: it is the source text, the fallback for every other language and the
// language the static export is prerendered in. Every other language is a separate chunk that
// loads when a learner uses it, so an English visitor does not download Russian (about 29 kB
// gzipped) and a Russian one downloads it once. Until it arrives, lookups fall back to English.
export const MESSAGES: Record<Locale, unknown> = { en, ru: undefined }

let version = 0

/** Bumped each time a dictionary arrives, so a component can re-render without a locale change. */
export function messagesVersion(): number {
  return version
}

export function messagesLoaded(locale: Locale): boolean {
  return MESSAGES[locale] !== undefined
}

export async function loadMessages(locale: Locale): Promise<void> {
  if (messagesLoaded(locale)) return
  // A literal path per language, so the bundler can split each into its own chunk.
  if (locale === 'ru') MESSAGES.ru = (await import('./ru')).ru
  version += 1
}
