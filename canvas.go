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
	onMouseDown func(MouseEvent)
	onMouseMove func(MouseEvent)
	onMouseUp   func(MouseEvent)
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
	widget := parent.Canvas(
		tk.Width(c.width),
		tk.Height(c.height),
		tk.Background(c.background.String()),
		tk.Highlightthickness(0),
		tk.Borderwidth(0),
	)
	c.widget = widget
	c.tkImage = tk.NewPhoto(tk.Width(c.width), tk.Height(c.height))
	widget.CreateImage(0, 0, tk.Image(c.tkImage), tk.Anchor("nw"))
	c.Redraw()

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

	if c.onMouseDown != nil {
		for _, binding := range []struct {
			sequence string
			button   MouseButton
		}{
			{"<Button-1>", MouseLeft},
			{"<Button-2>", MouseMiddle},
			{"<Button-3>", MouseRight},
		} {
			button := binding.button
			tk.Bind(widget.Window, binding.sequence, tk.Command(func(event *tk.Event) {
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

	return mountedWidget{window: widget.Window, expandX: c.expand, expandY: c.expand}
}
