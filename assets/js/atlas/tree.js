// Family tree: ancestors above, descendants below, click to refocus, find the path between any two people.
import { copyLink } from './share.js';
const $ = (id) => document.getElementById(id);
const esc = (s) => String(s ?? '').replace(/[&<>"']/g, (c) => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' }[c]));
const fmtYear = (y) => (y < 0 ? `${-y} BC` : y === 0 ? '1 BC' : `${y} AD`);
const SVGNS = 'http://www.w3.org/2000/svg';

const data = await fetch('/api/atlas/people').then((r) => r.json());
const people = data.people, groups = data.groups;
const byId = new Map(people.map((p) => [p.id, p]));
const groupColor = (g) => (g ? `hsl(${(g * 47) % 360} 62% 58%)` : '#8f96b2');

// adjacency: parents, children, partners (symmetrical)
const parents = new Map(), kids = new Map(), spouses = new Map();
const add = (m, k, v) => { if (!byId.has(k) || !byId.has(v)) return; if (!m.has(k)) m.set(k, new Set()); m.get(k).add(v); };
for (const p of people) {
  if (p.f) { add(parents, p.id, p.f); add(kids, p.f, p.id); }
  if (p.m) { add(parents, p.id, p.m); add(kids, p.m, p.id); }
  for (const c of p.c || []) { add(kids, p.id, c); add(parents, c, p.id); }
  for (const x of p.p || []) { add(spouses, p.id, x); add(spouses, x, p.id); }
}
const get = (m, id) => [...(m.get(id) || [])];

function bestByName(name) {
  return people.filter((p) => p.n === name).sort((a, b) => b.vc - a.vc)[0];
}
const startId = $('tree').dataset.id || (bestByName('Jesus Christ') || people[0]).id;

// ---------------------------------------------------------------- svg scaffolding
const root = $('tree');
root.innerHTML = '';
const svg = document.createElementNS(SVGNS, 'svg');
svg.setAttribute('width', '100%'); svg.setAttribute('height', '100%'); svg.style.cursor = 'grab'; svg.style.touchAction = 'none';
const defs = document.createElementNS(SVGNS, 'defs');
defs.innerHTML = '<filter id="glow" x="-50%" y="-50%" width="200%" height="200%"><feGaussianBlur stdDeviation="5" result="b"/><feMerge><feMergeNode in="b"/><feMergeNode in="SourceGraphic"/></feMerge></filter>';
svg.appendChild(defs);
const vp = document.createElementNS(SVGNS, 'g'); svg.appendChild(vp);
root.appendChild(svg);
let tx = 0, ty = 0, sc = 1;
const applyT = () => vp.setAttribute('transform', `translate(${tx} ${ty}) scale(${sc})`);

let drag = null;
svg.addEventListener('pointerdown', (e) => { if (e.target.closest('.node')) return; drag = { x: e.clientX, y: e.clientY, tx, ty }; svg.setPointerCapture(e.pointerId); svg.style.cursor = 'grabbing'; });
svg.addEventListener('pointermove', (e) => { if (!drag) return; tx = drag.tx + e.clientX - drag.x; ty = drag.ty + e.clientY - drag.y; applyT(); });
svg.addEventListener('pointerup', () => { drag = null; svg.style.cursor = 'grab'; });
svg.addEventListener('wheel', (e) => {
  e.preventDefault();
  const r = svg.getBoundingClientRect(), mx = e.clientX - r.left, my = e.clientY - r.top;
  const k = Math.exp(-e.deltaY * 0.0015), ns = Math.min(2.5, Math.max(0.08, sc * k));
  tx = mx - (mx - tx) * (ns / sc); ty = my - (my - ty) * (ns / sc); sc = ns; applyT();
}, { passive: false });

const NW = 132, NH = 40, GX = 18, GY = 86;

// ---------------------------------------------------------------- layout
// Tidy layout: leaves get consecutive x, parents are centred above their children.
function layout(rootId, next, maxDepth, cap) {
  const nodes = [], edges = [];
  let leaf = 0;
  const seenPath = new Set();
  function rec(id, depth) {
    const n = { id, depth, x: 0, kids: [] };
    nodes.push(n);
    seenPath.add(id);
    const ks = depth < maxDepth && nodes.length < cap ? next(id).filter((k) => !seenPath.has(k)) : [];
    for (const k of ks) { if (nodes.length >= cap) break; const c = rec(k, depth + 1); n.kids.push(c); edges.push([n, c]); }
    seenPath.delete(id);
    if (n.kids.length) n.x = (n.kids[0].x + n.kids[n.kids.length - 1].x) / 2; else n.x = leaf++ * (NW + GX);
    return n;
  }
  const r = rec(rootId, 0);
  return { root: r, nodes, edges };
}

let focus = null, mode = 'tree', pathNodes = null;

function nodeEl(p, x, y, cls, delay) {
  const g = document.createElementNS(SVGNS, 'g');
  g.setAttribute('class', 'node ' + cls); g.setAttribute('transform', `translate(${x - NW / 2} ${y - NH / 2})`);
  g.style.animationDelay = delay + 'ms';
  const c = groupColor(p.gr);
  const yrs = p.b != null ? `${fmtYear(p.b)}${p.d != null ? ' – ' + fmtYear(p.d) : ''}` : (p.y0 || p.y1) ? `fl. ${fmtYear(p.y0)}` : '';
  g.innerHTML = `<rect width="${NW}" height="${NH}" rx="9" fill="#0f1426" stroke="${c}" stroke-width="${cls.includes('focus') ? 2.5 : 1.3}"/>
    <rect width="5" height="${NH}" rx="2" fill="${c}"/>
    <text x="14" y="17" font-size="12.5" fill="#f3ecd6" font-weight="600">${esc(p.n.length > 17 ? p.n.slice(0, 16) + '…' : p.n)}</text>
    <text x="14" y="31" font-size="9.5" fill="#9a9fb5">${esc(yrs)}${p.g === 'Female' ? ' ♀' : p.g === 'Male' ? ' ♂' : ''}</text>
    <title>${esc(p.n)}</title>`;
  g.setAttribute('tabindex', '0'); g.setAttribute('role', 'link'); g.setAttribute('aria-label', `${p.n}${yrs ? ', ' + yrs : ''}. Press Enter to focus the tree on them.`);
  g.addEventListener('click', () => { setFocus(p.id, true); });
  g.addEventListener('keydown', (e) => { if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); setFocus(p.id, true); } });
  return g;
}
function edgePath(ax, ay, bx, by) { const my = (ay + by) / 2; return `M${ax} ${ay} C${ax} ${my}, ${bx} ${my}, ${bx} ${by}`; }
function pathEl(d, color, cls, delay) {
  const e = document.createElementNS(SVGNS, 'path'); e.setAttribute('d', d); e.setAttribute('fill', 'none'); e.setAttribute('stroke', color); e.setAttribute('class', cls || 'edge');
  e.style.animationDelay = delay + 'ms'; return e;
}

