# Building the Rosaline Calculator

The Calculator is Rosaline v0.13.0's complete layout and presentation tutorial.
It combines a grid keypad, layered background, centered card, aligned dynamic
labels, flexible spacing, separators, sizing, keyboard input, shortcuts, a
custom theme, and tested application logic.

Run it from the project root:

```bash
CGO_ENABLED=0 go run ./examples/calculator
```

## Keep calculation separate from presentation

The `calculator` struct stores the current display, pending operator, left-hand
number, and friendly status text. Methods such as `digit`, `choose`, `equals`,
and `clear` change that ordinary Go state.

```go
type calculator struct {
	display    string
	expression string
	status     string
	left       float64
	operator   string
	replace    bool
	error      bool
}
```

This separation makes the arithmetic testable without opening a window. The
example includes tests for normal calculations, chained operations, decimal
editing, percentages, sign changes, backspace, division by zero, and reset.

## Build one reusable button grid

Small helpers connect button labels to model methods. The keypad is then one
flat list in reading order:

```go
keypad := rosaline.Grid(4,
	button("AC", model.clear),
	button("±", model.sign),
	button("%", model.percent),
	operation("÷", "÷"),
	// More rows...
).Gap(9).Expand()
```

Four columns describe the relationship directly. Rosaline calculates rows,
equal column widths, expansion, and cell placement.

## Make the display stable and readable

The number display is both dynamic and styled:

```go
rosaline.MinSize(
	rosaline.LabelFunc(func() string { return model.display }).
		FontSize(34).
		Bold().
		TextAlign(rosaline.AlignEnd),
	320, 60,
)
```

`MinSize` prevents short and long results from changing the panel's layout.
Right-aligned large bold text gives the display familiar calculator behavior.

## Use flexible space for status rows

The header and footer both use `Spring`:

```go
rosaline.Row(
	rosaline.Label("ROSALINE"),
	rosaline.Spring(),
	rosaline.Label("CALCULATOR"),
)
```

The labels remain at opposite edges regardless of panel width. This is more
adaptable than guessing a fixed spacer width.

## Layer and center the interface

The final screen is a canvas background with a card above it:

```go
rosaline.Stack(
	background,
	rosaline.Center(
		rosaline.Size(panel, 390, 570),
	),
).Expand()
```

The first Stack child paints the entire background. `Center` keeps the sized
card in the middle while leaving its surroundings visible. The alignment
metadata works without adding an opaque wrapper between the two layers.

## Share behavior between mouse and keyboard

Button callbacks call the calculator model directly. The window's key handler
uses the same methods for digits and printable operators. Named keys use
standalone shortcuts:

```go
Shortcuts: rosaline.Shortcuts(
	rosaline.Shortcut("Enter", model.equals),
	rosaline.Shortcut("Escape", model.clear),
	rosaline.Shortcut("Backspace", model.backspace),
),
```

This gives keyboard users complete operation without duplicating arithmetic
logic.

## Theme with semantic colors

The example changes `Background`, `Surface`, `Primary`, `Text`, `Muted`, and
`Border`. Cards, separators, labels, buttons, and layouts then use those
semantic colors automatically. The program does not style individual native
backend objects.

## Explore the complete source

The runnable application is in
[`examples/calculator/main.go`](../examples/calculator/main.go), with model
tests in [`examples/calculator/main_test.go`](../examples/calculator/main_test.go).
Read the focused [Layout and Presentation guide](LAYOUT_AND_PRESENTATION.md)
first when you want the smallest example of each individual feature.
