# MetaBible

An interactive atlas of scripture, built on the cleaned
[Theographic](https://github.com/robertrouse/theographic-metadata) data in `data/json`
(books, chapters, verses, people, places, events, periods, groups).

Go + [templ](https://templ.guide) on the server, MapLibre GL and three.js in the browser.
There is no database: everything is loaded from `data/json` into memory at startup.

## Run

```sh
go run .            # serves on http://localhost:3000
```

During development `air` also works (`templ generate && go build -o ./tmp/main .`, see `.air.toml`).
After editing a `.templ` file run `templ generate`.

```sh
go test ./pkg/...   # data layer + every route
```

## What is where

| Page | Route | What it does |
|---|---|---|
| Story | `/story` | The whole Bible in order: era intros, each event with its verses, people and places, optional read-aloud. `?s=<step>` links to a moment |
| Time Atlas | `/atlas` | 3D terrain map with a 4004 BC – AD 100 timeline; `?year=&c=lng,lat,zoom,pitch,bearing&person=` |
| Journeys | `/journeys` | Curated routes (`data/journeys.json`) and a derived itinerary for any person; `?j=<id>&s=<stop>` or `?person=<id>` |
| Family | `/tree`, `/tree/{id}` | Ancestors and descendants; `?to=<id>` shows the family path between two people |
| Cosmos | `/cosmos` | 3D constellation of people and places, time on the vertical axis; `?id=` |
| Canon | `/canon` | Book skyline, shared-people arcs, word tide |
| Lifespans | `/lifespans` | Who was alive when; `?year=&person=` |
| Migrations | `/migrations` | Heatmap of where the text is per century, plus a trail per tribe; `?b=&mode=&g=` |
| Browse | `/books` `/people` `/places` `/events` `/concordance` | Filterable lists |
| Dossiers | `/person/{id}` `/place/{id}` `/event/{id}` `/verse/{id}` `/chapter/{id}` `/book/{id}` | Detail pages. Chapters are a reading view whose map follows the verse in view |

Press **Ctrl-K** (or `/`) anywhere to search people, places, events, books and verse references
(`John 3:16`). Hovering a link to a person, place, verse, book or chapter shows a peek card.

## Code map

- `pkg/atlas` loads the JSON and builds the indexes, the dossier/browse view models and the
  JSON API under `/api/atlas/*` (`world`, `people`, `journeys`, `journey/{id}`, `cosmos`, `canon`,
  `word`, `search`, `peek/{kind}/{id}`, `story`, `migration`).
- `pkg/controllers/atlas.go` registers every page and API route.
- `pkg/pages/*.templ` are the page templates (`atlas.templ` shell and immersive pages,
  `dossier.templ`, `browse.templ`).
- `assets/js/atlas/*.js` are one ES module per page; `assets/js/site/site.js` is the shared
  palette and peek cards; `assets/css/metabible.css` is the single stylesheet.
- `assets/vendor` bundles three.js, MapLibre GL, the fonts and the map label glyphs, so the only
  network access the site needs is for map tiles and elevation data (Esri, CARTO, AWS terrain tiles).
- `data/journeys.json` holds the hand-written tours. Add one there; each stop needs a name,
  coordinates and a verse reference (`Gen.12.1`), and is checked at startup.

## Notes on the data

- Birth and death are known for only ~75 people. For everyone else the active years are estimated
  from the densest cluster of verses they appear in, so lifespans and the "as they appear in the text"
  bars are approximate.
- Routes for arbitrary people are derived from the places named in their verses, in order of the text.
- Tribe trails and the heatmap count places named in verses dated to each century. They show where
  the text is, not population.

The old Mongo-backed code (`pkg/store`, `pkg/models`, `docker-compose.yml`) is no longer used by
the running site; `cmd/conversion` still uses the models to clean the source data.
