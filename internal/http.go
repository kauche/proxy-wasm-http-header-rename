package internal

import (
	"fmt"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

var _ types.HttpContext = (*httpContext)(nil)

type httpContext struct {
	types.DefaultHttpContext

	configuration *pluginConfiguration
}

func (c *httpContext) OnHttpRequestHeaders(_ int, _ bool) types.Action {
	for _, requestHeaderToRename := range c.configuration.requestHeadersToRename {
		if err := c.renameRequestHeader(requestHeaderToRename.header); err != nil {
			setErrorHTTPResponseWithLog("failed to rename the header: %s", err)
			return types.ActionPause
		}
	}

	return types.ActionContinue
}

func (c *httpContext) renameRequestHeader(h headerValue) error {
	value, err := proxywasm.GetHttpRequestHeader(h.key)
	if err != nil {
		if err == types.ErrorStatusNotFound {
			return nil
		}

		return fmt.Errorf("failed to get the original header, `%s`: %w", h.key, err)
	}

	newValue := value
	if h.prefix != "" {
		newValue = h.prefix + value
	}

	if err := proxywasm.ReplaceHttpRequestHeader(h.value, newValue); err != nil {
		return fmt.Errorf("failed to set the new header, `%s`: %w", h.value, err)
	}

	if err := proxywasm.RemoveHttpRequestHeader(h.key); err != nil {
		return fmt.Errorf("failed to delete the original header, `%s`: %w", h.key, err)
	}

	return nil
}

func setErrorHTTPResponseWithLog(format string, args ...interface{}) {
	proxywasm.LogErrorf(format, args...)
	if err := proxywasm.SendHttpResponse(500, nil, []byte(`{"error": "internal server error"}`), -1); err != nil {
		proxywasm.LogErrorf("failed to set the http error response: %s", err)
	}
}
