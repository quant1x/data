use std::fmt::Write as FmtWrite;
use std::io::{Cursor, Read, Write};
use std::sync::atomic::{AtomicU32, Ordering};
use std::sync::{Arc, OnceLock, RwLock};

use anyhow::{anyhow, Context, Result};
use byteorder::{LittleEndian, ReadBytesExt};
use chrono::{Datelike, Duration as ChronoDuration, Local, NaiveDateTime, Timelike};
use encoding_rs::{BIG5, GB18030};
use flate2::read::ZlibDecoder;
use flate2::write::ZlibEncoder;
use flate2::Compression;

use crate::exchange::{MARKET_ID_SHANGHAI, MARKET_ID_SHENZHEN};

/// 连接池及时间计算使用的常量。
const TM_H_WIDTH: i32 = 1_000_000;
const TM_M_WIDTH: i32 = 10_000;

static SEQ_ID: AtomicU32 = AtomicU32::new(0);

type BaseUnitFn = Arc<dyn Fn(u16, &str) -> f64 + Send + Sync>;
static BASE_UNIT_FN: OnceLock<RwLock<BaseUnitFn>> = OnceLock::new();

fn base_unit_registry() -> &'static RwLock<BaseUnitFn> {
    BASE_UNIT_FN.get_or_init(|| RwLock::new(Arc::new(default_base_unit_impl)))
}

fn default_base_unit_impl(market_id: u16, code: &str) -> f64 {
    let mut unit = 100.0;
    if market_id == MARKET_ID_SHANGHAI {
        if code.starts_with("51") {
            unit = 1000.0;
        }
    } else if market_id == MARKET_ID_SHENZHEN {
        if code.starts_with("159") {
            unit = 1000.0;
        }
    }
    unit
}

pub fn base_unit(market_id: u16, code: &str) -> f64 {
    let guard = base_unit_registry()
        .read()
        .expect("base unit lock poisoned");
    guard(market_id, code)
}

pub fn register_base_unit_function<F>(func: F)
where
    F: Fn(u16, &str) -> f64 + Send + Sync + 'static,
{
    let mut guard = base_unit_registry()
        .write()
        .expect("base unit lock poisoned");
    *guard = Arc::new(func);
}

pub fn zlib_compress(src: &[u8]) -> Result<Vec<u8>> {
    let mut encoder = ZlibEncoder::new(Vec::new(), Compression::default());
    encoder.write_all(src)?;
    Ok(encoder.finish()?)
}

/// Decompresses zlib-compressed data into a byte vector.
///
/// # Arguments
/// * `compressed` - A slice containing the zlib-compressed data
///
/// # Returns
/// A `Result` containing:
/// * `Ok(Vec<u8>)` - The decompressed byte vector on success
/// * `Err` - If decompression fails
///
/// # Errors
/// Returns an error if:
/// * The input data is not valid zlib-compressed data
/// * The decompression process fails
///
/// # Examples
/// ```
/// let compressed_data = vec![...]; // zlib-compressed bytes
/// let decompressed = zlib_decompress(&compressed_data)?;
/// ```
pub fn zlib_decompress(compressed: &[u8]) -> Result<Vec<u8>> {
    let mut decoder = ZlibDecoder::new(compressed);
    let mut out = Vec::new();
    decoder.read_to_end(&mut out)?;
    Ok(out)
}

pub fn hex_string_to_bytes(hex_str: &str) -> Result<Vec<u8>> {
    let filtered: String = hex_str.chars().filter(|c| !c.is_whitespace()).collect();
    if filtered.len() % 2 != 0 {
        return Err(anyhow!("hex string has odd length"));
    }
    (0..filtered.len())
        .step_by(2)
        .map(|i| u8::from_str_radix(&filtered[i..i + 2], 16).context("parse hex"))
        .collect()
}

pub fn bytes_to_hex_string(bytes: &[u8]) -> String {
    let mut result = String::with_capacity(bytes.len() * 3);
    for (index, byte) in bytes.iter().enumerate() {
        if index > 0 {
            result.push(' ');
        }
        let _ = write!(&mut result, "{:02X}", byte);
    }
    result
}

