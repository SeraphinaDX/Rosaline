// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import "testing"

func TestRadioGroupNormalizesChoicesAndSelection(t *testing.T) {
	value := "missing"
	group := RadioGroup(&value,
		Choice("", "low"),
		Choice("High", "high"),
		Choice("Duplicate", "high"),
	)
	if value != "low" || group.Selected() != "low" {
		t.Fatalf("selection = %q, want low", value)
	}
	choices := group.Choices()
	if len(choices) != 2 || choices[0].Label() != "low" || choices[1].Value() != "high" {
		t.Fatalf("choices = %#v", choices)
	}
	choices[0] = Choice("Changed", "changed")
	if group.Choices()[0].Value() != "low" {
		t.Fatal("modifying Choices result changed the group")
	}
}

func TestRadioGroupSelectAndReplaceChoices(t *testing.T) {
	value := "one"
	group := RadioGroup(&value, Choice("One", "one"), Choice("Two", "two"))
	group.Select("two")
	group.Select("missing")
	if value != "two" {
		t.Fatalf("selection = %q, want two", value)
	}
	group.SetChoices(Choice("Three", "three"), Choice("Four", "four"))
	if value != "three" || len(group.Choices()) != 2 {
		t.Fatalf("replacement = %q, %#v", value, group.Choices())
	}
	group.SetChoices()
	if group.Selected() != "" {
		t.Fatalf("empty choices retained %q", group.Selected())
	}
}

func TestRadioGroupOptionsAndNilPointer(t *testing.T) {
	group := RadioGroup(nil, Choice("One", "one")).
		Horizontal().Vertical().Horizontal().Focus().OnChange(func(string) {})
	if group.value == nil || !group.horizontal || !group.focus || group.onChange == nil {
		t.Fatalf("group options were not preserved: %#v", group)
	}
}

func TestChoiceUsesFriendlyLabel(t *testing.T) {
	if got := Choice("", "").Label(); got != "Choice" {
		t.Fatalf("empty choice label = %q, want Choice", got)
	}
}
