package precise

import (
	"github.com/google/go-cmp/cmp"
	"github.com/ilnar/rate-limiter-assignment/limiter"
	"testing"
	"time"
)

const fakeUser = "fakeUser"

func TestCreate(t *testing.T) {
	if _, err := Create(0, 1*time.Minute); err == nil {
		t.Error("Expected failure due to zero limit")
	}
	if _, err := Create(1, 0*time.Minute); err == nil {
		t.Error("Expected failure due to zero duration")
	}
	if _, err := Create(2, 1*time.Minute); err != nil {
		t.Errorf("Error creating rate limiter: %v", err)
	}
}

func TestDetectWrongTime(t *testing.T) {
	l, err := Create(2, 1*time.Minute)
	if err != nil {
		t.Fatalf("Error creating rate limiter: %v", err)
	}
	if v, err := l(fakeUser, time.Unix(30, 0)); err != nil || !v.Allow {
		t.Fatalf("Error making first request, want it to be allowed, got: %v (error=%v)", v, err)
	}
	got, err := l(fakeUser, time.Unix(29, 0))
	if err == nil {
		t.Errorf("Expected error, got: %v", got)
	}
}

func TestEnforceLimitPerUser(t *testing.T) {
	l, err := Create(1, 1*time.Minute)
	if err != nil {
		t.Fatalf("Error creating rate limiter: %v", err)
	}

	firstUser := "firstUser"
	if v, err := l(firstUser, time.Unix(30, 0)); err != nil || !v.Allow {
		t.Errorf("Error making request for user %v, want it to be allowed, got: %v (error=%v)", firstUser, v, err)
	}
	if v, err := l(firstUser, time.Unix(31, 0)); err != nil || v.Allow {
		t.Errorf("Want second request for user %v to fail, got: %v (error=%v)", firstUser, v, err)
	}

	secondUser := "secondUser"
	if v, err := l(secondUser, time.Unix(30, 0)); err != nil || !v.Allow {
		t.Errorf("Error making request for user %v, want it to be allowed, got: %v (error=%v)", secondUser, v, err)
	}
	if v, err := l(secondUser, time.Unix(31, 0)); err != nil || v.Allow {
		t.Errorf("Want second request for user %v to fail, got: %v (error=%v)", secondUser, v, err)
	}
}

func TestHandleSimultaneousRequests(t *testing.T) {
	l, err := Create(1, 1*time.Minute)
	if err != nil {
		t.Fatalf("Error creating rate limiter: %v", err)
	}
	if v, err := l(fakeUser, time.Unix(30, 0)); err != nil || !v.Allow {
		t.Errorf("Error making first request, want it to be allowed, got: %v (error=%v)", v, err)
	}
	if v, err := l(fakeUser, time.Unix(30, 0)); err != nil || v.Allow {
		t.Errorf("Want second request to fail, got: %v (error=%v)", v, err)
	}
}

func TestThrottleCorrectly(t *testing.T) {
	l, err := Create(2, 1*time.Minute)
	if err != nil {
		t.Fatalf("Error creating rate limiter: %v", err)
	}
	// Rate spike in the middle of minute.
	if v, err := l(fakeUser, time.Unix(29, 0)); err != nil || !v.Allow {
		t.Fatalf("Error making first request, want it to be allowed, got: %v (error=%v)", v, err)
	}
	if v, err := l(fakeUser, time.Unix(30, 0)); err != nil || !v.Allow {
		t.Fatalf("Error making second request, want it to be allowed, got: %v (error=%v)", v, err)
	}
	want := limiter.Verdict{Allow: false, RetryIn: 30 * time.Second}
	got, err := l(fakeUser, time.Unix(59, 0))
	if !cmp.Equal(want, got) {
		t.Errorf("Unexpected result, want %v; got %v", want, got)
	}
	if v, err := l(fakeUser, time.Unix(90, 0)); err != nil || !v.Allow {
		t.Errorf("Error making last request, want it to be allowed, got: %v (error=%v)", v, err)
	}
}

func TestNewRequestAfterLongDelay(t *testing.T) {
	l, err := Create(2, 1*time.Minute)
	if err != nil {
		t.Fatalf("Error creating rate limiter: %v", err)
	}
	if v, err := l(fakeUser, time.Unix(29, 0)); err != nil || !v.Allow {
		t.Fatalf("Error making first request, want it to be allowed, got: %v (error=%v)", v, err)
	}
	if v, err := l(fakeUser, time.Unix(30, 0)); err != nil || !v.Allow {
		t.Fatalf("Error making second request, want it to be allowed, got: %v (error=%v)", v, err)
	}
	if v, err := l(fakeUser, time.Unix(1030, 0)); err != nil || !v.Allow {
		t.Errorf("Error making delayed request, want it to be allowed, got: %v (error=%v)", v, err)
	}
}
