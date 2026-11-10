package atlas

import (
	"encoding/json"
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
