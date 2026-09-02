// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import (
	"math"
	"testing"
)

func TestSliderNormalizesRangeAndValue(t *testing.T) {
	value := 150.0
	slider := Slider(&value, 100, 0)
	minimum, maximum := slider.Bounds()
	if minimum != 0 || maximum != 100 || value != 100 {
		t.Fatalf("slider = range %v..%v value %v", minimum, maximum, value)
	}
	slider.SetRange(5, 5)
	minimum, maximum = slider.Bounds()
	if minimum != 5 || maximum != 6 || value != 6 {
		t.Fatalf("equal range = %v..%v value %v", minimum, maximum, value)
	}
}

func TestSliderStepAndSetValue(t *testing.T) {
	value := 0.0
	slider := Slider(&value, 0, 10).Step(0.5)
	slider.SetValue(4.26)
	if value != 4.5 || slider.Value() != 4.5 {
		t.Fatalf("stepped value = %v, want 4.5", value)
	}
	slider.SetValue(math.NaN())
	if value != 0 {
		t.Fatalf("NaN value = %v, want minimum", value)
	}
}

func TestSliderSafeDefaultsAndBuilderOptions(t *testing.T) {
	slider := Slider(nil, math.NaN(), math.Inf(1)).
		Length(340).Vertical().Horizontal().Vertical().Focus().OnChange(func(float64) {})
	minimum, maximum := slider.Bounds()
	if slider.value == nil || minimum != 0 || maximum != 100 || slider.length != 340 || !slider.vertical || !slider.focus || slider.onChange == nil {
		t.Fatalf("slider options were not preserved: %#v", slider)
	}
	slider.Step(-1)
	if slider.step != 0 {
		t.Fatalf("invalid step = %v, want continuous", slider.step)
	}
}
