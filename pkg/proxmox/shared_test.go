package proxmox

import (
	"testing"

	"github.com/proxmox-oidc-credential-helper/proxmox-oidc-credential-helper/pkg/config"
	"github.com/stretchr/testify/assert"
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
