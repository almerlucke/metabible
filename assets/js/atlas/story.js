// Story mode: the Bible in order, one event at a time, with the map following along.
/* global maplibregl */
import { createMap, FONT, esc, arc } from './map.js';

const $ = (id) => document.getElementById(id);
const story = await fetch('/api/atlas/story').then((r) => r.json());
const { eras, steps } = story;
const eraOf = new Map(eras.map((e) => [e.id, e]));
const evStep = new Map(steps.map((s, i) => [s.id, i]).filter(([id]) => id));
const qs = new URLSearchParams(location.search);

let idx = Math.min(steps.length - 1, Math.max(0, +(qs.get('s') || 0)));
let playing = qs.has('play'), speed = 1, narrate = false, token = 0, timer = null, ready = false;

const map = createMap('map', { center: [36, 33], zoom: 4.6, pitch: 45, bearing: -5 });
const empty = { type: 'FeatureCollection', features: [] };
map.on('error', (e) => console.log('maplibre error', e && e.error && e.error.message));
map.on('load', () => {
  for (const id of ['trail', 'now', 'link']) map.addSource(id, { type: 'geojson', data: empty });
  map.addLayer({ id: 'link', type: 'line', source: 'link', layout: { 'line-cap': 'round' }, paint: { 'line-color': '#f5c451', 'line-width': 2.2, 'line-opacity': 0.75, 'line-dasharray': [2, 1.5] } });
  map.addLayer({ id: 'trail', type: 'circle', source: 'trail', paint: { 'circle-radius': 5, 'circle-color': ['get', 'color'], 'circle-opacity': ['get', 'o'], 'circle-stroke-color': '#0a0d1a', 'circle-stroke-width': 1 } });
  map.addLayer({ id: 'glow', type: 'circle', source: 'now', paint: { 'circle-radius': 30, 'circle-color': ['get', 'color'], 'circle-opacity': 0.3, 'circle-blur': 0.8 } });
  map.addLayer({ id: 'now', type: 'circle', source: 'now', paint: { 'circle-radius': 8, 'circle-color': '#fff', 'circle-stroke-color': ['get', 'color'], 'circle-stroke-width': 3 } });
  map.addLayer({ id: 'now-l', type: 'symbol', source: 'now', layout: { 'text-field': ['get', 'n'], 'text-font': FONT, 'text-size': 14, 'text-offset': [0, 1.4], 'text-anchor': 'top', 'text-allow-overlap': true }, paint: { 'text-color': '#fff', 'text-halo-color': 'rgba(5,7,16,1)', 'text-halo-width': 2 } });
  ready = true; show(idx, false);
});

// ------------------------------------------------------------------ chrome
$('sera').innerHTML = `<h3>The story</h3>` + eras.map((e) => `<button class="era" data-i="${e.start}" data-id="${e.id}" style="--c:${e.c}"><i></i><span>${esc(e.n)}</span><small>${esc(`${e.from < 0 ? -e.from + ' BC' : e.from + ' AD'}`)}</small></button>`).join('');
$('sera').querySelectorAll('.era').forEach((b) => b.onclick = () => { stop(); go(+b.dataset.i); });
$('scard').innerHTML = '<div class="prog" id="sprog"></div><div id="sbody"></div><div class="ctl" id="sctl"></div>';
$('sprog').innerHTML = eras.map((e, k) => { const end = k + 1 < eras.length ? eras[k + 1].start : steps.length; return `<div class="seg" data-i="${e.start}" style="flex:${end - e.start};--c:${e.c}" title="${esc(e.n)}"><i></i></div>`; }).join('');
$('sprog').querySelectorAll('.seg').forEach((s) => s.onclick = (ev) => { const r = s.getBoundingClientRect(); const e0 = eras.findIndex((e) => e.start === +s.dataset.i); const end = e0 + 1 < eras.length ? eras[e0 + 1].start : steps.length; stop(); go(Math.floor(+s.dataset.i + ((ev.clientX - r.left) / r.width) * (end - +s.dataset.i))); });
const hasTTS = 'speechSynthesis' in window;
$('sctl').innerHTML = `<button class="mb-btn icon" id="sprev" title="Previous (←)">◀</button><button class="mb-btn icon" id="splay" title="Play / pause (space)">▶</button><button class="mb-btn icon" id="snext" title="Next (→)">▶|</button>
  <button class="mb-btn sp on" data-s="1">1×</button><button class="mb-btn sp" data-s="1.6">1.6×</button><button class="mb-btn sp" data-s="2.5">2.5×</button>
  ${hasTTS ? '<button class="mb-btn" id="snar" title="Read aloud">🔊 Read aloud</button>' : ''}<span class="grow"></span><button class="mb-btn" id="sshare" title="Copy a link to this moment">🔗 Share</button><span class="mb-muted" id="scount"></span>`;
