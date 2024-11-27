package ratelimiter

type CacheRegistro struct {
	Registros []string
}

func (i *CacheRegistro) gravar(registro string) error {

	if len(i.Registros) == 0 || !i.contem(registro) {
		i.Registros = append(i.Registros, registro)
		return nil
	}

	return nil
}

func (i *CacheRegistro) contem(registro string) bool {
	for _, ip := range i.Registros {
		if ip == registro {
			return true
		}
	}
	return false
}
