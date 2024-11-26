package main

type ipRequest struct {
	IPRequests []string
}

func (i *ipRequest) gravar(registro string) error {

	if len(i.IPRequests) == 0 || !i.contem(registro) {
		i.IPRequests = append(i.IPRequests, registro)
		return nil
	}

	return nil
}

func (i *ipRequest) contem(registro string) bool {
	for _, ip := range i.IPRequests {
		if ip == registro {
			return true
		}
	}
	return false
}
