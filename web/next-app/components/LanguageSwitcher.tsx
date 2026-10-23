'use client'

import { PUBLISHED_LOCALES, LOCALE_NAMES, setLocale, useLocale, useT } from '@/lib/i18n'

/**
 * The language choice, as two buttons.
 *
 * Buttons rather than a <select> because there are two options, both should be visible, and a
 * select's popup is hard to hit on a phone. The names are written in their own language and
 * carry a `lang` attribute, so a learner who landed in the wrong language can still find
 * theirs and a screen reader pronounces each correctly.
 */
export default function LanguageSwitcher() {
  const t = useT()
  const locale = useLocale()
  return (
    <div role="group" aria-label={t('common.language.label')} className="flex flex-wrap gap-2">
      {PUBLISHED_LOCALES.map(l => (
        <button
          type="button"
          key={l}
          lang={l}
          aria-pressed={l === locale}
          onClick={() => setLocale(l)}
          className={`rounded-none px-4 h-9 text-xs font-mono border ${
            l === locale
              ? 'border-mathua-blue text-mathua-blue bg-mathua-blue-faint'
              : 'border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue'
          }`}
        >
          {LOCALE_NAMES[l]}
        </button>
      ))}
    </div>
  )
}
