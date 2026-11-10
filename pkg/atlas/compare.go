package atlas

import (
	"fmt"
	"sort"
	"strings"
)

// ---------------------------------------------------------------- family graph

type famEdge struct {
	To   string
	Kind string // "up" (to a parent), "down" (to a child), "sp" (to a spouse)
}

func (a *Atlas) family() map[string][]famEdge {
	a.famOnce.Do(func() {
		g := map[string][]famEdge{}
		seen := map[[3]string]bool{}
		add := func(from, to, kind string) {
			if a.PersonByID[from] == nil || a.PersonByID[to] == nil {
				return
			}
			k := [3]string{from, to, kind}
			if !seen[k] {
				seen[k] = true
				g[from] = append(g[from], famEdge{to, kind})
			}
		}
		for _, p := range a.People {
			for _, par := range []string{p.Father, p.Mother} {
				if par != "" {
					add(p.ID, par, "up")
					add(par, p.ID, "down")
				}
			}
			for _, c := range p.Children {
				add(p.ID, c, "down")
				add(c, p.ID, "up")
			}
			for _, s := range p.Partners {
				add(p.ID, s, "sp")
				add(s, p.ID, "sp")
			}
		}
		a.fam = g
	})
	return a.fam
}

type RelStep struct {
	Person Ref    `json:"p"`
	Via    string `json:"via,omitempty"`
}

type Relation struct {
	Steps  []RelStep `json:"steps"`
	Label  string    `json:"label"`  // what B is to A
	Phrase string    `json:"phrase"` // a sentence
}

// FamilyPath finds the shortest family connection between two people.
func (a *Atlas) FamilyPath(from, to string) *Relation {
	if from == to || a.PersonByID[from] == nil || a.PersonByID[to] == nil {
		return nil
	}
	g := a.family()
	type prev struct {
		id   string
		kind string
	}
	pv := map[string]prev{from: {}}
	q := []string{from}
	for len(q) > 0 {
		id := q[0]
		q = q[1:]
		if id == to {
			break
		}
		for _, e := range g[id] {
			if _, ok := pv[e.To]; !ok {
				pv[e.To] = prev{id, e.Kind}
				q = append(q, e.To)
			}
		}
	}
	if _, ok := pv[to]; !ok {
		return nil
	}
	var rev []RelStep
	for cur := to; cur != from; cur = pv[cur].id {
		rev = append(rev, RelStep{Person: *a.personRef(cur), Via: pv[cur].kind})
	}
	rev = append(rev, RelStep{Person: *a.personRef(from)})
	r := &Relation{}
	for i := len(rev) - 1; i >= 0; i-- {
		r.Steps = append(r.Steps, rev[i])
	}
	ups, downs, sp := 0, 0, 0
	order := ""
	for _, s := range r.Steps[1:] {
		switch s.Via {
		case "up":
			ups++
			order += "u"
		case "down":
			downs++
			order += "d"
		case "sp":
			sp++
			order += "s"
		}
	}
	A, B := a.PersonByID[from], a.PersonByID[to]
	r.Label = relationLabel(order, ups, downs, sp, B.Gender)
	switch {
	case r.Label == "":
		r.Phrase = fmt.Sprintf("%s and %s are connected through %d steps of family.", A.Name, B.Name, len(r.Steps)-1)
	default:
		r.Phrase = fmt.Sprintf("%s is %s's %s.", B.Name, A.Name, r.Label)
	}
	return r
}

func greatPrefix(n int) string {
	switch {
	case n <= 0:
		return ""
	case n <= 3:
		return strings.Repeat("great-", n)
	}
	return fmt.Sprintf("%s-great-", ordinal(n))
}

func ordinal(n int) string {
	suf := "th"
	if n%100 < 11 || n%100 > 13 {
		switch n % 10 {
		case 1:
			suf = "st"
		case 2:
			suf = "nd"
		case 3:
			suf = "rd"
		}
	}
	return fmt.Sprintf("%d%s", n, suf)
}

func gendered(g, male, female, neutral string) string {
	switch g {
	case "Male":
		return male
	case "Female":
		return female
	}
	return neutral
}

