// Small helpers for shareable links.
export async function copyLink(url, btn, label = window.tx('js.share.copy')) {
  try { await navigator.clipboard.writeText(url); if (btn) btn.textContent = window.tx('js.share.copied'); }
  catch (e) { prompt(window.tx('js.share.prompt'), url); }
  if (btn) setTimeout(() => { btn.textContent = label; }, 1800);
}
export function setQuery(params) {
  const q = new URLSearchParams();
  for (const [k, v] of Object.entries(params)) if (v !== undefined && v !== null && v !== '') q.set(k, v);
  const s = q.toString();
  history.replaceState(null, '', location.pathname + (s ? '?' + s : ''));
}
export const debounce = (fn, ms) => { let t; return (...a) => { clearTimeout(t); t = setTimeout(() => fn(...a), ms); }; };
