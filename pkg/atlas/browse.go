package atlas

import (
	"html"
	"strings"
)

// ---------------------------------------------------------------- books

type BookCard struct {
	ID       string
	Name     string
	Chapters int
	Verses   int
	People   int
	Places   int
	Heat     []float64
	CIDs     []string
}

type BookGroup struct {
	Testament string
	Div       string
	Books     []BookCard
}

func (a *Atlas) BooksBrowse() []BookGroup {
	var out []BookGroup
	for _, b := range a.Books {
		if len(out) == 0 || out[len(out)-1].Div != b.Div {
			out = append(out, BookGroup{Testament: b.Testament, Div: b.Div})
		}
		g := &out[len(out)-1]
		maxPe := 1
		for _, c := range b.Chapters {
			if c.People > maxPe {
				maxPe = c.People
			}
		}
		pe, pl := map[string]bool{}, map[string]bool{}
		bc := BookCard{ID: b.ID, Name: b.Name, Chapters: len(b.Chapters), Verses: b.Verses}
		for _, c := range b.Chapters {
			bc.Heat = append(bc.Heat, float64(c.People)/float64(maxPe))
			bc.CIDs = append(bc.CIDs, c.ID)
			for _, v := range a.chapterVerses[c.ID] {
				for _, p := range v.People {
					pe[p] = true
				}
				for _, p := range v.Places {
					pl[p] = true
				}
			}
		}
		bc.People, bc.Places = len(pe), len(pl)
		g.Books = append(g.Books, bc)
	}
	return out
}

// ---------------------------------------------------------------- events

type EventCard struct {
	ID      string
	Title   string
	Date    string
	Year    float64
	Text    string
	Ref     string
	VerseID string
	Places  []Ref
	People  []Ref
	Days    int
}

type EraGroup struct {
	Era    Era
	Events []EventCard
}

func (a *Atlas) EventsTimeline() []EraGroup {
	var out []EraGroup
	idx := map[string]int{}
	for _, e := range Eras {
		idx[e.ID] = len(out)
		out = append(out, EraGroup{Era: e})
	}
	for _, e := range a.Events {
		era := eraAt(int(e.Year))
		pe := e.People
		if len(pe) > 4 {
			pe = pe[:4]
		}
		pl := e.Places
		if len(pl) > 3 {
			pl = pl[:3]
		}
		c := EventCard{ID: e.ID, Title: e.Title, Date: e.Date, Year: e.Year, Text: trunc(e.Text, 150), Ref: e.Ref, VerseID: e.VerseID,
			Places: a.placeRefs(pl), People: a.personRefs(pe), Days: e.Days}
		g := &out[idx[era.ID]]
		g.Events = append(g.Events, c)
	}
	var res []EraGroup
	for _, g := range out {
		if len(g.Events) > 0 {
			res = append(res, g)
		}
	}
	return res
}

// ---------------------------------------------------------------- concordance

type ConcHit struct {
	ID    string
	Label string
	HTML  string
}

type ConcBook struct {
	ID    string
	Name  string
	Count int
	Pct   float64
}

type ConcResult struct {
	Query  string
	Terms  []string
	Total  int
	Hits   []ConcHit
	Shown  int
	Books  []ConcBook
	BookID string
	Capped bool
}

func parseTerms(q string) []string {
	q = strings.ToLower(strings.TrimSpace(q))
	var terms []string
	for {
		i := strings.Index(q, `"`)
		if i < 0 {
			break
		}
		j := strings.Index(q[i+1:], `"`)
		if j < 0 {
			break
		}
		if phrase := strings.TrimSpace(q[i+1 : i+1+j]); phrase != "" {
			terms = append(terms, phrase)
		}
		q = q[:i] + " " + q[i+j+2:]
	}
	terms = append(terms, strings.Fields(q)...)
	return terms
}

func isLetter(b byte) bool {
	return b >= 'a' && b <= 'z' || b >= '0' && b <= '9' || b >= 0x80
}

// indexWord finds term in lower starting at from, only where it begins a word.
func indexWord(lower, term string, from int) int {
	for from <= len(lower) {
		i := strings.Index(lower[from:], term)
		if i < 0 {
			return -1
		}
		i += from
		if i == 0 || !isLetter(lower[i-1]) {
			return i
		}
		from = i + 1
	}
	return -1
}

func highlight(text string, terms []string) string {
	lower := strings.ToLower(text)
	if len(lower) != len(text) {
		return html.EscapeString(text)
	}
	mark := make([]bool, len(text))
	for _, t := range terms {
		for from := 0; ; {
			i := indexWord(lower, t, from)
			if i < 0 {
				break
			}
			end := i + len(t)
			for end < len(lower) && isLetter(lower[end]) { // highlight the whole word: love -> loved
				end++
			}
			for k := i; k < end; k++ {
				mark[k] = true
			}
			from = end
		}
	}
	var b strings.Builder
	for i := 0; i < len(text); {
		j := i
		for j < len(text) && mark[j] == mark[i] {
			j++
		}
		seg := html.EscapeString(text[i:j])
		if mark[i] {
			b.WriteString("<mark>" + seg + "</mark>")
		} else {
			b.WriteString(seg)
		}
		i = j
	}
	return b.String()
}

// Concordance finds verses containing every term (a quoted phrase counts as one term).
func (a *Atlas) Concordance(q, bookID string, limit int) *ConcResult {
	terms := parseTerms(q)
	res := &ConcResult{Query: q, Terms: terms, BookID: bookID}
	if len(terms) == 0 || len(strings.Join(terms, "")) < 2 {
		return res
	}
	counts := map[string]int{}
	for i, t := range a.lowerText {
		ok := true
		for _, term := range terms {
			if indexWord(t, term, 0) < 0 {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}
		v := a.Verses[i]
		counts[v.Book]++
		res.Total++
		if bookID != "" && v.Book != bookID {
			continue
		}
		if len(res.Hits) < limit {
			res.Hits = append(res.Hits, ConcHit{ID: v.ID, Label: a.verseLabel(v), HTML: highlight(v.Text, terms)})
		} else {
			res.Capped = true
		}
	}
	max := 1
	for _, n := range counts {
		if n > max {
			max = n
		}
	}
	for _, b := range a.Books {
		res.Books = append(res.Books, ConcBook{ID: b.ID, Name: b.Name, Count: counts[b.ID], Pct: float64(counts[b.ID]) / float64(max)})
	}
	res.Shown = len(res.Hits)
	return res
}
