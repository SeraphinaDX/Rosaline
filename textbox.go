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
}

// TextBox creates a single-line input. The box updates value as the user types.
// Pass a pointer with &, as in TextBox(&name).
func TextBox(value *string) *TextBoxWidget {
	if value == nil {
		value = new(string)
	}
	return &TextBoxWidget{value: value, columns: 28}
}

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
	lastValue := *t.value
	syncValue := func() {
		current := entry.Textvariable()
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
			lastValue = *t.value
		}
	})
	ctx.addFocusable(entry.Window, t.focus)

	return mountedWidget{window: entry.Window}
}
