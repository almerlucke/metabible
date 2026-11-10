package main

import (
	"context"
	"fmt"
	"github.com/almerlucke/metabible/config"
	"github.com/almerlucke/metabible/pkg/atlas"
	"github.com/almerlucke/metabible/pkg/controllers"
	"github.com/almerlucke/metabible/pkg/models"
	"github.com/almerlucke/metabible/pkg/store"
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
	at, err := atlas.Load("data/json", "data/journeys.json")
	if err != nil {
		log.Panic(err)
	}
	controllers.MountAtlas(mux, at) // "/", /atlas, /journeys, /tree, /cosmos, /canon, /lifespans, /api/atlas/*
	mux.Handle("/books", &controllers.BooksPageController{Store: st})
	mux.Handle("/people", &controllers.PeoplePageController{Store: st})
	mux.Handle("/places", &controllers.PlacesPageController{Store: st})
	mux.Handle("/events", &controllers.EventsPageController{Store: st})
	mux.Handle("/concordance", &controllers.ConcordancePageController{Store: st})

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", mux)
}
