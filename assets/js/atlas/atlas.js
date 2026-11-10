// Time Atlas: scrub or play 4004 BC → AD 100 over a 3D relief map.
/* global maplibregl */
import { createMap, setBasemap, setTerrain, BASEMAPS, kindColorExpr, FONT, esc, fmtYear, arc } from './map.js';

const T0 = -4100, T1 = 100, BIN = 10;
const $ = (id) => document.getElementById(id);
const qs = new URLSearchParams(location.search);

const [world, peopleData] = await Promise.all([
  fetch('/api/atlas/world').then((r) => r.json()),
  fetch('/api/atlas/people').then((r) => r.json()),
]);
const { places, events, eras } = world;
const personName = new Map(peopleData.people.map((p) => [p.id, p.n]));
const placeById = new Map(places.map((p) => [p.id, p]));
for (const e of events) { e.pl = e.pl || []; e.pe = e.pe || []; e.pr = e.pr || []; e.po = e.po || []; }
events.sort((a, b) => a.y - b.y);

// ---------------------------------------------------------------- time scale
// Slider position s∈[0,1] maps to years through a density-weighted CDF, so busy
// stretches (Exodus, Kings, Gospels) get more room than the empty millennia.
const nb = Math.ceil((T1 - T0) / BIN);
const w = new Float64Array(nb).fill(0.4);
for (const e of events) {
  const i = Math.min(nb - 1, Math.max(0, Math.floor((e.y - T0) / BIN)));
  for (let k = -2; k <= 2; k++) if (w[i + k] !== undefined) w[i + k] += (3 - Math.abs(k)) * 1.6;
}
const cum = new Float64Array(nb + 1);
for (let i = 0; i < nb; i++) cum[i + 1] = cum[i] + w[i];
const total = cum[nb];
function yearToS(y) {
  const f = Math.min(nb - 1e-9, Math.max(0, (y - T0) / BIN)), i = Math.floor(f);
  return (cum[i] + w[i] * (f - i)) / total;
}
function sToYear(s) {
  const t = Math.min(0.999999, Math.max(0, s)) * total;
  let lo = 0, hi = nb;
  while (lo < hi - 1) { const m = (lo + hi) >> 1; if (cum[m] <= t) lo = m; else hi = m; }
  return T0 + (lo + (t - cum[lo]) / w[lo]) * BIN;
}
const eraAt = (y) => eras.find((e) => y >= e.from && y < e.to) || eras[eras.length - 1];
function latestEventIdx(y) { // last event with e.y <= y, or -1
  let lo = 0, hi = events.length;
  while (lo < hi) { const m = (lo + hi) >> 1; if (events[m].y <= y) lo = m + 1; else hi = m; }
  return lo - 1;
}
const placeY0 = places.filter((p) => p.vc > 0).map((p) => p.y0).sort((a, b) => a - b);
function placesSoFar(y) { let lo = 0, hi = placeY0.length; while (lo < hi) { const m = (lo + hi) >> 1; if (placeY0[m] <= y) lo = m + 1; else hi = m; } return lo; }

// ---------------------------------------------------------------- state
let s = 0, year = T0, playing = false, speed = 1, follow = true, lastIdx = -2, lastFly = 0;
const RUN_SECONDS = 260;
let personRoute = null, personInfo = null;

// ---------------------------------------------------------------- map
const startCenter = [35.3, 31.9];
const map = createMap('map', { center: [35.6, 31.6], zoom: 5.3, pitch: 52, bearing: -8 });
map.on('error', (e) => console.log('maplibre error', e && e.error && e.error.message));
map.addControl(new maplibregl.ScaleControl({ unit: 'metric' }), 'bottom-left');

const placeFC = { type: 'FeatureCollection', features: places.filter((p) => p.vc > 0).map((p) => ({
  type: 'Feature', geometry: { type: 'Point', coordinates: [p.lng, p.lat] }, properties: { id: p.id, n: p.n, k: p.k, vc: p.vc, y0: p.y0 } })) };

