// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import (
	"fmt"

	tk "modernc.org/tk9.0"
)

// App describes a Rosaline application window.
type App struct {
	Title   string
	Width   int
	Height  int
	Padding int
	Theme   Theme
	Menu    *AppMenuBar
	Timers  []*Timer
	Content Widget
}

var activeContext *mountContext

// Run opens a window containing content using beginner-friendly defaults.
func Run(content Widget) {
	RunApp(App{Content: content})
}

// RunApp opens an application window and runs its event loop.
func RunApp(app App) {
	if app.Title == "" {
		app.Title = "Rosaline"
	}
	if app.Width <= 0 {
		app.Width = 640
	}
	if app.Height <= 0 {
		app.Height = 420
	}
	if app.Padding < 0 {
		app.Padding = 0
	}
	if app.Padding == 0 {
		app.Padding = 18
	}
	if app.Theme == (Theme{}) {
		app.Theme = DefaultTheme
	}
	if app.Content == nil {
		app.Content = Label("This Rosaline window has no content.")
	}

	ctx := &mountContext{theme: app.Theme}
	activeContext = ctx
	mountedTimers := mountTimers(ctx, app.Timers)
	defer func() {
		ctx.closed = true
		for _, timer := range mountedTimers {
			timer.unmount(ctx)
		}
		activeContext = nil
	}()
	tk.App.WmTitle(app.Title)
	tk.WmGeometry(tk.App, fmt.Sprintf("%dx%d", app.Width, app.Height))
	tk.App.Configure(tk.Background(app.Theme.Background.String()))
	if app.Menu != nil {
		app.Menu.mount(ctx)
	}

	root := app.Content.mount(ctx, tk.App)
	options := []tk.Opt{
		tk.Fill("both"),
		tk.Padx(app.Padding),
		tk.Pady(app.Padding),
	}
	if root.expandX || root.expandY {
		options = append(options, tk.Expand(true))
	}
	tk.Pack(append([]tk.Opt{root.window}, options...)...)

	ctx.refresh()
	ctx.installFocusTraversal()
	tk.App.Center()
	if ctx.initialFocus != nil {
		tk.Focus(ctx.initialFocus)
	}
	for _, timer := range mountedTimers {
		timer.begin()
	}
	tk.App.Wait()
}

// Quit closes the Rosaline application.
func Quit() {
	if activeContext == nil {
		return
	}
	activeContext.closed = true
	tk.Destroy(tk.App)
}
