package ratelimiter

import (
	"net/http"
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
	"github.com/tkuchiki/faketime"
)

const requisicaoPorSegundoIP = 5
const requisicaoPorSegundoToken = 10
const tempoSegundosBloqueadoIP = 120
const tempoSegundosBloqueadoToken = 60
const tempoSegundosExpiracaoToken = 60

func inicializarRateLimiter() RateLimiter {
	cache := &CacheRegistro{}
	return *NewRateLimiter(cache)
}

func inicializarRequestTeste(ip string, porta string) http.Request {
	req, _ := http.NewRequest("", "", nil)
	req.RemoteAddr = ip + ":" + porta
	return *req
}

func inicializarRequestTesteComToken(ip string, porta string, token string) http.Request {
	req, _ := http.NewRequest("", "", nil)
	req.Header.Add("API_KEY", token)
	req.RemoteAddr = ip + ":" + porta
	return *req
}

func TestGravarRegistro(t *testing.T) {

	ratelimiter := inicializarRateLimiter()
	req := inicializarRequestTeste("175.890.789.123", "1234")
	ratelimiter.Controlar(&req, requisicaoPorSegundoIP, requisicaoPorSegundoToken, tempoSegundosBloqueadoIP, tempoSegundosBloqueadoToken, tempoSegundosExpiracaoToken)
	req = inicializarRequestTeste("175.890.789.124", "1234")
	ratelimiter.Controlar(&req, requisicaoPorSegundoIP, requisicaoPorSegundoToken, tempoSegundosBloqueadoIP, tempoSegundosBloqueadoToken, tempoSegundosExpiracaoToken)
	req = inicializarRequestTeste("175.890.789.125", "1234")
	ratelimiter.Controlar(&req, requisicaoPorSegundoIP, requisicaoPorSegundoToken, tempoSegundosBloqueadoIP, tempoSegundosBloqueadoToken, tempoSegundosExpiracaoToken)

	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.890.789.123").Id == "175.890.789.123")
	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.890.789.124").Id == "175.890.789.124")
	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.890.789.125").Id == "175.890.789.125")
	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.890.789.480") == nil)

}

func TestBloquearRegistroQuandoUltrapassarRequisicaoPorSegundo(t *testing.T) {

	const totalRequisicao int = 6

	ratelimiter := inicializarRateLimiter()
	req := inicializarRequestTeste("175.890.789.123", "1234")

	for i := 0; i < totalRequisicao; i++ {
		ratelimiter.Controlar(&req, requisicaoPorSegundoIP, requisicaoPorSegundoToken, tempoSegundosBloqueadoIP, tempoSegundosBloqueadoToken, tempoSegundosExpiracaoToken)
	}

	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.890.789.123").TotalRequests == requisicaoPorSegundoIP)
	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.890.789.123").Bloqueado == true)

}

func TestNaoBloquearRegistroQuandoNaoUltrapassarRequisicaoPorSegundo(t *testing.T) {

	const totalRequisicao int = 4

	ratelimiter := inicializarRateLimiter()
	req := inicializarRequestTeste("175.890.789.123", "1234")

	for i := 0; i < totalRequisicao; i++ {
		ratelimiter.Controlar(&req, requisicaoPorSegundoIP, requisicaoPorSegundoToken, tempoSegundosBloqueadoIP, tempoSegundosBloqueadoToken, tempoSegundosExpiracaoToken)
	}

	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.890.789.123").TotalRequests == requisicaoPorSegundoIP-1)
	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.890.789.123").Bloqueado == false)
}

func TestRemoverRegistroAposTempoBloqueio(t *testing.T) {

	ratelimiter := inicializarRateLimiter()
	req := inicializarRequestTeste("175.890.789.131", "1234")

	ratelimiter.Controlar(&req, requisicaoPorSegundoIP, requisicaoPorSegundoToken, tempoSegundosBloqueadoIP, tempoSegundosBloqueadoToken, tempoSegundosExpiracaoToken)

	registroIncluido := ratelimiter.controlaRateLimit.buscar("175.890.789.131")
	assert.Equal(t, true, registroIncluido.Id == "175.890.789.131")

	timeFake := faketime.NewFaketime(2000, time.February, 10, 9, 0, 0, 0, time.UTC)
	defer timeFake.Undo()

	timeFake.Do()

	ratelimiter.controlaRateLimit.remover()
	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar(registroIncluido.Id) == nil)

}

