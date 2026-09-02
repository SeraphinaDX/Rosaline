// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import (
	"fmt"

	tk "modernc.org/tk9.0"
)

// App describes Rosaline's primary application window.
type App struct {
	Title     string
	Width     int
	Height    int
	Padding   int
	Theme     Theme
	Menu      *AppMenuBar
	Timers    []*Timer
	Tasks     []*Task
	Shortcuts []KeyShortcut
	OnKeyDown func(KeyEvent)
	OnKeyUp   func(KeyEvent)
	Content   Widget
}

var activeContext *mountContext

// Run opens a window containing content using beginner-friendly defaults.
func Run(content Widget) {
	RunApp(App{Content: content})
}

// RunApp opens the primary application window and runs its event loop.
func RunApp(app App) {
	options := (&Window{options: WindowOptions{
		Title: app.Title, Width: app.Width, Height: app.Height,
		Padding: app.Padding, Theme: app.Theme, Menu: app.Menu,
		Timers: app.Timers, Tasks: app.Tasks, Shortcuts: app.Shortcuts,
		OnKeyDown: app.OnKeyDown, OnKeyUp: app.OnKeyUp, Content: app.Content,
	}}).resolvedOptions()

	session := &applicationSession{
		main:    mainWindowHandle,
		windows: make(map[*Window]struct{}),
	}
	activeSession = session
	mainWindowHandle.options = options
	mainWindowHandle.open = true
	mainWindowHandle.session = session
	mainWindowHandle.native = tk.App
	ctx := &mountContext{theme: options.Theme, session: session, owner: mainWindowHandle}
	mainWindowHandle.ctx = ctx
	activeContext = ctx
	mountedTimers := mountTimers(ctx, options.Timers)
	mainWindowHandle.mountedTimers = mountedTimers
	mountedTasks := mountTasks(ctx, options.Tasks)
	mainWindowHandle.mountedTasks = mountedTasks
	defer func() {
		session.closing = true
		session.closeSecondaryWindows(false, false)
		ctx.abandon()
		mountedTasks.unmount()
		for _, timer := range mountedTimers {
			timer.unmount(ctx)
		}
		mainWindowHandle.open = false
		mainWindowHandle.session = nil
		mainWindowHandle.ctx = nil
		mainWindowHandle.native = nil
		mainWindowHandle.mountedTimers = nil
		mainWindowHandle.mountedTasks = nil
		activeSession = nil
		activeContext = nil
	}()
	tk.App.WmTitle(options.Title)
	tk.WmGeometry(tk.App, fmt.Sprintf("%dx%d", options.Width, options.Height))
	tk.App.Configure(tk.Background(options.Theme.Background.String()))
	mountWindowContent(ctx, tk.App, options)
	tk.WmProtocol(tk.App, tk.WM_DELETE_WINDOW, Quit)
	tk.App.Center()
	if ctx.initialFocus != nil {
		tk.Focus(ctx.initialFocus)
	}
	for _, timer := range mountedTimers {
		timer.begin()
	}
	mountedTasks.begin()
	tk.App.Wait()
}

// Quit closes the Rosaline application.
func Quit() {
	if activeSession == nil || activeSession.closing {
		return
	}
	activeSession.closing = true
	activeSession.closeSecondaryWindows(true, true)
	if activeSession.main != nil && activeSession.main.ctx != nil {
		ctx := activeSession.main.ctx
		activeSession.main.mountedTasks.unmount()
		for _, timer := range activeSession.main.mountedTimers {
			timer.unmount(ctx)
		}
		ctx.release()
	}
	tk.Destroy(tk.App)
}

func mountWindowContent(ctx *mountContext, parent *tk.Window, options WindowOptions) {
	if options.Menu != nil {
		options.Menu.mount(ctx, parent)
	}
	mountKeyboard(ctx, parent, options.OnKeyDown, options.OnKeyUp, options.Shortcuts)
	root := options.Content.mount(ctx, parent)
	packOptions := rootPackOptions(root, options.Padding)
	tk.Pack(append([]tk.Opt{root.window}, packOptions...)...)
	ctx.refresh()
	ctx.installFocusTraversal()
}

func rootPackOptions(root mountedWidget, padding int) []tk.Opt {
	options := []tk.Opt{tk.Padx(padding), tk.Pady(padding)}
	if root.aligned {
		options = append(options, tk.Anchor(stickyAnchor(root.sticky)), tk.Fill(stickyFill(root.sticky)))
	} else {
		options = append(options, tk.Fill("both"))
	}
	if root.expandX || root.expandY {
		options = append(options, tk.Expand(true))
	}
	return options
}
