# Rosaline Architecture

Rosaline's architecture exists to keep complexity out of beginner programs.

```text
Application code
    -> Rosaline widgets and layout
    -> Rosaline drawing API
    -> private platform backend
    -> Linux, Windows, or macOS
```

## Public API rule

If an ordinary feature requires a beginner to understand the backend, event
loop, renderer, or platform integration, the public API is too complicated.

## Retained widgets, direct drawing

Normal controls such as labels and buttons are retained widgets. Rosaline owns
their lifetime and layout. Canvas drawing is direct: an application receives a
`DrawingCanvas` and issues drawing operations in order.

This combination supports both ordinary utilities and drawing-heavy programs.

## Backend boundary

The current backend uses the CGo-free `modernc.org/tk9.0` package. Its platform
libraries are self-contained, so Rosaline applications build with
`CGO_ENABLED=0` and do not need a C compiler or system Tk development package.

Backend types never appear in Rosaline's public API. A future native Wayland,
software, or GPU backend can therefore be added without changing application
code.

## State and refresh

`State[T]` is deliberately small and idiomatic. Rosaline refreshes dynamic
labels after Rosaline callbacks. This is enough for early applications without
inventing a new language or reflection-based binding system.

Later milestones will add a UI-safe scheduling method for goroutines and more
granular invalidation.

## Form values

Text boxes, text areas, and checkboxes bind to pointers to ordinary Go
variables. Rosaline flushes the latest control values before a button callback
and refreshes controls after Rosaline events. This gives beginners predictable
two-way values without reflection, generated code, or a separate binding
language.

Focus order follows the order in which interactive widgets are mounted.
Rosaline handles Tab and Shift+Tab traversal privately so application code does
not need platform event handling for normal forms.

## Canvas input and redraw

Canvas mouse callbacks receive Rosaline's own `MouseEvent`; backend event types
never cross the public API boundary. Coordinates are converted directly into
the same top-left coordinate system used by drawing operations.

The canvas drawing callback is the retained description of the picture.
Rosaline clears and reruns it after mouse input. Applications update ordinary
Go drawing state inside event callbacks instead of modifying backend objects.
`Redraw` exposes the same operation to buttons and other UI callbacks.

## Images and files

Rosaline decodes image files through Go's standard `image.Image` interface and
pure-Go format packages before handing pixels to the private backend. Invalid
or missing files therefore return normal Go errors with useful filenames rather
than failing inside the event loop.

File dialogs choose paths only. Applications continue to read and write data
with normal Go packages such as `os` and `io`, keeping file formats and
application policy outside the GUI layer.

## Menus and application callbacks

Buttons, menu commands, and menu shortcuts all use `func()` callbacks. Rosaline
flushes bound form values before callbacks and refreshes dynamic widgets after
them. `Quit` marks the current UI context closed before destroying the window,
so callback cleanup does not try to refresh destroyed backend widgets.

Menu and dialog backend types remain private. Scroll areas similarly expose a
single composable Rosaline widget while their canvas, scrollbars, wheel events,
and platform-specific behavior stay internal.
