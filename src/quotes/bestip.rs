use anyhow::{anyhow, Result};

use crate::quotes::base_consts::QuotesError;
use crate::quotes::server::Server;

pub const HOST_HQ: &str = "HQ";
pub const HOST_EX: &str = "EX";
pub const HOST_GP: &str = "GP";

/// TODO: 迁移 `bestip.go` 的测速逻辑。
pub fn get_fast_host(_key: &str) -> Result<Vec<Server>> {
    Err(anyhow!(
        QuotesError::Unimplemented("bestip::get_fast_host",)
    ))
}
