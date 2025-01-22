package controllers

import (
	"github.com/almerlucke/metabible/pkg/pages"
	"github.com/almerlucke/metabible/pkg/store"
	"net/http"
)

type VersePageController struct {
	Store *store.Store
}

func (vc *VersePageController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	verse, err := vc.Store.GetVerseByID(ctx, r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if verse == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	chapter, err := vc.Store.GetChapterByID(ctx, verse.Chapter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if chapter == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	book, err := vc.Store.GetBookByID(ctx, verse.Book)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if book == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var prevVerse string
	var nextVerse string

	if verse.VerseNum < len(chapter.Verses) {
		nextVerse = chapter.Verses[verse.VerseNum]
	}

	if verse.VerseNum > 1 {
		prevVerse = chapter.Verses[verse.VerseNum-2]
	}

	pages.VersePage(book, chapter, verse, prevVerse, nextVerse).Render(ctx, w)
}
