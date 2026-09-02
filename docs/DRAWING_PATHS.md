# Paths and Bézier Curves

A `Path` describes a reusable shape made from lines and curves. Paths are
useful for icons, diagrams, smooth outlines, custom controls, and drawing
applications.

## Complete heart example

```go
package main

import "github.com/SeraphinaDX/Rosaline"

func main() {
	heart := rosaline.NewPath().
		MoveTo(160, 80).
		CubicTo(80, 20, 30, 130, 160, 240).
		CubicTo(290, 130, 240, 20, 160, 80).
		Close()

	rosaline.Run(
		rosaline.Canvas(func(c *rosaline.DrawingCanvas) {
			c.Clear(rosaline.White)
			c.FillPath(heart, rosaline.SoftRose)
			c.StrokePath(heart, 4, rosaline.Rose)
		}).Size(320, 280),
	)
}
```

The path stores geometry but not color. This lets the same path be filled or
outlined several times with different styles.

## Building a path

Every path begins with `NewPath`:

```go
shape := rosaline.NewPath()
```

Then add commands:

- `MoveTo(x, y)` begins a new section without drawing a line.
- `LineTo(x, y)` adds a straight line.
- `QuadraticTo(controlX, controlY, x, y)` adds a curve with one control point.
- `CubicTo(control1X, control1Y, control2X, control2Y, x, y)` adds a curve with
  two control points.
- `Close()` connects the current point to the section's starting point.

Methods return the same path, so they can be chained as shown in the example.

## Quadratic and cubic curves

A control point pulls a curve toward itself without necessarily lying on the
curve. Quadratic curves have one control point and are good for simple bends:

```go
wave := rosaline.NewPath().
	MoveTo(20, 100).
	QuadraticTo(100, 20, 180, 100)
```

Cubic curves have two control points and give independent control over how the
curve leaves its start and approaches its end:

```go
wave.CubicTo(220, 180, 280, 20, 340, 100)
```

## Filling and stroking

`FillPath` fills the inside of a path. Open sections are closed for filling.
`StrokePath` draws its outline:

```go
c.FillPath(shape, rosaline.SoftRose)
c.StrokePath(shape, 3, rosaline.Rose)
```

The number `3` is the stroke width in pixels.

## Paths in Paint

The Paint example stores mouse positions as Go values, creates a path from
each stroke, and then draws it with `StrokePath`. This is a useful pattern:
application state remains normal Go data, while paths are created when the
canvas renders.

Run it with:

```bash
CGO_ENABLED=0 go run ./examples/paint
```

## Common mistakes

### Forgetting MoveTo

Begin each separate shape section with `MoveTo`. It establishes the point from
which the first line or curve starts.

### Putting color in the path

A path only describes geometry. Pass color to `FillPath` or `StrokePath` so
the shape can be reused.

### Rebuilding a fixed path every redraw

If a shape never changes, create it once outside the canvas drawing callback,
as the heart example does.

## Go concepts used here

- Struct pointers
- Method chaining
- Reusing values
- Separating data from presentation
