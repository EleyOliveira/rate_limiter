package ratelimiter

import (
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
)

const requisicaoPorSegundo = 5
const totalSegundosBloqueado = 60
const totalSegundosExpiracaoToken = 60

func InicializarRateLimiter() RateLimiter {
	cache := &CacheRegistro{}
	return *NewRateLimiter(cache)
}

func TestGravarRegistro(t *testing.T) {

	ratelimiter := InicializarRateLimiter()
	ratelimiter.Controlar("175.890.789.123", requisicaoPorSegundo, totalSegundosBloqueado, totalSegundosExpiracaoToken)
	ratelimiter.Controlar("175.890.789.124", requisicaoPorSegundo, totalSegundosBloqueado, totalSegundosExpiracaoToken)
	ratelimiter.Controlar("175.890.789.125", requisicaoPorSegundo, totalSegundosBloqueado, totalSegundosExpiracaoToken)

	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.890.789.123").Id == "175.890.789.123")
	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.890.789.124").Id == "175.890.789.124")
	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.890.789.125").Id == "175.890.789.125")
	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.890.789.480") == nil)

}

func TestBloquearRegistroQuandoUltrapassarRequisicaoPorSegundo(t *testing.T) {

	ratelimiter := InicializarRateLimiter()
	ratelimiter.Controlar("175.890.789.123", requisicaoPorSegundo, totalSegundosBloqueado, totalSegundosExpiracaoToken)
	ratelimiter.Controlar("175.890.789.123", requisicaoPorSegundo, totalSegundosBloqueado, totalSegundosExpiracaoToken)
	ratelimiter.Controlar("175.890.789.123", requisicaoPorSegundo, totalSegundosBloqueado, totalSegundosExpiracaoToken)
	ratelimiter.Controlar("175.890.789.123", requisicaoPorSegundo, totalSegundosBloqueado, totalSegundosExpiracaoToken)
	ratelimiter.Controlar("175.890.789.123", requisicaoPorSegundo, totalSegundosBloqueado, totalSegundosExpiracaoToken)
	ratelimiter.Controlar("175.890.789.123", requisicaoPorSegundo, totalSegundosBloqueado, totalSegundosExpiracaoToken)

	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.890.789.123").TotalRequests == requisicaoPorSegundo)
	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.890.789.123").Bloqueado == true)

}

func TestNaoBloquearRegistroQuandoNaoUltrapassarRequisicaoPorSegundo(t *testing.T) {

	ratelimiter := InicializarRateLimiter()
	ratelimiter.Controlar("175.890.789.123", requisicaoPorSegundo, totalSegundosBloqueado, totalSegundosExpiracaoToken)
	ratelimiter.Controlar("175.890.789.123", requisicaoPorSegundo, totalSegundosBloqueado, totalSegundosExpiracaoToken)
	ratelimiter.Controlar("175.890.789.123", requisicaoPorSegundo, totalSegundosBloqueado, totalSegundosExpiracaoToken)
	ratelimiter.Controlar("175.890.789.123", requisicaoPorSegundo, totalSegundosBloqueado, totalSegundosExpiracaoToken)

	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.890.789.123").TotalRequests == 4)
	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.890.789.123").Bloqueado == false)
}

func TestRemoverRegistroAposTempoBloqueio(t *testing.T) {

	ratelimiter := InicializarRateLimiter()
	ratelimiter.Controlar("175.890.789.131", requisicaoPorSegundo, totalSegundosBloqueado, totalSegundosExpiracaoToken)

	registro := ratelimiter.controlaRateLimit.buscar("175.890.789.131")
	registro.FinalControle = registro.FinalControle.Add(time.Minute * time.Duration(totalSegundosBloqueado+1))
	ratelimiter.controlaRateLimit.remover()

	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.890.789.131") == nil)

}

func TestNaoRemoverRegistroAntesTempoBloqueio(t *testing.T) {

	ratelimiter := InicializarRateLimiter()
	ratelimiter.Controlar("175.890.789.132", requisicaoPorSegundo, totalSegundosBloqueado, totalSegundosExpiracaoToken)

	ratelimiter.controlaRateLimit.remover()

	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.890.789.132").Id == "175.890.789.132")

}