import { describe, it, expect, vi, beforeAll } from 'vitest'
import { act, fireEvent, render, screen, waitFor } from '@testing-library/react'
import HomePage from '../app/page'
import DomainTable from '../components/DomainTable'
import GlossarySearch from '../components/GlossarySearch'
import GlossaryEntryCard from '../components/GlossaryEntry'
import { ConceptGraphSection, StudentModelSection } from '../app/how-it-works/sections'
import { resetLocaleForTests } from '../lib/i18n/locale'
import { getTranslator } from '../lib/i18n'
import { GLOSSARY_ENTRIES } from '../lib/glossary/data'

// Russian rendering of the landing, How it works, the domain table, the glossary chrome and the
// concept graph. The English suite pins the locale to en, so each test switches it first.

const reactFlowMocks = vi.hoisted(() => ({
  zoomIn: vi.fn(),
  zoomOut: vi.fn(),
  fitView: vi.fn(),
  getViewport: vi.fn(() => ({ x: 0, y: 0, zoom: 1 })),
  setViewport: vi.fn(),
  setCenter: vi.fn(),
}))

vi.mock('@xyflow/react', async importOriginal => {
  const actual = await importOriginal<typeof import('@xyflow/react')>()
  return { ...actual, useReactFlow: () => reactFlowMocks }
})

vi.mock('next/navigation', () => ({
  useRouter: () => ({ push: vi.fn(), replace: vi.fn() }),
  usePathname: () => '/',
  useSearchParams: () => new URLSearchParams(),
}))

vi.mock('../lib/auth', async importOriginal => {
  const mod = await importOriginal<typeof import('../lib/auth')>()
  return { ...mod, ensureGuestId: vi.fn(), ensureGuestToken: vi.fn(async () => 'guest-token') }
})

import ConceptGraphFlow from '../components/ConceptGraphFlow'

beforeAll(() => {
  Object.defineProperty(globalThis, 'ResizeObserver', {
    writable: true,
    value: class {
      observe() {}
      unobserve() {}
      disconnect() {}
    },
  })
  Object.defineProperty(globalThis, 'DOMMatrixReadOnly', {
    writable: true,
    value: class {
      m22 = 1
    },
  })
})

const noKeyIds = () => expect(document.body.textContent ?? '').not.toMatch(/\b(explore|common)\.[a-z]/)

describe('landing in Russian', () => {
  it('shows the hero, the domains and the level names in Russian', async () => {
    resetLocaleForTests('ru')
    localStorage.clear()
    render(<HomePage />)
    expect(await screen.findByRole('button', { name: 'Начать' })).toBeInTheDocument()
    expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent(
      'Адаптивная платформа для изучения математики',
    )
    // The coverage table uses the shared domain names (glossary section 5).
    expect(screen.getAllByText('Дроби').length).toBeGreaterThan(0)
    expect(screen.getAllByText('Введение в алгебру').length).toBeGreaterThan(0)
    expect(screen.getByText('Подмастерье')).toBeInTheDocument()
    // Plural link text: "Все 1250 тем по разделам".
    expect(screen.getByText(/^Все \d+ тем по разделам →$/)).toBeInTheDocument()
    // State chips are uppercase Russian state names.
    expect(screen.getByText('ЗАКРЕПЛЯЕТСЯ')).toBeInTheDocument()
    noKeyIds()
  })
})

describe('How it works in Russian', () => {
  it('keeps the numbers and picks the plural forms', () => {
    resetLocaleForTests('ru')
    render(
      <>
        <ConceptGraphSection />
        <StudentModelSection />
      </>,
    )
    expect(screen.getByText(/^Сейчас граф включает \d+ тем, которые охватывают \d+ раздел/)).toBeInTheDocument()
    // 5 and 7 take the genitive plural, 8 / 18 / 20 seconds are "секунд".
    expect(screen.getAllByText('5 верных ответов подряд').length).toBeGreaterThan(0)
    expect(screen.getAllByText('8 секунд').length).toBeGreaterThan(0)
    expect(screen.getAllByText('Сложение однозначных чисел').length).toBeGreaterThan(0)
    noKeyIds()
  })

  it('declines "entries" like a Russian noun', () => {
    const t = getTranslator('ru')
    expect(t('explore.glossary.search.entries', { count: 1 })).toBe('1 запись')
    expect(t('explore.glossary.search.entries', { count: 2 })).toBe('2 записи')
    expect(t('explore.glossary.search.entries', { count: 5 })).toBe('5 записей')
    expect(t('explore.glossary.search.entries', { count: 11 })).toBe('11 записей')
    expect(t('explore.glossary.search.entries', { count: 21 })).toBe('21 запись')
    expect(t('explore.glossary.search.entries', { count: 22 })).toBe('22 записи')
  })
})

