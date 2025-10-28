use crate::internal::helpers::{get_datetime_from_u32, int_to_f64, sequence_id};
use crate::proto::STD_MSG_XDXR_INFO;
use crate::quotes::message::{Message, StdRequestHeader, StdResponseHeader};
use anyhow::Result;
use byteorder::{LittleEndian, ReadBytesExt};
use std::any::Any;
use std::io::{Cursor, Read};

fn get_xdxr_category_name(category: i32) -> &'static str {
    match category {
        1 => "除权除息",
        2 => "送配股上市",
        3 => "非流通股上市",
        4 => "未知股本变动",
        5 => "股本变化",
        6 => "增发新股",
        7 => "股份回购",
        8 => "增发新股上市",
        9 => "转配股上市",
        10 => "可转债上市",
        11 => "扩缩股",
        12 => "非流通股缩股",
        13 => "送认购权证",
        14 => "送认沽权证",
        _ => "未知",
    }
}

#[derive(Debug, Clone, Default)]
pub struct XdxrInfoRequest {
    pub market: u8,
    pub code: [u8; 6],
}

#[derive(Debug, Clone, Default)]
pub struct XdxrInfo {
    pub date: String,
    pub category: i32,
    pub name: String,
    pub fen_hong: f64,        // 分红
    pub pei_gu_jia: f64,      // 配股价
    pub song_zhuan_gu: f64,   // 送转股
    pub pei_gu: f64,          // 配股
    pub suo_gu: f64,          // 缩股
    pub qian_liu_tong: f64,   // 前流通
    pub hou_liu_tong: f64,    // 后流通
    pub qian_zong_gu_ben: f64, // 前总股本
    pub hou_zong_gu_ben: f64,  // 后总股本
    pub fen_shu: f64,         // 份数
    pub xing_quan_jia: f64,   // 行权价
}

impl XdxrInfo {
    pub fn is_capital_change(&self) -> bool {
        match self.category {
            1 | 11 | 12 | 13 | 14 => false,
            _ => {
                if self.hou_liu_tong > 0.0 && self.hou_zong_gu_ben > 0.0 {
                    true
                } else {
                    false
                }
            }
        }
    }

    pub fn adjust(&self) -> Box<dyn Fn(f64) -> f64 + Send + Sync> {
        let song_zhuan_gu = self.song_zhuan_gu;
        let pei_gu = self.pei_gu;
        let suo_gu = self.suo_gu;
        let fen_hong = self.fen_hong;
        let pei_gu_jia = self.pei_gu_jia;

        let xdxr_gu_shu = (song_zhuan_gu + pei_gu - suo_gu) / 10.0;
        let xdxr_fen_hong = (pei_gu_jia * pei_gu - fen_hong) / 10.0;

        Box::new(move |p: f64| (p + xdxr_fen_hong) / (1.0 + xdxr_gu_shu))
    }
}

#[derive(Debug, Clone)]
pub struct XdxrInfoPackage {
    req_header: StdRequestHeader,
    resp_header: StdResponseHeader,
    request: XdxrInfoRequest,
    reply: Vec<XdxrInfo>,
}

impl XdxrInfoPackage {
    pub fn new() -> Self {
        let mut req_header = StdRequestHeader::default();
        req_header.zip_flag = 0x0c; // FlagNotZipped
        req_header.seq_id = sequence_id();
        req_header.packet_type = 0x01;
        req_header.pkg_len1 = 0x000b;
        req_header.pkg_len2 = 0x000b;
        req_header.method = STD_MSG_XDXR_INFO;

        Self {
            req_header,
            resp_header: StdResponseHeader::default(),
            request: XdxrInfoRequest::default(),
            reply: Vec::new(),
        }
    }

    pub fn set_params(&mut self, req: XdxrInfoRequest) {
        self.request = req;
    }
}

impl Message for XdxrInfoPackage {
    fn serialize(&mut self) -> Result<Vec<u8>> {
        self.req_header.seq_id = sequence_id();
        let mut out = self.req_header.to_bytes()?;

        // Add content hex "0100"
        out.extend_from_slice(&[0x01, 0x00]);

        // Write request
        out.push(self.request.market);
        out.extend_from_slice(&self.request.code);

        Ok(out)
    }

