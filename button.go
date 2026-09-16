// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import tk "modernc.org/tk9.0"

// ButtonWidget is a clickable button.
type ButtonWidget struct {
	text    string
	onClick func()
	primary bool
	enabled bool
	button  *tk.ButtonWidget
}

// Button creates a button. onClick runs when the user activates it.
func Button(text string, onClick func()) *ButtonWidget {
	if onClick == nil {
		onClick = func() {}
	}
	return &ButtonWidget{text: text, onClick: onClick, enabled: true}
}

// OnClick replaces the function run when the button is activated.
func (b *ButtonWidget) OnClick(handler func()) *ButtonWidget {
	if b == nil {
		return b
	}
	if handler == nil {
		handler = func() {}
	}
	b.onClick = handler
	return b
}

// Text returns the button's current label.
func (b *ButtonWidget) Text() string {
	if b == nil {
		return ""
	}
	return b.text
}

// SetText replaces the button label immediately.
func (b *ButtonWidget) SetText(text string) {
	if b == nil {
		return
	}
	b.text = text
	if b.button != nil {
		b.button.Configure(tk.Txt(text))
	}
}

// SetEnabled enables or disables the button immediately.
func (b *ButtonWidget) SetEnabled(enabled bool) {
	if b == nil {
		return
	}
	b.enabled = enabled
	if b.button != nil {
		b.button.Configure(tk.State(enabledState(enabled)))
	}
}

// Enabled reports whether the button can currently be activated.
func (b *ButtonWidget) Enabled() bool { return b != nil && b.enabled }

// Primary gives a button the theme's primary color.
func (b *ButtonWidget) Primary() *ButtonWidget {
	b.primary = true
	return b
}

func (b *ButtonWidget) mount(ctx *mountContext, parent *tk.Window) mountedWidget {
	background := ctx.theme.Surface
	foreground := ctx.theme.Text
	active := ctx.theme.Border
	if b.primary {
		background = ctx.theme.Primary
		foreground = White
		active = ctx.theme.Primary
	}
	button := parent.Button(
		tk.Txt(b.text),
		tk.Command(func() {
			ctx.flush()
			b.onClick()
			ctx.refresh()
		}),
		tk.Background(background.String()),
		tk.Foreground(foreground.String()),
		tk.Activebackground(active.String()),
		tk.Activeforeground(foreground.String()),
		tk.Padx(12),
		tk.Pady(7),
		tk.Relief("flat"),
		tk.Borderwidth(0),
		tk.State(enabledState(b.enabled)),
		takeFocusOption(true),
	)
	b.button = button
	ctx.addCleanup(func() { b.button = nil })
	ctx.addFocusable(button.Window, false)
	return mountedWidget{window: button.Window}
}

func enabledState(enabled bool) string {
	if enabled {
		return "normal"
	}
	return "disabled"
}
