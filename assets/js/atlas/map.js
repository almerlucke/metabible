// Shared MapLibre helpers: basemaps, 3D terrain, kind colours.
/* global maplibregl */

export const BASEMAPS = {
  satellite: {
    label: 'Satellite',
    tiles: ['https://server.arcgisonline.com/ArcGIS/rest/services/World_Imagery/MapServer/tile/{z}/{y}/{x}'],
    maxzoom: 17, attribution: 'Imagery © Esri, Maxar, Earthstar Geographics',
  },
  physical: {
    label: 'Physical',
    tiles: ['https://server.arcgisonline.com/ArcGIS/rest/services/World_Physical_Map/MapServer/tile/{z}/{y}/{x}'],
    maxzoom: 8, attribution: 'Tiles © Esri, US National Park Service',
  },
  topo: {
    label: 'Topographic',
    tiles: ['https://server.arcgisonline.com/ArcGIS/rest/services/World_Topo_Map/MapServer/tile/{z}/{y}/{x}'],
    maxzoom: 17, attribution: 'Tiles © Esri',
  },
  dark: {
    label: 'Night',
    tiles: ['https://basemaps.cartocdn.com/dark_nolabels/{z}/{x}/{y}.png'],
    maxzoom: 17, attribution: '© OpenStreetMap contributors © CARTO',
  },
};

export const KIND_COLOR = {
  City: '#f5c451', Region: '#6fb1ff', Mountain: '#d9a066', Water: '#4fc3f7',
  Valley: '#8bc34a', Island: '#ce93d8', Landmark: '#ff8a65', Path: '#ffffff', '': '#b8bed4',
};
export const kindColorExpr = () => ['match', ['get', 'k'],
  'City', KIND_COLOR.City, 'Region', KIND_COLOR.Region, 'Mountain', KIND_COLOR.Mountain, 'Water', KIND_COLOR.Water,
  'Valley', KIND_COLOR.Valley, 'Island', KIND_COLOR.Island, 'Landmark', KIND_COLOR.Landmark, KIND_COLOR['']];

export const FONT = ['Open Sans Regular', 'Arial Unicode MS Regular'];

export function createMap(container, { center = [35.2, 32], zoom = 5, pitch = 50, bearing = 0, basemap = 'satellite', terrain = true, cooperative = false, controls = true } = {}) {
  const bm = BASEMAPS[basemap];
  const map = new maplibregl.Map({
    container,
    center, zoom, pitch, bearing,
    maxPitch: 75,
    cooperativeGestures: cooperative,
    attributionControl: { compact: true },
    style: {
      version: 8,
      glyphs: '/assets/vendor/glyphs/{fontstack}/{range}.pbf',
      sources: {
        base: { type: 'raster', tiles: bm.tiles, tileSize: 256, maxzoom: bm.maxzoom, attribution: bm.attribution },
        dem2: { type: 'raster-dem', tiles: ['https://s3.amazonaws.com/elevation-tiles-prod/terrarium/{z}/{x}/{y}.png'], encoding: 'terrarium', tileSize: 256, maxzoom: 12 },
        dem: { type: 'raster-dem', tiles: ['https://s3.amazonaws.com/elevation-tiles-prod/terrarium/{z}/{x}/{y}.png'], encoding: 'terrarium', tileSize: 256, maxzoom: 12 },
      },
      layers: [
        { id: 'base', type: 'raster', source: 'base', paint: { 'raster-saturation': -0.15, 'raster-brightness-max': 0.92 } },
        { id: 'hill', type: 'hillshade', source: 'dem2', paint: { 'hillshade-exaggeration': 0.35, 'hillshade-shadow-color': '#000010', 'hillshade-highlight-color': '#fff6d8', 'hillshade-accent-color': '#101030' } },
      ],
      sky: { 'sky-color': '#07102a', 'horizon-color': '#27355f', 'fog-color': '#0a0e1f', 'sky-horizon-blend': 0.6, 'horizon-fog-blend': 0.8, 'fog-ground-blend': 0.9 },
    },
  });
  if (controls) map.addControl(new maplibregl.NavigationControl({ visualizePitch: true }), 'bottom-right');
  map.on('style.load', () => { if (terrain) map.setTerrain({ source: 'dem', exaggeration: 1.6 }); });
  return map;
}

export function setBasemap(map, name) {
  const bm = BASEMAPS[name];
  if (!bm) return;
  const src = map.getSource('base');
  // MapLibre can't swap tiles in place: replace the source+layer under the hillshade.
  map.removeLayer('base');
  map.removeSource('base');
  map.addSource('base', { type: 'raster', tiles: bm.tiles, tileSize: 256, maxzoom: bm.maxzoom, attribution: bm.attribution });
  map.addLayer({ id: 'base', type: 'raster', source: 'base', paint: { 'raster-saturation': name === 'dark' ? 0 : -0.15 } }, 'hill');
  void src;
}

export function setTerrain(map, on) {
  map.setTerrain(on ? { source: 'dem', exaggeration: 1.6 } : null);
  if (!on) map.easeTo({ pitch: 0, duration: 800 });
}

export const esc = (s) => String(s ?? '').replace(/[&<>"']/g, (c) => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' }[c]));

export function fmtYear(y) {
  y = Math.round(y);
  return y < 0 ? `${-y} BC` : y === 0 ? '1 BC' : `${y} AD`;
}

// Quadratic bezier "arc" between two lng/lat points, bowed to the left of travel.
export function arc(a, b, bow = 0.18, n = 28) {
  const dx = b[0] - a[0], dy = b[1] - a[1];
  const mx = (a[0] + b[0]) / 2 - dy * bow, my = (a[1] + b[1]) / 2 + dx * bow;
  const pts = [];
  for (let i = 0; i <= n; i++) {
    const t = i / n, u = 1 - t;
    pts.push([u * u * a[0] + 2 * u * t * mx + t * t * b[0], u * u * a[1] + 2 * u * t * my + t * t * b[1]]);
  }
  return pts;
}

export function ease(t) { return t < 0.5 ? 2 * t * t : 1 - Math.pow(-2 * t + 2, 2) / 2; }
