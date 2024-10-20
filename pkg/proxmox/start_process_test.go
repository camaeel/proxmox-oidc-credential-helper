package proxmox

import (
	"encoding/json"
	"github.com/camaeel/proxmox-oidc-credential-helper/pkg/config"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetOidcURL(t *testing.T) {
	cfg := config.Config{
		Realm:        "testRealm",
		CallbackPort: 1234,
		CallbackPath: "/callback",
	}
	expected := "http://proxmox.example.com/some-url-with-redirect"
	expectedRequestBody := map[string]string{
		"redirect-url": "http://localhost:1234/callback",
		"realm":        "testRealm",
	}

	// setup fake http server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api2/json/access/openid/auth-url", r.URL.Path)

		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		parsedBody := map[string]string{}
		err = json.Unmarshal(body, &parsedBody)
		assert.NoError(t, err)
		assert.Equal(t, expectedRequestBody, parsedBody)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data": "http://proxmox.example.com/some-url-with-redirect"}`))
	}))
	defer mockServer.Close()

	cfg.ProxmoxURL = mockServer.URL

	output, err := GetOidcURL(cfg)
	assert.NoError(t, err)
	assert.Equal(t, expected, output)
}
