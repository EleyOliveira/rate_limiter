package ratelimiter

import (
	"time"
)

type CacheRegistro struct {
	Registros []*Registro
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
		if item.FinalControle.Before(time.Now().Add(time.Minute * time.Duration(item.TempoBloqueado))) {
			registros = append(registros, item)
		}
	}

	i.Registros = registros

}
