package controllers

import (
	"github.com/almerlucke/metabible/pkg/models"
	"github.com/almerlucke/metabible/pkg/store"
	"github.com/almerlucke/metabible/pkg/templates"
	"net/http"
)

type VersePageController struct {
	Store *store.Store
}

func (vc *VersePageController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		verse *models.Verse
		err   error
		ref   string
		ctx   = r.Context()
	)

	ref = r.PathValue("reference")

	if r.URL.Query().Get("osisRef") == "true" {
		verse, err = vc.Store.GetVerseByOsisRef(ctx, ref)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if verse == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	templates.VersePage(verse).Render(r.Context(), w)
}
