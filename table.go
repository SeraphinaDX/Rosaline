// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import (
	"fmt"
	"slices"
	"strings"

	tk "modernc.org/tk9.0"
)

const (
	defaultTableHeight      = 10
	defaultTableColumnWidth = 160
)

// TableWidget displays rows of text under named columns.
type TableWidget struct {
	columns    []string
	rows       [][]string
	widths     []int
	height     int
	selected   int
	expand     bool
	onSelect   func(row int, values []string)
	onActivate func(row int, values []string)
	tree       *tk.TTreeviewWidget
	itemIDs    []string
	ctx        *mountContext
}

// Table creates an empty table with the supplied column headings. Add data
// with SetRows. Empty headings receive friendly names, and a table created
// without headings gets one column named "Value".
func Table(columns ...string) *TableWidget {
	if len(columns) == 0 {
		columns = []string{"Value"}
	}
	clean := append([]string(nil), columns...)
	for index, column := range clean {
		if strings.TrimSpace(column) == "" {
			clean[index] = fmt.Sprintf("Column %d", index+1)
		}
	}
	return &TableWidget{
		columns:  clean,
		widths:   defaultTableWidths(len(clean)),
		height:   defaultTableHeight,
		selected: -1,
	}
}

// SetRows replaces every row and returns the table so setup calls can be
// chained. Each row is copied, padded with empty cells when short, and trimmed
// when it contains more values than the table has columns.
func (t *TableWidget) SetRows(rows ...[]string) *TableWidget {
	if t == nil {
		return t
	}
	oldIndex, oldValues, oldOK := t.Selected()
	t.rows = normalizeTableRows(rows, len(t.columns))
	t.normalizeSelection()
	if t.tree != nil {
		t.insertRows()
		t.applySelection()
		newIndex, newValues, newOK := t.Selected()
		if oldIndex != newIndex || oldOK != newOK || !slices.Equal(oldValues, newValues) {
			t.notifySelection()
		}
	}
	return t
}

// Columns returns a copy of the table's column headings.
func (t *TableWidget) Columns() []string {
	if t == nil {
		return nil
	}
	return append([]string(nil), t.columns...)
}

// Rows returns a deep copy of the table data.
func (t *TableWidget) Rows() [][]string {
	if t == nil {
		return nil
	}
	return copyTableRows(t.rows)
}

// ColumnWidth sets one column's preferred width in pixels. Invalid column
// indices and non-positive widths are ignored.
func (t *TableWidget) ColumnWidth(column, pixels int) *TableWidget {
	if t == nil || column < 0 || column >= len(t.widths) || pixels <= 0 {
		return t
	}
	t.widths[column] = pixels
	if t.tree != nil {
		t.tree.Column(t.columnID(column), tk.Width(pixels))
	}
	return t
}

// Height sets the preferred number of visible rows.
func (t *TableWidget) Height(rows int) *TableWidget {
	if t == nil || rows <= 0 {
		return t
	}
	t.height = rows
	if t.tree != nil {
		t.tree.Configure(tk.Height(rows))
	}
	return t
}

// Expand asks the table to use available layout space.
func (t *TableWidget) Expand() *TableWidget {
	if t != nil {
		t.expand = true
	}
	return t
}

// OnSelect runs when the selected row changes. The values slice is a copy and
// is safe for the application to keep or modify.
func (t *TableWidget) OnSelect(handler func(row int, values []string)) *TableWidget {
	if t != nil {
		t.onSelect = handler
	}
	return t
}

// OnActivate runs when a row is double-clicked or activated with Enter. The
// values slice is a copy and is safe for the application to keep or modify.
func (t *TableWidget) OnActivate(handler func(row int, values []string)) *TableWidget {
	if t != nil {
		t.onActivate = handler
	}
	return t
}

// Select changes the selected row. Invalid indices clear the selection. When
// the table is mounted, a changed selection also runs OnSelect.
func (t *TableWidget) Select(row int) {
	if t == nil {
		return
	}
	oldSelected := t.selected
	if row < 0 || row >= len(t.rows) {
		t.selected = -1
	} else {
		t.selected = row
	}
	if t.tree != nil {
		t.applySelection()
		if oldSelected != t.selected {
			t.notifySelection()
		}
	}
}

// Selected returns the selected row index and a copy of its values. ok is
// false when the table is empty or has no selection.
func (t *TableWidget) Selected() (row int, values []string, ok bool) {
	if t == nil || t.selected < 0 || t.selected >= len(t.rows) {
		return -1, nil, false
	}
	return t.selected, append([]string(nil), t.rows[t.selected]...), true
}

func defaultTableWidths(count int) []int {
	widths := make([]int, count)
	for index := range widths {
		widths[index] = defaultTableColumnWidth
	}
	return widths
}

func normalizeTableRows(rows [][]string, columns int) [][]string {
	result := make([][]string, len(rows))
	for rowIndex, row := range rows {
		result[rowIndex] = make([]string, columns)
		copy(result[rowIndex], row)
	}
	return result
}

func copyTableRows(rows [][]string) [][]string {
	result := make([][]string, len(rows))
	for index, row := range rows {
		result[index] = append([]string(nil), row...)
	}
	return result
}

