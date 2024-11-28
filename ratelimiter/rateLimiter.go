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

func (r *RateLimiter) Controlar(id string) {

	registro := r.controlaRateLimit.buscar(id)

	if registro.Bloqueado {
		return
	}

	if registro.Id == "" {
		novoRegistro := Registro{
			Id:            id,
			FinalControle: time.Now().Add(time.Second * 1),
			Bloqueado:     false,
			TotalRequests: 0,
		}
		r.controlaRateLimit.gravar(novoRegistro)
		return
	}

	if registro.FinalControle.Before(time.Now()) {
		if registro.TotalRequests < 5 {
			registro.TotalRequests++
		} else {
			registro.Bloqueado = true
		}
	}
}
