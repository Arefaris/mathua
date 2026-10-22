// Compile-time key checking for `t()`.
//
// The key type is derived from the English dictionary, so adding a message there makes its key
// valid with no second list to maintain. A group is a nested object; a plural message is an
// object whose keys are all CLDR categories and which has `other`.
import type { EnglishMessages } from './messages'

type PluralKey = 'zero' | 'one' | 'two' | 'few' | 'many' | 'other'

type IsPlural<T> = T extends { other: string } ? ([keyof T] extends [PluralKey] ? true : false) : false

type Paths<T, Prefix extends string = ''> = {
  [K in keyof T & string]: T[K] extends string
    ? `${Prefix}${K}`
    : IsPlural<T[K]> extends true
      ? `${Prefix}${K}`
      : Paths<T[K], `${Prefix}${K}.`>
}[keyof T & string]

export type MessageKey = Paths<EnglishMessages>
