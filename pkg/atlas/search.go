package atlas

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

type searchEntry struct {
	Type  string
	ID    string
	Label string
	Aux   string // the year, date, kind or testament that follows the type in the result line
	HasY  bool
	Y     int
	URL   string
	Lower string
	W     int
}

type SearchHit struct {
	Type  string `json:"t"`
	ID    string `json:"id"`
	Label string `json:"l"`
	Sub   string `json:"s,omitempty"`
	URL   string `json:"u"`
}

func (a *Atlas) buildSearch() {
	for _, p := range a.People {
		l := strings.ToLower(p.Name)
		if p.AlsoCall != "" {
			l += " " + strings.ToLower(p.AlsoCall)
		}
		a.nameIdx = append(a.nameIdx, searchEntry{Type: "person", ID: p.ID, Label: p.Name, HasY: p.Y0 != 0 || p.Y1 != 0, Y: p.Y0, URL: "/person/" + p.ID, Lower: l, W: p.Verses})
	}
	for _, p := range a.Places {
		a.nameIdx = append(a.nameIdx, searchEntry{Type: "place", ID: p.ID, Label: p.Name, Aux: p.Kind, URL: "/place/" + p.ID, Lower: strings.ToLower(p.Name), W: p.Verses})
	}
	for _, e := range a.Events {
		a.nameIdx = append(a.nameIdx, searchEntry{Type: "event", ID: e.ID, Label: e.Title, Aux: e.Date, URL: "/atlas?event=" + e.ID, Lower: strings.ToLower(e.Title), W: 20})
	}
	for _, b := range a.Books {
		a.nameIdx = append(a.nameIdx, searchEntry{Type: "book", ID: b.ID, Label: b.Name, Aux: b.Testament, URL: "/book/" + b.ID, Lower: strings.ToLower(b.Name), W: 100000})
	}
}

// subFor writes the grey line after a search result ("person · 2056 BC") in a language.
func (a *Atlas) subFor(lang string, e searchEntry) string {
	switch e.Type {
	case "person":
		if e.HasY {
			return a.tx(lang, "sub.person.year", "y", a.formatYear(lang, e.Y))
		}
		return a.tx(lang, "sub.person")
	case "place":
		if e.Aux != "" {
			return a.kindWords(lang, strings.ToLower(e.Aux))
		}
		return a.tx(lang, "sub.place")
	case "event":
		return a.tx(lang, "sub.event", "d", a.locDate(lang, e.Aux))
	case "book":
		return a.tx(lang, "sub.book", "t", a.lookup(lang, "testament", e.Aux))
	}
	return e.Type
}

var refRe = regexp.MustCompile(`^\s*([1-3]?\s*[a-zA-Z]+(?:\s+of\s+[a-zA-Z]+)?)\s+(\d+)(?:\s*[:.]\s*(\d+))?\s*$`)

func (a *Atlas) findBook(name string) *Book {
	n := strings.ToLower(strings.ReplaceAll(name, " ", ""))
	if n == "" {
		return nil
	}
	var best *Book
	for _, b := range a.Books {
		if strings.ToLower(b.Short) == n || strings.ToLower(b.Osis) == n {
			return b
		}
		for _, name := range a.allBookNames(b) {
			bn := strings.ToLower(strings.ReplaceAll(name, " ", ""))
			if bn == n {
				return b
			}
			if best == nil && strings.HasPrefix(bn, n) {
				best = b
			}
		}
	}
	return best
}

func (a *Atlas) Search(q string, limit int, lang string) []SearchHit {
	q = strings.TrimSpace(q)
	if q == "" {
		return nil
	}
	var hits []SearchHit

	if m := refRe.FindStringSubmatch(q); m != nil {
		if b := a.findBook(m[1]); b != nil {
			ch, _ := strconv.Atoi(m[2])
			if m[3] != "" {
				vn, _ := strconv.Atoi(m[3])
				osis := fmt.Sprintf("%s.%d.%d", b.Osis, ch, vn)
				if v, ok := a.VerseByOsis[strings.ToLower(osis)]; ok {
					hits = append(hits, SearchHit{"verse", v.ID, fmt.Sprintf("%s %d:%d", a.BookName(b, lang), ch, vn), trunc(a.VerseIn(v, lang), 80), "/verse/" + v.ID})
				}
			} else {
				for _, c := range b.Chapters {
					if c.Num == ch {
						hits = append(hits, SearchHit{"chapter", c.ID, fmt.Sprintf("%s %d", a.BookName(b, lang), ch), a.tx(lang, "sub.chapter"), "/chapter/" + c.ID})
					}
				}
			}
		}
	}

	lq := strings.ToLower(q)
	type scored struct {
		e     searchEntry
		score int
	}
	var sc []scored
	for _, e := range a.nameIdx {
		idx := strings.Index(e.Lower, lq)
		if idx < 0 {
			continue
		}
		s := 1
		if idx == 0 {
			s = 3
		} else if idx > 0 && !unicode.IsLetter(rune(e.Lower[idx-1])) {
			s = 2
		}
		if e.Lower == lq {
			s = 5
		}
		sc = append(sc, scored{e, s})
	}
	sort.SliceStable(sc, func(i, j int) bool {
		if sc[i].score != sc[j].score {
			return sc[i].score > sc[j].score
		}
		return sc[i].e.W > sc[j].e.W
	})
	for _, s := range sc {
		if len(hits) >= limit {
			break
		}
		label := s.e.Label
		if s.e.Type == "book" {
			label = a.BookName(a.BookByID[s.e.ID], lang)
		}
		hits = append(hits, SearchHit{s.e.Type, s.e.ID, label, a.subFor(lang, s.e), s.e.URL})
	}

	if len(q) >= 4 && len(hits) < limit {
		n := 0
		for i, t := range a.lowerFor(lang) {
			if strings.Contains(t, lq) {
				v := a.Verses[i]
				b := a.BookByID[v.Book]
				bn := ""
				if b != nil {
					bn = a.BookName(b, lang)
				}
				hits = append(hits, SearchHit{"verse", v.ID, fmt.Sprintf("%s %d:%d", bn, a.ChapterNum[v.Chapter], v.Num), trunc(a.VerseIn(v, lang), 90), "/verse/" + v.ID})
				n++
				if n >= 4 || len(hits) >= limit+4 {
					break
				}
			}
		}
	}
	return hits
}

