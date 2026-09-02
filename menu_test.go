// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import "testing"

func TestShortcutSequence(t *testing.T) {
	tests := []struct {
		input string
		want  string
		ok    bool
	}{
		{"Ctrl+O", "<Control-o>", true},
		{"Control+Shift+S", "<Control-Shift-S>", true},
		{"Alt+F4", "<Alt-F4>", true},
		{"Ctrl+", "", false},
		{"O", "", false},
		{"Super+O", "", false},
	}
	for _, test := range tests {
		got, ok := shortcutSequence(test.input)
		if got != test.want || ok != test.ok {
			t.Errorf("shortcutSequence(%q) = %q, %v; want %q, %v", test.input, got, ok, test.want, test.ok)
		}
	}
}

func TestMenuConstruction(t *testing.T) {
	action := MenuItem("Open", nil).Shortcut("Ctrl+O")
	menu := Menu("File", action, MenuSeparator(), nil)
	bar := MenuBar(menu, nil)
	if action.shortcut != "Ctrl+O" || len(menu.entries) != 2 || len(bar.menus) != 1 {
		t.Fatal("menu constructors did not preserve their options")
	}
}
