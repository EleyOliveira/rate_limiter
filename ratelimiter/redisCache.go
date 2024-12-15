package ratelimiter

import (
	"context"
	"encoding/json"

	"github.com/EleyOliveira/rate_limiter/internal/infra/database"
	"github.com/redis/go-redis/v9"
)

type CacheRedis struct {
	Registros []*Registro
	Tokens    []*Token
}

func (i *CacheRedis) gravar(registro Registro) error {

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

func (i *CacheRedis) buscar(id string) (*Registro, error) {
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

func (i *CacheRedis) buscarTodos() ([]*Registro, error) {

	var cursor uint64
	var keys []string

	cliente := database.ObterRedisClienteIP()

	for {
		keys, cursor, _ = cliente.Scan(context.Background(), cursor, "*", 10).Result()

		for _, key := range keys {
			val, _ := cliente.Get(context.Background(), key).Result()

			var registro Registro

			json.Unmarshal([]byte(val), &registro)

			i.Registros = append(i.Registros, &registro)
		}

		if cursor == 0 {
			break
		}
	}

	return i.Registros, nil
}

func (i *CacheRedis) remover(ids []string) {

	cliente := database.ObterRedisClienteIP()

	for _, id := range ids {
		cliente.Del(context.Background(), id)
	}
}

func (i *CacheRedis) gravarToken(token Token) error {

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

func (i *CacheRedis) buscarToken(id string) (*Token, error) {

	val, err := database.ObterRedisClienteToken().Get(context.Background(), id).Result()

	if err == redis.Nil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var token Token
	err = json.Unmarshal([]byte(val), &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (i *CacheRedis) buscarTokenTodos() ([]*Token, error) {
	var cursor uint64
	var keys []string

	cliente := database.ObterRedisClienteToken()

	for {
		keys, cursor, _ = cliente.Scan(context.Background(), cursor, "*", 10).Result()

		for _, key := range keys {
			val, _ := cliente.Get(context.Background(), key).Result()

			var token Token

			json.Unmarshal([]byte(val), &token)

			i.Tokens = append(i.Tokens, &token)
		}

		if cursor == 0 {
			break
		}
	}

	return i.Tokens, nil
}

func (i *CacheRedis) removerToken(ids []string) {

	cliente := database.ObterRedisClienteToken()

	for _, id := range ids {
		cliente.Del(context.Background(), id)
	}
}
