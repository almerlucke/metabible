// Themes: follow an idea through the books, the eras and the map.
/* global maplibregl */
import { createMap, esc } from './map.js';
import { copyLink } from './share.js';

const $ = (id) => document.getElementById(id);
const list = await fetch('/api/atlas/themes').then((r) => r.json());
const qs = new URLSearchParams(location.search);
let current = qs.get('t') || list[0].id, map = null, ready = false;

$('thlist').innerHTML = list.map((t) => `<button class="thitem" data-id="${t.id}"><b>${esc(t.n)}</b><small>${t.vc.toLocaleString()}</small></button>`).join('');
$('thlist').querySelectorAll('.thitem').forEach((b) => b.onclick = () => select(b.dataset.id));
$('thmain').innerHTML = '<div id="thhead"></div><div class="thcols"><div id="thbody"></div><div id="thside"></div></div>';

async function select(id) {
  current = id; history.replaceState(null, '', '?t=' + id);
  $('thlist').querySelectorAll('.thitem').forEach((b) => b.classList.toggle('on', b.dataset.id === id));
  const d = await fetch('/api/atlas/theme/' + id).then((r) => r.json());
  const max = Math.max(1, ...d.books);
  const emax = Math.max(1, ...d.eras);
  $('thhead').innerHTML = `<h2>${esc(d.n)}</h2><p class="lead2">${esc(d.b)} <span class="mb-muted">${d.vc.toLocaleString()} verses</span></p>
    <div class="actions"><a class="mb-btn" href="/concordance?search=${encodeURIComponent((list.find((t) => t.id === id) || {}).n || '')}">Search the concordance</a><button class="mb-btn" id="thshare">🔗 Copy link</button></div>`;
  $('thshare').onclick = () => copyLink(location.href, $('thshare'));
  $('thbody').innerHTML = `<section><h3>Across the books</h3><div class="bookbars one" role="img" aria-label="Verses per book">${d.bookNames.map((n, i) => `<div class="bb" title="${esc(n)}: ${d.books[i]}"><i class="a" style="height:${(d.books[i] / max) * 100}%"></i></div>`).join('')}</div>
      <div class="bblbl"><span>Genesis</span><span>Malachi</span><span>Matthew</span><span>Revelation</span></div></section>
    <section><h3>Across the eras</h3><div class="erabars">${d.eraNames.map((n, i) => `<div class="eb"><span>${esc(n)}</span><div class="bar"><i style="width:${(d.eras[i] / emax) * 100}%"></i></div><small>${d.eras[i]}</small></div>`).join('')}</div></section>
    <section><h3>Some verses</h3>${d.sample.map((v) => `<p class="quote"><a href="/verse/${v.id}">${esc(v.l)}</a> ${v.h}</p>`).join('')}</section>`;
  $('thside').innerHTML = `<div class="card"><div class="mb-tag">Where it is named</div><div id="thmap" class="mb-minimap"></div></div>
    <div class="card"><div class="mb-tag">People most often alongside</div><div class="chiprow">${d.people.map((p) => `<a class="mb-chip" href="/person/${p.ID || p.id}">${esc(p.Name || p.name)}</a>`).join('')}</div></div>`;
  drawMap(d);
}
function drawMap(d) {
  if (map) { map.remove(); map = null; }
  const pts = d.places;
  if (!pts.length) { $('thmap').style.display = 'none'; return; }
  map = createMap('thmap', { center: [pts[0].lng, pts[0].lat], zoom: 4.5, pitch: 0, terrain: false, cooperative: true, controls: false });
  const fc = { type: 'FeatureCollection', features: pts.map((p, i) => ({ type: 'Feature', geometry: { type: 'Point', coordinates: [p.lng, p.lat] }, properties: { id: p.id, n: p.n, w: d.placeW[i] } })) };
  map.on('load', () => {
    map.addSource('t', { type: 'geojson', data: fc });
    map.addLayer({ id: 'glow', type: 'circle', source: 't', paint: { 'circle-radius': ['interpolate', ['linear'], ['get', 'w'], 1, 7, 40, 24], 'circle-color': '#f5c451', 'circle-opacity': 0.28, 'circle-blur': 0.7 } });
    map.addLayer({ id: 'dots', type: 'circle', source: 't', paint: { 'circle-radius': ['interpolate', ['linear'], ['get', 'w'], 1, 3, 40, 8], 'circle-color': '#fff3b0', 'circle-stroke-color': '#2a1030', 'circle-stroke-width': 1 } });
    const b = pts.reduce((bb, p) => bb.extend([p.lng, p.lat]), new maplibregl.LngLatBounds([pts[0].lng, pts[0].lat], [pts[0].lng, pts[0].lat]));
    map.fitBounds(b, { padding: 30, maxZoom: 6.5, duration: 0 });
    map.on('click', 'dots', (e) => { const f = e.features[0].properties; new maplibregl.Popup({ offset: 8 }).setLngLat(e.features[0].geometry.coordinates).setHTML(`<a href="/place/${f.id}"><b>${esc(f.n)}</b></a> · ${f.w}`).addTo(map); });
  });
}
select(list.find((t) => t.id === current) ? current : list[0].id);
