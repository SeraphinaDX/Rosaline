# Building a Preferences Application

The Preferences example is Rosaline v0.6.0's larger-interface tutorial. It
combines tabs, selectable lists, form controls, buttons, dynamic labels, and a
canvas preview in one application.

Run the finished program from the Rosaline source tree:

```bash
CGO_ENABLED=0 go run ./examples/preferences
```

The complete source is in
[`examples/preferences/main.go`](../examples/preferences/main.go).

## Why the settings are ordinary Go values

The application starts with normal variables:

```go
selectedPalette := 0
selectedFont := 1
displayName := "Seraphina"
autosave := true
showTips := true
autoUpdates := true
status := "Ready"
```

Text boxes and checkboxes bind directly to the string and Boolean variables.
The selected list indices are integers because they identify entries in Go
slices. A save function can read every current value without querying backend
widgets or decoding a framework-specific model.

## One list controls a preview

The color choices live in a slice of small structs. Each struct keeps the name,
description, and colors that belong together:

```go
type palette struct {
	name        string
	description string
	background  rosaline.Color
	accent      rosaline.Color
	text        rosaline.Color
}
```

The theme list changes `selectedPalette`, updates a status string, and redraws
the canvas:

```go
themeList := rosaline.List("Rosaline", "Lavender", "Ocean", "Forest").
	OnSelect(func(index int, value string) {
		selectedPalette = index
		status = "Selected the " + value + " palette"
		preview.Redraw()
	})
```

This is a useful division of responsibility. The callback changes application
state; the canvas reads that state when drawing. The list does not need to know
how the preview is rendered.

## Tabs group related decisions

The application has Appearance, Editor, Updates, and About pages:

```go
tabs := rosaline.Tabs(
	rosaline.Tab("Appearance", appearanceUI),
	rosaline.Tab("Editor", editorUI),
	rosaline.Tab("Updates", updatesUI),
	rosaline.Tab("About", aboutUI),
).Expand().OnChange(func(index int, title string) {
	status = fmt.Sprintf("Opened %s (tab %d)", title, index+1)
})
```

The real example constructs each layout inline, but naming the page widgets as
above is equally valid and can make a much larger program easier to navigate.

## Restore is one reusable operation

The Restore button resets variables, selections, the visible tab, and the
preview together:

```go
restoreDefaults := func() {
	displayName = "Seraphina"
	autosave = true
	showTips = true
	autoUpdates = true
	selectedPalette = 0
	selectedFont = 1
	themeList.Select(selectedPalette)
	fontList.Select(selectedFont)
	tabs.Select(0)
	preview.Redraw()
	status = "Restored the default preferences"
}
```

Keeping that operation in one function prevents a menu item, keyboard shortcut,
or second button from implementing slightly different reset behavior later.

## Saving in a real application

The demonstration shows the chosen values in a message. A real preferences
application would encode the same variables with Go's `encoding/json` package
and write them with `os.WriteFile`. Rosaline deliberately leaves application
data formats and filesystem policy to ordinary Go code.

## Features to try next

- Add a new palette struct and one matching list item.
- Add a checkbox to the Editor page.
- Use `OnActivate` to save a palette when the user presses Enter.
- Store the preferences in a JSON file and load them before constructing the
  controls.
- Add validation so an empty display name cannot be saved.

## Go concepts used here

- structs and slices
- pointers to variables
- closures over shared state
- multiple return values
- reusable functions
- formatted strings

