// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import "testing"

func TestTabsSelectFirstPageByDefault(t *testing.T) {
	tabs := Tabs(Tab("General", Label("General")), Tab("Drawing", Label("Drawing")))
	index, title, ok := tabs.Selected()
	if !ok || index != 0 || title != "General" {
		t.Fatalf("Selected() = %d, %q, %v; want 0, General, true", index, title, ok)
	}
}

func TestTabsSelectPage(t *testing.T) {
	tabs := Tabs(Tab("General", nil), Tab("Drawing", nil))
	tabs.Select(1)
	index, title, ok := tabs.Selected()
	if !ok || index != 1 || title != "Drawing" {
		t.Fatalf("Selected() = %d, %q, %v; want 1, Drawing, true", index, title, ok)
	}
	tabs.Select(20)
	index, _, _ = tabs.Selected()
	if index != 1 {
		t.Fatalf("invalid selection changed index to %d", index)
	}
}

func TestTabsIgnoreNilPagesAndCopyPageSlice(t *testing.T) {
	general := Tab("General", nil)
	tabs := Tabs(nil, general, nil)
	pages := tabs.Pages()
	if len(pages) != 1 || pages[0] != general {
		t.Fatalf("Pages() = %#v, want General page", pages)
	}
	pages[0] = nil
	if len(tabs.Pages()) != 1 || tabs.Pages()[0] != general {
		t.Fatal("modifying Pages result changed tabs")
	}
}

func TestEmptyTabsHaveNoSelection(t *testing.T) {
	index, title, ok := Tabs().Selected()
	if ok || index != -1 || title != "" {
		t.Fatalf("Selected() = %d, %q, %v; want -1, empty, false", index, title, ok)
	}
}

func TestTabUsesFriendlyDefaultTitle(t *testing.T) {
	page := Tab("", nil)
	if page.Title() != "Tab" {
		t.Fatalf("Title() = %q, want Tab", page.Title())
	}
}

func TestTabsBuilderOptions(t *testing.T) {
	tabs := Tabs(Tab("General", nil)).Expand().OnChange(func(int, string) {})
	if !tabs.expand || tabs.onChange == nil {
		t.Fatal("tabs builder options were not retained")
	}
}
