# Rosaline Quick Start

This guide creates a small counter application. It introduces a window, text,
layout, buttons, callbacks, and changing state without assuming prior GUI
experience.

## 1. Create the project

```bash
mkdir rosaline-counter
cd rosaline-counter
go mod init example.com/rosaline-counter
go get github.com/SeraphinaDX/Rosaline
```

Create `main.go`:

```go
package main

import (
	"fmt"

	"github.com/SeraphinaDX/Rosaline"
)

func main() {
	count := rosaline.NewState(0)

	rosaline.RunApp(rosaline.App{
		Title:  "My First Rosaline App",
		Width:  420,
		Height: 240,
		Content: rosaline.Column(
			rosaline.Label("My counter"),
			rosaline.LabelFunc(func() string {
				return fmt.Sprintf("Count: %d", count.Get())
			}),
			rosaline.Button("Add one", func() {
				count.Update(func(n int) int { return n + 1 })
			}).Primary(),
		).Gap(12),
	})
}
```

Run it:

```bash
CGO_ENABLED=0 go run .
```

## 2. How the application works

`NewState(0)` creates an integer whose first value is zero:

```go
count := rosaline.NewState(0)
```

Go infers that this is `State[int]` because `0` is an integer.

`RunApp` opens the window and starts its event loop. Rosaline handles window
initialization and platform details:

```go
rosaline.RunApp(rosaline.App{
	Title: "My First Rosaline App",
	// ...
})
```

`Column` puts its children from top to bottom:

```go
rosaline.Column(
	rosaline.Label("First"),
	rosaline.Label("Second"),
)
```

`LabelFunc` recalculates its text after a Rosaline event. The anonymous
function reads the latest count:

```go
rosaline.LabelFunc(func() string {
	return fmt.Sprintf("Count: %d", count.Get())
})
```

The button receives another anonymous function. Rosaline calls it when the
button is activated:

```go
rosaline.Button("Add one", func() {
	count.Update(func(n int) int { return n + 1 })
})
```

## Go concepts used here

- Variables
- Function calls
- Struct literals
- Anonymous functions
- Closures
- Generic types (inferred automatically)

You do not need to master all of those terms before changing the example.
Experiment with the title, starting number, button text, and window size.

## 3. Draw something

Add a Canvas anywhere a normal widget can go:

```go
rosaline.Canvas(func(c *rosaline.DrawingCanvas) {
	c.Clear(rosaline.White)
	c.FillCircle(100, 100, 50, rosaline.Rose)
	c.Rect(25, 25, 150, 150, 3, rosaline.Black)
}).Size(300, 220)
```

The drawing callback receives a canvas. Coordinates start at the top-left:
positive X moves right and positive Y moves down.

## Common mistakes

### The import cannot be found

Run this inside your application folder:

```bash
go get github.com/SeraphinaDX/Rosaline
```

### Linux says it cannot open a display

Run the program from inside your graphical desktop session. The first backend
uses X11 or XWayland. Native Wayland support is planned.

### A label does not change

Use `LabelFunc` for changing text and `Label` for fixed text.

## Next steps

- Build a form with [TEXT_INPUT.md](TEXT_INPUT.md) and
  [CHECKBOX.md](CHECKBOX.md).
- Run `go run ./examples/forms` to see input, validation, and submission used
  together.
- Run `go run ./examples/canvas` from the Rosaline source tree.
- Read `docs/ROADMAP.md` to see which features come next.
- Build a two-button counter with Add and Subtract actions.
