// SPDX-License-Identifier: LGPL-3.0-or-later

package main

import (
	"fmt"
	"strings"

	"github.com/SeraphinaDX/Rosaline"
)

func main() {
	draft := "Rosaline makes graphical Go applications feel approachable."
	savedDraft := draft
	status := "Ready"

	var editorWindow *rosaline.Window
	var previewWindow *rosaline.Window

	saveDraft := func() {
		savedDraft = draft
		status = fmt.Sprintf("Saved %d characters", len(savedDraft))
		editorWindow.SetTitle("Editor — Saved")
	}

	editorWindow = rosaline.NewWindow(rosaline.WindowOptions{
		Title:  "Editor",
		Width:  620,
		Height: 430,
		Parent: rosaline.MainWindow(),
		Menu: rosaline.MenuBar(
			rosaline.Menu("File",
				rosaline.MenuItem("Save", saveDraft).Shortcut("Ctrl+S"),
				rosaline.MenuSeparator(),
				rosaline.MenuItem("Close", func() {
					editorWindow.Close()
				}).Shortcut("Ctrl+W"),
			),
		),
		Content: rosaline.Column(
			rosaline.Label("Document editor").Color(rosaline.Rose),
			rosaline.Label("Changes appear in the preview window immediately."),
			rosaline.TextArea(&draft).
				Size(62, 14).
				Focus().
				OnChange(func(string) {
					status = "Draft has unsaved changes"
					editorWindow.SetTitle("Editor — Unsaved")
				}),
			rosaline.Row(
				rosaline.Button("Save", saveDraft).Primary(),
				rosaline.Button("Show preview", func() {
					previewWindow.Show()
				}),
				rosaline.Button("Close", func() {
					editorWindow.Close()
				}),
			).Gap(8),
		).Gap(10).Expand(),
		OnClose: func() {
			status = "Editor closed"
		},
	})

	previewWindow = rosaline.NewWindow(rosaline.WindowOptions{
		Title:  "Live Preview",
		Width:  520,
		Height: 300,
		Parent: editorWindow,
		Content: rosaline.Column(
			rosaline.Label("Live preview").Color(rosaline.Rose),
			rosaline.LabelFunc(func() string {
				text := strings.TrimSpace(draft)
				if text == "" {
					return "Nothing to preview yet."
				}
				return text
			}),
			rosaline.LabelFunc(func() string {
				return fmt.Sprintf("%d characters · %d saved", len(draft), len(savedDraft))
			}),
			rosaline.Button("Close preview", func() {
				previewWindow.Close()
			}),
		).Gap(12).Expand(),
		OnClose: func() {
			status = "Preview closed"
		},
	})

	var aboutWindow *rosaline.Window
	aboutWindow = rosaline.NewWindow(rosaline.WindowOptions{
		Title:  "About Project Desk",
		Width:  440,
		Height: 280,
		Parent: rosaline.MainWindow(),
		Content: rosaline.Column(
			rosaline.Label("Project Desk").Color(rosaline.Rose),
			rosaline.Label("A multiple-window application built with Rosaline."),
			rosaline.Label("Pure Go · Linux first-class · LGPL-3.0-or-later"),
			rosaline.Button("Close", func() {
				aboutWindow.Close()
			}).Primary(),
		).Gap(12),
	})

	openEditor := func() {
		status = "Editor opened"
		editorWindow.Show()
	}
	openPreview := func() {
		status = "Preview opened"
		previewWindow.Show()
	}
	openAbout := func() {
		status = "About window opened"
		aboutWindow.Show()
	}

	rosaline.RunApp(rosaline.App{
		Title:   "Rosaline Project Desk",
		Width:   720,
		Height:  430,
		Padding: 18,
		Menu: rosaline.MenuBar(
			rosaline.Menu("File",
				rosaline.MenuItem("Quit", rosaline.Quit).Shortcut("Ctrl+Q"),
			),
			rosaline.Menu("Window",
				rosaline.MenuItem("Editor", openEditor).Shortcut("Ctrl+E"),
				rosaline.MenuItem("Preview", openPreview).Shortcut("Ctrl+P"),
				rosaline.MenuItem("About", openAbout),
			),
		),
		Content: rosaline.Column(
			rosaline.Label("Project Desk").Color(rosaline.Rose),
			rosaline.Label("Open the same window repeatedly—Rosaline focuses it instead of making duplicates."),
			rosaline.Row(
				rosaline.Button("Open editor", openEditor).Primary(),
				rosaline.Button("Open preview", openPreview),
				rosaline.Button("About", openAbout),
			).Gap(8),
			rosaline.LabelFunc(func() string {
				return fmt.Sprintf(
					"Editor open: %t · Preview open: %t",
					editorWindow.IsOpen(),
					previewWindow.IsOpen(),
				)
			}),
			rosaline.LabelFunc(func() string { return status }),
			rosaline.Spacer(1, 12),
			rosaline.Label("Try this:"),
			rosaline.Label("1. Open Preview—the Editor opens automatically because it is the parent."),
			rosaline.Label("2. Type in Editor and watch both windows update."),
			rosaline.Label("3. Close Editor and its Preview closes safely too."),
		).Gap(12).Expand(),
	})
}
