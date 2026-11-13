// Package atlas loads the cleaned bible metadata into compact in-memory
// indexes and serves it as JSON for the interactive front-end (time map,
// journeys, family tree, cosmos, canon explorer, search and peek cards).
// It deliberately does not depend on Mongo.
package atlas

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type Place struct {
	ID      string   `json:"id"`
	Name    string   `json:"n"`
	Orig    string   `json:"-"` // the source name, when Name is a translation
	Lat     float64  `json:"lat"`
	Lng     float64  `json:"lng"`
	Kind    string   `json:"k"`
	Sub     string   `json:"s,omitempty"`
	Y0      int      `json:"y0"`
	Y1      int      `json:"y1"`
	Verses  int      `json:"vc"`
	HasGeo  bool     `json:"-"`
	Comment string   `json:"-"`
	Dict    string   `json:"-"`
	DictMT  bool     `json:"-"` // Dict is a machine translation
	Vs      []string `json:"-"`
	Born    []string `json:"-"`
	Died    []string `json:"-"`
}

type Person struct {
	ID       string   `json:"id"`
	Name     string   `json:"n"`
	Orig     string   `json:"-"` // the source name, when Name is a translation
	Gender   string   `json:"g,omitempty"`
	Y0       int      `json:"y0"`
	Y1       int      `json:"y1"`
	Born     *int     `json:"b,omitempty"`
	Died     *int     `json:"d,omitempty"`
	Father   string   `json:"f,omitempty"`
	Mother   string   `json:"m,omitempty"`
	Partners []string `json:"p,omitempty"`
	Children []string `json:"c,omitempty"`
	Group    int      `json:"gr"`
	Disambig string   `json:"dis,omitempty"`
	Verses   int      `json:"vc"`
	Dict     string   `json:"-"`
	DictMT   bool     `json:"-"` // Dict is a machine translation
	AlsoCall string   `json:"-"`
	Vs       []string `json:"-"`
	Events   []string `json:"-"`
	BirthPl  string   `json:"-"`
	DeathPl  string   `json:"-"`
}

type Event struct {
	ID      string   `json:"id"`
	Title   string   `json:"t"`
	Orig    string   `json:"-"`              // the source title, when Title is a translation
	Life    bool     `json:"life,omitempty"` // a birth, lifetime or death entry (not a happening in itself)
	Year    float64  `json:"y"`
	Days    int      `json:"d"`
	Places  []string `json:"pl,omitempty"`
	People  []string `json:"pe,omitempty"`
	Pred    []string `json:"pr,omitempty"`
	PartOf  []string `json:"po,omitempty"`
	Ref     string   `json:"ref,omitempty"`
	VerseID string   `json:"vid,omitempty"`
	Text    string   `json:"tx,omitempty"`
	Date    string   `json:"dt"`
	Sort    float64  `json:"-"`
	Notes   string   `json:"-"`
	Vs      []string `json:"-"`
}

type Verse struct {
	ID      string
	Book    string
	Chapter string
	Osis    string
	VID     string
	Num     int
	Text    string
	Rich    string
	People  []string
	Places  []string
	Events  []string
	Year    int
}

type Book struct {
	ID        string        `json:"id"`
	Name      string        `json:"n"`
	Short     string        `json:"s"`
	Order     int           `json:"o"`
	Div       string        `json:"div"`
	Testament string        `json:"t"`
	Osis      string        `json:"osis"`
	Chapters  []ChapterStat `json:"ch"`
	Verses    int           `json:"vc"`
}

type ChapterStat struct {
	ID     string `json:"id"`
	Num    int    `json:"n"`
	Verses int    `json:"v"`
	People int    `json:"pe"`
	Places int    `json:"pl"`
	Events int    `json:"ev"`
}

type Group struct {
	ID      string   `json:"id"`
	Name    string   `json:"n"`
	Members []string `json:"-"`
}

// lazy holds what is worked out on first use. It is a pointer so that a per-language view
// of the atlas (see View) starts with empty caches of its own.
type lazy struct {
	migOnce    sync.Once
	mig        *Migration
	storyMu    sync.Mutex
	stories    map[string]*Story
	cosmosOnce sync.Once
	cosmos     *Cosmos
	progMu     sync.Mutex
	prog       map[string]*ProgressIndex
	bookMu     sync.Mutex
	bookCache  map[string][]*Book
	themeOnce  sync.Once
	themeIdx   map[string][]int
	themeInfos []ThemeInfo
	famOnce    sync.Once
	fam        map[string][]famEdge
}

