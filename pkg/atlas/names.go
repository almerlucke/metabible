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
		if strings.Count(base, ".") != 1 { // nl.sample.json is a preview
			continue
		}
		raw, err := os.ReadFile(f)
		m := map[string]string{}
		if err == nil && json.Unmarshal(raw, &m) == nil && len(m) > 0 {
			a.easton[strings.TrimSuffix(base, ".json")] = m
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
	if nt == nil && titles == nil && dict == nil && a.artwork[lang] == nil {
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
			c.Dict, c.DictMT = d, true
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
			c.Dict, c.DictMT = d, true
		}
		v.Places[i] = &c
		v.PlaceByID[c.ID] = &c
	}
	if titles != nil {
		moved := map[*Event]*Event{}
		v.Events = make([]*Event, len(a.Events))
		v.EventByID = make(map[string]*Event, len(a.Events))
		for i, e := range a.Events {
			c := *e
			if t := titles[e.ID]; t != "" {
				c.Orig, c.Title = e.Title, t
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
	v.nameIdx = nil
	v.buildSearch()
	a.views.m[lang] = v
	return v
}
