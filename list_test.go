// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import "testing"

func TestListSelectsFirstItemByDefault(t *testing.T) {
	list := List("Rose", "Violet", "Peony")
	index, value, ok := list.Selected()
	if !ok || index != 0 || value != "Rose" {
		t.Fatalf("Selected() = %d, %q, %v; want 0, Rose, true", index, value, ok)
	}
}

func TestEmptyListHasNoSelection(t *testing.T) {
	index, value, ok := List().Selected()
	if ok || index != -1 || value != "" {
		t.Fatalf("Selected() = %d, %q, %v; want -1, empty, false", index, value, ok)
	}
}

func TestListSelectionAndItemsAreIndependentCopies(t *testing.T) {
	items := []string{"Rose", "Violet"}
	list := List(items...)
	items[0] = "Changed"
	list.Select(1)

	index, value, ok := list.Selected()
	if !ok || index != 1 || value != "Violet" {
		t.Fatalf("Selected() = %d, %q, %v", index, value, ok)
	}
	copy := list.Items()
	copy[1] = "Changed"
	if _, value, _ := list.Selected(); value != "Violet" {
		t.Fatalf("modifying Items result changed selection to %q", value)
	}
}

func TestListSetItemsKeepsSelectionInRange(t *testing.T) {
	list := List("A", "B", "C")
	list.Select(2)
	list.SetItems("One", "Two")

	index, value, ok := list.Selected()
	if !ok || index != 1 || value != "Two" {
		t.Fatalf("Selected() = %d, %q, %v; want 1, Two, true", index, value, ok)
	}

	list.SetItems()
	if _, _, ok := list.Selected(); ok {
		t.Fatal("empty replacement should clear selection")
	}

	list.SetItems("Again")
	index, value, ok = list.Selected()
	if !ok || index != 0 || value != "Again" {
		t.Fatalf("Selected() = %d, %q, %v; want 0, Again, true", index, value, ok)
	}
}

func TestInvalidListSelectionClearsSelection(t *testing.T) {
	list := List("A", "B")
	list.Select(9)
	if _, _, ok := list.Selected(); ok {
		t.Fatal("invalid index should clear selection")
	}
}

func TestListBuilderOptions(t *testing.T) {
	list := List("A").Size(40, 12).Expand().OnSelect(func(int, string) {}).OnActivate(func(int, string) {})
	if list.columns != 40 || list.rows != 12 || !list.expand {
		t.Fatalf("list options = columns %d, rows %d, expand %v", list.columns, list.rows, list.expand)
	}
	if list.onSelect == nil || list.onActivate == nil {
		t.Fatal("list callbacks were not retained")
	}
}
