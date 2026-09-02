# Building an Image Viewer

The image-viewer example combines Rosaline's image loading, dynamic widgets,
scrolling, menus, keyboard shortcuts, open/save dialogs, confirmation dialogs,
and error messages in one application.

Run it from the extracted Rosaline source root:

```bash
CGO_ENABLED=0 go run ./examples/imageviewer
```

## Application state

The viewer keeps three pieces of state:

```go
viewer := rosaline.Image(nil)
status := "No image loaded"
currentPath := ""
```

The image widget changes with `SetImage`. The status appears in a `LabelFunc`,
and `currentPath` remembers which file can be copied.

## Opening and validating an image

The Open action first asks for a path, then decodes the file:

```go
path, ok := rosaline.OpenFileDialog(options)
if !ok {
	return
}

picture, err := rosaline.LoadImage(path)
if err != nil {
	rosaline.Error("Could not open image", err.Error())
	return
}

viewer.SetImage(picture)
currentPath = path
```

Cancellation is not an error. A selected but invalid image is, so those cases
are handled separately.

## Saving a copy

The save dialog chooses a destination. The example then uses `os.ReadFile` and
`os.WriteFile` to copy the original bytes. Rosaline handles user interaction;
the Go standard library handles file data.

## Closing safely

The Close action uses `Confirm` before clearing the widget:

```go
if rosaline.Confirm("Close image?", "Remove the current image?") {
	viewer.SetImage(nil)
	currentPath = ""
}
```

## Keeping large images usable

The image widget is placed inside a scroll area:

```go
rosaline.Scroll(viewer).Size(820, 520).Expand()
```

Images retain their real pixel size. The viewport provides both scroll
directions rather than shrinking or cropping the image silently.

## Menus and shared callbacks

The same `openImage`, `saveCopy`, and `closeImage` functions are passed to menu
items. `Shortcut` connects Ctrl+O, Ctrl+Shift+S, and Ctrl+W without a separate
keyboard event system in the application.

The full source is in
[`examples/imageviewer/main.go`](../examples/imageviewer/main.go).

## Try changing it

- Add Previous and Next actions for images in the same folder.
- Display the decoded format in the status label.
- Remember the last-used directory in a normal Go variable.
- Add a toolbar with buttons that reuse the menu callbacks.

This example is the groundwork for image editors, sprite tools, document
viewers, and asset browsers.
