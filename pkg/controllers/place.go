package controllers

import (
	"github.com/almerlucke/metabible/pkg/pages"
	"github.com/almerlucke/metabible/pkg/store"
	"net/http"
)

type PlacePageController struct {
	Store *store.Store
}

func (pc *PlacePageController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	place, err := pc.Store.GetPlaceByID(ctx, r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if place == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	pages.PlacePage(place).Render(ctx, w)
}
