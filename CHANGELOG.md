# Changelog

Rosaline follows semantic versioning while its public API develops toward 1.0.

## 0.14.1

- Add `AskSaveChanges`, a standard three-way prompt that distinguishes saving,
  discarding, and cancelling an action.
- Fix the Notepad's New, Open, Quit, and window-close flows so No explicitly
  continues without saving while Cancel keeps the current document open.
- Keep the document open when Save As is cancelled or when a save fails.
- Add non-visual tests for every save-prompt result.
- Update the dialog, text-editing, and Notepad documentation with the safer
  pattern.

## 0.14.0

- Expand `TextArea` into a beginner-friendly document editor with `Text`,
  `SetText`, `Append`, `Clear`, and full-window `Expand` operations.
- Add automatic vertical text scrolling while keeping backend widgets private.
- Add saved-text checkpoints through `Modified` and `MarkSaved`, including
  natural clean-state restoration when undo returns to the saved content.
- Add safe undo, redo, cut, copy, paste, and select-all methods suitable for
  direct use as menu callbacks.
- Add exact wrapping `FindNext`, `ReplaceSelection`, and `ReplaceAll` methods.
- Add backend-neutral `TextPosition`, `Cursor`, and `OnCursorMove` APIs.
- Add cancellable `OnCloseRequest` callbacks to primary and secondary windows
  so applications can protect unsaved work.
- Add focused non-visual tests for text document state, replacement, searching,
  cursor-index parsing, safe unmounted commands, and cancelled window closure.
- Add a complete Text Editing guide and a polished Notepad application with
  files, menus, shortcuts, a find window, live statistics, and tested model
  logic.

## 0.13.0

- Add beginner-friendly `Grid` layouts with automatic rows, equal-width
  columns, gaps, padding, and optional equal-row expansion.
- Add `Stack` for layered interfaces, with later widgets displayed above
  earlier backgrounds.
- Add `Align` and `Center` with start, center, end, and stretch positions that
  work in windows, rows, columns, grids, and stacks without opaque wrappers.
- Add adaptive `Spring` space alongside the existing fixed `Spacer`.
- Add horizontal or vertical themed `Separator` widgets with configurable
  thickness.
- Add semantic `Card` surfaces with theme borders, padding, and expansion.
- Add `Size` and `MinSize` wrappers for preferred and minimum pixel dimensions.
- Add label `FontSize`, `Bold`, and `TextAlign` presentation methods.
- Correct row and column layout expansion to respect horizontal and vertical
  axes independently.
- Add non-visual tests for layout defaults, normalization, alignment metadata,
  axis behavior, label presentation, sizing, cards, stacks, and separators.
- Add a complete layout guide and a polished Calculator application combining
  every major presentation primitive with keyboard input and tested arithmetic
  logic.

## 0.12.0

- Add window-owned `Task` background work with the beginner-friendly
  `Background` constructor.
- Add standard `context.Context` cancellation and reusable `Start`, `Cancel`,
  `Running`, and `AutoStart` controls.
- Add percentage and message updates through `TaskReporter.Report` and
  GUI-thread `OnProgress` callbacks.
- Add `TaskReporter.Post` for safely transferring arbitrary Go results back to
  application state and mounted widgets.
- Add GUI-thread `OnDone` error handling, panic-to-error conversion, progress
  inspection, and automatic dynamic-widget refresh.
- Tie every task to `App.Tasks` or `WindowOptions.Tasks`, cancel it on window
  close, and discard callbacks from an old or closed window lifetime.
- Poll background queues only while work is active and never call the window
  backend from worker goroutines.
- Add thorough non-visual tests for progress, posting order, cancellation,
  panics, ownership, remounting, and late-callback safety.
- Add a complete background-work guide and Background Bloom application
  combining responsive image generation, progress, cancellation, shortcuts,
  dialogs, image display, and PNG export.

## 0.11.0

- Add backend-neutral `Key`, `KeyEvent`, and friendly constants for common
  navigation, editing, modifier, and function keys.
- Add window-wide `OnKeyDown` and `OnKeyUp` callbacks to `App` and
  `WindowOptions` without replacing native text or control behavior.
- Add focused canvas key-down and key-up callbacks with automatic redraw,
  initial focus, click-to-focus, Tab traversal, and a visible focus ring.
- Add standalone `Shortcut` values and the convenient `Shortcuts` collector so
  applications can bind commands without creating a menu.
- Add `Primary` shortcuts and event modifiers that mean Control on Linux and
  Windows and Command on macOS.
- Extend menu shortcuts with `Primary`, named unmodified keys such as Escape
  and F1, platform-appropriate accelerator labels, and shared parsing.
- Add non-visual tests for key conversion, platform modifiers, shortcut
  parsing and display, construction, and canvas keyboard options.
- Add a complete keyboard-input guide and Keyboard Garden application tutorial
  combining movement, modifiers, releases, shortcuts, dialogs, and PNG export.

## 0.10.0

- Add pointer-bound `RadioGroup` controls with separate visible labels and
  stored values.
- Add horizontal and vertical radio layouts, callbacks, initial focus,
  programmatic selection, and dynamic choice replacement.
- Treat each radio group as one logical Tab stop while preserving native
  keyboard movement between its choices.
- Add read-only pointer-bound `ComboBox` controls with callbacks, width,
  initial focus, programmatic selection, and dynamic option replacement.
- Add pointer-bound numeric `Slider` controls with normalized ranges, optional
  steps, sizing, direction, focus, callbacks, and programmatic control.
- Add determinate and indeterminate `ProgressBar` controls with custom
  maximums, sizing, direction, busy animation, and start/stop controls.
