# Rosaline

Rosaline is a small, beginner-friendly graphics and GUI library for Go. It is
designed for people who know a little Go and want to make a real graphical
program without first learning a large framework.

Rosaline is currently an early `v0.1.0` foundation. The public API is small on
purpose and will grow through well-documented, tested features.

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

Rosaline requires Go 1.24 or newer. No C compiler or separately installed GUI
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

## Included in v0.1.0

- Application windows
- Labels and dynamic labels
- Buttons and message dialogs
- Rows, columns, spacing, padding, and expansion
- Simple state values
- A first-class canvas with lines, rectangles, circles, and text
- Semantic colors and themes
- Runnable hello, counter, and canvas examples
- Unit tests for non-visual core behavior

## Run the examples

From the extracted project root:

```bash
CGO_ENABLED=0 go run ./examples/hello
CGO_ENABLED=0 go run ./examples/counter
CGO_ENABLED=0 go run ./examples/canvas
```

## Project status

The API is experimental until v1.0. The next milestones add text input,
checkboxes, images, keyboard and mouse events, canvas redraw/input, scrolling,
menus, and file dialogs.

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
