package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	echoserverpb "github.com/110y/echoserver/echoserver/api/v1"
)

func TestE2E(t *testing.T) {
	t.Parallel()

	host := "upstream"

	req, err := createHTTPRequest(host)
	if err != nil {
		t.Errorf("failed to create a http request: %s", err)
		return
	}

	req.Header.Set("original-header-1", "original-header-value-1")
	req.Header.Set("original-header-2", "original-header-value-2")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("failed to send first http request: %s", err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		t.Errorf("invalid http status code for the first http request: %s", res.Status)
		return
	}

	echores := new(echoserverpb.EchoResponse)
	if err = json.NewDecoder(res.Body).Decode(echores); err != nil {
		t.Errorf("failed to marshal the first response to json: %s", err)
		return
	}

	newHeader1, ok := echores.Headers["new-header-1"]
	if !ok {
		t.Error("new-header-1 header is not found")
		return
	}
	if len(newHeader1.Value) != 1 {
		t.Errorf("new-header-1 header has invalid number of values: %d", len(newHeader1.Value))
		return
	}
	if newHeader1.Value[0] != "original-header-value-1" {
		t.Errorf("the value for new-header-1 is expected to be `original-header-value1`, but got `%s`", newHeader1.Value[0])
		return
	}

	newHeader2, ok := echores.Headers["new-header-2"]
	if !ok {
		t.Error("new-header-2 header is not found")
		return
	}
	if len(newHeader1.Value) != 1 {
		t.Errorf("new-header-2 header has invalid number of values: %d", len(newHeader1.Value))
		return
	}
	if newHeader2.Value[0] != "original-header-value-2" {
		t.Errorf("the value for new-header-2 is expected to be `original-header-value1`, but got `%s`", newHeader2.Value[0])
		return
	}

	_, ok = echores.Headers["original-header-1"]
	if ok {
		t.Error("original-header-1 header should be removed")
		return
	}

	_, ok = echores.Headers["original-header-2"]
	if ok {
		t.Error("original-header-2 header should be removed")
		return
	}
}

func createHTTPRequest(host string) (*http.Request, error) {
	addr := os.Getenv("ENVOY_ADDRESS")
	if addr == "" {
		addr = "localhost:9090"
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/", addr), bytes.NewBuffer([]byte(`{"message":"hello"}`)))
	if err != nil {
		return nil, err
	}

	req.Host = host
	req.Header.Set("content-type", "application/json")

	return req, nil
}
