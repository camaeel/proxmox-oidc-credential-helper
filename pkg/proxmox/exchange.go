package proxmox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/proxmox-oidc-credential-helper/proxmox-oidc-credential-helper/pkg/callback"
	"github.com/proxmox-oidc-credential-helper/proxmox-oidc-credential-helper/pkg/config"
)

func ExchangeCallbackResultForTicket(cfg config.Config, result callback.CallbackResult) (Ticket, error) {
	body, err := json.Marshal(map[string]string{
		"code":         result.Code,
		"state":        result.State,
		"redirect-url": getCallbackUrl(cfg),
	})
	if err != nil {
		return Ticket{}, err
	}
	bodyReader := bytes.NewReader(body)
	requestURL := fmt.Sprintf("%s/api2/json/access/openid/login", cfg.ProxmoxURL)
	req, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)
	if err != nil {
		return Ticket{}, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return Ticket{}, err
	}
	if res.StatusCode != http.StatusOK {
		return Ticket{}, fmt.Errorf("wrong response status: %d", res.StatusCode)
	}
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return Ticket{}, err
	}

	ticket := Ticket{}
	err = json.Unmarshal(respBody, &ticket)
	return ticket, err
}
