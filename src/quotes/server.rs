use std::fmt;

/// 服务器主机信息，对标 Go 版本 `Server`。
#[derive(Clone, Debug, Default)]
pub struct Server {
    pub source: String,
    pub name: String,
    pub host: String,
    pub port: u16,
    pub cross_time: i64,
}

impl Server {
    /// 返回 `host:port` 形式的字符串；当 host 或 port 无效时返回空字符串。
    pub fn addr(&self) -> String {
        if self.host.is_empty() || self.port == 0 {
            String::new()
        } else {
            format!("{}:{}", self.host, self.port)
        }
    }
}

impl fmt::Display for Server {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "{}", self.addr())
    }
}
