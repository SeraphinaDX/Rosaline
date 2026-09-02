# Progress Bars

A `ProgressBar` communicates that work is underway. It supports known numeric
progress and an indeterminate busy animation.

## Determinate progress

Use determinate progress when the application knows how much work is complete:

```go
package main

import "github.com/SeraphinaDX/Rosaline"

func main() {
	completion := 35.0
	bar := rosaline.ProgressBar(&completion)

	rosaline.Run(
		rosaline.Column(
			bar,
			rosaline.Button("Advance", func() {
				bar.SetValue(bar.Value() + 10)
			}),
		),
	)
}
```

Run it in a new module:

```bash
go mod init example.com/progress-demo
go get github.com/SeraphinaDX/Rosaline
CGO_ENABLED=0 go run .
```

The default range is zero through 100. Values below zero or above the maximum
are clamped.

## A different maximum

```go
currentPage := 3.0
bar := rosaline.ProgressBar(&currentPage).Maximum(12)
```

`Max` returns the configured upper bound. A non-positive, infinite, or
not-a-number maximum safely uses 100.

## Busy progress

Use `Busy` when an operation has started but cannot report a meaningful
percentage:

```go
bar := rosaline.ProgressBar(nil).Busy()
```

`Busy` changes the bar to indeterminate mode and starts its animation. Passing
`nil` is convenient because a busy bar does not need an application value.

```go
bar.Stop()
bar.Start()
```

`Stop` pauses the animation and `Start` resumes it. `Running` reports its state.
`Determinate` stops the animation and switches back to the numeric value;
`IsBusy` reports which mode is active.

## Size and direction

```go
bar := rosaline.ProgressBar(&value).
	Length(320).
	Vertical()
```

The default is a 260-pixel horizontal bar. `Horizontal` switches a vertical
bar back. Progress bars display status rather than accepting input, so they do
not participate in keyboard focus order.

## Updating during real work

Button and timer callbacks may assign to the bound variable normally. Rosaline
refreshes the bar after those callbacks:

```go
progress := 0.0

timer := rosaline.Every(100*time.Millisecond, func() {
	progress++
})
```

Keep expensive work out of a UI callback. Rosaline's current public API does
not yet offer a general goroutine-to-UI scheduling method, so do not update
widgets directly from a goroutine. For operations that can be divided into
small quick steps, an application timer is a safe fit.

## Common mistakes

- Use a `float64` pointer for determinate progress.
- Call `Busy` before `Start`; `Start` only resumes a bar already in busy mode.
- Use `SetValue`, a Rosaline callback, or shared pointer binding so the screen
  refreshes after a value changes.
- Do not use a progress bar as a slider; it intentionally does not accept user
  input.

## Go concepts used here

- numeric variables and pointers
- nil pointers as an optional input
- callbacks
- method chaining

See both progress modes in the
[Task Settings application](TASK_SETTINGS_APPLICATION.md).
