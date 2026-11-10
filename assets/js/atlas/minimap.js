// Small embedded maps on the dossier pages. Config comes from data-config on .mb-minimap.
/* global maplibregl */
import { createMap, kindColorExpr, FONT, esc } from './map.js';

document.querySelectorAll('.mb-minimap').forEach((el) => {
  if (el.dataset.ready) return; el.dataset.ready = '1';
  const cfg = JSON.parse(el.dataset.config || '{}');
  const pts = cfg.points || [];
  if (!pts.length) { el.style.display = 'none'; return; }
  const wide = el.classList.contains('wide');
  const first = cfg.focus || pts[0];
  const map = createMap(el, { center: [first.lng, first.lat], zoom: cfg.zoom || 6, pitch: wide ? (cfg.pitch || 0) : 0, terrain: wide, cooperative: true, controls: wide });
  const fc = { type: 'FeatureCollection', features: pts.map((p) => ({ type: 'Feature', geometry: { type: 'Point', coordinates: [p.lng, p.lat] }, properties: { id: p.id || '', n: p.n, k: p.k || '', focus: !!(cfg.focus && cfg.focus.id === p.id) } })) };
  map.on('load', () => {
    map.addSource('pts', { type: 'geojson', data: fc });
    map.addLayer({ id: 'halo', type: 'circle', source: 'pts', filter: ['get', 'focus'], paint: { 'circle-radius': 22, 'circle-color': '#ffd24a', 'circle-opacity': 0.3, 'circle-blur': 0.7 } });
    map.addLayer({ id: 'dots', type: 'circle', source: 'pts', paint: { 'circle-radius': ['case', ['get', 'focus'], 8, 5.5], 'circle-color': ['case', ['get', 'focus'], '#ffffff', kindColorExpr()], 'circle-stroke-color': cfg.color || '#f5c451', 'circle-stroke-width': ['case', ['get', 'focus'], 3, 1.2] } });
    map.addLayer({ id: 'lbl', type: 'symbol', source: 'pts', layout: { 'text-field': ['get', 'n'], 'text-font': FONT, 'text-size': 12, 'text-offset': [0, 1.1], 'text-anchor': 'top', 'text-optional': true }, paint: { 'text-color': '#fff6dc', 'text-halo-color': 'rgba(5,7,16,.95)', 'text-halo-width': 1.5 } });
    if (!cfg.focus && pts.length > 1) {
      const b = pts.reduce((bb, p) => bb.extend([p.lng, p.lat]), new maplibregl.LngLatBounds([pts[0].lng, pts[0].lat], [pts[0].lng, pts[0].lat]));
      map.fitBounds(b, { padding: 36, maxZoom: 7.5, duration: 0 });
    }
    map.on('click', 'dots', (e) => {
      const f = e.features[0].properties;
      new maplibregl.Popup({ offset: 8 }).setLngLat(e.features[0].geometry.coordinates).setHTML(f.id ? `<a href="/place/${f.id}"><b>${esc(f.n)}</b></a>` : `<b>${esc(f.n)}</b>`).addTo(map);
    });
    map.on('mouseenter', 'dots', () => { map.getCanvas().style.cursor = 'pointer'; });
    map.on('mouseleave', 'dots', () => { map.getCanvas().style.cursor = ''; });
  });
  if (cfg.link) {
    const a = document.createElement('a'); a.href = cfg.link; a.className = 'mb-btn mm-open'; a.textContent = tx('js.minimap.open'); el.appendChild(a);
  }
});
