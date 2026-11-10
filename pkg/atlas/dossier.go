package atlas

import (
	"encoding/json"
	"fmt"
	"html"
	"regexp"
	"sort"
	"strings"
)

// ---------------------------------------------------------------- rich text

var richLink = regexp.MustCompile(`\[([^\]]*)\]\(([^)]*)\)`)
var richHref = regexp.MustCompile(`^/(person|place|verse|event|chapter|book)/([A-Za-z0-9]+)$`)

var emphRe = regexp.MustCompile(`(^|[^A-Za-z0-9])_([^_]+)_([^A-Za-z0-9]|$)`)
var boldRe = regexp.MustCompile(`\*\*([^*]+)\*\*`)

// plainText escapes a text run and renders the light markup used in the data (_italic_, **bold**).
func plainText(s string) string {
	e := html.EscapeString(s)
	e = emphRe.ReplaceAllString(e, "$1<em>$2</em>$3")
	return boldRe.ReplaceAllString(e, "<strong>$1</strong>")
}

// RichInline turns the data's markdown links into safe anchors; anything that is
// not a clean internal link is reduced to its label.
func RichInline(s string) string {
	var b strings.Builder
	last := 0
	for _, m := range richLink.FindAllStringSubmatchIndex(s, -1) {
		b.WriteString(plainText(s[last:m[0]]))
		label, href := s[m[2]:m[3]], s[m[4]:m[5]]
		if h := richHref.FindStringSubmatch(href); h != nil {
			fmt.Fprintf(&b, `<a class="k-%s" href="%s">%s</a>`, h[1], h[0], html.EscapeString(label))
		} else {
			b.WriteString(html.EscapeString(label))
		}
		last = m[1]
	}
	b.WriteString(plainText(s[last:]))
	return b.String()
}

// Rich renders multi-paragraph text.
func Rich(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	var b strings.Builder
	for _, p := range strings.Split(s, "\n\n") {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		b.WriteString("<p>" + strings.ReplaceAll(RichInline(p), "\n", " ") + "</p>")
	}
	return b.String()
}

// ---------------------------------------------------------------- shared bits

type Ref struct {
	ID   string
	Name string
	Sub  string
}

type VerseRef struct {
	ID      string
	Label   string
	Excerpt string
}

type BookVerses struct {
	Book   string
	BookID string
	Verses []VerseRef
}

type EventRef struct {
	ID    string
	Title string
	Date  string
	Year  float64
}

func (a *Atlas) verseRef(v *Verse) VerseRef {
	return VerseRef{ID: v.ID, Label: a.verseLabel(v), Excerpt: trunc(v.Text, 110)}
}

func (a *Atlas) groupVerses(ids []string, limit int) ([]BookVerses, int) {
	var vs []*Verse
	for _, id := range ids {
		if v := a.VerseByID[id]; v != nil {
			vs = append(vs, v)
		}
	}
	sort.Slice(vs, func(i, j int) bool { return vs[i].VID < vs[j].VID })
	var out []BookVerses
	for _, v := range vs {
		b := a.BookByID[v.Book]
		if b == nil {
			continue
		}
		if len(out) == 0 || out[len(out)-1].BookID != b.ID {
			out = append(out, BookVerses{Book: b.Name, BookID: b.ID})
		}
		g := &out[len(out)-1]
		if limit <= 0 || len(g.Verses) < limit {
			g.Verses = append(g.Verses, a.verseRef(v))
		}
	}
	return out, len(vs)
}

func (a *Atlas) personRef(id string) *Ref {
	if p := a.PersonByID[id]; p != nil {
		return &Ref{ID: p.ID, Name: p.Name}
	}
	return nil
}

func (a *Atlas) personRefs(ids []string) []Ref {
	var out []Ref
	seen := map[string]bool{}
	for _, id := range ids {
		if p := a.PersonByID[id]; p != nil && !seen[id] {
			seen[id] = true
			out = append(out, Ref{ID: p.ID, Name: p.Name})
		}
	}
	return out
}

func (a *Atlas) placeRefs(ids []string) []Ref {
	var out []Ref
	seen := map[string]bool{}
	for _, id := range ids {
		if p := a.PlaceByID[id]; p != nil && !seen[id] {
			seen[id] = true
			out = append(out, Ref{ID: p.ID, Name: p.Name, Sub: strings.ToLower(p.Kind)})
		}
	}
	return out
}

