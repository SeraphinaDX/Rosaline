// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
)

func TestLoadImage(t *testing.T) {
	path := filepath.Join(t.TempDir(), "rose.png")
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	pixels := image.NewRGBA(image.Rect(0, 0, 7, 5))
	pixels.Set(2, 3, color.RGBA{R: 196, G: 63, B: 122, A: 255})
	if err := png.Encode(file, pixels); err != nil {
		file.Close()
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}

	picture, err := LoadImage(path)
	if err != nil {
		t.Fatal(err)
	}
	if picture.Width() != 7 || picture.Height() != 5 {
		t.Fatalf("image size = %dx%d, want 7x5", picture.Width(), picture.Height())
	}
	if picture.Format() != "png" || picture.Path() != path {
		t.Fatalf("unexpected picture metadata: format=%q path=%q", picture.Format(), picture.Path())
	}
}

func TestLoadImageErrorsAreHelpful(t *testing.T) {
	path := filepath.Join(t.TempDir(), "not-an-image.txt")
	if err := os.WriteFile(path, []byte("hello"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := LoadImage(path)
	if err == nil || !strings.Contains(err.Error(), path) {
		t.Fatalf("LoadImage error = %v, want filename", err)
	}
}

func TestLoadImageFS(t *testing.T) {
	var encoded bytes.Buffer
	pixels := image.NewRGBA(image.Rect(0, 0, 6, 4))
	pixels.Set(2, 1, color.RGBA{R: 196, G: 63, B: 122, A: 255})
	if err := png.Encode(&encoded, pixels); err != nil {
		t.Fatal(err)
	}

	picture, err := LoadImageFS(fstest.MapFS{
		"assets/rose.png": &fstest.MapFile{Data: encoded.Bytes()},
	}, "assets/rose.png")
	if err != nil {
		t.Fatal(err)
	}
	if picture.Width() != 6 || picture.Height() != 4 || picture.Path() != "assets/rose.png" || picture.Format() != "png" {
		t.Fatalf("unexpected embedded picture metadata: %#v", picture)
	}
}

func TestNewPictureAndImageWidget(t *testing.T) {
	picture := NewPicture(image.NewRGBA(image.Rect(0, 0, 9, 4)))
	if picture.Width() != 9 || picture.Height() != 4 || picture.Image() == nil {
		t.Fatalf("picture size = %dx%d, want 9x4", picture.Width(), picture.Height())
	}
	clicked := 0
	widget := Image(nil).Placeholder("Choose a picture").Fit(180, 120).Expand().OnClick(func() { clicked++ })
	widget.SetImage(picture)
	if widget.Picture() != picture || widget.placeholder != "Choose a picture" || widget.fitWidth != 180 || widget.fitHeight != 120 || !widget.expand || widget.onClick == nil {
		t.Fatal("image widget options were not preserved")
	}
	if clicked != 0 {
		t.Fatal("image click handler ran before an event")
	}
}

func TestFitImagePreservesAspectRatio(t *testing.T) {
	source := image.NewRGBA(image.Rect(0, 0, 8, 4))
	for y := range 4 {
		for x := range 8 {
			source.Set(x, y, color.RGBA{R: 196, G: 63, B: 122, A: 255})
		}
	}
	fitted := fitImage(source, 20, 20)
	if fitted.Bounds().Dx() != 20 || fitted.Bounds().Dy() != 20 {
		t.Fatalf("fitted size = %v, want 20x20", fitted.Bounds())
	}
	center := color.NRGBAModel.Convert(fitted.At(10, 10)).(color.NRGBA)
	edge := color.NRGBAModel.Convert(fitted.At(10, 1)).(color.NRGBA)
	if center.A == 0 || edge.A != 0 {
		t.Fatalf("fit did not center the 2:1 image: center=%#v edge=%#v", center, edge)
	}
}
