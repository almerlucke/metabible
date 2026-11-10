package atlas

import (
	"strings"
	"sync"
	"testing"
)

var (
	once sync.Once
	atl  *Atlas
	lerr error
)

func load(t *testing.T) *Atlas {
	t.Helper()
	once.Do(func() { atl, lerr = Load("../../data/json", "../../data/journeys.json") })
	if lerr != nil {
		t.Fatalf("load: %v", lerr)
	}
	return atl
}

func TestParseYear(t *testing.T) {
	cases := []struct {
		in    string
		year  float64
		label string
	}{
		{"-4003", -4003, "4003 BC"},
		{"0060", 60, "60 AD"},
		{"33", 33, "33 AD"},
		{"-1446", -1446, "1446 BC"},
	}
	for _, c := range cases {
		y, l := parseYear(c.in)
		if y != c.year || l != c.label {
			t.Errorf("parseYear(%q) = %v, %q; want %v, %q", c.in, y, l, c.year, c.label)
		}
	}
	// a dated year moves forward in time within the year
	y, _ := parseYear("0033-04-03")
	if y <= 33 || y >= 34 {
		t.Errorf("parseYear dated = %v, want within (33,34)", y)
	}
}

func TestParseDays(t *testing.T) {
	cases := map[string]int{"6D": 6, "2Y": 730, "3M": 90, "1W": 7, "": 0, "1M2D": 32}
	for in, want := range cases {
		if got := parseDays(in); got != want {
			t.Errorf("parseDays(%q) = %d, want %d", in, got, want)
		}
	}
}

func TestFormatYear(t *testing.T) {
	if FormatYear(-1) != "1 BC" || FormatYear(1) != "1 AD" || FormatYear(0) != "1 BC" {
		t.Errorf("unexpected: %s %s %s", FormatYear(-1), FormatYear(1), FormatYear(0))
	}
	if Pct(-4100) != 0 || Pct(100) != 100 || Pct(-9999) != 0 {
		t.Errorf("Pct bounds wrong")
	}
}

func TestRichIsSafe(t *testing.T) {
	got := RichInline(`<script>x</script> [God](/person/rec123) and [bad](javascript:alert(1)) _of_ [dead](/person/daughter_of_lot_-_older_985)`)
	if strings.Contains(got, "<script>") || strings.Contains(got, "javascript:") {
		t.Errorf("unsafe output: %s", got)
	}
	if !strings.Contains(got, `<a class="k-person" href="/person/rec123">God</a>`) {
		t.Errorf("missing link: %s", got)
	}
	if !strings.Contains(got, "<em>of</em>") {
		t.Errorf("missing emphasis: %s", got)
	}
	if strings.Contains(got, "daughter_of_lot") {
		t.Errorf("malformed href should be reduced to its label: %s", got)
	}
}

func TestParseTermsAndWordMatch(t *testing.T) {
	terms := parseTerms(`love "one another" Grace`)
	want := []string{"one another", "love", "grace"}
	if len(terms) != 3 || terms[0] != want[0] || terms[1] != want[1] || terms[2] != want[2] {
		t.Errorf("terms = %v, want %v", terms, want)
	}
	if indexWord("he is gone", "one", 0) != -1 {
		t.Error("'one' must not match inside 'gone'")
	}
	if indexWord("one and loved", "love", 0) < 0 {
		t.Error("'love' should match the start of 'loved'")
	}
	if got := highlight("He loved them", []string{"love"}); got != "He <mark>loved</mark> them" {
		t.Errorf("whole-word highlight = %q", got)
	}
	if got := highlight("Love one another", []string{"love"}); got != "<mark>Love</mark> one another" {
		t.Errorf("highlight = %q", got)
	}
}

func TestConcordance(t *testing.T) {
	a := load(t)
	r := a.Concordance("love one another", "", 100)
	if r.Total < 10 || len(r.Books) != len(a.Books) {
		t.Fatalf("total=%d books=%d", r.Total, len(r.Books))
	}
	for _, h := range r.Hits {
		if !strings.Contains(h.HTML, "<mark>") {
			t.Errorf("hit without highlight: %s", h.Label)
		}
	}
	if got := a.Concordance("x", "", 10); got.Total != 0 {
		t.Error("too-short queries should return nothing")
	}
}

func TestSearch(t *testing.T) {
	a := load(t)
	hits := a.Search("John 3:16", 10)
	if len(hits) == 0 || hits[0].Type != "verse" || hits[0].Label != "John 3:16" {
		t.Fatalf("verse reference not resolved: %+v", hits)
	}
	hits = a.Search("abraham", 10)
	if len(hits) == 0 || hits[0].Label != "Abraham" || hits[0].Type != "person" {
		t.Errorf("expected Abraham first: %+v", hits)
	}
}

func TestActiveYears(t *testing.T) {
	a := load(t)
	for _, p := range a.People {
		if p.Name == "David" {
			if p.Y0 > -900 || p.Y0 < -1200 {
				t.Errorf("David estimated years %d..%d", p.Y0, p.Y1)
			}
			return
		}
	}
	t.Fatal("David not found")
}

