package controllers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/almerlucke/metabible/pkg/atlas"
)

func TestRoutes(t *testing.T) {
	a, err := atlas.Load("../../data/json", "../../data/journeys.json")
	if err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	MountAtlas(mux, a)

	person, place, event, verse, book := a.People[5].ID, a.Places[5].ID, a.Events[5].ID, a.Verses[500].ID, a.Books[3].ID
	chapter := a.Books[3].Chapters[1].ID
	paths := []string{
		"/", "/atlas", "/journeys", "/tree", "/tree/" + person, "/cosmos", "/canon", "/lifespans",
		"/books", "/people", "/places", "/events", "/concordance", "/concordance?search=grace",
		"/person/" + person, "/place/" + place, "/event/" + event, "/verse/" + verse, "/chapter/" + chapter, "/book/" + book,
		"/story", "/migrations", "/api/atlas/migration", "/api/atlas/story", "/api/atlas/world", "/api/atlas/people", "/api/atlas/journeys", "/api/atlas/journey/" + person, "/api/atlas/canon",
		"/api/atlas/word?q=love", "/api/atlas/search?q=moses", "/api/atlas/peek/person/" + person, "/api/atlas/peek/verse/" + verse,
	}
	for _, p := range paths {
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, httptest.NewRequest("GET", p, nil))
		if rec.Code != 200 {
			t.Errorf("GET %s = %d", p, rec.Code)
		}
	}
	for _, p := range []string{"/person/nope", "/place/nope", "/event/nope", "/verse/nope", "/chapter/nope", "/book/nope"} {
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, httptest.NewRequest("GET", p, nil))
		if rec.Code != 404 {
			t.Errorf("GET %s = %d, want 404", p, rec.Code)
		}
	}
	// pages must not leak unescaped search input
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("GET", "/concordance?search=%3Cscript%3Ealert(1)%3C/script%3E", nil))
	if strings.Contains(rec.Body.String(), "<script>alert(1)") {
		t.Error("concordance echoed unescaped input")
	}
}
