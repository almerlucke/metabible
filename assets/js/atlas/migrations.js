// Peoples & places: a heatmap of where the text is, century by century, and a trail for each tribe or group.
/* global maplibregl */
import { createMap, FONT, esc, arc } from './map.js';

const $ = (id) => document.getElementById(id);
const m = await fetch('/api/atlas/migration').then((r) => r.json());
const qs = new URLSearchParams(location.search);
const fmt = (y) => (y < 0 ? `${-y} BC` : y === 0 ? '1 BC' : `${y} AD`);
const bucketLabel = (b) => { const a = m.from + b * m.step; return `${fmt(a)} – ${fmt(a + m.step)}`; };

// weights per bucket
const byBucket = Array.from({ length: m.nb }, () => new Map());
for (const [pi, b, n] of m.counts) byBucket[b].set(pi, n);

let bucket = Math.min(m.nb - 1, Math.max(0, +(qs.get('b') ?? 26)));
let mode = qs.get('mode') === 'trail' ? 'trail' : 'heat', group = +(qs.get('g') || 0), playing = false, ready = false, last = performance.now();
if (!m.groups.find((g) => g.Index === group || g.g === group)) group = m.groups[0].g;

const map = createMap('map', { center: [36, 32.5], zoom: 5.1, pitch: 40, bearing: 0 });
const empty = { type: 'FeatureCollection', features: [] };

map.on('load', () => {
  for (const id of ['heat', 'tall', 'tdone', 'tpts', 'tnow']) map.addSource(id, { type: 'geojson', data: empty });
  map.addLayer({ id: 'heat', type: 'heatmap', source: 'heat', paint: {
    'heatmap-weight': ['interpolate', ['linear'], ['get', 'w'], 0, 0, 1, 0.12, 12, 0.5, 60, 1],
    'heatmap-intensity': ['interpolate', ['linear'], ['zoom'], 3, 0.8, 9, 2],
    'heatmap-radius': ['interpolate', ['linear'], ['zoom'], 3, 14, 6, 34, 9, 70],
    'heatmap-opacity': 0.85,
    'heatmap-color': ['interpolate', ['linear'], ['heatmap-density'], 0, 'rgba(0,0,0,0)', 0.15, 'rgba(60,40,150,.45)', 0.35, 'rgba(180,50,120,.6)', 0.6, 'rgba(245,140,60,.75)', 0.85, 'rgba(255,214,110,.9)', 1, 'rgba(255,250,210,.98)'] } });
  map.addLayer({ id: 'hp', type: 'circle', source: 'heat', minzoom: 4, filter: ['>=', ['get', 'w'], 3], paint: { 'circle-radius': ['interpolate', ['linear'], ['get', 'w'], 3, 3, 40, 8], 'circle-color': '#fff3b0', 'circle-opacity': 0.8, 'circle-stroke-color': '#2a1030', 'circle-stroke-width': 1 } });
  map.addLayer({ id: 'hl', type: 'symbol', source: 'heat', minzoom: 5, filter: ['>=', ['get', 'w'], 4], layout: { 'text-field': ['get', 'n'], 'text-font': FONT, 'text-size': 12, 'text-offset': [0, 1], 'text-anchor': 'top', 'text-optional': true, 'symbol-sort-key': ['-', 0, ['get', 'w']] }, paint: { 'text-color': '#fff6dc', 'text-halo-color': 'rgba(5,7,16,.95)', 'text-halo-width': 1.5 } });
  map.addLayer({ id: 'tall', type: 'line', source: 'tall', layout: { 'line-cap': 'round' }, paint: { 'line-color': '#7df9ff', 'line-width': 2, 'line-opacity': 0.25, 'line-dasharray': [1.5, 2] } });
  map.addLayer({ id: 'tdone', type: 'line', source: 'tdone', layout: { 'line-cap': 'round', 'line-join': 'round' }, paint: { 'line-color': '#7df9ff', 'line-width': 3.5, 'line-opacity': 0.9 } });
  map.addLayer({ id: 'tpts', type: 'circle', source: 'tpts', paint: { 'circle-radius': ['interpolate', ['linear'], ['get', 'w'], 3, 4, 80, 9], 'circle-color': ['case', ['get', 'past'], '#7df9ff', '#10142a'], 'circle-stroke-color': '#7df9ff', 'circle-stroke-width': 1.6 } });
  map.addLayer({ id: 'tglow', type: 'circle', source: 'tnow', paint: { 'circle-radius': 28, 'circle-color': '#7df9ff', 'circle-opacity': 0.3, 'circle-blur': 0.8 } });
  map.addLayer({ id: 'tnow', type: 'circle', source: 'tnow', paint: { 'circle-radius': 9, 'circle-color': '#fff', 'circle-stroke-color': '#7df9ff', 'circle-stroke-width': 3 } });
  map.addLayer({ id: 'tlbl', type: 'symbol', source: 'tpts', layout: { 'text-field': ['get', 'lbl'], 'text-font': FONT, 'text-size': 11.5, 'text-offset': [0, 1.2], 'text-anchor': 'top', 'text-optional': true }, paint: { 'text-color': '#d8f7ff', 'text-halo-color': 'rgba(5,7,16,.95)', 'text-halo-width': 1.4 } });
  map.on('click', 'hp', (e) => { const f = e.features[0].properties; new maplibregl.Popup({ offset: 8 }).setLngLat(e.features[0].geometry.coordinates).setHTML(`<b>${esc(f.n)}</b><br><span class="mb-muted">${f.w} mentions in this window</span><br><a href="/place/${f.id}">Open place →</a>`).addTo(map); });
  map.on('mouseenter', 'hp', () => { map.getCanvas().style.cursor = 'pointer'; });
  map.on('mouseleave', 'hp', () => { map.getCanvas().style.cursor = ''; });
  ready = true; update(true);
});