function renderTree() {
  vp.innerHTML = '';
  const up = layout(focus, (id) => get(parents, id).sort((a, b) => (byId.get(b).vc - byId.get(a).vc)), 9, 260);
  const down = layout(focus, (id) => get(kids, id), 5, 420);
  // shift so both trees are centred on the same x for focus
  const ox = -down.root.x, oxu = -up.root.x;
  const bb = { x0: 1e9, x1: -1e9, y0: 1e9, y1: -1e9 };
  const grow = (x, y) => { bb.x0 = Math.min(bb.x0, x - NW / 2); bb.x1 = Math.max(bb.x1, x + NW / 2); bb.y0 = Math.min(bb.y0, y - NH / 2); bb.y1 = Math.max(bb.y1, y + NH / 2); };
  const eg = document.createElementNS(SVGNS, 'g'), ng = document.createElementNS(SVGNS, 'g'); vp.append(eg, ng);
  for (const [a, b] of up.edges) { eg.appendChild(pathEl(edgePath(a.x + oxu, -a.depth * GY - 0, b.x + oxu, -b.depth * GY), '#f5c45166', 'edge', b.depth * 90)); }
  for (const [a, b] of down.edges) { eg.appendChild(pathEl(edgePath(a.x + ox, a.depth * GY, b.x + ox, b.depth * GY), groupColor(byId.get(b.id).gr) + '99', 'edge', b.depth * 90)); }
  for (const n of up.nodes) { if (n.depth === 0) continue; const x = n.x + oxu, y = -n.depth * GY; grow(x, y); ng.appendChild(nodeEl(byId.get(n.id), x, y, 'anc', n.depth * 90)); }
  for (const n of down.nodes) { if (n.depth === 0) continue; const x = n.x + ox, y = n.depth * GY; grow(x, y); ng.appendChild(nodeEl(byId.get(n.id), x, y, 'desc', n.depth * 90)); }
  // focus + spouses
  grow(0, 0); ng.appendChild(nodeEl(byId.get(focus), 0, 0, 'focus', 0));
  const sp = get(spouses, focus).slice(0, 4);
  const side = [(NW + 36), -(NW + 36), 2 * (NW + 36), -2 * (NW + 36)];
  sp.forEach((id, i) => { const x = side[i]; grow(x, 0); eg.appendChild(pathEl(`M0 0 L${x} 0`, '#e0644f99', 'edge sp', 0)); ng.appendChild(nodeEl(byId.get(id), x, 0, 'spouse', 150)); });
  fit(bb);
}
function fit(bb) {
  const r = svg.getBoundingClientRect(), w = r.width - (innerWidth > 900 ? 340 : 0), h = r.height;
  const bw = bb.x1 - bb.x0 + 80, bh = bb.y1 - bb.y0 + 80;
  sc = Math.min(1, Math.max(0.12, Math.min(w / bw, h / bh)));
  tx = (innerWidth > 900 ? 330 : 0) + w / 2 - ((bb.x0 + bb.x1) / 2) * sc; ty = h / 2 - ((bb.y0 + bb.y1) / 2) * sc; applyT();
}

