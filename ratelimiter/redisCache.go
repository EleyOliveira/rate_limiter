package ratelimiter

import (
	"context"
	"encoding/json"
	"errors"
	"time"

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

	err = conectarRedis().Set(context.Background(), registro.Id, data, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (i *CacheRegistro) buscar(id string) (*Registro, error) {
	val, err := conectarRedis().Get(context.Background(), id).Result()

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
}

func (i *CacheRegistro) gravarToken(token Token) error {

	data, err := json.Marshal(token)
	if err != nil {
		return err
	}

	err = conectarRedis().Set(context.Background(), token.Id, data, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (i *CacheRegistro) buscarToken(id string) (*Token, error) {

	val, err := conectarRedis().Get(context.Background(), id).Result()

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
}

func conectarRedis() *redis.Client {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return rdb
}
