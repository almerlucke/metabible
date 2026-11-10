package atlas

import (
	"fmt"
	"sort"
	"strings"
)

// Issue is something suspicious in the data. Errors are hard inconsistencies
// (dangling references, impossible dates); warnings are things worth a look.
type Issue struct {
	Severity string `json:"severity"` // "error" or "warning"
	Kind     string `json:"kind"`
	Subject  string `json:"subject"`
	ID       string `json:"id,omitempty"`
	Message  string `json:"message"`
}

func isLifespanEvent(title string) bool {
	return strings.HasPrefix(title, "Death of") || strings.HasPrefix(title, "Birth of") || strings.HasPrefix(title, "Lifetime of")
}

// Validate cross-checks the loaded data and returns what it finds, errors first.
func (a *Atlas) Validate() []Issue {
	var out []Issue
	add := func(sev, kind, subject, id, format string, args ...any) {
		out = append(out, Issue{sev, kind, subject, id, fmt.Sprintf(format, args...)})
	}

	// events against their predecessors
	for _, e := range a.Events {
		for _, pid := range e.Pred {
			p := a.EventByID[pid]
			if p == nil {
				add("error", "event-ref", e.Title, e.ID, "predecessor %s does not exist", pid)
				continue
			}
			d := e.Year - p.Year
			switch {
			case d < -10:
				add("error", "event-order", e.Title, e.ID, "starts %.0f years before its predecessor %q (%s vs %s)", -d, p.Title, e.Date, p.Date)
			case d < -3:
				add("warning", "event-order", e.Title, e.ID, "starts %.0f years before its predecessor %q (%s vs %s)", -d, p.Title, e.Date, p.Date)
			case d > 400 && !isLifespanEvent(e.Title):
				add("warning", "event-gap", e.Title, e.ID, "starts %.0f years after its predecessor %q", d, p.Title)
			}
		}
		for _, pl := range e.Places {
			if a.PlaceByID[pl] == nil {
				add("error", "event-ref", e.Title, e.ID, "place %s does not exist", pl)
			}
		}
		for _, pe := range e.People {
			if a.PersonByID[pe] == nil {
				add("warning", "event-ref", e.Title, e.ID, "participant %s does not exist", pe)
			}
		}
		if e.Date == "" {
			add("error", "event-date", e.Title, e.ID, "no usable start date")
		}
	}

	// places
	for _, p := range a.Places {
		if !p.HasGeo {
			continue
		}
		if p.Lat < -90 || p.Lat > 90 || p.Lng < -180 || p.Lng > 180 {
			add("error", "place-coords", p.Name, p.ID, "coordinates out of range (%v, %v)", p.Lat, p.Lng)
		} else if p.Lat == 0 && p.Lng == 0 {
			add("error", "place-coords", p.Name, p.ID, "coordinates are 0,0")
		} else if p.Lng < -15 || p.Lng > 100 || p.Lat < 0 || p.Lat > 60 {
			add("warning", "place-far", p.Name, p.ID, "far from the biblical world (%.1f, %.1f)", p.Lat, p.Lng)
		}
	}

	// people
	for _, p := range a.People {
		if p.Born != nil && p.Died != nil && *p.Born > *p.Died {
			add("error", "person-dates", p.Name, p.ID, "born %d, died %d", *p.Born, *p.Died)
		}
		for _, cid := range p.Children {
			c := a.PersonByID[cid]
			if c == nil {
				add("warning", "person-ref", p.Name, p.ID, "child %s does not exist", cid)
				continue
			}
			if p.Born != nil && c.Born != nil && *c.Born < *p.Born {
				add("error", "person-order", c.Name, c.ID, "born %d, before their parent %s (%d)", *c.Born, p.Name, *p.Born)
			}
		}
		if p.Father != "" && a.PersonByID[p.Father] == nil {
			add("warning", "person-ref", p.Name, p.ID, "father %s does not exist", p.Father)
		}
	}

	// verses (a missing event is reported once, with how many verses point at it)
	missingEvents := map[string]int{}
	for _, v := range a.Verses {
		for _, id := range v.Places {
			if a.PlaceByID[id] == nil {
				add("error", "verse-ref", v.Osis, v.ID, "place %s does not exist", id)
			}
		}
		for _, id := range v.People {
			if a.PersonByID[id] == nil {
				add("error", "verse-ref", v.Osis, v.ID, "person %s does not exist", id)
			}
		}
		for _, id := range v.Events {
			if a.EventByID[id] == nil {
				missingEvents[id]++
			}
		}
	}
	for id, n := range missingEvents {
		add("warning", "verse-event", id, id, "event is referenced by %d verses but is not in the events table", n)
	}

	// curated journeys: each stop should sit close to the place it links to
	for _, j := range a.Journeys {
		for _, s := range j.Stops {
			if s.Place == "" {
				if !s.Sea {
					add("warning", "journey-stop", j.Title+": "+s.Name, j.ID, "stop is not linked to a known place")
				}
				continue
			}
			if p := a.PlaceByID[s.Place]; p != nil {
				if km := haversine(s.Lat, s.Lng, p.Lat, p.Lng); km > 40 && strings.EqualFold(s.Name, p.Name) {
					add("warning", "journey-stop", j.Title+": "+s.Name, j.ID, "stop is %.0f km from the place it links to", km)
				}
			}
		}
	}

	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Severity != out[j].Severity {
			return out[i].Severity == "error"
		}
		if out[i].Kind != out[j].Kind {
			return out[i].Kind < out[j].Kind
		}
		return out[i].Subject < out[j].Subject
	})
	return out
}
