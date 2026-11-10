package controllers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/almerlucke/metabible/pkg/atlas"
)

func TestRoutes(t *testing.T) {
	a, err := atlas.Load("../../data/json", "../../data/journeys.json")
	if err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	MountAtlas(mux, a)
	h := Middleware(mux)

	person, place, event, verse, book := a.People[5].ID, a.Places[5].ID, a.Events[5].ID, a.Verses[500].ID, a.Books[3].ID
	chapter := a.Books[3].Chapters[1].ID
	paths := []string{
		"/", "/atlas", "/journeys", "/tree", "/tree/" + person, "/cosmos", "/canon", "/lifespans",
		"/og.png?t=Hello", "/gallery", "/credits", "/books", "/people", "/places", "/events", "/concordance", "/concordance?search=grace",
		"/person/" + person, "/place/" + place, "/event/" + event, "/verse/" + verse, "/chapter/" + chapter, "/book/" + book,
		"/story", "/migrations", "/compare", "/themes", "/progress", "/quiz", "/api/atlas/themes", "/api/atlas/progress-index", "/export/journey/exodus.gpx", "/export/journey/exodus.kml", "/api/atlas/compare?kind=person&a=" + person + "&b=" + person, "/api/atlas/theme/love", "/api/atlas/migration", "/api/atlas/story", "/api/atlas/validate", "/api/atlas/world", "/api/atlas/people", "/api/atlas/journeys", "/api/atlas/journey/" + person, "/api/atlas/canon",
		"/api/atlas/word?q=love", "/api/atlas/search?q=moses", "/api/atlas/peek/person/" + person, "/api/atlas/peek/verse/" + verse,
	}
	for _, p := range paths {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest("GET", p, nil))
		if rec.Code != 200 {
			t.Errorf("GET %s = %d", p, rec.Code)
		}
		if ct := rec.Header().Get("Content-Type"); (strings.HasPrefix(ct, "text/") || strings.HasPrefix(ct, "application/json")) && !utf8.Valid(rec.Body.Bytes()) {
			t.Errorf("GET %s is not valid UTF-8", p)
		}
	}
	for _, p := range []string{"/person/nope", "/place/nope", "/event/nope", "/verse/nope", "/chapter/nope", "/book/nope"} {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest("GET", p, nil))
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

func TestGzipAndOG(t *testing.T) {
	a, err := atlas.Load("../../data/json", "../../data/journeys.json")
	if err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	MountAtlas(mux, a)
	h := Middleware(mux)

	req := httptest.NewRequest("GET", "/books", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Header().Get("Content-Encoding") != "gzip" {
		t.Error("html should be compressed")
	}
	plain := httptest.NewRecorder()
	h.ServeHTTP(plain, httptest.NewRequest("GET", "/books", nil))
	if plain.Header().Get("Content-Encoding") != "" || !strings.Contains(plain.Body.String(), "<title>Books") {
		t.Error("clients without gzip get plain html")
	}
	if rec.Body.Len() >= plain.Body.Len() {
		t.Errorf("gzip did not shrink the page: %d vs %d", rec.Body.Len(), plain.Body.Len())
	}
	// share card meta uses the request host
	r2 := httptest.NewRequest("GET", "/credits", nil)
	r2.Host = "example.org"
	p2 := httptest.NewRecorder()
	h.ServeHTTP(p2, r2)
	if !strings.Contains(p2.Body.String(), `content="http://example.org/og.png?t=Credits"`) {
		t.Error("og:image should be absolute")
	}
	png := httptest.NewRecorder()
	h.ServeHTTP(png, httptest.NewRequest("GET", "/og.png?t=Hi", nil))
	if png.Header().Get("Content-Type") != "image/png" || png.Body.Len() < 1000 {
		t.Error("og image")
	}
}

func TestLanguageSwitch(t *testing.T) {
	a, err := atlas.Load("../../data/json", "../../data/journeys.json")
	if err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	MountAtlas(mux, a)
	h := Middleware(mux)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/books?x=1&lang=nl", nil))
	if rec.Code != 303 || rec.Header().Get("Location") != "/books?x=1" {
		t.Fatalf("switch: %d %s", rec.Code, rec.Header().Get("Location"))
	}
	var cookie *http.Cookie
	for _, c := range rec.Result().Cookies() {
		if c.Name == "mb_lang" {
			cookie = c
		}
	}
	if cookie == nil || cookie.Value != "nl" {
		t.Fatal("language cookie not set")
	}
	req := httptest.NewRequest("GET", "/books", nil)
	req.AddCookie(cookie)
	page := httptest.NewRecorder()
	h.ServeHTTP(page, req)
	body := page.Body.String()
	for _, want := range []string{`<html lang="nl"`, "Boeken", "Oude Testament"} {
		if !strings.Contains(body, want) {
			t.Errorf("dutch page lacks %q", want)
		}
	}
	// the browser's translation file
	js := httptest.NewRecorder()
	h.ServeHTTP(js, httptest.NewRequest("GET", "/i18n/nl.js", nil))
	if !strings.Contains(js.Body.String(), `window.MB_LANG="nl"`) || !strings.Contains(js.Body.String(), "v.Chr.") {
		t.Errorf("i18n script: %.120s", js.Body.String())
	}
	// APIs follow the cookie too
	api := httptest.NewRequest("GET", "/api/atlas/themes", nil)
	api.AddCookie(cookie)
	ar := httptest.NewRecorder()
	h.ServeHTTP(ar, api)
	if !strings.Contains(ar.Body.String(), "Verbond") {
		t.Errorf("themes api not in dutch: %.100s", ar.Body.String())
	}
}