func eventRefs(es []*Event) []EventRef {
	out := make([]EventRef, 0, len(es))
	seen := map[string]bool{}
	for _, e := range es {
		if !seen[e.ID] {
			seen[e.ID] = true
			out = append(out, EventRef{e.ID, e.Title, e.Date, e.Year})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Year < out[j].Year })
	return out
}

func eraAt(y int) Era {
	for _, e := range Eras {
		if y >= e.From && y < e.To {
			return e
		}
	}
	return Eras[len(Eras)-1]
}

// Pct places a year on the 4100 BC – AD 100 scale (0–100).
func Pct(y int) float64 {
	v := float64(y+4100) / 4200 * 100
	if v < 0 {
		return 0
	}
	if v > 100 {
		return 100
	}
	return v
}

type MapPoint struct {
	ID   string  `json:"id,omitempty"`
	Name string  `json:"n"`
	Lat  float64 `json:"lat"`
	Lng  float64 `json:"lng"`
	Kind string  `json:"k,omitempty"`
	Year int     `json:"y,omitempty"`
}

type MapConfig struct {
	Points []MapPoint `json:"points"`
	Route  bool       `json:"route,omitempty"`
	Focus  *MapPoint  `json:"focus,omitempty"`
	Zoom   float64    `json:"zoom,omitempty"`
	Pitch  float64    `json:"pitch,omitempty"`
	Color  string     `json:"color,omitempty"`
	Link   string     `json:"link,omitempty"`
}

func (m MapConfig) JSON() string {
	b, _ := json.Marshal(m)
	return string(b)
}

func (a *Atlas) points(ids []string) []MapPoint {
	var out []MapPoint
	seen := map[string]bool{}
	for _, id := range ids {
		if p := a.PlaceByID[id]; p != nil && p.HasGeo && !seen[id] {
			seen[id] = true
			out = append(out, MapPoint{ID: p.ID, Name: p.Name, Lat: p.Lat, Lng: p.Lng, Kind: p.Kind})
		}
	}
	return out
}

// ---------------------------------------------------------------- person

type PersonDossier struct {
	P          *Person
	Disambig   string
	GroupName  string
	Era        Era
	HasLife    bool
	LifeFrom   int
	LifeTo     int
	LifeExact  bool
	Pct0, Pct1 float64
	BirthPlace *Ref
	DeathPlace *Ref
	Father     *Ref
	Mother     *Ref
	Partners   []Ref
	Siblings   []Ref
	Children   []Ref
	Dict       string
	Events     []EventRef
	Verses     []BookVerses
	VerseCount int
	Map        MapConfig
	HasMap     bool
	TopPlaces  []Ref
	Wiki       *WikiEntry
	Art        []Art
	Gallery    []GalleryItem
	Aliases    []string
	WikiPhoto  *GalleryItem
}

