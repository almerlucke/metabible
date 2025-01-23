package controllers

import (
	"github.com/almerlucke/metabible/pkg/models"
	"github.com/almerlucke/metabible/pkg/pages"
	"github.com/almerlucke/metabible/pkg/store"
	"net/http"
	"slices"
)

type ConcordancePageController struct {
	Store *store.Store
}

func (cc *ConcordancePageController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		verses []*models.Verse
		items  []*models.ConcordanceItem
		err    error
	)

	ctx := r.Context()
	search := r.URL.Query().Get("search")

	if len(search) > 2 {
		verses, err = cc.Store.VerseSearch(ctx, search)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		items = make([]*models.ConcordanceItem, len(verses))

		for i, verse := range verses {
			book, err := cc.Store.GetBookByID(ctx, verse.Book)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			items[i] = &models.ConcordanceItem{
				Verse:      verse,
				Book:       book,
				ChapterNum: verse.ChapterNum(),
			}
		}
	}

	slices.SortFunc(items, func(a *models.ConcordanceItem, b *models.ConcordanceItem) int {
		return a.SortIndex() - b.SortIndex()
	})

	pages.ConcordancePage(items, search).Render(ctx, w)
}
