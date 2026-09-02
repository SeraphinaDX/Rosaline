// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import (
	tk "modernc.org/tk9.0"
)

// CanvasWidget is a custom 2D drawing surface.
type CanvasWidget struct {
	draw        func(*DrawingCanvas)
	width       int
	height      int
	background  Color
	expand      bool
	focus       bool
	onMouseDown func(MouseEvent)
	onMouseMove func(MouseEvent)
	onMouseUp   func(MouseEvent)
	onKeyDown   func(KeyEvent)
	onKeyUp     func(KeyEvent)
	widget      *tk.CanvasWidget
	tkImage     *tk.Img
	redrawCount uint64
}

// Canvas creates a 2D drawing surface.
func Canvas(draw func(*DrawingCanvas)) *CanvasWidget {
	if draw == nil {
		draw = func(*DrawingCanvas) {}
	}
	return &CanvasWidget{
		draw:       draw,
		width:      480,
		height:     300,
		background: White,
	}
}

// Size sets the canvas's initial size in pixels.
func (c *CanvasWidget) Size(width, height int) *CanvasWidget {
	if width > 0 {
		c.width = width
	}
	if height > 0 {
		c.height = height
	}
	return c
}

// Background sets the canvas background.
func (c *CanvasWidget) Background(color Color) *CanvasWidget {
	c.background = color
	return c
}

// Expand asks the canvas to use available layout space.
func (c *CanvasWidget) Expand() *CanvasWidget {
	c.expand = true
	return c
}

// Focus asks Rosaline to give this canvas keyboard focus when the window
// opens. A canvas with a key handler also participates in Tab focus order.
func (c *CanvasWidget) Focus() *CanvasWidget {
	if c != nil {
		c.focus = true
	}
	return c
}

// OnMouseDown runs when a mouse button is pressed over the canvas.
func (c *CanvasWidget) OnMouseDown(handler func(MouseEvent)) *CanvasWidget {
	c.onMouseDown = handler
	return c
}

// OnMouseMove runs when the pointer moves over the canvas. Event Button is
// MouseNone for normal movement and identifies the held button while dragging.
func (c *CanvasWidget) OnMouseMove(handler func(MouseEvent)) *CanvasWidget {
	c.onMouseMove = handler
	return c
}

// OnMouseUp runs when a mouse button is released over the canvas.
func (c *CanvasWidget) OnMouseUp(handler func(MouseEvent)) *CanvasWidget {
	c.onMouseUp = handler
	return c
}

// OnKeyDown runs when a key is pressed while the canvas has focus. Clicking a
// keyboard-enabled canvas gives it focus.
func (c *CanvasWidget) OnKeyDown(handler func(KeyEvent)) *CanvasWidget {
	if c != nil {
		c.onKeyDown = handler
	}
	return c
}

// OnKeyUp runs when a key is released while the canvas has focus.
func (c *CanvasWidget) OnKeyUp(handler func(KeyEvent)) *CanvasWidget {
	if c != nil {
		c.onKeyUp = handler
	}
	return c
}

// Redraw clears the canvas and runs its drawing function again. Call Redraw
// from Rosaline callbacks after changing drawing state. Mouse callbacks redraw
// automatically, so they normally do not need to call it themselves.
func (c *CanvasWidget) Redraw() {
	if c.widget == nil || c.tkImage == nil {
		return
	}
	pixels, background := renderDrawing(c.width, c.height, c.background, c.draw)
	c.background = background
	c.widget.Configure(tk.Background(background.String()))
	if err := c.tkImage.Put(pixels); err != nil {
		oldImage := c.tkImage
		c.widget.Delete("all")
		c.tkImage = tk.NewPhoto(tk.Data(pixels))
		c.widget.CreateImage(0, 0, tk.Image(c.tkImage), tk.Anchor("nw"))
		oldImage.Delete()
	}
	c.redrawCount++
}