func (t *TableWidget) normalizeSelection() {
	switch {
	case len(t.rows) == 0:
		t.selected = -1
	case t.selected < 0:
		t.selected = 0
	case t.selected >= len(t.rows):
		t.selected = len(t.rows) - 1
	}
}

func (t *TableWidget) columnID(index int) string {
	return fmt.Sprintf("column%d", index)
}

func (t *TableWidget) columnIDs() []string {
	ids := make([]string, len(t.columns))
	for index := range ids {
		ids[index] = t.columnID(index)
	}
	return ids
}

func (t *TableWidget) insertRows() {
	if t.tree == nil {
		return
	}
	if len(t.itemIDs) != 0 {
		items := make([]any, len(t.itemIDs))
		for index, item := range t.itemIDs {
			items[index] = item
		}
		t.tree.Delete(items...)
	}
	t.itemIDs = make([]string, 0, len(t.rows))
	for _, row := range t.rows {
		item := t.tree.Insert("", tk.END, tk.Values(row))
		t.itemIDs = append(t.itemIDs, item)
	}
}

func (t *TableWidget) applySelection() {
	if t.tree == nil {
		return
	}
	if t.selected < 0 || t.selected >= len(t.itemIDs) {
		t.tree.Selection("set", []string{})
		return
	}
	item := t.itemIDs[t.selected]
	t.tree.Selection("set", item)
	t.tree.Focus(item)
	t.tree.See(item)
}

func (t *TableWidget) readSelection() bool {
	if t.tree == nil {
		return false
	}
	selected := -1
	items := t.tree.Selection("")
	if len(items) != 0 {
		selected = slices.Index(t.itemIDs, items[0])
	}
	changed := selected != t.selected
	t.selected = selected
	return changed
}

func (t *TableWidget) notifySelection() {
	if t.ctx != nil {
		t.ctx.flush()
	}
	if t.onSelect != nil {
		if row, values, ok := t.Selected(); ok {
			t.onSelect(row, values)
		}
	}
	if t.ctx != nil {
		t.ctx.refresh()
	}
}

func (t *TableWidget) activateSelection() {
	if t.ctx != nil {
		t.ctx.flush()
	}
	if t.onActivate != nil {
		if row, values, ok := t.Selected(); ok {
			t.onActivate(row, values)
		}
	}
	if t.ctx != nil {
		t.ctx.refresh()
	}
}

func (t *TableWidget) mount(ctx *mountContext, parent *tk.Window) mountedWidget {
	t.ctx = ctx
	ctx.addCleanup(func() {
		t.tree = nil
		t.itemIDs = nil
		t.ctx = nil
	})
	frame := parent.Frame(
		tk.Background(ctx.theme.Background.String()),
		tk.Borderwidth(0),
	)

	var xScroll, yScroll *tk.TScrollbarWidget
	t.tree = frame.TTreeview(
		tk.Columns(strings.Join(t.columnIDs(), " ")),
		tk.Show("headings"),
		tk.Selectmode("browse"),
		tk.Height(t.height),
		takeFocusOption(true),
		tk.Xscrollcommand(func(event *tk.Event) { event.ScrollSet(xScroll) }),
		tk.Yscrollcommand(func(event *tk.Event) { event.ScrollSet(yScroll) }),
	)
	xScroll = frame.TScrollbar(
		tk.Orient("horizontal"),
		takeFocusOption(false),
		tk.Command(func(event *tk.Event) { event.Xview(t.tree) }),
	)
	yScroll = frame.TScrollbar(
		tk.Orient("vertical"),
		takeFocusOption(false),
		tk.Command(func(event *tk.Event) { event.Yview(t.tree) }),
	)

	for index, heading := range t.columns {
		column := t.columnID(index)
		t.tree.Column(column, tk.Anchor("w"), tk.Width(t.widths[index]), tk.Stretch(true))
		t.tree.Heading(column, tk.Txt(heading), tk.Anchor("w"))
	}
	t.insertRows()
	t.applySelection()

	tk.Grid(t.tree, tk.Row(0), tk.Column(0), tk.Sticky("nsew"))
	tk.Grid(yScroll, tk.Row(0), tk.Column(1), tk.Sticky("ns"))
	tk.Grid(xScroll, tk.Row(1), tk.Column(0), tk.Sticky("ew"))
	tk.GridRowConfigure(frame.Window, 0, tk.Weight(1))
	tk.GridColumnConfigure(frame.Window, 0, tk.Weight(1))
	ctx.addFocusable(t.tree.Window, false)

	tk.Bind(t.tree.Window, "<<TreeviewSelect>>", tk.Command(func() {
		if t.readSelection() {
			t.notifySelection()
		}
	}))
	tk.Bind(t.tree.Window, "<Double-Button-1>", tk.Command(func(event *tk.Event) {
		item := t.tree.IdentifyItem(event.X, event.Y)
		if item == "" {
			return
		}
		if index := slices.Index(t.itemIDs, item); index >= 0 {
			t.selected = index
			t.applySelection()
			t.activateSelection()
		}
	}))
	tk.Bind(t.tree.Window, "<Return>", tk.Command(func(event *tk.Event) {
		t.readSelection()
		t.activateSelection()
		event.SetReturnCodeBreak()
	}))

	return mountedWidget{window: frame.Window, expandX: t.expand, expandY: t.expand}
}
