/* Mathua service worker — hand-rolled for the Next.js static export.
 *
 * Strategy:
 *   - navigations (HTML): network-first, cache fallback, then /offline.html
 *   - hashed static assets + images + fonts: cache-first (immutable)
 *   - GET /api/*: network-first with cache fallback (read-only endpoints;
 *     non-GET requests are never intercepted so answer submissions fail
 *     loudly instead of silently succeeding from a stale cache)
 *
 * VERSION is stamped at image build time (Dockerfile rewrites the line below
 * with a UTC timestamp) so every deploy ships byte-different sw.js: browsers
 * install the new worker, and activate purges all caches from older builds.
 * A static version would let clients run stale hashed JS indefinitely — old
 * chunks 404 on fresh images while this worker keeps serving its cached
 * copies (cache-first), pinning the client on pre-fix rendering forever.
 */
const VERSION = 'v2-local'
const SHELL_CACHE = `mathua-shell-${VERSION}`
const ASSET_CACHE = `mathua-assets-${VERSION}`
const DATA_CACHE = `mathua-data-${VERSION}`

// '/offline' resolves to offline.html on both the Go server (clean-URL
// rewrite) and static hosts; '.html' would 301 on some hosts, and a
// redirected cached response is invalid for navigations.
const PRECACHE_URLS = ['/', '/offline']

const ASSET_PATTERN = /\/_next\/static\/|\/diagrams\/|\/icons\/|\.(?:png|jpe?g|gif|svg|ico|webp|woff2?|ttf|css|js|webmanifest)$/i

self.addEventListener('install', event => {
  event.waitUntil(
    caches
      .open(SHELL_CACHE)
      .then(cache => cache.addAll(PRECACHE_URLS))
      .then(() => self.skipWaiting())
  )
})

self.addEventListener('activate', event => {
  event.waitUntil(
    caches
      .keys()
      .then(keys =>
        Promise.all(
          keys
            .filter(key => ![SHELL_CACHE, ASSET_CACHE, DATA_CACHE].includes(key))
            .map(key => caches.delete(key))
        )
      )
      .then(() => self.clients.claim())
  )
})

// An API answer depends on the language it was asked in: concept labels and course names are
// served translated. The browser cache would key on the URL alone, so a Russian learner who
// went offline could be handed the English copy. The language the page asked for is part of
// the key; the key is never sent anywhere, so it can be any URL that stays unique.
function dataKey(request) {
  const lang = request.headers && request.headers.get ? request.headers.get('Accept-Language') : null
  if (!lang) return request
  const keyed = new URL(request.url)
  keyed.searchParams.set('__lang', lang)
  return keyed.href
}

async function networkFirst(request, cacheName, fallbackUrls) {
  const cache = await caches.open(cacheName)
  const key = dataKey(request)
  try {
    const response = await fetch(request)
    if (response.ok) cache.put(key, response.clone())
    return response
  } catch (err) {
    const cached = await cache.match(key)
    if (cached) return cached
    for (const url of fallbackUrls ?? []) {
      const fallback = await caches.match(url)
      if (fallback) return fallback
    }
    throw err
  }
}

// Cloudflare injects a per-request Precursor token into the HTML:
// `window.__CF$cv$params = {r, t, u, ut, i}` followed by a script that appends
// /cdn-cgi/challenge-platform/scripts/precursor/main.js. The token is only valid
// for the response it arrived with. Storing that page and replaying it later
// makes the injected script post a dead token, and Cloudflare answers 400 to a
// challenge the page can no longer complete — which is the `h/undefined/...`
// request in the console, repeated on every visit because the cache never
// refreshes. Store navigation HTML only when it carries no token, and never
// serve a cached page that does.
const CHALLENGE_MARKER = '__CF$cv$params'

async function hasChallengeToken(response) {
  try {
    const type = response.headers.get('content-type') || ''
    if (!type.includes('text/html')) return false
    const text = await response.clone().text()
    return text.includes(CHALLENGE_MARKER)
  } catch {
    // If the body cannot be read, treat the page as uncacheable: a missed
    // offline page is a smaller harm than replaying a stale challenge.
    return true
  }
}

async function navigationNetworkFirst(request) {
  const cache = await caches.open(SHELL_CACHE)
  try {
    const response = await fetch(request)
    if (response.ok && !(await hasChallengeToken(response))) {
      cache.put(request, response.clone())
    }
    return response
  } catch (err) {
    const cached = await cache.match(request)
    if (cached && !(await hasChallengeToken(cached))) return cached
    const offline = await caches.match('/offline')
    if (offline) return offline
    throw err
  }
}

async function cacheFirst(request, cacheName) {
  const cache = await caches.open(cacheName)
  const cached = await cache.match(request)
  if (cached) return cached
  const response = await fetch(request)
  if (response.ok) cache.put(request, response.clone())
  return response
}

self.addEventListener('fetch', event => {
  const { request } = event
  if (request.method !== 'GET') return

  const url = new URL(request.url)
  if (url.origin !== self.location.origin) return

  if (request.mode === 'navigate') {
    event.respondWith(navigationNetworkFirst(request))
    return
  }

  if (ASSET_PATTERN.test(url.pathname)) {
    event.respondWith(cacheFirst(request, ASSET_CACHE))
    return
  }

  if (url.pathname.startsWith('/api/')) {
    event.respondWith(networkFirst(request, DATA_CACHE))
  }
})
