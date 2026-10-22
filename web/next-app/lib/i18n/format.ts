// Display formatting that follows the locale.
//
// For text the interface writes itself: a streak count, a date in a heading. It must never be
// applied to a number inside a problem or an answer. The answer field and the grader expect a
// dot as the decimal separator, so "2.5" stays "2.5" in Russian.
import type { Locale } from './locale'

const numberFormats = new Map<string, Intl.NumberFormat>()

export function formatNumber(value: number, locale: Locale, options?: Intl.NumberFormatOptions): string {
  const cacheKey = locale + JSON.stringify(options ?? {})
  let f = numberFormats.get(cacheKey)
  if (!f) {
    f = new Intl.NumberFormat(locale, options)
    numberFormats.set(cacheKey, f)
  }
  return f.format(value)
}

export function formatDate(value: Date | string | number, locale: Locale, options?: Intl.DateTimeFormatOptions): string {
  const d = value instanceof Date ? value : new Date(value)
  if (isNaN(d.getTime())) return ''
  return new Intl.DateTimeFormat(locale, options ?? { year: 'numeric', month: 'short', day: 'numeric' }).format(d)
}
