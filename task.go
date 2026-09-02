// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"
)

const (
	backgroundPollInterval  = 16 * time.Millisecond
	backgroundEventBuffer   = 64
	backgroundEventsPerTick = 256
)

// TaskProgress is a progress update from background work. Percent is always
// between zero and 100. Message is optional application-defined status text.
type TaskProgress struct {
	Percent float64
	Message string
}

// TaskReporter safely sends information from a background Task to Rosaline's
// GUI thread. Applications receive one from Background rather than creating it.
type TaskReporter struct {
	run *taskRun
}

// Report sends percentage progress and optional status text to OnProgress. It
// returns false after the task is cancelled or its window closes.
func (r *TaskReporter) Report(percent float64, message ...string) bool {
	if r == nil || r.run == nil {
		return false
	}
	if math.IsNaN(percent) || math.IsInf(percent, 0) {
		percent = 0
	}
	percent = max(0, min(100, percent))
	text := ""
	if len(message) > 0 {
		text = message[0]
	}
	return r.run.send(taskEvent{
		kind:     taskProgressEvent,
		progress: TaskProgress{Percent: percent, Message: text},
	})
}

// Post schedules callback on Rosaline's GUI thread. Use it when background
// work has produced a result that should change application state or a widget.
// It returns false after the task is cancelled or its window closes.
func (r *TaskReporter) Post(callback func()) bool {
	if r == nil || r.run == nil || callback == nil {
		return false
	}
	return r.run.send(taskEvent{kind: taskPostEvent, callback: callback})
}

// Task represents reusable background work owned by an App or WindowOptions.
// Create one with Background.
type Task struct {
	mu         sync.Mutex
	work       func(context.Context, *TaskReporter) error
	onProgress func(TaskProgress)
	onDone     func(error)
	progress   TaskProgress
	autoStart  bool
	requested  bool
	running    bool
	generation uint64
	group      *taskGroup
	run        *taskRun
}

type taskEventKind uint8

const (
	taskProgressEvent taskEventKind = iota
	taskPostEvent
	taskDoneEvent
)

type taskEvent struct {
	kind       taskEventKind
	generation uint64
	progress   TaskProgress
	callback   func()
	err        error
}

type taskRun struct {
	ctx        context.Context
	cancel     context.CancelFunc
	generation uint64
	events     chan taskEvent
	lifecycle  <-chan struct{}
}

// Background creates stopped background work. The work function runs in a Go
// goroutine after Start. Include the Task in App.Tasks or WindowOptions.Tasks.
func Background(work func(context.Context, *TaskReporter) error) *Task {
	if work == nil {
		work = func(context.Context, *TaskReporter) error { return nil }
	}
	return &Task{work: work}
}

// OnProgress sets the GUI-thread callback for Report updates.
func (t *Task) OnProgress(callback func(TaskProgress)) *Task {
	if t != nil {
		t.mu.Lock()
		t.onProgress = callback
		t.mu.Unlock()
	}
	return t
}

// OnDone sets the GUI-thread callback invoked when work finishes. Cancellation
// is reported with context.Canceled and can be checked with errors.Is.
func (t *Task) OnDone(callback func(error)) *Task {
	if t != nil {
		t.mu.Lock()
		t.onDone = callback
		t.mu.Unlock()
	}
	return t
}

// AutoStart starts the task when its window opens. A reusable secondary window
// starts the task again each time it is opened.
func (t *Task) AutoStart() *Task {
	if t == nil {
		return t
	}
	t.mu.Lock()
	t.autoStart = true
	t.mu.Unlock()
	t.Start()
	return t
}

// Start begins stopped work. Calling Start while the task is running has no
// effect. Before RunApp, it queues the task to begin when its window opens.
// Call task controls from Rosaline callbacks, not background goroutines.
func (t *Task) Start() {
	if t == nil {
		return
	}
	t.mu.Lock()
	if t.running || t.requested {
		t.mu.Unlock()
		return
	}
	t.requested = true
	group := t.group
	if group != nil && group.started {
		t.launchLocked(group)
	}
	t.mu.Unlock()
	if group != nil && group.started {
		group.wake()
	}
}

// Cancel asks running work to stop through its context. Work should observe
// ctx.Done or return when Report or Post returns false.
func (t *Task) Cancel() {
	if t == nil {
		return
	}
	t.mu.Lock()
	t.requested = false
	if t.run != nil {
		t.run.cancel()
	}
	t.mu.Unlock()
}

// Running reports whether work is running or queued to start with its window.
func (t *Task) Running() bool {
	if t == nil {
		return false
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.running || t.requested
}

// Progress returns the most recently delivered progress update.
func (t *Task) Progress() TaskProgress {
	if t == nil {
		return TaskProgress{}
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.progress
}

func (t *Task) mount(group *taskGroup) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.group != nil && t.group != group {
		return false
	}
	t.group = group
	if t.autoStart {
		t.requested = true
	}
	return true
}

func (t *Task) begin(group *taskGroup) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.group != group || !t.requested || t.running {
		return false
	}
	t.launchLocked(group)
	return true
}

