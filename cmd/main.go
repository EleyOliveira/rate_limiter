package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/EleyOliveira/rate_limiter/ratelimiter"
	"github.com/google/uuid"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "text/plain")

		IPRequisicao, err := obterIPRequest(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			panic(err)
		}

		registro := &ratelimiter.CacheRegistro{}

		requisicaoPorSegundo := 5
		totalTempoBloqueado := 2
		totalSegundosExpiracaoToken := 60

		id := IPRequisicao

		token := obterTokenRequest(r)
		if token != "" {
			requisicaoPorSegundo = 10
			totalTempoBloqueado = 3
			id = token
		}
		ratelimiter := ratelimiter.NewRateLimiter(registro)
		ratelimiter.Controlar(id, requisicaoPorSegundo, totalTempoBloqueado, totalSegundosExpiracaoToken)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, IPRequisicao, "\n", obterTokenRequest(r))
		fmt.Println(IPRequisicao, "\n", obterTokenRequest(r))

	})

	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		novoToken := uuid.New().String()
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, novoToken)
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
