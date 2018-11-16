package facebook

import (
	"hindsight/config"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFacebookInit(t *testing.T) {
	provider := new(config.ConfigProvider)
	if cfg, err := config.Init(provider); cfg.Facebook_disable_test {
		log.Println("Facebook test disabled due to configuration")
		return
	} else {
		assert.Empty(t, err)
	}

	err := Init()
	assert.Empty(t, err)
}