func trunc(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n]) + "…"
}

var mdLink = regexp.MustCompile(`\[([^\]]*)\]\([^)]*\)`)

func plain(s string) string {
	s = mdLink.ReplaceAllString(s, "$1")
	s = strings.ReplaceAll(s, "\n", " ")
	return strings.Join(strings.Fields(s), " ")
}

// Peek is the small card shown when hovering a link.
type Peek struct {
	Type  string      `json:"t"`
	Title string      `json:"title"`
	Sub   string      `json:"sub,omitempty"`
	Text  string      `json:"text,omitempty"`
	Facts []string    `json:"facts,omitempty"`
	Links [][2]string `json:"links,omitempty"`
}

func (a *Atlas) verseLabel(v *Verse, lang string) string {
	b := a.BookByID[v.Book]
	if b == nil {
		return v.Osis
	}
	return fmt.Sprintf("%s %d:%d", a.BookName(b, lang), a.ChapterNum[v.Chapter], v.Num)
}

func (a *Atlas) PeekFor(kind, id, lang string) *Peek {
	switch kind {
	case "person":
		p := a.PersonByID[id]
		if p == nil {
			return nil
		}
		pk := &Peek{Type: "person", Title: p.Name, Text: trunc(plain(p.Dict), 260)}
		if p.Gender != "" {
			pk.Sub = a.lookup(lang, "gender", p.Gender)
		}
		if p.Y0 != 0 || p.Y1 != 0 {
			pk.Facts = append(pk.Facts, a.tx(lang, "peek.appears", "a", a.formatYear(lang, p.Y0), "b", a.formatYear(lang, p.Y1)))
		}
		pk.Facts = append(pk.Facts, a.txn(lang, "count.verse", p.Verses))
		if len(p.Children) > 0 {
			pk.Facts = append(pk.Facts, a.txn(lang, "peek.children", len(p.Children)))
		}
		if f := a.PersonByID[p.Father]; f != nil {
			pk.Facts = append(pk.Facts, a.tx(lang, "peek.sonof", "name", f.Name))
		}
		pk.Links = [][2]string{{a.tx(lang, "peek.link.journey"), "/journeys?person=" + p.ID}, {a.tx(lang, "peek.link.tree"), "/tree/" + p.ID}}
		return pk
	case "place":
		p := a.PlaceByID[id]
		if p == nil {
			return nil
		}
		pk := &Peek{Type: "place", Title: p.Name, Sub: a.kindWords(lang, strings.TrimSpace(p.Kind+" "+p.Sub)), Text: trunc(plain(p.Dict), 260)}
		if pk.Text == "" {
			pk.Text = p.Comment
		}
		pk.Facts = append(pk.Facts, a.txn(lang, "count.verse", p.Verses))
		if p.HasGeo {
			pk.Facts = append(pk.Facts, fmt.Sprintf("%.2f°, %.2f°", p.Lat, p.Lng))
			pk.Links = [][2]string{{a.tx(lang, "peek.link.atlas"), fmt.Sprintf("/atlas?place=%s", p.ID)}}
		}
		return pk
	case "verse":
		v := a.VerseByID[id]
		if v == nil {
			return nil
		}
		pk := &Peek{Type: "verse", Title: a.verseLabel(v, lang), Text: a.VerseIn(v, lang)}
		pk.Facts = append(pk.Facts, a.tx(lang, "peek.verse.counts", "p", a.txn(lang, "count.person", len(v.People)), "l", a.txn(lang, "count.place", len(v.Places))))
		return pk
	case "book":
		b := a.BookByID[id]
		if b == nil {
			return nil
		}
		return &Peek{Type: "book", Title: a.BookName(b, lang), Sub: a.lookup(lang, "div", b.Div) + " · " + a.lookup(lang, "testament", b.Testament),
			Facts: []string{a.tx(lang, "peek.book.counts", "c", a.txn(lang, "count.chapter", len(b.Chapters)), "v", a.txn(lang, "count.verse", b.Verses))}}
	case "chapter":
		for _, b := range a.Books {
			for _, c := range b.Chapters {
				if c.ID == id {
					return &Peek{Type: "chapter", Title: fmt.Sprintf("%s %d", a.BookName(b, lang), c.Num),
						Facts: []string{a.tx(lang, "peek.chapter.counts", "v", a.txn(lang, "count.verse", c.Verses), "p", a.txn(lang, "count.person", c.People), "l", a.txn(lang, "count.place", c.Places))}}
				}
			}
		}
	case "event":
		e := a.EventByID[id]
		if e == nil {
			return nil
		}
		return &Peek{Type: "event", Title: e.Title, Sub: a.locDate(lang, e.Date), Text: a.verseText(e.VerseID, lang, e.Text), Links: [][2]string{{a.tx(lang, "peek.link.atlas"), "/atlas?event=" + e.ID}}}
	}
	return nil
}
