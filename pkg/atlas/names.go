package atlas

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Spellings of people and places per language live in data/names/<lang>.json, made by
// `go run ./cmd/names`: {"people": {"<id>": "Izak"}, "places": {...}}. What a file leaves out
// keeps its source (English) name.

type nameTable struct {
	People map[string]string `json:"people"`
	Places map[string]string `json:"places"`
}

// wikiTable is data/enrichment.<lang>.json (cmd/wikilang): the summary of the same Wikipedia
// article in the language.
type wikiTable struct {
	Places map[string]WikiEntry `json:"places"`
	People map[string]WikiEntry `json:"people"`
}

type viewSet struct {
	mu sync.Mutex
	m  map[string]*Atlas
}

func (a *Atlas) loadNames(dir string) {
	files, _ := filepath.Glob(filepath.Join(dir, "..", "names", "*.json"))
	a.names = map[string]*nameTable{}
	a.titles = map[string]map[string]string{}
	a.easton = map[string]map[string]string{}
	a.artwork = map[string]map[string]string{}
	a.edited = map[string]map[string]bool{}
	a.wikis = map[string]*wikiTable{}
	wfiles, _ := filepath.Glob(filepath.Join(dir, "..", "enrichment.*.json"))
	for _, f := range wfiles {
		var t wikiTable
		raw, err := os.ReadFile(f)
		if err != nil || json.Unmarshal(raw, &t) != nil || len(t.Places)+len(t.People) == 0 {
			continue
		}
		code := strings.TrimSuffix(strings.TrimPrefix(filepath.Base(f), "enrichment."), ".json")
		a.wikis[code] = &t
	}
	afiles, _ := filepath.Glob(filepath.Join(dir, "..", "artwork", "*.json"))
	for _, f := range afiles {
		raw, err := os.ReadFile(f)
		m := map[string]string{}
		if err == nil && json.Unmarshal(raw, &m) == nil && len(m) > 0 {
			a.artwork[strings.TrimSuffix(filepath.Base(f), ".json")] = m
		}
	}
	efiles, _ := filepath.Glob(filepath.Join(dir, "..", "easton", "*.json"))
	for _, f := range efiles {
		base := filepath.Base(f)
		if strings.Count(base, ".") != 1 { // nl.hand.json is read below, nl.sample.json is a preview
			continue
		}
		raw, err := os.ReadFile(f)
		m := map[string]string{}
		if err == nil && json.Unmarshal(raw, &m) == nil && len(m) > 0 {
			a.easton[strings.TrimSuffix(base, ".json")] = m
		}
	}
	// <lang>.hand.json holds texts translated by hand; they replace the machine translation
	hfiles, _ := filepath.Glob(filepath.Join(dir, "..", "easton", "*.hand.json"))
	for _, f := range hfiles {
		lang := strings.TrimSuffix(filepath.Base(f), ".hand.json")
		raw, err := os.ReadFile(f)
		m := map[string]string{}
		if err != nil || json.Unmarshal(raw, &m) != nil {
			continue
		}
		if a.easton[lang] == nil {
			a.easton[lang] = map[string]string{}
		}
		a.edited[lang] = map[string]bool{}
		for id, text := range m {
			if text = strings.TrimSpace(text); text != "" {
				a.easton[lang][id], a.edited[lang][id] = text, true
			}
		}
	}
	tfiles, _ := filepath.Glob(filepath.Join(dir, "..", "titles", "*.json"))
	for _, f := range tfiles {
		raw, err := os.ReadFile(f)
		m := map[string]string{}
		if err == nil && json.Unmarshal(raw, &m) == nil && len(m) > 0 {
			a.titles[strings.TrimSuffix(filepath.Base(f), ".json")] = m
		}
	}
	for _, f := range files {
		base := filepath.Base(f)
		if strings.Count(base, ".") != 1 { // nl.glue.json is input for the tool, not a table
			continue
		}
		raw, err := os.ReadFile(f)
		if err != nil {
			continue
		}
		var t nameTable
		if json.Unmarshal(raw, &t) != nil || len(t.People)+len(t.Places) == 0 {
			continue
		}
		a.names[strings.TrimSuffix(base, ".json")] = &t
	}
}

