// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import "testing"

func TestAppWindowOptionsPreserveLifecycleCallbacks(t *testing.T) {
	opened := false
	closed := false
	app := App{
		OnOpen:  func() { opened = true },
		OnClose: func() { closed = true },
		Content: Label("Hello"),
	}
	options := appWindowOptions(app)
	if options.OnOpen == nil || options.OnClose == nil || options.Content != app.Content {
		t.Fatal("application lifecycle callbacks or content were not preserved")
	}
	options.OnOpen()
	options.OnClose()
	if !opened || !closed {
		t.Fatal("preserved lifecycle callbacks did not run")
	}
}
