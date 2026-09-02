# Rosaline

Rosaline is a small, beginner-friendly graphics and GUI library for Go. It is
designed for people who know a little Go and want to make a real graphical
program without first learning a large framework.

Rosaline is currently at `v0.8.0`. The public API is small on purpose and grows
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

The first build can take a while because the CGo-free window and AVIF backends
must be compiled. Go caches them, so later builds are normally much faster.

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

## Advanced drawing and image export

Paths can be reused, transformed, clipped, filled, and outlined:

```go
heart := rosaline.NewPath().
	MoveTo(160, 80).
	CubicTo(80, 20, 30, 130, 160, 240).
	CubicTo(290, 130, 240, 20, 160, 80).
	Close()

canvas := rosaline.Canvas(func(c *rosaline.DrawingCanvas) {
	c.Clear(rosaline.White)
	c.FillPath(heart, rosaline.SoftRose)
	c.StrokePath(heart, 4, rosaline.Rose)
})
```

The same drawing engine produces off-screen images:

```go
picture := canvas.Picture()
err := picture.SavePNG("heart.png")
err = picture.SaveAVIF("heart.avif")
```

See [docs/DRAWING_PATHS.md](docs/DRAWING_PATHS.md),
[docs/TRANSFORMS_AND_CLIPPING.md](docs/TRANSFORMS_AND_CLIPPING.md), and
[docs/IMAGE_EXPORT.md](docs/IMAGE_EXPORT.md).

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

## Tabs and selectable lists

Larger applications can group related pages and present scrollable choices
without leaving Rosaline's small composable API:

```go
themes := rosaline.List("Rosaline", "Lavender", "Ocean").
	OnSelect(func(index int, value string) {
		fmt.Println("selected", value)
	})

preferences := rosaline.Tabs(
	rosaline.Tab("Appearance", themes),
	rosaline.Tab("About", rosaline.Label("Built with Rosaline")),
).Expand()
```

See [docs/LISTS.md](docs/LISTS.md), [docs/TABS.md](docs/TABS.md), and the
complete [Preferences application](docs/PREFERENCES_APPLICATION.md).

## Tables made from ordinary Go data

Tables use normal strings and slices rather than a framework-specific model:

```go
files := rosaline.Table("Name", "Type", "Size").
	SetRows(
		[]string{"README.md", "Markdown", "8 KB"},
		[]string{"picture.png", "Image", "2.4 MB"},
	).
	OnActivate(func(row int, values []string) {
		fmt.Println("opened", values[0])
	})
```

Selection, keyboard activation, column sizing, dynamic replacement, and both
scrollbars are built in. See [docs/TABLES.md](docs/TABLES.md) and the complete
[File Browser application](docs/FILE_BROWSER.md).

## Nested data with simple trees

Trees use ordinary node pointers, labels, and optional application values:

```go
documents := rosaline.Node("Documents",
	rosaline.Node("Notes.txt"),
	rosaline.Node("Ideas.txt"),
).Expanded()

folders := rosaline.Tree(documents).
	OnActivate(func(node *rosaline.TreeNode) {
		fmt.Println("open", node.Value())
	})
```

Selection, activation, expansion callbacks, dynamic children, native keyboard
navigation, and both scrollbars are built in. See
[docs/TREES.md](docs/TREES.md) and the upgraded
[File Browser application](docs/FILE_BROWSER.md).

## Included in v0.8.0

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
- Reusable paths with straight, quadratic, and cubic Bézier sections
- Translate, rotate, scale, Push/Pop, and transformed clipping
- Canvas clicks, pointer movement, dragging, and button-release events
- Automatic and manually requested canvas redraws
- CGo-free loading and display of PNG, JPEG, GIF, BMP, TIFF, WebP, and AVIF
- Off-screen drawing with the same API as visible canvases
- PNG and CGo-free AVIF image export
- Horizontal and vertical scroll areas
- Native open, save, message, error, and confirmation dialogs
- Menu bars with working keyboard shortcuts
- Repeating and one-shot application timers
- Start, stop, restart, and running-state timer controls
- Frame-rate-based canvas animation
- Native tabbed interfaces with selection callbacks and programmatic selection
- Scrollable single-selection lists with selection and activation callbacks
- Dynamic list replacement and safe programmatic selection
- Focus traversal that automatically skips controls on hidden tab pages
- Native multi-column tables built from ordinary `[][]string` data
- Table selection, activation, column sizing, dynamic rows, and two-axis
  scrolling
- Native trees with nested nodes, labels, and application-defined values
- Tree selection, activation, expansion callbacks, dynamic roots and children,
  and two-axis scrolling
- Semantic colors and themes
- A saveable Paint application with menus, shortcuts, and PNG/AVIF output
- A complete Preferences application combining tabs, lists, form controls, and
  a live canvas preview
- A complete File Browser using lazy folder trees, tables, and Go's standard
  filesystem APIs
- Runnable hello, counter, canvas, forms, drawing, paint, image-viewer, and
  animation examples, plus the Preferences and File Browser applications
- Unit tests for non-visual core behavior

## Run the examples

From the extracted project root:

```bash
CGO_ENABLED=0 go run ./examples/hello
CGO_ENABLED=0 go run ./examples/counter
CGO_ENABLED=0 go run ./examples/canvas
CGO_ENABLED=0 go run ./examples/drawing
CGO_ENABLED=0 go run ./examples/forms
CGO_ENABLED=0 go run ./examples/paint
CGO_ENABLED=0 go run ./examples/imageviewer
CGO_ENABLED=0 go run ./examples/animation
CGO_ENABLED=0 go run ./examples/preferences
CGO_ENABLED=0 go run ./examples/filebrowser
```

## Project status

The API is experimental until v1.0. The next milestones add multiple windows,
radio buttons, general keyboard events, accessibility groundwork, and deeper
Linux display testing.

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

Dependencies keep their own licenses; see
[THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md).

Copyright (C) 2026 Britney Lozza and Rosaline contributors.
