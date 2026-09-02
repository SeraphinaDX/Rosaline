// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
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

func TestNewPictureAndImageWidget(t *testing.T) {
	picture := NewPicture(image.NewRGBA(image.Rect(0, 0, 9, 4)))
	if picture.Width() != 9 || picture.Height() != 4 || picture.Image() == nil {
		t.Fatalf("picture size = %dx%d, want 9x4", picture.Width(), picture.Height())
	}
	widget := Image(nil).Placeholder("Choose a picture").Expand()
	widget.SetImage(picture)
	if widget.Picture() != picture || widget.placeholder != "Choose a picture" || !widget.expand {
		t.Fatal("image widget options were not preserved")
	}
}
