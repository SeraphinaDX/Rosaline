// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import tk "modernc.org/tk9.0"

const defaultCardPadding = 14

// GridLayout arranges widgets in automatically filled rows and columns.
type GridLayout struct {
	columns  int
	children []Widget
	gap      int
	padding  int
	expand   bool
}

// Grid creates a layout with the requested number of columns. Children fill
// each row from left to right. Invalid column counts safely use one column.
func Grid(columns int, children ...Widget) *GridLayout {
	return &GridLayout{
		columns:  max(1, columns),
		children: cleanWidgets(children),
		gap:      8,
	}
}

// Gap sets the space between grid cells in pixels.
func (g *GridLayout) Gap(pixels int) *GridLayout {
	if g != nil {
		g.gap = max(0, pixels)
	}
	return g
}

// Padding sets the space around the inside of the grid in pixels.
func (g *GridLayout) Padding(pixels int) *GridLayout {
	if g != nil {
		g.padding = max(0, pixels)
	}
	return g
}

// Expand asks the grid and its equal-sized cells to use available space.
func (g *GridLayout) Expand() *GridLayout {
	if g != nil {
		g.expand = true
	}
	return g
}

func (g *GridLayout) mount(ctx *mountContext, parent *tk.Window) mountedWidget {
	frame := parent.Frame(tk.Background(ctx.theme.Background.String()), tk.Borderwidth(0))
	inner := frame.Frame(tk.Background(ctx.theme.Background.String()), tk.Borderwidth(0))
	tk.Pack(inner, tk.Fill("both"), tk.Expand(true), tk.Padx(g.padding), tk.Pady(g.padding))

	rows := (len(g.children) + g.columns - 1) / g.columns
	for column := range g.columns {
		tk.GridColumnConfigure(inner, column, tk.Weight(1), tk.Uniform("rosaline-grid-columns"))
	}
	for row := range rows {
		if g.expand {
			tk.GridRowConfigure(inner, row, tk.Weight(1), tk.Uniform("rosaline-grid-rows"))
		}
	}

	halfGap := g.gap / 2
	for index, child := range g.children {
		mounted := child.mount(ctx, inner.Window)
		sticky := mounted.sticky
		if !mounted.aligned {
			sticky = "nsew"
		}
		tk.Grid(
			mounted.window,
			tk.Row(index/g.columns),
			tk.Column(index%g.columns),
			tk.Padx(halfGap),
			tk.Pady(halfGap),
			tk.Sticky(sticky),
		)
	}
	return mountedWidget{window: frame.Window, expandX: g.expand, expandY: g.expand}
}

// StackLayout layers widgets in the same space. Later children appear above
// earlier children. Align or Center can position an overlay without stretching
// its visible content.
type StackLayout struct {
	children []Widget
	expand   bool
}

// Stack layers ordinary widgets in one shared area.
func Stack(children ...Widget) *StackLayout {
	return &StackLayout{children: cleanWidgets(children)}
}

// Expand asks the stack to use available space.
func (s *StackLayout) Expand() *StackLayout {
	if s != nil {
		s.expand = true
	}
	return s
}

func (s *StackLayout) mount(ctx *mountContext, parent *tk.Window) mountedWidget {
	frame := parent.Frame(tk.Background(ctx.theme.Background.String()), tk.Borderwidth(0))
	tk.GridColumnConfigure(frame, 0, tk.Weight(1))
	tk.GridRowConfigure(frame, 0, tk.Weight(1))
	for _, child := range s.children {
		mounted := child.mount(ctx, frame.Window)
		sticky := mounted.sticky
		if !mounted.aligned {
			sticky = "nsew"
		}
		tk.Grid(mounted.window, tk.Row(0), tk.Column(0), tk.Sticky(sticky))
		mounted.window.Raise(nil)
	}
	return mountedWidget{window: frame.Window, expandX: s.expand, expandY: s.expand}
}

// Alignment controls where Align places content on one axis.
type Alignment uint8

const (
	// AlignStart places content at the left or top edge.
	AlignStart Alignment = iota
	// AlignCenter places content in the center of an axis.
	AlignCenter
	// AlignEnd places content at the right or bottom edge.
	AlignEnd
	// AlignStretch stretches content across an axis.
	AlignStretch
)

// AlignmentBox positions content within all the space its parent gives it.
type AlignmentBox struct {
	content    Widget
	horizontal Alignment
	vertical   Alignment
}

