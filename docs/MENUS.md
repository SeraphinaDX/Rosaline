# Menus

Rosaline applications can attach native menu bars to `App`.

## Complete example

```go
package main

import "github.com/SeraphinaDX/Rosaline"

func main() {
	rosaline.RunApp(rosaline.App{
		Title: "Menu Example",
		Menu: rosaline.MenuBar(
			rosaline.Menu("File",
				rosaline.MenuItem("Open…", func() {
					rosaline.Message("Open", "Open was selected.")
				}).Shortcut("Ctrl+O"),
				rosaline.MenuSeparator(),
				rosaline.MenuItem("Quit", rosaline.Quit).
					Shortcut("Ctrl+Q"),
			),
		),
		Content: rosaline.Label("Try the File menu."),
	})
}
```

`MenuBar` contains menus, and each `Menu` contains items or separators. Menu
callbacks are ordinary `func()` values, just like button callbacks.

## Shortcuts

`Shortcut` both displays and binds the shortcut:

```go
rosaline.MenuItem("Save As…", saveAs).Shortcut("Ctrl+Shift+S")
```

Supported modifiers are `Ctrl`, `Control`, `Shift`, and `Alt`. Combine them
with a letter or named key such as `F4`. Menu shortcuts flush current form
values before running and refresh dynamic widgets afterward.

## Sharing actions

Keep an action in one function when a button, menu item, and shortcut should do
the same thing:

```go
save := func() {
	// Save the document.
}

button := rosaline.Button("Save", save)
item := rosaline.MenuItem("Save", save).Shortcut("Ctrl+S")
```

This avoids two versions of the same application logic.

## Closing the application

`rosaline.Quit` safely closes the main window and event loop. It can be used as
a menu callback without wrapping it in another function.

## Go concepts used here

- Functions stored in variables
- Reusing callbacks
- Variadic functions
- Method chaining
