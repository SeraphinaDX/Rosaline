// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import (
	"strings"

	tk "modernc.org/tk9.0"
)

// RadioChoice is one labeled value accepted by RadioGroup. Create choices
// with Choice rather than filling this type manually.
type RadioChoice struct {
	label string
	value string
}

// Choice creates one radio-group choice. Empty labels use the value, or the
// friendly label "Choice" when both strings are empty.
func Choice(label, value string) RadioChoice {
	if strings.TrimSpace(label) == "" {
		label = value
		if strings.TrimSpace(label) == "" {
			label = "Choice"
		}
	}
	return RadioChoice{label: label, value: value}
}

// Label returns the text displayed beside this choice.
func (c RadioChoice) Label() string { return c.label }

// Value returns the value stored in the RadioGroup's bound Go string.
func (c RadioChoice) Value() string { return c.value }

// RadioGroupWidget displays mutually exclusive choices bound to a Go string.
type RadioGroupWidget struct {
	value      *string
	choices    []RadioChoice
	horizontal bool
	focus      bool
	onChange   func(string)
	frame      *tk.FrameWidget
	variable   *tk.VariableOpt
	buttons    []*tk.TRadiobuttonWidget
	ctx        *mountContext
	focusSlot  int
}

// RadioGroup creates a vertical group of choices. It updates value whenever
// the user selects one. An unavailable value safely selects the first choice.
func RadioGroup(value *string, choices ...RadioChoice) *RadioGroupWidget {
	if value == nil {
		value = new(string)
	}
	group := &RadioGroupWidget{
		value:     value,
		choices:   normalizeRadioChoices(choices),
		focusSlot: -1,
	}
	group.normalizeValue()
	return group
}

// Horizontal arranges the choices from left to right.
func (r *RadioGroupWidget) Horizontal() *RadioGroupWidget {
	if r != nil && !r.horizontal {
		r.horizontal = true
		r.rebuild()
	}
	return r
}

// Vertical arranges the choices from top to bottom. This is the default.
func (r *RadioGroupWidget) Vertical() *RadioGroupWidget {
	if r != nil && r.horizontal {
		r.horizontal = false
		r.rebuild()
	}
	return r
}

// Focus asks Rosaline to focus the selected choice when the window opens.
func (r *RadioGroupWidget) Focus() *RadioGroupWidget {
	if r != nil {
		r.focus = true
	}
	return r
}

// OnChange runs after the selected value changes through the UI or Select.
func (r *RadioGroupWidget) OnChange(handler func(string)) *RadioGroupWidget {
	if r != nil {
		r.onChange = handler
	}
	return r
}

// Choices returns a copy of the configured choices.
func (r *RadioGroupWidget) Choices() []RadioChoice {
	if r == nil {
		return nil
	}
	return append([]RadioChoice(nil), r.choices...)
}

// Selected returns the currently bound choice value.
func (r *RadioGroupWidget) Selected() string {
	if r == nil || r.value == nil {
		return ""
	}
	return *r.value
}

// Select changes the selection. Values not present in the group are ignored.
func (r *RadioGroupWidget) Select(value string) {
	if r == nil || !radioChoiceContains(r.choices, value) || *r.value == value {
		return
	}
	*r.value = value
	if r.variable != nil {
		r.variable.Set(value)
		r.updateFocusSlot()
		if r.onChange != nil {
			r.onChange(value)
		}
		if r.ctx != nil {
			r.ctx.refresh()
		}
	}
}

// SetChoices replaces every choice. Duplicate values are ignored. If the old
// value is unavailable, the first replacement is selected.
func (r *RadioGroupWidget) SetChoices(choices ...RadioChoice) {
	if r == nil {
		return
	}
	oldValue := *r.value
	r.choices = normalizeRadioChoices(choices)
	r.normalizeValue()
	if r.frame != nil {
		if r.variable != nil {
			r.variable.Set(*r.value)
		}
		r.rebuild()
		if oldValue != *r.value && r.onChange != nil {
			r.onChange(*r.value)
		}
		if r.ctx != nil {
			r.ctx.refresh()
		}
	}
}

