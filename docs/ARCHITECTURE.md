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
