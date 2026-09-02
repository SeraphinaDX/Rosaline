# Menus

Rosaline applications can attach native menu bars to `App` and to secondary
windows created with `NewWindow`.

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
				}).Shortcut("Primary+O"),
				rosaline.MenuSeparator(),
				rosaline.MenuItem("Quit", rosaline.Quit).
					Shortcut("Primary+Q"),
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
rosaline.MenuItem("Save As…", saveAs).Shortcut("Primary+Shift+S")
```

`Primary` means Control on Linux and Windows and Command on macOS. Explicit
`Ctrl`, `Control`, `Shift`, `Alt`, `Option`, `Cmd`, `Command`, and `Super`
modifiers are also supported. Combine them with a letter or named key such as
`F4`. Menu shortcuts flush current form values before running and refresh
dynamic widgets afterward.

Menu items and standalone `rosaline.Shortcut` values share the same spelling
and platform handling. See [KEYBOARD_INPUT.md](KEYBOARD_INPUT.md) for the full
list and shortcuts that do not need a menu.

## Sharing actions

Keep an action in one function when a button, menu item, and shortcut should do
the same thing:

```go
save := func() {
	// Save the document.
}

button := rosaline.Button("Save", save)
item := rosaline.MenuItem("Save", save).Shortcut("Primary+S")
```

This avoids two versions of the same application logic.

## Secondary-window menus

Pass a menu bar through `WindowOptions` just as you would through `App`:

```go
editor := rosaline.NewWindow(rosaline.WindowOptions{
	Title: "Editor",
	Menu: rosaline.MenuBar(
		rosaline.Menu("File",
			rosaline.MenuItem("Save", save).Shortcut("Primary+S"),
		),
	),
	Content: editorContent,
})
```

The menu and its shortcuts belong to that secondary window. See
[MULTIPLE_WINDOWS.md](MULTIPLE_WINDOWS.md) for its complete lifecycle.

## Closing the application

`rosaline.Quit` safely closes the main window, every secondary window, and the
event loop. A secondary window's Close item should call that window's `Close`
method when the rest of the application should remain open.

## Go concepts used here

- Functions stored in variables
- Reusing callbacks
- Variadic functions
- Method chaining
