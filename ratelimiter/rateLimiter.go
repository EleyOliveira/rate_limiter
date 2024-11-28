package ratelimiter

import "time"

type RateLimiter struct {
	controlaRateLimit ControlaCache
}

type Registro struct {
	Id            string
	FinalControle time.Time
	Bloqueado     bool
	TotalRequests int
}

func NewRateLimiter(controlaRateLimit ControlaCache) *RateLimiter {
	return &RateLimiter{
		controlaRateLimit: controlaRateLimit,
	}
}

func (r *RateLimiter) Controlar(id string, requisicaoPorSegundo int) {

	registro := r.controlaRateLimit.buscar(id)

	if registro == nil {
		novoRegistro := Registro{
			Id:            id,
			FinalControle: time.Now().Add(time.Second * 1),
			Bloqueado:     false,
			TotalRequests: 1,
		}
		r.controlaRateLimit.gravar(novoRegistro)
		return
	}

	if registro.Bloqueado {
		return
	}

	if registro.FinalControle.After(time.Now()) {
		if registro.TotalRequests < requisicaoPorSegundo {
			registro.TotalRequests++
		} else {
			registro.Bloqueado = true
		}
	}
}
