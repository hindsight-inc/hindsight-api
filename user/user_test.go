package user

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	user := User{Username: "u", Password: "p"}
	assert.NotEmpty(t, user)
}
