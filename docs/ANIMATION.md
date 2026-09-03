# Canvas Animation

Animation is a small loop: update ordinary Go variables, redraw the canvas,
then repeat. Rosaline's `Animate` helper creates the repeating timer for that
loop.

## Smallest animation

This complete program moves a circle across a canvas:

```go
package main

import "github.com/SeraphinaDX/Rosaline"

func main() {
	x := 20.0

	var canvas *rosaline.CanvasWidget
	canvas = rosaline.Canvas(func(c *rosaline.DrawingCanvas) {
		c.Clear(rosaline.White)
		c.FillCircle(x, 80, 14, rosaline.Rose)
	}).Size(500, 160)

	animation := rosaline.Animate(60, func() {
		x += 2
		if x > 500 {
			x = 0
		}
		canvas.Redraw()
	})

	rosaline.RunApp(rosaline.App{
		Title:   "Moving Circle",
		Timers:  []*rosaline.Timer{animation},
		Content: canvas,
	})
}
```

`Animate(60, ...)` asks for 60 frames per second. It is a target rather than a
promise: a busy computer or expensive drawing function may produce fewer
frames, while the interface remains responsive.

## Why canvas is declared first

The animation callback needs to call `canvas.Redraw`, but the canvas drawing
function also needs the changing `x` variable. Declaring the pointer first
lets both callbacks refer to the same canvas:

```go
var canvas *rosaline.CanvasWidget
canvas = rosaline.Canvas(...)
```

This is a normal Go closure pattern. Rosaline does not keep a separate scene or
animation language.

## Bouncing at an edge

A moving value normally has a position and a speed:

```go
x += speed
if x-radius <= 0 || x+radius >= width {
	speed = -speed
}
```

Negating `speed` reverses the direction. The complete animation example uses
the same idea on both axes and adds orbiting petals.

Run it from the project root:

```bash
CGO_ENABLED=0 go run ./examples/animation
```

The source is in [`examples/animation/main.go`](../examples/animation/main.go).

## Pausing and resuming

An animation is an ordinary `Timer`, so buttons can control it:

```go
rosaline.Button("Pause", animation.Stop)
rosaline.Button("Resume", animation.Start)
```

The drawing stays visible while the timer is stopped. Starting it continues
with the same Go state.

## Choosing a frame rate

- 60 FPS is a good default for smooth motion.
- 30 FPS reduces drawing work and is enough for many utilities.
- Lower rates are useful for slow visual changes.

Keep each frame callback short. Store the scene in Go variables and let the
canvas drawing function render those variables.

## Games and fixed simulation steps

For a game, measure the real time between drawing frames but update the model
in small, consistent steps. This makes movement and collision behavior stable
when frame delivery varies:

```go
accumulator += elapsed.Seconds()
for accumulator >= fixedStep {
	update(fixedStep)
	accumulator -= fixedStep
}
canvas.Redraw()
```

Cap unusually large elapsed durations so opening a dialog or moving a window
does not make the game jump forward. The complete
[Starshower application](STARSHOWER_APPLICATION.md) demonstrates the pattern
with held keys, movement, firing, collisions, pause, and restart.

## Common mistakes

### Forgetting Redraw

Timer callbacks refresh `LabelFunc` and other dynamic widgets automatically,
but a canvas is only redrawn when requested:

```go
animation := rosaline.Animate(60, func() {
	x++
	canvas.Redraw()
})
```

### Drawing only inside the timer callback

Put drawing commands in the canvas drawing function. The timer should change
state and request a redraw. This keeps the picture reproducible whenever the
canvas needs to draw again.

### Doing expensive work every frame

Loading files or processing large images in every frame will make animation
uneven. Load resources once, then reuse them while drawing.

## Go concepts used here

- Closures
- Floating-point variables
- Assignment and arithmetic
- Conditional statements
- Function values used as callbacks
