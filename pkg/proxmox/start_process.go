package proxmox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/camaeel/proxmox-oidc-credential-helper/pkg/config"
	"io"
	"net/http"
)

func GetOidcURL(cfg config.Config) (string, error) {
	body, err := json.Marshal(map[string]string{
		"redirect-url": getCallbackUrl(cfg),
		"realm":        cfg.Realm,
	})
	if err != nil {
		return "", err
	}
	bodyReader := bytes.NewReader(body)
	requestURL := fmt.Sprintf("%s/api2/json/access/openid/auth-url", cfg.ProxmoxURL)
	req, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return "", err
	}
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("wrong response status: %d", res.StatusCode)
	}

	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	output := map[string]string{}
	err = json.Unmarshal(respBody, &output)
	if err != nil {
		return "", err
	}

	return output["data"], nil
}
