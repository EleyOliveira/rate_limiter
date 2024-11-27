package ratelimiter

type ControlaCache interface {
	gravar(registro string) error
	contem(registro string) bool
}
