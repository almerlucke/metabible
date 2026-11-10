package og

import (
	"bytes"
	"image/png"
	"testing"
)

func TestRender(t *testing.T) {
	b := Render("The Judgeship of Jair and the Ammonite Oppression That Followed After Him", "A subtitle that is rather long and keeps going so that it must wrap onto a second line of text")
	img, err := png.Decode(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	if img.Bounds().Dx() != W || img.Bounds().Dy() != H {
		t.Errorf("size %v", img.Bounds())
	}
	if len(Render("", "")) == 0 {
		t.Error("empty input should still render")
	}
}
