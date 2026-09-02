// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import (
	"math"
	"strconv"

	tk "modernc.org/tk9.0"
)

const defaultSliderLength = 260

// SliderWidget displays a numeric slider bound to a Go float64.
type SliderWidget struct {
	value    *float64
	minimum  float64
	maximum  float64
	step     float64
	length   int
	vertical bool
	focus    bool
	onChange func(float64)
	scale    *tk.TScaleWidget
	variable *tk.VariableOpt
	ctx      *mountContext
}

// Slider creates a horizontal numeric slider. Reversed bounds are swapped,
// equal bounds receive a safe one-unit range, and value is clamped.
func Slider(value *float64, minimum, maximum float64) *SliderWidget {
	if value == nil {
		value = new(float64)
	}
	minimum, maximum = normalizeNumericRange(minimum, maximum)
	slider := &SliderWidget{
		value:   value,
		minimum: minimum,
		maximum: maximum,
		length:  defaultSliderLength,
	}
	*slider.value = slider.normalizeValue(*slider.value)
	return slider
}

// Step rounds values to a positive interval measured from the minimum. Zero
// or invalid steps leave the slider continuous.
func (s *SliderWidget) Step(step float64) *SliderWidget {
	if s == nil {
		return s
	}
	if !finiteNumber(step) || step <= 0 {
		step = 0
	}
	s.step = step
	s.SetValue(*s.value)
	return s
}

// SetRange changes the numeric bounds and clamps the current value.
func (s *SliderWidget) SetRange(minimum, maximum float64) *SliderWidget {
	if s == nil {
		return s
	}
	s.minimum, s.maximum = normalizeNumericRange(minimum, maximum)
	if s.scale != nil {
		s.scale.Configure(tk.From(s.minimum), tk.To(s.maximum))
	}
	s.SetValue(*s.value)
	return s
}

// Bounds returns the normalized minimum and maximum.
func (s *SliderWidget) Bounds() (minimum, maximum float64) {
	if s == nil {
		return 0, 0
	}
	return s.minimum, s.maximum
}

// Length sets the preferred slider length in pixels.
func (s *SliderWidget) Length(pixels int) *SliderWidget {
	if s != nil && pixels > 0 {
		s.length = pixels
		if s.scale != nil {
			s.scale.Configure(tk.Length(pixels))
		}
	}
	return s
}

// Vertical arranges the slider from bottom to top.
func (s *SliderWidget) Vertical() *SliderWidget {
	if s != nil && !s.vertical {
		s.vertical = true
		if s.scale != nil {
			s.scale.Configure(tk.Orient("vertical"))
		}
	}
	return s
}

// Horizontal arranges the slider from left to right. This is the default.
func (s *SliderWidget) Horizontal() *SliderWidget {
	if s != nil && s.vertical {
		s.vertical = false
		if s.scale != nil {
			s.scale.Configure(tk.Orient("horizontal"))
		}
	}
	return s
}

// Focus asks Rosaline to focus this slider when the window opens.
func (s *SliderWidget) Focus() *SliderWidget {
	if s != nil {
		s.focus = true
	}
	return s
}

// OnChange runs after the value changes through the UI or SetValue.
func (s *SliderWidget) OnChange(handler func(float64)) *SliderWidget {
	if s != nil {
		s.onChange = handler
	}
	return s
}

// SetValue clamps and optionally rounds a new value.
func (s *SliderWidget) SetValue(value float64) {
	if s == nil {
		return
	}
	value = s.normalizeValue(value)
	oldValue := *s.value
	*s.value = value
	if s.variable != nil {
		s.variable.Set(formatNumber(value))
		if oldValue != value && s.onChange != nil {
			s.onChange(value)
		}
		if s.ctx != nil {
			s.ctx.refresh()
		}
	}
}

// Value returns the currently bound numeric value.
func (s *SliderWidget) Value() float64 {
	if s == nil || s.value == nil {
		return 0
	}
	return *s.value
}

func finiteNumber(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0)
}

func normalizeNumericRange(minimum, maximum float64) (float64, float64) {
	if !finiteNumber(minimum) {
		minimum = 0
	}
	if !finiteNumber(maximum) {
		maximum = 100
	}
	if maximum < minimum {
		minimum, maximum = maximum, minimum
	}
	if maximum == minimum {
		maximum = minimum + 1
	}
	return minimum, maximum
}

func formatNumber(value float64) string {
	return strconv.FormatFloat(value, 'g', -1, 64)
}

func parseNumber(value string, fallback float64) float64 {
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil || !finiteNumber(parsed) {
		return fallback
	}
	return parsed
}

func (s *SliderWidget) normalizeValue(value float64) float64 {
	if !finiteNumber(value) {
		value = s.minimum
	}
	value = max(s.minimum, min(s.maximum, value))
	if s.step > 0 {
		value = s.minimum + math.Round((value-s.minimum)/s.step)*s.step
		value = max(s.minimum, min(s.maximum, value))
	}
	return value
}

func (s *SliderWidget) mount(ctx *mountContext, parent *tk.Window) mountedWidget {
	*s.value = s.normalizeValue(*s.value)
	s.ctx = ctx
	s.variable = tk.Variable(formatNumber(*s.value))
	orientation := "horizontal"
	if s.vertical {
		orientation = "vertical"
	}
	s.scale = parent.TScale(
		tk.From(s.minimum),
		tk.To(s.maximum),
		tk.Orient(orientation),
		tk.Length(s.length),
		s.variable,
		takeFocusOption(true),
	)

	syncValue := func() {
		current := s.normalizeValue(parseNumber(s.scale.Get(), *s.value))
		oldValue := *s.value
		*s.value = current
		if parseNumber(s.variable.Get(), current) != current {
			s.variable.Set(formatNumber(current))
		}
		if oldValue != current && s.onChange != nil {
			s.onChange(current)
		}
	}
	s.scale.Configure(tk.Command(func() {
		syncValue()
		ctx.refresh()
	}))
	ctx.flushes = append(ctx.flushes, syncValue)
	ctx.refreshes = append(ctx.refreshes, func() {
		*s.value = s.normalizeValue(*s.value)
		if s.variable != nil && parseNumber(s.variable.Get(), *s.value) != *s.value {
			s.variable.Set(formatNumber(*s.value))
		}
	})
	ctx.addFocusable(s.scale.Window, s.focus)
	ctx.addCleanup(func() {
		s.scale = nil
		s.variable = nil
		s.ctx = nil
	})
	return mountedWidget{window: s.scale.Window}
}