const eventFC = { type: 'FeatureCollection', features: [] };
const chainFC = { type: 'FeatureCollection', features: [] };
for (const e of events) {
  const pl = e.pl.map((id) => placeById.get(id)).filter(Boolean);
  for (const p of pl) eventFC.features.push({ type: 'Feature', geometry: { type: 'Point', coordinates: [p.lng, p.lat] }, properties: { eid: e.id, t: e.t, y: e.y } });
}
const evById = new Map(events.map((e) => [e.id, e]));
const centroid = (e) => {
  const pl = e.pl.map((id) => placeById.get(id)).filter(Boolean);
  if (!pl.length) return null;
  return [pl.reduce((a, p) => a + p.lng, 0) / pl.length, pl.reduce((a, p) => a + p.lat, 0) / pl.length];
};
for (const e of events) {
  const to = centroid(e);
  if (!to) continue;
  for (const pid of e.pr || []) {
    const pe = evById.get(pid); const from = pe && centroid(pe);
    if (!from) continue;
    const d = Math.hypot(to[0] - from[0], to[1] - from[1]);
    if (d < 0.05) continue;
    chainFC.features.push({ type: 'Feature', geometry: { type: 'LineString', coordinates: arc(from, to, 0.16) }, properties: { y: e.y } });
  }
}
const empty = { type: 'FeatureCollection', features: [] };

map.on('load', () => {
  map.addSource('places', { type: 'geojson', data: placeFC });
  map.addSource('events', { type: 'geojson', data: eventFC });
  map.addSource('chains', { type: 'geojson', data: chainFC });
  map.addSource('pulse', { type: 'geojson', data: empty });
  map.addSource('route', { type: 'geojson', data: empty });
  map.addSource('routepts', { type: 'geojson', data: empty });

  map.addLayer({ id: 'chains', type: 'line', source: 'chains', filter: ['<=', ['get', 'y'], year],
    paint: { 'line-color': '#f5c451', 'line-opacity': 0.5, 'line-width': ['interpolate', ['linear'], ['zoom'], 3, 0.8, 8, 2], 'line-dasharray': [2, 2] } });
  map.addLayer({ id: 'route', type: 'line', source: 'route', layout: { 'line-cap': 'round', 'line-join': 'round' },
    paint: { 'line-color': '#7df9ff', 'line-width': 3, 'line-opacity': 0.9, 'line-blur': 0.5 } });
  map.addLayer({ id: 'place-dots', type: 'circle', source: 'places', filter: ['<=', ['get', 'y0'], year],
    paint: {
      'circle-color': kindColorExpr(),
      'circle-radius': ['interpolate', ['linear'], ['zoom'],
        3, ['interpolate', ['linear'], ['get', 'vc'], 1, 1.8, 60, 3.4, 600, 6],
        9, ['interpolate', ['linear'], ['get', 'vc'], 1, 4, 60, 7, 600, 12]],
      'circle-opacity': 0.9, 'circle-stroke-color': '#0a0d1a', 'circle-stroke-width': 0.8,
    } });
  for (const [id, min, vcMin] of [['l1', 3.5, 60], ['l2', 5.5, 12], ['l3', 7.2, 0]]) {
    map.addLayer({ id: 'place-' + id, type: 'symbol', source: 'places', minzoom: min,
      filter: ['all', ['<=', ['get', 'y0'], year], ['>=', ['get', 'vc'], vcMin]],
      layout: { 'text-field': ['get', 'n'], 'text-font': FONT, 'text-size': 12, 'text-offset': [0, 0.9], 'text-anchor': 'top', 'symbol-sort-key': ['-', 0, ['get', 'vc']], 'text-optional': true },
      paint: { 'text-color': '#f6efd6', 'text-halo-color': 'rgba(6,8,18,.95)', 'text-halo-width': 1.4 } });
  }
  map.addLayer({ id: 'event-dots', type: 'circle', source: 'events', filter: ['<=', ['get', 'y'], year],
    paint: { 'circle-color': 'rgba(0,0,0,0)', 'circle-stroke-color': '#ffe08a', 'circle-stroke-width': 1.6, 'circle-radius': ['interpolate', ['linear'], ['zoom'], 3, 5, 9, 11], 'circle-stroke-opacity': 0.8 } });
  map.addLayer({ id: 'pulse', type: 'circle', source: 'pulse',
    paint: { 'circle-color': '#ffd24a', 'circle-radius': ['get', 'r'], 'circle-opacity': ['get', 'o'], 'circle-blur': 0.6, 'circle-pitch-alignment': 'map' } });
  map.addLayer({ id: 'routepts', type: 'circle', source: 'routepts', filter: ['<=', ['get', 'y'], year],
    paint: { 'circle-color': '#7df9ff', 'circle-radius': 4.5, 'circle-stroke-color': '#06223a', 'circle-stroke-width': 1.5 } });

  map.on('click', 'place-dots', (e) => {
    const f = e.features[0].properties;
    new maplibregl.Popup({ offset: 10 }).setLngLat(e.features[0].geometry.coordinates)
      .setHTML(`<b>${esc(f.n)}</b><br><span class="mb-muted">${esc(f.k || 'place')} · ${f.vc} verses</span><br><a href="/place/${f.id}">Open place →</a>`).addTo(map);
  });
  map.on('click', 'event-dots', (e) => {
    const f = e.features[0].properties; const ev = evById.get(f.eid);
    if (ev) showEventPopup(ev, e.features[0].geometry.coordinates);
  });
  for (const l of ['place-dots', 'event-dots']) {
    map.on('mouseenter', l, () => { map.getCanvas().style.cursor = 'pointer'; });
    map.on('mouseleave', l, () => { map.getCanvas().style.cursor = ''; });
  }
  map.on('dragstart', () => { if (playing) { follow = false; syncLayers(); } });

  ready = true;
  applyYear(true);
  boot();
});
let ready = false;

