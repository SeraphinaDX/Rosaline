# Tables

A `Table` displays rows of text beneath named columns. It includes native row
selection, keyboard navigation, activation, and both scrollbars automatically.

## Smallest complete example

Create a new folder containing this `main.go`:

```go
package main

import (
	"fmt"

	"github.com/SeraphinaDX/Rosaline"
)

func main() {
	files := rosaline.Table("Name", "Type", "Size").
		SetRows(
			[]string{"README.md", "Markdown", "8 KB"},
			[]string{"picture.png", "Image", "2.4 MB"},
		).
		OnActivate(func(row int, values []string) {
			rosaline.Message("Selected file", fmt.Sprintf(
				"Row %d contains %s",
				row,
				values[0],
			))
		})

	rosaline.Run(files)
}
```

Then run it:

```bash
go mod init example.com/table-demo
go get github.com/SeraphinaDX/Rosaline
CGO_ENABLED=0 go run .
```

The headings are ordinary strings. Each row is a `[]string` containing one
value for each column.

## Creating rows in a loop

Most applications receive data from files, a database, or a web service. Build
a slice of rows and pass it to `SetRows` with `...`:

```go
rows := make([][]string, len(people))
for index, person := range people {
	rows[index] = []string{person.Name, person.Email}
}

table.SetRows(rows...)
```

The `...` passes every row in the slice as an argument. This is normal Go and
does not require a Rosaline data model.

## Selection and activation

Use `OnSelect` for lightweight previews or status text:

```go
table.OnSelect(func(row int, values []string) {
	status = "Selected " + values[0]
})
```

Use `OnActivate` for the stronger action associated with double-clicking or
pressing Enter:

```go
table.OnActivate(func(row int, values []string) {
	openFile(values[0])
})
```

Both callbacks receive the zero-based row index and a copy of the visible row.
Changing the callback's slice cannot change the table.

## Reading and changing selection

`Selected` safely reports whether a row is selected:

```go
row, values, ok := table.Selected()
if ok {
	fmt.Println(row, values)
}
```

Call `table.Select(2)` to select the third row. An invalid index clears the
selection. A non-empty table selects its first row by default.

## Updating data

Replace every row with one call:

```go
table.SetRows(
	[]string{"January", "31 days"},
	[]string{"February", "28 days"},
)
```

Rosaline keeps the old selection index when possible and moves it to the last
row when the new data is shorter. `Rows` returns a deep copy when application
code needs to inspect all current data.

## Comfortable sizing

```go
table := rosaline.Table("Name", "Type", "Size").
	ColumnWidth(0, 320).
	ColumnWidth(1, 140).
	ColumnWidth(2, 100).
	Height(16).
	Expand()
```

Column indices start at zero. Widths use pixels; height is the preferred number
of visible rows. Users can still resize native column dividers. `Expand` lets
the table use extra space from a row, column, tab, or window.

## Uneven rows

Short rows are padded with empty cells. Extra values are trimmed:

```go
table := rosaline.Table("Name", "Email")
table.SetRows(
	[]string{"Ana"},
	[]string{"Britney", "britney@example.com", "not displayed"},
)
```

This keeps every returned row aligned with the table's headings and avoids
index errors in callbacks.

## Keyboard behavior

- Up and Down move between rows.
- Page Up and Page Down move through larger tables.
- Home and End move to the beginning or end.
- Enter activates the selected row.
- Tab and Shift+Tab move to neighboring controls.

## Common mistakes

- Remember that row and column indices start at zero.
- Check `ok` before using values returned by `Selected`.
- Keep non-display data in your application structs instead of hiding it in an
  extra table column.
- Use `SetRows(rows...)`, including the final `...`, when `rows` is a
  `[][]string`.

## Go concepts used here

- slices and slices of slices
- variadic arguments
- loops and indices
- multiple return values
- callbacks and closures

See the [File Browser](FILE_BROWSER.md) for a complete application built around
a dynamic table.

