package pages

import (
	"fmt"
	"math"

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

func yr(y int) string { return atlas.FormatYear(y) }

func eras() []atlas.Era { return atlas.Eras }

func href(kind, id string) templ.SafeURL { return templ.SafeURL("/" + kind + "/" + id) }

func plural(n int, one, many string) string {
	if n == 1 {
		return fmt.Sprintf("%d %s", n, one)
	}
	return fmt.Sprintf("%d %s", n, many)
}
