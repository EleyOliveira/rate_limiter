package ratelimiter

import (
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
	"github.com/tkuchiki/faketime"
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

	timeFake := faketime.NewFaketime(2000, time.February, 10, 9, 0, 0, 0, time.UTC)
	defer timeFake.Undo()

	timeFake.Do()

	ratelimiter.controlaRateLimit.remover()
	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.890.789.131") == nil)

}

func TestNaoRemoverRegistroAntesTempoBloqueio(t *testing.T) {

	ratelimiter := InicializarRateLimiter()
	ratelimiter.Controlar("175.890.789.132", requisicaoPorSegundo, totalSegundosBloqueado, totalSegundosExpiracaoToken)

	ratelimiter.controlaRateLimit.remover()

	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.890.789.132").Id == "175.890.789.132")

}

func TestGravarToken(t *testing.T) {

	ratelimiter := InicializarRateLimiter()

	token, err := ratelimiter.GerarToken(totalSegundosExpiracaoToken)
	assert.Equal(t, true, err == nil)

	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscarToken(token).Id == token)

}
