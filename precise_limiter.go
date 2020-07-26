// Package preciselimiter contains simple implementation of high precision rate limiter.
package preciselimiter

import (
	"fmt"
	"time"
)

// Verdict indicates if request should be discarded.
type Verdict struct {
	Allow   bool
	RetryIn time.Duration
}

// Limiter defines interface allow to check if a new request exceeds limit.
type Limiter func(userID string, ts time.Time) (Verdict, error)

// Create creates a new precise rate limiter.
func Create(limit int, duration time.Duration) (Limiter, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("limit must be greater than zero, got %v", limit)
	}
	pl := preciseLimiter{limit, duration, map[string][]time.Time{}}
	return pl.processRequest, nil
}

type preciseLimiter struct {
	limit    int
	duration time.Duration

	requestTimesPerUser map[string][]time.Time
}

func (pl *preciseLimiter) processRequest(userID string, t time.Time) (Verdict, error) {
	if _, found := pl.requestTimesPerUser[userID]; !found {
		pl.requestTimesPerUser[userID] = []time.Time{t}
		return Verdict{Allow: true, RetryIn: 0}, nil
	}

	rts := pl.requestTimesPerUser[userID]
	// Check that we are not travelling back in time.
	if l := len(rts); l > 0 && t.Before(rts[l-1]) {
		return Verdict{}, fmt.Errorf("request time is not monotonic, want it to be after %v; got %v", rts[l-1], t)
	}

	// Truncate the list of timestamps.
	oldest := t.Add(-pl.duration)
	pl.requestTimesPerUser[userID] = rts[findFirstNewerThan(rts, oldest):]
	// Check size.
	if len(rts) >= pl.limit {
		return Verdict{Allow: false, RetryIn: rts[0].Sub(oldest)}, nil
	}

	// Add new request timestamp.
	pl.requestTimesPerUser[userID] = append(pl.requestTimesPerUser[userID], t)
	return Verdict{Allow: true, RetryIn: 0}, nil
}

func findFirstNewerThan(ts []time.Time, t time.Time) int {
	// TODO: Implement binary search here.
	for i, ct := range ts {
		if ct.After(t) {
			return i
		}
	}
	return len(ts)
}
