// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import tk "modernc.org/tk9.0"

// ButtonWidget is a clickable button.
type ButtonWidget struct {
	text    string
	onClick func()
	primary bool
}

// Button creates a button. onClick runs when the user activates it.
func Button(text string, onClick func()) *ButtonWidget {
	if onClick == nil {
		onClick = func() {}
	}
	return &ButtonWidget{text: text, onClick: onClick}
}

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
	)
	return mountedWidget{window: button.Window}
}
