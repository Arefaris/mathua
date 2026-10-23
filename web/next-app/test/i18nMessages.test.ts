import { describe, it, expect } from 'vitest'
import { mkdtempSync, mkdirSync, writeFileSync, rmSync } from 'node:fs'
import { tmpdir } from 'node:os'
import { join } from 'node:path'
import { checkAll, checkMessages } from '../scripts/i18n-check.mjs'

// The real dictionaries: every published language must have the same keys as English, and
// every message the same {placeholders}. A key present in only one language renders as its raw
// id on screen.
describe('the shipped dictionaries', () => {
  it('agree on keys, placeholders and plural forms', () => {
    expect(checkMessages()).toEqual([])
  })
})

// The checker itself, against fixtures: a checker that never fails proves nothing.
describe('i18n-check', () => {
  function withDictionaries(dicts: Record<string, Record<string, unknown>>) {
    const dir = mkdtempSync(join(tmpdir(), 'i18n-'))
    for (const [locale, tree] of Object.entries(dicts)) {
      mkdirSync(join(dir, locale))
      writeFileSync(join(dir, locale, 'a.json'), JSON.stringify(tree))
    }
    return dir
  }
  const runAll = (dicts: Record<string, Record<string, unknown>>, options = {}) => {
    const dir = withDictionaries(dicts)
    try {
      return checkAll(dir, options)
    } finally {
      rmSync(dir, { recursive: true, force: true })
    }
  }
  const run = (en: Record<string, unknown>, ru: Record<string, unknown>) => runAll({ en, ru }).problems
  const plural = {
    en: { one: '{count} topic', other: '{count} topics' },
    ru: { one: '{count} тема', few: '{count} темы', many: '{count} тем', other: '{count} темы' },
  }

  it('passes matching dictionaries', () => {
    expect(run({ a: 'Hello {name}', n: plural.en }, { a: 'Привет, {name}', n: plural.ru })).toEqual([])
  })

  it('reports a key missing from ru', () => {
    expect(run({ a: 'x', b: 'y' }, { a: 'х' })).toContain('ru is missing a.b')
  })

  it('reports a key only ru has', () => {
    expect(run({ a: 'x' }, { a: 'х', b: 'у' })).toContain('ru has a.b, which English does not')
  })

  it('reports a dropped placeholder and an invented one', () => {
    const problems = run({ a: 'Level {level}', b: 'Plain' }, { a: 'Уровень', b: 'Просто {x}' })
    expect(problems).toContain('ru a.a: missing placeholder {level}')
    expect(problems).toContain('ru a.b: unknown placeholder {x}')
  })

  it('requires the Russian plural forms', () => {
    const problems = run({ n: plural.en }, { n: { one: '{count} тема', other: '{count} темы' } })
    expect(problems).toContain('ru a.n: plural lacks few')
    expect(problems).toContain('ru a.n: plural lacks many')
  })

  it('reports a plural on one side only', () => {
    expect(run({ n: plural.en }, { n: '{count} тем' })).toContain(
      'ru a.n: plural in one language and plain text in the other',
    )
  })

  it('enforces Russian typography', () => {
    const problems = run({ a: 'x', b: 'y', c: 'z' }, { a: 'Нажмите "ОК"', b: 'Загрузка...', c: 'А - Б' })
    expect(problems.join('\n')).toMatch(/a\.a: straight double quote/)
    expect(problems.join('\n')).toMatch(/a\.b: three dots/)
    expect(problems.join('\n')).toMatch(/a\.c: spaced hyphen/)
  })

  it('lets an unpublished language be partial but not wrong', () => {
    const { problems, progress } = runAll(
      { en: { a: 'x {n}', b: 'y' }, uk: { a: 'ікс' } },
      { locales: ['en', 'uk'], published: ['en'] },
    )
    expect(problems).not.toContain('uk is missing a.b')
    expect(problems).toContain('uk a.a: missing placeholder {n}')
    expect(progress).toEqual(['uk: 1/2 keys translated (not published)'])
  })

  it('holds a published language to every key', () => {
    const { problems } = runAll({ en: { a: 'x', b: 'y' }, uk: { a: 'ікс' } }, { locales: ['en', 'uk'], published: ['en', 'uk'] })
    expect(problems).toContain('uk is missing a.b')
  })

  it('rejects letters of the neighbouring alphabet', () => {
    const problems = runAll({ en: { a: 'x', b: 'y' }, ru: { a: 'Нові' }, uk: { b: 'Объём' } }).problems.join('\n')
    expect(problems).toMatch(/ru a\.a: letters of another alphabet/)
    expect(problems).toMatch(/uk a\.b: letters of another alphabet/)
  })

  it('requires the typographic apostrophe in Ukrainian', () => {
    const problems = runAll({ en: { a: 'x' }, uk: { a: "п'ять" } }).problems
    expect(problems).toContain('uk a.a: straight apostrophe, use ’')
    expect(runAll({ en: { a: 'x' }, uk: { a: 'п’ять' } }).problems).toEqual([])
  })
})
