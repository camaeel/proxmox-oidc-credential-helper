package proxmox

import (
	"github.com/camaeel/proxmox-oidc-credential-helper/pkg/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetCallbackUrl(t *testing.T) {
	cfg := config.Config{
		CallbackPort: 1234,
		CallbackPath: "/callback",
	}
	expected := "http://localhost:1234/callback"

	output := getCallbackUrl(cfg)

	assert.Equal(t, expected, output)
}
