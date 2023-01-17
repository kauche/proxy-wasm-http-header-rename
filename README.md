# proxy-wasm-http-header-rename

A [proxy-wasm](https://github.com/proxy-wasm/spec) compliant WebAssembly module for renaming HTTP Headers.

## Usage

1. Download the latest WebAssembly module binary from the [release page](https://github.com/kauche/proxy-wasm-http-header-rename/releases).

2. Configure the proxy to use the WebAssembly module like below (this assumes [Envoy](https://www.envoyproxy.io/) as the proxy):

```yaml
listeners:
  - name: example
    filter_chains:
      - filters:
          - name: envoy.filters.network.http_connection_manager
            typed_config:
              # ...
              http_filters:
                - name: envoy.filters.http.wasm
                  typed_config:
                    '@type': type.googleapis.com/udpa.type.v1.TypedStruct
                    type_url: type.googleapis.com/envoy.extensions.filters.http.wasm.v3.Wasm
                    value:
                      config:
                        vm_config:
                          runtime: envoy.wasm.runtime.v8
                          code:
                            local:
                              filename: /etc/envoy/proxy-wasm-http-header-rename.wasm
                        configuration:
                          "@type": type.googleapis.com/google.protobuf.StringValue
                          value: |
                            {
                              "request_headers_to_rename": [
                                {
                                  "header": {
                                    "key": "original-header-name",
                                    "value": "new-header-name"
                                  }
                                }
                              ]
                            }
                - name: envoy.filters.http.router
                  typed_config:
                    '@type': type.googleapis.com/envoy.extensions.filters.http.router.v3.Router
# ...
```

## Motivation

For now, Envoy does not support renaming HTTP Headers natively as described in [this issue](https://github.com/envoyproxy/envoy/issues/8947). So, we can use this WebAssembly module to rename HTTP Headers until the issue is resolved.
