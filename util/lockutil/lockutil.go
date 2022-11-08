// Package lockutil provides a sync.Lock implementation that ensures a certain
// interval of time has passed since the last time a lock was acquired.
package lockutil

import (
	"sync"
	"time"
)

// WithGuaranteedTimeSinceLastRelease returns a sync.Locker based on another
// locker that guarantees the lock can only be acquired a certain period after
// the most recent Unlock.
func WithGuaranteedTimeSinceLastRelease(l sync.Locker, minTimeSinceUnlock time.Duration) sync.Locker {
	return &timeRestrictedLock{
		l,
		minTimeSinceUnlock,
		time.Time{},
		time.Now,
	}
}

type timeRestrictedLock struct {
	underlying         sync.Locker
	minTimeSinceUnlock time.Duration
	lastReleased       time.Time
	now                func() time.Time
}

func (l *timeRestrictedLock) Lock() {
	l.underlying.Lock()
	remaining := l.minTimeSinceUnlock - l.now().Sub(l.lastReleased)
	if remaining < 0 {
		return
	}
	time.Sleep(remaining)
}

func (l *timeRestrictedLock) Unlock() {
	l.lastReleased = l.now()
	l.underlying.Unlock()
}

// LockedValue holds a value of a particular type with a sync.Mutex.
type LockedValue[T any] struct {
	mu  sync.Mutex
	val T
}

// NewLockedValue returns a new locked value that holds the provided value.
func NewLockedValue[T any](v T) *LockedValue[T] { return &LockedValue[T]{val: v} }

func (lv *LockedValue[T]) Store(v T) {
	lv.mu.Lock()
	defer lv.mu.Unlock()
	lv.val = v
}

func (lv *LockedValue[T]) Load() T {
	lv.mu.Lock()
	defer lv.mu.Unlock()
	return lv.val
}

// Update calls the callback with the old value while the lock is held and then
// sets the value to the result of the callback. Care must be taken by the
// callback to avoid deadlock.
func (lv *LockedValue[T]) Update(update func(old T) T) {
	lv.mu.Lock()
	defer lv.mu.Unlock()
	lv.val = update(lv.val)
}

// Observe calls the callback with the current value while the lock is held.
// Care must be taken by the callback to avoid deadlock.
func (lv *LockedValue[T]) Observe(fn func(val T)) {
	lv.mu.Lock()
	defer lv.mu.Unlock()
	fn(lv.val)
}

// Compute calls the callback with the current value while the lock is held and
// computes a value. Care must be taken by the callback to avoid deadlock.
func Compute[T, O any](lv *LockedValue[T], fn func(val T) O) O {
	var out O
	lv.Observe(func(val T) {
		out = fn(val)
	})
	return out
}

// GetMapValue returns the value of a key within a locked map.
func GetMapValue[K comparable, V any](lv *LockedValue[map[K]V], key K) V {
	var out V
	lv.Observe(func(val map[K]V) {
		out = val[key]
	})
	return out
}

// DeleteMapValue returns the value of a key within a locked map.
func DeleteMapValue[K comparable, V any](lv *LockedValue[map[K]V], key K) {
	lv.Update(func(m map[K]V) map[K]V {
		delete(m, key)
		return m
	})
}

// SetMapValue sets the value of a value within a locked map.
func SetMapValue[K comparable, V any](lv *LockedValue[map[K]V], key K, value V) {
	lv.Update(func(m map[K]V) map[K]V {
		m[key] = value
		return m
	})
}

// LoadOrStoreMapValue sets a value in a locked map
func LoadOrStoreMapValue[K comparable, V any](
	lv *LockedValue[map[K]V],
	key K,
	value V) (actual V, loaded bool) {

	lv.Update(func(m map[K]V) map[K]V {
		actual, loaded = m[key]
		if !loaded {
			actual = value
			m[key] = actual
		}
		return m
	})
	return
}
