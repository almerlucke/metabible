package controllers

import (
	"github.com/almerlucke/metabible/pkg/pages"
	"github.com/almerlucke/metabible/pkg/store"
	"net/http"
)

type BookPageController struct {
	Store *store.Store
}

func (bc *BookPageController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	book, err := bc.Store.GetBookByID(ctx, r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if book == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	books, err := bc.Store.GetAllBooks(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	chapters, err := bc.Store.GetChaptersForBook(ctx, book)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	excerpts := make([]string, len(chapters))

	for i, chapter := range chapters {
		verse, err := bc.Store.GetVerseByID(ctx, chapter.Verses[0])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		if verse == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		excerpts[i] = verse.ShortExcerpt()
	}

	var prevBook string
	var nextBook string

	if book.BookOrder < len(books) {
		nextBook = books[book.BookOrder].ID
	}

	if book.BookOrder > 1 {
		prevBook = books[book.BookOrder-2].ID
	}

	pages.BookPage(book, chapters, excerpts, prevBook, nextBook).Render(ctx, w)
}
