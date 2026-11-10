package atlas

import (
	"compress/gzip"
	"encoding/json"
	"net/http"
	"sort"
	"strings"

	"github.com/almerlucke/metabible/pkg/i18n"
)

type Era struct {
	ID    string `json:"id"`
	Name  string `json:"n"`
	From  int    `json:"from"`
	To    int    `json:"to"`
	Color string `json:"c"`
}

var Eras = []Era{
	{"primeval", "Primeval history", -4100, -2100, "#7e57c2"},
	{"patriarchs", "Patriarchs", -2100, -1700, "#26a69a"},
	{"egypt", "Egypt & Exodus", -1700, -1400, "#ef6c00"},
	{"conquest", "Conquest & Judges", -1400, -1050, "#c0392b"},
	{"kingdom", "United Kingdom", -1050, -930, "#f5c451"},
	{"divided", "Divided Kingdom", -930, -586, "#5c6bc0"},
	{"exile", "Exile & Return", -586, -400, "#8d6e63"},
	{"silence", "Between the Testaments", -400, -5, "#607d8b"},
	{"gospels", "Life of Jesus", -5, 33, "#ffe082"},
	{"church", "The Early Church", 33, 100, "#4fc3f7"},
}

type gzw struct {
	http.ResponseWriter
	gz *gzip.Writer
}

func (g gzw) Write(b []byte) (int, error) { return g.gz.Write(b) }

func writeJSON(w http.ResponseWriter, r *http.Request, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if w.Header().Get("Cache-Control") == "" {
		w.Header().Set("Cache-Control", "public, max-age=300")
	}
	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Add("Vary", "Accept-Encoding")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		_ = json.NewEncoder(gzw{w, gz}).Encode(v)
		return
	}
	_ = json.NewEncoder(w).Encode(v)
}

// Mount registers all /api/atlas routes.
func (a *Atlas) Mount(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/atlas/world", a.handleWorld)
	mux.HandleFunc("GET /api/atlas/people", a.handlePeople)
	mux.HandleFunc("GET /api/atlas/journeys", a.handleJourneys)
	mux.HandleFunc("GET /api/atlas/journey/{id}", a.handlePersonJourney)
	mux.HandleFunc("GET /api/atlas/cosmos", a.handleCosmos)
	mux.HandleFunc("GET /api/atlas/story", a.handleStory)
	mux.HandleFunc("GET /api/atlas/validate", a.handleValidate)
	mux.HandleFunc("GET /api/atlas/compare", a.handleCompare)
	mux.HandleFunc("GET /export/journey/{file}", a.handleExport)
	mux.HandleFunc("GET /api/atlas/progress-index", a.handleProgressIndex)
	mux.HandleFunc("GET /api/atlas/themes", a.handleThemes)
	mux.HandleFunc("GET /api/atlas/theme/{id}", a.handleTheme)
	mux.HandleFunc("GET /api/atlas/relation", a.handleRelation)
	mux.HandleFunc("GET /api/atlas/migration", a.handleMigration)
	mux.HandleFunc("GET /api/atlas/canon", a.handleCanon)
	mux.HandleFunc("GET /api/atlas/word", a.handleWord)
	mux.HandleFunc("GET /api/atlas/search", a.handleSearch)
	mux.HandleFunc("GET /api/atlas/peek/{kind}/{id}", a.handlePeek)
	mux.HandleFunc("GET /api/atlas/person/{id}", a.handlePerson)
}

func (a *Atlas) handleWorld(w http.ResponseWriter, r *http.Request) {
	places := make([]*Place, 0, len(a.Places))
	for _, p := range a.Places {
		if p.HasGeo {
			places = append(places, p)
		}
	}
	writeJSON(w, r, map[string]any{"places": places, "events": a.Events, "eras": a.LocalEras(i18n.FromContext(r.Context()))})
}

func (a *Atlas) handlePeople(w http.ResponseWriter, r *http.Request) {
	out := make([]*Person, 0, len(a.People))
	for _, p := range a.People {
		if p.Verses > 0 || len(p.Children) > 0 || p.Father != "" {
			out = append(out, p)
		}
	}
	lang := i18n.FromContext(r.Context())
	names := make([]string, len(a.Groups))
	for i, g := range a.Groups {
		names[i] = a.lookup(lang, "group", g.Name)
	}
	writeJSON(w, r, map[string]any{"people": out, "groups": names})
}

func (a *Atlas) handleJourneys(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, r, a.JourneysFor(i18n.FromContext(r.Context())))
}

func (a *Atlas) handlePersonJourney(w http.ResponseWriter, r *http.Request) {
	j := a.PersonJourney(r.PathValue("id"), i18n.FromContext(r.Context()))
	if j == nil {
		http.NotFound(w, r)
		return
	}
	writeJSON(w, r, j)
}

func (a *Atlas) handleCosmos(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, r, a.Cosmos())
}

func (a *Atlas) handleSearch(w http.ResponseWriter, r *http.Request) {
	hits := a.Search(r.URL.Query().Get("q"), 14, i18n.FromContext(r.Context()))
	if hits == nil {
		hits = []SearchHit{}
	}
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, r, hits)
}

func (a *Atlas) handlePeek(w http.ResponseWriter, r *http.Request) {
	pk := a.PeekFor(r.PathValue("kind"), r.PathValue("id"), i18n.FromContext(r.Context()))
	if pk == nil {
		http.NotFound(w, r)
		return
	}
	writeJSON(w, r, pk)
}

