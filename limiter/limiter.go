package limiter

import "time"

// Verdict indicates if request should be discarded.
type Verdict struct {
	Allow   bool
	RetryIn time.Duration
}

// Limiter defines interface allow to check if a new request exceeds limit.
type Limiter func(userID string, ts time.Time) (Verdict, error)
