package atlas

import "strings"

var eraIntros = map[string]string{
	"primeval":   "Before there were nations there was a garden. Creation, a fall, a flood, and a tower: the first chapters of the story are told in a world with almost no map, only genealogies stretching across two thousand years.",
	"patriarchs": "A man leaves Mesopotamia for a land he has never seen, and a family begins to grow into a people. Abraham, Isaac and Jacob wander between Haran, Canaan and Egypt, and the story acquires its first real places.",
	"egypt":      "Joseph's family goes down to Egypt in famine and stays for centuries. Then a deliverer arrives, a sea divides, and a nation walks through the wilderness to a mountain where it receives its law.",
	"conquest":   "Across the Jordan the land has to be taken and then lived in. Joshua's campaigns give way to the Judges: cycles of failure, rescue and rescue again, until the people ask for a king.",
	"kingdom":    "Saul, David and Solomon: three kings in about a century. Jerusalem becomes a capital, a temple rises on Mount Moriah, and the kingdom reaches its greatest extent.",
	"divided":    "The kingdom splits in two. Prophets such as Elijah, Isaiah and Jeremiah speak into a long decline, as Assyria and then Babylon grow stronger, until both kingdoms fall.",
	"exile":      "Jerusalem burns and its people are carried to Babylon. Then Persia changes the map, and groups of exiles return to rebuild the city, the temple and the wall.",
	"silence":    "Between the last prophet and the first gospel there are some four hundred years that the narrative does not tell. Empires change, and the world is being prepared for what comes next.",
	"gospels":    "In a small province of the Roman Empire a child is born in Bethlehem. Around the Sea of Galilee, and then on a last road to Jerusalem, a ministry unfolds that the rest of the Bible looks back on.",
	"church":     "From a room in Jerusalem the message crosses seas and borders. Peter, Paul and many others carry it to Antioch, Asia Minor, Greece and Rome, and John closes the story on the island of Patmos.",
}

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
func (a *Atlas) Story() *Story {
	st := &Story{}
	cur := ""
	for _, e := range a.Events {
		era := eraAt(int(e.Year))
		if era.ID != cur {
			cur = era.ID
			st.Eras = append(st.Eras, StoryEra{Era: era, Intro: eraIntros[era.ID], Start: len(st.Steps)})
			st.Steps = append(st.Steps, StoryStep{Kind: "era", Era: era.ID, Title: era.Name, Year: float64(era.From), Date: FormatYear(era.From) + " – " + FormatYear(era.To), Text: eraIntros[era.ID]})
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
