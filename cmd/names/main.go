// Command names builds data/names/<lang>.json: the spelling of people and places in a
// Bible translation, found by aligning each name with the words of the verses it appears in.
//
//	go run ./cmd/names -lang nl            # writes data/names/nl.json
//	go run ./cmd/names -lang nl -report    # also lists what it was unsure about
//
// A word of the translation is a candidate when it is written with a capital letter in the
// middle of sentences (a name), and it wins for an entity when it turns up in nearly all of the
// entity's verses, rarely elsewhere, and looks like the English name. Whatever it cannot place
// with confidence is left out, and the site then shows the English name.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"github.com/almerlucke/metabible/pkg/atlas"
)

var wordRe = regexp.MustCompile(`\p{L}[\p{L}’'\-]*`)

type stat struct{ capMid, lowMid, initial int }

func main() {
	lang := flag.String("lang", "nl", "language code (reads data/bible/<lang>.json)")
	report := flag.Bool("report", false, "print the uncertain and missing names")
	minScore := flag.Float64("min", 0.5, "minimum score to accept a name")
	dbg := flag.String("debug", "", "print the result for this person id")
	flag.Parse()

	a, err := atlas.Load("data/json", "")
	check(err)
	var bible map[string]string
	raw, err := os.ReadFile(fmt.Sprintf("data/bible/%s.json", *lang))
	check(err)
	check(json.Unmarshal(raw, &bible))

	// words of every verse, in their original case; which words are names
	verses := map[string][]string{} // verse record id -> distinct words
	stats := map[string]*stat{}
	df := map[string]int{}
	for _, v := range a.Verses {
		text, ok := bible[v.VID]
		if !ok {
			continue
		}
		toks := wordRe.FindAllStringIndex(text, -1)
		seen := map[string]bool{}
		var ws []string
		for i, loc := range toks {
			w := strings.Trim(text[loc[0]:loc[1]], "-’'")
			if w == "" {
				continue
			}
			low := strings.ToLower(w)
			s := stats[low]
			if s == nil {
				s = &stat{}
				stats[low] = s
			}
			initial := i == 0
			if !initial && loc[0] > 1 {
				switch text[:loc[0]][len(text[:loc[0]])-1] {
				case ' ':
					prev := strings.TrimRight(text[:loc[0]], " ")
					if prev != "" && strings.ContainsAny(prev[len(prev)-1:], ".!?") {
						initial = true
					}
				}
			}
			up := unicode.IsUpper([]rune(w)[0])
			switch {
			case initial:
				s.initial++
			case up:
				s.capMid++
			default:
				s.lowMid++
			}
			if !seen[w] {
				seen[w] = true
				ws = append(ws, w)
				df[w]++
			}
		}
		verses[v.ID] = ws
	}
	isName := func(w string) bool {
		s := stats[strings.ToLower(w)]
		if s == nil || !unicode.IsUpper([]rune(w)[0]) {
			return false
		}
		if s.capMid > 0 {
			return float64(s.capMid) >= 0.9*float64(s.capMid+s.lowMid)
		}
		return s.lowMid == 0 && s.initial <= 40 // only ever seen opening a verse: a rare name
	}

	// best picks the word that spells the entity's name in the translation. A word that looks
	// like the English name wins (Isaac / Izak); otherwise a word that is in nearly all of the
	// entity's verses and few others does (John / Johannes).
	best := func(name string, vs []string) (string, float64, int) {
		var ids []string
		for _, id := range vs {
			if _, ok := verses[id]; ok {
				ids = append(ids, id)
			}
		}
		n := len(ids)
		if n == 0 {
			return "", 0, 0
		}
		count := map[string]int{}
		for _, id := range ids {
			for _, w := range verses[id] {
				if isName(w) {
					count[w]++
				}
			}
		}
		var sw, dw string
		var ss, ds float64
		for w, c := range count {
			dice := 2 * float64(c) / float64(n+df[w])
			sim := similarity(name, w)
			if sim >= 0.8 && (n <= 2 || float64(c) >= 0.3*float64(n)) {
				// the unfolded spelling counts too, so Hanan beats the look-alike Anan
				if sc := sim + 0.3*dice + 0.6*plain(name, w); sc > ss || (sc == ss && w < sw) {
					sw, ss = w, sc
				}
			}
			if sim >= 0.4 && n >= 2 && float64(c) >= 0.6*float64(n) && dice >= 0.5 {
				if dice > ds || (dice == ds && w < dw) {
					dw, ds = w, dice
				}
			}
		}
		if sw != "" {
			return sw, 1, n
		}
		if dw != "" {
			return dw, 0.6 + 0.4*ds, n
		}
		// nothing convincing: report the likeliest for the review list
		var lw string
		var ls float64
		for w, c := range count {
			if sc := similarity(name, w) + 0.3*2*float64(c)/float64(n+df[w]); sc > ls {
				lw, ls = w, sc
			}
		}
		return lw, 0, n
	}

	out := struct {
		People map[string]string `json:"people"`
		Places map[string]string `json:"places"`
	}{map[string]string{}, map[string]string{}}
	var unsure, missing []string
	var glue struct {
		Words     map[string]string `json:"words"`
		Overrides map[string]string `json:"overrides"` // whole English names
		Names     map[string]string `json:"names"`     // single English names that cannot be found by alignment
		Spelling  map[string]string `json:"spelling"`  // modern spelling of words the Bible text writes without a diaeresis (Azie -> Azië)
	}
	if raw, err := os.ReadFile(fmt.Sprintf("data/names/%s.glue.json", *lang)); err == nil {
		check(json.Unmarshal(raw, &glue))
	}
	parenRe := regexp.MustCompile(`^(.*?)\s*\(([^)]*)\)\s*$`)
	known := map[string]map[string]int{} // English single word -> local spelling -> verses
	vote := func(en, local string, n int) {
		if known[en] == nil {
			known[en] = map[string]int{}
		}
		known[en][local] += n
	}
	type ent struct {
		kind, id, name string
		vs             []string
	}
	var ents []ent
	for _, p := range a.People {
		if p.ID == *dbg {
			w, sc, n := best(p.Name, p.Vs)
			fmt.Fprintf(os.Stderr, "debug %q %q %.2f %d\n", p.Name, w, sc, n)
		}
		ents = append(ents, ent{"person", p.ID, p.Name, p.Vs})
	}
	for _, p := range a.Places {
		ents = append(ents, ent{"place", p.ID, p.Name, p.Vs})
	}
	isSingle := func(s string) bool { return !strings.ContainsAny(strings.TrimSpace(s), " '’") }
	// pass 1: single words are aligned with the text
	single := map[string]string{} // entity id -> local spelling of its base name
	for _, e := range ents {
		base := e.name
		if m := parenRe.FindStringSubmatch(e.name); m != nil {
			base = m[1]
		}
		if glue.Overrides[e.name] != "" || !isSingle(base) {
			continue
		}
		w, s, n := best(base, e.vs)
		if nm := glue.Names[base]; nm != "" {
			w, s = nm, 1
			if n == 0 {
				n = 1
			}
		}
		// a genitive ("Simons") is not the name
		if len(w) > 3 && strings.HasSuffix(w, "s") && !strings.HasSuffix(base, "s") && isName(w[:len(w)-1]) && df[w[:len(w)-1]]*2 >= df[w] {
			w = w[:len(w)-1]
		}
		switch {
		case w == "":
			missing = append(missing, fmt.Sprintf("%s\t%s\t(no verses in this translation)", e.kind, e.name))
		case s < *minScore:
			unsure = append(unsure, fmt.Sprintf("%s\t%s\t-> %s\t%.2f (%d verses)", e.kind, e.name, w, s, n))
		default:
			single[e.id] = w
			vote(base, w, n+1)
			if s < 1 {
				unsure = append(unsure, fmt.Sprintf("%s\t%s\t=> %s\t%.2f (%d verses)", e.kind, e.name, w, s, n))
			}
		}
	}
	spelling := func(en string) (string, bool) {
		if nm := glue.Names[en]; nm != "" {
			return nm, true
		}
		var bw string
		bn := 0
		for w, n := range known[en] {
			if n > bn || (n == bn && w < bw) {
				bw, bn = w, n
			}
		}
		return bw, bw != ""
	}
	// pass 2: names of several words (and the bit in brackets) are put together word by word
	// compose puts a name of several words together: glue words from the glue file, names from
	// the alignment. ok is false when a plain English word is left over.
	compose := func(en string, capFirst bool) (res string, changed, ok bool) {
		fields := strings.Fields(en)
		ok = true
		for i, f := range fields {
			core := strings.Trim(f, ",.")
			if core == "" {
				continue
			}
			if w, has := glue.Words[strings.ToLower(core)]; has {
				if w == "" {
					fields[i] = ""
				} else {
					if i == 0 && capFirst && unicode.IsUpper([]rune(core)[0]) {
						r := []rune(w)
						r[0] = unicode.ToUpper(r[0])
						w = string(r)
					}
					fields[i] = strings.Replace(f, core, w, 1)
				}
				changed = true
			} else if w, has := spelling(core); has && w != core {
				fields[i] = strings.Replace(f, core, w, 1)
				changed = true
			} else if !unicode.IsUpper([]rune(core)[0]) && !unicode.IsDigit([]rune(core)[0]) {
				ok = false
			}
		}
		var keep []string
		for _, f := range fields {
			if f != "" {
				keep = append(keep, f)
			}
		}
		return strings.Join(keep, " "), changed, ok
	}
	for _, e := range ents {
		dst := out.People
		if e.kind == "place" {
			dst = out.Places
		}
		if o := glue.Overrides[e.name]; o != "" {
			dst[e.id] = o
			continue
		}
		base, suffix := e.name, ""
		if m := parenRe.FindStringSubmatch(e.name); m != nil {
			base, suffix = m[1], m[2]
		}
		name, changed := single[e.id], single[e.id] != "" && single[e.id] != base
		if single[e.id] == "" {
			name = base
			if !isSingle(base) {
				var ok bool
				if name, changed, ok = compose(base, true); !ok {
					continue
				}
			}
		}
		if suffix != "" {
			sx, ch, ok := compose(suffix, false)
			if !ok {
				continue
			}
			changed = changed || ch
			name += " (" + sx + ")"
		}
		if changed {
			dst[e.id] = name
		}
	}

	// the Bible text spells Azie, Syrie... without the diaeresis; the site uses the modern spelling
	wordRe2 := regexp.MustCompile(`[\p{L}]+`)
	modern := func(m map[string]string) {
		for id, name := range m {
			m[id] = wordRe2.ReplaceAllStringFunc(name, func(w string) string {
				if r, ok := glue.Spelling[w]; ok {
					return r
				}
				return w
			})
		}
	}
	modern(out.People)
	modern(out.Places)

	if *report {
		sort.Strings(unsure)
		fmt.Println("== unsure (=> accepted, -> rejected)")
		for _, l := range unsure {
			fmt.Println(l)
		}
		fmt.Println("== not attempted")
		for _, l := range missing {
			fmt.Println(l)
		}
	}
	buf, err := json.MarshalIndent(out, "", " ")
	check(err)
	check(os.MkdirAll("data/names", 0o755))
	check(os.WriteFile(fmt.Sprintf("data/names/%s.json", *lang), append(buf, '\n'), 0o644))
	fmt.Fprintf(os.Stderr, "%d people, %d places renamed; %d unsure, %d not attempted\n", len(out.People), len(out.Places), len(unsure), len(missing))
}