// ------------------------------------------------------------------ panel
$('mpanel').innerHTML = `<h3>Peoples &amp; places</h3>
  <div class="seg2"><button class="mb-btn" data-m="heat">🔥 Where the text is</button><button class="mb-btn" data-m="trail">👣 Follow a tribe</button></div>
  <div class="mb-muted note" id="mnote"></div>
  <div class="when"><b id="mwhen"></b><div class="mb-muted" id="mera"></div></div>
  <div class="slide"><button class="mb-btn icon" id="mplay" title="Play">▶</button><input type="range" id="mslide" min="0" max="${m.nb - 1}" step="1" value="${bucket}"></div>
  <div class="ticks"><span>4100 BC</span><span>2000 BC</span><span>1 BC</span></div>
  <div id="mgroup" style="display:none"><div class="mb-tag" style="margin:.7rem 0 .3rem">Tribe or group</div><select id="msel" class="mb-input">${m.groups.map((g) => `<option value="${g.g}">${esc(g.n)}</option>`).join('')}</select></div>
  <div class="mb-tag" style="margin:.9rem 0 .3rem" id="mlisttitle">Most mentioned now</div><div id="mlist"></div>
  <div class="mb-muted foot">Counts come from the places named in verses dated to each century, so this shows where the text is, not census data.</div>`;
$('msel').value = group;
$('mpanel').querySelectorAll('.seg2 .mb-btn').forEach((b) => b.onclick = () => { mode = b.dataset.m; update(true); });
$('mslide').oninput = (e) => { bucket = +e.target.value; update(false); };
$('msel').onchange = (e) => { group = +e.target.value; update(true); };
$('mplay').onclick = () => { playing = !playing; $('mplay').textContent = playing ? '❚❚' : '▶'; last = performance.now(); acc = 0; };
let acc = 0;
function loop(now) {
  const dt = (now - last) / 1000; last = now;
  if (playing) { acc += dt; if (acc > 1.1) { acc = 0; bucket = bucket >= m.nb - 1 ? 0 : bucket + 1; $('mslide').value = bucket; update(false); } }
  requestAnimationFrame(loop);
}
requestAnimationFrame(loop);

// ------------------------------------------------------------------ drawing
const grp = () => m.groups.find((g) => g.g === group);
function heatData() {
  const w = new Map();
  for (const b of [bucket - 1, bucket, bucket + 1]) if (b >= 0 && b < m.nb) for (const [pi, n] of byBucket[b]) w.set(pi, (w.get(pi) || 0) + n * (b === bucket ? 1 : 0.5));
  const feats = [], top = [];
  for (const [pi, n] of w) { const p = m.places[pi]; feats.push({ type: 'Feature', geometry: { type: 'Point', coordinates: [p.lng, p.lat] }, properties: { id: p.id, n: p.n, w: Math.round(n) } }); top.push([p, Math.round(n)]); }
  top.sort((a, b) => b[1] - a[1]);
  return { fc: { type: 'FeatureCollection', features: feats }, top: top.slice(0, 9) };
}

