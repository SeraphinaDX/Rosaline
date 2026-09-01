// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import tk "modernc.org/tk9.0"

// Widget is anything Rosaline can place in a window or layout.
// Applications normally use constructor functions such as Label, Button,
// Column, Row, and Canvas rather than implementing Widget themselves.
type Widget interface {
	mount(*mountContext, *tk.Window) mountedWidget
}

type mountedWidget struct {
	window  *tk.Window
	expandX bool
	expandY bool
}

type mountContext struct {
	theme     Theme
	refreshes []func()
}

func (c *mountContext) refresh() {
	for _, refresh := range c.refreshes {
		refresh()
	}
}
