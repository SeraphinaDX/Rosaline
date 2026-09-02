// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import (
	"bytes"
	"fmt"
	"image/png"
	"os"

	"github.com/gen2brain/avif"
)

// AVIFOptions controls AVIF image export. Zero values use Rosaline's
// high-quality beginner-friendly defaults.
type AVIFOptions struct {
	Quality  int
	Speed    int
	Lossless bool
}

// SavePNG writes the picture as a PNG image.
func (p *Picture) SavePNG(path string) error {
	if p == nil || p.pixels == nil {
		return fmt.Errorf("rosaline: cannot save an empty picture")
	}
	var output bytes.Buffer
	if err := png.Encode(&output, p.pixels); err != nil {
		return fmt.Errorf("rosaline: could not encode PNG: %w", err)
	}
	if err := os.WriteFile(path, output.Bytes(), 0o644); err != nil {
		return fmt.Errorf("rosaline: could not save PNG %q: %w", path, err)
	}
	return nil
}

// SaveAVIF writes the picture as an AVIF image. Without options, Rosaline uses
// quality 90 and speed 8. AVIF encoding is CGo-free and needs no external tool.
func (p *Picture) SaveAVIF(path string, options ...AVIFOptions) error {
	if p == nil || p.pixels == nil {
		return fmt.Errorf("rosaline: cannot save an empty picture")
	}

	settings := AVIFOptions{Quality: 90, Speed: 8}
	if len(options) > 0 {
		settings = options[0]
		if settings.Quality <= 0 {
			settings.Quality = 90
		}
		if settings.Quality > 100 {
			settings.Quality = 100
		}
		if settings.Speed < 0 {
			settings.Speed = 8
		}
		if settings.Speed > 10 {
			settings.Speed = 10
		}
	}

	var output bytes.Buffer
	err := avif.Encode(&output, p.pixels, avif.Options{
		Quality:      settings.Quality,
		QualityAlpha: settings.Quality,
		Speed:        settings.Speed,
		Lossless:     settings.Lossless,
	})
	if err != nil {
		return fmt.Errorf("rosaline: could not encode AVIF: %w", err)
	}
	if err := os.WriteFile(path, output.Bytes(), 0o644); err != nil {
		return fmt.Errorf("rosaline: could not save AVIF %q: %w", path, err)
	}
	return nil
}

// Render draws an off-screen picture. It uses the same DrawingCanvas API as a
// visible Canvas widget and defaults to a white background.
func Render(width, height int, draw func(*DrawingCanvas)) *Picture {
	pixels, _ := renderDrawing(width, height, White, draw)
	return NewPicture(pixels)
}
