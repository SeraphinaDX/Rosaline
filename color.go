// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import (
	"fmt"
	"strconv"
	"strings"
)

// Color stores a red, green, blue, and alpha component.
type Color struct {
	R, G, B, A uint8
}

// RGB creates an opaque color.
func RGB(r, g, b uint8) Color { return Color{R: r, G: g, B: b, A: 255} }

// RGBA creates a color with an alpha component.
func RGBA(r, g, b, a uint8) Color { return Color{R: r, G: g, B: b, A: a} }

// Hex parses #RGB, #RRGGBB, or #RRGGBBAA. Invalid values return black.
// Use ParseHex when an invalid value should be reported as an error.
func Hex(value string) Color {
	color, err := ParseHex(value)
	if err != nil {
		return RGB(0, 0, 0)
	}
	return color
}

// ParseHex parses #RGB, #RRGGBB, or #RRGGBBAA.
func ParseHex(value string) (Color, error) {
	s := strings.TrimPrefix(strings.TrimSpace(value), "#")
	if len(s) == 3 {
		s = strings.Repeat(string(s[0]), 2) + strings.Repeat(string(s[1]), 2) + strings.Repeat(string(s[2]), 2)
	}
	if len(s) != 6 && len(s) != 8 {
		return Color{}, fmt.Errorf("rosaline: color %q must use #RGB, #RRGGBB, or #RRGGBBAA", value)
	}

	parse := func(part string) (uint8, error) {
		n, err := strconv.ParseUint(part, 16, 8)
		return uint8(n), err
	}

	r, err := parse(s[0:2])
	if err != nil {
		return Color{}, fmt.Errorf("rosaline: invalid color %q: %w", value, err)
	}
	g, err := parse(s[2:4])
	if err != nil {
		return Color{}, fmt.Errorf("rosaline: invalid color %q: %w", value, err)
	}
	b, err := parse(s[4:6])
	if err != nil {
		return Color{}, fmt.Errorf("rosaline: invalid color %q: %w", value, err)
	}
	a := uint8(255)
	if len(s) == 8 {
		a, err = parse(s[6:8])
		if err != nil {
			return Color{}, fmt.Errorf("rosaline: invalid color %q: %w", value, err)
		}
	}
	return RGBA(r, g, b, a), nil
}

func (c Color) String() string { return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B) }

var (
	Black       = RGB(0, 0, 0)
	White       = RGB(255, 255, 255)
	Rose        = Hex("#d64f8c")
	SoftRose    = Hex("#f4a6c8")
	Transparent = RGBA(0, 0, 0, 0)
)