// similarity compares two names after folding the spelling differences English and the
// translation's language usually have (ph/f, c/k, th/t, doubled vowels...). 0..1.
func similarity(en, local string) float64 {
	a, b := fold(en), fold(local)
	if a == "" || b == "" {
		return 0
	}
	d := lev(a, b)
	m := max(len([]rune(a)), len([]rune(b)))
	return 1 - float64(d)/float64(m)
}

var folder = strings.NewReplacer("ph", "f", "th", "t", "ck", "k", "ch", "", "h", "", "c", "k", "z", "s", "w", "v", "y", "i", "j", "i", "q", "k", "x", "ks",
	"aa", "a", "ee", "e", "oo", "o", "uu", "u", "ii", "i", "ou", "o", "ai", "e", "ei", "e", "ij", "i", "ae", "e", "oe", "u", "-", "", "'", "", "’", "", " ", "")

func fold(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder
	for _, r := range s {
		switch r {
		case 'ë', 'é', 'è', 'ê':
			r = 'e'
		case 'ï', 'î':
			r = 'i'
		case 'ö', 'ô':
			r = 'o'
		case 'ü', 'û':
			r = 'u'
		case 'á', 'à', 'â':
			r = 'a'
		}
		b.WriteRune(r)
	}
	s = folder.Replace(b.String())
	// collapse doubled letters
	var out []rune
	for _, r := range s {
		if len(out) > 0 && out[len(out)-1] == r {
			continue
		}
		out = append(out, r)
	}
	return string(out)
}

// plain compares two names letter for letter, ignoring only case. 0..1.
func plain(a, b string) float64 {
	a, b = strings.ToLower(a), strings.ToLower(b)
	m := max(len([]rune(a)), len([]rune(b)))
	if m == 0 {
		return 0
	}
	return 1 - float64(lev(a, b))/float64(m)
}

func lev(a, b string) int {
	ra, rb := []rune(a), []rune(b)
	prev := make([]int, len(rb)+1)
	for j := range prev {
		prev[j] = j
	}
	for i := 1; i <= len(ra); i++ {
		cur := make([]int, len(rb)+1)
		cur[0] = i
		for j := 1; j <= len(rb); j++ {
			c := 1
			if ra[i-1] == rb[j-1] {
				c = 0
			}
			cur[j] = min(prev[j]+1, cur[j-1]+1, prev[j-1]+c)
		}
		prev = cur
	}
	return prev[len(rb)]
}

func check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
