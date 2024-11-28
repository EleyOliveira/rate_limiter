package ratelimiter

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

const requisicaoPorSegundo = 5

func InicializarRateLimiter() RateLimiter {
	cache := &CacheRegistro{}
	return *NewRateLimiter(cache)
}

func TestGravarRegistro(t *testing.T) {

	ratelimiter := InicializarRateLimiter()
	ratelimiter.Controlar("175.890.789.123", requisicaoPorSegundo)
	ratelimiter.Controlar("175.890.789.124", requisicaoPorSegundo)
	ratelimiter.Controlar("175.890.789.125", requisicaoPorSegundo)

	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.890.789.123").Id == "175.890.789.123")
	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.890.789.124").Id == "175.890.789.124")
	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.890.789.125").Id == "175.890.789.125")
	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.890.789.480") == nil)

}

func TestBloquearRegistroQuandoUltrapassarRequisicaoPorSegundo(t *testing.T) {

	ratelimiter := InicializarRateLimiter()
	ratelimiter.Controlar("175.890.789.123", requisicaoPorSegundo)
	ratelimiter.Controlar("175.890.789.123", requisicaoPorSegundo)
	ratelimiter.Controlar("175.890.789.123", requisicaoPorSegundo)
	ratelimiter.Controlar("175.890.789.123", requisicaoPorSegundo)
	ratelimiter.Controlar("175.890.789.123", requisicaoPorSegundo)
	ratelimiter.Controlar("175.890.789.123", requisicaoPorSegundo)

	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.890.789.123").TotalRequests == requisicaoPorSegundo)
	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.890.789.123").Bloqueado == true)

}

func TestNaoBloquearRegistroQuandoNaoUltrapassarRequisicaoPorSegundo(t *testing.T) {

	ratelimiter := InicializarRateLimiter()
	ratelimiter.Controlar("175.890.789.123", requisicaoPorSegundo)
	ratelimiter.Controlar("175.890.789.123", requisicaoPorSegundo)
	ratelimiter.Controlar("175.890.789.123", requisicaoPorSegundo)
	ratelimiter.Controlar("175.890.789.123", requisicaoPorSegundo)

	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.890.789.123").TotalRequests == 4)
	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.890.789.123").Bloqueado == false)
}

func TestRemoverRegistroAposTempoBloqueio(t *testing.T) {

	ratelimiter := InicializarRateLimiter()
	ratelimiter.Controlar("175.890.789.131", requisicaoPorSegundo)

	ratelimiter.controlaRateLimit.remover()

	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.890.789.123") == nil)

}
