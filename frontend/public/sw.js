// Onatar service worker — PRD RNF-03 (offline parcial).
// Caches the rules content (GET /api/v1/content) stale-while-revalidate and
// falls back to the app shell for navigation requests.

const CACHE = 'onatar-v1'
const CONTENT_URL = '/api/v1/content'

self.addEventListener('install', (event) => {
  event.waitUntil(
    caches
      .open(CACHE)
      .then((cache) => cache.addAll(['/', '/index.html']))
      .then(() => self.skipWaiting())
  )
})

self.addEventListener('activate', (event) => {
  event.waitUntil(
    caches
      .keys()
      .then((keys) => Promise.all(keys.filter((k) => k !== CACHE).map((k) => caches.delete(k))))
      .then(() => self.clients.claim())
  )
})

self.addEventListener('fetch', (event) => {
  const { request } = event
  if (request.method !== 'GET') return

  const url = new URL(request.url)
  if (url.origin !== self.location.origin) return

  if (url.pathname === CONTENT_URL) {
    // stale-while-revalidate: serve cache immediately, refresh in background.
    event.respondWith(
      caches.open(CACHE).then(async (cache) => {
        const cached = await cache.match(request)
        const network = fetch(request)
          .then((resp) => {
            if (resp.ok) cache.put(request, resp.clone())
            return resp
          })
          .catch(() => cached)
        return cached || network
      })
    )
    return
  }

  if (request.mode === 'navigate') {
    event.respondWith(
      fetch(request)
        .then((resp) => {
          cachePut(CACHE, request, resp.clone())
          return resp
        })
        .catch(() => caches.match('/index.html'))
    )
  }
})

function cachePut(cacheName, request, response) {
  caches.open(cacheName).then((cache) => cache.put(request, response))
}
