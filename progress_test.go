// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import (
	"math"
	"testing"
)

func TestProgressBarClampsValue(t *testing.T) {
	value := 125.0
	progress := ProgressBar(&value)
	if value != 100 || progress.Value() != 100 {
		t.Fatalf("value = %v, want 100", value)
	}

	progress.SetValue(-10)
	if value != 0 {
		t.Fatalf("value = %v, want 0", value)
	}

	progress.SetValue(math.NaN())
	if value != 0 {
		t.Fatalf("NaN value = %v, want 0", value)
	}
}

func TestProgressBarMaximum(t *testing.T) {
	value := 75.0
	progress := ProgressBar(&value).Maximum(50)
	if progress.Max() != 50 || value != 50 {
		t.Fatalf("max/value = %v/%v, want 50/50", progress.Max(), value)
	}

	progress.Maximum(-1)
	if progress.Max() != 100 {
		t.Fatalf("invalid maximum = %v, want 100", progress.Max())
	}
}

func TestProgressBarBusyControls(t *testing.T) {
	progress := ProgressBar(nil).Busy()
	if !progress.IsBusy() || !progress.Running() {
		t.Fatal("Busy did not start indeterminate progress")
	}

	progress.Stop()
	if progress.Running() {
		t.Fatal("Stop left progress running")
	}
	progress.Start()
	if !progress.Running() {
		t.Fatal("Start did not resume progress")
	}
	progress.Determinate()
	if progress.IsBusy() || progress.Running() {
		t.Fatal("Determinate left progress in busy mode")
	}
	progress.Start()
	if progress.Running() {
		t.Fatal("Start should not animate determinate progress")
	}
}

func TestProgressBarOptions(t *testing.T) {
	progress := ProgressBar(nil).Length(320).Vertical()
	if progress.length != 320 || !progress.vertical {
		t.Fatalf("options = length %d, vertical %v", progress.length, progress.vertical)
	}
	progress.Length(0).Horizontal()
	if progress.length != 320 || progress.vertical {
		t.Fatalf("updated options = length %d, vertical %v", progress.length, progress.vertical)
	}
}

func TestNilProgressBarIsSafe(t *testing.T) {
	var progress *ProgressBarWidget
	progress.Maximum(10).Length(10).Vertical().Horizontal().Busy().Stop().Start().Determinate()
	progress.SetValue(10)
	if progress.Value() != 0 || progress.Max() != 0 || progress.IsBusy() || progress.Running() {
		t.Fatal("nil progress bar accessors should return zero values")
	}
}
