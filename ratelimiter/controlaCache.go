package ratelimiter

type ControlaCache interface {
	gravar(registro Registro) error
	buscar(id string) (*Registro, error)
	remover()
	gravarToken(token Token) error
	buscarToken(id string) (*Token, error)
	removerToken()
}
