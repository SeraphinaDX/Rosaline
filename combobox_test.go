// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import "testing"

func TestComboBoxNormalizesOptionsAndSelection(t *testing.T) {
	value := "missing"
	combo := ComboBox(&value, "Fast", "Balanced", "Fast")
	if value != "Fast" || combo.Selected() != "Fast" {
		t.Fatalf("selection = %q, want Fast", value)
	}
	options := combo.Options()
	if len(options) != 2 || options[1] != "Balanced" {
		t.Fatalf("options = %#v", options)
	}
	options[0] = "Changed"
	if combo.Options()[0] != "Fast" {
		t.Fatal("modifying Options result changed the combo box")
	}
}

func TestComboBoxSelectAndReplaceOptions(t *testing.T) {
	value := "One"
	combo := ComboBox(&value, "One", "Two")
	combo.Select("Two")
	combo.Select("Missing")
	if value != "Two" {
		t.Fatalf("selection = %q, want Two", value)
	}
	combo.SetOptions("Three", "Four")
	if value != "Three" || len(combo.Options()) != 2 {
		t.Fatalf("replacement = %q, %#v", value, combo.Options())
	}
	combo.SetOptions()
	if combo.Selected() != "" {
		t.Fatalf("empty options retained %q", combo.Selected())
	}
}

func TestComboBoxBuilderOptionsAndNilPointer(t *testing.T) {
	combo := ComboBox(nil, "One").Width(36).Focus().OnChange(func(string) {})
	if combo.value == nil || combo.columns != 36 || !combo.focus || combo.onChange == nil {
		t.Fatalf("combo options were not preserved: %#v", combo)
	}
}
