// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import (
	"fmt"

	tk "modernc.org/tk9.0"
)

// WindowOptions describes a secondary application window. Every field is
// optional; Rosaline supplies the same friendly defaults used by RunApp.
type WindowOptions struct {
	Title     string
	Width     int
	Height    int
	Padding   int
	Theme     Theme
	Menu      *AppMenuBar
	Timers    []*Timer
	Shortcuts []KeyShortcut
	OnKeyDown func(KeyEvent)
	OnKeyUp   func(KeyEvent)
	Content   Widget
	Parent    *Window
	OnClose   func()
}

// Window is a reusable secondary-window handle. Create one with NewWindow.
type Window struct {
	options       WindowOptions
	primary       bool
	open          bool
	closing       bool
	session       *applicationSession
	ctx           *mountContext
	native        *tk.Window
	mountedTimers []*Timer
}

type applicationSession struct {
	main       *Window
	windows    map[*Window]struct{}
	closing    bool
	refreshing bool
}

var (
	mainWindowHandle = &Window{primary: true}
	activeSession    *applicationSession

	showWindowBackend  = func(window *Window) { window.mount() }
	closeWindowBackend = func(window *Window) {
		if window.native != nil {
			tk.Destroy(window.native)
		}
	}
	focusWindowBackend = func(window *Window) {
		if window.native == nil {
			return
		}
		tk.WmDeiconify(window.native)
		if window.primary {
			window.native.Raise(nil)
		} else {
			window.native.Raise(tk.App)
		}
		tk.Focus(window.native)
	}
	setWindowTitleBackend = func(window *Window, title string) {
		if window.native != nil {
			window.native.WmTitle(title)
		}
	}
)

// NewWindow creates a reusable secondary window. Creation does not display it;
// call Show from a Rosaline callback while the application is running.
func NewWindow(options WindowOptions) *Window {
	return &Window{options: options}
}

// MainWindow returns the primary application-window handle. It can be used as
// WindowOptions.Parent when a secondary window should belong to the main one.
func MainWindow() *Window {
	return mainWindowHandle
}

// Show opens the window. If it is already open, Show raises and focuses the
// existing window rather than creating a duplicate.
func (w *Window) Show() *Window {
	if w == nil || activeSession == nil || activeSession.closing || w.closing {
		return w
	}
	if w.primary {
		focusWindowBackend(w)
		return w
	}
	if w.open {
		focusWindowBackend(w)
		return w
	}
	if parent := w.options.Parent; parent != nil {
		if parent.closing {
			return w
		}
		if !parent.primary && !parent.IsOpen() {
			parent.Show()
			if !parent.IsOpen() {
				return w
			}
		}
	}

	w.session = activeSession
	w.open = true
	w.session.windows[w] = struct{}{}
	showWindowBackend(w)
	return w
}

// Close closes the window and its open child windows. The same Window can be
// shown again later. Closing MainWindow closes the entire application.
func (w *Window) Close() {
	if w == nil {
		return
	}
	if w.primary {
		Quit()
		return
	}
	w.close(true, true)
}

// Focus raises and focuses an open window. Closed windows are unchanged.
func (w *Window) Focus() *Window {
	if w != nil && w.open {
		focusWindowBackend(w)
	}
	return w
}

// SetTitle changes the title now and keeps it when the window is reopened.
// Empty titles use "Rosaline".
func (w *Window) SetTitle(title string) *Window {
	if w == nil {
		return w
	}
	if title == "" {
		title = "Rosaline"
	}
	w.options.Title = title
	if w.open {
		setWindowTitleBackend(w, title)
	}
	return w
}

// IsOpen reports whether the window is currently displayed.
func (w *Window) IsOpen() bool {
	return w != nil && w.open
}

func (w *Window) resolvedOptions() WindowOptions {
	options := w.options
	if options.Title == "" {
		options.Title = "Rosaline"
	}
	if options.Width <= 0 {
		options.Width = 640
	}
	if options.Height <= 0 {
		options.Height = 420
	}
	if options.Padding <= 0 {
		options.Padding = 18
	}
	if options.Theme == (Theme{}) {
		switch {
		case options.Parent != nil && options.Parent.ctx != nil:
			options.Theme = options.Parent.ctx.theme
		case w.session != nil && w.session.main != nil && w.session.main.ctx != nil:
			options.Theme = w.session.main.ctx.theme
		default:
			options.Theme = DefaultTheme
		}
	}
	if options.Content == nil {
		options.Content = Label("This Rosaline window has no content.")
	}
	return options
}

func (w *Window) mount() {
	options := w.resolvedOptions()
	top := tk.App.Toplevel()
	w.native = top.Window
	tk.WmWithdraw(w.native)
	w.native.WmTitle(options.Title)
	tk.WmGeometry(w.native, fmt.Sprintf("%dx%d", options.Width, options.Height))
	w.native.Configure(tk.Background(options.Theme.Background.String()))
	if parent := w.options.Parent; parent != nil && parent.open && parent.native != nil {
		tk.WmTransient(w.native, parent.native)
	}

	w.ctx = &mountContext{theme: options.Theme, session: w.session, owner: w}
	w.mountedTimers = mountTimers(w.ctx, options.Timers)
	mountWindowContent(w.ctx, w.native, options)
	tk.WmProtocol(w.native, tk.WM_DELETE_WINDOW, w.Close)
	w.native.Center()
	tk.WmDeiconify(w.native)
	if w.ctx.initialFocus != nil {
		tk.Focus(w.ctx.initialFocus)
	}
	for _, timer := range w.mountedTimers {
		timer.begin()
	}
}

func (w *Window) close(destroy, notify bool) {
	if w == nil || !w.open || w.closing {
		return
	}
	w.closing = true
	w.open = false
	session := w.session

	if session != nil {
		children := make([]*Window, 0)
		for candidate := range session.windows {
			if candidate != w && candidate.options.Parent == w && candidate.open {
				children = append(children, candidate)
			}
		}
		for _, child := range children {
			child.close(destroy, notify)
		}
	}

	if destroy {
		for _, timer := range w.mountedTimers {
			timer.unmount(w.ctx)
		}
		if w.ctx != nil {
			w.ctx.release()
		}
	} else {
		if w.ctx != nil {
			w.ctx.abandon()
		}
		for _, timer := range w.mountedTimers {
			timer.unmount(w.ctx)
		}
	}
	if destroy {
		closeWindowBackend(w)
	}
	if session != nil {
		delete(session.windows, w)
	}
	w.native = nil
	w.ctx = nil
	w.mountedTimers = nil
	w.session = nil
	if notify && w.options.OnClose != nil {
		w.options.OnClose()
	}
	w.closing = false
	if session != nil && !session.closing {
		session.refreshAll()
	}
}

func (s *applicationSession) refreshAll() {
	if s == nil || s.closing || s.refreshing {
		return
	}
	s.refreshing = true
	defer func() { s.refreshing = false }()
	if s.main != nil && s.main.ctx != nil {
		s.main.ctx.refreshLocal()
	}
	for window := range s.windows {
		if window.ctx != nil {
			window.ctx.refreshLocal()
		}
	}
}

func (s *applicationSession) closeSecondaryWindows(destroy, notify bool) {
	if s == nil {
		return
	}
	windows := make([]*Window, 0, len(s.windows))
	for window := range s.windows {
		windows = append(windows, window)
	}
	for _, window := range windows {
		window.close(destroy, notify)
	}
}
