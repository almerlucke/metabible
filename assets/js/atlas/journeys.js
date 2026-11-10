// Journeys: animated routes with narration. Curated tours plus derived itineraries for any person.
/* global maplibregl */
import { createMap, setBasemap, FONT, esc, fmtYear, arc, ease } from './map.js';
import { copyLink } from './share.js';

const $ = (id) => document.getElementById(id);
const qs = new URLSearchParams(location.search);
const curated = await fetch('/api/atlas/journeys').then((r) => r.json());

const map = createMap('map', { center: [35.5, 32], zoom: 5, pitch: 50, bearing: 0 });
map.on('error', (e) => console.log('maplibre error', e && e.error && e.error.message));

const empty = { type: 'FeatureCollection', features: [] };
let ready = false;
let cur = null;          // current journey
let segs = [];           // polylines per segment (index i = from stop i-1 to i)
let idx = -1;            // current stop (-1 = intro)
let playing = false, token = 0, headAnim = null;
const compare = new Set();

map.on('load', () => {
  for (const id of ['cmp', 'done', 'all', 'stops', 'head']) map.addSource(id, { type: 'geojson', data: empty });
  map.addLayer({ id: 'cmp', type: 'line', source: 'cmp', layout: { 'line-cap': 'round' }, paint: { 'line-color': ['get', 'color'], 'line-width': 2.5, 'line-opacity': 0.8 } });
  map.addLayer({ id: 'all', type: 'line', source: 'all', layout: { 'line-cap': 'round', 'line-join': 'round' }, paint: { 'line-color': ['get', 'color'], 'line-width': ['interpolate', ['linear'], ['zoom'], 3, 2, 8, 3.5], 'line-opacity': 0.7 } });
  map.addLayer({ id: 'done-glow', type: 'line', source: 'done', layout: { 'line-cap': 'round', 'line-join': 'round' }, paint: { 'line-color': ['get', 'color'], 'line-width': 10, 'line-opacity': 0.25, 'line-blur': 6 } });
  map.addLayer({ id: 'done', type: 'line', source: 'done', layout: { 'line-cap': 'round', 'line-join': 'round' }, paint: { 'line-color': ['get', 'color'], 'line-width': 3.5 } });
  map.addLayer({ id: 'stops', type: 'circle', source: 'stops', paint: {
    'circle-radius': ['case', ['get', 'cur'], 11, 8], 'circle-color': ['case', ['get', 'reached'], ['get', 'color'], '#10142a'],
    'circle-stroke-color': ['get', 'color'], 'circle-stroke-width': 2, 'circle-pitch-alignment': 'map' } });
  map.addLayer({ id: 'stop-num', type: 'symbol', source: 'stops', layout: { 'text-field': ['get', 'i'], 'text-font': FONT, 'text-size': 11, 'text-allow-overlap': true, 'text-pitch-alignment': 'viewport' },
    paint: { 'text-color': ['case', ['get', 'reached'], '#0a0d1a', '#f6efd6'] } });
  map.addLayer({ id: 'stop-name', type: 'symbol', source: 'stops', minzoom: 3, layout: { 'text-field': ['get', 'n'], 'text-font': FONT, 'text-size': 12.5, 'text-offset': [0, 1.6], 'text-anchor': 'top', 'text-optional': true, 'text-pitch-alignment': 'viewport' },
    paint: { 'text-color': '#fff6dc', 'text-halo-color': 'rgba(5,7,16,.95)', 'text-halo-width': 1.6 } });
  map.addLayer({ id: 'head-glow', type: 'circle', source: 'head', paint: { 'circle-radius': 26, 'circle-color': ['get', 'color'], 'circle-opacity': 0.28, 'circle-blur': 0.8 } });
  map.addLayer({ id: 'head', type: 'circle', source: 'head', paint: { 'circle-radius': 6.5, 'circle-color': '#ffffff', 'circle-stroke-color': ['get', 'color'], 'circle-stroke-width': 3 } });
  map.on('click', 'stops', (e) => { const i = e.features[0].properties.i; if (cur) { stop(); goTo(+i - 1, false); } });
  map.on('mouseenter', 'stops', () => { map.getCanvas().style.cursor = 'pointer'; });
  map.on('mouseleave', 'stops', () => { map.getCanvas().style.cursor = ''; });
  ready = true;
  boot();
});

