# Changing Controls at Runtime

Rosaline normally keeps application state in ordinary Go variables. Sometimes
an event also needs to change a control directly—for example, changing a status
label or disabling a button after it is clicked.

## Complete example

```go
package main

import rosaline "github.com/SeraphinaDX/Rosaline"

func main() {
	status := rosaline.Label("Ready")
	save := rosaline.Button("Save", nil).Primary()

	save.OnClick(func() {
		status.SetText("Saved")
		save.SetText("Saved!")
		save.SetEnabled(false)
	})

	rosaline.Run(rosaline.Column(status, save))
}
```

The variables hold pointers to the same controls that Rosaline displays. A
setter used inside a callback updates the mounted interface immediately.

## Available methods

- `LabelWidget`: `Text`, `SetText`
- `ButtonWidget`: `OnClick`, `Text`, `SetText`, `Enabled`, `SetEnabled`
- `TextBoxWidget`: `Text`, `SetText`, `Enabled`, `SetEnabled`, `Focus`
- `CheckBoxWidget`: `Text`, `SetText`, `Checked`, `SetChecked`, `Enabled`,
  `SetEnabled`

`SetText` on a text box and `SetChecked` on a checkbox also update the ordinary
Go value passed to the constructor. When the control is mounted and the value
changes, its `OnChange` callback runs.

## Rosaline Studio components

Rosaline Studio assigns a component name to every designed control. Generated
event methods can retrieve them through `app.Widgets()`:

```go
func (app *Application) SaveButtonClick() {
	app.Widgets().StatusLabel.SetText("Saved")
	app.Widgets().SaveButton.SetEnabled(false)
}
```

The component name shown in Studio becomes the Go field name. Names therefore
start with a capital letter and contain no spaces, such as `StatusLabel` or
`SaveButton`.

## Common mistakes

Do not create a second control inside the handler. This changes a new object
that is not part of the window:

```go
rosaline.Label("Ready").SetText("Saved")
```

Keep the displayed widget in a variable—or use Studio's generated component
reference—and call the setter on that same object.
