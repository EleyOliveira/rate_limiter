package ratelimiter

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestGravarRegistro(t *testing.T) {

	ipRequest := &IPRequest{}

	ratelimiter := NewRateLimiter(ipRequest)
	ratelimiter.Controlar("175.890.789.123")
	ratelimiter.Controlar("175.890.789.124")
	ratelimiter.Controlar("175.890.789.125")

	assert.Equal(t, true, ratelimiter.controlaRateLimit.contem("175.890.789.123"))
	assert.Equal(t, true, ratelimiter.controlaRateLimit.contem("175.890.789.124"))
	assert.Equal(t, true, ratelimiter.controlaRateLimit.contem("175.890.789.125"))
	assert.Equal(t, false, ratelimiter.controlaRateLimit.contem("175.890.789.480"))
}
