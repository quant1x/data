use std::collections::HashMap;
use std::fs;
use std::path::Path;
use std::sync::{Arc, Mutex, Once};
use std::time::{Duration, SystemTime};

use serde::{Deserialize, Serialize};

use crate::quotes::cmd_security_list::Security;

/// Security target information
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SecurityTarget {
    pub code: String,
    pub vol_unit: u16,
    pub decimal_point: i8,
    pub name: String,
}

impl From<&Security> for SecurityTarget {
    fn from(security: &Security) -> Self {
        Self {
            code: security.code.clone(),
            vol_unit: security.vol_unit,
            decimal_point: security.decimal_point,
            name: security.name.clone(),
        }
    }
}

lazy_static::lazy_static! {
    static ref CACHE_SECURITY_CODE_LIST: String = {
        // This would need to be set from environment/cache path
        // For now, use a default path
        "targets.csv".to_string()
    };

    static ref MAP_STOCK_LIST: Arc<Mutex<HashMap<String, SecurityTarget>>> = Arc::new(Mutex::new(HashMap::new()));
    static ref ONCE_STOCK_LIST: Once = Once::new();
}

/// Read cached security list from CSV file
fn read_cache_security_list() -> Result<(), Box<dyn std::error::Error>> {
    let filename = CACHE_SECURITY_CODE_LIST.as_str();
    if !Path::new(filename).exists() {
        return Ok(());
    }

    let content = fs::read_to_string(filename)?;
    let mut rdr = csv::Reader::from_reader(content.as_bytes());
    let mut list = Vec::new();

    for result in rdr.deserialize() {
        let record: SecurityTarget = result?;
        list.push(record);
    }

    if list.is_empty() {
        return Ok(());
    }

    let mut map = MAP_STOCK_LIST.lock().unwrap();
    for item in list {
        map.insert(item.code.clone(), item);
    }

    Ok(())
}

/// Write security list to CSV cache file
fn write_cache_security_list(list: &[SecurityTarget]) -> Result<(), Box<dyn std::error::Error>> {
    let filename = CACHE_SECURITY_CODE_LIST.as_str();

    // Sort by code
    let mut sorted_list = list.to_vec();
    sorted_list.sort_by(|a, b| a.code.cmp(&b.code));

    let mut wtr = csv::Writer::from_path(filename)?;
    for item in sorted_list {
        wtr.serialize(&item)?;
    }
    wtr.flush()?;

    Ok(())
}

/// Lazy load stock list with caching
fn lazy_load_stock_list() -> Result<(), Box<dyn std::error::Error>> {
    let filename = CACHE_SECURITY_CODE_LIST.as_str();
    let mut b_updated = false;

    if !Path::new(filename).exists() {
        // File doesn't exist, needs to be created
        b_updated = true;
    } else {
        // File exists, load it first
        read_cache_security_list()?;

        // Check if file needs updating (simplified - check if older than 1 day)
        if let Ok(metadata) = fs::metadata(filename) {
            if let Ok(modified) = metadata.modified() {
                if let Ok(duration) = SystemTime::now().duration_since(modified) {
                    if duration > Duration::from_secs(24 * 60 * 60) {
                        b_updated = true;
                    }
                }
            }
        }
    }

    if !b_updated {
        return Ok(());
    }

    let list = get_security_list()?;
    if list.is_empty() {
        return Ok(());
    }

    // Update cache
    let mut map = MAP_STOCK_LIST.lock().unwrap();
    let mut new_list = Vec::new();

    for stock in list {
        let target = SecurityTarget::from(&stock);
        let code = target.code.clone();
        map.insert(code, target.clone());
        new_list.push(target);
    }

    // Write to cache file
    write_cache_security_list(&new_list)?;

    Ok(())
}

/// Get security list from exchanges
fn get_security_list() -> Result<Vec<Security>, Box<dyn std::error::Error>> {
    let std_api = crate::quotes::api::new_std_api()?;
    let mut all_list = Vec::new();

    const SECURITY_LIST_MAX: u16 = 1000;

    // Get Shanghai market securities (market = 1)
    let mut start: u16 = 0;
    loop {
        let reply = std_api.security_list(1, start)?;
        let mut securities = reply.list.clone();

        for security in &mut securities {
            security.code = format!("sh{}", security.code);
        }

        if !securities.is_empty() {
            all_list.extend(securities);
        }

        if reply.count < SECURITY_LIST_MAX {
            break;
        }
        start += reply.count;
    }

    // Get Shenzhen market securities (market = 0)
    start = 0;
    loop {
        let reply = std_api.security_list(0, start)?;
        let mut securities = reply.list.clone();

        for security in &mut securities {
            security.code = format!("sz{}", security.code);
        }

        if !securities.is_empty() {
            all_list.extend(securities);
        }

        if reply.count < SECURITY_LIST_MAX {
            break;
        }
        start += reply.count;
    }

    Ok(all_list)
}

/// Checkout security info by code
pub fn checkout_security_info(security_code: &str) -> Option<SecurityTarget> {
    ONCE_STOCK_LIST.call_once(|| {
        lazy_load_stock_list().unwrap_or_else(|e| {
            eprintln!("Failed to load stock list: {}", e);
        });
    });

    // Correct security code format if needed
    let corrected_code = correct_security_code(security_code);

    let map = MAP_STOCK_LIST.lock().unwrap();
    map.get(&corrected_code).cloned()
}

/// Get security base unit for price calculations
pub fn security_base_unit(market_id: u32, code: &str) -> f64 {
    let security_code = get_security_code(market_id, code);
    if let Some(security_info) = checkout_security_info(&security_code) {
        10_f64.powi(security_info.decimal_point as i32)
    } else {
        100.0
    }
}

/// Correct security code format
fn correct_security_code(code: &str) -> String {
    // Simplified implementation - in real code this would normalize the format
    code.to_string()
}

/// Get security code with market prefix
fn get_security_code(market_id: u32, code: &str) -> String {
    match market_id {
        1 => format!("sh{}", code), // Shanghai
        0 => format!("sz{}", code), // Shenzhen
        _ => code.to_string(),
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_checkout_security_info() {
        // This test would require a mock or actual server
        // For now, just test the structure
        let _result = checkout_security_info("000001");
        // Assert would depend on actual data
    }

    #[test]
    fn test_security_base_unit() {
        // Test with default value when security not found
        let unit = security_base_unit(1, "nonexistent");
        assert_eq!(unit, 100.0);
    }
}