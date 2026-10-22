import { describe, it, expect } from 'vitest'
import { mkdtempSync, mkdirSync, writeFileSync, rmSync } from 'node:fs'
import { tmpdir } from 'node:os'
import { join } from 'node:path'
import { checkMessages } from '../scripts/i18n-check.mjs'

// The real dictionaries: en and ru must have the same keys, and every message the same
// {placeholders}. A key present in only one language renders as its raw id on screen.
describe('en and ru dictionaries', () => {
  it('agree on keys, placeholders and plural forms', () => {
    expect(checkMessages()).toEqual([])
  })
})

// The checker itself, against fixtures: a checker that never fails proves nothing.
describe('i18n-check', () => {
  function withDictionaries(en: Record<string, unknown>, ru: Record<string, unknown>) {
    const dir = mkdtempSync(join(tmpdir(), 'i18n-'))
    mkdirSync(join(dir, 'en'))
    mkdirSync(join(dir, 'ru'))
    writeFileSync(join(dir, 'en', 'a.json'), JSON.stringify(en))
    writeFileSync(join(dir, 'ru', 'a.json'), JSON.stringify(ru))
    return dir
  }
  const run = (en: Record<string, unknown>, ru: Record<string, unknown>) => {
    const dir = withDictionaries(en, ru)
    try {
      return checkMessages(dir)
    } finally {
      rmSync(dir, { recursive: true, force: true })
    }
  }
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
    expect(run({ a: 'x' }, { a: 'х', b: 'у' })).toContain('en is missing a.b')
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
      'a.n: plural in one language and plain text in the other',
    )
  })

  it('enforces Russian typography', () => {
    const problems = run({ a: 'x', b: 'y', c: 'z' }, { a: 'Нажмите "ОК"', b: 'Загрузка...', c: 'А - Б' })
    expect(problems.join('\n')).toMatch(/a\.a: straight double quote/)
    expect(problems.join('\n')).toMatch(/a\.b: three dots/)
    expect(problems.join('\n')).toMatch(/a\.c: spaced hyphen/)
  })
})