func (t *Task) launchLocked(group *taskGroup) {
	t.generation++
	ctx, cancel := context.WithCancel(context.Background())
	run := &taskRun{
		ctx:        ctx,
		cancel:     cancel,
		generation: t.generation,
		events:     make(chan taskEvent, backgroundEventBuffer),
		lifecycle:  group.lifecycle,
	}
	t.requested = false
	t.running = true
	t.progress = TaskProgress{}
	t.run = run
	work := t.work
	go executeTask(work, run)
}

func executeTask(work func(context.Context, *TaskReporter) error, run *taskRun) {
	var err error
	func() {
		defer func() {
			if recovered := recover(); recovered != nil {
				err = fmt.Errorf("rosaline: background task panicked: %v", recovered)
			}
		}()
		err = work(run.ctx, &TaskReporter{run: run})
	}()
	if err == nil && run.ctx.Err() != nil {
		err = run.ctx.Err()
	}
	event := taskEvent{kind: taskDoneEvent, generation: run.generation, err: err}
	select {
	case run.events <- event:
	case <-run.lifecycle:
	}
}

func (r *taskRun) send(event taskEvent) bool {
	if r.ctx.Err() != nil {
		return false
	}
	select {
	case <-r.lifecycle:
		return false
	default:
	}
	event.generation = r.generation
	select {
	case r.events <- event:
		return true
	case <-r.ctx.Done():
		return false
	case <-r.lifecycle:
		return false
	}
}

func (t *Task) unmount(group *taskGroup) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.group != group {
		return
	}
	if t.run != nil {
		t.run.cancel()
	}
	t.generation++
	t.group = nil
	t.run = nil
	t.running = false
	t.requested = t.autoStart
}

func (t *Task) active(group *taskGroup) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.group == group && (t.running || t.requested)
}

func (t *Task) drain(group *taskGroup) bool {
	t.mu.Lock()
	if t.group != group || t.run == nil {
		t.mu.Unlock()
		return false
	}
	run := t.run
	t.mu.Unlock()

	refreshed := false
	for range backgroundEventsPerTick {
		select {
		case event := <-run.events:
			if t.handle(group, run, event) {
				refreshed = true
			}
		default:
			return refreshed
		}
	}
	return refreshed
}

func (t *Task) handle(group *taskGroup, run *taskRun, event taskEvent) bool {
	t.mu.Lock()
	if t.group != group || t.run != run || event.generation != t.generation {
		t.mu.Unlock()
		return false
	}

	switch event.kind {
	case taskProgressEvent:
		t.progress = event.progress
		callback := t.onProgress
		t.mu.Unlock()
		if callback != nil {
			callback(event.progress)
		}
	case taskPostEvent:
		callback := event.callback
		t.mu.Unlock()
		if callback != nil {
			callback()
		}
	case taskDoneEvent:
		t.running = false
		t.run = nil
		callback := t.onDone
		t.mu.Unlock()
		if callback != nil {
			callback(event.err)
		}
	default:
		t.mu.Unlock()
		return false
	}
	return true
}

type taskGroup struct {
	ctx        *mountContext
	tasks      []*Task
	lifecycle  chan struct{}
	afterID    string
	generation uint64
	started    bool
	closed     bool
}

func mountTasks(ctx *mountContext, tasks []*Task) *taskGroup {
	group := &taskGroup{ctx: ctx, lifecycle: make(chan struct{})}
	seen := make(map[*Task]bool, len(tasks))
	for _, task := range tasks {
		if task == nil || seen[task] {
			continue
		}
		seen[task] = true
		if task.mount(group) {
			group.tasks = append(group.tasks, task)
		}
	}
	return group
}

func (g *taskGroup) begin() {
	if g == nil || g.closed || g.started {
		return
	}
	g.started = true
	active := false
	for _, task := range g.tasks {
		active = task.begin(g) || active
	}
	if active {
		g.wake()
	}
}

func (g *taskGroup) wake() {
	if g == nil || g.closed || !g.started || g.afterID != "" || g.ctx == nil || g.ctx.closed {
		return
	}
	g.generation++
	generation := g.generation
	g.afterID = scheduleTimer(backgroundPollInterval, func() {
		if g.closed || generation != g.generation || g.ctx == nil || g.ctx.closed {
			return
		}
		g.afterID = ""
		if g.drain() {
			g.ctx.refresh()
		}
		if g.anyActive() {
			g.wake()
		}
	})
}

func (g *taskGroup) drain() bool {
	if g == nil || g.closed {
		return false
	}
	refreshed := false
	for _, task := range g.tasks {
		refreshed = task.drain(g) || refreshed
	}
	return refreshed
}

func (g *taskGroup) anyActive() bool {
	if g == nil || g.closed {
		return false
	}
	for _, task := range g.tasks {
		if task.active(g) {
			return true
		}
	}
	return false
}

func (g *taskGroup) unmount() {
	if g == nil || g.closed {
		return
	}
	g.closed = true
	g.generation++
	if g.afterID != "" && g.ctx != nil && !g.ctx.closed {
		cancelTimer(g.afterID)
	}
	g.afterID = ""
	close(g.lifecycle)
	for _, task := range g.tasks {
		task.unmount(g)
	}
	g.tasks = nil
}
