# Rosaline

Rosaline is a small, beginner-friendly graphics and GUI library for Go. It is
designed for people who know a little Go and want to make a real graphical
program without first learning a large framework.

Rosaline is currently at `v0.13.0`. The public API is small on purpose and grows
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

## Layouts that stay readable

Build equal-cell grids, layered interfaces, centered cards, and adaptive rows
from ordinary widgets:

```go
panel := rosaline.Card(
	rosaline.Column(
		rosaline.Label("Quick actions").FontSize(24).Bold(),
		rosaline.Separator(),
		rosaline.Grid(2,
			rosaline.Button("Open", open),
			rosaline.Button("Save", save).Primary(),
		).Gap(8),
		rosaline.Row(
			rosaline.Label("Ready"),
			rosaline.Spring(),
			rosaline.Label("2 actions"),
		),
	).Gap(12),
)

rosaline.Run(
	rosaline.Stack(background, rosaline.Center(panel)).Expand(),
)
```

`Align`, `Size`, and `MinSize` cover positioning and sizing without manual
coordinates. See [Layout and Presentation](docs/LAYOUT_AND_PRESENTATION.md) and
the complete [Calculator application](docs/CALCULATOR_APPLICATION.md).

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

## Slow work without a frozen window

Background tasks do ordinary Go work in a goroutine while Rosaline safely
delivers progress and results to the GUI thread:

```go
progress := 0.0

task := rosaline.Background(func(ctx context.Context, report *rosaline.TaskReporter) error {
	for step := 1; step <= 100; step++ {
		if !report.Report(float64(step), "Working...") {
			return ctx.Err()
		}
	}
	return nil
}).OnProgress(func(update rosaline.TaskProgress) {
	progress = update.Percent
}).AutoStart()

rosaline.RunApp(rosaline.App{
	Tasks:   []*rosaline.Task{task},
	Content: rosaline.ProgressBar(&progress),
})
```

Tasks support cancellation, reusable starts, posted result callbacks, normal Go
errors, and automatic window-lifetime cleanup. See
[Background Tasks](docs/BACKGROUND_TASKS.md) and the complete
[Background Bloom application](docs/BACKGROUND_BLOOM_APPLICATION.md).

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

## Multiple windows without another event loop

Secondary windows are reusable Go handles:

```go
var about *rosaline.Window
about = rosaline.NewWindow(rosaline.WindowOptions{
	Title:  "About",
	Parent: rosaline.MainWindow(),
	Content: rosaline.Button("Close", func() {
		about.Close()
	}),
})

openAbout := rosaline.Button("About", func() {
	about.Show()
})
```

Calling `Show` twice focuses the existing window. Windows can close and reopen,
share normal Go state, own menus and timers, and form safe parent-child
relationships. See [docs/MULTIPLE_WINDOWS.md](docs/MULTIPLE_WINDOWS.md) and the
complete [Project Desk example](examples/windows/main.go).

## Everyday controls with ordinary Go values

Radio groups, combo boxes, sliders, and progress bars bind directly to normal
Go variables:

```go
priority := "normal"
category := "Documentation"
completion := 35.0

rosaline.Run(
	rosaline.Column(
		rosaline.ComboBox(&category, "Documentation", "Development"),
		rosaline.RadioGroup(&priority,
			rosaline.Choice("Low", "low"),
			rosaline.Choice("Normal", "normal"),
			rosaline.Choice("High", "high"),
		).Horizontal(),
		rosaline.Slider(&completion, 0, 100).Step(5),
		rosaline.ProgressBar(&completion),
	)
)
```

The slider and progress bar share one pointer, so they stay synchronized after
Rosaline events without a binding language. Options and choices can be replaced
while the application is running. Progress bars also support an indeterminate
busy mode. See [Radio Groups](docs/RADIO_GROUPS.md),
[Combo Boxes](docs/COMBO_BOXES.md), [Sliders](docs/SLIDERS.md),
[Progress Bars](docs/PROGRESS_BARS.md), and the complete
[Task Settings application](docs/TASK_SETTINGS_APPLICATION.md).

## Keyboard input without platform code

Windows and canvases receive the same small `KeyEvent`, and applications can
define shortcuts without building a menu:

