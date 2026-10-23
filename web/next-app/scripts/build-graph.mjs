// Build the landing graph artifacts from the source concept corpus.
//
// Expensive work (layout) happens here once per build, never on the client:
//   public/data/concepts.positions.json  quantized [x,y,z] map (poster + scene)
//   public/data/concepts.graph.json      nodes + indexed edges + domain counts
//   data/graph.meta.json                 tiny static stats for page sections
//
// Translated labels, when the content side has supplied them (see label-overlay.mjs):
//   public/data/concepts.graph.ru.json   the same payload with Russian labels
//   data/concepts.labels.ru.json         flat id -> label, bundled by lib/conceptData.ts
// With no overlay the first is removed and the second is `{}`, so a checkout without
// translations builds what it always built.
//
// Usage: npm run graph:build
import { readFileSync, writeFileSync, mkdirSync, rmSync } from 'node:fs'
import { dirname, join } from 'node:path'
import { fileURLToPath } from 'node:url'
import { layoutDAG3D } from './layoutDAG3D.ts'
import { applyLabelOverlay, buildGraphMeta, buildGraphPayload, graphPayloadFile } from '../lib/graphPayload.ts'
import { readLabelOverlay } from './label-overlay.mjs'
import { DOMAIN_ORDER } from '../lib/graphDomains.ts'

const root = join(dirname(fileURLToPath(import.meta.url)), '..')
const concepts = JSON.parse(readFileSync(join(root, 'data', 'concepts.json'), 'utf8'))
// The app lives in web/next-app, two levels below the repository root that holds data/i18n.
const overlayDir = join(root, '..', '..', 'data', 'i18n')
// Every translated language in lib/i18n/locale.ts LOCALES; lib/conceptData.ts imports one
// labels file per language, so each is written even when its overlay does not exist yet.
const TRANSLATED_LOCALES = ['ru', 'uk']

const t0 = performance.now()
const input = concepts.map(c => ({
  id: c.id,
  prerequisites: c.prerequisites ?? [],
  domain: c.domain,
}))
const { positions, importance } = layoutDAG3D(input, { domainOrder: [...DOMAIN_ORDER] })
const layoutMs = performance.now() - t0

const round2 = n => Math.round(n * 100) / 100
const quantized = {}
for (const [id, p] of Object.entries(positions)) {
  quantized[id] = [round2(p[0]), round2(p[1]), round2(p[2])]
}

const payload = buildGraphPayload(concepts, importance)
const meta = buildGraphMeta(concepts)

mkdirSync(join(root, 'public', 'data'), { recursive: true })
writeFileSync(join(root, 'public', 'data', 'concepts.positions.json'), JSON.stringify(quantized) + '\n')
writeFileSync(join(root, 'public', 'data', 'concepts.graph.json'), JSON.stringify(payload) + '\n')
writeFileSync(join(root, 'data', 'graph.meta.json'), JSON.stringify(meta, null, 2) + '\n')

// Translated labels. Positions are keyed by id and do not change with the language, so the
// positions file is shared.
const known = new Set(concepts.map(c => c.id))
const summary = []
for (const locale of TRANSLATED_LOCALES) {
  const overlay = readLabelOverlay(join(overlayDir, locale, 'concepts.json'))
  const graphPath = join(root, 'public', 'data', graphPayloadFile(locale))
  const labels = {}
  if (overlay) {
    for (const [id, label] of Object.entries(overlay)) {
      if (known.has(id) && label.trim() !== '') labels[id] = label
    }
    const stale = Object.keys(overlay).filter(id => !known.has(id))
    if (stale.length) console.warn(`graph:build ${locale}: ${stale.length} overlay ids are not in the corpus (first: ${stale[0]})`)
    const localized = buildGraphPayload(applyLabelOverlay(concepts, labels), importance)
    writeFileSync(graphPath, JSON.stringify(localized) + '\n')
    summary.push(`${locale} labels for ${Object.keys(labels).length}`)
  } else {
    // A stale artifact from an earlier overlay would keep shipping labels nobody maintains.
    rmSync(graphPath, { force: true })
    summary.push(`no ${locale} overlay`)
  }
  writeFileSync(join(root, 'data', `concepts.labels.${locale}.json`), JSON.stringify(labels) + '\n')
}

console.log(
  `graph:build ${concepts.length} nodes / ${payload.edges.length} edges in ${layoutMs.toFixed(0)}ms ` +
    `(layout), ${Object.keys(quantized).length} positions -> public/data/ + data/graph.meta.json; ` +
    summary.join('; ')
)