// relationLabel names what B is to A given the shape of the path: u = steps to a parent, d = to a child, s = marriage.
func relationLabel(order string, ups, downs, sp int, bGender string) string {
	switch {
	case sp == 0 && downs == 0: // ancestor
		switch ups {
		case 1:
			return gendered(bGender, "father", "mother", "parent")
		case 2:
			return gendered(bGender, "grandfather", "grandmother", "grandparent")
		}
		return greatPrefix(ups-2) + gendered(bGender, "grandfather", "grandmother", "grandparent")
	case sp == 0 && ups == 0: // descendant
		switch downs {
		case 1:
			return gendered(bGender, "son", "daughter", "child")
		case 2:
			return gendered(bGender, "grandson", "granddaughter", "grandchild")
		}
		return greatPrefix(downs-2) + gendered(bGender, "grandson", "granddaughter", "grandchild")
	case sp == 0 && ups == 1 && downs == 1:
		return gendered(bGender, "brother", "sister", "sibling")
	case sp == 0 && ups == 1 && downs >= 2: // B descends from A's parent
		return greatPrefix(downs-2) + gendered(bGender, "nephew", "niece", "nephew or niece")
	case sp == 0 && downs == 1 && ups >= 2:
		return greatPrefix(ups-2) + gendered(bGender, "uncle", "aunt", "uncle or aunt")
	case sp == 0 && ups >= 2 && downs >= 2:
		deg, rem := ups, downs
		if downs < ups {
			deg = downs
		}
		deg--
		rem = ups - downs
		if rem < 0 {
			rem = -rem
		}
		s := ordinal(deg) + " cousin"
		if rem > 0 {
			s += fmt.Sprintf(" %d× removed", rem)
		}
		return s
	case sp == 1 && ups+downs == 0:
		return gendered(bGender, "husband", "wife", "spouse")
	}
	return ""
}

// ---------------------------------------------------------------- compare

type CmpSide struct {
	ID       string  `json:"id"`
	Name     string  `json:"n"`
	Kind     string  `json:"kind"`
	From     int     `json:"from"`
	To       int     `json:"to"`
	HasYears bool    `json:"hasYears"`
	Exact    bool    `json:"exact"`
	Verses   int     `json:"vc"`
	Books    []int   `json:"books"` // mentions per book, canonical order
	Lat      float64 `json:"lat,omitempty"`
	Lng      float64 `json:"lng,omitempty"`
	Sub      string  `json:"sub,omitempty"`
}

type CompareResult struct {
	Kind         string     `json:"kind"`
	A            CmpSide    `json:"a"`
	B            CmpSide    `json:"b"`
	SharedVerses []VerseRef `json:"sharedVerses"`
	SharedCount  int        `json:"sharedCount"`
	SharedEvents []EventRef `json:"sharedEvents"`
	SharedPlaces []Ref      `json:"sharedPlaces"`
	SharedPeople []Ref      `json:"sharedPeople"`
	Relation     *Relation  `json:"relation,omitempty"`
	DistanceKm   int        `json:"distanceKm,omitempty"`
	Story        []string   `json:"story"`
	BookNames    []string   `json:"bookNames"`
	QuoteVerse   string     `json:"quoteVerse,omitempty"`
	QuoteText    string     `json:"quoteText,omitempty"`
	QuoteLabel   string     `json:"quoteLabel,omitempty"`
}

func (a *Atlas) bookCounts(verseIDs []string) []int {
	idx := map[string]int{}
	for i, b := range a.Books {
		idx[b.ID] = i
	}
	out := make([]int, len(a.Books))
	for _, id := range verseIDs {
		if v := a.VerseByID[id]; v != nil {
			if i, ok := idx[v.Book]; ok {
				out[i]++
			}
		}
	}
	return out
}

func (a *Atlas) personSide(p *Person) CmpSide {
	s := CmpSide{ID: p.ID, Name: p.Name, Kind: "person", Verses: p.Verses, Books: a.bookCounts(p.Vs), Sub: p.Disambig}
	from, to := p.Y0, p.Y1
	if p.Born != nil {
		from = *p.Born
		s.Exact = p.Died != nil
		if p.Died != nil {
			to = *p.Died
		} else {
			to = from + 70
		}
	}
	if from != 0 || to != 0 || p.Born != nil {
		if !s.Exact && to-from < 70 {
			m := (from + to) / 2
			from, to = m-35, m+35
		}
		s.HasYears, s.From, s.To = true, from, to
	}
	return s
}

func (a *Atlas) placeSide(p *Place) CmpSide {
	return CmpSide{ID: p.ID, Name: p.Name, Kind: "place", Verses: p.Verses, Books: a.bookCounts(p.Vs), From: p.Y0, To: p.Y1, HasYears: p.Verses > 0,
		Lat: p.Lat, Lng: p.Lng, Sub: strings.TrimSpace(p.Kind + " " + p.Sub)}
}

func intersect(x, y []string) []string {
	m := map[string]bool{}
	for _, v := range y {
		m[v] = true
	}
	var out []string
	for _, v := range x {
		if m[v] {
			out = append(out, v)
		}
	}
	return out
}

