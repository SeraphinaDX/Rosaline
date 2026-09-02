# Combo Boxes

A `ComboBox` is a compact drop-down for choosing one string from a list. It is
read-only, so the bound value always matches one of the supplied options.

## Smallest complete example

```go
package main

import "github.com/SeraphinaDX/Rosaline"

func main() {
	color := "Rose"

	rosaline.Run(
		rosaline.Column(
			rosaline.Label("Accent color"),
			rosaline.ComboBox(&color, "Rose", "Violet", "Ocean"),
			rosaline.LabelFunc(func() string {
				return "Selected: " + color
			}),
		),
	)
}
```

Run it in a new module:

```bash
go mod init example.com/combo-demo
go get github.com/SeraphinaDX/Rosaline
CGO_ENABLED=0 go run .
```

The combo box reads the initial `color` and writes each new selection back to
the same variable.

## Width and initial focus

`Width` uses approximate text columns rather than pixels:

```go
rosaline.ComboBox(&country, countries...).Width(32).Focus()
```

`Focus` asks Rosaline to focus the control when the window opens. If several
controls request focus, the first mounted one wins.

## Reacting to a change

`OnChange` receives the newly selected option:

```go
combo.OnChange(func(selected string) {
	status = "Changed to " + selected
})
```

Rosaline updates the bound string first, runs the callback, and then refreshes
dynamic widgets.

## Programmatic selection

```go
combo.Select("Violet")
fmt.Println(combo.Selected())
```

`Select` ignores values that are not present. `Options` returns a copy of the
current option strings.

## Replacing options

Use `SetOptions` for choices that depend on another part of the application:

```go
combo.SetOptions("Planning", "Research", "Documentation")
```

Rosaline preserves the old value when it is still present. Otherwise it uses
the first replacement. An empty list clears the bound string. Duplicate option
strings are removed automatically.

## ComboBox or List?

Use a `ComboBox` when space is limited and the user only needs one choice. Use
a `List` when seeing several choices simultaneously or supporting a separate
activation action is important. Use a `RadioGroup` for a small decision whose
choices should always remain visible.

## Common mistakes

- Pass a string pointer: `ComboBox(&color, ...)`.
- Supply the options after the pointer, not as a separate slice without `...`.
- Keep the widget in a variable before calling `Select` or `SetOptions`.
- Use `OnChange` for selection changes; a read-only combo box does not accept
  arbitrary typed text.

## Go concepts used here

- string variables and pointers
- slices and `...`
- callbacks and closures
- method chaining

See dynamic combo-box options in the
[Task Settings application](TASK_SETTINGS_APPLICATION.md).