func (a *Atlas) PersonDossier(id string) *PersonDossier {
	p := a.PersonByID[id]
	if p == nil {
		return nil
	}
	d := &PersonDossier{P: p, Disambig: p.Disambig, Dict: Rich(p.Dict), Wiki: a.wikiPerson(id), Aliases: splitAliases(p.AlsoCall, p.Name)}
	if arts := a.artByPerson[id]; len(arts) > 0 {
		d.Art = arts[:min(len(arts), 24)]
		d.Gallery = a.ArtGallery(d.Art)
	}
	d.WikiPhoto = a.WikiPhotoItem("person", id, p.Name, d.Wiki)
	if p.Group > 0 && p.Group <= len(a.Groups) {
		d.GroupName = a.Groups[p.Group-1].Name
	}
	from, to := p.Y0, p.Y1
	if p.Born != nil {
		from = *p.Born
		d.LifeExact = p.Died != nil
		if p.Died != nil {
			to = *p.Died
		} else {
			to = from + 70
		}
	}
	if from != 0 || to != 0 || p.Born != nil {
		if !d.LifeExact && to-from < 70 {
			m := (from + to) / 2
			from, to = m-35, m+35
		}
		d.HasLife, d.LifeFrom, d.LifeTo = true, from, to
		d.Pct0, d.Pct1 = Pct(from), Pct(to)
		d.Era = eraAt((from + to) / 2)
	}
	d.BirthPlace = refOrNil(a.placeRefs([]string{p.BirthPl}))
	d.DeathPlace = refOrNil(a.placeRefs([]string{p.DeathPl}))
	d.Father, d.Mother = a.personRef(p.Father), a.personRef(p.Mother)
	d.Partners = a.personRefs(p.Partners)
	d.Children = a.personRefs(p.Children)
	sib := map[string]bool{}
	var sibIDs []string
	for _, par := range []string{p.Father, p.Mother} {
		if pp := a.PersonByID[par]; pp != nil {
			for _, c := range pp.Children {
				if c != id && !sib[c] {
					sib[c] = true
					sibIDs = append(sibIDs, c)
				}
			}
		}
	}
	d.Siblings = a.personRefs(sibIDs)
	d.Verses, d.VerseCount = a.groupVerses(p.Vs, 0)

	evs := append([]*Event{}, a.personEvents[id]...)
	for _, eid := range p.Events {
		if e := a.EventByID[eid]; e != nil {
			evs = append(evs, e)
		}
	}
	d.Events = eventRefs(evs)

	// places mentioned alongside, most frequent first, for the mini-map
	count := map[string]int{}
	for _, vid := range p.Vs {
		if v := a.VerseByID[vid]; v != nil && len(v.Places) <= 4 {
			for _, pl := range v.Places {
				count[pl]++
			}
		}
	}
	var ids []string
	for k := range count {
		if pl := a.PlaceByID[k]; pl != nil && pl.HasGeo {
			ids = append(ids, k)
		}
	}
	sort.Slice(ids, func(i, j int) bool {
		if count[ids[i]] != count[ids[j]] {
			return count[ids[i]] > count[ids[j]]
		}
		return ids[i] < ids[j]
	})
	if len(ids) > 25 {
		ids = ids[:25]
	}
	d.TopPlaces = a.placeRefs(ids[:min(len(ids), 10)])
	if len(ids) > 0 {
		d.HasMap = true
		d.Map = MapConfig{Points: a.points(ids), Zoom: 5, Color: "#f5c451", Link: "/journeys?person=" + id}
	}
	return d
}

func refOrNil(r []Ref) *Ref {
	if len(r) == 0 {
		return nil
	}
	return &r[0]
}

// ---------------------------------------------------------------- place

type PlaceDossier struct {
	P          *Place
	Dict       string
	Comment    string
	Kind       string
	Events     []EventRef
	Born       []Ref
	Died       []Ref
	Nearby     []Ref
	People     []Ref
	Verses     []BookVerses
	VerseCount int
	FirstYear  int
	LastYear   int
	Map        MapConfig
	Wiki       *WikiEntry
	WikiPhoto  *GalleryItem
}

func (a *Atlas) PlaceDossier(id string) *PlaceDossier {
	p := a.PlaceByID[id]
	if p == nil {
		return nil
	}
	d := &PlaceDossier{P: p, Dict: Rich(p.Dict), Comment: p.Comment, Kind: strings.TrimSpace(p.Kind + " " + p.Sub), FirstYear: p.Y0, LastYear: p.Y1, Wiki: a.wikiPlace(id)}
	d.WikiPhoto = a.WikiPhotoItem("place", id, p.Name, d.Wiki)
	d.Events = eventRefs(a.placeEvents[id])
	d.Born, d.Died = a.personRefs(p.Born), a.personRefs(p.Died)
	d.Verses, d.VerseCount = a.groupVerses(p.Vs, 0)

	// people who appear most often in verses naming this place
	count := map[string]int{}
	for _, vid := range p.Vs {
		if v := a.VerseByID[vid]; v != nil && len(v.People) <= 8 {
			for _, pe := range v.People {
				count[pe]++
			}
		}
	}
	var pids []string
	for k := range count {
		pids = append(pids, k)
	}
	sort.Slice(pids, func(i, j int) bool {
		if count[pids[i]] != count[pids[j]] {
			return count[pids[i]] > count[pids[j]]
		}
		return pids[i] < pids[j]
	})
	if len(pids) > 14 {
		pids = pids[:14]
	}
	d.People = a.personRefs(pids)

	if p.HasGeo {
		type nb struct {
			p *Place
			d float64
		}
		var near []nb
		for _, o := range a.Places {
			if o.ID == id || !o.HasGeo || o.Verses == 0 {
				continue
			}
			if km := haversine(p.Lat, p.Lng, o.Lat, o.Lng); km < 80 {
				near = append(near, nb{o, km})
			}
		}
		sort.Slice(near, func(i, j int) bool { return near[i].p.Verses > near[j].p.Verses })
		if len(near) > 16 {
			near = near[:16]
		}
		pts := []MapPoint{{ID: p.ID, Name: p.Name, Lat: p.Lat, Lng: p.Lng, Kind: p.Kind}}
		for _, n := range near {
			d.Nearby = append(d.Nearby, Ref{ID: n.p.ID, Name: n.p.Name, Sub: fmt.Sprintf("%d km", int(n.d))})
			pts = append(pts, MapPoint{ID: n.p.ID, Name: n.p.Name, Lat: n.p.Lat, Lng: n.p.Lng, Kind: n.p.Kind})
		}
		d.Map = MapConfig{Points: pts, Focus: &pts[0], Zoom: 8, Pitch: 50, Link: "/atlas?place=" + id}
	}
	return d
}

