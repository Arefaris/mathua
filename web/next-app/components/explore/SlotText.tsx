import { Fragment, type ReactNode } from 'react'

/**
 * A translated sentence with live elements inside it.
 *
 * A sentence such as "Sign in or {link} to track your progress" cannot be split into three
 * strings: Russian puts the link where English does not, and a translator has to see the whole
 * sentence. The message carries `{name}` markers, and `slots` says what each one becomes. A
 * marker with no slot is printed as written, so a missing slot shows up on screen instead of
 * silently dropping words.
 */
export default function SlotText({ text, slots }: { text: string; slots: Record<string, ReactNode> }) {
  const parts = text.split(/\{(\w+)\}/)
  return (
    <>
      {parts.map((part, i) => {
        // `split` with a capture group alternates text, name, text, name...
        if (i % 2 === 0) return part === '' ? null : <Fragment key={i}>{part}</Fragment>
        return <Fragment key={i}>{part in slots ? slots[part] : `{${part}}`}</Fragment>
      })}
    </>
  )
}