pub fn sequence_id() -> u32 {
    SEQ_ID.fetch_add(1, Ordering::SeqCst) + 1
}

pub fn utf8_to_gbk(text: &[u8]) -> String {
    let mut slice = text;
    if let Some(pos) = text.iter().position(|&b| b == 0x00) {
        slice = &text[..pos];
    }
    let (decoded, _, _) = GB18030.decode(slice);
    decoded
        .chars()
        .filter(|&c| c != '\u{0}')
        .collect::<String>()
}

pub fn v1_utf8_to_gbk(text: &[u8]) -> String {
    let (decoded, _, _) = GB18030.decode(text);
    decoded
        .chars()
        .filter(|&c| c != '\u{0}')
        .collect::<String>()
}

pub fn decode_gbk(bytes: &[u8]) -> Result<Vec<u8>> {
    let (decoded, _, had_errors) = GB18030.decode(bytes);
    if had_errors {
        return Err(anyhow!("gbk decode error"));
    }
    Ok(decoded.into_owned().into_bytes())
}

pub fn encode_gbk(bytes: &[u8]) -> Result<Vec<u8>> {
    let text = std::str::from_utf8(bytes)?;
    let (encoded, _, had_errors) = GB18030.encode(text);
    if had_errors {
        return Err(anyhow!("gbk encode error"));
    }
    Ok(encoded.into_owned())
}

pub fn decode_big5(bytes: &[u8]) -> Result<Vec<u8>> {
    let (decoded, _, had_errors) = BIG5.decode(bytes);
    if had_errors {
        return Err(anyhow!("big5 decode error"));
    }
    Ok(decoded.into_owned().into_bytes())
}

pub fn encode_big5(bytes: &[u8]) -> Result<Vec<u8>> {
    let text = std::str::from_utf8(bytes)?;
    let (encoded, _, had_errors) = BIG5.encode(text);
    if had_errors {
        return Err(anyhow!("big5 encode error"));
    }
    Ok(encoded.into_owned())
}

pub fn is_nan_or_inf(value: f64) -> bool {
    value.is_nan() || value.is_infinite()
}

pub fn number_to_f64<T>(value: T) -> f64
where
    T: Into<f64>,
{
    value.into()
}

pub fn decimal(value: f64, digits: Option<i32>) -> f64 {
    let mut digits = digits.unwrap_or(2);
    if digits < 0 {
        digits = 0;
    }

    let mut value = value;
    if value.is_nan() || value.is_infinite() {
        value = 0.0;
    }

    let half = if value.is_sign_negative() { -5.0 } else { 5.0 };
    let digits_i32 = digits as i32;
    let n10 = 10f64.powi(digits_i32);
    let nj1 = 10f64.powi(digits_i32 + 1);

    ((value * nj1 + half) / 10.0).trunc() / n10
}

pub fn int_to_f64<T>(integer: T) -> f64
where
    T: Into<i64>,
{
    let ivol = integer.into();
    let log_point = (ivol >> 24) as i32;
    let hleax = ((ivol >> 16) & 0xff) as i32;
    let lheax = ((ivol >> 8) & 0xff) as i32;
    let lleax = (ivol & 0xff) as i32;

    let dw_ecx = log_point * 2 - 0x7f;
    let dw_edx = log_point * 2 - 0x86;
    let dw_esi = log_point * 2 - 0x8e;
    let dw_eax = log_point * 2 - 0x96;

    let tmp_eax = if dw_ecx < 0 { -dw_ecx } else { dw_ecx };
    let mut dbl_xmm6 = 2f64.powf(tmp_eax as f64);
    if dw_ecx < 0 {
        dbl_xmm6 = 1.0 / dbl_xmm6;
    }

    let dbl_xmm4 = if hleax > 0x80 {
        let tmpdbl_xmm3 = 2f64.powf((dw_edx + 1) as f64);
        let mut dbl_xmm0 = 2f64.powf(dw_edx as f64) * 128.0;
        dbl_xmm0 += (hleax & 0x7f) as f64 * tmpdbl_xmm3;
        dbl_xmm0
    } else if dw_edx >= 0 {
        2f64.powf(dw_edx as f64) * hleax as f64
    } else {
        (1.0 / 2f64.powf(dw_edx as f64)) * hleax as f64
    };

    let mut dbl_xmm3 = 2f64.powf(dw_esi as f64) * lheax as f64;
    let mut dbl_xmm1 = 2f64.powf(dw_eax as f64) * lleax as f64;
    if hleax & 0x80 != 0 {
        dbl_xmm3 *= 2.0;
        dbl_xmm1 *= 2.0;
    }

    dbl_xmm6 + dbl_xmm4 + dbl_xmm3 + dbl_xmm1
}

