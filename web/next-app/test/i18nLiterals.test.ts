import { describe, it, expect } from 'vitest'
import { findLiterals, scanSource } from '../scripts/i18n-literals.mjs'

// The real UI: every text a learner can see goes through the dictionaries. A failure here
// usually means an upstream merge brought in new English copy; move it into
// lib/i18n/messages (en and every translation), or list a deliberate exception in
// scripts/i18n-literals.allow.json.
describe('the learner-facing UI', () => {
  it('has no untranslated text', () => {
    const { findings, unused } = findLiterals()
    expect(findings).toEqual([])
    expect(unused).toEqual([])
  })
})

// The checker itself, against fixtures: a checker that never fails proves nothing.
describe('i18n-literals', () => {
  const texts = (src: string) => scanSource(src).map((f: { text: string }) => f.text)

  it('finds English copy in JSX text, text attributes, expressions and toasts', () => {
    const src = `
      export function A({ ok }) {
        toast.error('Could not save')
        return (
          <div aria-label="Report a problem" title={'Close the panel'}>
            <p>Could not load what to study next.</p>
            {ok ? 'Deselect everything' : 'Select everything'}
          </div>
        )
      }`
    expect(texts(src)).toEqual([
      'Could not save',
      'Report a problem',
      'Close the panel',
      'Could not load what to study next.',
      'Deselect everything',
      'Select everything',
    ])
  })

  it('leaves keys, code and names alone', () => {
    const src = `
      export function B({ t }) {
        return (
          <div className="font-mono text-xs" aria-label={t('learn.report.label')}>
            {t('common.status.loading')}
            <style>{'.a { color: red; }'}</style>
            <code>simplify(expected - answer) == 0</code>
            <span>arith.add.single</span>
            <span>difficulty</span>
            <span>= (0.7 × days_since_last_seen)</span>
            <span>Mathua</span>
            <span>you@example.com</span>
          </div>
        )
      }`
    expect(texts(src)).toEqual([])
  })
})