function showEventPopup(ev, lngLat) {
  const people = (ev.pe || []).slice(0, 6).map((id) => `<a class="mb-chip" href="/person/${id}">${esc(personName.get(id) || '…')}</a>`).join('');
  new maplibregl.Popup({ offset: 12, maxWidth: '320px' }).setLngLat(lngLat).setHTML(
    `<div class="mb-tag">${esc(ev.dt)}</div><b>${esc(ev.t)}</b>` +
    (ev.tx ? `<div class="mb-verse" style="font-size:.85rem;margin:.35rem 0">“${esc(ev.tx)}”</div>` : '') +
    (ev.vid ? `<a href="/verse/${ev.vid}">${esc(ev.ref)} →</a><br>` : '') + people).addTo(map);
}

// ---------------------------------------------------------------- UI chrome
$('yearbox').innerHTML = '<div class="yr"></div><div class="era"></div><div class="cnt mb-muted"></div>';
$('layers').innerHTML = `
  <h3>Layers</h3>
  <div class="row"><label>Map</label><select id="bm" class="mb-input">${Object.entries(BASEMAPS).map(([k, v]) => `<option value="${k}">${v.label}</option>`).join('')}</select></div>
  <label class="chk"><input type="checkbox" id="t-terrain" checked> 3D terrain</label>
  <label class="chk"><input type="checkbox" id="t-chains" checked> Cause &amp; effect arcs</label>
  <label class="chk"><input type="checkbox" id="t-events" checked> Event rings</label>
  <label class="chk"><input type="checkbox" id="t-follow" checked> Auto-camera</label>
  <div class="row" style="margin-top:.5rem"><input id="psearch" class="mb-input" placeholder="Follow a person…" autocomplete="off"></div>
  <div id="psres"></div><div id="pcur"></div>`;
