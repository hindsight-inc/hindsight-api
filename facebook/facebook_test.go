package facebook

import (
	"log"
	"testing"
	"github.com/stretchr/testify/assert"
	"hindsight/config"
)

func TestFacebookInit(t *testing.T) {
	if cfg, err := config.Init(); cfg.Facebook_disable_test {
		log.Println("Facebook test disabled due to configuration")
		return
	} else {
		assert.Empty(t, err)
	}

	err := Init()
	assert.Empty(t, err)
}
