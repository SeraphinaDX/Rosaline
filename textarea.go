// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import tk "modernc.org/tk9.0"

// TextAreaWidget is a multiline text input bound to a Go string.
type TextAreaWidget struct {
	value    *string
	columns  int
	lines    int
	onChange func(string)
	focus    bool
}

// TextArea creates a multiline input. The area updates value as the user types.
// Pass a pointer with &, as in TextArea(&notes).
func TextArea(value *string) *TextAreaWidget {
	if value == nil {
		value = new(string)
	}
	return &TextAreaWidget{value: value, columns: 40, lines: 6}
}

// Size sets the preferred width in text columns and height in text lines.
func (t *TextAreaWidget) Size(columns, lines int) *TextAreaWidget {
	if columns > 0 {
		t.columns = columns
	}
	if lines > 0 {
		t.lines = lines
	}
	return t
}

// OnChange runs after the user changes the value.
func (t *TextAreaWidget) OnChange(handler func(string)) *TextAreaWidget {
	t.onChange = handler
	return t
}

// Focus asks Rosaline to give this text area focus when the window opens.
// If several widgets request focus, the first one wins.
func (t *TextAreaWidget) Focus() *TextAreaWidget {
	t.focus = true
	return t
}

func (t *TextAreaWidget) mount(ctx *mountContext, parent *tk.Window) mountedWidget {
	area := parent.Text(
		tk.Width(t.columns),
		tk.Height(t.lines),
		tk.Wrap("word"),
		tk.Undo(true),
		tk.Background(ctx.theme.Surface.String()),
		tk.Foreground(ctx.theme.Text.String()),
		tk.Insertbackground(ctx.theme.Text.String()),
		tk.Selectbackground(ctx.theme.Primary.String()),
		tk.Relief("solid"),
		tk.Borderwidth(1),
		tk.Highlightthickness(1),
		tk.Highlightbackground(ctx.theme.Border.String()),
		tk.Highlightcolor(ctx.theme.Primary.String()),
		takeFocusOption(true),
	)
	if *t.value != "" {
		area.Insert("1.0", *t.value)
	}

	lastValue := *t.value
	syncValue := func() {
		current := area.Text()
		*t.value = current
		if current != lastValue {
			lastValue = current
			if t.onChange != nil {
				t.onChange(current)
			}
		}
	}

	syncAndRefresh := func() {
		syncValue()
		ctx.refresh()
	}
	tk.Bind(area.Window, "<KeyRelease>", tk.Command(syncAndRefresh))
	tk.Bind(area.Window, "<FocusOut>", tk.Command(syncAndRefresh))

	ctx.flushes = append(ctx.flushes, syncValue)
	ctx.refreshes = append(ctx.refreshes, func() {
		if area.Text() != *t.value {
			area.Delete("1.0", "end")
			if *t.value != "" {
				area.Insert("1.0", *t.value)
			}
			lastValue = *t.value
		}
	})
	ctx.addFocusable(area.Window, t.focus)

	return mountedWidget{window: area.Window}
}
