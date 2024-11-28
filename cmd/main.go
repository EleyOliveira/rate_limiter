package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/EleyOliveira/rate_limiter/ratelimiter"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")

		IPRequisicao, err := obterIPRequest(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			panic(err)
		}

		registro := &ratelimiter.CacheRegistro{}

		ratelimiter := ratelimiter.NewRateLimiter(registro)
		ratelimiter.Controlar(IPRequisicao, 5)

		fmt.Fprintln(w, IPRequisicao, "\n", obterTokenRequest(r))
		fmt.Println(IPRequisicao, "\n", obterTokenRequest(r))

	})
	http.ListenAndServe(":8080", nil)
}

func obterIPRequest(r *http.Request) (string, error) {
	addr, _, err := net.SplitHostPort(r.RemoteAddr)

	if err != nil {
		return "", err
	}

	return addr, nil
}

func obterTokenRequest(r *http.Request) string {
	token := r.Header.Get("API_KEY")
	return token
}
