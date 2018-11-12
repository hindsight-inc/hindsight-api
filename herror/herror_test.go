package herror

import (
	"net/http"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	msg := "Test only"
	code, h := Bad(DomainTopicCreate, ReasonNonexistentEntry, msg)
	assert.Equal(t, code, http.StatusBadRequest)
	assert.NotEmpty(t, h)
	// TODO: validate contents of h
}