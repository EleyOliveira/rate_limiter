package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/EleyOliveira/rate_limiter/ratelimiter"
	"github.com/joho/godotenv"
)

type Configuracao struct {
	requisicoesPorSegundoIP       int
	requisicoesPorSegundoToken    int
	tempoBloqueioEmSegundosIP     int
	tempoBloqueioEmSegundosToken  int
	tempoEmSegundosExpiracaoToken int
}

func main() {

	configuracao := carregarConfiguracao()

	registro := &ratelimiter.CacheRegistro{}
	ratelimiter := ratelimiter.NewRateLimiter(registro)
	ratelimiter.InicializarLimpezaRegistro(1 * time.Minute)
	ratelimiter.InicializarLimpezaToken(1 * time.Minute)

	http.Handle("/", rateLimiterMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintln(w, "Sucesso")

	}), ratelimiter, configuracao))

	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {

		token, err := ratelimiter.GerarToken(configuracao.tempoEmSegundosExpiracaoToken)
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

func carregarConfiguracao() Configuracao {

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	requisicoesPorSegundoIP, err := validarConfiguracao("REQUISICOES_POR_SEGUNDO_IP")
	if err != nil {
		panic(err)
	}

	requisicoesPorSegundoToken, err := validarConfiguracao("REQUISICOES_POR_SEGUNDO_TOKEN")
	if err != nil {
		panic(err)
	}

	tempoBloqueioEmSegundosIP, err := validarConfiguracao("TEMPO_BLOQUEIO_EM_SEGUNDO_IP")
	if err != nil {
		panic(err)
	}

	tempoBloqueioEmSegundosToken, err := validarConfiguracao("TEMPO_BLOQUEIO_EM_SEGUNDO_TOKEN")
	if err != nil {
		panic(err)
	}

	tempoEmSegundosExpiracaoToken, err := validarConfiguracao("TEMPO_EM_SEGUNDOS_EXPIRACAO_TOKEN")
	if err != nil {
		panic(err)
	}

	return Configuracao{
		requisicoesPorSegundoIP,
		requisicoesPorSegundoToken,
		tempoBloqueioEmSegundosIP,
		tempoBloqueioEmSegundosToken,
		tempoEmSegundosExpiracaoToken,
	}

}

func validarConfiguracao(configuracao string) (int, error) {
	valorConfiguracao := os.Getenv(configuracao)

	if valorConfiguracao == "" {
		return 0, fmt.Errorf("configuração %s não encontrada", configuracao)
	}

	valor, err := strconv.Atoi(valorConfiguracao)
	if err != nil {
		return 0, fmt.Errorf("o valor da configuração %s inválida", configuracao)
	}

	return valor, nil

}

func rateLimiterMiddleware(next http.Handler, rl *ratelimiter.RateLimiter, config Configuracao) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		statusCode, err := rl.Controlar(r, config.requisicoesPorSegundoIP,
			config.requisicoesPorSegundoToken, config.tempoBloqueioEmSegundosIP,
			config.tempoBloqueioEmSegundosToken, config.tempoEmSegundosExpiracaoToken)

		if err != nil {
			w.WriteHeader(statusCode)
			fmt.Fprintln(w, err.Error())
			return
		}

		next.ServeHTTP(w, r)
	})
}