$('timeline').innerHTML = `
  <button class="mb-btn icon" id="play" title="Play / pause (space)">▶</button>
  <div class="speeds">${[0.5, 1, 2, 4].map((x) => `<button class="mb-btn sp${x === 1 ? ' on' : ''}" data-s="${x}">${x}×</button>`).join('')}</div>
  <div class="tl-wrap"><canvas id="tl"></canvas><div id="tltip"></div></div>`;
function syncLayers() { $('t-follow').checked = follow; }

$('bm').onchange = (e) => setBasemap(map, e.target.value);
$('t-terrain').onchange = (e) => setTerrain(map, e.target.checked);
$('t-chains').onchange = (e) => map.setLayoutProperty('chains', 'visibility', e.target.checked ? 'visible' : 'none');
$('t-events').onchange = (e) => map.setLayoutProperty('event-dots', 'visibility', e.target.checked ? 'visible' : 'none');
$('t-follow').onchange = (e) => { follow = e.target.checked; };
$('play').onclick = () => setPlaying(!playing);
document.querySelectorAll('.sp').forEach((b) => { b.onclick = () => { speed = +b.dataset.s; document.querySelectorAll('.sp').forEach((x) => x.classList.toggle('on', x === b)); }; });
addEventListener('keydown', (e) => {
  if (e.target.tagName === 'INPUT' || e.target.tagName === 'SELECT') return;
  if (e.code === 'Space') { e.preventDefault(); setPlaying(!playing); }
  else if (e.key === 'ArrowRight') stepEvent(1);
  else if (e.key === 'ArrowLeft') stepEvent(-1);
});
function setPlaying(p) { playing = p; $('play').textContent = p ? '❚❚' : '▶'; last = performance.now(); }
function stepEvent(d) {
  const i = Math.max(0, Math.min(events.length - 1, latestEventIdx(year) + d));
  setS(yearToS(events[i].y + 0.01)); flyToEvent(events[i], true);
}

// ---------------------------------------------------------------- timeline canvas
const tl = $('tl'), tctx = tl.getContext('2d');
let tlW = 0, tlH = 0, hoverS = null;
function sizeTl() {
  const r = tl.getBoundingClientRect(); const dpr = devicePixelRatio || 1;
  tlW = r.width; tlH = r.height; tl.width = tlW * dpr; tl.height = tlH * dpr; tctx.setTransform(dpr, 0, 0, dpr, 0, 0); drawTl();
}
function drawTl() {
  const c = tctx; c.clearRect(0, 0, tlW, tlH);
  const bandY = 18, bandH = 26;
  c.font = '600 10px Inter, sans-serif'; c.textBaseline = 'middle';
  for (const e of eras) {
    const x0 = yearToS(e.from) * tlW, x1 = yearToS(e.to) * tlW;
    c.fillStyle = e.c + '55'; c.fillRect(x0, bandY, Math.max(1, x1 - x0 - 1), bandH);
    c.fillStyle = e.c; c.fillRect(x0, bandY, Math.max(1, x1 - x0 - 1), 2);
    if (x1 - x0 > 56) { c.fillStyle = '#f3ecd6'; c.save(); c.beginPath(); c.rect(x0, bandY, x1 - x0 - 2, bandH); c.clip(); c.fillText(e.n, x0 + 5, bandY + bandH / 2 + 1); c.restore(); }
  }
  if (personInfo) { // lifespan of the followed person
    const x0 = yearToS(personInfo.y0) * tlW, x1 = yearToS(Math.max(personInfo.y1, personInfo.y0 + 1)) * tlW;
    c.fillStyle = '#7df9ff'; c.fillRect(x0, 8, Math.max(3, x1 - x0), 4);
  }
  c.strokeStyle = 'rgba(245,196,81,.55)'; c.lineWidth = 1; c.beginPath();
  for (const e of events) { const x = Math.round(yearToS(e.y) * tlW) + .5; c.moveTo(x, bandY + bandH + 3); c.lineTo(x, bandY + bandH + (e.pl.length ? 14 : 9)); }
  c.stroke();
  c.fillStyle = 'rgba(236,232,220,.6)'; c.font = '10px Inter, sans-serif'; c.textBaseline = 'alphabetic';
  for (let y = -4000; y <= 100; y += y < -1000 ? 1000 : 500) { const x = yearToS(y) * tlW; c.fillText(y === 0 ? '1' : fmtYear(y).replace(' BC', 'BC').replace(' AD', 'AD'), Math.min(tlW - 30, Math.max(0, x - 8)), tlH - 3); }
  const px = s * tlW;
  c.strokeStyle = '#ffe08a'; c.lineWidth = 2; c.beginPath(); c.moveTo(px, 4); c.lineTo(px, tlH - 14); c.stroke();
  c.fillStyle = '#ffe08a'; c.beginPath(); c.arc(px, 6, 5, 0, 7); c.fill();
  if (hoverS != null) { c.strokeStyle = 'rgba(255,255,255,.35)'; c.lineWidth = 1; const hx = hoverS * tlW; c.beginPath(); c.moveTo(hx, 4); c.lineTo(hx, tlH - 14); c.stroke(); }
}
let dragging = false;
function pointerS(e) { const r = tl.getBoundingClientRect(); return Math.min(1, Math.max(0, (e.clientX - r.left) / r.width)); }
tl.addEventListener('pointerdown', (e) => { dragging = true; tl.setPointerCapture(e.pointerId); setS(pointerS(e)); });
tl.addEventListener('pointermove', (e) => {
  hoverS = pointerS(e); if (dragging) setS(hoverS);
  const tip = $('tltip'); tip.style.display = 'block'; tip.style.left = hoverS * tlW + 'px'; tip.textContent = fmtYear(sToYear(hoverS)); drawTl();
});
tl.addEventListener('pointerup', () => { dragging = false; });
tl.addEventListener('pointerleave', () => { hoverS = null; $('tltip').style.display = 'none'; drawTl(); });
addEventListener('resize', sizeTl);

