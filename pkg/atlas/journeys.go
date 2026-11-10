package atlas

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strings"
)

type Stop struct {
	Place   string  `json:"pl,omitempty"`
	Name    string  `json:"n"`
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
	Year    int     `json:"y,omitempty"`
	Ref     string  `json:"ref,omitempty"`
	Label   string  `json:"lab,omitempty"`
	VerseID string  `json:"vid,omitempty"`
	Text    string  `json:"tx,omitempty"`
	Note    string  `json:"note,omitempty"`
	Kind    string  `json:"k,omitempty"`
	Sea     bool    `json:"sea,omitempty"`
}

type Journey struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Subtitle string `json:"sub,omitempty"`
	Intro    string `json:"intro,omitempty"`
	Color    string `json:"color"`
	Person   string `json:"person,omitempty"`
	Era      string `json:"era,omitempty"`
	Stops    []Stop `json:"stops"`
	Km       int    `json:"km"`
}

type rawStop struct {
	Name string  `json:"name"`
	Lat  float64 `json:"lat"`
	Lng  float64 `json:"lng"`
	Ref  string  `json:"ref"`
	Note string  `json:"note"`
	Year int     `json:"year"`
	Sea  bool    `json:"sea"`
}

type rawJourney struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Subtitle string    `json:"subtitle"`
	Intro    string    `json:"intro"`
	Color    string    `json:"color"`
	Person   string    `json:"person"`
	Era      string    `json:"era"`
	Stops    []rawStop `json:"stops"`
}

func haversine(lat1, lng1, lat2, lng2 float64) float64 {
	const r = 6371.0
	p1, p2 := lat1*math.Pi/180, lat2*math.Pi/180
	dp := (lat2 - lat1) * math.Pi / 180
	dl := (lng2 - lng1) * math.Pi / 180
	h := math.Sin(dp/2)*math.Sin(dp/2) + math.Cos(p1)*math.Cos(p2)*math.Sin(dl/2)*math.Sin(dl/2)
	return 2 * r * math.Asin(math.Sqrt(h))
}

func (a *Atlas) nearestPlace(name string, lat, lng float64) *Place {
	var best *Place
	bd := 120.0 // km
	ln := strings.ToLower(name)
	for _, p := range a.Places {
		if !p.HasGeo {
			continue
		}
		pn := strings.ToLower(p.Name)
		if pn != ln && !strings.HasPrefix(pn, ln) && !strings.HasPrefix(ln, pn) {
			continue
		}
		d := haversine(lat, lng, p.Lat, p.Lng)
		if d < bd {
			bd, best = d, p
		}
	}
	return best
}

func loadJourneys(a *Atlas, file string) ([]*Journey, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var raws []rawJourney
	if err := json.Unmarshal(b, &raws); err != nil {
		return nil, fmt.Errorf("journeys: %w", err)
	}
	var out []*Journey
	for _, r := range raws {
		j := &Journey{ID: r.ID, Title: r.Title, Subtitle: r.Subtitle, Intro: r.Intro, Color: r.Color, Person: r.Person, Era: r.Era}
		for _, s := range r.Stops {
			st := Stop{Name: s.Name, Lat: s.Lat, Lng: s.Lng, Note: s.Note, Year: s.Year, Ref: s.Ref, Sea: s.Sea}
			if p := a.nearestPlace(s.Name, s.Lat, s.Lng); p != nil {
				st.Place = p.ID
				st.Kind = p.Kind
			}
			if s.Ref != "" {
				if v, ok := a.VerseByOsis[strings.ToLower(s.Ref)]; ok {
					st.VerseID = v.ID
					st.Text = v.Text
					st.Label = a.verseLabel(v)
					if st.Year == 0 {
						st.Year = v.Year
					}
				} else {
					return nil, fmt.Errorf("journey %s: unknown verse ref %q", r.ID, s.Ref)
				}
			}
			j.Stops = append(j.Stops, st)
		}
		for i := 1; i < len(j.Stops); i++ {
			j.Km += int(haversine(j.Stops[i-1].Lat, j.Stops[i-1].Lng, j.Stops[i].Lat, j.Stops[i].Lng))
		}
		out = append(out, j)
	}
	return out, nil
}

// PersonJourney derives an itinerary from the places mentioned in the verses a
// person appears in, in canonical order. It is an approximation, but covers
// everyone who is mentioned alongside places.
func (a *Atlas) PersonJourney(id string) *Journey {
	p := a.PersonByID[id]
	if p == nil {
		return nil
	}
	j := &Journey{ID: "person-" + id, Title: p.Name, Subtitle: "Places mentioned alongside " + p.Name + ", in order of the text", Color: "#f5c451", Person: id}
	last := ""
	for _, v := range a.Verses {
		if len(v.Places) == 0 || len(v.Places) > 4 {
			continue
		}
		in := false
		for _, pid := range v.People {
			if pid == id {
				in = true
				break
			}
		}
		if !in {
			continue
		}
		for _, plid := range v.Places {
			pl := a.PlaceByID[plid]
			if pl == nil || !pl.HasGeo || plid == last {
				continue
			}
			last = plid
			j.Stops = append(j.Stops, Stop{Place: plid, Name: pl.Name, Lat: pl.Lat, Lng: pl.Lng, Year: v.Year,
				Ref: v.Osis, Label: a.verseLabel(v), VerseID: v.ID, Text: v.Text, Kind: pl.Kind})
		}
	}
	if len(j.Stops) > 250 {
		j.Stops = j.Stops[:250]
	}
	for i := 1; i < len(j.Stops); i++ {
		j.Km += int(haversine(j.Stops[i-1].Lat, j.Stops[i-1].Lng, j.Stops[i].Lat, j.Stops[i].Lng))
	}
	return j
}
