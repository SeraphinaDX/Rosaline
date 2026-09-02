// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import (
	"testing"

	tk "modernc.org/tk9.0"
)

func TestHeldMouseButton(t *testing.T) {
	tests := []struct {
		name  string
		state tk.Modifier
		want  MouseButton
	}{
		{"none", tk.ModifierNone, MouseNone},
		{"left", tk.ModifierButton1, MouseLeft},
		{"middle", tk.ModifierButton2, MouseMiddle},
		{"right", tk.ModifierButton3, MouseRight},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := heldMouseButton(test.state); got != test.want {
				t.Fatalf("heldMouseButton(%v) = %v, want %v", test.state, got, test.want)
			}
		})
	}
}

func TestMouseEventConversion(t *testing.T) {
	event := &tk.Event{
		X:     42,
		Y:     27,
		State: tk.ModifierShift | tk.ModifierControl | tk.ModifierAlt,
	}

	got := mouseEvent(event, MouseLeft, true)
	if got.X != 42 || got.Y != 27 || got.Button != MouseLeft || !got.Dragging {
		t.Fatalf("unexpected mouse event: %#v", got)
	}
	if !got.Shift || !got.Control || !got.Alt {
		t.Fatalf("modifier keys were not preserved: %#v", got)
	}
}

func TestCanvasMouseOptions(t *testing.T) {
	down := 0
	move := 0
	up := 0
	canvas := Canvas(nil).
		OnMouseDown(func(MouseEvent) { down++ }).
		OnMouseMove(func(MouseEvent) { move++ }).
		OnMouseUp(func(MouseEvent) { up++ })

	if canvas.onMouseDown == nil || canvas.onMouseMove == nil || canvas.onMouseUp == nil {
		t.Fatal("canvas mouse handlers were not stored")
	}
	if down != 0 || move != 0 || up != 0 {
		t.Fatal("canvas handlers ran before an event")
	}
}

func TestCanvasRedrawBeforeMountIsSafe(t *testing.T) {
	Canvas(nil).Redraw()
}
