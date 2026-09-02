// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import "testing"

func TestTextBoxOptions(t *testing.T) {
	value := "Rosaline"
	changes := 0
	submits := 0

	box := TextBox(&value).
		Placeholder("Your name").
		Password().
		Width(36).
		OnChange(func(string) { changes++ }).
		OnSubmit(func(string) { submits++ }).
		Focus()

	if box.value != &value {
		t.Fatal("TextBox did not keep the supplied value pointer")
	}
	if box.placeholder != "Your name" || !box.password || box.columns != 36 || !box.focus {
		t.Fatalf("TextBox options were not preserved: %#v", box)
	}
	if box.onChange == nil || box.onSubmit == nil || changes != 0 || submits != 0 {
		t.Fatal("TextBox event handlers were not stored correctly")
	}
}

func TestTextInputsAcceptNilPointers(t *testing.T) {
	if TextBox(nil).value == nil {
		t.Fatal("TextBox(nil) should create safe internal storage")
	}
	if TextArea(nil).value == nil {
		t.Fatal("TextArea(nil) should create safe internal storage")
	}
	if CheckBox("Test", nil).value == nil {
		t.Fatal("CheckBox(nil) should create safe internal storage")
	}
}

func TestTextAreaSize(t *testing.T) {
	area := TextArea(new(string)).Size(52, 9).Focus()
	if area.columns != 52 || area.lines != 9 || !area.focus {
		t.Fatalf("TextArea options were not preserved: %#v", area)
	}
}

func TestBoolValue(t *testing.T) {
	for _, input := range []string{"1", "true", "TRUE", " yes ", "on"} {
		if !boolValue(input) {
			t.Errorf("boolValue(%q) = false, want true", input)
		}
	}
	for _, input := range []string{"", "0", "false", "no", "off", "anything"} {
		if boolValue(input) {
			t.Errorf("boolValue(%q) = true, want false", input)
		}
	}
}
