package auth

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"hindsight/config"
)

func TestAuthMiddleware(t *testing.T) {
	_, err := config.Init()
	assert.Empty(t, err)

	authMiddleware := GetMiddleware()
	assert.NotEmpty(t, authMiddleware)
}
