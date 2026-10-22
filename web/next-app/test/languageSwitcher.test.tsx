import { describe, it, expect, vi } from 'vitest'
import { act, fireEvent, render, screen } from '@testing-library/react'
import { useEffect } from 'react'
import LanguageSwitcher from '../components/LanguageSwitcher'
import { LocaleRoot, setLocale, useLocale, useT } from '../lib/i18n'

function Label() {
  const t = useT()
  const locale = useLocale()
  return <p data-testid="label">{`${locale}:${t('common.language.label')}`}</p>
}

describe('LanguageSwitcher', () => {
  it('shows both languages by their own names and marks the current one', () => {
    render(<LanguageSwitcher />)
    expect(screen.getByRole('button', { name: 'English' })).toHaveAttribute('aria-pressed', 'true')
    expect(screen.getByRole('button', { name: 'Русский' })).toHaveAttribute('aria-pressed', 'false')
    expect(screen.getByRole('button', { name: 'Русский' })).toHaveAttribute('lang', 'ru')
  })

  it('switches the interface immediately, without a reload', () => {
    render(
      <>
        <LanguageSwitcher />
        <Label />
      </>,
    )
    expect(screen.getByTestId('label')).toHaveTextContent('en:Language')
    fireEvent.click(screen.getByRole('button', { name: 'Русский' }))
    expect(screen.getByTestId('label')).toHaveTextContent('ru:Язык')
    expect(screen.getByRole('button', { name: 'Русский' })).toHaveAttribute('aria-pressed', 'true')
    expect(document.documentElement.lang).toBe('ru')
    expect(localStorage.getItem('mathua_locale')).toBe('ru')
  })
})

describe('LocaleRoot', () => {
  it('remounts the page on a switch so localized data is fetched again', () => {
    let mounts = 0
    function Page() {
      useEffect(() => {
        mounts += 1
      }, [])
      return null
    }
    render(
      <LocaleRoot>
        <Page />
      </LocaleRoot>,
    )
    expect(mounts).toBe(1)
    act(() => setLocale('ru'))
    expect(mounts).toBe(2)
    act(() => setLocale('ru'))
    expect(mounts).toBe(2)
  })
})

// The page is prerendered in English and hydrated in the learner's language. A helper that
// formats text outside a hook must agree with the hook for that first pass, or React finds the
// DOM different from what it rendered and throws the server markup away.
describe('hydration in Russian', () => {
  it('matches the prerendered English markup on the first client render, then shows Russian', async () => {
    const { renderToString } = await import('react-dom/server')
    const { hydrateRoot } = await import('react-dom/client')
    const { domainName } = await import('../lib/domainName')
    const { resetLocaleForTests } = await import('../lib/i18n/locale')

    function Probe() {
      const t = useT()
      return (
        <p>
          {t('common.language.label')} / {domainName('arithmetic')}
        </p>
      )
    }
    const tree = (
      <LocaleRoot>
        <Probe />
      </LocaleRoot>
    )

    resetLocaleForTests('en')
    const spy = vi.spyOn(console, 'error').mockImplementation(() => {})
    const html = renderToString(tree)
    expect(html.replace(/<!-- -->/g, '')).toContain('Language / Arithmetic')
    spy.mockClear()

    resetLocaleForTests('ru')
    const container = document.createElement('div')
    container.innerHTML = html
    document.body.appendChild(container)
    await act(async () => {
      hydrateRoot(container, tree)
    })

    expect(spy.mock.calls.map(c => String(c[0]))).toEqual([])
    expect(container.textContent).toBe('Язык / Арифметика')
    spy.mockRestore()
    container.remove()
  })
})
