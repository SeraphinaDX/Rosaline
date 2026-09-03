// SPDX-License-Identifier: LGPL-3.0-or-later

package main

import "testing"

func TestDocumentNameAndTitle(t *testing.T) {
	document := &document{}
	if document.name() != "Untitled" || document.title(false) != "Untitled — Rosaline Notepad" {
		t.Fatalf("untitled document = %q, %q", document.name(), document.title(false))
	}
	document.path = "/tmp/rose garden.txt"
	if document.name() != "rose garden.txt" {
		t.Fatalf("document name = %q", document.name())
	}
	if title := document.title(true); title != "rose garden.txt * — Rosaline Notepad" {
		t.Fatalf("modified title = %q", title)
	}
}

func TestTextStats(t *testing.T) {
	for _, test := range []struct {
		text         string
		words, lines int
	}{
		{text: "", words: 0, lines: 1},
		{text: "one rose", words: 2, lines: 1},
		{text: "one\ntwo three\n", words: 3, lines: 3},
	} {
		words, lines := textStats(test.text)
		if words != test.words || lines != test.lines {
			t.Errorf("textStats(%q) = %d, %d; want %d, %d",
				test.text, words, lines, test.words, test.lines)
		}
	}
}
