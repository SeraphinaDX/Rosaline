// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import "testing"

func TestTextAreaDocumentStateWithoutWindow(t *testing.T) {
	text := "A rose"
	area := TextArea(&text)
	if area.Text() != "A rose" || area.Modified() {
		t.Fatalf("initial text = %q, modified %v", area.Text(), area.Modified())
	}

	area.Append(" garden")
	if area.Text() != "A rose garden" || !area.Modified() {
		t.Fatalf("appended text = %q, modified %v", area.Text(), area.Modified())
	}

	area.MarkSaved()
	if area.Modified() {
		t.Fatal("MarkSaved did not clear modified state")
	}

	area.SetText("New page")
	if area.Text() != "New page" || !area.Modified() {
		t.Fatalf("replacement text = %q, modified %v", area.Text(), area.Modified())
	}

	area.Clear()
	if area.Text() != "" || !area.Modified() {
		t.Fatalf("cleared text = %q, modified %v", area.Text(), area.Modified())
	}
}

func TestTextAreaReplaceAllAndFindWithoutWindow(t *testing.T) {
	text := "rose, rose, violet"
	area := TextArea(&text)
	if !area.FindNext("violet") || area.FindNext("missing") || area.FindNext("") {
		t.Fatal("FindNext returned an unexpected result")
	}
	if count := area.ReplaceAll("rose", "peony"); count != 2 {
		t.Fatalf("ReplaceAll count = %d, want 2", count)
	}
	if text != "peony, peony, violet" || !area.Modified() {
		t.Fatalf("ReplaceAll text = %q, modified %v", text, area.Modified())
	}
	if count := area.ReplaceAll("", "x"); count != 0 {
		t.Fatalf("empty search replaced %d values", count)
	}
}

func TestTextAreaOptionsAndSafeUnmountedCommands(t *testing.T) {
	area := TextArea(nil).Size(-1, 0).Expand().Focus()
	area.OnChange(func(string) {}).OnCursorMove(func(TextPosition) {})
	area.Undo()
	area.Redo()
	area.Cut()
	area.Copy()
	area.Paste()
	area.SelectAll()

	if area.columns != 40 || area.lines != 6 || !area.expand || !area.focus {
		t.Fatalf("options = columns %d, lines %d, expand %v, focus %v",
			area.columns, area.lines, area.expand, area.focus)
	}
	if position := area.Cursor(); position != (TextPosition{Line: 1}) {
		t.Fatalf("initial cursor = %#v", position)
	}
}

func TestParseTextPosition(t *testing.T) {
	for input, want := range map[string]TextPosition{
		"1.0":     {Line: 1, Column: 0},
		"12.34":   {Line: 12, Column: 34},
		"bad":     {Line: 1},
		"0.2":     {Line: 1},
		"2.-1":    {Line: 1},
		"2.other": {Line: 1},
	} {
		if got := parseTextPosition(input); got != want {
			t.Errorf("parseTextPosition(%q) = %#v, want %#v", input, got, want)
		}
	}
}
