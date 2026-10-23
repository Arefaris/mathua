// How much of Mathua is translated, per language, and the progress block in README.md.
//
// Progress is measured in English words, not in files or keys: a course name and a lesson of
// two thousand words are not the same amount of work, and counting them as one item each would
// report the project nearly done while every lesson is still English. Formulas, code and markup
// are not words and are not counted.
//
// Areas, and how a translation of each is found:
//   Interface        web dictionaries and the next-up cards     (see i18n-sources.mjs)
//   Topic names      data/i18n/<locale>/concepts.json
//   Courses          data/i18n/<locale>/courses.json, destinations.json
//   Lesson steps     data/lessons/kp/<locale>/<id>.json
//   Lesson texts     data/lessons/<corpus>/<locale>/<same path>.md   (teaching, algebrica, authored).
//                    Lessons adapted from OpenStax are not translated (OpenStax forbids feeding its
//                    books to generative AI without permission; see CLAUDE.md), so they are left
//                    out of the total and reported apart.
//   Worked solutions data/lessons/solutions/<locale>/<id>.json
//   Problem texts    the Question and Explanation templates in internal/generator. They have no
//                    translation mechanism yet (TRANSLATION_PLAN.md stage 4), so they count as 0
//                    until it lands; the counter is wired here when it does.
//
// Usage (from web/next-app):
//   npm run i18n:progress            rewrite the block in README.md
//   npm run i18n:progress -- --check fail if the block is out of date (also a test)
import { existsSync, readdirSync, readFileSync, statSync, writeFileSync } from 'node:fs'
import { dirname, join, relative, sep } from 'node:path'
import { fileURLToPath } from 'node:url'
import { englishSources, translatedKeys } from './i18n-sources.mjs'

const here = dirname(fileURLToPath(import.meta.url))
const repo = join(here, '..', '..', '..')
const readme = join(repo, 'README.md')
const LOCALES = [
  { code: 'uk', name: 'Ukrainian' },
  { code: 'ru', name: 'Russian' },
]
const START = '<!-- i18n-progress:start -->'
const END = '<!-- i18n-progress:end -->'