type Atlas struct {
	*lazy

	Places      []*Place
	People      []*Person
	Events      []*Event
	Books       []*Book
	Groups      []*Group
	Verses      []*Verse
	PlaceByID   map[string]*Place
	PersonByID  map[string]*Person
	EventByID   map[string]*Event
	BookByID    map[string]*Book
	VerseByID   map[string]*Verse
	VerseByOsis map[string]*Verse
	verseByVID  map[string]*Verse
	ChapterNum  map[string]int
	GroupIdx    map[string]int
	Journeys    []*Journey

	lowerText  []string
	bibles     map[string]map[string]string
	lowerBible map[string][]string
	nameIdx    []searchEntry
	names      map[string]*nameTable        // spellings of people and places per language (data/names)
	titles     map[string]map[string]string // event titles per language (data/titles)
	easton     map[string]map[string]string // translated dictionary texts per language (data/easton)
	edited     map[string]map[string]bool   // dictionary texts that were translated by hand (no machine-translation note)
	artwork    map[string]map[string]string // engraving titles per language, by cleaned English title (data/artwork)
	lang       string                       // the language of a per-language view (empty on the source atlas)
	views      *viewSet
	base       *Atlas // set on a per-language view

	tr            Translator
	Enrich        *Enrichment
	artFixes      []artFix
	bookLangs     []string // languages with translated book names, for typing a reference in any of them
	artByEvent    map[string][]Art
	artByPerson   map[string][]Art
	chapterVerses map[string][]*Verse
	chapterBook   map[string]*Book
	personEvents  map[string][]*Event
	placeEvents   map[string][]*Event
	eventSucc     map[string][]*Event
}

var yearRe = regexp.MustCompile(`^(-?\d+)(?:-(\d+)-(\d+))?`)

func atoiSafe(s string) int {
	n, _ := strconv.Atoi(strings.TrimLeft(s, "0"))
	return n
}

func parseYear(s string) (float64, string) {
	m := yearRe.FindStringSubmatch(s)
	if m == nil {
		return 0, ""
	}
	neg := strings.HasPrefix(m[1], "-")
	digits := strings.TrimPrefix(m[1], "-")
	y := float64(atoiSafe(digits))
	if neg {
		y = -y
	}
	frac := 0.0
	if m[2] != "" {
		mo := float64(atoiSafe(m[2]))
		d := float64(atoiSafe(m[3]))
		frac = ((mo-1)*30.4 + d) / 365.0
		if neg {
			// moving later in a BC year goes towards the next (less negative) number
			y += frac
		} else {
			y += frac
		}
	}
	return y, FormatYear(int(y))
}

func FormatYear(y int) string {
	if y < 0 {
		return fmt.Sprintf("%d BC", -y)
	}
	if y == 0 {
		return "1 BC"
	}
	return fmt.Sprintf("%d AD", y)
}

func parseDays(d string) int {
	re := regexp.MustCompile(`(\d+(?:\.\d+)?)([YMWD])`)
	total := 0.0
	for _, m := range re.FindAllStringSubmatch(d, -1) {
		n, _ := strconv.ParseFloat(m[1], 64)
		switch m[2] {
		case "Y":
			total += n * 365
		case "M":
			total += n * 30
		case "W":
			total += n * 7
		case "D":
			total += n
		}
	}
	return int(total)
}

func readJSON[T any](dir, name string) ([]T, error) {
	b, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		return nil, err
	}
	var out []T
	return out, json.Unmarshal(b, &out)
}

type rawPlace struct {
	ID          string   `json:"id"`
	Display     string   `json:"displayTitle"`
	Lat         string   `json:"latitude"`
	Lng         string   `json:"longitude"`
	FeatureType string   `json:"featureType"`
	SubType     string   `json:"featureSubType"`
	Verses      []string `json:"verses"`
	Comment     string   `json:"comment"`
	Dict        string   `json:"dictText"`
	Born        []string `json:"peopleBorn"`
	Died        []string `json:"peopleDied"`
}

type rawPerson struct {
	ID        string   `json:"id"`
	Display   string   `json:"displayTitle"`
	Gender    string   `json:"gender"`
	MinYear   int      `json:"minYear"`
	MaxYear   int      `json:"maxYear"`
	BirthYear []string `json:"birthYear"`
	DeathYear []string `json:"deathYear"`
	Father    string   `json:"father"`
	Mother    string   `json:"mother"`
	Partners  []string `json:"partners"`
	Children  []string `json:"children"`
	MemberOf  []string `json:"memberOf"`
	Verses    []string `json:"verses"`
	Dict      string   `json:"dictText"`
	AlsoCall  string   `json:"alsoCalled"`
	Events    []string `json:"events"`
	Disambig  string   `json:"disambiguation"`
	BirthPl   []string `json:"birthPlace"`
	DeathPl   []string `json:"deathPlace"`
}

