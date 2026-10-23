// Detect translations whose English original has changed since they were written.
//
// A merge from upstream can reword an English sentence without touching its key: the key still
// exists, every check passes, and the Russian and Ukrainian text keep saying what the old English
// said. Nothing on screen looks wrong. This script records, per language, a fingerprint of the
// English text each translation was made from, in i18n/sources/<locale>.json at the repository
// root, and reports every translation whose English no longer matches its fingerprint.
//
// Covered: the web dictionaries, concept labels, course and destination text, lesson-step shards,
// lesson texts and the next-up card wording in internal/server/localize.go.
//
// Usage (from web/next-app):
//   npm run i18n:sources                       check every language
//   npm run i18n:sources -- stamp uk --all     record the current English for everything uk has
//   npm run i18n:sources -- stamp ru web:learn.page.loadFailed.title concepts:frac.add.diff
//
// Stamp only what you have just translated or reviewed against the new English: stamping is the
// statement "this translation says what the English says now".
//
// It is also a test (test/i18nSources.test.ts), so `npm test` fails on a stale translation.
import { createHash } from 'node:crypto'
import { existsSync, mkdirSync, readdirSync, readFileSync, statSync, writeFileSync } from 'node:fs'
import { dirname, join, relative, sep } from 'node:path'
import { fileURLToPath } from 'node:url'
import { load } from './i18n-check.mjs'

const here = dirname(fileURLToPath(import.meta.url))
const repo = join(here, '..', '..', '..')
const messagesDir = join(here, '..', 'lib', 'i18n', 'messages')
const sourcesDir = join(repo, 'i18n', 'sources')
const TRANSLATED = ['ru', 'uk']
// Lesson texts: a translation of data/lessons/<corpus>/<path>.md is <corpus>/<locale>/<path>.md.
const LESSON_CORPORA = ['teaching', 'algebrica', 'authored']
const LOCALE_DIRS = new Set(['en', ...TRANSLATED])

/** Lesson files under a corpus dir, relative to it with forward slashes. */
function lessonFiles(corpusDir, skipLocales) {
  const out = []
  const walk = dir => {
    for (const name of readdirSync(dir)) {
      const full = join(dir, name)
      if (statSync(full).isDirectory()) {
        if (skipLocales && dir === corpusDir && LOCALE_DIRS.has(name)) continue
        walk(full)
      } else if (name.endsWith('.md') && name !== 'README.md') {
        out.push(relative(corpusDir, full).split(sep).join('/'))
      }
    }
  }
  if (existsSync(corpusDir)) walk(corpusDir)
  return out
}

const readJSON = path => JSON.parse(readFileSync(path, 'utf8'))
const fingerprint = value => createHash('sha1').update(JSON.stringify(value)).digest('hex').slice(0, 12)

/** English originals, `section:key -> value`. */
export function englishSources() {
  const out = new Map()

  for (const [key, e] of load(messagesDir, 'en').out) out.set(`web:${key}`, e.texts)

  for (const file of readdirSync(join(repo, 'data', 'concepts'))) {
    if (!file.endsWith('.json') || file === 'enrichment.json') continue
    for (const c of readJSON(join(repo, 'data', 'concepts', file))) out.set(`concepts:${c.id}`, c.label)
  }
  for (const c of readJSON(join(repo, 'data', 'courses.json'))) out.set(`courses:${c.id}`, [c.name, c.description])
  for (const d of readJSON(join(repo, 'data', 'destinations.json'))) out.set(`destinations:${d.id}`, [d.name, d.description])

  const kpDir = join(repo, 'data', 'lessons', 'kp')
  for (const file of readdirSync(kpDir)) {
    if (!file.endsWith('.json')) continue
    const steps = readJSON(join(kpDir, file)).map(s => [s.label, s.section, s.subgoals])
    out.set(`kp:${file.replace(/\.json$/, '')}`, steps)
  }

  // A lesson's whole English text; line endings are normalised so a Windows checkout
  // fingerprints the same as any other.
  for (const corpus of LESSON_CORPORA) {
    const dir = join(repo, 'data', 'lessons', corpus)
    for (const rel of lessonFiles(dir, true)) {
      out.set(`lessons:${corpus}/${rel}`, readFileSync(join(dir, rel), 'utf8').replace(/\r\n/g, '\n'))
    }
  }

  // Each card is the source span from its badge to its call to action, so a change to any of
  // the three strings (or to the format of its detail) changes the fingerprint.
  for (const rel of ['internal/scheduler/recommend.go', 'internal/engine/recommend.go']) {
    const src = readFileSync(join(repo, rel), 'utf8').replace(/\r\n/g, '\n')
    for (const m of src.matchAll(/Badge\s*[:=]\s*"([^"]+)"[\s\S]*?CTA\s*[:=]\s*"[^"]*"/g)) {
      out.set(`cards:${m[1]}`, m[0].replace(/\s+/g, ' '))
    }
  }
  return out
}

