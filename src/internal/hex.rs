use std::fmt::Write;

/// 将形如 `"01 0A"` 的十六进制字符串转换为字节数组。
pub fn hex_string_to_bytes(hex_str: &str) -> Option<Vec<u8>> {
    let filtered: String = hex_str.chars().filter(|c| !c.is_whitespace()).collect();
    if filtered.len() % 2 != 0 {
        return None;
    }
    (0..filtered.len())
        .step_by(2)
        .map(|i| u8::from_str_radix(&filtered[i..i + 2], 16).ok())
        .collect()
}

/// 将字节数组格式化为带空格分隔的十六进制字符串。
pub fn bytes_to_hex_string(bytes: &[u8]) -> String {
    let mut result = String::with_capacity(bytes.len() * 3);
    for (index, byte) in bytes.iter().enumerate() {
        if index > 0 {
            result.push(' ');
        }
        write!(&mut result, "{:02X}", byte).expect("write to string");
    }
    result
}
