package controllers

import (
	"github.com/almerlucke/metabible/pkg/pages"
	"github.com/almerlucke/metabible/pkg/store"
	"net/http"
)

type ChapterPageController struct {
	Store *store.Store
}

func (cc *ChapterPageController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	chapter, err := cc.Store.GetChapterByID(ctx, r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if chapter == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	verses, err := cc.Store.GetVersesForChapter(ctx, chapter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	book, err := cc.Store.GetBookByID(ctx, chapter.Book)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	if book == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var prevChapter string
	var nextChapter string

	if chapter.ChapterNum < book.ChapterCount {
		nextChapter = book.Chapters[chapter.ChapterNum]
	}

	if chapter.ChapterNum > 1 {
		prevChapter = book.Chapters[chapter.ChapterNum-2]
	}

	pages.ChapterPage(book, chapter, verses, prevChapter, nextChapter).Render(ctx, w)
}
