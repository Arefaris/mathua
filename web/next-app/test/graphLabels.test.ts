import { describe, it, expect, vi } from 'vitest'
import { mkdtempSync, writeFileSync, rmSync } from 'node:fs'
import { tmpdir } from 'node:os'
import { join } from 'node:path'
import { applyLabelOverlay, buildGraphPayload, graphPayloadFile } from '../lib/graphPayload'
import type { ConceptRecord } from '../lib/conceptData'
import { parseLabelOverlay, readLabelOverlay } from '../scripts/label-overlay.mjs'

const corpus: ConceptRecord[] = [
  { id: 'a.one', label: 'One', domain: 'arithmetic', prerequisites: [] },
  { id: 'a.two', label: 'Two', domain: 'arithmetic', prerequisites: ['a.one'] },
]

describe('the label overlay', () => {
  it('is optional: a missing file reads as no overlay, not as an error', () => {
    expect(readLabelOverlay(join(tmpdir(), 'mathua-no-such-dir', 'concepts.json'))).toBeNull()
  })

  it('reads a flat id -> label file', () => {
    const dir = mkdtempSync(join(tmpdir(), 'overlay-'))
    try {
      const file = join(dir, 'concepts.json')
      writeFileSync(file, JSON.stringify({ 'a.one': 'Один' }))
      expect(readLabelOverlay(file)).toEqual({ 'a.one': 'Один' })
    } finally {
      rmSync(dir, { recursive: true, force: true })
    }
  })

  it('rejects a file that exists but is wrong, so a broken overlay cannot ship silently', () => {
    expect(() => parseLabelOverlay('{nope')).toThrow(/not valid JSON/)
    expect(() => parseLabelOverlay('[]')).toThrow(/expected an object/)
    expect(() => parseLabelOverlay('{"a.one": 3}')).toThrow(/a\.one is not a string/)
  })

  it('swaps labels and keeps the English one where the overlay is silent or blank', () => {
    const ru = applyLabelOverlay(corpus, { 'a.one': 'Один', 'a.two': '  ' })
    expect(ru.map(c => c.label)).toEqual(['Один', 'Two'])
    // The input is not mutated.
    expect(corpus[0].label).toBe('One')
  })

  it('leaves ids, domains and edges untouched, so positions can be shared', () => {
    const en = buildGraphPayload(corpus, { 'a.one': 1, 'a.two': 2 })
    const ru = buildGraphPayload(applyLabelOverlay(corpus, { 'a.one': 'Один' }), { 'a.one': 1, 'a.two': 2 })
    expect(ru.edges).toEqual(en.edges)
    expect(ru.nodes.map(n => n[0])).toEqual(en.nodes.map(n => n[0]))
    expect(ru.nodes[0][1]).toBe('Один')
  })

  it('names the artifact for each language', () => {
    expect(graphPayloadFile('en')).toBe('concepts.graph.json')
    expect(graphPayloadFile('ru')).toBe('concepts.graph.ru.json')
  })
})

describe('bundled concept labels follow the locale', () => {
  it('reads the overlay when the locale is Russian and falls back to English per concept', async () => {
    vi.resetModules()
    vi.doMock('../data/concepts.labels.ru.json', () => ({ default: { 'frac.add.word': 'Сложение дробей: задачи' } }))
    // Reset modules gives conceptData its own copy of the locale store, so drive that copy.
    const data = await import('../lib/conceptData')
    const { resetLocaleForTests } = await import('../lib/i18n/locale')
    const known = data.concepts.find(c => c.id === 'frac.add.word')!
    const other = data.concepts.find(c => c.id !== 'frac.add.word')!
    const english = known.label
    resetLocaleForTests('ru')
    expect(known.label).toBe('Сложение дробей: задачи')
    // No translation for this one: English, not an empty string and not the id.
    expect(other.label).not.toBe('')
    expect(other.label).not.toBe(other.id)
    resetLocaleForTests('en')
    expect(known.label).toBe(english)
    vi.doUnmock('../data/concepts.labels.ru.json')
  })
})
