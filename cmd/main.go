package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		registro := &ratelimiter.CacheRegistro{}

		ratelimiter := ratelimiter.NewRateLimiter(registro)
		statusCode, err := ratelimiter.Controlar(r, configuracao.requisicoesPorSegundoIP,
			configuracao.requisicoesPorSegundoToken, configuracao.tempoBloqueioEmSegundosIP,
			configuracao.tempoBloqueioEmSegundosToken, configuracao.tempoEmSegundosExpiracaoToken)
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
