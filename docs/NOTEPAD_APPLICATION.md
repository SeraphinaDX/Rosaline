# Building the Rosaline Notepad

The Notepad is Rosaline v0.14.0's complete text-editing tutorial. It combines
an expanding editor, files, menus, shortcuts, undo and redo, the clipboard,
find and replace, a secondary tool window, saved-state tracking, dynamic window
titles, cursor position, document statistics, and safe closing.

Run it from the project root:

```bash
CGO_ENABLED=0 go run ./examples/notepad
```

## Keep document data ordinary

The application keeps only the text and current path in its document model:

```go
type document struct {
	text string
	path string
}
```

The path is empty for an untitled document. A small `name` method turns that
into either `Untitled` or the final part of the saved path. Rosaline does not
require a document interface or framework-specific file model.

## Let the editor own editing details

The document text is passed by pointer and the editor fills its card:

```go
editor = rosaline.TextArea(&document.text).
	Size(88, 30).
	Expand().
	Focus()
```

The application reads `editor.Text()` when saving. Undo history, selection,
clipboard integration, cursor indices, and the native scrollbar stay inside
the widget.

## Use the standard library for files

Opening is normal Go:

```go
contents, err := os.ReadFile(path)
if err != nil {
	rosaline.Error("Could not open file", err.Error())
	return
}
editor.SetText(string(contents))
editor.MarkSaved()
```

Saving is equally direct:

```go
err := os.WriteFile(document.path, []byte(editor.Text()), 0o644)
```

Rosaline's file dialogs choose a path. The standard library owns the actual
file format and I/O, which keeps the skills learned here useful in every Go
program.

## Share commands between menus and logic

The Edit menu passes editor methods directly:

```go
rosaline.Menu("Edit",
	rosaline.MenuItem("Undo", editor.Undo).Shortcut("Primary+Z"),
	rosaline.MenuItem("Redo", editor.Redo).Shortcut("Primary+Shift+Z"),
	rosaline.MenuItem("Cut", editor.Cut).Shortcut("Primary+X"),
	rosaline.MenuItem("Copy", editor.Copy).Shortcut("Primary+C"),
	rosaline.MenuItem("Paste", editor.Paste).Shortcut("Primary+V"),
)
```

Mouse menu choices and keyboard shortcuts therefore use exactly the same
behavior. `Primary` means Control on Linux and Windows and Command on macOS.

## Put find tools in a reusable window

The find window contains two text boxes and four buttons. `FindNext` selects a
match inside the main editor; `ReplaceSelection` changes that match;
`ReplaceAll` returns a count for the status bar.

Calling `findWindow.Show()` repeatedly focuses the existing window instead of
creating duplicates. Its `Parent` is `rosaline.MainWindow()`, so it follows the
Notepad's lifetime without creating another event loop.

## Treat saved text as a checkpoint

`editor.Modified()` drives both the title's `*` marker and close protection.
After a successful save, `MarkSaved` moves the checkpoint to the current text.
This is safer than maintaining a Boolean that can become incorrect after undo.

The same `canDiscard` function protects New, Open, and Quit:

```go
canDiscard := func() bool {
	if !editor.Modified() {
		return true
	}
	return rosaline.Confirm("Unsaved changes", "Discard this document?")
}
```

The application also gives it to `App.OnCloseRequest`, so the desktop close
button follows the same rule.

## Build a live status bar

The bottom row combines normal Go string processing with `Cursor`:

```go
words := len(strings.Fields(editor.Text()))
position := editor.Cursor()
status := fmt.Sprintf("%d words · Ln %d, Col %d",
	words, position.Line, position.Column+1)
```

`LabelFunc` reevaluates this after Rosaline callbacks. No observable-property
or binding language is needed.

## Test logic without a window

The example tests document naming, modified-title formatting, word counts, and
line counts without opening a GUI. Rosaline's own tests cover the text area's
unmounted document operations and close-request behavior.

Keeping these calculations in small functions makes them fast and dependable
while leaving native widget behavior to focused integration testing.

## Explore the complete source

The runnable application is in
[`examples/notepad/main.go`](../examples/notepad/main.go), with model tests in
[`examples/notepad/main_test.go`](../examples/notepad/main_test.go). Read the
focused [Text Editing guide](TEXT_EDITING.md) first for smaller examples of
each feature.
