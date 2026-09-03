// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	tk "modernc.org/tk9.0"
)

// TextPosition identifies a place in a text area. Lines begin at 1 and
// columns begin at 0, matching the way most editors display cursor positions.
type TextPosition struct {
	Line   int
	Column int
}

// TextAreaWidget is a multiline text input bound to a Go string.
type TextAreaWidget struct {
	value        *string
	cleanValue   string
	columns      int
	lines        int
	onChange     func(string)
	onCursorMove func(TextPosition)
	focus        bool
	expand       bool
	area         *tk.TextWidget
	ctx          *mountContext
	cursor       TextPosition
}

// TextArea creates a multiline input. The area updates value as the user types.
// Pass a pointer with &, as in TextArea(&notes).
func TextArea(value *string) *TextAreaWidget {
	if value == nil {
		value = new(string)
	}
	return &TextAreaWidget{
		value:      value,
		cleanValue: *value,
		columns:    40,
		lines:      6,
		cursor:     TextPosition{Line: 1},
	}
}

// Size sets the preferred width in text columns and height in text lines.
func (t *TextAreaWidget) Size(columns, lines int) *TextAreaWidget {
	if t == nil {
		return t
	}
	if columns > 0 {
		t.columns = columns
	}
	if lines > 0 {
		t.lines = lines
	}
	return t
}

// Expand asks the text area to use all available horizontal and vertical
// space. This is useful for editors and other document-style applications.
func (t *TextAreaWidget) Expand() *TextAreaWidget {
	if t != nil {
		t.expand = true
	}
	return t
}

// OnChange runs after the user or a mounted editing method changes the text.
func (t *TextAreaWidget) OnChange(handler func(string)) *TextAreaWidget {
	if t != nil {
		t.onChange = handler
	}
	return t
}

// OnCursorMove runs after keyboard or mouse input moves the insertion cursor.
func (t *TextAreaWidget) OnCursorMove(handler func(TextPosition)) *TextAreaWidget {
	if t != nil {
		t.onCursorMove = handler
	}
	return t
}

// Focus asks Rosaline to give this text area focus when the window opens.
// If several widgets request focus, the first one wins.
func (t *TextAreaWidget) Focus() *TextAreaWidget {
	if t != nil {
		t.focus = true
	}
	return t
}

// Text returns the current text. It is also available through the pointer
// originally passed to TextArea.
func (t *TextAreaWidget) Text() string {
	if t == nil || t.value == nil {
		return ""
	}
	return *t.value
}

// SetText replaces all text and begins a fresh undo history. Call MarkSaved
// afterward when the replacement came from opening or creating a document.
func (t *TextAreaWidget) SetText(text string) {
	if t == nil || t.value == nil {
		return
	}
	old := *t.value
	*t.value = text
	if t.area != nil {
		t.replaceNative(text, true)
		if old != text {
			t.notifyChange(old)
		}
		t.refresh()
	}
}

// Append adds text at the end of the document.
func (t *TextAreaWidget) Append(text string) {
	if t == nil || t.value == nil || text == "" {
		return
	}
	if t.area == nil {
		*t.value += text
		return
	}
	t.area.Insert("end-1c", text)
	t.syncFromNative(true)
	t.refresh()
}

// Clear removes all text. When mounted, the change can be undone.
func (t *TextAreaWidget) Clear() {
	if t == nil || t.value == nil || *t.value == "" {
		return
	}
	if t.area == nil {
		*t.value = ""
		return
	}
	t.area.Replace("1.0", "end-1c", "")
	t.syncFromNative(true)
	t.refresh()
}

// Modified reports whether the current text differs from the last value
// recorded by MarkSaved. A newly created text area begins unmodified.
func (t *TextAreaWidget) Modified() bool {
	return t != nil && t.value != nil && *t.value != t.cleanValue
}

// MarkSaved records the current text as the clean saved version. Subsequent
// edits make Modified return true until the text matches this version again.
func (t *TextAreaWidget) MarkSaved() {
	if t == nil || t.value == nil {
		return
	}
	if t.area != nil {
		t.syncFromNative(false)
		t.area.SetModified(false)
	}
	t.cleanValue = *t.value
	t.refresh()
}

