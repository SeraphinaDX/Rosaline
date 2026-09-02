# Sliders

A `Slider` lets the user choose a number within a range. It binds directly to
an ordinary Go `float64`.

## Smallest complete example

```go
package main

import (
	"fmt"

	"github.com/SeraphinaDX/Rosaline"
)

func main() {
	volume := 40.0

	rosaline.Run(
		rosaline.Column(
			rosaline.LabelFunc(func() string {
				return fmt.Sprintf("Volume: %.0f", volume)
			}),
			rosaline.Slider(&volume, 0, 100),
		),
	)
}
```

Run it in a new module:

```bash
go mod init example.com/slider-demo
go get github.com/SeraphinaDX/Rosaline
CGO_ENABLED=0 go run .
```

The three numeric arguments are the bound value pointer, minimum, and maximum.
Rosaline clamps an initial value outside the range.

## Fixed steps

The default slider is continuous. `Step` rounds values to a useful interval:

```go
rosaline.Slider(&rating, 1, 10).Step(1)
rosaline.Slider(&opacity, 0, 1).Step(0.05)
```

Steps are measured from the minimum. A zero, negative, infinite, or not-a-number
step safely restores continuous movement.

## Reacting to changes

```go
slider.OnChange(func(value float64) {
	status = fmt.Sprintf("New value: %.1f", value)
})
```

The callback runs after the bound value has changed. It also runs when
`SetValue` changes a mounted slider to a different value.

## Reading and changing a slider

```go
slider := rosaline.Slider(&zoom, 25, 400).Step(25)
slider.SetValue(125)

minimum, maximum := slider.Bounds()
fmt.Println(slider.Value(), minimum, maximum)
```

`SetValue` applies the same clamping and step rounding as mouse or keyboard
input. `SetRange` changes both bounds and immediately normalizes the current
value.

## Size, direction, and focus

```go
rosaline.Slider(&level, 0, 100).
	Length(320).
	Vertical().
	Focus()
```

Sliders are horizontal by default. `Length` is measured in pixels. `Vertical`
and `Horizontal` change the direction, and `Focus` requests initial keyboard
focus. Platform-native arrow keys work without application event handling.

## Share a value with a progress bar

Two controls can bind to the same pointer:

```go
completion := 25.0

rosaline.Column(
	rosaline.Slider(&completion, 0, 100).Step(5),
	rosaline.ProgressBar(&completion),
)
```

Rosaline refreshes the progress bar after slider events. There is no binding
language or extra state object; both widgets simply use the same Go variable.

## Safe range behavior

- Reversed bounds are swapped.
- Equal bounds become a one-unit range.
- Invalid bounds use safe defaults.
- Infinite and not-a-number values become the minimum.

These rules prevent malformed data from reaching the native control, but an
application should still validate user-facing settings according to its own
requirements.

## Go concepts used here

- `float64` variables and pointers
- multiple return values
- callbacks and closures
- formatted strings

See sliders working with progress bars in the
[Task Settings application](TASK_SETTINGS_APPLICATION.md).
