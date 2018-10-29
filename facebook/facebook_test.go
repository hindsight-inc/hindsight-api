package facebook

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestFacebookInit(t *testing.T) {
	err := Init()
	assert.Empty(t, err)
}
