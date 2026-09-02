// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import (
	"strconv"
	"strings"

	tk "modernc.org/tk9.0"
)

// ScrollWidget displays content inside a viewport with horizontal and vertical
// scrollbars.
type ScrollWidget struct {
	content Widget
	width   int
	height  int
	expand  bool
}

// Scroll creates a scrollable viewport around content.
func Scroll(content Widget) *ScrollWidget {
	if content == nil {
		content = Label("This scroll area has no content.")
	}
	return &ScrollWidget{content: content, width: 480, height: 300}
}

// Size sets the viewport's preferred size in pixels.
func (s *ScrollWidget) Size(width, height int) *ScrollWidget {
	if width > 0 {
		s.width = width
	}
	if height > 0 {
		s.height = height
	}
	return s
}

// Expand asks the scroll area to use available layout space.
func (s *ScrollWidget) Expand() *ScrollWidget {
	s.expand = true
	return s
}

func (s *ScrollWidget) mount(ctx *mountContext, parent *tk.Window) mountedWidget {
	frame := parent.Frame(
		tk.Background(ctx.theme.Background.String()),
		tk.Borderwidth(0),
	)

	var canvas *tk.CanvasWidget
	xScroll := frame.Scrollbar(
		tk.Orient("horizontal"),
		takeFocusOption(false),
		tk.Command(func(event *tk.Event) { event.Xview(canvas) }),
	)
	yScroll := frame.Scrollbar(
		tk.Orient("vertical"),
		takeFocusOption(false),
		tk.Command(func(event *tk.Event) { event.Yview(canvas) }),
	)
	canvas = frame.Canvas(
		tk.Width(s.width),
		tk.Height(s.height),
		tk.Background(ctx.theme.Surface.String()),
		tk.Highlightthickness(1),
		tk.Highlightbackground(ctx.theme.Border.String()),
		tk.Borderwidth(0),
		tk.Xscrollcommand(func(event *tk.Event) { event.ScrollSet(xScroll) }),
		tk.Yscrollcommand(func(event *tk.Event) { event.ScrollSet(yScroll) }),
	)

	inner := canvas.Frame(
		tk.Background(ctx.theme.Surface.String()),
		tk.Borderwidth(0),
	)
	content := s.content.mount(ctx, inner.Window)
	packOptions := []tk.Opt{content.window, tk.Anchor("nw")}
	if content.expandX || content.expandY {
		packOptions = append(packOptions, tk.Fill("both"), tk.Expand(true))
	}
	tk.Pack(packOptions...)
	canvas.CreateWindow(0, 0, tk.ItemWindow(inner.Window), tk.Anchor("nw"))

	updateRegion := func() {
		if bounds := canvas.Bbox("all"); len(bounds) == 4 {
			canvas.Configure(tk.Scrollregion(strings.Join(bounds, " ")))
		}
	}
	tk.Bind(inner.Window, "<Configure>", tk.Command(updateRegion))

	tk.Grid(canvas, tk.Row(0), tk.Column(0), tk.Sticky("nsew"))
	tk.Grid(yScroll, tk.Row(0), tk.Column(1), tk.Sticky("ns"))
	tk.Grid(xScroll, tk.Row(1), tk.Column(0), tk.Sticky("ew"))
	tk.GridRowConfigure(frame.Window, 0, tk.Weight(1))
	tk.GridColumnConfigure(frame.Window, 0, tk.Weight(1))

	bindScrollWheel(canvas)
	updateRegion()

	return mountedWidget{window: frame.Window, expandX: s.expand, expandY: s.expand}
}

func bindScrollWheel(canvas *tk.CanvasWidget) {
	scroll := func(event *tk.Event, units int) {
		command := &tk.Event{Args: []string{"scroll", strconv.Itoa(units), "units"}}
		if event.State&tk.ModifierShift != 0 {
			command.Xview(canvas)
		} else {
			command.Yview(canvas)
		}
		event.SetReturnCodeBreak()
	}

	var bindTree func(*tk.Window)
	bindTree = func(window *tk.Window) {
		tk.Bind(window, "<MouseWheel>", tk.Command(func(event *tk.Event) {
			switch {
			case event.Delta > 0:
				scroll(event, -3)
			case event.Delta < 0:
				scroll(event, 3)
			}
		}))
		tk.Bind(window, "<Button-4>", tk.Command(func(event *tk.Event) {
			scroll(event, -3)
		}))
		tk.Bind(window, "<Button-5>", tk.Command(func(event *tk.Event) {
			scroll(event, 3)
		}))
		for _, child := range tk.WinfoChildren(window) {
			bindTree(child)
		}
	}

	bindTree(canvas.Window)
}
