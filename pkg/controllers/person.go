package controllers

import (
	"github.com/almerlucke/metabible/pkg/models"
	"github.com/almerlucke/metabible/pkg/pages"
	"github.com/almerlucke/metabible/pkg/store"
	"net/http"
	"slices"
)

type PersonPageController struct {
	Store *store.Store
}

func (pc *PersonPageController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	person, err := pc.Store.GetPersonByID(ctx, r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if person == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var info pages.PersonInfo

	info.Person = person

	if len(person.BirthPlace) > 0 {
		if info.BirthPlace, err = pc.Store.GetPlaceByID(ctx, person.BirthPlace[0]); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if len(person.BirthYear) > 0 {
		if info.BirthYear, err = pc.Store.GetPeriodByID(ctx, person.BirthYear[0]); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if len(person.DeathPlace) > 0 {
		if info.DeathPlace, err = pc.Store.GetPlaceByID(ctx, person.DeathPlace[0]); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if len(person.DeathYear) > 0 {
		if info.DeathYear, err = pc.Store.GetPeriodByID(ctx, person.DeathYear[0]); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if info.Siblings, err = pc.Store.GetPeopleInList(ctx, person.Siblings); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if info.Children, err = pc.Store.GetPeopleInList(ctx, person.Children); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if info.Partners, err = pc.Store.GetPeopleInList(ctx, person.Partners); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if info.Father, err = pc.Store.GetPersonByID(ctx, person.Father); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if info.Mother, err = pc.Store.GetPersonByID(ctx, person.Mother); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if info.Timeline, err = pc.Store.GetEventsInList(ctx, person.Timeline); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if verses, err := pc.Store.GetVersesInList(ctx, person.Verses); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		items := make([]*models.ConcordanceItem, len(verses))

		for i, verse := range verses {
			book, err := pc.Store.GetBookByID(ctx, verse.Book)
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

		slices.SortFunc(items, func(a *models.ConcordanceItem, b *models.ConcordanceItem) int {
			return a.SortIndex() - b.SortIndex()
		})

		info.Verses = items
	}

	pages.PersonPage(&info).Render(ctx, w)
}
