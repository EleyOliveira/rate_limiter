package ratelimiter

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/EleyOliveira/rate_limiter/internal/infra/database"
	"github.com/redis/go-redis/v9"
)

type CacheRegistro struct {
	Registros []*Registro
	Tokens    []*Token
}

func (i *CacheRegistro) gravar(registro Registro) error {

	data, err := json.Marshal(registro)
	if err != nil {
		return err
	}

	err = database.ObterRedisClienteIP().Set(context.Background(), registro.Id, data, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (i *CacheRegistro) buscar(id string) (*Registro, error) {
	val, err := database.ObterRedisClienteIP().Get(context.Background(), id).Result()

	if err == redis.Nil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	var registro Registro
	err = json.Unmarshal([]byte(val), &registro)
	if err != nil {
		return nil, err
	}

	return &registro, nil
}

func (i *CacheRegistro) remover() {
	var cursor uint64
	var keys []string

	cliente := database.ObterRedisClienteIP()

	for {
		keys, cursor, _ = cliente.Scan(context.Background(), cursor, "*", 10).Result()

		for _, key := range keys {
			val, _ := cliente.Get(context.Background(), key).Result()

			var registro Registro

			json.Unmarshal([]byte(val), &registro)

			if registro.FinalControle.Add(time.Second * time.Duration(registro.TempoBloqueado)).Before(time.Now()) {
				cliente.Del(context.Background(), registro.Id)
			}
		}

		if cursor == 0 {
			break
		}
	}
}

func (i *CacheRegistro) gravarToken(token Token) error {

	data, err := json.Marshal(token)
	if err != nil {
		return err
	}

	err = database.ObterRedisClienteToken().Set(context.Background(), token.Id, data, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (i *CacheRegistro) buscarToken(id string) (*Token, error) {

	val, err := database.ObterRedisClienteToken().Get(context.Background(), id).Result()

	if err == redis.Nil {
		return nil, errors.New("Token não encontrado")
	}

	if err != nil {
		return nil, err
	}

	var token Token
	err = json.Unmarshal([]byte(val), &token)
	if err != nil {
		return nil, err
	}

	if token.ExpiraEm.Before(time.Now()) {
		return nil, errors.New("Token expirado")
	}

	return &token, nil
}

func (i *CacheRegistro) removerToken() {
	var cursor uint64
	var keys []string

	cliente := database.ObterRedisClienteToken()

	for {
		keys, cursor, _ = cliente.Scan(context.Background(), cursor, "*", 10).Result()

		for _, key := range keys {
			val, _ := cliente.Get(context.Background(), key).Result()

			var token Token

			json.Unmarshal([]byte(val), &token)

			if token.ExpiraEm.Before(time.Now()) {
				cliente.Del(context.Background(), token.Id)
			}
		}

		if cursor == 0 {
			break
		}
	}
}
