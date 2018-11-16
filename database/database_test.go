package database

import (
	"hindsight/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabase(t *testing.T) {
	provider := new(config.ConfigProvider)
	_, err := config.Init(provider)
	assert.Empty(t, err)

	db := Init()
	assert.NotEmpty(t, db, "Database should not be empty.")
}
