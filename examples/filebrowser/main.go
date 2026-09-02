// SPDX-License-Identifier: LGPL-3.0-or-later

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/SeraphinaDX/Rosaline"
)

type browserEntry struct {
	path      string
	name      string
	kind      string
	size      string
	modified  string
	directory bool
}

func main() {
	workingDirectory, err := os.Getwd()
	if err != nil {
		workingDirectory = "."
	}

	currentPath := workingDirectory
	loadedPath := workingDirectory
	status := "Loading files…"
	var entries []browserEntry

	table := rosaline.Table("Name", "Type", "Size", "Modified").
		ColumnWidth(0, 360).
		ColumnWidth(1, 130).
		ColumnWidth(2, 100).
		ColumnWidth(3, 180).
		Height(18).
		Expand()

	var loadDirectory func(string)
	loadDirectory = func(path string) {
		absolute, err := filepath.Abs(path)
		if err != nil {
			rosaline.Error("Could not open folder", err.Error())
			currentPath = loadedPath
			return
		}
		absolute = filepath.Clean(absolute)

		directoryEntries, err := os.ReadDir(absolute)
		if err != nil {
			rosaline.Error("Could not open folder", err.Error())
			currentPath = loadedPath
			return
		}

		sort.SliceStable(directoryEntries, func(left, right int) bool {
			if directoryEntries[left].IsDir() != directoryEntries[right].IsDir() {
				return directoryEntries[left].IsDir()
			}
			return strings.ToLower(directoryEntries[left].Name()) < strings.ToLower(directoryEntries[right].Name())
		})

		entries = nil
		if parent := filepath.Dir(absolute); parent != absolute {
			entries = append(entries, browserEntry{
				path:      parent,
				name:      "..",
				kind:      "Parent folder",
				size:      "—",
				modified:  "—",
				directory: true,
			})
		}

		for _, entry := range directoryEntries {
			item := browserEntry{
				path:      filepath.Join(absolute, entry.Name()),
				name:      entry.Name(),
				directory: entry.IsDir(),
			}
			info, infoErr := entry.Info()
			if infoErr != nil {
				item.kind = "Unavailable"
				item.size = "—"
				item.modified = "—"
			} else {
				item.kind = fileKind(info)
				if info.IsDir() {
					item.size = "—"
				} else {
					item.size = formatSize(info.Size())
				}
				item.modified = info.ModTime().Format("2006-01-02 15:04")
			}
			entries = append(entries, item)
		}

		rows := make([][]string, len(entries))
		for index, entry := range entries {
			rows[index] = []string{entry.name, entry.kind, entry.size, entry.modified}
		}
		table.SetRows(rows...)
		loadedPath = absolute
		currentPath = absolute
		status = fmt.Sprintf("%d items in %s", len(directoryEntries), absolute)
	}

	table.OnSelect(func(row int, values []string) {
		if row >= 0 && row < len(entries) {
			status = entries[row].path
		}
	}).OnActivate(func(row int, values []string) {
		if row < 0 || row >= len(entries) {
			return
		}
		entry := entries[row]
		if entry.directory {
			loadDirectory(entry.path)
			return
		}
		rosaline.Message("File details", fmt.Sprintf(
			"%s\n\nType: %s\nSize: %s\nModified: %s\n\n%s",
			entry.name,
			entry.kind,
			entry.size,
			entry.modified,
			entry.path,
		))
	})

	goUp := func() {
		loadDirectory(filepath.Dir(loadedPath))
	}
	goHome := func() {
		home, err := os.UserHomeDir()
		if err != nil {
			rosaline.Error("Could not find home folder", err.Error())
			return
		}
		loadDirectory(home)
	}
	refresh := func() {
		loadDirectory(loadedPath)
	}

	loadDirectory(workingDirectory)

	rosaline.RunApp(rosaline.App{
		Title:   "Rosaline File Browser",
		Width:   940,
		Height:  620,
		Padding: 12,
		Menu: rosaline.MenuBar(
			rosaline.Menu("File",
				rosaline.MenuItem("Quit", rosaline.Quit).Shortcut("Ctrl+Q"),
			),
			rosaline.Menu("Go",
				rosaline.MenuItem("Parent folder", goUp).Shortcut("Alt+Up"),
				rosaline.MenuItem("Home folder", goHome).Shortcut("Alt+Home"),
				rosaline.MenuItem("Refresh", refresh).Shortcut("Ctrl+R"),
			),
			rosaline.Menu("Help",
				rosaline.MenuItem("About", func() {
					rosaline.Message("About", "File Browser built with Rosaline v0.7.0")
				}),
			),
		),
		Content: rosaline.Column(
			rosaline.Label("File Browser").Color(rosaline.Rose),
			rosaline.Row(
				rosaline.Button("Up", goUp),
				rosaline.Button("Home", goHome),
				rosaline.Button("Refresh", refresh),
				rosaline.TextBox(&currentPath).
					Width(72).
					OnSubmit(loadDirectory),
			).Gap(8),
			table,
			rosaline.LabelFunc(func() string { return status }),
		).Gap(10).Expand(),
	})
}

func fileKind(info os.FileInfo) string {
	switch {
	case info.IsDir():
		return "Folder"
	case info.Mode()&os.ModeSymlink != 0:
		return "Symbolic link"
	case info.Mode().IsRegular():
		return "File"
	default:
		return "Special file"
	}
}

func formatSize(bytes int64) string {
	const unit = int64(1024)
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	value := float64(bytes)
	labels := []string{"KB", "MB", "GB", "TB"}
	for _, label := range labels {
		value /= float64(unit)
		if value < float64(unit) || label == labels[len(labels)-1] {
			return fmt.Sprintf("%.1f %s", value, label)
		}
	}
	return fmt.Sprintf("%d B", bytes)
}