```go
canvas := rosaline.Canvas(draw).
	Focus().
	OnKeyDown(func(event rosaline.KeyEvent) {
		if event.Is(rosaline.KeyRight) {
			x += 10
		}
	})

rosaline.RunApp(rosaline.App{
	Shortcuts: rosaline.Shortcuts(
		rosaline.Shortcut("Primary+S", save),
		rosaline.Shortcut("F1", showHelp),
	),
	Content: canvas,
})
```

`Primary` follows the platform convention: Control on Linux and Windows,
Command on macOS. Canvas key callbacks redraw automatically, and window key
handlers observe input without breaking normal text editing. See
[Keyboard Input and Shortcuts](docs/KEYBOARD_INPUT.md) and the complete
[Keyboard Garden application](docs/KEYBOARD_GARDEN_APPLICATION.md).

## Included in v0.13.0

- Application windows
- Labels and dynamic labels with font size, bold, and text alignment
- Buttons and message dialogs
- Rows, columns, spacing, padding, and expansion
- Equal-column grids with automatic rows, gaps, padding, and expansion
- Layered stacks with metadata-driven alignment and centered overlays
- Adaptive springs, themed separators and cards, and preferred or minimum
  sizing wrappers
- Simple state values
- Single-line and multiline text input
- Password display, placeholders, change events, and Enter submission
- Checkboxes bound to Go Boolean variables
- Radio groups with separate labels and values, vertical or horizontal layout,
  callbacks, programmatic selection, and dynamic choices
- Read-only combo boxes with callbacks, programmatic selection, width control,
  and dynamic options
- Numeric sliders with safe ranges, optional steps, horizontal or vertical
  direction, focus, callbacks, and programmatic control
- Determinate and indeterminate progress bars with custom maximums, size,
  direction, and start/stop controls
- Tab and Shift+Tab keyboard navigation
- Backend-neutral window and canvas key-down and key-up events
- Friendly named key constants, printable text, and modifier fields
- Menu-free application shortcuts with cross-platform `Primary` modifiers
- Keyboard-enabled canvases with initial focus, click-to-focus, a visible focus
  ring, and automatic redraw
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
- Window-owned background tasks with standard Go context cancellation
- GUI-thread progress, completion, and posted result callbacks
- Safe task restart, auto-start, panic conversion, and late-callback cleanup
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
- Reusable secondary windows with simple show, close, focus, title, and state
  controls
- Window-specific content, menus, shortcuts, themes, focus traversal, and
  timers
- Parent-child window lifecycles, automatic parent opening, cascading closure,
  and duplicate prevention
- Automatic dynamic-widget refresh across every open window
- Semantic colors and themes
- A saveable Paint application with menus, shortcuts, and PNG/AVIF output
- A complete Preferences application combining tabs, lists, form controls, and
  a live canvas preview
- A complete File Browser using lazy folder trees, tables, and Go's standard
  filesystem APIs
- A complete Project Desk application combining an editor, live child preview,
  About window, shared state, menus, and dynamic titles
- A complete Task Settings application combining everyday controls, shared Go
  state, validation, dynamic choices, and both progress modes
- A complete Keyboard Garden combining canvas input, modifiers, releases,
  standalone shortcuts, drawing, dialogs, and PNG export
- A complete Background Bloom combining responsive image generation, progress,
  cancellation, result posting, shortcuts, dialogs, and PNG export
- A complete Calculator combining grids, stacks, cards, alignment, dynamic
  typography, adaptive spacing, keyboard input, shortcuts, and tested logic
- Runnable hello, counter, canvas, forms, drawing, paint, image-viewer, and
  animation examples, plus the Preferences, File Browser, and Project Desk
  applications, Task Settings, Keyboard Garden, Background Bloom, and the
  Calculator
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
CGO_ENABLED=0 go run ./examples/windows
CGO_ENABLED=0 go run ./examples/tasksettings
CGO_ENABLED=0 go run ./examples/keyboard
CGO_ENABLED=0 go run ./examples/background
CGO_ENABLED=0 go run ./examples/calculator
```

## Project status

The API is experimental until v1.0. The next milestones add accessibility
groundwork, custom widgets, and deeper Linux display testing.

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
