// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import (
	"time"

	ttk "modernc.org/tk9.0"
)

const defaultProgressBarLength = 260

// ProgressBarWidget displays determinate progress or a busy animation.
type ProgressBarWidget struct {
	value    *float64
	maximum  float64
	length   int
	vertical bool
	busy     bool
	running  bool
	bar      *ttk.TProgressbarWidget
	variable *ttk.VariableOpt
	ctx      *mountContext
}

// ProgressBar creates a horizontal determinate progress bar bound to a Go
// float64. Values are clamped between zero and 100 by default.
func ProgressBar(value *float64) *ProgressBarWidget {
	if value == nil {
		value = new(float64)
	}
	progress := &ProgressBarWidget{
		value:   value,
		maximum: 100,
		length:  defaultProgressBarLength,
	}
	*progress.value = progress.normalizeValue(*progress.value)
	return progress
}

// Maximum changes the upper bound. Invalid or non-positive values use 100.
func (p *ProgressBarWidget) Maximum(maximum float64) *ProgressBarWidget {
	if p == nil {
		return p
	}
	if !finiteNumber(maximum) || maximum <= 0 {
		maximum = 100
	}
	p.maximum = maximum
	p.SetValue(*p.value)
	if p.bar != nil {
		p.bar.Configure(ttk.Maximum(p.maximum))
	}
	return p
}

// Max returns the current upper bound.
func (p *ProgressBarWidget) Max() float64 {
	if p == nil {
		return 0
	}
	return p.maximum
}

// Length sets the preferred progress bar length in pixels.
func (p *ProgressBarWidget) Length(pixels int) *ProgressBarWidget {
	if p != nil && pixels > 0 {
		p.length = pixels
		if p.bar != nil {
			p.bar.Configure(ttk.Length(pixels))
		}
	}
	return p
}

// Vertical arranges progress from bottom to top.
func (p *ProgressBarWidget) Vertical() *ProgressBarWidget {
	if p != nil && !p.vertical {
		p.vertical = true
		if p.bar != nil {
			p.bar.Configure(ttk.Orient("vertical"))
		}
	}
	return p
}

// Horizontal arranges progress from left to right. This is the default.
func (p *ProgressBarWidget) Horizontal() *ProgressBarWidget {
	if p != nil && p.vertical {
		p.vertical = false
		if p.bar != nil {
			p.bar.Configure(ttk.Orient("horizontal"))
		}
	}
	return p
}

// Busy switches to an indeterminate animation and starts it. Use this when
// work is happening but its completion percentage is unknown.
func (p *ProgressBarWidget) Busy() *ProgressBarWidget {
	if p == nil {
		return p
	}
	wasRunning := p.busy && p.running
	p.busy = true
	p.running = true
	if p.bar != nil {
		p.bar.Configure(ttk.Mode("indeterminate"))
		if !wasRunning {
			p.bar.Start(50 * time.Millisecond)
		}
	}
	return p
}

// Determinate returns to percentage-style progress and stops any busy
// animation.
func (p *ProgressBarWidget) Determinate() *ProgressBarWidget {
	if p == nil {
		return p
	}
	if p.bar != nil {
		p.bar.Stop()
		p.bar.Configure(ttk.Mode("determinate"))
	}
	p.busy = false
	p.running = false
	p.syncNativeValue()
	return p
}

// Start resumes a busy progress bar. It has no effect in determinate mode.
func (p *ProgressBarWidget) Start() *ProgressBarWidget {
	if p != nil && p.busy && !p.running {
		p.running = true
		if p.bar != nil {
			p.bar.Start(50 * time.Millisecond)
		}
	}
	return p
}

// Stop pauses a busy progress bar. It keeps the bar in busy mode so Start can
// resume it later.
func (p *ProgressBarWidget) Stop() *ProgressBarWidget {
	if p != nil && p.busy && p.running {
		p.running = false
		if p.bar != nil {
			p.bar.Stop()
		}
	}
	return p
}

// IsBusy reports whether the bar is in indeterminate mode.
func (p *ProgressBarWidget) IsBusy() bool {
	return p != nil && p.busy
}

// Running reports whether the busy animation is currently running.
func (p *ProgressBarWidget) Running() bool {
	return p != nil && p.busy && p.running
}

// SetValue changes determinate progress and clamps it to the current maximum.
func (p *ProgressBarWidget) SetValue(value float64) {
	if p == nil {
		return
	}
	value = p.normalizeValue(value)
	*p.value = value
	p.syncNativeValue()
	if p.ctx != nil {
		p.ctx.refresh()
	}
}

// Value returns the currently bound determinate value.
func (p *ProgressBarWidget) Value() float64 {
	if p == nil || p.value == nil {
		return 0
	}
	return *p.value
}

func (p *ProgressBarWidget) normalizeValue(value float64) float64 {
	if !finiteNumber(value) {
		return 0
	}
	return max(0, min(p.maximum, value))
}

func (p *ProgressBarWidget) syncNativeValue() {
	if p == nil || p.variable == nil || p.value == nil {
		return
	}
	*p.value = p.normalizeValue(*p.value)
	if parseNumber(p.variable.Get(), *p.value) != *p.value {
		p.variable.Set(formatNumber(*p.value))
	}
}

func (p *ProgressBarWidget) mount(ctx *mountContext, parent *ttk.Window) mountedWidget {
	*p.value = p.normalizeValue(*p.value)
	p.ctx = ctx
	p.variable = ttk.Variable(formatNumber(*p.value))
	orientation := "horizontal"
	if p.vertical {
		orientation = "vertical"
	}
	mode := "determinate"
	if p.busy {
		mode = "indeterminate"
	}
	p.bar = parent.TProgressbar(
		p.variable,
		ttk.Maximum(p.maximum),
		ttk.Length(p.length),
		ttk.Orient(orientation),
		ttk.Mode(mode),
	)
	if p.busy && p.running {
		p.bar.Start(50 * time.Millisecond)
	}
	ctx.refreshes = append(ctx.refreshes, p.syncNativeValue)
	ctx.addCleanup(func() {
		if p.bar != nil && p.busy && p.running {
			p.bar.Stop()
		}
		p.bar = nil
		p.variable = nil
		p.ctx = nil
	})
	return mountedWidget{window: p.bar.Window}
}
