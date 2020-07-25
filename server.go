package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var httpAddr = flag.String("http_address", ":8080", "TCP network address to listen on for HTTP requests.")

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	throttledHandler := func(w http.ResponseWriter, r *http.Request) {
		// Apply rate limiter.
		handler(w, r)
	}
	http.HandleFunc("/", throttledHandler)
	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}
