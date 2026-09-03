// SPDX-License-Identifier: LGPL-3.0-or-later

// The Notepad example combines Rosaline's text-editing and application APIs.
// Run it with: go run ./examples/notepad
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/SeraphinaDX/Rosaline"
)

type document struct {
	text string
	path string
}

func (d *document) name() string {
	if d.path == "" {
		return "Untitled"
	}
	return filepath.Base(d.path)
}

func (d *document) title(modified bool) string {
	mark := ""
	if modified {
		mark = " *"
	}
	return d.name() + mark + " — Rosaline Notepad"
}

func textStats(text string) (words, lines int) {
	words = len(strings.Fields(text))
	lines = 1
	if text != "" {
		lines += strings.Count(text, "\n")
	}
	return words, lines
}

func main() {
	document := &document{}
	status := "Ready"

	var editor *rosaline.TextAreaWidget
	updateTitle := func() {
		rosaline.MainWindow().SetTitle(document.title(editor != nil && editor.Modified()))
	}

	editor = rosaline.TextArea(&document.text).
		Size(88, 30).
		Expand().
		Focus().
		OnChange(func(string) {
			status = "Unsaved changes"
			updateTitle()
		})

	var save, saveAs func() bool
	confirmChanges := func(action string) bool {
		if !editor.Modified() {
			return true
		}
		switch rosaline.AskSaveChanges(
			"Unsaved changes",
			"Save changes to "+document.name()+" before "+action+"?",
		) {
		case rosaline.SaveChanges:
			return save != nil && save()
		case rosaline.DiscardChanges:
			return true
		default:
			return false
		}
	}

	newDocument := func() {
		if !confirmChanges("creating a new document") {
			return
		}
		editor.SetText("")
		editor.MarkSaved()
		document.path = ""
		status = "New document"
		updateTitle()
	}

	openDocument := func() {
		path, ok := rosaline.OpenFileDialog(rosaline.FileDialogOptions{
			Title: "Open Text File",
			Filters: []rosaline.FileFilter{
				{Name: "Text files", Extensions: []string{".txt", ".md", ".go"}},
				{Name: "All files", Extensions: []string{"*"}},
			},
		})
		if !ok {
			return
		}
		contents, err := os.ReadFile(path)
		if err != nil {
			rosaline.Error("Could not open file", err.Error())
			return
		}
		if !confirmChanges("opening " + filepath.Base(path)) {
			return
		}
		editor.SetText(string(contents))
		editor.MarkSaved()
		document.path = path
		status = "Opened " + filepath.Base(path)
		updateTitle()
	}

	save = func() bool {
		if document.path == "" {
			return saveAs()
		}
		if err := os.WriteFile(document.path, []byte(editor.Text()), 0o644); err != nil {
			rosaline.Error("Could not save file", err.Error())
			status = "Save failed"
			return false
		}
		editor.MarkSaved()
		status = "Saved " + filepath.Base(document.path)
		updateTitle()
		return true
	}
	saveAs = func() bool {
		path, ok := rosaline.SaveFileDialog(rosaline.FileDialogOptions{
			Title:            "Save Text File",
			InitialFile:      document.name(),
			DefaultExtension: ".txt",
			Filters: []rosaline.FileFilter{
				{Name: "Text files", Extensions: []string{".txt"}},
				{Name: "Markdown", Extensions: []string{".md"}},
				{Name: "All files", Extensions: []string{"*"}},
			},
		})
		if !ok {
			return false
		}
		oldPath := document.path
		document.path = path
		if !save() {
			document.path = oldPath
			updateTitle()
			return false
		}
		return true
	}

	findText := ""
	replacement := ""
	findNext := func() {
		if findText == "" {
			status = "Enter text to find"
			return
		}
		if editor.FindNext(findText) {
			status = "Found “" + findText + "”"
		} else {
			status = "No match for “" + findText + "”"
		}
	}
	replaceOne := func() {
		if editor.ReplaceSelection(replacement) {
			findNext()
			status = "Replaced one match"
			return
		}
		findNext()
	}
	replaceEvery := func() {
		count := editor.ReplaceAll(findText, replacement)
		status = fmt.Sprintf("Replaced %d matches", count)
		updateTitle()
	}

	var findWindow *rosaline.Window
	findWindow = rosaline.NewWindow(rosaline.WindowOptions{
		Title:  "Find and Replace",
		Width:  430,
		Height: 235,
		Parent: rosaline.MainWindow(),
		Content: rosaline.Column(
			rosaline.Label("Find").Bold(),
			rosaline.TextBox(&findText).Placeholder("Text to find").OnSubmit(func(string) {
				findNext()
			}).Focus(),
			rosaline.Label("Replace with").Bold(),
			rosaline.TextBox(&replacement).Placeholder("Replacement text"),
			rosaline.Row(
				rosaline.Button("Find Next", findNext).Primary(),
				rosaline.Button("Replace", replaceOne),
				rosaline.Button("Replace All", replaceEvery),
				rosaline.Button("Close", func() { findWindow.Close() }),
			).Gap(8),
		).Gap(8),
	})

	showAbout := func() {
		rosaline.Message(
			"About Rosaline Notepad",
			"A small text editor built entirely with Rosaline.\n\nIt demonstrates files, menus, shortcuts, editing, find and replace, cursor position, and safe closing.",
		)
	}

	menu := rosaline.MenuBar(
		rosaline.Menu("File",
			rosaline.MenuItem("New", newDocument).Shortcut("Primary+N"),
			rosaline.MenuItem("Open…", openDocument).Shortcut("Primary+O"),
			rosaline.MenuSeparator(),
			rosaline.MenuItem("Save", func() { save() }).Shortcut("Primary+S"),
			rosaline.MenuItem("Save As…", func() { saveAs() }).Shortcut("Primary+Shift+S"),
			rosaline.MenuSeparator(),
			rosaline.MenuItem("Quit", rosaline.Quit).Shortcut("Primary+Q"),
		),
		rosaline.Menu("Edit",
			rosaline.MenuItem("Undo", editor.Undo).Shortcut("Primary+Z"),
			rosaline.MenuItem("Redo", editor.Redo).Shortcut("Primary+Shift+Z"),
			rosaline.MenuSeparator(),
			rosaline.MenuItem("Cut", editor.Cut).Shortcut("Primary+X"),
			rosaline.MenuItem("Copy", editor.Copy).Shortcut("Primary+C"),
			rosaline.MenuItem("Paste", editor.Paste).Shortcut("Primary+V"),
			rosaline.MenuItem("Select All", editor.SelectAll).Shortcut("Primary+A"),
			rosaline.MenuSeparator(),
			rosaline.MenuItem("Find and Replace…", func() { findWindow.Show() }).Shortcut("Primary+F"),
			rosaline.MenuItem("Find Next", findNext).Shortcut("F3"),
		),
		rosaline.Menu("Help",
			rosaline.MenuItem("About", showAbout).Shortcut("F1"),
		),
	)

	theme := rosaline.DefaultTheme
	theme.Background = rosaline.Hex("#f7edf4")
	theme.Surface = rosaline.Hex("#fffafd")
	theme.Primary = rosaline.Hex("#b83f77")
	theme.Border = rosaline.Hex("#d8afc5")

	rosaline.RunApp(rosaline.App{
		Title:          document.title(false),
		Width:          900,
		Height:         680,
		Padding:        12,
		Theme:          theme,
		Menu:           menu,
		OnCloseRequest: func() bool { return confirmChanges("closing") },
		Content: rosaline.Column(
			rosaline.Row(
				rosaline.Label("ROSALINE").Bold().Color(theme.Primary),
				rosaline.Label("NOTEPAD").Color(theme.Muted),
				rosaline.Spring(),
				rosaline.LabelFunc(document.name).Bold(),
			).Gap(10),
			rosaline.Card(editor).Padding(5).Expand(),
			rosaline.Row(
				rosaline.LabelFunc(func() string { return status }).Color(theme.Muted),
				rosaline.Spring(),
				rosaline.LabelFunc(func() string {
					words, lines := textStats(editor.Text())
					position := editor.Cursor()
					return fmt.Sprintf("%d words · %d lines · Ln %d, Col %d",
						words, lines, position.Line, position.Column+1)
				}).Color(theme.Muted),
			).Gap(10),
		).Gap(9).Expand(),
	})
}