func TestJourneys(t *testing.T) {
	a := load(t)
	if len(a.Journeys) < 8 {
		t.Fatalf("expected the curated journeys, got %d", len(a.Journeys))
	}
	for _, j := range a.Journeys {
		if len(j.Stops) < 2 || j.Km <= 0 {
			t.Errorf("journey %s: stops=%d km=%d", j.ID, len(j.Stops), j.Km)
		}
		for _, s := range j.Stops {
			if s.VerseID == "" || s.Text == "" {
				t.Errorf("journey %s stop %s has no resolved verse", j.ID, s.Name)
			}
		}
	}
	var paul *Person
	for _, p := range a.People {
		if p.Name == "Paul" {
			paul = p
		}
	}
	if paul == nil {
		t.Fatal("Paul not found")
	}
	if j := a.PersonJourney(paul.ID); j == nil || len(j.Stops) < 10 {
		t.Errorf("derived journey for Paul too short: %+v", j)
	}
}

func TestDossiersDoNotPanic(t *testing.T) {
	a := load(t)
	for i, p := range a.People {
		if i%20 == 0 && a.PersonDossier(p.ID) == nil {
			t.Fatalf("person dossier nil: %s", p.ID)
		}
	}
	for _, p := range a.Places {
		if a.PlaceDossier(p.ID) == nil {
			t.Fatalf("place dossier nil: %s", p.ID)
		}
	}
	for _, e := range a.Events {
		if a.EventDossier(e.ID) == nil {
			t.Fatalf("event dossier nil: %s", e.ID)
		}
	}
	for i, v := range a.Verses {
		if i%100 == 0 && a.VerseDossier(v.ID) == nil {
			t.Fatalf("verse dossier nil: %s", v.ID)
		}
	}
	for _, b := range a.Books {
		if a.BookDossier(b.ID) == nil {
			t.Fatalf("book dossier nil: %s", b.ID)
		}
		for i, c := range b.Chapters {
			if i%7 == 0 && a.ChapterDossier(c.ID) == nil {
				t.Fatalf("chapter dossier nil: %s %d", b.Name, c.Num)
			}
		}
	}
	if len(a.BooksBrowse()) == 0 || len(a.EventsTimeline()) == 0 {
		t.Error("browse data empty")
	}
}

func TestCosmosLayout(t *testing.T) {
	a := load(t)
	c := a.Cosmos()
	n := len(c.IDs)
	if n < 1000 || len(c.X) != n || len(c.Y) != n || len(c.Z) != n || len(c.Edges)%3 != 0 {
		t.Fatalf("cosmos shape: n=%d x=%d edges=%d", n, len(c.X), len(c.Edges))
	}
	for k := 0; k < len(c.Edges); k += 3 {
		if int(c.Edges[k]) >= n || int(c.Edges[k+1]) >= n {
			t.Fatalf("edge out of range")
		}
	}
}

func TestEraCoverage(t *testing.T) {
	for y := -4100; y <= 100; y += 25 {
		e := eraAt(y)
		if y < e.From && e.ID != Eras[0].ID && y >= e.To {
			t.Errorf("year %d mapped to wrong era %s", y, e.ID)
		}
	}
}

func TestMigrationAndStory(t *testing.T) {
	a := load(t)
	m := a.Migration()
	if len(m.Places) < 500 || len(m.Counts) == 0 || len(m.Groups) < 5 {
		t.Fatalf("migration: places=%d counts=%d groups=%d", len(m.Places), len(m.Counts), len(m.Groups))
	}
	for _, c := range m.Counts {
		if c[0] >= len(m.Places) || c[1] < 0 || c[1] >= m.N || c[2] <= 0 {
			t.Fatalf("bad count triple %v", c)
		}
	}
	st := a.Story()
	if len(st.Eras) != len(Eras) && len(st.Eras) < 8 {
		t.Errorf("story eras = %d", len(st.Eras))
	}
	for i, e := range st.Eras {
		if st.Steps[e.Start].Kind != "era" {
			t.Errorf("era %d does not start with an intro step", i)
		}
		if e.Intro == "" {
			t.Errorf("era %s has no intro text", e.ID)
		}
	}
	if len(st.Steps) != len(a.Events)+len(st.Eras) {
		t.Errorf("steps = %d, want %d", len(st.Steps), len(a.Events)+len(st.Eras))
	}
}

