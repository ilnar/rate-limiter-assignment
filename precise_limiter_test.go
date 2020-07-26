package preciselimiter

import (
	"github.com/google/go-cmp/cmp"
	"testing"
	"time"
)

const fakeUser = "fakeUser"

func TestThrottle(t *testing.T) {
	l, err := Create(2, 1*time.Minute)
	if err != nil {
		t.Fatalf("Error creating rate limiter: %v", err)
	}
	// Rate spike in the middle of minute.
	if v, err := l(fakeUser, time.Unix(29, 0)); err != nil || !v.Allow {
		t.Fatalf("Error making first request, want it to be allowed, got: %v (error=%v)", v, err)
	}
	if v, err := l(fakeUser, time.Unix(30, 0)); err != nil || !v.Allow {
		t.Fatalf("Error making first request, want it to be allowed, got: %v (error=%v)", v, err)
	}
	want := Verdict{Allow: false, RetryIn: 30 * time.Second}
	got, err := l(fakeUser, time.Unix(59, 0))
	if !cmp.Equal(want, got) {
		t.Errorf("Unexpected result, want %v; got %v", want, got)
	}

}
