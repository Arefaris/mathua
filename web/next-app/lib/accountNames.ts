import { getLocale } from './i18n/locale'
import { getTranslator } from './i18n/react'

// Names the server sends as English words: a learner's level ("Math Architect"), a league tier
// ("Bronze") and the source of an attempt ("quiz"). They are display text here and identifiers
// there, so they are translated by lookup and shown as sent when the dictionary does not know
// them. A new tier added on the server therefore appears in English instead of as a blank.
//
// The caller must be a component that reads the locale (`useT`), or it will not re-render when
// the learner switches.

/** "Math Architect" becomes "mathArchitect", the form the dictionary keys use. */
function keyOf(name: string): string {
  const words = name.trim().split(/\s+/)
  return words.map((w, i) => (i === 0 ? w.toLowerCase() : w.charAt(0).toUpperCase() + w.slice(1).toLowerCase())).join('')
}

function lookup(group: string, name: string): string {
  return getTranslator(getLocale()).dynamic(`account.${group}.${keyOf(name)}`, undefined, name)
}

export const levelName = (level: string) => lookup('level', level)
export const tierName = (tier: string) => lookup('tier', tier)
export const sourceName = (source: string) => lookup('source', source)
