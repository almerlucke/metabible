package main

import (
	"context"
	"fmt"
	"github.com/almerlucke/metabible/config"
	"github.com/almerlucke/metabible/pkg/controllers"
	"github.com/almerlucke/metabible/pkg/models"
	"github.com/almerlucke/metabible/pkg/pages"
	"github.com/almerlucke/metabible/pkg/store"
	"log"
	"net/http"
	"regexp"
	"time"
)

func seed(st *store.Store, ctx context.Context) {
	err := store.Seed[*models.Book](st, ctx, "data/json/books.json", "books")
	if err != nil {
		log.Panic(err)
	}
	err = store.Seed[*models.Chapter](st, ctx, "data/json/chapters.json", "chapters")
	if err != nil {
		log.Panic(err)
	}
	err = store.Seed[*models.Easton](st, ctx, "data/json/easton.json", "easton")
	if err != nil {
		log.Panic(err)
	}
	err = store.Seed[*models.Event](st, ctx, "data/json/events.json", "events")
	if err != nil {
		log.Panic(err)
	}
	err = store.Seed[*models.Person](st, ctx, "data/json/people.json", "people")
	if err != nil {
		log.Panic(err)
	}
	err = store.Seed[*models.PeopleGroup](st, ctx, "data/json/peopleGroups.json", "people_groups")
	if err != nil {
		log.Panic(err)
	}
	err = store.Seed[*models.Period](st, ctx, "data/json/periods.json", "periods")
	if err != nil {
		log.Panic(err)
	}
	err = store.Seed[*models.Place](st, ctx, "data/json/places.json", "places")
	if err != nil {
		log.Panic(err)
	}
	err = store.Seed[*models.Verse](st, ctx, "data/json/verses.json", "verses")
	if err != nil {
		log.Panic(err)
	}
}

type Handler struct {
	st         *store.Store
	verseRegex *regexp.Regexp
}

func NewHandler(st *store.Store) *Handler {
	return &Handler{
		st:         st,
		verseRegex: regexp.MustCompile(`[0-9]?[a-z]+#(.+\.[0-9]+\.[0-9]+)`),
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	places, err := h.st.GetAllPlaces(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pages.Index(places).Render(r.Context(), w)
}

func main() {
	cfg := config.MustLoad("config/dev.json")

	st, err := store.New(cfg.Mongo)
	if err != nil {
		log.Panic(err)
	}

	defer st.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	seed(st, ctx)

	mux := http.NewServeMux()

	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	mux.Handle("/", &controllers.PlacesPageController{Store: st})
	mux.Handle("/places", &controllers.PlacesPageController{Store: st})
	mux.Handle("/people", &controllers.PeoplePageController{Store: st})
	mux.Handle("/place/{id}", &controllers.PlacePageController{Store: st})
	mux.Handle("/person/{id}", &controllers.PersonPageController{Store: st})
	mux.Handle("/verse/{id}", &controllers.VersePageController{Store: st})
	mux.Handle("/chapter/{id}", &controllers.ChapterPageController{Store: st})
	mux.Handle("/book/{id}", &controllers.BookPageController{Store: st})

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", mux)
}
