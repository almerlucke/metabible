// Package og draws the 1200×630 preview images shown when a page link is shared.
package og

import (
	"bytes"
	"hash/fnv"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"strings"
	"sync"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

const (
	W = 1200
	H = 630
)

var (
	once       sync.Once
	title, sub font.Face
	brand      font.Face
)

func faces() {
	once.Do(func() {
		bold, _ := opentype.Parse(gobold.TTF)
		reg, _ := opentype.Parse(goregular.TTF)
		title, _ = opentype.NewFace(bold, &opentype.FaceOptions{Size: 76, DPI: 72, Hinting: font.HintingFull})
		sub, _ = opentype.NewFace(reg, &opentype.FaceOptions{Size: 34, DPI: 72, Hinting: font.HintingFull})
		brand, _ = opentype.NewFace(bold, &opentype.FaceOptions{Size: 28, DPI: 72, Hinting: font.HintingFull})
	})
}

func wrap(face font.Face, text string, maxW int, maxLines int) []string {
	words := strings.Fields(text)
	var lines []string
	cur := ""
	for _, w := range words {
		try := strings.TrimSpace(cur + " " + w)
		if font.MeasureString(face, try).Ceil() > maxW && cur != "" {
			lines = append(lines, cur)
			cur = w
		} else {
			cur = try
		}
	}
	if cur != "" {
		lines = append(lines, cur)
	}
	if len(lines) > maxLines {
		lines = lines[:maxLines]
		lines[maxLines-1] = strings.TrimRight(lines[maxLines-1], " .,;") + "…"
	}
	return lines
}

func text(img draw.Image, face font.Face, s string, x, y int, c color.Color) {
	d := &font.Drawer{Dst: img, Src: image.NewUniform(c), Face: face, Dot: fixed.P(x, y)}
	d.DrawString(s)
}

// Render draws a night-sky card with the page title and an optional subtitle.
func Render(t, s string) []byte {
	faces()
	img := image.NewRGBA(image.Rect(0, 0, W, H))
	// vertical gradient, deep navy to near black, with a warm glow at the upper right
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			k := float64(y) / H
			r, g, b := 14-8*k, 18-8*k, 40-22*k
			dx, dy := float64(x-980), float64(y-120)
			glow := math.Exp(-(dx*dx + dy*dy) / (2 * 260 * 260))
			r += 70 * glow
			g += 48 * glow
			b += 10 * glow
			img.SetRGBA(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), 255})
		}
	}
	// stars, placed deterministically from the title
	h := fnv.New32a()
	h.Write([]byte(t))
	seed := h.Sum32()
	next := func() uint32 { seed = seed*1664525 + 1013904223; return seed >> 8 }
	for i := 0; i < 150; i++ {
		x, y := int(next()%W), int(next()%H)
		a := uint8(60 + next()%160)
		c := color.RGBA{245, 214, 130, a}
		if i%3 == 0 {
			c = color.RGBA{200, 215, 255, a}
		}
		img.Set(x, y, c)
		if i%9 == 0 {
			img.Set(x+1, y, c)
			img.Set(x, y+1, c)
		}
	}
	gold := color.RGBA{245, 196, 81, 255}
	// frame and brand
	for i := 0; i < 3; i++ {
		for x := 40; x < W-40; x++ {
			img.Set(x, 40+i, gold)
			img.Set(x, H-41-i, gold)
		}
		for y := 40; y < H-40; y++ {
			img.Set(40+i, y, gold)
			img.Set(W-41-i, y, gold)
		}
	}
	text(img, brand, "M E T A B I B L E", 90, 120, gold)

	lines := wrap(title, t, W-180, 3)
	y := 260
	if len(lines) == 1 {
		y = 300
	}
	for _, l := range lines {
		text(img, title, l, 90, y, color.RGBA{255, 232, 160, 255})
		y += 92
	}
	if s != "" {
		for i, l := range wrap(sub, s, W-180, 2) {
			text(img, sub, l, 90, y+20+i*46, color.RGBA{190, 196, 220, 255})
		}
	}
	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	return buf.Bytes()
}
