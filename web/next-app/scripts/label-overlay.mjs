// Read the translated concept labels that the content side maintains.
//
// `data/i18n/ru/concepts.json` at the repository root is a flat object, concept id to label.
// It is owned by the translation work, not by this app, and it is allowed not to exist: a
// checkout with no Russian content must still build, and must build exactly what it built
// before. A file that exists but is malformed is a different case. It fails the build, because
// shipping English labels under a Russian toggle without anyone noticing is the failure this
// check is here to prevent.
import { existsSync, readFileSync } from 'node:fs'

/** Parse and validate overlay text. Throws with a message naming what is wrong. */
export function parseLabelOverlay(text, where = 'label overlay') {
  let data
  try {
    data = JSON.parse(text)
  } catch (err) {
    throw new Error(`${where}: not valid JSON (${err.message})`)
  }
  if (data === null || typeof data !== 'object' || Array.isArray(data)) {
    throw new Error(`${where}: expected an object of concept id to label`)
  }
  for (const [id, label] of Object.entries(data)) {
    if (typeof label !== 'string') throw new Error(`${where}: label for ${id} is not a string`)
  }
  return data
}

/** The overlay, or null when the file does not exist. */
export function readLabelOverlay(path) {
  if (!existsSync(path)) return null
  return parseLabelOverlay(readFileSync(path, 'utf8'), path)
}
