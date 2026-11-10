package atlas

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Image is an open-licensed picture with the attribution its license needs.
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

// WikiEntry is a short Wikipedia summary (CC BY-SA 4.0) with an optional free image.
type WikiEntry struct {
	Title   string `json:"title"`
	Summary string `json:"summary"`
	Desc    string `json:"desc,omitempty"`
	Page    string `json:"page"`
	Image   *Image `json:"image,omitempty"`
}

// Art is a public-domain engraving (Gustave Doré) matched to an event or to people.
type Art struct {
	Image
	Event  string   `json:"event,omitempty"`
	People []string `json:"people,omitempty"`
	Score  float64  `json:"score"`
}

// Enrichment is the optional data/enrichment.json produced by cmd/enrich.
type Enrichment struct {
	Generated string               `json:"generated"`
	Sources   map[string]string    `json:"sources"`
	Places    map[string]WikiEntry `json:"places"`
	People    map[string]WikiEntry `json:"people"`
	Art       []Art                `json:"art"`
}

// dropDuplicates removes pictures that data/art-duplicates.json (made by cmd/arthash) lists as copies of another.
func dropDuplicates(art []Art, dupFile string) ([]Art, error) {
	b, err := os.ReadFile(dupFile)
	if errors.Is(err, os.ErrNotExist) {
		return art, nil
	}
	if err != nil {
		return nil, err
	}
	var d struct {
		Duplicates []struct {
			Title string `json:"title"`
		} `json:"duplicates"`
	}
	if err := json.Unmarshal(b, &d); err != nil {
		return nil, err
	}
	skip := map[string]bool{}
	for _, x := range d.Duplicates {
		skip[x.Title] = true
	}
	var out []Art
	for _, x := range art {
		if !skip[x.Title] {
			out = append(out, x)
		}
	}
	return out, nil
}

func (a *Atlas) loadEnrichment(file string) error {
	b, err := os.ReadFile(file)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	var en Enrichment
	if err := json.Unmarshal(b, &en); err != nil {
		return err
	}
	var err2 error
	if en.Art, err2 = dropDuplicates(en.Art, filepath.Join(filepath.Dir(file), "art-duplicates.json")); err2 != nil {
		return err2
	}
	var kept []Art
	for _, art := range en.Art {
		drop := false
		for _, f := range a.artFixes {
			if f.Duplicate && f.Title != "" && strings.Contains(strings.ToLower(art.Title), strings.ToLower(f.Title)) {
				drop = true
			}
		}
		if !drop {
			kept = append(kept, art)
		}
	}
	en.Art = kept
	for i := range en.Art {
		for _, f := range a.artFixes {
			if !f.Duplicate && f.Title != "" && strings.Contains(strings.ToLower(en.Art[i].Title), strings.ToLower(f.Title)) {
				if f.Event != "" && a.EventByID[f.Event] == nil {
					return fmt.Errorf("corrections: art fix points at unknown event %s", f.Event)
				}
				en.Art[i].Event = f.Event
			}
		}
	}
	// link engravings to people by the distinctive names in their titles, besides what the fetch stored
	keys := a.artNameKeys()
	for i := range en.Art {
		have := map[string]bool{}
		for _, id := range en.Art[i].People {
			have[id] = true
		}
		for _, id := range peopleInTitle(en.Art[i].Title, keys) {
			if !have[id] {
				en.Art[i].People = append(en.Art[i].People, id)
			}
		}
	}
	a.Enrich = &en
	a.artByEvent = map[string][]Art{}
	a.artByPerson = map[string][]Art{}
	for _, art := range en.Art {
		if art.Event != "" {
			a.artByEvent[art.Event] = append(a.artByEvent[art.Event], art)
		}
		for _, p := range art.People {
			a.artByPerson[p] = append(a.artByPerson[p], art)
		}
	}
	return nil
}

func (a *Atlas) wikiPerson(id string) *WikiEntry {
	if a.Enrich == nil {
		return nil
	}
	if w, ok := a.Enrich.People[id]; ok {
		return &w
	}
	return nil
}

func (a *Atlas) wikiPlace(id string) *WikiEntry {
	if a.Enrich == nil {
		return nil
	}
	if w, ok := a.Enrich.Places[id]; ok {
		return &w
	}
	return nil
}

// Credit is one line on the credits page.
type Credit struct {
	Title   string
	Page    string
	Artist  string
	License string
	Kind    string
}

// Credits lists every image and text source in use, for the credits page.
func (a *Atlas) Credits() (images []Credit, texts []Credit) {
	if a.Enrich == nil {
		return nil, nil
	}
	for id, w := range a.Enrich.Places {
		if p := a.PlaceByID[id]; p != nil {
			texts = append(texts, Credit{Title: p.Name + " — " + w.Title, Page: w.Page, License: "CC BY-SA 4.0", Kind: "Wikipedia"})
		}
		if w.Image != nil {
			images = append(images, Credit{Title: w.Title, Page: w.Image.Page, Artist: w.Image.Artist, License: w.Image.License, Kind: "place"})
		}
	}
	for id, w := range a.Enrich.People {
		if p := a.PersonByID[id]; p != nil {
			texts = append(texts, Credit{Title: p.Name + " — " + w.Title, Page: w.Page, License: "CC BY-SA 4.0", Kind: "Wikipedia"})
		}
		if w.Image != nil {
			images = append(images, Credit{Title: w.Title, Page: w.Image.Page, Artist: w.Image.Artist, License: w.Image.License, Kind: "person"})
		}
	}
	for _, art := range a.Enrich.Art {
		images = append(images, Credit{Title: art.Title, Page: art.Page, Artist: art.Artist, License: art.License, Kind: "engraving"})
	}
	sort.Slice(images, func(i, j int) bool { return images[i].Title < images[j].Title })
	sort.Slice(texts, func(i, j int) bool { return texts[i].Title < texts[j].Title })
	return images, texts
}
