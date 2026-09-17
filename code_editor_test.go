// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import "testing"

func TestGoCodeEditorUsesSafeEditorDefaults(t *testing.T) {
	source := "package main\n"
	editor := GoCodeEditor(&source).ReadOnly()
	if !editor.goSyntax || !editor.monospace || !editor.IsReadOnly() {
		t.Fatalf("code editor options were not retained: %#v", editor)
	}
	editor.Append("func main() {}\n")
	editor.Clear()
	if source != "package main\n" {
		t.Fatalf("read-only editor changed source to %q", source)
	}
	editor.SetText("package example\n")
	if source != "package example\n" {
		t.Fatal("programmatic SetText should work while read-only")
	}
	editor.SetReadOnly(false)
	if editor.IsReadOnly() {
		t.Fatal("SetReadOnly(false) did not make the editor writable")
	}
}

func TestGoSyntaxSpansRecognizeCommonTokens(t *testing.T) {
	source := "package main\n\n// hello\nvar name string = `Rosaline`\nvar count = 42\n"
	spans := goSyntaxSpans(source)
	want := map[string]bool{
		goKeywordTag: false,
		goCommentTag: false,
		goTypeTag:    false,
		goStringTag:  false,
		goNumberTag:  false,
	}
	for _, span := range spans {
		if _, exists := want[span.tag]; exists {
			want[span.tag] = true
		}
		if span.start < 0 || span.end <= span.start || span.end > len(source) {
			t.Fatalf("invalid syntax span: %#v", span)
		}
	}
	for tag, found := range want {
		if !found {
			t.Fatalf("Go syntax did not produce %s: %#v", tag, spans)
		}
	}
}

func TestTextOffsetIndexCountsUnicodeColumns(t *testing.T) {
	source := "package main\nvar flower = \"🌹\"\n"
	offset := len("package main\nvar flower = \"🌹")
	if got := textOffsetIndex(source, offset); got != "2.15" {
		t.Fatalf("unicode text index = %q, want 2.15", got)
	}
}

func TestGoToIsSafeBeforeMount(t *testing.T) {
	editor := GoCodeEditor(nil)
	editor.GoTo(12, 4)
}
