package controllers

import (
	"github.com/almerlucke/metabible/pkg/pages"
	"github.com/almerlucke/metabible/pkg/store"
	"net/http"
)

type PersonPageController struct {
	Store *store.Store
}

func (pc *PersonPageController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	personLookup := r.PathValue("personLookup")

	person, err := pc.Store.GetPersonByLookup(ctx, personLookup)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if person == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	pages.PersonPage(person).Render(r.Context(), w)
}
