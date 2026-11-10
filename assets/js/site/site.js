// Shared behaviour for every page: dark theme, command palette (Ctrl-K or "/"),
// and hover "peek" cards for any link to a person, place, verse, book or chapter.
(function () {
  'use strict';
  document.documentElement.setAttribute('data-theme', 'dark');

  const esc = (s) => String(s).replace(/[&<>"']/g, (c) => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' }[c]));

  // ---------- command palette ----------
  let pal, input, list, hits = [], sel = 0, timer, seq = 0;
  const TYPE = { person: 'Person', place: 'Place', event: 'Event', book: 'Book', verse: 'Verse', chapter: 'Chapter' };

  function buildPalette() {
    pal = document.createElement('div');
    pal.id = 'mb-palette';
    pal.setAttribute('role', 'dialog'); pal.setAttribute('aria-modal', 'true'); pal.setAttribute('aria-label', 'Search');
    pal.innerHTML = '<div class="box"><input type="text" role="combobox" aria-expanded="true" aria-controls="mb-results" aria-label="Search" placeholder="Search people, places, events, books… or type John 3:16" autocomplete="off" spellcheck="false"><ul id="mb-results" role="listbox"></ul><div class="hint">↑↓ navigate · ↵ open · esc close</div></div>';
    document.body.appendChild(pal);
    input = pal.querySelector('input');
    list = pal.querySelector('ul');
    pal.addEventListener('mousedown', (e) => { if (e.target === pal) closePalette(); });
    input.addEventListener('input', () => { clearTimeout(timer); timer = setTimeout(run, 120); });
    input.addEventListener('keydown', (e) => {
      if (e.key === 'ArrowDown') { sel = Math.min(hits.length - 1, sel + 1); render(); e.preventDefault(); }
      else if (e.key === 'ArrowUp') { sel = Math.max(0, sel - 1); render(); e.preventDefault(); }
      else if (e.key === 'Enter' && hits[sel]) { go(hits[sel]); }
      else if (e.key === 'Escape') closePalette();
    });
  }
  function go(h) {
    closePalette();
    if (document.startViewTransition && !h.u.includes('?')) document.startViewTransition(() => { location.href = h.u; });
    else location.href = h.u;
  }
  async function run() {
    const q = input.value.trim();
    if (!q) { hits = []; render(); return; }
    const my = ++seq;
    try {
      const r = await fetch('/api/atlas/search?q=' + encodeURIComponent(q));
      const data = await r.json();
      if (my !== seq) return;
      hits = data; sel = 0; render();
    } catch (e) { /* ignore */ }
  }
  function render() {
    list.innerHTML = hits.map((h, i) =>
      `<li id="mb-opt-${i}" role="option" aria-selected="${i === sel}" class="${i === sel ? 'sel' : ''}" data-i="${i}"><span class="t">${TYPE[h.t] || h.t}</span><span class="l">${esc(h.l)}</span><span class="s">${esc(h.s || '')}</span></li>`).join('');
    list.querySelectorAll('li').forEach((li) => {
      li.addEventListener('mousemove', () => { sel = +li.dataset.i; list.querySelectorAll('li').forEach((x, j) => x.classList.toggle('sel', j === sel)); });
      li.addEventListener('mousedown', (e) => { e.preventDefault(); go(hits[+li.dataset.i]); });
    });
    const s = list.querySelector('.sel'); if (s) { s.scrollIntoView({ block: 'nearest' }); input.setAttribute('aria-activedescendant', s.id); }
  }
  function openPalette() { if (!pal) buildPalette(); pal.classList.add('open'); input.value = ''; hits = []; render(); input.focus(); }
  function closePalette() { if (pal) pal.classList.remove('open'); }
  window.mbOpenSearch = openPalette;

  document.addEventListener('keydown', (e) => {
    const tag = (e.target.tagName || '').toLowerCase();
    const typing = tag === 'input' || tag === 'textarea' || e.target.isContentEditable;
    if ((e.key === 'k' && (e.ctrlKey || e.metaKey)) || (e.key === '/' && !typing)) { e.preventDefault(); openPalette(); }
  });
  document.addEventListener('click', (e) => { if (e.target.closest('[data-mb-search]')) { e.preventDefault(); openPalette(); } });

  // ---------- peek cards ----------
  const cache = new Map();
  let card, hoverTimer, curKey = null;
  const LINK = /^\/(person|place|verse|book|chapter)\/([A-Za-z0-9]+)\/?$/;

  function ensureCard() {
    if (!card) { card = document.createElement('div'); card.id = 'mb-peek'; document.body.appendChild(card); }
    return card;
  }
  function showCard(a, data) {
    const c = ensureCard();
    c.innerHTML = `<div class="ti">${esc(data.title)}</div>` +
      (data.sub ? `<div class="su">${esc(data.sub)}</div>` : '') +
      (data.text ? `<div class="tx ${data.t === 'verse' ? 'verse' : ''}">${esc(data.text)}</div>` : '') +
      (data.facts && data.facts.length ? `<div class="fa">${data.facts.map(esc).join(' · ')}</div>` : '');
    const r = a.getBoundingClientRect();
    const w = 320, h = c.offsetHeight || 140;
    let x = Math.min(Math.max(8, r.left), innerWidth - w - 8);
    let y = r.bottom + 8;
    if (y + h > innerHeight - 8) y = Math.max(8, r.top - h - 8);
    c.style.left = x + 'px'; c.style.top = y + 'px';
    requestAnimationFrame(() => c.classList.add('show'));
  }
  function hideCard() { clearTimeout(hoverTimer); curKey = null; if (card) card.classList.remove('show'); }

  document.addEventListener('mouseover', (e) => {
    const a = e.target.closest && e.target.closest('a[href]');
    if (!a || a.closest('.mb-nav') || a.closest('#mb-palette')) return;
    const m = LINK.exec(a.getAttribute('href') || '');
    if (!m) return;
    const key = m[1] + '/' + m[2];
    if (curKey === key) return;
    curKey = key;
    clearTimeout(hoverTimer);
    hoverTimer = setTimeout(async () => {
      let d = cache.get(key);
      if (!d) {
        try { const r = await fetch('/api/atlas/peek/' + key); if (!r.ok) return; d = await r.json(); cache.set(key, d); } catch (err) { return; }
      }
      if (curKey === key) showCard(a, d);
    }, 280);
  });
  document.addEventListener('mouseout', (e) => {
    const a = e.target.closest && e.target.closest('a[href]');
    if (a && !a.contains(e.relatedTarget)) hideCard();
  });
  document.addEventListener('scroll', hideCard, true);

  // ---------- presentation mode (P to toggle, Esc to leave) ----------
  function setPresent(on) {
    document.body.classList.toggle('present', on);
    let hint = document.querySelector('.present-hint');
    if (!hint) { hint = document.createElement('button'); hint.className = 'present-hint'; hint.type = 'button'; hint.textContent = 'Exit presentation (Esc)'; hint.addEventListener('click', () => setPresent(false)); document.body.appendChild(hint); }
    try { const u = new URL(location.href); if (on) u.searchParams.set('present', '1'); else u.searchParams.delete('present'); history.replaceState(null, '', u); } catch (e) { /* ignore */ }
    setTimeout(() => window.dispatchEvent(new Event('resize')), 60);
  }
  document.addEventListener('click', (e) => { if (e.target.closest('[data-mb-present]')) setPresent(!document.body.classList.contains('present')); });
  document.addEventListener('keydown', (e) => {
    const t = (e.target.tagName || '').toLowerCase();
    if (t === 'input' || t === 'textarea' || t === 'select' || e.ctrlKey || e.metaKey || e.altKey) return;
    if (e.key === 'p' || e.key === 'P') setPresent(!document.body.classList.contains('present'));
    else if (e.key === 'Escape' && document.body.classList.contains('present')) setPresent(false);
  });
  if (new URLSearchParams(location.search).has('present')) document.addEventListener('DOMContentLoaded', () => setPresent(true));

  // ---------- accessibility fallbacks for generated markup ----------
  function labelThings() {
    document.querySelectorAll('button:not([aria-label])').forEach((b) => { if (!b.textContent.trim() && b.title) b.setAttribute('aria-label', b.title); });
    document.querySelectorAll('button.icon:not([aria-label])').forEach((b) => { if (b.title) b.setAttribute('aria-label', b.title); });
    document.querySelectorAll('canvas:not([aria-label])').forEach((c) => { c.setAttribute('role', 'img'); c.setAttribute('aria-label', 'Interactive visualisation: ' + document.title.replace(/ · MetaBible$/, '')); });
    document.querySelectorAll('input:not([aria-label]):not([type=hidden]):not([type=checkbox]):not([type=file])').forEach((i) => { const p = i.getAttribute('placeholder'); if (p && !i.labels.length) i.setAttribute('aria-label', p); });
  }
  let lt; const mo = new MutationObserver(() => { clearTimeout(lt); lt = setTimeout(labelThings, 200); });
  document.addEventListener('DOMContentLoaded', () => { labelThings(); mo.observe(document.body, { childList: true, subtree: true }); });

  // ---------- offline support ----------
  if ('serviceWorker' in navigator && location.protocol !== 'file:') {
    addEventListener('load', () => navigator.serviceWorker.register('/sw.js').catch(() => { /* optional */ }));
  }

  // ---------- lightbox: any element with data-lb opens its picture; arrows move through its group ----------
  let lb, lbItems = [], lbIdx = 0, lbFrom = null;
  const lbEsc = esc;
  function lbBuild() {
    lb = document.createElement('div');
    lb.id = 'mb-lb'; lb.setAttribute('role', 'dialog'); lb.setAttribute('aria-modal', 'true'); lb.setAttribute('aria-label', 'Image viewer');
    lb.innerHTML = '<button class="x" type="button" aria-label="Close">✕</button><button class="nav prev" type="button" aria-label="Previous image">‹</button>' +
      '<figure><img alt=""><figcaption><div class="t"></div><div class="c"></div><div class="links"></div></figcaption></figure>' +
      '<button class="nav next" type="button" aria-label="Next image">›</button><div class="count"></div>';
    document.body.appendChild(lb);
    lb.addEventListener('click', (e) => { if (e.target === lb || e.target.closest('.x')) lbClose(); else if (e.target.closest('.prev')) lbGo(-1); else if (e.target.closest('.next')) lbGo(1); });
  }
  function lbShow() {
    const el = lbItems[lbIdx], d = el.dataset;
    const img = lb.querySelector('img'); img.alt = d.title || ''; img.removeAttribute('src'); img.src = d.full;
    lb.querySelector('.t').textContent = d.title || '';
    lb.querySelector('.c').innerHTML = d.page ? `<a href="${lbEsc(d.page)}" target="_blank" rel="noopener">${lbEsc(d.credit || 'Source')}</a>` : lbEsc(d.credit || '');
    let links = []; try { links = JSON.parse(d.links || '[]'); } catch (e) { /* ignore */ }
    lb.querySelector('.links').innerHTML = links.map((l) => `<a class="mb-chip" href="${lbEsc(l.u)}"${/^https?:/.test(l.u) ? ' target="_blank" rel="noopener"' : ''}>${lbEsc(l.l)}</a>`).join('');
    lb.querySelector('.count').textContent = lbItems.length > 1 ? `${lbIdx + 1} / ${lbItems.length}` : '';
    lb.classList.toggle('single', lbItems.length < 2);
    for (const k of [lbIdx - 1, lbIdx + 1]) { const n = lbItems[k]; if (n) { const p = new Image(); p.src = n.dataset.full; } } // preload neighbours
  }
  function lbOpen(el) {
    if (!lb) lbBuild();
    const group = el.closest('[data-lb-group]');
    lbItems = group ? [...group.querySelectorAll('[data-lb]')].filter((x) => !x.hidden) : [el];
    lbIdx = Math.max(0, lbItems.indexOf(el)); lbFrom = el;
    lbShow(); lb.classList.add('open'); document.body.style.overflow = 'hidden';
    lb.querySelector('.x').focus();
  }
  function lbClose() { if (!lb) return; lb.classList.remove('open'); document.body.style.overflow = ''; if (lbFrom) lbFrom.focus(); }
  function lbGo(d) { if (lbItems.length < 2) return; lbIdx = (lbIdx + d + lbItems.length) % lbItems.length; lbShow(); }
  document.addEventListener('click', (e) => {
    const el = e.target.closest('[data-lb]');
    if (!el || e.target.closest('a')) return;
    e.preventDefault(); lbOpen(el);
  });
  document.addEventListener('keydown', (e) => {
    if (lb && lb.classList.contains('open')) {
      if (e.key === 'Escape') { lbClose(); e.stopPropagation(); } else if (e.key === 'ArrowRight') lbGo(1); else if (e.key === 'ArrowLeft') lbGo(-1);
      return;
    }
    const el = e.target.closest && e.target.closest('[data-lb][tabindex]');
    if (el && (e.key === 'Enter' || e.key === ' ')) { e.preventDefault(); lbOpen(el); }
  }, true);
})();

