// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import "testing"

func TestScrollOptions(t *testing.T) {
	scroll := Scroll(nil).Size(720, 480).Expand()
	if scroll.content == nil || scroll.width != 720 || scroll.height != 480 || !scroll.expand {
		t.Fatalf("unexpected scroll options: %#v", scroll)
	}
}
