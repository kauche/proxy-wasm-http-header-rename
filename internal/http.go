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
		if err := c.renameRequestHeader(requestHeaderToRename.header.key, requestHeaderToRename.header.value); err != nil {
			setErrorHTTPResponseWithLog("failed to rename the header: %s", err)
			return types.ActionPause
		}
	}

	return types.ActionContinue
}

func (c *httpContext) renameRequestHeader(origName, newName string) error {
	value, err := proxywasm.GetHttpRequestHeader(origName)
	if err != nil {
		if err == types.ErrorStatusNotFound {
			return nil
		}

		return fmt.Errorf("failed to get the original header, `%s`: %w", origName, err)
	}

	if err := proxywasm.ReplaceHttpRequestHeader(newName, value); err != nil {
		return fmt.Errorf("failed to set the new header, `%s`: %w", newName, err)
	}

	if err := proxywasm.RemoveHttpRequestHeader(origName); err != nil {
		return fmt.Errorf("failed to delete the original header, `%s`: %w", origName, err)
	}

	return nil
}

func setErrorHTTPResponseWithLog(format string, args ...interface{}) {
	proxywasm.LogErrorf(format, args...)
	if err := proxywasm.SendHttpResponse(500, nil, []byte(`{"error": "internal server error"}`), -1); err != nil {
		proxywasm.LogErrorf("failed to set the http error response: %s", err)
	}
}
