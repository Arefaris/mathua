import { describe, it, expect } from 'vitest'
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
