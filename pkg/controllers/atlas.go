package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/almerlucke/metabible/pkg/atlas"
	"github.com/almerlucke/metabible/pkg/i18n"
	"github.com/almerlucke/metabible/pkg/og"
	"github.com/almerlucke/metabible/pkg/pages"
)

// MountAtlas registers the immersive atlas pages and the JSON API behind them.
// It has no dependency on Mongo.
// Bundle holds the translations; Middleware applies the language chosen for each request.
var Bundle = i18n.MustLoad()

func MountAtlas(mux *http.ServeMux, a *atlas.Atlas) {
	a.Mount(mux)
	pages.SetBundle(Bundle)
	a.SetTranslations(Bundle)

	// the texts the browser needs, as a script: /i18n/nl.js
	mux.HandleFunc("GET /i18n/{file}", func(w http.ResponseWriter, r *http.Request) {
		lang := strings.TrimSuffix(r.PathValue("file"), ".js")
		if !Bundle.Has(lang) {
			lang = i18n.Default
		}
		b, _ := json.Marshal(Bundle.Client(lang))
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		w.Header().Set("Cache-Control", "public, max-age=300")
		fmt.Fprintf(w, "window.MB_LANG=%q;window.MB_I18N=%s;\n", lang, b)
	})

	stats := pages.HomeStats{Verses: len(a.Verses), People: len(a.People), Places: len(a.Places), Events: len(a.Events)}
	mux.Handle("GET /{$}", templ.Handler(pages.HomePage(stats)))
	mux.Handle("GET /migrations", templ.Handler(pages.MigrationsPage()))
	mux.Handle("GET /compare", templ.Handler(pages.ComparePage()))
	mux.Handle("GET /themes", templ.Handler(pages.ThemesPage()))
	mux.Handle("GET /progress", templ.Handler(pages.ProgressPage()))
	mux.Handle("GET /quiz", templ.Handler(pages.QuizPage()))
	mux.Handle("GET /story", templ.Handler(pages.StoryPage()))
	mux.Handle("GET /atlas", templ.Handler(pages.AtlasPage()))
	mux.Handle("GET /journeys", templ.Handler(pages.JourneysPage()))
	mux.Handle("GET /tree", templ.Handler(pages.TreePage("")))
	mux.HandleFunc("GET /tree/{id}", func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(pages.TreePage(r.PathValue("id"))).ServeHTTP(w, r)
	})
	mux.Handle("GET /cosmos", templ.Handler(pages.CosmosPage()))
	mux.Handle("GET /canon", templ.Handler(pages.CanonPage()))
	mux.Handle("GET /lifespans", templ.Handler(pages.LifespansPage()))

	mux.HandleFunc("GET /og.png", func(w http.ResponseWriter, r *http.Request) {
		t, s := r.URL.Query().Get("t"), r.URL.Query().Get("s")
		if len([]rune(t)) > 120 {
			t = string([]rune(t)[:120])
		}
		if len([]rune(s)) > 160 {
			s = string([]rune(s)[:160])
		}
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		_, _ = w.Write(og.Render(t, s))
	})
	mux.HandleFunc("GET /gallery", func(w http.ResponseWriter, r *http.Request) {
		items := a.AllGallery()
		var eng, ppl, plc int
		for _, it := range items {
			switch it.Kind {
			case "engraving":
				eng++
			case "person":
				ppl++
			case "place":
				plc++
			}
		}
		templ.Handler(pages.GalleryPage(items, eng, ppl, plc)).ServeHTTP(w, r)
	})
	mux.HandleFunc("GET /credits", func(w http.ResponseWriter, r *http.Request) {
		images, texts := a.Credits()
		templ.Handler(pages.CreditsPage(images, texts)).ServeHTTP(w, r)
	})

	// browse pages
	mux.Handle("GET /books", templ.Handler(pages.BooksPage(a.BooksBrowse())))
	mux.HandleFunc("GET /events", func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(pages.EventsPage(a.EventsTimeline(i18n.FromContext(r.Context())))).ServeHTTP(w, r)
	})
	mux.Handle("GET /people", templ.Handler(pages.PeoplePage()))
	mux.Handle("GET /places", templ.Handler(pages.PlacesPage()))
	mux.HandleFunc("GET /concordance", func(w http.ResponseWriter, r *http.Request) {
		res := a.Concordance(r.URL.Query().Get("search"), r.URL.Query().Get("book"), 400, i18n.FromContext(r.Context()))
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
		d := a.PersonDossier(r.PathValue("id"), i18n.FromContext(r.Context()))
		detail(w, r, pages.PersonPage(d), d != nil)
	})
	mux.HandleFunc("GET /place/{id}", func(w http.ResponseWriter, r *http.Request) {
		d := a.PlaceDossier(r.PathValue("id"), i18n.FromContext(r.Context()))
		detail(w, r, pages.PlacePage(d), d != nil)
	})
	mux.HandleFunc("GET /event/{id}", func(w http.ResponseWriter, r *http.Request) {
		d := a.EventDossier(r.PathValue("id"), i18n.FromContext(r.Context()))
		detail(w, r, pages.EventPage(d), d != nil)
	})
	mux.HandleFunc("GET /verse/{id}", func(w http.ResponseWriter, r *http.Request) {
		d := a.VerseDossier(r.PathValue("id"), i18n.FromContext(r.Context()))
		detail(w, r, pages.VersePage(d), d != nil)
	})
	mux.HandleFunc("GET /chapter/{id}", func(w http.ResponseWriter, r *http.Request) {
		d := a.ChapterDossier(r.PathValue("id"), i18n.FromContext(r.Context()))
		detail(w, r, pages.ChapterPage(d), d != nil)
	})
	mux.HandleFunc("GET /book/{id}", func(w http.ResponseWriter, r *http.Request) {
		d := a.BookDossier(r.PathValue("id"), i18n.FromContext(r.Context()))
		detail(w, r, pages.BookPage(d), d != nil)
	})
}