// ---------------------------------------------------------------- path between two people
function findPath(a, b) {
  const prev = new Map([[a, null]]); const q = [a];
  while (q.length) {
    const id = q.shift(); if (id === b) break;
    const nb = [...get(parents, id).map((x) => [x, 'up']), ...get(kids, id).map((x) => [x, 'down']), ...get(spouses, id).map((x) => [x, 'sp'])];
    for (const [n, kind] of nb) if (!prev.has(n)) { prev.set(n, [id, kind]); q.push(n); }
  }
  if (!prev.has(b)) return null;
  const out = []; let cur = b;
  while (cur !== a) { const [pid, kind] = prev.get(cur); out.push({ id: cur, via: kind }); cur = pid; }
  out.push({ id: a, via: null }); return out.reverse();
}
function renderPath(path) {
  vp.innerHTML = '';
  const eg = document.createElementNS(SVGNS, 'g'), ng = document.createElementNS(SVGNS, 'g'); vp.append(eg, ng);
  const pts = [];
  const bb = { x0: 1e9, x1: -1e9, y0: 1e9, y1: -1e9 };
  if (path.length > 12) {
    // long chains wrap into a serpentine grid so they stay readable
    const perRow = Math.max(6, Math.round(Math.sqrt(path.length * 2.4)));
    path.forEach((s, i) => { const r = Math.floor(i / perRow), c = i % perRow; pts.push({ x: (r % 2 ? perRow - 1 - c : c) * (NW + 30), y: r * 84 }); });
  } else {
    let gen = 0, x = 0;
    path.forEach((s, i) => { if (i > 0) { x += NW + 34; gen += s.via === 'up' ? -1 : s.via === 'down' ? 1 : 0; } pts.push({ x, y: gen * 70 }); });
  }
  pts.forEach((p, i) => {
    bb.x0 = Math.min(bb.x0, p.x - NW / 2); bb.x1 = Math.max(bb.x1, p.x + NW / 2); bb.y0 = Math.min(bb.y0, p.y - NH / 2); bb.y1 = Math.max(bb.y1, p.y + NH / 2);
    if (i > 0) { const q = pts[i - 1]; const same = Math.abs(q.y - p.y) < 1 && path.length > 12; const dir = p.x >= q.x ? 1 : -1;
      const d = path.length > 12 ? (same ? `M${q.x + dir * NW / 2} ${q.y} L${p.x - dir * NW / 2} ${p.y}` : `M${q.x} ${q.y + NH / 2} C${q.x} ${q.y + 60}, ${p.x} ${p.y - 60}, ${p.x} ${p.y - NH / 2}`)
        : `M${q.x + NW / 2} ${q.y} C${(q.x + p.x) / 2} ${q.y}, ${(q.x + p.x) / 2} ${p.y}, ${p.x - NW / 2} ${p.y}`;
      eg.appendChild(pathEl(d, path[i].via === 'sp' ? '#e0644f' : '#f5c451', 'edge hot', i * 40)); }
    ng.appendChild(nodeEl(byId.get(path[i].id), p.x, p.y, i === 0 || i === path.length - 1 ? 'focus' : 'anc', i * 40));
  });
  fit(bb);
}
function describe(path) {
  const ups = path.filter((s) => s.via === 'up').length, downs = path.filter((s) => s.via === 'down').length, sp = path.filter((s) => s.via === 'sp').length;
  const parts = []; if (ups) parts.push(`${ups} generation${ups > 1 ? 's' : ''} up`); if (downs) parts.push(`${downs} down`); if (sp) parts.push(`${sp} marriage${sp > 1 ? 's' : ''}`);
  return `${path.length - 1} steps: ${parts.join(', ')}`;
}

