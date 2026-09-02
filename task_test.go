// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import (
	"context"
	"errors"
	"math"
	"strings"
	"testing"
	"time"
)

func mountedTestTask(task *Task) *taskGroup {
	group := mountTasks(&mountContext{}, []*Task{task})
	group.started = true
	if !task.begin(group) {
		panic("test task was not queued to start")
	}
	return group
}

func drainTestTask(t *testing.T, task *Task, group *taskGroup) {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for task.Running() {
		group.drain()
		if time.Now().After(deadline) {
			t.Fatal("background task did not finish")
		}
		time.Sleep(time.Millisecond)
	}
	group.drain()
}

func TestTaskStartsBeforeMountAndCanBeCancelled(t *testing.T) {
	task := Background(nil)
	if task.Running() {
		t.Fatal("new background task should be stopped")
	}

	task.Start()
	if !task.Running() {
		t.Fatal("Start should queue an unmounted task")
	}
	task.Start()
	task.Cancel()
	if task.Running() {
		t.Fatal("Cancel should clear a queued task")
	}
}

func TestTaskDeliversProgressPostsAndDoneInOrder(t *testing.T) {
	var calls []string
	task := Background(func(_ context.Context, report *TaskReporter) error {
		if !report.Report(125, "Finishing") {
			t.Fatal("Report unexpectedly failed")
		}
		if !report.Post(func() { calls = append(calls, "post") }) {
			t.Fatal("Post unexpectedly failed")
		}
		return nil
	}).OnProgress(func(progress TaskProgress) {
		calls = append(calls, "progress")
		if progress.Percent != 100 || progress.Message != "Finishing" {
			t.Fatalf("unexpected progress: %#v", progress)
		}
	}).OnDone(func(err error) {
		calls = append(calls, "done")
		if err != nil {
			t.Fatalf("unexpected task error: %v", err)
		}
	}).AutoStart()

	group := mountedTestTask(task)
	drainTestTask(t, task, group)
	if got := strings.Join(calls, ","); got != "progress,post,done" {
		t.Fatalf("callback order = %q", got)
	}
	if progress := task.Progress(); progress.Percent != 100 || progress.Message != "Finishing" {
		t.Fatalf("Progress() = %#v", progress)
	}
	group.unmount()
}

func TestTaskNormalizesInvalidProgress(t *testing.T) {
	var values []float64
	task := Background(func(_ context.Context, report *TaskReporter) error {
		report.Report(-10)
		report.Report(math.Inf(1))
		return nil
	}).OnProgress(func(progress TaskProgress) {
		values = append(values, progress.Percent)
	}).AutoStart()

	group := mountedTestTask(task)
	drainTestTask(t, task, group)
	if len(values) != 2 || values[0] != 0 || values[1] != 0 {
		t.Fatalf("normalized progress = %v", values)
	}
	group.unmount()
}

func TestTaskCancellationReachesWorkAndDone(t *testing.T) {
	started := make(chan struct{})
	var doneErr error
	task := Background(func(ctx context.Context, _ *TaskReporter) error {
		close(started)
		<-ctx.Done()
		return ctx.Err()
	}).OnDone(func(err error) {
		doneErr = err
	}).AutoStart()

	group := mountedTestTask(task)
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("task did not start")
	}
	task.Cancel()
	drainTestTask(t, task, group)
	if !errors.Is(doneErr, context.Canceled) {
		t.Fatalf("done error = %v, want context.Canceled", doneErr)
	}
	group.unmount()
}

func TestTaskConvertsPanicsToErrors(t *testing.T) {
	var doneErr error
	task := Background(func(context.Context, *TaskReporter) error {
		panic("broken flower")
	}).OnDone(func(err error) {
		doneErr = err
	}).AutoStart()

	group := mountedTestTask(task)
	drainTestTask(t, task, group)
	if doneErr == nil || !strings.Contains(doneErr.Error(), "broken flower") {
		t.Fatalf("panic error = %v", doneErr)
	}
	group.unmount()
}

