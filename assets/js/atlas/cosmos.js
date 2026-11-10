// Cosmos: every person and place as a star. Height is time (4004 BC at the bottom, AD 100 at the top);
// the horizontal plane is a force layout of who appears with whom.
import * as THREE from 'three';
import { OrbitControls } from 'three/addons/controls/OrbitControls.js';

import { copyLink } from './share.js';
const $ = (id) => document.getElementById(id);
const esc = (s) => String(s ?? '').replace(/[&<>"']/g, (c) => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' }[c]));
const fmtYear = (y) => window.mbYear(y);

const ERAS = [
  [tx('js.era.short.primeval'), -4100, -2100, '#9575cd'], [tx('js.era.short.patriarchs'), -2100, -1700, '#26a69a'], [tx('js.era.short.exodus'), -1700, -1400, '#ef6c00'],
  [tx('js.era.short.judges'), -1400, -1050, '#e5533d'], [tx('js.era.short.kingdom'), -1050, -930, '#f5c451'], [tx('js.era.short.divided'), -930, -586, '#7986cb'],
  [tx('js.era.short.exile'), -586, -400, '#a1887f'], [tx('js.era.short.between'), -400, -5, '#78909c'], [tx('js.era.short.jesus'), -5, 33, '#ffe082'], [tx('js.era.short.church'), 33, 100, '#4fc3f7'],
];
const eraColor = (y) => (ERAS.find((e) => y >= e[1] && y < e[2]) || ERAS[ERAS.length - 1])[3];

$('cpanel').innerHTML = '<div class="mb-tag">' + tx('js.cosmos.loading.the.heavens') + '</div>';
const c = await fetch('/api/atlas/cosmos').then((r) => r.json());
const N = c.ids.length;

// ---------------------------------------------------------------- scene
const host = $('cosmos');
const renderer = new THREE.WebGLRenderer({ antialias: true });
// weaker devices get a lighter scene: no connection web, fewer background stars, no retina scaling
const weak = (navigator.hardwareConcurrency || 8) <= 4 || (navigator.deviceMemory || 8) <= 4 || matchMedia('(prefers-reduced-motion: reduce)').matches;
renderer.setPixelRatio(weak ? 1 : Math.min(devicePixelRatio, 2));
host.appendChild(renderer.domElement);
const scene = new THREE.Scene();
scene.background = new THREE.Color(0x05060f);
scene.fog = new THREE.FogExp2(0x05060f, 0.0032);
const camera = new THREE.PerspectiveCamera(55, 1, 0.5, 2000);
camera.position.set(190, 10, 210);
const controls = new OrbitControls(camera, renderer.domElement);
controls.target.set(0, -40, 0); controls.enableDamping = true; controls.dampingFactor = 0.06;
controls.autoRotate = !weak; controls.autoRotateSpeed = 0.35; controls.maxDistance = 600; controls.minDistance = 8;
controls.addEventListener('start', () => { controls.autoRotate = false; $('c-rot').checked = false; });

function resize() { const w = host.clientWidth, h = host.clientHeight; renderer.setSize(w, h); camera.aspect = w / h; camera.updateProjectionMatrix(); }
addEventListener('resize', resize); resize();

// background dust
{
  const g = new THREE.BufferGeometry(), n = weak ? 700 : 2500, p = new Float32Array(n * 3);
  for (let i = 0; i < n; i++) { const r = 300 + Math.random() * 600, t = Math.random() * 6.283, u = Math.acos(2 * Math.random() - 1); p[i * 3] = r * Math.sin(u) * Math.cos(t); p[i * 3 + 1] = r * Math.cos(u) * 0.8; p[i * 3 + 2] = r * Math.sin(u) * Math.sin(t); }
  g.setAttribute('position', new THREE.BufferAttribute(p, 3));
  scene.add(new THREE.Points(g, new THREE.PointsMaterial({ size: 1.2, color: 0x8892c8, transparent: true, opacity: 0.55, depthWrite: false, fog: false })));
}

// ---------------------------------------------------------------- time axis + era discs
const yOf = (yr) => (yr / 4100) * 90;
{
  const axis = new THREE.BufferGeometry().setFromPoints([new THREE.Vector3(0, yOf(-4100), 0), new THREE.Vector3(0, yOf(100), 0)]);
  scene.add(new THREE.Line(axis, new THREE.LineBasicMaterial({ color: 0xf5c451, transparent: true, opacity: 0.45 })));
  for (const [, from, , col] of ERAS) {
    const pts = []; for (let i = 0; i <= 96; i++) { const a = (i / 96) * Math.PI * 2; pts.push(new THREE.Vector3(Math.cos(a) * 118, yOf(from), Math.sin(a) * 118)); }
    scene.add(new THREE.Line(new THREE.BufferGeometry().setFromPoints(pts), new THREE.LineBasicMaterial({ color: new THREE.Color(col), transparent: true, opacity: 0.22 })));
  }
}

// ---------------------------------------------------------------- stars
const pos = new Float32Array(N * 3), col = new Float32Array(N * 3), size = new Float32Array(N), alpha = new Float32Array(N);
const baseSize = new Float32Array(N);
const maxG = c.groups.length;
let colorMode = 'era';
function colorFor(i) {
  if (colorMode === 'type') return new THREE.Color(c.t[i] === 0 ? '#f5c451' : '#6fb1ff');
  if (colorMode === 'tribe') return c.g[i] ? new THREE.Color(`hsl(${(c.g[i] * 47) % 360}, 70%, 62%)`) : new THREE.Color(c.t[i] === 1 ? '#4a6fa5' : '#6f7590');
  return new THREE.Color(eraColor(c.yr[i]));
}
for (let i = 0; i < N; i++) {
  pos[i * 3] = c.x[i]; pos[i * 3 + 1] = c.y[i]; pos[i * 3 + 2] = c.z[i];
  baseSize[i] = (c.t[i] === 1 ? 2.6 : 2.2) + c.w[i] * 1.6; size[i] = baseSize[i]; alpha[i] = 0.85;
}
function recolor() { for (let i = 0; i < N; i++) { const k = colorFor(i); col.set([k.r, k.g, k.b], i * 3); } geo.attributes.color.needsUpdate = true; }
const geo = new THREE.BufferGeometry();
geo.setAttribute('position', new THREE.BufferAttribute(pos, 3));
geo.setAttribute('color', new THREE.BufferAttribute(col, 3));
geo.setAttribute('size', new THREE.BufferAttribute(size, 1));
geo.setAttribute('alpha', new THREE.BufferAttribute(alpha, 1));
const mat = new THREE.ShaderMaterial({
  transparent: true, depthWrite: false, blending: THREE.AdditiveBlending,
  uniforms: { scale: { value: 520 } },
  vertexShader: `attribute float size; attribute float alpha; varying vec3 vC; varying float vA; uniform float scale;
    void main(){ vC = color; vA = alpha; vec4 mv = modelViewMatrix * vec4(position,1.0); gl_PointSize = size * scale / -mv.z; gl_Position = projectionMatrix * mv; }`,
  vertexColors: true,
  fragmentShader: `varying vec3 vC; varying float vA;
    void main(){ float d = length(gl_PointCoord - 0.5); float a = smoothstep(0.5, 0.05, d); float core = smoothstep(0.18, 0.0, d); gl_FragColor = vec4(vC + core*0.6, a * vA); }`,
});
const points = new THREE.Points(geo, mat);
points.frustumCulled = false;
scene.add(points);
recolor();

// edges: faint all-over web plus highlight lines for the selection
const E = c.e; const adj = Array.from({ length: N }, () => []);
for (let k = 0; k < E.length; k += 3) { adj[E[k]].push([E[k + 1], E[k + 2]]); adj[E[k + 1]].push([E[k], E[k + 2]]); }
const ePos = new Float32Array((E.length / 3) * 6), eCol = new Float32Array((E.length / 3) * 6);
for (let k = 0, j = 0; k < E.length; k += 3, j += 6) {
  const a = E[k], b = E[k + 1], kind = E[k + 2];
  ePos.set([c.x[a], c.y[a], c.z[a], c.x[b], c.y[b], c.z[b]], j);
  const cc = kind === 1 ? [0.95, 0.77, 0.3] : kind === 2 ? [0.3, 0.5, 0.9] : [0.55, 0.55, 0.75];
  eCol.set([...cc, ...cc], j);
}
const eGeo = new THREE.BufferGeometry();
eGeo.setAttribute('position', new THREE.BufferAttribute(ePos, 3)); eGeo.setAttribute('color', new THREE.BufferAttribute(eCol, 3));
const web = new THREE.LineSegments(eGeo, new THREE.LineBasicMaterial({ vertexColors: true, transparent: true, opacity: 0.07, blending: THREE.AdditiveBlending, depthWrite: false }));
scene.add(web);
const selLines = new THREE.LineSegments(new THREE.BufferGeometry(), new THREE.LineBasicMaterial({ color: 0xffe08a, transparent: true, opacity: 0.8, blending: THREE.AdditiveBlending, depthWrite: false }));
scene.add(selLines);
const ring = new THREE.Mesh(new THREE.RingGeometry(0.94, 1.0, 64), new THREE.MeshBasicMaterial({ color: 0xffe08a, side: THREE.DoubleSide, transparent: true, opacity: 0.9, depthWrite: false }));
ring.visible = false; scene.add(ring);

// ---------------------------------------------------------------- labels
const labelHost = document.createElement('div'); labelHost.className = 'clabels'; host.appendChild(labelHost);
const topIdx = [...Array(N).keys()].sort((a, b) => c.w[b] - c.w[a]).slice(0, 38);
const labels = topIdx.map((i) => { const d = document.createElement('div'); d.className = 'clabel'; d.textContent = c.names[i]; labelHost.appendChild(d); return { i, d }; });
const eraLabels = ERAS.map(([n, from, to, color]) => { const d = document.createElement('div'); d.className = 'clabel era'; d.textContent = n; d.style.color = color; labelHost.appendChild(d); return { d, y: yOf((from + to) / 2) }; });
const tmp = new THREE.Vector3();
function project(x, y, z, el, maxD) {
  tmp.set(x, y, z).project(camera);
  const dist = camera.position.distanceTo(new THREE.Vector3(x, y, z));
  if (tmp.z > 1 || dist > maxD) { el.style.display = 'none'; return; }
  el.style.display = 'block';
  el.style.transform = `translate(${(tmp.x * 0.5 + 0.5) * host.clientWidth}px, ${(-tmp.y * 0.5 + 0.5) * host.clientHeight}px) translate(-50%, -140%)`;
  el.style.opacity = Math.max(0.15, 1 - dist / maxD);
}

// ---------------------------------------------------------------- interaction
const ray = new THREE.Raycaster(); ray.params.Points.threshold = 1.6;
const mouse = new THREE.Vector2(); let hover = -1, sel = -1, moved = false, downAt = null;
renderer.domElement.addEventListener('pointermove', (e) => {
  const r = renderer.domElement.getBoundingClientRect(); mouse.set(((e.clientX - r.left) / r.width) * 2 - 1, -((e.clientY - r.top) / r.height) * 2 + 1);
  moved = true; tipPos = [e.clientX - r.left, e.clientY - r.top];
});
renderer.domElement.addEventListener('pointerdown', (e) => { downAt = [e.clientX, e.clientY]; });
renderer.domElement.addEventListener('pointerup', (e) => { if (downAt && Math.hypot(e.clientX - downAt[0], e.clientY - downAt[1]) < 5) { if (hover >= 0) select(hover, true); else if (sel >= 0) select(-1); } downAt = null; });
let tipPos = [0, 0];
const tip = document.createElement('div'); tip.className = 'ctip'; host.appendChild(tip);

function pick() {
  ray.setFromCamera(mouse, camera);
  const hits = ray.intersectObject(points);
  let best = -1, bd = 1e9;
  for (const h of hits) { if (alpha[h.index] < 0.2) continue; const d = h.distanceToRay * (1 - Math.min(0.7, c.w[h.index] / 8)); if (d < bd) { bd = d; best = h.index; } }
  return best;
}

function select(i, fly) {
  sel = i;
  history.replaceState(null, '', i >= 0 ? '/cosmos?id=' + c.ids[i] : '/cosmos');
  if (i < 0) { selLines.visible = false; ring.visible = false; $('cinfo').style.display = 'none'; applyFilters(); return; }
  const pts = []; for (const [j] of adj[i]) pts.push(c.x[i], c.y[i], c.z[i], c.x[j], c.y[j], c.z[j]);
  selLines.geometry.setAttribute('position', new THREE.BufferAttribute(new Float32Array(pts), 3)); selLines.visible = true;
  ring.visible = true; ring.position.set(c.x[i], c.y[i], c.z[i]);
  const nb = adj[i].map(([j]) => j).sort((a, b) => c.w[b] - c.w[a]);
  const url = (c.t[i] === 0 ? '/person/' : '/place/') + c.ids[i];
  $('cinfo').style.display = 'block';
  $('cinfo').innerHTML = `<div class="mb-tag">${c.t[i] === 0 ? (c.g[i] ? esc(c.groups[c.g[i] - 1]) : tx('js.person')) : tx('js.place')} · ${fmtYear(c.yr[i])}</div>
    <h3><a href="${url}">${esc(c.names[i])}</a></h3>
    <div class="links"><a href="${url}">${tx('js.cosmos.open.page')}</a>${c.t[i] === 0 ? ` · <a href="/tree/${c.ids[i]}">${tx('js.cosmos.family')}</a> · <a href="/journeys?person=${c.ids[i]}">${tx('js.cosmos.journey')}</a>` : ` · <a href="/atlas?place=${c.ids[i]}">${tx('js.cosmos.on.the.atlas')}</a>`}</div>
    <div class="mb-tag" style="margin:.7rem 0 .3rem">${nb.length} connections</div>
    <div>${nb.slice(0, 28).map((j) => `<a href="#" class="mb-chip nb" data-i="${j}">${esc(c.names[j])}</a>`).join('')}</div>`;
  $('cinfo').querySelectorAll('.nb').forEach((a) => a.addEventListener('click', (e) => { e.preventDefault(); select(+a.dataset.i, true); }));
  if (fly) flyTo(i);
  applyFilters();
}
let fly = null;
function flyTo(i) {
  const tgt = new THREE.Vector3(c.x[i], c.y[i], c.z[i]);
  const dir = camera.position.clone().sub(controls.target).normalize();
  fly = { t0: performance.now(), from: controls.target.clone(), to: tgt, cf: camera.position.clone(), ct: tgt.clone().add(dir.multiplyScalar(70)) };
  controls.autoRotate = false; $('c-rot').checked = false;
}

// ---------------------------------------------------------------- panel
const filt = { people: true, places: true, y0: -4100, y1: 100 };
$('cpanel').innerHTML = `<h3>${tx('js.cosmos.cosmos')}</h3>
  <div class="mb-muted" style="font-size:.78rem;margin-bottom:.6rem">${tx('js.cosmos.hint', { n: N.toLocaleString(mbLang), l: (E.length / 3).toLocaleString(mbLang) })}</div>
  <input id="csearch" class="mb-input" placeholder="${tx('js.cosmos.find.a.star')}" autocomplete="off"><div id="cres"></div>
  <div class="row" style="margin-top:.6rem"><label>${tx('js.cosmos.colour')}</label><select id="c-color" class="mb-input"><option value="era">${tx('js.cosmos.by.era')}</option><option value="tribe">${tx('js.cosmos.by.tribe')}</option><option value="type">${tx('js.cosmos.people.places')}</option></select></div>
  <label class="chk"><input type="checkbox" id="c-people" checked> ${tx('js.canon.people')}</label>
  <label class="chk"><input type="checkbox" id="c-places" checked> ${tx('js.canon.places')}</label>
  <label class="chk"><input type="checkbox" id="c-web" ${weak ? "" : "checked"}> ${tx('js.cosmos.connections')}</label>
  <label class="chk"><input type="checkbox" id="c-rot" ${weak ? "" : "checked"}> ${tx('js.cosmos.slow.rotation')}</label>
  <div style="margin-top:.5rem" class="mb-muted">${tx('js.cosmos.time.window')} <b id="c-range">${tx('js.cosmos.4100.bc.ad.100')}</b></div>
  <button class="mb-btn" id="cshare" style="width:100%;margin:.6rem 0">${tx('js.compare.copy.link')}</button>
  <input type="range" id="c-y0" min="-4100" max="100" value="-4100" step="10" style="width:100%"><input type="range" id="c-y1" min="-4100" max="100" value="100" step="10" style="width:100%">`;
$('cshare').onclick = () => copyLink(location.href, $('cshare'));
$('c-color').onchange = (e) => { colorMode = e.target.value; recolor(); };
$('c-people').onchange = (e) => { filt.people = e.target.checked; applyFilters(); };
$('c-places').onchange = (e) => { filt.places = e.target.checked; applyFilters(); };
web.visible = !weak;
$('c-web').onchange = (e) => { web.visible = e.target.checked; };
$('c-rot').onchange = (e) => { controls.autoRotate = e.target.checked; };
const onRange = () => { let a = +$('c-y0').value, b = +$('c-y1').value; if (a > b) [a, b] = [b, a]; filt.y0 = a; filt.y1 = b; $('c-range').textContent = `${fmtYear(a)} – ${fmtYear(b)}`; applyFilters(); };
$('c-y0').oninput = onRange; $('c-y1').oninput = onRange;
let st;
$('csearch').addEventListener('input', (e) => {
  clearTimeout(st); const q = e.target.value.trim().toLowerCase();
  if (q.length < 2) { $('cres').innerHTML = ''; return; }
  st = setTimeout(() => {
    const hits = []; for (let i = 0; i < N && hits.length < 400; i++) if (c.names[i].toLowerCase().includes(q)) hits.push(i);
    hits.sort((a, b) => c.w[b] - c.w[a]);
    $('cres').innerHTML = hits.slice(0, 7).map((i) => `<div class="jr" data-i="${i}">${esc(c.names[i])} <span class="mb-muted">${c.t[i] ? 'place' : 'person'}</span></div>`).join('');
    $('cres').querySelectorAll('.jr').forEach((d) => d.onclick = () => { $('cres').innerHTML = ''; $('csearch').value = ''; select(+d.dataset.i, true); });
  }, 100);
});

function applyFilters() {
  const nbSet = sel >= 0 ? new Set([sel, ...adj[sel].map((x) => x[0])]) : null;
  for (let i = 0; i < N; i++) {
    const vis = (c.t[i] === 0 ? filt.people : filt.places) && c.yr[i] >= filt.y0 && c.yr[i] <= filt.y1;
    let a = vis ? 0.85 : 0.0, s = baseSize[i];
    if (nbSet) { if (nbSet.has(i)) { a = vis ? 1 : 0; s *= i === sel ? 2.2 : 1.5; } else a *= 0.18; }
    alpha[i] = a; size[i] = s;
  }
  geo.attributes.alpha.needsUpdate = true; geo.attributes.size.needsUpdate = true;
}
applyFilters();

// ---------------------------------------------------------------- loop
const clock = new THREE.Clock();
function frame() {
  const dt = clock.getDelta();
  if (fly) {
    const k = Math.min(1, (performance.now() - fly.t0) / 1400), e = k < 0.5 ? 2 * k * k : 1 - Math.pow(-2 * k + 2, 2) / 2;
    controls.target.lerpVectors(fly.from, fly.to, e); camera.position.lerpVectors(fly.cf, fly.ct, e); if (k >= 1) fly = null;
  }
  controls.update();
  if (moved) {
    moved = false; const h = pick();
    if (h !== hover) { hover = h; renderer.domElement.style.cursor = h >= 0 ? 'pointer' : 'grab'; }
  }
  if (hover >= 0) { tip.style.display = 'block'; tip.style.transform = `translate(${tipPos[0] + 14}px, ${tipPos[1] + 14}px)`; tip.innerHTML = `<b>${esc(c.names[hover])}</b><br><span class="mb-muted">${c.t[hover] ? 'place' : 'person'} · ${fmtYear(c.yr[hover])}</span>`; } else tip.style.display = 'none';
  if (ring.visible) { ring.quaternion.copy(camera.quaternion); const d = camera.position.distanceTo(ring.position), s = Math.max(1.4, d * 0.03) * (1 + Math.sin(performance.now() / 300) * 0.06) + c.w[sel] * 0.3; ring.scale.setScalar(s); }
  const showLabels = $('c-labels') ? $('c-labels').checked : true;
  for (const l of labels) { if (!showLabels || alpha[l.i] < 0.2) l.d.style.display = 'none'; else project(c.x[l.i], c.y[l.i] + 1, c.z[l.i], l.d, 420); }
  for (const l of eraLabels) project(0, l.y, 0, l.d, 500);
  renderer.render(scene, camera);
  requestAnimationFrame(frame);
  void dt;
}
frame();
const q = new URLSearchParams(location.search).get('id');
if (q) { const i = c.ids.indexOf(q); if (i >= 0) select(i, true); }
