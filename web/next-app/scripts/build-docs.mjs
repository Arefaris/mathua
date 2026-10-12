// Build the docs content module from the repository's markdown.
//
// The website's documentation pages render the same markdown a contributor reads in the repo,
// so there is one source and the pages cannot drift from it. Before this, the pages were
// hand-written paraphrases and had drifted on exactly the facts a reader would check first:
// the in-app architecture page called Postgres "a stub" and the diagnostic "binary search".
//
//   content/docs.generated.json   slug -> { title, description, source, body }
//
// Usage:
//   npm run docs:build    regenerate the module
//   npm run docs:check    fail if the module is out of date (CI)
//
// The output is committed, like web/next-app/data/concepts.json, because the type checker and
// the unit tests import it and they run before a build. Committing generated content is only
// safe with the check: `docs:check` regenerates in memory and diffs, so an edit to a repo doc
// that was not regenerated fails the build rather than shipping a stale page.
import { readFileSync, writeFileSync, mkdirSync, existsSync } from 'node:fs'
import { dirname, join } from 'node:path'
import { fileURLToPath } from 'node:url'

const here = dirname(fileURLToPath(import.meta.url))
const appRoot = join(here, '..')
const repoRoot = join(appRoot, '..', '..')
const outFile = join(appRoot, 'content', 'docs.generated.json')

// The docs the website publishes, and the metadata each page needs.
//
// This list is the boundary: a doc not named here does not reach the web. code-defects.md,
// lesson-media.md, mastery-semantics.md and the mobile checklist stay in the repo because they
// describe known defects, corpus audits and internal design, not the product.
//
// `description` is written here rather than derived, because a meta description is editorial
// and the first sentence of an engineering doc is usually not one.
const DOCS = [
  {
    slug: 'architecture',
    source: 'docs/architecture.md',
    description: 'The five-layer stack — UI, API, core engine, grading and storage — and the computerised adaptive diagnostic.',
  },
  {
    slug: 'administration',
    source: 'docs/administration.md',
    description: 'Roles, the permission model, staff invitations, the last-owner invariant and the audit trail.',
  },
  {
    slug: 'system-design',
    source: 'docs/system-design.md',
    description: 'Full system design: schema, scheduler, SM-2, grading, the API surface, and the trade-offs behind each.',
  },
  {
    slug: 'math-input',
    source: 'docs/math-input.md',
    description: 'How a typed answer reaches a verdict: MathLive, MathJSON, the plain-text shim, and the mobile keypad.',
  },
  {
    slug: 'efficacy',
    source: 'docs/efficacy.md',
    description: 'How Mathua measures learning, and the honest limits of those measurements.',
  },
  {
    slug: 'licenses',
    source: 'THIRD_PARTY_NOTICES.md',
    description: 'Every third-party library, content source and font Mathua ships, with its license.',
  },
]

function titleOf(markdown, fallback) {
  const m = markdown.match(/^#\s+(.+)$/m)
  return m ? m[1].trim() : fallback
}

// Drop the leading H1. The page renders its own heading, and two titles on one page is worse
// than one.
function withoutTitle(markdown) {
  return markdown.replace(/^#\s+.*(\r?\n)+/, '')
}

export function build() {
  const out = {}
  const missing = []
  for (const doc of DOCS) {
    const path = join(repoRoot, doc.source)
    if (!existsSync(path)) {
      // A missing source is fatal, not skipped: a silently absent page is the failure this
      // whole arrangement exists to prevent.
      missing.push(doc.source)
      continue
    }
    const raw = readFileSync(path, 'utf8')
    out[doc.slug] = {
      title: titleOf(raw, doc.slug),
      description: doc.description,
      source: doc.source,
      body: withoutTitle(raw).trimEnd() + '\n',
    }
  }
  if (missing.length) {
    throw new Error(
      `docs source not found: ${missing.join(', ')}. ` +
      `Every name in build-docs.mjs must exist; there is no silent skip.`,
    )
  }
  return out
}

function serialise(data) {
  return JSON.stringify(data, null, 2) + '\n'
}

const check = process.argv.includes('--check')
const next = serialise(build())

if (check) {
  const current = existsSync(outFile) ? readFileSync(outFile, 'utf8') : ''
  if (current !== next) {
    console.error(
      'FAIL: content/docs.generated.json is out of date with the repo docs.\n' +
      '      Run `npm run docs:build` and commit the result.',
    )
    process.exit(1)
  }
  console.log(`OK: docs content matches ${DOCS.length} repo source(s)`)
} else {
  mkdirSync(dirname(outFile), { recursive: true })
  writeFileSync(outFile, next)
  console.log(`OK: wrote ${outFile} from ${DOCS.length} source(s)`)
}
