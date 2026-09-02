# Images

Rosaline loads common image files with pure Go and displays them with the
`Image` widget. PNG, JPEG, GIF, BMP, TIFF, WebP, and AVIF are supported.

## Complete example

```go
package main

import "github.com/SeraphinaDX/Rosaline"

func main() {
	picture, err := rosaline.LoadImage("rose.png")
	if err != nil {
		panic(err)
	}

	rosaline.Run(
		rosaline.Column(
			rosaline.Label("My picture"),
			rosaline.Image(picture),
		),
	)
}
```

`LoadImage` returns a `Picture` and an error. The error includes the filename
and the underlying problem, such as a missing file or unsupported contents.
Applications should normally show it with `rosaline.Error` rather than panic.

## Picture information

Pictures expose useful metadata:

```go
width := picture.Width()
height := picture.Height()
format := picture.Format()
path := picture.Path()
pixels := picture.Image()
```

Width and height are measured in pixels. Format is a short name such as `png`
or `jpeg`. Path is the filename passed to `LoadImage`.
`Image()` returns the standard-library `image.Image` value for applications
that need to inspect or process the pixels.

## Changing the displayed image

An image widget may begin empty and receive a picture later:

```go
viewer := rosaline.Image(nil).
	Placeholder("Choose an image to begin.")

picture, err := rosaline.LoadImage(path)
if err != nil {
	rosaline.Error("Could not open image", err.Error())
	return
}
viewer.SetImage(picture)
```

Call `SetImage(nil)` to clear it. `Picture()` returns the currently displayed
picture.

## Using Go-generated images

`NewPicture` accepts Go's standard `image.Image` interface:

```go
pixels := image.NewRGBA(image.Rect(0, 0, 320, 200))
picture := rosaline.NewPicture(pixels)
widget := rosaline.Image(picture)
```

This makes Rosaline compatible with standard-library drawing code and other
pure-Go image packages.

## Saving pictures

Pictures can be written as PNG or AVIF:

```go
err := picture.SavePNG("rose.png")
err = picture.SaveAVIF("rose.avif")
```

See [IMAGE_EXPORT.md](IMAGE_EXPORT.md) for off-screen drawing, AVIF quality
options, save dialogs, and complete programs.

## Large images

Wrap an image in `Scroll` so it remains usable when it is larger than the
window:

```go
rosaline.Scroll(rosaline.Image(picture)).
	Size(700, 500).
	Expand()
```

See [SCROLLING.md](SCROLLING.md) and the complete
[IMAGE_VIEWER.md](IMAGE_VIEWER.md) tutorial.

## Go concepts used here

- Errors and early returns
- Struct methods
- Interfaces from the standard library
- Multiple return values
