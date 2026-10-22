import type { Locale } from '../locale'
import { en } from './en'
import { ru } from './ru'

export type EnglishMessages = typeof en

export const MESSAGES: Record<Locale, unknown> = { en, ru }
