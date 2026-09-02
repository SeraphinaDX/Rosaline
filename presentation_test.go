// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import "testing"

func TestGridUsesSafeDefaultsAndFiltersChildren(t *testing.T) {
	first := Label("First")
	grid := Grid(0, first, nil, Label("Second"))
	if grid.columns != 1 {
		t.Fatalf("columns = %d, want 1", grid.columns)
	}
	if len(grid.children) != 2 || grid.children[0] != first {
		t.Fatalf("children = %#v", grid.children)
	}
	if grid.gap != 8 {
		t.Fatalf("default gap = %d, want 8", grid.gap)
	}

	grid.Gap(-4).Padding(-5).Expand()
	if grid.gap != 0 || grid.padding != 0 || !grid.expand {
		t.Fatalf("configured grid = gap %d, padding %d, expand %v", grid.gap, grid.padding, grid.expand)
	}
}

func TestStackFiltersChildrenAndExpands(t *testing.T) {
	first := Label("Background")
	stack := Stack(first, nil, Label("Overlay")).Expand()
	if len(stack.children) != 2 || stack.children[0] != first || !stack.expand {
		t.Fatalf("stack = %#v", stack)
	}
}

func TestAlignmentStickyValues(t *testing.T) {
	tests := []struct {
		horizontal Alignment
		vertical   Alignment
		want       string
	}{
		{AlignStart, AlignStart, "nw"},
		{AlignCenter, AlignCenter, ""},
		{AlignEnd, AlignEnd, "se"},
		{AlignStretch, AlignStretch, "nsew"},
		{AlignEnd, AlignCenter, "e"},
		{AlignCenter, AlignStart, "n"},
	}
	for _, test := range tests {
		if got := alignmentSticky(test.horizontal, test.vertical); got != test.want {
			t.Errorf("sticky(%d, %d) = %q, want %q", test.horizontal, test.vertical, got, test.want)
		}
	}
}

func TestStickyPackValues(t *testing.T) {
	tests := []struct {
		sticky string
		anchor string
		fill   string
	}{
		{"", "center", "none"},
		{"nw", "nw", "none"},
		{"se", "se", "none"},
		{"ew", "center", "x"},
		{"ns", "center", "y"},
		{"nsew", "center", "both"},
	}
	for _, test := range tests {
		if anchor := stickyAnchor(test.sticky); anchor != test.anchor {
			t.Errorf("anchor(%q) = %q, want %q", test.sticky, anchor, test.anchor)
		}
		if fill := stickyFill(test.sticky); fill != test.fill {
			t.Errorf("fill(%q) = %q, want %q", test.sticky, fill, test.fill)
		}
	}
}

func TestAlignAndCenterNormalizeValues(t *testing.T) {
	invalid := Alignment(99)
	aligned := Align(nil, invalid, AlignEnd)
	if aligned.content == nil || aligned.horizontal != AlignStart || aligned.vertical != AlignEnd {
		t.Fatalf("aligned = %#v", aligned)
	}

	centered := Center(Label("Hello"))
	if centered.horizontal != AlignCenter || centered.vertical != AlignCenter {
		t.Fatalf("centered = %#v", centered)
	}
}

func TestSpringAndSeparatorOptions(t *testing.T) {
	if Spring() == nil {
		t.Fatal("Spring returned nil")
	}

	separator := Separator()
	if separator.vertical || separator.thickness != 1 {
		t.Fatalf("separator defaults = vertical %v, thickness %d", separator.vertical, separator.thickness)
	}
	separator.Vertical().Thickness(4)
	if !separator.vertical || separator.thickness != 4 {
		t.Fatalf("vertical separator = %#v", separator)
	}
	separator.Horizontal().Thickness(0)
	if separator.vertical || separator.thickness != 1 {
		t.Fatalf("normalized separator = %#v", separator)
	}
}

func TestCardUsesFriendlyDefaults(t *testing.T) {
	card := Card(nil)
	if card.content == nil || card.padding != defaultCardPadding || card.expand {
		t.Fatalf("card defaults = %#v", card)
	}
	card.Padding(-1).Expand()
	if card.padding != 0 || !card.expand {
		t.Fatalf("configured card = %#v", card)
	}
}

func TestSizeBoxesNormalizeDimensions(t *testing.T) {
	preferred := Size(nil, 0, -5)
	if preferred.content == nil || preferred.width != 1 || preferred.height != 1 || preferred.mode != preferredSize {
		t.Fatalf("preferred size = %#v", preferred)
	}
	preferred.Expand()
	if !preferred.expand {
		t.Fatal("preferred size did not expand")
	}

	minimum := MinSize(Label("Content"), 320, 180)
	if minimum.width != 320 || minimum.height != 180 || minimum.mode != minimumSize {
		t.Fatalf("minimum size = %#v", minimum)
	}
}

func TestBoxExpansionUsesOnlyItsMainAxis(t *testing.T) {
	tests := []struct {
		direction direction
		mounted   mountedWidget
		anchor    string
		fill      string
		expand    bool
	}{
		{vertical, mountedWidget{}, "nw", "x", false},
		{vertical, mountedWidget{expandX: true}, "nw", "x", false},
		{vertical, mountedWidget{expandY: true}, "nw", "both", true},
		{horizontal, mountedWidget{}, "nw", "y", false},
		{horizontal, mountedWidget{expandY: true}, "nw", "y", false},
		{horizontal, mountedWidget{expandX: true}, "nw", "both", true},
		{vertical, mountedWidget{expandX: true, expandY: true, sticky: "e", aligned: true}, "e", "none", true},
		{horizontal, mountedWidget{expandX: true, expandY: true, sticky: "ns", aligned: true}, "center", "y", true},
	}
	for _, test := range tests {
		anchor, fill, expand := boxChildLayout(test.direction, test.mounted)
		if anchor != test.anchor || fill != test.fill || expand != test.expand {
			t.Errorf("layout(%d, %#v) = %q, %q, %v; want %q, %q, %v", test.direction, test.mounted, anchor, fill, expand, test.anchor, test.fill, test.expand)
		}
	}
}
