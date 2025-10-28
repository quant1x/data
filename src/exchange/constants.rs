pub type MarketType = u16;

pub const MARKET_ID_SHENZHEN: MarketType = 0;
pub const MARKET_ID_SHANGHAI: MarketType = 1;
pub const MARKET_ID_BEIJING: MarketType = 2;
pub const MARKET_ID_HONG_KONG: MarketType = 21;
pub const MARKET_ID_USA: MarketType = 22;

pub const MARKET_FLAG_SHANGHAI: &str = "sh";
pub const MARKET_FLAG_SHENZHEN: &str = "sz";
pub const MARKET_FLAG_BEIJING: &str = "bj";
pub const MARKET_FLAG_HONG_KONG: &str = "hk";
pub const MARKET_FLAG_USA: &str = "us";

pub const MARKET_FLAGS: [&str; 10] = ["sh", "sz", "bj", "hk", "us", "SH", "SZ", "BJ", "HK", "US"];
pub const MARKET_A_SHARE_FLAGS: [&str; 6] = ["sh", "sz", "bj", "SH", "SZ", "BJ"];

pub const MARKET_CN_FIRST_DATE: &str = "19901219";
pub const MARKET_CH_FIRST_LISTTIME: &str = "1990-12-19";

pub const STOCK_DELISTING: &str = "DELISTING";
