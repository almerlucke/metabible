package atlas

import (
	"encoding/json"
	"github.com/almerlucke/metabible/pkg/i18n"
	"html"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
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
// source; a translation gets its links by name (linkNames).
func (a *Atlas) richIn(v *Verse, lang string) string {
	if t, ok := a.bibles[lang][v.VID]; ok && t != "" {
		return a.linkNames(t, v, lang)
	}
	return RichInline(mdOf(v))
}

var nameNote = regexp.MustCompile(`\s*\([^)]*\)\s*$`)

type nameLink struct{ name, class, href string }

// linkNames escapes the text of a translated verse and links the people and places the
// verse is about wherever their name in the language occurs (whole words, an -s genitive
// included). The source marks these in the English text only, so the names are found by
// their translated name, restricted to the ones the verse lists.
func (a *Atlas) linkNames(text string, v *Verse, lang string) string {
	nt := a.names[lang]
	if nt == nil || (len(v.People) == 0 && len(v.Places) == 0) {
		return html.EscapeString(text)
	}
	var cands []nameLink
	add := func(class, kind, id, name string) {
		if name = strings.TrimSpace(nameNote.ReplaceAllString(name, "")); name != "" {
			cands = append(cands, nameLink{name, class, "/" + kind + "/" + id})
		}
	}
	for _, id := range v.People {
		add("k-person", "person", id, nt.People[id])
	}
	for _, id := range v.Places {
		add("k-place", "place", id, nt.Places[id])
	}
	// the longest name first, so "Jeruzalem" wins over a shorter name inside it
	sort.SliceStable(cands, func(i, j int) bool { return len(cands[i].name) > len(cands[j].name) })
	type span struct {
		from, to int
		c        *nameLink
	}
	var spans []span
	taken := func(f, t int) bool {
		for _, s := range spans {
			if f < s.to && t > s.from {
				return true
			}
		}
		return false
	}
	matched := make([]bool, len(cands))
	for i := range cands {
		c := &cands[i]
		for from := 0; ; {
			k := strings.Index(text[from:], c.name)
			if k < 0 {
				break
			}
			f := from + k
			t := f + len(c.name)
			from = t
			if r, _ := utf8.DecodeLastRuneInString(text[:f]); f > 0 && (unicode.IsLetter(r) || unicode.IsDigit(r)) {
				continue
			}
			if t < len(text) && text[t] == 's' && (t+1 == len(text) || !isLetterAt(text, t+1)) {
				t++
			}
			if t < len(text) && isLetterAt(text, t) {
				continue
			}
			matched[i] = true
			if taken(f, t) {
				continue
			}
			spans = append(spans, span{f, t, c})
		}
	}
	// The source text spells many names differently from the modern table (Kanaan/Kanaän,
	// Kohath/Kahath, Elzafan/Elizafan): for a name not found, link the capitalised word of
	// the verse that looks most like it.
	var words []span
	for i := 0; i < len(text); {
		r, n := utf8.DecodeRuneInString(text[i:])
		if !unicode.IsLetter(r) {
			i += n
			continue
		}
		j := i
		for j < len(text) {
			r2, n2 := utf8.DecodeRuneInString(text[j:])
			if !unicode.IsLetter(r2) {
				break
			}
			j += n2
		}
		if unicode.IsUpper(r) {
			words = append(words, span{from: i, to: j})
		}
		i = j
	}
	for i := range cands {
		if matched[i] {
			continue
		}
		c := &cands[i]
		head := foldName(strings.Fields(c.name)[0])
		if len(head) < 3 {
			continue
		}
		best, score := -1, 0.0
		for k, w := range words {
			if taken(w.from, w.to) {
				continue
			}
			x := foldName(text[w.from:w.to])
			sim := nameSimilarity(head, x)
			if (x[0] == head[0] && sim >= 0.72) || sim >= 0.85 {
				if sim > score {
					best, score = k, sim
				}
			}
		}
		if best >= 0 {
			spans = append(spans, span{words[best].from, words[best].to, c})
		}
	}
	sort.Slice(spans, func(i, j int) bool { return spans[i].from < spans[j].from })
	var b strings.Builder
	last := 0
	for _, s := range spans {
		b.WriteString(html.EscapeString(text[last:s.from]))
		b.WriteString(`<a class="` + s.c.class + `" href="` + s.c.href + `">` + html.EscapeString(text[s.from:s.to]) + `</a>`)
		last = s.to
	}
	b.WriteString(html.EscapeString(text[last:]))
	return b.String()
}

func isLetterAt(s string, i int) bool {
	r, _ := utf8.DecodeRuneInString(s[i:])
	return unicode.IsLetter(r)
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

var nameFolder = strings.NewReplacer("ë", "e", "ï", "i", "ö", "o", "ü", "u", "é", "e", "è", "e", "ê", "e", "á", "a", "à", "a",
	"c", "k", "z", "s", "ph", "f", "th", "t", "ch", "k", "y", "j", "w", "v")

// foldName reduces a name to a rough sound, so spelling variants of one name come close.
func foldName(s string) string {
	s = nameFolder.Replace(strings.ToLower(s))
	var b []rune
	for _, r := range s {
		if len(b) == 0 || b[len(b)-1] != r {
			b = append(b, r)
		}
	}
	return string(b)
}

// nameSimilarity is 1 for equal names, lower the more edits separate them.
func nameSimilarity(a, b string) float64 {
	x, y := []rune(a), []rune(b)
	if len(x) == 0 || len(y) == 0 {
		return 0
	}
	prev := make([]int, len(y)+1)
	for j := range prev {
		prev[j] = j
	}
	for i := 1; i <= len(x); i++ {
		cur := make([]int, len(y)+1)
		cur[0] = i
		for j := 1; j <= len(y); j++ {
			c := 1
			if x[i-1] == y[j-1] {
				c = 0
			}
			cur[j] = min(prev[j]+1, cur[j-1]+1, prev[j-1]+c)
		}
		prev = cur
	}
	m := max(len(x), len(y))
	return 1 - float64(prev[len(y)])/float64(m)
}
