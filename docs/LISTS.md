# Lists

A `List` presents several text choices in a compact, scrollable control. It
selects the first item by default and uses the platform's familiar mouse and
keyboard behavior.

## Smallest complete example

Create a new folder containing this `main.go`:

```go
package main

import "github.com/SeraphinaDX/Rosaline"

func main() {
	choice := "Rose"

	rosaline.Run(
		rosaline.Column(
			rosaline.LabelFunc(func() string {
				return "Selected: " + choice
			}),
			rosaline.List("Rose", "Violet", "Peony").
				OnSelect(func(index int, value string) {
					choice = value
				}),
		),
	)
}
```

Then run it:

```bash
go mod init example.com/list-demo
go get github.com/SeraphinaDX/Rosaline
CGO_ENABLED=0 go run .
```

The callback receives both the zero-based item `index` and its `value`. The
example only needs the value, so it leaves `index` unused.

## Selection and activation

`OnSelect` runs when the highlighted item changes. `OnActivate` represents a
stronger action: double-clicking an item or pressing Enter while it is selected.

```go
list := rosaline.List("New document", "Open document", "Examples").
	OnSelect(func(index int, value string) {
		fmt.Println("selected", index, value)
	}).
	OnActivate(func(index int, value string) {
		rosaline.Message("Open", value)
	})
```

Use selection for previews and descriptions. Use activation when the item
should open, launch, or confirm something.

## Reading and changing selection

`Selected` safely reports whether the list has a selection:

```go
index, value, ok := list.Selected()
if ok {
	fmt.Println(index, value)
}
```

`Select(2)` selects the third item. An out-of-range index clears the selection.
Programmatic changes run `OnSelect` after the control has been mounted.

## Replacing the items

`SetItems` is useful for search results, recent files, or other changing data:

```go
list.SetItems("January", "February", "March")
```

Rosaline keeps the old selection index when possible, moves it to the last
available item when the new list is shorter, and selects the first item when an
empty list becomes non-empty. `Items` returns a copy, so changing the returned
slice cannot accidentally change the widget.

## Size, scrolling, and expansion

```go
rosaline.List(items...).
	Size(36, 12).
	Expand()
```

The first size is an approximate width in characters; the second is the number
of visible rows. A vertical scrollbar appears beside the list. `Expand` lets a
surrounding row, column, or tab give it extra space.

Arrow keys move through items, Page Up and Page Down move farther, Home and End
jump to the ends, Enter activates the current item, and Tab moves to the next
control.

## Common mistakes

- Remember that indices start at zero: index `0` is the first item.
- Check `ok` from `Selected` before using its index or value.
- Use `OnSelect` for lightweight preview work. Reserve slow work for an
  activation button or `OnActivate`.
- Keep the returned widget when you need `Select`, `SetItems`, or `Selected`
  later.

## Go concepts used here

- slices and variadic arguments
- zero-based indices
- multiple return values
- callbacks and closures

See the lists working with tabs, form controls, and a canvas preview in the
[Preferences application](PREFERENCES_APPLICATION.md).