- Add focused non-visual tests for every new control and dynamic focus-slot
  support for controls whose native children can be replaced.
- Add four complete feature guides and a Task Settings application tutorial
  combining shared Go values, validation, changing choices, and both progress
  modes.

## 0.9.0

- Add reusable secondary windows with `NewWindow` and `WindowOptions`.
- Add safe `Show`, `Close`, `Focus`, `SetTitle`, and `IsOpen` controls.
- Focus an existing open window rather than creating accidental duplicates.
- Add explicit main-window, parent, and child relationships with cascading
  closure and automatic parent opening.
- Give each window independent content, menus, shortcuts, themes, focus order,
  and timers while refreshing shared Go state across every open window.
- Safely detach window-owned controls and timers during close and allow the
  same window handle to reopen later.
- Add non-visual window lifecycle and cross-window refresh tests.
- Add a complete Multiple Windows guide and Project Desk application.

## 0.8.0

- Add native trees built from simple `Node` values and child pointers.
- Add selection, activation, expansion, sizing, and programmatic control.
- Add safe node inspection and dynamic root or child replacement.
- Add automatic horizontal and vertical scrollbars and native keyboard tree
  navigation.
- Add thorough non-visual Tree API tests, including cycle and selection safety.
- Add a complete Trees guide with a beginner-friendly lazy-loading pattern.
- Upgrade the File Browser with Home and Filesystem folder trees that load one
  directory only when the user opens it.

## 0.7.0

- Add native multi-column tables with beginner-friendly string headings and
  `[][]string` rows.
- Add safe table selection, programmatic selection, and deep-copy inspection.
- Add mouse and keyboard selection and activation callbacks.
- Add dynamic row replacement with predictable selection preservation.
- Add configurable column widths, visible-row height, layout expansion, and
  automatic horizontal and vertical scrollbars.
- Add thorough non-visual Table API tests.
- Add a complete Tables guide.
- Add a File Browser application and tutorial using ordinary Go filesystem
  APIs, menus, shortcuts, buttons, an address box, and a dynamic table.

## 0.6.1

- Fix `invalid command true` when switching tabs by passing Tcl's required
  numeric `0` and `1` values to every widget's `-takefocus` option.
- Add a regression test covering both enabled and disabled focus options.

## 0.6.0

- Add native tabbed interfaces with `Tab`, `Tabs`, selection queries,
  programmatic selection, expansion, and change callbacks.
- Add scrollable single-selection lists with configurable size and expansion.
- Add list selection and activation callbacks for mouse and keyboard use.
- Add safe list inspection, programmatic selection, and dynamic item
  replacement.
- Keep Tab and Shift+Tab navigation out of inactive tab pages.
- Add complete Lists and Tabs guides.
- Add a Preferences application tutorial combining tabs, lists, form controls,
  dynamic labels, buttons, and a live canvas preview.

## 0.5.0

- Add reusable paths with lines, quadratic curves, cubic Bézier curves, and
  closed shapes.
- Add path filling and stroking.
- Add translate, rotate, scale, reset, and nested push/pop transforms.
- Add transformed rectangular clipping with push/pop restoration.
- Use one pure-Go renderer for visible canvases and off-screen pictures.
- Add `Render` and `CanvasWidget.Picture` for off-screen drawing.
- Add PNG and CGo-free AVIF export with quality, speed, and lossless options.
- Add AVIF loading and display to the existing image API and image viewer.
- Upgrade Paint with menus, shortcuts, unsaved-work protection, and PNG/AVIF
  saving.
- Add an advanced drawing example and complete feature and application guides.

## 0.4.1

- Add app-owned repeating timers with `Every`.
- Add one-shot delayed callbacks with `After`.
- Add frame-rate-based canvas animation timers with `Animate`.
- Add `Start`, `Stop`, `Restart`, and `Running` timer controls.
- Refresh dynamic widgets automatically after timer callbacks.
- Add a polished animation example and complete timer and animation guides.

## 0.4.0

- Add pure-Go image loading for PNG, JPEG, GIF, BMP, TIFF, and WebP.
- Add dynamic image widgets backed by Go's standard `image.Image` interface.
- Add two-axis `Scroll` containers with cross-platform wheel handling.
- Add open and save file dialogs with typed options and normalized filters.
- Add error and Boolean confirmation dialogs.
- Add native menu bars, reusable menu actions, separators, and working Ctrl,
  Shift, and Alt shortcuts.
- Add `Quit` for safely closing the application from callbacks.
- Add an image-viewer example and complete feature and application guides.

## 0.3.0

- Add backend-neutral `MouseEvent` and `MouseButton` types.
- Add canvas mouse-down, mouse-move, and mouse-up callbacks.
- Report canvas coordinates, dragging state, mouse button, and Shift, Control,
  and Alt modifiers.
- Redraw canvases automatically after mouse callbacks.
- Add public `CanvasWidget.Redraw` for drawing-state changes from buttons and
  other Rosaline callbacks.
- Add a runnable paint example and a complete canvas-input guide.

## 0.2.0

- Add pointer-bound `TextBox` with placeholder, password, width, focus,
  `OnChange`, and `OnSubmit` options.
- Add pointer-bound multiline `TextArea` with size, focus, and `OnChange`.
- Add pointer-bound `CheckBox` with focus and `OnChange`.
- Add predictable Tab and Shift+Tab navigation between interactive controls.
- Flush current form values before button callbacks.
- Add a complete forms example and feature-by-feature input documentation.
- Raise the minimum Go version to 1.25, matching the pure-Go backend.

## 0.1.0

- Initial application, layout, label, button, state, theme, dialog, and canvas
  foundation.
- Add hello, counter, and canvas examples plus the first Quick Start.
