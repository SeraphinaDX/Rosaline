// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import "testing"

func TestMenuConstruction(t *testing.T) {
	action := MenuItem("Open", nil).Shortcut("Ctrl+O")
	menu := Menu("File", action, MenuSeparator(), nil)
	bar := MenuBar(menu, nil)
	if action.shortcut != "Ctrl+O" || len(menu.entries) != 2 || len(bar.menus) != 1 {
		t.Fatal("menu constructors did not preserve their options")
	}
}
