package atlas

import (
	"os"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/almerlucke/metabible/pkg/i18n"
)

var (
	once sync.Once
	atl  *Atlas
	lerr error
)

func load(t *testing.T) *Atlas {
	t.Helper()
	once.Do(func() {
		atl, lerr = Load("../../data/json", "../../data/journeys.json")
		if lerr == nil {
			atl.SetTranslations(i18n.MustLoad())
		}
	})
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
	r := a.Concordance("love one another", "", 100, "en")
	if r.Total < 10 || len(r.Books) != len(a.Books) {
		t.Fatalf("total=%d books=%d", r.Total, len(r.Books))
	}
	for _, h := range r.Hits {
		if !strings.Contains(h.HTML, "<mark>") {
			t.Errorf("hit without highlight: %s", h.Label)
		}
	}
	if got := a.Concordance("x", "", 10, "en"); got.Total != 0 {
		t.Error("too-short queries should return nothing")
	}
}

func TestSearch(t *testing.T) {
	a := load(t)
	hits := a.Search("John 3:16", 10, "en")
	if len(hits) == 0 || hits[0].Type != "verse" || hits[0].Label != "John 3:16" {
		t.Fatalf("verse reference not resolved: %+v", hits)
	}
	hits = a.Search("abraham", 10, "en")
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
	if j := a.PersonJourney(paul.ID, "en"); j == nil || len(j.Stops) < 10 {
		t.Errorf("derived journey for Paul too short: %+v", j)
	}
}

func TestDossiersDoNotPanic(t *testing.T) {
	a := load(t)
	for i, p := range a.People {
		if i%20 == 0 && a.PersonDossier(p.ID, "en") == nil {
			t.Fatalf("person dossier nil: %s", p.ID)
		}
	}
	for _, p := range a.Places {
		if a.PlaceDossier(p.ID, "en") == nil {
			t.Fatalf("place dossier nil: %s", p.ID)
		}
	}
	for _, e := range a.Events {
		if a.EventDossier(e.ID, "en") == nil {
			t.Fatalf("event dossier nil: %s", e.ID)
		}
	}
	for i, v := range a.Verses {
		if i%100 == 0 && a.VerseDossier(v.ID, "en") == nil {
			t.Fatalf("verse dossier nil: %s", v.ID)
		}
	}
	for _, b := range a.Books {
		if a.BookDossier(b.ID, "en") == nil {
			t.Fatalf("book dossier nil: %s", b.ID)
		}
		for i, c := range b.Chapters {
			if i%7 == 0 && a.ChapterDossier(c.ID, "en") == nil {
				t.Fatalf("chapter dossier nil: %s %d", b.Name, c.Num)
			}
		}
	}
	if len(a.BooksBrowse("en")) == 0 || len(a.EventsTimeline("en")) == 0 {
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
	st := a.Story("en")
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
	if r := a.FamilyPath(abraham.ID, jacob.ID, "en"); r == nil || r.Label != "grandson" {
		t.Errorf("Abraham -> Jacob: %+v", r)
	}
	if r := a.FamilyPath(jacob.ID, esau.ID, "en"); r == nil || r.Label != "brother" {
		t.Errorf("Jacob -> Esau: %+v", r)
	}
	if r := a.FamilyPath(adam.ID, jesus.ID, "en"); r == nil || len(r.Steps) < 30 {
		t.Errorf("Adam -> Jesus should be a long chain: %+v", r)
	}
	if a.FamilyPath(abraham.ID, abraham.ID, "en") != nil {
		t.Error("no path to oneself")
	}
	cmp := a.Compare("person", abraham.ID, jacob.ID, "en")
	if cmp == nil || cmp.SharedCount == 0 || len(cmp.Story) < 3 || cmp.Relation == nil {
		t.Errorf("compare people: %+v", cmp)
	}
	if a.Compare("place", a.Places[0].ID, a.Places[1].ID, "en") == nil {
		t.Error("compare places returned nil")
	}
	if a.Compare("nope", "x", "y", "en") != nil {
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
	list := a.Themes("en")
	if len(list) < 15 {
		t.Fatalf("themes: %d", len(list))
	}
	for _, th := range list {
		if th.Verses < 20 {
			t.Errorf("theme %s matches only %d verses", th.ID, th.Verses)
		}
		d := a.Theme(th.ID, "en")
		if d == nil || len(d.Sample) == 0 || len(d.Books) != len(a.Books) {
			t.Errorf("theme %s detail incomplete", th.ID)
		}
	}
	if a.Theme("nope", "en") != nil {
		t.Error("unknown theme should be nil")
	}
	if !matchStem("a good wife", "wife$") || matchStem("housewife", "wife$") || !matchStem("he forgave them", "forgav") && !matchStem("he forgave them", "forgave") {
		t.Error("stem matching")
	}
}

func TestExportAndProgress(t *testing.T) {
	a := load(t)
	j := a.journeyByKey("exodus", "en")
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
	p := a.ProgressIndex("en")
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
			if strings.Contains(strings.ToLower(art.Title), "expelled from paradise") {
				t.Errorf("duplicate engraving should have been dropped: %q", art.Title)
			}
			if strings.Contains(strings.ToLower(art.Title), "driven out of eden") && art.Event != "recj3SioaPwrDtpGT" {
				t.Errorf("art fix not applied to %q: event %q", art.Title, art.Event)
			}
		}
	}
}

