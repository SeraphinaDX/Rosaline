# Dialogs

Rosaline provides small functions for messages, errors, confirmations, and
choosing files. The dialogs use the desktop's normal window system.

## Messages and errors

```go
rosaline.Message("Saved", "Your document was saved.")
rosaline.Error("Could not save", err.Error())
```

Both functions wait until the user dismisses the dialog.

## Asking for confirmation

`Confirm` returns a Boolean:

```go
if !rosaline.Confirm("Delete drawing?", "This cannot be undone.") {
	return
}

deleteDrawing()
```

The result is true only when the user chooses Yes.

## Opening a file

```go
path, ok := rosaline.OpenFileDialog(rosaline.FileDialogOptions{
	Title: "Open Image",
	Filters: []rosaline.FileFilter{
		{Name: "Images", Extensions: []string{"png", "jpg", "webp"}},
		{Name: "All files", Extensions: []string{"*"}},
	},
})
if !ok {
	return
}
```

`ok` is false when the user cancels. Checking it immediately keeps cancellation
separate from real file errors.

Extensions may be written as `png`, `.png`, or `*.png`; Rosaline normalizes
them. A filter may contain several extensions.

## Choosing a save location

```go
path, ok := rosaline.SaveFileDialog(rosaline.FileDialogOptions{
	Title:            "Save Drawing",
	InitialFile:      "drawing.png",
	DefaultExtension: ".png",
	Filters: []rosaline.FileFilter{
		{Name: "PNG images", Extensions: []string{"png"}},
	},
})
if !ok {
	return
}
```

Save dialogs confirm before returning the path of an existing file.

## File dialog options

Every field is optional:

- `Title` sets the dialog title.
- `InitialDirectory` chooses the first folder shown.
- `InitialFile` suggests a filename.
- `DefaultExtension` adds an extension when needed.
- `Filters` groups selectable file types.

The dialog only chooses a path. Reading or writing remains ordinary Go with
packages such as `os`, which keeps Rosaline applications idiomatic.

## Go concepts used here

- Multiple return values
- Boolean checks
- Struct literals
- Slices
- Error handling
