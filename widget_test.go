// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import (
	"testing"

	tk "modernc.org/tk9.0"
)

func TestRelativeFocusSkipsHiddenWidgets(t *testing.T) {
	first := &tk.Window{}
	hidden := &tk.Window{}
	last := &tk.Window{}
	ctx := mountContext{focusables: []focusableWidget{
		{window: first},
		{window: hidden, visible: func() bool { return false }},
		{window: last},
	}}

	if got := ctx.relativeFocus(0, 1); got != last {
		t.Fatalf("relativeFocus forward = %p, want %p", got, last)
	}
	if got := ctx.relativeFocus(2, -1); got != first {
		t.Fatalf("relativeFocus backward = %p, want %p", got, first)
	}
}

func TestFocusConditionTracksChangingVisibility(t *testing.T) {
	ctx := mountContext{}
	active := false
	window := &tk.Window{}
	ctx.withFocusCondition(func() bool { return active }, func() {
		ctx.addFocusable(window, true)
	})

	if ctx.initialFocus != nil {
		t.Fatal("hidden widget unexpectedly became the initial focus")
	}
	if ctx.focusables[0].isVisible() {
		t.Fatal("widget should start hidden")
	}
	active = true
	if !ctx.focusables[0].isVisible() {
		t.Fatal("widget should become visible when its condition changes")
	}
}
