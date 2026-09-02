# Building a Saveable Paint Application

The v0.5 Paint example combines canvas input, paths, menus, file dialogs,
confirmation dialogs, and PNG/AVIF export into a real application.

Run it from the source root:

```bash
CGO_ENABLED=0 go run ./examples/paint
```

Its complete source is in [`examples/paint/main.go`](../examples/paint/main.go).

## Application state

Paint keeps the artwork as Go data:

```go
type point struct {
	x float64
	y float64
}

type stroke struct {
	points []point
}
```

The application stores a slice of strokes. This matters because the canvas can
redraw the complete picture at any time and export the same data off-screen.

## Turning strokes into paths

Each stored stroke becomes a path while drawing:

```go
path := rosaline.NewPath().MoveTo(points[0].x, points[0].y)
for _, point := range points[1:] {
	path.LineTo(point.x, point.y)
}
c.StrokePath(path, 6, rosaline.Rose)
```

The mouse callbacks only collect data. The drawing callback decides how that
data looks.

## One save function per format

The PNG and AVIF commands each:

1. Ask for a filename with `SaveFileDialog`.
2. Return immediately if the user cancels.
3. Render `canvas.Picture()`.
4. Save it with `SavePNG` or `SaveAVIF`.
5. Show a useful error or success message.

Keeping these steps in functions allows both menu items and buttons to call the
same code:

```go
rosaline.MenuItem("Save as PNG…", savePNG).Shortcut("Ctrl+S")
rosaline.Button("Save PNG", savePNG)
```

This avoids duplicated behavior.

## Protecting unsaved work

The `dirty` Boolean becomes true when a new stroke starts and false after a
successful save. Before clearing, Paint asks for confirmation when unsaved
work exists.

This is a small but important application pattern: remember whether user data
has changed, and protect it before destructive actions.

## Features working together

Paint demonstrates:

- Mouse down, drag, and release events
- Slices and custom structs
- Paths and stroke rendering
- Dynamic status labels
- Menus and keyboard shortcuts
- Save dialogs and file filters
- PNG and AVIF output
- Confirmation and error dialogs

No special project architecture is required. The program remains one readable
Go file while still behaving like a genuine desktop utility.

## Next improvements to try

Good exercises for extending Paint are:

- Add several brush colors.
- Store a width in each stroke.
- Add Undo by removing the final stroke.
- Add a line or rectangle tool.
- Export at twice the canvas resolution with `Render` and `Scale`.

## Go concepts used here

- Structs and slices
- Closures
- Shared callback functions
- Boolean state
- Early returns
- Error handling
- Separation of state and drawing
