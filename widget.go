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

type focusableWidget struct {
	window  *tk.Window
	visible func() bool
}

type mountContext struct {
	theme          Theme
	session        *applicationSession
	owner          *Window
	flushes        []func()
	refreshes      []func()
	cleanups       []func()
	focusables     []focusableWidget
	focusCondition func() bool
	initialFocus   *tk.Window
	closed         bool
}

// takeFocusOption must use Tcl's literal 0 and 1 values. The -takefocus
// option treats other values, including "true" and "false", as Tcl scripts.
func takeFocusOption(enabled bool) tk.Opt {
	if enabled {
		return tk.Takefocus(1)
	}
	return tk.Takefocus(0)
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
	if c.session != nil {
		c.session.refreshAll()
		return
	}
	c.refreshLocal()
}

func (c *mountContext) refreshLocal() {
	if c == nil || c.closed {
		return
	}
	for _, refresh := range c.refreshes {
		refresh()
	}
}

func (c *mountContext) addCleanup(cleanup func()) {
	if c != nil && cleanup != nil {
		c.cleanups = append(c.cleanups, cleanup)
	}
}

func (c *mountContext) release() {
	if c == nil || c.closed {
		return
	}
	c.flush()
	for index := len(c.cleanups) - 1; index >= 0; index-- {
		c.cleanups[index]()
	}
	c.abandon()
}

func (c *mountContext) abandon() {
	if c == nil {
		return
	}
	c.closed = true
	c.flushes = nil
	c.refreshes = nil
	c.cleanups = nil
	c.focusables = nil
	c.initialFocus = nil
}

func (c *mountContext) addFocusable(window *tk.Window, initial bool) int {
	focusable := focusableWidget{window: window, visible: c.focusCondition}
	c.focusables = append(c.focusables, focusable)
	if initial && c.initialFocus == nil && focusable.isVisible() {
		c.initialFocus = window
	}
	return len(c.focusables) - 1
}

func (f focusableWidget) isVisible() bool {
	return f.window != nil && (f.visible == nil || f.visible())
}

func (c *mountContext) updateFocusable(index int, window *tk.Window) {
	if c == nil || index < 0 || index >= len(c.focusables) {
		return
	}
	c.focusables[index].window = window
}

func (c *mountContext) withFocusCondition(condition func() bool, mount func()) {
	previous := c.focusCondition
	if previous == nil {
		c.focusCondition = condition
	} else {
		c.focusCondition = func() bool {
			return previous() && condition()
		}
	}
	mount()
	c.focusCondition = previous
}

func (c *mountContext) relativeFocus(from, step int) *tk.Window {
	count := len(c.focusables)
	if count == 0 {
		return nil
	}
	for offset := 1; offset <= count; offset++ {
		index := (from + step*offset) % count
		if index < 0 {
			index += count
		}
		if c.focusables[index].isVisible() {
			return c.focusables[index].window
		}
	}
	return nil
}

func (c *mountContext) installFocusTraversal() {
	if len(c.focusables) < 2 {
		return
	}

	for i, focusable := range c.focusables {
		c.bindFocusTraversal(focusable.window, i)
	}
}

func (c *mountContext) bindFocusTraversal(window *tk.Window, index int) {
	if c == nil || window == nil {
		return
	}
	forward := tk.Command(func(event *tk.Event) {
		c.flush()
		if next := c.relativeFocus(index, 1); next != nil {
			tk.Focus(next)
		}
		event.SetReturnCodeBreak()
	})
	backward := tk.Command(func(event *tk.Event) {
		c.flush()
		if previous := c.relativeFocus(index, -1); previous != nil {
			tk.Focus(previous)
		}
		event.SetReturnCodeBreak()
	})

	tk.Bind(window, "<Tab>", forward)
	tk.Bind(window, "<Shift-Key-Tab>", backward)
	tk.Bind(window, "<ISO_Left_Tab>", backward)
}
