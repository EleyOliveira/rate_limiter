package ratelimiter

type CacheSlice struct {
	Registros []*Registro
	Tokens    []*Token
}

func (i *CacheSlice) gravar(registro Registro) error {

	i.Registros = append(i.Registros, &registro)
	return nil

}

func (i *CacheSlice) buscar(id string) (*Registro, error) {
	for _, item := range i.Registros {
		if item.Id == id {
			return item, nil
		}
	}
	return nil, nil
}

func (i *CacheSlice) buscarTodos() ([]*Registro, error) {
	return i.Registros, nil
}

func (i *CacheSlice) remover(ids []string) {
	var registros []*Registro
	for _, item := range i.Registros {
		contains := func(s []string, str string) bool {
			for _, v := range s {
				if v == str {
					return true
				}
			}
			return false
		}
		if !contains(ids, item.Id) {
			registros = append(registros, item)
		}
	}
	i.Registros = registros
}

func (i *CacheSlice) gravarToken(token Token) error {

	i.Tokens = append(i.Tokens, &token)
	return nil

}

func (i *CacheSlice) buscarToken(id string) (*Token, error) {
	for _, item := range i.Tokens {
		if item.Id == id {
			return item, nil
		}
	}
	return nil, nil
}

func (i *CacheSlice) buscarTokenTodos() ([]*Token, error) {
	return i.Tokens, nil
}

func (i *CacheSlice) removerToken(ids []string) {
	var tokens []*Token
	for _, item := range i.Tokens {
		contains := func(s []string, str string) bool {
			for _, v := range s {
				if v == str {
					return true
				}
			}
			return false
		}
		if !contains(ids, item.Id) {
			tokens = append(tokens, item)
		}
	}
	i.Tokens = tokens
}
