// Public surface of the i18n layer. Import from '@/lib/i18n', not from the files behind it.
export {
  DEFAULT_LOCALE,
  LOCALES,
  PUBLISHED_LOCALES,
  isPublishedLocale,
  LOCALE_NAMES,
  LOCALE_STORAGE_KEY,
  getLocale,
  isLocale,
  setLocale,
  type Locale,
} from './locale'
export { translate, pluralForm, type Params } from './translate'
export { useT, useLocale, tNow, getTranslator, LocaleRoot, type TFunction } from './react'
export { formatNumber, formatDate } from './format'
export type { MessageKey } from './keys'
