// Get all anchor tags on the page
var block = document.getElementById('placeBlock');

block.innerHTML = block.innerHTML.replaceAll(/[0-9]?[a-z]+#(.+?\.[0-9]+\.[0-9]+)/gm, "verse/$1?osisRef=true")