describe('domain table in Russian', () => {
  const rows = [
    {
      domain: 'fractions',
      label: 'Дроби',
      total: 12,
      mastered: 3,
      dueForReview: 2,
      inProgress: 1,
      nextId: 'frac.add.same',
      nextLabel: 'Сложение дробей',
      concepts: [{ id: 'frac.add.same', label: 'Сложение дробей', status: 'dueForReview' as const }],
    },
  ]

  it('translates the counts, the call to action and the status groups', () => {
    resetLocaleForTests('ru')
    render(<DomainTable rows={rows} onToggle={() => {}} expanded="fractions" />)
    expect(screen.getAllByText(/пора повторить: 2/).length).toBeGreaterThan(0)
    expect(screen.getByText('Продолжить: Сложение дробей')).toBeInTheDocument()
    expect(screen.getByText('Скрыть темы')).toBeInTheDocument()
    expect(screen.getByText('Пора повторить (1)')).toBeInTheDocument()
    noKeyIds()
  })

  it('counts the concepts with the shared plural', () => {
    resetLocaleForTests('ru')
    render(<DomainTable rows={rows} onToggle={() => {}} expanded={null} />)
    expect(screen.getByText('1 тема')).toBeInTheDocument()
  })
})

describe('glossary chrome in Russian', () => {
  it('translates the controls, the categories and the result count', () => {
    resetLocaleForTests('ru')
    render(
      <GlossarySearch
        query=""
        category="all"
        resultCount={21}
        rangeStart={1}
        rangeEnd={21}
        mode="text"
        mathLatex=""
        onQueryChange={() => {}}
        onCategoryChange={() => {}}
        onModeChange={() => {}}
        onMathChange={() => {}}
        onClear={() => {}}
      />,
    )
    expect(screen.getByPlaceholderText(/^Символ, название или значение/)).toBeInTheDocument()
    expect(screen.getByText('21 запись')).toBeInTheDocument()
    fireEvent.click(screen.getByRole('button', { name: 'Теги' }))
    expect(screen.getByRole('button', { name: 'Греческие буквы' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Теория чисел и высшая математика' })).toBeInTheDocument()
    noKeyIds()
  })

  it('shows a page range with the genitive plural', () => {
    resetLocaleForTests('ru')
    render(
      <GlossarySearch
        query=""
        category="all"
        resultCount={52}
        rangeStart={21}
        rangeEnd={40}
        mode="text"
        mathLatex=""
        onQueryChange={() => {}}
        onCategoryChange={() => {}}
        onModeChange={() => {}}
        onMathChange={() => {}}
        onClear={() => {}}
      />,
    )
    expect(screen.getByText('Показано 21–40 из 52 записей')).toBeInTheDocument()
  })

  it('frames an entry in Russian and leaves its content alone', () => {
    resetLocaleForTests('ru')
    const entry = GLOSSARY_ENTRIES.find(e => e.commonMistakes && e.commonMistakes.length > 0)!
    render(<GlossaryEntryCard entry={entry} onSelectRelated={() => {}} />)
    expect(screen.getByText('Типичная ошибка:')).toBeInTheDocument()
    expect(screen.getByText(entry.definition)).toBeInTheDocument()
  })
})

describe('concept graph in Russian', () => {
  const concepts = [
    { id: 'count', label: 'Счёт', domain: 'arithmetic', prerequisites: [] },
    { id: 'add', label: 'Сложение', domain: 'fractions', prerequisites: ['count'] },
  ]

  it('translates the controls and the node description', async () => {
    resetLocaleForTests('ru')
    const { container } = render(
      <ConceptGraphFlow
        concepts={concepts}
        conceptStatuses={{ count: 'mastered', add: 'learning' }}
        conceptProgress={{ add: 0.4 }}
      />,
    )
    expect(screen.getByPlaceholderText('Найти тему…')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Открыть список тем' })).toBeInTheDocument()
    const node = container.querySelector('.concept-node[aria-label^="Сложение"]')
    expect(node?.getAttribute('aria-label')).toBe('Сложение, раздел «Дроби», Изучается, 40% до освоения')
    expect(node?.querySelector('.node-progress-tip')?.textContent).toBe('40% до освоения')
    fireEvent.click(screen.getByRole('button', { name: 'Открыть список тем' }))
    await waitFor(() => expect(screen.getByRole('dialog', { name: 'Список тем' })).toBeInTheDocument())
    noKeyIds()
  })

  it('shows the empty state in Russian', () => {
    resetLocaleForTests('ru')
    render(<ConceptGraphFlow concepts={[]} />)
    expect(screen.getByText('Нет тем для показа')).toBeInTheDocument()
  })
})

// A switch while the page is open re-renders every subscribed component.
describe('switching the language in place', () => {
  it('re-renders a translated component without remounting it', () => {
    resetLocaleForTests('en')
    render(<DomainTable rows={rows0()} onToggle={() => {}} expanded={null} />)
    expect(screen.getByText('3 concepts')).toBeInTheDocument()
    act(() => resetLocaleForTests('ru'))
    expect(screen.getByText('3 темы')).toBeInTheDocument()
  })
})

function rows0() {
  return [
    {
      domain: 'algebra',
      label: 'Algebra',
      total: 3,
      mastered: 0,
      dueForReview: 0,
      inProgress: 0,
      nextId: null,
      nextLabel: null,
      concepts: [
        { id: 'a', label: 'A', status: 'unseen' as const },
        { id: 'b', label: 'B', status: 'unseen' as const },
        { id: 'c', label: 'C', status: 'unseen' as const },
      ],
    },
  ]
}
