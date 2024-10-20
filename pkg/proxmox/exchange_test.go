package proxmox

import (
	"encoding/json"
	"github.com/camaeel/proxmox-oidc-credential-helper/pkg/callback"
	"github.com/camaeel/proxmox-oidc-credential-helper/pkg/config"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExchangeCallbackResultForTicket(t *testing.T) {
	cfg := config.Config{
		CallbackPort: 1234,
		CallbackPath: "/callback",
	}
	input := callback.CallbackResult{
		Code:  "code1234",
		State: "state1111",
	}
	expectedRequestBody := map[string]string{
		"code":         input.Code,
		"state":        input.State,
		"redirect-url": "http://localhost:1234/callback",
	}

	expectedOutput := Ticket{
		Data: TicketData{
			Ticket:              "1234",
			Username:            "user1234@realm1",
			CSRFPreventionToken: "csrftpokne111",
		},
	}

	// setup fake http server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api2/json/access/openid/login", r.URL.Path)

		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		parsedBody := map[string]string{}
		err = json.Unmarshal(body, &parsedBody)
		assert.NoError(t, err)
		assert.Equal(t, expectedRequestBody, parsedBody)

		w.WriteHeader(http.StatusOK)
		data, err := json.Marshal(expectedOutput)
		assert.NoError(t, err)
		_, err = w.Write(data)
		assert.NoError(t, err)
	}))
	defer mockServer.Close()

	cfg.ProxmoxURL = mockServer.URL

	output, err := ExchangeCallbackResultForTicket(cfg, input)

	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, output)

}
