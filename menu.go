// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import (
	"strings"

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
	add(*mountContext, *tk.MenuWidget, *tk.Window)
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

// Shortcut displays and binds a keyboard shortcut such as "Primary+O" or
// "Primary+Shift+S".
func (m *MenuAction) Shortcut(shortcut string) *MenuAction {
	m.shortcut = strings.TrimSpace(shortcut)
	return m
}

// MenuSeparator inserts a dividing line between menu commands.
func MenuSeparator() MenuEntry {
	return menuSeparator{}
}

func (m *AppMenuBar) mount(ctx *mountContext, window *tk.Window) {
	menuBar := window.Menu(tk.Tearoff(false))
	for _, appMenu := range m.menus {
		dropdown := menuBar.Menu(tk.Tearoff(false))
		for _, entry := range appMenu.entries {
			entry.add(ctx, dropdown, window)
		}
		menuBar.AddCascade(tk.Lbl(appMenu.text), tk.Mnu(dropdown))
	}
	window.Configure(tk.Mnu(menuBar))
}

func (m *MenuAction) add(ctx *mountContext, menu *tk.MenuWidget, window *tk.Window) {
	invoke := func() {
		ctx.flush()
		m.onClick()
		ctx.refresh()
	}
	options := []tk.Opt{tk.Lbl(m.text), tk.Command(invoke)}
	if _, ok := shortcutSequence(m.shortcut); ok {
		options = append(options, tk.Accelerator(shortcutDisplay(m.shortcut)))
	}
	menu.AddCommand(options...)
	bindShortcut(ctx, window, m.shortcut, m.onClick)
}

func (menuSeparator) add(_ *mountContext, menu *tk.MenuWidget, _ *tk.Window) {
	menu.AddSeparator()
}
