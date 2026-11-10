// Package i18n holds the site's translations. Each language is one JSON file in locales/
// (flat "dotted.key": "text"). English is the source and the fallback: a key missing in
// another language shows the English text, so an unfinished language still works.
//
// Texts may contain {placeholders}, filled from key/value pairs: T("hi", "name", "Ada").
// Plural forms are separate keys: "verses.one" and "verses.other" (see Plural).
// Keys starting with "js." are also sent to the browser (see Client).
package i18n

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

//go:embed locales/*.json
var files embed.FS

const Default = "en"

// Bundle is the set of loaded languages.
type Bundle struct {
	langs map[string]map[string]string
}

// Lang describes a language for the switcher.
type Lang struct {
	Code string
	Name string // in that language, e.g. "Nederlands"
}

// Load reads every embedded locale file.
func Load() (*Bundle, error) {
	entries, err := files.ReadDir("locales")
	if err != nil {
		return nil, err
	}
	b := &Bundle{langs: map[string]map[string]string{}}
	for _, e := range entries {
		code := strings.TrimSuffix(e.Name(), ".json")
		raw, err := files.ReadFile("locales/" + e.Name())
		if err != nil {
			return nil, err
		}
		m := map[string]string{}
		if err := json.Unmarshal(raw, &m); err != nil {
			return nil, fmt.Errorf("locale %s: %w", code, err)
		}
		b.langs[code] = m
	}
	if b.langs[Default] == nil {
		return nil, fmt.Errorf("i18n: no %s locale", Default)
	}
	return b, nil
}

// MustLoad is Load for program start-up.
func MustLoad() *Bundle {
	b, err := Load()
	if err != nil {
		panic(err)
	}
	return b
}

// Has reports whether a language is available.
func (b *Bundle) Has(code string) bool { _, ok := b.langs[code]; return ok }

// Codes lists the language codes.
func (b *Bundle) Codes() []string {
	var out []string
	for _, l := range b.Langs() {
		out = append(out, l.Code)
	}
	return out
}

// Langs lists the languages, English first.
func (b *Bundle) Langs() []Lang {
	var out []Lang
	for code, m := range b.langs {
		name := m["_name"]
		if name == "" {
			name = code
		}
		out = append(out, Lang{code, name})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Code == Default || out[j].Code == Default {
			return out[i].Code == Default
		}
		return out[i].Code < out[j].Code
	})
	return out
}

// Keys returns every key of a language (used by tests).
func (b *Bundle) Keys(code string) []string {
	var ks []string
	for k := range b.langs[code] {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	return ks
}

// Raw returns the text for a key without filling placeholders; ok is false if no language has it.
func (b *Bundle) Raw(lang, key string) (string, bool) {
	if s, ok := b.langs[lang][key]; ok {
		return s, true
	}
	s, ok := b.langs[Default][key]
	return s, ok
}

// T translates a key, filling {placeholders} from the key/value pairs.
func (b *Bundle) T(lang, key string, kv ...any) string {
	s, ok := b.Raw(lang, key)
	if !ok {
		return key // a missing key is visible, not silent
	}
	for i := 0; i+1 < len(kv); i += 2 {
		s = strings.ReplaceAll(s, "{"+fmt.Sprint(kv[i])+"}", fmt.Sprint(kv[i+1]))
	}
	return s
}

// Plural picks key+".one" for 1 and key+".other" otherwise, and fills {n} (plus any pairs).
func (b *Bundle) Plural(lang, key string, n int, kv ...any) string {
	form := ".other"
	if n == 1 {
		form = ".one"
	}
	return b.T(lang, key+form, append([]any{"n", n}, kv...)...)
}

// Client returns the texts the browser needs (keys starting with "js."), English filling any gaps.
func (b *Bundle) Client(lang string) map[string]string {
	out := map[string]string{}
	for k, v := range b.langs[Default] {
		if strings.HasPrefix(k, "js.") {
			out[k] = v
		}
	}
	for k, v := range b.langs[lang] {
		if strings.HasPrefix(k, "js.") {
			out[k] = v
		}
	}
	return out
}

// ---------------------------------------------------------------- request language

type ctxKey struct{}

func WithLang(ctx context.Context, lang string) context.Context {
	return context.WithValue(ctx, ctxKey{}, lang)
}

// FromContext is the language chosen for this request.
func FromContext(ctx context.Context) string {
	if l, ok := ctx.Value(ctxKey{}).(string); ok && l != "" {
		return l
	}
	return Default
}

const CookieName = "mb_lang"

// Detect chooses the language: ?lang=, then the cookie, then the browser's Accept-Language.
func (b *Bundle) Detect(r *http.Request) string {
	if l := r.URL.Query().Get("lang"); b.Has(l) {
		return l
	}
	if c, err := r.Cookie(CookieName); err == nil && b.Has(c.Value) {
		return c.Value
	}
	for _, part := range strings.Split(r.Header.Get("Accept-Language"), ",") {
		tag := strings.ToLower(strings.TrimSpace(strings.SplitN(part, ";", 2)[0]))
		if b.Has(tag) {
			return tag
		}
		if i := strings.IndexAny(tag, "-_"); i > 0 && b.Has(tag[:i]) {
			return tag[:i]
		}
	}
	return Default
}

// Middleware stores the language in the request context. A ?lang= link switches the language,
// remembers it in a cookie and redirects to the same page without the parameter.
func (b *Bundle) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if l := r.URL.Query().Get("lang"); l != "" && b.Has(l) && r.Method == http.MethodGet && !strings.HasPrefix(r.URL.Path, "/api/") {
			http.SetCookie(w, &http.Cookie{Name: CookieName, Value: l, Path: "/", MaxAge: 365 * 24 * 3600, SameSite: http.SameSiteLaxMode})
			q := r.URL.Query()
			q.Del("lang")
			u := r.URL.Path
			if enc := q.Encode(); enc != "" {
				u += "?" + enc
			}
			http.Redirect(w, r, u, http.StatusSeeOther)
			return
		}
		lang := b.Detect(r)
		w.Header().Add("Vary", "Cookie")
		w.Header().Add("Vary", "Accept-Language")
		next.ServeHTTP(w, r.WithContext(WithLang(r.Context(), lang)))
	})
}
