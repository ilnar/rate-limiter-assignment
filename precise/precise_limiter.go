// Package precise contains simple implementation of high precision rate limiter.
package precise

import (
	"fmt"
	"github.com/ilnar/rate-limiter-assignment/limiter"
	"time"
)

// Create creates a new precise rate limiter.
func Create(limit int, duration time.Duration) (limiter.Limiter, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("want limit to be greater than zero, got %v", limit)
	}

	if duration.Nanoseconds() <= 0 {
		return nil, fmt.Errorf("want duration to be greater than zero, got %v", duration)
	}

	pl := preciseLimiter{limit, duration, map[string][]time.Time{}}
	return pl.processRequest, nil
}

type preciseLimiter struct {
	limit    int
	duration time.Duration

	requestTimesPerUser map[string][]time.Time
}

func (pl *preciseLimiter) processRequest(userID string, t time.Time) (limiter.Verdict, error) {
	if _, found := pl.requestTimesPerUser[userID]; !found {
		pl.requestTimesPerUser[userID] = []time.Time{t}
		return limiter.Verdict{Allow: true}, nil
	}

	rts := pl.requestTimesPerUser[userID]
	// Check that we are not travelling back in time.
	if l := len(rts); l > 0 && t.Before(rts[l-1]) {
		return limiter.Verdict{}, fmt.Errorf("request time is not monotonic, want it to be after %v; got %v", rts[l-1], t)
	}

	// Truncate the list of timestamps.
	oldest := t.Add(-pl.duration)
	rts = rts[findFirstNewerThan(rts, oldest):]
	pl.requestTimesPerUser[userID] = rts

	// Check size.
	if len(rts) >= pl.limit {
		return limiter.Verdict{Allow: false, RetryIn: rts[0].Sub(oldest)}, nil
	}

	// Add new request timestamp.
	pl.requestTimesPerUser[userID] = append(pl.requestTimesPerUser[userID], t)
	return limiter.Verdict{Allow: true}, nil
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
