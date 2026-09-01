// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import tk "modernc.org/tk9.0"

type direction uint8

const (
	vertical direction = iota
	horizontal
)

// Box arranges child widgets in a row or column.
type Box struct {
	direction direction
	children  []Widget
	gap       int
	padding   int
	expand    bool
}

// Column arranges widgets from top to bottom.
func Column(children ...Widget) *Box {
	return &Box{direction: vertical, children: cleanWidgets(children), gap: 8}
}

// Row arranges widgets from left to right.
func Row(children ...Widget) *Box {
	return &Box{direction: horizontal, children: cleanWidgets(children), gap: 8}
}

func cleanWidgets(widgets []Widget) []Widget {
	result := make([]Widget, 0, len(widgets))
	for _, widget := range widgets {
		if widget != nil {
			result = append(result, widget)
		}
	}
	return result
}

// Gap sets the space between children in pixels.
func (b *Box) Gap(pixels int) *Box {
	if pixels < 0 {
		pixels = 0
	}
	b.gap = pixels
	return b
}

// Padding sets the space inside the layout in pixels.
func (b *Box) Padding(pixels int) *Box {
	if pixels < 0 {
		pixels = 0
	}
	b.padding = pixels
	return b
}

// Expand asks the layout to use available window space.
func (b *Box) Expand() *Box {
	b.expand = true
	return b
}

func (b *Box) mount(ctx *mountContext, parent *tk.Window) mountedWidget {
	frame := parent.Frame(tk.Background(ctx.theme.Background.String()), tk.Borderwidth(0))
	inner := frame.Frame(tk.Background(ctx.theme.Background.String()), tk.Borderwidth(0))
	tk.Pack(
		inner,
		tk.Fill("both"),
		tk.Expand(true),
		tk.Padx(b.padding),
		tk.Pady(b.padding),
	)
	side := "top"
	fill := "x"
	if b.direction == horizontal {
		side = "left"
		fill = "y"
	}
	halfGap := b.gap / 2
	for _, child := range b.children {
		mounted := child.mount(ctx, inner.Window)
		opts := []tk.Opt{
			tk.Side(side),
			tk.Anchor("nw"),
			tk.Padx(halfGap),
			tk.Pady(halfGap),
		}
		if mounted.expandX || mounted.expandY {
			opts = append(opts, tk.Fill("both"), tk.Expand(true))
		} else {
			opts = append(opts, tk.Fill(fill))
		}
		tk.Pack(append([]tk.Opt{mounted.window}, opts...)...)
	}
	return mountedWidget{window: frame.Window, expandX: b.expand, expandY: b.expand}
}

type spacerWidget struct {
	width, height int
}

// Spacer inserts a fixed amount of empty space.
func Spacer(width, height int) Widget { return &spacerWidget{width: width, height: height} }

func (s *spacerWidget) mount(ctx *mountContext, parent *tk.Window) mountedWidget {
	frame := parent.Frame(
		tk.Width(max(0, s.width)),
		tk.Height(max(0, s.height)),
		tk.Background(ctx.theme.Background.String()),
		tk.Borderwidth(0),
	)
	return mountedWidget{window: frame.Window}
}
