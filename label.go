// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import tk "modernc.org/tk9.0"

// LabelWidget displays text.
type LabelWidget struct {
	text  func() string
	color *Color
}

// Label creates a label with fixed text.
func Label(text string) *LabelWidget {
	return &LabelWidget{text: func() string { return text }}
}

// LabelFunc creates a label whose text is recalculated after Rosaline events.
// It is useful for counters and other small pieces of changing text.
func LabelFunc(text func() string) *LabelWidget {
	if text == nil {
		text = func() string { return "" }
	}
	return &LabelWidget{text: text}
}

// Color sets this label's text color.
func (l *LabelWidget) Color(color Color) *LabelWidget {
	l.color = &color
	return l
}

func (l *LabelWidget) mount(ctx *mountContext, parent *tk.Window) mountedWidget {
	color := ctx.theme.Text
	if l.color != nil {
		color = *l.color
	}
	label := parent.Label(
		tk.Txt(l.text()),
		tk.Foreground(color.String()),
		tk.Background(ctx.theme.Background.String()),
		tk.Anchor("w"),
		tk.Borderwidth(0),
	)
	ctx.refreshes = append(ctx.refreshes, func() {
		label.Configure(tk.Txt(l.text()))
	})
	return mountedWidget{window: label.Window}
}
