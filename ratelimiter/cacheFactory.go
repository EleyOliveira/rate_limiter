package ratelimiter

import (
	"errors"
)

func ObterCache(cacheTipo string) (ControlaCache, error) {
	if cacheTipo == "slice" {
		return &CacheSlice{}, nil
	}

	if cacheTipo == "redis" {
		return &CacheRedis{}, nil
	}

	return nil, errors.New("tipo de cache não encontrado")
}
