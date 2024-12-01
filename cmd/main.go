package main

import (
	"fmt"
	"net/http"

	"github.com/EleyOliveira/rate_limiter/ratelimiter"
	"github.com/google/uuid"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		registro := &ratelimiter.CacheRegistro{}

		ratelimiter := ratelimiter.NewRateLimiter(registro)
		ratelimiter.Controlar(r, 5, 5, 60)

		w.WriteHeader(http.StatusOK)

	})

	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		novoToken := uuid.New().String()
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, novoToken)
	})

	http.ListenAndServe(":8080", nil)
}
