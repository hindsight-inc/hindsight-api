package config

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	//t.Log("Test: Config")
	c, err := Init()

	assert.Empty(t, err)
	assert.NotEmpty(t, c)
}

func TestConfig(t *testing.T) {
	c := Shared()
	assert.NotEmpty(t, c.MySQL_database)
	assert.NotEmpty(t, c.JWT_Realm)
	assert.NotEmpty(t, c.Facebook_app_id)
}

func TestSecret(t *testing.T) {
	c := Shared()
	assert.NotEmpty(t, c.MySQL_password)
	assert.NotEmpty(t, c.JWT_Key)
	assert.NotEmpty(t, c.Facebook_app_secret)
	assert.NotEmpty(t, c.Facebook_access_token)
}
