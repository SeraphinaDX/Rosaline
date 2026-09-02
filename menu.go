// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import (
	"strings"
	"unicode/utf8"

	tk "modernc.org/tk9.0"
)

// AppMenuBar is a window's top-level menu bar.
type AppMenuBar struct {
	menus []*AppMenu
}

// AppMenu is one named drop-down menu in a menu bar.
type AppMenu struct {
	text    string
	entries []MenuEntry
}

// MenuEntry is an item or separator accepted by Menu.
type MenuEntry interface {
	add(*mountContext, *tk.MenuWidget)
}

// MenuAction is a clickable command inside a menu.
type MenuAction struct {
	text     string
	shortcut string
	onClick  func()
}

type menuSeparator struct{}

// MenuBar creates a top-level menu bar.
func MenuBar(menus ...*AppMenu) *AppMenuBar {
	clean := make([]*AppMenu, 0, len(menus))
	for _, menu := range menus {
		if menu != nil {
			clean = append(clean, menu)
		}
	}
	return &AppMenuBar{menus: clean}
}

// Menu creates one named drop-down menu.
func Menu(text string, entries ...MenuEntry) *AppMenu {
	clean := make([]MenuEntry, 0, len(entries))
	for _, entry := range entries {
		if entry != nil {
			clean = append(clean, entry)
		}
	}
	return &AppMenu{text: text, entries: clean}
}

// MenuItem creates a clickable menu command.
func MenuItem(text string, onClick func()) *MenuAction {
	if onClick == nil {
		onClick = func() {}
	}
	return &MenuAction{text: text, onClick: onClick}
}

// Shortcut displays and binds a keyboard shortcut such as "Ctrl+O" or
// "Ctrl+Shift+S".
func (m *MenuAction) Shortcut(shortcut string) *MenuAction {
	m.shortcut = strings.TrimSpace(shortcut)
	return m
}

// MenuSeparator inserts a dividing line between menu commands.
func MenuSeparator() MenuEntry {
	return menuSeparator{}
}

func (m *AppMenuBar) mount(ctx *mountContext) {
	menuBar := tk.Menu(tk.Tearoff(false))
	for _, appMenu := range m.menus {
		dropdown := menuBar.Menu(tk.Tearoff(false))
		for _, entry := range appMenu.entries {
			entry.add(ctx, dropdown)
		}
		menuBar.AddCascade(tk.Lbl(appMenu.text), tk.Mnu(dropdown))
	}
	tk.App.Configure(tk.Mnu(menuBar))
}

func (m *MenuAction) add(ctx *mountContext, menu *tk.MenuWidget) {
	invoke := func() {
		ctx.flush()
		m.onClick()
		ctx.refresh()
	}
	options := []tk.Opt{tk.Lbl(m.text), tk.Command(invoke)}
	if m.shortcut != "" {
		options = append(options, tk.Accelerator(m.shortcut))
	}
	menu.AddCommand(options...)

	if sequence, ok := shortcutSequence(m.shortcut); ok {
		tk.Bind(tk.App, sequence, tk.Command(func(event *tk.Event) {
			invoke()
			event.SetReturnCodeBreak()
		}))
	}
}

func (menuSeparator) add(_ *mountContext, menu *tk.MenuWidget) {
	menu.AddSeparator()
}

func shortcutSequence(shortcut string) (string, bool) {
	parts := strings.Split(shortcut, "+")
	if len(parts) < 2 {
		return "", false
	}

	modifiers := make([]string, 0, len(parts)-1)
	hasShift := false
	for _, part := range parts[:len(parts)-1] {
		switch strings.ToLower(strings.TrimSpace(part)) {
		case "ctrl", "control":
			modifiers = append(modifiers, "Control")
		case "shift":
			modifiers = append(modifiers, "Shift")
			hasShift = true
		case "alt":
			modifiers = append(modifiers, "Alt")
		default:
			return "", false
		}
	}

	key := strings.TrimSpace(parts[len(parts)-1])
	if key == "" {
		return "", false
	}
	if utf8.RuneCountInString(key) == 1 {
		if hasShift {
			key = strings.ToUpper(key)
		} else {
			key = strings.ToLower(key)
		}
	}
	return "<" + strings.Join(append(modifiers, key), "-") + ">", true
}