// View is the atlas as a visitor of the language sees it: the same data with people and
// places under their local names. Without a name table for the language it is the atlas itself.
// Views are built once and share everything but the names (and what is cached from them).
func (a *Atlas) View(lang string) *Atlas {
	if a.base != nil {
		return a.base.View(lang)
	}
	nt, titles := a.names[lang], a.titles[lang]
	dict := a.easton[lang]
	wiki := a.wikis[lang]
	if nt == nil && titles == nil && dict == nil && a.artwork[lang] == nil && a.bibles[lang] == nil && wiki == nil {
		return a
	}
	if nt == nil {
		nt = &nameTable{}
	}
	a.views.mu.Lock()
	defer a.views.mu.Unlock()
	if v := a.views.m[lang]; v != nil {
		return v
	}
	v := new(Atlas)
	*v = *a
	v.lazy = &lazy{}
	v.base = a
	v.lang = lang
	v.People = make([]*Person, len(a.People))
	v.PersonByID = make(map[string]*Person, len(a.People))
	for i, p := range a.People {
		c := *p
		if n := nt.People[p.ID]; n != "" {
			c.Orig, c.Name = p.Name, n
		}
		if d := dict[p.ID]; d != "" {
			c.Dict, c.DictMT = d, !a.edited[lang][p.ID]
		}
		v.People[i] = &c
		v.PersonByID[c.ID] = &c
	}
	v.Places = make([]*Place, len(a.Places))
	v.PlaceByID = make(map[string]*Place, len(a.Places))
	for i, p := range a.Places {
		c := *p
		if n := nt.Places[p.ID]; n != "" {
			c.Orig, c.Name = p.Name, n
		}
		if d := dict[p.ID]; d != "" {
			c.Dict, c.DictMT = d, !a.edited[lang][p.ID]
		}
		v.Places[i] = &c
		v.PlaceByID[c.ID] = &c
	}
	if titles != nil || a.bibles[lang] != nil {
		moved := map[*Event]*Event{}
		v.Events = make([]*Event, len(a.Events))
		v.EventByID = make(map[string]*Event, len(a.Events))
		for i, e := range a.Events {
			c := *e
			if t := titles[e.ID]; t != "" {
				c.Orig, c.Title = e.Title, t
			}
			c.Text = a.verseText(e.VerseID, lang, e.Text) // the verse the event points at, in the language
			if v := a.VerseByID[e.VerseID]; v != nil {
				c.Ref = a.verseLabel(v, lang) // "Genesis 1:1" with the local book name, instead of the OSIS "Gen.1.1"
			}
			moved[e], v.Events[i], v.EventByID[c.ID] = &c, &c, &c
		}
		remap := func(m map[string][]*Event) map[string][]*Event {
			out := make(map[string][]*Event, len(m))
			for k, es := range m {
				ne := make([]*Event, len(es))
				for i, e := range es {
					ne[i] = moved[e]
				}
				out[k] = ne
			}
			return out
		}
		v.personEvents, v.placeEvents, v.eventSucc = remap(a.personEvents), remap(a.placeEvents), remap(a.eventSucc)
	}
	if wiki != nil && a.Enrich != nil {
		// the article of the language replaces the English text; the photos stay
		en := *a.Enrich
		en.Places = make(map[string]WikiEntry, len(a.Enrich.Places))
		for id, w := range a.Enrich.Places {
			if l, ok := wiki.Places[id]; ok {
				w.Title, w.Summary, w.Desc, w.Page = l.Title, l.Summary, l.Desc, l.Page
			}
			en.Places[id] = w
		}
		en.People = make(map[string]WikiEntry, len(a.Enrich.People))
		for id, w := range a.Enrich.People {
			if l, ok := wiki.People[id]; ok {
				w.Title, w.Summary, w.Desc, w.Page = l.Title, l.Summary, l.Desc, l.Page
			}
			en.People[id] = w
		}
		v.Enrich = &en
	}
	v.nameIdx = nil
	v.buildSearch()
	a.views.m[lang] = v
	return v
}