type rawEvent struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	StartDate string   `json:"startDate"`
	Duration  string   `json:"duration"`
	Locations []string `json:"locations"`
	PlacesFV  []string `json:"placesFromVerses"`
	Particip  []string `json:"participants"`
	Pred      []string `json:"predecessor"`
	PartOf    []string `json:"partOf"`
	Verses    []string `json:"verses"`
	Notes     string   `json:"notes"`
	SortKey   float64  `json:"sortKey"`
	VerseSort string   `json:"verseSort"`
}

type rawVerse struct {
	ID      string   `json:"id"`
	Book    string   `json:"book"`
	Chapter string   `json:"chapter"`
	Osis    string   `json:"osisRef"`
	VID     string   `json:"verseId"`
	Num     int      `json:"verseNum"`
	Text    string   `json:"verseText"`
	MdText  string   `json:"mdText"`
	People  []string `json:"people"`
	Places  []string `json:"places"`
	Events  []string `json:"eventsDescribed"`
	Year    int      `json:"yearNum"`
}

type rawBook struct {
	ID        string   `json:"id"`
	Name      string   `json:"bookName"`
	Short     string   `json:"shortName"`
	Order     int      `json:"bookOrder"`
	Div       string   `json:"bookDiv"`
	Testament string   `json:"testament"`
	Osis      string   `json:"osisName"`
	Chapters  []string `json:"chapters"`
	Verses    int      `json:"verseCount"`
}

type rawChapter struct {
	ID  string `json:"id"`
	Num int    `json:"chapterNum"`
}

type rawGroup struct {
	ID      string   `json:"id"`
	Name    string   `json:"groupName"`
	Members []string `json:"members"`
}

type rawPeriod struct {
	ID      string `json:"id"`
	IsoYear int    `json:"isoYear"`
}

