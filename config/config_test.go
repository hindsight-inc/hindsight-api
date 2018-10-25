package config

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	//t.Log("Test: Config")
	c, err := Init()

	assert.Empty(t, err)
	assert.NotEmpty(t, c)
	assert.NotEmpty(t, c.Facebook_app_id)
	assert.NotEmpty(t, c.Facebook_app_secret)
}
