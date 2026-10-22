import { describe, it, expect, afterEach, vi } from 'vitest'
import { createResolver, pluralForm } from '../lib/i18n/translate'
import {
  LOCALE_STORAGE_KEY,
  getLocale,
  getLocaleEpoch,
  resetLocaleForTests,
  setLocale,
  subscribeLocale,
} from '../lib/i18n/locale'
import { formatNumber } from '../lib/i18n/format'
import { getTranslator } from '../lib/i18n'
import { authedFetch } from '../lib/auth'
import { createInflightCache, createSessionCache } from '../lib/requestCache'

// Small dictionaries that name every form, so the tests do not depend on what the product says.
const fixtures = {
  en: {
    greet: 'Hello, {name}!',
    onlyEn: 'Only English',
    topic: { one: '{count} topic', other: '{count} topics' },
    none: { zero: 'No topics', one: '{count} topic', other: '{count} topics' },
    group: { title: 'Group' },
  },
  ru: {
    greet: 'Здравствуйте, {name}!',
    topic: { one: '{count} тема', few: '{count} темы', many: '{count} тем', other: '{count} темы' },
    none: { zero: 'Тем нет', one: '{count} тема', few: '{count} темы', many: '{count} тем', other: '{count} темы' },
  },
}
const resolve = createResolver(fixtures)
const t = (locale: 'en' | 'ru', key: string, params?: Record<string, string | number>) =>
  resolve(locale, key, params)

afterEach(() => {
  try {
    localStorage.clear()
  } catch {
    /* storage unavailable */
  }
  resetLocaleForTests('en')
})

describe('t(): interpolation', () => {
  it('fills {name} from the params', () => {
    expect(t('en', 'greet', { name: 'Ada' })).toBe('Hello, Ada!')
    expect(t('ru', 'greet', { name: 'Ада' })).toBe('Здравствуйте, Ада!')
  })

  it('leaves an unknown placeholder visible rather than printing nothing', () => {
    expect(t('en', 'greet')).toBe('Hello, {name}!')
    expect(t('en', 'greet', { other: 1 })).toBe('Hello, {name}!')
  })

  it('does not reformat a number it is given', () => {
    // A decimal inside a problem must reach the screen with its dot.
    expect(t('ru', 'greet', { name: 2.5 })).toBe('Здравствуйте, 2.5!')
  })
})

describe('t(): Russian plurals', () => {
  it.each([
    [1, '1 тема'],
    [2, '2 темы'],
    [3, '3 темы'],
    [4, '4 темы'],
    [5, '5 тем'],
    [11, '11 тем'],
    [12, '12 тем'],
    [14, '14 тем'],
    [21, '21 тема'],
    [22, '22 темы'],
    [25, '25 тем'],
    [101, '101 тема'],
    [111, '111 тем'],
    [112, '112 тем'],
  ])('%i -> %s', (count, expected) => {
    expect(t('ru', 'topic', { count })).toBe(expected)
  })

  it('uses the "other" form for a fraction', () => {
    expect(t('ru', 'topic', { count: 1.5 })).toBe('1.5 темы')
  })

  it('reports the CLDR category', () => {
    expect(pluralForm('ru', 1)).toBe('one')
    expect(pluralForm('ru', 2)).toBe('few')
    expect(pluralForm('ru', 5)).toBe('many')
    expect(pluralForm('ru', 21)).toBe('one')
    expect(pluralForm('ru', 22)).toBe('few')
    expect(pluralForm('ru', 11)).toBe('many')
  })
})

describe('t(): English plurals', () => {
  it('picks one and other', () => {
    expect(t('en', 'topic', { count: 1 })).toBe('1 topic')
    expect(t('en', 'topic', { count: 2 })).toBe('2 topics')
    expect(t('en', 'topic', { count: 0 })).toBe('0 topics')
  })

  it('uses an explicit zero form when the message has one', () => {
    expect(t('en', 'none', { count: 0 })).toBe('No topics')
    expect(t('ru', 'none', { count: 0 })).toBe('Тем нет')
    expect(t('ru', 'none', { count: 5 })).toBe('5 тем')
  })

  it('falls back to "other" when no count is given', () => {
    expect(t('en', 'topic')).toBe('{count} topics')
  })
})

describe('t(): fallback', () => {
  it('uses English for a key Russian lacks', () => {
    expect(t('ru', 'onlyEn')).toBe('Only English')
    expect(t('ru', 'group.title')).toBe('Group')
  })

  it('uses the English plural rules for English text shown in Russian', () => {
    const only = createResolver({ en: { n: { one: '{count} item', other: '{count} items' } }, ru: {} })
    expect(only('ru', 'n', { count: 2 })).toBe('2 items')
    expect(only('ru', 'n', { count: 5 })).toBe('5 items')
  })

  it('returns undefined when neither language has the key, so translate() can show the key', () => {
    expect(t('ru', 'nope.missing')).toBeUndefined()
    expect(t('en', 'group')).toBeUndefined()
  })
})

