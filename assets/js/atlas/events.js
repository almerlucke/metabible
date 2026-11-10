// Live filter for the events timeline.
const q = document.getElementById('evq');
const cards = [...document.querySelectorAll('.evc')];
const sections = [...document.querySelectorAll('section.era')];
q.addEventListener('input', () => {
  const t = q.value.trim().toLowerCase();
  for (const c of cards) c.hidden = t && !c.dataset.q.toLowerCase().includes(t);
  for (const s of sections) s.hidden = ![...s.querySelectorAll('.evc')].some((c) => !c.hidden);
});
