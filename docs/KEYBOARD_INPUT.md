# Keyboard Input and Shortcuts

Rosaline provides friendly key events for windows and canvases, plus shortcuts
that work without creating a menu. Application code never sees backend event
types or platform key names.

## Smallest complete example

This program shows each key as it is pressed and released:

```go
package main

import "github.com/SeraphinaDX/Rosaline"

func main() {
	status := "Press a key"
	text := ""

	rosaline.RunApp(rosaline.App{
		Title: "Keyboard Example",
		OnKeyDown: func(event rosaline.KeyEvent) {
			status = "Pressed " + event.Key.String()
		},
		OnKeyUp: func(event rosaline.KeyEvent) {
			status = "Released " + event.Key.String()
		},
		Content: rosaline.Column(
			rosaline.LabelFunc(func() string { return status }),
			rosaline.TextBox(&text).Placeholder("Typing still works").Focus(),
		),
	})
}
```

Run it in a new module:

```bash
go mod init example.com/keyboard-demo
go get github.com/SeraphinaDX/Rosaline
CGO_ENABLED=0 go run .
```

Window key handlers observe keyboard activity without replacing the platform's
normal control behavior. The text box therefore continues to edit text.

## Key and Text

Every `KeyEvent` has two related fields:

```go
event.Key
event.Text
```

`Key` identifies the key. Named constants make navigation keys easy to read:

```go
if event.Key == rosaline.KeyLeft {
	moveLeft()
}
```

`Text` is the text produced by a key press. For example, pressing Shift+A has
key `rosaline.Key("a")` and text `"A"`. Arrow and modifier keys produce no
text. On some platforms release events do not include text, so use `Key` when
matching releases.

`event.Is` is a shorter equivalent comparison:

```go
if event.Is(rosaline.KeyEscape) {
	closeDialog()
}
```

Rosaline defines constants for common keys:

- `KeyEnter`, `KeyEscape`, `KeySpace`, `KeyTab`, and `KeyBackspace`
- `KeyDelete`, `KeyInsert`, `KeyHome`, and `KeyEnd`
- `KeyPageUp` and `KeyPageDown`
- `KeyLeft`, `KeyUp`, `KeyRight`, and `KeyDown`
- `KeyShift`, `KeyControl`, `KeyAlt`, `KeySuper`, and `KeyCapsLock`
- `KeyF1` through `KeyF12`

Printable keys use lowercase `Key` values:

```go
if event.Key == rosaline.Key("w") {
	moveUp()
}
```

## Modifier keys

Events include four convenient Boolean fields:

```go
event.Shift
event.Control
event.Alt
event.Primary
```

`Primary` means Control on Linux and Windows and Command on macOS. It is useful
when an application's main action should follow each platform's convention.
On macOS, `Alt` represents the Option key.

## Canvas keyboard input

Canvas handlers run only while that canvas has keyboard focus:

```go
canvas := rosaline.Canvas(draw).
	Focus().
	OnKeyDown(func(event rosaline.KeyEvent) {
		if event.Is(rosaline.KeyRight) {
			x += 10
		}
	}).
	OnKeyUp(func(event rosaline.KeyEvent) {
		fmt.Println("released", event.Key)
	})
```

`Focus` asks for initial focus before the window opens. Calling it later from a
button or menu callback immediately restores focus to the canvas. A canvas with
either key handler automatically joins Tab focus order, shows a focus ring, and
receives focus when clicked. Rosaline redraws it after a key callback unless
the callback already called `Redraw`.

For smooth game controls, set a Boolean on key down and clear it on key up. The
complete [Starshower application](STARSHOWER_APPLICATION.md) uses this pattern
for turning, thrust, and firing without relying on platform key-repeat timing.

If the same window also has an application key handler, a canvas event can be
observed by both. Keep canvas movement in the canvas handler and window-wide
behavior in the application handler so one press does not perform the same
action twice.

## Standalone shortcuts

Shortcuts are useful when a key combination should run one command regardless
of which normal control has focus:

```go
rosaline.RunApp(rosaline.App{
	Shortcuts: rosaline.Shortcuts(
		rosaline.Shortcut("Primary+S", save),
		rosaline.Shortcut("Primary+Shift+S", saveAs),
		rosaline.Shortcut("F1", showHelp),
	),
	Content: content,
})
```

`Shortcuts` is a small convenience function that avoids a slice literal. This
is equivalent:

```go
Shortcuts: []rosaline.KeyShortcut{
	rosaline.Shortcut("Escape", cancel),
}
```

Shortcuts flush pointer-bound form values before their callback and refresh
dynamic controls afterward, exactly like buttons and menu commands.

## Shortcut names

Supported modifiers are:

- `Primary`
- `Ctrl` or `Control`
- `Shift`
- `Alt`
- `Option`
- `Cmd` or `Command`
- `Super`

Named keys include Enter, Escape, Space, Backspace, Delete, Insert, Home, End,
PageUp, PageDown, the arrow keys, Tab, and F1 through F24. `Plus` and `Minus`
are available when those characters would make the combination hard to read.
`Option`, `Cmd`, and `Command` are macOS-specific; prefer `Primary` when one
shortcut should follow every platform's convention.

A bare or Shift-only printable shortcut such as `Shortcut("S", save)` is
rejected so normal typing remains safe. Use `Primary+S`. Unmodified named keys
such as Escape and F1 are supported. Tab and Shift+Tab are reserved for moving
keyboard focus; use another combination for application commands.

## Secondary windows

Each secondary window owns its key handlers and shortcuts:

```go
editor := rosaline.NewWindow(rosaline.WindowOptions{
	Title: "Editor",
	Shortcuts: rosaline.Shortcuts(
		rosaline.Shortcut("Primary+S", saveEditor),
		rosaline.Shortcut("Escape", closeEditor),
	),
	OnKeyDown: func(event rosaline.KeyEvent) {
		lastEditorKey = event.Key.String()
	},
	Content: editorContent,
})
```

Closing and reopening the window remounts its keyboard bindings safely. The
shortcut only applies while that window is open and active.

## Common mistakes

- Use `Key` for navigation and releases; use `Text` for the character a user
  typed.
- Write `rosaline.Key("w")`, not `rosaline.KeyW`, for printable letters.
- Call `Focus` when a canvas should receive keys immediately.
- Use `Primary` for a conventional cross-platform shortcut.
- When canvas letter controls overlap a shortcut, check `Primary`, `Control`,
  and `Alt` before handling the letter.
- Do not perform slow work directly inside a key callback.
- Avoid assigning the same combination to multiple shortcuts in one window;
  the last binding wins.

## Go concepts used here

- named string types and constants
- struct fields
- callbacks and closures
- Boolean conditions
- slices and variadic functions

See every part working together in the
[Keyboard Garden application](KEYBOARD_GARDEN_APPLICATION.md).
