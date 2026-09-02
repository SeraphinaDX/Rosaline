// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import (
	"fmt"
	"testing"

	tk "modernc.org/tk9.0"
)

func TestTakeFocusOptionsUseTclNumericBooleans(t *testing.T) {
	if got := fmt.Sprint(takeFocusOption(true)); got != "-takefocus 1" {
		t.Fatalf("enabled takefocus option = %q, want numeric 1", got)
	}
	if got := fmt.Sprint(takeFocusOption(false)); got != "-takefocus 0" {
		t.Fatalf("disabled takefocus option = %q, want numeric 0", got)
	}
}

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

func TestFocusableSlotCanFollowAReplacedNativeControl(t *testing.T) {
	first := &tk.Window{}
	replacement := &tk.Window{}
	last := &tk.Window{}
	ctx := mountContext{}
	ctx.addFocusable(first, false)
	slot := ctx.addFocusable(nil, false)
	ctx.addFocusable(last, false)

	if got := ctx.relativeFocus(0, 1); got != last {
		t.Fatalf("empty slot returned %p, want last widget %p", got, last)
	}
	ctx.updateFocusable(slot, replacement)
	if got := ctx.relativeFocus(0, 1); got != replacement {
		t.Fatalf("updated slot returned %p, want replacement %p", got, replacement)
	}

	ctx.updateFocusable(-1, first)
	ctx.updateFocusable(99, first)
}

func TestMountContextReleasesControlsInReverseOrder(t *testing.T) {
	ctx := mountContext{}
	var order []string
	ctx.flushes = append(ctx.flushes, func() { order = append(order, "flush") })
	ctx.addCleanup(func() { order = append(order, "first") })
	ctx.addCleanup(func() { order = append(order, "second") })

	ctx.release()
	ctx.release()
	got := fmt.Sprint(order)
	if got != "[flush second first]" {
		t.Fatalf("release order = %s, want [flush second first]", got)
	}
	if !ctx.closed || ctx.flushes != nil || ctx.refreshes != nil || ctx.cleanups != nil {
		t.Fatal("released context retained active callbacks")
	}
}
