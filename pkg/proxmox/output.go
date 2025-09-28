package proxmox

import (
	"encoding/json"
	"fmt"

	"github.com/proxmox-oidc-credential-helper/proxmox-oidc-credential-helper/pkg/config"
)

func OutputTicket(cfg config.Config, ticket Ticket) (string, error) {
	switch cfg.OutputFormat {
	case "text":
		return fmt.Sprintf("export PROXMOX_VE_AUTH_TICKET='%s'\nexport PROXMOX_VE_CSRF_PREVENTION_TOKEN='%s'\n", ticket.Data.Ticket, ticket.Data.CSRFPreventionToken), nil
	case "json":
		data, err := json.Marshal(ticket)
		return string(data), err
	}
	return "", fmt.Errorf("unknown output format: %s", cfg.OutputFormat)
}
