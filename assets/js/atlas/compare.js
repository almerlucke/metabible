// Compare two people or two places, with an automatically written story of the connection.
import { copyLink } from './share.js';
const $ = (id) => document.getElementById(id);
const esc = (s) => String(s ?? '').replace(/[&<>"']/g, (c) => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' }[c]));
const fmt = (y) => { y = Math.round(y); return y < 0 ? `${-y} BC` : y === 0 ? '1 BC' : `${y} AD`; };
const pct = (y) => Math.min(100, Math.max(0, ((y + 4100) / 4200) * 100));
const qs = new URLSearchParams(location.search);

let kind = qs.get('kind') === 'place' ? 'place' : 'person';
const pick = { a: qs.get('a') || '', b: qs.get('b') || '' };
const names = { a: '', b: '' };

const PRESETS = {
  person: [['Abraham', 'Sarah'], ['Moses', 'Joshua'], ['David', 'Solomon'], ['Paul', 'Simon Peter'], ['Jacob (Israel)', 'Esau (Edom)']],
  place: [['Jerusalem', 'Babylon'], ['Bethlehem', 'Nazareth'], ['Egypt', 'Canaan'], ['Corinth', 'Ephesus']],
};

$('cmpui').innerHTML = `<div class="kindsw"><button class="mb-btn" data-k="person">People</button><button class="mb-btn" data-k="place">Places</button></div>
  <div class="pickers">
    <div class="picker"><input id="pa" class="mb-input" placeholder="First…" autocomplete="off" aria-label="First"><div class="res" id="ra"></div></div>
    <button class="mb-btn" id="swap" title="Swap" aria-label="Swap">⇄</button>
    <div class="picker"><input id="pb" class="mb-input" placeholder="Second…" autocomplete="off" aria-label="Second"><div class="res" id="rb"></div></div>
  </div>
  <div class="presets" id="presets"></div>`;

function drawKind() {
  $('cmpui').querySelectorAll('.kindsw .mb-btn').forEach((b) => b.classList.toggle('on', b.dataset.k === kind));
  $('presets').innerHTML = '<span class="mb-muted">Try:</span> ' + PRESETS[kind].map((p, i) => `<button class="mb-chip pre" data-i="${i}">${esc(p[0])} &amp; ${esc(p[1])}</button>`).join('');
  $('presets').querySelectorAll('.pre').forEach((b) => b.onclick = async () => {
    const [x, y] = PRESETS[kind][+b.dataset.i];
    const [ia, ib] = await Promise.all([find(x), find(y)]);
    if (ia && ib) { setPick('a', ia); setPick('b', ib); go(); }
  });
}
async function find(name) {
  const hits = await fetch('/api/atlas/search?q=' + encodeURIComponent(name.replace(/\s*\(.*\)/, ''))).then((r) => r.json());
  return hits.find((h) => h.t === kind && h.l === name) || hits.find((h) => h.t === kind);
}
function setPick(side, h) { pick[side] = h.id; names[side] = h.l; $('p' + side).value = h.l; $('r' + side).innerHTML = ''; }

for (const side of ['a', 'b']) {
  let t;
  $('p' + side).addEventListener('input', (e) => {
    clearTimeout(t); const q = e.target.value.trim();
    if (q.length < 2) { $('r' + side).innerHTML = ''; return; }
    t = setTimeout(async () => {
      const hits = (await fetch('/api/atlas/search?q=' + encodeURIComponent(q)).then((r) => r.json())).filter((h) => h.t === kind).slice(0, 6);
      $('r' + side).innerHTML = hits.map((h) => `<div class="jr" data-id="${h.id}" data-l="${esc(h.l)}">${esc(h.l)} <span class="mb-muted">${esc((h.s || '').replace(/^(person|place) · /, ''))}</span></div>`).join('');
      $('r' + side).querySelectorAll('.jr').forEach((d) => d.onclick = () => { setPick(side, { id: d.dataset.id, l: d.dataset.l }); if (pick.a && pick.b) go(); });
    }, 130);
  });
}
$('swap').onclick = () => { [pick.a, pick.b] = [pick.b, pick.a]; [names.a, names.b] = [names.b, names.a]; $('pa').value = names.a; $('pb').value = names.b; if (pick.a && pick.b) go(); };
$('cmpui').querySelectorAll('.kindsw .mb-btn').forEach((b) => b.onclick = () => { kind = b.dataset.k; pick.a = pick.b = ''; $('pa').value = $('pb').value = ''; $('cmpout').innerHTML = ''; history.replaceState(null, '', '/compare?kind=' + kind); drawKind(); });
drawKind();

async function go() {
  history.replaceState(null, '', `/compare?kind=${kind}&a=${pick.a}&b=${pick.b}`);
  $('cmpout').innerHTML = '<p class="mb-muted">Comparing…</p>';
  const r = await fetch(`/api/atlas/compare?kind=${kind}&a=${pick.a}&b=${pick.b}`).then((x) => (x.ok ? x.json() : null));
  if (!r) { $('cmpout').innerHTML = '<p class="mb-muted">Could not compare those.</p>'; return; }
  names.a = r.a.n; names.b = r.b.n; $('pa').value = r.a.n; $('pb').value = r.b.n;
  render(r);
}

function span(s, color) {
  if (!s.hasYears) return '<div class="cspan none">no dates</div>';
  return `<div class="cspan"><i style="left:${pct(s.from)}%;width:${Math.max(1, pct(s.to) - pct(s.from))}%;background:${color}"></i></div><div class="yrs">${fmt(s.from)} – ${fmt(s.to)}${s.exact ? '' : ' <small>(approx.)</small>'}</div>`;
}
function bars(r) {
  const max = Math.max(1, ...r.a.books, ...r.b.books);
  return `<div class="bookbars" role="img" aria-label="Mentions per book for both">${r.bookNames.map((n, i) => `<div class="bb" title="${esc(n)}: ${esc(r.a.n)} ${r.a.books[i]}, ${esc(r.b.n)} ${r.b.books[i]}"><i class="a" style="height:${(r.a.books[i] / max) * 100}%"></i><i class="b" style="height:${(r.b.books[i] / max) * 100}%"></i></div>`).join('')}</div>
    <div class="bblbl"><span>Genesis</span><span>Malachi</span><span>Matthew</span><span>Revelation</span></div>`;
}
const ref = (kindName, x) => `<a class="mb-chip" href="/${kindName}/${x.id}">${esc(x.Name || x.name || x.n)}</a>`;

function render(r) {
  const link = (s) => `<a href="/${r.kind}/${s.id}">${esc(s.n)}</a>`;
  const rel = r.relation;
  let h = `<div class="cheads"><div class="side a"><div class="mb-tag">${esc(r.a.sub || r.kind)}</div><h2>${link(r.a)}</h2>${span(r.a, '#f5c451')}<div class="mb-muted">${r.a.vc} verses</div></div>
    <div class="vs">vs</div><div class="side b"><div class="mb-tag">${esc(r.b.sub || r.kind)}</div><h2>${link(r.b)}</h2>${span(r.b, '#7df9ff')}<div class="mb-muted">${r.b.vc} verses</div></div></div>`;
  h += `<section class="cstory"><h3>Their story</h3>${r.story.map((p) => `<p>${esc(p)}</p>`).join('')}
    ${r.quoteText ? `<blockquote class="mb-verse">“${esc(r.quoteText)}” <a href="/verse/${r.quoteVerse}">${esc(r.quoteLabel)}</a></blockquote>` : ''}
    <div class="actions"><button class="mb-btn" id="cshare">🔗 Copy link</button>${r.kind === 'person' ? `<a class="mb-btn" href="/tree/${r.a.id}?to=${r.b.id}">🌳 Show on the family tree</a>` : ''}
    ${r.kind === 'place' ? `<a class="mb-btn" href="/atlas?place=${r.a.id}">⏳ On the time atlas</a>` : ''}</div></section>`;
  if (rel) h += `<section><h3>Family connection</h3><div class="relchain">${rel.steps.map((s, i) => `${i ? `<span class="via">${s.via === 'up' ? 'child of' : s.via === 'down' ? 'parent of' : 'married to'}</span>` : ''}<a class="mb-chip" href="/person/${s.p.ID || s.p.id}">${esc(s.p.Name || s.p.name || s.p.n)}</a>`).join('')}</div></section>`;
  h += `<section><h3>Where they appear in the text</h3><div class="legend"><i class="a"></i>${esc(r.a.n)} <i class="b"></i>${esc(r.b.n)}</div>${bars(r)}</section>`;
  if (r.sharedPlaces && r.sharedPlaces.length) h += `<section><h3>Places both are named at</h3>${r.sharedPlaces.map((x) => ref('place', x)).join('')}</section>`;
  if (r.sharedPeople && r.sharedPeople.length) h += `<section><h3>People who appear at both</h3>${r.sharedPeople.map((x) => ref('person', x)).join('')}</section>`;
  if (r.sharedEvents && r.sharedEvents.length) h += `<section><h3>Shared events</h3><ul class="rows">${r.sharedEvents.map((e) => `<li><a href="/event/${e.ID || e.id}"><span class="when">${esc(e.Date || e.date)}</span> ${esc(e.Title || e.title)}</a></li>`).join('')}</ul></section>`;
  if (r.sharedVerses && r.sharedVerses.length) h += `<section><h3>Named together (${r.sharedCount} verses)</h3><div class="vlist">${r.sharedVerses.map((v) => `<a class="vref" href="/verse/${v.ID || v.id}" title="${esc(v.Excerpt || v.excerpt)}">${esc(v.Label || v.label)}</a>`).join('')}</div></section>`;
  $('cmpout').innerHTML = h;
  $('cshare').onclick = () => copyLink(location.href, $('cshare'));
}
if (pick.a && pick.b) go();
