package limiter

import (
	"github.com/ilnar/rate-limiter-assignment/policy"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TODO: Add tests for internal errors.
// TODO: Add tests for multiple limiters.

const (
	fakeAPIKey   = "fakeApiKey"
	fakeUsername = "fakeUsername"
	fakePolicy   = policy.Smooth
)

func fakeAPIKeyResolver(apiKey string) (username string, found bool) {
	if apiKey == fakeAPIKey {
		username = fakeUsername
		found = true
	}
	return
}
func fakePolicyFinder(username string) (p policy.Policy, found bool) {
	if username == fakeUsername {
		p = fakePolicy
		found = true
	}
	return
}

func TestSuccess(t *testing.T) {
	timer := func() time.Time { return time.Unix(1, 0) }
	limiters := Limiters{
		fakePolicy: func(string, time.Time) (Verdict, error) { return Verdict{Allow: true}, nil },
	}

	var gotReq *http.Request
	wrappedHandler := func(_ http.ResponseWriter, r *http.Request) { gotReq = r }

	h := CreateHandler(fakeAPIKeyResolver, fakePolicyFinder, limiters, wrappedHandler, timer)
	wantReq := httptest.NewRequest("GET", "http://example.com/foo?api_key="+fakeAPIKey, nil)
	w := httptest.NewRecorder()
	h.Handle(w, wantReq)
	if wantReq != gotReq {
		t.Errorf("Want req %v; got: %v", wantReq, gotReq)
	}
}

func TestApiKeyFailure(t *testing.T) {
	timer := func() time.Time { return time.Unix(1, 0) }
	limiters := Limiters{
		fakePolicy: func(string, time.Time) (Verdict, error) { return Verdict{Allow: true}, nil },
	}
	wrappedHandler := func(http.ResponseWriter, *http.Request) {}

	h := CreateHandler(fakeAPIKeyResolver, fakePolicyFinder, limiters, wrappedHandler, timer)
	r := httptest.NewRequest("GET", "http://example.com/foo?", nil)
	w := httptest.NewRecorder()
	h.Handle(w, r)
	want := http.StatusUnauthorized
	if w.Code != want {
		t.Errorf("Wrong HTTP status code, want %v, got: %v", want, w.Code)
	}
}

func TestRateLimit(t *testing.T) {
	timer := func() time.Time { return time.Unix(1, 0) }
	limiters := Limiters{
		fakePolicy: func(string, time.Time) (Verdict, error) {
			return Verdict{Allow: false, RetryIn: 1 * time.Second}, nil
		},
	}
	wrappedHandler := func(http.ResponseWriter, *http.Request) {}

	h := CreateHandler(fakeAPIKeyResolver, fakePolicyFinder, limiters, wrappedHandler, timer)
	r := httptest.NewRequest("GET", "http://example.com/foo?api_key="+fakeAPIKey, nil)
	w := httptest.NewRecorder()
	h.Handle(w, r)
	want := http.StatusTooManyRequests
	if w.Code != want {
		t.Errorf("Wrong HTTP status code, want %v, got: %v", want, w.Code)
	}
}
