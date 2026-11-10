// arthash finds engravings in data/enrichment.json that are the same picture uploaded more than
// once (or cropped or recolored). It downloads a tiny thumbnail of each, computes a perceptual
// "difference hash", and writes the copies to data/art-duplicates.json for the site to skip.
//
//	go run ./cmd/arthash
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"math/bits"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type art struct {
	URL   string  `json:"url"`
	Title string  `json:"title"`
	Event string  `json:"event"`
	Score float64 `json:"score"`
}

type enrichment struct {
	Art []art `json:"art"`
}

type Dup struct {
	Title string `json:"title"`
	Of    string `json:"of"`
	Dist  int    `json:"distance"`
}

var client = &http.Client{Timeout: 30 * time.Second}

func fetch(u string) (image.Image, error) {
	for attempt := 0; attempt < 6; attempt++ {
		req, _ := http.NewRequest("GET", u, nil)
		req.Header.Set("User-Agent", "MetaBibleEnrich/1.0 (open-source Bible atlas)")
		resp, err := client.Do(req)
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}
		if resp.StatusCode == 429 || resp.StatusCode >= 500 {
			wait := 15 * time.Second
			if ra, err := strconv.Atoi(resp.Header.Get("Retry-After")); err == nil {
				wait = time.Duration(ra+2) * time.Second
			}
			resp.Body.Close()
			log.Printf("rate limited, waiting %s", wait)
			time.Sleep(wait)
			continue
		}
		if resp.StatusCode != 200 {
			resp.Body.Close()
			return nil, fmt.Errorf("%s", resp.Status)
		}
		b, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		img, _, err := image.Decode(strings.NewReader(string(b)))
		return img, err
	}
	return nil, fmt.Errorf("giving up")
}

// dhash: 9x8 grayscale samples, each bit says whether a pixel is brighter than its right neighbour.
func dhash(img image.Image) uint64 {
	b := img.Bounds()
	var g [8][9]float64
	for y := 0; y < 8; y++ {
		for x := 0; x < 9; x++ {
			// average a small block so noise and tiny crops matter less
			x0, x1 := b.Min.X+x*b.Dx()/9, b.Min.X+(x+1)*b.Dx()/9
			y0, y1 := b.Min.Y+y*b.Dy()/8, b.Min.Y+(y+1)*b.Dy()/8
			var sum, n float64
			for yy := y0; yy < y1 || yy == y0; yy++ {
				for xx := x0; xx < x1 || xx == x0; xx++ {
					r, gg, bb, _ := img.At(xx, yy).RGBA()
					sum += 0.299*float64(r) + 0.587*float64(gg) + 0.114*float64(bb)
					n++
				}
			}
			g[y][x] = sum / n
		}
	}
	var h uint64
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			h <<= 1
			if g[y][x] > g[y][x+1] {
				h |= 1
			}
		}
	}
	return h
}

func thumbURL(u string) string {
	p, err := url.Parse(u)
	if err != nil {
		return u
	}
	q := p.Query()
	q.Set("width", "96")
	p.RawQuery = q.Encode()
	return p.String()
}

func main() {
	in := flag.String("in", "data/enrichment.json", "enrichment file")
	out := flag.String("out", "data/art-duplicates.json", "output")
	maxDist := flag.Int("distance", 9, "largest Hamming distance (of 64 bits) still counted as the same picture")
	flag.Parse()

	b, err := os.ReadFile(*in)
	if err != nil {
		log.Fatal(err)
	}
	var e enrichment
	if err := json.Unmarshal(b, &e); err != nil {
		log.Fatal(err)
	}
	type hashed struct {
		a art
		h uint64
	}
	var hs []hashed
	for i, a := range e.Art {
		img, err := fetch(thumbURL(a.URL))
		if err != nil {
			log.Printf("skip %q: %v", a.Title, err)
			continue
		}
		hs = append(hs, hashed{a, dhash(img)})
		if i%20 == 0 {
			log.Printf("%d/%d hashed", i, len(e.Art))
		}
		time.Sleep(250 * time.Millisecond)
	}
	// keep a picture that already has an event match over one that has none, then the shorter title
	sort.SliceStable(hs, func(i, j int) bool {
		if (hs[i].a.Event != "") != (hs[j].a.Event != "") {
			return hs[i].a.Event != ""
		}
		return len(hs[i].a.Title) < len(hs[j].a.Title)
	})
	var kept []hashed
	var dups []Dup
	for _, c := range hs {
		dup := false
		for _, k := range kept {
			if d := bits.OnesCount64(c.h ^ k.h); d <= *maxDist {
				dups = append(dups, Dup{Title: c.a.Title, Of: k.a.Title, Dist: d})
				dup = true
				break
			}
		}
		if !dup {
			kept = append(kept, c)
		}
	}
	sort.Slice(dups, func(i, j int) bool { return dups[i].Title < dups[j].Title })
	j, _ := json.MarshalIndent(map[string]any{"generated": time.Now().UTC().Format("2006-01-02"), "duplicates": dups}, "", " ")
	if err := os.WriteFile(*out, j, 0o644); err != nil {
		log.Fatal(err)
	}
	log.Printf("%d pictures, %d unique, %d duplicates -> %s", len(hs), len(kept), len(dups), *out)
}
