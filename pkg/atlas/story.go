package atlas

import "strings"

type StoryPoint struct {
	ID   string  `json:"id"`
	Name string  `json:"n"`
	Lat  float64 `json:"lat"`
	Lng  float64 `json:"lng"`
	Kind string  `json:"k,omitempty"`
}

type StoryPerson struct {
	ID   string `json:"id"`
	Name string `json:"n"`
}

type StoryStep struct {
	Kind     string        `json:"kind"` // "era" or "event"
	Era      string        `json:"era"`
	ID       string        `json:"id,omitempty"`
	Title    string        `json:"t"`
	Date     string        `json:"dt,omitempty"`
	Year     float64       `json:"y"`
	Text     string        `json:"tx,omitempty"`
	Ref      string        `json:"ref,omitempty"`
	VerseID  string        `json:"vid,omitempty"`
	Places   []StoryPoint  `json:"pl,omitempty"`
	People   []StoryPerson `json:"pe,omitempty"`
	Blurb    string        `json:"bl,omitempty"`
	BlurbOf  string        `json:"blo,omitempty"`
	BlurbURL string        `json:"blu,omitempty"`
	Pred     []string      `json:"pr,omitempty"`
}

type StoryEra struct {
	Era
	Intro string `json:"intro"`
	Start int    `json:"start"`
}

type Story struct {
	Eras  []StoryEra  `json:"eras"`
	Steps []StoryStep `json:"steps"`
}

func firstSentences(s string, max int) string {
	s = plain(s)
	if len(s) <= max {
		return s
	}
	cut := s[:max]
	if i := strings.LastIndexAny(cut, ".;"); i > max/2 {
		return cut[:i+1]
	}
	return trunc(s, max)
}

// Story builds the chronological walk: an introduction for each era followed by its events.
func (a *Atlas) Story(lang string) *Story {
	st := &Story{}
	eras := a.LocalEras(lang)
	cur := ""
	for _, e := range a.Events {
		era := eraAt(int(e.Year))
		if era.ID != cur {
			cur = era.ID
			var local Era
			for _, le := range eras {
				if le.ID == era.ID {
					local = le
				}
			}
			intro := a.tx(lang, "era."+era.ID+".intro")
			st.Eras = append(st.Eras, StoryEra{Era: local, Intro: intro, Start: len(st.Steps)})
			st.Steps = append(st.Steps, StoryStep{Kind: "era", Era: era.ID, Title: local.Name, Year: float64(era.From), Date: FormatYear(era.From) + " – " + FormatYear(era.To), Text: intro})
		}
		s := StoryStep{Kind: "event", Era: era.ID, ID: e.ID, Title: e.Title, Date: e.Date, Year: e.Year, Ref: e.Ref, VerseID: e.VerseID, Pred: e.Pred}

		// the account: up to three of the event's verses, in order
		var vs []*Verse
		for _, vid := range e.Vs {
			if v := a.VerseByID[vid]; v != nil {
				vs = append(vs, v)
			}
		}
		for i := 1; i < len(vs); i++ {
			for j := i; j > 0 && vs[j].VID < vs[j-1].VID; j-- {
				vs[j], vs[j-1] = vs[j-1], vs[j]
			}
		}
		var parts []string
		for i, v := range vs {
			if i >= 2 {
				break
			}
			parts = append(parts, v.Text)
		}
		s.Text = trunc(strings.Join(parts, " "), 420)
		if len(vs) > 0 {
			s.Ref, s.VerseID = a.verseLabel(vs[0]), vs[0].ID
		}

		for i, id := range e.Places {
			if i >= 4 {
				break
			}
			if p := a.PlaceByID[id]; p != nil && p.HasGeo {
				s.Places = append(s.Places, StoryPoint{ID: p.ID, Name: p.Name, Lat: p.Lat, Lng: p.Lng, Kind: p.Kind})
			}
		}
		for _, id := range e.People {
			if p := a.PersonByID[id]; p != nil && len(s.People) < 5 {
				s.People = append(s.People, StoryPerson{ID: p.ID, Name: p.Name})
			}
		}
		// a short "know more": the first participant (or place) with a dictionary entry
		for _, id := range e.People {
			if p := a.PersonByID[id]; p != nil && p.Dict != "" && p.Name != "God" && p.Name != "Holy Spirit" {
				s.Blurb, s.BlurbOf, s.BlurbURL = firstSentences(p.Dict, 260), p.Name, "/person/"+p.ID
				break
			}
		}
		if s.Blurb == "" {
			for _, id := range e.Places {
				if p := a.PlaceByID[id]; p != nil && p.Dict != "" {
					s.Blurb, s.BlurbOf, s.BlurbURL = firstSentences(p.Dict, 260), p.Name, "/place/"+p.ID
					break
				}
			}
		}
		st.Steps = append(st.Steps, s)
	}
	return st
}
