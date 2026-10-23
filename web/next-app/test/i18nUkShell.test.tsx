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
import { getTranslator, setLocale } from '../lib/i18n'

// The shared chrome and the pages around it, read in Ukrainian. Ukrainian is registered but not
// yet offered in the language switcher, so every test pins the locale itself.

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

describe('the shared chrome in Ukrainian', () => {
  it('names the header controls and the navigation menu in Ukrainian', () => {
    resetLocaleForTests('uk')
    render(<Header />)
    fireEvent.click(screen.getByRole('button', { name: 'Відкрити меню навігації' }))
    const menu = screen.getByRole('menu', { name: 'Навігація сайтом' })
    expect(within(menu).getByRole('menuitem', { name: 'Рейтинг' })).toBeInTheDocument()
    expect(within(menu).getByRole('menuitem', { name: 'Глосарій' })).toBeInTheDocument()
    expect(within(menu).getByRole('menuitem', { name: 'Зв’язатися' })).toBeInTheDocument()
    expect(within(menu).getByRole('menuitem', { name: 'Вхід' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Змінити тему' })).toBeInTheDocument()
    expect(screen.getByText('бета')).toHaveAttribute('title', expect.stringContaining('бета-версії'))
  })

  it('shows a link list the page passes in exactly as given', () => {
    resetLocaleForTests('uk')
    render(<Header links={[{ label: 'Свій пункт', href: '/x' }]} />)
    expect(screen.getAllByText('Свій пункт').length).toBeGreaterThan(0)
  })

  it('labels the tab bar and its five tabs', () => {
    resetLocaleForTests('uk')
    render(<BottomTabs />)
    const bar = screen.getByRole('navigation', { name: 'Основне меню' })
    for (const name of ['Головна', 'Вчитися', 'Розділи', 'Профіль', 'Граф']) {
      expect(within(bar).getByRole('link', { name })).toBeInTheDocument()
    }
  })

  it('translates the footer', () => {
    resetLocaleForTests('uk')
    render(<Footer />)
    expect(screen.getByRole('link', { name: 'Розділи' })).toBeInTheDocument()
    expect(screen.getByText('Адаптивна система навчання математики з відкритим кодом.')).toBeInTheDocument()
    expect(screen.getByText('Ліцензія MIT')).toBeInTheDocument()
    expect(screen.getByRole('link', { name: 'Нотатка автора' })).toBeInTheDocument()
    expect(document.documentElement.lang).toBe('uk')
  })

  it('translates the profile sidebar, including a level the server sends in English', () => {
    resetLocaleForTests('uk')
    render(
      <SidebarProvider>
        <AppSidebar name="Ada" studentId="s1" username="ada" level="Scholar" streak={5} />
      </SidebarProvider>,
    )
    expect(screen.getByText('Навігація')).toBeInTheDocument()
    expect(screen.getByText('Матеріали')).toBeInTheDocument()
    expect(screen.getByText('Рейтинг')).toBeInTheDocument()
    expect(screen.getByText(/Науковець/)).toBeInTheDocument()
    expect(screen.getByText(/5 дн\./)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Змінити тему' })).toBeInTheDocument()
  })

  it('translates the hide/show control and the confirm dialog', () => {
    resetLocaleForTests('uk')
    render(
      <>
        <ChromeToggle />
        <ConfirmDialog open title="Заголовок" confirmLabel="Так" busy onConfirm={() => {}} onCancel={() => {}}>
          <span>x</span>
        </ConfirmDialog>
      </>,
    )
    expect(screen.getByRole('button', { name: 'Сховати навігацію' })).toBeInTheDocument()
    expect(screen.getByText('Виконується…')).toBeInTheDocument()
    expect(screen.getAllByText('Скасувати').length).toBeGreaterThan(0)
  })

  it('labels the table of contents', () => {
    resetLocaleForTests('uk')
    render(<TableOfContents headings={[{ id: 'a', text: 'Розділ A' }]} />)
    expect(screen.getAllByLabelText('Зміст').length).toBeGreaterThan(0)
    expect(screen.getAllByText('На цій сторінці').length).toBeGreaterThan(0)
  })
})

describe('pages and states in Ukrainian', () => {
  it('translates the 404 page', () => {
    resetLocaleForTests('uk')
    render(<NotFound />)
    expect(screen.getByRole('heading', { name: 'Ця сторінка не визначена.' })).toBeInTheDocument()
    expect(screen.getByRole('link', { name: '← На головну' })).toHaveAttribute('href', '/')
    expect(screen.getByRole('button', { name: 'Назад' })).toBeInTheDocument()
  })

  it('translates the contact page and keeps its links in the sentence', () => {
    resetLocaleForTests('uk')
    render(<ContactPage />)
    expect(screen.getByRole('heading', { name: 'У питанні помилка' })).toBeInTheDocument()
    expect(screen.getByLabelText('Ваше повідомлення')).toBeInTheDocument()
    const x = screen.getByRole('link', { name: '@wisdomsglow' })
    expect(x).toHaveAttribute('href', 'https://x.com/wisdomsglow')
    expect(screen.getByRole('link', { name: 'Сервер Discord Mathua' })).toHaveAttribute(
      'href',
      'https://discord.gg/qXPeQ2jys',
    )
    // The sentence is whole: the link sits inside the paragraph, not in a fragment of its own.
    expect(x.closest('p')?.textContent).toMatch(/^Напишіть мені в X: @wisdomsglow\./)
  })

  it('uses the right Ukrainian form after "щонайменше"', () => {
    resetLocaleForTests('uk')
    render(<ContactBox />)
    fireEvent.change(screen.getByLabelText('Ваше повідомлення'), { target: { value: 'коротко' } })
    expect(screen.getByText('щонайменше 10 символів')).toBeInTheDocument()
  })

  it.each([
    [1, 'щонайменше 1 символ'],
    [2, 'щонайменше 2 символи'],
    [5, 'щонайменше 5 символів'],
    [11, 'щонайменше 11 символів'],
    [21, 'щонайменше 21 символ'],
    [22, 'щонайменше 22 символи'],
  ])('counts %i characters with the right form', (count, expected) => {
    expect(getTranslator('uk')('shell.contact.box.minLength', { count })).toBe(expected)
  })

  it('translates the docs index and says the articles are English', () => {
    resetLocaleForTests('uk')
    render(<DocsIndexPage />)
    expect(screen.getByRole('heading', { name: 'Документація' })).toBeInTheDocument()
    expect(screen.getByText('Архітектура')).toBeInTheDocument()
    expect(screen.getByText('Самі статті — англійською.')).toBeInTheDocument()
    // The links stay inside the sentence.
    expect(screen.getByRole('link', { name: 'Як це працює' })).toHaveAttribute('href', '/how-it-works')
  })

  it('never prints a raw key id', () => {
    resetLocaleForTests('uk')
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
