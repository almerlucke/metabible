// Gallery page: filter the mosaic by kind and title.
const tiles = [...document.querySelectorAll('.mosaic .tile')];
let kind = '', q = '';
function apply() {
  for (const t of tiles) t.hidden = (kind && t.dataset.kind !== kind) || (q && !t.dataset.q.includes(q));
}
document.querySelectorAll('#galtabs button').forEach((b) => b.addEventListener('click', () => {
  kind = b.dataset.f; document.querySelectorAll('#galtabs button').forEach((x) => x.classList.toggle('on', x === b)); apply();
}));
const input = document.getElementById('galq');
if (input) input.addEventListener('input', () => { q = input.value.trim().toLowerCase(); apply(); });
