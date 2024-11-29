package ratelimiter

import (
	"time"
)

type CacheRegistro struct {
	Registros []*Registro
	Tokens    []*Token
}

func (i *CacheRegistro) gravar(registro Registro) error {

	i.Registros = append(i.Registros, &registro)
	return nil

}

func (i *CacheRegistro) buscar(id string) *Registro {
	for _, item := range i.Registros {
		if item.Id == id {
			return item
		}
	}
	return nil
}

func (i *CacheRegistro) remover() {
	var registros []*Registro
	for _, item := range i.Registros {
		if item.FinalControle.Before(time.Now().Add(time.Second * time.Duration(item.TempoBloqueado))) {
			registros = append(registros, item)
		}
	}
	i.Registros = registros
}

func (i *CacheRegistro) gravarToken(token Token) error {

	i.Tokens = append(i.Tokens, &token)
	return nil

}

func (i *CacheRegistro) buscarToken(id string) *Token {
	for _, item := range i.Tokens {
		if item.Id == id {
			return item
		}
	}
	return nil
}