// ---------------------------------------------------------------- panel
const panel = $('tpanel');
panel.style.cssText = 'left:16px;top:16px;bottom:16px;width:310px;padding:1rem';
let info = null;
async function setFocus(id, push) {
  focus = id; mode = 'tree'; pathNodes = null;
  if (push) history.pushState(null, '', '/tree/' + id);
  renderTree(); renderPanel();
  info = await fetch('/api/atlas/person/' + id).then((r) => r.json()).catch(() => null);
  if (focus === id) renderPanel();
}
function chips(ids) { return ids.map((i) => `<a class="mb-chip t" data-id="${i}" href="/tree/${i}">${esc(byId.get(i).n)}</a>`).join(''); }
function renderPanel() {
  const p = byId.get(focus);
  const fam = [['Parents', get(parents, focus)], ['Partners', get(spouses, focus)], ['Siblings', sibs(focus)], ['Children', get(kids, focus)]].filter(([, l]) => l.length);
  panel.innerHTML = `<h3>Family tree</h3>
    <input id="tsearch" class="mb-input" placeholder="Find a person…" autocomplete="off"><div id="tres"></div>
    <div class="focus"><div class="mb-tag" style="color:${groupColor(p.gr)}">${esc(p.gr ? groups[p.gr - 1] : 'person')}${p.g ? ' · ' + esc(p.g) : ''}</div>
      <h2>${esc(p.n)}</h2><div class="mb-muted" style="font-size:.78rem">${p.b != null ? 'Born ' + fmtYear(p.b) : ''}${p.d != null ? ' · died ' + fmtYear(p.d) : ''}${p.b == null && (p.y0 || p.y1) ? 'In the text ' + fmtYear(p.y0) + ' – ' + fmtYear(p.y1) : ''} · ${p.vc} verses</div>
      ${info && info.text ? `<p class="bio">${esc(info.text)}</p>` : ''}
      <div class="links"><a href="/person/${p.id}">Full page</a> · <a href="/journeys?person=${p.id}">Journey</a> · <a href="/lifespans?person=${p.id}">Lifespans</a></div></div>
    ${fam.map(([t, l]) => `<div class="mb-tag" style="margin:.7rem 0 .25rem">${t}</div>${chips(l.slice(0, 24))}${l.length > 24 ? `<span class="mb-muted"> +${l.length - 24}</span>` : ''}`).join('')}
    <div class="mb-tag" style="margin:1rem 0 .3rem">Path to another person</div>
    <input id="psearch" class="mb-input" placeholder="Relate ${esc(p.n)} to…" autocomplete="off"><div id="pres"></div>
    <div class="presets"><button class="mb-btn" data-a="Adam" data-b="Jesus Christ">Adam → Jesus</button><button class="mb-btn" data-a="Abraham" data-b="David">Abraham → David</button><button class="mb-btn" data-a="Noah" data-b="Moses">Noah → Moses</button></div>
    <div id="pathinfo"></div>`;
  panel.querySelectorAll('a.t').forEach((a) => a.addEventListener('click', (e) => { e.preventDefault(); setFocus(a.dataset.id, true); }));
  hookSearch($('tsearch'), $('tres'), (id) => setFocus(id, true));
  hookSearch($('psearch'), $('pres'), (id) => showPath(focus, id));
  panel.querySelectorAll('.presets button').forEach((b) => b.onclick = () => { const A = bestByName(b.dataset.a), B = bestByName(b.dataset.b); if (A && B) showPath(A.id, B.id); });
}
function sibs(id) { const s = new Set(); for (const par of get(parents, id)) for (const k of get(kids, par)) if (k !== id) s.add(k); return [...s]; }
function showPath(a, b) {
  const path = findPath(a, b);
  const el = $('pathinfo');
  if (!path) { el.innerHTML = '<p class="mb-muted">No family connection recorded between them.</p>'; return; }
  mode = 'path'; pathNodes = path; renderPath(path);
  history.replaceState(null, '', `/tree/${a}?to=${b}`);
  el.innerHTML = `<p style="margin:.6rem 0 .3rem"><b>${esc(byId.get(a).n)}</b> → <b>${esc(byId.get(b).n)}</b><br><span class="mb-muted">${describe(path)}</span></p>
    <ol class="plist">${path.map((s) => `<li><a class="t" data-id="${s.id}" href="/tree/${s.id}">${esc(byId.get(s.id).n)}</a>${s.via ? `<small> ${s.via === 'up' ? '(parent of previous)' : s.via === 'down' ? '(child of previous)' : '(married)'}</small>` : ''}</li>`).join('')}</ol>
    <button class="mb-btn" id="pback">Back to tree</button> <button class="mb-btn" id="pshare">🔗 Copy link</button>`;
  el.querySelectorAll('a.t').forEach((x) => x.addEventListener('click', (e) => { e.preventDefault(); setFocus(x.dataset.id, true); }));
  $('pback').onclick = () => setFocus(focus, false);
  $('pshare').onclick = () => copyLink(location.href, $('pshare'));
}
function hookSearch(input, box, pick) {
  let t;
  input.addEventListener('input', () => {
    clearTimeout(t); const q = input.value.trim().toLowerCase();
    if (q.length < 2) { box.innerHTML = ''; return; }
    t = setTimeout(() => {
      const hits = people.filter((p) => p.n.toLowerCase().includes(q)).sort((a, b) => (b.n.toLowerCase().startsWith(q) - a.n.toLowerCase().startsWith(q)) || b.vc - a.vc).slice(0, 7);
      box.innerHTML = hits.map((h) => `<div class="jr" data-id="${h.id}">${esc(h.n)} <span class="mb-muted">${h.b != null ? fmtYear(h.b) : h.y0 ? fmtYear(h.y0) : ''}</span></div>`).join('');
      box.querySelectorAll('.jr').forEach((d) => d.onclick = () => { box.innerHTML = ''; input.value = ''; pick(d.dataset.id); });
    }, 100);
  });
}
addEventListener('popstate', () => { const m = location.pathname.match(/\/tree\/(.+)$/); if (m && byId.has(m[1])) setFocus(m[1], false); });
addEventListener('resize', () => { if (mode === 'tree') renderTree(); });

const first = byId.has(startId) ? startId : people[0].id;
await setFocus(first, false);
const toId = new URLSearchParams(location.search).get('to');
if (toId && byId.has(toId)) showPath(first, toId);