// handlePerson returns details used by the family tree side panel.
func (a *Atlas) handlePerson(w http.ResponseWriter, r *http.Request) {
	p := a.PersonByID[r.PathValue("id")]
	if p == nil {
		http.NotFound(w, r)
		return
	}
	var verses []map[string]string
	for i, id := range p.Vs {
		if i >= 3 {
			break
		}
		if v := a.VerseByID[id]; v != nil {
			verses = append(verses, map[string]string{"id": v.ID, "l": a.verseLabel(v), "t": a.VerseIn(v, i18n.FromContext(r.Context()))})
		}
	}
	writeJSON(w, r, map[string]any{"id": p.ID, "text": trunc(plain(p.Dict), 600), "verses": verses})
}

type arc struct {
	A      int `json:"a"`
	B      int `json:"b"`
	Shared int `json:"s"`
}

func (a *Atlas) handleCanon(w http.ResponseWriter, r *http.Request) {
	bookIdx := map[string]int{}
	for i, b := range a.Books {
		bookIdx[b.ID] = i
	}
	personBooks := map[string]map[int]bool{}
	for _, v := range a.Verses {
		bi, ok := bookIdx[v.Book]
		if !ok {
			continue
		}
		for _, p := range v.People {
			m := personBooks[p]
			if m == nil {
				m = map[int]bool{}
				personBooks[p] = m
			}
			m[bi] = true
		}
	}
	pair := map[[2]int]int{}
	for _, m := range personBooks {
		if len(m) > 12 { // ubiquitous names (God, Israel …) would drown the signal
			continue
		}
		var bs []int
		for b := range m {
			bs = append(bs, b)
		}
		sort.Ints(bs)
		for i := 0; i < len(bs); i++ {
			for j := i + 1; j < len(bs); j++ {
				pair[[2]int{bs[i], bs[j]}]++
			}
		}
	}
	var arcs []arc
	for k, s := range pair {
		if s >= 8 {
			arcs = append(arcs, arc{k[0], k[1], s})
		}
	}
	sort.Slice(arcs, func(i, j int) bool { return arcs[i].Shared > arcs[j].Shared })
	if len(arcs) > 400 {
		arcs = arcs[:400]
	}
	writeJSON(w, r, map[string]any{"books": a.Books, "arcs": arcs})
}

func (a *Atlas) handleWord(w http.ResponseWriter, r *http.Request) {
	q := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))
	if len(q) < 2 {
		writeJSON(w, r, map[string]any{"q": q, "total": 0, "books": [][]int{}})
		return
	}
	chIdx := map[string][2]int{}
	for bi, b := range a.Books {
		for ci, c := range b.Chapters {
			chIdx[c.ID] = [2]int{bi, ci}
		}
	}
	counts := make([][]int, len(a.Books))
	for i, b := range a.Books {
		counts[i] = make([]int, len(b.Chapters))
	}
	total := 0
	for i, t := range a.lowerText {
		n := strings.Count(t, q)
		if n == 0 {
			continue
		}
		if k, ok := chIdx[a.Verses[i].Chapter]; ok {
			counts[k[0]][k[1]] += n
			total += n
		}
	}
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, r, map[string]any{"q": q, "total": total, "books": counts})
}

func (a *Atlas) handleStory(w http.ResponseWriter, r *http.Request) {
	lang := i18n.FromContext(r.Context())
	a.storyMu.Lock()
	st := a.stories[lang]
	if st == nil {
		if a.stories == nil {
			a.stories = map[string]*Story{}
		}
		st = a.Story(lang)
		a.stories[lang] = st
	}
	a.storyMu.Unlock()
	writeJSON(w, r, st)
}

func (a *Atlas) handleMigration(w http.ResponseWriter, r *http.Request) {
	a.migOnce.Do(func() { a.mig = a.Migration() })
	lang := i18n.FromContext(r.Context())
	m := *a.mig
	m.Groups = make([]MigGroup, len(a.mig.Groups))
	for i, g := range a.mig.Groups {
		g.Name = a.lookup(lang, "group", g.Name)
		m.Groups[i] = g
	}
	writeJSON(w, r, &m)
}

// handleValidate reports data problems found by Validate (for the maintainer).
func (a *Atlas) handleValidate(w http.ResponseWriter, r *http.Request) {
	issues := a.Validate()
	if issues == nil {
		issues = []Issue{}
	}
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, r, issues)
}

func (a *Atlas) handleCompare(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	res := a.Compare(q.Get("kind"), q.Get("a"), q.Get("b"), i18n.FromContext(r.Context()))
	if res == nil {
		http.NotFound(w, r)
		return
	}
	writeJSON(w, r, res)
}

// handleRelation returns the family path between two people.
func (a *Atlas) handleRelation(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	rel := a.FamilyPath(q.Get("a"), q.Get("b"), i18n.FromContext(r.Context()))
	if rel == nil {
		http.NotFound(w, r)
		return
	}
	writeJSON(w, r, rel)
}

func (a *Atlas) handleThemes(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, r, a.Themes(i18n.FromContext(r.Context())))
}

func (a *Atlas) handleTheme(w http.ResponseWriter, r *http.Request) {
	t := a.Theme(r.PathValue("id"), i18n.FromContext(r.Context()))
	if t == nil {
		http.NotFound(w, r)
		return
	}
	writeJSON(w, r, t)
}

func (a *Atlas) handleProgressIndex(w http.ResponseWriter, r *http.Request) {
	a.progOnce.Do(func() { a.prog = a.ProgressIndex() })
	writeJSON(w, r, a.prog)
}
