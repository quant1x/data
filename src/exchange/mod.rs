pub mod classify;
pub mod constants;
pub mod flags;
pub mod index;
pub mod limit;
pub mod symbol;

pub use classify::{
    assert_block_by_security_code, assert_code, assert_etf_by_market_and_code,
    assert_index_by_market_and_code, assert_index_by_security_code,
    assert_stock_by_market_and_code, assert_stock_by_security_code, correct_security_code,
    TargetKind,
};
pub use constants::{
    MarketType, MARKET_FLAG_BEIJING, MARKET_FLAG_HONG_KONG, MARKET_FLAG_SHANGHAI,
    MARKET_FLAG_SHENZHEN, MARKET_FLAG_USA, MARKET_ID_BEIJING, MARKET_ID_HONG_KONG,
    MARKET_ID_SHANGHAI, MARKET_ID_SHENZHEN, MARKET_ID_USA,
};
pub use flags::get_market_flag;
pub use index::{index_list, A_SHARE_INDEX_LIST};
pub use limit::{limit_up, market_limit};
pub use symbol::{
    detect_market, get_market, get_market_flag_by_symbol, get_market_id, get_security_code,
};