func TestDropDuplicates(t *testing.T) {
	dir := t.TempDir()
	f := dir + "/art-duplicates.json"
	if err := os.WriteFile(f, []byte(`{"duplicates":[{"title":"B","of":"A","distance":2}]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := dropDuplicates([]Art{{Image: Image{Title: "A"}}, {Image: Image{Title: "B"}}, {Image: Image{Title: "C"}}}, f)
	if err != nil || len(got) != 2 || got[0].Title != "A" || got[1].Title != "C" {
		t.Errorf("got %+v, %v", got, err)
	}
	if same, _ := dropDuplicates([]Art{{}}, dir+"/missing.json"); len(same) != 1 {
		t.Error("a missing file means no duplicates")
	}
}

func TestArtNameKeys(t *testing.T) {
	a := load(t)
	keys := a.artNameKeys()
	jesus, peter := "", ""
	for _, p := range a.People {
		if p.Name == "Jesus Christ" {
			jesus = p.ID
		}
		if p.Name == "Simon Peter" {
			peter = p.ID
		}
	}
	if keys["jesus"] != jesus || keys["christ"] != jesus || keys["peter"] != peter {
		t.Errorf("distinctive names not mapped: jesus=%q christ=%q peter=%q", keys["jesus"], keys["christ"], keys["peter"])
	}
	for _, generic := range []string{"king", "angel", "god", "lord", "israel"} {
		if _, ok := keys[generic]; ok {
			t.Errorf("%q must never identify a person", generic)
		}
	}
	if got := peopleInTitle("Christ Before Pilate", keys); len(got) == 0 || got[0] != jesus {
		t.Errorf("Christ Before Pilate -> %v", got)
	}
}

// ATLAS_REPORT=1 go test ./pkg/atlas -run TestArtReport -v  lists who has engravings and why.
func TestArtReport(t *testing.T) {
	if os.Getenv("ATLAS_REPORT") == "" {
		t.Skip("set ATLAS_REPORT=1 to print")
	}
	a := load(t)
	type row struct {
		name string
		n    int
	}
	var rows []row
	for id, arts := range a.artByPerson {
		rows = append(rows, row{a.PersonByID[id].Name, len(arts)})
	}
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].n > rows[j].n || rows[i].n == rows[j].n && rows[i].name < rows[j].name
	})
	t.Logf("%d people have engravings", len(rows))
	for _, r := range rows {
		t.Logf("%-24s %d", r.name, r.n)
	}
	for _, name := range []string{"Jesus Christ", "Simon Peter", "Mary (Mother of Jesus)", "John the Baptist", "Pontius Pilate"} {
		for _, p := range a.People {
			if p.Name == name {
				for _, art := range a.artByPerson[p.ID] {
					t.Logf("  %s <- %s", name, art.Title)
				}
			}
		}
	}
}

func TestGallery(t *testing.T) {
	if got := sized("https://commons.wikimedia.org/wiki/Special:FilePath/A%20b.jpg?width=960", 480); !strings.HasSuffix(got, "width=480") {
		t.Errorf("FilePath sizing: %s", got)
	}
	if got := sized("https://upload.wikimedia.org/wikipedia/commons/thumb/8/86/X.jpg/500px-X.jpg", 1280); !strings.Contains(got, "/1280px-X.jpg") {
		t.Errorf("thumb sizing: %s", got)
	}
	a := load(t)
	items := a.AllGallery()
	if a.Enrich != nil && len(items) < 100 {
		t.Fatalf("gallery too small: %d", len(items))
	}
	for _, it := range items {
		if it.Thumb == "" || it.Full == "" || it.Title == "" || it.Credit == "" || it.Page == "" {
			t.Errorf("incomplete gallery item: %+v", it)
		}
		if !strings.HasPrefix(it.LinksJSON(), "[") {
			t.Errorf("links json: %s", it.LinksJSON())
		}
	}
}

func TestCleanArtTitle(t *testing.T) {
	cases := map[string]string{
		"Abraham Journeying into the Land of Canaan (89393798)":   "Abraham Journeying into the Land of Canaan",
		"Gustave Dore - John the Baptist baptizes Jesus":          "John the Baptist baptizes Jesus",
		"Gustave Doré, St Peter Denying Christ":                   "St Peter Denying Christ",
		"Gustave Doré - Christ on the Cross - Google Art Project": "Christ on the Cross",
		"By Dore, Gustave; La Sainte Trinite":                     "By Dore, Gustave; La Sainte Trinite",
		"Gustave Doré":                                            "Engraving",
	}
	for in, want := range cases {
		if got := cleanArtTitle(in); got != want {
			t.Errorf("%q -> %q, want %q", in, got, want)
		}
	}
}

func TestSplitAliases(t *testing.T) {
	got := splitAliases("Lord,Son,Lamb,JESUS,LORD,lord,Jesus Christ ,,Messiah", "Jesus Christ")
	want := []string{"Lord", "Son", "Lamb", "JESUS", "Messiah"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("got %v, want %v", got, want)
		}
	}
	if len(splitAliases("", "X")) != 0 {
		t.Error("empty input has no aliases")
	}
}

func TestDutch(t *testing.T) {
	a := load(t)
	var abraham, jacob *Person
	for _, p := range a.People {
		if p.Name == "Abraham" {
			abraham = p
		}
		if p.Name == "Jacob (Israel)" {
			jacob = p
		}
	}
	rel := a.FamilyPath(abraham.ID, jacob.ID, "nl")
	if rel == nil || rel.Label != "de kleinzoon" || !strings.Contains(rel.Phrase, "is de kleinzoon van Abraham") {
		t.Errorf("dutch relation: %+v", rel)
	}
	cmp := a.Compare("person", abraham.ID, jacob.ID, "nl")
	if cmp == nil || !strings.Contains(strings.Join(cmp.Story, " "), "v.Chr.") || strings.Contains(strings.Join(cmp.Story, " "), " BC") {
		t.Errorf("dutch story should use v.Chr.: %v", cmp.Story)
	}
	st := a.Story("nl")
	if !strings.HasPrefix(st.Eras[0].Name, "Oergeschiedenis") || !strings.Contains(st.Eras[0].Intro, "tuin") {
		t.Errorf("dutch story era: %+v", st.Eras[0])
	}
	js := a.JourneysFor("nl")
	if js[0].Title != "Abrahams roeping" || js[1].Stops[0].Name != "Rameses" || !strings.Contains(js[0].Stops[1].Note, "Haran") {
		t.Errorf("dutch journey: %s / %+v", js[0].Title, js[1].Stops[0])
	}
	if a.JourneysFor("en")[0].Title != "Abraham's Calling" {
		t.Error("english journeys unchanged")
	}
	if th := a.Themes("nl"); th[0].Name != "Verbond" {
		t.Errorf("dutch themes: %v", th[0])
	}
	hits := a.Search("abraham", 3, "nl")
	if len(hits) == 0 || !strings.HasPrefix(hits[0].Sub, "persoon") || !strings.Contains(hits[0].Sub, "v.Chr.") {
		t.Errorf("dutch search sub: %+v", hits)
	}
	pk := a.PeekFor("person", abraham.ID, "nl")
	if pk == nil || !strings.Contains(strings.Join(pk.Facts, " "), "verzen") {
		t.Errorf("dutch peek: %+v", pk)
	}
	if got := a.locDate("nl", "1446 BC – 30 AD"); got != "1446 v.Chr. – 30 n.Chr." {
		t.Errorf("locDate: %q", got)
	}
}

func TestDutchBible(t *testing.T) {
	a := load(t)
	if !a.HasBible("nl") {
		t.Fatal("no Dutch Bible loaded")
	}
	v := a.VerseByOsis["gen.1.1"]
	if got := a.VerseIn(v, "nl"); !strings.HasPrefix(got, "In den beginne") {
		t.Errorf("Genesis 1:1 = %q", got)
	}
	if got := a.VerseIn(v, "en"); !strings.HasPrefix(got, "In the beginning") {
		t.Errorf("English changed: %q", got)
	}
	// Job 38:39 is numbered differently in the Statenvertaling: it falls back to English
	if v := a.VerseByOsis["job.38.39"]; v != nil && a.VerseIn(v, "nl") != v.Text {
		t.Error("missing verse did not fall back")
	}
	// concordance searches the Dutch text
	r := a.Concordance("in den beginne", "", 10, "nl")
	if r.Total == 0 {
		t.Error("Dutch concordance found nothing")
	}
	if a.Concordance("in den beginne", "", 10, "en").Total != 0 {
		t.Error("English concordance matched Dutch words")
	}
	d := a.VerseDossier(v.ID, "nl")
	if !strings.Contains(d.HTML, "In den beginne") {
		t.Errorf("verse page not Dutch: %s", d.HTML)
	}
	covered := 0
	for _, v := range a.Verses {
		if a.VerseIn(v, "nl") != v.Text {
			covered++
		}
	}
	if covered < 30000 {
		t.Errorf("only %d verses differ from KJV", covered)
	}
}

func TestDutchBookNames(t *testing.T) {
	a := load(t)
	b := a.findBook("richteren")
	if b == nil || b.Name != "Judges" {
		t.Fatalf("Dutch book name not understood: %+v", b)
	}
	if got := a.BookName(b, "nl"); got != "Richteren" {
		t.Errorf("BookName = %q", got)
	}
	if got := a.BookName(b, "en"); got != "Judges" {
		t.Errorf("English changed: %q", got)
	}
	hits := a.Search("Openbaring 1:1", 5, "nl")
	if len(hits) == 0 || !strings.HasPrefix(hits[0].Label, "Openbaring") {
		t.Errorf("reference in Dutch: %+v", hits)
	}
	hits = a.Search("psalmen", 5, "nl")
	if len(hits) == 0 || hits[0].Label != "Psalmen" {
		t.Errorf("book search in Dutch: %+v", hits)
	}
	if got := a.verseLabel(a.VerseByOsis["john.3.16"], "nl"); got != "Johannes 3:16" {
		t.Errorf("label = %q", got)
	}
	for _, b := range a.BooksIn("nl") {
		if b.Name == "" {
			t.Fatal("empty name")
		}
	}
	if a.BooksIn("nl")[18].Name != "Psalmen" || a.Books[18].Name != "Psalms" {
		t.Error("BooksIn must copy, not rename the shared books")
	}
}

func TestDutchNames(t *testing.T) {
	a := load(t)
	nl := a.View("nl")
	if nl == a {
		t.Fatal("no Dutch view")
	}
	if a.View("en") != a || a.View("xx") != a {
		t.Error("languages without a name table must use the atlas itself")
	}
	byName := func(v *Atlas, name string) *Person {
		for _, p := range v.People {
			if p.Name == name && p.Verses > 20 {
				return p
			}
		}
		return nil
	}
	if byName(a, "Isaac") == nil || byName(nl, "Izak") == nil || byName(nl, "Isaac") != nil {
		t.Error("Isaac should be Izak in the Dutch view only")
	}
	var jer *Place
	for _, p := range nl.Places {
		if p.Orig == "Jerusalem" {
			jer = p
		}
	}
	if jer == nil || jer.Name != "Jeruzalem" || nl.PlaceByID[jer.ID] != jer || a.PlaceByID[jer.ID].Name != "Jerusalem" {
		t.Errorf("Jerusalem: %+v", jer)
	}
	// the English name still finds the person
	found := false
	for _, h := range nl.Search("isaac", 10, "nl") {
		if h.Label == "Izak" {
			found = true
		}
	}
	if !found {
		t.Error("searching the English name in Dutch should find Izak")
	}
	if hits := nl.Search("jeruzalem", 3, "nl"); len(hits) == 0 || hits[0].Label != "Jeruzalem" {
		t.Errorf("Dutch place search: %+v", hits)
	}
	d := nl.PersonDossier(byName(nl, "Izak").ID, "nl")
	if d.P.Name != "Izak" || d.Father.Name != "Abraham" {
		t.Errorf("dossier: %q father %v", d.P.Name, d.Father)
	}
	// a view must not share caches with the base
	if nl.Story("nl") == a.Story("en") {
		t.Error("story shared between views")
	}
}

func TestDutchEventTitles(t *testing.T) {
	a := load(t)
	nl := a.View("nl")
	same := 0
	for _, e := range a.Events {
		d := nl.EventByID[e.ID]
		if d == nil || d == e {
			t.Fatalf("event %s missing from the Dutch view", e.ID)
		}
		if d.Title == e.Title {
			same++
		}
		if d.Orig != e.Title && d.Orig != "" {
			t.Errorf("Orig %q vs %q", d.Orig, e.Title)
		}
	}
	if same > 10 {
		t.Errorf("%d event titles untranslated", same)
	}
	if hits := nl.Search("sodom verwoest", 3, "nl"); len(hits) == 0 || hits[0].Type != "event" {
		t.Errorf("Dutch title search: %+v", hits)
	}
	if hits := nl.Search("sodom destroyed", 3, "nl"); len(hits) == 0 {
		t.Error("English title should still be found")
	}
	// the cross references of a person point at the translated events
	for _, p := range nl.People {
		if es := nl.personEvents[p.ID]; len(es) > 0 {
			if nl.EventByID[es[0].ID] != es[0] {
				t.Fatal("personEvents still points at the source events")
			}
			break
		}
	}
	lives := 0
	for _, e := range a.Events {
		if e.Life {
			lives++
		}
	}
	if lives < 100 {
		t.Errorf("only %d lifespan events flagged", lives)
	}
}