pub fn format_time0(time_stamp: &str) -> Result<String> {
    format_time_common(time_stamp)
}

pub fn time_from_str(time_stamp: &str) -> Result<String> {
    format_time_common(time_stamp)
}

fn format_time_common(time_stamp: &str) -> Result<String> {
    let length = time_stamp.len();
    if length < 6 {
        return Err(anyhow!("time stamp too short"));
    }
    let t1 = time_stamp[..length - 6].parse::<i64>()?;
    let mut tm = format!("{:02}:", t1);
    let tmp = &time_stamp[length - 6..length - 4];
    let mut n = tmp.parse::<i64>()?;
    if n < 60 {
        tm.push_str(&format!("{:02}:", tmp));
        let tmp2 = &time_stamp[length - 4..];
        let f = tmp2.parse::<f64>()?;
        tm.push_str(&format!("{:06.3}", (f * 60.0) / 10_000.0));
    } else {
        let tmp2 = &time_stamp[length - 6..];
        let f = tmp2.parse::<f64>()?;
        tm.push_str(&format!("{:02}:", ((f * 60.0) / 1_000_000.0) as i64));
        n = f as i64;
        tm.push_str(&format!(
            "{:06.3}",
            ((n * 60) % 1_000_000) as f64 * 60.0 / 1_000_000.0
        ));
    }
    Ok(tm)
}

pub fn time_from_int(stamp: i32) -> String {
    let mut h = stamp / TM_H_WIDTH;
    let tmp1 = stamp % TM_H_WIDTH;
    let m1 = tmp1 / TM_M_WIDTH;
    let tmp2 = tmp1 % TM_M_WIDTH;
    let m;
    let st;
    if m1 < 60 {
        m = m1;
        let tmp3 = tmp2 * 60;
        st = tmp3 as f64 / TM_M_WIDTH as f64;
    } else {
        h += 1;
        let tmp3 = tmp1;
        m = tmp3 / TM_H_WIDTH;
        let tmp4 = (tmp3 % TM_H_WIDTH) * 60;
        st = tmp4 as f64 / TM_H_WIDTH as f64;
    }
    format!("{:02}:{:02}:{:06.3}", h, m, st)
}

pub fn get_datetime_from_u32(
    category: i32,
    zipday: u32,
    tminutes: u16,
) -> (i32, i32, i32, i32, i32) {
    let mut hour = 15;
    let mut minute = 0;
    let year;
    let month;
    let day;
    if category < 4 || category == 7 || category == 8 {
        year = ((zipday >> 11) + 2004) as i32;
        month = ((zipday % 2048) / 100) as i32;
        day = ((zipday % 2048) % 100) as i32;
        hour = (tminutes / 60) as i32;
        minute = (tminutes % 60) as i32;
    } else {
        year = (zipday / 10_000) as i32;
        month = ((zipday % 10_000) / 100) as i32;
        day = (zipday % 100) as i32;
    }
    (year, month, day, hour, minute)
}

