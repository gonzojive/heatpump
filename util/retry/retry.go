// Package retry provides retry loop facilities.
package retry

import (
	"context"
	"fmt"
	"time"
)

const defaultMinTimeBetweenAttempts = time.Second * 15

type Config struct {
	// IsRetriable reports if the error should be retried or not.
	IsRetriable func(err error) bool
	RawErrors   bool

	Sleep func(ctx context.Context, prevAttemptStart time.Time) error
}

type loop struct {
	cfg       *Config
	attempt   int
	lastStart time.Time
}

func (cfg *Config) Start(ctx context.Context, fn func(ctx context.Context) error) error {
	isRetriable := cfg.IsRetriable
	if isRetriable == nil {
		isRetriable = func(err error) bool { return true }
	}
	sleep := cfg.Sleep
	if sleep == nil {
		sleep = func(ctx context.Context, prevAttemptStart time.Time) error {
			select {
			case <-ctx.Done():
				return fmt.Errorf("cancelation occurred while retrying: %w", ctx.Err())
			default:
			}
			later := prevAttemptStart.Add(defaultMinTimeBetweenAttempts)
			now := time.Now()
			if now.Before(later) {
				return nil
			}
			sleepTime := later.Sub(now)
			timer := time.NewTimer(sleepTime)
			defer timer.Stop()
			select {
			case <-timer.C:
				return nil
			case <-ctx.Done():
				return fmt.Errorf("cancelation occurred while retrying: %w", ctx.Err())
			}
		}
	}

	loop := &loop{cfg: cfg}
	for {
		loop.attempt++
		loop.lastStart = time.Now()

		err := fn(ctx)
		if err == nil {
			return nil
		}
		if !isRetriable(err) {
			if cfg.RawErrors {
				return err
			}
			return fmt.Errorf("retry %d failed with error %w", loop.attempt+1, err)
		}

		if err := sleep(ctx, loop.lastStart); err != nil {
			return err
		}
	}
}