// Align positions content independently on the horizontal and vertical axes.
func Align(content Widget, horizontal, vertical Alignment) *AlignmentBox {
	return &AlignmentBox{
		content:    safeWidget(content),
		horizontal: normalizeAlignment(horizontal),
		vertical:   normalizeAlignment(vertical),
	}
}

// Center centers content horizontally and vertically.
func Center(content Widget) *AlignmentBox {
	return Align(content, AlignCenter, AlignCenter)
}

func normalizeAlignment(alignment Alignment) Alignment {
	if alignment > AlignStretch {
		return AlignStart
	}
	return alignment
}

func alignmentSticky(horizontal, vertical Alignment) string {
	sticky := ""
	switch normalizeAlignment(vertical) {
	case AlignStart:
		sticky += "n"
	case AlignEnd:
		sticky += "s"
	case AlignStretch:
		sticky += "ns"
	}
	switch normalizeAlignment(horizontal) {
	case AlignStart:
		sticky += "w"
	case AlignEnd:
		sticky += "e"
	case AlignStretch:
		sticky += "ew"
	}
	return sticky
}

func stickyAnchor(sticky string) string {
	left := containsDirection(sticky, 'w') && !containsDirection(sticky, 'e')
	right := containsDirection(sticky, 'e') && !containsDirection(sticky, 'w')
	top := containsDirection(sticky, 'n') && !containsDirection(sticky, 's')
	bottom := containsDirection(sticky, 's') && !containsDirection(sticky, 'n')
	switch {
	case top && left:
		return "nw"
	case top && right:
		return "ne"
	case bottom && left:
		return "sw"
	case bottom && right:
		return "se"
	case top:
		return "n"
	case bottom:
		return "s"
	case left:
		return "w"
	case right:
		return "e"
	default:
		return "center"
	}
}

func stickyFill(sticky string) string {
	horizontal := containsDirection(sticky, 'w') && containsDirection(sticky, 'e')
	vertical := containsDirection(sticky, 'n') && containsDirection(sticky, 's')
	switch {
	case horizontal && vertical:
		return "both"
	case horizontal:
		return "x"
	case vertical:
		return "y"
	default:
		return "none"
	}
}

func containsDirection(sticky string, direction byte) bool {
	for index := range len(sticky) {
		if sticky[index] == direction {
			return true
		}
	}
	return false
}

func (a *AlignmentBox) mount(ctx *mountContext, parent *tk.Window) mountedWidget {
	mounted := a.content.mount(ctx, parent)
	mounted.expandX = true
	mounted.expandY = true
	mounted.sticky = alignmentSticky(a.horizontal, a.vertical)
	mounted.aligned = true
	return mounted
}

// SpringWidget is flexible empty space that absorbs extra room in a Row or
// Column. Use Spacer when the empty space should have a fixed size.
type SpringWidget struct{}

// Spring creates flexible empty space.
func Spring() *SpringWidget { return &SpringWidget{} }

func (s *SpringWidget) mount(ctx *mountContext, parent *tk.Window) mountedWidget {
	frame := parent.Frame(tk.Background(ctx.theme.Background.String()), tk.Borderwidth(0))
	return mountedWidget{window: frame.Window, expandX: true, expandY: true}
}

// SeparatorWidget is a thin themed dividing line.
type SeparatorWidget struct {
	vertical  bool
	thickness int
}

// Separator creates a one-pixel horizontal dividing line.
func Separator() *SeparatorWidget { return &SeparatorWidget{thickness: 1} }

// Vertical changes the separator to a vertical line.
func (s *SeparatorWidget) Vertical() *SeparatorWidget {
	if s != nil {
		s.vertical = true
	}
	return s
}

// Horizontal changes the separator to a horizontal line. This is the default.
func (s *SeparatorWidget) Horizontal() *SeparatorWidget {
	if s != nil {
		s.vertical = false
	}
	return s
}

// Thickness changes the line thickness in pixels. Invalid values use one.
func (s *SeparatorWidget) Thickness(pixels int) *SeparatorWidget {
	if s != nil {
		s.thickness = max(1, pixels)
	}
	return s
}

