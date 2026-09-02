// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import (
	"image"
	"os"
	"path/filepath"
	"testing"
)

func TestSavePNGRoundTrip(t *testing.T) {
	picture := exportTestPicture()
	path := filepath.Join(t.TempDir(), "drawing.png")
	if err := picture.SavePNG(path); err != nil {
		t.Fatalf("SavePNG: %v", err)
	}

	loaded, err := LoadImage(path)
	if err != nil {
		t.Fatalf("LoadImage: %v", err)
	}
	if loaded.Format() != "png" {
		t.Fatalf("format = %q, want png", loaded.Format())
	}
	if loaded.Width() != 64 || loaded.Height() != 64 {
		t.Fatalf("size = %dx%d, want 64x64", loaded.Width(), loaded.Height())
	}
}

func TestSaveAVIFRoundTrip(t *testing.T) {
	picture := exportTestPicture()
	path := filepath.Join(t.TempDir(), "drawing.avif")
	if err := picture.SaveAVIF(path, AVIFOptions{Quality: 90, Speed: 10}); err != nil {
		t.Fatalf("SaveAVIF: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat AVIF: %v", err)
	}
	if info.Size() == 0 {
		t.Fatal("saved AVIF is empty")
	}

	loaded, err := LoadImage(path)
	if err != nil {
		t.Fatalf("LoadImage AVIF: %v", err)
	}
	if loaded.Format() != "avif" {
		t.Fatalf("format = %q, want avif", loaded.Format())
	}
	if loaded.Width() != 64 || loaded.Height() != 64 {
		t.Fatalf("size = %dx%d, want 64x64", loaded.Width(), loaded.Height())
	}
}

func TestSaveEmptyPictureReturnsError(t *testing.T) {
	var picture *Picture
	if err := picture.SavePNG(filepath.Join(t.TempDir(), "empty.png")); err == nil {
		t.Fatal("SavePNG should reject an empty picture")
	}
	if err := picture.SaveAVIF(filepath.Join(t.TempDir(), "empty.avif")); err == nil {
		t.Fatal("SaveAVIF should reject an empty picture")
	}
}

func exportTestPicture() *Picture {
	pixels := image.NewNRGBA(image.Rect(0, 0, 64, 64))
	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			pixels.Set(x, y, nativeColor(RGB(uint8(x*4), uint8(y*4), 160)))
		}
	}
	return NewPicture(pixels)
}
