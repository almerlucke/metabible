// People browser: filter, sort and page through everyone in the text.
const $ = (id) => document.getElementById(id);
const esc = (s) => String(s ?? '').replace(/[&<>"']/g, (c) => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' }[c]));
const fmtYear = (y) => window.mbYear(y);
const ERAS = [[tx('js.era.name.primeval'), -4100, -2100, '#9575cd'], [tx('js.era.name.patriarchs'), -2100, -1700, '#26a69a'], [tx('js.era.name.egypt'), -1700, -1400, '#ef6c00'], [tx('js.era.name.conquest'), -1400, -1050, '#e5533d'],
  [tx('js.era.name.kingdom'), -1050, -930, '#f5c451'], [tx('js.era.name.divided'), -930, -586, '#7986cb'], [tx('js.era.name.exile'), -586, -400, '#a1887f'], [tx('js.era.name.silence'), -400, -5, '#78909c'], [tx('js.era.name.gospels'), -5, 33, '#ffe082'], [tx('js.era.name.church'), 33, 100, '#4fc3f7']];
const eraOf = (y) => ERAS.findIndex((e) => y >= e[1] && y < e[2]);
const pct = (y) => Math.min(100, Math.max(0, ((y + 4100) / 4200) * 100));
const groupColor = (g) => (g ? `hsl(${(g * 47) % 360} 62% 58%)` : '#8f96b2');

const data = await fetch('/api/atlas/people').then((r) => r.json());
const groups = data.groups;
const people = data.people.map((p) => {
  const exact = p.b != null && p.d != null;
  let a = p.b != null ? p.b : p.y0, b = p.d != null ? p.d : Math.max(p.y1, a + 1);
  if (p.b != null && p.d == null) b = Math.min(b, a + 120);
  let lo = Math.min(a, b), hi = Math.max(a, b);
  const known = p.b != null || p.y0 || p.y1;
  const wide = !exact && hi - lo > 140;
  if (!exact && hi - lo < 70) { const m = (lo + hi) / 2; lo = m - 35; hi = m + 35; }
  return { ...p, lo, hi, exact, known: known && !wide, mid: (lo + hi) / 2, lname: p.n.toLowerCase() };
});

// controls
$('pgroup').innerHTML += groups.map((g, i) => `<option value="${i + 1}">${esc(g)}</option>`).join('');
$('pera').innerHTML += ERAS.map((e, i) => `<option value="${i}">${e[0]}</option>`).join('');
const letters = [...new Set(people.map((p) => p.n[0].toUpperCase()))].filter((l) => /[A-Z]/.test(l)).sort();
$('paz').innerHTML = '<a href="#" data-l="" class="on">' + tx('js.people.all') + '</a>' + letters.map((l) => `<a href="#" data-l="${l}">${l}</a>`).join('');
let letter = '', shown = 0, list = [];
const PAGE = 96;

function apply() {
  const q = $('pq').value.trim().toLowerCase(), g = +$('pgroup').value, e = $('pera').value, gen = $('pgender').value, sort = $('psort').value;
  list = people.filter((p) => (!q || p.lname.includes(q) || (p.dis || '').toLowerCase().includes(q))
    && (!g || p.gr === g) && (e === '' || (p.known && eraOf(p.mid) === +e)) && (!gen || p.g === gen) && (!letter || p.n[0].toUpperCase() === letter));
  if (sort === 'az') list.sort((a, b) => a.lname.localeCompare(b.lname));
  else if (sort === 'time') list = list.filter((p) => p.known).sort((a, b) => a.mid - b.mid);
  else list.sort((a, b) => b.vc - a.vc);
  $('pcount').textContent = `${list.length.toLocaleString()} ${list.length === 1 ? 'person' : 'people'}`;
  $('pgrid').innerHTML = ''; shown = 0; more();
}
function card(p) {
  const era = p.known ? ERAS[Math.max(0, eraOf(p.mid))] : null;
  const life = era ? `<div class="mini-life"><i class="${p.exact ? 'exact' : ''}" style="left:${pct(p.lo)}%;width:${Math.max(1.2, pct(p.hi) - pct(p.lo))}%;background:${era[3]}"></i></div><div class="yrs">${fmtYear(p.lo)} – ${fmtYear(p.hi)}</div>` : '<div class="mini-life empty"></div><div class="yrs">' + tx('js.people.no.date') + '</div>';
  return `<a class="pcard" href="/person/${p.id}" style="--g:${groupColor(p.gr)}">
    <div class="nm">${esc(p.n)}</div>
    <div class="ds">${esc(p.dis || '')}</div>${life}
    <div class="ft"><span>${p.gr ? esc(groups[p.gr - 1]) : ''}</span><span>${tx.n('js.verses', p.vc)}</span></div></a>`;
}
function more() {
  const next = list.slice(shown, shown + PAGE);
  $('pgrid').insertAdjacentHTML('beforeend', next.map(card).join(''));
  shown += next.length;
  $('pmore').innerHTML = shown < list.length ? `<button class="mb-btn" id="pmorebtn">${tx('js.people.showmore', { n: (list.length - shown).toLocaleString(mbLang) })}</button>` : '';
  if ($('pmorebtn')) $('pmorebtn').onclick = more;
}
const io = new IntersectionObserver((en) => { if (en[0].isIntersecting && shown < list.length) more(); }, { rootMargin: '600px' });
io.observe($('pmore'));
for (const id of ['pq', 'pgroup', 'pera', 'pgender', 'psort']) $(id).addEventListener(id === 'pq' ? 'input' : 'change', apply);
$('paz').addEventListener('click', (e) => { const a = e.target.closest('a'); if (!a) return; e.preventDefault(); letter = a.dataset.l; $('paz').querySelectorAll('a').forEach((x) => x.classList.toggle('on', x === a)); apply(); });
apply();
