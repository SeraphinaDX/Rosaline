// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import "testing"

func TestLabelPresentationOptions(t *testing.T) {
	label := Label("Rosaline").
		FontSize(28).
		Bold().
		TextAlign(AlignEnd)
	if label.fontSize != 28 || !label.bold || label.alignment != AlignEnd {
		t.Fatalf("label presentation = %#v", label)
	}

	label.FontSize(-4).TextAlign(Alignment(99))
	if label.fontSize != 0 || label.alignment != AlignStart {
		t.Fatalf("normalized label presentation = %#v", label)
	}
}

func TestLabelAnchors(t *testing.T) {
	if labelAnchor(AlignStart) != "w" || labelAnchor(AlignCenter) != "center" || labelAnchor(AlignEnd) != "e" {
		t.Fatal("label anchors did not match start, center, and end")
	}
	if labelAnchor(AlignStretch) != "w" {
		t.Fatal("stretch text alignment should use the start anchor")
	}
	label := Label("Stretch").TextAlign(AlignStretch)
	if label.alignment != AlignStart {
		t.Fatal("stretch text alignment should normalize to start")
	}
}