$('sprev').onclick = () => { stop(); go(idx - 1); };
$('snext').onclick = () => { stop(); go(idx + 1); };
$('splay').onclick = () => (playing ? stop() : start());
document.querySelectorAll('.sp').forEach((b) => b.onclick = () => { speed = +b.dataset.s; document.querySelectorAll('.sp').forEach((x) => x.classList.toggle('on', x === b)); });
if (hasTTS) $('snar').onclick = () => { narrate = !narrate; $('snar').classList.toggle('on', narrate); if (!narrate) speechSynthesis.cancel(); else if (playing) { clearTimeout(timer); speak(); } };
$('sshare').onclick = async () => { const u = location.origin + '/story?s=' + idx; try { await navigator.clipboard.writeText(u); $('sshare').textContent = '✓ Copied'; } catch (e) { prompt('Copy this link', u); } setTimeout(() => { $('sshare').textContent = '🔗 Share'; }, 1800); };
addEventListener('keydown', (e) => {
  if (e.target.tagName === 'INPUT') return;
  if (e.key === 'ArrowRight') { stop(); go(idx + 1); } else if (e.key === 'ArrowLeft') { stop(); go(idx - 1); } else if (e.code === 'Space') { e.preventDefault(); playing ? stop() : start(); }
});

// ------------------------------------------------------------------ stepping
function go(i) { idx = Math.max(0, Math.min(steps.length - 1, i)); show(idx, true); }
function start() { playing = true; $('splay').textContent = '❚❚'; schedule(); }
function stop() { playing = false; token++; clearTimeout(timer); if (hasTTS) speechSynthesis.cancel(); const b = $('splay'); if (b) b.textContent = '▶'; }
function dwell(s) { return Math.min(11000, Math.max(3800, 2600 + ((s.tx || '').length + (s.bl || '').length * 0.3) * 38)) / speed; }
function schedule() {
  clearTimeout(timer); const my = ++token;
  if (!playing) return;
  if (narrate && hasTTS) { speak(); return; }
  timer = setTimeout(() => { if (my === token && playing) next(); }, dwell(steps[idx]));
}
function next() { if (idx >= steps.length - 1) { stop(); return; } go(idx + 1); }
function speak() {
  const s = steps[idx], my = token;
  speechSynthesis.cancel();
  const u = new SpeechSynthesisUtterance(`${s.t}. ${s.tx || ''}`); u.rate = 0.95 * Math.min(speed, 1.6);
  u.onend = () => { if (my === token && playing) timer = setTimeout(next, 700); };
  speechSynthesis.speak(u);
}

const colorOf = (s) => (eraOf.get(s.era) || eras[0]).c;
function centroid(s) { return s.pl && s.pl.length ? [s.pl.reduce((a, p) => a + p.lng, 0) / s.pl.length, s.pl.reduce((a, p) => a + p.lat, 0) / s.pl.length] : null; }

function show(i, animate) {
  const s = steps[i];
  history.replaceState(null, '', '?s=' + i);
  renderCard(s, i);
  $('sera').querySelectorAll('.era').forEach((b) => b.classList.toggle('on', b.dataset.id === s.era));
  // progress
  $('sprog').querySelectorAll('.seg').forEach((seg) => {
    const e0 = eras.findIndex((e) => e.start === +seg.dataset.i), start = eras[e0].start, end = e0 + 1 < eras.length ? eras[e0 + 1].start : steps.length;
    seg.firstElementChild.style.width = Math.max(0, Math.min(1, (i + 1 - start) / (end - start))) * 100 + '%';
  });
  if (ready) drawMap(s, i, animate);
  if (playing) schedule();
}

