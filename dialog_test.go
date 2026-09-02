// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import "testing"

func TestNormalizeExtension(t *testing.T) {
	tests := map[string]string{
		"png":    ".png",
		".JPG":   ".jpg",
		"*.WebP": ".webp",
		" *.* ":  "*",
		"*":      "*",
		"":       "",
	}
	for input, want := range tests {
		if got := normalizeExtension(input); got != want {
			t.Errorf("normalizeExtension(%q) = %q, want %q", input, got, want)
		}
	}
}
