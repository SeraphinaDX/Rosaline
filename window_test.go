// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import "testing"

func installFakeWindowSession(t *testing.T) *applicationSession {
	t.Helper()
	originalSession := activeSession
	originalShow := showWindowBackend
	originalClose := closeWindowBackend
	originalFocus := focusWindowBackend
	originalTitle := setWindowTitleBackend

	main := &Window{primary: true, open: true}
	session := &applicationSession{main: main, windows: make(map[*Window]struct{})}
	main.session = session
	main.ctx = &mountContext{theme: DefaultTheme, session: session, owner: main}
	activeSession = session

	t.Cleanup(func() {
		activeSession = originalSession
		showWindowBackend = originalShow
		closeWindowBackend = originalClose
		focusWindowBackend = originalFocus
		setWindowTitleBackend = originalTitle
	})
	return session
}

func TestWindowShowFocusCloseAndReopen(t *testing.T) {
	session := installFakeWindowSession(t)
	shows := 0
	closes := 0
	focuses := 0
	titles := 0
	closedCallbacks := 0
	showWindowBackend = func(*Window) { shows++ }
	closeWindowBackend = func(*Window) { closes++ }
	focusWindowBackend = func(*Window) { focuses++ }
	setWindowTitleBackend = func(_ *Window, title string) {
		if title != "Renamed" {
			t.Fatalf("title = %q, want Renamed", title)
		}
		titles++
	}

	window := NewWindow(WindowOptions{OnClose: func() { closedCallbacks++ }})
	window.Show()
	if !window.IsOpen() || shows != 1 {
		t.Fatalf("first Show = open %v, backend calls %d", window.IsOpen(), shows)
	}
	if _, ok := session.windows[window]; !ok {
		t.Fatal("open window was not registered with its application")
	}

	window.Show()
	if shows != 1 || focuses != 1 {
		t.Fatalf("second Show = %d creates, %d focuses; want 1, 1", shows, focuses)
	}
	window.SetTitle("Renamed")
	if titles != 1 || window.options.Title != "Renamed" {
		t.Fatal("SetTitle did not update the open window and retained options")
	}

	window.Close()
	if window.IsOpen() || closes != 1 || closedCallbacks != 1 {
		t.Fatalf("Close = open %v, closes %d, callbacks %d", window.IsOpen(), closes, closedCallbacks)
	}
	window.Close()
	if closes != 1 || closedCallbacks != 1 {
		t.Fatal("closing an already closed window should have no effect")
	}

	window.Show()
	if !window.IsOpen() || shows != 2 {
		t.Fatalf("reopened window = open %v, backend calls %d", window.IsOpen(), shows)
	}
}

func TestShowingChildOpensParentAndClosingParentClosesChild(t *testing.T) {
	installFakeWindowSession(t)
	shows := 0
	closes := 0
	showWindowBackend = func(*Window) { shows++ }
	closeWindowBackend = func(*Window) { closes++ }
	focusWindowBackend = func(*Window) {}

	parent := NewWindow(WindowOptions{Title: "Parent"})
	child := NewWindow(WindowOptions{Title: "Child", Parent: parent})
	child.Show()
	if !parent.IsOpen() || !child.IsOpen() || shows != 2 {
		t.Fatalf("child Show = parent %v, child %v, creates %d", parent.IsOpen(), child.IsOpen(), shows)
	}

	parent.Close()
	if parent.IsOpen() || child.IsOpen() || closes != 2 {
		t.Fatalf("parent Close = parent %v, child %v, closes %d", parent.IsOpen(), child.IsOpen(), closes)
	}
}

func TestWindowShowOutsideRunningApplicationIsSafe(t *testing.T) {
	originalSession := activeSession
	activeSession = nil
	t.Cleanup(func() { activeSession = originalSession })

	window := NewWindow(WindowOptions{})
	window.Show().Focus()
	window.Close()
	if window.IsOpen() {
		t.Fatal("window opened without a running application")
	}
}

func TestSecondaryWindowInheritsApplicationThemeAndDefaults(t *testing.T) {
	session := installFakeWindowSession(t)
	theme := DefaultTheme
	theme.Primary = Hex("#b88cff")
	session.main.ctx.theme = theme

	window := NewWindow(WindowOptions{})
	window.session = session
	options := window.resolvedOptions()
	if options.Title != "Rosaline" || options.Width != 640 || options.Height != 420 || options.Padding != 18 {
		t.Fatalf("defaults = %q %dx%d padding %d", options.Title, options.Width, options.Height, options.Padding)
	}
	if options.Theme != theme || options.Content == nil {
		t.Fatal("window did not inherit the application theme and default content")
	}
}

func TestWindowEventsRefreshEveryOpenWindow(t *testing.T) {
	session := installFakeWindowSession(t)
	mainRefreshes := 0
	childRefreshes := 0
	session.main.ctx.refreshes = []func(){func() { mainRefreshes++ }}

	child := &Window{open: true, session: session}
	child.ctx = &mountContext{
		session:   session,
		owner:     child,
		refreshes: []func(){func() { childRefreshes++ }},
	}
	session.windows[child] = struct{}{}
	child.ctx.refresh()

	if mainRefreshes != 1 || childRefreshes != 1 {
		t.Fatalf("refreshes = main %d, child %d; want 1 each", mainRefreshes, childRefreshes)
	}
}
