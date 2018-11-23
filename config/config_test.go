package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	//t.Log("Test: Config")
	provider := new(ViperProvider)
	c, err := Init(provider)

	assert.Empty(t, err)
	assert.NotEmpty(t, c)
}

func TestConfig(t *testing.T) {
	c := Shared()
	assert.NotEmpty(t, c.HTTPPort)
	assert.NotEmpty(t, c.MySQLDatabase)
	assert.NotEmpty(t, c.JWTRealm)
	assert.NotEmpty(t, c.FacebookAppID)
	//assert.NotEmpty(t, c.Facebook_disable_test)
}

func TestSecret(t *testing.T) {
	c := Shared()
	assert.NotEmpty(t, c.MySQLPassword)
	assert.NotEmpty(t, c.JWTKey)
	assert.NotEmpty(t, c.FacebookAppSecret)
	assert.NotEmpty(t, c.FacebookAccessToken)
}
