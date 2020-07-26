package main

import (
	"flag"
	"fmt"
	"github.com/ilnar/rate-limiter-assignment/auth"
	"github.com/ilnar/rate-limiter-assignment/limiter"
	"github.com/ilnar/rate-limiter-assignment/policy"
	"github.com/ilnar/rate-limiter-assignment/precise"
	"log"
	"net/http"
	"time"
)

var httpAddr = flag.String("http_address", ":8080", "TCP network address to listen on for HTTP requests.")

// TODO: Extend here.
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK, come in, %s!", r.URL.Path[1:])
}

// TODO: Unhardcode.
func limitersConfig() limiter.Limiters {
	smooth, err := precise.Create(1, 1*time.Second)
	if err != nil {
		log.Fatalf("Error creating high throughput rate limiter: %v", err)
	}
	spiky, err := precise.Create(10, 1*time.Minute)
	if err != nil {
		log.Fatalf("Error creating low throughput rate limiter: %v", err)
	}
	return limiter.Limiters{
		policy.Default: spiky,
		policy.Smooth:  smooth,
		policy.Spiky:   spiky,
	}
}

func main() {
	h := limiter.CreateHandler(auth.AuthenticateByAPIKey, policy.FindByUsername, limitersConfig(), handler, time.Now)
	http.HandleFunc("/", h.Handle)
	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}
