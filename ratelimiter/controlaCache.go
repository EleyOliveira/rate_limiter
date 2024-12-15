package ratelimiter

type ControlaCache interface {
	gravar(registro Registro) error
	buscar(id string) (*Registro, error)
	buscarTodos() ([]*Registro, error)
	remover(ids []string)
	gravarToken(token Token) error
	buscarToken(id string) (*Token, error)
	buscarTokenTodos() ([]*Token, error)
	removerToken(ids []string)
}
