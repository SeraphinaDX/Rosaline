# Background Tasks

Background tasks keep slow file, network, or computation work from freezing a
window. Rosaline runs the work in a Go goroutine and brings its progress and
results back to the GUI thread safely.

## A complete first task

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/SeraphinaDX/Rosaline"
)

func main() {
	progress := 0.0
	status := "Ready"

	task := rosaline.Background(func(ctx context.Context, report *rosaline.TaskReporter) error {
		for step := 1; step <= 100; step++ {
			if !report.Report(float64(step), fmt.Sprintf("Step %d of 100", step)) {
				return ctx.Err()
			}
			time.Sleep(20 * time.Millisecond)
		}
		return nil
	}).OnProgress(func(update rosaline.TaskProgress) {
		progress = update.Percent
		status = update.Message
	}).OnDone(func(err error) {
		if err != nil {
			status = "Stopped: " + err.Error()
			return
		}
		status = "Finished"
	}).AutoStart()

	rosaline.RunApp(rosaline.App{
		Title: "Background Task",
		Tasks: []*rosaline.Task{task},
		Content: rosaline.Column(
			rosaline.ProgressBar(&progress),
			rosaline.LabelFunc(func() string { return status }),
			rosaline.Row(
				rosaline.Button("Start", task.Start).Primary(),
				rosaline.Button("Cancel", task.Cancel),
			),
		),
	})
}
```

`Background` creates a stopped `Task`. `App.Tasks` gives that task to the
window, which is what makes its lifetime safe. `AutoStart` asks it to begin when
the window opens.

## The two sides of a task

Code in the `Background` function runs away from the GUI thread. It may use
normal Go packages and slow operations, but it must not directly change a
Rosaline widget or UI variable.

`OnProgress`, `OnDone`, and functions passed to `Post` run on the GUI thread.
They may safely update ordinary application state and widgets. Rosaline
refreshes dynamic widgets after delivering them.

## Reporting progress

```go
report.Report(45, "Loading page 9")
```

`Report` accepts a percentage and optional message. Percentages are clamped to
zero through 100. It returns `false` when the task is cancelled or its window
closes, which gives a loop a simple exit:

```go
if !report.Report(percent) {
	return ctx.Err()
}
```

`OnProgress` receives both values:

```go
task.OnProgress(func(update rosaline.TaskProgress) {
	progress = update.Percent
	status = update.Message
})
```

`task.Progress()` returns the most recently delivered update when another
Rosaline callback needs to inspect it.

## Sending a result to the interface

Use `Post` for a result that is not just percentage progress:

```go
task := rosaline.Background(func(ctx context.Context, report *rosaline.TaskReporter) error {
	contents, err := os.ReadFile("notes.txt")
	if err != nil {
		return err
	}

	report.Post(func() {
		text = string(contents)
		status = "Notes loaded"
	})
	return nil
})
```

The posted function runs later on the GUI thread. This allows tasks to deliver
any normal Go value without `any`, type assertions, or a framework-specific
result object. `Post` also returns `false` after cancellation or window close.

## Starting, cancelling, and repeating

- `Start()` starts stopped work. While work is already running it does nothing.
- `Cancel()` cancels the task's standard Go `context.Context`.
- `Running()` reports whether work is active or waiting for its window to open.
- `AutoStart()` starts when the window opens. Reopened secondary windows run an
  auto-start task again.

A completed or cancelled task can be started again. Call these controls from
buttons, menus, shortcuts, timers, or other Rosaline callbacks—not from the
background function.

## Cancellation with normal Go context

For a long loop, inspect the context:

```go
select {
case <-ctx.Done():
	return ctx.Err()
default:
}
```

Pass the same context to standard-library or third-party operations that accept
one. `OnDone` receives `context.Canceled` after cancellation:

```go
task.OnDone(func(err error) {
	if errors.Is(err, context.Canceled) {
		status = "Cancelled"
	}
})
```

Cancellation is cooperative. If work ignores its context and never reports or
posts, Go cannot safely force it to stop. Its window will still close promptly,
and Rosaline will discard callbacks arriving after closure.

## Errors and panics

Return normal Go errors from the worker. Rosaline gives them to `OnDone` on the
GUI thread, where a dialog is safe:

```go
task.OnDone(func(err error) {
	if err != nil && !errors.Is(err, context.Canceled) {
		rosaline.Error("Import failed", err.Error())
	}
})
```

Rosaline also converts a panic in background work into an error instead of
crashing the GUI event loop. A panic still represents a programming bug and
should be fixed; the conversion keeps the application able to explain it.

## Secondary windows

Tasks can belong to a reusable window:

```go
window := rosaline.NewWindow(rosaline.WindowOptions{
	Title: "Importer",
	Tasks: []*rosaline.Task{importTask},
	Content: content,
})
```

Closing that window cancels its active tasks before destroying its widgets.
One `Task` can belong to only one open window at a time. Create separate tasks
when two windows need independent work.

## Common mistakes

### Forgetting App.Tasks

Creating a task does not give it a window lifetime:

```go
rosaline.RunApp(rosaline.App{
	Tasks: []*rosaline.Task{task},
	Content: content,
})
```

### Updating a widget in background work

Do not call `SetImage`, `SetValue`, `Redraw`, dialogs, or other widget methods
inside the `Background` function. Put that change inside `report.Post`,
`OnProgress`, or `OnDone`.

### Using a task for animation

Use `Animate` or `Every` for short, repeated GUI updates. Use `Background` for
work that would otherwise block the event loop.

## Go concepts used here

- goroutines managed behind a small API
- `context.Context` cancellation
- errors and `errors.Is`
- callbacks and captured variables
- safely transferring values between goroutines

See the complete [Background Bloom application](BACKGROUND_BLOOM_APPLICATION.md)
for image generation, progress, cancellation, shortcuts, saving, and error
handling in one program.
