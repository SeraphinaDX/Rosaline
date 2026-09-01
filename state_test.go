// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import "testing"

func TestState(t *testing.T) {
	state := NewState(2)
	state.Update(func(value int) int { return value * 3 })
	if got := state.Get(); got != 6 {
		t.Fatalf("State.Get() = %d, want 6", got)
	}
	state.Set(9)
	if got := state.Get(); got != 9 {
		t.Fatalf("State.Get() = %d, want 9", got)
	}
}
