package ratelimiter

import (
	"errors"
	"fmt"
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
	Id       string
	ExpiraEm time.Time
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

	tokenRequest := obterTokenRequest(request)

	mutex.Lock()
	defer mutex.Unlock()

	var token *Token

	if tokenRequest != "" {
		var err error
		token, err = r.controlaRateLimit.buscarToken(tokenRequest)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		requisicaoPorSegundo = requisicoesPorSegundoToken
		totalMinutosBloqueado = tempoBloqueioEmSegundosToken
	}

	ip, err := obterIPRequest(request)
	if err != nil {
		panic(err)
	}

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
	fmt.Println(registro.FinalControle)
	fmt.Println(time.Now().Truncate(time.Second).Add(time.Second * 1))
	fmt.Println(registro.Bloqueado)
	fmt.Println(registro.TotalRequests)

	if registro.FinalControle.Before(time.Now().Truncate(time.Second).Add(time.Second*1)) && !registro.Bloqueado {
		registro.FinalControle = time.Now().Add(time.Second * 1)
		registro.TotalRequests = 0
	}

	atualizarRegistro(registro, requisicaoPorSegundo)

	if registro.Bloqueado {
		if token != nil {
			if token.ExpiraEm.After(time.Now()) && registro.FinalControle.Add(time.Second*time.Duration(registro.TempoBloqueado)).Before(time.Now()) {
				registro.Bloqueado = false
				registro.TotalRequests = 1
				return http.StatusOK, nil
			}
		} else {
			r.controlaRateLimit.remover()
		}

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
	mutex.Lock()
	defer mutex.Unlock()

	token := Token{
		Id:       uuid.New().String(),
		ExpiraEm: time.Now().Add(time.Second * time.Duration(totalSegundosExpiracaoToken)),
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

func (r *RateLimiter) InicializarLimpezaRegistro(intervalo time.Duration) {

	go func() {
		ticker := time.NewTicker(intervalo)
		defer ticker.Stop()

		for {
			for range ticker.C {
				r.removerRegistro()
			}
		}
	}()
}

func (r *RateLimiter) InicializarLimpezaToken(intervalo time.Duration) {

	go func() {
		ticker := time.NewTicker(intervalo)
		defer ticker.Stop()

		for {
			for range ticker.C {
				r.removerToken()
			}
		}
	}()
}

func (r *RateLimiter) removerRegistro() {
	mutex.Lock()
	defer mutex.Unlock()
	r.controlaRateLimit.remover()
}

func (r *RateLimiter) removerToken() {
	mutex.Lock()
	defer mutex.Unlock()
	r.controlaRateLimit.removerToken()
}
