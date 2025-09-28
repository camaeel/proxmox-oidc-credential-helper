package proxmox

import (
	"fmt"

	"github.com/proxmox-oidc-credential-helper/proxmox-oidc-credential-helper/pkg/config"
)

func getCallbackUrl(cfg config.Config) string {
	return fmt.Sprintf("http://localhost:%d%s", cfg.CallbackPort, cfg.CallbackPath)
}