// ---------------------------------------------------------------- feed
function renderFeed(idx) {
  const items = [];
  for (let i = idx; i >= 0 && items.length < 6; i--) items.push(events[i]);
  $('feed').innerHTML = `<div class="mb-tag" style="margin-bottom:.5rem">Unfolding events</div>` + (items.length ? items.map((e, k) => {
    const pl = e.pl.slice(0, 3).map((id) => placeById.get(id)).filter(Boolean).map((p) => `<a class="mb-chip" href="/place/${p.id}">📍 ${esc(p.n)}</a>`).join('');
    const pe = (e.pe || []).slice(0, 4).map((id) => `<a class="mb-chip" href="/person/${id}">${esc(personName.get(id) || '…')}</a>`).join('');
    return `<article class="ev${k === 0 ? ' new' : ''}" data-id="${e.id}"><div class="mb-tag">${esc(e.dt)}</div><div class="t">${esc(e.t)}</div>` +
      (e.tx ? `<div class="mb-verse">“${esc(e.tx.length > 160 ? e.tx.slice(0, 160) + '…' : e.tx)}”</div>` : '') +
      (e.vid ? `<a class="ref" href="/verse/${e.vid}">${esc(e.ref)}</a>` : '') + `<div>${pl}${pe}</div></article>`;
  }).join('') : '<div class="mb-muted">In the beginning…</div>');
  $('feed').querySelectorAll('article.ev').forEach((a) => a.addEventListener('click', (ev) => {
    if (ev.target.closest('a')) return;
    const e = evById.get(a.dataset.id); if (e) { setPlaying(false); flyToEvent(e, true); }
  }));
}

