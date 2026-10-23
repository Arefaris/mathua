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

// Ukrainian rendering of the landing, How it works, the domain table, the glossary chrome and the
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

describe('landing in Ukrainian', () => {
  it('shows the hero, the domains and the level names in Ukrainian', async () => {
    resetLocaleForTests('uk')
    localStorage.clear()
    render(<HomePage />)
    expect(await screen.findByRole('button', { name: 'Почати' })).toBeInTheDocument()
    expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent(
      'Адаптивна платформа для вивчення математики',
    )
    // The coverage table uses the shared domain names (glossary section 5).
    expect(screen.getAllByText('Дроби').length).toBeGreaterThan(0)
    expect(screen.getAllByText('Вступ до алгебри').length).toBeGreaterThan(0)
    expect(screen.getByText('Підмайстер')).toBeInTheDocument()
    // The link text does not depend on the number: "Усі теми за предметами: 1250 →".
    expect(screen.getByText(/^Усі теми за предметами: \d+ →$/)).toBeInTheDocument()
    // State chips are uppercase Ukrainian state names.
    expect(screen.getByText('ЗАКРІПЛЮЄТЬСЯ')).toBeInTheDocument()
    noKeyIds()
  })
})

describe('How it works in Ukrainian', () => {
  it('keeps the numbers and picks the plural forms', () => {
    resetLocaleForTests('uk')
    render(
      <>
        <ConceptGraphSection />
        <StudentModelSection />
      </>,
    )
    expect(screen.getByText(/^Наразі в графі є \d+ тем і \d+ розділ/)).toBeInTheDocument()
    // 5 and 7 take the genitive plural, 8 / 18 / 20 seconds are "секунд".
    expect(screen.getAllByText('5 правильних відповідей поспіль').length).toBeGreaterThan(0)
    expect(screen.getAllByText('8 секунд').length).toBeGreaterThan(0)
    expect(screen.getAllByText('Додавання одноцифрових чисел').length).toBeGreaterThan(0)
    noKeyIds()
  })

  it('declines "entries" like a Ukrainian noun', () => {
    const t = getTranslator('uk')
    expect(t('explore.glossary.search.entries', { count: 1 })).toBe('1 запис')
    expect(t('explore.glossary.search.entries', { count: 2 })).toBe('2 записи')
    expect(t('explore.glossary.search.entries', { count: 5 })).toBe('5 записів')
    expect(t('explore.glossary.search.entries', { count: 11 })).toBe('11 записів')
    expect(t('explore.glossary.search.entries', { count: 21 })).toBe('21 запис')
    expect(t('explore.glossary.search.entries', { count: 22 })).toBe('22 записи')
  })
})

describe('domain table in Ukrainian', () => {
  const rows = [
    {
      domain: 'fractions',
      label: 'Дроби',
      total: 12,
      mastered: 3,
      dueForReview: 2,
      inProgress: 1,
      nextId: 'frac.add.same',
      nextLabel: 'Додавання дробів',
      concepts: [{ id: 'frac.add.same', label: 'Додавання дробів', status: 'dueForReview' as const }],
    },
  ]

  it('translates the counts, the call to action and the status groups', () => {
    resetLocaleForTests('uk')
    render(<DomainTable rows={rows} onToggle={() => {}} expanded="fractions" />)
    expect(screen.getAllByText(/повторити: 2/).length).toBeGreaterThan(0)
    expect(screen.getByText('Продовжити: Додавання дробів')).toBeInTheDocument()
    expect(screen.getByText('Сховати теми')).toBeInTheDocument()
    expect(screen.getByText('Час повторити (1)')).toBeInTheDocument()
    noKeyIds()
  })

  it('counts the concepts with the shared plural', () => {
    resetLocaleForTests('uk')
    render(<DomainTable rows={rows} onToggle={() => {}} expanded={null} />)
    expect(screen.getByText('1 тема')).toBeInTheDocument()
  })
})

describe('glossary chrome in Ukrainian', () => {
  it('translates the controls, the categories and the result count', () => {
    resetLocaleForTests('uk')
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
    expect(screen.getByPlaceholderText(/^Шукайте символ/)).toBeInTheDocument()
    expect(screen.getByText('21 запис')).toBeInTheDocument()
    fireEvent.click(screen.getByRole('button', { name: 'Теги' }))
    expect(screen.getByRole('button', { name: 'Грецькі літери' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Теорія чисел і вища математика' })).toBeInTheDocument()
    noKeyIds()
  })

  it('shows a page range with the genitive plural', () => {
    resetLocaleForTests('uk')
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
    expect(screen.getByText('Показано 21–40 із 52 записів')).toBeInTheDocument()
  })

  it('frames an entry in Ukrainian and leaves its content alone', () => {
    resetLocaleForTests('uk')
    const entry = GLOSSARY_ENTRIES.find(e => e.commonMistakes && e.commonMistakes.length > 0)!
    render(<GlossaryEntryCard entry={entry} onSelectRelated={() => {}} />)
    expect(screen.getByText('Типова помилка:')).toBeInTheDocument()
    expect(screen.getByText(entry.definition)).toBeInTheDocument()
  })
})

describe('concept graph in Ukrainian', () => {
  const concepts = [
    { id: 'count', label: 'Лічба', domain: 'arithmetic', prerequisites: [] },
    { id: 'add', label: 'Додавання', domain: 'fractions', prerequisites: ['count'] },
  ]

  it('translates the controls and the node description', async () => {
    resetLocaleForTests('uk')
    const { container } = render(
      <ConceptGraphFlow
        concepts={concepts}
        conceptStatuses={{ count: 'mastered', add: 'learning' }}
        conceptProgress={{ add: 0.4 }}
      />,
    )
    expect(screen.getByPlaceholderText('Пошук тем…')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Відкрити список тем' })).toBeInTheDocument()
    const node = container.querySelector('.concept-node[aria-label^="Додавання"]')
    expect(node?.getAttribute('aria-label')).toBe('Додавання, розділ «Дроби», Вивчається, 40% до засвоєння')
    expect(node?.querySelector('.node-progress-tip')?.textContent).toBe('40% до засвоєння')
    fireEvent.click(screen.getByRole('button', { name: 'Відкрити список тем' }))
    await waitFor(() => expect(screen.getByRole('dialog', { name: 'Список тем' })).toBeInTheDocument())
    noKeyIds()
  })

  it('shows the empty state in Ukrainian', () => {
    resetLocaleForTests('uk')
    render(<ConceptGraphFlow concepts={[]} />)
    expect(screen.getByText('Немає тем для показу')).toBeInTheDocument()
  })
})

// A switch while the page is open re-renders every subscribed component.
describe('switching the language in place', () => {
  it('re-renders a translated component without remounting it', () => {
    resetLocaleForTests('en')
    render(<DomainTable rows={rows0()} onToggle={() => {}} expanded={null} />)
    expect(screen.getByText('3 concepts')).toBeInTheDocument()
    act(() => resetLocaleForTests('uk'))
    expect(screen.getByText('3 теми')).toBeInTheDocument()
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
