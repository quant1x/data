use std::fmt;
use std::sync::Arc;
use std::time::Duration;

use super::server::Server;

/// 配置通用连接参数，对标 Go 版本 `Options`。
pub struct Options {
    pub connection_timeout: Duration,
    pub read_timeout: Duration,
    pub write_timeout: Duration,
    pub max_retry_times: usize,
    pub retry_duration: Duration,
    pub release_address: Option<Arc<dyn Fn(&Server) + Send + Sync>>, // 归还服务器地址回调
}

impl Default for Options {
    fn default() -> Self {
        Self {
            connection_timeout: Duration::from_secs(10),
            read_timeout: Duration::from_secs(5),
            write_timeout: Duration::from_secs(5),
            max_retry_times: 3,
            retry_duration: Duration::from_secs(1),
            release_address: None,
        }
    }
}

impl Clone for Options {
    fn clone(&self) -> Self {
        Self {
            connection_timeout: self.connection_timeout,
            read_timeout: self.read_timeout,
            write_timeout: self.write_timeout,
            max_retry_times: self.max_retry_times,
            retry_duration: self.retry_duration,
            release_address: self.release_address.clone(),
        }
    }
}

impl fmt::Debug for Options {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        f.debug_struct("Options")
            .field("connection_timeout", &self.connection_timeout)
            .field("read_timeout", &self.read_timeout)
            .field("write_timeout", &self.write_timeout)
            .field("max_retry_times", &self.max_retry_times)
            .field("retry_duration", &self.retry_duration)
            .finish()
    }
}
