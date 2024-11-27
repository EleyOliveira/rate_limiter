package ratelimiter

import "time"

type RateLimiter struct {
	controlaRateLimit ControlaCache
}

type Registro struct {
	Id            string
	FinalControle time.Time
	Bloqueado     bool
}

var registros []Registro

func NewRateLimiter(controlaRateLimit ControlaCache) *RateLimiter {
	return &RateLimiter{
		controlaRateLimit: controlaRateLimit,
	}
}

func (r *RateLimiter) Controlar(registro string) {
	r.controlaRateLimit.gravar(registro)
}

func GravarRegistro(id string) {

	registro := Registro{
		Id:            id,
		FinalControle: time.Now().Add(time.Second * 1),
		Bloqueado:     false,
	}

	registros = append(registros, registro)
}

func ExisteRegistro(id string) bool {
	for _, reg := range registros {
		if reg.Id == id {
			return true
		}
	}
	return false
}