    fn unserialize(&mut self, _header: &StdResponseHeader, data: &[u8]) -> Result<()> {
        self.resp_header = _header.clone();

        let mut cursor = Cursor::new(data);

        // Skip unknown 9 bytes
        cursor.set_position(cursor.position() + 9);

        let count = cursor.read_u16::<LittleEndian>()?;
        self.reply = Vec::with_capacity(count as usize);

        for _ in 0..count {
            let mut raw_data = [0u8; 29];
            cursor.read_exact(&mut raw_data)?;

            let market = raw_data[0] as i32;
            let code_bytes = &raw_data[1..7];
            let code = String::from_utf8_lossy(code_bytes)
                .trim_end_matches('\0')
                .to_string();
            let unknown = raw_data[7] as i32;

            let date_bytes = &raw_data[8..12];
            let date = u32::from_le_bytes(date_bytes.try_into().unwrap());

            let category = raw_data[12] as i32;
            let data_bytes = &raw_data[13..29];

            let (year, month, day, _, _) = get_datetime_from_u32(9, date, 0);

            let mut xdxr = XdxrInfo {
                date: format!("{:04}-{:02}-{:02}", year, month, day),
                category,
                name: get_xdxr_category_name(category).to_string(),
                ..Default::default()
            };

            match category {
                1 => {
                    // 除权除息
                    let mut data_cursor = Cursor::new(data_bytes);
                    xdxr.fen_hong = data_cursor.read_f32::<LittleEndian>()? as f64;
                    xdxr.pei_gu_jia = data_cursor.read_f32::<LittleEndian>()? as f64;
                    xdxr.song_zhuan_gu = data_cursor.read_f32::<LittleEndian>()? as f64;
                    xdxr.pei_gu = data_cursor.read_f32::<LittleEndian>()? as f64;
                }
                11 | 12 => {
                    // 扩缩股/非流通股缩股
                    let mut data_cursor = Cursor::new(data_bytes);
                    data_cursor.set_position(8); // Skip 8 bytes
                    xdxr.suo_gu = data_cursor.read_f32::<LittleEndian>()? as f64;
                }
                13 | 14 => {
                    // 送认购权证/送认沽权证
                    let mut data_cursor = Cursor::new(data_bytes);
                    xdxr.xing_quan_jia = data_cursor.read_f32::<LittleEndian>()? as f64;
                    data_cursor.set_position(8); // Skip 4 bytes
                    xdxr.fen_shu = data_cursor.read_f32::<LittleEndian>()? as f64;
                }
                _ => {
                    // 其他类型 - 股本变化
                    let mut data_cursor = Cursor::new(data_bytes);
                    let qian_liu_tong = data_cursor.read_u32::<LittleEndian>()?;
                    let qian_zong_gu_ben = data_cursor.read_u32::<LittleEndian>()?;
                    let hou_liu_tong = data_cursor.read_u32::<LittleEndian>()?;
                    let hou_zong_gu_ben = data_cursor.read_u32::<LittleEndian>()?;

                    xdxr.qian_liu_tong = get_v(qian_liu_tong);
                    xdxr.qian_zong_gu_ben = get_v(qian_zong_gu_ben);
                    xdxr.hou_liu_tong = get_v(hou_liu_tong);
                    xdxr.hou_zong_gu_ben = get_v(hou_zong_gu_ben);
                }
            }

            self.reply.push(xdxr);
        }

        Ok(())
    }

    fn reply(&self) -> &(dyn Any + Send + Sync) {
        &self.reply
    }
}

pub fn new_xdxr_info() -> Result<XdxrInfoPackage> {
    Ok(XdxrInfoPackage::new())
}

