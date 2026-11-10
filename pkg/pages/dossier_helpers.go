package pages

import (
	"context"
	"fmt"
	"math"
	"net/url"

	"github.com/a-h/templ"
	"github.com/almerlucke/metabible/pkg/atlas"
)

func segStyle(a, b float64) templ.SafeCSS {
	return templ.SafeCSS(fmt.Sprintf("left:%.2f%%;width:%.2f%%", a, math.Max(b-a, 0.8)))
}

func eraSegStyle(e atlas.Era) templ.SafeCSS {
	a, b := atlas.Pct(e.From), atlas.Pct(e.To)
	return templ.SafeCSS(fmt.Sprintf("left:%.2f%%;width:%.2f%%;background:%s", a, b-a, e.Color))
}

func colorStyle(c string) templ.SafeCSS { return templ.SafeCSS("color:" + c) }

func heatStyle(h float64) templ.SafeCSS {
	return templ.SafeCSS(fmt.Sprintf("background:rgba(245,196,81,%.2f)", 0.08+0.9*math.Pow(h, 0.6)))
}

func imgCredit(i *atlas.Image) string {
	a := i.Artist
	if len([]rune(a)) > 48 {
		a = string([]rune(a)[:48]) + "…"
	}
	if a == "" {
		return i.License
	}
	return a + " · " + i.License
}

func artCredit(a atlas.Art) string {
	if a.Artist != "" {
		return "Gustave Doré · " + a.License
	}
	return a.License
}

func yr(y int) string { return atlas.FormatYear(y) }

func eras() []atlas.Era { return atlas.Eras }

func href(kind, id string) templ.SafeURL { return templ.SafeURL("/" + kind + "/" + id) }

func plural(n int, one, many string) string {
	if n == 1 {
		return fmt.Sprintf("%d %s", n, one)
	}
	return fmt.Sprintf("%d %s", n, many)
}

type baseKey struct{}

// WithBase stores the site's base URL (scheme and host) for building absolute links.
func WithBase(ctx context.Context, base string) context.Context {
	return context.WithValue(ctx, baseKey{}, base)
}

func baseURL(ctx context.Context) string {
	if b, ok := ctx.Value(baseKey{}).(string); ok {
		return b
	}
	return ""
}

// ogImage is the absolute URL of the share-card image for a page title.
func ogImage(ctx context.Context, title string) string {
	return baseURL(ctx) + "/og.png?t=" + url.QueryEscape(title)
}
