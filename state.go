// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import "sync"

// State stores a value that can be safely read and changed.
// LabelFunc and button callbacks are enough for basic reactive interfaces.
type State[T any] struct {
	mu    sync.RWMutex
	value T
}

// NewState creates state with an initial value.
func NewState[T any](initial T) *State[T] { return &State[T]{value: initial} }

// Get returns the current value.
func (s *State[T]) Get() T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.value
}

// Set replaces the current value.
func (s *State[T]) Set(value T) {
	s.mu.Lock()
	s.value = value
	s.mu.Unlock()
}

// Update calculates and stores a new value from the current value.
func (s *State[T]) Update(change func(T) T) {
	if change == nil {
		return
	}
	s.mu.Lock()
	s.value = change(s.value)
	s.mu.Unlock()
}
