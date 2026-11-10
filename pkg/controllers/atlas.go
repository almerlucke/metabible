package controllers

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/almerlucke/metabible/pkg/atlas"
	"github.com/almerlucke/metabible/pkg/pages"
)

// MountAtlas registers the immersive atlas pages and the JSON API behind them.
// It has no dependency on Mongo.
func MountAtlas(mux *http.ServeMux, a *atlas.Atlas) {
	a.Mount(mux)

	stats := pages.HomeStats{Verses: len(a.Verses), People: len(a.People), Places: len(a.Places), Events: len(a.Events)}
	mux.Handle("GET /{$}", templ.Handler(pages.HomePage(stats)))
	mux.Handle("GET /atlas", templ.Handler(pages.AtlasPage()))
	mux.Handle("GET /journeys", templ.Handler(pages.JourneysPage()))
	mux.Handle("GET /tree", templ.Handler(pages.TreePage("")))
	mux.HandleFunc("GET /tree/{id}", func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(pages.TreePage(r.PathValue("id"))).ServeHTTP(w, r)
	})
	mux.Handle("GET /cosmos", templ.Handler(pages.CosmosPage()))
	mux.Handle("GET /canon", templ.Handler(pages.CanonPage()))
	mux.Handle("GET /lifespans", templ.Handler(pages.LifespansPage()))

	// browse pages
	mux.Handle("GET /books", templ.Handler(pages.BooksPage(a.BooksBrowse())))
	mux.Handle("GET /events", templ.Handler(pages.EventsPage(a.EventsTimeline())))
	mux.Handle("GET /people", templ.Handler(pages.PeoplePage()))
	mux.Handle("GET /places", templ.Handler(pages.PlacesPage()))
	mux.HandleFunc("GET /concordance", func(w http.ResponseWriter, r *http.Request) {
		res := a.Concordance(r.URL.Query().Get("search"), r.URL.Query().Get("book"), 400)
		templ.Handler(pages.ConcordancePage(res)).ServeHTTP(w, r)
	})

	// dossier pages, served from the in-memory atlas (no Mongo needed)
	detail := func(w http.ResponseWriter, r *http.Request, c templ.Component, ok bool) {
		if !ok {
			http.NotFound(w, r)
			return
		}
		templ.Handler(c).ServeHTTP(w, r)
	}
	mux.HandleFunc("GET /person/{id}", func(w http.ResponseWriter, r *http.Request) {
		d := a.PersonDossier(r.PathValue("id"))
		detail(w, r, pages.PersonPage(d), d != nil)
	})
	mux.HandleFunc("GET /place/{id}", func(w http.ResponseWriter, r *http.Request) {
		d := a.PlaceDossier(r.PathValue("id"))
		detail(w, r, pages.PlacePage(d), d != nil)
	})
	mux.HandleFunc("GET /event/{id}", func(w http.ResponseWriter, r *http.Request) {
		d := a.EventDossier(r.PathValue("id"))
		detail(w, r, pages.EventPage(d), d != nil)
	})
	mux.HandleFunc("GET /verse/{id}", func(w http.ResponseWriter, r *http.Request) {
		d := a.VerseDossier(r.PathValue("id"))
		detail(w, r, pages.VersePage(d), d != nil)
	})
	mux.HandleFunc("GET /chapter/{id}", func(w http.ResponseWriter, r *http.Request) {
		d := a.ChapterDossier(r.PathValue("id"))
		detail(w, r, pages.ChapterPage(d), d != nil)
	})
	mux.HandleFunc("GET /book/{id}", func(w http.ResponseWriter, r *http.Request) {
		d := a.BookDossier(r.PathValue("id"))
		detail(w, r, pages.BookPage(d), d != nil)
	})
}