// ---------------------------------------------------------------- event

type EventVerse struct {
	ID    string
	Label string
	Text  string
}

type EventDossier struct {
	E            *Event
	Duration     string
	Notes        string
	Before       []EventRef
	After        []EventRef
	PartOf       []EventRef
	Parts        []EventRef
	Participants []Ref
	Places       []Ref
	Verses       []EventVerse
	Prev, Next   string
	Era          Era
	Map          MapConfig
	HasMap       bool
	Art          []Art
	Gallery      []GalleryItem
}

func humanDuration(days int) string {
	switch {
	case days <= 0:
		return ""
	case days == 1:
		return "1 day"
	case days >= 365:
		y := float64(days) / 365
		if y == float64(int(y)) {
			return fmt.Sprintf("%d years", int(y))
		}
		return fmt.Sprintf("%.1f years", y)
	case days >= 30:
		return fmt.Sprintf("%d months", days/30)
	}
	return fmt.Sprintf("%d days", days)
}

func (a *Atlas) EventDossier(id string) *EventDossier {
	e := a.EventByID[id]
	if e == nil {
		return nil
	}
	d := &EventDossier{E: e, Duration: humanDuration(e.Days), Notes: e.Notes, Era: eraAt(int(e.Year)), Art: a.artByEvent[id]}
	d.Gallery = a.ArtGallery(d.Art)
	for _, pid := range e.Pred {
		if o := a.EventByID[pid]; o != nil {
			d.Before = append(d.Before, EventRef{o.ID, o.Title, o.Date, o.Year})
		}
	}
	d.After = eventRefs(a.eventSucc[id])
	for _, pid := range e.PartOf {
		if o := a.EventByID[pid]; o != nil {
			d.PartOf = append(d.PartOf, EventRef{o.ID, o.Title, o.Date, o.Year})
		}
	}
	for _, o := range a.Events {
		for _, pid := range o.PartOf {
			if pid == id {
				d.Parts = append(d.Parts, EventRef{o.ID, o.Title, o.Date, o.Year})
			}
		}
	}
	d.Participants = a.personRefs(e.People)
	d.Places = a.placeRefs(e.Places)
	var vs []*Verse
	for _, vid := range e.Vs {
		if v := a.VerseByID[vid]; v != nil {
			vs = append(vs, v)
		}
	}
	sort.Slice(vs, func(i, j int) bool { return vs[i].VID < vs[j].VID })
	for _, v := range vs {
		d.Verses = append(d.Verses, EventVerse{v.ID, a.verseLabel(v), v.Text})
	}
	for i, o := range a.Events {
		if o.ID == id {
			if i > 0 {
				d.Prev = a.Events[i-1].ID
			}
			if i < len(a.Events)-1 {
				d.Next = a.Events[i+1].ID
			}
		}
	}
	if pts := a.points(e.Places); len(pts) > 0 {
		d.HasMap = true
		d.Map = MapConfig{Points: pts, Zoom: 7, Pitch: 50, Color: d.Era.Color, Link: "/atlas?event=" + id}
	}
	return d
}

// ---------------------------------------------------------------- verse + chapter

type CtxVerse struct {
	ID      string
	Num     int
	HTML    string
	Current bool
}

