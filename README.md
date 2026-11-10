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

### Data checks and corrections

On startup the server logs how many errors and warnings the data check found; the details are at
`/api/atlas/validate`. `go test` fails on any *error* (dangling references, impossible dates).
Known errors in the source data are fixed in `data/corrections.json` (people's birth/death years,
event dates), each with a note on why, so the cleaned JSON stays untouched.

### Open content

`go run ./cmd/enrich` fetches Wikipedia summaries and photographs (Wikimedia Commons, free licenses
only) for the most prominent people and places, plus Gustave Doré's public-domain engravings matched
to events and people, into `data/enrichment.json` with the author and license of everything. It is
polite (a few requests per second at most, honours `Retry-After`, sends no personal data) and resumes
where it stopped. If the file is missing the site simply has no photos.

### Other

- Journeys download as GPX or KML (`/export/journey/<id>.gpx`).
- Pages print cleanly (maps and controls hidden).
- The site is an installable app (`/assets/manifest.webmanifest`) and keeps working offline for pages
  already visited; bump `VERSION` in `assets/sw.js` when static assets change.
- Links shared from the site show a generated preview card (`/og.png`).

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
| Compare | `/compare` | Two people or places side by side, with a generated story of how they connect; `?kind=&a=&b=` |
| Themes | `/themes` | Follow an idea (covenant, exile, prayer, water…) across books, eras and the map; `?t=` |
| Quiz | `/quiz` | A daily round of five questions made from the data, or a random one |
| My reading | `/progress` | Chapters you have read, places met, time covered. Kept in your browser only; export/import as JSON |
| Browse | `/books` `/people` `/places` `/events` `/concordance` | Filterable lists |
| Credits | `/credits` | Sources and licenses, including every image |
| Dossiers | `/person/{id}` `/place/{id}` `/event/{id}` `/verse/{id}` `/chapter/{id}` `/book/{id}` | Detail pages. Chapters are a reading view whose map follows the verse in view |

Press **P** for a full-screen presentation mode (Esc leaves it), and **Ctrl-K** (or `/`) anywhere to search people, places, events, books and verse references
(`John 3:16`). Hovering a link to a person, place, verse, book or chapter shows a peek card.

## Languages

The interface is available in English and Dutch, switched with the EN/NL links in the top bar (the choice is
remembered in a cookie; the browser's language is used the first time).

- All text lives in `pkg/i18n/locales/<code>.json`, flat `"dotted.key": "text"` files embedded in the binary.
  English is the source and the fallback, so an unfinished language still works.
- `{placeholders}` are filled in by name; counted texts use `key.one` / `key.other`.
  Keys starting with `js.` are also sent to the browser (`/i18n/<code>.js`), where `tx('js.key', {…})` reads them.
- Sentences the server writes itself (relationship stories, era introductions, journey notes, theme names,
  hover cards, search results) use keys too, so they follow the language.
- **Add a language:** copy `en.json` to `<code>.json`, set `_name`, translate the values, and add nothing else.
  `go test ./pkg/i18n` lists any missing keys and checks the placeholders match.
- **Bible text:** extra translations are `data/bible/<code>.json`, `{"01001001": "text"}` keyed by the verse id of
  the source data (`data/bible/nl.json` is the Statenvertaling, 1637, public domain, converted from the
  [open-bibles](https://github.com/seven1m/open-bibles) Zefania file). Verses a translation lacks fall back to the King James text.
  Verse pages, the reading view, concordance, search, hover cards, story and journeys follow the language; a translated verse page
  also shows the King James text with its person and place links. A translation is plain text (no inline links).
- Bible book names are keys too (`book.<OSIS>`), and a reference can be typed in any language (`Richteren 3:4`, `Judges 3:4`).
- Still English everywhere: the dictionary descriptions of people and places (Easton), people/place names, and image titles. Translating those is separate work (a name table).

## Code map

- `pkg/atlas` loads the JSON and builds the indexes, the dossier/browse view models and the
  JSON API under `/api/atlas/*` (`world`, `people`, `journeys`, `journey/{id}`, `cosmos`, `canon`,
  `word`, `search`, `peek/{kind}/{id}`, `story`, `migration`, `compare`, `relation`, `themes`, `theme/{id}`,
  `progress-index`, `validate`).
- `pkg/controllers/atlas.go` registers every page and API route; `pkg/i18n` holds the translations.
- `pkg/pages/*.templ` are the page templates (`atlas.templ` shell and immersive pages,
  `dossier.templ`, `browse.templ`).
- `assets/js/atlas/*.js` are one ES module per page; `assets/js/site/site.js` is the shared
  palette and peek cards; `assets/css/metabible.css` is the single stylesheet.
- `assets/vendor` bundles three.js, MapLibre GL, the fonts and the map label glyphs, so the only
  network access the site needs is for map tiles and elevation data (Esri, CARTO, AWS terrain tiles).
- `data/journeys.json` holds the hand-written tours. Add one there; each stop needs a name,
  coordinates and a verse reference (`Gen.12.1`), and is checked at startup.

## License

The data and everything derived from it follow the source's license, Creative Commons
Attribution-ShareAlike 4.0 (Theographic). Wikipedia text is CC BY-SA 4.0 as well. See `/credits`.

## Notes on the data

- Birth and death are known for only ~75 people. For everyone else the active years are estimated
  from the densest cluster of verses they appear in, so lifespans and the "as they appear in the text"
  bars are approximate.
- Routes for arbitrary people are derived from the places named in their verses, in order of the text.
- Tribe trails and the heatmap count places named in verses dated to each century. They show where
  the text is, not population.

The old Mongo-backed code (`pkg/store`, `pkg/models`, `docker-compose.yml`) is no longer used by
the running site; `cmd/conversion` still uses the models to clean the source data.
