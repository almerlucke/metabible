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
	nt := a.names[lang]
	if nt == nil {
		return a
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
	v.People = make([]*Person, len(a.People))
	v.PersonByID = make(map[string]*Person, len(a.People))
	for i, p := range a.People {
		c := *p
		if n := nt.People[p.ID]; n != "" {
			c.Orig, c.Name = p.Name, n
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
		v.Places[i] = &c
		v.PlaceByID[c.ID] = &c
	}
	v.nameIdx = nil
	v.buildSearch()
	a.views.m[lang] = v
	return v
}