/** English words in a text, with formulas, code, markup and links removed. */
export function countWords(text) {
  const prose = String(text)
    .replace(/```[\s\S]*?```/g, ' ')
    .replace(/`[^`]*`/g, ' ')
    .replace(/\$\$[\s\S]*?\$\$|\$[^$\n]*\$|\\\([\s\S]*?\\\)|\\\[[\s\S]*?\\\]/g, ' ')
    .replace(/<[^>]+>/g, ' ')
    .replace(/!?\[([^\]]*)\]\([^)]*\)/g, '$1')
    .replace(/https?:\/\/\S+/g, ' ')
    .replace(/%[sdvqf]|\{\w+\}/g, ' ')
  return (prose.match(/[A-Za-z]+(?:'[a-z]+)?/g) ?? []).length
}

function walk(dir, out = []) {
  for (const name of readdirSync(dir)) {
    const full = join(dir, name)
    if (statSync(full).isDirectory()) walk(full, out)
    else out.push(full)
  }
  return out
}

const isLocaleDir = (corpusDir, file) => LOCALES.some(l => relative(corpusDir, file).split(sep)[0] === l.code)

/** Areas measured through the fingerprinted sources: `section:key` -> words. */
function sourceAreas() {
  const words = new Map()
  for (const [key, value] of englishSources()) {
    // A card's source is a span of Go code; only its string literals are copy.
    const text = key.startsWith('cards:') ? [...value.matchAll(/"([^"]*)"/g)].map(m => m[1]).join(' ') : JSON.stringify(value)
    words.set(key, countWords(text))
  }
  return words
}

// The attribution line an OpenStax-derived lesson opens with (the same test as in
// internal/i18n/overlay_data_test.go).
const OPENSTAX = /^> Content (sourced|adapted) from \[OpenStax/m

const AREA_OF = { web: 'Interface', cards: 'Interface', concepts: 'Topic names', courses: 'Courses', destinations: 'Courses', kp: 'Lesson steps' }
const AREAS = ['Interface', 'Topic names', 'Courses', 'Lesson steps', 'Lesson texts', 'Worked solutions', 'Problem texts']

/** `{ area: { total, [locale]: translated } }`, all in words. */
export function measure() {
  const result = Object.fromEntries(AREAS.map(a => [a, Object.fromEntries([['total', 0], ...LOCALES.map(l => [l.code, 0])])]))
  const add = (area, words, translatedIn) => {
    result[area].total += words
    for (const l of LOCALES) if (translatedIn(l.code)) result[area][l.code] += words
  }

  const words = sourceAreas()
  const have = Object.fromEntries(LOCALES.map(l => [l.code, translatedKeys(l.code)]))
  for (const [key, n] of words) {
    const area = AREA_OF[key.split(':')[0]]
    if (area) add(area, n, loc => have[loc].has(key))
  }

  const lessons = join(repo, 'data', 'lessons')
  result.excluded = { files: 0, words: 0 }
  for (const corpus of ['teaching', 'algebrica', 'authored']) {
    const dir = join(lessons, corpus)
    for (const file of walk(dir)) {
      if (!file.endsWith('.md') || file.endsWith('README.md') || isLocaleDir(dir, file)) continue
      const text = readFileSync(file, 'utf8')
      if (OPENSTAX.test(text)) {
        result.excluded.files += 1
        result.excluded.words += countWords(text)
        continue
      }
      const rel = relative(dir, file)
      add('Lesson texts', countWords(text), loc => existsSync(join(dir, loc, rel)))
    }
  }

  const solutions = join(lessons, 'solutions')
  for (const name of readdirSync(solutions)) {
    if (!name.endsWith('.json')) continue
    const strings = JSON.stringify(JSON.parse(readFileSync(join(solutions, name), 'utf8')), (_, v) => v)
    const text = [...strings.matchAll(/"((?:[^"\\]|\\.)*)"/g)].map(m => m[1]).join(' ')
    add('Worked solutions', countWords(text), loc => existsSync(join(solutions, loc, name)))
  }

  for (const file of walk(join(repo, 'internal', 'generator'))) {
    if (!file.endsWith('.go') || file.endsWith('_test.go')) continue
    for (const line of readFileSync(file, 'utf8').split('\n')) {
      if (!/\b(Question|Explanation)\s*:/.test(line)) continue
      const text = [...line.matchAll(/"((?:[^"\\]|\\.)*)"/g)].map(m => m[1]).join(' ')
      add('Problem texts', countWords(text), () => false)
    }
  }
  return result
}

const pct = (part, total) => (total === 0 ? 0 : (100 * part) / total)

function label(p) {
  if (p === 0) return '0%'
  if (p < 1) return '<1%'
  if (p > 99 && p < 100) return '99%'
  return `${Math.floor(p)}%`
}

function bar(p) {
  const width = 20
  const filled = p > 0 && p < 100 / width ? 1 : Math.floor((p / 100) * width)
  return '█'.repeat(filled) + '░'.repeat(width - filled)
}

function badgeColor(p) {
  if (p >= 90) return 'brightgreen'
  if (p >= 50) return 'green'
  if (p >= 20) return 'yellow'
  if (p >= 5) return 'orange'
  return 'red'
}

/** The Markdown between the markers in README.md. */
export function renderBlock(result = measure()) {
  const totals = Object.fromEntries(['total', ...LOCALES.map(l => l.code)].map(k => [k, AREAS.reduce((s, a) => s + result[a][k], 0)]))
  const overall = Object.fromEntries(LOCALES.map(l => [l.code, pct(totals[l.code], totals.total)]))
  const badges = LOCALES.map(
    l => `![${l.name}](https://img.shields.io/badge/${l.name}-${encodeURIComponent(label(overall[l.code]))}-${badgeColor(overall[l.code])}?style=flat-square)`,
  ).join(' ')
  const cell = (part, total) => `\`${bar(pct(part, total))}\` ${label(pct(part, total))}`
  const rows = AREAS.map(a => `| ${a} | ${result[a].total.toLocaleString('en-US')} | ${LOCALES.map(l => cell(result[a][l.code], result[a].total)).join(' | ')} |`)
  return [
    START,
    badges,
    '',
    `| Area | English words | ${LOCALES.map(l => l.name).join(' | ')} |`,
    `|---|--:|---|---|`,
    ...rows,
    `| **Overall** | **${totals.total.toLocaleString('en-US')}** | ${LOCALES.map(l => `**${cell(totals[l.code], totals.total)}**`).join(' | ')} |`,
    '',
    ...(result.excluded?.files
      ? [
          `_Not counted: ${result.excluded.files} lessons adapted from OpenStax (${result.excluded.words.toLocaleString('en-US')} words). ` +
            'OpenStax does not allow its books to be fed to generative AI without permission, so these lessons stay in English._',
          '',
        ]
      : []),
    '_Measured in English words, without formulas and code. Regenerate with `npm run i18n:progress` in `web/next-app`._',
    END,
  ].join('\n')
}

/** Returns the README with the block replaced, or throws when the markers are missing. */
export function updatedReadme(text = readFileSync(readme, 'utf8'), block = renderBlock()) {
  const nl = text.includes('\r\n') ? '\r\n' : '\n'
  const s = text.indexOf(START)
  const e = text.indexOf(END)
  if (s < 0 || e < s) throw new Error(`README.md has no ${START} … ${END} block`)
  return text.slice(0, s) + block.replace(/\n/g, nl) + text.slice(e + END.length)
}

export function readmeIsCurrent() {
  const text = readFileSync(readme, 'utf8')
  return updatedReadme(text) === text
}

if (process.argv[1] && fileURLToPath(import.meta.url) === process.argv[1]) {
  if (process.argv.includes('--check')) {
    if (!readmeIsCurrent()) {
      console.error('i18n:progress: the README progress block is out of date; run npm run i18n:progress')
      process.exit(1)
    }
    console.log('i18n:progress: README is current')
  } else {
    writeFileSync(readme, updatedReadme())
    console.log(renderBlock())
  }
}
