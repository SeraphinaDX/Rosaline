# Multiple Windows

Rosaline applications can open secondary windows without creating another
event loop. A `Window` keeps its own content, title, menu, focus order, theme,
and timers while sharing ordinary Go state with the rest of the application.

## Smallest complete example

Create a new folder containing this `main.go`:

```go
package main

import "github.com/SeraphinaDX/Rosaline"

func main() {
	var about *rosaline.Window
	about = rosaline.NewWindow(rosaline.WindowOptions{
		Title:  "About",
		Width:  420,
		Height: 260,
		Parent: rosaline.MainWindow(),
		Content: rosaline.Column(
			rosaline.Label("My first multi-window app"),
			rosaline.Button("Close", func() {
				about.Close()
			}),
		).Gap(12),
	})

	rosaline.Run(
		rosaline.Button("About", func() {
			about.Show()
		}),
	)
}
```

Then run it:

```bash
go mod init example.com/window-demo
go get github.com/SeraphinaDX/Rosaline
CGO_ENABLED=0 go run .
```

`NewWindow` creates a reusable Go handle but does not display anything yet.
`Show` is called from the button callback after Rosaline's event loop is
running.

## Showing without duplicates

Call `Show` whenever the user requests a tool window:

```go
settingsWindow.Show()
```

If the window is closed, Rosaline opens it. If it is already open, Rosaline
raises and focuses the existing window. Applications do not need a separate
Boolean or duplicate-window check.

`Focus` only raises an open window:

```go
settingsWindow.Focus()
```

`IsOpen` is useful for status text and commands that depend on window state:

```go
if settingsWindow.IsOpen() {
	fmt.Println("settings are visible")
}
```

## Closing and reopening

`Close` destroys the native window and safely detaches its focus controls and
timers:

```go
settingsWindow.Close()
```

The Go handle, application variables, and widget descriptions remain. A later
`Show` remounts them in a new native window. Calling `Close` more than once is
safe.

Use `OnClose` when other application state should change regardless of whether
the user pressed a Close button or the window manager's close control:

```go
settingsWindow := rosaline.NewWindow(rosaline.WindowOptions{
	Title: "Settings",
	OnClose: func() {
		status = "Settings closed"
	},
	Content: settingsContent,
})
```

## Parent and child windows

Use `MainWindow()` when a tool belongs to the primary application window:

```go
about := rosaline.NewWindow(rosaline.WindowOptions{
	Parent: rosaline.MainWindow(),
	Content: rosaline.Label("About this application"),
})
```

A secondary window can also parent another secondary window:

```go
preview := rosaline.NewWindow(rosaline.WindowOptions{
	Parent: editor,
	Content: previewContent,
})
```

Showing `preview` automatically opens `editor` first when needed. Closing
`editor` also closes its preview and any other open descendants. Omit `Parent`
when a window should be independent of the other secondary windows.

## Sharing normal Go state

Windows can read and update the same variables:

```go
message := "Hello"

editor := rosaline.NewWindow(rosaline.WindowOptions{
	Content: rosaline.TextBox(&message),
})

preview := rosaline.NewWindow(rosaline.WindowOptions{
	Content: rosaline.LabelFunc(func() string {
		return message
	}),
})
```

After a Rosaline callback, dynamic widgets refresh in every open window. No
special cross-window binding or message bus is required.

## Changing a title

`SetTitle` updates an open window immediately and remembers the new title when
the window is reopened:

```go
editor.SetTitle("Editor — Unsaved")
```

An empty title uses `Rosaline`.

## Window-specific menus, themes, and timers

`WindowOptions` deliberately resembles `App`:

```go
tool := rosaline.NewWindow(rosaline.WindowOptions{
	Title:   "Tool",
	Width:   600,
	Height:  420,
	Padding: 16,
	Theme:   customTheme,
	Menu:    toolMenu,
	Timers:  []*rosaline.Timer{toolTimer},
	Content: toolContent,
})
```

An omitted theme inherits the parent or primary application theme. Shortcuts
in a secondary window's menu belong to that window. Its timers begin when it
opens, detach when it closes, and resume on reopening when they are still in a
running state.

## Closing the whole application

Use `rosaline.Quit` to close the primary window and every open secondary
window. `rosaline.MainWindow().Close()` has the same effect. Normal secondary
window buttons should call their own window's `Close` method instead.

## Common mistakes

- Call `Show` from a Rosaline callback, not before `Run` or `RunApp` begins.
- Use `var window *rosaline.Window` followed by assignment when the window's
  own callbacks need to call `window.Close()`.
- Use `MainWindow()` as `Parent`; do not create a second event loop.
- Reuse one `Window` handle. Repeated `Show` already prevents duplicates.
- Do not mount the same Widget value in two open windows. Create one widget for
  each location and let them share ordinary Go variables instead.
- Put a timer in only one open window at a time; one `Timer` cannot belong to
  two window contexts simultaneously.

## Complete application

Run the Project Desk example from Rosaline's source root:

```bash
CGO_ENABLED=0 go run ./examples/windows
```

It combines a primary window, editor, child preview, About window, shared text,
dynamic titles, menus, shortcuts, cross-window refresh, duplicate prevention,
and cascading closure. The source is
[`examples/windows/main.go`](../examples/windows/main.go).

## Go concepts used here

- pointers and reusable handles
- structs and optional fields
- callbacks and closures
- shared variables
- parent and child relationships
