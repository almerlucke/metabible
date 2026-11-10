// Quiz: questions made from the data. A daily round is the same for everyone (seeded by the date).
const $ = (id) => document.getElementById(id);
const esc = (s) => String(s ?? '').replace(/[&<>"']/g, (c) => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' }[c]));
const fmt = (y) => window.mbYear(y);
const qs = new URLSearchParams(location.search);

const [world, peopleData] = await Promise.all([fetch('/api/atlas/world').then((r) => r.json()), fetch('/api/atlas/people').then((r) => r.json())]);
const events = world.events.filter((e) => e.t && e.y != null);
for (const e of events) { e.pl = e.pl || []; e.pe = e.pe || []; }
const placeById = new Map(world.places.map((p) => [p.id, p]));
const people = peopleData.people;
const byId = new Map(people.map((p) => [p.id, p]));

// seeded random
function rng(seed) { let a = seed >>> 0; return () => { a |= 0; a = (a + 0x6D2B79F5) | 0; let t = Math.imul(a ^ (a >>> 15), 1 | a); t = (t + Math.imul(t ^ (t >>> 7), 61 | t)) ^ t; return ((t ^ (t >>> 14)) >>> 0) / 4294967296; }; }
const today = new Date().toISOString().slice(0, 10);
const daily = qs.get('mode') !== 'random';
let rand = rng(daily ? [...today].reduce((a, c) => a * 31 + c.charCodeAt(0), 7) : Date.now());
const pickOne = (a) => a[Math.floor(rand() * a.length)];
const shuffle = (a) => { a = a.slice(); for (let i = a.length - 1; i > 0; i--) { const j = Math.floor(rand() * (i + 1)); [a[i], a[j]] = [a[j], a[i]]; } return a; };
const sample = (a, n, ok = () => true) => shuffle(a.filter(ok)).slice(0, n);

// ---------------------------------------------------------------- question makers
const prominent = people.filter((p) => p.vc >= 20);
const exactLives = people.filter((p) => p.b != null && p.d != null && p.vc >= 8);
const datedEvents = events.filter((e) => e.tx && !/^(Birth|Death|Lifetime) of/.test(e.t));

function qWhichFirst() {
  for (let k = 0; k < 40; k++) {
    const [x, y] = sample(datedEvents, 2);
    if (x && y && Math.abs(x.y - y.y) > 40) {
      const first = x.y < y.y ? x : y;
      return { kind: 'choice', text: tx('js.quiz.q.first'), options: shuffle([x.t, y.t]), answer: first.t, explain: tx('js.quiz.e.first', { a: first.t, da: mbDate(first.dt), db: mbDate(first === x ? y.dt : x.dt) }), link: `/event/${first.id}` };
    }
  }
}
function qWhoAlive() {
  const a = pickOne(exactLives); if (!a) return null;
  const year = Math.round(a.b + (a.d - a.b) * (0.3 + rand() * 0.4));
  const others = sample(exactLives, 3, (p) => p.id !== a.id && (p.d < year - 5 || p.b > year + 5));
  if (others.length < 3) return null;
  return { kind: 'choice', text: tx('js.quiz.q.alive', { y: fmt(year) }), options: shuffle([a.n, ...others.map((o) => o.n)]), answer: a.n, explain: tx('js.quiz.e.alive', { n: a.n, a: fmt(a.b), b: fmt(a.d) }), link: `/person/${a.id}` };
}
function qFather() {
  const kids = prominent.filter((p) => p.f && byId.get(p.f) && byId.get(p.f).vc >= 6);
  const c = pickOne(kids); if (!c) return null;
  const f = byId.get(c.f);
  const others = sample(prominent, 3, (p) => p.g === 'Male' && p.id !== f.id && p.id !== c.id);
  return { kind: 'choice', text: tx('js.quiz.q.father', { n: c.n }), options: shuffle([f.n, ...others.map((o) => o.n)]), answer: f.n, explain: tx('js.quiz.e.father', { f: f.n, c: c.n }), link: `/person/${c.id}` };
}
function qWhere() {
  const evs = datedEvents.filter((e) => e.pl.length === 1 && placeById.get(e.pl[0]) && placeById.get(e.pl[0]).vc >= 8);
  const e = pickOne(evs); if (!e) return null;
  const p = placeById.get(e.pl[0]);
  const others = sample(world.places, 3, (x) => x.id !== p.id && x.k === p.k && x.vc >= 25);
  if (others.length < 3) return null;
  return { kind: 'choice', text: tx('js.quiz.q.where', { t: e.t }), hint: e.tx ? `“${e.tx.slice(0, 140)}${e.tx.length > 140 ? '…' : ''}” ${e.ref || ''}` : '', options: shuffle([p.n, ...others.map((o) => o.n)]), answer: p.n, explain: tx('js.quiz.e.where', { t: e.t, p: p.n }), link: `/place/${p.id}` };
}
function qOrder() {
  const picks = sample(datedEvents, 4).sort((a, b) => a.y - b.y);
  if (picks.length < 4 || picks[3].y - picks[0].y < 150 || picks.some((p, i) => i && p.y - picks[i - 1].y < 25)) return qOrder();
  return { kind: 'order', text: tx('js.quiz.q.order'), items: shuffle(picks.map((p) => p.t)), answer: picks.map((p) => p.t), explain: picks.map((p) => `${p.t} (${mbDate(p.dt)})`).join(' → ') };
}
function qLifespan() {
  const a = pickOne(exactLives); if (!a) return null;
  const real = a.d - a.b, opts = new Set([real]);
  while (opts.size < 4) { const d = Math.max(20, real + Math.round((rand() - 0.5) * 140 / 5) * 5); opts.add(d); }
  return { kind: 'choice', text: tx('js.quiz.q.lifespan', { n: a.n }), options: shuffle([...opts].map((n) => tx('js.quiz.years', { n }))), answer: tx('js.quiz.years', { n: real }), explain: tx('js.quiz.e.lifespan', { n: a.n, a: fmt(a.b), b: fmt(a.d) }), link: `/person/${a.id}` };
}

function makeRound() {
  const makers = [qWhichFirst, qWhoAlive, qFather, qWhere, qOrder, qLifespan];
  const qs2 = [];
  for (const m of shuffle(makers)) { const q = m(); if (q) qs2.push(q); if (qs2.length === 5) break; }
  return qs2;
}

// ---------------------------------------------------------------- state and UI
const round = makeRound();
let i = 0, score = 0, results = [], chosen = [];
const statKey = 'mb.quiz';
function stats() { try { return JSON.parse(localStorage.getItem(statKey)) || {}; } catch (e) { return {}; } }
function saveStats(s) { try { localStorage.setItem(statKey, JSON.stringify(s)); } catch (e) { /* ignore */ } }

function start() {
  const st = stats();
  const doneToday = daily && st.lastDaily === today;
  $('quiz').innerHTML = `<p class="lead mb-muted">${daily ? tx('js.quiz.daily', { date: today }) : tx('js.quiz.random')}</p>
    ${st.streak ? `<p class="mb-muted">${tx.n('js.quiz.streak', st.streak, { s: st.best || 0 })}</p>` : ''}
    ${doneToday ? `<p>${tx('js.quiz.you.already.played.today')} <b>${st.lastScore}/5</b>.</p>` : ''}
    <div class="actions"><button class="mb-btn on" id="go">${doneToday ? tx('js.quiz.again') : tx('js.quiz.start')}</button><a class="mb-btn" href="/quiz?mode=${daily ? 'random' : 'daily'}">${daily ? tx('js.quiz.random.instead') : tx('js.quiz.today')}</a></div>`;
  $('go').onclick = () => { i = 0; score = 0; results = []; show(); };
}

function show() {
  const q = round[i];
  chosen = [];
  let h = `<div class="qprog">${round.map((_, k) => `<i class="${k < i ? (results[k] ? 'ok' : 'bad') : k === i ? 'cur' : ''}"></i>`).join('')}</div>
    <h2 class="qtext" id="qtext">${esc(q.text)}</h2>${q.hint ? `<p class="mb-verse">${esc(q.hint)}</p>` : ''}<div id="qbody"></div><div id="qfeed" aria-live="polite"></div>`;
  $('quiz').innerHTML = h;
  if (q.kind === 'choice') {
    $('qbody').innerHTML = `<div class="qopts" role="group" aria-labelledby="qtext">${q.options.map((o, k) => `<button class="qopt" data-o="${esc(o)}">${esc(o)}</button>`).join('')}</div>`;
    $('qbody').querySelectorAll('.qopt').forEach((b) => b.onclick = () => answer(b.dataset.o === q.answer, b));
  } else {
    $('qbody').innerHTML = `<p class="mb-muted">${tx('js.quiz.click.the.events.in')}</p><div class="qopts" id="qitems">${q.items.map((o) => `<button class="qopt" data-o="${esc(o)}">${esc(o)}</button>`).join('')}</div><div id="qpicked" class="qpicked"></div>`;
    $('qitems').querySelectorAll('.qopt').forEach((b) => b.onclick = () => {
      if (b.disabled) return; chosen.push(b.dataset.o); b.disabled = true; b.classList.add('used');
      $('qpicked').innerHTML = chosen.map((c, k) => `<span class="mb-chip">${k + 1}. ${esc(c)}</span>`).join('');
      if (chosen.length === q.items.length) answer(chosen.every((c, k) => c === q.answer[k]), null);
    });
  }
}

function answer(ok, btn) {
  const q = round[i];
  results[i] = ok; if (ok) score++;
  document.querySelectorAll('.qopt').forEach((b) => { b.disabled = true; if (q.kind === 'choice') { if (b.dataset.o === q.answer) b.classList.add('right'); else if (b === btn) b.classList.add('wrong'); } });
  $('qfeed').innerHTML = `<p class="${ok ? 'good' : 'miss'}">${ok ? tx('js.quiz.correct') : tx('js.quiz.wrong')} ${esc(q.explain)} ${q.link ? `<a href="${q.link}">${tx('js.quiz.read.more')}</a>` : ''}</p><button class="mb-btn on" id="next">${i + 1 < round.length ? tx('js.quiz.next') : tx('js.quiz.score')}</button>`;
  $('next').focus();
  $('next').onclick = () => { i++; if (i < round.length) show(); else finish(); };
}

function finish() {
  const st = stats();
  if (daily) {
    const yesterday = new Date(Date.now() - 864e5).toISOString().slice(0, 10);
    if (st.lastDaily !== today) st.streak = st.lastDaily === yesterday ? (st.streak || 0) + 1 : 1;
    st.lastDaily = today; st.lastScore = score;
  }
  st.best = Math.max(st.best || 0, score); st.played = (st.played || 0) + 1; saveStats(st);
  const grid = results.map((r) => (r ? '🟩' : '🟥')).join('');
  const share = `${tx('js.quiz.share', { date: daily ? today : '', score })}\n${grid}\n${location.origin}/quiz`;
  $('quiz').innerHTML = `<h2>${score}/5</h2><p class="bigemoji">${grid}</p><p>${tx('js.quiz.msg' + score)}</p>
    ${st.streak ? `<p class="mb-muted">${tx.n('js.quiz.streak2', st.streak)}</p>` : ''}
    <div class="actions"><button class="mb-btn on" id="sh">${tx('js.quiz.copy.result')}</button><a class="mb-btn" href="/quiz?mode=random">${tx('js.quiz.another.round')}</a></div>`;
  $('sh').onclick = async () => { try { await navigator.clipboard.writeText(share); $('sh').textContent = tx('js.quiz.copied'); } catch (e) { prompt(tx('js.quiz.copyprompt'), share); } };
}
start();
