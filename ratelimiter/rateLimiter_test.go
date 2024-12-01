package ratelimiter

import (
	"net/http"
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
	"github.com/tkuchiki/faketime"
)

const requisicaoPorSegundo = 5
const totalSegundosBloqueado = 60
const totalSegundosExpiracaoToken = 60

func inicializarRateLimiter() RateLimiter {
	cache := &CacheRegistro{}
	return *NewRateLimiter(cache)
}

func inicializarRequestTeste(ip string, porta string) http.Request {
	req, err := http.NewRequest("GET", "http://teste.com", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.RemoteAddr = "175.890.789.132:1234"
	return *req
}

func TestGravarRegistro(t *testing.T) {

	ratelimiter := inicializarRateLimiter()
	//ratelimiter.Controlar("175.890.789.123", requisicaoPorSegundo, totalSegundosBloqueado, totalSegundosExpiracaoToken)
	//ratelimiter.Controlar("175.890.789.124", requisicaoPorSegundo, totalSegundosBloqueado, totalSegundosExpiracaoToken)
	//ratelimiter.Controlar("175.890.789.125", requisicaoPorSegundo, totalSegundosBloqueado, totalSegundosExpiracaoToken)

	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.890.789.123").Id == "175.890.789.123")
	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.890.789.124").Id == "175.890.789.124")
	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.890.789.125").Id == "175.890.789.125")
	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.890.789.480") == nil)

}

func TestBloquearRegistroQuandoUltrapassarRequisicaoPorSegundo(t *testing.T) {

	ratelimiter := inicializarRateLimiter()
	//ratelimiter.Controlar("175.890.789.123", requisicaoPorSegundo, totalSegundosBloqueado, totalSegundosExpiracaoToken)
	//ratelimiter.Controlar("175.890.789.123", requisicaoPorSegundo, totalSegundosBloqueado, totalSegundosExpiracaoToken)
	//ratelimiter.Controlar("175.890.789.123", requisicaoPorSegundo, totalSegundosBloqueado, totalSegundosExpiracaoToken)
	//ratelimiter.Controlar("175.890.789.123", requisicaoPorSegundo, totalSegundosBloqueado, totalSegundosExpiracaoToken)
	//ratelimiter.Controlar("175.890.789.123", requisicaoPorSegundo, totalSegundosBloqueado, totalSegundosExpiracaoToken)
	//ratelimiter.Controlar("175.890.789.123", requisicaoPorSegundo, totalSegundosBloqueado, totalSegundosExpiracaoToken)

	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.890.789.123").TotalRequests == requisicaoPorSegundo)
	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.890.789.123").Bloqueado == true)

}

func TestNaoBloquearRegistroQuandoNaoUltrapassarRequisicaoPorSegundo(t *testing.T) {

	ratelimiter := inicializarRateLimiter()
	//ratelimiter.Controlar("175.890.789.123", requisicaoPorSegundo, totalSegundosBloqueado, totalSegundosExpiracaoToken)
	//ratelimiter.Controlar("175.890.789.123", requisicaoPorSegundo, totalSegundosBloqueado, totalSegundosExpiracaoToken)
	//ratelimiter.Controlar("175.890.789.123", requisicaoPorSegundo, totalSegundosBloqueado, totalSegundosExpiracaoToken)
	//ratelimiter.Controlar("175.890.789.123", requisicaoPorSegundo, totalSegundosBloqueado, totalSegundosExpiracaoToken)

	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.890.789.123").TotalRequests == 4)
	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.890.789.123").Bloqueado == false)
}

func TestRemoverRegistroAposTempoBloqueio(t *testing.T) {

	ratelimiter := inicializarRateLimiter()
	//ratelimiter.Controlar("175.890.789.131", requisicaoPorSegundo, totalSegundosBloqueado, totalSegundosExpiracaoToken)

	timeFake := faketime.NewFaketime(2000, time.February, 10, 9, 0, 0, 0, time.UTC)
	defer timeFake.Undo()

	timeFake.Do()

	ratelimiter.controlaRateLimit.remover()
	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.890.789.131") == nil)

}

func TestNaoRemoverRegistroAntesTempoBloqueio(t *testing.T) {

	ratelimiter := inicializarRateLimiter()
	req, err := http.NewRequest("GET", "http://teste.com", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.RemoteAddr = "175.890.789.132:1234"
	ratelimiter.Controlar(req, requisicaoPorSegundo, totalSegundosBloqueado, totalSegundosExpiracaoToken)

	ratelimiter.controlaRateLimit.remover()

	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.890.789.132").Id == "175.890.789.132")

}

func TestGravarToken(t *testing.T) {

	ratelimiter := inicializarRateLimiter()

	token, err := ratelimiter.GerarToken(totalSegundosExpiracaoToken)
	assert.Equal(t, true, err == nil)

	tokenGravado, err := ratelimiter.controlaRateLimit.buscarToken(token)

	assert.Equal(t, true, err == nil)
	assert.Equal(t, true, tokenGravado.Utilizado)
	assert.Equal(t, true, tokenGravado.Id == token)

}

func TestRetornarMensagemParaTokenExpirado(t *testing.T) {

	ratelimiter := inicializarRateLimiter()

	token, err := ratelimiter.GerarToken(totalSegundosExpiracaoToken)
	assert.Equal(t, true, err == nil)

	timeFake := faketime.NewFaketime(2500, time.February, 10, 9, 0, 0, 0, time.UTC)
	defer timeFake.Undo()

	timeFake.Do()

	tokenGravado, err := ratelimiter.controlaRateLimit.buscarToken(token)

	assert.Equal(t, true, tokenGravado == nil)
	assert.Equal(t, true, err.Error() == "Token expirado")
}

func TestRetornarMensagemParaTokenJaUtilizado(t *testing.T) {

	ratelimiter := inicializarRateLimiter()

	token, err := ratelimiter.GerarToken(totalSegundosExpiracaoToken)
	assert.Equal(t, true, err == nil)

	tokenGravado, err := ratelimiter.controlaRateLimit.buscarToken(token)

	assert.Equal(t, true, tokenGravado.Utilizado)
	assert.Equal(t, true, err == nil)

	tokenGravado, err = ratelimiter.controlaRateLimit.buscarToken(token)
	assert.Equal(t, true, tokenGravado == nil)
	assert.Equal(t, true, err.Error() == "Token já utilizado")

}

func TestRemoverToke(t *testing.T) {

	ratelimiter := inicializarRateLimiter()

	token, err := ratelimiter.GerarToken(totalSegundosExpiracaoToken)
	assert.Equal(t, true, err == nil)

	tokenGravado, err := ratelimiter.controlaRateLimit.buscarToken(token)

	assert.Equal(t, true, err == nil)
	assert.Equal(t, true, tokenGravado.Id == token)

	timeFake := faketime.NewFaketime(2500, time.February, 10, 9, 0, 0, 0, time.UTC)
	defer timeFake.Undo()

	timeFake.Do()

	ratelimiter.controlaRateLimit.removerToken()
	tokenRemovido, err := ratelimiter.controlaRateLimit.buscarToken(token)

	assert.Equal(t, true, tokenRemovido == nil)
	assert.Equal(t, true, err.Error() == "Token não encontrado")

}
