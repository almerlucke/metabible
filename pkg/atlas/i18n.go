package atlas

import (
	"regexp"
	"strings"

	"github.com/almerlucke/metabible/pkg/i18n"
)

// Translator is what the data layer needs from the translation bundle to write sentences in a language.
type Translator interface {
	T(lang, key string, kv ...any) string
	Plural(lang, key string, n int, kv ...any) string
	Raw(lang, key string) (string, bool)
	Codes() []string
}

// SetTranslations lets the data layer build texts (relationship stories, era introductions…) in the visitor's language.
func (a *Atlas) SetTranslations(t Translator) {
	a.tr = t
	a.bookLangs = nil
	for _, c := range t.Codes() {
		if c != i18n.Default {
			a.bookLangs = append(a.bookLangs, c)
		}
	}
	a.nameIdx = nil
	a.buildSearch() // books are found by their name in any language
}

// tx translates a key; without a translator (tests of the bare data layer) the key is returned.
func (a *Atlas) tx(lang, key string, kv ...any) string {
	if a.tr == nil {
		return key
	}
	return a.tr.T(lang, key, kv...)
}

func (a *Atlas) txn(lang, key string, n int, kv ...any) string {
	if a.tr == nil {
		return key
	}
	return a.tr.Plural(lang, key, n, kv...)
}

// ---------------------------------------------------------------- language helpers

var (
	bcRe = regexp.MustCompile(`(\d+) BC\b`)
	adRe = regexp.MustCompile(`(\d+) AD\b`)
)

// locDate turns the English dates built into the data ("1446 BC") into the language's form.
func (a *Atlas) locDate(lang, s string) string {
	if a.tr == nil || lang == i18n.Default || s == "" {
		return s
	}
	bc, ad := a.tx(lang, "date.bc"), a.tx(lang, "date.ad")
	s = bcRe.ReplaceAllString(s, "${1} "+strings.ReplaceAll(bc, "$", "$$"))
	return adRe.ReplaceAllString(s, "${1} "+strings.ReplaceAll(ad, "$", "$$"))
}

// formatYear is FormatYear in a language.
func (a *Atlas) formatYear(lang string, y int) string { return a.locDate(lang, FormatYear(y)) }

// LocalEras are the eras with their names in a language.
func (a *Atlas) LocalEras(lang string) []Era {
	out := make([]Era, len(Eras))
	for i, e := range Eras {
		out[i] = e
		if a.tr != nil {
			if s, ok := a.tr.Raw(lang, "era."+e.ID+".name"); ok {
				out[i].Name = s
			}
		}
	}
	return out
}

func slugKey(s string) string {
	return strings.ReplaceAll(strings.ToLower(strings.TrimSpace(s)), " ", "-")
}

// lookup translates a data value through keys like gender.male; unknown values stay as they are.
func (a *Atlas) lookup(lang, group, value string) string {
	if a.tr == nil || value == "" {
		return value
	}
	if s, ok := a.tr.Raw(lang, group+"."+slugKey(value)); ok {
		return s
	}
	return value
}

func (a *Atlas) kindWords(lang, k string) string {
	words := strings.Fields(k)
	for i, w := range words {
		words[i] = a.lookup(lang, "kind", w)
	}
	return strings.Join(words, " ")
}

func (a *Atlas) joinAnd(lang string, s []string) string {
	switch len(s) {
	case 0:
		return ""
	case 1:
		return s[0]
	}
	return strings.Join(s[:len(s)-1], ", ") + " " + a.tx(lang, "word.and") + " " + s[len(s)-1]
}
