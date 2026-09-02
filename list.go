// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import tk "modernc.org/tk9.0"

// ListWidget displays a scrollable list with one selected item.
type ListWidget struct {
	items      []string
	selected   int
	columns    int
	rows       int
	expand     bool
	onSelect   func(index int, value string)
	onActivate func(index int, value string)
	listbox    *tk.ListboxWidget
	ctx        *mountContext
}

// List creates a single-selection list. The first item is selected by default.
func List(items ...string) *ListWidget {
	selected := -1
	if len(items) != 0 {
		selected = 0
	}
	return &ListWidget{
		items:    append([]string(nil), items...),
		selected: selected,
		columns:  28,
		rows:     8,
	}
}

// Size sets the list's approximate width in characters and height in rows.
func (l *ListWidget) Size(columns, rows int) *ListWidget {
	if columns > 0 {
		l.columns = columns
	}
	if rows > 0 {
		l.rows = rows
	}
	if l.listbox != nil {
		l.listbox.Configure(tk.Width(l.columns), tk.Height(l.rows))
	}
	return l
}

// Expand asks the list to use available layout space.
func (l *ListWidget) Expand() *ListWidget {
	l.expand = true
	return l
}

// OnSelect runs when the selected item changes.
func (l *ListWidget) OnSelect(handler func(index int, value string)) *ListWidget {
	l.onSelect = handler
	return l
}

// OnActivate runs when an item is double-clicked or activated with Enter.
func (l *ListWidget) OnActivate(handler func(index int, value string)) *ListWidget {
	l.onActivate = handler
	return l
}

// Items returns a copy of the list's current items.
func (l *ListWidget) Items() []string {
	if l == nil {
		return nil
	}
	return append([]string(nil), l.items...)
}

// SetItems replaces every item. Selection remains at the same index when
// possible, moves to the last available item when necessary, and selects the
// first item when changing an empty list into a non-empty list.
func (l *ListWidget) SetItems(items ...string) {
	if l == nil {
		return
	}
	oldIndex, oldValue, oldOK := l.Selected()
	l.items = append([]string(nil), items...)
	l.normalizeSelection()
	if l.listbox != nil {
		l.listbox.Delete(0, tk.END)
		l.insertItems()
		l.applySelection()
		newIndex, newValue, newOK := l.Selected()
		if oldIndex != newIndex || oldValue != newValue || oldOK != newOK {
			l.notifySelection()
		}
	}
}

// Select changes the selected item. Invalid indices clear the selection.
// When the list is mounted, changing selection also runs OnSelect.
func (l *ListWidget) Select(index int) {
	if l == nil {
		return
	}
	oldSelected := l.selected
	if index < 0 || index >= len(l.items) {
		l.selected = -1
	} else {
		l.selected = index
	}
	if l.listbox != nil {
		l.applySelection()
		if l.selected != oldSelected {
			l.notifySelection()
		}
	}
}

// Selected returns the selected index and value. ok is false when the list is
// empty or has no selection.
func (l *ListWidget) Selected() (index int, value string, ok bool) {
	if l == nil || l.selected < 0 || l.selected >= len(l.items) {
		return -1, "", false
	}
	return l.selected, l.items[l.selected], true
}

func (l *ListWidget) normalizeSelection() {
	switch {
	case len(l.items) == 0:
		l.selected = -1
	case l.selected < 0:
		l.selected = 0
	case l.selected >= len(l.items):
		l.selected = len(l.items) - 1
	}
}

func (l *ListWidget) insertItems() {
	if l.listbox == nil || len(l.items) == 0 {
		return
	}
	items := make([]any, len(l.items))
	for index, item := range l.items {
		items[index] = item
	}
	l.listbox.Insert(tk.END, items...)
}

func (l *ListWidget) applySelection() {
	if l.listbox == nil {
		return
	}
	l.listbox.SelectionClear(0, tk.END)
	if l.selected < 0 || l.selected >= len(l.items) {
		return
	}
	l.listbox.SelectionSet(l.selected)
	l.listbox.Activate(l.selected)
	l.listbox.See(l.selected)
}

func (l *ListWidget) readSelection() bool {
	if l.listbox == nil {
		return false
	}
	indices := l.listbox.Curselection()
	if len(indices) == 0 || indices[0] < 0 || indices[0] >= len(l.items) {
		return false
	}
	changed := indices[0] != l.selected
	l.selected = indices[0]
	return changed
}

func (l *ListWidget) notifySelection() {
	if l.ctx != nil {
		l.ctx.flush()
	}
	if l.onSelect != nil {
		if index, value, ok := l.Selected(); ok {
			l.onSelect(index, value)
		}
	}
	if l.ctx != nil {
		l.ctx.refresh()
	}
}

func (l *ListWidget) activateSelection() {
	if l.ctx != nil {
		l.ctx.flush()
	}
	if l.onActivate != nil {
		if index, value, ok := l.Selected(); ok {
			l.onActivate(index, value)
		}
	}
	if l.ctx != nil {
		l.ctx.refresh()
	}
}

func (l *ListWidget) mount(ctx *mountContext, parent *tk.Window) mountedWidget {
	l.ctx = ctx
	frame := parent.Frame(
		tk.Background(ctx.theme.Background.String()),
		tk.Borderwidth(0),
	)
	var scrollbar *tk.TScrollbarWidget
	l.listbox = frame.Listbox(
		tk.Width(l.columns),
		tk.Height(l.rows),
		tk.Selectmode("browse"),
		tk.Exportselection(false),
		takeFocusOption(true),
		tk.Background(ctx.theme.Surface.String()),
		tk.Foreground(ctx.theme.Text.String()),
		tk.Selectbackground(ctx.theme.Primary.String()),
		tk.Selectforeground(White.String()),
		tk.Highlightcolor(ctx.theme.Primary.String()),
		tk.Highlightbackground(ctx.theme.Border.String()),
		tk.Relief("flat"),
		tk.Borderwidth(1),
		tk.Yscrollcommand(func(event *tk.Event) { event.ScrollSet(scrollbar) }),
	)
	scrollbar = frame.TScrollbar(
		tk.Orient("vertical"),
		tk.Command(func(event *tk.Event) { event.Yview(l.listbox) }),
	)

	l.insertItems()
	l.applySelection()
	tk.Pack(scrollbar, tk.Side("right"), tk.Fill("y"))
	tk.Pack(l.listbox, tk.Side("left"), tk.Fill("both"), tk.Expand(true))
	ctx.addFocusable(l.listbox.Window, false)

	tk.Bind(l.listbox, "<<ListboxSelect>>", tk.Command(func() {
		if l.readSelection() {
			l.notifySelection()
		}
	}))
	activate := tk.Command(func(event *tk.Event) {
		l.readSelection()
		l.activateSelection()
		if event != nil {
			event.SetReturnCodeBreak()
		}
	})
	tk.Bind(l.listbox, "<Double-Button-1>", activate)
	tk.Bind(l.listbox, "<Return>", activate)

	return mountedWidget{window: frame.Window, expandX: l.expand, expandY: l.expand}
}
