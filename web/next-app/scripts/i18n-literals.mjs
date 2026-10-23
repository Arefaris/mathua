// Find English text a learner can see that does not go through the dictionaries.
//
// Every upstream merge brings new copy, and git merges a brand-new English sentence into a
// translated component without a conflict: the page simply shows one English line among Russian
// ones and nothing fails. This script is the alarm. It parses each component with the TypeScript
// compiler (so code is never mistaken for text) and reports:
//
//   - JSX text:                    <p>Could not load</p>
//   - text attributes:             aria-label="Report a problem", title=, placeholder=, alt=
//   - string literals in a JSX expression container: {'Select everything'}, {ok ? 'Yes' : 'No'}
//   - toast messages:              toast.error('Saved')
//
// A finding is fixed by moving the text into the dictionaries (en and every translation), or,
// when the English is deliberate (a brand, a licence-required credit), by listing it in
// scripts/i18n-literals.allow.json with the reason.
//
// It is also a test (test/i18nLiterals.test.ts), so `npm test` fails on a new literal.
//
// Usage: npm run i18n:literals
import { readdirSync, readFileSync, statSync } from 'node:fs'
import { dirname, join, relative, sep } from 'node:path'
import { fileURLToPath } from 'node:url'
import ts from 'typescript'

const here = dirname(fileURLToPath(import.meta.url))
const appRoot = join(here, '..')
const allowFile = join(here, 'i18n-literals.allow.json')

const ROOTS = ['app', 'components', 'lib']

// Out of scope for translation, by decision (see TRANSLATION_PLAN.md): staff tools, the long
// English articles, legal text and the i18n layer itself.
const SKIP = [
  'app/admin/',
  'app/docs/',
  'app/privacy/',
  'app/licenses/',
  'lib/i18n/',
  'components/ReportReplay.tsx',
  'components/LabQuestionCard.tsx',
  'components/AdminNav.tsx',
  'components/SystemDesignFlow/',
]

const TEXT_ATTRIBUTES = new Set(['aria-label', 'title', 'placeholder', 'alt', 'aria-description', 'aria-roledescription'])

// Text is reported only if it reads as English words: two or more Latin letters in a row, after
// removing what is not prose. Brand and product names and units stay English in every language.
const NOT_PROSE = [
  /\\\(.*?\\\)|\\\[.*?\\\]|\$[^$]*\$/g, // LaTeX
  /https?:\/\/\S+/g,
  /\b(Mathua|Math Academy|Discord|GitHub|Google|Apple|Microsoft|Facebook|KaTeX|MathLive|LaTeX|SymPy|OpenStax|Algebrica|DLMF|DiceBear|ORCCA|XP|LV|MIT|CC BY(-NC)?(-SA)?( \d\.\d)?|CSV|JSON|SVG|API|OK|ID)\b/g,
]
const ENGLISH = /[A-Za-z]{2,}/

// Elements whose content is code or markup, never copy.
const CODE_ELEMENTS = new Set(['style', 'script', 'code', 'pre', 'kbd', 'samp'])

