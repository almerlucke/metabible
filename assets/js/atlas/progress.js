// "My reading": what you have read, the places you have met, the centuries you have covered.
/* global maplibregl */
import { createMap, FONT, esc } from './map.js';

const $ = (id) => document.getElementById(id);
const idx = await fetch('/api/atlas/progress-index').then((r) => r.json());
const fmt = (y) => window.mbYear(y);
const pct = (y) => Math.min(100, Math.max(0, ((y + 4100) / 4200) * 100));

function load() { try { return JSON.parse(localStorage.getItem('mb.progress')) || { read: {} }; } catch (e) { return { read: {} }; } }
function save(p) { try { localStorage.setItem('mb.progress', JSON.stringify(p)); } catch (e) { /* ignore */ } }
let map = null;

function render() {
  const p = load(); p.read = p.read || {};
  const all = idx.books.flatMap((b) => b.ch.map((c) => ({ ...c, book: b })));
  const done = all.filter((c) => p.read[c.id]);
  const met = new Set(done.flatMap((c) => c.pl || []));
  const people = done.reduce((a, c) => a + c.pe, 0);
  const lastCh = all.find((c) => c.id === p.last);
  const ys = done.filter((c) => c.y0 || c.y1);
  const from = ys.length ? Math.min(...ys.map((c) => c.y0)) : null, to = ys.length ? Math.max(...ys.map((c) => c.y1)) : null;

  let h = `<div class="pstats">
    <div><b>${done.length}</b><span>${tx('js.progress.chaptersread', { n: all.length })}</span></div>
    <div><b>${Math.round((done.length / all.length) * 100)}%</b><span>${tx('js.progress.of.the.bible')}</span></div>
    <div><b>${met.size}</b><span>${tx('js.progress.places.met')}</span></div>
    <div><b>${from == null ? '—' : fmt(from) + ' – ' + fmt(to)}</b><span>${tx('js.progress.time.covered')}</span></div></div>`;
  if (lastCh) h += `<p class="cont"><a class="mb-btn on" href="/chapter/${lastCh.id}">Continue: ${esc(lastCh.book.n)} ${lastCh.n} →</a></p>`;
  else h += '<p class="mb-muted">' + tx('js.progress.nothing.yet.open.any') + '</p>';
  if (from != null) h += `<div class="cover" title="${tx('js.progress.time.covered2')}"><i style="left:${pct(from)}%;width:${Math.max(1, pct(to) - pct(from))}%"></i></div><div class="bblbl"><span>${tx('js.migrations.4100.bc')}</span><span>${tx('js.migrations.2000.bc')}</span><span>${tx('js.progress.ad.100')}</span></div>`;
  h += '<h3>' + tx('js.progress.book.by.book') + '</h3><div class="pbooks">' + idx.books.map((b) => {
    const n = b.ch.filter((c) => p.read[c.id]).length;
    return `<div class="pb"><a class="bn" href="/book/${b.id}">${esc(b.n)}</a><div class="cells">${b.ch.map((c) => `<a href="/chapter/${c.id}" class="${p.read[c.id] ? 'on' : ''}" title="${esc(b.n)} ${c.n}"></a>`).join('')}</div><small>${n}/${b.ch.length}</small></div>`;
  }).join('') + '</div>';
  h += `<h3>${tx('js.progress.places.you.have.met')}</h3><div id="pmap" class="mb-minimap wide"></div>
    <h3>${tx('js.progress.keep.or.move.your')}</h3><div class="actions"><button class="mb-btn" id="pexp">${tx('js.progress.download')}</button><label class="mb-btn">${tx('js.progress.import')}<input type="file" id="pimp" accept="application/json" hidden></label><button class="mb-btn" id="preset">${tx('js.progress.reset')}</button></div>
    <p class="mb-muted">${tx('js.progress.your.reading.is.stored')}</p>`;
  $('prog').innerHTML = h;

  $('pexp').onclick = () => { const a = document.createElement('a'); a.href = URL.createObjectURL(new Blob([JSON.stringify(load(), null, 1)], { type: 'application/json' })); a.download = 'metabible-reading.json'; a.click(); };
  $('pimp').onchange = async (e) => { try { const j = JSON.parse(await e.target.files[0].text()); if (j && typeof j.read === 'object') { const cur = load(); cur.read = { ...cur.read, ...j.read }; cur.last = j.last || cur.last; save(cur); render(); } } catch (err) { alert(tx('js.progress.badfile')); } };
  $('preset').onclick = () => { if (confirm(tx('js.progress.confirm'))) { save({ read: {} }); render(); } };
  drawMap([...met].map((i) => idx.places[i]));
}

function drawMap(pts) {
  if (map) { map.remove(); map = null; }
  const el = $('pmap');
  if (!pts.length) { el.style.display = 'none'; return; }
  map = createMap('pmap', { center: [pts[0].lng, pts[0].lat], zoom: 4.5, pitch: 35, cooperative: true });
  map.on('load', () => {
    map.addSource('p', { type: 'geojson', data: { type: 'FeatureCollection', features: pts.map((p) => ({ type: 'Feature', geometry: { type: 'Point', coordinates: [p.lng, p.lat] }, properties: { id: p.id, n: p.n } })) } });
    map.addLayer({ id: 'g', type: 'circle', source: 'p', paint: { 'circle-radius': 14, 'circle-color': '#7df9ff', 'circle-opacity': 0.25, 'circle-blur': 0.8 } });
    map.addLayer({ id: 'd', type: 'circle', source: 'p', paint: { 'circle-radius': 4.5, 'circle-color': '#fff', 'circle-stroke-color': '#7df9ff', 'circle-stroke-width': 1.5 } });
    map.addLayer({ id: 'l', type: 'symbol', source: 'p', minzoom: 5, layout: { 'text-field': ['get', 'n'], 'text-font': FONT, 'text-size': 11.5, 'text-offset': [0, 1], 'text-anchor': 'top', 'text-optional': true }, paint: { 'text-color': '#e8fbff', 'text-halo-color': 'rgba(5,7,16,.95)', 'text-halo-width': 1.4 } });
    const b = pts.reduce((bb, p) => bb.extend([p.lng, p.lat]), new maplibregl.LngLatBounds([pts[0].lng, pts[0].lat], [pts[0].lng, pts[0].lat]));
    map.fitBounds(b, { padding: 50, maxZoom: 7, duration: 0 });
  });
}
render();
