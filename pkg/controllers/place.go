package controllers

import (
	"context"
	"github.com/almerlucke/metabible/pkg/store"
	"github.com/almerlucke/metabible/pkg/templates"
	"net/http"
	"time"
)

type PlacePageController struct {
	Store *store.Store
}

func (pc *PlacePageController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	placeLookup := r.PathValue("placeLookup")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	place, err := pc.Store.GetPlaceByLookup(ctx, placeLookup)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if place == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	templates.PlacePage(place).Render(r.Context(), w)
}
