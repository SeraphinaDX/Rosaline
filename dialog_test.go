// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import (
	"testing"

	tk "modernc.org/tk9.0"
)

func TestAskSaveChangesMapsEveryDialogResult(t *testing.T) {
	original := showMessageBox
	t.Cleanup(func() { showMessageBox = original })

	for response, want := range map[string]SaveDecision{
		"yes":    SaveChanges,
		"no":     DiscardChanges,
		"cancel": CancelChanges,
		"":       CancelChanges,
	} {
		showMessageBox = func(...tk.Opt) string { return response }
		if got := AskSaveChanges("Unsaved", "Save first?"); got != want {
			t.Errorf("response %q = %v, want %v", response, got, want)
		}
	}
}

func TestNormalizeExtension(t *testing.T) {
	tests := map[string]string{
		"png":    ".png",
		".JPG":   ".jpg",
		"*.WebP": ".webp",
		" *.* ":  "*",
		"*":      "*",
		"":       "",
	}
	for input, want := range tests {
		if got := normalizeExtension(input); got != want {
			t.Errorf("normalizeExtension(%q) = %q, want %q", input, got, want)
		}
	}
}
