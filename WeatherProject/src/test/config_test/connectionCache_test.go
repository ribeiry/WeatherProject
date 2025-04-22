package config_test

import (
	"WeatherProject/config"
	"testing"

	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
)

func TestInitCache(t *testing.T) {
	config.InitCache()

	assert.NotNil(t, config.CachingInMemory, "O cache deveria estar inicializado")

	// Opcional: inserir e recuperar do cache para garantir funcionamento
	config.CachingInMemory.Set("chave", "valor", cache.DefaultExpiration)
	val, found := config.CachingInMemory.Get("chave")

	assert.True(t, found, "A chave deveria estar presente no cache")
	assert.Equal(t, "valor", val)
}
