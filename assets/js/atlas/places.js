// Places browser: a map with a filterable list beside it.
/* global maplibregl */
import { createMap, kindColorExpr, KIND_COLOR, FONT, esc } from './map.js';

const $ = (id) => document.getElementById(id);
const world = await fetch('/api/atlas/world').then((r) => r.json());
const places = world.places.filter((p) => p.vc > 0).sort((a, b) => b.vc - a.vc);
const KINDS = ['City', 'Region', 'Mountain', 'Water', 'Valley', 'Island', 'Landmark', ''];
const kindName = (k) => k || 'Other';
const active = new Set(KINDS);
let query = '';

$('plist').innerHTML = `<h3>Places</h3>
  <input id="psearch" class="mb-input" placeholder="Search ${places.length.toLocaleString()} places…" autocomplete="off">
  <div class="kinds">${KINDS.map((k) => `<button class="mb-chip on" data-k="${k}" style="--c:${KIND_COLOR[k]}"><i></i>${kindName(k)}</button>`).join('')}</div>
  <div class="mb-muted" id="pcnt" style="font-size:.78rem;margin:.3rem 0 .5rem"></div>
  <div id="pitems" class="mb-scroll"></div>`;

const map = createMap('map', { center: [36, 33], zoom: 4.6, pitch: 35, bearing: 0 });
const fc = { type: 'FeatureCollection', features: places.map((p) => ({ type: 'Feature', geometry: { type: 'Point', coordinates: [p.lng, p.lat] }, properties: { id: p.id, n: p.n, k: p.k || '', vc: p.vc } })) };
let ready = false;
const popup = new maplibregl.Popup({ offset: 10, closeButton: true });

map.on('load', () => {
  map.addSource('p', { type: 'geojson', data: fc });
  map.addLayer({ id: 'dots', type: 'circle', source: 'p', paint: { 'circle-color': kindColorExpr(), 'circle-opacity': 0.9, 'circle-stroke-color': '#0a0d1a', 'circle-stroke-width': 0.8,
    'circle-radius': ['interpolate', ['linear'], ['zoom'], 3, ['interpolate', ['linear'], ['get', 'vc'], 1, 2, 60, 4, 600, 7], 9, ['interpolate', ['linear'], ['get', 'vc'], 1, 4, 60, 8, 600, 13]] } });
  for (const [id, min, vcMin] of [['l1', 3.5, 60], ['l2', 5.5, 12], ['l3', 7.2, 0]]) {
    map.addLayer({ id: 'lbl-' + id, type: 'symbol', source: 'p', minzoom: min, filter: ['>=', ['get', 'vc'], vcMin],
      layout: { 'text-field': ['get', 'n'], 'text-font': FONT, 'text-size': 12, 'text-offset': [0, 0.9], 'text-anchor': 'top', 'symbol-sort-key': ['-', 0, ['get', 'vc']], 'text-optional': true },
      paint: { 'text-color': '#f6efd6', 'text-halo-color': 'rgba(6,8,18,.95)', 'text-halo-width': 1.4 } });
  }
  map.on('click', 'dots', (e) => show(e.features[0].properties.id, e.features[0].geometry.coordinates));
  map.on('mouseenter', 'dots', () => { map.getCanvas().style.cursor = 'pointer'; });
  map.on('mouseleave', 'dots', () => { map.getCanvas().style.cursor = ''; });
  ready = true; refresh();
});

function show(id, coords) {
  const p = places.find((x) => x.id === id); if (!p) return;
  popup.setLngLat(coords || [p.lng, p.lat]).setHTML(`<b>${esc(p.n)}</b><br><span class="mb-muted">${esc(kindName(p.k))} · ${p.vc} verses</span><br><a href="/place/${p.id}">Open place →</a> · <a href="/atlas?place=${p.id}">time atlas</a>`).addTo(map);
}

function matches() {
  const q = query.toLowerCase();
  return places.filter((p) => active.has(p.k || '') && (!q || p.n.toLowerCase().includes(q)));
}
function refresh() {
  const ms = matches();
  $('pcnt').textContent = `${ms.length.toLocaleString()} places`;
  $('pitems').innerHTML = ms.slice(0, 120).map((p) => `<div class="pi" data-id="${p.id}"><i style="background:${KIND_COLOR[p.k || '']}"></i><span>${esc(p.n)}</span><small>${p.vc}</small></div>`).join('') + (ms.length > 120 ? `<div class="mb-muted" style="padding:.4rem">…refine the search to see more</div>` : '');
  $('pitems').querySelectorAll('.pi').forEach((el) => el.onclick = () => { const p = places.find((x) => x.id === el.dataset.id); map.flyTo({ center: [p.lng, p.lat], zoom: Math.max(map.getZoom(), 7.5), pitch: 50, duration: 1400 }); setTimeout(() => show(p.id), 900); });
  if (ready) {
    const filter = (ms.length === places.length) ? null : ['in', ['get', 'id'], ['literal', ms.map((p) => p.id)]];
    map.setFilter('dots', filter);
    for (const [id, vcMin] of [['l1', 60], ['l2', 12], ['l3', 0]]) map.setFilter('lbl-' + id, filter ? ['all', filter, ['>=', ['get', 'vc'], vcMin]] : ['>=', ['get', 'vc'], vcMin]);
  }
}
let t;
$('psearch').addEventListener('input', (e) => { clearTimeout(t); t = setTimeout(() => { query = e.target.value.trim(); refresh(); }, 120); });
$('plist').querySelectorAll('.kinds button').forEach((b) => b.onclick = () => { const k = b.dataset.k; if (active.has(k)) active.delete(k); else active.add(k); b.classList.toggle('on', active.has(k)); refresh(); });
refresh();
