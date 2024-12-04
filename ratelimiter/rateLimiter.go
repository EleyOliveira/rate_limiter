package ratelimiter

import (
	"errors"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

type RateLimiter struct {
	controlaRateLimit ControlaCache
}

type Registro struct {
	Id             string
	FinalControle  time.Time
	Bloqueado      bool
	TotalRequests  int
	TempoBloqueado int
}

type Token struct {
	Id        string
	ExpiraEm  time.Time
	Utilizado bool
}

var mutex sync.Mutex

func NewRateLimiter(controlaRateLimit ControlaCache) *RateLimiter {
	return &RateLimiter{
		controlaRateLimit: controlaRateLimit,
	}
}

func (r *RateLimiter) Controlar(request *http.Request, requisicoesPorSegundoIP int,
	requisicoesPorSegundoToken int, tempoBloqueioEmSegundosIP int, tempoBloqueioEmSegundosToken int,
	tempoEmSegundosExpiracaoToken int) (int, error) {

	requisicaoPorSegundo := requisicoesPorSegundoIP
	totalMinutosBloqueado := tempoBloqueioEmSegundosIP

	if obterTokenRequest(request) != "" {
		requisicaoPorSegundo = requisicoesPorSegundoToken
		totalMinutosBloqueado = tempoBloqueioEmSegundosToken
	}

	ip, err := obterIPRequest(request)
	if err != nil {
		panic(err)
	}

	mutex.Lock()
	defer mutex.Unlock()
	registro := r.controlaRateLimit.buscar(ip)

	if registro == nil {
		novoRegistro := Registro{
			Id:             ip,
			FinalControle:  time.Now().Add(time.Second * 1),
			Bloqueado:      false,
			TotalRequests:  1,
			TempoBloqueado: totalMinutosBloqueado,
		}

		r.controlaRateLimit.gravar(novoRegistro)

		return http.StatusOK, nil
	}

	atualizarRegistro(registro, requisicaoPorSegundo)

	if registro.Bloqueado {
		return http.StatusTooManyRequests, errors.New("you have reached the maximum number of requests or actions allowed within a certain time frame")
	}

	return http.StatusOK, nil
}

func atualizarRegistro(registro *Registro, requisicaoPorSegundo int) {

	if registro.FinalControle.After(time.Now()) {
		if registro.TotalRequests < requisicaoPorSegundo {
			registro.TotalRequests++
		} else {
			registro.Bloqueado = true
		}
	}
}

func (r *RateLimiter) GerarToken(totalSegundosExpiracaoToken int) (string, error) {
	token := Token{
		Id:        uuid.New().String(),
		ExpiraEm:  time.Now().Add(time.Second * time.Duration(totalSegundosExpiracaoToken)),
		Utilizado: false,
	}

	if err := r.controlaRateLimit.gravarToken(token); err != nil {
		return "", err
	}

	return token.Id, nil
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
