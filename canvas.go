// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import (
	"fmt"

	tk "modernc.org/tk9.0"
)

// TextStyle controls text drawn on a Canvas.
type TextStyle struct {
	Color Color
	Size  int
}

// DrawingCanvas provides Rosaline's beginner-friendly 2D drawing operations.
type DrawingCanvas struct {
	widget     *tk.CanvasWidget
	background Color
}

// Clear removes existing shapes and changes the canvas background.
func (c *DrawingCanvas) Clear(color Color) {
	c.background = color
	c.widget.Delete("all")
	c.widget.Configure(tk.Background(color.String()))
}

// FillRect draws a filled rectangle.
func (c *DrawingCanvas) FillRect(x, y, width, height float64, color Color) {
	c.widget.CreateRectangle(x, y, x+width, y+height,
		tk.Fill(color.String()), tk.Outline(color.String()))
}

// Rect draws the outline of a rectangle.
func (c *DrawingCanvas) Rect(x, y, width, height, stroke float64, color Color) {
	c.widget.CreateRectangle(x, y, x+width, y+height,
		tk.Fill(""), tk.Outline(color.String()), tk.Width(stroke))
}

// Line draws a line.
func (c *DrawingCanvas) Line(x1, y1, x2, y2, stroke float64, color Color) {
	c.widget.CreateLine(x1, y1, x2, y2, tk.Fill(color.String()), tk.Width(stroke))
}

// FillCircle draws a filled circle.
func (c *DrawingCanvas) FillCircle(x, y, radius float64, color Color) {
	c.widget.CreateOval(x-radius, y-radius, x+radius, y+radius,
		tk.Fill(color.String()), tk.Outline(color.String()))
}

// Circle draws the outline of a circle.
func (c *DrawingCanvas) Circle(x, y, radius, stroke float64, color Color) {
	c.widget.CreateOval(x-radius, y-radius, x+radius, y+radius,
		tk.Fill(""), tk.Outline(color.String()), tk.Width(stroke))
}

// Text draws text from its top-left corner.
func (c *DrawingCanvas) Text(text string, x, y float64, style TextStyle) {
	if style.Color.A == 0 {
		style.Color = Black
	}
	if style.Size <= 0 {
		style.Size = 14
	}
	c.widget.CreateText(x, y,
		tk.Txt(text),
		tk.Fill(style.Color.String()),
		tk.Anchor("nw"),
		tk.Font(fmt.Sprintf("TkDefaultFont %d", style.Size)),
	)
}

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
	drawing     *DrawingCanvas
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
	if c.drawing == nil {
		return
	}

	c.drawing.widget.Delete("all")
	c.drawing.widget.Configure(tk.Background(c.background.String()))
	c.drawing.background = c.background
	c.draw(c.drawing)
	c.background = c.drawing.background
	c.redrawCount++
}

func (c *CanvasWidget) mount(ctx *mountContext, parent *tk.Window) mountedWidget {
	widget := parent.Canvas(
		tk.Width(c.width),
		tk.Height(c.height),
		tk.Background(c.background.String()),
		tk.Highlightthickness(0),
		tk.Borderwidth(0),
	)
	c.drawing = &DrawingCanvas{widget: widget, background: c.background}
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
