package atlas

import (
	"encoding/json"
	"net/url"
	"regexp"
	"sort"
	"strings"
)

// GalleryLink is a link shown beside an image in the lightbox.
type GalleryLink struct {
	Label string `json:"l"`
	URL   string `json:"u"`
}

// GalleryItem is one picture ready for a mosaic tile and the lightbox.
type GalleryItem struct {
	Thumb  string
	Full   string
	Title  string
	Page   string // the Commons file page
	Credit string
	Kind   string // "engraving", "person" or "place"
	Links  []GalleryLink
}

func (g GalleryItem) LinksJSON() string {
	if g.Links == nil {
		return "[]"
	}
	b, _ := json.Marshal(g.Links)
	return string(b)
}

var thumbPx = regexp.MustCompile(`/\d+px-`)

// sized returns the image URL at about width w (Commons thumbnails come in fixed steps).
func sized(u string, w int) string {
	switch {
	case strings.Contains(u, "Special:FilePath"):
		p, err := url.Parse(u)
		if err != nil {
			return u
		}
		q := p.Query()
		q.Set("width", itoa(w))
		p.RawQuery = q.Encode()
		return p.String()
	case strings.Contains(u, "upload.wikimedia.org") && thumbPx.MatchString(u):
		return thumbPx.ReplaceAllString(u, "/"+itoa(w)+"px-")
	}
	return u
}

func itoa(n int) string {
	b, _ := json.Marshal(n)
	return string(b)
}

func creditOf(i Image, doré bool) string {
	artist := i.Artist
	if doré {
		artist = "Gustave Doré"
	}
	if len([]rune(artist)) > 60 {
		artist = string([]rune(artist)[:60]) + "…"
	}
	switch {
	case artist != "" && i.License != "":
		return artist + " · " + i.License
	case i.License != "":
		return i.License
	}
	return artist
}

var (
	artIDRe     = regexp.MustCompile(`\s*\(\d{5,}\)`)
	artPrefixRe = regexp.MustCompile(`(?i)^(by\s+)?gustave\s+dor[eé]\s*(\([^)]*\))?\s*[,;:\-–—]*\s*`)
	artSuffixRe = regexp.MustCompile(`(?i)\s*[-–—]\s*google art project$|\s*(retouched|retusche|cropped)\s*$`)
)

// cleanArtTitle turns a Commons file name into something readable.
func cleanArtTitle(t string) string {
	t = strings.ReplaceAll(t, "_", " ")
	t = artIDRe.ReplaceAllString(t, "")
	t = artPrefixRe.ReplaceAllString(t, "")
	t = artSuffixRe.ReplaceAllString(t, "")
	t = strings.TrimSpace(strings.Trim(t, "-–—,; "))
	if t == "" {
		return "Engraving"
	}
	return strings.ToUpper(t[:1]) + t[1:]
}

func (a *Atlas) artItem(art Art) GalleryItem {
	it := GalleryItem{Thumb: sized(art.URL, 480), Full: sized(art.URL, 1280), Title: cleanArtTitle(art.Title), Page: art.Page, Credit: creditOf(art.Image, true), Kind: "engraving"}
	if art.Event != "" {
		if e := a.EventByID[art.Event]; e != nil {
			it.Links = append(it.Links, GalleryLink{"Event: " + e.Title, "/event/" + e.ID})
		}
	}
	for i, id := range art.People {
		if i >= 4 {
			break
		}
		if p := a.PersonByID[id]; p != nil {
			it.Links = append(it.Links, GalleryLink{p.Name, "/person/" + p.ID})
		}
	}
	return it
}

// ArtGallery turns engravings into gallery items.
func (a *Atlas) ArtGallery(arts []Art) []GalleryItem {
	out := make([]GalleryItem, 0, len(arts))
	for _, art := range arts {
		out = append(out, a.artItem(art))
	}
	return out
}

func (a *Atlas) wikiItem(kind, id, name string, w WikiEntry) *GalleryItem {
	if w.Image == nil {
		return nil
	}
	return &GalleryItem{Thumb: w.Image.URL, Full: sized(w.Image.URL, 1280), Title: name, Page: w.Image.Page, Credit: creditOf(*w.Image, false), Kind: kind,
		Links: []GalleryLink{{name, "/" + kind + "/" + id}, {"Wikipedia: " + w.Title, w.Page}}}
}

// WikiPhotoItem is a single Wikipedia photo as a lightbox item.
func (a *Atlas) WikiPhotoItem(kind, id, name string, w *WikiEntry) *GalleryItem {
	if w == nil {
		return nil
	}
	return a.wikiItem(kind, id, name, *w)
}

// AllGallery collects every image on the site: engravings, then photos of people and places.
func (a *Atlas) AllGallery() []GalleryItem {
	if a.Enrich == nil {
		return nil
	}
	var eng, people, places []GalleryItem
	for _, art := range a.Enrich.Art {
		eng = append(eng, a.artItem(art))
	}
	for id, w := range a.Enrich.People {
		if p := a.PersonByID[id]; p != nil {
			if it := a.wikiItem("person", id, p.Name, w); it != nil {
				people = append(people, *it)
			}
		}
	}
	for id, w := range a.Enrich.Places {
		if p := a.PlaceByID[id]; p != nil {
			if it := a.wikiItem("place", id, p.Name, w); it != nil {
				places = append(places, *it)
			}
		}
	}
	byTitle := func(s []GalleryItem) {
		sort.Slice(s, func(i, j int) bool { return s[i].Title < s[j].Title })
	}
	byTitle(eng)
	byTitle(people)
	byTitle(places)
	// engravings that are linked to something come first
	sort.SliceStable(eng, func(i, j int) bool { return len(eng[i].Links) > len(eng[j].Links) })
	return append(append(append([]GalleryItem{}, eng...), people...), places...)
}
