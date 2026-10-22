import { Fragment, type ReactNode } from 'react'

/**
 * A translated sentence with elements inside it.
 *
 * A message such as "{signIn} to join a league." needs a link in the middle of the sentence, and
 * the link sits in a different place in Russian than in English. The translator keeps the word
 * order and writes `{signIn}` where the element goes; this splits the string on those tokens and
 * puts the matching element in each. A token with no slot is left as written, so a typo in a
 * dictionary shows up on screen instead of vanishing.
 */
export default function RichText({ text, slots }: { text: string; slots: Record<string, ReactNode> }) {
  return (
    <>
      {text.split(/(\{\w+\})/g).map((part, i) => {
        const name = /^\{(\w+)\}$/.exec(part)?.[1]
        return name !== undefined && name in slots ? <Fragment key={i}>{slots[name]}</Fragment> : part
      })}
    </>
  )
}
