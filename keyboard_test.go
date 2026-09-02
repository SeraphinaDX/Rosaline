// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import (
	"testing"

	tk "modernc.org/tk9.0"
)

func TestKeyFromKeysym(t *testing.T) {
	tests := map[string]Key{
		"Return":       KeyEnter,
		"KP_Enter":     KeyEnter,
		"Escape":       KeyEscape,
		"Prior":        KeyPageUp,
		"ISO_Left_Tab": KeyTab,
		"Shift_L":      KeyShift,
		"A":            Key("a"),
		"é":            Key("é"),
		"F12":          KeyF12,
	}
	for keysym, want := range tests {
		if got := keyFromKeysym(keysym); got != want {
			t.Errorf("keyFromKeysym(%q) = %q, want %q", keysym, got, want)
		}
	}
}

func TestKeyEventConversionOnLinux(t *testing.T) {
	event := &tk.Event{
		Keysym: "A",
		Char:   "A",
		State:  tk.ModifierShift | tk.ModifierControl | tk.ModifierAlt,
	}
	got := keyEventForOS(event, "linux")
	if got.Key != Key("a") || got.Text != "A" || !got.Shift || !got.Control || !got.Alt || !got.Primary {
		t.Fatalf("unexpected Linux key event: %#v", got)
	}
	if !got.Is(Key("a")) || got.Is(KeyEnter) {
		t.Fatalf("KeyEvent.Is returned the wrong result: %#v", got)
	}
}

func TestKeyEventConversionOnDarwin(t *testing.T) {
	command := keyEventForOS(&tk.Event{Keysym: "s", State: tk.ModifierMod1}, "darwin")
	if !command.Primary || command.Alt || command.Control {
		t.Fatalf("unexpected Command event: %#v", command)
	}
	option := keyEventForOS(&tk.Event{Keysym: "x", State: tk.ModifierMod2}, "darwin")
	if !option.Alt || option.Primary {
		t.Fatalf("unexpected Option event: %#v", option)
	}
}

func TestShortcutSequence(t *testing.T) {
	tests := []struct {
		input string
		goos  string
		want  string
		ok    bool
	}{
		{"Primary+S", "linux", "<Control-s>", true},
		{"Primary+S", "darwin", "<Command-s>", true},
		{"Control+Shift+S", "linux", "<Control-Shift-S>", true},
		{"Alt+F4", "linux", "<Alt-F4>", true},
		{"Alt+F4", "darwin", "<Option-F4>", true},
		{"Ctrl+Enter", "linux", "<Control-Return>", true},
		{"Escape", "linux", "<Escape>", true},
		{"F1", "linux", "<F1>", true},
		{"Ctrl+Plus", "linux", "<Control-plus>", true},
		{"Ctrl+", "linux", "", false},
		{"O", "linux", "", false},
		{"Shift+O", "linux", "", false},
		{"Tab", "linux", "", false},
		{"Shift+Tab", "linux", "", false},
		{"Ctrl+Tab", "linux", "<Control-Tab>", true},
		{"Super+O", "linux", "<Mod4-o>", true},
		{"Super+O", "darwin", "<Command-o>", true},
		{"Command+O", "darwin", "<Command-o>", true},
		{"Command+O", "linux", "", false},
		{"Option+O", "linux", "", false},
		{"Ctrl+Mystery", "linux", "", false},
	}
	for _, test := range tests {
		got, ok := shortcutSequenceForOS(test.input, test.goos)
		if got != test.want || ok != test.ok {
			t.Errorf("shortcutSequenceForOS(%q, %q) = %q, %v; want %q, %v", test.input, test.goos, got, ok, test.want, test.ok)
		}
	}
}

func TestShortcutConstructionAndDisplay(t *testing.T) {
	called := 0
	shortcut := Shortcut(" Primary+S ", func() { called++ })
	values := Shortcuts(shortcut)
	if shortcut.Keys() != "Primary+S" || len(values) != 1 || called != 0 {
		t.Fatalf("unexpected shortcut construction: %#v, %#v", shortcut, values)
	}
	values[0] = Shortcut("Escape", nil)
	if shortcut.Keys() != "Primary+S" {
		t.Fatal("modifying Shortcuts result changed the source shortcut")
	}
	if got := shortcutDisplayForOS("Primary+Alt+S", "darwin"); got != "Command+Option+S" {
		t.Fatalf("Darwin display = %q", got)
	}
	if got := shortcutDisplayForOS("Primary+Alt+S", "linux"); got != "Ctrl+Alt+S" {
		t.Fatalf("Linux display = %q", got)
	}
	if got := shortcutDisplayForOS("Super+S", "darwin"); got != "Command+S" {
		t.Fatalf("Darwin Super display = %q", got)
	}
}

func TestNilKeyEventIsSafe(t *testing.T) {
	if got := keyEventForOS(nil, "linux"); got != (KeyEvent{}) {
		t.Fatalf("nil key event = %#v", got)
	}
}

func TestWindowKeyboardOptionsArePreserved(t *testing.T) {
	shortcut := Shortcut("Escape", func() {})
	window := NewWindow(WindowOptions{
		Shortcuts: []KeyShortcut{shortcut},
		OnKeyDown: func(KeyEvent) {},
		OnKeyUp:   func(KeyEvent) {},
	})
	options := window.resolvedOptions()
	if len(options.Shortcuts) != 1 || options.Shortcuts[0].Keys() != "Escape" || options.OnKeyDown == nil || options.OnKeyUp == nil {
		t.Fatalf("resolved options lost keyboard settings: %#v", options)
	}
}
