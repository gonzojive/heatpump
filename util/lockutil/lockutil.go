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
