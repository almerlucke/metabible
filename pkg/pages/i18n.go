package pages

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	"github.com/almerlucke/metabible/pkg/i18n"
)

var bundle *i18n.Bundle

// SetBundle gives the templates their translations. Call once at start-up.
func SetBundle(b *i18n.Bundle) { bundle = b }

func langOf(ctx context.Context) string { return i18n.FromContext(ctx) }

// t translates a key for the request's language.
func t(ctx context.Context, key string, kv ...any) string {
	if bundle == nil {
		return key
	}
	return bundle.T(i18n.FromContext(ctx), key, kv...)
}

// tn translates a counted key ("key.one" / "key.other"), filling {n}.
func tn(ctx context.Context, key string, n int, kv ...any) string {
	if bundle == nil {
		return key
	}
	return bundle.Plural(i18n.FromContext(ctx), key, n, kv...)
}

func langs() []i18n.Lang {
	if bundle == nil {
		return nil
	}
	return bundle.Langs()
}

var (
	bcRe = regexp.MustCompile(`(\d+) BC\b`)
	adRe = regexp.MustCompile(`(\d+) AD\b`)
)

// locDate converts the English dates built into the data ("1446 BC", "AD" years) to the request's language.
func locDate(ctx context.Context, s string) string {
	if bundle == nil || i18n.FromContext(ctx) == i18n.Default {
		return s
	}
	bc, ad := t(ctx, "date.bc"), t(ctx, "date.ad")
	s = bcRe.ReplaceAllString(s, "${1} "+strings.ReplaceAll(bc, "$", "$$"))
	return adRe.ReplaceAllString(s, "${1} "+strings.ReplaceAll(ad, "$", "$$"))
}

// count writes "3 verses" / "1 verse" using the keys count.<noun>.one / count.<noun>.other.
func count(ctx context.Context, noun string, n int) string { return tn(ctx, "count."+noun, n) }

// yrl writes a year ("1446 BC") in the request's language.
func yrl(ctx context.Context, y int) string {
	switch {
	case y < 0:
		return strconvI(-y) + " " + t(ctx, "date.bc")
	case y == 0:
		return "1 " + t(ctx, "date.bc")
	}
	return strconvI(y) + " " + t(ctx, "date.ad")
}

func strconvI(n int) string { return strings.TrimSpace(strings.Join(strings.Fields(itoa(n)), "")) }

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b []byte
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	if neg {
		b = append([]byte{'-'}, b...)
	}
	return string(b)
}

// slugKey makes "Poetry-Wisdom" / "Old Testament" into "poetry-wisdom" / "old-testament" for key lookups.
func slugKey(s string) string {
	return strings.ReplaceAll(strings.ToLower(strings.TrimSpace(s)), " ", "-")
}

// lookup translates data values (gender, place kind, book group…) through keys like gender.male;
// a value without a key is shown as it is.
func lookup(ctx context.Context, group, value string) string {
	if value == "" || bundle == nil {
		return value
	}
	key := group + "." + slugKey(value)
	if s, ok := bundle.Raw(i18n.FromContext(ctx), key); ok {
		return s
	}
	return value
}

func genderName(ctx context.Context, g string) string    { return lookup(ctx, "gender", g) }
func divName(ctx context.Context, d string) string       { return lookup(ctx, "div", d) }
func testamentName(ctx context.Context, d string) string { return lookup(ctx, "testament", d) }

// kindName translates a place kind such as "City", "Water River" or "river" word by word.
func kindName(ctx context.Context, k string) string {
	words := strings.Fields(k)
	for i, w := range words {
		words[i] = lookup(ctx, "kind", w)
	}
	return strings.Join(words, " ")
}

// eraName is an era's name in the request's language.
func eraName(ctx context.Context, id, fallback string) string {
	if bundle == nil {
		return fallback
	}
	if s, ok := bundle.Raw(i18n.FromContext(ctx), "era."+id+".name"); ok {
		return s
	}
	return fallback
}

// durText writes a length of time given in days ("6 days", "1 day", "2 years").
func durText(ctx context.Context, days int) string {
	switch {
	case days <= 0:
		return ""
	case days == 1:
		return tn(ctx, "duration.day", 1)
	case days >= 365:
		y := float64(days) / 365
		if y == float64(int(y)) {
			return tn(ctx, "duration.year", int(y))
		}
		return t(ctx, "duration.year.other", "n", strconv.FormatFloat(y, 'f', 1, 64))
	case days >= 30:
		return tn(ctx, "duration.month", days/30)
	}
	return tn(ctx, "duration.day", days)
}
