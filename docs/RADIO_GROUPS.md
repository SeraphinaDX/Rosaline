# Radio Groups

A `RadioGroup` presents a small set of choices when exactly one must be
selected. It binds the selected choice to an ordinary Go string.

## Smallest complete example

Create a new folder containing this `main.go`:

```go
package main

import "github.com/SeraphinaDX/Rosaline"

func main() {
	priority := "normal"

	rosaline.Run(
		rosaline.Column(
			rosaline.Label("Priority"),
			rosaline.RadioGroup(
				&priority,
				rosaline.Choice("Low", "low"),
				rosaline.Choice("Normal", "normal"),
				rosaline.Choice("High", "high"),
			),
			rosaline.LabelFunc(func() string {
				return "Selected value: " + priority
			}),
		),
	)
}
```

Then run it:

```bash
go mod init example.com/radio-demo
go get github.com/SeraphinaDX/Rosaline
CGO_ENABLED=0 go run .
```

`&priority` gives the group access to the string so it can update that string
when the user makes a choice.

## Labels and values

Each choice has visible text and a stored value:

```go
rosaline.Choice("Highest priority", "high")
```

The user sees `Highest priority`, while the bound Go string receives `high`.
Keeping those separate lets an application use friendly text without making
its saved data depend on that wording.

## Horizontal choices

Radio groups are vertical by default. Use `Horizontal` for a short set that
fits comfortably on one line:

```go
group := rosaline.RadioGroup(
	&size,
	rosaline.Choice("Small", "s"),
	rosaline.Choice("Medium", "m"),
	rosaline.Choice("Large", "l"),
).Horizontal()
```

`Vertical` switches it back. `Focus` asks Rosaline to focus the selected
choice when the window opens. A whole group occupies one place in Rosaline's
Tab order; the platform's normal radio-button keys move within it.

## Reacting to a change

`OnChange` receives the newly stored value:

```go
group.OnChange(func(value string) {
	fmt.Println("priority is now", value)
})
```

Rosaline updates the bound string before calling the function. Dynamic labels
and other bound controls refresh after the callback.

## Reading and changing the selection

Keep the returned widget when another callback must control it:

```go
group := rosaline.RadioGroup(&priority,
	rosaline.Choice("Low", "low"),
	rosaline.Choice("High", "high"),
)

group.Select("high")
fmt.Println(group.Selected())
```

`Select` ignores a value that is not one of the choices. `Choices` returns a
copy suitable for inspection without exposing the group's internal slice.

## Replacing choices

Applications can replace every choice at runtime:

```go
group.SetChoices(
	rosaline.Choice("Today", "today"),
	rosaline.Choice("This week", "week"),
)
```

If the previous value is still available, Rosaline keeps it. Otherwise it
selects the first new choice. An empty replacement clears the bound string.
Duplicate values are ignored because two visible choices must not represent
the same selection.

## Common mistakes

- Pass a pointer such as `&priority`, not just `priority`.
- Use `Choice(label, value)`; the label is what users see and the value is what
  the Go variable stores.
- Use a `ComboBox` instead when the choice list is too long to display at once.
- Keep the widget in a variable when you need `Select` or `SetChoices` later.

## Go concepts used here

- string variables and pointers
- variadic function arguments
- callbacks and closures
- return values

See radio groups working with other controls in the
[Task Settings application](TASK_SETTINGS_APPLICATION.md).
