// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import (
	"strings"

	tk "modernc.org/tk9.0"
)

// CheckBoxWidget is a labeled checkbox bound to a Go bool.
type CheckBoxWidget struct {
	text     string
	value    *bool
	onChange func(bool)
	focus    bool
	enabled  bool
	variable *tk.VariableOpt
	checkBox *tk.CheckbuttonWidget
	ctx      *mountContext
}

// CheckBox creates a checkbox. It updates value when the user toggles it.
// Pass a pointer with &, as in CheckBox("Updates", &updates).
func CheckBox(text string, value *bool) *CheckBoxWidget {
	if value == nil {
		value = new(bool)
	}
	return &CheckBoxWidget{text: text, value: value, enabled: true}
}

// Checked reports the current bound value.
func (c *CheckBoxWidget) Checked() bool {
	return c != nil && c.value != nil && *c.value
}

// SetChecked changes the value. When mounted, OnChange runs for a change.
func (c *CheckBoxWidget) SetChecked(checked bool) {
	if c == nil || c.value == nil {
		return
	}
	old := *c.value
	*c.value = checked
	if c.variable != nil {
		c.variable.Set(checked)
		if old != checked && c.onChange != nil {
			c.onChange(checked)
		}
		if c.ctx != nil {
			c.ctx.refresh()
		}
	}
}

// SetText replaces the checkbox label immediately.
func (c *CheckBoxWidget) SetText(text string) {
	if c == nil {
		return
	}
	c.text = text
	if c.checkBox != nil {
		c.checkBox.Configure(tk.Txt(text))
	}
}

// Text returns the current checkbox label.
func (c *CheckBoxWidget) Text() string {
	if c == nil {
		return ""
	}
	return c.text
}

// SetEnabled enables or disables the checkbox immediately.
func (c *CheckBoxWidget) SetEnabled(enabled bool) {
	if c == nil {
		return
	}
	c.enabled = enabled
	if c.checkBox != nil {
		c.checkBox.Configure(tk.State(enabledState(enabled)))
	}
}

// Enabled reports whether the checkbox can currently be toggled.
func (c *CheckBoxWidget) Enabled() bool { return c != nil && c.enabled }

// OnChange runs after the user toggles the checkbox.
func (c *CheckBoxWidget) OnChange(handler func(bool)) *CheckBoxWidget {
	c.onChange = handler
	return c
}

// Focus asks Rosaline to give this checkbox focus when the window opens.
// If several widgets request focus, the first one wins.
func (c *CheckBoxWidget) Focus() *CheckBoxWidget {
	c.focus = true
	return c
}

func (c *CheckBoxWidget) mount(ctx *mountContext, parent *tk.Window) mountedWidget {
	variable := tk.Variable(*c.value)
	c.variable = variable
	c.ctx = ctx
	lastValue := *c.value
	var checkBox *tk.CheckbuttonWidget

	syncValue := func() {
		current := boolValue(variable.Get())
		*c.value = current
		if current != lastValue {
			lastValue = current
			if c.onChange != nil {
				c.onChange(current)
			}
		}
	}

	checkBox = parent.Checkbutton(
		tk.Txt(c.text),
		variable,
		tk.Onvalue("true"),
		tk.Offvalue("false"),
		tk.Command(func() {
			syncValue()
			ctx.refresh()
		}),
		tk.Background(ctx.theme.Background.String()),
		tk.Foreground(ctx.theme.Text.String()),
		tk.Activebackground(ctx.theme.Background.String()),
		tk.Activeforeground(ctx.theme.Text.String()),
		tk.Selectcolor(ctx.theme.Surface.String()),
		tk.Highlightcolor(ctx.theme.Primary.String()),
		tk.Highlightbackground(ctx.theme.Background.String()),
		takeFocusOption(true),
		tk.Anchor("w"),
		tk.State(enabledState(c.enabled)),
	)
	c.checkBox = checkBox
	ctx.addCleanup(func() {
		c.variable = nil
		c.checkBox = nil
		c.ctx = nil
	})

	ctx.flushes = append(ctx.flushes, syncValue)
	ctx.refreshes = append(ctx.refreshes, func() {
		if boolValue(variable.Get()) != *c.value {
			variable.Set(*c.value)
			lastValue = *c.value
		}
	})
	ctx.addFocusable(checkBox.Window, c.focus)

	return mountedWidget{window: checkBox.Window}
}

func boolValue(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}
