// Check that every translated dictionary says the same things as the English one.
//
// A translation can be wrong in ways no type checker sees: a key that exists in one language
// and not the other renders as the raw key id, and a dropped `{count}` renders a sentence that
// reads fine and shows no number. This script compares each dictionary with English key by key
// and placeholder by placeholder.
//
// A language that learners can pick (PUBLISHED_LOCALES in lib/i18n/locale.ts) must be complete.
// One that is still being translated may miss keys, which fall back to English, but what it
// does have is held to the same rules, so it cannot be published with errors already in it.
//
// It is also a test (test/i18nMessages.test.ts imports `checkMessages`), so `npm test` fails on
// a mismatch and no separate CI step has to remember to run it.
//
// Usage: npm run i18n:check
import { existsSync, readdirSync, readFileSync } from 'node:fs'
import { dirname, join } from 'node:path'
import { fileURLToPath } from 'node:url'

const here = dirname(fileURLToPath(import.meta.url))
const defaultDir = join(here, '..', 'lib', 'i18n', 'messages')
const localeFile = join(here, '..', 'lib', 'i18n', 'locale.ts')

const PLURAL_FORMS = ['zero', 'one', 'two', 'few', 'many', 'other']
// Categories each language's rules can select (CLDR). A plural message must carry every one it
// needs, or a count that selects the missing one silently falls through to `other`.
const REQUIRED_FORMS = {
  en: ['one', 'other'],
  ru: ['one', 'few', 'many', 'other'],
  uk: ['one', 'few', 'many', 'other'],
}

// Letters of a neighbouring alphabet are a slip, not a translation: Ukrainian і, ї, є, ґ in
// Russian text, and Russian ы, э, ъ, ё in Ukrainian.
const FOREIGN_LETTERS = {
  ru: /[ґєіїҐЄІЇ]/,
  uk: /[ыэъёЫЭЪЁ]/,
}

/** Reads a `const NAME: readonly Locale[] = [...]` list out of lib/i18n/locale.ts. */
export function localeList(name) {
  const src = readFileSync(localeFile, 'utf8')
  const m = src.match(new RegExp(`export const ${name}: readonly Locale\\[\\] = \\[([^\\]]*)\\]`))
  if (!m) throw new Error(`i18n:check cannot find ${name} in lib/i18n/locale.ts`)
  return [...m[1].matchAll(/'([a-z]+)'/g)].map(x => x[1])
}

function isPlural(node) {
  if (node === null || typeof node !== 'object' || Array.isArray(node)) return false
  const keys = Object.keys(node)
  return keys.includes('other') && keys.every(k => PLURAL_FORMS.includes(k) && typeof node[k] === 'string')
}

function placeholders(text) {
  return new Set([...text.matchAll(/\{(\w+)\}/g)].map(m => m[1]))
}

/** Flatten one dictionary to `key -> { plural, texts[], forms[] }`. */
export function flatten(tree, prefix, out, problems, where) {
  for (const [name, node] of Object.entries(tree)) {
    const key = prefix ? `${prefix}.${name}` : name
    if (typeof node === 'string') {
      out.set(key, { plural: false, texts: [node], forms: [] })
    } else if (isPlural(node)) {
      out.set(key, { plural: true, texts: Object.values(node), forms: Object.keys(node) })
    } else if (node !== null && typeof node === 'object' && !Array.isArray(node)) {
      flatten(node, key, out, problems, where)
    } else {
      problems.push(`${where}: ${key} is not a string, a plural object or a group`)
    }
  }
}

export function load(dir, locale, problems = []) {
  const out = new Map()
  if (!existsSync(join(dir, locale))) return { out, files: [] }
  const files = readdirSync(join(dir, locale)).filter(f => f.endsWith('.json')).sort()
  for (const file of files) {
    const area = file.replace(/\.json$/, '')
    const tree = JSON.parse(readFileSync(join(dir, locale, file), 'utf8'))
    flatten(tree, area, out, problems, `${locale}/${file}`)
  }
  return { out, files }
}

