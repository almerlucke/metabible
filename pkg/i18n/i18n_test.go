package i18n

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"sort"
	"strings"
	"testing"
)

var placeholder = regexp.MustCompile(`\{[a-z0-9_]+\}`)

func phs(s string) []string {
	m := placeholder.FindAllString(s, -1)
	sort.Strings(m)
	return m
}

// Every language must have every English key, with the same {placeholders}; and no stray keys.
func TestLocalesAreComplete(t *testing.T) {
	b := MustLoad()
	en := b.Keys(Default)
	for _, l := range b.Langs() {
		if l.Code == Default {
			continue
		}
		have := map[string]bool{}
		for _, k := range b.Keys(l.Code) {
			have[k] = true
		}
		var missing []string
		for _, k := range en {
			if !have[k] {
				missing = append(missing, k)
				continue
			}
			e, _ := b.Raw(Default, k)
			x, _ := b.Raw(l.Code, k)
			if strings.Join(phs(e), ",") != strings.Join(phs(x), ",") {
				t.Errorf("%s %s: placeholders differ: %q vs %q", l.Code, k, e, x)
			}
		}
		if len(missing) > 0 {
			sample := missing
			if len(sample) > 15 {
				sample = sample[:15]
			}
			t.Errorf("%s is missing %d keys, e.g. %v", l.Code, len(missing), sample)
		}
		for k := range have {
			if _, ok := b.Raw(Default, k); !ok {
				t.Errorf("%s has key %q that English lacks", l.Code, k)
			}
		}
	}
}

func TestPluralAndFallback(t *testing.T) {
	b := MustLoad()
	if got := b.Plural("nl", "count.verse", 1); got != "1 vers" {
		t.Errorf("one: %q", got)
	}
	if got := b.Plural("nl", "count.verse", 3); got != "3 verzen" {
		t.Errorf("other: %q", got)
	}
	if got := b.T("nl", "definitely.missing.key"); got != "definitely.missing.key" {
		t.Errorf("missing keys show themselves: %q", got)
	}
	if got := b.T("xx", "nav.story"); got != "Story" {
		t.Errorf("unknown language falls back to English: %q", got)
	}
}

func TestDetect(t *testing.T) {
	b := MustLoad()
	cases := []struct{ url, cookie, accept, want string }{
		{"/", "", "", "en"},
		{"/?lang=nl", "", "en", "nl"},
		{"/", "nl", "en", "nl"},
		{"/", "", "nl-NL,nl;q=0.9,en;q=0.8", "nl"},
		{"/", "", "fr-FR,fr;q=0.9", "en"},
		{"/?lang=zz", "", "", "en"},
	}
	for _, c := range cases {
		r := httptest.NewRequest("GET", c.url, nil)
		if c.cookie != "" {
			r.AddCookie(&http.Cookie{Name: CookieName, Value: c.cookie})
		}
		if c.accept != "" {
			r.Header.Set("Accept-Language", c.accept)
		}
		if got := b.Detect(r); got != c.want {
			t.Errorf("%+v -> %s", c, got)
		}
	}
}