func TestFamilyAndRelations(t *testing.T) {
	a := load(t)
	find := func(name string) *Person {
		var best *Person
		for _, p := range a.People {
			if p.Name == name && (best == nil || p.Verses > best.Verses) {
				best = p
			}
		}
		if best == nil {
			t.Fatalf("%s not found", name)
		}
		return best
	}
	abraham, jacob, esau, adam, jesus := find("Abraham"), find("Jacob (Israel)"), find("Esau (Edom)"), find("Adam"), find("Jesus Christ")
	if r := a.FamilyPath(abraham.ID, jacob.ID); r == nil || r.Label != "grandson" {
		t.Errorf("Abraham -> Jacob: %+v", r)
	}
	if r := a.FamilyPath(jacob.ID, esau.ID); r == nil || r.Label != "brother" {
		t.Errorf("Jacob -> Esau: %+v", r)
	}
	if r := a.FamilyPath(adam.ID, jesus.ID); r == nil || len(r.Steps) < 30 {
		t.Errorf("Adam -> Jesus should be a long chain: %+v", r)
	}
	if a.FamilyPath(abraham.ID, abraham.ID) != nil {
		t.Error("no path to oneself")
	}
	cmp := a.Compare("person", abraham.ID, jacob.ID)
	if cmp == nil || cmp.SharedCount == 0 || len(cmp.Story) < 3 || cmp.Relation == nil {
		t.Errorf("compare people: %+v", cmp)
	}
	if a.Compare("place", a.Places[0].ID, a.Places[1].ID) == nil {
		t.Error("compare places returned nil")
	}
	if a.Compare("nope", "x", "y") != nil {
		t.Error("unknown kind should be nil")
	}
}

func TestRelationLabel(t *testing.T) {
	cases := []struct {
		order         string
		u, d, s       int
		gender, label string
	}{
		{"u", 1, 0, 0, "Male", "father"},
		{"uu", 2, 0, 0, "Female", "grandmother"},
		{"uuuuu", 5, 0, 0, "Male", "great-great-great-grandfather"},
		{"uuuuuuu", 7, 0, 0, "Male", "5th-great-grandfather"},
		{"dd", 0, 2, 0, "Male", "grandson"},
		{"ud", 1, 1, 0, "Female", "sister"},
		{"udd", 1, 2, 0, "Male", "nephew"},
		{"uud", 2, 1, 0, "Female", "aunt"},
		{"uudd", 2, 2, 0, "", "1st cousin"},
		{"uuudd", 3, 2, 0, "", "1st cousin 1× removed"},
		{"s", 0, 0, 1, "Female", "wife"},
	}
	for _, c := range cases {
		if got := relationLabel(c.order, c.u, c.d, c.s, c.gender); got != c.label {
			t.Errorf("%s: got %q, want %q", c.order, got, c.label)
		}
	}
}

func TestThemes(t *testing.T) {
	a := load(t)
	list := a.Themes()
	if len(list) < 15 {
		t.Fatalf("themes: %d", len(list))
	}
	for _, th := range list {
		if th.Verses < 20 {
			t.Errorf("theme %s matches only %d verses", th.ID, th.Verses)
		}
		d := a.Theme(th.ID)
		if d == nil || len(d.Sample) == 0 || len(d.Books) != len(a.Books) {
			t.Errorf("theme %s detail incomplete", th.ID)
		}
	}
	if a.Theme("nope") != nil {
		t.Error("unknown theme should be nil")
	}
	if !matchStem("a good wife", "wife$") || matchStem("housewife", "wife$") || !matchStem("he forgave them", "forgav") && !matchStem("he forgave them", "forgave") {
		t.Error("stem matching")
	}
}

func TestExportAndProgress(t *testing.T) {
	a := load(t)
	j := a.journeyByKey("exodus")
	if j == nil {
		t.Fatal("exodus journey missing")
	}
	g := gpx(j)
	if !strings.Contains(g, "<rte>") || strings.Count(g, "<wpt") != len(j.Stops) || !strings.Contains(g, "Mount Sinai") {
		t.Errorf("gpx looks wrong")
	}
	k := kml(j)
	if !strings.Contains(k, "<LineString>") || strings.Count(k, "<Placemark>") != len(j.Stops)+1 {
		t.Errorf("kml looks wrong")
	}
	if strings.Contains(gpx(&Journey{Title: "A & <B>", Stops: []Stop{{Name: "x<y"}}}), "<B>") {
		t.Error("xml must be escaped")
	}
	p := a.ProgressIndex()
	if len(p.Books) != 66 || len(p.Places) < 500 {
		t.Errorf("progress index: books=%d places=%d", len(p.Books), len(p.Places))
	}
	n := 0
	for _, b := range p.Books {
		n += len(b.Chapters)
	}
	if n != 1189 {
		t.Errorf("chapters = %d, want 1189", n)
	}
}

func TestCorrectionsApplied(t *testing.T) {
	a := load(t)
	j := a.PersonByID["recsyU2HqpyE9PeTo"] // Jehoram of Judah
	if j == nil || j.Born == nil || *j.Born != -884 || j.Died == nil || *j.Died != -844 {
		t.Errorf("Jehoram corrections not applied: %+v", j)
	}
	if s := a.PersonByID["recwcML6FMulTblKH"]; s == nil || s.Born != nil { // Samson: bad birth year cleared
		t.Errorf("Samson birth year should be cleared")
	}
	if a.Enrich != nil {
		for _, art := range a.Enrich.Art {
			if strings.Contains(strings.ToLower(art.Title), "driven out of eden") && art.Event != "recj3SioaPwrDtpGT" {
				t.Errorf("art fix not applied to %q: event %q", art.Title, art.Event)
			}
		}
	}
}