function flyToEvent(e, force) {
  const pl = e.pl.map((id) => placeById.get(id)).filter(Boolean);
  if (!pl.length || !ready) return;
  const lngs = pl.map((p) => p.lng), lats = pl.map((p) => p.lat);
  const spread = Math.max(Math.max(...lngs) - Math.min(...lngs), Math.max(...lats) - Math.min(...lats));
  const zoom = spread > 8 ? 4.3 : spread > 3 ? 5.4 : spread > 0.8 ? 6.6 : 7.6;
  map.easeTo({ center: [(Math.min(...lngs) + Math.max(...lngs)) / 2, (Math.min(...lats) + Math.max(...lats)) / 2], zoom, pitch: 55, bearing: ((e.y * 7) % 40) - 20, duration: force ? 1800 : 2600, essential: true });
  lastFly = performance.now();
}

// ---------------------------------------------------------------- person route
let psTimer;
$('psearch').addEventListener('input', (e) => {
  clearTimeout(psTimer);
  const q = e.target.value.trim();
  if (q.length < 2) { $('psres').innerHTML = ''; return; }
  psTimer = setTimeout(async () => {
    const hits = (await fetch('/api/atlas/search?q=' + encodeURIComponent(q)).then((r) => r.json())).filter((h) => h.t === 'person').slice(0, 5);
    $('psres').innerHTML = hits.map((h) => `<div class="ps" data-id="${h.id}">${esc(h.l)} <span class="mb-muted">${esc(h.s.replace('person · ', ''))}</span></div>`).join('');
    $('psres').querySelectorAll('.ps').forEach((d) => d.onclick = () => { followPerson(d.dataset.id); $('psres').innerHTML = ''; $('psearch').value = ''; });
  }, 160);
});
async function followPerson(id) {
  const j = await fetch('/api/atlas/journey/' + id).then((r) => (r.ok ? r.json() : null));
  const p = peopleData.people.find((x) => x.id === id);
  if (!j || !map.getSource('route')) return;
  personInfo = p ? { y0: p.y0, y1: p.y1, name: p.n } : null;
  const stops = j.stops.filter((x) => x.lat != null);
  // route line is revealed by year through the 'y' property on each segment
  const segs = [];
  for (let i = 1; i < stops.length; i++) {
    const a = [stops[i - 1].lng, stops[i - 1].lat], b = [stops[i].lng, stops[i].lat];
    if (a[0] === b[0] && a[1] === b[1]) continue;
    segs.push({ type: 'Feature', geometry: { type: 'LineString', coordinates: arc(a, b, 0.12) }, properties: { y: stops[i].y } });
  }
  map.getSource('route').setData({ type: 'FeatureCollection', features: segs });
  map.setFilter('route', ['<=', ['get', 'y'], year]);
  map.getSource('routepts').setData({ type: 'FeatureCollection', features: stops.map((x) => ({ type: 'Feature', geometry: { type: 'Point', coordinates: [x.lng, x.lat] }, properties: { y: x.y, n: x.n } })) });
  $('pcur').innerHTML = `<div class="mb-chip" style="color:#7df9ff;border-color:#7df9ff">Following ${esc(j.title)} · ${stops.length} stops <a href="#" id="pclear">✕</a></div> <a class="mb-chip" href="/journeys?person=${id}">Open journey →</a>`;
  $('pclear').onclick = (e) => { e.preventDefault(); personInfo = null; $('pcur').innerHTML = ''; map.getSource('route').setData(empty); map.getSource('routepts').setData(empty); drawTl(); };
  if (p) { setS(yearToS(Math.max(T0, p.y0))); }
  if (stops.length) map.fitBounds(stops.reduce((b, x) => b.extend([x.lng, x.lat]), new maplibregl.LngLatBounds([stops[0].lng, stops[0].lat], [stops[0].lng, stops[0].lat])), { padding: 120, maxZoom: 7, pitch: 50, duration: 2200 });
  drawTl();
}

