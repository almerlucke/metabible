package controllers

import (
	"github.com/almerlucke/metabible/pkg/pages"
	"github.com/almerlucke/metabible/pkg/store"
	"net/http"
)

type BooksPageController struct {
	Store *store.Store
}

func (bc *BooksPageController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	books, err := bc.Store.GetAllBooks(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pages.BooksPage(books).Render(ctx, w)
}
