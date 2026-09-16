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

// AppMenu is one named drop-down menu or nested submenu.
type AppMenu struct {
	text    string
	entries []MenuEntry
}

// MenuEntry is an item, separator, or submenu accepted by Menu.
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

// Menu creates one named drop-down menu. A Menu can also be passed to another
// Menu to create a cascading submenu.
func Menu(text string, entries ...MenuEntry) *AppMenu {
	return &AppMenu{text: text, entries: cleanMenuEntries(entries)}
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
		appMenu.add(ctx, menuBar, window)
	}
	window.Configure(tk.Mnu(menuBar))
}

// add also makes Menu usable as an entry inside another Menu, producing a
// native cascading submenu with the same item and separator API.
func (m *AppMenu) add(ctx *mountContext, menu *tk.MenuWidget, window *tk.Window) {
	if m == nil {
		return
	}
	dropdown := menu.Menu(tk.Tearoff(false))
	for _, entry := range m.entries {
		entry.add(ctx, dropdown, window)
	}
	menu.AddCascade(tk.Lbl(m.text), tk.Mnu(dropdown))
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

func cleanMenuEntries(entries []MenuEntry) []MenuEntry {
	clean := make([]MenuEntry, 0, len(entries))
	for _, entry := range entries {
		if entry != nil {
			clean = append(clean, entry)
		}
	}
	return clean
}

func mountPopupMenu(ctx *mountContext, parent *tk.Window, entries []MenuEntry) *tk.MenuWidget {
	if ctx == nil || parent == nil || len(entries) == 0 {
		return nil
	}
	menu := parent.Menu(tk.Tearoff(false))
	for _, entry := range entries {
		entry.add(ctx, menu, parent)
	}
	return menu
}
