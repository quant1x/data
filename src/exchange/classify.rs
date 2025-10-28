use crate::exchange::constants::{MarketType, MARKET_ID_SHANGHAI, MARKET_ID_SHENZHEN};
use crate::exchange::symbol::detect_market;

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum TargetKind {
    Stock,
    Index,
    Block,
    Etf,
}

pub fn assert_index_by_market_and_code(market_id: MarketType, symbol: &str) -> bool {
    let trimmed = symbol.trim_start_matches(char::from(0));
    match market_id {
        MARKET_ID_SHANGHAI => {
            trimmed.starts_with("000") || trimmed.starts_with("880") || trimmed.starts_with("881")
        }
        MARKET_ID_SHENZHEN => trimmed.starts_with("399"),
        _ => false,
    }
}

pub fn assert_index_by_security_code(security_code: &str) -> bool {
    let (market_id, _, code) = detect_market(security_code);
    assert_index_by_market_and_code(market_id, &code)
}

pub fn assert_block_by_security_code(security_code: &mut String) -> bool {
    let (market_id, flag, code) = detect_market(security_code);
    if market_id != MARKET_ID_SHANGHAI {
        return false;
    }
    if !(code.starts_with("880") || code.starts_with("881")) {
        return false;
    }
    *security_code = format!("{flag}{code}");
    true
}

pub fn assert_etf_by_market_and_code(market_id: MarketType, symbol: &str) -> bool {
    let trimmed = symbol.trim_start_matches(char::from(0));
    match market_id {
        MARKET_ID_SHANGHAI => trimmed.starts_with("51"),
        MARKET_ID_SHENZHEN => trimmed.starts_with("159"),
        _ => false,
    }
}

pub fn assert_stock_by_market_and_code(market_id: MarketType, symbol: &str) -> bool {
    let trimmed = symbol.trim_start_matches(char::from(0));
    match market_id {
        MARKET_ID_SHANGHAI => {
            trimmed.starts_with("60") || trimmed.starts_with("68") || trimmed.starts_with("51")
        }
        MARKET_ID_SHENZHEN => {
            trimmed.starts_with("00") || trimmed.starts_with("30") || trimmed.starts_with("159")
        }
        _ => false,
    }
}

pub fn assert_stock_by_security_code(security_code: &str) -> bool {
    let (market_id, _, code) = detect_market(security_code);
    assert_stock_by_market_and_code(market_id, &code)
}

pub fn correct_security_code(security_code: &str) -> String {
    if security_code.trim().is_empty() {
        return String::new();
    }
    let (_, flag, code) = detect_market(security_code);
    format!("{flag}{code}")
}

pub fn assert_code(security_code: &str) -> TargetKind {
    let (market_id, _, code) = detect_market(security_code);
    if market_id == MARKET_ID_SHANGHAI && (code.starts_with("880") || code.starts_with("881")) {
        return TargetKind::Block;
    }
    if market_id == MARKET_ID_SHANGHAI && code.starts_with("000") {
        return TargetKind::Index;
    }
    if market_id == MARKET_ID_SHENZHEN && code.starts_with("399") {
        return TargetKind::Index;
    }
    if market_id == MARKET_ID_SHANGHAI && code.starts_with("51") {
        return TargetKind::Etf;
    }
    if market_id == MARKET_ID_SHENZHEN && code.starts_with("159") {
        return TargetKind::Etf;
    }
    TargetKind::Stock
}
