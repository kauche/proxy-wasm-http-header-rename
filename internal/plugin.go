package internal

import (
	"errors"
	"fmt"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
	"github.com/tidwall/gjson"
)

var _ types.PluginContext = (*pluginContext)(nil)

type pluginContext struct {
	types.DefaultPluginContext

	configuration *pluginConfiguration
}

func (c *pluginContext) NewHttpContext(_ uint32) types.HttpContext {
	return &httpContext{
		configuration: c.configuration,
	}
}

func (c *pluginContext) OnPluginStart(_ int) types.OnPluginStartStatus {
	config, err := getPluginConfiguration()
	if err != nil {
		proxywasm.LogErrorf("failed to get the plugin configuration: %s", err)
		return types.OnPluginStartStatusFailed
	}

	c.configuration = config

	return types.OnPluginStartStatusOK
}

func getPluginConfiguration() (*pluginConfiguration, error) {
	config, err := proxywasm.GetPluginConfiguration()
	if err != nil {
		if err == types.ErrorStatusNotFound {
			return nil, errors.New("the plugin configuration is not found")
		}

		return nil, fmt.Errorf("failed to get the plugin configuration: %w", err)
	}

	if len(config) == 0 {
		return nil, errors.New("the plugin configuration is empty")
	}

	if !gjson.ValidBytes(config) {
		return nil, errors.New("the plugin configuration is not valid JSON")
	}

	jsonConfig := gjson.ParseBytes(config)

	requestHeadersToRename := jsonConfig.Get(configKeyRequestHeadersToRename).Array()
	if len(requestHeadersToRename) == 0 {
		return nil, errors.New("the request headers to rename are not found")
	}

	headersToRename := make([]requestHeaderToRename, len(requestHeadersToRename))

	for i, header := range requestHeadersToRename {
		h := header.Get(configKeyHeader)

		key := h.Get(configKeyKey).String()
		if key == "" {
			return nil, errors.New("the header key for renaming is empty")
		}

		value := h.Get(configKeyValue).String()
		if value == "" {
			return nil, errors.New("the header value for renaming is empty")
		}

		prefix := h.Get(configKeyPrefix).String()

		headersToRename[i] = requestHeaderToRename{
			header: headerValue{
				key:    key,
				value:  value,
				prefix: prefix,
			},
		}
	}

	return &pluginConfiguration{
		requestHeadersToRename: headersToRename,
	}, nil
}