// ------------------------------------------------------------ geometry
function buildSegs(j) {
  segs = [null];
  for (let i = 1; i < j.stops.length; i++) {
    const a = [j.stops[i - 1].lng, j.stops[i - 1].lat], b = [j.stops[i].lng, j.stops[i].lat];
    segs.push(a[0] === b[0] && a[1] === b[1] ? [a, b] : arc(a, b, j.stops[i].sea || j.stops[i - 1].sea ? 0.04 : 0.11, 40));
  }
}
const line = (coords, color) => ({ type: 'Feature', geometry: { type: 'LineString', coordinates: coords }, properties: { color } });
function pointOn(poly, f) {
  const x = f * (poly.length - 1), i = Math.min(poly.length - 2, Math.floor(x)), t = x - i;
  return [poly[i][0] + (poly[i + 1][0] - poly[i][0]) * t, poly[i][1] + (poly[i + 1][1] - poly[i][1]) * t];
}
function donePath(upTo, frac) { // coordinates drawn so far: all segs ≤ upTo fully, plus partial seg upTo+1
  const out = [];
  for (let i = 1; i <= upTo; i++) out.push(...segs[i]);
  if (frac > 0 && segs[upTo + 1]) { const p = segs[upTo + 1], n = Math.max(1, Math.floor(frac * (p.length - 1))); out.push(...p.slice(0, n + 1)); out.push(pointOn(p, frac)); }
  return out;
}
function drawProgress(upTo, frac = 0) {
  if (!cur) return;
  const color = cur.color, path = donePath(upTo, frac);
  map.getSource('done').setData(path.length > 1 ? { type: 'FeatureCollection', features: [line(path, color)] } : empty);
  map.getSource('stops').setData({ type: 'FeatureCollection', features: cur.stops.map((s, i) => ({ type: 'Feature', geometry: { type: 'Point', coordinates: [s.lng, s.lat] },
    properties: { i: i + 1, n: s.n, color, reached: i <= upTo, cur: i === upTo } })) });
}
function setHead(coord) { map.getSource('head').setData(coord ? { type: 'FeatureCollection', features: [{ type: 'Feature', geometry: { type: 'Point', coordinates: coord }, properties: { color: cur.color } }] } : empty); }
function bounds(coords) { return coords.reduce((b, c) => b.extend(c), new maplibregl.LngLatBounds(coords[0], coords[0])); }
const PAD = () => (innerWidth > 900 ? { top: 90, bottom: 280, left: 380, right: 80 } : { top: 80, bottom: 300, left: 30, right: 30 });

// ------------------------------------------------------------ journeys list
function renderList() {
  $('jlist').innerHTML = `<h3>Journeys</h3>
    <input id="jsearch" class="mb-input" placeholder="Trace any person…" autocomplete="off"><div id="jres"></div>
    <div class="mb-tag" style="margin:.9rem 0 .4rem">Curated tours</div>` +
    curated.map((j) => `<div class="jc${cur && cur.id === j.id ? ' on' : ''}" data-id="${j.id}" style="--c:${j.color}">
      <div class="dot"></div><div class="tx"><b>${esc(j.title)}</b><span>${esc(j.sub || '')}</span><small>${j.stops.length} stops · ${j.km.toLocaleString()} km</small></div>
      <label class="cmp" title="Overlay on map"><input type="checkbox" data-cmp="${j.id}" ${compare.has(j.id) ? 'checked' : ''}><span>⧉</span></label></div>`).join('');
  $('jlist').querySelectorAll('.jc').forEach((el) => el.addEventListener('click', (e) => { if (e.target.closest('.cmp')) return; open(curated.find((j) => j.id === el.dataset.id)); }));
  $('jlist').querySelectorAll('[data-cmp]').forEach((c) => c.addEventListener('change', () => { c.checked ? compare.add(c.dataset.cmp) : compare.delete(c.dataset.cmp); drawCompare(); }));
  let t;
  $('jsearch').addEventListener('input', (e) => {
    clearTimeout(t); const q = e.target.value.trim();
    if (q.length < 2) { $('jres').innerHTML = ''; return; }
    t = setTimeout(async () => {
      const hits = (await fetch('/api/atlas/search?q=' + encodeURIComponent(q)).then((r) => r.json())).filter((h) => h.t === 'person').slice(0, 6);
      $('jres').innerHTML = hits.map((h) => `<div class="jr" data-id="${h.id}">${esc(h.l)} <span class="mb-muted">${esc(h.s.replace('person · ', ''))}</span></div>`).join('');
      $('jres').querySelectorAll('.jr').forEach((d) => d.onclick = () => openPerson(d.dataset.id));
    }, 160);
  });
}
function drawCompare() {
  if (!ready) return;
  const feats = [];
  for (const id of compare) { const j = curated.find((x) => x.id === id); if (!j) continue; const c = []; for (let i = 1; i < j.stops.length; i++) c.push(...arc([j.stops[i - 1].lng, j.stops[i - 1].lat], [j.stops[i].lng, j.stops[i].lat], 0.11, 30)); feats.push(line(c, j.color)); }
  map.getSource('cmp').setData({ type: 'FeatureCollection', features: feats });
  if (feats.length > 1) map.fitBounds(bounds(feats.flatMap((f) => f.geometry.coordinates)), { padding: PAD(), pitch: 30, duration: 1600 });
}

