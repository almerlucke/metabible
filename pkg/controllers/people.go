package controllers

import (
	"github.com/almerlucke/metabible/pkg/pages"
	"github.com/almerlucke/metabible/pkg/store"
	"log"
	"net/http"
)

type PeoplePageController struct {
	Store *store.Store
}

func (pc *PeoplePageController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	search := r.URL.Query().Get("search")

	people, err := pc.Store.PeopleSearch(ctx, search)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error searching people: %v", err)
		return
	}

	pages.PeoplePage(people, search).Render(ctx, w)
}
