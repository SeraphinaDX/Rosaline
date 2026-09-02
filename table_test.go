// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import (
	"slices"
	"testing"
)

func TestTableUsesFriendlyColumns(t *testing.T) {
	table := Table("Name", "", "Size")
	if got, want := table.Columns(), []string{"Name", "Column 2", "Size"}; !slices.Equal(got, want) {
		t.Fatalf("Columns() = %#v, want %#v", got, want)
	}
	if got := Table().Columns(); !slices.Equal(got, []string{"Value"}) {
		t.Fatalf("Table() columns = %#v, want Value", got)
	}
}

func TestTableRowsAreNormalizedAndCopied(t *testing.T) {
	first := []string{"README.md", "Markdown"}
	table := Table("Name", "Type", "Size").SetRows(
		first,
		[]string{"photo.png", "Image", "2 MB", "ignored"},
	)
	first[0] = "changed"

	rows := table.Rows()
	want := [][]string{
		{"README.md", "Markdown", ""},
		{"photo.png", "Image", "2 MB"},
	}
	if !equalTableRows(rows, want) {
		t.Fatalf("Rows() = %#v, want %#v", rows, want)
	}
	rows[0][0] = "changed again"
	if table.Rows()[0][0] != "README.md" {
		t.Fatal("modifying Rows result changed the table")
	}
}

func TestTableSelectsFirstRowAndReturnsCopy(t *testing.T) {
	table := Table("Name").SetRows([]string{"Rose"}, []string{"Violet"})
	row, values, ok := table.Selected()
	if !ok || row != 0 || !slices.Equal(values, []string{"Rose"}) {
		t.Fatalf("Selected() = %d, %#v, %v", row, values, ok)
	}
	values[0] = "changed"
	if _, values, _ := table.Selected(); values[0] != "Rose" {
		t.Fatal("modifying Selected values changed the table")
	}

	table.Select(1)
	row, values, ok = table.Selected()
	if !ok || row != 1 || values[0] != "Violet" {
		t.Fatalf("Selected() after Select = %d, %#v, %v", row, values, ok)
	}
}

func TestTableSetRowsKeepsSelectionInRange(t *testing.T) {
	table := Table("Value").SetRows([]string{"A"}, []string{"B"}, []string{"C"})
	table.Select(2)
	table.SetRows([]string{"One"}, []string{"Two"})
	row, values, ok := table.Selected()
	if !ok || row != 1 || values[0] != "Two" {
		t.Fatalf("Selected() = %d, %#v, %v; want row 1", row, values, ok)
	}

	table.SetRows()
	if _, _, ok := table.Selected(); ok {
		t.Fatal("empty table should have no selection")
	}
	table.SetRows([]string{"Again"})
	if row, _, ok := table.Selected(); !ok || row != 0 {
		t.Fatalf("refilled table selection = %d, %v; want row 0", row, ok)
	}
}

func TestInvalidTableSelectionClearsSelection(t *testing.T) {
	table := Table("Value").SetRows([]string{"A"})
	table.Select(9)
	if _, _, ok := table.Selected(); ok {
		t.Fatal("invalid row should clear selection")
	}
}

func TestTableBuilderOptions(t *testing.T) {
	table := Table("Name", "Size").
		ColumnWidth(0, 320).
		ColumnWidth(9, 500).
		Height(18).
		Expand().
		OnSelect(func(int, []string) {}).
		OnActivate(func(int, []string) {})

	if table.widths[0] != 320 || table.widths[1] != defaultTableColumnWidth {
		t.Fatalf("column widths = %#v", table.widths)
	}
	if table.height != 18 || !table.expand || table.onSelect == nil || table.onActivate == nil {
		t.Fatalf("table options were not retained: %#v", table)
	}
}

func equalTableRows(left, right [][]string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if !slices.Equal(left[index], right[index]) {
			return false
		}
	}
	return true
}
