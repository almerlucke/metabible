package atlas

import (
	"encoding/json"
	"github.com/almerlucke/metabible/pkg/i18n"
	"html"
	"os"
	"path/filepath"
	"strings"
)

// Extra Bible translations live in data/bible/<lang>.json as {"01001001": "text", ...}
// keyed by the verseId of the source data. A verse a translation lacks (versification
// differs a little between traditions) simply falls back to the King James text.

// loadBibles reads every data/bible/*.json next to the data directory.
func (a *Atlas) loadBibles(dir string) {
	files, _ := filepath.Glob(filepath.Join(dir, "..", "bible", "*.json"))
	a.bibles = map[string]map[string]string{}
	a.lowerBible = map[string][]string{}
	for _, f := range files {
		lang := strings.TrimSuffix(filepath.Base(f), ".json")
		raw, err := os.ReadFile(f)
		if err != nil {
			continue
		}
		m := map[string]string{}
		if json.Unmarshal(raw, &m) != nil || len(m) == 0 {
			continue
		}
		a.bibles[lang] = m
		low := make([]string, len(a.Verses))
		for i, v := range a.Verses {
			low[i] = strings.ToLower(a.VerseIn(v, lang))
		}
		a.lowerBible[lang] = low
	}
}

// VerseIn is the text of a verse in the given language (King James when there is none).
func (a *Atlas) VerseIn(v *Verse, lang string) string {
	if t, ok := a.bibles[lang][v.VID]; ok && t != "" {
		return t
	}
	return v.Text
}

// HasBible reports whether a full translation exists for the language.
func (a *Atlas) HasBible(lang string) bool { return len(a.bibles[lang]) > 0 }

// lowerFor is the lowercase verse index used for searching in the language.
func (a *Atlas) lowerFor(lang string) []string {
	if l, ok := a.lowerBible[lang]; ok {
		return l
	}
	return a.lowerText
}

// richIn is the verse as safe HTML. English keeps the inline person/place links of the
// source; a translation is plain text (the people and places are listed beside it).
func (a *Atlas) richIn(v *Verse, lang string) string {
	if t, ok := a.bibles[lang][v.VID]; ok && t != "" {
		return html.EscapeString(t)
	}
	return RichInline(mdOf(v))
}

// verseText is the text of the verse an event, journey stop or story step points at.
func (a *Atlas) verseText(id, lang, fallback string) string {
	if v := a.VerseByID[id]; v != nil {
		return a.VerseIn(v, lang)
	}
	return fallback
}

// themeMark shows a theme's verse. The stems that found it are English, so only the
// King James text is highlighted.
func (a *Atlas) themeMark(v *Verse, stems []string, lang string) string {
	if !a.HasBible(lang) {
		return highlightStems(v.Text, stems)
	}
	return html.EscapeString(a.VerseIn(v, lang))
}

// BookName is the name of a book in the language (the source name when there is none).
func (a *Atlas) BookName(b *Book, lang string) string {
	if b == nil {
		return ""
	}
	if lang == i18n.Default || a.tr == nil {
		return b.Name
	}
	return a.trOr(lang, "book."+b.Osis, b.Name)
}

// BooksIn is the list of books with names in the language, for the JSON APIs.
func (a *Atlas) BooksIn(lang string) []*Book {
	if lang == i18n.Default || a.tr == nil {
		return a.Books
	}
	a.bookMu.Lock()
	defer a.bookMu.Unlock()
	if out, ok := a.bookCache[lang]; ok {
		return out
	}
	out := make([]*Book, len(a.Books))
	for i, b := range a.Books {
		c := *b
		c.Name = a.BookName(b, lang)
		out[i] = &c
	}
	if a.bookCache == nil {
		a.bookCache = map[string][]*Book{}
	}
	a.bookCache[lang] = out
	return out
}

// allBookNames lists a book's name in every language, so a reference can be typed in any of them.
func (a *Atlas) allBookNames(b *Book) []string {
	names := []string{b.Name}
	for _, lang := range a.bookLangs {
		if n := a.BookName(b, lang); n != b.Name {
			names = append(names, n)
		}
	}
	return names
}
