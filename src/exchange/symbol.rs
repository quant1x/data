use crate::exchange::constants::{
    MarketType, MARKET_FLAG_BEIJING, MARKET_FLAG_HONG_KONG, MARKET_FLAG_SHANGHAI,
    MARKET_FLAG_SHENZHEN, MARKET_FLAG_USA, MARKET_ID_BEIJING, MARKET_ID_HONG_KONG,
    MARKET_ID_SHANGHAI, MARKET_ID_SHENZHEN, MARKET_ID_USA,
};

const MARKET_FLAGS_LOWER: [&str; 5] = [
    MARKET_FLAG_SHANGHAI,
    MARKET_FLAG_SHENZHEN,
    MARKET_FLAG_BEIJING,
    MARKET_FLAG_HONG_KONG,
    MARKET_FLAG_USA,
];

const SH_PREFIXES: [&str; 9] = ["50", "51", "60", "68", "90", "110", "113", "132", "204"];
const SH_SECONDARY_PREFIXES: [&str; 4] = ["5", "6", "7", "9"];
const SZ_PREFIXES: [&str; 11] = [
    "00", "12", "13", "15", "16", "18", "20", "30", "39", "115", "1318",
];
const BJ_PREFIXES: [&str; 2] = ["4", "8"];

fn starts_with_any(value: &str, prefixes: &[&str]) -> bool {
    prefixes.iter().any(|prefix| value.starts_with(prefix))
}

fn ends_with_any(value: &str, suffixes: &[&str]) -> bool {
    suffixes.iter().any(|suffix| value.ends_with(suffix))
}

fn map_flag(flag: &str) -> &'static str {
    match flag {
        "sz" => MARKET_FLAG_SHENZHEN,
        "bj" => MARKET_FLAG_BEIJING,
        "hk" => MARKET_FLAG_HONG_KONG,
        "us" => MARKET_FLAG_USA,
        _ => MARKET_FLAG_SHANGHAI,
    }
}

fn truncate_symbol(symbol: &str, len: usize) -> String {
    symbol.chars().take(len).collect()
}

pub fn get_security_code(market: MarketType, symbol: &str) -> String {
    let trimmed = symbol.trim();
    if trimmed.is_empty() {
        return String::new();
    }

    match market {
        MARKET_ID_USA => format!("{MARKET_FLAG_USA}{trimmed}"),
        MARKET_ID_HONG_KONG => {
            format!("{MARKET_FLAG_HONG_KONG}{}", truncate_symbol(trimmed, 5))
        }
        MARKET_ID_BEIJING => format!("{MARKET_FLAG_BEIJING}{}", truncate_symbol(trimmed, 6)),
        MARKET_ID_SHENZHEN => format!("{MARKET_FLAG_SHENZHEN}{}", truncate_symbol(trimmed, 6)),
        _ => format!("{MARKET_FLAG_SHANGHAI}{}", truncate_symbol(trimmed, 6)),
    }
}

pub fn get_market(symbol: &str) -> &'static str {
    let trimmed = symbol.trim();
    if trimmed.is_empty() {
        return MARKET_FLAG_SHANGHAI;
    }

    let lower = trimmed.to_ascii_lowercase();
    if lower.len() >= 2 {
        if starts_with_any(&lower, &MARKET_FLAGS_LOWER) {
            return map_flag(&lower[..2]);
        }
        if ends_with_any(&lower, &MARKET_FLAGS_LOWER) {
            let len = lower.len();
            return map_flag(&lower[len - 2..]);
        }
    }

    if starts_with_any(&lower, &SH_PREFIXES)
        || starts_with_any(&lower, &SH_SECONDARY_PREFIXES)
        || lower.starts_with("88")
    {
        return MARKET_FLAG_SHANGHAI;
    }
    if starts_with_any(&lower, &SZ_PREFIXES) {
        return MARKET_FLAG_SHENZHEN;
    }
    if starts_with_any(&lower, &BJ_PREFIXES) {
        return MARKET_FLAG_BEIJING;
    }

    MARKET_FLAG_SHANGHAI
}

pub fn get_market_flag_by_symbol(symbol: &str) -> &'static str {
    get_market(symbol)
}

fn flag_to_market_id(flag: &str) -> MarketType {
    match flag {
        MARKET_FLAG_SHENZHEN => MARKET_ID_SHENZHEN,
        MARKET_FLAG_BEIJING => MARKET_ID_BEIJING,
        MARKET_FLAG_HONG_KONG => MARKET_ID_HONG_KONG,
        MARKET_FLAG_USA => MARKET_ID_USA,
        _ => MARKET_ID_SHANGHAI,
    }
}

pub fn get_market_id(symbol: &str) -> MarketType {
    flag_to_market_id(get_market(symbol))
}

pub fn detect_market(symbol: &str) -> (MarketType, String, String) {
    let trimmed = symbol.trim();
    if trimmed.is_empty() {
        return (
            MARKET_ID_SHANGHAI,
            MARKET_FLAG_SHANGHAI.to_string(),
            String::new(),
        );
    }

    let market_flag = get_market(trimmed);
    let lower = trimmed.to_ascii_lowercase();
    let mut code = trimmed.to_string();

    if lower.len() >= market_flag.len() && lower.starts_with(market_flag) {
        code = trimmed[market_flag.len()..].to_string();
        if code.starts_with('.') {
            code = code[1..].to_string();
        }
    } else if lower.len() >= market_flag.len() && lower.ends_with(market_flag) {
        let suffix_start = trimmed.len().saturating_sub(market_flag.len());
        if suffix_start > 0 && trimmed.as_bytes()[suffix_start - 1] == b'.' {
            code = trimmed[..suffix_start - 1].to_string();
        } else {
            code = trimmed[..suffix_start].to_string();
        }
    }

    (
        flag_to_market_id(market_flag),
        market_flag.to_string(),
        code,
    )
}
