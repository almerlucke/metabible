// Who was alive? A zoomable timeline of the Bible's most prominent people, with a year cursor.
const $ = (id) => document.getElementById(id);
const esc = (s) => String(s ?? '').replace(/[&<>"']/g, (c) => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' }[c]));
const fmtYear = (y) => window.mbYear(y);
const qs = new URLSearchParams(location.search);
import { copyLink, setQuery } from './share.js';

const [pd, world] = await Promise.all([fetch('/api/atlas/people').then((r) => r.json()), fetch('/api/atlas/world').then((r) => r.json())]);
const eras = world.eras, groups = pd.groups;
// span: known life (birth–death) preferred, otherwise first–last year they appear in the text
const all = pd.people.filter((p) => p.y0 || p.y1 || p.b != null).map((p) => {
  const a = p.b != null ? p.b : p.y0;
  let b = p.d != null ? p.d : Math.max(p.y1, a + 1);
  if (p.b != null && p.d == null) b = Math.min(b, a + 120);
  let lo = Math.min(a, b), hi = Math.max(a, b);
  const exact = p.b != null && p.d != null;
  if (!exact && hi - lo < 70) { const m = (lo + hi) / 2; lo = m - 35; hi = m + 35; } // typical lifetime around where they appear
  return { ...p, a: lo, b: hi, exact };
}).filter((p) => p.exact || p.b - p.a <= 140); // wide spans are tribes, nations or deities, not lives
const TOP = 420;
const shown = [...all].sort((x, y) => y.vc - x.vc).slice(0, TOP).sort((x, y) => x.a - y.a);
const byId = new Map(all.map((p) => [p.id, p]));
const groupColor = (g) => (g ? `hsl(${(g * 47) % 360} 62% 58%)` : '#9aa3c7');

// lane packing
function packLanes(list) {
  const ends = [];
  for (const p of list) {
    const start = p.a, end = Math.max(p.b, p.a + 90);
    let l = ends.findIndex((e) => e < start - 15);
    if (l < 0) { l = ends.length; ends.push(0); }
    ends[l] = end; p.lane = l;
  }
  return ends.length;
}
let lanes = packLanes(shown);

$('lifespans').innerHTML = `<canvas id="lc"></canvas>
  <div class="mb-panel" id="lpanel"></div>
  <button class="mb-btn" id="lmin" type="button" aria-label="${tx('js.lifespans.hide.the.panel')}" title="${tx('js.lifespans.hide.or.show.the')}">‹</button>
  <div class="ctip" id="ltip"></div>`;
const cv = $('lc'), ctx = cv.getContext('2d');
// the chart runs a little past AD 100 so people at the very end have room, and labels can never be cut off
const MIN_Y = -4100, MAX_Y = 260;
let W = 0, H = 0, vx0 = MIN_Y, vx1 = MAX_Y, year = qs.get('year') ? +qs.get('year') : -1000, sel = null, hover = null;
const TOPPAD = 34, BOT = 54, LANE = 11;
// the info panel sits over the left of the screen on wide displays: the timeline starts beside it, not under it
let collapsed = false;
const gut = () => (collapsed || innerWidth <= 900 ? 0 : 304);
const xOf = (y) => gut() + ((y - vx0) / (vx1 - vx0)) * (W - gut()), yearAt = (x) => vx0 + ((x - gut()) / (W - gut())) * (vx1 - vx0);
const alive = (p) => p.a <= year && p.b >= year;

function size() { const r = cv.parentElement.getBoundingClientRect(), d = devicePixelRatio || 1; W = r.width; H = r.height; cv.width = W * d; cv.height = H * d; cv.style.width = W + 'px'; cv.style.height = H + 'px'; ctx.setTransform(d, 0, 0, d, 0, 0); draw(); }
function laneY(l) { const avail = H - TOPPAD - BOT, lh = Math.min(LANE + 4, avail / lanes); return TOPPAD + l * lh + lh / 2; }

function draw() {
  ctx.clearRect(0, 0, W, H);
  ctx.save(); ctx.beginPath(); ctx.rect(gut(), 0, W - gut(), H); ctx.clip(); // nothing is ever drawn under the panel's margin
  // era bands
  ctx.font = '600 11px Inter, sans-serif'; ctx.textBaseline = 'middle';
  for (const e of eras) {
    const x0 = Math.max(0, xOf(e.from)), x1 = Math.min(W, xOf(e.to)); if (x1 < 0 || x0 > W) continue;
    ctx.fillStyle = e.c + '14'; ctx.fillRect(x0, TOPPAD - 8, x1 - x0, H - TOPPAD - BOT + 12);
    ctx.fillStyle = e.c; ctx.fillRect(x0, H - BOT + 8, Math.max(0, x1 - x0 - 1), 3);
    if (x1 - x0 > 70) { ctx.fillStyle = '#d9d3bf'; ctx.fillText(e.n, x0 + 6, H - BOT + 24); }
  }
  // ruler
  const span = vx1 - vx0, step = span > 2500 ? 500 : span > 1000 ? 250 : span > 500 ? 100 : span > 200 ? 50 : 20;
  ctx.fillStyle = 'rgba(236,232,220,.55)'; ctx.font = '10.5px Inter, sans-serif'; ctx.strokeStyle = 'rgba(245,196,81,.10)'; ctx.lineWidth = 1;
  for (let y = Math.ceil(vx0 / step) * step; y <= vx1; y += step) { const x = Math.round(xOf(y)) + .5; ctx.beginPath(); ctx.moveTo(x, TOPPAD - 8); ctx.lineTo(x, H - BOT + 6); ctx.stroke(); ctx.fillText(y === 0 ? '1' : fmtYear(y).replace(' ', ''), x + 3, 14); }
  // contemporaries of selection
  const rel = sel ? (p) => p.a <= sel.b && p.b >= sel.a : null;
  if (sel) { ctx.fillStyle = 'rgba(125,249,255,.07)'; ctx.fillRect(xOf(sel.a), TOPPAD - 8, Math.max(2, xOf(sel.b) - xOf(sel.a)), H - TOPPAD - BOT + 12); }
  // bars
  for (const p of shown) {
    const x0 = xOf(p.a), x1 = Math.max(x0 + 4, xOf(p.b)); if (x1 < 0 || x0 > W) continue;
    const y = laneY(p.lane), isAlive = alive(p), isRel = rel ? rel(p) : isAlive, isSel = sel && sel.id === p.id;
    ctx.globalAlpha = isSel ? 1 : isRel ? 0.95 : 0.22;
    ctx.fillStyle = isSel ? '#7df9ff' : groupColor(p.gr);
    const h = 6; ctx.beginPath(); ctx.roundRect(x0, y - h / 2, x1 - x0, h, 3); ctx.fill();
    if (!p.exact) { ctx.fillStyle = '#0a0d1a'; ctx.globalAlpha *= 0.5; for (let xx = x0 + 3; xx < x1 - 2; xx += 6) ctx.fillRect(xx, y - 1, 3, 2); }
    ctx.globalAlpha = isRel || isSel ? 1 : 0.35;
    if (x1 - x0 > 55 || isSel || hover === p) { ctx.fillStyle = isSel ? '#fff' : '#f3ecd6'; ctx.font = (isSel ? '600 ' : '') + '10.5px Inter, sans-serif';
      // put the label on the other side of the bar when it would run off the right edge
      const tw = ctx.measureText(p.n).width;
      if (x1 + 5 + tw > W - 6 && x0 - 5 - tw > gut() + 4) { ctx.textAlign = 'right'; ctx.fillText(p.n, x0 - 5, y); ctx.textAlign = 'left'; }
      else ctx.fillText(p.n, Math.min(x1 + 5, Math.max(4, W - 6 - tw)), y); }
    ctx.globalAlpha = 1;
  }
  // year cursor
  const cx = xOf(year);
  ctx.strokeStyle = '#ffe08a'; ctx.lineWidth = 2; ctx.beginPath(); ctx.moveTo(cx, TOPPAD - 14); ctx.lineTo(cx, H - BOT + 8); ctx.stroke();
  ctx.fillStyle = '#ffe08a'; ctx.beginPath(); ctx.moveTo(cx - 7, TOPPAD - 20); ctx.lineTo(cx + 7, TOPPAD - 20); ctx.lineTo(cx, TOPPAD - 8); ctx.fill();
  ctx.font = '600 12px Inter, sans-serif'; ctx.fillStyle = '#ffe08a'; ctx.textAlign = cx > W - 80 ? 'right' : 'left'; ctx.fillText(fmtYear(year), cx + (cx > W - 80 ? -10 : 10), TOPPAD - 18); ctx.textAlign = 'left';
  ctx.restore();
}

function panel() {
  setQuery({ year: Math.round(year), person: sel ? sel.id : '' });
  const p = $('lpanel');
  const aliveNow = all.filter(alive).sort((a, b) => b.vc - a.vc);
  const era = eras.find((e) => year >= e.from && year < e.to) || eras[eras.length - 1];
  let html = `<div class="mb-tag" style="color:${era.c}">${esc(era.n)}</div><div class="yr">${fmtYear(year)}</div>
    <div class="mb-muted" style="font-size:.8rem">${tx('js.lifespans.summary', { n: aliveNow.length, k: aliveNow.filter((x) => x.exact).length })}</div>`;
  if (sel) {
    const contemps = all.filter((x) => x.id !== sel.id && x.a <= sel.b && x.b >= sel.a).map((x) => ({ x, ov: Math.min(x.b, sel.b) - Math.max(x.a, sel.a) })).sort((m, n) => n.x.vc - m.x.vc).slice(0, 14);
    html += `<div class="sel"><div class="mb-tag" style="color:#7df9ff">${tx('js.lifespans.selected')}</div><h3><a href="/person/${sel.id}">${esc(sel.n)}</a></h3>
      <div class="mb-muted" style="font-size:.8rem">${fmtYear(sel.a)} – ${fmtYear(sel.b)}${sel.exact ? '' : ' ' + tx('js.lifespans.approx')}</div>
      <div class="links"><a href="/tree/${sel.id}">${tx('js.cosmos.family')}</a> · <a href="/journeys?person=${sel.id}">${tx('js.cosmos.journey')}</a> · <a href="/cosmos?id=${sel.id}">${tx('js.cosmos.cosmos')}</a> <a href="#" id="lclear" style="float:right">${tx('js.lifespans.clear')}</a></div>
      <div class="mb-tag" style="margin:.7rem 0 .3rem">${tx('js.lifespans.contemporaries')}</div>` + contemps.map(({ x }) => `<a class="mb-chip c" href="#" data-id="${x.id}">${esc(x.n)}</a>`).join('') + '</div>';
  } else {
    html += `<div class="mb-tag" style="margin:.8rem 0 .3rem">${tx('js.lifespans.most.prominent.then')}</div>` + aliveNow.slice(0, 16).map((x) => `<a class="mb-chip c" href="#" data-id="${x.id}" style="border-color:${groupColor(x.gr)}66">${esc(x.n)}</a>`).join('');
  }
  html += '<div style="margin-top:.9rem"><button class="mb-btn" id="lshare" style="width:100%">' + tx('js.compare.copy.link') + '</button></div>';
  p.innerHTML = html;
  $('lshare').onclick = () => copyLink(location.href, $('lshare'));
  p.querySelectorAll('a.c').forEach((a) => a.addEventListener('click', (e) => { e.preventDefault(); choose(byId.get(a.dataset.id)); }));
  const c = $('lclear'); if (c) c.onclick = (e) => { e.preventDefault(); sel = null; panel(); draw(); };
}
function choose(p) {
  if (!p) return; sel = p; year = Math.round((p.a + p.b) / 2);
  if (!shown.includes(p)) { shown.push(p); shown.sort((x, y) => x.a - y.a); lanes = packLanes(shown); }
  // keep the person in view
  if (p.a < vx0 || p.b > vx1) { const m = (p.a + p.b) / 2, half = Math.max(300, (p.b - p.a) * 2); vx0 = Math.max(MIN_Y, m - half); vx1 = Math.min(MAX_Y, m + half); }
  panel(); draw();
}

// ---------------------------------------------------------------- interaction
const tip = $('ltip');
function hit(mx, my) {
  let best = null, bd = 9;
  if (mx < gut()) return null;
  for (const p of shown) { const x0 = xOf(p.a), x1 = Math.max(x0 + 4, xOf(p.b)); if (mx < x0 - 3 || mx > x1 + 3) continue; const d = Math.abs(laneY(p.lane) - my); if (d < bd) { bd = d; best = p; } }
  return best;
}
let drag = null;
cv.addEventListener('pointerdown', (e) => { drag = { x: e.clientX, v0: vx0, v1: vx1, moved: false }; cv.setPointerCapture(e.pointerId); });
cv.addEventListener('pointermove', (e) => {
  const r = cv.getBoundingClientRect(), mx = e.clientX - r.left, my = e.clientY - r.top;
  if (drag) {
    const dx = e.clientX - drag.x; if (Math.abs(dx) > 3) drag.moved = true;
    if (drag.moved) { const dy = -(dx / (W - gut())) * (drag.v1 - drag.v0); let a = drag.v0 + dy, b = drag.v1 + dy; if (a < MIN_Y) { b += MIN_Y - a; a = MIN_Y; } if (b > MAX_Y) { a -= b - MAX_Y; b = MAX_Y; } vx0 = a; vx1 = b; draw(); return; }
  }
  const h = hit(mx, my); if (h !== hover) { hover = h; draw(); }
  cv.style.cursor = h ? 'pointer' : 'crosshair';
  if (h) { tip.style.display = 'block'; tip.style.transform = `translate(${mx + 14}px,${my + 14}px)`; tip.innerHTML = `<b>${esc(h.n)}</b><br><span class="mb-muted">${fmtYear(h.a)} – ${fmtYear(h.b)} · ${h.vc} verses</span>`; } else tip.style.display = 'none';
});
cv.addEventListener('pointerup', (e) => {
  const r = cv.getBoundingClientRect(), mx = e.clientX - r.left, my = e.clientY - r.top;
  if (drag && !drag.moved) { const h = hit(mx, my); if (h) choose(h); else { year = Math.round(yearAt(mx)); panel(); draw(); } }
  drag = null;
});
cv.addEventListener('pointerleave', () => { tip.style.display = 'none'; hover = null; });
cv.addEventListener('wheel', (e) => {
  e.preventDefault();
  const r = cv.getBoundingClientRect(), mx = e.clientX - r.left, anchor = yearAt(mx), k = Math.exp(e.deltaY * 0.0015);
  let a = anchor - (anchor - vx0) * k, b = anchor + (vx1 - anchor) * k;
  if (b - a > MAX_Y - MIN_Y) { a = MIN_Y; b = MAX_Y; } if (b - a < 40) return; a = Math.max(MIN_Y, a); b = Math.min(MAX_Y, b);
  vx0 = a; vx1 = b; draw();
}, { passive: false });
addEventListener('keydown', (e) => { if (e.key === 'ArrowLeft') { year -= e.shiftKey ? 100 : 10; panel(); draw(); } else if (e.key === 'ArrowRight') { year += e.shiftKey ? 100 : 10; panel(); draw(); } });
addEventListener('resize', size);
$('lmin').onclick = () => { collapsed = !collapsed; $('lifespans').classList.toggle('min', collapsed); $('lmin').textContent = collapsed ? '›' : '‹'; $('lmin').setAttribute('aria-label', collapsed ? tx('js.lifespans.show.panel') : tx('js.lifespans.hide.the.panel')); draw(); };
size(); panel();
const pid = qs.get('person'); if (pid && byId.get(pid)) choose(byId.get(pid));
