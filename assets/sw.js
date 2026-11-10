// MetaBible service worker: the shell and libraries work offline; pages and data refresh when online.
// Bump VERSION when changing this file or when static assets change shape.
const VERSION = 'mb-v1';
const STATIC = VERSION + '-static', DATA = VERSION + '-data', PAGES = VERSION + '-pages';
const PRECACHE = ['/', '/assets/css/metabible.css', '/assets/js/site/site.js', '/assets/vendor/fonts/fonts.css',
  '/assets/vendor/maplibre/maplibre-gl.js', '/assets/vendor/maplibre/maplibre-gl.css', '/assets/manifest.webmanifest'];

self.addEventListener('install', (e) => { e.waitUntil(caches.open(STATIC).then((c) => c.addAll(PRECACHE)).then(() => self.skipWaiting())); });
self.addEventListener('activate', (e) => {
  e.waitUntil(caches.keys().then((ks) => Promise.all(ks.filter((k) => !k.startsWith(VERSION)).map((k) => caches.delete(k)))).then(() => self.clients.claim()));
});

async function cacheFirst(req, name) {
  const c = await caches.open(name), hit = await c.match(req);
  if (hit) return hit;
  const res = await fetch(req); if (res.ok) c.put(req, res.clone()); return res;
}
async function staleWhileRevalidate(req, name) {
  const c = await caches.open(name), hit = await c.match(req);
  const net = fetch(req).then((res) => { if (res.ok) c.put(req, res.clone()); return res; }).catch(() => hit);
  return hit || net;
}
async function networkFirst(req, name) {
  const c = await caches.open(name);
  try { const res = await fetch(req); if (res.ok) c.put(req, res.clone()); return res; }
  catch (e) { return (await c.match(req)) || (await caches.match('/')) || Response.error(); }
}

self.addEventListener('fetch', (e) => {
  const req = e.request, url = new URL(req.url);
  if (req.method !== 'GET' || url.origin !== location.origin || req.headers.has('range')) return;
  const p = url.pathname;
  if (p.startsWith('/assets/vendor/')) e.respondWith(cacheFirst(req, STATIC));
  else if (p.startsWith('/assets/') || p.startsWith('/i18n/')) e.respondWith(staleWhileRevalidate(req, STATIC));
  else if (p.startsWith('/api/atlas/') && !p.includes('/search') && !p.includes('/word') && !p.includes('/validate')) e.respondWith(staleWhileRevalidate(req, DATA));
  else if (req.mode === 'navigate') e.respondWith(networkFirst(req, PAGES));
});