function update(refit) {
  $('mwhen').textContent = bucketLabel(bucket);
  $('mpanel').querySelectorAll('.seg2 .mb-btn').forEach((b) => b.classList.toggle('on', b.dataset.m === mode));
  $('mgroup').style.display = mode === 'trail' ? 'block' : 'none';
  $('mnote').textContent = mode === 'heat' ? 'The glow marks where places are named in verses dated around this century (±100 years).' : 'The line joins the centre of the places named alongside members of this group, century by century.';
  history.replaceState(null, '', `?b=${bucket}&mode=${mode}&g=${group}`);
  const { fc, top } = heatData();
  if (!ready) return;
  map.getSource('heat').setData(mode === 'heat' ? fc : empty);
  for (const l of ['heat', 'hp', 'hl']) map.setLayoutProperty(l, 'visibility', mode === 'heat' ? 'visible' : 'none');
  for (const l of ['tall', 'tdone', 'tpts', 'tglow', 'tnow', 'tlbl']) map.setLayoutProperty(l, 'visibility', mode === 'trail' ? 'visible' : 'none');

  if (mode === 'heat') {
    $('mlisttitle').textContent = 'Most mentioned now';
    $('mlist').innerHTML = top.length ? top.map(([p, n]) => `<div class="mi" data-ll="${p.lng},${p.lat}"><span>${esc(p.n)}</span><small>${n}</small></div>`).join('') : '<div class="mb-muted">Almost no places are named in this century.</div>';
    $('mlist').querySelectorAll('.mi').forEach((el) => el.onclick = () => { const [x, y] = el.dataset.ll.split(',').map(Number); map.flyTo({ center: [x, y], zoom: 7.5, pitch: 50, duration: 1400 }); });
    return;
  }
  // trail
  const g = grp(), pts = g.trail.slice().sort((a, b) => a.b - b.b);
  const line = [];
  for (let i = 1; i < pts.length; i++) line.push(...arc([pts[i - 1].lng, pts[i - 1].lat], [pts[i].lng, pts[i].lat], 0.12, 24));
  map.getSource('tall').setData({ type: 'FeatureCollection', features: pts.length > 1 ? [{ type: 'Feature', geometry: { type: 'LineString', coordinates: line }, properties: {} }] : [] });
  const upTo = pts.filter((p) => p.b <= bucket);
  const done = [];
  for (let i = 1; i < upTo.length; i++) done.push(...arc([upTo[i - 1].lng, upTo[i - 1].lat], [upTo[i].lng, upTo[i].lat], 0.12, 24));
  map.getSource('tdone').setData({ type: 'FeatureCollection', features: done.length ? [{ type: 'Feature', geometry: { type: 'LineString', coordinates: done }, properties: {} }] : [] });
  map.getSource('tpts').setData({ type: 'FeatureCollection', features: pts.map((p) => ({ type: 'Feature', geometry: { type: 'Point', coordinates: [p.lng, p.lat] }, properties: { w: p.w, past: p.b <= bucket, lbl: `${fmt(m.from + p.b * m.step + 50)}` } })) });
  const cur = upTo[upTo.length - 1];
  map.getSource('tnow').setData(cur ? { type: 'FeatureCollection', features: [{ type: 'Feature', geometry: { type: 'Point', coordinates: [cur.lng, cur.lat] }, properties: {} }] } : empty);
  $('mlisttitle').textContent = g.n;
  $('mlist').innerHTML = pts.map((p) => `<div class="mi${p.b === (cur && cur.b) ? ' on' : ''}" data-b="${p.b}"><span>${fmt(m.from + p.b * m.step)} <small>${esc(p.top.join(', '))}</small></span><small>${p.w}</small></div>`).join('');
  $('mlist').querySelectorAll('.mi').forEach((el) => el.onclick = () => { bucket = +el.dataset.b; $('mslide').value = bucket; update(false); });
  if (refit && pts.length) map.fitBounds(pts.reduce((bb, p) => bb.extend([p.lng, p.lat]), new maplibregl.LngLatBounds([pts[0].lng, pts[0].lat], [pts[0].lng, pts[0].lat])), { padding: { top: 90, bottom: 90, left: 380, right: 80 }, maxZoom: 7, duration: 1600 });
  else if (cur && !refit) map.easeTo({ center: [cur.lng, cur.lat], duration: 700 });
}