pub fn get_datetime_now(category: i32, lasttime: &str) -> Result<(i32, i32, i32, i32, i32)> {
    let mut utime = NaiveDateTime::parse_from_str(lasttime, "%Y-%m-%d %H:%M:%S")?;
    match category {
        0 => utime += ChronoDuration::minutes(5),
        1 => utime += ChronoDuration::minutes(15),
        2 => utime += ChronoDuration::minutes(30),
        3 => utime += ChronoDuration::hours(1),
        4 => utime += ChronoDuration::days(1),
        5 => utime += ChronoDuration::days(7),
        6 => utime += ChronoDuration::days(30),
        7 | 8 => utime += ChronoDuration::minutes(1),
        9 => utime += ChronoDuration::days(1),
        10 => utime += ChronoDuration::days(90),
        11 => utime += ChronoDuration::days(365),
        _ => {}
    }

    let mut hour;
    let mut minute;
    if category < 4 || category == 7 || category == 8 {
        if (utime.hour() >= 15 && utime.minute() > 0) || (utime.hour() > 15) {
            utime += ChronoDuration::days(1);
            utime += ChronoDuration::minutes(30);
            hour = ((utime.hour() + 18) % 24) as i32;
        } else {
            hour = utime.hour() as i32;
        }
        minute = utime.minute() as i32;
    } else {
        let now = Local::now().naive_local();
        if utime.and_utc().timestamp() > now.and_utc().timestamp() {
            utime = now;
        }
        hour = utime.hour() as i32;
        minute = utime.minute() as i32;
        if utime.hour() > 15 {
            hour = 15;
            minute = 0;
        }
    }

    Ok((
        utime.year(),
        utime.month() as i32,
        utime.day() as i32,
        hour,
        minute,
    ))
}

pub fn get_time(data: &[u8], pos: &mut usize) -> Result<(u16, u16)> {
    if *pos + 2 > data.len() {
        return Err(anyhow!("buffer underflow"));
    }
    let mut cursor = Cursor::new(&data[*pos..*pos + 2]);
    let sec = cursor.read_u16::<LittleEndian>()?;
    *pos += 2;
    Ok((sec / 60, sec % 60))
}

pub fn get_datetime(
    category: i32,
    data: &[u8],
    pos: &mut usize,
) -> Result<(i32, i32, i32, i32, i32)> {
    if category < 4 || category == 7 || category == 8 {
        if *pos + 4 > data.len() {
            return Err(anyhow!("buffer underflow"));
        }
        let mut cursor = Cursor::new(&data[*pos..*pos + 2]);
        let zipday = cursor.read_u16::<LittleEndian>()?;
        *pos += 2;
        let mut cursor = Cursor::new(&data[*pos..*pos + 2]);
        let tminutes = cursor.read_u16::<LittleEndian>()?;
        *pos += 2;
        Ok(get_datetime_from_u32(category, zipday as u32, tminutes))
    } else {
        if *pos + 4 > data.len() {
            return Err(anyhow!("buffer underflow"));
        }
        let mut cursor = Cursor::new(&data[*pos..*pos + 4]);
        let zipday = cursor.read_u32::<LittleEndian>()?;
        *pos += 4;
        Ok(get_datetime_from_u32(category, zipday, 0))
    }
}

pub fn decode_varint(data: &[u8], pos: &mut usize) -> Result<i32> {
    if *pos >= data.len() {
        return Err(anyhow!("buffer underflow"));
    }

    let mut pos_byte = 6;
    let mut b_data = data[*pos];
    let mut value = (b_data & 0x3f) as i32;
    let sign = (b_data & 0x40) > 0;

    if (b_data & 0x80) > 0 {
        loop {
            *pos += 1;
            if *pos >= data.len() {
                return Err(anyhow!("buffer underflow"));
            }
            b_data = data[*pos];
            value += ((b_data & 0x7f) as i32) << pos_byte;
            pos_byte += 7;
            if (b_data & 0x80) == 0 {
                break;
            }
        }
    }

    *pos += 1;
    if sign {
        value = -value;
    }
    Ok(value)
}
