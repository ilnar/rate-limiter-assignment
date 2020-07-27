package limiter

import (
	"fmt"
	"github.com/ilnar/rate-limiter-assignment/policy"
	"log"
	"net/http"
	"time"
)

const API_KEY_QUERY_VAR = "api_key"

// APIKeyResolver defines an interface for getting username by API key.
type APIKeyResolver func(apiKey string) (string, bool)

// PolicyFinder defines an interface for getting Policy by username.
type PolicyFinder func(username string) (policy.Policy, bool)

// Limiter defines a type for storing different implementations of Limiter.
type Limiters map[policy.Policy]Limiter

// WrappedHandler defines an interface for the request handler being called by limiter in case of success.
type WrappedHandler func(http.ResponseWriter, *http.Request)

// Timer defines an interface for getting timestamp of a request.
type Timer func() time.Time

// Handler provides an HTTP handler supporting rate limiting.
type Handler struct {
	apiKeyResolver APIKeyResolver
	policyFinder   PolicyFinder
	limiters       Limiters
	handler        WrappedHandler
	timer          Timer
}

// CreateHandler creates a new handler.
func CreateHandler(ar APIKeyResolver, pf PolicyFinder, ls Limiters, wh WrappedHandler, t Timer) *Handler {
	return &Handler{apiKeyResolver: ar, policyFinder: pf, limiters: ls, handler: wh, timer: t}
}

// Handle examines request URI, identifies corresponding user and applies rate limiting.
// It does not read request body.
func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	username, err := h.getUser(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Invalid API KEY")
		return
	}

	l, err := h.getLimiter(username)
	if err != nil {
		log.Printf("Error finding rate limiter for user %q: %v", username, err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal error")
		return
	}

	t := h.timer()
	v, err := l(username, t)
	if err != nil {
		log.Printf("Error calling rate limiter for user %q and time %v: %v", username, t, err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal error")
		return
	}
	if !v.Allow {
		w.WriteHeader(http.StatusTooManyRequests)
		fmt.Fprintf(w, "Rate limit exceeded. Try again in %v seconds", v.RetryIn.Seconds())
		return
	}
	h.handler(w, r)
}

func (h Handler) getUser(r *http.Request) (string, error) {
	q := r.URL.Query()
	if len(q[API_KEY_QUERY_VAR]) != 1 {
		return "", fmt.Errorf("missing API key in query: %v", q)
	}
	apiKey := q[API_KEY_QUERY_VAR][0]
	username, found := h.apiKeyResolver(apiKey)
	if !found {
		return "", fmt.Errorf("user not found for API key: %q", apiKey)
	}
	return username, nil
}

func (h Handler) getLimiter(username string) (Limiter, error) {
	p, found := h.policyFinder(username)
	if !found {
		return nil, fmt.Errorf("policy not found for user: %v", username)
	}

	l, found := h.limiters[p]
	if !found {
		return nil, fmt.Errorf("limiter not found for policy: %v", p)
	}

	return l, nil
}