// ---------------------------------------------------------------- driver
function setS(v) { s = Math.min(1, Math.max(0, v)); year = sToYear(s); applyYear(); }
function applyYear(force) {
  const e = eraAt(year);
  const yb = $('yearbox');
  yb.querySelector('.yr').textContent = fmtYear(year);
  const era = yb.querySelector('.era'); era.textContent = e.n; era.style.color = e.c;
  const idx = latestEventIdx(year);
  yb.querySelector('.cnt').textContent = `${idx + 1} of ${events.length} events · ${placesSoFar(year)} places`;
  if (ready) {
    map.setFilter('place-dots', ['<=', ['get', 'y0'], year]);
    for (const [id, vcMin] of [['l1', 60], ['l2', 12], ['l3', 0]]) map.setFilter('place-' + id, ['all', ['<=', ['get', 'y0'], year], ['>=', ['get', 'vc'], vcMin]]);
    map.setFilter('event-dots', ['<=', ['get', 'y'], year]);
    map.setFilter('chains', ['<=', ['get', 'y'], year]);
    if (personInfo) { map.setFilter('route', ['<=', ['get', 'y'], year]); map.setFilter('routepts', ['<=', ['get', 'y'], year]); }
  }
  if (idx !== lastIdx || force) {
    const forward = idx > lastIdx; lastIdx = idx;
    renderFeed(idx);
    const ev = events[idx];
    if (ev && ready && follow && forward && performance.now() - lastFly > 2400) flyToEvent(ev, false);
  }
  drawTl();
}

function pulses() {
  // recent events glow outwards for ~25 "story-years"
  const feats = [];
  const t = performance.now() / 1000;
  for (let i = latestEventIdx(year); i >= 0 && i > latestEventIdx(year) - 6; i--) {
    const e = events[i]; const age = year - e.y; if (age > 30) break;
    const k = Math.max(0, 1 - age / 30);
    for (const id of e.pl.slice(0, 2)) {
      const p = placeById.get(id); if (!p) continue;
      const ph = (t * 0.8 + i * 0.37) % 1;
      feats.push({ type: 'Feature', geometry: { type: 'Point', coordinates: [p.lng, p.lat] }, properties: { r: 6 + ph * 26 * (0.4 + k), o: 0.4 * k * (1 - ph) } });
    }
  }
  map.getSource('pulse').setData({ type: 'FeatureCollection', features: feats });
}

let last = performance.now();
function loop(now) {
  const dt = Math.min(0.1, (now - last) / 1000); last = now;
  if (playing && !dragging) {
    setS(s + dt * speed / RUN_SECONDS);
    if (s >= 1) setPlaying(false);
  }
  if (ready) pulses();
  requestAnimationFrame(loop);
}

// initial UI state does not wait for the map
(function initUi() {
  const y = qs.get('year'), eid = qs.get('event'), pid = qs.get('place');
  if (y) s = yearToS(+y);
  else if (eid && evById.get(eid)) s = yearToS(evById.get(eid).y + 0.05);
  else if (pid && placeById.get(pid)) s = yearToS(Math.max(T0, placeById.get(pid).y0) + 5);
  else s = 0.001;
  year = sToYear(s);
  applyYear(true);
})();

function boot() {
  sizeTl();
  const y = qs.get('year'), eid = qs.get('event'), pid = qs.get('place'), person = qs.get('person');
  if (y) setS(yearToS(+y));
  if (eid && evById.get(eid)) {
    const e = evById.get(eid); setS(yearToS(e.y + 0.05)); flyToEvent(e, true);
    const c = centroid(e); if (c) setTimeout(() => showEventPopup(e, c), 1400);
  } else if (pid && placeById.get(pid)) {
    const p = placeById.get(pid); setS(yearToS(Math.max(T0, p.y0) + 5));
    map.jumpTo({ center: [p.lng, p.lat], zoom: 8, pitch: 55 });
  } else if (!y && !person) {
    setS(0.001);
  }
  if (person) followPerson(person);
  if (qs.get('play')) { setS(0); setPlaying(true); }
}
requestAnimationFrame(loop);
// the canvas needs layout before the first size
requestAnimationFrame(sizeTl);
void startCenter;
