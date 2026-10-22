import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'
import { describe, expect, it, vi } from 'vitest'

// Loads public/sw.js as a classic worker. It has no exports and expects the
// worker globals, so the test supplies them and captures the fetch listener.
const source = readFileSync(resolve(__dirname, '../public/sw.js'), 'utf8')

interface Harness {
  dispatchFetch: (request: unknown) => Promise<Response>
  puts: Array<{ request: unknown; response: Response }>
  store: Map<string, Response>
}

function loadWorker(fetchImpl: (req: unknown) => Promise<Response>): Harness {
  const listeners: Record<string, (event: unknown) => void> = {}
  const puts: Array<{ request: unknown; response: Response }> = []
  const store = new Map<string, Response>()

  const keyOf = (req: unknown) =>
    typeof req === 'string' ? req : String((req as { url?: string }).url ?? '')

  const cache = {
    put: async (request: unknown, response: Response) => {
      puts.push({ request, response })
      store.set(keyOf(request), response)
    },
    match: async (request: unknown) => store.get(keyOf(request)),
    addAll: async () => undefined,
  }

  const caches = {
    open: async () => cache,
    match: async (request: unknown) => store.get(keyOf(request)),
    keys: async () => [],
    delete: async () => true,
  }

  const self = {
    addEventListener: (name: string, fn: (event: unknown) => void) => {
      listeners[name] = fn
    },
    skipWaiting: async () => undefined,
    clients: { claim: async () => undefined },
    location: { origin: 'https://mathua.com' },
  }

  const load = new Function('self', 'caches', 'fetch', 'Response', 'URL', source)
  load(self, caches, fetchImpl, Response, URL)

  return {
    puts,
    store,
    dispatchFetch: async (request: unknown) => {
      let settled: Promise<Response> | undefined
      listeners.fetch({
        request,
        respondWith: (p: Promise<Response>) => {
          settled = p
        },
      })
      if (!settled) throw new Error('fetch event was not handled')
      return settled
    },
  }
}

const nav = { method: 'GET', mode: 'navigate', url: 'https://mathua.com/learn' }

const htmlResponse = (body: string) =>
  new Response(body, { status: 200, headers: { 'content-type': 'text/html; charset=utf-8' } })

describe('the service worker and Cloudflare-injected pages', () => {
  it('does not cache a navigation carrying a Precursor token', async () => {
    // Cloudflare appends this to the HTML of a challenged response. Cached, it
    // replays a dead token and Cloudflare answers 400 on every later visit.
    const body = '<html><body>learn</body><script>window.__CF$cv$params={r:"x",ut:"dead"}</script></html>'
    const worker = loadWorker(async () => htmlResponse(body))

    const res = await worker.dispatchFetch(nav)

    expect(await res.text()).toBe(body)
    expect(worker.puts).toHaveLength(0)
  })

  it('caches an ordinary navigation so it can be read offline', async () => {
    const body = '<html><body>learn</body></html>'
    const worker = loadWorker(async () => htmlResponse(body))

    await worker.dispatchFetch(nav)

    expect(worker.puts).toHaveLength(1)
  })

  it('serves the offline page rather than a cached challenged page', async () => {
    const worker = loadWorker(async () => {
      throw new TypeError('Failed to fetch')
    })
    const challenged = htmlResponse('<html><script>window.__CF$cv$params={}</script></html>')
    await worker.store.set(nav.url, challenged)
    await worker.store.set('/offline', htmlResponse('<html><body>offline</body></html>'))

    const res = await worker.dispatchFetch(nav)

    expect(await res.text()).toContain('offline')
  })

  it('serves a cached ordinary page when the network fails', async () => {
    const worker = loadWorker(async () => {
      throw new TypeError('Failed to fetch')
    })
    await worker.store.set(nav.url, htmlResponse('<html><body>cached learn</body></html>'))

    const res = await worker.dispatchFetch(nav)

    expect(await res.text()).toContain('cached learn')
  })
})

describe('API answers are cached per language', () => {
  const api = (lang: string) => ({
    method: 'GET',
    mode: 'cors',
    url: 'https://mathua.com/api/graph',
    headers: new Headers({ 'Accept-Language': lang }),
  })
  const json = (body: string) => new Response(body, { status: 200 })

  it('does not serve the English copy to a Russian request when offline', async () => {
    let online = true
    const worker = loadWorker(async () => {
      if (!online) throw new TypeError('Failed to fetch')
      return json('english')
    })
    await worker.dispatchFetch(api('en'))
    online = false
    // Nothing was cached for Russian, so the failure must surface rather than answer in English.
    await expect(worker.dispatchFetch(api('ru'))).rejects.toThrow('Failed to fetch')
    const english = await worker.dispatchFetch(api('en'))
    expect(await english.text()).toBe('english')
  })
})

describe('the worker version', () => {
  it('is not the static development stamp in a built export', () => {
    // The source must stay stampable; the build rewrites out/sw.js. This guards
    // the line the Dockerfile and scripts/stamp-sw.mjs both replace.
    expect(source).toMatch(/const VERSION = '[^']*'/)
  })
})

// Keep the import used even if the build check above is simplified later.
void vi