// ------------------------------------------------------------ opening + stepping
async function openPerson(id) {
  const j = await fetch('/api/atlas/journey/' + id).then((r) => (r.ok ? r.json() : null));
  if (!j || j.stops.length < 2) { $('jstep').innerHTML = '<div class="mb-muted">Not enough places mentioned alongside this person to draw a route.</div>'; $('jstep').style.display = 'block'; return; }
  j.intro = 'An approximate itinerary: the places named in the verses where this person appears, in the order of the text.';
  $('jres').innerHTML = ''; open(j);
}
function open(j) {
  stop(); cur = j; idx = -1; buildSegs(j); renderList();
  const all = j.stops.flatMap((s) => [[s.lng, s.lat]]);
  map.getSource('all').setData({ type: 'FeatureCollection', features: [line(segs.slice(1).flat(), j.color)] });
  drawProgress(-1); setHead(null);
  map.fitBounds(bounds(all), { padding: PAD(), pitch: 45, bearing: 0, duration: 2200, maxZoom: 8 });
  renderCard();
  syncUrl();
}
function syncUrl() {
  if (!cur) return;
  const base = cur.id.startsWith('person-') ? '?person=' + cur.person : '?j=' + cur.id;
  history.replaceState(null, '', base + (idx >= 0 ? '&s=' + (idx + 1) : ''));
}
function stop() { playing = false; token++; if (headAnim) cancelAnimationFrame(headAnim); renderControls(); }

async function goTo(i, animate) {
  if (!cur) return;
  const my = ++token; i = Math.max(-1, Math.min(cur.stops.length - 1, i));
  if (animate && i === idx + 1 && i > 0) {
    const a = cur.stops[i - 1], b = cur.stops[i];
    const d = segs[i];
    map.fitBounds(bounds(d), { padding: PAD(), pitch: 55, bearing: ((b.lng - a.lng) * 6) % 40, duration: 1700, maxZoom: 7.5 });
    const dur = Math.min(4200, Math.max(1800, Math.hypot(b.lng - a.lng, b.lat - a.lat) * 420));
    await new Promise((res) => {
      const t0 = performance.now();
      const step = (now) => {
        if (my !== token) return res();
        const f = Math.min(1, (now - t0) / dur), e = ease(f);
        setHead(pointOn(d, e)); drawProgress(i - 1, e);
        if (f < 1) headAnim = requestAnimationFrame(step); else res();
      };
      headAnim = requestAnimationFrame(step);
    });
    if (my !== token) return;
  } else if (i >= 0) {
    map.easeTo({ center: [cur.stops[i].lng, cur.stops[i].lat], zoom: Math.max(map.getZoom(), 6.2), pitch: 55, duration: 1400 });
  }
  idx = i; syncUrl(); drawProgress(i);
  setHead(i >= 0 ? [cur.stops[i].lng, cur.stops[i].lat] : null);
  renderCard();
  if (playing && my === token) {
    if (i >= cur.stops.length - 1) { playing = false; renderControls(); return; }
    const s = cur.stops[i];
    await new Promise((r) => setTimeout(r, Math.min(7000, 2200 + ((s.note || '').length + (s.tx || '').length) * 22)));
    if (playing && my === token) goTo(i + 1, true);
  }
}

