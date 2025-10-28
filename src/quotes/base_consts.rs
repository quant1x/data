use thiserror::Error;

/// 与 Go 版本保持一致的市场标识。
#[derive(Debug, Copy, Clone, PartialEq, Eq)]
pub enum TdxMarket {
    Unknown,
}

pub const DEFAULT_RETRY_TIMES: usize = 3;
pub const MESSAGE_HEADER_BYTES: usize = 0x10;
pub const MESSAGE_MAX_BYTES: usize = 1 << 15;

pub const POOL_INITED: usize = 1;
pub const POOL_MAX: usize = 10;
pub const CONN_TIMEOUT: u64 = 10;
pub const RECV_TIMEOUT: u64 = 5;

#[derive(Debug, Error)]
pub enum QuotesError {
    #[error("more than 8M data")]
    BadData,
    #[error("invalid server address")]
    InvalidServerAddress,
    #[error("feature not yet implemented: {0}")]
    Unimplemented(&'static str),
}
