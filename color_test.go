// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import "testing"

func TestParseHex(t *testing.T) {
	tests := []struct {
		input string
		want  Color
	}{
		{"#f08", RGBA(255, 0, 136, 255)},
		{"#112233", RGBA(17, 34, 51, 255)},
		{"11223344", RGBA(17, 34, 51, 68)},
	}
	for _, test := range tests {
		got, err := ParseHex(test.input)
		if err != nil {
			t.Fatalf("ParseHex(%q): %v", test.input, err)
		}
		if got != test.want {
			t.Errorf("ParseHex(%q) = %#v, want %#v", test.input, got, test.want)
		}
	}
}

func TestParseHexRejectsInvalidValue(t *testing.T) {
	if _, err := ParseHex("pink"); err == nil {
		t.Fatal("ParseHex accepted an invalid color")
	}
}
