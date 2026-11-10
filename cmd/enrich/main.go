// enrich pulls open, properly licensed content for the most prominent places and
// people (Wikipedia text under CC BY-SA, Wikimedia Commons images with free
// licenses) and Gustave Doré's public-domain Bible engravings, and writes them
// to data/enrichment.json together with their sources and licenses.
//
// It only sends public names and coordinates, one request at a time, with a
// generic client name. Re-running resumes: entries already present are kept.
//
//	go run ./cmd/enrich            # places, people and Doré
//	go run ./cmd/enrich -limit 30  # quick trial
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/almerlucke/metabible/pkg/atlas"
)

const ua = "MetaBibleEnrich/1.0 (open-source Bible atlas)"

type Image struct {
	URL     string `json:"url"`
	Page    string `json:"page"`
	Title   string `json:"title,omitempty"`
	License string `json:"license"`
	Artist  string `json:"artist,omitempty"`
	Credit  string `json:"credit,omitempty"`
	W       int    `json:"w,omitempty"`
	H       int    `json:"h,omitempty"`
}

type Entry struct {
	Title   string `json:"title"`
	Summary string `json:"summary"`
	Desc    string `json:"desc,omitempty"`
	Page    string `json:"page"`
	Image   *Image `json:"image,omitempty"`
}

type Art struct {
	Image
	Event  string   `json:"event,omitempty"`
	People []string `json:"people,omitempty"`
	Score  float64  `json:"score"`
}

type Output struct {
	Generated string            `json:"generated"`
	Sources   map[string]string `json:"sources"`
	Places    map[string]Entry  `json:"places"`
	People    map[string]Entry  `json:"people"`
	Art       []Art             `json:"art"`
}

var client = &http.Client{Timeout: 25 * time.Second}
var lastCall time.Time

func get(u string, v any) error {
	for attempt := 0; attempt < 8; attempt++ {
		gap := 1100 * time.Millisecond // the Action API is strict
		if strings.Contains(u, "/api/rest_v1/") {
			gap = 300 * time.Millisecond
		}
		if d := time.Since(lastCall); d < gap {
			time.Sleep(gap - d)
		}
		lastCall = time.Now()
		req, _ := http.NewRequest("GET", u, nil)
		req.Header.Set("User-Agent", ua)
		resp, err := client.Do(req)
		if err != nil {
			time.Sleep(time.Second * time.Duration(attempt+1))
			continue
		}
		b, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode == 429 || resp.StatusCode >= 500 {
			wait := 10 * time.Second
			if ra, err := strconv.Atoi(resp.Header.Get("Retry-After")); err == nil && ra > 0 {
				wait = time.Duration(ra+2) * time.Second
			}
			if wait > 3*time.Minute {
				wait = 3 * time.Minute
			}
			log.Printf("rate limited (%d), waiting %s", resp.StatusCode, wait)
			time.Sleep(wait)
			continue
		}
		if resp.StatusCode == 404 {
			return os.ErrNotExist
		}
		if resp.StatusCode != 200 {
			return fmt.Errorf("%s: %s", u, resp.Status)
		}
		return json.Unmarshal(b, v)
	}
	return fmt.Errorf("giving up on %s", u)
}

// ---------------------------------------------------------------- wikipedia

