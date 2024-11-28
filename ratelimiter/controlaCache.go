package ratelimiter

type ControlaCache interface {
	gravar(registro Registro) error
	buscar(id string) Registro
}
