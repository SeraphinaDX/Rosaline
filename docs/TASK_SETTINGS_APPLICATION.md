# Building a Task Settings Application

The Task Settings example combines Rosaline v0.10.0's radio groups, combo
boxes, sliders, and progress bars with text input, buttons, validation, and
dynamic labels.

Run the finished program from the Rosaline source tree:

```bash
CGO_ENABLED=0 go run ./examples/tasksettings
```

The complete source is in
[`examples/tasksettings/main.go`](../examples/tasksettings/main.go).

## Start with ordinary Go values

The current settings are normal variables:

```go
taskName := "Write the Rosaline guide"
category := "Documentation"
priority := "normal"
completion := 35.0
status := "Ready"
```

The controls receive pointers to the values they edit. Save logic can read the
variables directly; it does not need to query a form object or decode a
framework-specific model.

## Use the right selection control

Category uses a combo box because it may grow into a longer list while
priority uses a radio group because its three choices should remain visible:

```go
categoryBox := rosaline.ComboBox(
	&category,
	"Documentation",
	"Development",
	"Design",
	"Testing",
)

priorityGroup := rosaline.RadioGroup(
	&priority,
	rosaline.Choice("Low", "low"),
	rosaline.Choice("Normal", "normal"),
	rosaline.Choice("High", "high"),
).Horizontal()
```

Notice that the radio labels use title case for people while their stored
values remain short lowercase strings suitable for files or application data.

## Share one value between controls

The completion slider and progress bar use the same pointer:

```go
completionSlider := rosaline.Slider(&completion, 0, 100).Step(5)
completionBar := rosaline.ProgressBar(&completion)
```

When a slider event changes `completion`, Rosaline refreshes both widgets. This
is ordinary pointer sharing: there is no extra synchronization syntax to learn.

## Show unknown work separately

The example also creates an indeterminate progress bar:

```go
busyBar := rosaline.ProgressBar(nil).Busy().Stop()
```

`Busy` selects indeterminate mode. The immediate `Stop` leaves it paused so the
example's Start and Stop buttons can demonstrate explicit control. A real
application would normally call `Busy` when work begins and `Stop` or
`Determinate` when that work ends.

## Replace choices without rebuilding the application

The Planning and Work buttons call `SetOptions` and `SetChoices`:

```go
categoryBox.SetOptions("Planning", "Research", "Documentation")
priorityGroup.SetChoices(
	rosaline.Choice("Someday", "someday"),
	rosaline.Choice("This week", "week"),
	rosaline.Choice("Today", "today"),
)
```

Each widget keeps the old value when possible. If it is not in the replacement
set, the first new value becomes selected. The pointer is therefore always
safe to read in the Save callback.

## Validate at the point of action

The text box accepts an empty task name while the user is editing. Save checks
the final input and gives a clear message:

```go
if strings.TrimSpace(taskName) == "" {
	rosaline.Message("Missing task name", "Please enter a name for the task.")
	return
}
```

This is usually friendlier than showing an error after every keystroke. More
complex applications can move validation into a reusable Go function and call
it from buttons, menus, or keyboard shortcuts.

## Features to try next

- Add a due-date string with another `TextBox`.
- Save the task struct as JSON using `encoding/json` and `os.WriteFile`.
- Add a `CheckBox` for reminders.
- Put several tasks in a `Table` and open their settings in a second window.
- Add a timer that advances a short simulated task in small steps.

## Go concepts used here

- variables and pointers
- callbacks and closures
- reusable functions
- string validation
- formatted strings
- sharing one value between several objects
