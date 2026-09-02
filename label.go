// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import tk "modernc.org/tk9.0"

// LabelWidget displays text.
type LabelWidget struct {
	text      func() string
	color     *Color
	fontSize  int
	bold      bool
	alignment Alignment
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

// FontSize sets the label text size in pixels. Non-positive values use the
// platform's normal interface size.
func (l *LabelWidget) FontSize(pixels int) *LabelWidget {
	if l != nil {
		l.fontSize = max(0, pixels)
	}
	return l
}

// Bold gives the label bold text.
func (l *LabelWidget) Bold() *LabelWidget {
	if l != nil {
		l.bold = true
	}
	return l
}

// TextAlign aligns text inside the label horizontally.
func (l *LabelWidget) TextAlign(alignment Alignment) *LabelWidget {
	if l != nil {
		alignment = normalizeAlignment(alignment)
		if alignment == AlignStretch {
			alignment = AlignStart
		}
		l.alignment = alignment
	}
	return l
}

func labelAnchor(alignment Alignment) string {
	switch normalizeAlignment(alignment) {
	case AlignCenter:
		return "center"
	case AlignEnd:
		return "e"
	default:
		return "w"
	}
}

func (l *LabelWidget) mount(ctx *mountContext, parent *tk.Window) mountedWidget {
	color := ctx.theme.Text
	if l.color != nil {
		color = *l.color
	}
	options := []tk.Opt{
		tk.Txt(l.text()),
		tk.Foreground(color.String()),
		tk.Background(ctx.theme.Background.String()),
		tk.Anchor(labelAnchor(l.alignment)),
		tk.Borderwidth(0),
	}
	if l.fontSize > 0 || l.bold {
		fontOptions := make([]tk.Opt, 0, 2)
		if l.fontSize > 0 {
			fontOptions = append(fontOptions, tk.Size(-l.fontSize))
		}
		if l.bold {
			fontOptions = append(fontOptions, tk.Weight("bold"))
		}
		font := tk.NewFont(fontOptions...)
		options = append(options, tk.Font(font))
		ctx.addCleanup(font.Delete)
	}
	label := parent.Label(options...)
	ctx.refreshes = append(ctx.refreshes, func() {
		label.Configure(tk.Txt(l.text()))
	})
	return mountedWidget{window: label.Window}
}
