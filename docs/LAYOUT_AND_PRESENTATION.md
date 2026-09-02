# Layout and Presentation

Rosaline layouts arrange ordinary widgets without coordinates. Rows and
columns remain the simplest tools; grids, stacks, alignment, flexible space,
cards, separators, sizing, and label typography handle more polished screens.

## Rows and columns

`Row` places children from left to right. `Column` places them from top to
bottom:

```go
rosaline.Column(
	rosaline.Label("Account"),
	rosaline.TextBox(&name),
	rosaline.Row(
		rosaline.Button("Cancel", cancel),
		rosaline.Button("Save", save).Primary(),
	),
).Gap(12).Padding(8)
```

`Gap` changes the space between children. `Padding` changes the space around
the inside edge. `Expand` asks the whole row or column to use available room.

Rows and columns now understand horizontal and vertical expansion separately.
A horizontal separator can fill a column's width without consuming its spare
height, and a flexible spring grows only along the row or column that owns it.

## Automatic grids

`Grid` takes a column count followed by ordinary widgets:

```go
rosaline.Grid(3,
	rosaline.Button("One", one),
	rosaline.Button("Two", two),
	rosaline.Button("Three", three),
	rosaline.Button("Four", four),
	rosaline.Button("Five", five),
	rosaline.Button("Six", six),
).Gap(8)
```

Children fill a row from left to right, then continue on the next row. Columns
have equal widths. A non-positive column count safely becomes one, and nil
children are ignored.

```go
grid.Padding(12).Expand()
```

`Padding` adds inside space. `Expand` lets the grid use extra room and gives its
rows equal shares of extra height. This makes button pads and dashboards easy
without calculating individual cell coordinates.

## Layering with Stack

`Stack` puts every child in the same area. Later children appear above earlier
children:

```go
rosaline.Stack(
	backgroundCanvas,
	rosaline.Center(
		rosaline.Card(rosaline.Label("Centered overlay")),
	),
).Expand()
```

The centered card keeps its natural size instead of stretching over the
background. Stack is useful for canvas overlays, empty-state messages, badges,
and controls floating above artwork.

Most opaque widgets cover the space they occupy. Use `Align` or `Center` when
an overlay should retain its natural size and reveal the lower layer around it.

## Alignment

`Align` positions one widget on two axes:

```go
rosaline.Align(
	rosaline.Button("Continue", next),
	rosaline.AlignEnd,
	rosaline.AlignCenter,
)
```

The first alignment is horizontal and the second is vertical. Each accepts:

- `AlignStart` for left or top
- `AlignCenter` for the center
- `AlignEnd` for right or bottom
- `AlignStretch` to fill that axis

`Center(widget)` is the friendly shorthand for centering both axes:

```go
rosaline.Center(rosaline.Label("Perfectly centered"))
```

Alignment works consistently as the window content or inside a row, column,
grid, or stack.

## Fixed and flexible space

`Spacer` creates a known number of empty pixels:

```go
rosaline.Spacer(20, 10)
```

`Spring` absorbs spare space:

```go
rosaline.Row(
	rosaline.Label("Document.txt"),
	rosaline.Spring(),
	rosaline.Button("Save", save).Primary(),
)
```

That row keeps the filename at the left and the Save button at the right. In a
column, the same spring pushes later content toward the bottom.

## Separators

`Separator` creates a horizontal one-pixel line using `Theme.Border`:

```go
rosaline.Separator()
```

It has three simple options:

```go
rosaline.Separator().Thickness(2)
rosaline.Separator().Vertical()
rosaline.Separator().Vertical().Horizontal()
```

Use horizontal separators in columns and vertical separators in rows.

## Cards

`Card` places content on `Theme.Surface`, draws a `Theme.Border`, and adds 14
pixels of padding:

```go
rosaline.Card(
	rosaline.Column(
		rosaline.Label("Storage").Bold(),
		rosaline.Label("42 GB available"),
	),
)
```

Customize its padding or let it grow:

```go
card.Padding(20).Expand()
```

Cards use semantic theme colors, so application layouts do not need hardcoded
light and dark variants.

## Preferred and minimum sizes

`Size` gives content a preferred pixel width and height:

```go
rosaline.Size(preview, 480, 320)
```

The child fills that area. Add `Expand` when the area may grow with its parent:

```go
rosaline.Size(preview, 480, 320).Expand()
```

`MinSize` preserves a child's natural requested size while requiring at least
the supplied dimensions:

```go
rosaline.MinSize(statusLabel, 300, 48)
```

This is useful for dynamic labels whose text should not make a window or panel
jump between sizes. Non-positive dimensions safely become one pixel.

## Label presentation

Fixed and dynamic labels share the same presentation methods:

```go
rosaline.LabelFunc(currentTotal).
	FontSize(32).
	Bold().
	TextAlign(rosaline.AlignEnd).
	Color(rosaline.Rose)
```

- `FontSize` uses pixels; a non-positive size restores the platform default.
- `Bold` enables bold text.
- `TextAlign` accepts `AlignStart`, `AlignCenter`, or `AlignEnd`.
- `Color` uses any Rosaline `Color`.

`TextAlign` controls text inside the label. `Align` controls where the entire
label widget sits in its parent; they solve different layout problems.

## A complete presentation sample

```go
package main

import "github.com/SeraphinaDX/Rosaline"

func main() {
	rosaline.RunApp(rosaline.App{
		Title:  "Layout Sample",
		Width:  620,
		Height: 420,
		Content: rosaline.Stack(
			rosaline.Canvas(func(c *rosaline.DrawingCanvas) {
				c.Clear(rosaline.Hex("#fff3fa"))
				c.FillCircle(80, 80, 70, rosaline.SoftRose)
			}).Size(600, 400).Expand(),
			rosaline.Center(
				rosaline.Card(
					rosaline.Column(
						rosaline.Label("Choose a flower").FontSize(24).Bold(),
						rosaline.Separator(),
						rosaline.Grid(2,
							rosaline.Button("Rose", func() {}),
							rosaline.Button("Violet", func() {}),
							rosaline.Button("Daisy", func() {}),
							rosaline.Button("Peony", func() {}).Primary(),
						).Gap(8),
						rosaline.Row(
							rosaline.Label("Four choices"),
							rosaline.Spring(),
							rosaline.Label("Ready"),
						),
					).Gap(12),
				).Padding(20),
			),
		).Expand(),
	})
}
```

## Common mistakes

- `Grid(4, ...)` means four columns, not four rows.
- Put the background first in a `Stack`; later widgets appear above it.
- Use `Spring`, not a very large `Spacer`, when space should adapt as the
  window changes size.
- Use `MinSize` when dynamic text may grow; use `Size` when the containing area
  should begin at a specific size.
- `TextAlign` aligns text inside a label. Use `Align` to position a complete
  widget.

## Go concepts used here

- variadic function arguments
- method chaining
- semantic constants
- composing values by nesting function calls
- callbacks shared between multiple controls

See the [Calculator application](CALCULATOR_APPLICATION.md) for these features
combined into a polished keyboard-friendly program.
