# Building the Background Bloom Application

Background Bloom generates a detailed image one row at a time while its window
stays responsive. It combines background work, progress, cancellation, image
display, file saving, shortcuts, dynamic labels, and error dialogs.

Run the complete example from the project root:

```bash
CGO_ENABLED=0 go run ./examples/background
```

## Keep interface state together

The application begins with normal Go values:

```go
progress := 0.0
status := "Preparing the first bloom..."
var finished *rosaline.Picture

preview := rosaline.Image(nil).
	Placeholder("Your generated bloom will appear here.").
	Expand()
```

The progress bar shares `progress`, the dynamic status label reads `status`,
and Save uses `finished`. No special state model is needed.

## Paint pixels away from the GUI

`paintBloom` receives standard Go context and the Rosaline reporter:

```go
pixels, err := paintBloom(ctx, report)
if err != nil {
	return err
}
```

It creates `image.RGBA` pixels using Go's `image`, `color`, and `math`
packages. Each row checks cancellation and reports its percentage. This is the
slow part of the application, so it belongs in `Background` rather than a
button callback.

## Transfer the finished picture safely

Creating pixel data is safe in the worker. Changing the mounted image widget
is a GUI operation, so the example posts that part:

```go
picture := rosaline.NewPicture(pixels)
report.Post(func() {
	finished = picture
	preview.SetImage(picture)
	status = "Bloom complete — ready to save"
})
```

This small boundary is the central background-work rule: produce data in the
worker, then apply it to the interface in a posted callback.

## Connect progress and completion

```go
task.OnProgress(func(update rosaline.TaskProgress) {
	progress = update.Percent
	status = update.Message
})
```

Rosaline delivers this callback on the GUI thread and refreshes the progress
bar and status label. `OnDone` then distinguishes cancellation, an error, and
success with ordinary Go error handling.

## Make controls predictable

Create Another first checks `Running`, so repeated clicks do not overlap two
renders. Cancel calls `task.Cancel`; the worker notices through its context.
Save checks that a finished picture exists before opening a file dialog.

The same functions power both buttons and shortcuts:

```go
Shortcuts: rosaline.Shortcuts(
	rosaline.Shortcut("Primary+R", start),
	rosaline.Shortcut("Primary+S", save),
	rosaline.Shortcut("Escape", task.Cancel),
),
```

Sharing callbacks keeps mouse and keyboard behavior identical.

## Give the task to the app

```go
rosaline.RunApp(rosaline.App{
	Tasks: []*rosaline.Task{task},
	Content: content,
})
```

That one ownership declaration lets Rosaline start the auto-start task only
after the interface is mounted, deliver its callbacks on the correct event
loop, and cancel it during application shutdown.

The complete, runnable source is in
[`examples/background/main.go`](../examples/background/main.go). Read it from
top to bottom after the focused [Background Tasks guide](BACKGROUND_TASKS.md);
every operation outside Rosaline remains ordinary Go.
