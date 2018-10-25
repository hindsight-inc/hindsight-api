package auth

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	authMiddleware := GetMiddleware()
	assert.NotEmpty(t, authMiddleware)
}
