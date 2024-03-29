package internal

const (
	configKeyRequestHeadersToRename = "request_headers_to_rename"
	configKeyHeader                 = "header"
	configKeyKey                    = "key"
	configKeyValue                  = "value"
	configKeyPrefix                 = "prefix"
)

type pluginConfiguration struct {
	requestHeadersToRename []requestHeaderToRename
}

type requestHeaderToRename struct {
	header headerValue
}

type headerValue struct {
	key    string
	value  string
	prefix string
}
