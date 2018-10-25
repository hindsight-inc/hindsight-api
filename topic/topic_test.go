package topic

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestTopic(t *testing.T) {
	topic := Topic{Title: "t", Content: "t"}
	assert.NotEmpty(t, topic)
}

