package controllers

import (
	"github.com/almerlucke/metabible/pkg/pages"
	"github.com/almerlucke/metabible/pkg/store"
	"log"
	"net/http"
)

type PlacesPageController struct {
	Store *store.Store
}

func (pc *PlacesPageController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	search := r.URL.Query().Get("search")

	places, err := pc.Store.PlaceSearch(ctx, search)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error searching places: %v", err)
		return
	}

	pages.PlacesPage(places, search).Render(ctx, w)
}
