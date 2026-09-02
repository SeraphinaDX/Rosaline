# Rosaline Architecture

Rosaline's architecture exists to keep complexity out of beginner programs.

```text
Application code
    -> Rosaline widgets and layout
    -> Rosaline drawing API
    -> pure-Go software renderer
    -> private window backend or image encoder
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

## Layout metadata

Rows, columns, grids, and stacks all mount the same private `mountedWidget`
description. Expansion stays separate on the horizontal and vertical axes, so
a widget can fill a column's width without consuming its remaining height.

Alignment is metadata rather than another native container. The receiving
layout translates start, center, end, and stretch into its geometry manager's
own positioning rules. This matters especially for `Stack`: a centered card
can sit above a canvas without an opaque alignment frame hiding the canvas.

`Card`, `Size`, and `MinSize` remain normal wrapper widgets. Cards temporarily
use the semantic surface as their children's layout background during mount;
they do not expose backend colors or handles. Preferred sizing uses a
non-propagating positioned child, while minimum sizing lets natural content
grow beyond the requested floor.

## One drawing engine

Visible canvas widgets and off-screen pictures share one pure-Go software
renderer. Paths, transforms, clipping, text, and primitive shapes therefore
produce the same pixels in a window, a PNG, and an AVIF export.

The window backend receives a finished pixel image rather than exposing its own
drawing objects. This clean boundary also means a future GPU display path can
be added without changing `DrawingCanvas` or exported artwork.

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

## Background-work boundary

Each window may own a group of `Task` values. Task work runs in ordinary Go
goroutines, but those goroutines never call the private backend. Progress,
posted functions, and completion errors enter a buffered per-run queue. A
small window-owned event-loop callback drains that queue only while work is
active, invokes application callbacks on the GUI thread, and performs one
coalesced refresh.

Each run has a generation and each task group has a window-lifetime signal.
Closing a window cancels the task context, closes that lifetime, removes the
event-loop callback, and ignores queued events from the old generation. A
worker that cooperates with cancellation exits promptly; even a misbehaving
worker cannot touch destroyed backend objects through Rosaline's reporter.

## Form values

Text boxes, text areas, checkboxes, radio groups, combo boxes, sliders, and
progress bars bind to pointers to ordinary Go variables. Rosaline flushes the
latest control values before a button callback and refreshes controls after
Rosaline events. This gives beginners predictable two-way values without
reflection, generated code, or a separate binding language. Two controls can
share one pointer—for example, a slider and progress bar—while application
logic continues to read and write a normal `float64`.

Focus order follows the order in which interactive widgets are mounted.
Rosaline handles Tab and Shift+Tab traversal privately so application code does
not need platform event handling for normal forms. Tab pages add a dynamic
visibility condition to their child controls, so keyboard traversal skips
controls on inactive pages without unmounting or recreating their Go state.
Radio groups reserve one logical focus slot even though they contain several
native buttons. Replacing their choices updates that slot to the selected
button, so dynamic groups remain predictable in Tab traversal.

## Keyboard boundary

Native key symbols and platform modifier masks are converted into Rosaline's
`KeyEvent` before application callbacks run. Named keys use stable `Key`
constants, printable keys use lowercase values, and produced text remains
separate so Shift and keyboard layouts do not blur key identity with typed
content.

Window handlers use the native toplevel binding tag. They can observe events
from focused child controls after normal control behavior, so adding an
`OnKeyDown` handler does not disable text editing. Canvas handlers bind to the
canvas itself and participate in Rosaline's existing focus traversal and
automatic-redraw lifecycle.

Standalone shortcuts and menu accelerators share one private parser. The
`Primary` modifier resolves to Control on Linux and Windows and Command on
macOS; application code remains identical. Shortcut callbacks use the same
flush, invoke, and cross-window refresh sequence as buttons and menus.

## Larger interfaces

Lists and tabs wrap native backend controls while exposing only Rosaline
widgets, strings, indices, and Go callbacks. A list owns a copy of its item
slice and normalizes selection when items change. A tab page mounts ordinary
Rosaline content inside a private page frame, preserving the same composition
and expansion rules used elsewhere.

Both controls keep their selection in normal Go fields. Backend selection
events update those fields before application callbacks run; programmatic
selection updates the backend after mounting. This keeps `Selected` useful
before, during, and after UI construction without exposing backend objects.

Tables follow the same rule while keeping their data deliberately simple.
Column headings are strings and rows are slices of strings. Rosaline copies and
normalizes that data, maps native item identifiers back to stable row indices,
and returns fresh copies to callbacks and query methods. Applications keep rich
domain objects in their own Go structs and only convert the fields they want to
display into table rows.

Trees keep hierarchy equally explicit. Each `TreeNode` owns a label, optional
string value, child pointers, and expansion state. Applications retain node
pointers when they need to select or update a branch. Rosaline privately maps
those pointers to native tree items, rejects cycles during dynamic child
replacement, and restores a safe parent selection when a selected branch is
removed. `SetChildren` lets applications lazy-load large trees without a
framework-specific model or background filesystem scan.

## Window lifecycle

One private application session owns the primary event loop and every open
secondary window. Each window receives its own mount context for controls,
focus traversal, menu shortcuts, and timers. A secondary window can therefore
close and remount without leaving native controls or scheduled callbacks
attached to its previous instance.

The session coordinates shared refreshes after callbacks, so a normal Go
variable changed in one window can update `LabelFunc` and other dynamic widgets
in every open window. Parent pointers remain ordinary Rosaline `Window`
handles; the backend mapping and native transient relationships stay private.
Closing a parent walks its registered descendants before destroying the
parent, while application quit closes the entire session through the same
lifecycle path.

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

File dialogs choose paths only. Applications continue to read and write their
own documents with normal Go packages such as `os` and `io`. Rosaline only owns
the image formats directly produced by its renderer: PNG and CGo-free AVIF.

`Render` creates a `Picture` without a window. `CanvasWidget.Picture` runs the
same drawing callback at the widget's configured size. Both routes feed the
same image encoders used by `Picture.SavePNG` and `Picture.SaveAVIF`.

## Menus and application callbacks

Buttons, menu commands, and menu shortcuts all use `func()` callbacks. Rosaline
flushes bound form values before callbacks and refreshes dynamic widgets after
them. `Quit` marks the current UI context closed before destroying the window,
so callback cleanup does not try to refresh destroyed backend widgets.

Menu and dialog backend types remain private. Scroll areas similarly expose a
single composable Rosaline widget while their canvas, scrollbars, wheel events,
and platform-specific behavior stay internal.
