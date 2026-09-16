// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import "testing"

func TestMenuConstruction(t *testing.T) {
	action := MenuItem("Open", nil).Shortcut("Ctrl+O")
	recent := Menu("Recent", MenuItem("Notes", nil))
	menu := Menu("File", action, recent, MenuSeparator(), nil)
	bar := MenuBar(menu, nil)
	if action.shortcut != "Ctrl+O" || len(menu.entries) != 3 || len(recent.entries) != 1 || len(bar.menus) != 1 {
		t.Fatal("menu constructors did not preserve their options")
	}
	var _ MenuEntry = recent
}

func TestContextMenuConstruction(t *testing.T) {
	action := MenuItem("Delete", func() {})
	canvas := Canvas(nil).ContextMenu(action, MenuSeparator(), nil)
	tree := Tree(Node("Root")).ContextMenu(action, nil).OnKeyDown(func(KeyEvent) {})
	if len(canvas.contextMenu) != 2 {
		t.Fatalf("canvas context menu has %d entries", len(canvas.contextMenu))
	}
	if len(tree.contextMenu) != 1 || tree.onKeyDown == nil {
		t.Fatal("tree editing callbacks were not retained")
	}
}
