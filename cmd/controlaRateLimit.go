package main

type ControlaRateLimit interface {
	gravar(registro string) error
	contem(registro string) bool
}