function drawMap(s, i, animate) {
  const color = colorOf(s);
  // fading trail of the last few places the story has visited
  const trail = [];
  for (let k = i - 1, n = 0; k >= 0 && n < 14; k--) {
    const t = steps[k]; if (!t.pl || !t.pl.length) continue;
    n++;
    for (const p of t.pl) trail.push({ type: 'Feature', geometry: { type: 'Point', coordinates: [p.lng, p.lat] }, properties: { color: colorOf(t), o: Math.max(0.12, 0.7 - n * 0.045) } });
  }
  map.getSource('trail').setData({ type: 'FeatureCollection', features: trail });
  const now = (s.pl || []).map((p) => ({ type: 'Feature', geometry: { type: 'Point', coordinates: [p.lng, p.lat] }, properties: { color, n: p.n } }));
  map.getSource('now').setData({ type: 'FeatureCollection', features: now });
  // cause-and-effect arc from the event this one follows
  const here = centroid(s), links = [];
  if (here) for (const pid of s.pr || []) { const pi = evStep.get(pid); const from = pi != null && centroid(steps[pi]); if (from && Math.hypot(from[0] - here[0], from[1] - here[1]) > 0.05) links.push({ type: 'Feature', geometry: { type: 'LineString', coordinates: arc(from, here, 0.15) }, properties: {} }); }
  map.getSource('link').setData({ type: 'FeatureCollection', features: links });

  const pad = innerWidth > 900 ? { top: 90, bottom: 290, left: 290, right: 60 } : { top: 80, bottom: 320, left: 20, right: 20 };
  const dur = animate ? 2200 : 0;
  if (s.pl && s.pl.length > 1) {
    const b = s.pl.reduce((bb, p) => bb.extend([p.lng, p.lat]), new maplibregl.LngLatBounds([s.pl[0].lng, s.pl[0].lat], [s.pl[0].lng, s.pl[0].lat]));
    map.fitBounds(b, { padding: pad, maxZoom: 7, pitch: 50, bearing: ((s.y * 7) % 30) - 15, duration: dur });
  } else if (s.pl && s.pl.length === 1) {
    map.easeTo({ center: [s.pl[0].lng, s.pl[0].lat], zoom: 6.3, pitch: 50, bearing: ((s.y * 7) % 30) - 15, duration: dur, padding: pad });
  } else if (s.kind === 'era' && animate) {
    map.easeTo({ center: [36, 33], zoom: 4.6, pitch: 45, bearing: -5, duration: dur, padding: { top: 0, bottom: 0, left: 0, right: 0 } });
  }
}

function renderCard(s, i) {
  const col = colorOf(s);
  let h = '';
  if (s.kind === 'era') {
    const e = eraOf.get(s.era);
    h = `<div class="mb-tag" style="color:${col}">Chapter ${eras.indexOf(e) + 1} of ${eras.length} · ${esc(s.dt)}</div><h2 class="eratitle">${esc(s.t)}</h2><p class="intro">${esc(s.tx)}</p>`;
  } else {
    h = `<div class="mb-tag" style="color:${col}">${esc(s.dt)} · ${esc(eraOf.get(s.era).n)}</div><h2><a href="/event/${s.id}">${esc(s.t)}</a></h2>` +
      (s.tx ? `<div class="mb-verse">“${esc(s.tx)}”</div>` : '') +
      (s.vid ? `<a class="ref" href="/verse/${s.vid}">${esc(s.ref)} →</a>` : '') +
      `<div class="chips">${(s.pl || []).map((p) => `<a class="mb-chip" href="/place/${p.id}">📍 ${esc(p.n)}</a>`).join('')}${(s.pe || []).map((p) => `<a class="mb-chip" href="/person/${p.id}">${esc(p.n)}</a>`).join('')}</div>` +
      (s.bl ? `<div class="know"><b>${esc(s.blo)}</b> — ${esc(s.bl)} <a href="${s.blu}">more →</a></div>` : '');
  }
  $('sbody').innerHTML = h;
  $('scount').textContent = `${i + 1} / ${steps.length}`;
}

if (playing) { playing = false; start(); }
