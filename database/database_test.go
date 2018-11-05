package database

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"hindsight/config"
)

func TestDatabase(t *testing.T) {
	_, err := config.Init()
	assert.Empty(t, err)

	db := Init()
	assert.NotEmpty(t, db, "Database should not be empty.")
}
