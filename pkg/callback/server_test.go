package callback

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOidcCallbackHandler(t *testing.T) {
	inputCode := "1234abc"
	inputState := "test223"

	cancelFunc := func() { panic("shouldn't happen") }
	resultChan := make(chan CallbackResult)
	handlerFunc := oidcCallbackHandlerFactory(cancelFunc, resultChan)

	// create test request
	req, err := http.NewRequest("GET", "/oidc/callback", nil)
	if err != nil {
		t.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("code", inputCode)
	q.Add("state", inputState)
	req.URL.RawQuery = q.Encode()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlerFunc)

	go func() {
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	}()
	output := <-resultChan

	assert.Equal(t, inputCode, output.Code)
	assert.Equal(t, inputState, output.State)
}