// Undo reverses the newest available edit. It safely does nothing when the
// text area is not mounted or no undo operation is available.
func (t *TextAreaWidget) Undo() {
	t.performNativeEdit(func() { t.area.Undo() })
}

// Redo reapplies the newest available undone edit. It safely does nothing
// when the text area is not mounted or no redo operation is available.
func (t *TextAreaWidget) Redo() {
	t.performNativeEdit(func() { t.area.Redo() })
}

// Cut copies the selected text to the platform clipboard and removes it.
func (t *TextAreaWidget) Cut() {
	t.performNativeEdit(func() { t.area.Cut() })
}

// Copy copies the selected text to the platform clipboard.
func (t *TextAreaWidget) Copy() {
	if t != nil && t.area != nil {
		t.area.Copy()
	}
}

// Paste inserts text from the platform clipboard at the cursor.
func (t *TextAreaWidget) Paste() {
	t.performNativeEdit(func() { t.area.Paste() })
}

// SelectAll selects all text in the area.
func (t *TextAreaWidget) SelectAll() {
	if t != nil && t.area != nil {
		t.area.SelectAll()
		t.area.MarkSet("insert", "end-1c")
		t.updateCursor(true)
		t.refresh()
	}
}

// Cursor returns the current insertion position. Before the text area is
// mounted, the cursor is at line 1, column 0.
func (t *TextAreaWidget) Cursor() TextPosition {
	if t == nil {
		return TextPosition{Line: 1}
	}
	return t.cursor
}

// FindNext selects the next exact occurrence of query, wrapping to the start
// when necessary. It returns false when query is empty or absent.
func (t *TextAreaWidget) FindNext(query string) bool {
	if t == nil || t.value == nil || query == "" {
		return false
	}
	if t.area == nil {
		return strings.Contains(*t.value, query)
	}

	start := t.area.Index("insert")
	if selection := t.area.TagRanges("sel"); len(selection) >= 2 {
		start = selection[1]
	}
	found := t.area.Search("--", query, start, "end")
	if found == "" && start != "1.0" {
		found = t.area.Search("--", query, "1.0", start)
	}
	if found == "" {
		return false
	}

	end := fmt.Sprintf("%s + %d chars", found, utf8.RuneCountInString(query))
	t.area.TagRemove("sel", "1.0", "end")
	t.area.TagAdd("sel", found, end)
	t.area.MarkSet("insert", end)
	t.area.See(found)
	t.updateCursor(true)
	t.refresh()
	return true
}

// ReplaceSelection replaces the selected text and reports whether a selection
// existed. The replacement becomes one normal undoable edit.
func (t *TextAreaWidget) ReplaceSelection(replacement string) bool {
	if t == nil || t.area == nil {
		return false
	}
	selection := t.area.TagRanges("sel")
	if len(selection) < 2 {
		return false
	}
	start := selection[0]
	t.area.Replace(start, selection[1], replacement)
	end := fmt.Sprintf("%s + %d chars", start, utf8.RuneCountInString(replacement))
	t.area.TagRemove("sel", "1.0", "end")
	t.area.MarkSet("insert", end)
	t.syncFromNative(true)
	t.refresh()
	return true
}

// ReplaceAll replaces every exact occurrence of old with replacement and
// returns the number of replacements. An empty old value changes nothing.
func (t *TextAreaWidget) ReplaceAll(old, replacement string) int {
	if t == nil || t.value == nil || old == "" {
		return 0
	}
	if t.area != nil {
		t.syncFromNative(false)
	}
	count := strings.Count(*t.value, old)
	if count == 0 {
		return 0
	}
	updated := strings.ReplaceAll(*t.value, old, replacement)
	if t.area == nil {
		*t.value = updated
		return count
	}
	t.area.Replace("1.0", "end-1c", updated)
	t.syncFromNative(true)
	t.refresh()
	return count
}

