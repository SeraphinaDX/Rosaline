# Text Editing

`TextArea` can grow from a small notes field into a complete document editor.
Its public methods use ordinary strings and callbacks; the native undo stack,
selection, clipboard, scrolling, and cursor indices remain private.

## Smallest editable document

```go
package main

import "github.com/SeraphinaDX/Rosaline"

func main() {
	text := "Start writing here."
	editor := rosaline.TextArea(&text).Expand().Focus()

	rosaline.RunApp(rosaline.App{
		Title:   "My Editor",
		Width:   700,
		Height:  500,
		Content: editor,
	})
}
```

`Expand` lets the control fill the window in both directions. A vertical
scrollbar is included automatically. The `text` variable always contains the
latest value after a Rosaline callback.

## Read and change text

The pointer remains the easiest way to share text with other controls. The
widget also provides explicit operations for editor commands:

```go
editor.Text()
editor.SetText("A new document")
editor.Append("\nAnother line")
editor.Clear()
```

`SetText` starts a fresh undo history, which makes it appropriate when opening
a different document. `Append` and mounted `Clear` are normal undoable edits.
All methods are safe before the window opens and after it closes.

## Track saved and unsaved text

A new text area treats its initial value as saved:

```go
if editor.Modified() {
	status = "Unsaved changes"
}
```

After writing a file successfully, mark that exact text as saved:

```go
if err := os.WriteFile(path, []byte(editor.Text()), 0o644); err != nil {
	rosaline.Error("Could not save", err.Error())
	return
}
editor.MarkSaved()
```

`Modified` compares the current text with the value recorded by `MarkSaved`.
Undoing all the way back to the saved content therefore clears the modified
state naturally.

When opening a file, replace the text and then mark it saved:

```go
editor.SetText(string(contents))
editor.MarkSaved()
```

## Undo, clipboard, and selection

Editor commands are ordinary methods:

```go
editor.Undo()
editor.Redo()
editor.Cut()
editor.Copy()
editor.Paste()
editor.SelectAll()
```

They operate on the mounted text area and safely do nothing before it opens.
Cut, copy, and paste use the desktop clipboard. Undo and redo safely do nothing
when their history is empty.

The same methods work directly as menu callbacks:

```go
rosaline.Menu("Edit",
	rosaline.MenuItem("Undo", editor.Undo).Shortcut("Primary+Z"),
	rosaline.MenuItem("Copy", editor.Copy).Shortcut("Primary+C"),
	rosaline.MenuItem("Paste", editor.Paste).Shortcut("Primary+V"),
)
```

## Find and replace

`FindNext` selects the next exact match and wraps to the beginning:

```go
if !editor.FindNext(search) {
	status = "No match"
}
```

The selected match can be replaced, or every match can be changed at once:

```go
editor.ReplaceSelection(replacement)
count := editor.ReplaceAll(search, replacement)
```

An empty search never matches. Searches are case-sensitive so their behavior
is direct and predictable.

## Cursor position

`Cursor` returns a small Go value:

```go
position := editor.Cursor()
fmt.Printf("line %d, column %d\n", position.Line, position.Column+1)
```

Lines begin at 1. Columns begin at 0 internally, so an editor status bar often
adds one when presenting the column to a person. `OnCursorMove` is available
when an application needs a callback rather than querying the value in a
`LabelFunc`.

## Protect unsaved work

`OnCloseRequest` runs before either a menu command or the window's close button
ends the application:

```go
rosaline.RunApp(rosaline.App{
	OnCloseRequest: func() bool {
		return !editor.Modified() || rosaline.Confirm(
			"Unsaved changes",
			"Discard this document?",
		)
	},
	Content: editor,
})
```

Return `true` to close or `false` to keep the window open. Secondary windows
support the same option in `WindowOptions`.

## Common mistakes

- Call `MarkSaved` only after a file was written successfully.
- Use `SetText` followed by `MarkSaved` when opening another document.
- Do not write files inside `OnChange`; it runs for every editing change.
- `ReplaceSelection` returns false when nothing is selected.
- Keep slow file processing in a background `Task`, but ordinary small text
  files can be read and written directly from menu callbacks.

## Go concepts used here

- Strings and pointers
- Struct values
- Boolean return values
- Methods used as callbacks
- Reading and writing files with `os`

See [Building the Rosaline Notepad](NOTEPAD_APPLICATION.md) for a complete
application using every feature together.
