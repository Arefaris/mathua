import type { Metadata } from 'next'
import ContactContent from './ContactContent'

export const metadata: Metadata = {
  title: 'Contact, Mathua',
  description:
    'How to report a problem with a question, a bug in the site, or your account, and who reads it.',
}

export default function ContactPage() {
  return <ContactContent />
}
