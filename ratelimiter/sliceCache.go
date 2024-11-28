package ratelimiter

type CacheRegistro struct {
	Registros []Registro
}

func (i *CacheRegistro) gravar(registro Registro) error {

	i.Registros = append(i.Registros, registro)
	return nil

}

func (i *CacheRegistro) buscar(id string) Registro {
	for _, item := range i.Registros {
		if item.Id == id {
			return item
		}
	}
	return Registro{}
}
