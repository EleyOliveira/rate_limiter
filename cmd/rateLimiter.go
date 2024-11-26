package main

type RateLimiter struct {
	controlaRateLimit ControlaRateLimit
}

func NewRateLimiter(controlaRateLimit ControlaRateLimit) *RateLimiter {
	return &RateLimiter{
		controlaRateLimit: controlaRateLimit,
	}
}

func (r *RateLimiter) Controlar(registro string) {
	r.controlaRateLimit.gravar(registro)
}
