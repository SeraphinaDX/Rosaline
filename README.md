# Rosaline

Rosaline is a small, beginner-friendly graphics and GUI library for Go. It is
designed for people who know a little Go and want to make a real graphical
program without first learning a large framework.

Rosaline is currently at `v0.4.1`. The public API is small on purpose and grows
through well-documented, tested features.

## Goals

- Beginner-friendly Go API
- Builds with `CGO_ENABLED=0`
- Linux is a first-class platform
- The same application code runs on Linux, Windows, and macOS
- Serious 2D drawing alongside normal GUI widgets
- Small, memorable public API
- Complete examples and feature-by-feature documentation

## Install

Rosaline uses GitHub's normal Go module support:

```bash
go get github.com/SeraphinaDX/Rosaline
```

Rosaline requires Go 1.25 or newer. No C compiler or separately installed GUI
toolkit is required. On Linux, a graphical desktop with X11 or XWayland is
currently required; native Wayland support is on the roadmap.

## Hello, Rosaline

```go
package main

import "github.com/SeraphinaDX/Rosaline"

func main() {
	rosaline.Run(
		rosaline.Column(
			rosaline.Label("Hello, world!"),
			rosaline.Button("Click me", func() {
				rosaline.Message("Rosaline", "It works!")
			}).Primary(),
		),
	)
}
```

Save that as `main.go`, then run:

```bash
CGO_ENABLED=0 go run .
```

See [docs/QUICK_START.md](docs/QUICK_START.md) for a guided first application.

## A small form

Form controls update normal Go variables through pointers:

```go
var name string
var subscribed bool

rosaline.Run(
	rosaline.Column(
		rosaline.Label("Your name"),
		rosaline.TextBox(&name).Placeholder("Type your name").Focus(),
		rosaline.CheckBox("Send me updates", &subscribed),
		rosaline.Button("Continue", func() {
			rosaline.Message("Hello", "Welcome, "+name+"!")
		}).Primary(),
	),
)
```

See [docs/TEXT_INPUT.md](docs/TEXT_INPUT.md),
[docs/CHECKBOX.md](docs/CHECKBOX.md), and [docs/FORMS.md](docs/FORMS.md).

## An interactive canvas

Canvas callbacks use the same coordinates as drawing commands:

```go
var x, y float64

canvas := rosaline.Canvas(func(c *rosaline.DrawingCanvas) {
	c.Clear(rosaline.White)
	c.FillCircle(x, y, 12, rosaline.Rose)
})

canvas.OnMouseDown(func(event rosaline.MouseEvent) {
	if event.Button == rosaline.MouseLeft {
		x, y = event.X, event.Y
	}
})
```

Rosaline redraws automatically after mouse callbacks. See
[docs/CANVAS_INPUT.md](docs/CANVAS_INPUT.md) for clicking, dragging, modifier
keys, manual redraws, and a complete paint program.

## A real image-viewer application

Rosaline v0.4 combines menus, dialogs, images, and scrolling:

```go
viewer := rosaline.Image(nil)

openImage := func() {
	path, ok := rosaline.OpenFileDialog(rosaline.FileDialogOptions{
		Title: "Open Image",
	})
	if !ok {
		return
	}
	picture, err := rosaline.LoadImage(path)
	if err != nil {
		rosaline.Error("Could not open image", err.Error())
		return
	}
	viewer.SetImage(picture)
}
```

See [docs/IMAGE_VIEWER.md](docs/IMAGE_VIEWER.md) for the complete application.

## Timers and animation

Timers belong to the application and automatically stop with its event loop:

```go
seconds := 0

clock := rosaline.Every(time.Second, func() {
	seconds++
})

rosaline.RunApp(rosaline.App{
	Timers: []*rosaline.Timer{clock},
	Content: rosaline.LabelFunc(func() string {
		return fmt.Sprintf("Running for %d seconds", seconds)
	}),
})
```

Use `After` for one delayed callback and `Animate` for a frame-rate-based
canvas loop. See [docs/TIMERS.md](docs/TIMERS.md) and
[docs/ANIMATION.md](docs/ANIMATION.md).

## Included in v0.4.1

- Application windows
- Labels and dynamic labels
- Buttons and message dialogs
- Rows, columns, spacing, padding, and expansion
- Simple state values
- Single-line and multiline text input
- Password display, placeholders, change events, and Enter submission
- Checkboxes bound to Go Boolean variables
- Tab and Shift+Tab keyboard navigation
- A first-class canvas with lines, rectangles, circles, and text
- Canvas clicks, pointer movement, dragging, and button-release events
- Automatic and manually requested canvas redraws
- Pure-Go loading and display of common image formats
- Horizontal and vertical scroll areas
- Native open, save, message, error, and confirmation dialogs
- Menu bars with working keyboard shortcuts
- Repeating and one-shot application timers
- Start, stop, restart, and running-state timer controls
- Frame-rate-based canvas animation
- Semantic colors and themes
- Runnable hello, counter, canvas, forms, paint, image-viewer, and animation
  examples
- Unit tests for non-visual core behavior

## Run the examples

From the extracted project root:

```bash
CGO_ENABLED=0 go run ./examples/hello
CGO_ENABLED=0 go run ./examples/counter
CGO_ENABLED=0 go run ./examples/canvas
CGO_ENABLED=0 go run ./examples/forms
CGO_ENABLED=0 go run ./examples/paint
CGO_ENABLED=0 go run ./examples/imageviewer
CGO_ENABLED=0 go run ./examples/animation
```

## Project status

The API is experimental until v1.0. The next milestones add paths, Bézier
curves, transforms, clipping, image export, radio buttons, and general keyboard
events.

Rosaline's backend is intentionally private. Application code only imports the
`rosaline` package, so the backend can improve without forcing beginners to
rewrite their programs.

## License

Rosaline is free software licensed under the
[GNU Lesser General Public License v3.0 or later](LICENSE).

Applications may use and link to Rosaline without being required to adopt the
LGPL. Modifications to Rosaline itself must remain available under the LGPL,
and distribution must follow the license's relinking and source-availability
requirements. The incorporated GNU GPL v3 text is included in
[LICENSE.GPL](LICENSE.GPL).

Copyright (C) 2026 Britney Lozza and Rosaline contributors.
