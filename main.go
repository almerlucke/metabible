package main

import (
	"context"
	"fmt"
	"github.com/almerlucke/metabible/config"
	"github.com/almerlucke/metabible/pkg/models"
	"github.com/almerlucke/metabible/pkg/store"
	"github.com/almerlucke/metabible/templates"
	"log"
	"net/http"
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
	st *store.Store
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	places, err := h.st.GetAllPlaces(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	templates.Index(places).Render(r.Context(), w)
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
	mux.Handle("/", &Handler{st: st})

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", mux)
}
