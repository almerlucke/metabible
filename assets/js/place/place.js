const place = JSON.parse(document.getElementById('place').textContent);
const lat = parseFloat(place.latitude);
const lng = parseFloat(place.longitude);

// Creating map options
var mapOptions = {
    center: [lat, lng],
    zoom: 5
}

// Creating a map object
var map = new L.map('map', mapOptions);

// Creating a Layer object
var layer = L.tileLayer('https://server.arcgisonline.com/ArcGIS/rest/services/World_Topo_Map/MapServer/tile/{z}/{y}/{x}', {
    attribution: 'Tiles &copy; Esri &mdash; Esri, DeLorme, NAVTEQ, TomTom, Intermap, iPC, USGS, FAO, NPS, NRCAN, GeoBase, Kadaster NL, Ordnance Survey, Esri Japan, METI, Esri China (Hong Kong), and the GIS User Community'
});

// Adding layer to the map
map.addLayer(layer);

L.marker([lat, lng]).addTo(map).bindPopup(place.displayTitle).openPopup();