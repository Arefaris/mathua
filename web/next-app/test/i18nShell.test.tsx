import { describe, it, expect, beforeAll, vi } from 'vitest'
import { act, fireEvent, render, screen, within } from '@testing-library/react'
import Header from '../components/Header'
import BottomTabs from '../components/BottomTabs'
import Footer from '../components/Footer'
import ConfirmDialog from '../components/ConfirmDialog'
import ChromeToggle from '../components/ChromeToggle'
import TableOfContents from '../components/TableOfContents'
import ContactBox from '../components/ContactBox'
import { AppSidebar } from '../components/app-sidebar'
import { SidebarProvider } from '../components/ui/sidebar'
import NotFound from '../app/not-found'
import ContactPage from '../app/contact/page'
import DocsIndexPage from '../app/docs/page'
import { resetLocaleForTests } from '../lib/i18n/locale'
import { setLocale } from '../lib/i18n'

beforeAll(() => {
  Object.defineProperty(window, 'matchMedia', {
    writable: true,
    value: () => ({
      matches: false,
      media: '',
      onchange: null,
      addListener: () => {},
      removeListener: () => {},
      addEventListener: () => {},
      removeEventListener: () => {},
      dispatchEvent: () => false,
    }),
  })
})

vi.mock('next/navigation', () => ({
  useRouter: () => ({ push: vi.fn() }),
  usePathname: () => '/learn',
}))
vi.mock('../hooks/useTheme', () => ({
  useTheme: () => ({ theme: 'dark', mounted: true, toggleTheme: vi.fn() }),
}))
vi.mock('../hooks/useAuthState', () => ({
  useAuthState: () => ({ loggedIn: false, user: null }),
}))
vi.mock('../hooks/useChrome', () => ({
  useChrome: () => ({ visible: true, toggle: () => {}, hasChrome: true, hasTabs: true }),
}))
vi.mock('../lib/api', async importOriginal => ({
  ...(await importOriginal<typeof import('../lib/api')>()),
  getSettings: vi.fn(async () => ({})),
  submitReport: vi.fn(async () => ({})),
}))

