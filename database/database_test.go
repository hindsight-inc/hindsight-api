package database

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestDatabase(t *testing.T) {
	db := Init()
	assert.NotEmpty(t, db)
}
