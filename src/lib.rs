use proxy_wasm::traits::{Context, HttpContext, RootContext};
use proxy_wasm::types::{Action, ContextType, LogLevel};
use serde::{Deserialize, Serialize};

#[cfg(not(all(target_arch = "wasm32", target_os = "unknown")))]

proxy_wasm::main! {{
    proxy_wasm::set_log_level(LogLevel::Trace);
        proxy_wasm::set_root_context(|_| -> Box<dyn RootContext> { Box::new(HTTPHeaderRenameRoot{
            config: Config { request_headers_to_rename: vec![] }
        })
    });
}}

struct HTTPHeaderRenameRoot {
    config: Config,
}

impl Context for HTTPHeaderRenameRoot {}

impl RootContext for HTTPHeaderRenameRoot {
    fn get_type(&self) -> Option<ContextType> {
        Some(ContextType::HttpContext)
    }

    fn on_configure(&mut self, _: usize) -> bool {
        if let Some(config_bytes) = self.get_plugin_configuration() {
            self.config = serde_json::from_reader(config_bytes.as_slice()).unwrap();
        }
        true
    }

    fn create_http_context(&self, _: u32) -> Option<Box<dyn HttpContext>> {
        Some(Box::new(HTTPHeadeRename {
            config: self.config.clone(),
        }))
    }
}

struct HTTPHeadeRename {
    config: Config,
}

impl Context for HTTPHeadeRename {}

impl HttpContext for HTTPHeadeRename {
    fn on_http_request_headers(&mut self, _: usize, _: bool) -> Action {
        for header in &self.config.request_headers_to_rename {
            let old_header = header.header.key.as_str();
            let header_value = self.get_http_request_header(old_header);
            if let Some(header_value) = header_value {
                let mut value = header_value.to_string();
                if let Some(prefix) = &header.header.prefix {
                    value = format!("{}{}", prefix, header_value);
                }

                self.set_http_request_header(header.header.value.as_str(), Some(&value));
                self.remove_http_request_header(old_header);
            }
        }

        Action::Continue
    }
}

#[derive(Serialize, Deserialize, Clone)]
struct Config {
    request_headers_to_rename: Vec<HTTPHeaderRenameRule>,
}

#[derive(Serialize, Deserialize, Clone)]
struct HTTPHeaderRenameRule {
    header: HTTPHeader,
}

#[derive(Serialize, Deserialize, Clone)]
struct HTTPHeader {
    key: String,
    value: String,
    prefix: Option<String>,
}