type VerseDossier struct {
	V          *Verse
	Book       *Book
	ChapterID  string
	ChapterNum int
	Label      string
	HTML       string
	Prev, Next string
	Context    []CtxVerse
	People     []Ref
	Places     []Ref
	Events     []EventRef
	Year       string
	Map        MapConfig
	HasMap     bool
	Art        []Art
	Gallery    []GalleryItem
}

func (a *Atlas) VerseDossier(id string) *VerseDossier {
	v := a.VerseByID[id]
	if v == nil {
		return nil
	}
	b := a.BookByID[v.Book]
	if b == nil {
		return nil
	}
	d := &VerseDossier{V: v, Book: b, ChapterID: v.Chapter, ChapterNum: a.ChapterNum[v.Chapter], Label: a.verseLabel(v), HTML: RichInline(mdOf(v)), Year: FormatYear(v.Year)}
	siblings := a.chapterVerses[v.Chapter]
	for i, o := range siblings {
		if o.ID != id {
			continue
		}
		if i > 0 {
			d.Prev = siblings[i-1].ID
		}
		if i < len(siblings)-1 {
			d.Next = siblings[i+1].ID
		}
		for j := max(0, i-4); j <= min(len(siblings)-1, i+4); j++ {
			d.Context = append(d.Context, CtxVerse{siblings[j].ID, siblings[j].Num, RichInline(mdOf(siblings[j])), j == i})
		}
	}
	d.People, d.Places = a.personRefs(v.People), a.placeRefs(v.Places)
	var evs []*Event
	for _, eid := range v.Events {
		if e := a.EventByID[eid]; e != nil {
			evs = append(evs, e)
		}
	}
	d.Events = eventRefs(evs)
	for _, e := range evs {
		d.Art = append(d.Art, a.artByEvent[e.ID]...)
	}
	d.Gallery = a.ArtGallery(d.Art)
	if pts := a.points(v.Places); len(pts) > 0 {
		d.HasMap = true
		d.Map = MapConfig{Points: pts, Zoom: 6, Pitch: 40, Link: "/atlas?year=" + fmt.Sprint(v.Year)}
	}
	return d
}

// mdOf prefers the rich text (with links); falls back to the plain verse text.
func mdOf(v *Verse) string {
	if v.Rich != "" {
		return v.Rich
	}
	return v.Text
}

type ReadVerse struct {
	ID      string       `json:"id"`
	Num     int          `json:"n"`
	HTML    string       `json:"-"`
	Places  []string     `json:"pl,omitempty"`
	Year    int          `json:"y"`
	Events  []EventRef   `json:"-"`
	Art     *Art         `json:"-"`
	ArtItem *GalleryItem `json:"-"`
}

type ChapterDossier struct {
	Book       *Book
	ID         string
	Num        int
	Verses     []ReadVerse
	Prev, Next string
	PrevLabel  string
	NextLabel  string
	Map        MapConfig
	Years      string
	People     []Ref
}

func (a *Atlas) ChapterDossier(id string) *ChapterDossier {
	b := a.chapterBook[id]
	vs := a.chapterVerses[id]
	if b == nil || len(vs) == 0 {
		return nil
	}
	d := &ChapterDossier{Book: b, ID: id, Num: a.ChapterNum[id]}
	var placeIDs []string
	people := map[string]int{}
	for _, v := range vs {
		rv := ReadVerse{ID: v.ID, Num: v.Num, HTML: RichInline(mdOf(v)), Places: v.Places, Year: v.Year}
		var evs []*Event
		for _, eid := range v.Events {
			if e := a.EventByID[eid]; e != nil {
				evs = append(evs, e)
			}
		}
		rv.Events = eventRefs(evs)
		// an engraving is shown once, under the first verse of its event
		for _, e := range evs {
			if e.VerseID == v.ID {
				if arts := a.artByEvent[e.ID]; len(arts) > 0 {
					art := arts[0]
					rv.Art = &art
					it := a.artItem(art)
					rv.ArtItem = &it
				}
			}
		}
		d.Verses = append(d.Verses, rv)
		placeIDs = append(placeIDs, v.Places...)
		for _, p := range v.People {
			people[p]++
		}
	}
	for i, c := range b.Chapters {
		if c.ID == id {
			if i > 0 {
				d.Prev, d.PrevLabel = b.Chapters[i-1].ID, fmt.Sprintf("%s %d", b.Name, b.Chapters[i-1].Num)
			}
			if i < len(b.Chapters)-1 {
				d.Next, d.NextLabel = b.Chapters[i+1].ID, fmt.Sprintf("%s %d", b.Name, b.Chapters[i+1].Num)
			}
		}
	}
	d.Years = FormatYear(vs[0].Year)
	if last := vs[len(vs)-1].Year; last != vs[0].Year {
		d.Years += " – " + FormatYear(last)
	}
	var pids []string
	for k := range people {
		pids = append(pids, k)
	}
	sort.Slice(pids, func(i, j int) bool {
		if people[pids[i]] != people[pids[j]] {
			return people[pids[i]] > people[pids[j]]
		}
		return pids[i] < pids[j]
	})
	d.People = a.personRefs(pids[:min(len(pids), 14)])
	d.Map = MapConfig{Points: a.points(placeIDs), Zoom: 6, Pitch: 45}
	return d
}

