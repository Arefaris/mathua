import type { Metadata } from 'next'
import NoteContent from './NoteContent'

export const metadata: Metadata = {
  title: "Creator's Note — Mathua",
  description: 'Why Mathua exists, from the person building it.',
}

// The page is a server component so it can export metadata; the text on it needs the learner's
// language, which only the browser knows, so the body lives in a client component.
export default function NotePage() {
  return <NoteContent />
}
