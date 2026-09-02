// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import (
	"strings"

	tk "modernc.org/tk9.0"
)

// Message displays a simple informational dialog.
func Message(title, text string) {
	tk.MessageBox(tk.Title(title), tk.Msg(text), tk.Type("ok"), tk.Icon("info"), tk.Parent(tk.App))
}

// Error displays an error dialog.
func Error(title, text string) {
	tk.MessageBox(tk.Title(title), tk.Msg(text), tk.Type("ok"), tk.Icon("error"), tk.Parent(tk.App))
}

// Confirm asks a yes-or-no question and reports whether the user chose Yes.
func Confirm(title, text string) bool {
	return tk.MessageBox(
		tk.Title(title),
		tk.Msg(text),
		tk.Type("yesno"),
		tk.Icon("question"),
		tk.Parent(tk.App),
	) == "yes"
}

// FileFilter describes one group of files in an open or save dialog.
type FileFilter struct {
	Name       string
	Extensions []string
}

// FileDialogOptions customizes an open or save dialog. Every field is
// optional; Rosaline supplies beginner-friendly defaults.
type FileDialogOptions struct {
	Title            string
	InitialDirectory string
	InitialFile      string
	DefaultExtension string
	Filters          []FileFilter
}

// OpenFileDialog asks the user to choose one existing file. ok is false when
// the user cancels the dialog.
func OpenFileDialog(options FileDialogOptions) (path string, ok bool) {
	tkOptions := fileDialogOptions(options, "Open File")
	selected := tk.GetOpenFile(tkOptions...)
	if len(selected) == 0 || selected[0] == "" {
		return "", false
	}
	return selected[0], true
}

// SaveFileDialog asks the user where to save a file. ok is false when the user
// cancels. Existing files require confirmation before they are returned.
func SaveFileDialog(options FileDialogOptions) (path string, ok bool) {
	tkOptions := fileDialogOptions(options, "Save File")
	tkOptions = append(tkOptions, tk.Confirmoverwrite(true))
	selected := tk.GetSaveFile(tkOptions...)
	if selected == "" {
		return "", false
	}
	return selected, true
}

func fileDialogOptions(options FileDialogOptions, defaultTitle string) []tk.Opt {
	title := options.Title
	if title == "" {
		title = defaultTitle
	}
	tkOptions := []tk.Opt{tk.Title(title), tk.Parent(tk.App)}
	if options.InitialDirectory != "" {
		tkOptions = append(tkOptions, tk.Initialdir(options.InitialDirectory))
	}
	if options.InitialFile != "" {
		tkOptions = append(tkOptions, tk.Initialfile(options.InitialFile))
	}
	if options.DefaultExtension != "" {
		tkOptions = append(tkOptions, tk.Defaultextension(normalizeExtension(options.DefaultExtension)))
	}
	if len(options.Filters) != 0 {
		filters := make([]tk.FileType, 0, len(options.Filters))
		for _, filter := range options.Filters {
			extensions := make([]string, 0, len(filter.Extensions))
			for _, extension := range filter.Extensions {
				if normalized := normalizeExtension(extension); normalized != "" {
					extensions = append(extensions, normalized)
				}
			}
			filters = append(filters, tk.FileType{
				TypeName:   filter.Name,
				Extensions: extensions,
			})
		}
		tkOptions = append(tkOptions, tk.Filetypes(filters))
	}
	return tkOptions
}

func normalizeExtension(extension string) string {
	extension = strings.TrimSpace(extension)
	if extension == "*" || extension == "*.*" {
		return "*"
	}
	extension = strings.TrimPrefix(extension, "*")
	if extension == "" || extension == "." {
		return ""
	}
	if !strings.HasPrefix(extension, ".") {
		extension = "." + extension
	}
	return strings.ToLower(extension)
}
