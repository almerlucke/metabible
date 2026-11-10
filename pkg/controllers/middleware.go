package controllers

import (
	"compress/gzip"
	"net/http"
	"strings"
	"sync"

	"github.com/almerlucke/metabible/pkg/pages"
)

var gzPool = sync.Pool{New: func() any { return gzip.NewWriter(nil) }}

type gzResponse struct {
	http.ResponseWriter
	gz      *gzip.Writer
	decided bool
	enabled bool
}

func compressible(ct string) bool {
	ct = strings.ToLower(ct)
	for _, p := range []string{"text/", "application/javascript", "application/json", "application/manifest", "image/svg", "application/gpx", "application/vnd.google-earth.kml"} {
		if strings.HasPrefix(ct, p) {
			return true
		}
	}
	return false
}

func (g *gzResponse) decide() {
	if g.decided {
		return
	}
	g.decided = true
	h := g.Header()
	if h.Get("Content-Encoding") == "" && compressible(h.Get("Content-Type")) {
		g.enabled = true
		h.Set("Content-Encoding", "gzip")
		h.Add("Vary", "Accept-Encoding")
		h.Del("Content-Length")
		g.gz = gzPool.Get().(*gzip.Writer)
		g.gz.Reset(g.ResponseWriter)
	}
}

func (g *gzResponse) WriteHeader(code int) {
	g.decide()
	g.ResponseWriter.WriteHeader(code)
}

func (g *gzResponse) Write(b []byte) (int, error) {
	g.decide()
	if g.enabled {
		return g.gz.Write(b)
	}
	return g.ResponseWriter.Write(b)
}

func (g *gzResponse) Flush() {
	if g.enabled {
		_ = g.gz.Flush()
	}
	if f, ok := g.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (g *gzResponse) close() {
	if g.enabled {
		_ = g.gz.Close()
		gzPool.Put(g.gz)
	}
}

// Middleware adds what every response needs: compression for text, long-lived caching for the
// bundled vendor files, and the request's base URL for absolute links in page metadata.
func Middleware(next http.Handler) http.Handler {
	next = Bundle.Middleware(next) // language detection sits inside, so it sees the final request
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		scheme := "http"
		if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
			scheme = "https"
		}
		r = r.WithContext(pages.WithBase(r.Context(), scheme+"://"+r.Host))

		p := r.URL.Path
		if strings.HasSuffix(p, ".webmanifest") {
			w.Header().Set("Content-Type", "application/manifest+json")
		}
		switch {
		case strings.HasPrefix(p, "/assets/vendor/"):
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		case strings.HasPrefix(p, "/assets/"):
			w.Header().Set("Cache-Control", "public, max-age=300")
		}

		// the JSON API compresses itself; leave ranges and non-GET alone
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") || r.Header.Get("Range") != "" || strings.HasPrefix(p, "/api/") {
			next.ServeHTTP(w, r)
			return
		}
		g := &gzResponse{ResponseWriter: w}
		defer g.close()
		next.ServeHTTP(g, r)
	})
}
