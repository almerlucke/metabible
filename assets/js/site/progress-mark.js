// Records that a chapter was read, in this browser only (localStorage "mb.progress").
(function () {
  var el = document.querySelector('.reader[data-chapter]');
  if (!el) return;
  var id = el.getAttribute('data-chapter');
  function load() { try { return JSON.parse(localStorage.getItem('mb.progress')) || { read: {} }; } catch (e) { return { read: {} }; } }
  function save(p) { try { localStorage.setItem('mb.progress', JSON.stringify(p)); } catch (e) { /* storage unavailable */ } }
  var p = load(); p.last = id; save(p); // remember where you were as soon as you open it
  function mark() { var q = load(); q.read = q.read || {}; if (!q.read[id]) q.read[id] = Date.now(); q.last = id; save(q); badge(); }
  function badge() { var b = document.getElementById('readbadge'); if (b) b.textContent = '✓ Marked as read'; }
  if ((load().read || {})[id]) badge();
  // read = reached the last verse, or stayed a while
  var verses = document.querySelectorAll('.verses .v');
  if (verses.length && 'IntersectionObserver' in window) {
    var io = new IntersectionObserver(function (en) { if (en[0].isIntersecting) { mark(); io.disconnect(); } }, { threshold: 0.6 });
    io.observe(verses[verses.length - 1]);
  }
  var btn = document.getElementById('readbtn');
  if (btn) btn.addEventListener('click', mark);
})();
