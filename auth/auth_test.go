package auth

import (
	"hindsight/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	provider := new(config.ViperProvider)
	_, err := config.Init(provider)
	assert.Empty(t, err)

	authMiddleware := GetMiddleware()
	assert.NotEmpty(t, authMiddleware)
}
