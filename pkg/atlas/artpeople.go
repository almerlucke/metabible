package atlas

import (
	"regexp"
	"strings"
)

var artWord = regexp.MustCompile(`[a-z]+`)

// Words that are part of many names or titles and never identify a person on their own.
var artNameStop = map[string]bool{
	"the": true, "of": true, "and": true, "son": true, "daughter": true, "father": true, "mother": true, "wife": true, "husband": true,
	"king": true, "queen": true, "high": true, "priest": true, "prophet": true, "apostle": true, "brother": true, "sister": true,
	"man": true, "woman": true, "child": true, "great": true, "first": true, "second": true, "elder": true, "younger": true,
	"judge": true, "judah": true, "israel": true, "lord": true, "god": true, "spirit": true, "holy": true, "angel": true,
	"sons": true, "levite": true, "pharaoh": true, "daughters": true, "servant": true, "captain": true, "baptist": true,
}

// artNameKeys maps a distinctive word to the one person it identifies: "jesus" and "christ" to Jesus Christ,
// "peter" to Simon Peter. A word counts when it is in only one prominent person's name, or when that person is
// by far the most mentioned of those who share it.
func (a *Atlas) artNameKeys() map[string]string {
	// everyone who shares a word counts, so a lesser-known namesake still makes the word ambiguous
	holders := map[string][]*Person{}
	for _, p := range a.People {
		if p.Name == "God" || p.Name == "Holy Spirit" {
			continue
		}
		base := strings.ToLower(parenRe.ReplaceAllString(p.Name, ""))
		seen := map[string]bool{}
		for _, w := range artWord.FindAllString(base, -1) {
			if len(w) >= 4 && !artNameStop[w] && !seen[w] {
				seen[w] = true
				holders[w] = append(holders[w], p)
			}
		}
	}
	keys := map[string]string{}
	for w, ps := range holders {
		best, second := ps[0], 0
		for _, p := range ps[1:] {
			if p.Verses > best.Verses {
				second, best = best.Verses, p
			} else if p.Verses > second {
				second = p.Verses
			}
		}
		if best.Verses >= 20 && (len(ps) == 1 || best.Verses >= 3*second) {
			keys[w] = best.ID
		}
	}
	return keys
}

var parenRe = regexp.MustCompile(`\s*\(.*?\)`)

// peopleInTitle finds the people an engraving's title names, by their distinctive names.
func peopleInTitle(title string, keys map[string]string) []string {
	var out []string
	seen := map[string]bool{}
	for _, w := range artWord.FindAllString(strings.ToLower(title), -1) {
		if id, ok := keys[w]; ok && !seen[id] {
			seen[id] = true
			out = append(out, id)
		}
	}
	return out
}