// ------------------------------------------------------------ narration card
function renderControls() {
  const b = $('jplay'); if (b) b.textContent = playing ? '❚❚' : '▶';
}
function renderCard() {
  const el = $('jstep'); el.style.display = 'block';
  const n = cur.stops.length;
  const prog = `<div class="prog">${cur.stops.map((s, i) => `<i class="${i <= idx ? 'on' : ''}${i === idx ? ' cur' : ''}" data-i="${i}" title="${esc(s.n)}" style="--c:${cur.color}"></i>`).join('')}</div>`;
  const ctl = `<div class="ctl"><button class="mb-btn icon" id="jprev">◀</button><button class="mb-btn icon" id="jplay">▶</button><button class="mb-btn icon" id="jnext">▶|</button><button class="mb-btn" id="jshare" style="margin-left:auto">🔗 Share</button></div>`;
  if (idx < 0) {
    el.innerHTML = `<div class="mb-tag" style="color:${cur.color}">${cur.era ? esc(cur.era) + ' · ' : ''}${n} stops · ${cur.km.toLocaleString()} km</div>
      <h2>${esc(cur.title)}</h2><div class="mb-muted" style="margin-bottom:.4rem">${esc(cur.sub || cur.subtitle || '')}</div>
      <p class="intro">${esc(cur.intro || '')}</p>${ctl}${prog}`;
  } else {
    const s = cur.stops[idx];
    let km = 0; for (let i = 1; i <= idx; i++) km += hav(cur.stops[i - 1], cur.stops[i]);
    el.innerHTML = `<div class="mb-tag" style="color:${cur.color}">Stop ${idx + 1} of ${n}${s.y ? ' · ' + fmtYear(s.y) : ''} · ${Math.round(km).toLocaleString()} km travelled</div>
      <h2>${s.pl ? `<a href="/place/${s.pl}">${esc(s.n)}</a>` : esc(s.n)}</h2>
      ${s.note ? `<p class="intro">${esc(s.note)}</p>` : ''}
      ${s.tx ? `<div class="mb-verse">“${esc(s.tx)}”</div>` : ''}
      ${s.vid ? `<a class="ref" href="/verse/${s.vid}">${esc(s.lab || s.ref)} →</a>` : ''}${ctl}${prog}`;
  }
  $('jshare').onclick = () => copyLink(location.href, $('jshare'), '🔗 Share');
  $('jprev').onclick = () => { stop(); goTo(idx - 1, false); };
  $('jnext').onclick = () => { stop(); goTo(idx + 1, true); };
  $('jplay').onclick = () => { if (playing) stop(); else { playing = true; renderControls(); goTo(idx < 0 ? 0 : (idx >= n - 1 ? 0 : idx + 1), idx >= 0 && idx < n - 1); } };
  el.querySelectorAll('.prog i').forEach((d) => d.onclick = () => { stop(); goTo(+d.dataset.i, false); });
  renderControls();
}
function hav(a, b) { const r = 6371, p = Math.PI / 180, dp = (b.lat - a.lat) * p, dl = (b.lng - a.lng) * p, h = Math.sin(dp / 2) ** 2 + Math.cos(a.lat * p) * Math.cos(b.lat * p) * Math.sin(dl / 2) ** 2; return 2 * r * Math.asin(Math.sqrt(h)); }

addEventListener('keydown', (e) => {
  if (e.target.tagName === 'INPUT') return;
  if (e.key === 'ArrowRight') { stop(); goTo(idx + 1, true); } else if (e.key === 'ArrowLeft') { stop(); goTo(idx - 1, false); }
  else if (e.code === 'Space' && cur) { e.preventDefault(); $('jplay').click(); }
});

function boot() {
  renderList();
  const jid = qs.get('j'), pid = qs.get('person');
  const stopNo = +(qs.get('s') || 0);
  if (pid) openPerson(pid).then(() => { if (stopNo) goTo(stopNo - 1, false); });
  else { const j = curated.find((x) => x.id === jid) || curated[0]; open(j); if (stopNo) goTo(stopNo - 1, false); else if (qs.get('play')) { playing = true; goTo(0, false); } }
}
void setBasemap;
