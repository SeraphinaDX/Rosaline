# Timers

Timers let an application run code later or at a regular interval without
blocking its window. Rosaline timers run on the GUI event loop, so timer
callbacks can safely update normal interface state.

## A repeating timer

This complete program counts the seconds since its window opened:

```go
package main

import (
	"fmt"
	"time"

	"github.com/SeraphinaDX/Rosaline"
)

func main() {
	seconds := 0

	clock := rosaline.Every(time.Second, func() {
		seconds++
	})

	rosaline.RunApp(rosaline.App{
		Title:  "Clock",
		Timers: []*rosaline.Timer{clock},
		Content: rosaline.LabelFunc(func() string {
			return fmt.Sprintf("Open for %d seconds", seconds)
		}),
	})
}
```

`Every` creates a repeating timer. Adding it to `App.Timers` gives the timer to
the application. It starts when the window opens and stops belonging to the
event loop when the window closes.

Rosaline refreshes dynamic widgets after each timer callback, so the
`LabelFunc` above updates automatically.

## Running something once

Use `After` when work should happen one time after a delay:

```go
message := "Please wait..."

ready := rosaline.After(2*time.Second, func() {
	message = "Ready!"
})
```

Include `ready` in `App.Timers` just like a repeating timer. After its callback
runs, `ready.Running()` becomes false.

## Starting and stopping

Every timer has four simple controls:

- `Start()` continues a stopped timer.
- `Stop()` pauses it and cancels its current wait.
- `Restart()` starts its full delay again.
- `Running()` reports whether it is started.

For example:

```go
rosaline.Button("Pause", func() {
	clock.Stop()
})

rosaline.Button("Continue", func() {
	clock.Start()
})
```

Calling `Start` on a timer that is already running has no effect. `Restart`
always resets the current wait, which is particularly useful for one-shot
timers such as inactivity notices.

## Durations are ordinary Go

Rosaline uses the standard library's `time.Duration` type:

```go
rosaline.Every(500*time.Millisecond, update)
rosaline.After(3*time.Second, finish)
rosaline.Every(2*time.Minute, saveDraft)
```

This is why timer programs import Go's `time` package. Durations shorter than
one millisecond are safely treated as one millisecond.

## Common mistakes

### Forgetting App.Timers

Creating a timer is not enough by itself. The application must own it:

```go
rosaline.RunApp(rosaline.App{
	Timers: []*rosaline.Timer{clock},
	Content: content,
})
```

This ownership is what keeps timer lifetime predictable.

Secondary windows can own timers through `WindowOptions.Timers`. Those timers
attach when the window opens and detach when it closes. See
[MULTIPLE_WINDOWS.md](MULTIPLE_WINDOWS.md) for window lifecycle details.

### Sleeping in a callback

Do not use `time.Sleep` in a button or timer callback. Sleeping blocks the GUI
event loop and makes the window unresponsive. Use `After` to schedule work for
later instead.

### Calling timer methods from a goroutine

Call `Start`, `Stop`, and `Restart` from Rosaline callbacks. Background work and
thread-safe delivery back to the GUI use
[`Background`](BACKGROUND_TASKS.md). A task may safely post a result to a timer
or other interface state through its reporter.

## Go concepts used here

- Importing a standard-library package
- `time.Duration` values
- Callback functions
- Slices of pointers
- Variables captured by functions
