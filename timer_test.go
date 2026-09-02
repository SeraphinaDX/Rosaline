// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import (
	"testing"
	"time"
)

func TestEveryStartsAndCanBeStopped(t *testing.T) {
	timer := Every(time.Second, nil)
	if !timer.Running() {
		t.Fatal("Every timer should start enabled")
	}

	timer.Stop()
	if timer.Running() {
		t.Fatal("Stop should stop a timer")
	}

	timer.Start()
	if !timer.Running() {
		t.Fatal("Start should start a stopped timer")
	}
}

func TestAfterIsOneShot(t *testing.T) {
	timer := After(time.Second, func() {})
	if timer.repeat {
		t.Fatal("After should create a one-shot timer")
	}
	if timer.interval != time.Second {
		t.Fatalf("After interval = %v, want %v", timer.interval, time.Second)
	}
}

func TestShortTimerIntervalsAreSafe(t *testing.T) {
	timer := Every(0, func() {})
	if timer.interval != minimumTimerInterval {
		t.Fatalf("interval = %v, want %v", timer.interval, minimumTimerInterval)
	}
}

func TestAnimateUsesFramesPerSecond(t *testing.T) {
	timer := Animate(20, func() {})
	want := 50 * time.Millisecond
	if timer.interval != want {
		t.Fatalf("Animate interval = %v, want %v", timer.interval, want)
	}
	if !timer.repeat {
		t.Fatal("Animate should create a repeating timer")
	}
}

func TestAnimateUsesSafeFrameRateLimits(t *testing.T) {
	defaultTimer := Animate(0, nil)
	if defaultTimer.interval != time.Second/defaultAnimationFPS {
		t.Fatalf("default interval = %v, want %v", defaultTimer.interval, time.Second/defaultAnimationFPS)
	}

	fastTimer := Animate(2000, nil)
	if fastTimer.interval != time.Millisecond {
		t.Fatalf("limited interval = %v, want %v", fastTimer.interval, time.Millisecond)
	}
}

func TestMountTimersIgnoresNilAndDuplicates(t *testing.T) {
	timer := Every(time.Second, func() {})
	ctx := &mountContext{}
	mounted := mountTimers(ctx, []*Timer{nil, timer, timer})
	if len(mounted) != 1 {
		t.Fatalf("mounted %d timers, want 1", len(mounted))
	}
	if timer.ctx != ctx {
		t.Fatal("timer should retain its application context")
	}
	timer.unmount(ctx)
}

func TestRepeatingTimerCallsBackRefreshesAndSchedulesAgain(t *testing.T) {
	originalSchedule := scheduleTimer
	originalCancel := cancelTimer
	t.Cleanup(func() {
		scheduleTimer = originalSchedule
		cancelTimer = originalCancel
	})

	var scheduled []func()
	scheduleTimer = func(_ time.Duration, script ...any) string {
		scheduled = append(scheduled, script[0].(func()))
		return "timer"
	}
	cancelTimer = func(string) {}

	ticks := 0
	refreshes := 0
	ctx := &mountContext{refreshes: []func(){func() { refreshes++ }}}
	timer := Every(time.Second, func() { ticks++ })
	timer.mount(ctx)
	timer.begin()

	if len(scheduled) != 1 {
		t.Fatalf("scheduled %d callbacks, want 1", len(scheduled))
	}
	scheduled[0]()
	if ticks != 1 {
		t.Fatalf("ticks = %d, want 1", ticks)
	}
	if refreshes != 1 {
		t.Fatalf("refreshes = %d, want 1", refreshes)
	}
	if len(scheduled) != 2 {
		t.Fatalf("scheduled %d callbacks after tick, want 2", len(scheduled))
	}
}

func TestOneShotTimerStopsAfterCallback(t *testing.T) {
	originalSchedule := scheduleTimer
	originalCancel := cancelTimer
	t.Cleanup(func() {
		scheduleTimer = originalSchedule
		cancelTimer = originalCancel
	})

	var callback func()
	scheduleTimer = func(_ time.Duration, script ...any) string {
		callback = script[0].(func())
		return "timer"
	}
	cancelTimer = func(string) {}

	calls := 0
	ctx := &mountContext{}
	timer := After(time.Second, func() { calls++ })
	timer.mount(ctx)
	timer.begin()
	callback()

	if calls != 1 {
		t.Fatalf("calls = %d, want 1", calls)
	}
	if timer.Running() {
		t.Fatal("one-shot timer should stop after its callback")
	}
}

func TestRestartCancelsCurrentWait(t *testing.T) {
	originalSchedule := scheduleTimer
	originalCancel := cancelTimer
	t.Cleanup(func() {
		scheduleTimer = originalSchedule
		cancelTimer = originalCancel
	})

	schedules := 0
	cancelled := ""
	scheduleTimer = func(_ time.Duration, _ ...any) string {
		schedules++
		return "timer-id"
	}
	cancelTimer = func(id string) { cancelled = id }

	ctx := &mountContext{}
	timer := After(time.Second, func() {})
	timer.mount(ctx)
	timer.begin()
	timer.Restart()

	if cancelled != "timer-id" {
		t.Fatalf("cancelled %q, want %q", cancelled, "timer-id")
	}
	if schedules != 2 {
		t.Fatalf("scheduled %d waits, want 2", schedules)
	}
}
