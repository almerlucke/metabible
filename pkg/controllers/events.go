package controllers

import (
	"github.com/almerlucke/metabible/pkg/pages"
	"github.com/almerlucke/metabible/pkg/store"
	"net/http"
)

type EventsPageController struct {
	Store *store.Store
}

func (ec *EventsPageController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	search := r.URL.Query().Get("search")

	events, err := ec.Store.EventSearch(ctx, search)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pages.EventsPage(events, search).Render(ctx, w)
}
