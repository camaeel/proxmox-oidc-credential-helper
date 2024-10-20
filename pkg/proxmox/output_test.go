package proxmox

import (
	"fmt"
	"github.com/camaeel/proxmox-oidc-credential-helper/pkg/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOutputTicket(t *testing.T) {
	input := Ticket{
		Data: TicketData{
			Ticket:              "1234ticket",
			CSRFPreventionToken: "csrf1token",
			Username:            "login@realm1",
		},
	}
	testCases := []struct {
		name           string
		outputFormat   string
		expectedOutput string
		expectedError  error
	}{
		{
			name:           "json",
			outputFormat:   "json",
			expectedOutput: "{\"data\":{\"ticket\":\"1234ticket\",\"CSRFPreventionToken\":\"csrf1token\",\"username\":\"login@realm1\"}}",
			expectedError:  nil,
		},
		{
			name:           "text",
			outputFormat:   "text",
			expectedOutput: fmt.Sprintf("export PROXMOX_VE_AUTH_TICKET='%s'\nexport PROXMOX_VE_CSRF_PREVENTION_TOKEN='%s'\n", input.Data.Ticket, input.Data.CSRFPreventionToken),
			expectedError:  nil,
		},
		{
			name:           "wrong",
			outputFormat:   "wrong",
			expectedOutput: "",
			expectedError:  fmt.Errorf("unknown output format: wrong"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := config.Config{
				OutputFormat: tc.outputFormat,
			}
			output, err := OutputTicket(cfg, input)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedOutput, output)
		})
	}

}
