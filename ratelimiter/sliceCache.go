package ratelimiter

import (
	"errors"
	"time"
)

type CacheRegistroSlice struct {
	Registros []*Registro
	Tokens    []*Token
}

func (i *CacheRegistroSlice) gravar(registro Registro) error {

	i.Registros = append(i.Registros, &registro)
	return nil

}

func (i *CacheRegistroSlice) buscar(id string) (*Registro, error) {
	for _, item := range i.Registros {
		if item.Id == id {
			return item, nil
		}
	}
	return nil, nil
}

func (i *CacheRegistroSlice) remover() {
	var registros []*Registro
	for _, item := range i.Registros {
		if item.FinalControle.Add(time.Second * time.Duration(item.TempoBloqueado)).After(time.Now()) {
			registros = append(registros, item)
		}
	}
	i.Registros = registros
}

func (i *CacheRegistroSlice) gravarToken(token Token) error {

	i.Tokens = append(i.Tokens, &token)
	return nil

}

func (i *CacheRegistroSlice) buscarToken(id string) (*Token, error) {
	for _, item := range i.Tokens {
		if item.Id == id {
			if item.ExpiraEm.Before(time.Now()) {
				return nil, errors.New("Token expirado")
			}

			return item, nil
		}
	}
	return nil, errors.New("Token não encontrado")
}

func (i *CacheRegistroSlice) removerToken() {
	var tokens []*Token
	for _, item := range i.Tokens {
		if item.ExpiraEm.After(time.Now()) {
			tokens = append(tokens, item)
		}
	}
	i.Tokens = tokens
}
