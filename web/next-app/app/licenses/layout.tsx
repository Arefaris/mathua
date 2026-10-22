'use client'

import Header from '../../components/Header'
import Footer from '../../components/Footer'
import BottomTabs from '../../components/BottomTabs'

export default function LicensesLayout({ children }: { children: React.ReactNode }) {
  return (
    <>
      <Header links={[
        { label: 'Domains', href: '/domains', labelKey: 'shell.nav.domains' },
        { label: 'Docs', href: '/docs', labelKey: 'shell.nav.docs' },
        { label: 'Architecture', href: '/docs/architecture', labelKey: 'shell.nav.architecture' },
        { label: 'Contributing', href: '/docs/contributing', labelKey: 'shell.nav.contributing' },
        { label: 'Note', href: '/note', labelKey: 'shell.nav.note' },
      ]} />
      <main>{children}</main>
      <Footer />
      <BottomTabs />
    </>
  )
}
