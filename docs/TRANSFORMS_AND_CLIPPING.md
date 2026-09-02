# Transforms and Clipping

Transforms move, resize, and rotate drawing without changing the original
coordinates. Clipping restricts drawing to a chosen region.

## Complete transformed drawing

```go
package main

import "github.com/SeraphinaDX/Rosaline"

func main() {
	rosaline.Run(
		rosaline.Canvas(func(c *rosaline.DrawingCanvas) {
			c.Clear(rosaline.White)
			c.ClipRect(40, 30, 320, 220)

			for angle := 0.0; angle < 360; angle += 30 {
				c.Push()
				c.Translate(200, 140)
				c.Rotate(angle)
				c.Scale(1.4, 0.7)
				c.FillCircle(90, 0, 18, rosaline.SoftRose)
				c.Pop()
			}

			c.ResetClip()
			c.Rect(40, 30, 320, 220, 3, rosaline.Rose)
		}).Size(400, 280),
	)
}
```

## Translation

`Translate` moves the coordinate origin:

```go
c.Translate(100, 50)
c.FillCircle(0, 0, 20, rosaline.Rose)
```

The circle appears at 100, 50 even though it is drawn at 0, 0.

## Rotation

`Rotate` uses clockwise degrees:

```go
c.Rotate(45)
```

Rotation happens around the current origin. Translate to an object's center
before rotating when it should spin in place.

## Scale

`Scale` multiplies horizontal and vertical coordinates:

```go
c.Scale(2, 2)   // twice as large
c.Scale(1, 0.5) // normal width, half height
```

Scaling a circle differently on each axis produces an ellipse. Stroke widths
scale with their shapes.

## Push and Pop

Transforms affect every later drawing command. `Push` saves the current state,
and `Pop` restores it:

```go
c.Push()
c.Translate(200, 100)
c.Rotate(30)
c.FillRect(-30, -20, 60, 40, rosaline.Rose)
c.Pop()

// Normal coordinates are restored here.
```

Calls may be nested. An unmatched `Pop` safely does nothing.

`ResetTransform` is useful when drawing should return directly to window
coordinates without popping saved states.

## Clipping

`ClipRect` limits later drawing to a rectangle:

```go
c.ClipRect(20, 20, 200, 100)
c.FillCircle(220, 70, 80, rosaline.Rose)
```

Only the part of the circle inside the rectangle is visible. `Clip` accepts a
`rosaline.Rect` value when keeping the rectangle in a variable is clearer:

```go
area := rosaline.Rect{X: 20, Y: 20, Width: 200, Height: 100}
c.Clip(area)
```

New clipping regions intersect existing ones. `ResetClip` removes all active
clips. Push and Pop also save and restore clipping regions.

## See everything together

The advanced drawing example combines paths, Bézier curves, nested transforms,
and clipping:

```bash
CGO_ENABLED=0 go run ./examples/drawing
```

## Common mistakes

### Rotating before translating

Transform order matters. To rotate around a center, translate to the center,
rotate, and then draw around 0, 0.

### Forgetting Pop

Pair each Push with a Pop. Keeping the pair close together makes transformed
drawing easier to read.

### Clipping before Clear

`Clear` always clears the complete canvas. Usually call it first, then apply
clipping for the scene.

## Go concepts used here

- Loops
- Floating-point arithmetic
- Struct literals
- Scopes and paired operations
