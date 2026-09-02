// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import tk "modernc.org/tk9.0"

// ComboBoxWidget displays a compact drop-down selection bound to a Go string.
type ComboBoxWidget struct {
	value    *string
	options  []string
	columns  int
	focus    bool
	onChange func(string)
	combo    *tk.TComboboxWidget
	ctx      *mountContext
}

// ComboBox creates a read-only drop-down. It updates value when the user picks
// an option. An unavailable value safely selects the first option.
func ComboBox(value *string, options ...string) *ComboBoxWidget {
	if value == nil {
		value = new(string)
	}
	combo := &ComboBoxWidget{
		value:   value,
		options: uniqueStrings(options),
		columns: 28,
	}
	combo.normalizeValue()
	return combo
}

// Width sets the preferred width in text columns.
func (c *ComboBoxWidget) Width(columns int) *ComboBoxWidget {
	if c != nil && columns > 0 {
		c.columns = columns
		if c.combo != nil {
			c.combo.Configure(tk.Width(columns))
		}
	}
	return c
}

// Focus asks Rosaline to focus this combo box when the window opens.
func (c *ComboBoxWidget) Focus() *ComboBoxWidget {
	if c != nil {
		c.focus = true
	}
	return c
}

// OnChange runs after the selected value changes through the UI or Select.
func (c *ComboBoxWidget) OnChange(handler func(string)) *ComboBoxWidget {
	if c != nil {
		c.onChange = handler
	}
	return c
}

// Options returns a copy of the current drop-down options.
func (c *ComboBoxWidget) Options() []string {
	if c == nil {
		return nil
	}
	return append([]string(nil), c.options...)
}

// Selected returns the currently bound option.
func (c *ComboBoxWidget) Selected() string {
	if c == nil || c.value == nil {
		return ""
	}
	return *c.value
}

// Select changes the selection. Values not present in the options are ignored.
func (c *ComboBoxWidget) Select(value string) {
	if c == nil || !stringContains(c.options, value) || *c.value == value {
		return
	}
	*c.value = value
	if c.combo != nil {
		c.applySelection()
		if c.onChange != nil {
			c.onChange(value)
		}
		if c.ctx != nil {
			c.ctx.refresh()
		}
	}
}

// SetOptions replaces every option. Duplicate strings are ignored. If the old
// value is unavailable, the first replacement is selected.
func (c *ComboBoxWidget) SetOptions(options ...string) {
	if c == nil {
		return
	}
	oldValue := *c.value
	c.options = uniqueStrings(options)
	c.normalizeValue()
	if c.combo != nil {
		c.combo.Configure(tk.Values(c.options))
		c.applySelection()
		if oldValue != *c.value && c.onChange != nil {
			c.onChange(*c.value)
		}
		if c.ctx != nil {
			c.ctx.refresh()
		}
	}
}

func uniqueStrings(values []string) []string {
	result := make([]string, 0, len(values))
	seen := make(map[string]bool, len(values))
	for _, value := range values {
		if seen[value] {
			continue
		}
		seen[value] = true
		result = append(result, value)
	}
	return result
}

func stringContains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func (c *ComboBoxWidget) normalizeValue() {
	if len(c.options) == 0 {
		*c.value = ""
		return
	}
	if !stringContains(c.options, *c.value) {
		*c.value = c.options[0]
	}
}

func (c *ComboBoxWidget) applySelection() {
	if c.combo == nil {
		return
	}
	for index, option := range c.options {
		if option == *c.value {
			c.combo.Current(index)
			return
		}
	}
	c.combo.Configure(tk.Textvariable(""))
}

func (c *ComboBoxWidget) mount(ctx *mountContext, parent *tk.Window) mountedWidget {
	c.normalizeValue()
	c.ctx = ctx
	c.combo = parent.TCombobox(
		tk.Values(c.options),
		tk.Textvariable(*c.value),
		tk.State("readonly"),
		tk.Width(c.columns),
		takeFocusOption(true),
	)
	c.applySelection()

	syncValue := func() {
		current := c.combo.Textvariable()
		if !stringContains(c.options, current) {
			return
		}
		oldValue := *c.value
		*c.value = current
		if oldValue != current && c.onChange != nil {
			c.onChange(current)
		}
	}
	tk.Bind(c.combo, "<<ComboboxSelected>>", tk.Command(func() {
		syncValue()
		ctx.refresh()
	}))
	ctx.flushes = append(ctx.flushes, syncValue)
	ctx.refreshes = append(ctx.refreshes, func() {
		c.normalizeValue()
		if c.combo != nil && c.combo.Textvariable() != *c.value {
			c.applySelection()
		}
	})
	ctx.addFocusable(c.combo.Window, c.focus)
	ctx.addCleanup(func() {
		c.combo = nil
		c.ctx = nil
	})
	return mountedWidget{window: c.combo.Window}
}
