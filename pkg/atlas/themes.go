package atlas

import (
	"sort"
	"strings"
)

// Theme is a hand-chosen topic found by matching word stems at the start of words.
// A stem ending in "$" must match a whole word. The lists are editorial and meant to be edited.
type ThemeDef struct {
	ID    string
	Name  string
	Blurb string
	Stems []string
}

var themeDefs = []ThemeDef{
	{"covenant", "Covenant", "Promises made and kept between God and people.", []string{"covenant", "oath", "swore", "sworn"}},
	{"creation", "Creation", "Beginnings: light, earth, sea, living things.", []string{"created", "creation", "creator", "firmament", "formed"}},
	{"exile", "Exile & return", "Captivity, scattering and coming home.", []string{"captiv", "exile", "carried away", "remnant", "return"}},
	{"prayer", "Prayer", "Calling out, asking, thanking.", []string{"pray", "supplicat", "intercess", "cried unto", "call upon"}},
	{"healing", "Healing", "Sickness made well, sight and strength restored.", []string{"heal", "whole$", "cured", "leper", "blind", "lame$"}},
	{"mercy", "Mercy & grace", "Kindness shown to those who do not deserve it.", []string{"mercy", "merciful", "grace", "gracious", "lovingkindness", "compassion"}},
	{"forgiveness", "Forgiveness", "Debts cancelled and wrongs let go.", []string{"forgiv", "forgave", "pardon"}},
	{"love", "Love", "Affection, loyalty and devotion.", []string{"love", "beloved"}},
	{"fear", "Fear & courage", "Fright, and the words that answer it.", []string{"fear not", "be not afraid", "afraid", "be strong", "courage", "dismayed"}},
	{"wisdom", "Wisdom", "Understanding, instruction and the fear of the Lord.", []string{"wisdom", "wise$", "understanding", "instruction", "prudent"}},
	{"justice", "Justice & judgment", "Rightness, verdicts and what is owed.", []string{"judgment", "judgement", "justice", "righteous", "judge$", "judges$"}},
	{"kingdom", "Kings & kingdom", "Thrones, rulers and the kingdom to come.", []string{"king$", "kings$", "kingdom", "throne", "reign", "anointed"}},
	{"temple", "Temple & worship", "Sacred places, altars, offerings and song.", []string{"temple", "tabernacle", "altar", "sacrifice", "offering", "praise", "worship"}},
	{"wilderness", "Wilderness", "Deserts and the places of testing.", []string{"wilderness", "desert", "wandered"}},
	{"shepherd", "Shepherds & flocks", "Sheep, pasture and those who tend them.", []string{"shepherd", "sheep$", "flock", "lamb$", "pasture"}},
	{"water", "Water", "Rivers, rain, wells, sea and thirst.", []string{"water$", "waters$", "river", "well$", "rain", "thirst", "sea$"}},
	{"light", "Light & darkness", "Morning, lamps, stars and the dark.", []string{"light$", "lamp", "darkness", "morning", "shine"}},
	{"bread", "Bread & feasting", "Meals, harvest, wine and the table.", []string{"bread", "feast", "harvest", "wine$", "supper", "banquet"}},
	{"peace", "Peace", "Rest, quiet, safety and shalom.", []string{"peace", "rest$", "quiet", "safe$", "safely", "shalom"}},
	{"hope", "Hope & faith", "Waiting, trusting and believing.", []string{"hope", "faith", "believ", "trust", "wait upon"}},
	{"journey", "Journeys", "Going out, going up and going down.", []string{"journey", "went up", "went down", "departed", "pilgrim"}},
	{"dreams", "Dreams & visions", "Night messages and things seen.", []string{"dream", "vision", "angel", "voice from"}},
	{"war", "War & battle", "Armies, swords and strongholds.", []string{"battle", "war$", "sword", "army", "armies", "captain", "fortress"}},
	{"death", "Death & resurrection", "Endings, graves and rising again.", []string{"death", "died", "grave", "buried", "resurrection", "rose again", "risen"}},
}

type ThemeInfo struct {
	ID     string `json:"id"`
	Name   string `json:"n"`
	Blurb  string `json:"b"`
	Verses int    `json:"vc"`
}

type ThemeDetail struct {
	ThemeInfo
	Books     []int      `json:"books"`
	BookNames []string   `json:"bookNames"`
	Eras      []int      `json:"eras"` // verses per era
	EraNames  []string   `json:"eraNames"`
	Sample    []ThemeHit `json:"sample"`
	Places    []MapPoint `json:"places"`
	PlaceW    []int      `json:"placeW"`
	People    []Ref      `json:"people"`
}

type ThemeHit struct {
	ID    string `json:"id"`
	Label string `json:"l"`
	HTML  string `json:"h"`
}

func matchStem(lower, stem string) bool {
	if strings.HasSuffix(stem, "$") {
		w := strings.TrimSuffix(stem, "$")
		for from := 0; ; {
			i := indexWord(lower, w, from)
			if i < 0 {
				return false
			}
			end := i + len(w)
			if end >= len(lower) || !isLetter(lower[end]) {
				return true
			}
			from = i + 1
		}
	}
	return indexWord(lower, stem, 0) >= 0
}

