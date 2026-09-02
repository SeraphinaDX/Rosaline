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
	theme        Theme
	flushes      []func()
	refreshes    []func()
	focusables   []*tk.Window
	initialFocus *tk.Window
	closed       bool
}

func (c *mountContext) flush() {
	if c.closed {
		return
	}
	for _, flush := range c.flushes {
		flush()
	}
}

func (c *mountContext) refresh() {
	if c.closed {
		return
	}
	for _, refresh := range c.refreshes {
		refresh()
	}
}

func (c *mountContext) addFocusable(window *tk.Window, initial bool) {
	c.focusables = append(c.focusables, window)
	if initial && c.initialFocus == nil {
		c.initialFocus = window
	}
}

func (c *mountContext) installFocusTraversal() {
	if len(c.focusables) < 2 {
		return
	}

	for i, window := range c.focusables {
		next := c.focusables[(i+1)%len(c.focusables)]
		previous := c.focusables[(i-1+len(c.focusables))%len(c.focusables)]

		forward := tk.Command(func(event *tk.Event) {
			c.flush()
			tk.Focus(next)
			event.SetReturnCodeBreak()
		})
		backward := tk.Command(func(event *tk.Event) {
			c.flush()
			tk.Focus(previous)
			event.SetReturnCodeBreak()
		})

		tk.Bind(window, "<Tab>", forward)
		tk.Bind(window, "<Shift-Key-Tab>", backward)
		tk.Bind(window, "<ISO_Left_Tab>", tk.Command(func(event *tk.Event) {
			c.flush()
			tk.Focus(previous)
			event.SetReturnCodeBreak()
		}))
	}
}
