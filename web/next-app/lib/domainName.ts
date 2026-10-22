import { domainLabel } from './graphDomains'
import { getLocale } from './i18n/locale'
import { getTranslator } from './i18n/react'

/**
 * A domain's name in the learner's language.
 *
 * `domainLabel` in graphDomains.ts stays the English source: build-graph.mjs imports that file
 * under plain Node, which cannot resolve the app's extensionless imports, and the tests read the
 * English names from it. This is the one the interface calls. An id the dictionary does not know
 * falls back to the English label, then to the id with underscores removed.
 *
 * The caller must be a component that reads the locale (`useT`), or it will not re-render when
 * the learner switches.
 */
export function domainName(domain: string): string {
  return getTranslator(getLocale()).dynamic(`common.domain.${domain}`, undefined, domainLabel(domain))
}
