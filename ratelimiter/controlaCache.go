package ratelimiter

type ControlaCache interface {
	gravar(registro Registro) error
	gravarToken(token Token) error
	buscar(id string) *Registro
	buscarToken(id string) (*Token, error)
	remover()
}