fn get_v(v: u32) -> f64 {
    if v == 0 {
        0.0
    } else {
        int_to_f64(v as i64)
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_xdxr_info_default() {
        let xdxr = XdxrInfo::default();
        assert_eq!(xdxr.date, "");
        assert_eq!(xdxr.category, 0);
        assert_eq!(xdxr.name, "");
        assert_eq!(xdxr.fen_hong, 0.0);
        assert_eq!(xdxr.pei_gu_jia, 0.0);
        assert_eq!(xdxr.song_zhuan_gu, 0.0);
        assert_eq!(xdxr.pei_gu, 0.0);
        assert_eq!(xdxr.suo_gu, 0.0);
        assert_eq!(xdxr.qian_liu_tong, 0.0);
        assert_eq!(xdxr.hou_liu_tong, 0.0);
        assert_eq!(xdxr.qian_zong_gu_ben, 0.0);
        assert_eq!(xdxr.hou_zong_gu_ben, 0.0);
        assert_eq!(xdxr.fen_shu, 0.0);
        assert_eq!(xdxr.xing_quan_jia, 0.0);
    }

    #[test]
    fn test_xdxr_info_is_capital_change() {
        let mut xdxr = XdxrInfo::default();

        // Category 1 (除权除息) should not be capital change
        xdxr.category = 1;
        assert!(!xdxr.is_capital_change());

        // Category 11 (扩缩股) should not be capital change
        xdxr.category = 11;
        assert!(!xdxr.is_capital_change());

        // Category 5 (股本变化) with valid circulation and total shares should be capital change
        xdxr.category = 5;
        xdxr.hou_liu_tong = 100.0;
        xdxr.hou_zong_gu_ben = 200.0;
        assert!(xdxr.is_capital_change());

        // Category 5 with zero circulation should not be capital change
        xdxr.hou_liu_tong = 0.0;
        xdxr.hou_zong_gu_ben = 200.0;
        assert!(!xdxr.is_capital_change());

        // Category 5 with zero total shares should not be capital change
        xdxr.hou_liu_tong = 100.0;
        xdxr.hou_zong_gu_ben = 0.0;
        assert!(!xdxr.is_capital_change());
    }

    #[test]
    fn test_xdxr_info_adjust_dividend() {
        // Test 除权除息 (category 1) adjustment
        let xdxr = XdxrInfo {
            category: 1,
            fen_hong: 2.0,        // 分红 2元
            pei_gu_jia: 10.0,     // 配股价 10元
            song_zhuan_gu: 5.0,   // 送转股 5股
            pei_gu: 3.0,          // 配股 3股
            suo_gu: 1.0,          // 缩股 1股
            ..Default::default()
        };

        let adjust_fn = xdxr.adjust();
        let result = adjust_fn(100.0);

        // Expected calculation:
        // xdxr_gu_shu = (5 + 3 - 1) / 10 = 0.7
        // xdxr_fen_hong = (10 * 3 - 2) / 10 = 2.8
        // result = (100 + 2.8) / (1 + 0.7) = 102.8 / 1.7 ≈ 60.47
        assert!((result - 60.47).abs() < 0.01);
    }

    #[test]
    fn test_xdxr_info_adjust_no_change() {
        // Test with all zeros - should return original price
        let xdxr = XdxrInfo {
            category: 1,
            fen_hong: 0.0,
            pei_gu_jia: 0.0,
            song_zhuan_gu: 0.0,
            pei_gu: 0.0,
            suo_gu: 0.0,
            ..Default::default()
        };

        let adjust_fn = xdxr.adjust();
        let result = adjust_fn(100.0);
        assert!((result - 100.0).abs() < 0.001);
    }

    #[test]
    fn test_xdxr_info_adjust_negative_price() {
        let xdxr = XdxrInfo {
            category: 1,
            fen_hong: 10.0,       // High dividend
            pei_gu_jia: 5.0,
            song_zhuan_gu: 0.0,
            pei_gu: 0.0,
            suo_gu: 0.0,
            ..Default::default()
        };

        let adjust_fn = xdxr.adjust();
        let result = adjust_fn(1.0); // Very low price
        // Should still produce a valid result
        assert!(result.is_finite());
        assert!(result >= 0.0);
    }

    #[test]
    fn test_xdxr_info_package_new() {
        let pkg = XdxrInfoPackage::new();
        assert_eq!(pkg.req_header.method, STD_MSG_XDXR_INFO);
        assert_eq!(pkg.req_header.zip_flag, 0x0c);
        assert_eq!(pkg.req_header.packet_type, 0x01);
        assert_eq!(pkg.req_header.pkg_len1, 0x000b);
        assert_eq!(pkg.req_header.pkg_len2, 0x000b);
        assert!(pkg.reply.is_empty());
    }

    #[test]
    fn test_xdxr_info_package_set_params() {
        let mut pkg = XdxrInfoPackage::new();
        let request = XdxrInfoRequest {
            market: 1,
            code: [b'0', b'0', b'0', b'0', b'0', b'1'],
        };

        pkg.set_params(request.clone());
        assert_eq!(pkg.request.market, request.market);
        assert_eq!(pkg.request.code, request.code);
    }

    #[test]
    fn test_xdxr_category_mapping() {
        assert_eq!(get_xdxr_category_name(1), "除权除息");
        assert_eq!(get_xdxr_category_name(2), "送配股上市");
        assert_eq!(get_xdxr_category_name(3), "非流通股上市");
        assert_eq!(get_xdxr_category_name(4), "未知股本变动");
        assert_eq!(get_xdxr_category_name(5), "股本变化");
        assert_eq!(get_xdxr_category_name(6), "增发新股");
        assert_eq!(get_xdxr_category_name(7), "股份回购");
        assert_eq!(get_xdxr_category_name(8), "增发新股上市");
        assert_eq!(get_xdxr_category_name(9), "转配股上市");
        assert_eq!(get_xdxr_category_name(10), "可转债上市");
        assert_eq!(get_xdxr_category_name(11), "扩缩股");
        assert_eq!(get_xdxr_category_name(12), "非流通股缩股");
        assert_eq!(get_xdxr_category_name(13), "送认购权证");
        assert_eq!(get_xdxr_category_name(14), "送认沽权证");
        assert_eq!(get_xdxr_category_name(99), "未知");
        assert_eq!(get_xdxr_category_name(0), "未知");
    }

    #[test]
    fn test_get_v() {
        assert_eq!(get_v(0), 0.0);
        // Test with a non-zero value - this would depend on int_to_f64 implementation
        assert!(get_v(1000) > 0.0);
        assert!(get_v(1000000) > 0.0);
    }

    #[test]
    fn test_xdxr_info_package_serialize() {
        let mut pkg = XdxrInfoPackage::new();
        let request = XdxrInfoRequest {
            market: 1,
            code: [b'0', b'0', b'0', b'0', b'0', b'1'],
        };
        pkg.set_params(request);

        let data = pkg.serialize().unwrap();

        // Check that data contains expected elements
        assert!(data.len() >= 16); // Header + content
        // The serialized data should contain the market and code bytes
        assert!(data.len() >= 25); // Header (16) + content (at least 9 bytes)
    }

    #[test]
    fn test_xdxr_info_request_default() {
        let req = XdxrInfoRequest::default();
        assert_eq!(req.market, 0);
        assert_eq!(req.code, [0u8; 6]);
    }

    #[test]
    fn test_xdxr_info_clone() {
        let xdxr = XdxrInfo {
            date: "2023-01-01".to_string(),
            category: 1,
            name: "除权除息".to_string(),
            fen_hong: 1.5,
            pei_gu_jia: 8.0,
            song_zhuan_gu: 2.0,
            pei_gu: 1.0,
            suo_gu: 0.0,
            qian_liu_tong: 1000000.0,
            hou_liu_tong: 1100000.0,
            qian_zong_gu_ben: 2000000.0,
            hou_zong_gu_ben: 2100000.0,
            fen_shu: 0.0,
            xing_quan_jia: 0.0,
        };

        let cloned = xdxr.clone();
        assert_eq!(xdxr.date, cloned.date);
        assert_eq!(xdxr.category, cloned.category);
        assert_eq!(xdxr.fen_hong, cloned.fen_hong);
        assert_eq!(xdxr.hou_liu_tong, cloned.hou_liu_tong);
    }

    #[test]
    fn test_xdxr_info_adjust_extreme_values() {
        // Test with very large values
        let xdxr = XdxrInfo {
            song_zhuan_gu: 1000.0,
            pei_gu: 500.0,
            suo_gu: 100.0,
            fen_hong: 100.0,
            pei_gu_jia: 50.0,
            ..Default::default()
        };

        let adjust_fn = xdxr.adjust();
        let result = adjust_fn(1000.0);
        assert!(result.is_finite());
        assert!(result > 0.0);

        // Test with very small values
        let xdxr_small = XdxrInfo {
            song_zhuan_gu: 0.001,
            pei_gu: 0.0005,
            suo_gu: 0.0001,
            fen_hong: 0.0001,
            pei_gu_jia: 0.005,
            ..Default::default()
        };

        let adjust_fn_small = xdxr_small.adjust();
        let result_small = adjust_fn_small(0.1);
        assert!(result_small.is_finite());
        assert!(result_small > 0.0);
    }

    #[test]
    fn test_new_xdxr_info() {
        let result = new_xdxr_info();
        assert!(result.is_ok());
        let pkg = result.unwrap();
        assert_eq!(pkg.req_header.method, STD_MSG_XDXR_INFO);
    }

    #[test]
    fn test_xdxr_info_debug() {
        let xdxr = XdxrInfo {
            date: "2023-01-01".to_string(),
            category: 1,
            name: "除权除息".to_string(),
            fen_hong: 1.5,
            ..Default::default()
        };

        let debug_str = format!("{:?}", xdxr);
        assert!(debug_str.contains("XdxrInfo"));
        assert!(debug_str.contains("2023-01-01"));
        assert!(debug_str.contains("除权除息"));
    }

    #[test]
    fn test_xdxr_info_package_debug() {
        let pkg = XdxrInfoPackage::new();
        let debug_str = format!("{:?}", pkg);
        assert!(debug_str.contains("XdxrInfoPackage"));
    }
}