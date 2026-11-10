// Reading view: the map follows whichever verse is in the reading band.
/* global maplibregl */
import { createMap, kindColorExpr, FONT, esc } from './map.js';

const data = JSON.parse(document.getElementById('readdata').textContent);
const points = data.points || [], byId = new Map(points.map((p) => [p.id, p]));
const verses = data.verses || [];
const placesOf = new Map(verses.map((v) => [v.id, (v.pl || []).map((id) => byId.get(id)).filter(Boolean)]));
const cap = document.getElementById('rcap');

const first = points[0];
const map = createMap('readmap', { center: [first.lng, first.lat], zoom: 5, pitch: 45, controls: true });
const feat = (p, active) => ({ type: 'Feature', geometry: { type: 'Point', coordinates: [p.lng, p.lat] }, properties: { id: p.id, n: p.n, k: p.k || '', active } });
let ready = false, current = null;

map.on('load', () => {
  map.addSource('all', { type: 'geojson', data: { type: 'FeatureCollection', features: points.map((p) => feat(p, false)) } });
  map.addSource('act', { type: 'geojson', data: { type: 'FeatureCollection', features: [] } });
  map.addLayer({ id: 'all', type: 'circle', source: 'all', paint: { 'circle-radius': 5, 'circle-color': kindColorExpr(), 'circle-opacity': 0.55, 'circle-stroke-color': '#0a0d1a', 'circle-stroke-width': 1 } });
  map.addLayer({ id: 'all-l', type: 'symbol', source: 'all', layout: { 'text-field': ['get', 'n'], 'text-font': FONT, 'text-size': 11.5, 'text-offset': [0, 1], 'text-anchor': 'top', 'text-optional': true }, paint: { 'text-color': '#e8e1c8', 'text-halo-color': 'rgba(5,7,16,.95)', 'text-halo-width': 1.4, 'text-opacity': 0.8 } });
  map.addLayer({ id: 'glow', type: 'circle', source: 'act', paint: { 'circle-radius': 26, 'circle-color': '#ffd24a', 'circle-opacity': 0.35, 'circle-blur': 0.8 } });
  map.addLayer({ id: 'act', type: 'circle', source: 'act', paint: { 'circle-radius': 8, 'circle-color': '#fff', 'circle-stroke-color': '#f5c451', 'circle-stroke-width': 3 } });
  map.addLayer({ id: 'act-l', type: 'symbol', source: 'act', layout: { 'text-field': ['get', 'n'], 'text-font': FONT, 'text-size': 14, 'text-offset': [0, 1.3], 'text-anchor': 'top', 'text-allow-overlap': true }, paint: { 'text-color': '#fff', 'text-halo-color': 'rgba(5,7,16,1)', 'text-halo-width': 2 } });
  const b = points.reduce((bb, p) => bb.extend([p.lng, p.lat]), new maplibregl.LngLatBounds([first.lng, first.lat], [first.lng, first.lat]));
  map.fitBounds(b, { padding: 60, maxZoom: 7, duration: 0 });
  ready = true;
  if (current) activate(current);
  map.on('click', 'all', (e) => { const f = e.features[0].properties; new maplibregl.Popup({ offset: 8 }).setLngLat(e.features[0].geometry.coordinates).setHTML(`<a href="/place/${f.id}"><b>${esc(f.n)}</b></a>`).addTo(map); });
});

function activate(id) {
  current = id; if (!ready) return;
  const ps = placesOf.get(id) || [];
  map.getSource('act').setData({ type: 'FeatureCollection', features: ps.map((p) => feat(p, true)) });
  if (ps.length) {
    cap.innerHTML = ps.map((p) => `<a href="/place/${p.id}">${esc(p.n)}</a>`).join(' · ');
    if (ps.length === 1) map.easeTo({ center: [ps[0].lng, ps[0].lat], zoom: Math.max(map.getZoom(), 6), duration: 900 });
    else map.fitBounds(ps.reduce((bb, p) => bb.extend([p.lng, p.lat]), new maplibregl.LngLatBounds([ps[0].lng, ps[0].lat], [ps[0].lng, ps[0].lat])), { padding: 90, maxZoom: 7, duration: 900 });
  }
}

// the verse nearest the upper-middle of the viewport is "current"
const els = [...document.querySelectorAll('.verses .v')];
const io = new IntersectionObserver((entries) => {
  for (const e of entries) e.target.classList.toggle('in', e.isIntersecting);
  const visible = els.filter((x) => x.classList.contains('in'));
  const withPlaces = visible.find((x) => (placesOf.get(x.dataset.id) || []).length);
  els.forEach((x) => x.classList.toggle('active', x === withPlaces));
  if (withPlaces && withPlaces.dataset.id !== current) activate(withPlaces.dataset.id);
}, { rootMargin: '-25% 0px -45% 0px' });
els.forEach((x) => io.observe(x));
els.forEach((x) => x.addEventListener('click', (e) => { if (e.target.closest('a')) return; activate(x.dataset.id); }));