describe('the locale store', () => {
  it('defaults to Russian when the browser asks for it, and English otherwise', () => {
    const lang = vi.spyOn(window.navigator, 'language', 'get')
    resetLocaleForTests(null)
    lang.mockReturnValue('ru-RU')
    expect(getLocale()).toBe('ru')
    resetLocaleForTests(null)
    lang.mockReturnValue('de-DE')
    expect(getLocale()).toBe('en')
    resetLocaleForTests(null)
    lang.mockReturnValue('en-US')
    expect(getLocale()).toBe('en')
    lang.mockRestore()
  })

  it('prefers a stored choice to the browser language', () => {
    const lang = vi.spyOn(window.navigator, 'language', 'get').mockReturnValue('ru')
    localStorage.setItem(LOCALE_STORAGE_KEY, 'en')
    resetLocaleForTests(null)
    expect(getLocale()).toBe('en')
    lang.mockRestore()
  })

  it('ignores a stored value that is not a locale', () => {
    localStorage.setItem(LOCALE_STORAGE_KEY, 'xx')
    resetLocaleForTests(null)
    expect(getLocale()).toBe('en')
  })

  it('persists a switch under mathua_locale and updates <html lang>', () => {
    setLocale('ru')
    expect(localStorage.getItem('mathua_locale')).toBe('ru')
    expect(document.documentElement.lang).toBe('ru')
    setLocale('en')
    expect(document.documentElement.lang).toBe('en')
  })

  it('still switches when storage throws', () => {
    const set = vi.spyOn(Storage.prototype, 'setItem').mockImplementation(() => {
      throw new Error('blocked')
    })
    setLocale('ru')
    expect(getLocale()).toBe('ru')
    set.mockRestore()
  })

  it('notifies subscribers and bumps the epoch only on a real change', () => {
    const listener = vi.fn()
    const off = subscribeLocale(listener)
    const before = getLocaleEpoch()
    setLocale('en')
    expect(listener).not.toHaveBeenCalled()
    setLocale('ru')
    expect(listener).toHaveBeenCalledTimes(1)
    expect(getLocaleEpoch()).toBe(before + 1)
    off()
    setLocale('en')
    expect(listener).toHaveBeenCalledTimes(1)
  })
})

describe('formatting', () => {
  it('uses the locale for display numbers', () => {
    expect(formatNumber(1234.5, 'en')).toBe('1,234.5')
    // Russian groups with a no-break space and separates decimals with a comma.
    expect(formatNumber(1234.5, 'ru').replace(/\s/g, ' ')).toBe('1 234,5')
  })
})

describe('the API contract: Accept-Language', () => {
  it('is sent with the current locale on every authedFetch call', async () => {
    const seen: string[] = []
    const fake = vi.fn(async (_input: unknown, init?: RequestInit) => {
      seen.push(new Headers(init?.headers).get('Accept-Language') ?? '')
      return new Response('{}')
    })
    vi.stubGlobal('fetch', fake)
    try {
      await authedFetch('/api/graph')
      setLocale('ru')
      await authedFetch('/api/graph')
      setLocale('en')
      await authedFetch('/api/graph')
      expect(seen).toEqual(['en', 'ru', 'en'])
    } finally {
      vi.unstubAllGlobals()
    }
  })

  it('does not override a value the caller set', async () => {
    const fake = vi.fn(async (_input: unknown, init?: RequestInit) => {
      expect(new Headers(init?.headers).get('Accept-Language')).toBe('en')
      return new Response('{}')
    })
    vi.stubGlobal('fetch', fake)
    try {
      setLocale('ru')
      await authedFetch('/x', { headers: { 'Accept-Language': 'en' } })
      expect(fake).toHaveBeenCalledTimes(1)
    } finally {
      vi.unstubAllGlobals()
    }
  })
})

describe('request caches are keyed on the locale', () => {
  it('does not answer a Russian request with the English entry', async () => {
    const once = createSessionCache()
    const run = vi.fn(async () => getLocale())
    expect(await once('graph', run)).toBe('en')
    setLocale('ru')
    expect(await once('graph', run)).toBe('ru')
    expect(run).toHaveBeenCalledTimes(2)
    // Switching back finds the English entry still there.
    setLocale('en')
    expect(await once('graph', run)).toBe('en')
    expect(run).toHaveBeenCalledTimes(2)
  })

  it('does not share an in-flight request across a switch', async () => {
    const dedupe = createInflightCache()
    const run = vi.fn(async () => getLocale())
    const a = dedupe('kp', run)
    setLocale('ru')
    const b = dedupe('kp', run)
    expect(run).toHaveBeenCalledTimes(2)
    expect(await a).toBe('en')
    expect(await b).toBe('ru')
  })
})

describe('the typed key', () => {
  it('rejects a key the English dictionary does not have, and shows the key if one slips through', () => {
    const tr = getTranslator('en')
    // @ts-expect-error -- not a key in common.json
    expect(tr('common.language.nope')).toBe('common.language.nope')
    expect(tr('common.language.label')).toBe('Language')
    expect(tr.dynamic('common.language.nope', undefined, 'fallback')).toBe('fallback')
    expect(tr.dynamic('common.language.label')).toBe('Language')
  })
})