func (a *Atlas) topByCount(ids []string, limit int) []string {
	count := map[string]int{}
	for _, id := range ids {
		count[id]++
	}
	var u []string
	for k := range count {
		u = append(u, k)
	}
	sort.Slice(u, func(i, j int) bool {
		if count[u[i]] != count[u[j]] {
			return count[u[i]] > count[u[j]]
		}
		return u[i] < u[j]
	})
	if len(u) > limit {
		u = u[:limit]
	}
	return u
}

func fmtSpan(s CmpSide) string {
	if !s.HasYears {
		return "an unknown time"
	}
	return fmt.Sprintf("%s – %s", FormatYear(s.From), FormatYear(s.To))
}

func (a *Atlas) Compare(kind, idA, idB string) *CompareResult {
	r := &CompareResult{Kind: kind}
	for _, b := range a.Books {
		r.BookNames = append(r.BookNames, b.Name)
	}
	switch kind {
	case "person":
		pa, pb := a.PersonByID[idA], a.PersonByID[idB]
		if pa == nil || pb == nil {
			return nil
		}
		r.A, r.B = a.personSide(pa), a.personSide(pb)
		shared := intersect(pa.Vs, pb.Vs)
		r.SharedCount = len(shared)
		r.SharedVerses, _ = flat(a.groupVerses(shared, 0))
		if len(r.SharedVerses) > 12 {
			r.SharedVerses = r.SharedVerses[:12]
		}
		// events both take part in
		inB := map[string]bool{}
		for _, e := range a.personEvents[idB] {
			inB[e.ID] = true
		}
		var both []*Event
		for _, e := range a.personEvents[idA] {
			if inB[e.ID] {
				both = append(both, e)
			}
		}
		r.SharedEvents = eventRefs(both)
		// places named alongside both of them
		placesOf := func(p *Person) []string {
			var ids []string
			for _, vid := range p.Vs {
				if v := a.VerseByID[vid]; v != nil && len(v.Places) <= 4 {
					ids = append(ids, v.Places...)
				}
			}
			return ids
		}
		pla, plb := placesOf(pa), placesOf(pb)
		sharedPl := intersect(pla, plb)
		r.SharedPlaces = a.placeRefs(a.topByCount(sharedPl, 8))
		r.Relation = a.FamilyPath(idA, idB)
		// distance between where each is mostly named
		ca, cb := a.centroid(a.topByCount(pla, 6)), a.centroid(a.topByCount(plb, 6))
		if ca != nil && cb != nil {
			r.DistanceKm = int(haversine(ca[0], ca[1], cb[0], cb[1]))
		}
		if len(shared) > 0 {
			if v := a.VerseByID[shared[0]]; v != nil {
				best := v
				for _, id := range shared {
					if o := a.VerseByID[id]; o != nil && len(o.Text) < len(best.Text) && len(o.Text) > 40 {
						best = o
					}
				}
				r.QuoteVerse, r.QuoteText, r.QuoteLabel = best.ID, best.Text, a.verseLabel(best)
			}
		}
		r.Story = a.personStory(r, pa, pb)
	case "place":
		pa, pb := a.PlaceByID[idA], a.PlaceByID[idB]
		if pa == nil || pb == nil {
			return nil
		}
		r.A, r.B = a.placeSide(pa), a.placeSide(pb)
		shared := intersect(pa.Vs, pb.Vs)
		r.SharedCount = len(shared)
		r.SharedVerses, _ = flat(a.groupVerses(shared, 0))
		if len(r.SharedVerses) > 12 {
			r.SharedVerses = r.SharedVerses[:12]
		}
		inB := map[string]bool{}
		for _, e := range a.placeEvents[idB] {
			inB[e.ID] = true
		}
		var both []*Event
		for _, e := range a.placeEvents[idA] {
			if inB[e.ID] {
				both = append(both, e)
			}
		}
		r.SharedEvents = eventRefs(both)
		peopleOf := func(p *Place) []string {
			var ids []string
			for _, vid := range p.Vs {
				if v := a.VerseByID[vid]; v != nil && len(v.People) <= 8 {
					ids = append(ids, v.People...)
				}
			}
			return ids
		}
		sp := intersect(peopleOf(pa), peopleOf(pb))
		r.SharedPeople = a.personRefs(a.topByCount(sp, 12))
		if pa.HasGeo && pb.HasGeo {
			r.DistanceKm = int(haversine(pa.Lat, pa.Lng, pb.Lat, pb.Lng))
		}
		r.Story = a.placeStory(r, pa, pb)
	default:
		return nil
	}
	if r.Story == nil {
		r.Story = []string{}
	}
	return r
}

