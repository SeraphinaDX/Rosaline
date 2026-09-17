// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import (
	"go/scanner"
	"go/token"
	"strings"
	"unicode/utf8"

	tk "modernc.org/tk9.0"
)

const (
	goKeywordTag = "rosaline-go-keyword"
	goStringTag  = "rosaline-go-string"
	goCommentTag = "rosaline-go-comment"
	goNumberTag  = "rosaline-go-number"
	goTypeTag    = "rosaline-go-type"
)

var goPredeclaredNames = map[string]bool{
	"any": true, "bool": true, "byte": true, "comparable": true,
	"complex64": true, "complex128": true, "error": true, "float32": true,
	"float64": true, "int": true, "int8": true, "int16": true, "int32": true,
	"int64": true, "rune": true, "string": true, "uint": true, "uint8": true,
	"uint16": true, "uint32": true, "uint64": true, "uintptr": true,
	"append": true, "cap": true, "clear": true, "close": true, "complex": true,
	"copy": true, "delete": true, "imag": true, "len": true, "make": true,
	"max": true, "min": true, "new": true, "panic": true, "print": true,
	"println": true, "real": true, "recover": true,
}

// GoCodeEditor creates a fixed-width TextArea configured for Go source code.
// It keeps the normal TextArea editing, undo, clipboard, and saved-state APIs.
func GoCodeEditor(value *string) *TextAreaWidget {
	return TextArea(value).Monospace().SetGoSyntax(true)
}

// Monospace displays text with the platform's standard fixed-width font.
func (t *TextAreaWidget) Monospace() *TextAreaWidget {
	if t == nil {
		return t
	}
	t.monospace = true
	if t.area != nil {
		t.area.Configure(tk.Font(tk.FixedFont))
	}
	return t
}

// ReadOnly prevents user edits while preserving selection, copying, scrolling,
// syntax coloring, and programmatic SetText calls.
func (t *TextAreaWidget) ReadOnly() *TextAreaWidget {
	return t.SetReadOnly(true)
}

// SetReadOnly changes whether the user can edit this text area.
func (t *TextAreaWidget) SetReadOnly(readOnly bool) *TextAreaWidget {
	if t == nil {
		return t
	}
	t.readOnly = readOnly
	if t.area != nil {
		t.applyReadOnlyState()
	}
	return t
}

// IsReadOnly reports whether user editing is disabled.
func (t *TextAreaWidget) IsReadOnly() bool {
	return t != nil && t.readOnly
}

// SetGoSyntax enables or disables lightweight Go syntax coloring. Enabling it
// also selects a fixed-width font and disables line wrapping.
func (t *TextAreaWidget) SetGoSyntax(enabled bool) *TextAreaWidget {
	if t == nil {
		return t
	}
	t.goSyntax = enabled
	if enabled {
		t.monospace = true
	}
	if t.area != nil {
		wrap := "word"
		if enabled {
			wrap = "none"
		}
		t.area.Configure(tk.Wrap(wrap))
		t.applyEditorPresentation()
	}
	return t
}

// GoTo moves the cursor to a one-based line and zero-based column, scrolls it
// into view, and focuses the editor. Invalid positions use the nearest safe
// line and column.
func (t *TextAreaWidget) GoTo(line, column int) {
	if t == nil || t.area == nil {
		return
	}
	if line < 1 {
		line = 1
	}
	if column < 0 {
		column = 0
	}
	index := textIndex(line, column)
	t.area.MarkSet("insert", index)
	t.area.TagRemove("sel", "1.0", "end")
	t.area.See(index)
	t.updateCursor(true)
	t.Focus()
}

func (t *TextAreaWidget) applyEditorPresentation() {
	if t == nil || t.area == nil {
		return
	}
	if t.monospace {
		t.area.Configure(tk.Font(tk.FixedFont))
	}
	if t.goSyntax {
		t.applyGoHighlighting()
	} else {
		for _, tag := range []string{goKeywordTag, goStringTag, goCommentTag, goNumberTag, goTypeTag} {
			t.area.TagRemove(tag, "1.0", "end")
		}
	}
	t.applyReadOnlyState()
}

func (t *TextAreaWidget) applyReadOnlyState() {
	if t == nil || t.area == nil {
		return
	}
	state := "normal"
	if t.readOnly {
		state = "disabled"
	}
	t.area.Configure(tk.State(state))
}

func (t *TextAreaWidget) applyGoHighlighting() {
	if t == nil || t.area == nil || t.value == nil {
		return
	}
	wasReadOnly := t.readOnly
	if wasReadOnly {
		t.area.Configure(tk.State("normal"))
	}
	t.area.TagConfigure(goKeywordTag, tk.Foreground(t.ctx.theme.Primary.String()))
	t.area.TagConfigure(goStringTag, tk.Foreground(t.ctx.theme.Success.String()))
	t.area.TagConfigure(goCommentTag, tk.Foreground(t.ctx.theme.Muted.String()))
	t.area.TagConfigure(goNumberTag, tk.Foreground(t.ctx.theme.Danger.String()))
	t.area.TagConfigure(goTypeTag, tk.Foreground(t.ctx.theme.Primary.String()))
	for _, tag := range []string{goKeywordTag, goStringTag, goCommentTag, goNumberTag, goTypeTag} {
		t.area.TagRemove(tag, "1.0", "end")
	}
	for _, span := range goSyntaxSpans(*t.value) {
		t.area.TagAdd(span.tag, textOffsetIndex(*t.value, span.start), textOffsetIndex(*t.value, span.end))
	}
	if wasReadOnly {
		t.area.Configure(tk.State("disabled"))
	}
}

type syntaxSpan struct {
	start int
	end   int
	tag   string
}

func goSyntaxSpans(source string) []syntaxSpan {
	files := token.NewFileSet()
	file := files.AddFile("source.go", -1, len(source))
	var lexer scanner.Scanner
	lexer.Init(file, []byte(source), nil, scanner.ScanComments)
	result := make([]syntaxSpan, 0)
	for {
		position, symbol, literal := lexer.Scan()
		if symbol == token.EOF {
			break
		}
		tag := ""
		switch {
		case symbol.IsKeyword():
			tag = goKeywordTag
		case symbol == token.STRING || symbol == token.CHAR:
			tag = goStringTag
		case symbol == token.COMMENT:
			tag = goCommentTag
		case symbol == token.INT || symbol == token.FLOAT || symbol == token.IMAG:
			tag = goNumberTag
		case symbol == token.IDENT && goPredeclaredNames[literal]:
			tag = goTypeTag
		}
		if tag == "" {
			continue
		}
		start := file.Offset(position)
		result = append(result, syntaxSpan{start: start, end: start + len(literal), tag: tag})
	}
	return result
}

func textOffsetIndex(source string, offset int) string {
	if offset < 0 {
		offset = 0
	}
	if offset > len(source) {
		offset = len(source)
	}
	prefix := source[:offset]
	line := strings.Count(prefix, "\n") + 1
	lineStart := strings.LastIndex(prefix, "\n") + 1
	return textIndex(line, utf8.RuneCountInString(prefix[lineStart:]))
}

func textIndex(line, column int) string {
	return strings.Join([]string{itoa(line), itoa(column)}, ".")
}

func itoa(value int) string {
	if value == 0 {
		return "0"
	}
	digits := [20]byte{}
	position := len(digits)
	for value > 0 {
		position--
		digits[position] = byte('0' + value%10)
		value /= 10
	}
	return string(digits[position:])
}
