# Canvas Mouse Input

Rosaline canvases can respond to clicks, pointer movement, dragging, and mouse
button releases. The event API stays independent of the private window
backend, so the same application code works on Linux, Windows, and macOS.

## Smallest example

This program draws a circle wherever the user clicks:

```go
package main

import "github.com/SeraphinaDX/Rosaline"

func main() {
	var x, y float64
	clicked := false

	canvas := rosaline.Canvas(func(c *rosaline.DrawingCanvas) {
		c.Clear(rosaline.White)
		if clicked {
			c.FillCircle(x, y, 12, rosaline.Rose)
		}
	})

	canvas.OnMouseDown(func(event rosaline.MouseEvent) {
		if event.Button == rosaline.MouseLeft {
			x = event.X
			y = event.Y
			clicked = true
		}
	})

	rosaline.Run(canvas)
}
```

The drawing function reads normal Go variables. A mouse callback changes those
variables, and Rosaline redraws the canvas automatically after the callback.

## Mouse callbacks

A canvas provides three mouse callbacks:

- `OnMouseDown` runs when a button is pressed.
- `OnMouseMove` runs when the pointer moves over the canvas.
- `OnMouseUp` runs when a button is released.

Each callback receives a `MouseEvent`:

```go
canvas.OnMouseMove(func(event rosaline.MouseEvent) {
	fmt.Println(event.X, event.Y)
})
```

`X` and `Y` are canvas coordinates measured from the top-left corner. X grows
to the right and Y grows downward, matching Rosaline's drawing coordinates.

## Buttons and dragging

`Button` is one of:

- `rosaline.MouseLeft`
- `rosaline.MouseMiddle`
- `rosaline.MouseRight`
- `rosaline.MouseNone`, used for movement without a held button

During `OnMouseMove`, `Dragging` is true when a mouse button is held. A simple
drag handler looks like this:

```go
canvas.OnMouseMove(func(event rosaline.MouseEvent) {
	if event.Dragging && event.Button == rosaline.MouseLeft {
		points = append(points, point{x: event.X, y: event.Y})
	}
})
```

The event also has `Shift`, `Control`, and `Alt` Boolean fields. These are true
when the corresponding keyboard modifier is held during the mouse event.

## Keyboard input

A canvas can also receive backend-neutral key events:

```go
canvas.Focus().OnKeyDown(func(event rosaline.KeyEvent) {
	if event.Is(rosaline.KeyRight) {
		x += 10
	}
})
```

Keyboard-enabled canvases participate in Tab order, receive focus when clicked,
show a focus ring, and redraw automatically after their key callbacks. See
[KEYBOARD_INPUT.md](KEYBOARD_INPUT.md) and the complete
[Keyboard Garden application](KEYBOARD_GARDEN_APPLICATION.md).

Calling `Focus()` again from a button or menu callback returns keyboard focus
to an already-open canvas. This is useful for games and other keyboard-driven
tools. The complete [Starshower application](STARSHOWER_APPLICATION.md) shows
held-key movement and firing.

## Redrawing from another control

Mouse callbacks redraw automatically. If a button or another Rosaline callback
changes drawing state, call `Redraw` yourself:

```go
rosaline.Button("Clear", func() {
	points = nil
	canvas.Redraw()
})
```

`Redraw` clears the old canvas shapes and runs the original drawing function
again. Call it from Rosaline callbacks, not directly from a background
goroutine.

## Building a drawing program

A small paint program needs three pieces of state:

1. A slice containing the finished strokes.
2. A Boolean that remembers whether drawing is active.
3. The points in the current stroke.

`OnMouseDown` starts a stroke, `OnMouseMove` adds points while the left button
is held, and `OnMouseUp` finishes it. The drawing function loops through all
stored points and connects neighboring points with `Line`.

Run the complete example from the Rosaline source tree:

```bash
CGO_ENABLED=0 go run ./examples/paint
```

The source is in [`examples/paint/main.go`](../examples/paint/main.go).

## Common mistakes

### Drawing only inside the event handler

Keep the drawing data in Go variables and render it in the canvas drawing
function. That way clearing and redrawing the window always produces the same
picture.

### Forgetting to check the button

Without a button check, right and middle clicks may start drawing too:

```go
if event.Button != rosaline.MouseLeft {
	return
}
```

### Calling Redraw from every mouse callback

Rosaline already redraws after mouse callbacks. Explicit `Redraw` is for
buttons and other callbacks that change the picture.

## Go concepts used here

- Structs and slices
- Boolean state
- Callback functions
- `for` loops
- Coordinates and floating-point numbers
