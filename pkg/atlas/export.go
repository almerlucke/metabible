package atlas

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"

	"github.com/almerlucke/metabible/pkg/i18n"
)

func xmlEsc(s string) string {
	var b bytes.Buffer
	_ = xml.EscapeText(&b, []byte(s))
	return b.String()
}

func (a *Atlas) journeyByKey(key, lang string) *Journey {
	if strings.HasPrefix(key, "person-") {
		return a.PersonJourney(strings.TrimPrefix(key, "person-"), lang)
	}
	for _, j := range a.JourneysFor(lang) {
		if j.ID == key {
			return j
		}
	}
	return nil
}

func gpx(j *Journey) string {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	b.WriteString(`<gpx version="1.1" creator="MetaBible" xmlns="http://www.topografix.com/GPX/1/1">` + "\n")
	fmt.Fprintf(&b, "<metadata><name>%s</name><desc>%s</desc></metadata>\n", xmlEsc(j.Title), xmlEsc(j.Intro))
	for i, s := range j.Stops {
		fmt.Fprintf(&b, `<wpt lat="%.6f" lon="%.6f"><name>%d. %s</name><desc>%s</desc></wpt>`+"\n", s.Lat, s.Lng, i+1, xmlEsc(s.Name), xmlEsc(strings.TrimSpace(s.Label+" "+s.Note)))
	}
	fmt.Fprintf(&b, "<rte><name>%s</name>\n", xmlEsc(j.Title))
	for i, s := range j.Stops {
		fmt.Fprintf(&b, `<rtept lat="%.6f" lon="%.6f"><name>%d. %s</name></rtept>`+"\n", s.Lat, s.Lng, i+1, xmlEsc(s.Name))
	}
	b.WriteString("</rte>\n</gpx>\n")
	return b.String()
}

func kml(j *Journey) string {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	b.WriteString(`<kml xmlns="http://www.opengis.net/kml/2.2"><Document>` + "\n")
	fmt.Fprintf(&b, "<name>%s</name><description>%s</description>\n", xmlEsc(j.Title), xmlEsc(j.Intro))
	for i, s := range j.Stops {
		desc := strings.TrimSpace(s.Label + " " + s.Note)
		if s.Text != "" {
			desc += "\n“" + s.Text + "”"
		}
		fmt.Fprintf(&b, "<Placemark><name>%d. %s</name><description>%s</description><Point><coordinates>%.6f,%.6f,0</coordinates></Point></Placemark>\n", i+1, xmlEsc(s.Name), xmlEsc(desc), s.Lng, s.Lat)
	}
	b.WriteString("<Placemark><name>" + xmlEsc(j.Title) + " (route)</name><LineString><tessellate>1</tessellate><coordinates>")
	for _, s := range j.Stops {
		fmt.Fprintf(&b, "%.6f,%.6f,0 ", s.Lng, s.Lat)
	}
	b.WriteString("</coordinates></LineString></Placemark>\n</Document></kml>\n")
	return b.String()
}

// handleExport serves a journey as GPX or KML: /export/journey/{key}.gpx
func (a *Atlas) handleExport(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("file")
	var key, ext string
	switch {
	case strings.HasSuffix(name, ".gpx"):
		key, ext = strings.TrimSuffix(name, ".gpx"), "gpx"
	case strings.HasSuffix(name, ".kml"):
		key, ext = strings.TrimSuffix(name, ".kml"), "kml"
	default:
		http.NotFound(w, r)
		return
	}
	j := a.journeyByKey(key, i18n.FromContext(r.Context()))
	if j == nil || len(j.Stops) == 0 {
		http.NotFound(w, r)
		return
	}
	body, ctype := gpx(j), "application/gpx+xml"
	if ext == "kml" {
		body, ctype = kml(j), "application/vnd.google-earth.kml+xml"
	}
	w.Header().Set("Content-Type", ctype+"; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.%s"`, strings.ReplaceAll(key, `"`, ""), ext))
	_, _ = w.Write([]byte(body))
}

// ---------------------------------------------------------------- reading progress

type ProgressChapter struct {
	ID     string `json:"id"`
	Num    int    `json:"n"`
	Places []int  `json:"pl,omitempty"`
	People int    `json:"pe"`
	Y0     int    `json:"y0"`
	Y1     int    `json:"y1"`
}

type ProgressBook struct {
	ID        string            `json:"id"`
	Name      string            `json:"n"`
	Testament string            `json:"t"`
	Chapters  []ProgressChapter `json:"ch"`
}

type ProgressIndex struct {
	Books  []ProgressBook `json:"books"`
	Places []MapPoint     `json:"places"`
}

// ProgressIndex lets the browser work out, from the chapters a reader has finished,
// which places they have met and which centuries they have covered.
func (a *Atlas) ProgressIndex(lang string) *ProgressIndex {
	a.progMu.Lock()
	defer a.progMu.Unlock()
	if p, ok := a.prog[lang]; ok {
		return p
	}
	if a.prog == nil {
		a.prog = map[string]*ProgressIndex{}
	}
	idx := a.buildProgressIndex(lang)
	a.prog[lang] = idx
	return idx
}

func (a *Atlas) buildProgressIndex(lang string) *ProgressIndex {
	idx := &ProgressIndex{}
	pidx := map[string]int{}
	for _, b := range a.Books {
		pb := ProgressBook{ID: b.ID, Name: a.BookName(b, lang), Testament: b.Testament}
		for _, c := range b.Chapters {
			pc := ProgressChapter{ID: c.ID, Num: c.Num}
			seen := map[string]bool{}
			vs := a.chapterVerses[c.ID]
			people := map[string]bool{}
			for i, v := range vs {
				if i == 0 {
					pc.Y0 = v.Year
				}
				pc.Y1 = v.Year
				for _, p := range v.People {
					people[p] = true
				}
				for _, pl := range v.Places {
					if seen[pl] {
						continue
					}
					seen[pl] = true
					p := a.PlaceByID[pl]
					if p == nil || !p.HasGeo {
						continue
					}
					k, ok := pidx[pl]
					if !ok {
						k = len(idx.Places)
						pidx[pl] = k
						idx.Places = append(idx.Places, MapPoint{ID: p.ID, Name: p.Name, Lat: p.Lat, Lng: p.Lng, Kind: p.Kind})
					}
					pc.Places = append(pc.Places, k)
				}
			}
			pc.People = len(people)
			pb.Chapters = append(pb.Chapters, pc)
		}
		idx.Books = append(idx.Books, pb)
	}
	return idx
}