// Picture renders the canvas into an off-screen Picture. The result can be
// saved as PNG or AVIF and is available even before the widget is mounted.
func (c *CanvasWidget) Picture() *Picture {
	if c == nil {
		return nil
	}
	pixels, _ := renderDrawing(c.width, c.height, c.background, c.draw)
	return NewPicture(pixels)
}

func (c *CanvasWidget) mount(ctx *mountContext, parent *tk.Window) mountedWidget {
	keyboardEnabled := c.focus || c.onKeyDown != nil || c.onKeyUp != nil
	highlightThickness := 0
	if keyboardEnabled {
		highlightThickness = 2
	}
	widget := parent.Canvas(
		tk.Width(c.width),
		tk.Height(c.height),
		tk.Background(c.background.String()),
		tk.Highlightthickness(highlightThickness),
		tk.Highlightcolor(ctx.theme.Primary.String()),
		tk.Highlightbackground(ctx.theme.Border.String()),
		tk.Borderwidth(0),
		takeFocusOption(keyboardEnabled),
	)
	c.widget = widget
	c.tkImage = tk.NewPhoto(tk.Width(c.width), tk.Height(c.height))
	widget.CreateImage(0, 0, tk.Image(c.tkImage), tk.Anchor("nw"))
	c.Redraw()
	ctx.addCleanup(func() {
		if c.tkImage != nil {
			c.tkImage.Delete()
		}
		c.tkImage = nil
		c.widget = nil
	})

	dispatch := func(handler func(MouseEvent), event MouseEvent) {
		if handler == nil {
			return
		}
		before := c.redrawCount
		handler(event)
		if c.redrawCount == before {
			c.Redraw()
		}
		ctx.refresh()
	}

	for _, binding := range []struct {
		sequence string
		button   MouseButton
	}{
		{"<Button-1>", MouseLeft},
		{"<Button-2>", MouseMiddle},
		{"<Button-3>", MouseRight},
	} {
		if c.onMouseDown != nil || (keyboardEnabled && binding.button == MouseLeft) {
			button := binding.button
			tk.Bind(widget.Window, binding.sequence, tk.Command(func(event *tk.Event) {
				if keyboardEnabled && button == MouseLeft {
					tk.Focus(widget.Window)
				}
				dispatch(c.onMouseDown, mouseEvent(event, button, true))
			}))
		}
	}

	if c.onMouseMove != nil {
		tk.Bind(widget.Window, "<Motion>", tk.Command(func(event *tk.Event) {
			button := heldMouseButton(event.State)
			dispatch(c.onMouseMove, mouseEvent(event, button, button != MouseNone))
		}))
	}

	if c.onMouseUp != nil {
		for _, binding := range []struct {
			sequence string
			button   MouseButton
		}{
			{"<ButtonRelease-1>", MouseLeft},
			{"<ButtonRelease-2>", MouseMiddle},
			{"<ButtonRelease-3>", MouseRight},
		} {
			button := binding.button
			tk.Bind(widget.Window, binding.sequence, tk.Command(func(event *tk.Event) {
				dispatch(c.onMouseUp, mouseEvent(event, button, false))
			}))
		}
	}

	dispatchKey := func(handler func(KeyEvent), event *tk.Event) {
		if handler == nil {
			return
		}
		before := c.redrawCount
		handler(keyEvent(event))
		if c.redrawCount == before {
			c.Redraw()
		}
		ctx.refresh()
	}
	if c.onKeyDown != nil {
		tk.Bind(widget.Window, "<KeyPress>", tk.Command(func(event *tk.Event) {
			dispatchKey(c.onKeyDown, event)
		}))
	}
	if c.onKeyUp != nil {
		tk.Bind(widget.Window, "<KeyRelease>", tk.Command(func(event *tk.Event) {
			dispatchKey(c.onKeyUp, event)
		}))
	}
	if keyboardEnabled {
		ctx.addFocusable(widget.Window, c.focus)
	}

	return mountedWidget{window: widget.Window, expandX: c.expand, expandY: c.expand}
}