// The chrome is on every page, so a string left in English here shows on every screen.
describe('the shared chrome in Russian', () => {
  it('names the header controls and the navigation menu in Russian', () => {
    resetLocaleForTests('ru')
    render(<Header />)
    fireEvent.click(screen.getByRole('button', { name: 'Открыть меню навигации' }))
    const menu = screen.getByRole('menu', { name: 'Навигация по сайту' })
    // Learn and Graph are tabs on a phone, so the compass lists the rest.
    expect(within(menu).getByRole('menuitem', { name: 'Рейтинг' })).toBeInTheDocument()
    expect(within(menu).getByRole('menuitem', { name: 'Глоссарий' })).toBeInTheDocument()
    expect(within(menu).getByRole('menuitem', { name: 'Связаться' })).toBeInTheDocument()
    expect(within(menu).getByRole('menuitem', { name: 'Вход' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Сменить тему' })).toBeInTheDocument()
    expect(screen.getByText('бета')).toHaveAttribute('title', expect.stringContaining('в стадии беты'))
  })

  it('shows a link list the page passes in exactly as given', () => {
    resetLocaleForTests('ru')
    render(<Header links={[{ label: 'Свой пункт', href: '/x' }]} />)
    expect(screen.getAllByText('Свой пункт').length).toBeGreaterThan(0)
  })

  it('labels the tab bar and its five tabs', () => {
    resetLocaleForTests('ru')
    render(<BottomTabs />)
    const bar = screen.getByRole('navigation', { name: 'Основное меню' })
    for (const name of ['Главная', 'Учиться', 'Разделы', 'Профиль', 'Граф']) {
      expect(within(bar).getByRole('link', { name })).toBeInTheDocument()
    }
  })

  it('translates the footer and lets a visitor change language there', () => {
    render(<Footer />)
    expect(screen.getByRole('link', { name: 'Domains' })).toBeInTheDocument()
    // A guest cannot open /settings, so the footer carries the switch as well.
    fireEvent.click(screen.getByRole('button', { name: 'Русский' }))
    expect(screen.getByRole('link', { name: 'Разделы' })).toBeInTheDocument()
    expect(screen.getByText('Адаптивная система обучения математике с открытым кодом.')).toBeInTheDocument()
    expect(screen.getByText('Лицензия MIT')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Русский' })).toHaveAttribute('aria-pressed', 'true')
    expect(document.documentElement.lang).toBe('ru')
  })

  it('translates the profile sidebar, including a level the server sends in English', () => {
    resetLocaleForTests('ru')
    render(
      <SidebarProvider>
        <AppSidebar name="Ada" studentId="s1" username="ada" level="Scholar" streak={5} />
      </SidebarProvider>,
    )
    expect(screen.getByText('Навигация')).toBeInTheDocument()
    expect(screen.getByText('Материалы')).toBeInTheDocument()
    expect(screen.getByText('Рейтинг')).toBeInTheDocument()
    expect(screen.getByText(/Учёный/)).toBeInTheDocument()
    expect(screen.getByText(/5 дн\./)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Сменить тему' })).toBeInTheDocument()
  })

  it('translates the hide/show control and the confirm dialog', () => {
    resetLocaleForTests('ru')
    render(
      <>
        <ChromeToggle />
        <ConfirmDialog open title="Заголовок" confirmLabel="Да" busy onConfirm={() => {}} onCancel={() => {}}>
          <span>x</span>
        </ConfirmDialog>
      </>,
    )
    expect(screen.getByRole('button', { name: 'Скрыть навигацию' })).toBeInTheDocument()
    expect(screen.getByText('Выполняется…')).toBeInTheDocument()
    expect(screen.getAllByText('Отмена').length).toBeGreaterThan(0)
  })

  it('labels the table of contents', () => {
    resetLocaleForTests('ru')
    render(<TableOfContents headings={[{ id: 'a', text: 'Раздел A' }]} />)
    expect(screen.getAllByLabelText('Содержание').length).toBeGreaterThan(0)
    expect(screen.getAllByText('На этой странице').length).toBeGreaterThan(0)
  })
})

describe('pages and states in Russian', () => {
  it('translates the 404 page', () => {
    resetLocaleForTests('ru')
    render(<NotFound />)
    expect(screen.getByRole('heading', { name: 'Эта страница не определена.' })).toBeInTheDocument()
    expect(screen.getByRole('link', { name: '← На главную' })).toHaveAttribute('href', '/')
    expect(screen.getByRole('button', { name: 'Назад' })).toBeInTheDocument()
  })

  it('translates the contact page and keeps its links in the sentence', () => {
    resetLocaleForTests('ru')
    render(<ContactPage />)
    expect(screen.getByRole('heading', { name: 'Ошибка в вопросе' })).toBeInTheDocument()
    expect(screen.getByLabelText('Ваше сообщение')).toBeInTheDocument()
    const x = screen.getByRole('link', { name: '@wisdomsglow' })
    expect(x).toHaveAttribute('href', 'https://x.com/wisdomsglow')
    expect(screen.getByRole('link', { name: 'сервер Mathua в Discord' })).toHaveAttribute(
      'href',
      'https://discord.gg/qXPeQ2jys',
    )
    // The sentence is whole: the link sits inside the paragraph, not in a fragment of its own.
    expect(x.closest('p')?.textContent).toMatch(/^Напишите мне в X: @wisdomsglow\./)
  })

  it('uses the right Russian form after "не менее"', () => {
    resetLocaleForTests('ru')
    render(<ContactBox />)
    fireEvent.change(screen.getByLabelText('Ваше сообщение'), { target: { value: 'коротко' } })
    expect(screen.getByText('не менее 10 символов')).toBeInTheDocument()
  })

  it('translates the docs index and says the articles are English', () => {
    resetLocaleForTests('ru')
    render(<DocsIndexPage />)
    expect(screen.getByRole('heading', { name: 'Документация' })).toBeInTheDocument()
    expect(screen.getByText('Архитектура')).toBeInTheDocument()
    expect(screen.getByText('Сами статьи пока доступны только на английском.')).toBeInTheDocument()
  })

  it('does not show the English-only notice to an English reader', () => {
    render(<DocsIndexPage />)
    expect(screen.queryByText('The articles themselves are in English.')).toBeNull()
  })

  it('never prints a raw key id', () => {
    resetLocaleForTests('ru')
    const { container } = render(
      <>
        <Header />
        <BottomTabs />
        <Footer />
        <ContactPage />
      </>,
    )
    expect(container.textContent).not.toMatch(/\bshell\.[a-z]+\./)
    act(() => setLocale('en'))
    expect(container.textContent).not.toMatch(/\bshell\.[a-z]+\./)
  })
})