func TestTaskUnmountCancelsAndDropsLateCallbacks(t *testing.T) {
	started := make(chan struct{})
	finished := make(chan struct{})
	callbacks := 0
	task := Background(func(ctx context.Context, report *TaskReporter) error {
		close(started)
		<-ctx.Done()
		if report.Post(func() { callbacks++ }) {
			t.Error("Post should fail after unmount cancellation")
		}
		close(finished)
		return ctx.Err()
	}).OnDone(func(error) {
		callbacks++
	})
	task.Start()

	group := mountedTestTask(task)
	<-started
	group.unmount()
	select {
	case <-finished:
	case <-time.After(time.Second):
		t.Fatal("unmounted work did not observe cancellation")
	}
	if task.Running() {
		t.Fatal("unmounted task should not remain running")
	}
	if callbacks != 0 {
		t.Fatalf("callbacks after unmount = %d", callbacks)
	}
}

func TestTaskGroupDeduplicatesAndRejectsSecondOwner(t *testing.T) {
	task := Background(nil)
	first := mountTasks(&mountContext{}, []*Task{task, nil, task})
	if len(first.tasks) != 1 {
		t.Fatalf("first group has %d tasks", len(first.tasks))
	}
	second := mountTasks(&mountContext{}, []*Task{task})
	if len(second.tasks) != 0 {
		t.Fatalf("second group has %d tasks", len(second.tasks))
	}
	first.unmount()
	second.unmount()
}

func TestAutoStartIsRestoredWhenWindowReopens(t *testing.T) {
	task := Background(nil).AutoStart()
	first := mountedTestTask(task)
	drainTestTask(t, task, first)
	first.unmount()

	second := mountTasks(&mountContext{}, []*Task{task})
	second.started = true
	if !task.begin(second) {
		t.Fatal("auto-start task should begin after remount")
	}
	drainTestTask(t, task, second)
	second.unmount()
}

func TestTaskGroupPollsWhileActiveAndRefreshesDeliveredEvents(t *testing.T) {
	originalSchedule := scheduleTimer
	originalCancel := cancelTimer
	t.Cleanup(func() {
		scheduleTimer = originalSchedule
		cancelTimer = originalCancel
	})

	var scheduled []func()
	scheduleTimer = func(interval time.Duration, script ...any) string {
		if interval != backgroundPollInterval {
			t.Fatalf("poll interval = %v, want %v", interval, backgroundPollInterval)
		}
		scheduled = append(scheduled, script[0].(func()))
		return "task-poll"
	}
	cancelTimer = func(string) {}

	refreshes := 0
	progressCalls := 0
	ctx := &mountContext{refreshes: []func(){func() { refreshes++ }}}
	task := Background(func(_ context.Context, report *TaskReporter) error {
		report.Report(50, "Halfway")
		return nil
	}).OnProgress(func(TaskProgress) {
		progressCalls++
	}).AutoStart()
	group := mountTasks(ctx, []*Task{task})
	group.begin()

	deadline := time.Now().Add(time.Second)
	next := 0
	for task.Running() {
		if next < len(scheduled) {
			callback := scheduled[next]
			next++
			callback()
		} else {
			time.Sleep(time.Millisecond)
		}
		if time.Now().After(deadline) {
			t.Fatal("task poll did not deliver completion")
		}
	}
	if progressCalls != 1 || refreshes == 0 {
		t.Fatalf("progress calls = %d, refreshes = %d", progressCalls, refreshes)
	}
	if group.afterID != "" {
		t.Fatalf("idle group kept scheduled callback %q", group.afterID)
	}
	group.unmount()
}

func TestTaskGroupUnmountCancelsItsPoll(t *testing.T) {
	originalSchedule := scheduleTimer
	originalCancel := cancelTimer
	t.Cleanup(func() {
		scheduleTimer = originalSchedule
		cancelTimer = originalCancel
	})

	cancelled := ""
	scheduleTimer = func(time.Duration, ...any) string { return "background-poll" }
	cancelTimer = func(id string) { cancelled = id }

	task := Background(func(ctx context.Context, _ *TaskReporter) error {
		<-ctx.Done()
		return ctx.Err()
	}).AutoStart()
	group := mountTasks(&mountContext{}, []*Task{task})
	group.begin()
	group.unmount()

	if cancelled != "background-poll" {
		t.Fatalf("cancelled poll = %q", cancelled)
	}
}
