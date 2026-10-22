import type { Metadata } from 'next'
import DocsIndex from './DocsIndex'

export const metadata: Metadata = {
  title: 'Documentation — Mathua',
  description: 'Architecture, system design, administration, math input and the contributing guide for Mathua, the open-source adaptive math learning engine.',
}

export default function DocsIndexPage() {
  return <DocsIndex />
}
