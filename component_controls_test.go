// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import "testing"

func TestComponentTextSettersWithoutWindow(t *testing.T) {
	label := LabelFunc(func() string { return "Before" })
	label.SetText("After")
	if label.Text() != "After" {
		t.Fatalf("label text = %q", label.Text())
	}

	button := Button("Save", nil)
	button.SetText("Saved")
	clicked := false
	button.OnClick(func() { clicked = true })
	button.onClick()
	if button.Text() != "Saved" {
		t.Fatalf("button text = %q", button.Text())
	}
	if !clicked {
		t.Fatal("replacement button click handler did not run")
	}

	checked := false
	checkBox := CheckBox("Ready", &checked)
	checkBox.SetText("Complete")
	checkBox.SetChecked(true)
	if checkBox.Text() != "Complete" || !checkBox.Checked() || !checked {
		t.Fatalf("checkbox = text %q, checked %v, value %v", checkBox.Text(), checkBox.Checked(), checked)
	}

	value := "old"
	textBox := TextBox(&value)
	textBox.SetText("new")
	if textBox.Text() != "new" || value != "new" {
		t.Fatalf("textbox = text %q, value %q", textBox.Text(), value)
	}
}

func TestInteractiveComponentsCanBeDisabledBeforeMount(t *testing.T) {
	button := Button("Save", nil)
	checkBox := CheckBox("Ready", nil)
	textBox := TextBox(nil)

	button.SetEnabled(false)
	checkBox.SetEnabled(false)
	textBox.SetEnabled(false)
	if button.Enabled() || checkBox.Enabled() || textBox.Enabled() {
		t.Fatal("disabled controls reported themselves enabled")
	}

	button.SetEnabled(true)
	checkBox.SetEnabled(true)
	textBox.SetEnabled(true)
	if !button.Enabled() || !checkBox.Enabled() || !textBox.Enabled() {
		t.Fatal("enabled controls reported themselves disabled")
	}
}
