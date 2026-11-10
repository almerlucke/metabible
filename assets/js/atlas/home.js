// Home hero: a slow-turning field of golden stars with a few drifting "routes".
import * as THREE from 'three';

const canvas = document.getElementById('hero');
const renderer = new THREE.WebGLRenderer({ canvas, antialias: true, alpha: true });
renderer.setPixelRatio(Math.min(devicePixelRatio, 2));
const scene = new THREE.Scene();
const camera = new THREE.PerspectiveCamera(60, 1, 0.1, 1000);
camera.position.z = 60;

function dotTexture() {
  const c = document.createElement('canvas'); c.width = c.height = 64;
  const g = c.getContext('2d'), grad = g.createRadialGradient(32, 32, 0, 32, 32, 32);
  grad.addColorStop(0, 'rgba(255,255,255,1)'); grad.addColorStop(0.35, 'rgba(255,255,255,.55)'); grad.addColorStop(1, 'rgba(255,255,255,0)');
  g.fillStyle = grad; g.fillRect(0, 0, 64, 64);
  return new THREE.CanvasTexture(c);
}
const DOT = dotTexture();
function stars(n, spread, size, color) {
  const g = new THREE.BufferGeometry();
  const p = new Float32Array(n * 3);
  for (let i = 0; i < n; i++) {
    const r = spread * Math.cbrt(Math.random());
    const t = Math.random() * Math.PI * 2, u = Math.acos(2 * Math.random() - 1);
    p[i * 3] = r * Math.sin(u) * Math.cos(t);
    p[i * 3 + 1] = r * Math.sin(u) * Math.sin(t) * 0.6;
    p[i * 3 + 2] = r * Math.cos(u);
  }
  g.setAttribute('position', new THREE.BufferAttribute(p, 3));
  const m = new THREE.PointsMaterial({ size: size * 2.2, map: DOT, color, transparent: true, opacity: 0.9, depthWrite: false, blending: THREE.AdditiveBlending, sizeAttenuation: true });
  return new THREE.Points(g, m);
}
const gold = stars(1400, 90, 0.55, 0xf5c451);
const white = stars(1800, 140, 0.3, 0xcfd8ff);
scene.add(gold, white);

// glowing arcs between random star pairs, reminiscent of routes
const arcs = new THREE.Group();
for (let i = 0; i < 14; i++) {
  const a = new THREE.Vector3((Math.random() - .5) * 120, (Math.random() - .5) * 50, (Math.random() - .5) * 60);
  const b = new THREE.Vector3((Math.random() - .5) * 120, (Math.random() - .5) * 50, (Math.random() - .5) * 60);
  const mid = a.clone().add(b).multiplyScalar(0.5); mid.y += 14 + Math.random() * 14;
  const curve = new THREE.QuadraticBezierCurve3(a, mid, b);
  const geo = new THREE.BufferGeometry().setFromPoints(curve.getPoints(60));
  arcs.add(new THREE.Line(geo, new THREE.LineBasicMaterial({ color: 0xf5c451, transparent: true, opacity: 0.18 })));
}
scene.add(arcs);

function resize() {
  const w = canvas.clientWidth, h = canvas.clientHeight;
  renderer.setSize(w, h, false);
  camera.aspect = w / h; camera.updateProjectionMatrix();
}
addEventListener('resize', resize); resize();

let mx = 0, my = 0;
addEventListener('pointermove', (e) => { mx = (e.clientX / innerWidth - .5); my = (e.clientY / innerHeight - .5); });
const reduce = matchMedia('(prefers-reduced-motion: reduce)').matches;
let t = 0;
function frame() {
  t += reduce ? 0 : 0.0006;
  gold.rotation.y = t; white.rotation.y = t * 0.5; arcs.rotation.y = t * 0.8;
  camera.position.x += (mx * 14 - camera.position.x) * 0.03;
  camera.position.y += (-my * 8 - camera.position.y) * 0.03;
  camera.lookAt(0, 0, 0);
  renderer.render(scene, camera);
  requestAnimationFrame(frame);
}
frame();
