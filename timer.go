// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import (
	"time"

	tk "modernc.org/tk9.0"
)

const (
	minimumTimerInterval = time.Millisecond
	defaultAnimationFPS  = 60
	maximumAnimationFPS  = 1000
)

var (
	scheduleTimer = tk.TclAfter
	cancelTimer   = tk.TclAfterCancel
)

// Timer runs a callback later or at a regular interval while its App is open.
// Create timers with Every, After, or Animate, then include them in App.Timers.
type Timer struct {
	interval   time.Duration
	callback   func()
	repeat     bool
	running    bool
	ctx        *mountContext
	afterID    string
	generation uint64
}

// Every creates a running timer that calls callback repeatedly. It begins when
// its App starts. Durations shorter than one millisecond use one millisecond.
func Every(interval time.Duration, callback func()) *Timer {
	return newTimer(interval, callback, true)
}

// After creates a running one-shot timer. It calls callback once after delay,
// then stops. It begins when its App starts.
func After(delay time.Duration, callback func()) *Timer {
	return newTimer(delay, callback, false)
}

// Animate creates a repeating timer measured in frames per second. Use it to
// update drawing state, then call CanvasWidget.Redraw from the frame callback.
// Invalid frame rates use 60 FPS; rates above 1000 FPS are limited to 1000.
func Animate(framesPerSecond int, frame func()) *Timer {
	if framesPerSecond <= 0 {
		framesPerSecond = defaultAnimationFPS
	}
	if framesPerSecond > maximumAnimationFPS {
		framesPerSecond = maximumAnimationFPS
	}
	return Every(time.Second/time.Duration(framesPerSecond), frame)
}

func newTimer(interval time.Duration, callback func(), repeat bool) *Timer {
	if interval < minimumTimerInterval {
		interval = minimumTimerInterval
	}
	if callback == nil {
		callback = func() {}
	}
	return &Timer{
		interval: interval,
		callback: callback,
		repeat:   repeat,
		running:  true,
	}
}

// Start starts a stopped timer. Calling Start on a running timer has no effect.
// Call timer methods from Rosaline callbacks, not background goroutines.
func (t *Timer) Start() {
	if t == nil || t.running {
		return
	}
	t.running = true
	if t.ctx != nil {
		t.schedule()
	}
}

// Stop pauses a timer. A stopped repeating timer can be continued with Start.
func (t *Timer) Stop() {
	if t == nil {
		return
	}
	t.running = false
	t.cancelScheduled()
}

// Restart resets the wait and starts the timer again from the beginning.
func (t *Timer) Restart() {
	if t == nil {
		return
	}
	t.cancelScheduled()
	t.running = true
	if t.ctx != nil {
		t.schedule()
	}
}

// Running reports whether the timer is started. Before RunApp, true means the
// timer is ready to begin as soon as its App opens.
func (t *Timer) Running() bool {
	return t != nil && t.running
}

func mountTimers(ctx *mountContext, timers []*Timer) []*Timer {
	mounted := make([]*Timer, 0, len(timers))
	seen := make(map[*Timer]bool, len(timers))
	for _, timer := range timers {
		if timer == nil || seen[timer] {
			continue
		}
		seen[timer] = true
		timer.mount(ctx)
		mounted = append(mounted, timer)
	}
	return mounted
}

func (t *Timer) mount(ctx *mountContext) {
	t.cancelScheduled()
	t.ctx = ctx
}

func (t *Timer) begin() {
	if t.running {
		t.schedule()
	}
}

func (t *Timer) unmount(ctx *mountContext) {
	if t.ctx != ctx {
		return
	}
	t.generation++
	t.afterID = ""
	t.ctx = nil
}

func (t *Timer) cancelScheduled() {
	t.generation++
	if t.afterID != "" && t.ctx != nil && !t.ctx.closed {
		cancelTimer(t.afterID)
	}
	t.afterID = ""
}

func (t *Timer) schedule() {
	if !t.running || t.ctx == nil || t.ctx.closed {
		return
	}

	t.generation++
	generation := t.generation
	t.afterID = scheduleTimer(t.interval, func() {
		if generation != t.generation || !t.running || t.ctx == nil || t.ctx.closed {
			return
		}
		t.afterID = ""
		if !t.repeat {
			t.running = false
		}

		t.callback()
		if t.ctx != nil {
			t.ctx.refresh()
		}

		if t.repeat && t.running && generation == t.generation {
			t.schedule()
		}
	})
}
