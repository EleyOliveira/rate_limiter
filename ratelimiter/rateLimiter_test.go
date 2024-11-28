package ratelimiter

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestGravarRegistro(t *testing.T) {

	cache := &CacheRegistro{}

	ratelimiter := NewRateLimiter(cache)
	ratelimiter.Controlar("175.890.789.123")
	ratelimiter.Controlar("175.890.789.124")
	ratelimiter.Controlar("175.890.789.125")

	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.890.789.123").Id == "175.890.789.123")
	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.890.789.124").Id == "175.890.789.124")
	assert.Equal(t, true, ratelimiter.controlaRateLimit.buscar("175.890.789.125").Id == "175.890.789.125")
	assert.Equal(t, false, ratelimiter.controlaRateLimit.buscar("175.890.789.480").Id == "175.890.789.480")
}