func normalizeRadioChoices(choices []RadioChoice) []RadioChoice {
	result := make([]RadioChoice, 0, len(choices))
	seen := make(map[string]bool, len(choices))
	for _, choice := range choices {
		choice = Choice(choice.label, choice.value)
		if seen[choice.value] {
			continue
		}
		seen[choice.value] = true
		result = append(result, choice)
	}
	return result
}

func radioChoiceContains(choices []RadioChoice, value string) bool {
	for _, choice := range choices {
		if choice.value == value {
			return true
		}
	}
	return false
}

func (r *RadioGroupWidget) normalizeValue() {
	if len(r.choices) == 0 {
		*r.value = ""
		return
	}
	if !radioChoiceContains(r.choices, *r.value) {
		*r.value = r.choices[0].value
	}
}

func (r *RadioGroupWidget) updateFocusSlot() {
	if r.ctx == nil || r.focusSlot < 0 {
		return
	}
	var selected *tk.Window
	for index, choice := range r.choices {
		if choice.value == *r.value && index < len(r.buttons) {
			selected = r.buttons[index].Window
			break
		}
	}
	r.ctx.updateFocusable(r.focusSlot, selected)
}

func (r *RadioGroupWidget) rebuild() {
	if r == nil || r.frame == nil {
		return
	}
	for _, button := range r.buttons {
		tk.Destroy(button.Window)
	}
	r.buttons = nil

	for _, choice := range r.choices {
		choice := choice
		button := r.frame.TRadiobutton(
			tk.Txt(choice.label),
			r.variable,
			tk.Value(choice.value),
			takeFocusOption(true),
			tk.Command(func() {
				oldValue := *r.value
				*r.value = r.variable.Get()
				r.updateFocusSlot()
				if oldValue != *r.value && r.onChange != nil {
					r.onChange(*r.value)
				}
				if r.ctx != nil {
					r.ctx.refresh()
				}
			}),
		)
		r.buttons = append(r.buttons, button)
		if r.focusSlot >= 0 {
			r.ctx.bindFocusTraversal(button.Window, r.focusSlot)
		}
		options := []tk.Opt{button, tk.Anchor("w"), tk.Padx(3), tk.Pady(3)}
		if r.horizontal {
			options = append(options, tk.Side("left"))
		} else {
			options = append(options, tk.Side("top"), tk.Fill("x"))
		}
		tk.Pack(options...)
	}
	r.updateFocusSlot()
}

func (r *RadioGroupWidget) mount(ctx *mountContext, parent *tk.Window) mountedWidget {
	r.normalizeValue()
	r.ctx = ctx
	r.frame = parent.Frame(
		tk.Background(ctx.theme.Background.String()),
		tk.Borderwidth(0),
	)
	r.variable = tk.Variable(*r.value)
	r.focusSlot = ctx.addFocusable(nil, false)
	r.rebuild()
	if r.focus && ctx.initialFocus == nil && r.focusSlot < len(ctx.focusables) {
		ctx.initialFocus = ctx.focusables[r.focusSlot].window
	}
	ctx.flushes = append(ctx.flushes, func() {
		if r.variable != nil {
			value := r.variable.Get()
			if radioChoiceContains(r.choices, value) {
				*r.value = value
			}
		}
	})
	ctx.refreshes = append(ctx.refreshes, func() {
		r.normalizeValue()
		if r.variable != nil && r.variable.Get() != *r.value {
			r.variable.Set(*r.value)
			r.updateFocusSlot()
		}
	})
	ctx.addCleanup(func() {
		r.frame = nil
		r.variable = nil
		r.buttons = nil
		r.ctx = nil
		r.focusSlot = -1
	})
	return mountedWidget{window: r.frame.Window}
}
