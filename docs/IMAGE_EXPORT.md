# Off-Screen Drawing and Image Export

Rosaline uses the same pure-Go drawing engine for visible canvases and exported
images. A saved picture therefore matches the canvas rather than relying on a
screenshot.

PNG and AVIF export work with `CGO_ENABLED=0` and need no external image tool.

## Save a visible canvas

This complete program saves its canvas when the button is clicked:

```go
package main

import "github.com/SeraphinaDX/Rosaline"

func main() {
	canvas := rosaline.Canvas(func(c *rosaline.DrawingCanvas) {
		c.Clear(rosaline.White)
		c.FillCircle(160, 100, 60, rosaline.Rose)
	}).Size(320, 200)

	rosaline.Run(
		rosaline.Column(
			canvas,
			rosaline.Button("Save PNG", func() {
				if err := canvas.Picture().SavePNG("rose.png"); err != nil {
					rosaline.Error("Could not save", err.Error())
				}
			}).Primary(),
		),
	)
}
```

`canvas.Picture()` renders the current drawing state into an off-screen
`Picture`. It works even before the canvas is placed in a window.

## Save AVIF

Replace the saving line with:

```go
err := canvas.Picture().SaveAVIF("rose.avif")
```

The default is quality 90 with a fast encoding setting. AVIF usually produces
smaller files than PNG for detailed artwork and photographs, but it takes more
CPU time to encode.

Advanced applications can provide options:

```go
err := picture.SaveAVIF("rose.avif", rosaline.AVIFOptions{
	Quality:  95,
	Speed:    6,
	Lossless: false,
})
```

- `Quality` ranges from 1 to 100.
- `Speed` ranges from 0 to 10; larger values encode faster.
- `Lossless` preserves pixels exactly and ignores quality.

The first build after AVIF support is downloaded can take longer because the
CGo-free codec is compiled and cached. Later builds are normally much faster.

## Off-screen drawing without a window

`Render` creates a picture directly:

```go
picture := rosaline.Render(800, 600, func(c *rosaline.DrawingCanvas) {
	c.Clear(rosaline.White)
	c.FillRect(40, 40, 300, 180, rosaline.Rose)
})

if err := picture.SavePNG("card.png"); err != nil {
	panic(err)
}
```

This is useful for thumbnail generators, chart exporters, tests, batch image
tools, and applications that need a full-resolution image different from their
preview size.

## Loading AVIF

`LoadImage` now recognizes AVIF alongside PNG, JPEG, GIF, BMP, TIFF, and WebP:

```go
picture, err := rosaline.LoadImage("painting.avif")
```

The normal `Image` widget can display the result.

## Choosing a filename

Use `SaveFileDialog` in interactive applications:

```go
path, ok := rosaline.SaveFileDialog(rosaline.FileDialogOptions{
	Title:            "Save Painting",
	InitialFile:      "painting.avif",
	DefaultExtension: ".avif",
	Filters: []rosaline.FileFilter{
		{Name: "AVIF images", Extensions: []string{"avif"}},
	},
})
if ok {
	err := canvas.Picture().SaveAVIF(path)
	// Handle err here.
}
```

## Common mistakes

### Ignoring save errors

Directories may be unwritable or disks may be full. Always check the returned
error and show it to the user.

### Saving the window instead of the canvas

Use `canvas.Picture()`. It exports the drawing at the canvas's exact pixel size
without including buttons, menus, or window decorations.

### Expecting AVIF to encode as quickly as PNG

AVIF performs substantially more compression work. Use a larger `Speed` value
when fast interactive saves matter more than the smallest possible file.

## Go concepts used here

- Error handling
- Method calls
- Optional variadic arguments
- Struct configuration
- Reusing callback functions without a window
