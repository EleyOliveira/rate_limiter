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

	err = conectarRedisBancoIP().Set(context.Background(), registro.Id, data, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (i *CacheRegistro) buscar(id string) (*Registro, error) {
	val, err := conectarRedisBancoIP().Get(context.Background(), id).Result()

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

	for {
		keys, cursor, _ = conectarRedisBancoIP().Scan(context.Background(), cursor, "*", 10).Result()

		for _, key := range keys {
			val, _ := conectarRedisBancoIP().Get(context.Background(), key).Result()

			var registro Registro

			json.Unmarshal([]byte(val), &registro)

			if registro.FinalControle.Add(time.Second * time.Duration(registro.TempoBloqueado)).Before(time.Now()) {
				conectarRedisBancoIP().Del(context.Background(), registro.Id)
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

	err = conectarRedisBancoToken().Set(context.Background(), token.Id, data, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (i *CacheRegistro) buscarToken(id string) (*Token, error) {

	val, err := conectarRedisBancoToken().Get(context.Background(), id).Result()

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

	for {
		keys, cursor, _ = conectarRedisBancoToken().Scan(context.Background(), cursor, "*", 10).Result()

		for _, key := range keys {
			val, _ := conectarRedisBancoToken().Get(context.Background(), key).Result()

			var token Token

			json.Unmarshal([]byte(val), &token)

			if token.ExpiraEm.Before(time.Now()) {
				conectarRedisBancoToken().Del(context.Background(), token.Id)
			}
		}

		if cursor == 0 {
			break
		}
	}
}

func conectarRedisBancoIP() *redis.Client {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return rdb
}

func conectarRedisBancoToken() *redis.Client {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	})

	return rdb
}