func (s *SeparatorWidget) mount(ctx *mountContext, parent *tk.Window) mountedWidget {
	options := []tk.Opt{
		tk.Background(ctx.theme.Border.String()),
		tk.Borderwidth(0),
	}
	if s.vertical {
		options = append(options, tk.Width(s.thickness))
		line := parent.Frame(options...)
		return mountedWidget{window: line.Window, expandY: true, sticky: "ns", aligned: true}
	}
	options = append(options, tk.Height(s.thickness))
	line := parent.Frame(options...)
	return mountedWidget{window: line.Window, expandX: true, sticky: "ew", aligned: true}
}

// CardWidget presents one widget on the theme's surface color with a border
// and comfortable padding.
type CardWidget struct {
	content Widget
	padding int
	expand  bool
}

// Card wraps content in a themed surface and border.
func Card(content Widget) *CardWidget {
	return &CardWidget{content: safeWidget(content), padding: defaultCardPadding}
}

// Padding changes the space inside the card in pixels.
func (c *CardWidget) Padding(pixels int) *CardWidget {
	if c != nil {
		c.padding = max(0, pixels)
	}
	return c
}

// Expand asks the card to use available space.
func (c *CardWidget) Expand() *CardWidget {
	if c != nil {
		c.expand = true
	}
	return c
}

func (c *CardWidget) mount(ctx *mountContext, parent *tk.Window) mountedWidget {
	frame := parent.Frame(
		tk.Background(ctx.theme.Surface.String()),
		tk.Highlightbackground(ctx.theme.Border.String()),
		tk.Highlightcolor(ctx.theme.Border.String()),
		tk.Highlightthickness(1),
		tk.Borderwidth(0),
	)
	mounted := mountOnSurface(ctx, frame.Window, c.content)
	tk.Pack(mounted.window, tk.Fill("both"), tk.Expand(true), tk.Padx(c.padding), tk.Pady(c.padding))
	return mountedWidget{
		window:  frame.Window,
		expandX: c.expand || mounted.expandX,
		expandY: c.expand || mounted.expandY,
	}
}

func mountOnSurface(ctx *mountContext, parent *tk.Window, content Widget) (mounted mountedWidget) {
	originalTheme := ctx.theme
	ctx.theme.Background = originalTheme.Surface
	defer func() { ctx.theme = originalTheme }()
	return content.mount(ctx, parent)
}

type sizeMode uint8

const (
	preferredSize sizeMode = iota
	minimumSize
)

// SizeBox gives one widget a preferred or minimum pixel size.
type SizeBox struct {
	content Widget
	width   int
	height  int
	mode    sizeMode
	expand  bool
}

// Size gives content a preferred pixel size. The content fills that area.
func Size(content Widget, width, height int) *SizeBox {
	return newSizeBox(content, width, height, preferredSize)
}

// MinSize preserves the content's natural size while requiring at least the
// supplied width and height in pixels.
func MinSize(content Widget, width, height int) *SizeBox {
	return newSizeBox(content, width, height, minimumSize)
}

func newSizeBox(content Widget, width, height int, mode sizeMode) *SizeBox {
	return &SizeBox{
		content: safeWidget(content),
		width:   max(1, width),
		height:  max(1, height),
		mode:    mode,
	}
}

// Expand asks the sized area to grow when more space is available.
func (s *SizeBox) Expand() *SizeBox {
	if s != nil {
		s.expand = true
	}
	return s
}

func (s *SizeBox) mount(ctx *mountContext, parent *tk.Window) mountedWidget {
	frame := parent.Frame(
		tk.Width(s.width),
		tk.Height(s.height),
		tk.Background(ctx.theme.Background.String()),
		tk.Borderwidth(0),
	)
	mounted := s.content.mount(ctx, frame.Window)
	if s.mode == minimumSize {
		tk.GridColumnConfigure(frame, 0, tk.Weight(1), tk.Minsize(s.width))
		tk.GridRowConfigure(frame, 0, tk.Weight(1), tk.Minsize(s.height))
		tk.Grid(mounted.window, tk.Row(0), tk.Column(0), tk.Sticky("nsew"))
	} else {
		tk.Place(mounted.window, tk.X(0), tk.Y(0), tk.Relwidth(1), tk.Relheight(1))
	}
	return mountedWidget{
		window:  frame.Window,
		expandX: s.expand || (s.mode == minimumSize && mounted.expandX),
		expandY: s.expand || (s.mode == minimumSize && mounted.expandY),
	}
}

func safeWidget(widget Widget) Widget {
	if widget == nil {
		return Spacer(0, 0)
	}
	return widget
}