// verseHits returns the indexes of verses that match any stem, and a rough strength (stems matched).
func (a *Atlas) themeMatches(d ThemeDef) (idx []int, strength []int) {
	for i, t := range a.lowerText {
		n := 0
		for _, s := range d.Stems {
			if matchStem(t, s) {
				n++
			}
		}
		if n > 0 {
			idx = append(idx, i)
			strength = append(strength, n)
		}
	}
	return
}

func (a *Atlas) themeList() {
	a.themeOnce.Do(func() {
		a.themeIdx = map[string][]int{}
		for _, d := range themeDefs {
			idx, _ := a.themeMatches(d)
			a.themeIdx[d.ID] = idx
		}
	})
}

func (a *Atlas) themeInfo(lang string, d ThemeDef) ThemeInfo {
	name, blurb := d.Name, d.Blurb
	if a.tr != nil {
		if s, ok := a.tr.Raw(lang, "theme."+d.ID+".name"); ok {
			name = s
		}
		if s, ok := a.tr.Raw(lang, "theme."+d.ID+".blurb"); ok {
			blurb = s
		}
	}
	return ThemeInfo{d.ID, name, blurb, len(a.themeIdx[d.ID])}
}

// Themes lists the themes with their names in a language.
func (a *Atlas) Themes(lang string) []ThemeInfo {
	a.themeList()
	out := make([]ThemeInfo, 0, len(themeDefs))
	for _, d := range themeDefs {
		out = append(out, a.themeInfo(lang, d))
	}
	return out
}

func (a *Atlas) Theme(id, lang string) *ThemeDetail {
	a.themeList()
	var def *ThemeDef
	for i := range themeDefs {
		if themeDefs[i].ID == id {
			def = &themeDefs[i]
		}
	}
	if def == nil {
		return nil
	}
	idx := a.themeIdx[id]
	d := &ThemeDetail{ThemeInfo: a.themeInfo(lang, *def)}
	bookIdx := map[string]int{}
	for i, b := range a.Books {
		bookIdx[b.ID] = i
		d.BookNames = append(d.BookNames, b.Name)
	}
	d.Books = make([]int, len(a.Books))
	for _, e := range a.LocalEras(lang) {
		d.EraNames = append(d.EraNames, e.Name)
	}
	d.Eras = make([]int, len(Eras))
	placeCount := map[string]int{}
	peopleCount := map[string]int{}
	type cand struct {
		i, score int
	}
	var cands []cand
	for _, i := range idx {
		v := a.Verses[i]
		if bi, ok := bookIdx[v.Book]; ok {
			d.Books[bi]++
		}
		era := eraAt(v.Year)
		for k, e := range Eras {
			if e.ID == era.ID {
				d.Eras[k]++
			}
		}
		if len(v.Places) <= 4 {
			for _, p := range v.Places {
				placeCount[p]++
			}
		}
		if len(v.People) <= 6 {
			for _, p := range v.People {
				peopleCount[p]++
			}
		}
		n := 0
		for _, s := range def.Stems {
			if matchStem(a.lowerText[i], strings.TrimSuffix(s, "$")+suffixOf(s)) {
				n++
			}
		}
		// prefer short verses that match several stems: easier to read as a sample
		cands = append(cands, cand{i, n*100 - len(v.Text)/4})
	}
	sort.Slice(cands, func(x, y int) bool {
		if cands[x].score != cands[y].score {
			return cands[x].score > cands[y].score
		}
		return cands[x].i < cands[y].i
	})
	// a sample spread over the canon rather than all from one book
	seenBook := map[string]int{}
	for _, c := range cands {
		v := a.Verses[c.i]
		if seenBook[v.Book] >= 2 || len(v.Text) < 30 {
			continue
		}
		seenBook[v.Book]++
		d.Sample = append(d.Sample, ThemeHit{v.ID, a.verseLabel(v), highlightStems(v.Text, def.Stems)})
		if len(d.Sample) >= 14 {
			break
		}
	}
	var pids []string
	for k := range placeCount {
		if p := a.PlaceByID[k]; p != nil && p.HasGeo {
			pids = append(pids, k)
		}
	}
	sort.Slice(pids, func(i, j int) bool {
		if placeCount[pids[i]] != placeCount[pids[j]] {
			return placeCount[pids[i]] > placeCount[pids[j]]
		}
		return pids[i] < pids[j]
	})
	if len(pids) > 120 {
		pids = pids[:120]
	}
	for _, id := range pids {
		p := a.PlaceByID[id]
		d.Places = append(d.Places, MapPoint{ID: p.ID, Name: p.Name, Lat: p.Lat, Lng: p.Lng, Kind: p.Kind})
		d.PlaceW = append(d.PlaceW, placeCount[id])
	}
	var peids []string
	for k := range peopleCount {
		peids = append(peids, k)
	}
	sort.Slice(peids, func(i, j int) bool {
		if peopleCount[peids[i]] != peopleCount[peids[j]] {
			return peopleCount[peids[i]] > peopleCount[peids[j]]
		}
		return peids[i] < peids[j]
	})
	d.People = a.personRefs(peids[:min(len(peids), 14)])
	return d
}

func suffixOf(s string) string {
	if strings.HasSuffix(s, "$") {
		return "$"
	}
	return ""
}

func highlightStems(text string, stems []string) string {
	var terms []string
	for _, s := range stems {
		terms = append(terms, strings.TrimSuffix(s, "$"))
	}
	return highlight(text, terms)
}
