package file

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestImage(t *testing.T) {
	image := Image{Title: "Test"}
	assert.NotEmpty(t, image)
}