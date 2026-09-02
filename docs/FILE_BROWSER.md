# Building a File Browser

The File Browser is Rosaline v0.7.0's complete table application. It combines a
dynamic table, filesystem data, an address box, buttons, menus, shortcuts,
status text, folder navigation, and error dialogs.

Run it from the Rosaline source tree:

```bash
CGO_ENABLED=0 go run ./examples/filebrowser
```

The complete source is in
[`examples/filebrowser/main.go`](../examples/filebrowser/main.go).

## Keep application data separate from displayed text

Each filesystem entry has values the table displays and values the application
needs for navigation:

```go
type browserEntry struct {
	path      string
	name      string
	kind      string
	size      string
	modified  string
	directory bool
}
```

The full path and `directory` flag do not need hidden table columns. They stay
in normal Go fields where application logic can use them directly.

## Turn entries into table rows

After `os.ReadDir` gathers the directory, one small loop creates the visible
rows:

```go
rows := make([][]string, len(entries))
for index, entry := range entries {
	rows[index] = []string{
		entry.name,
		entry.kind,
		entry.size,
		entry.modified,
	}
}
table.SetRows(rows...)
```

The row at index `n` describes the application entry at index `n`. That simple
relationship makes selection and activation easy to understand.

## Open a selected folder

The table's activation callback checks the row index and uses the matching
entry:

```go
table.OnActivate(func(row int, values []string) {
	if row < 0 || row >= len(entries) {
		return
	}
	entry := entries[row]
	if entry.directory {
		loadDirectory(entry.path)
		return
	}
	rosaline.Message("File details", entry.path)
})
```

Double-click and Enter share the same callback, so keyboard and mouse behavior
cannot drift apart.

## One function owns directory changes

`loadDirectory` is responsible for every navigation route:

- opening a table row
- entering a path in the address box
- pressing Up or Home
- using the Go menu
- refreshing the current directory

It resolves the path, reads entries, creates rows, updates the table, and sets
the status text. Errors restore the last valid address and display a Rosaline
error dialog. Centralizing that work prevents each button from implementing a
slightly different form of navigation.

## Why the example sorts before displaying

`os.ReadDir` returns entries sorted by filename. The example uses
`sort.SliceStable` to move folders before files while preserving a predictable
case-insensitive name order inside each group. Sorting is ordinary application
logic; the Table widget only displays the order it receives.

## Formatting belongs near the application

The `formatSize` function converts byte counts into readable KB, MB, GB, or TB
values. `fileKind` describes folders, symbolic links, normal files, and special
files. These are file-browser decisions rather than general GUI behavior, so
they live in the example instead of making Rosaline larger.

## Features to try next

- Show an image preview when the selected file is an image.
- Add a checkbox that includes or hides dotfiles.
- Remember the last directory in a JSON preferences file.
- Add a second tab containing detailed information about the selected entry.
- Add a Tree widget later for a folder sidebar.

## Go concepts used here

- structs and slices
- filesystem paths
- directory iteration
- stable sorting
- closures over shared application state
- early returns for error handling
- reusable functions