func (t *TextAreaWidget) mount(ctx *mountContext, parent *tk.Window) mountedWidget {
	frame := parent.Frame(
		tk.Background(ctx.theme.Background.String()),
		tk.Borderwidth(0),
	)
	var scrollbar *tk.TScrollbarWidget
	area := frame.Text(
		tk.Width(t.columns),
		tk.Height(t.lines),
		tk.Wrap("word"),
		tk.Undo(true),
		tk.Autoseparators(true),
		tk.Exportselection(0),
		tk.Background(ctx.theme.Surface.String()),
		tk.Foreground(ctx.theme.Text.String()),
		tk.Insertbackground(ctx.theme.Text.String()),
		tk.Selectbackground(ctx.theme.Primary.String()),
		tk.Relief("solid"),
		tk.Borderwidth(1),
		tk.Highlightthickness(1),
		tk.Highlightbackground(ctx.theme.Border.String()),
		tk.Highlightcolor(ctx.theme.Primary.String()),
		tk.Yscrollcommand(func(event *tk.Event) { event.ScrollSet(scrollbar) }),
		takeFocusOption(true),
	)
	scrollbar = frame.TScrollbar(
		tk.Orient("vertical"),
		takeFocusOption(false),
		tk.Command(func(event *tk.Event) { event.Yview(area) }),
	)
	tk.Pack(scrollbar, tk.Side("right"), tk.Fill("y"))
	tk.Pack(area, tk.Side("left"), tk.Fill("both"), tk.Expand(true))
	t.area = area
	t.ctx = ctx
	if *t.value != "" {
		area.Insert("1.0", *t.value)
	}
	area.EditReset()
	area.SetModified(false)
	t.updateCursor(false)

	syncAndRefresh := func() {
		t.syncFromNative(true)
		t.refresh()
	}
	tk.Bind(area.Window, "<KeyRelease>", tk.Command(syncAndRefresh))
	tk.Bind(area.Window, "<FocusOut>", tk.Command(syncAndRefresh))
	tk.Bind(area.Window, "<ButtonRelease-1>", tk.Command(syncAndRefresh))

	ctx.flushes = append(ctx.flushes, func() { t.syncFromNative(true) })
	ctx.refreshes = append(ctx.refreshes, func() {
		if t.area == nil || t.area.Text() == *t.value {
			return
		}
		t.replaceNative(*t.value, true)
	})
	ctx.addCleanup(func() {
		t.area = nil
		t.ctx = nil
	})
	ctx.addFocusable(area.Window, t.focus)

	return mountedWidget{window: frame.Window, expandX: t.expand, expandY: t.expand}
}

func (t *TextAreaWidget) replaceNative(text string, resetUndo bool) {
	if t == nil || t.area == nil {
		return
	}
	t.area.Replace("1.0", "end-1c", text)
	if resetUndo {
		t.area.EditReset()
	}
	t.area.SetModified(text != t.cleanValue)
	t.updateCursor(true)
}

func (t *TextAreaWidget) syncFromNative(notify bool) {
	if t == nil || t.area == nil || t.value == nil {
		return
	}
	old := *t.value
	*t.value = t.area.Text()
	t.updateCursor(notify)
	if notify {
		t.notifyChange(old)
	}
}

func (t *TextAreaWidget) notifyChange(old string) {
	if t != nil && t.value != nil && old != *t.value && t.onChange != nil {
		t.onChange(*t.value)
	}
}

func (t *TextAreaWidget) updateCursor(notify bool) {
	if t == nil || t.area == nil {
		return
	}
	position := parseTextPosition(t.area.Index("insert"))
	changed := position != t.cursor
	t.cursor = position
	if notify && changed && t.onCursorMove != nil {
		t.onCursorMove(position)
	}
}

func parseTextPosition(index string) TextPosition {
	lineText, columnText, ok := strings.Cut(index, ".")
	if !ok {
		return TextPosition{Line: 1}
	}
	line, lineErr := strconv.Atoi(lineText)
	column, columnErr := strconv.Atoi(columnText)
	if lineErr != nil || columnErr != nil || line < 1 || column < 0 {
		return TextPosition{Line: 1}
	}
	return TextPosition{Line: line, Column: column}
}

func (t *TextAreaWidget) performNativeEdit(edit func()) {
	if t == nil || t.area == nil || edit == nil {
		return
	}
	// Tk reports an empty undo or redo stack as an error. These public methods
	// intentionally behave like normal editor commands and safely do nothing.
	func() {
		defer func() { _ = recover() }()
		edit()
	}()
	t.syncFromNative(true)
	t.refresh()
}

func (t *TextAreaWidget) refresh() {
	if t != nil && t.ctx != nil {
		t.ctx.refresh()
	}
}
