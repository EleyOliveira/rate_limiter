package main

import (
	"fmt"
	"net/http"

	"github.com/EleyOliveira/rate_limiter/ratelimiter"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		registro := &ratelimiter.CacheRegistro{}

		ratelimiter := ratelimiter.NewRateLimiter(registro)
		statusCode, err := ratelimiter.Controlar(r, 5, 5, 60)
		if err != nil {
			w.WriteHeader(statusCode)
			fmt.Fprintln(w, err.Error())
			return
		}

		w.WriteHeader(statusCode)
		fmt.Fprintln(w, "Sucesso")

	})

	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {

		registro := &ratelimiter.CacheRegistro{}

		ratelimiter := ratelimiter.NewRateLimiter(registro)
		token, err := ratelimiter.GerarToken(60)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, token)
	})

	http.ListenAndServe(":8080", nil)
}