func flat(groups []BookVerses, n int) ([]VerseRef, int) {
	var out []VerseRef
	for _, g := range groups {
		out = append(out, g.Verses...)
	}
	return out, n
}

func (a *Atlas) centroid(placeIDs []string) *[2]float64 {
	var la, ln float64
	n := 0
	for _, id := range placeIDs {
		if p := a.PlaceByID[id]; p != nil && p.HasGeo {
			la += p.Lat
			ln += p.Lng
			n++
		}
	}
	if n == 0 {
		return nil
	}
	return &[2]float64{la / float64(n), ln / float64(n)}
}

func overlapSentence(an, bn string, as, bs CmpSide) string {
	if !as.HasYears || !bs.HasYears {
		return ""
	}
	lo, hi := max(as.From, bs.From), min(as.To, bs.To)
	if hi > lo {
		return fmt.Sprintf("%s (%s) and %s (%s) were alive at the same time for about %d years.", an, fmtSpan(as), bn, fmtSpan(bs), hi-lo)
	}
	gap := lo - hi
	first, second, f, s := an, bn, as, bs
	if bs.To < as.From {
		first, second, f, s = bn, an, bs, as
	}
	return fmt.Sprintf("%s (%s) lived about %d years before %s (%s).", first, fmtSpan(f), gap, second, fmtSpan(s))
}

func (a *Atlas) personStory(r *CompareResult, pa, pb *Person) []string {
	var out []string
	if s := overlapSentence(pa.Name, pb.Name, r.A, r.B); s != "" {
		out = append(out, s)
	}
	if r.Relation != nil {
		out = append(out, r.Relation.Phrase+fmt.Sprintf(" The connection runs through %d steps of family.", len(r.Relation.Steps)-1))
	} else {
		out = append(out, fmt.Sprintf("The text records no family connection between %s and %s.", pa.Name, pb.Name))
	}
	if r.SharedCount > 0 {
		out = append(out, fmt.Sprintf("They are named together in %s.", plural2(r.SharedCount, "verse", "verses")))
	}
	if len(r.SharedEvents) > 0 {
		var t []string
		for i, e := range r.SharedEvents {
			if i >= 3 {
				break
			}
			t = append(t, e.Title)
		}
		out = append(out, "They take part in the same events: "+strings.Join(t, "; ")+".")
	}
	if len(r.SharedPlaces) > 0 {
		var t []string
		for i, p := range r.SharedPlaces {
			if i >= 4 {
				break
			}
			t = append(t, p.Name)
		}
		out = append(out, "Their stories meet at "+joinAnd(t)+".")
	}
	if r.DistanceKm > 0 {
		out = append(out, fmt.Sprintf("The places they are most associated with lie about %d km apart.", r.DistanceKm))
	}
	return out
}

func (a *Atlas) placeStory(r *CompareResult, pa, pb *Place) []string {
	var out []string
	if r.DistanceKm > 0 {
		out = append(out, fmt.Sprintf("%s and %s are about %d km apart as the crow flies; on foot at roughly 30 km a day that is %d days' journey.", pa.Name, pb.Name, r.DistanceKm, max(1, r.DistanceKm*13/10/30)))
	}
	if s := overlapSentence(pa.Name, pb.Name, r.A, r.B); s != "" {
		out = append(out, strings.Replace(s, "were alive at the same time", "are named in the text over the same centuries", 1))
	}
	if r.SharedCount > 0 {
		out = append(out, fmt.Sprintf("They are named together in %s.", plural2(r.SharedCount, "verse", "verses")))
	}
	if len(r.SharedPeople) > 0 {
		var t []string
		for i, p := range r.SharedPeople {
			if i >= 5 {
				break
			}
			t = append(t, p.Name)
		}
		out = append(out, "People who appear at both: "+joinAnd(t)+".")
	}
	if len(r.SharedEvents) > 0 {
		out = append(out, fmt.Sprintf("Both are settings for %s.", plural2(len(r.SharedEvents), "shared event", "shared events")))
	}
	return out
}

func plural2(n int, one, many string) string {
	if n == 1 {
		return fmt.Sprintf("%d %s", n, one)
	}
	return fmt.Sprintf("%d %s", n, many)
}

func joinAnd(s []string) string {
	switch len(s) {
	case 0:
		return ""
	case 1:
		return s[0]
	}
	return strings.Join(s[:len(s)-1], ", ") + " and " + s[len(s)-1]
}