/** The keys a language has translated, in the same `section:key` form. */
export function translatedKeys(locale) {
  const keys = new Set()
  if (existsSync(join(messagesDir, locale))) {
    for (const key of load(messagesDir, locale).out.keys()) keys.add(`web:${key}`)
  }
  const overlay = join(repo, 'data', 'i18n', locale)
  if (existsSync(join(overlay, 'concepts.json'))) {
    for (const id of Object.keys(readJSON(join(overlay, 'concepts.json')))) keys.add(`concepts:${id}`)
  }
  for (const section of ['courses', 'destinations']) {
    if (existsSync(join(overlay, `${section}.json`))) {
      for (const id of Object.keys(readJSON(join(overlay, `${section}.json`)))) keys.add(`${section}:${id}`)
    }
  }
  const kp = join(repo, 'data', 'lessons', 'kp', locale)
  if (existsSync(kp)) {
    for (const file of readdirSync(kp)) if (file.endsWith('.json')) keys.add(`kp:${file.replace(/\.json$/, '')}`)
  }
  for (const corpus of LESSON_CORPORA) {
    for (const rel of lessonFiles(join(repo, 'data', 'lessons', corpus, locale), false)) keys.add(`lessons:${corpus}/${rel}`)
  }
  const localize = readFileSync(join(repo, 'internal', 'server', 'localize.go'), 'utf8').replace(/\r\n/g, '\n')
  const block = localize.match(new RegExp(`i18n\\.${locale.toUpperCase()}: \\{\\n([\\s\\S]*?)\\n\\t\\},`))
  if (block) for (const m of block[1].matchAll(/^\t\t"([^"]+)":\s*\{/gm)) keys.add(`cards:${m[1]}`)
  return keys
}

function lockPath(locale) {
  return join(sourcesDir, `${locale}.json`)
}

function readLock(locale) {
  return existsSync(lockPath(locale)) ? readJSON(lockPath(locale)) : {}
}

function writeLock(locale, lock) {
  mkdirSync(sourcesDir, { recursive: true })
  const sorted = Object.fromEntries(Object.keys(lock).sort().map(k => [k, lock[k]]))
  writeFileSync(lockPath(locale), JSON.stringify(sorted, null, 2) + '\n')
}

/**
 * Per language: `stale` (the English changed since the translation was stamped) and `unstamped`
 * (translated, but no fingerprint recorded) must both be empty. `orphaned` lists fingerprints
 * of things no longer translated or no longer in English; they are dropped on the next stamp.
 */
export function checkSources() {
  const en = englishSources()
  const report = {}
  for (const locale of TRANSLATED) {
    const lock = readLock(locale)
    const have = translatedKeys(locale)
    const stale = []
    const unstamped = []
    for (const key of [...have].sort()) {
      if (!en.has(key)) continue // an orphaned translation is the other checks' business
      if (!(key in lock)) unstamped.push(key)
      else if (lock[key] !== fingerprint(en.get(key))) stale.push(key)
    }
    const orphaned = Object.keys(lock).filter(k => !have.has(k) || !en.has(k))
    report[locale] = { stale, unstamped, orphaned, translated: have.size }
  }
  return report
}

export function stamp(locale, which) {
  if (!TRANSLATED.includes(locale)) throw new Error(`unknown locale ${locale}`)
  const en = englishSources()
  const have = translatedKeys(locale)
  const lock = readLock(locale)
  for (const k of Object.keys(lock)) if (!have.has(k) || !en.has(k)) delete lock[k]
  const keys = which === 'all' ? [...have] : which
  let n = 0
  for (const key of keys) {
    if (!have.has(key)) throw new Error(`${locale} has no translation for ${key}`)
    if (!en.has(key)) throw new Error(`${key} has no English original`)
    lock[key] = fingerprint(en.get(key))
    n += 1
  }
  writeLock(locale, lock)
  return n
}

if (process.argv[1] && fileURLToPath(import.meta.url) === process.argv[1]) {
  const [cmd, locale, ...rest] = process.argv.slice(2)
  if (cmd === 'stamp') {
    const which = rest.includes('--all') ? 'all' : rest
    if (!locale || (which !== 'all' && which.length === 0)) {
      console.error('usage: i18n-sources stamp <locale> (--all | section:key ...)')
      process.exit(2)
    }
    console.log(`i18n:sources: stamped ${stamp(locale, which)} ${locale} translation(s)`)
  } else {
    let failed = false
    for (const [loc, r] of Object.entries(checkSources())) {
      console.log(`i18n:sources: ${loc}: ${r.translated} translated, ${r.stale.length} stale, ${r.unstamped.length} unstamped`)
      for (const k of r.stale) console.error(`  stale     ${loc} ${k}`)
      for (const k of r.unstamped) console.error(`  unstamped ${loc} ${k}`)
      if (r.orphaned.length) console.log(`  ${r.orphaned.length} orphaned fingerprint(s), dropped on the next stamp`)
      failed ||= r.stale.length > 0 || r.unstamped.length > 0
    }
    if (failed) {
      console.error('Retranslate the stale entries against the current English, then stamp them.')
      process.exit(1)
    }
  }
}
