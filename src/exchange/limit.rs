use crate::exchange::constants::MARKET_FLAG_BEIJING;
use crate::exchange::symbol::detect_market;
use crate::internal::helpers::decimal;

pub fn market_limit(security_code: &str) -> f64 {
    let (_, flag, code) = detect_market(security_code);
    if flag == MARKET_FLAG_BEIJING {
        return 0.30;
    }

    let normalized_code = code.trim_start_matches(char::from(0));
    if normalized_code.starts_with("30") || normalized_code.starts_with("68") {
        return 0.20;
    }
    0.10
}

pub fn limit_up(security_code: &str, price: f64) -> f64 {
    let limit = market_limit(security_code);
    let last_close = decimal(price, None);
    decimal(last_close * (1.0 + limit), None)
}
