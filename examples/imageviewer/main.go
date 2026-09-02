// SPDX-License-Identifier: LGPL-3.0-or-later

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/SeraphinaDX/Rosaline"
)

func main() {
	viewer := rosaline.Image(nil).Placeholder("Choose File > Open to view an image.")
	status := "No image loaded"
	currentPath := ""

	imageFilters := []rosaline.FileFilter{
		{Name: "Images", Extensions: []string{"png", "jpg", "jpeg", "gif", "bmp", "tiff", "webp", "avif"}},
		{Name: "All files", Extensions: []string{"*"}},
	}

	openImage := func() {
		path, ok := rosaline.OpenFileDialog(rosaline.FileDialogOptions{
			Title:   "Open Image",
			Filters: imageFilters,
		})
		if !ok {
			return
		}

		picture, err := rosaline.LoadImage(path)
		if err != nil {
			rosaline.Error("Could not open image", err.Error())
			return
		}
		viewer.SetImage(picture)
		currentPath = path
		status = fmt.Sprintf("%s — %d × %d", filepath.Base(path), picture.Width(), picture.Height())
	}

	saveCopy := func() {
		if currentPath == "" {
			rosaline.Message("Save a copy", "Open an image first.")
			return
		}
		path, ok := rosaline.SaveFileDialog(rosaline.FileDialogOptions{
			Title:            "Save Image Copy",
			InitialFile:      filepath.Base(currentPath),
			DefaultExtension: filepath.Ext(currentPath),
			Filters:          imageFilters,
		})
		if !ok {
			return
		}

		data, err := os.ReadFile(currentPath)
		if err == nil {
			err = os.WriteFile(path, data, 0o644)
		}
		if err != nil {
			rosaline.Error("Could not save image", err.Error())
			return
		}
		rosaline.Message("Image saved", "Saved a copy to:\n"+path)
	}

	closeImage := func() {
		if currentPath == "" {
			return
		}
		if !rosaline.Confirm("Close image?", "Remove the current image from the viewer?") {
			return
		}
		viewer.SetImage(nil)
		currentPath = ""
		status = "No image loaded"
	}

	rosaline.RunApp(rosaline.App{
		Title:   "Rosaline Image Viewer",
		Width:   900,
		Height:  650,
		Padding: 12,
		Menu: rosaline.MenuBar(
			rosaline.Menu("File",
				rosaline.MenuItem("Open…", openImage).Shortcut("Ctrl+O"),
				rosaline.MenuItem("Save Copy…", saveCopy).Shortcut("Ctrl+Shift+S"),
				rosaline.MenuItem("Close", closeImage).Shortcut("Ctrl+W"),
				rosaline.MenuSeparator(),
				rosaline.MenuItem("Quit", rosaline.Quit).Shortcut("Ctrl+Q"),
			),
			rosaline.Menu("Help",
				rosaline.MenuItem("About", func() {
					rosaline.Message("About", "Image Viewer built with Rosaline v0.5.0")
				}),
			),
		),
		Content: rosaline.Column(
			rosaline.LabelFunc(func() string { return status }).Color(rosaline.Rose),
			rosaline.Scroll(viewer).Size(820, 520).Expand(),
		).Gap(10).Expand(),
	})
}
