# Building the Keyboard Garden

Keyboard Garden is Rosaline v0.11.0's complete keyboard-input tutorial. It
combines focused canvas events, key modifiers, releases, application
shortcuts, drawing, dynamic labels, dialogs, and PNG export.

Run it from the Rosaline source tree:

```bash
CGO_ENABLED=0 go run ./examples/keyboard
```

The complete source is in
[`examples/keyboard/main.go`](../examples/keyboard/main.go).

## Keep the garden in ordinary Go values

The cursor and planted flowers are plain application data:

```go
type flower struct {
	x float64
	y float64
}

x, y := 350.0, 210.0
flowers := make([]flower, 0)
```

Keyboard callbacks change these values. The canvas drawing function only reads
them and recreates the current picture. This is the same state-and-redraw
pattern used by Rosaline's mouse-driven Paint example.

## Give the canvas focus

```go
canvas := rosaline.Canvas(draw).
	Size(700, 400).
	Expand().
	Focus()
```

`Focus` makes the garden ready for arrow keys as soon as the window opens. The
canvas also receives focus when clicked and displays a theme-colored focus
ring, so users can tell where keyboard input will go.

## Handle navigation and printable keys together

The key-down callback accepts arrows and WASD:

```go
switch event.Key {
case rosaline.KeyLeft, rosaline.Key("a"):
	x -= step
case rosaline.KeyRight, rosaline.Key("d"):
	x += step
case rosaline.KeyUp, rosaline.Key("w"):
	y -= step
case rosaline.KeyDown, rosaline.Key("s"):
	y += step
}
```

Both forms produce the same movement without converting backend key names.
Operating systems normally repeat key-down events while a key is held, so the
cursor continues moving.

The example first ignores letter movement while a command modifier is held:

```go
if event.Primary || event.Control || event.Alt {
	return
}
```

That lets `Primary+S` reach the Save shortcut without also treating its `S` as
movement. This distinction matters when canvas controls and window shortcuts
share printable keys.

## Use modifiers for a variation

Shift changes the movement distance:

```go
step := 12.0
if event.Shift {
	step = 30
}
```

This is easier to understand than separately registering every possible
Shift-and-arrow combination. `Control`, `Alt`, and `Primary` are available for
the same pattern.

## Let action keys modify data

Space records the cursor position and Backspace removes the newest flower:

```go
case rosaline.KeySpace:
	flowers = append(flowers, flower{x: x, y: y})
case rosaline.KeyBackspace, rosaline.KeyDelete:
	if len(flowers) > 0 {
		flowers = flowers[:len(flowers)-1]
	}
```

Rosaline redraws the canvas automatically when the callback returns. The
drawing function loops over the updated slice and renders every flower.

## Observe releases

The example uses `OnKeyUp` to update its status line:

```go
canvas.OnKeyUp(func(event rosaline.KeyEvent) {
	lastKey = event.Key.String() + " released"
})
```

Games can instead set and clear Boolean movement flags on down and up events.
Utilities often only need `OnKeyDown`.

## Share commands between buttons and shortcuts

Reset, Save, and Help are named functions because both a button and a shortcut
use each one:

```go
Shortcuts: rosaline.Shortcuts(
	rosaline.Shortcut("Primary+R", reset),
	rosaline.Shortcut("Primary+S", save),
	rosaline.Shortcut("F1", help),
)
```

The visible buttons call the same `reset`, `save`, and `help` functions. There
is only one implementation of each action to test and maintain.

## Export exactly what the canvas draws

The Save command chooses a path and renders the current drawing off-screen:

```go
path, ok := rosaline.SaveFileDialog(options)
if !ok {
	return
}
if err := canvas.Picture().SavePNG(path); err != nil {
	rosaline.Error("Could not save garden", err.Error())
}
```

The focus ring belongs to the interactive widget rather than its drawing, so
it does not appear in the exported PNG.

## Features to try next

- Add different flowers for number keys 1, 2, and 3.
- Add a `ComboBox` that chooses the current flower color.
- Track held arrow keys and use an animation timer for frame-based movement.
- Add an Undo shortcut that removes the latest flower.
- Open the help text in a secondary window instead of a message dialog.

## Go concepts used here

- structs and slices
- `switch` statements
- named constants
- callbacks sharing application state
- clamping numeric values
- error handling
- reusable functions
