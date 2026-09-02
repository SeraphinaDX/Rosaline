# Changelog

Rosaline follows semantic versioning while its public API develops toward 1.0.

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
