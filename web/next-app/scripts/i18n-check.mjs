// Check that the Russian dictionary says the same things as the English one.
//
// A translation can be wrong in ways no type checker sees: a key that exists in one language
// and not the other renders as the raw key id, and a dropped `{count}` renders a sentence that
// reads fine and shows no number. This script compares the two dictionaries key by key and
// placeholder by placeholder.
//
// It is also a test (test/i18nMessages.test.ts imports `checkMessages`), so `npm test` fails on
// a mismatch and no separate CI step has to remember to run it.
//
// Usage: npm run i18n:check
import { readdirSync, readFileSync } from 'node:fs'
import { dirname, join } from 'node:path'
import { fileURLToPath } from 'node:url'

const here = dirname(fileURLToPath(import.meta.url))
const defaultDir = join(here, '..', 'lib', 'i18n', 'messages')

const PLURAL_FORMS = ['zero', 'one', 'two', 'few', 'many', 'other']
// Categories each language's rules can select (CLDR). A plural message must carry every one it
// needs, or a count that selects the missing one silently falls through to `other`.
const REQUIRED_FORMS = { en: ['one', 'other'], ru: ['one', 'few', 'many', 'other'] }

function isPlural(node) {
  if (node === null || typeof node !== 'object' || Array.isArray(node)) return false
  const keys = Object.keys(node)
  return keys.includes('other') && keys.every(k => PLURAL_FORMS.includes(k) && typeof node[k] === 'string')
}

function placeholders(text) {
  return new Set([...text.matchAll(/\{(\w+)\}/g)].map(m => m[1]))
}

/** Flatten one dictionary to `key -> { plural, texts[], forms[] }`. */
function flatten(tree, prefix, out, problems, where) {
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

function load(dir, locale, problems) {
  const out = new Map()
  const files = readdirSync(join(dir, locale)).filter(f => f.endsWith('.json')).sort()
  for (const file of files) {
    const area = file.replace(/\.json$/, '')
    const tree = JSON.parse(readFileSync(join(dir, locale, file), 'utf8'))
    flatten(tree, area, out, problems, `${locale}/${file}`)
  }
  return { out, files }
}

// Russian typography rules from docs/i18n/glossary-ru.md: guillemets for quotes, an ellipsis
// character, a spaced em dash. Checked here because they are mechanical and easy to type wrong.
function typographyProblems(key, texts) {
  const found = []
  for (const text of texts) {
    if (/"/.test(text)) found.push(`ru ${key}: straight double quote, use «ёлочки»`)
    if (/\.\.\./.test(text)) found.push(`ru ${key}: three dots, use the character …`)
    if (/\S - \S/.test(text)) found.push(`ru ${key}: spaced hyphen, use the dash —`)
  }
  return found
}

/** Returns a list of human-readable problems; empty means the dictionaries agree. */
export function checkMessages(dir = defaultDir) {
  const problems = []
  const en = load(dir, 'en', problems)
  const ru = load(dir, 'ru', problems)

  for (const f of en.files) if (!ru.files.includes(f)) problems.push(`ru/${f} is missing`)
  for (const f of ru.files) if (!en.files.includes(f)) problems.push(`en/${f} is missing`)

  for (const [key, e] of en.out) {
    const r = ru.out.get(key)
    if (!r) {
      problems.push(`ru is missing ${key}`)
      continue
    }
    if (e.plural !== r.plural) {
      problems.push(`${key}: plural in one language and plain text in the other`)
      continue
    }
    const enHolders = new Set(e.texts.flatMap(t => [...placeholders(t)]))
    const ruHolders = new Set(r.texts.flatMap(t => [...placeholders(t)]))
    const missing = [...enHolders].filter(p => !ruHolders.has(p))
    const extra = [...ruHolders].filter(p => !enHolders.has(p))
    if (missing.length) problems.push(`ru ${key}: missing placeholder ${missing.map(p => `{${p}}`).join(', ')}`)
    if (extra.length) problems.push(`ru ${key}: unknown placeholder ${extra.map(p => `{${p}}`).join(', ')}`)
    if (e.plural) {
      for (const need of REQUIRED_FORMS.en) if (!e.forms.includes(need)) problems.push(`en ${key}: plural lacks ${need}`)
      for (const need of REQUIRED_FORMS.ru) if (!r.forms.includes(need)) problems.push(`ru ${key}: plural lacks ${need}`)
      if (!enHolders.has('count')) problems.push(`en ${key}: plural message does not use {count}`)
    }
    if (r.texts.some(t => t.trim() === '')) problems.push(`ru ${key}: empty text`)
    problems.push(...typographyProblems(key, r.texts))
  }
  for (const key of ru.out.keys()) {
    if (!en.out.has(key)) problems.push(`en is missing ${key}`)
  }
  for (const [key, e] of en.out) {
    if (e.texts.some(t => t.trim() === '')) problems.push(`en ${key}: empty text`)
  }
  return problems
}

if (process.argv[1] && fileURLToPath(import.meta.url) === process.argv[1]) {
  const problems = checkMessages()
  if (problems.length) {
    console.error(`i18n:check found ${problems.length} problem(s):\n  ${problems.join('\n  ')}`)
    process.exit(1)
  }
  console.log('i18n:check: en and ru agree')
}
