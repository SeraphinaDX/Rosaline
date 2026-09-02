// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import tk "modernc.org/tk9.0"

// MouseButton identifies a mouse button without exposing platform details.
type MouseButton uint8

const (
	// MouseNone means that no mouse button is pressed.
	MouseNone MouseButton = iota
	// MouseLeft is the primary mouse button.
	MouseLeft
	// MouseMiddle is the middle mouse button.
	MouseMiddle
	// MouseRight is the secondary mouse button.
	MouseRight
)

// MouseEvent describes mouse input on a Canvas.
// X and Y are measured from the canvas's top-left corner.
type MouseEvent struct {
	X        float64
	Y        float64
	Button   MouseButton
	Dragging bool
	Shift    bool
	Control  bool
	Alt      bool
}

func mouseEvent(event *tk.Event, button MouseButton, dragging bool) MouseEvent {
	return MouseEvent{
		X:        float64(event.X),
		Y:        float64(event.Y),
		Button:   button,
		Dragging: dragging,
		Shift:    event.State&tk.ModifierShift != 0,
		Control:  event.State&tk.ModifierControl != 0,
		Alt:      event.State&tk.ModifierAlt != 0,
	}
}

func heldMouseButton(state tk.Modifier) MouseButton {
	switch {
	case state&tk.ModifierButton1 != 0:
		return MouseLeft
	case state&tk.ModifierButton2 != 0:
		return MouseMiddle
	case state&tk.ModifierButton3 != 0:
		return MouseRight
	default:
		return MouseNone
	}
}
