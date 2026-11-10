// Canon: all 66 books as a skyline of chapters, an arc diagram of shared people, and a word "tide".
const $ = (id) => document.getElementById(id);
const esc = (s) => String(s ?? '').replace(/[&<>"']/g, (c) => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' }[c]));

const data = await fetch('/api/atlas/canon').then((r) => r.json());
const books = data.books, arcs = data.arcs;
const divs = [...new Set(books.map((b) => b.div))];
const divColor = (d) => `hsl(${(divs.indexOf(d) * 360) / divs.length + 15} 65% 62%)`;
const maxCh = Math.max(...books.map((b) => b.ch.length));

const root = $('canon');
root.innerHTML = `<div class="canon-wrap">
  <div class="mb-tag">${tx('js.canon.the.whole.canon.at')}</div><h1 class="mb-display" style="margin:.2rem 0 .3rem;color:var(--mb-gold2)">${tx('js.canon.the.canon')}</h1>
  <p class="mb-muted" style="max-width:720px">${tx('js.canon.every.chapter.of.every')}</p>
  <div class="canon-controls">
    <div class="seg" id="lens"><button class="mb-btn on" data-m="pe">${tx('js.canon.people')}</button><button class="mb-btn" data-m="pl">${tx('js.canon.places')}</button><button class="mb-btn" data-m="ev">${tx('js.canon.events')}</button><button class="mb-btn" data-m="v">${tx('js.canon.verses')}</button><button class="mb-btn" data-m="w">${tx('js.canon.word')}</button></div>
    <input id="word" class="mb-input" placeholder="${tx('js.canon.type.a.word.love')}" style="max-width:280px">
    <span id="wordinfo" class="mb-muted"></span>
  </div>
  <div class="canon-legend" id="legend"></div>
  <div class="canvas-wrap"><canvas id="sky"></canvas></div>
  <h2 class="mb-display" style="margin:2.2rem 0 .3rem;color:var(--mb-gold2)">${tx('js.canon.who.connects.the.books')}</h2>
  <p class="mb-muted" style="max-width:720px">${tx('js.canon.arcs.join.books.that')}</p>
  <div class="canvas-wrap"><canvas id="arcs"></canvas></div><div id="arcinfo" class="mb-muted" style="min-height:2.6rem"></div>
</div><div class="ctip" id="tip" style="position:fixed"></div>`;
$('legend').innerHTML = divs.map((d) => `<span><i style="background:${divColor(d)}"></i>${esc(d)}</span>`).join('');

let lens = 'pe', wordCounts = null;
const tip = $('tip');

// ---------------------------------------------------------------- skyline heatmap
const sky = $('sky'), sctx = sky.getContext('2d');
const ROW = 15, LABEL = 120, PADT = 8;
let SW = 0, cell = 0, maxVal = 1;
const val = (b, bi, ci) => {
  const c = b.ch[ci];
  if (lens === 'w') return wordCounts ? wordCounts[bi][ci] : 0;
  return c[lens];
};
function drawSky() {
  const dpr = devicePixelRatio || 1, w = sky.parentElement.clientWidth;
  SW = w; cell = (w - LABEL - 10) / maxCh;
  sky.width = w * dpr; sky.height = (books.length * ROW + PADT * 2) * dpr; sky.style.height = books.length * ROW + PADT * 2 + 'px'; sctx.setTransform(dpr, 0, 0, dpr, 0, 0);
  maxVal = 1; books.forEach((b, bi) => b.ch.forEach((_, ci) => { maxVal = Math.max(maxVal, val(b, bi, ci)); }));
  sctx.clearRect(0, 0, w, sky.height);
  sctx.font = '11px Inter, sans-serif'; sctx.textBaseline = 'middle';
  books.forEach((b, bi) => {
    const y = PADT + bi * ROW;
    sctx.fillStyle = divColor(b.div); sctx.textAlign = 'right'; sctx.fillText(b.n, LABEL - 8, y + ROW / 2);
    b.ch.forEach((c, ci) => {
      const v = val(b, bi, ci) / maxVal, t = Math.pow(v, 0.6);
      sctx.fillStyle = lens === 'w' ? `rgba(255,105,180,${0.07 + t * 0.93})` : lens === 'pl' ? `rgba(111,177,255,${0.07 + t * 0.93})` : lens === 'ev' ? `rgba(240,110,80,${0.07 + t * 0.93})` : `rgba(245,196,81,${0.07 + t * 0.93})`;
      sctx.fillRect(LABEL + ci * cell, y + 1, Math.max(1, cell - 1), ROW - 2);
    });
  });
}
sky.addEventListener('mousemove', (e) => {
  const r = sky.getBoundingClientRect(), x = e.clientX - r.left - LABEL, y = e.clientY - r.top - PADT;
  const bi = Math.floor(y / ROW), ci = Math.floor(x / cell);
  const b = books[bi];
  if (!b || x < 0 || ci >= b.ch.length) { tip.style.display = 'none'; return; }
  const c = b.ch[ci];
  tip.style.display = 'block'; tip.style.left = e.clientX + 14 + 'px'; tip.style.top = e.clientY + 14 + 'px';
  tip.innerHTML = `<b>${esc(b.n)} ${c.n}</b><br><span class="mb-muted">${tx.n('js.verses', c.v)} · ${tx.n('js.people', c.pe)} · ${tx.n('js.places', c.pl)} · ${tx.n('js.events', c.ev)}${wordCounts ? ` · “${esc($('word').value)}” ×${wordCounts[bi][ci]}` : ''}</span>`;
  sky.style.cursor = 'pointer';
});
sky.addEventListener('mouseleave', () => { tip.style.display = 'none'; });
sky.addEventListener('click', (e) => {
  const r = sky.getBoundingClientRect(), x = e.clientX - r.left - LABEL, y = e.clientY - r.top - PADT;
  const b = books[Math.floor(y / ROW)], ci = Math.floor(x / cell);
  if (b && x >= 0 && b.ch[ci]) location.href = '/chapter/' + b.ch[ci].id;
  else if (b && x < 0) location.href = '/book/' + b.id;
});

$('lens').querySelectorAll('button').forEach((b) => b.onclick = () => {
  lens = b.dataset.m; $('lens').querySelectorAll('button').forEach((x) => x.classList.toggle('on', x === b));
  if (lens === 'w' && !wordCounts) $('word').focus();
  drawSky();
});
let wt;
$('word').addEventListener('input', () => {
  clearTimeout(wt);
  wt = setTimeout(async () => {
    const q = $('word').value.trim();
    if (q.length < 2) { wordCounts = null; $('wordinfo').textContent = ''; drawSky(); return; }
    const r = await fetch('/api/atlas/word?q=' + encodeURIComponent(q)).then((x) => x.json());
    wordCounts = r.books; lens = 'w'; $('lens').querySelectorAll('button').forEach((x) => x.classList.toggle('on', x.dataset.m === 'w'));
    // top chapter summary
    const flat = []; r.books.forEach((cs, bi) => cs.forEach((n, ci) => { if (n) flat.push([n, bi, ci]); }));
    flat.sort((a, b) => b[0] - a[0]);
    $('wordinfo').innerHTML = `<b>${r.total.toLocaleString()}</b> ${tx('js.canon.occpeak')} ` + flat.slice(0, 3).map(([n, bi, ci]) => `<a href="/chapter/${books[bi].ch[ci].id}">${esc(books[bi].n)} ${books[bi].ch[ci].n}</a> (${n})`).join(', ');
    drawSky();
  }, 220);
});

// ---------------------------------------------------------------- arcs
const ac = $('arcs'), actx = ac.getContext('2d');
let AW = 0, AH = 0, hov = -1;
const bx = (i) => 40 + (i * (AW - 80)) / (books.length - 1);
function drawArcs() {
  const dpr = devicePixelRatio || 1; AW = ac.parentElement.clientWidth; AH = Math.min(520, AW * 0.42) + 110;
  ac.width = AW * dpr; ac.height = AH * dpr; ac.style.height = AH + 'px'; actx.setTransform(dpr, 0, 0, dpr, 0, 0);
  actx.clearRect(0, 0, AW, AH);
  const base = AH - 100; const maxS = Math.max(...arcs.map((a) => a.s));
  for (const a of [...arcs].reverse()) {
    const x1 = bx(a.a), x2 = bx(a.b), r = (x2 - x1) / 2, on = hov < 0 || hov === a.a || hov === a.b;
    actx.beginPath(); actx.ellipse((x1 + x2) / 2, base, r, Math.min(r, base - 10) * 0.95, 0, Math.PI, 2 * Math.PI);
    const g = actx.createLinearGradient(x1, 0, x2, 0); g.addColorStop(0, divColor(books[a.a].div)); g.addColorStop(1, divColor(books[a.b].div));
    actx.strokeStyle = g; actx.globalAlpha = on ? (hov < 0 ? 0.12 + 0.5 * (a.s / maxS) : 0.85) : 0.025; actx.lineWidth = 0.6 + 3.2 * (a.s / maxS); actx.stroke();
  }
  actx.globalAlpha = 1;
  books.forEach((b, i) => {
    actx.fillStyle = divColor(b.div); actx.beginPath(); actx.arc(bx(i), base + 8, hov === i ? 6 : 3.4, 0, 7); actx.fill();
    actx.save(); actx.translate(bx(i) + 3, base + 18); actx.rotate(Math.PI / 2.6); actx.font = (hov === i ? '600 ' : '') + '10.5px Inter, sans-serif'; actx.fillStyle = hov === i ? '#fff' : divColor(b.div); actx.textAlign = 'left'; actx.fillText(b.n, 0, 0); actx.restore();
  });
}
ac.addEventListener('mousemove', (e) => {
  const r = ac.getBoundingClientRect(), x = e.clientX - r.left, y = e.clientY - r.top, base = AH - 100;
  let best = -1; if (y > base - 30) { let bd = 12; books.forEach((_, i) => { const d = Math.abs(bx(i) - x); if (d < bd) { bd = d; best = i; } }); }
  if (best !== hov) {
    hov = best; drawArcs();
    if (best < 0) $('arcinfo').innerHTML = '';
    else {
      const rel = arcs.filter((a) => a.a === best || a.b === best).map((a) => [a.a === best ? a.b : a.a, a.s]).sort((a, b) => b[1] - a[1]).slice(0, 6);
      $('arcinfo').innerHTML = `<b style="color:var(--mb-gold2)">${esc(books[best].n)}</b> ${tx('js.canon.shares')} ` + (rel.length ? rel.map(([i, s]) => `${esc(books[i].n)} (${s})`).join(' · ') : tx('js.canon.nobook'));
    }
  }
  ac.style.cursor = best >= 0 ? 'pointer' : 'default';
});
ac.addEventListener('click', () => { if (hov >= 0) location.href = '/book/' + books[hov].id; });
ac.addEventListener('mouseleave', () => { hov = -1; drawArcs(); $('arcinfo').innerHTML = ''; });

function drawAll() { drawSky(); drawArcs(); }
addEventListener('resize', drawAll);
drawAll();