// Load reads the json files in dir (data/json) and the curated journeys file.
func Load(dir string, journeysFile string) (*Atlas, error) {
	a := &Atlas{
		lazy:      &lazy{},
		views:     &viewSet{m: map[string]*Atlas{}},
		PlaceByID: map[string]*Place{}, PersonByID: map[string]*Person{},
		EventByID: map[string]*Event{}, BookByID: map[string]*Book{},
		VerseByID: map[string]*Verse{}, VerseByOsis: map[string]*Verse{}, verseByVID: map[string]*Verse{},
		ChapterNum: map[string]int{}, GroupIdx: map[string]int{},
	}

	rbooks, err := readJSON[rawBook](dir, "books.json")
	if err != nil {
		return nil, err
	}
	rchap, err := readJSON[rawChapter](dir, "chapters.json")
	if err != nil {
		return nil, err
	}
	rverses, err := readJSON[rawVerse](dir, "verses.json")
	if err != nil {
		return nil, err
	}
	rplaces, err := readJSON[rawPlace](dir, "places.json")
	if err != nil {
		return nil, err
	}
	rpeople, err := readJSON[rawPerson](dir, "people.json")
	if err != nil {
		return nil, err
	}
	revents, err := readJSON[rawEvent](dir, "events.json")
	if err != nil {
		return nil, err
	}
	rgroups, err := readJSON[rawGroup](dir, "peopleGroups.json")
	if err != nil {
		return nil, err
	}
	rperiods, err := readJSON[rawPeriod](dir, "periods.json")
	if err != nil {
		return nil, err
	}

	periodYear := map[string]int{}
	for _, p := range rperiods {
		periodYear[p.ID] = p.IsoYear
	}
	for _, c := range rchap {
		a.ChapterNum[c.ID] = c.Num
	}

	for _, b := range rbooks {
		bk := &Book{ID: b.ID, Name: b.Name, Short: b.Short, Order: b.Order, Div: b.Div, Testament: b.Testament, Osis: b.Osis, Verses: b.Verses}
		for _, cid := range b.Chapters {
			bk.Chapters = append(bk.Chapters, ChapterStat{ID: cid, Num: a.ChapterNum[cid]})
		}
		sort.Slice(bk.Chapters, func(i, j int) bool { return bk.Chapters[i].Num < bk.Chapters[j].Num })
		a.Books = append(a.Books, bk)
		a.BookByID[bk.ID] = bk
	}
	sort.Slice(a.Books, func(i, j int) bool { return a.Books[i].Order < a.Books[j].Order })

	for i, g := range rgroups {
		a.Groups = append(a.Groups, &Group{ID: g.ID, Name: g.Name, Members: g.Members})
		a.GroupIdx[g.ID] = i + 1
	}

	for _, v := range rverses {
		vv := &Verse{ID: v.ID, Book: v.Book, Chapter: v.Chapter, Osis: v.Osis, VID: v.VID, Num: v.Num, Text: v.Text, Rich: v.MdText,
			People: v.People, Places: v.Places, Events: v.Events, Year: v.Year}
		a.Verses = append(a.Verses, vv)
		a.VerseByID[vv.ID] = vv
		a.VerseByOsis[strings.ToLower(vv.Osis)] = vv
		a.verseByVID[vv.VID] = vv
	}
	sort.Slice(a.Verses, func(i, j int) bool { return a.Verses[i].VID < a.Verses[j].VID })

	// per place first/last year mention and per person verse stats
	y0 := map[string]int{}
	y1 := map[string]int{}
	for _, v := range a.Verses {
		for _, pid := range v.Places {
			if cur, ok := y0[pid]; !ok || v.Year < cur {
				y0[pid] = v.Year
			}
			if cur, ok := y1[pid]; !ok || v.Year > cur {
				y1[pid] = v.Year
			}
		}
	}

	for _, p := range rplaces {
		lat, e1 := strconv.ParseFloat(p.Lat, 64)
		lng, e2 := strconv.ParseFloat(p.Lng, 64)
		pl := &Place{ID: p.ID, Name: p.Display, Lat: lat, Lng: lng, Kind: p.FeatureType, Sub: p.SubType,
			HasGeo: e1 == nil && e2 == nil, Verses: len(p.Verses), Y0: y0[p.ID], Y1: y1[p.ID],
			Comment: p.Comment, Dict: p.Dict, Vs: p.Verses, Born: p.Born, Died: p.Died}
		a.Places = append(a.Places, pl)
		a.PlaceByID[pl.ID] = pl
	}

	for _, p := range rpeople {
		pe := &Person{ID: p.ID, Name: p.Display, Gender: p.Gender, Y0: p.MinYear, Y1: p.MaxYear,
			Father: p.Father, Mother: p.Mother, Partners: p.Partners, Children: p.Children,
			Verses: len(p.Verses), Dict: p.Dict, AlsoCall: p.AlsoCall, Vs: p.Verses, Events: p.Events, Disambig: p.Disambig}
		if len(p.BirthPl) > 0 {
			pe.BirthPl = p.BirthPl[0]
		}
		if len(p.DeathPl) > 0 {
			pe.DeathPl = p.DeathPl[0]
		}
		if len(p.BirthYear) > 0 {
			if y, ok := periodYear[p.BirthYear[0]]; ok {
				yy := y
				pe.Born = &yy
			}
		}
		if len(p.DeathYear) > 0 {
			if y, ok := periodYear[p.DeathYear[0]]; ok {
				yy := y
				pe.Died = &yy
			}
		}
		if len(p.MemberOf) > 0 {
			pe.Group = a.GroupIdx[p.MemberOf[0]]
		}
		a.People = append(a.People, pe)
		a.PersonByID[pe.ID] = pe
	}

	a.estimateActiveYears()

	for _, e := range revents {
		y, dt := parseYear(e.StartDate)
		ev := &Event{ID: e.ID, Title: e.Title, Life: isLifespanEvent(e.Title), Year: y, Date: dt, Days: parseDays(e.Duration),
			Pred: e.Pred, PartOf: e.PartOf, People: e.Particip, Sort: e.SortKey, Notes: e.Notes, Vs: e.Verses}
		seen := map[string]bool{}
		for _, id := range append(append([]string{}, e.Locations...), e.PlacesFV...) {
			if pl, ok := a.PlaceByID[id]; ok && pl.HasGeo && !seen[id] {
				seen[id] = true
				ev.Places = append(ev.Places, id)
			}
		}
		if len(e.Locations) > 0 {
			// explicit locations rank first, then placesFromVerses only if none
			ev.Places = nil
			seen = map[string]bool{}
			for _, id := range e.Locations {
				if pl, ok := a.PlaceByID[id]; ok && pl.HasGeo && !seen[id] {
					seen[id] = true
					ev.Places = append(ev.Places, id)
				}
			}
		}
		if len(e.Verses) > 0 {
			// the verse list is unordered: prefer the verse matching verseSort, else the earliest
			var best *Verse
			if v, ok := a.verseByVID[e.VerseSort]; ok {
				best = v
			}
			for _, vid := range e.Verses {
				if v, ok := a.VerseByID[vid]; ok && (best == nil || (e.VerseSort == "" && v.VID < best.VID)) {
					best = v
				}
			}
			if v := best; v != nil {
				ev.Ref = v.Osis
				ev.VerseID = v.ID
				ev.Text = v.Text
			}
		}
		a.Events = append(a.Events, ev)
		a.EventByID[ev.ID] = ev
	}
	sort.Slice(a.Events, func(i, j int) bool { return a.Events[i].Sort < a.Events[j].Sort })

	if err := a.applyCorrections(filepath.Join(filepath.Dir(journeysFile), "corrections.json")); err != nil {
		return nil, err
	}
	if err := a.loadEnrichment(filepath.Join(filepath.Dir(journeysFile), "enrichment.json")); err != nil {
		return nil, fmt.Errorf("enrichment: %w", err)
	}

	// book / chapter statistics
	type cs struct {
		people, places map[string]bool
		events         map[string]bool
		verses         int
	}
	stats := map[string]*cs{}
	for _, v := range a.Verses {
		s := stats[v.Chapter]
		if s == nil {
			s = &cs{map[string]bool{}, map[string]bool{}, map[string]bool{}, 0}
			stats[v.Chapter] = s
		}
		s.verses++
		for _, p := range v.People {
			s.people[p] = true
		}
		for _, p := range v.Places {
			s.places[p] = true
		}
		for _, e := range v.Events {
			s.events[e] = true
		}
	}
	for _, b := range a.Books {
		for i := range b.Chapters {
			if s := stats[b.Chapters[i].ID]; s != nil {
				b.Chapters[i].Verses = s.verses
				b.Chapters[i].People = len(s.people)
				b.Chapters[i].Places = len(s.places)
				b.Chapters[i].Events = len(s.events)
			}
		}
	}

	a.chapterVerses = map[string][]*Verse{}
	a.chapterBook = map[string]*Book{}
	for _, v := range a.Verses {
		a.chapterVerses[v.Chapter] = append(a.chapterVerses[v.Chapter], v)
	}
	for _, b := range a.Books {
		for _, c := range b.Chapters {
			a.chapterBook[c.ID] = b
		}
	}
	a.personEvents = map[string][]*Event{}
	a.placeEvents = map[string][]*Event{}
	a.eventSucc = map[string][]*Event{}
	for _, e := range a.Events {
		for _, id := range e.People {
			a.personEvents[id] = append(a.personEvents[id], e)
		}
		for _, id := range e.Places {
			a.placeEvents[id] = append(a.placeEvents[id], e)
		}
		for _, id := range e.Pred {
			a.eventSucc[id] = append(a.eventSucc[id], e)
		}
	}

	a.lowerText = make([]string, len(a.Verses))
	for i, v := range a.Verses {
		a.lowerText[i] = strings.ToLower(v.Text)
	}
	a.loadBibles(dir)
	a.loadNames(dir)
	a.buildSearch()

	if journeysFile != "" {
		js, err := loadJourneys(a, journeysFile)
		if err != nil {
			return nil, err
		}
		a.Journeys = js
	}
	return a, nil
}

// estimateActiveYears replaces the raw first/last-mention years (which stretch
// from a person's own time to every later book that mentions them) with the
// densest window of years holding most of their verses, or their known
// birth/death when available.
func (a *Atlas) estimateActiveYears() {
	years := map[string][]int{}
	for _, v := range a.Verses {
		for _, id := range v.People {
			years[id] = append(years[id], v.Year)
		}
	}
	for _, p := range a.People {
		ys := years[p.ID]
		if len(ys) == 0 {
			continue
		}
		sort.Ints(ys)
		k := (len(ys)*5 + 9) / 10
		if k < 1 {
			k = 1
		}
		best, bi := 1<<30, 0
		for i := 0; i+k-1 < len(ys); i++ {
			if d := ys[i+k-1] - ys[i]; d < best {
				best, bi = d, i
			}
		}
		p.Y0, p.Y1 = ys[bi], ys[bi+k-1]
	}
}