type summary struct {
	Type        string `json:"type"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Extract     string `json:"extract"`
	Thumbnail   *struct {
		Source string `json:"source"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	} `json:"thumbnail"`
	OriginalImage *struct {
		Source string `json:"source"`
	} `json:"originalimage"`
	ContentURLs struct {
		Desktop struct {
			Page string `json:"page"`
		} `json:"desktop"`
	} `json:"content_urls"`
	Coordinates *struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"coordinates"`
}

func wikiSummary(title string) (*summary, error) {
	var s summary
	err := get("https://en.wikipedia.org/api/rest_v1/page/summary/"+url.PathEscape(strings.ReplaceAll(title, " ", "_")), &s)
	if err != nil {
		return nil, err
	}
	if s.Type != "standard" || s.Extract == "" {
		return nil, os.ErrNotExist
	}
	return &s, nil
}

func wikiSearch(q string) []string {
	var r struct {
		Query struct {
			Search []struct {
				Title string `json:"title"`
			} `json:"search"`
		} `json:"query"`
	}
	u := "https://en.wikipedia.org/w/api.php?" + url.Values{"action": {"query"}, "list": {"search"}, "srsearch": {q}, "srlimit": {"4"}, "format": {"json"}}.Encode()
	if get(u, &r) != nil {
		return nil
	}
	var out []string
	for _, s := range r.Query.Search {
		out = append(out, s.Title)
	}
	return out
}

var biblical = regexp.MustCompile(`(?i)\b(bible|biblical|hebrew bible|old testament|new testament|gospel|apostle|prophet|patriarch|israel|israelite|judah|jewish|torah|christian|jesus|canaan|levant|ancient)\b`)
var parenRe = regexp.MustCompile(`\s*\(.*?\)`)

func baseName(s string) string { return strings.TrimSpace(parenRe.ReplaceAllString(s, "")) }

func haversine(lat1, lng1, lat2, lng2 float64) float64 {
	p := math.Pi / 180
	dp, dl := (lat2-lat1)*p, (lng2-lng1)*p
	h := math.Sin(dp/2)*math.Sin(dp/2) + math.Cos(lat1*p)*math.Cos(lat2*p)*math.Sin(dl/2)*math.Sin(dl/2)
	return 12742 * math.Asin(math.Sqrt(h))
}

func tooPlain(s string) string {
	s = strings.Join(strings.Fields(s), " ")
	if len([]rune(s)) > 700 {
		r := []rune(s)[:700]
		cut := string(r)
		if i := strings.LastIndex(cut, ". "); i > 300 {
			return cut[:i+1]
		}
		return cut + "…"
	}
	return s
}

// ---------------------------------------------------------------- commons images

type extMeta struct {
	License, Artist, Credit string
	Free                    bool
}

var unwanted = regexp.MustCompile(`(?i)(flag_of|coat_of_arms|coat of arms|seal_of|logo|locator|location_map|emblem|icon|symbol)`)

var tagRe = regexp.MustCompile(`<[^>]*>`)

func cleanMeta(s string) string {
	return strings.TrimSpace(html.UnescapeString(strings.Join(strings.Fields(tagRe.ReplaceAllString(s, "")), " ")))
}

var freeRe = regexp.MustCompile(`(?i)^(cc[ -]by|cc0|public domain|pd|attribution|gfdl|cc-by)`)

func fileMeta(file string) (*extMeta, string, error) {
	var r struct {
		Query struct {
			Pages map[string]struct {
				ImageInfo []struct {
					ExtMetadata map[string]struct {
						Value any `json:"value"`
					} `json:"extmetadata"`
				} `json:"imageinfo"`
			} `json:"pages"`
		} `json:"query"`
	}
	u := "https://commons.wikimedia.org/w/api.php?" + url.Values{"action": {"query"}, "titles": {"File:" + file}, "prop": {"imageinfo"}, "iiprop": {"extmetadata"}, "format": {"json"}}.Encode()
	if err := get(u, &r); err != nil {
		return nil, "", err
	}
	for _, p := range r.Query.Pages {
		if len(p.ImageInfo) == 0 {
			continue
		}
		m := p.ImageInfo[0].ExtMetadata
		val := func(k string) string {
			if v, ok := m[k]; ok {
				if s, ok := v.Value.(string); ok {
					return cleanMeta(s)
				}
			}
			return ""
		}
		em := &extMeta{License: val("LicenseShortName"), Artist: val("Artist"), Credit: val("Credit")}
		em.Free = em.License != "" && freeRe.MatchString(em.License) && val("NonFree") == ""
		return em, "https://commons.wikimedia.org/wiki/File:" + url.PathEscape(strings.ReplaceAll(file, " ", "_")), nil
	}
	return nil, "", os.ErrNotExist
}

// fileFromURL extracts the Commons file name from an upload.wikimedia.org URL.
func fileFromURL(u string) string {
	parts := strings.Split(u, "/")
	if i := indexOf(parts, "thumb"); i >= 0 && len(parts) > i+4 {
		return mustUnescape(parts[i+3])
	}
	return mustUnescape(parts[len(parts)-1])
}
func indexOf(a []string, s string) int {
	for i, v := range a {
		if v == s {
			return i
		}
	}
	return -1
}
func mustUnescape(s string) string {
	if u, err := url.PathUnescape(s); err == nil {
		return u
	}
	return s
}

func imageFor(s *summary) *Image {
	if s.Thumbnail == nil {
		return nil
	}
	src := s.Thumbnail.Source
	if s.OriginalImage != nil {
		src = s.OriginalImage.Source
	}
	file := fileFromURL(src)
	meta, page, err := fileMeta(file)
	if err != nil || !meta.Free {
		return nil
	}
	return &Image{URL: s.Thumbnail.Source, Page: page, License: meta.License, Artist: meta.Artist, Credit: meta.Credit, W: s.Thumbnail.Width, H: s.Thumbnail.Height}
}

// ---------------------------------------------------------------- batched wikipedia queries

type page struct {
	Title       string  `json:"title"`
	Extract     string  `json:"extract"`
	Description string  `json:"description"`
	FullURL     string  `json:"fullurl"`
	PageImage   string  `json:"pageimage"`
	Missing     *string `json:"missing"`
	Thumbnail   *struct {
		Source string `json:"source"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	} `json:"thumbnail"`
	Coordinates []struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"coordinates"`
	PageProps map[string]any `json:"pageprops"`
}

type batchResp struct {
	Query struct {
		Normalized []struct{ From, To string } `json:"normalized"`
		Redirects  []struct{ From, To string } `json:"redirects"`
		Pages      map[string]page             `json:"pages"`
	} `json:"query"`
}

var restCache = map[string]*page{}

// restPage reads one article through the REST summary endpoint (cached and generous with limits).
func restPage(title string) *page {
	if p, ok := restCache[title]; ok {
		return p
	}
	var r struct {
		Type        string `json:"type"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Extract     string `json:"extract"`
		Thumbnail   *struct {
			Source string `json:"source"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"thumbnail"`
		OriginalImage *struct {
			Source string `json:"source"`
		} `json:"originalimage"`
		ContentURLs struct {
			Desktop struct {
				Page string `json:"page"`
			} `json:"desktop"`
		} `json:"content_urls"`
		Coordinates *struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		} `json:"coordinates"`
	}
	err := get("https://en.wikipedia.org/api/rest_v1/page/summary/"+url.PathEscape(strings.ReplaceAll(title, " ", "_"))+"?redirect=true", &r)
	if err != nil || r.Type != "standard" || r.Extract == "" {
		restCache[title] = nil
		return nil
	}
	p := &page{Title: r.Title, Extract: r.Extract, Description: r.Description, FullURL: r.ContentURLs.Desktop.Page}
	if r.Thumbnail != nil {
		p.Thumbnail = &struct {
			Source string `json:"source"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		}{r.Thumbnail.Source, r.Thumbnail.Width, r.Thumbnail.Height}
		src := r.Thumbnail.Source
		if r.OriginalImage != nil {
			src = r.OriginalImage.Source
		}
		p.PageImage = fileFromURL(src)
	}
	if r.Coordinates != nil {
		p.Coordinates = []struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		}{{r.Coordinates.Lat, r.Coordinates.Lon}}
	}
	restCache[title] = p
	return p
}

// licensedImages looks up licenses for many Commons files at once.
func licensedImages(files []string) map[string]*Image {
	out := map[string]*Image{}
	for len(files) > 0 {
		n := len(files)
		if n > 40 {
			n = 40
		}
		chunk := files[:n]
		files = files[n:]
		var titles []string
		for _, f := range chunk {
			titles = append(titles, "File:"+f)
		}
		var r struct {
			Query struct {
				Normalized []struct{ From, To string } `json:"normalized"`
				Pages      map[string]struct {
					Title     string `json:"title"`
					ImageInfo []struct {
						ExtMetadata map[string]struct {
							Value any `json:"value"`
						} `json:"extmetadata"`
					} `json:"imageinfo"`
				} `json:"pages"`
			} `json:"query"`
		}
		u := "https://commons.wikimedia.org/w/api.php?" + url.Values{"action": {"query"}, "titles": {strings.Join(titles, "|")}, "prop": {"imageinfo"}, "iiprop": {"extmetadata"}, "format": {"json"}}.Encode()
		if err := get(u, &r); err != nil {
			log.Printf("license batch failed: %v", err)
			continue
		}
		norm := map[string]string{}
		for _, t := range titles {
			norm[t] = t
		}
		for _, nn := range r.Query.Normalized {
			for k, v := range norm {
				if v == nn.From {
					norm[k] = nn.To
				}
			}
		}
		byTitle := map[string]*Image{}
		for _, p := range r.Query.Pages {
			if len(p.ImageInfo) == 0 {
				continue
			}
			m := p.ImageInfo[0].ExtMetadata
			val := func(k string) string {
				if v, ok := m[k]; ok {
					if s, ok := v.Value.(string); ok {
						return cleanMeta(s)
					}
				}
				return ""
			}
			lic := val("LicenseShortName")
			if lic == "" || !freeRe.MatchString(lic) || val("NonFree") != "" {
				continue
			}
			byTitle[p.Title] = &Image{Page: "https://commons.wikimedia.org/wiki/" + url.PathEscape(strings.ReplaceAll(p.Title, " ", "_")), License: lic, Artist: val("Artist"), Credit: val("Credit")}
		}
		for i, f := range chunk {
			if img := byTitle[norm[titles[i]]]; img != nil {
				out[f] = img
			}
		}
	}
	return out
}

func attachImages(entries map[string]*Entry, pageOf map[string]*page) {
	var files []string
	owner := map[string][]string{}
	for id, p := range pageOf {
		if p.Thumbnail == nil || p.PageImage == "" || unwanted.MatchString(p.PageImage) {
			continue
		}
		files = append(files, p.PageImage)
		owner[p.PageImage] = append(owner[p.PageImage], id)
	}
	lic := licensedImages(files)
	for f, img := range lic {
		for _, id := range owner[f] {
			p := pageOf[id]
			c := *img
			c.URL, c.W, c.H = p.Thumbnail.Source, p.Thumbnail.Width, p.Thumbnail.Height
			entries[id].Image = &c
		}
	}
}

// ---------------------------------------------------------------- places and people

func placeCandidates(p *atlas.Place) []string {
	name := baseName(p.Name)
	return []string{name, name + " (city)", name + " (ancient city)", name + " (biblical city)", name + " (biblical region)", name + " (Bible)", name + " (town)"}
}

func personCandidates(p *atlas.Person) []string {
	name := baseName(p.Name)
	return []string{name + " (biblical figure)", name + " (Bible)", name + " (biblical)", name + " (prophet)", name + " (apostle)", name + " (patriarch)", name + " (king)", name}
}

// enrichBatch resolves items one by one (first accepted candidate wins), and attaches licensed
// images and reports results in groups so progress is saved as it goes.
func enrichBatch[T any](items []T, id func(T) string, cands func(T) []string, accept func(T, *page) bool, done map[string]Entry, emit func(map[string]*Entry)) {
	var todo []T
	for _, it := range items {
		if _, ok := done[id(it)]; !ok {
			todo = append(todo, it)
		}
	}
	log.Printf("%d to look up", len(todo))
	for start := 0; start < len(todo); start += 40 {
		end := start + 40
		if end > len(todo) {
			end = len(todo)
		}
		entries := map[string]*Entry{}
		pageOf := map[string]*page{}
		for _, it := range todo[start:end] {
			for _, c := range cands(it) {
				if p := restPage(c); p != nil && accept(it, p) {
					entries[id(it)] = &Entry{Title: p.Title, Summary: tooPlain(p.Extract), Desc: p.Description, Page: p.FullURL}
					pageOf[id(it)] = p
					log.Printf("  %s -> %s", id(it), p.Title)
					break
				}
			}
		}
		attachImages(entries, pageOf)
		emit(entries)
		log.Printf("progress %d/%d", end, len(todo))
	}
}

func acceptPlace(p *atlas.Place, pg *page) bool {
	if len(pg.Coordinates) == 0 {
		return false
	}
	return haversine(p.Lat, p.Lng, pg.Coordinates[0].Lat, pg.Coordinates[0].Lon) <= 150
}

func acceptPerson(p *atlas.Person, pg *page) bool {
	name := baseName(p.Name)
	if name == "" || name == "God" || len(strings.Fields(name)) == 0 {
		return false
	}
	first := strings.ToLower(strings.Fields(name)[0])
	txt := strings.ToLower(pg.Extract + " " + pg.Title)
	return biblical.MatchString(pg.Extract) && (strings.Contains(txt, strings.ToLower(name)) || strings.Contains(strings.ToLower(pg.Title), first))
}

// ---------------------------------------------------------------- Doré

var stop = map[string]bool{"the": true, "of": true, "and": true, "a": true, "an": true, "in": true, "to": true, "his": true, "her": true, "from": true, "by": true, "with": true, "at": true, "on": true, "for": true, "is": true, "into": true, "up": true, "out": true, "he": true, "she": true, "it": true, "as": true, "who": true, "their": true, "them": true, "or": true}
var wordRe = regexp.MustCompile(`[a-z]+`)
var numPrefix = regexp.MustCompile(`^\d+\.\s*`)
var extRe = regexp.MustCompile(`(?i)\.(jpe?g|png|gif)$`)

func tokens(s string) map[string]bool {
	t := map[string]bool{}
	for _, w := range wordRe.FindAllString(strings.ToLower(s), -1) {
		if !stop[w] && len(w) > 2 {
			t[strings.TrimSuffix(strings.TrimSuffix(w, "s"), "ing")] = true
		}
	}
	return t
}

func categoryFiles(cat string) []string {
	var out []string
	cont := ""
	for {
		q := url.Values{"action": {"query"}, "list": {"categorymembers"}, "cmtitle": {cat}, "cmlimit": {"500"}, "cmtype": {"file"}, "format": {"json"}}
		if cont != "" {
			q.Set("cmcontinue", cont)
		}
		var r struct {
			Cont  struct{ Cmcontinue string } `json:"continue"`
			Query struct {
				Members []struct{ Title string } `json:"categorymembers"`
			} `json:"query"`
		}
		if get("https://commons.wikimedia.org/w/api.php?"+q.Encode(), &r) != nil {
			break
		}
		for _, m := range r.Query.Members {
			out = append(out, strings.TrimPrefix(m.Title, "File:"))
		}
		if r.Cont.Cmcontinue == "" {
			break
		}
		cont = r.Cont.Cmcontinue
	}
	return out
}

func enrichDore(a *atlas.Atlas, limit int) []Art {
	files := append(categoryFiles("Category:Art depicting the Old Testament by Gustave Doré"), categoryFiles("Category:Art depicting the New Testament by Gustave Doré")...)
	// collapse variants of the same picture, preferring jpg
	best := map[string]string{}
	for _, f := range files {
		if !extRe.MatchString(f) {
			continue
		}
		stem := strings.ToLower(strings.Join(strings.Fields(parenRe.ReplaceAllString(extRe.ReplaceAllString(numPrefix.ReplaceAllString(f, ""), ""), "")), " "))
		stem = strings.TrimSpace(stem)
		if len(stem) < 6 || strings.Contains(stem, "antonio") || strings.Contains(stem, "dore bible") {
			continue
		}
		if cur, ok := best[stem]; !ok || (strings.HasSuffix(strings.ToLower(f), ".jpg") && !strings.HasSuffix(strings.ToLower(cur), ".jpg")) {
			best[stem] = f
		}
	}
	var stems []string
	for s := range best {
		stems = append(stems, s)
	}
	sort.Strings(stems)
	if limit > 0 && len(stems) > limit {
		stems = stems[:limit]
	}
	log.Printf("Doré: %d files, %d distinct pictures", len(files), len(stems))

	type ev struct {
		id   string
		tok  map[string]bool
		name string
	}
	var evs []ev
	for _, e := range a.Events {
		evs = append(evs, ev{e.ID, tokens(e.Title), e.Title})
	}
	var chosen []string
	for _, stem := range stems {
		chosen = append(chosen, best[stem])
	}
	lic := licensedImages(chosen)
	var out []Art
	for _, stem := range stems {
		file := best[stem]
		meta := lic[file]
		if meta == nil {
			continue
		}
		title := strings.Join(strings.Fields(extRe.ReplaceAllString(numPrefix.ReplaceAllString(file, ""), "")), " ")
		art := Art{Image: Image{URL: "https://commons.wikimedia.org/wiki/Special:FilePath/" + url.PathEscape(file) + "?width=960", Page: meta.Page, Title: title, License: meta.License, Artist: meta.Artist, Credit: meta.Credit}}
		tt := tokens(title)
		bestScore, bestID := 0.0, ""
		for _, e := range evs {
			shared := 0
			for w := range tt {
				if e.tok[w] {
					shared++
				}
			}
			if shared < 2 {
				continue
			}
			den := math.Min(float64(len(tt)), float64(len(e.tok)))
			if sc := float64(shared) / den; sc > bestScore {
				bestScore, bestID = sc, e.id
			}
		}
		if bestScore >= 0.6 {
			art.Event, art.Score = bestID, bestScore
		}
		// people named in the picture's title (prominent people only)
		for _, p := range a.People {
			if p.Verses >= 25 && tt[strings.ToLower(strings.TrimSuffix(strings.TrimSuffix(baseName(p.Name), "s"), "ing"))] && len(baseName(p.Name)) > 3 {
				art.People = append(art.People, p.ID)
			}
		}
		out = append(out, art)
	}
	return out
}

// ---------------------------------------------------------------- main

func main() {
	limit := flag.Int("limit", 0, "only process this many places/people/pictures (0 = all prominent ones)")
	nPlaces := flag.Int("places", 450, "how many of the most mentioned places")
	nPeople := flag.Int("people", 500, "how many of the most mentioned people")
	dore := flag.Bool("dore", true, "also fetch Doré engravings")
	out := flag.String("out", "data/enrichment.json", "output file")
	flag.Parse()

	a, err := atlas.Load("data/json", "data/journeys.json")
	if err != nil {
		log.Fatal(err)
	}
	o := Output{Sources: map[string]string{
		"Wikipedia text":    "https://en.wikipedia.org (CC BY-SA 4.0)",
		"Wikimedia Commons": "https://commons.wikimedia.org (per-file free licenses, recorded with each image)",
		"Gustave Doré":      "Doré's Bible illustrations (1866), public domain, via Wikimedia Commons",
	}, Places: map[string]Entry{}, People: map[string]Entry{}}
	if b, err := os.ReadFile(*out); err == nil {
		_ = json.Unmarshal(b, &o)
		if o.Places == nil {
			o.Places = map[string]Entry{}
		}
		if o.People == nil {
			o.People = map[string]Entry{}
		}
	}
	save := func() {
		o.Generated = time.Now().UTC().Format("2006-01-02")
		b, _ := json.MarshalIndent(o, "", " ")
		_ = os.WriteFile(*out, b, 0o644)
	}

	if *dore {
		o.Art = enrichDore(a, *limit)
		save()
	}
	places := append([]*atlas.Place{}, a.Places...)
	sort.Slice(places, func(i, j int) bool { return places[i].Verses > places[j].Verses })
	var pl []*atlas.Place
	for _, p := range places {
		if len(pl) >= *nPlaces || (*limit > 0 && len(pl) >= *limit) {
			break
		}
		if p.HasGeo && p.Verses >= 3 {
			pl = append(pl, p)
		}
	}
	enrichBatch(pl, func(p *atlas.Place) string { return p.ID }, placeCandidates, acceptPlace, o.Places, func(m map[string]*Entry) {
		for id, e := range m {
			o.Places[id] = *e
		}
		save()
	})
	log.Printf("places: %d enriched", len(o.Places))

	people := append([]*atlas.Person{}, a.People...)
	sort.Slice(people, func(i, j int) bool { return people[i].Verses > people[j].Verses })
	var pe []*atlas.Person
	for _, p := range people {
		if len(pe) >= *nPeople || (*limit > 0 && len(pe) >= *limit) {
			break
		}
		if p.Verses >= 3 {
			pe = append(pe, p)
		}
	}
	enrichBatch(pe, func(p *atlas.Person) string { return p.ID }, personCandidates, acceptPerson, o.People, func(m map[string]*Entry) {
		for id, e := range m {
			o.People[id] = *e
		}
		save()
	})
	log.Printf("people: %d enriched", len(o.People))

	withImg := 0
	for _, e := range o.Places {
		if e.Image != nil {
			withImg++
		}
	}
	for _, e := range o.People {
		if e.Image != nil {
			withImg++
		}
	}
	log.Printf("done: %d places, %d people (%d with free images), %d Doré pictures", len(o.Places), len(o.People), withImg, len(o.Art))
}