// ReadPayload is the per-verse place index used by the reading-view map.
type ReadPayload struct {
	Verses []ReadVerse `json:"verses"`
	Points []MapPoint  `json:"points"`
}

func (d *ChapterDossier) ReadData() ReadPayload {
	return ReadPayload{d.Verses, d.Map.Points}
}

// ---------------------------------------------------------------- book

type BookChapter struct {
	ID      string
	Num     int
	Excerpt string
	People  int
	Places  int
	Heat    float64
}

type BookDossier struct {
	B          *Book
	Chapters   []BookChapter
	TopPeople  []Ref
	TopPlaces  []Ref
	Prev, Next *Ref
	PeopleN    int
	PlacesN    int
	FirstYear  string
}

func (a *Atlas) BookDossier(id string) *BookDossier {
	b := a.BookByID[id]
	if b == nil {
		return nil
	}
	d := &BookDossier{B: b}
	maxPe := 1
	for _, c := range b.Chapters {
		if c.People > maxPe {
			maxPe = c.People
		}
	}
	pc, plc := map[string]int{}, map[string]int{}
	for _, c := range b.Chapters {
		bc := BookChapter{ID: c.ID, Num: c.Num, People: c.People, Places: c.Places, Heat: float64(c.People) / float64(maxPe)}
		if vs := a.chapterVerses[c.ID]; len(vs) > 0 {
			bc.Excerpt = trunc(vs[0].Text, 120)
			for _, v := range vs {
				for _, p := range v.People {
					pc[p]++
				}
				for _, p := range v.Places {
					plc[p]++
				}
			}
		}
		d.Chapters = append(d.Chapters, bc)
	}
	top := func(m map[string]int, n int) []string {
		var ids []string
		for k := range m {
			ids = append(ids, k)
		}
		sort.Slice(ids, func(i, j int) bool {
			if m[ids[i]] != m[ids[j]] {
				return m[ids[i]] > m[ids[j]]
			}
			return ids[i] < ids[j]
		})
		return ids[:min(len(ids), n)]
	}
	d.PeopleN, d.PlacesN = len(pc), len(plc)
	d.TopPeople = a.personRefs(top(pc, 16))
	d.TopPlaces = a.placeRefs(top(plc, 16))
	for i, o := range a.Books {
		if o.ID == id {
			if i > 0 {
				d.Prev = &Ref{ID: a.Books[i-1].ID, Name: a.Books[i-1].Name}
			}
			if i < len(a.Books)-1 {
				d.Next = &Ref{ID: a.Books[i+1].ID, Name: a.Books[i+1].Name}
			}
		}
	}
	return d
}

// splitAliases turns the data's "Lord,Son,Lamb" into separate names, once each (ignoring case)
// and without repeating the person's own name.
func splitAliases(raw, self string) []string {
	seen := map[string]bool{strings.ToLower(strings.TrimSpace(parenRe.ReplaceAllString(self, ""))): true}
	var out []string
	for _, part := range strings.FieldsFunc(raw, func(r rune) bool { return r == ',' || r == ';' }) {
		part = strings.TrimSpace(part)
		k := strings.ToLower(part)
		if part == "" || seen[k] {
			continue
		}
		seen[k] = true
		out = append(out, part)
	}
	return out
}