// One token with no spaces is copy only if it reads like a word a person would show ("Reload");
// anything with code punctuation (an id like arith.add, a call, an address, a path) or in all
// lower case (an identifier like `difficulty`) is code shown as code.
function isCodeToken(text) {
  if (/\s/.test(text)) return false
  return /[._/@()=:<>{}[\]#$%^*\\|`~→]/.test(text) || /^[a-z0-9-]+$/.test(text)
}

// A formula or a status line shown as code: an operator, a snake_case identifier, or `key: value`.
function isFormula(text) {
  return /[=×^]|\w_\w|^[a-z]+: \S+$/.test(text)
}

function isProse(text) {
  if (isCodeToken(text) || isFormula(text)) return false
  let t = text
  for (const re of NOT_PROSE) t = t.replace(re, ' ')
  // ALL-CAPS tokens are constants (DECAYING), not words a learner reads.
  t = t.replace(/\b[A-Z][A-Z0-9_]+\b/g, ' ')
  return /[A-Za-z]*[a-z][A-Za-z]*/.test(t) && (/[a-z]{3,}/i.test(t) || /\b[A-Z][a-z]\b/.test(t))
}

function insideCodeElement(node) {
  for (let p = node.parent; p; p = p.parent) {
    if (ts.isJsxElement(p) && CODE_ELEMENTS.has(p.openingElement.tagName.getText())) return true
  }
  return false
}

function walk(dir, out) {
  for (const name of readdirSync(dir)) {
    const full = join(dir, name)
    if (statSync(full).isDirectory()) walk(full, out)
    else if (/\.(tsx|ts)$/.test(name) && !/\.(test|spec)\.tsx?$/.test(name) && !name.endsWith('.d.ts')) out.push(full)
  }
  return out
}

function scanFile(file) {
  return scanSource(readFileSync(file, 'utf8'), file)
}

/** The literals in one source text; exported for the checker's own tests. */
export function scanSource(src, file = 'source.tsx') {
  const sf = ts.createSourceFile(file, src, ts.ScriptTarget.Latest, true, file.endsWith('.tsx') ? ts.ScriptKind.TSX : ts.ScriptKind.TS)
  const found = []
  const report = (node, text) => {
    const clean = text.replace(/\s+/g, ' ').trim()
    if (!clean || !isProse(clean) || insideCodeElement(node)) return
    const { line } = sf.getLineAndCharacterOfPosition(node.getStart(sf))
    found.push({ line: line + 1, text: clean })
  }
  // String literals inside a JSX expression, following conditionals and `||`/`??` fallbacks,
  // but not into calls: `t('a.key')` is a key, not text.
  const literalsIn = (expr, cb) => {
    if (!expr) return
    if (ts.isStringLiteral(expr) || ts.isNoSubstitutionTemplateLiteral(expr)) cb(expr, expr.text)
    else if (ts.isTemplateExpression(expr)) cb(expr, expr.head.text + expr.templateSpans.map(s => ' ' + s.literal.text).join(''))
    else if (ts.isConditionalExpression(expr)) {
      literalsIn(expr.whenTrue, cb)
      literalsIn(expr.whenFalse, cb)
    } else if (ts.isBinaryExpression(expr) && [ts.SyntaxKind.BarBarToken, ts.SyntaxKind.QuestionQuestionToken, ts.SyntaxKind.AmpersandAmpersandToken].includes(expr.operatorToken.kind)) {
      literalsIn(expr.left, cb)
      literalsIn(expr.right, cb)
    } else if (ts.isParenthesizedExpression(expr)) literalsIn(expr.expression, cb)
  }
  const visit = node => {
    if (ts.isJsxText(node)) {
      report(node, node.text)
    } else if (ts.isJsxAttribute(node) && TEXT_ATTRIBUTES.has(node.name.getText(sf))) {
      const init = node.initializer
      if (init && ts.isStringLiteral(init)) report(init, init.text)
      else if (init && ts.isJsxExpression(init)) literalsIn(init.expression, report)
    } else if (ts.isJsxExpression(node) && node.parent && (ts.isJsxElement(node.parent) || ts.isJsxFragment(node.parent))) {
      literalsIn(node.expression, report)
    } else if (ts.isCallExpression(node)) {
      const callee = node.expression.getText(sf)
      if (/^toast(\.\w+)?$/.test(callee) && node.arguments.length) literalsIn(node.arguments[0], report)
    }
    ts.forEachChild(node, visit)
  }
  visit(sf)
  return found
}

function loadAllow() {
  try {
    return JSON.parse(readFileSync(allowFile, 'utf8'))
  } catch {
    return []
  }
}

/** Returns findings not covered by the allow list, and allow entries that matched nothing. */
export function findLiterals(root = appRoot) {
  const allow = loadAllow()
  const used = new Set()
  const findings = []
  for (const dir of ROOTS) {
    for (const file of walk(join(root, dir), [])) {
      const rel = relative(root, file).split(sep).join('/')
      if (SKIP.some(s => rel.startsWith(s))) continue
      for (const f of scanFile(file)) {
        const i = allow.findIndex(a => a.file === rel && a.text === f.text)
        if (i >= 0) {
          used.add(i)
          continue
        }
        findings.push({ file: rel, ...f })
      }
    }
  }
  const unused = allow.filter((_, i) => !used.has(i))
  return { findings, unused }
}

if (process.argv[1] && fileURLToPath(import.meta.url) === process.argv[1]) {
  const { findings, unused } = findLiterals()
  for (const a of unused) console.warn(`i18n:literals: allow entry no longer matches: ${a.file}: ${JSON.stringify(a.text)}`)
  if (findings.length) {
    console.error(`i18n:literals found ${findings.length} untranslated text(s):`)
    for (const f of findings) console.error(`  ${f.file}:${f.line}: ${JSON.stringify(f.text)}`)
    process.exit(1)
  }
  console.log('i18n:literals: no untranslated text in the learner-facing UI')
}