func TestNaoRemoverRegistroAntesTempoBloqueio(t *testing.T) {

	ratelimiter := inicializarRateLimiter()
	req := inicializarRequestTeste("175.890.789.132", "1234")

	ratelimiter.Controlar(&req, requisicaoPorSegundoIP, requisicaoPorSegundoToken, tempoSegundosBloqueadoIP, tempoSegundosBloqueadoToken, tempoSegundosExpiracaoToken)

	ratelimiter.controlaRateLimit.remover()

	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.890.789.132").Id == "175.890.789.132")

}

func TestGravarToken(t *testing.T) {

	ratelimiter := inicializarRateLimiter()

	token, err := ratelimiter.GerarToken(tempoSegundosExpiracaoToken)
	assert.Equal(t, true, err == nil)

	tokenGravado, err := ratelimiter.controlaRateLimit.buscarToken(token)

	assert.Equal(t, true, err == nil)
	assert.Equal(t, true, tokenGravado.Utilizado)
	assert.Equal(t, true, tokenGravado.Id == token)

}

func TestRetornarMensagemParaTokenExpirado(t *testing.T) {

	ratelimiter := inicializarRateLimiter()

	token, err := ratelimiter.GerarToken(tempoSegundosExpiracaoToken)
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

	token, err := ratelimiter.GerarToken(tempoSegundosExpiracaoToken)
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

	token, err := ratelimiter.GerarToken(tempoSegundosExpiracaoToken)
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

func TestBloquearPorSegundosRequisicaoToken(t *testing.T) {

	const totalRequisicao int = 11

	ratelimiter := inicializarRateLimiter()
	req := inicializarRequestTesteComToken("175.456.879.120", "3421", "token_teste_123")

	for i := 0; i < totalRequisicao; i++ {
		ratelimiter.Controlar(&req, requisicaoPorSegundoIP, requisicaoPorSegundoToken, tempoSegundosBloqueadoIP, tempoSegundosBloqueadoToken, tempoSegundosExpiracaoToken)
	}

	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.456.879.120").TotalRequests == requisicaoPorSegundoToken)
	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.456.879.120").Bloqueado == true)

}

func TestNaoBloquearPorSegundosRequisicaoToken(t *testing.T) {

	const totalRequisicao int = 7

	ratelimiter := inicializarRateLimiter()
	req := inicializarRequestTesteComToken("175.456.879.120", "3421", "token_teste_123")

	for i := 0; i < totalRequisicao; i++ {
		ratelimiter.Controlar(&req, requisicaoPorSegundoIP, requisicaoPorSegundoToken, tempoSegundosBloqueadoIP, tempoSegundosBloqueadoToken, tempoSegundosExpiracaoToken)
	}

	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.456.879.120").TotalRequests < requisicaoPorSegundoToken)
	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.456.879.120").Bloqueado == false)

}

func TestRetornarMensagemTotalRequisicaoPorSegundoExcedida(t *testing.T) {

	const totalRequisicao int = 6

	ratelimiter := inicializarRateLimiter()
	req := inicializarRequestTeste("175.890.789.123", "1234")

	for i := 0; i < totalRequisicao; i++ {
		ratelimiter.Controlar(&req, requisicaoPorSegundoIP, requisicaoPorSegundoToken, tempoSegundosBloqueadoIP, tempoSegundosBloqueadoToken, tempoSegundosExpiracaoToken)
	}

	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.890.789.123").TotalRequests == requisicaoPorSegundoIP)
	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.890.789.123").Bloqueado == true)

	statusCode, err := ratelimiter.Controlar(&req, requisicaoPorSegundoIP, requisicaoPorSegundoToken, tempoSegundosBloqueadoIP, tempoSegundosBloqueadoToken, tempoSegundosExpiracaoToken)

	assert.Equal(t, true, statusCode == http.StatusTooManyRequests)
	assert.Equal(t, true, err.Error() == "you have reached the maximum number of requests or actions allowed within a certain time frame")

}
