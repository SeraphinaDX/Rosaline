# CheckBox

A checkbox lets the user turn an option on or off. Rosaline binds it to an
ordinary Go `bool`.

## Complete example

```go
package main

import "github.com/SeraphinaDX/Rosaline"

func main() {
	var notifications bool

	rosaline.Run(
		rosaline.Column(
			rosaline.CheckBox("Enable notifications", &notifications),
			rosaline.Button("Save", func() {
				if notifications {
					rosaline.Message("Settings", "Notifications are enabled.")
				} else {
					rosaline.Message("Settings", "Notifications are disabled.")
				}
			}).Primary(),
		),
	)
}
```

`&notifications` gives the checkbox access to the Boolean variable. The value
is `true` while checked and `false` while unchecked.

## Reacting immediately

Most forms only read the value when a button is clicked. Use `OnChange` when
something should happen as soon as the user toggles the box:

```go
rosaline.CheckBox("Dark background", &dark).OnChange(func(checked bool) {
	if checked {
		status = "Dark background selected"
	} else {
		status = "Light background selected"
	}
})
```

The handler receives the new value. `LabelFunc` values refresh after it runs.

Use `Focus()` if the checkbox should receive keyboard focus when the window
opens. Space toggles a focused checkbox, Tab moves forward, and Shift+Tab moves
backward.

## Go concepts used here

- Boolean values
- Pointers with `&`
- `if` and `else`
- Callback functions

See [FORMS.md](FORMS.md) for a larger example.