// Typography rules from docs/i18n/glossary-*.md: guillemets for quotes, an ellipsis character,
// a spaced em dash and, in Ukrainian, the typographic apostrophe. Checked here because they are
// mechanical and easy to type wrong.
function typographyProblems(locale, key, texts) {
  const found = []
  for (const text of texts) {
    if (/"/.test(text)) found.push(`${locale} ${key}: straight double quote, use guillemets « »`)
    if (/\.\.\./.test(text)) found.push(`${locale} ${key}: three dots, use the character …`)
    if (/\S - \S/.test(text)) found.push(`${locale} ${key}: spaced hyphen, use the dash —`)
    if (FOREIGN_LETTERS[locale]?.test(text)) found.push(`${locale} ${key}: letters of another alphabet in ${JSON.stringify(text)}`)
    if (locale === 'uk' && /\p{L}'\p{L}/u.test(text)) found.push(`uk ${key}: straight apostrophe, use ’`)
  }
  return found
}

function compare(locale, en, tr, published, problems) {
  if (published) {
    for (const f of en.files) if (!tr.files.includes(f)) problems.push(`${locale}/${f} is missing`)
  }
  for (const f of tr.files) if (!en.files.includes(f)) problems.push(`${locale}/${f} has no English file`)

  let missingKeys = 0
  for (const [key, e] of en.out) {
    const r = tr.out.get(key)
    if (!r) {
      missingKeys += 1
      if (published) problems.push(`${locale} is missing ${key}`)
      continue
    }
    if (e.plural !== r.plural) {
      problems.push(`${locale} ${key}: plural in one language and plain text in the other`)
      continue
    }
    const enHolders = new Set(e.texts.flatMap(t => [...placeholders(t)]))
    const trHolders = new Set(r.texts.flatMap(t => [...placeholders(t)]))
    const missing = [...enHolders].filter(p => !trHolders.has(p))
    const extra = [...trHolders].filter(p => !enHolders.has(p))
    if (missing.length) problems.push(`${locale} ${key}: missing placeholder ${missing.map(p => `{${p}}`).join(', ')}`)
    if (extra.length) problems.push(`${locale} ${key}: unknown placeholder ${extra.map(p => `{${p}}`).join(', ')}`)
    if (e.plural) {
      for (const need of REQUIRED_FORMS[locale] ?? ['other']) {
        if (!r.forms.includes(need)) problems.push(`${locale} ${key}: plural lacks ${need}`)
      }
    }
    if (r.texts.some(t => t.trim() === '')) problems.push(`${locale} ${key}: empty text`)
    problems.push(...typographyProblems(locale, key, r.texts))
  }
  for (const key of tr.out.keys()) {
    if (!en.out.has(key)) problems.push(`${locale} has ${key}, which English does not`)
  }
  return missingKeys
}

/**
 * Compares every translated dictionary with English. `problems` must be empty; `progress`
 * reports how far each unpublished language has come.
 *
 * The real dictionaries take their language lists from lib/i18n/locale.ts. A fixture directory
 * (the tests) takes the languages it contains, all published, unless `options` says otherwise.
 */
export function checkAll(dir = defaultDir, options = {}) {
  const problems = []
  const progress = []
  const dirs = readdirSync(dir, { withFileTypes: true }).filter(d => d.isDirectory()).map(d => d.name)
  const real = dir === defaultDir
  const locales = options.locales ?? (real ? localeList('LOCALES') : dirs)
  const published = new Set(options.published ?? (real ? localeList('PUBLISHED_LOCALES') : locales))

  for (const p of published) if (!locales.includes(p)) problems.push(`PUBLISHED_LOCALES names ${p}, which LOCALES does not`)
  for (const d of dirs) if (!locales.includes(d)) problems.push(`messages/${d} is not a locale in LOCALES`)

  const en = load(dir, 'en', problems)
  for (const [key, e] of en.out) {
    if (e.texts.some(t => t.trim() === '')) problems.push(`en ${key}: empty text`)
    if (e.plural) {
      for (const need of REQUIRED_FORMS.en) if (!e.forms.includes(need)) problems.push(`en ${key}: plural lacks ${need}`)
      if (!e.texts.some(t => placeholders(t).has('count'))) problems.push(`en ${key}: plural message does not use {count}`)
    }
  }
  for (const locale of locales.filter(l => l !== 'en')) {
    if (!dirs.includes(locale)) {
      problems.push(`messages/${locale} does not exist`)
      continue
    }
    const tr = load(dir, locale, problems)
    const missing = compare(locale, en, tr, published.has(locale), problems)
    if (!published.has(locale)) {
      progress.push(`${locale}: ${en.out.size - missing}/${en.out.size} keys translated (not published)`)
    }
  }
  return { problems, progress }
}

/** Returns a list of human-readable problems; empty means the dictionaries agree. */
export function checkMessages(dir = defaultDir, options = {}) {
  return checkAll(dir, options).problems
}

if (process.argv[1] && fileURLToPath(import.meta.url) === process.argv[1]) {
  const { problems, progress } = checkAll()
  for (const line of progress) console.log(`i18n:check: ${line}`)
  if (problems.length) {
    console.error(`i18n:check found ${problems.length} problem(s):\n  ${problems.join('\n  ')}`)
    process.exit(1)
  }
  console.log('i18n:check: every dictionary agrees with English')
}
