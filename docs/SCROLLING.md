# Scrolling

`Scroll` places any Rosaline widget inside a viewport with horizontal and
vertical scrollbars.

## Smallest example

```go
rosaline.Run(
	rosaline.Scroll(
		rosaline.Column(
			rosaline.Label("Line 1"),
			rosaline.Label("Line 2"),
			rosaline.Label("Line 3"),
		).Gap(12),
	).Size(400, 240),
)
```

`Size(width, height)` sets the preferred viewport size in pixels. The content
keeps its natural size and can extend beyond the viewport.

## Filling the window

Use `Expand` when the scroll area should grow and shrink with its layout:

```go
rosaline.Column(
	rosaline.Label("Document"),
	rosaline.Scroll(document).Expand(),
).Expand()
```

The outer layout and the scroll area both expand here. The fixed label keeps
its normal height while the scroll area receives the remaining space.

## Mouse and keyboard behavior

- The mouse wheel scrolls vertically.
- Shift+wheel scrolls horizontally.
- The visible scrollbars can be dragged or clicked.
- Scrollbars appear even when all content currently fits, keeping layout
  predictable as dynamic content changes size.

Linux wheel buttons and the mouse-wheel events used by Windows and macOS are
handled privately by Rosaline.

## Typical uses

Scroll areas are useful for large images, long forms, document previews,
tool palettes, and drawing surfaces. See [IMAGE_VIEWER.md](IMAGE_VIEWER.md) for
a complete large-image example.

## Go concepts used here

- Nested function calls
- Method chaining
- Layout composition
