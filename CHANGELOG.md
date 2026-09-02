# Changelog

Rosaline follows semantic versioning while its public API develops toward 1.0.

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
