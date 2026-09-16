// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import tk "modernc.org/tk9.0"

// TextBoxWidget is a single-line text input bound to a Go string.
type TextBoxWidget struct {
	value       *string
	placeholder string
	password    bool
	columns     int
	onChange    func(string)
	onSubmit    func(string)
	focus       bool
	enabled     bool
	lastValue   string
	entry       *tk.EntryWidget
	ctx         *mountContext
}

// TextBox creates a single-line input. The box updates value as the user types.
// Pass a pointer with &, as in TextBox(&name).
func TextBox(value *string) *TextBoxWidget {
	if value == nil {
		value = new(string)
	}
	return &TextBoxWidget{value: value, columns: 28, enabled: true, lastValue: *value}
}

// Text returns the current bound text.
func (t *TextBoxWidget) Text() string {
	if t == nil || t.value == nil {
		return ""
	}
	return *t.value
}

// SetText replaces the text. When mounted, OnChange runs for a changed value.
func (t *TextBoxWidget) SetText(text string) {
	if t == nil || t.value == nil {
		return
	}
	old := *t.value
	*t.value = text
	t.lastValue = text
	if t.entry != nil {
		t.entry.Configure(tk.Textvariable(text))
		if old != text && t.onChange != nil {
			t.onChange(text)
		}
		if t.ctx != nil {
			t.ctx.refresh()
		}
	}
}

// SetEnabled enables or disables text editing immediately.
func (t *TextBoxWidget) SetEnabled(enabled bool) {
	if t == nil {
		return
	}
	t.enabled = enabled
	if t.entry != nil {
		t.entry.Configure(tk.State(enabledState(enabled)))
	}
}

// Enabled reports whether the text box currently accepts input.
func (t *TextBoxWidget) Enabled() bool { return t != nil && t.enabled }

// Placeholder shows a hint while the text box is empty.
func (t *TextBoxWidget) Placeholder(text string) *TextBoxWidget {
	t.placeholder = text
	return t
}

// Password hides typed characters. The bound Go string still contains the
// real value so the application can validate or submit it.
func (t *TextBoxWidget) Password() *TextBoxWidget {
	t.password = true
	return t
}

// Width sets the preferred width in text columns.
func (t *TextBoxWidget) Width(columns int) *TextBoxWidget {
	if columns > 0 {
		t.columns = columns
	}
	return t
}

// OnChange runs after the user changes the value.
func (t *TextBoxWidget) OnChange(handler func(string)) *TextBoxWidget {
	t.onChange = handler
	return t
}

// OnSubmit runs when the user presses Enter while the text box has focus.
func (t *TextBoxWidget) OnSubmit(handler func(string)) *TextBoxWidget {
	t.onSubmit = handler
	return t
}

// Focus asks Rosaline to give this text box focus when the window opens.
// If several widgets request focus, the first one wins.
func (t *TextBoxWidget) Focus() *TextBoxWidget {
	t.focus = true
	if t.entry != nil {
		tk.Focus(t.entry.Window)
	}
	return t
}

func (t *TextBoxWidget) mount(ctx *mountContext, parent *tk.Window) mountedWidget {
	options := []tk.Opt{
		tk.Textvariable(*t.value),
		tk.Width(t.columns),
		tk.Background(ctx.theme.Surface.String()),
		tk.Foreground(ctx.theme.Text.String()),
		tk.Insertbackground(ctx.theme.Text.String()),
		tk.Selectbackground(ctx.theme.Primary.String()),
		tk.Relief("solid"),
		tk.Borderwidth(1),
		tk.Highlightthickness(1),
		tk.Highlightbackground(ctx.theme.Border.String()),
		tk.Highlightcolor(ctx.theme.Primary.String()),
		tk.State(enabledState(t.enabled)),
		takeFocusOption(true),
	}
	if t.placeholder != "" {
		options = append(options,
			tk.Placeholder(t.placeholder),
			tk.Placeholderforeground(ctx.theme.Muted.String()),
		)
	}
	if t.password {
		options = append(options, tk.Show("*"))
	}

	entry := parent.Entry(options...)
	t.entry = entry
	t.ctx = ctx
	t.lastValue = *t.value
	ctx.addCleanup(func() {
		t.entry = nil
		t.ctx = nil
	})
	syncValue := func() {
		current := entry.Textvariable()
		*t.value = current
		if current != t.lastValue {
			t.lastValue = current
			if t.onChange != nil {
				t.onChange(current)
			}
		}
	}

	syncAndRefresh := func() {
		syncValue()
		ctx.refresh()
	}
	tk.Bind(entry.Window, "<KeyRelease>", tk.Command(syncAndRefresh))
	tk.Bind(entry.Window, "<FocusOut>", tk.Command(syncAndRefresh))
	tk.Bind(entry.Window, "<Return>", tk.Command(func(event *tk.Event) {
		syncValue()
		if t.onSubmit != nil {
			t.onSubmit(*t.value)
			ctx.refresh()
		}
		event.SetReturnCodeBreak()
	}))

	ctx.flushes = append(ctx.flushes, syncValue)
	ctx.refreshes = append(ctx.refreshes, func() {
		if entry.Textvariable() != *t.value {
			entry.Configure(tk.Textvariable(*t.value))
			t.lastValue = *t.value
		}
	})
	ctx.addFocusable(entry.Window, t.focus)

	return mountedWidget{window: entry.Window}
}
